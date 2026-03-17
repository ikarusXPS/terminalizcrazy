package collab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// CollabClient represents a collaboration client
type CollabClient struct {
	serverURL  string
	conn       *websocket.Conn
	user       *User
	roomID     string
	shareCode  string
	users      map[string]*User
	messages   []Message
	connected  bool
	onMessage  func(*Message)
	onConnect  func()
	onDisconnect func()
	mu         sync.RWMutex
	done       chan struct{}
}

// NewClient creates a new collaboration client
func NewClient(serverURL string) *CollabClient {
	return &CollabClient{
		serverURL: serverURL,
		users:     make(map[string]*User),
		messages:  []Message{},
		done:      make(chan struct{}),
	}
}

// SetMessageHandler sets the callback for incoming messages
func (c *CollabClient) SetMessageHandler(handler func(*Message)) {
	c.onMessage = handler
}

// SetConnectHandler sets the callback for connection events
func (c *CollabClient) SetConnectHandler(handler func()) {
	c.onConnect = handler
}

// SetDisconnectHandler sets the callback for disconnection events
func (c *CollabClient) SetDisconnectHandler(handler func()) {
	c.onDisconnect = handler
}

// CreateRoom creates a new room and connects
func (c *CollabClient) CreateRoom(userID, userName string) (string, error) {
	// Create room via HTTP
	reqBody, _ := json.Marshal(map[string]string{
		"user_id":   userID,
		"user_name": userName,
	})

	resp, err := http.Post(
		fmt.Sprintf("%s/create", c.serverURL),
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create room: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		RoomID    string `json:"room_id"`
		ShareCode string `json:"share_code"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	c.roomID = result.RoomID
	c.shareCode = result.ShareCode
	c.user = &User{
		ID:   userID,
		Name: userName,
		Role: RoleOwner,
	}

	// Connect via WebSocket
	if err := c.connect(); err != nil {
		return "", err
	}

	// Send join message
	c.sendJoin()

	return result.ShareCode, nil
}

// JoinRoom joins an existing room by share code
func (c *CollabClient) JoinRoom(shareCode, userID, userName string) error {
	// Verify room exists via HTTP
	reqBody, _ := json.Marshal(JoinRequest{
		ShareCode: shareCode,
		UserName:  userName,
	})

	resp, err := http.Post(
		fmt.Sprintf("%s/join", c.serverURL),
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return fmt.Errorf("failed to join room: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("room not found")
	}

	var result struct {
		RoomID string `json:"room_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	c.roomID = result.RoomID
	c.shareCode = shareCode
	c.user = &User{
		ID:   userID,
		Name: userName,
		Role: RoleEditor,
	}

	// Connect via WebSocket
	if err := c.connect(); err != nil {
		return err
	}

	// Send join message
	c.sendJoin()

	return nil
}

// connect establishes WebSocket connection
func (c *CollabClient) connect() error {
	// Parse server URL and convert to WebSocket URL
	u, err := url.Parse(c.serverURL)
	if err != nil {
		return err
	}

	scheme := "ws"
	if u.Scheme == "https" {
		scheme = "wss"
	}

	wsURL := fmt.Sprintf("%s://%s/ws", scheme, u.Host)

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("websocket connection failed: %w", err)
	}

	c.conn = conn
	c.connected = true
	c.done = make(chan struct{})

	go c.readPump()
	go c.pingPump()

	if c.onConnect != nil {
		c.onConnect()
	}

	return nil
}

// sendJoin sends a join message
func (c *CollabClient) sendJoin() {
	msg := &Message{
		Type:     MsgTypeJoin,
		RoomID:   c.roomID,
		UserID:   c.user.ID,
		UserName: c.user.Name,
	}
	c.SendMessage(msg)
}

// readPump reads messages from WebSocket
func (c *CollabClient) readPump() {
	defer func() {
		c.connected = false
		c.conn.Close()
		if c.onDisconnect != nil {
			c.onDisconnect()
		}
	}()

	for {
		select {
		case <-c.done:
			return
		default:
			_, data, err := c.conn.ReadMessage()
			if err != nil {
				return
			}

			var msg Message
			if err := json.Unmarshal(data, &msg); err != nil {
				continue
			}

			c.handleMessage(&msg)
		}
	}
}

// pingPump sends periodic pings
func (c *CollabClient) pingPump() {
	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.done:
			return
		case <-ticker.C:
			if c.conn != nil {
				c.conn.WriteMessage(websocket.PingMessage, nil)
			}
		}
	}
}

