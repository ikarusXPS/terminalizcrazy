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

	-- Agent plans table
	CREATE TABLE IF NOT EXISTS agent_plans (
		id TEXT PRIMARY KEY,
		session_id TEXT NOT NULL,
		goal TEXT NOT NULL,
		status TEXT DEFAULT 'pending',
		current_task INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
	);

	-- Agent tasks table
	CREATE TABLE IF NOT EXISTS agent_tasks (
		id TEXT PRIMARY KEY,
		plan_id TEXT NOT NULL,
		sequence INTEGER NOT NULL,
		description TEXT NOT NULL,
		command TEXT NOT NULL,
		status TEXT DEFAULT 'pending',
		output TEXT,
		error TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		completed_at DATETIME,
		FOREIGN KEY (plan_id) REFERENCES agent_plans(id) ON DELETE CASCADE
	);

	-- Workflows table (for Phase 3)
	CREATE TABLE IF NOT EXISTS workflows (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		description TEXT,
		steps TEXT NOT NULL,
		variables TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Workspaces table (for Phase 2)
	CREATE TABLE IF NOT EXISTS workspaces (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		layout TEXT NOT NULL,
		panes_json TEXT,
		active_pane TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_agent_plans_session ON agent_plans(session_id);
	CREATE INDEX IF NOT EXISTS idx_agent_tasks_plan ON agent_tasks(plan_id);
	CREATE INDEX IF NOT EXISTS idx_workflows_name ON workflows(name);
	CREATE INDEX IF NOT EXISTS idx_workspaces_name ON workspaces(name);
	CREATE INDEX IF NOT EXISTS idx_workspaces_created ON workspaces(created_at);
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
	var errs []error

	tables := []struct {
		query string
		key   string
	}{
		{"SELECT COUNT(*) FROM sessions", "sessions"},
		{"SELECT COUNT(*) FROM messages", "messages"},
		{"SELECT COUNT(*) FROM command_history", "commands"},
		{"SELECT COUNT(*) FROM agent_plans", "plans"},
	}

	for _, t := range tables {
		var count int64
		if err := s.db.QueryRow(t.query).Scan(&count); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", t.key, err))
			continue
		}
		stats[t.key] = count
	}

	if len(errs) > 0 {
		return stats, fmt.Errorf("stats errors: %v", errs)
	}
	return stats, nil
}

// --- Agent Plan Methods ---

