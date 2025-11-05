package credentials

import (
	"encoding/base64"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	"github.com/chrilep/LaunchRDP/logging"
)

// Windows DPAPI structures and functions
var (
	crypt32                = syscall.NewLazyDLL("crypt32.dll")
	kernel32               = syscall.NewLazyDLL("kernel32.dll")
	procCryptProtectData   = crypt32.NewProc("CryptProtectData")
	procCryptUnprotectData = crypt32.NewProc("CryptUnprotectData")
	procLocalFree          = kernel32.NewProc("LocalFree")
)

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

// StoreCredential stores a credential in Windows Credential Manager as Domain credential
// Uses cmdkey format for domain credentials: cmdkey /add:TERMSRV/hostname /user:username /pass:password
func (cm *CredentialManager) StoreCredential(hostname, username, password string) error {
	debug := false

	logging.Log(debug, "StoreCredential called for hostname:", hostname, "username:", username)

	// Skip invalid hostnames
	if hostname == "" {
		logging.Log(true, "ERROR: Invalid hostname provided")
		return fmt.Errorf("invalid hostname: %s", hostname)
	}

	// Use /add for domain credentials (Windows-Anmeldeinformationen)
	cmdArgs := []string{"/add:TERMSRV/" + hostname, "/user:" + username, "/pass:" + password}
	// SECURITY: Never log the actual cmdkey args as they contain the password in plaintext
	logging.Log(debug, "Executing cmdkey for TERMSRV/"+hostname+" with user "+username)

	cmd := exec.Command("cmdkey", cmdArgs...)

	// Clear password from memory as soon as possible
	for i := range cmdArgs {
		if strings.Contains(cmdArgs[i], "/pass:") {
			cmdArgs[i] = "/pass:***CLEARED***"
		}
	}
	// Clear the original password parameter
	password = ""

	logging.Log(debug, "Running cmdkey command...")
	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))

	if err != nil {
		logging.Log(true, "ERROR: cmdkey failed with error:", err, "output:", outputStr)
		return fmt.Errorf("failed to store credential for %s: %w, output: %s", hostname, err, outputStr)
	}

	logging.Log(debug, "cmdkey completed successfully, output:", outputStr)
	logging.Log(debug, "Successfully stored domain credential for", hostname, "with user", username)
	return nil
}

// StoreGenericCredential stores a generic RDP credential for a username
func (cm *CredentialManager) StoreGenericCredential(username, password string) error {
	debug := false

	// Use a generic target that works with RDP
	target := fmt.Sprintf("TERMSRV/RDP_%s", username)

	cmd := exec.Command("cmdkey", "/generic:"+target, "/user:"+username, "/pass:"+password)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to store generic credential: %w, output: %s", err, string(output))
	}

	logging.Log(debug, "Successfully stored generic credential for user", username, "at target", target)
	return nil
}

// TestCredentialStorage tests if credential storage is working
func (cm *CredentialManager) TestCredentialStorage() error {
	testTarget := "TERMSRV/LaunchRDP_Test"
	testUser := "testuser"
	testPass := "testpass"

	// Store test credential
	cmd := exec.Command("cmdkey", "/generic:"+testTarget, "/user:"+testUser, "/pass:"+testPass)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to store test credential: %w, output: %s", err, string(output))
	}

	// Delete test credential
	cmd = exec.Command("cmdkey", "/delete:"+testTarget)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete test credential: %w, output: %s", err, string(output))
	}

	return nil
} // DeleteCredential deletes a credential from Windows Credential Manager
func (cm *CredentialManager) DeleteCredential(hostname string) error {
	// Use /delete for domain credentials
	cmd := exec.Command("cmdkey", "/delete:"+hostname)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete credential: %w, output: %s", err, string(output))
	}

	return nil
}

// ListCredentials lists stored credentials (for verification)
func (cm *CredentialManager) ListCredentials() ([]string, error) {
	cmd := exec.Command("cmdkey", "/list")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list credentials: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	var credentials []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Target: TERMSRV/") {
			// Extract hostname from "Target: TERMSRV/hostname"
			target := strings.TrimPrefix(line, "Target: TERMSRV/")
			credentials = append(credentials, target)
		}
	}

	return credentials, nil
}

// HasCredential checks if a credential exists for the given hostname
func (cm *CredentialManager) HasCredential(hostname string) bool {
	credentials, err := cm.ListCredentials()
	if err != nil {
		return false
	}

	for _, cred := range credentials {
		if cred == hostname {
			return true
		}
	}

	return false
}

// EncryptPasswordDPAPI encrypts a password using Windows DPAPI (most secure for Windows)
// DPAPI (Data Protection API) ties encryption to the current user + machine
// Only the same user on the same machine can decrypt the data
func (cm *CredentialManager) EncryptPasswordDPAPI(password string) (string, error) {
	debug := false

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
	debug := false

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