// handleMessage processes incoming messages
func (c *CollabClient) handleMessage(msg *Message) {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch msg.Type {
	case MsgTypeSync:
		if data, ok := msg.Data.(map[string]interface{}); ok {
			// Parse sync data
			if users, ok := data["users"].([]interface{}); ok {
				c.users = make(map[string]*User)
				for _, u := range users {
					if userData, ok := u.(map[string]interface{}); ok {
						user := &User{
							ID:    getString(userData, "id"),
							Name:  getString(userData, "name"),
							Color: getString(userData, "color"),
						}
						if role, ok := userData["role"].(string); ok {
							user.Role = UserRole(role)
						}
						c.users[user.ID] = user
					}
				}
			}
			if messages, ok := data["messages"].([]interface{}); ok {
				c.messages = make([]Message, 0, len(messages))
				for _, m := range messages {
					if msgData, ok := m.(map[string]interface{}); ok {
						parsedMsg := parseMessageData(msgData)
						c.messages = append(c.messages, parsedMsg)
					}
				}
			}
		}

	case MsgTypeUserList:
		if data, ok := msg.Data.([]interface{}); ok {
			c.users = make(map[string]*User)
			for _, u := range data {
				if userData, ok := u.(map[string]interface{}); ok {
					user := &User{
						ID:    getString(userData, "id"),
						Name:  getString(userData, "name"),
						Color: getString(userData, "color"),
					}
					if role, ok := userData["role"].(string); ok {
						user.Role = UserRole(role)
					}
					c.users[user.ID] = user
				}
			}
		}

	case MsgTypeJoin:
		user := &User{
			ID:       msg.UserID,
			Name:     msg.UserName,
			JoinedAt: msg.Timestamp,
		}
		c.users[msg.UserID] = user

	case MsgTypeLeave:
		delete(c.users, msg.UserID)

	case MsgTypeChat, MsgTypeCommand, MsgTypeOutput:
		c.messages = append(c.messages, *msg)

	case MsgTypeTyping:
		if user, ok := c.users[msg.UserID]; ok {
			user.IsTyping = msg.Content == "true"
		}
	}

	// Call message handler
	if c.onMessage != nil {
		c.onMessage(msg)
	}
}

// SendMessage sends a message to the server
func (c *CollabClient) SendMessage(msg *Message) error {
	if !c.connected || c.conn == nil {
		return fmt.Errorf("not connected")
	}

	msg.Timestamp = time.Now()
	msg.RoomID = c.roomID
	if c.user != nil {
		msg.UserID = c.user.ID
		msg.UserName = c.user.Name
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return c.conn.WriteMessage(websocket.TextMessage, data)
}

// SendChat sends a chat message
func (c *CollabClient) SendChat(content string) error {
	return c.SendMessage(&Message{
		Type:    MsgTypeChat,
		Content: content,
	})
}

// SendCommand sends a command message
func (c *CollabClient) SendCommand(command string) error {
	return c.SendMessage(&Message{
		Type:    MsgTypeCommand,
		Command: command,
	})
}

// SendOutput sends command output
func (c *CollabClient) SendOutput(content string) error {
	return c.SendMessage(&Message{
		Type:    MsgTypeOutput,
		Content: content,
	})
}

// SendTyping sends typing indicator
func (c *CollabClient) SendTyping(isTyping bool) error {
	content := "false"
	if isTyping {
		content = "true"
	}
	return c.SendMessage(&Message{
		Type:    MsgTypeTyping,
		Content: content,
	})
}

// Disconnect closes the connection
func (c *CollabClient) Disconnect() {
	if c.done != nil {
		close(c.done)
	}
	if c.conn != nil {
		c.conn.Close()
	}
	c.connected = false
}

// IsConnected returns connection status
func (c *CollabClient) IsConnected() bool {
	return c.connected
}

// GetShareCode returns the current share code
func (c *CollabClient) GetShareCode() string {
	return c.shareCode
}

// GetRoomID returns the current room ID
func (c *CollabClient) GetRoomID() string {
	return c.roomID
}

// GetUsers returns the current user list
func (c *CollabClient) GetUsers() []*User {
	c.mu.RLock()
	defer c.mu.RUnlock()

	users := make([]*User, 0, len(c.users))
	for _, u := range c.users {
		users = append(users, u)
	}
	return users
}

// GetMessages returns the message history
func (c *CollabClient) GetMessages() []Message {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return append([]Message{}, c.messages...)
}

// GetUser returns the current user
func (c *CollabClient) GetUser() *User {
	return c.user
}

// Helper functions
func getString(data map[string]interface{}, key string) string {
	if v, ok := data[key].(string); ok {
		return v
	}
	return ""
}

func parseMessageData(data map[string]interface{}) Message {
	msg := Message{
		Content:  getString(data, "content"),
		Command:  getString(data, "command"),
		UserID:   getString(data, "user_id"),
		UserName: getString(data, "user_name"),
	}
	if typeStr, ok := data["type"].(string); ok {
		msg.Type = MessageType(typeStr)
	}
	return msg
}
