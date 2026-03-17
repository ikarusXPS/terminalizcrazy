package collab

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for local use
	},
}

// Client represents a WebSocket client
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	user   *User
	roomID string
}

// Hub manages all rooms and clients
type Hub struct {
	rooms      map[string]*Room
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	mu         sync.RWMutex
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]*Room),
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				// Remove user from room
				if client.roomID != "" && client.user != nil {
					if room, ok := h.rooms[client.roomID]; ok {
						delete(room.Users, client.user.ID)

						// Notify others
						h.broadcastToRoom(client.roomID, &Message{
							Type:      MsgTypeLeave,
							UserID:    client.user.ID,
							UserName:  client.user.Name,
							Timestamp: time.Now(),
						})

						// Send updated user list
						h.sendUserList(client.roomID)

						// Close room if empty
						if len(room.Users) == 0 {
							delete(h.rooms, client.roomID)
						}
					}
				}
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.broadcastToRoom(msg.RoomID, msg)
		}
	}
}

// CreateRoom creates a new collaboration room
func (h *Hub) CreateRoom(ownerID, ownerName string) (*Room, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	roomID := generateID(8)
	shareCode := generateShareCode()

	room := &Room{
		ID:        roomID,
		ShareCode: shareCode,
		OwnerID:   ownerID,
		Users:     make(map[string]*User),
		Messages:  []Message{},
		CreatedAt: time.Now(),
	}

	// Add owner as first user
	room.Users[ownerID] = &User{
		ID:       ownerID,
		Name:     ownerName,
		Role:     RoleOwner,
		Color:    GetUserColor(0),
		JoinedAt: time.Now(),
	}

	h.rooms[roomID] = room
	return room, nil
}

// JoinRoom joins an existing room by share code
func (h *Hub) JoinRoom(shareCode, userID, userName string) (*Room, *User, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Find room by share code
	var room *Room
	for _, r := range h.rooms {
		if r.ShareCode == shareCode {
			room = r
			break
		}
	}

	if room == nil {
		return nil, nil, fmt.Errorf("room not found")
	}

	// Check if user already exists
	if _, exists := room.Users[userID]; exists {
		return nil, nil, fmt.Errorf("user already in room")
	}

	// Add user
	user := &User{
		ID:       userID,
		Name:     userName,
		Role:     RoleEditor, // Default role for joiners
		Color:    GetUserColor(len(room.Users)),
		JoinedAt: time.Now(),
	}

	room.Users[userID] = user
	return room, user, nil
}

// GetRoom returns a room by ID
func (h *Hub) GetRoom(roomID string) *Room {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.rooms[roomID]
}

// GetRoomByShareCode returns a room by share code
func (h *Hub) GetRoomByShareCode(shareCode string) *Room {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, r := range h.rooms {
		if r.ShareCode == shareCode {
			return r
		}
	}
	return nil
}

