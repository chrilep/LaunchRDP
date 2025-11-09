package credentials

import (
	"encoding/base64"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/chrilep/LaunchRDP/app/logging"
)

// Windows DPAPI and Credential Manager structures and functions
var (
	crypt32                = syscall.NewLazyDLL("crypt32.dll")
	kernel32               = syscall.NewLazyDLL("kernel32.dll")
	advapi32               = syscall.NewLazyDLL("advapi32.dll")
	procCryptProtectData   = crypt32.NewProc("CryptProtectData")
	procCryptUnprotectData = crypt32.NewProc("CryptUnprotectData")
	procLocalFree          = kernel32.NewProc("LocalFree")
	procCredWriteW         = advapi32.NewProc("CredWriteW")
	procCredDeleteW        = advapi32.NewProc("CredDeleteW")
)

// Windows Credential structures
const (
	CRED_TYPE_GENERIC             = 0x1 // Generic credential type - works for RDP
	CRED_TYPE_DOMAIN_PASSWORD     = 0x2 // Domain password - more restrictive
	CRED_PERSIST_LOCAL_MACHINE    = 0x2
	CRED_PERSIST_ENTERPRISE       = 0x3
	CRED_MAX_CREDENTIAL_BLOB_SIZE = 512
)

type credential struct {
	Flags              uint32
	Type               uint32
	TargetName         *uint16
	Comment            *uint16
	LastWritten        syscall.Filetime
	CredentialBlobSize uint32
	CredentialBlob     *byte
	Persist            uint32
	AttributeCount     uint32
	Attributes         uintptr
	TargetAlias        *uint16
	UserName           *uint16
}

type dataBlob struct {
	cbData uint32
	pbData *byte
}

func newBlob(d []byte) *dataBlob {
	if len(d) == 0 {
		return &dataBlob{}
	}
	return &dataBlob{
		pbData: &d[0],
		cbData: uint32(len(d)),
	}
}

func (b *dataBlob) toByteArray() []byte {
	d := make([]byte, b.cbData)
	copy(d, (*[1 << 30]byte)(unsafe.Pointer(b.pbData))[:b.cbData])
	return d
}

// CredentialManager handles Windows credential store operations
type CredentialManager struct{}

// NewCredentialManager creates a new credential manager
func NewCredentialManager() *CredentialManager {
	return &CredentialManager{}
}

