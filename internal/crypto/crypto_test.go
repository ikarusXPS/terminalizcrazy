package crypto

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyManager_EncryptDecrypt(t *testing.T) {
	km, err := NewKeyManager()
	require.NoError(t, err)

	tests := []struct {
		name      string
		plaintext string
	}{
		{"empty string", ""},
		{"simple key", "sk-ant-api03-test123"},
		{"long key", "AIzaSyAfD8afWpo9AfcflRchlcKhmiYfx-cDXXE-with-extra-data"},
		{"special chars", "sk-test!@#$%^&*()_+-=[]{}|;':\",./<>?"},
		{"unicode", "sk-test-日本語-émoji-🔐"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := km.Encrypt(tt.plaintext)
			require.NoError(t, err)

			if tt.plaintext == "" {
				assert.Equal(t, "", encrypted)
				return
			}

			// Verify it has the prefix
			assert.True(t, strings.HasPrefix(encrypted, EncryptedPrefix))

			// Decrypt and verify
			decrypted, err := km.Decrypt(encrypted)
			require.NoError(t, err)
			assert.Equal(t, tt.plaintext, decrypted)
		})
	}
}

func TestKeyManager_DecryptPlaintext(t *testing.T) {
	km, err := NewKeyManager()
	require.NoError(t, err)

	// Plaintext values should be returned as-is (backward compatibility)
	plaintext := "sk-ant-api03-plaintext"
	result, err := km.Decrypt(plaintext)
	require.NoError(t, err)
	assert.Equal(t, plaintext, result)
}

func TestKeyManager_DifferentEncryptions(t *testing.T) {
	km, err := NewKeyManager()
	require.NoError(t, err)

	plaintext := "sk-ant-api03-test123"

	// Encrypt twice - should produce different ciphertexts (random nonce)
	enc1, err := km.Encrypt(plaintext)
	require.NoError(t, err)

	enc2, err := km.Encrypt(plaintext)
	require.NoError(t, err)

	assert.NotEqual(t, enc1, enc2, "same plaintext should produce different ciphertexts")

	// Both should decrypt to the same value
	dec1, err := km.Decrypt(enc1)
	require.NoError(t, err)

	dec2, err := km.Decrypt(enc2)
	require.NoError(t, err)

	assert.Equal(t, plaintext, dec1)
	assert.Equal(t, plaintext, dec2)
}

func TestKeyManagerWithPassword(t *testing.T) {
	km1, err := NewKeyManagerWithPassword("test-password-123")
	require.NoError(t, err)

	km2, err := NewKeyManagerWithPassword("test-password-123")
	require.NoError(t, err)

	plaintext := "sk-ant-api03-secret"

	// Encrypt with km1
	encrypted, err := km1.Encrypt(plaintext)
	require.NoError(t, err)

	// Decrypt with km2 (same password)
	decrypted, err := km2.Decrypt(encrypted)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestKeyManagerWithPassword_WrongPassword(t *testing.T) {
	km1, err := NewKeyManagerWithPassword("correct-password")
	require.NoError(t, err)

	km2, err := NewKeyManagerWithPassword("wrong-password")
	require.NoError(t, err)

	plaintext := "sk-ant-api03-secret"

	// Encrypt with correct password
	encrypted, err := km1.Encrypt(plaintext)
	require.NoError(t, err)

	// Decrypt with wrong password should fail
	_, err = km2.Decrypt(encrypted)
	assert.Error(t, err)
}

func TestKeyManagerWithPassword_EmptyPassword(t *testing.T) {
	_, err := NewKeyManagerWithPassword("")
	assert.Error(t, err)
}

func TestIsEncrypted(t *testing.T) {
	tests := []struct {
		value    string
		expected bool
	}{
		{"", false},
		{"sk-ant-api03-plain", false},
		{"enc:", true},
		{"enc:abcdef123456", true},
		{"ENC:not-encrypted", false}, // case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsEncrypted(tt.value))
		})
	}
}

func TestKeyManager_InvalidCiphertext(t *testing.T) {
	km, err := NewKeyManager()
	require.NoError(t, err)

	tests := []struct {
		name  string
		value string
	}{
		{"invalid base64", "enc:not-valid-base64!!!"},
		{"too short", "enc:YWJj"}, // "abc" in base64
		{"corrupted", "enc:YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXo="}, // random data
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := km.Decrypt(tt.value)
			assert.Error(t, err)
		})
	}
}