// broadcastToRoom sends a message to all clients in a room
func (h *Hub) broadcastToRoom(roomID string, msg *Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.roomID == roomID {
			select {
			case client.send <- data:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}

// sendUserList sends the current user list to all clients in a room
func (h *Hub) sendUserList(roomID string) {
	room := h.rooms[roomID]
	if room == nil {
		return
	}

	users := make([]*User, 0, len(room.Users))
	for _, u := range room.Users {
		users = append(users, u)
	}

	msg := &Message{
		Type:      MsgTypeUserList,
		RoomID:    roomID,
		Timestamp: time.Now(),
		Data:      users,
	}

	h.broadcastToRoom(roomID, msg)
}

// AddMessage adds a message to room history and broadcasts it
func (h *Hub) AddMessage(roomID string, msg *Message) {
	h.mu.Lock()
	room := h.rooms[roomID]
	if room != nil {
		msg.Timestamp = time.Now()
		room.Messages = append(room.Messages, *msg)

		// Keep only last 100 messages
		if len(room.Messages) > 100 {
			room.Messages = room.Messages[len(room.Messages)-100:]
		}
	}
	h.mu.Unlock()

	h.broadcastToRoom(roomID, msg)
}

// Server handles HTTP and WebSocket connections
type Server struct {
	hub    *Hub
	port   int
	server *http.Server
}

// NewServer creates a new collaboration server
func NewServer(port int) *Server {
	hub := NewHub()

	s := &Server{
		hub:  hub,
		port: port,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.handleWebSocket)
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/create", s.handleCreate)
	mux.HandleFunc("/join", s.handleJoin)

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return s
}

// Start starts the server
func (s *Server) Start() error {
	go s.hub.Run()
	log.Printf("Collaboration server starting on port %d", s.port)
	return s.server.ListenAndServe()
}

// Stop stops the server
func (s *Server) Stop() error {
	return s.server.Close()
}

// GetHub returns the hub
func (s *Server) GetHub() *Hub {
	return s.hub
}

// handleHealth handles health check requests
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// handleCreate handles room creation requests
func (s *Server) handleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID   string `json:"user_id"`
		UserName string `json:"user_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	room, err := s.hub.CreateRoom(req.UserID, req.UserName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"room_id":    room.ID,
		"share_code": room.ShareCode,
	})
}

// handleJoin handles join requests
func (s *Server) handleJoin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req JoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	room := s.hub.GetRoomByShareCode(req.ShareCode)
	if room == nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"room_id": room.ID,
		"status":  "ok",
	})
}

// handleWebSocket handles WebSocket connections
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:  s.hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	s.hub.register <- client

	go client.writePump()
	go client.readPump()
}

// readPump handles incoming messages from client
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(65536)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}

		c.handleMessage(&msg)
	}
}

// writePump handles outgoing messages to client
func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage processes incoming messages
func (c *Client) handleMessage(msg *Message) {
	switch msg.Type {
	case MsgTypeJoin:
		c.handleJoin(msg)
	case MsgTypeChat, MsgTypeCommand, MsgTypeOutput:
		c.handleContent(msg)
	case MsgTypeTyping:
		c.handleTyping(msg)
	}
}

// handleJoin processes join messages
func (c *Client) handleJoin(msg *Message) {
	room := c.hub.GetRoom(msg.RoomID)
	if room == nil {
		c.sendError("Room not found")
		return
	}

	userID := msg.UserID
	if userID == "" {
		userID = generateID(8)
	}

	// Get or create user
	c.hub.mu.Lock()
	user, exists := room.Users[userID]
	if !exists {
		user = &User{
			ID:       userID,
			Name:     msg.UserName,
			Role:     RoleEditor,
			Color:    GetUserColor(len(room.Users)),
			JoinedAt: time.Now(),
		}
		room.Users[userID] = user
	}
	c.hub.mu.Unlock()

	c.user = user
	c.roomID = msg.RoomID

	// Send sync data
	c.hub.mu.RLock()
	users := make([]*User, 0, len(room.Users))
	for _, u := range room.Users {
		users = append(users, u)
	}
	syncData := &SyncData{
		Users:    users,
		Messages: room.Messages,
		RoomID:   room.ID,
	}
	c.hub.mu.RUnlock()

	syncMsg := &Message{
		Type:      MsgTypeSync,
		RoomID:    room.ID,
		Timestamp: time.Now(),
		Data:      syncData,
	}

	data, _ := json.Marshal(syncMsg)
	c.send <- data

	// Notify others
	c.hub.broadcastToRoom(room.ID, &Message{
		Type:      MsgTypeJoin,
		UserID:    user.ID,
		UserName:  user.Name,
		RoomID:    room.ID,
		Timestamp: time.Now(),
	})

	// Send updated user list
	c.hub.sendUserList(room.ID)
}

// handleContent processes content messages (chat, command, output)
func (c *Client) handleContent(msg *Message) {
	if c.roomID == "" || c.user == nil {
		return
	}

	msg.RoomID = c.roomID
	msg.UserID = c.user.ID
	msg.UserName = c.user.Name

	c.hub.AddMessage(c.roomID, msg)
}

// handleTyping processes typing indicator messages
func (c *Client) handleTyping(msg *Message) {
	if c.roomID == "" || c.user == nil {
		return
	}

	c.hub.mu.Lock()
	c.user.IsTyping = msg.Content == "true"
	c.hub.mu.Unlock()

	c.hub.broadcastToRoom(c.roomID, &Message{
		Type:      MsgTypeTyping,
		UserID:    c.user.ID,
		UserName:  c.user.Name,
		Content:   msg.Content,
		RoomID:    c.roomID,
		Timestamp: time.Now(),
	})
}

// sendError sends an error message to the client
func (c *Client) sendError(errMsg string) {
	msg := &Message{
		Type:      MsgTypeError,
		Content:   errMsg,
		Timestamp: time.Now(),
	}
	data, _ := json.Marshal(msg)
	c.send <- data
}

// generateID generates a random ID
func generateID(length int) string {
	bytes := make([]byte, length/2)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// generateShareCode generates a human-friendly share code
func generateShareCode() string {
	// Format: XXXX-XXXX (easy to share verbally)
	bytes := make([]byte, 4)
	rand.Read(bytes)
	code := hex.EncodeToString(bytes)
	return fmt.Sprintf("%s-%s", code[:4], code[4:])
}