// AgentPlan represents a stored agent plan
type AgentPlan struct {
	ID          string
	SessionID   string
	Goal        string
	Status      string
	CurrentTask int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// AgentTask represents a stored agent task
type AgentTask struct {
	ID          string
	PlanID      string
	Sequence    int
	Description string
	Command     string
	Status      string
	Output      string
	Error       string
	CreatedAt   time.Time
	CompletedAt *time.Time
}

// SaveAgentPlan saves an agent plan
func (s *Storage) SaveAgentPlan(id, sessionID, goal, status string, currentTask int) (*AgentPlan, error) {
	now := time.Now()

	_, err := s.db.Exec(
		`INSERT INTO agent_plans (id, session_id, goal, status, current_task, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(id) DO UPDATE SET status = ?, current_task = ?, updated_at = ?`,
		id, sessionID, goal, status, currentTask, now, now,
		status, currentTask, now,
	)
	if err != nil {
		return nil, err
	}

	return &AgentPlan{
		ID:          id,
		SessionID:   sessionID,
		Goal:        goal,
		Status:      status,
		CurrentTask: currentTask,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// GetAgentPlan retrieves an agent plan by ID
func (s *Storage) GetAgentPlan(id string) (*AgentPlan, error) {
	var plan AgentPlan
	err := s.db.QueryRow(
		"SELECT id, session_id, goal, status, current_task, created_at, updated_at FROM agent_plans WHERE id = ?",
		id,
	).Scan(&plan.ID, &plan.SessionID, &plan.Goal, &plan.Status, &plan.CurrentTask, &plan.CreatedAt, &plan.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &plan, nil
}

// ListAgentPlans returns agent plans for a session
func (s *Storage) ListAgentPlans(sessionID string, limit int) ([]AgentPlan, error) {
	rows, err := s.db.Query(
		`SELECT id, session_id, goal, status, current_task, created_at, updated_at
		 FROM agent_plans WHERE session_id = ?
		 ORDER BY created_at DESC LIMIT ?`,
		sessionID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []AgentPlan
	for rows.Next() {
		var plan AgentPlan
		if err := rows.Scan(&plan.ID, &plan.SessionID, &plan.Goal, &plan.Status, &plan.CurrentTask, &plan.CreatedAt, &plan.UpdatedAt); err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}

	return plans, rows.Err()
}

// UpdateAgentPlanStatus updates an agent plan's status
func (s *Storage) UpdateAgentPlanStatus(id, status string, currentTask int) error {
	_, err := s.db.Exec(
		"UPDATE agent_plans SET status = ?, current_task = ?, updated_at = ? WHERE id = ?",
		status, currentTask, time.Now(), id,
	)
	return err
}

// DeleteAgentPlan deletes an agent plan and its tasks
func (s *Storage) DeleteAgentPlan(id string) error {
	_, err := s.db.Exec("DELETE FROM agent_plans WHERE id = ?", id)
	return err
}

// --- Agent Task Methods ---

// SaveAgentTask saves an agent task
func (s *Storage) SaveAgentTask(id, planID string, sequence int, description, command, status string) (*AgentTask, error) {
	now := time.Now()

	_, err := s.db.Exec(
		`INSERT INTO agent_tasks (id, plan_id, sequence, description, command, status, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(id) DO UPDATE SET status = ?, command = ?`,
		id, planID, sequence, description, command, status, now,
		status, command,
	)
	if err != nil {
		return nil, err
	}

	return &AgentTask{
		ID:          id,
		PlanID:      planID,
		Sequence:    sequence,
		Description: description,
		Command:     command,
		Status:      status,
		CreatedAt:   now,
	}, nil
}

// GetAgentTasks retrieves tasks for a plan
func (s *Storage) GetAgentTasks(planID string) ([]AgentTask, error) {
	rows, err := s.db.Query(
		`SELECT id, plan_id, sequence, description, command, status, output, error, created_at, completed_at
		 FROM agent_tasks WHERE plan_id = ?
		 ORDER BY sequence ASC`,
		planID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []AgentTask
	for rows.Next() {
		var task AgentTask
		var output, errMsg sql.NullString
		var completedAt sql.NullTime

		if err := rows.Scan(&task.ID, &task.PlanID, &task.Sequence, &task.Description, &task.Command,
			&task.Status, &output, &errMsg, &task.CreatedAt, &completedAt); err != nil {
			return nil, err
		}

		if output.Valid {
			task.Output = output.String
		}
		if errMsg.Valid {
			task.Error = errMsg.String
		}
		if completedAt.Valid {
			task.CompletedAt = &completedAt.Time
		}

		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

// UpdateAgentTask updates an agent task
func (s *Storage) UpdateAgentTask(id, status, output, errMsg string) error {
	var completedAt interface{}
	if status == "completed" || status == "failed" || status == "skipped" {
		completedAt = time.Now()
	}

	_, err := s.db.Exec(
		"UPDATE agent_tasks SET status = ?, output = ?, error = ?, completed_at = ? WHERE id = ?",
		status, output, errMsg, completedAt, id,
	)
	return err
}

// --- Workflow Methods ---

// Workflow represents a stored workflow
type Workflow struct {
	ID          string
	Name        string
	Description string
	Steps       string // JSON encoded
	Variables   string // JSON encoded
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// SaveWorkflow saves a workflow
func (s *Storage) SaveWorkflow(name, description, steps, variables string) (*Workflow, error) {
	id := fmt.Sprintf("wf-%d", time.Now().UnixNano()%1000000)
	now := time.Now()

	_, err := s.db.Exec(
		`INSERT INTO workflows (id, name, description, steps, variables, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(name) DO UPDATE SET description = ?, steps = ?, variables = ?, updated_at = ?`,
		id, name, description, steps, variables, now, now,
		description, steps, variables, now,
	)
	if err != nil {
		return nil, err
	}

	return &Workflow{
		ID:          id,
		Name:        name,
		Description: description,
		Steps:       steps,
		Variables:   variables,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// GetWorkflow retrieves a workflow by name
func (s *Storage) GetWorkflow(name string) (*Workflow, error) {
	var wf Workflow
	var desc, vars sql.NullString

	err := s.db.QueryRow(
		"SELECT id, name, description, steps, variables, created_at, updated_at FROM workflows WHERE name = ?",
		name,
	).Scan(&wf.ID, &wf.Name, &desc, &wf.Steps, &vars, &wf.CreatedAt, &wf.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if desc.Valid {
		wf.Description = desc.String
	}
	if vars.Valid {
		wf.Variables = vars.String
	}

	return &wf, nil
}

// ListWorkflows returns all workflows
func (s *Storage) ListWorkflows() ([]Workflow, error) {
	rows, err := s.db.Query(
		"SELECT id, name, description, steps, variables, created_at, updated_at FROM workflows ORDER BY name ASC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []Workflow
	for rows.Next() {
		var wf Workflow
		var desc, vars sql.NullString

		if err := rows.Scan(&wf.ID, &wf.Name, &desc, &wf.Steps, &vars, &wf.CreatedAt, &wf.UpdatedAt); err != nil {
			return nil, err
		}

		if desc.Valid {
			wf.Description = desc.String
		}
		if vars.Valid {
			wf.Variables = vars.String
		}

		workflows = append(workflows, wf)
	}

	return workflows, rows.Err()
}

// DeleteWorkflow deletes a workflow
func (s *Storage) DeleteWorkflow(name string) error {
	_, err := s.db.Exec("DELETE FROM workflows WHERE name = ?", name)
	return err
}

// --- Data Retention Methods (GDPR Compliance) ---

// RetentionStats holds statistics about data cleanup
type RetentionStats struct {
	MessagesDeleted       int64
	CommandsDeleted       int64
	AgentPlansDeleted     int64
	SessionsDeleted       int64
}

// CleanupOldMessages deletes messages older than the specified number of days
func (s *Storage) CleanupOldMessages(retentionDays int) (int64, error) {
	if retentionDays <= 0 {
		return 0, nil // 0 means keep forever
	}

	result, err := s.db.Exec(`
		DELETE FROM messages
		WHERE created_at < datetime('now', '-' || ? || ' days')
	`, retentionDays)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup messages: %w", err)
	}

	return result.RowsAffected()
}

// CleanupOldCommandHistory deletes command history older than the specified number of days
func (s *Storage) CleanupOldCommandHistory(retentionDays int) (int64, error) {
	if retentionDays <= 0 {
		return 0, nil // 0 means keep forever
	}

	result, err := s.db.Exec(`
		DELETE FROM command_history
		WHERE executed_at < datetime('now', '-' || ? || ' days')
	`, retentionDays)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup command history: %w", err)
	}

	return result.RowsAffected()
}

// CleanupOldAgentPlans deletes agent plans older than the specified number of days
func (s *Storage) CleanupOldAgentPlans(retentionDays int) (int64, error) {
	if retentionDays <= 0 {
		return 0, nil // 0 means keep forever
	}

	// Delete tasks first (foreign key constraint)
	_, err := s.db.Exec(`
		DELETE FROM agent_tasks
		WHERE plan_id IN (
			SELECT id FROM agent_plans
			WHERE created_at < datetime('now', '-' || ? || ' days')
		)
	`, retentionDays)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup agent tasks: %w", err)
	}

	// Delete plans
	result, err := s.db.Exec(`
		DELETE FROM agent_plans
		WHERE created_at < datetime('now', '-' || ? || ' days')
	`, retentionDays)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup agent plans: %w", err)
	}

	return result.RowsAffected()
}

// CleanupEmptySessions deletes sessions that have no messages
func (s *Storage) CleanupEmptySessions() (int64, error) {
	result, err := s.db.Exec(`
		DELETE FROM sessions
		WHERE id NOT IN (SELECT DISTINCT session_id FROM messages)
	`)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup empty sessions: %w", err)
	}

	return result.RowsAffected()
}

// RunRetentionCleanup performs a full retention cleanup based on the provided settings
func (s *Storage) RunRetentionCleanup(messageRetentionDays, commandRetentionDays, planRetentionDays int) (*RetentionStats, error) {
	stats := &RetentionStats{}

	// Cleanup messages
	messagesDeleted, err := s.CleanupOldMessages(messageRetentionDays)
	if err != nil {
		return stats, fmt.Errorf("message cleanup failed: %w", err)
	}
	stats.MessagesDeleted = messagesDeleted

	// Cleanup command history
	commandsDeleted, err := s.CleanupOldCommandHistory(commandRetentionDays)
	if err != nil {
		return stats, fmt.Errorf("command history cleanup failed: %w", err)
	}
	stats.CommandsDeleted = commandsDeleted

	// Cleanup agent plans
	plansDeleted, err := s.CleanupOldAgentPlans(planRetentionDays)
	if err != nil {
		return stats, fmt.Errorf("agent plan cleanup failed: %w", err)
	}
	stats.AgentPlansDeleted = plansDeleted

	// Cleanup empty sessions
	sessionsDeleted, err := s.CleanupEmptySessions()
	if err != nil {
		return stats, fmt.Errorf("empty session cleanup failed: %w", err)
	}
	stats.SessionsDeleted = sessionsDeleted

	return stats, nil
}

// ExportUserData exports all user data for GDPR data portability (Art. 20)
func (s *Storage) ExportUserData() (map[string]interface{}, error) {
	export := make(map[string]interface{})

	// Export sessions (no limit for export)
	sessions, err := s.ListSessions(10000)
	if err != nil {
		return nil, fmt.Errorf("failed to export sessions: %w", err)
	}
	export["sessions"] = sessions

	// Export all messages
	var messages []map[string]interface{}
	rows, err := s.db.Query(`
		SELECT id, session_id, role, content, command, success, created_at
		FROM messages ORDER BY created_at
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to export messages: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, sessionID, role, content string
		var command sql.NullString
		var success sql.NullBool
		var createdAt time.Time

		if err := rows.Scan(&id, &sessionID, &role, &content, &command, &success, &createdAt); err != nil {
			return nil, err
		}

		msg := map[string]interface{}{
			"id":         id,
			"session_id": sessionID,
			"role":       role,
			"content":    content,
			"created_at": createdAt.Format(time.RFC3339),
		}
		if command.Valid {
			msg["command"] = command.String
		}
		if success.Valid {
			msg["success"] = success.Bool
		}
		messages = append(messages, msg)
	}
	export["messages"] = messages

	// Export command history
	var commands []map[string]interface{}
	cmdRows, err := s.db.Query(`
		SELECT id, session_id, command, output, exit_code, duration_ms, executed_at
		FROM command_history ORDER BY executed_at
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to export command history: %w", err)
	}
	defer cmdRows.Close()

	for cmdRows.Next() {
		var id, sessionID, command string
		var output sql.NullString
		var exitCode, durationMs int
		var executedAt time.Time

		if err := cmdRows.Scan(&id, &sessionID, &command, &output, &exitCode, &durationMs, &executedAt); err != nil {
			return nil, err
		}

		cmd := map[string]interface{}{
			"id":          id,
			"session_id":  sessionID,
			"command":     command,
			"exit_code":   exitCode,
			"duration_ms": durationMs,
			"executed_at": executedAt.Format(time.RFC3339),
		}
		if output.Valid {
			cmd["output"] = output.String
		}
		commands = append(commands, cmd)
	}
	export["command_history"] = commands

	// Export agent plans (all sessions)
	var plans []map[string]interface{}
	planRows, err := s.db.Query(`
		SELECT id, session_id, goal, status, current_task, created_at, updated_at
		FROM agent_plans ORDER BY created_at
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to export agent plans: %w", err)
	}
	defer planRows.Close()

	for planRows.Next() {
		var id, sessionID, goal, status string
		var currentTask int
		var createdAt, updatedAt time.Time

		if err := planRows.Scan(&id, &sessionID, &goal, &status, &currentTask, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		plan := map[string]interface{}{
			"id":           id,
			"session_id":   sessionID,
			"goal":         goal,
			"status":       status,
			"current_task": currentTask,
			"created_at":   createdAt.Format(time.RFC3339),
			"updated_at":   updatedAt.Format(time.RFC3339),
		}
		plans = append(plans, plan)
	}
	export["agent_plans"] = plans

	// Add export metadata
	export["_metadata"] = map[string]interface{}{
		"exported_at": time.Now().Format(time.RFC3339),
		"format":      "terminalizcrazy-gdpr-export-v1",
	}

	return export, nil
}

// DeleteAllUserData deletes all user data (GDPR right to erasure - Art. 17)
func (s *Storage) DeleteAllUserData() error {
	// Order matters due to foreign key constraints
	tables := []string{
		"agent_tasks",
		"agent_plans",
		"messages",
		"command_history",
		"sessions",
		"workflows",
		"workspaces",
	}

	for _, table := range tables {
		if _, err := s.db.Exec("DELETE FROM " + table); err != nil {
			return fmt.Errorf("failed to delete from %s: %w", table, err)
		}
	}

	// Vacuum to reclaim space
	if _, err := s.db.Exec("VACUUM"); err != nil {
		return fmt.Errorf("failed to vacuum database: %w", err)
	}

	return nil
}
