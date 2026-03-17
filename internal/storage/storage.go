package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// Storage handles persistent data storage
type Storage struct {
	db      *sql.DB
	dbPath  string
}

// Message represents a chat message in storage
type Message struct {
	ID        int64
	SessionID string
	Role      string
	Content   string
	Command   string
	Success   bool
	CreatedAt time.Time
}

// Session represents a terminal session
type Session struct {
	ID        string
	Name      string
	WorkDir   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CommandHistory represents a command in history
type CommandHistory struct {
	ID        int64
	Command   string
	Output    string
	Success   bool
	Duration  int64 // milliseconds
	CreatedAt time.Time
}

// New creates a new Storage instance
func New(dataDir string) (*Storage, error) {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	dbPath := filepath.Join(dataDir, "terminalizcrazy.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	s := &Storage{
		db:     db,
		dbPath: dbPath,
	}

	if err := s.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return s, nil
}

// migrate creates the database schema
func (s *Storage) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		work_dir TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id TEXT NOT NULL,
		role TEXT NOT NULL,
		content TEXT NOT NULL,
		command TEXT,
		success INTEGER DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS command_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		command TEXT NOT NULL,
		output TEXT,
		success INTEGER DEFAULT 1,
		duration_ms INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_messages_session ON messages(session_id);
	CREATE INDEX IF NOT EXISTS idx_messages_created ON messages(created_at);
	CREATE INDEX IF NOT EXISTS idx_command_history_created ON command_history(created_at);
	CREATE INDEX IF NOT EXISTS idx_command_history_command ON command_history(command);
	`

	_, err := s.db.Exec(schema)
	return err
}

// Close closes the database connection
func (s *Storage) Close() error {
	return s.db.Close()
}

// --- Session Methods ---

// CreateSession creates a new session
func (s *Storage) CreateSession(id, name, workDir string) (*Session, error) {
	now := time.Now()

	_, err := s.db.Exec(
		"INSERT INTO sessions (id, name, work_dir, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		id, name, workDir, now, now,
	)
	if err != nil {
		return nil, err
	}

	return &Session{
		ID:        id,
		Name:      name,
		WorkDir:   workDir,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// GetSession retrieves a session by ID
func (s *Storage) GetSession(id string) (*Session, error) {
	var session Session
	err := s.db.QueryRow(
		"SELECT id, name, work_dir, created_at, updated_at FROM sessions WHERE id = ?",
		id,
	).Scan(&session.ID, &session.Name, &session.WorkDir, &session.CreatedAt, &session.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// ListSessions returns recent sessions
func (s *Storage) ListSessions(limit int) ([]Session, error) {
	rows, err := s.db.Query(
		"SELECT id, name, work_dir, created_at, updated_at FROM sessions ORDER BY updated_at DESC LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var session Session
		if err := rows.Scan(&session.ID, &session.Name, &session.WorkDir, &session.CreatedAt, &session.UpdatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, rows.Err()
}

// UpdateSessionTimestamp updates the session's updated_at timestamp
func (s *Storage) UpdateSessionTimestamp(id string) error {
	_, err := s.db.Exec("UPDATE sessions SET updated_at = ? WHERE id = ?", time.Now(), id)
	return err
}

// DeleteSession deletes a session and its messages
func (s *Storage) DeleteSession(id string) error {
	_, err := s.db.Exec("DELETE FROM sessions WHERE id = ?", id)
	return err
}

// --- Message Methods ---

// SaveMessage saves a message to a session
func (s *Storage) SaveMessage(sessionID, role, content, command string, success bool) (*Message, error) {
	now := time.Now()

	result, err := s.db.Exec(
		"INSERT INTO messages (session_id, role, content, command, success, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		sessionID, role, content, command, success, now,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Update session timestamp
	s.UpdateSessionTimestamp(sessionID)

	return &Message{
		ID:        id,
		SessionID: sessionID,
		Role:      role,
		Content:   content,
		Command:   command,
		Success:   success,
		CreatedAt: now,
	}, nil
}

// GetSessionMessages retrieves messages for a session
func (s *Storage) GetSessionMessages(sessionID string, limit int) ([]Message, error) {
	rows, err := s.db.Query(
		`SELECT id, session_id, role, content, command, success, created_at
		 FROM messages WHERE session_id = ?
		 ORDER BY created_at ASC LIMIT ?`,
		sessionID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		var command sql.NullString
		if err := rows.Scan(&msg.ID, &msg.SessionID, &msg.Role, &msg.Content, &command, &msg.Success, &msg.CreatedAt); err != nil {
			return nil, err
		}
		if command.Valid {
			msg.Command = command.String
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

// --- Command History Methods ---

// SaveCommand saves a command to history
func (s *Storage) SaveCommand(command, output string, success bool, durationMs int64) (*CommandHistory, error) {
	now := time.Now()

	result, err := s.db.Exec(
		"INSERT INTO command_history (command, output, success, duration_ms, created_at) VALUES (?, ?, ?, ?, ?)",
		command, output, success, durationMs, now,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &CommandHistory{
		ID:        id,
		Command:   command,
		Output:    output,
		Success:   success,
		Duration:  durationMs,
		CreatedAt: now,
	}, nil
}

// GetCommandHistory retrieves command history
func (s *Storage) GetCommandHistory(limit int) ([]CommandHistory, error) {
	rows, err := s.db.Query(
		`SELECT id, command, output, success, duration_ms, created_at
		 FROM command_history
		 ORDER BY created_at DESC LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []CommandHistory
	for rows.Next() {
		var cmd CommandHistory
		var output sql.NullString
		var duration sql.NullInt64
		if err := rows.Scan(&cmd.ID, &cmd.Command, &output, &cmd.Success, &duration, &cmd.CreatedAt); err != nil {
			return nil, err
		}
		if output.Valid {
			cmd.Output = output.String
		}
		if duration.Valid {
			cmd.Duration = duration.Int64
		}
		history = append(history, cmd)
	}

	return history, rows.Err()
}

