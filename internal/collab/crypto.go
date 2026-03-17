package collab

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"sync"
)

// CryptoSession handles end-to-end encryption for a collaboration session
type CryptoSession struct {
	privateKey *ecdh.PrivateKey
	publicKey  *ecdh.PublicKey
	sharedKey  []byte
	gcm        cipher.AEAD
	peerKey    *ecdh.PublicKey
	isReady    bool
	mu         sync.RWMutex
}

// EncryptedMessage represents an encrypted message
type EncryptedMessage struct {
	Nonce      string `json:"nonce"`      // Base64-encoded nonce
	Ciphertext string `json:"ciphertext"` // Base64-encoded encrypted data
}

// KeyExchangePayload represents the public key exchange payload
type KeyExchangePayload struct {
	PublicKey string `json:"public_key"` // Base64-encoded public key
}

// NewCryptoSession creates a new crypto session
func NewCryptoSession() (*CryptoSession, error) {
	// Generate ECDH key pair using P-256 curve
	curve := ecdh.P256()
	privateKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %w", err)
	}

	return &CryptoSession{
		privateKey: privateKey,
		publicKey:  privateKey.PublicKey(),
	}, nil
}

// GetPublicKey returns the base64-encoded public key for key exchange
func (cs *CryptoSession) GetPublicKey() string {
	return base64.StdEncoding.EncodeToString(cs.publicKey.Bytes())
}

// GetKeyExchangePayload returns the key exchange payload
func (cs *CryptoSession) GetKeyExchangePayload() *KeyExchangePayload {
	return &KeyExchangePayload{
		PublicKey: cs.GetPublicKey(),
	}
}

// SetPeerKey sets the peer's public key and derives the shared secret
func (cs *CryptoSession) SetPeerKey(peerKeyB64 string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Decode peer's public key
	peerKeyBytes, err := base64.StdEncoding.DecodeString(peerKeyB64)
	if err != nil {
		return fmt.Errorf("failed to decode peer public key: %w", err)
	}

	// Parse public key
	curve := ecdh.P256()
	peerKey, err := curve.NewPublicKey(peerKeyBytes)
	if err != nil {
		return fmt.Errorf("failed to parse peer public key: %w", err)
	}

	cs.peerKey = peerKey

	// Perform ECDH to derive shared secret
	sharedSecret, err := cs.privateKey.ECDH(peerKey)
	if err != nil {
		return fmt.Errorf("failed to derive shared secret: %w", err)
	}

	// Derive AES-256 key from shared secret using SHA-256
	hash := sha256.Sum256(sharedSecret)
	cs.sharedKey = hash[:]

	// Initialize AES-GCM
	block, err := aes.NewCipher(cs.sharedKey)
	if err != nil {
		return fmt.Errorf("failed to create AES cipher: %w", err)
	}

	cs.gcm, err = cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	cs.isReady = true
	return nil
}

// IsReady returns true if the crypto session is ready for encryption/decryption
func (cs *CryptoSession) IsReady() bool {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.isReady
}

// Encrypt encrypts plaintext and returns an EncryptedMessage
func (cs *CryptoSession) Encrypt(plaintext []byte) (*EncryptedMessage, error) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	if !cs.isReady {
		return nil, fmt.Errorf("crypto session not ready - key exchange not completed")
	}

	// Generate random nonce
	nonce := make([]byte, cs.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt
	ciphertext := cs.gcm.Seal(nil, nonce, plaintext, nil)

	return &EncryptedMessage{
		Nonce:      base64.StdEncoding.EncodeToString(nonce),
		Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
	}, nil
}

