package workspace

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// SQLiteStorage implements Storage interface for SQLite
type SQLiteStorage struct {
	db *sql.DB
}

// NewSQLiteStorage creates a new SQLite storage for workspaces
func NewSQLiteStorage(db *sql.DB) (*SQLiteStorage, error) {
	s := &SQLiteStorage{db: db}

	// Create workspaces table if it doesn't exist
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate workspaces table: %w", err)
	}

	return s, nil
}

// migrate creates the workspaces table
func (s *SQLiteStorage) migrate() error {
	schema := `
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

	CREATE INDEX IF NOT EXISTS idx_workspaces_name ON workspaces(name);
	CREATE INDEX IF NOT EXISTS idx_workspaces_created ON workspaces(created_at);
	`

	_, err := s.db.Exec(schema)
	return err
}

// SaveWorkspace saves a workspace to the database
func (s *SQLiteStorage) SaveWorkspace(w *Workspace) error {
	panesJSON, err := json.Marshal(w.Panes)
	if err != nil {
		return fmt.Errorf("failed to marshal panes: %w", err)
	}

	_, err = s.db.Exec(`
		INSERT INTO workspaces (id, name, description, layout, panes_json, active_pane, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			description = excluded.description,
			layout = excluded.layout,
			panes_json = excluded.panes_json,
			active_pane = excluded.active_pane,
			updated_at = excluded.updated_at
	`,
		w.ID, w.Name, w.Description, string(w.Layout), string(panesJSON),
		w.ActivePane, w.CreatedAt, w.UpdatedAt,
	)

	return err
}

// GetWorkspace retrieves a workspace by ID
func (s *SQLiteStorage) GetWorkspace(id string) (*Workspace, error) {
	var w Workspace
	var panesJSON, description, activePaneSQL sql.NullString
	var layout string

	err := s.db.QueryRow(`
		SELECT id, name, description, layout, panes_json, active_pane, created_at, updated_at
		FROM workspaces WHERE id = ?
	`, id).Scan(
		&w.ID, &w.Name, &description, &layout, &panesJSON, &activePaneSQL, &w.CreatedAt, &w.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrWorkspaceNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace: %w", err)
	}

	w.Layout = LayoutType(layout)

	if description.Valid {
		w.Description = description.String
	}
	if activePaneSQL.Valid {
		w.ActivePane = activePaneSQL.String
	}

	if panesJSON.Valid {
		var panes []PaneState
		if err := json.Unmarshal([]byte(panesJSON.String), &panes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal panes: %w", err)
		}
		w.Panes = panes
	} else {
		w.Panes = []PaneState{}
	}

	return &w, nil
}

// ListWorkspaces returns all workspaces
func (s *SQLiteStorage) ListWorkspaces() ([]*Workspace, error) {
	rows, err := s.db.Query(`
		SELECT id, name, description, layout, panes_json, active_pane, created_at, updated_at
		FROM workspaces ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list workspaces: %w", err)
	}
	defer rows.Close()

	var workspaces []*Workspace
	for rows.Next() {
		var w Workspace
		var panesJSON, description, activePaneSQL sql.NullString
		var layout string

		if err := rows.Scan(
			&w.ID, &w.Name, &description, &layout, &panesJSON, &activePaneSQL, &w.CreatedAt, &w.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan workspace: %w", err)
		}

		w.Layout = LayoutType(layout)

		if description.Valid {
			w.Description = description.String
		}
		if activePaneSQL.Valid {
			w.ActivePane = activePaneSQL.String
		}

		if panesJSON.Valid {
			var panes []PaneState
			if err := json.Unmarshal([]byte(panesJSON.String), &panes); err != nil {
				return nil, fmt.Errorf("failed to unmarshal panes: %w", err)
			}
			w.Panes = panes
		} else {
			w.Panes = []PaneState{}
		}

		workspaces = append(workspaces, &w)
	}

	return workspaces, rows.Err()
}

// DeleteWorkspace deletes a workspace by ID
func (s *SQLiteStorage) DeleteWorkspace(id string) error {
	result, err := s.db.Exec("DELETE FROM workspaces WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete workspace: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrWorkspaceNotFound
	}

	return nil
}

// InMemoryStorage implements Storage interface for in-memory storage (for testing)
type InMemoryStorage struct {
	workspaces map[string]*Workspace
}

// NewInMemoryStorage creates a new in-memory storage
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		workspaces: make(map[string]*Workspace),
	}
}

// SaveWorkspace saves a workspace to memory
func (s *InMemoryStorage) SaveWorkspace(w *Workspace) error {
	s.workspaces[w.ID] = w.Clone()
	return nil
}

// GetWorkspace retrieves a workspace by ID
func (s *InMemoryStorage) GetWorkspace(id string) (*Workspace, error) {
	if w, ok := s.workspaces[id]; ok {
		return w.Clone(), nil
	}
	return nil, ErrWorkspaceNotFound
}

// ListWorkspaces returns all workspaces
func (s *InMemoryStorage) ListWorkspaces() ([]*Workspace, error) {
	workspaces := make([]*Workspace, 0, len(s.workspaces))
	for _, w := range s.workspaces {
		workspaces = append(workspaces, w.Clone())
	}

	// Sort by creation time
	for i := 0; i < len(workspaces)-1; i++ {
		for j := i + 1; j < len(workspaces); j++ {
			if workspaces[i].CreatedAt.After(workspaces[j].CreatedAt) {
				workspaces[i], workspaces[j] = workspaces[j], workspaces[i]
			}
		}
	}

	return workspaces, nil
}

// DeleteWorkspace deletes a workspace by ID
func (s *InMemoryStorage) DeleteWorkspace(id string) error {
	if _, ok := s.workspaces[id]; !ok {
		return ErrWorkspaceNotFound
	}
	delete(s.workspaces, id)
	return nil
}

// Utility function to convert workspace to storage format
func workspaceToStorageFormat(w *Workspace) (map[string]interface{}, error) {
	panesJSON, err := json.Marshal(w.Panes)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":          w.ID,
		"name":        w.Name,
		"description": w.Description,
		"layout":      string(w.Layout),
		"panes_json":  string(panesJSON),
		"active_pane": w.ActivePane,
		"created_at":  w.CreatedAt.Format(time.RFC3339),
		"updated_at":  w.UpdatedAt.Format(time.RFC3339),
	}, nil
}
