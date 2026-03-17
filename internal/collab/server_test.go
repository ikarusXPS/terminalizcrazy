package collab

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHub(t *testing.T) {
	hub := NewHub()

	assert.NotNil(t, hub)
	assert.NotNil(t, hub.rooms)
	assert.NotNil(t, hub.clients)
	assert.NotNil(t, hub.register)
	assert.NotNil(t, hub.unregister)
	assert.NotNil(t, hub.broadcast)
}

func TestHub_CreateRoom(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	room, err := hub.CreateRoom("owner-1", "Owner Name")

	require.NoError(t, err)
	assert.NotNil(t, room)
	assert.NotEmpty(t, room.ID)
	assert.NotEmpty(t, room.ShareCode)
	assert.Equal(t, "owner-1", room.OwnerID)
	assert.Len(t, room.Users, 1)

	owner := room.Users["owner-1"]
	require.NotNil(t, owner)
	assert.Equal(t, "Owner Name", owner.Name)
	assert.Equal(t, RoleOwner, owner.Role)
	assert.NotEmpty(t, owner.Color)
}

func TestHub_JoinRoom(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	room, _ := hub.CreateRoom("owner-1", "Owner")

	t.Run("join existing room", func(t *testing.T) {
		joinedRoom, user, err := hub.JoinRoom(room.ShareCode, "user-2", "Joiner")

		require.NoError(t, err)
		assert.Equal(t, room.ID, joinedRoom.ID)
		assert.Equal(t, "user-2", user.ID)
		assert.Equal(t, "Joiner", user.Name)
		assert.Equal(t, RoleEditor, user.Role)
	})

	t.Run("join nonexistent room", func(t *testing.T) {
		_, _, err := hub.JoinRoom("invalid-code", "user-3", "User3")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("user already in room", func(t *testing.T) {
		_, _, err := hub.JoinRoom(room.ShareCode, "user-2", "Duplicate")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already in room")
	})
}

func TestHub_GetRoom(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	room, _ := hub.CreateRoom("owner-1", "Owner")

	t.Run("by ID", func(t *testing.T) {
		found := hub.GetRoom(room.ID)
		assert.NotNil(t, found)
		assert.Equal(t, room.ID, found.ID)
	})

	t.Run("nonexistent ID", func(t *testing.T) {
		found := hub.GetRoom("nonexistent")
		assert.Nil(t, found)
	})
}

func TestHub_GetRoomByShareCode(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	room, _ := hub.CreateRoom("owner-1", "Owner")

	t.Run("existing share code", func(t *testing.T) {
		found := hub.GetRoomByShareCode(room.ShareCode)
		assert.NotNil(t, found)
		assert.Equal(t, room.ShareCode, found.ShareCode)
	})

	t.Run("nonexistent share code", func(t *testing.T) {
		found := hub.GetRoomByShareCode("xxxx-xxxx")
		assert.Nil(t, found)
	})
}

func TestHub_AddMessage(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	room, _ := hub.CreateRoom("owner-1", "Owner")

	msg := &Message{
		Type:    MsgTypeChat,
		Content: "Hello everyone!",
		UserID:  "owner-1",
	}

	hub.AddMessage(room.ID, msg)

	// Give some time for the message to be processed
	time.Sleep(10 * time.Millisecond)

	updatedRoom := hub.GetRoom(room.ID)
	require.Len(t, updatedRoom.Messages, 1)
	assert.Equal(t, "Hello everyone!", updatedRoom.Messages[0].Content)
	assert.NotZero(t, updatedRoom.Messages[0].Timestamp)
}

func TestHub_AddMessage_TrimHistory(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	room, _ := hub.CreateRoom("owner-1", "Owner")

	// Add more than 100 messages
	for i := 0; i < 105; i++ {
		hub.AddMessage(room.ID, &Message{
			Type:    MsgTypeChat,
			Content: "Message",
		})
	}

	time.Sleep(20 * time.Millisecond)

	updatedRoom := hub.GetRoom(room.ID)
	assert.LessOrEqual(t, len(updatedRoom.Messages), 100)
}

func TestNewServer(t *testing.T) {
	server := NewServer(8080)

	assert.NotNil(t, server)
	assert.Equal(t, 8080, server.port)
	assert.NotNil(t, server.hub)
	assert.NotNil(t, server.server)
}

func TestServer_GetHub(t *testing.T) {
	server := NewServer(8080)

	hub := server.GetHub()

	assert.NotNil(t, hub)
	assert.Same(t, server.hub, hub)
}

func TestServer_handleHealth(t *testing.T) {
	server := NewServer(8080)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	server.handleHealth(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "ok", response["status"])
}

func TestServer_handleCreate(t *testing.T) {
	server := NewServer(8080)
	go server.hub.Run()

	t.Run("valid request", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"user_id":   "user-1",
			"user_name": "Test User",
		})

		req := httptest.NewRequest("POST", "/create", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		server.handleCreate(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotEmpty(t, response["room_id"])
		assert.NotEmpty(t, response["share_code"])
	})

	t.Run("invalid method", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/create", nil)
		w := httptest.NewRecorder()

		server.handleCreate(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/create", bytes.NewReader([]byte("invalid")))
		w := httptest.NewRecorder()

		server.handleCreate(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestServer_handleJoin(t *testing.T) {
	server := NewServer(8080)
	go server.hub.Run()

	// Create a room first
	room, _ := server.hub.CreateRoom("owner-1", "Owner")

	t.Run("valid request", func(t *testing.T) {
		body, _ := json.Marshal(JoinRequest{
			ShareCode: room.ShareCode,
			UserName:  "Joiner",
		})

		req := httptest.NewRequest("POST", "/join", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		server.handleJoin(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, room.ID, response["room_id"])
		assert.Equal(t, "ok", response["status"])
	})

	t.Run("invalid method", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/join", nil)
		w := httptest.NewRecorder()

		server.handleJoin(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/join", bytes.NewReader([]byte("invalid")))
		w := httptest.NewRecorder()

		server.handleJoin(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("room not found", func(t *testing.T) {
		body, _ := json.Marshal(JoinRequest{
			ShareCode: "invalid-code",
			UserName:  "Joiner",
		})

		req := httptest.NewRequest("POST", "/join", bytes.NewReader(body))
		w := httptest.NewRecorder()

		server.handleJoin(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestMessage(t *testing.T) {
	msg := Message{
		Type:      MsgTypeChat,
		RoomID:    "room-1",
		UserID:    "user-1",
		UserName:  "Test User",
		Content:   "Hello!",
		Command:   "echo hello",
		Timestamp: time.Now(),
		Data:      map[string]string{"key": "value"},
	}

	assert.Equal(t, MsgTypeChat, msg.Type)
	assert.Equal(t, "room-1", msg.RoomID)
	assert.Equal(t, "user-1", msg.UserID)
	assert.Equal(t, "Test User", msg.UserName)
	assert.Equal(t, "Hello!", msg.Content)
	assert.Equal(t, "echo hello", msg.Command)
	assert.NotZero(t, msg.Timestamp)
	assert.NotNil(t, msg.Data)
}

func TestRoom(t *testing.T) {
	room := Room{
		ID:        "room-123",
		ShareCode: "abcd-efgh",
		OwnerID:   "owner-1",
		Users:     make(map[string]*User),
		Messages:  []Message{},
		CreatedAt: time.Now(),
	}

	room.Users["user-1"] = &User{
		ID:   "user-1",
		Name: "User 1",
		Role: RoleOwner,
	}

	assert.Equal(t, "room-123", room.ID)
	assert.Equal(t, "abcd-efgh", room.ShareCode)
	assert.Equal(t, "owner-1", room.OwnerID)
	assert.Len(t, room.Users, 1)
	assert.NotZero(t, room.CreatedAt)
}

func TestJoinRequest(t *testing.T) {
	req := JoinRequest{
		ShareCode: "abcd-efgh",
		UserName:  "Test User",
	}

	assert.Equal(t, "abcd-efgh", req.ShareCode)
	assert.Equal(t, "Test User", req.UserName)
}

func TestUserColors_AllValid(t *testing.T) {
	assert.Len(t, UserColors, 8)

	for i, color := range UserColors {
		assert.NotEmpty(t, color, "color at index %d should not be empty", i)
		assert.True(t, color[0] == '#', "color should start with #")
	}
}

func TestMessageTypeConstants(t *testing.T) {
	types := []struct {
		msgType  MessageType
		expected string
	}{
		{MsgTypeJoin, "join"},
		{MsgTypeLeave, "leave"},
		{MsgTypeUserList, "user_list"},
		{MsgTypeError, "error"},
		{MsgTypeShareCode, "share_code"},
		{MsgTypeRoomClosed, "room_closed"},
		{MsgTypeChat, "chat"},
		{MsgTypeCommand, "command"},
		{MsgTypeOutput, "output"},
		{MsgTypeTyping, "typing"},
		{MsgTypeCursor, "cursor"},
		{MsgTypeSync, "sync"},
	}

	for _, tt := range types {
		assert.Equal(t, MessageType(tt.expected), tt.msgType)
	}
}

func TestUserRoleConstants(t *testing.T) {
	assert.Equal(t, UserRole("owner"), RoleOwner)
	assert.Equal(t, UserRole("editor"), RoleEditor)
	assert.Equal(t, UserRole("viewer"), RoleViewer)
}

func TestUser_Fields(t *testing.T) {
	now := time.Now()
	user := User{
		ID:       "user-123",
		Name:     "Test User",
		Role:     RoleEditor,
		Color:    "#FF6B6B",
		IsTyping: true,
		JoinedAt: now,
	}

	assert.Equal(t, "user-123", user.ID)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, RoleEditor, user.Role)
	assert.Equal(t, "#FF6B6B", user.Color)
	assert.True(t, user.IsTyping)
	assert.Equal(t, now, user.JoinedAt)
}

func TestSyncData_Fields(t *testing.T) {
	users := []*User{
		{ID: "user-1", Name: "User 1"},
		{ID: "user-2", Name: "User 2"},
	}
	messages := []Message{
		{Type: MsgTypeChat, Content: "Hello"},
	}

	sync := SyncData{
		Users:    users,
		Messages: messages,
		RoomID:   "room-123",
	}

	assert.Equal(t, "room-123", sync.RoomID)
	assert.Len(t, sync.Users, 2)
	assert.Len(t, sync.Messages, 1)
}