// Decrypt decrypts an EncryptedMessage and returns the plaintext
func (cs *CryptoSession) Decrypt(msg *EncryptedMessage) ([]byte, error) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	if !cs.isReady {
		return nil, fmt.Errorf("crypto session not ready - key exchange not completed")
	}

	// Decode nonce
	nonce, err := base64.StdEncoding.DecodeString(msg.Nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to decode nonce: %w", err)
	}

	// Decode ciphertext
	ciphertext, err := base64.StdEncoding.DecodeString(msg.Ciphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	// Decrypt
	plaintext, err := cs.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// EncryptString encrypts a string
func (cs *CryptoSession) EncryptString(plaintext string) (*EncryptedMessage, error) {
	return cs.Encrypt([]byte(plaintext))
}

// DecryptString decrypts to a string
func (cs *CryptoSession) DecryptString(msg *EncryptedMessage) (string, error) {
	plaintext, err := cs.Decrypt(msg)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// EncryptJSON encrypts a JSON-serializable object
func (cs *CryptoSession) EncryptJSON(obj interface{}) (*EncryptedMessage, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return cs.Encrypt(data)
}

// DecryptJSON decrypts to a JSON object
func (cs *CryptoSession) DecryptJSON(msg *EncryptedMessage, target interface{}) error {
	plaintext, err := cs.Decrypt(msg)
	if err != nil {
		return err
	}
	return json.Unmarshal(plaintext, target)
}

// SecureMessage wraps a Message with encryption
type SecureMessage struct {
	Type      MessageType       `json:"type"`
	Encrypted *EncryptedMessage `json:"encrypted,omitempty"`
	KeyExchange *KeyExchangePayload `json:"key_exchange,omitempty"`
}

// EncryptMessage encrypts a collaboration message
func (cs *CryptoSession) EncryptMessage(msg *Message) (*SecureMessage, error) {
	encrypted, err := cs.EncryptJSON(msg)
	if err != nil {
		return nil, err
	}

	return &SecureMessage{
		Type:      msg.Type,
		Encrypted: encrypted,
	}, nil
}

// DecryptMessage decrypts a secure message
func (cs *CryptoSession) DecryptMessage(secure *SecureMessage) (*Message, error) {
	if secure.Encrypted == nil {
		return nil, fmt.Errorf("message has no encrypted content")
	}

	var msg Message
	if err := cs.DecryptJSON(secure.Encrypted, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

// CreateKeyExchangeMessage creates a key exchange message
func (cs *CryptoSession) CreateKeyExchangeMessage() *SecureMessage {
	return &SecureMessage{
		Type:        MsgTypeKeyExchange,
		KeyExchange: cs.GetKeyExchangePayload(),
	}
}

// HandleKeyExchange processes a key exchange message
func (cs *CryptoSession) HandleKeyExchange(secure *SecureMessage) error {
	if secure.KeyExchange == nil {
		return fmt.Errorf("message has no key exchange payload")
	}
	return cs.SetPeerKey(secure.KeyExchange.PublicKey)
}

// MsgTypeKeyExchange is the message type for key exchange
const MsgTypeKeyExchange MessageType = "key_exchange"

// CryptoManager manages crypto sessions for multiple peers
type CryptoManager struct {
	sessions map[string]*CryptoSession
	mu       sync.RWMutex
}

// NewCryptoManager creates a new crypto manager
func NewCryptoManager() *CryptoManager {
	return &CryptoManager{
		sessions: make(map[string]*CryptoSession),
	}
}

// GetOrCreateSession gets or creates a crypto session for a peer
func (cm *CryptoManager) GetOrCreateSession(peerID string) (*CryptoSession, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if session, ok := cm.sessions[peerID]; ok {
		return session, nil
	}

	session, err := NewCryptoSession()
	if err != nil {
		return nil, err
	}

	cm.sessions[peerID] = session
	return session, nil
}

// GetSession gets an existing session
func (cm *CryptoManager) GetSession(peerID string) *CryptoSession {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.sessions[peerID]
}

// RemoveSession removes a session
func (cm *CryptoManager) RemoveSession(peerID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.sessions, peerID)
}

// HasSession checks if a session exists
func (cm *CryptoManager) HasSession(peerID string) bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	_, ok := cm.sessions[peerID]
	return ok
}

// IsSessionReady checks if a session is ready for encryption
func (cm *CryptoManager) IsSessionReady(peerID string) bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	session, ok := cm.sessions[peerID]
	return ok && session.IsReady()
}

// EncryptForPeer encrypts a message for a specific peer
func (cm *CryptoManager) EncryptForPeer(peerID string, msg *Message) (*SecureMessage, error) {
	session := cm.GetSession(peerID)
	if session == nil {
		return nil, fmt.Errorf("no session for peer: %s", peerID)
	}
	return session.EncryptMessage(msg)
}

// DecryptFromPeer decrypts a message from a specific peer
func (cm *CryptoManager) DecryptFromPeer(peerID string, secure *SecureMessage) (*Message, error) {
	session := cm.GetSession(peerID)
	if session == nil {
		return nil, fmt.Errorf("no session for peer: %s", peerID)
	}
	return session.DecryptMessage(secure)
}

// GenerateRoomKey generates a random room-wide encryption key
func GenerateRoomKey() ([]byte, error) {
	key := make([]byte, 32) // 256-bit key
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("failed to generate room key: %w", err)
	}
	return key, nil
}

// DeriveRoomKey derives a room key from a password/passphrase
func DeriveRoomKey(password string, salt []byte) []byte {
	// Simple PBKDF-like derivation (for demo - use proper PBKDF2 in production)
	combined := append([]byte(password), salt...)
	hash := sha256.Sum256(combined)
	return hash[:]
}
