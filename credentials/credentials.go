package credentials

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
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
	logging.Log(true, "StoreCredential called for hostname:", hostname, "username:", username)

	// Skip invalid hostnames
	if hostname == "" {
		logging.Log(true, "ERROR: Invalid hostname provided")
		return fmt.Errorf("invalid hostname: %s", hostname)
	}

	// Use /add for domain credentials (Windows-Anmeldeinformationen)
	cmdArgs := []string{"/add:TERMSRV/" + hostname, "/user:" + username, "/pass:" + password}
	logging.Log(true, "Executing cmdkey with args:", strings.Join(cmdArgs, " "))

	cmd := exec.Command("cmdkey", cmdArgs...)

	logging.Log(true, "Running cmdkey command...")
	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))

	if err != nil {
		logging.Log(true, "ERROR: cmdkey failed with error:", err, "output:", outputStr)
		return fmt.Errorf("failed to store credential for %s: %w, output: %s", hostname, err, outputStr)
	}

	logging.Log(true, "cmdkey completed successfully, output:", outputStr)
	fmt.Printf("Successfully stored domain credential for %s with user %s\n", hostname, username)
	return nil
}

// getEncryptionKey generates a machine-specific encryption key
func (cm *CredentialManager) getEncryptionKey() ([]byte, error) {
	// Use machine hostname as part of the key generation
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	// Create a 32-byte key from hostname (pad or truncate)
	key := []byte(hostname + "LaunchRDP_Secret_Key_12345")
	if len(key) > 32 {
		key = key[:32]
	} else {
		for len(key) < 32 {
			key = append(key, 'X')
		}
	}
	return key, nil
}

