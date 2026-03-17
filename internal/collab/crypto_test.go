package collab

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCryptoSession(t *testing.T) {
	session, err := NewCryptoSession()

	require.NoError(t, err)
	assert.NotNil(t, session)
	assert.NotNil(t, session.privateKey)
	assert.NotNil(t, session.publicKey)
	assert.False(t, session.isReady)
}

func TestCryptoSession_GetPublicKey(t *testing.T) {
	session, _ := NewCryptoSession()

	pubKey := session.GetPublicKey()

	assert.NotEmpty(t, pubKey)
	// P-256 public key should be 65 bytes when uncompressed, base64 encoded
	assert.True(t, len(pubKey) > 40)
}

func TestCryptoSession_GetKeyExchangePayload(t *testing.T) {
	session, _ := NewCryptoSession()

	payload := session.GetKeyExchangePayload()

	assert.NotNil(t, payload)
	assert.NotEmpty(t, payload.PublicKey)
	assert.Equal(t, session.GetPublicKey(), payload.PublicKey)
}

func TestCryptoSession_SetPeerKey(t *testing.T) {
	session1, _ := NewCryptoSession()
	session2, _ := NewCryptoSession()

	// Exchange keys
	err := session1.SetPeerKey(session2.GetPublicKey())
	require.NoError(t, err)
	assert.True(t, session1.IsReady())

	err = session2.SetPeerKey(session1.GetPublicKey())
	require.NoError(t, err)
	assert.True(t, session2.IsReady())
}

func TestCryptoSession_SetPeerKey_InvalidBase64(t *testing.T) {
	session, _ := NewCryptoSession()

	err := session.SetPeerKey("not-valid-base64!@#$")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "decode")
	assert.False(t, session.IsReady())
}

func TestCryptoSession_SetPeerKey_InvalidKey(t *testing.T) {
	session, _ := NewCryptoSession()

	// Valid base64 but invalid public key bytes
	err := session.SetPeerKey("aGVsbG8gd29ybGQ=") // "hello world" in base64

	assert.Error(t, err)
	assert.False(t, session.IsReady())
}

func TestCryptoSession_IsReady(t *testing.T) {
	session, _ := NewCryptoSession()

	assert.False(t, session.IsReady())

	// Set up key exchange
	peer, _ := NewCryptoSession()
	session.SetPeerKey(peer.GetPublicKey())

	assert.True(t, session.IsReady())
}

func TestCryptoSession_EncryptDecrypt(t *testing.T) {
	// Set up two sessions with key exchange
	session1, _ := NewCryptoSession()
	session2, _ := NewCryptoSession()

	session1.SetPeerKey(session2.GetPublicKey())
	session2.SetPeerKey(session1.GetPublicKey())

	t.Run("encrypt and decrypt bytes", func(t *testing.T) {
		plaintext := []byte("Hello, World!")

		encrypted, err := session1.Encrypt(plaintext)
		require.NoError(t, err)
		assert.NotNil(t, encrypted)
		assert.NotEmpty(t, encrypted.Nonce)
		assert.NotEmpty(t, encrypted.Ciphertext)

		// Decrypt with peer session
		decrypted, err := session2.Decrypt(encrypted)
		require.NoError(t, err)
		assert.Equal(t, plaintext, decrypted)
	})

	t.Run("encrypt and decrypt string", func(t *testing.T) {
		plaintext := "Hello, secure world!"

		encrypted, err := session1.EncryptString(plaintext)
		require.NoError(t, err)

		decrypted, err := session2.DecryptString(encrypted)
		require.NoError(t, err)
		assert.Equal(t, plaintext, decrypted)
	})

	t.Run("encrypt and decrypt JSON", func(t *testing.T) {
		type TestData struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		}

		original := TestData{Name: "test", Value: 42}

		encrypted, err := session1.EncryptJSON(original)
		require.NoError(t, err)

		var decrypted TestData
		err = session2.DecryptJSON(encrypted, &decrypted)
		require.NoError(t, err)
		assert.Equal(t, original, decrypted)
	})
}