// StoreCredential stores a credential in Windows Credential Manager using native API
// Uses domain credential format: TERMSRV/hostname with CRED_TYPE_DOMAIN_PASSWORD
func (cm *CredentialManager) StoreCredential(hostname, username, password string) error {
	debug := false

	logging.Log(debug, "Input - hostname:", hostname, "username:", username, "password length:", len(password))

	// Skip invalid hostnames
	if hostname == "" {
		logging.Log(true, "ERROR: Invalid hostname provided")
		return fmt.Errorf("invalid hostname: %s", hostname)
	}

	if username == "" {
		logging.Log(true, "ERROR: Invalid username provided")
		return fmt.Errorf("invalid username")
	}

	if password == "" {
		logging.Log(true, "ERROR: Invalid password provided (empty)")
		return fmt.Errorf("invalid password: empty")
	}

	targetString := "TERMSRV/" + hostname
	logging.Log(debug, "Target string:", targetString)

	targetName, err := syscall.UTF16PtrFromString(targetString)
	if err != nil {
		logging.Log(true, "ERROR: Failed to convert target name to UTF16:", err)
		return fmt.Errorf("failed to convert target name: %v", err)
	}
	logging.Log(debug, "Target name converted to UTF16 successfully")

	// For CRED_TYPE_DOMAIN_PASSWORD, UserName must be in format DOMAIN\Username
	// If username doesn't contain backslash, assume local machine
	formattedUsername := username
	if !containsBackslash(username) {
		// No domain specified, use hostname as domain (for local accounts)
		formattedUsername = hostname + "\\" + username
		logging.Log(debug, "No domain in username, formatted as:", formattedUsername)
	} else {
		logging.Log(debug, "Username already contains domain:", formattedUsername)
	}

	userNamePtr, err := syscall.UTF16PtrFromString(formattedUsername)
	if err != nil {
		logging.Log(true, "ERROR: Failed to convert username to UTF16:", err)
		return fmt.Errorf("failed to convert username: %v", err)
	}
	logging.Log(debug, "Username converted to UTF16 successfully")

	// For CRED_TYPE_DOMAIN_PASSWORD, password must be UTF-16 encoded
	passwordUTF16, err := syscall.UTF16FromString(password)
	if err != nil {
		logging.Log(true, "ERROR: Failed to convert password to UTF16:", err)
		return fmt.Errorf("failed to convert password: %v", err)
	}

	// IMPORTANT: Remove null terminator! syscall.UTF16FromString adds a null terminator,
	// but CRED_TYPE_DOMAIN_PASSWORD does NOT want it in CredentialBlob
	// The documentation explicitly states: "do not include a trailing zero character"
	if len(passwordUTF16) > 0 && passwordUTF16[len(passwordUTF16)-1] == 0 {
		passwordUTF16 = passwordUTF16[:len(passwordUTF16)-1]
		logging.Log(debug, "Removed null terminator from UTF16 password")
	}

	// Log first few characters for debugging (only in hex to avoid exposing password)
	if len(passwordUTF16) >= 4 {
		logging.Log(debug, "Password UTF16 first 4 chars (hex):", fmt.Sprintf("%04x %04x %04x %04x",
			passwordUTF16[0], passwordUTF16[1], passwordUTF16[2], passwordUTF16[3]))
	}

	// Convert UTF-16 to bytes (WITHOUT null terminator)
	passwordBytes := make([]byte, len(passwordUTF16)*2)
	for i, r := range passwordUTF16 {
		passwordBytes[i*2] = byte(r)
		passwordBytes[i*2+1] = byte(r >> 8)
	}
	logging.Log(debug, "Password converted to UTF16, bytes length:", len(passwordBytes), "(without null terminator)")

	if len(passwordBytes) > CRED_MAX_CREDENTIAL_BLOB_SIZE {
		logging.Log(true, "ERROR: Password too long:", len(passwordBytes), "max:", CRED_MAX_CREDENTIAL_BLOB_SIZE)
		return fmt.Errorf("password too long (max %d bytes)", CRED_MAX_CREDENTIAL_BLOB_SIZE)
	}

	// Use CRED_TYPE_DOMAIN_PASSWORD for Windows Login Info
	cred := &credential{
		Type:               CRED_TYPE_DOMAIN_PASSWORD,
		TargetName:         targetName,
		CredentialBlobSize: uint32(len(passwordBytes)),
		CredentialBlob:     &passwordBytes[0],
		Persist:            CRED_PERSIST_LOCAL_MACHINE,
		UserName:           userNamePtr,
	}

	logging.Log(debug, "Credential struct created:")
	logging.Log(debug, "  Type:", cred.Type, "(CRED_TYPE_DOMAIN_PASSWORD)")
	logging.Log(debug, "  CredentialBlobSize:", cred.CredentialBlobSize)
	logging.Log(debug, "  Persist:", cred.Persist, "(CRED_PERSIST_LOCAL_MACHINE)")

	logging.Log(debug, "Calling CredWriteW...")
	ret, _, err := procCredWriteW.Call(
		uintptr(unsafe.Pointer(cred)),
		0,
	)

	logging.Log(debug, "CredWriteW returned - ret:", ret, "err:", err)

	if ret == 0 {
		logging.Log(true, "ERROR: CredWriteW failed - return value:", ret)
		logging.Log(true, "ERROR: System error:", err)
		logging.Log(true, "ERROR: Error code:", err.(syscall.Errno))
		return fmt.Errorf("failed to store credential for %s: %w (code: %d)", hostname, err, err.(syscall.Errno))
	}

	logging.Log(debug, "SUCCESS: Credential stored successfully")
	return nil
}

// containsBackslash checks if string contains a backslash
func containsBackslash(s string) bool {
	for _, c := range s {
		if c == '\\' {
			return true
		}
	}
	return false
}

// DeleteCredential deletes a credential from Windows Credential Manager using native API
func (cm *CredentialManager) DeleteCredential(hostname string) error {
	debug := true
	logging.Log(debug, "=== DeleteCredential START ===")
	logging.Log(debug, "Deleting credential for hostname:", hostname)

	targetString := "TERMSRV/" + hostname
	logging.Log(debug, "Target string:", targetString)

	targetName, err := syscall.UTF16PtrFromString(targetString)
	if err != nil {
		logging.Log(true, "ERROR: Failed to convert target name to UTF16:", err)
		return fmt.Errorf("failed to convert target name: %v", err)
	}

	logging.Log(debug, "Calling CredDeleteW...")
	ret, _, err := procCredDeleteW.Call(
		uintptr(unsafe.Pointer(targetName)),
		uintptr(CRED_TYPE_DOMAIN_PASSWORD),
		0,
	)

	logging.Log(debug, "CredDeleteW returned - ret:", ret, "err:", err)

	if ret == 0 {
		logging.Log(true, "ERROR: CredDeleteW failed - return value:", ret)
		logging.Log(true, "ERROR: System error:", err)
		if errno, ok := err.(syscall.Errno); ok {
			logging.Log(true, "ERROR: Error code:", errno)
		}
		return fmt.Errorf("failed to delete credential: %w", err)
	}

	logging.Log(debug, "SUCCESS: Credential deleted successfully")
	logging.Log(debug, "=== DeleteCredential END ===")
	return nil
}

