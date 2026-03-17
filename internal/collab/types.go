package collab

import (
	"time"
)

// MessageType represents the type of collaboration message
type MessageType string

const (
	// Control messages
	MsgTypeJoin       MessageType = "join"
	MsgTypeLeave      MessageType = "leave"
	MsgTypeUserList   MessageType = "user_list"
	MsgTypeError      MessageType = "error"
	MsgTypeShareCode  MessageType = "share_code"
	MsgTypeRoomClosed MessageType = "room_closed"

	// Content messages
	MsgTypeChat       MessageType = "chat"
	MsgTypeCommand    MessageType = "command"
	MsgTypeOutput     MessageType = "output"
	MsgTypeTyping     MessageType = "typing"
	MsgTypeCursor     MessageType = "cursor"
	MsgTypeSync       MessageType = "sync"
)

// UserRole represents user permissions in a room
type UserRole string

const (
	RoleOwner  UserRole = "owner"
	RoleEditor UserRole = "editor"
	RoleViewer UserRole = "viewer"
)

// User represents a connected user
type User struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Role     UserRole `json:"role"`
	Color    string   `json:"color"`
	IsTyping bool     `json:"is_typing"`
	JoinedAt time.Time `json:"joined_at"`
}

// Message represents a collaboration message
type Message struct {
	Type      MessageType `json:"type"`
	RoomID    string      `json:"room_id,omitempty"`
	UserID    string      `json:"user_id,omitempty"`
	UserName  string      `json:"user_name,omitempty"`
	Content   string      `json:"content,omitempty"`
	Command   string      `json:"command,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
}

// Room represents a collaboration room
type Room struct {
	ID        string           `json:"id"`
	ShareCode string           `json:"share_code"`
	OwnerID   string           `json:"owner_id"`
	Users     map[string]*User `json:"users"`
	Messages  []Message        `json:"messages"`
	CreatedAt time.Time        `json:"created_at"`
}

// JoinRequest represents a request to join a room
type JoinRequest struct {
	ShareCode string `json:"share_code"`
	UserName  string `json:"user_name"`
}

// SyncData represents full room state for syncing
type SyncData struct {
	Users    []*User   `json:"users"`
	Messages []Message `json:"messages"`
	RoomID   string    `json:"room_id"`
}

// UserColors for distinguishing users
var UserColors = []string{
	"#FF6B6B", // Red
	"#4ECDC4", // Teal
	"#45B7D1", // Blue
	"#96CEB4", // Green
	"#FFEAA7", // Yellow
	"#DDA0DD", // Plum
	"#98D8C8", // Mint
	"#F7DC6F", // Gold
}

// GetUserColor returns a color for a user based on index
func GetUserColor(index int) string {
	return UserColors[index%len(UserColors)]
}
