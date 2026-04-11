package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// EncryptedPrefix marks encrypted values in config
	EncryptedPrefix = "enc:"

	// SaltSize for PBKDF2
	SaltSize = 16

	// NonceSize for AES-GCM
	NonceSize = 12

	// KeySize for AES-256
	KeySize = 32

	// PBKDF2Iterations for key derivation
	PBKDF2Iterations = 100000
)

// KeyManager handles encryption/decryption of sensitive config values
type KeyManager struct {
	masterKey []byte
}

// NewKeyManager creates a key manager with machine-specific master key
func NewKeyManager() (*KeyManager, error) {
	masterKey, err := deriveMasterKey()
	if err != nil {
		return nil, fmt.Errorf("failed to derive master key: %w", err)
	}

	return &KeyManager{
		masterKey: masterKey,
	}, nil
}

// NewKeyManagerWithPassword creates a key manager with password-based master key
func NewKeyManagerWithPassword(password string) (*KeyManager, error) {
	if password == "" {
		return nil, errors.New("password cannot be empty")
	}

	// Use a fixed salt for password-based key (stored separately if needed)
	salt := []byte("terminalizcrazy-config-salt-v1")
	key := pbkdf2.Key([]byte(password), salt, PBKDF2Iterations, KeySize, sha256.New)

	return &KeyManager{
		masterKey: key,
	}, nil
}

// Encrypt encrypts a plaintext value and returns base64-encoded ciphertext with prefix
func (km *KeyManager) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	// Generate random salt and nonce
	salt := make([]byte, SaltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	nonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Derive encryption key from master key + salt
	key := pbkdf2.Key(km.masterKey, salt, PBKDF2Iterations, KeySize, sha256.New)

	// Create AES-GCM cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Encrypt
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	// Combine: salt + nonce + ciphertext
	combined := make([]byte, SaltSize+NonceSize+len(ciphertext))
	copy(combined[:SaltSize], salt)
	copy(combined[SaltSize:SaltSize+NonceSize], nonce)
	copy(combined[SaltSize+NonceSize:], ciphertext)

	// Encode and add prefix
	return EncryptedPrefix + base64.StdEncoding.EncodeToString(combined), nil
}

// Decrypt decrypts a value (handles both encrypted and plaintext)
func (km *KeyManager) Decrypt(value string) (string, error) {
	if value == "" {
		return "", nil
	}

	// If not encrypted, return as-is (backward compatibility)
	if !strings.HasPrefix(value, EncryptedPrefix) {
		return value, nil
	}

	// Remove prefix and decode
	encoded := strings.TrimPrefix(value, EncryptedPrefix)
	combined, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("failed to decode: %w", err)
	}

	if len(combined) < SaltSize+NonceSize+16 { // 16 = minimum GCM tag size
		return "", errors.New("ciphertext too short")
	}

	// Extract components
	salt := combined[:SaltSize]
	nonce := combined[SaltSize : SaltSize+NonceSize]
	ciphertext := combined[SaltSize+NonceSize:]

	// Derive decryption key
	key := pbkdf2.Key(km.masterKey, salt, PBKDF2Iterations, KeySize, sha256.New)

	// Create AES-GCM cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// IsEncrypted checks if a value is encrypted
func IsEncrypted(value string) bool {
	return strings.HasPrefix(value, EncryptedPrefix)
}

// deriveMasterKey creates a machine-specific master key
func deriveMasterKey() ([]byte, error) {
	// Collect machine-specific entropy
	var entropy strings.Builder

	// Hostname
	hostname, _ := os.Hostname()
	entropy.WriteString(hostname)

	// User home directory
	home, _ := os.UserHomeDir()
	entropy.WriteString(home)

	// OS and architecture
	entropy.WriteString(runtime.GOOS)
	entropy.WriteString(runtime.GOARCH)

	// Environment-based secret (optional, adds extra security)
	envSecret := os.Getenv("TERMINALIZCRAZY_MASTER_KEY")
	if envSecret != "" {
		entropy.WriteString(envSecret)
	}

	// Fixed application salt
	appSalt := []byte("terminalizcrazy-master-key-v1-2024")

	// Derive key using PBKDF2
	key := pbkdf2.Key([]byte(entropy.String()), appSalt, PBKDF2Iterations, KeySize, sha256.New)

	return key, nil
}