// EncryptPasswordDPAPI encrypts a password using Windows DPAPI (most secure for Windows)
// DPAPI (Data Protection API) ties encryption to the current user + machine
// Only the same user on the same machine can decrypt the data
func (cm *CredentialManager) EncryptPasswordDPAPI(password string) (string, error) {
	debug := true

	if password == "" {
		return "", nil
	}

	logging.Log(debug, "Encrypting password with Windows DPAPI (native)")

	// Convert password to bytes
	passwordBytes := []byte(password)
	// Clear password from memory immediately
	password = ""

	// Create data blob for input
	dataIn := newBlob(passwordBytes)

	// Output blob for encrypted data
	var dataOut dataBlob

	// Call CryptProtectData
	ret, _, err := procCryptProtectData.Call(
		uintptr(unsafe.Pointer(dataIn)),   // pDataIn
		0,                                 // szDataDescr (optional)
		0,                                 // pOptionalEntropy (optional)
		0,                                 // pvReserved
		0,                                 // pPromptStruct (optional)
		0,                                 // dwFlags
		uintptr(unsafe.Pointer(&dataOut)), // pDataOut
	)

	if ret == 0 {
		logging.Log(true, "ERROR: CryptProtectData failed:", err)
		return "", fmt.Errorf("CryptProtectData failed: %v", err)
	}

	// Convert encrypted data to base64
	encryptedBytes := dataOut.toByteArray()
	encrypted := base64.StdEncoding.EncodeToString(encryptedBytes)

	// Free the memory allocated by CryptProtectData
	syscall.SyscallN(procLocalFree.Addr(), uintptr(unsafe.Pointer(dataOut.pbData)))

	// Clear sensitive data from memory
	for i := range passwordBytes {
		passwordBytes[i] = 0
	}

	logging.Log(true, "Password encrypted successfully with DPAPI, length:", len(encrypted))
	return encrypted, nil
}

// DecryptPasswordDPAPI decrypts a password using Windows DPAPI
func (cm *CredentialManager) DecryptPasswordDPAPI(encryptedPassword string) (string, error) {
	debug := true

	if encryptedPassword == "" {
		return "", nil
	}

	logging.Log(debug, "Decrypting password with Windows DPAPI")
	return cm.decryptWithDPAPI(encryptedPassword)
}

// decryptWithDPAPI performs pure DPAPI decryption
func (cm *CredentialManager) decryptWithDPAPI(encryptedPassword string) (string, error) {
	// Decode base64
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}

	// Create data blob for input
	dataIn := newBlob(encryptedBytes)

	// Output blob for decrypted data
	var dataOut dataBlob

	// Call CryptUnprotectData
	ret, _, err := procCryptUnprotectData.Call(
		uintptr(unsafe.Pointer(dataIn)),   // pDataIn
		0,                                 // ppszDataDescr (optional)
		0,                                 // pOptionalEntropy (optional)
		0,                                 // pvReserved
		0,                                 // pPromptStruct (optional)
		0,                                 // dwFlags
		uintptr(unsafe.Pointer(&dataOut)), // pDataOut
	)

	if ret == 0 {
		return "", fmt.Errorf("CryptUnprotectData failed: %v", err)
	}

	// Convert decrypted data to string
	decryptedBytes := dataOut.toByteArray()
	decrypted := string(decryptedBytes)

	// Free the memory allocated by CryptUnprotectData
	syscall.SyscallN(procLocalFree.Addr(), uintptr(unsafe.Pointer(dataOut.pbData)))

	return decrypted, nil
}

// MigratePasswordToDPAPI migrates a password from legacy AES to DPAPI encryption

