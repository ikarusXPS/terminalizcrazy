package collab

import (
	"testing"
	"time"
)

func TestGenerateShareCode(t *testing.T) {
	code := generateShareCode()

	if len(code) != 9 { // XXXX-XXXX = 9 chars
		t.Errorf("Expected share code length 9, got %d", len(code))
	}

	if code[4] != '-' {
		t.Error("Expected dash in middle of share code")
	}
}

func TestGenerateID(t *testing.T) {
	id := generateID(8)

	if len(id) != 8 {
		t.Errorf("Expected ID length 8, got %d", len(id))
	}

	// Generate another and ensure they're different
	id2 := generateID(8)
	if id == id2 {
		t.Error("Generated IDs should be unique")
	}
}

func TestGetUserColor(t *testing.T) {
	// Test cycling through colors
	for i := 0; i < len(UserColors)*2; i++ {
		color := GetUserColor(i)
		expected := UserColors[i%len(UserColors)]
		if color != expected {
			t.Errorf("Expected color %s at index %d, got %s", expected, i, color)
		}
	}
}

func TestHubCreateRoom(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	room, err := hub.CreateRoom("user1", "Test User")
	if err != nil {
		t.Fatalf("Failed to create room: %v", err)
	}

	if room.ID == "" {
		t.Error("Room ID should not be empty")
	}

	if room.ShareCode == "" {
		t.Error("Share code should not be empty")
	}

	if room.OwnerID != "user1" {
		t.Errorf("Expected owner ID 'user1', got '%s'", room.OwnerID)
	}

	if len(room.Users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(room.Users))
	}

	owner := room.Users["user1"]
	if owner == nil {
		t.Fatal("Owner should be in users map")
	}

	if owner.Role != RoleOwner {
		t.Errorf("Expected owner role, got %s", owner.Role)
	}
}

func TestHubJoinRoom(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create room
	room, _ := hub.CreateRoom("owner", "Owner")

	// Join room
	joinedRoom, user, err := hub.JoinRoom(room.ShareCode, "user2", "Joiner")
	if err != nil {
		t.Fatalf("Failed to join room: %v", err)
	}

	if joinedRoom.ID != room.ID {
		t.Error("Joined room ID should match")
	}

	if user.ID != "user2" {
		t.Errorf("Expected user ID 'user2', got '%s'", user.ID)
	}

	if user.Role != RoleEditor {
		t.Errorf("Expected editor role for joiner, got %s", user.Role)
	}

	if len(joinedRoom.Users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(joinedRoom.Users))
	}
}

func TestHubJoinInvalidRoom(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	_, _, err := hub.JoinRoom("invalid-code", "user", "User")
	if err == nil {
		t.Error("Should fail to join with invalid share code")
	}
}

func TestHubGetRoom(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	room, _ := hub.CreateRoom("owner", "Owner")

	// Get by ID
	retrieved := hub.GetRoom(room.ID)
	if retrieved == nil {
		t.Fatal("Should find room by ID")
	}

	if retrieved.ID != room.ID {
		t.Error("Retrieved room ID should match")
	}

	// Get by share code
	retrievedByCode := hub.GetRoomByShareCode(room.ShareCode)
	if retrievedByCode == nil {
		t.Fatal("Should find room by share code")
	}

	if retrievedByCode.ShareCode != room.ShareCode {
		t.Error("Retrieved room share code should match")
	}
}

func TestMessageTypes(t *testing.T) {
	// Verify message type constants
	types := []MessageType{
		MsgTypeJoin,
		MsgTypeLeave,
		MsgTypeChat,
		MsgTypeCommand,
		MsgTypeOutput,
		MsgTypeTyping,
	}

	for _, msgType := range types {
		if msgType == "" {
			t.Error("Message type should not be empty")
		}
	}
}

func TestUserRoles(t *testing.T) {
	roles := []UserRole{RoleOwner, RoleEditor, RoleViewer}

	for _, role := range roles {
		if role == "" {
			t.Error("Role should not be empty")
		}
	}
}

func TestRoomMessages(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	room, _ := hub.CreateRoom("owner", "Owner")

	// Add message
	msg := &Message{
		Type:    MsgTypeChat,
		Content: "Hello",
		UserID:  "owner",
	}

	hub.AddMessage(room.ID, msg)

	// Verify message was added
	updatedRoom := hub.GetRoom(room.ID)
	if len(updatedRoom.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(updatedRoom.Messages))
	}

	if updatedRoom.Messages[0].Content != "Hello" {
		t.Error("Message content should match")
	}
}

func TestUser(t *testing.T) {
	user := &User{
		ID:       "test-id",
		Name:     "Test User",
		Role:     RoleEditor,
		Color:    "#FF6B6B",
		IsTyping: false,
		JoinedAt: time.Now(),
	}

	if user.ID != "test-id" {
		t.Error("User ID should match")
	}

	if user.Name != "Test User" {
		t.Error("User name should match")
	}

	if user.Role != RoleEditor {
		t.Error("User role should match")
	}
}

func TestSyncData(t *testing.T) {
	sync := &SyncData{
		RoomID: "room-123",
		Users: []*User{
			{ID: "u1", Name: "User 1"},
			{ID: "u2", Name: "User 2"},
		},
		Messages: []Message{
			{Type: MsgTypeChat, Content: "Hello"},
		},
	}

	if sync.RoomID != "room-123" {
		t.Error("Room ID should match")
	}

	if len(sync.Users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(sync.Users))
	}

	if len(sync.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(sync.Messages))
	}
}