func TestCryptoSession_Encrypt_NotReady(t *testing.T) {
	session, _ := NewCryptoSession()

	encrypted, err := session.Encrypt([]byte("test"))

	assert.Error(t, err)
	assert.Nil(t, encrypted)
	assert.Contains(t, err.Error(), "not ready")
}

func TestCryptoSession_Decrypt_NotReady(t *testing.T) {
	session, _ := NewCryptoSession()

	msg := &EncryptedMessage{
		Nonce:      "dGVzdA==",
		Ciphertext: "dGVzdA==",
	}

	plaintext, err := session.Decrypt(msg)

	assert.Error(t, err)
	assert.Nil(t, plaintext)
	assert.Contains(t, err.Error(), "not ready")
}

func TestCryptoSession_Decrypt_InvalidNonce(t *testing.T) {
	session1, _ := NewCryptoSession()
	session2, _ := NewCryptoSession()
	session1.SetPeerKey(session2.GetPublicKey())
	session2.SetPeerKey(session1.GetPublicKey())

	msg := &EncryptedMessage{
		Nonce:      "invalid-base64!",
		Ciphertext: "dGVzdA==",
	}

	plaintext, err := session1.Decrypt(msg)

	assert.Error(t, err)
	assert.Nil(t, plaintext)
	assert.Contains(t, err.Error(), "nonce")
}

func TestCryptoSession_Decrypt_InvalidCiphertext(t *testing.T) {
	session1, _ := NewCryptoSession()
	session2, _ := NewCryptoSession()
	session1.SetPeerKey(session2.GetPublicKey())
	session2.SetPeerKey(session1.GetPublicKey())

	msg := &EncryptedMessage{
		Nonce:      "dGVzdDEyMzQ1Njc4OTAxMg==", // 12 bytes base64 encoded
		Ciphertext: "invalid!base64@#$",
	}

	plaintext, err := session1.Decrypt(msg)

	assert.Error(t, err)
	assert.Nil(t, plaintext)
	assert.Contains(t, err.Error(), "ciphertext")
}

func TestCryptoSession_EncryptMessage(t *testing.T) {
	session1, _ := NewCryptoSession()
	session2, _ := NewCryptoSession()
	session1.SetPeerKey(session2.GetPublicKey())
	session2.SetPeerKey(session1.GetPublicKey())

	msg := &Message{
		Type:    MsgTypeChat,
		Content: "Hello!",
		UserID:  "user-1",
	}

	secure, err := session1.EncryptMessage(msg)

	require.NoError(t, err)
	assert.NotNil(t, secure)
	assert.Equal(t, MsgTypeChat, secure.Type)
	assert.NotNil(t, secure.Encrypted)
}

func TestCryptoSession_DecryptMessage(t *testing.T) {
	session1, _ := NewCryptoSession()
	session2, _ := NewCryptoSession()
	session1.SetPeerKey(session2.GetPublicKey())
	session2.SetPeerKey(session1.GetPublicKey())

	original := &Message{
		Type:    MsgTypeChat,
		Content: "Secret message",
		UserID:  "user-1",
	}

	secure, _ := session1.EncryptMessage(original)

	decrypted, err := session2.DecryptMessage(secure)

	require.NoError(t, err)
	assert.Equal(t, original.Type, decrypted.Type)
	assert.Equal(t, original.Content, decrypted.Content)
	assert.Equal(t, original.UserID, decrypted.UserID)
}

func TestCryptoSession_DecryptMessage_NoEncrypted(t *testing.T) {
	session, _ := NewCryptoSession()
	peer, _ := NewCryptoSession()
	session.SetPeerKey(peer.GetPublicKey())

	secure := &SecureMessage{
		Type:      MsgTypeChat,
		Encrypted: nil,
	}

	msg, err := session.DecryptMessage(secure)

	assert.Error(t, err)
	assert.Nil(t, msg)
	assert.Contains(t, err.Error(), "no encrypted content")
}