// EncryptPasswordForUserEdit encrypts a password for user credential editing (JSON storage only)
// This is used when user edits credentials - only saves to JSON, not to CredStore
func (cm *CredentialManager) EncryptPasswordForUserEdit(password string) (string, error) {
	logging.Log(true, "EncryptPasswordForUserEdit: Encrypting password for JSON storage only")

	if password == "" {
		return "", nil
	}

	// Always use DPAPI for new password storage
	encrypted, err := cm.EncryptPasswordDPAPI(password)
	if err != nil {
		logging.Log(true, "EncryptPasswordForUserEdit: DPAPI encryption failed:", err)
		return "", fmt.Errorf("failed to encrypt password for user edit: %v", err)
	}

	logging.Log(true, "EncryptPasswordForUserEdit: Password encrypted successfully for JSON storage")
	return encrypted, nil
}

// StoreCredentialForHostEdit stores credentials for host editing (JSON + CredStore)
// This is used when user edits hosts - saves to both JSON and Windows CredStore
func (cm *CredentialManager) StoreCredentialForHostEdit(hostname, username, encryptedPassword string) error {
	logging.Log(true, "StoreCredentialForHostEdit: Processing host edit for:", hostname, "user:", username)

	if encryptedPassword == "" {
		logging.Log(true, "StoreCredentialForHostEdit: No password provided, skipping CredStore update")
		return nil
	}

	// Decrypt the password first
	password, err := cm.DecryptPasswordDPAPI(encryptedPassword)
	if err != nil {
		logging.Log(true, "StoreCredentialForHostEdit: Failed to decrypt password:", err)
		return fmt.Errorf("failed to decrypt password for host edit: %v", err)
	}

	// Store in Windows CredStore
	logging.Log(true, "StoreCredentialForHostEdit: Storing credential in Windows CredStore")
	storeErr := cm.StoreCredential(hostname, username, password)

	// Clear password from memory immediately after use
	password = ""

	if storeErr != nil {
		logging.Log(true, "StoreCredentialForHostEdit: Failed to store in CredStore:", storeErr)
		return fmt.Errorf("failed to store credential in CredStore: %v", storeErr)
	}

	logging.Log(true, "StoreCredentialForHostEdit: Host credential edit completed successfully")
	return nil
}

// UpdateUserCredentials updates user credentials, reusing existing password if new password is empty
// Returns the (possibly unchanged) encrypted password and whether migration occurred
func (cm *CredentialManager) UpdateUserCredentials(oldEncryptedPassword, newPlaintextPassword string) (string, bool, error) {
	logging.Log(true, "UpdateUserCredentials: Processing user credential update")

	// If no new password provided, return the old encrypted password unchanged
	if newPlaintextPassword == "" {
		logging.Log(true, "UpdateUserCredentials: No new password provided, keeping existing password")
		return oldEncryptedPassword, false, nil
	}

	// If new password provided, encrypt it
	logging.Log(true, "UpdateUserCredentials: New password provided, encrypting with DPAPI")
	newEncrypted, err := cm.EncryptPasswordDPAPI(newPlaintextPassword)
	if err != nil {
		logging.Log(true, "UpdateUserCredentials: Failed to encrypt new password:", err)
		return "", false, fmt.Errorf("failed to encrypt new password: %v", err)
	}

	// Clear the plaintext password from memory
	newPlaintextPassword = ""

	logging.Log(true, "UpdateUserCredentials: New password encrypted successfully")
	return newEncrypted, false, nil
}

// UpdateUserCredentialsWithMigration updates user credentials
// Returns the encrypted password, whether migration occurred, and any error
func (cm *CredentialManager) UpdateUserCredentialsWithMigration(oldEncryptedPassword, newPlaintextPassword string) (string, bool, error) {
	logging.Log(true, "UpdateUserCredentialsWithMigration: Processing user credential update")

	// If no new password provided, return old password unchanged
	if newPlaintextPassword == "" {
		logging.Log(true, "UpdateUserCredentialsWithMigration: No new password provided, keeping existing password")
		return oldEncryptedPassword, false, nil
	}

	// If new password provided, encrypt it with DPAPI
	logging.Log(true, "UpdateUserCredentialsWithMigration: New password provided, encrypting with DPAPI")
	newEncrypted, err := cm.EncryptPasswordDPAPI(newPlaintextPassword)
	if err != nil {
		logging.Log(true, "UpdateUserCredentialsWithMigration: Failed to encrypt new password:", err)
		return "", false, fmt.Errorf("failed to encrypt new password: %v", err)
	}

	// Clear the plaintext password from memory
	newPlaintextPassword = ""

	logging.Log(true, "UpdateUserCredentialsWithMigration: New password encrypted successfully")
	return newEncrypted, false, nil
}