// SearchCommands searches command history
func (s *Storage) SearchCommands(query string, limit int) ([]CommandHistory, error) {
	rows, err := s.db.Query(
		`SELECT id, command, output, success, duration_ms, created_at
		 FROM command_history
		 WHERE command LIKE ?
		 ORDER BY created_at DESC LIMIT ?`,
		"%"+query+"%", limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []CommandHistory
	for rows.Next() {
		var cmd CommandHistory
		var output sql.NullString
		var duration sql.NullInt64
		if err := rows.Scan(&cmd.ID, &cmd.Command, &output, &cmd.Success, &duration, &cmd.CreatedAt); err != nil {
			return nil, err
		}
		if output.Valid {
			cmd.Output = output.String
		}
		if duration.Valid {
			cmd.Duration = duration.Int64
		}
		history = append(history, cmd)
	}

	return history, rows.Err()
}

// GetUniqueCommands returns unique commands (for autocomplete)
func (s *Storage) GetUniqueCommands(limit int) ([]string, error) {
	rows, err := s.db.Query(
		`SELECT DISTINCT command FROM command_history
		 ORDER BY created_at DESC LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []string
	for rows.Next() {
		var cmd string
		if err := rows.Scan(&cmd); err != nil {
			return nil, err
		}
		commands = append(commands, cmd)
	}

	return commands, rows.Err()
}

// ClearCommandHistory clears all command history
func (s *Storage) ClearCommandHistory() error {
	_, err := s.db.Exec("DELETE FROM command_history")
	return err
}

// GetStats returns storage statistics
func (s *Storage) GetStats() (map[string]int64, error) {
	stats := make(map[string]int64)

	var count int64

	// Sessions count
	s.db.QueryRow("SELECT COUNT(*) FROM sessions").Scan(&count)
	stats["sessions"] = count

	// Messages count
	s.db.QueryRow("SELECT COUNT(*) FROM messages").Scan(&count)
	stats["messages"] = count

	// Commands count
	s.db.QueryRow("SELECT COUNT(*) FROM command_history").Scan(&count)
	stats["commands"] = count

	return stats, nil
}