func TestCryptoSession_CreateKeyExchangeMessage(t *testing.T) {
	session, _ := NewCryptoSession()

	msg := session.CreateKeyExchangeMessage()

	assert.NotNil(t, msg)
	assert.Equal(t, MsgTypeKeyExchange, msg.Type)
	assert.NotNil(t, msg.KeyExchange)
	assert.Equal(t, session.GetPublicKey(), msg.KeyExchange.PublicKey)
}

func TestCryptoSession_HandleKeyExchange(t *testing.T) {
	session1, _ := NewCryptoSession()
	session2, _ := NewCryptoSession()

	keyExchangeMsg := session1.CreateKeyExchangeMessage()

	err := session2.HandleKeyExchange(keyExchangeMsg)

	require.NoError(t, err)
	assert.True(t, session2.IsReady())
}

func TestCryptoSession_HandleKeyExchange_NoPayload(t *testing.T) {
	session, _ := NewCryptoSession()

	msg := &SecureMessage{
		Type:        MsgTypeKeyExchange,
		KeyExchange: nil,
	}

	err := session.HandleKeyExchange(msg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no key exchange payload")
}

func TestNewCryptoManager(t *testing.T) {
	cm := NewCryptoManager()

	assert.NotNil(t, cm)
	assert.NotNil(t, cm.sessions)
	assert.Empty(t, cm.sessions)
}

func TestCryptoManager_GetOrCreateSession(t *testing.T) {
	cm := NewCryptoManager()

	t.Run("create new session", func(t *testing.T) {
		session, err := cm.GetOrCreateSession("peer-1")

		require.NoError(t, err)
		assert.NotNil(t, session)
		assert.True(t, cm.HasSession("peer-1"))
	})

	t.Run("get existing session", func(t *testing.T) {
		session1, _ := cm.GetOrCreateSession("peer-1")
		session2, _ := cm.GetOrCreateSession("peer-1")

		assert.Same(t, session1, session2)
	})
}

func TestCryptoManager_GetSession(t *testing.T) {
	cm := NewCryptoManager()

	// Non-existent session
	assert.Nil(t, cm.GetSession("nonexistent"))

	// Create and retrieve
	cm.GetOrCreateSession("peer-1")
	session := cm.GetSession("peer-1")
	assert.NotNil(t, session)
}

func TestCryptoManager_RemoveSession(t *testing.T) {
	cm := NewCryptoManager()
	cm.GetOrCreateSession("peer-1")

	assert.True(t, cm.HasSession("peer-1"))

	cm.RemoveSession("peer-1")

	assert.False(t, cm.HasSession("peer-1"))

	// Removing non-existent should not panic
	cm.RemoveSession("peer-2")
}

func TestCryptoManager_HasSession(t *testing.T) {
	cm := NewCryptoManager()

	assert.False(t, cm.HasSession("peer-1"))

	cm.GetOrCreateSession("peer-1")

	assert.True(t, cm.HasSession("peer-1"))
}

func TestCryptoManager_IsSessionReady(t *testing.T) {
	cm := NewCryptoManager()

	// Non-existent session
	assert.False(t, cm.IsSessionReady("peer-1"))

	// Create session but not ready
	cm.GetOrCreateSession("peer-1")
	assert.False(t, cm.IsSessionReady("peer-1"))

	// Set up key exchange
	peerSession, _ := NewCryptoSession()
	session := cm.GetSession("peer-1")
	session.SetPeerKey(peerSession.GetPublicKey())

	assert.True(t, cm.IsSessionReady("peer-1"))
}

func TestCryptoManager_EncryptForPeer(t *testing.T) {
	cm := NewCryptoManager()

	t.Run("no session", func(t *testing.T) {
		msg := &Message{Type: MsgTypeChat, Content: "test"}
		secure, err := cm.EncryptForPeer("nonexistent", msg)

		assert.Error(t, err)
		assert.Nil(t, secure)
		assert.Contains(t, err.Error(), "no session")
	})

	t.Run("with session", func(t *testing.T) {
		session, _ := cm.GetOrCreateSession("peer-1")
		peerSession, _ := NewCryptoSession()
		session.SetPeerKey(peerSession.GetPublicKey())

		msg := &Message{Type: MsgTypeChat, Content: "Hello"}
		secure, err := cm.EncryptForPeer("peer-1", msg)

		require.NoError(t, err)
		assert.NotNil(t, secure)
	})
}

func TestCryptoManager_DecryptFromPeer(t *testing.T) {
	cm1 := NewCryptoManager()
	cm2 := NewCryptoManager()

	t.Run("no session", func(t *testing.T) {
		secure := &SecureMessage{Type: MsgTypeChat}
		msg, err := cm1.DecryptFromPeer("nonexistent", secure)

		assert.Error(t, err)
		assert.Nil(t, msg)
	})

	t.Run("with session", func(t *testing.T) {
		session1, _ := cm1.GetOrCreateSession("peer-2")
		session2, _ := cm2.GetOrCreateSession("peer-1")

		// Key exchange
		session1.SetPeerKey(session2.GetPublicKey())
		session2.SetPeerKey(session1.GetPublicKey())

		// Encrypt with cm1, decrypt with cm2
		original := &Message{Type: MsgTypeChat, Content: "Secret"}
		secure, _ := cm1.EncryptForPeer("peer-2", original)

		decrypted, err := cm2.DecryptFromPeer("peer-1", secure)

		require.NoError(t, err)
		assert.Equal(t, original.Content, decrypted.Content)
	})
}

func TestGenerateRoomKey(t *testing.T) {
	key, err := GenerateRoomKey()

	require.NoError(t, err)
	assert.Len(t, key, 32) // 256 bits = 32 bytes
}

func TestGenerateRoomKey_Uniqueness(t *testing.T) {
	key1, _ := GenerateRoomKey()
	key2, _ := GenerateRoomKey()

	assert.NotEqual(t, key1, key2)
}

func TestDeriveRoomKey(t *testing.T) {
	password := "test-password"
	salt := []byte("test-salt-12345")

	key := DeriveRoomKey(password, salt)

	assert.Len(t, key, 32) // 256 bits = 32 bytes

	// Same inputs should produce same key
	key2 := DeriveRoomKey(password, salt)
	assert.Equal(t, key, key2)

	// Different password should produce different key
	key3 := DeriveRoomKey("different-password", salt)
	assert.NotEqual(t, key, key3)

	// Different salt should produce different key
	key4 := DeriveRoomKey(password, []byte("different-salt"))
	assert.NotEqual(t, key, key4)
}

func TestEncryptedMessage(t *testing.T) {
	msg := EncryptedMessage{
		Nonce:      "test-nonce",
		Ciphertext: "test-ciphertext",
	}

	assert.Equal(t, "test-nonce", msg.Nonce)
	assert.Equal(t, "test-ciphertext", msg.Ciphertext)
}

func TestKeyExchangePayload(t *testing.T) {
	payload := KeyExchangePayload{
		PublicKey: "test-public-key",
	}

	assert.Equal(t, "test-public-key", payload.PublicKey)
}

func TestSecureMessage(t *testing.T) {
	t.Run("with encrypted content", func(t *testing.T) {
		msg := SecureMessage{
			Type: MsgTypeChat,
			Encrypted: &EncryptedMessage{
				Nonce:      "nonce",
				Ciphertext: "cipher",
			},
		}

		assert.Equal(t, MsgTypeChat, msg.Type)
		assert.NotNil(t, msg.Encrypted)
	})

	t.Run("with key exchange", func(t *testing.T) {
		msg := SecureMessage{
			Type: MsgTypeKeyExchange,
			KeyExchange: &KeyExchangePayload{
				PublicKey: "key",
			},
		}

		assert.Equal(t, MsgTypeKeyExchange, msg.Type)
		assert.NotNil(t, msg.KeyExchange)
	})
}