// EncryptPassword encrypts a password using AES
func (cm *CredentialManager) EncryptPassword(password string) (string, error) {
	if password == "" {
		return "", nil
	}

	key, err := cm.getEncryptionKey()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(password), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptPassword decrypts a password using AES
func (cm *CredentialManager) DecryptPassword(encryptedPassword string) (string, error) {
	if encryptedPassword == "" {
		return "", nil
	}

	key, err := cm.getEncryptionKey()
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// StoreGenericCredential stores a generic RDP credential for a username
func (cm *CredentialManager) StoreGenericCredential(username, password string) error {
	// Use a generic target that works with RDP
	target := fmt.Sprintf("TERMSRV/RDP_%s", username)

	cmd := exec.Command("cmdkey", "/generic:"+target, "/user:"+username, "/pass:"+password)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to store generic credential: %w, output: %s", err, string(output))
	}

	fmt.Printf("Successfully stored generic credential for user %s at target %s\n", username, target)
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
	if password == "" {
		return "", nil
	}

	logging.Log(true, "Encrypting password with Windows DPAPI (native)")

	// Convert password to bytes
	passwordBytes := []byte(password)

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

	logging.Log(true, "Password encrypted successfully with DPAPI, length:", len(encrypted))
	return encrypted, nil
}

// DecryptPasswordDPAPI decrypts a password using Windows DPAPI with automatic migration from legacy AES
// Returns both the decrypted password and a flag indicating if migration is needed
func (cm *CredentialManager) DecryptPasswordDPAPI(encryptedPassword string) (string, error) {
	if encryptedPassword == "" {
		return "", nil
	}

	logging.Log(true, "Attempting to decrypt password with Windows DPAPI (native)")

	// First try DPAPI decryption
	decrypted, err := cm.decryptWithDPAPI(encryptedPassword)
	if err == nil {
		logging.Log(true, "Password decrypted successfully with DPAPI - no migration needed")
		return decrypted, nil
	}

	// If DPAPI fails, try legacy AES decryption (migration support)
	logging.Log(true, "DPAPI decryption failed, attempting legacy AES decryption for migration")
	decrypted, err = cm.DecryptPassword(encryptedPassword)
	if err == nil {
		logging.Log(true, "Password decrypted with legacy AES - MIGRATION NEEDED")
		return decrypted, nil
	}

	// Both methods failed
	logging.Log(true, "ERROR: Both DPAPI and legacy AES decryption failed")
	return "", fmt.Errorf("failed to decrypt password with both DPAPI and legacy AES")
}

// DecryptPasswordWithMigration decrypts and automatically migrates legacy passwords
// Returns decrypted password and the new DPAPI-encrypted version if migration occurred
func (cm *CredentialManager) DecryptPasswordWithMigration(encryptedPassword string) (string, string, bool, error) {
	if encryptedPassword == "" {
		return "", "", false, nil
	}

	logging.Log(true, "DecryptPasswordWithMigration: Attempting DPAPI decryption")

	// First try DPAPI decryption
	decrypted, err := cm.decryptWithDPAPI(encryptedPassword)
	if err == nil {
		logging.Log(true, "DecryptPasswordWithMigration: DPAPI decryption successful - no migration needed")
		return decrypted, encryptedPassword, false, nil
	}

	// If DPAPI fails, try legacy AES decryption and migrate
	logging.Log(true, "DecryptPasswordWithMigration: DPAPI failed, trying legacy AES for migration")
	decrypted, err = cm.DecryptPassword(encryptedPassword)
	if err == nil {
		logging.Log(true, "DecryptPasswordWithMigration: Legacy AES successful - migrating to DPAPI")

		// Migrate to DPAPI
		newEncrypted, migrateErr := cm.EncryptPasswordDPAPI(decrypted)
		if migrateErr != nil {
			logging.Log(true, "DecryptPasswordWithMigration: Migration failed:", migrateErr)
			// Still return the decrypted password even if migration fails
			return decrypted, encryptedPassword, false, nil
		}

		logging.Log(true, "DecryptPasswordWithMigration: Migration to DPAPI successful")
		return decrypted, newEncrypted, true, nil
	}

	// Both methods failed
	logging.Log(true, "DecryptPasswordWithMigration: Both DPAPI and legacy AES failed")
	return "", "", false, fmt.Errorf("failed to decrypt password with both DPAPI and legacy AES")
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
// This function should be called when a password is successfully decrypted with legacy AES
func (cm *CredentialManager) MigratePasswordToDPAPI(plainTextPassword string) (string, error) {
	logging.Log(true, "Migrating password from legacy AES to DPAPI")

	// Encrypt with new DPAPI method
	dpapiEncrypted, err := cm.EncryptPasswordDPAPI(plainTextPassword)
	if err != nil {
		logging.Log(true, "ERROR: Failed to migrate password to DPAPI:", err)
		return "", fmt.Errorf("failed to migrate password to DPAPI: %v", err)
	}

	logging.Log(true, "Password migrated successfully to DPAPI")
	return dpapiEncrypted, nil
}

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
	password, newEncrypted, migrated, err := cm.DecryptPasswordWithMigration(encryptedPassword)
	if err != nil {
		logging.Log(true, "StoreCredentialForHostEdit: Failed to decrypt password:", err)
		return fmt.Errorf("failed to decrypt password for host edit: %v", err)
	}

	// Store in Windows CredStore
	logging.Log(true, "StoreCredentialForHostEdit: Storing credential in Windows CredStore")
	if err := cm.StoreCredential(hostname, username, password); err != nil {
		logging.Log(true, "StoreCredentialForHostEdit: Failed to store in CredStore:", err)
		return fmt.Errorf("failed to store credential in CredStore: %v", err)
	}

	if migrated {
		logging.Log(true, "StoreCredentialForHostEdit: Password was migrated during host edit - caller should update JSON")
		// The caller needs to know about the migration to update the JSON
		return fmt.Errorf("MIGRATION_NEEDED:%s", newEncrypted)
	}

	logging.Log(true, "StoreCredentialForHostEdit: Host credential edit completed successfully")
	return nil
}
