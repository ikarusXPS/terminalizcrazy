package workspace

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// MaxWorkspaces is the maximum number of workspaces allowed
const MaxWorkspaces = 10

// Manager handles workspace creation, switching, and management
type Manager struct {
	workspaces     map[string]*Workspace
	activeID       string
	storage        Storage
	onChange       func(*Workspace)
	onSwitch       func(old, new *Workspace)
	mu             sync.RWMutex
	width          int
	height         int
}

// Storage is the interface for workspace persistence
type Storage interface {
	SaveWorkspace(w *Workspace) error
	GetWorkspace(id string) (*Workspace, error)
	ListWorkspaces() ([]*Workspace, error)
	DeleteWorkspace(id string) error
}

// ManagerConfig holds configuration for the workspace manager
type ManagerConfig struct {
	Storage Storage
	Width   int
	Height  int
}

// NewManager creates a new workspace manager
func NewManager(config ManagerConfig) (*Manager, error) {
	m := &Manager{
		workspaces: make(map[string]*Workspace),
		storage:    config.Storage,
		width:      config.Width,
		height:     config.Height,
	}

	// Load workspaces from storage if available
	if config.Storage != nil {
		workspaces, err := config.Storage.ListWorkspaces()
		if err == nil && len(workspaces) > 0 {
			for _, w := range workspaces {
				m.workspaces[w.ID] = w
			}
			// Set the first workspace as active
			m.activeID = workspaces[0].ID
		}
	}

	// Create default workspace if none exist
	if len(m.workspaces) == 0 {
		defaultWs := m.createDefaultWorkspace()
		m.workspaces[defaultWs.ID] = defaultWs
		m.activeID = defaultWs.ID
	}

	return m, nil
}

// createDefaultWorkspace creates the default workspace with quad layout
func (m *Manager) createDefaultWorkspace() *Workspace {
	id := generateID()
	w := NewWorkspace(id, "Default", LayoutQuad)
	w.Panes = DefaultPanesForLayout(LayoutQuad)

	// Apply layout to set positions
	if m.width > 0 && m.height > 0 {
		ApplyLayoutToWorkspace(w, m.width, m.height)
	}

	// Set first pane as active
	if len(w.Panes) > 0 {
		w.ActivePane = w.Panes[0].ID
	}

	return w
}

// CreateWorkspace creates a new workspace with the specified layout
func (m *Manager) CreateWorkspace(name string, layout LayoutType) (*Workspace, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.workspaces) >= MaxWorkspaces {
		return nil, ErrMaxWorkspacesReached
	}

	id := generateID()
	w := NewWorkspace(id, name, layout)
	w.Panes = DefaultPanesForLayout(layout)

	// Apply layout
	if m.width > 0 && m.height > 0 {
		ApplyLayoutToWorkspace(w, m.width, m.height)
	}

	// Set first pane as active
	if len(w.Panes) > 0 {
		w.ActivePane = w.Panes[0].ID
	}

	m.workspaces[id] = w

	// Persist if storage available
	if m.storage != nil {
		if err := m.storage.SaveWorkspace(w); err != nil {
			delete(m.workspaces, id)
			return nil, fmt.Errorf("failed to save workspace: %w", err)
		}
	}

	return w, nil
}

// GetWorkspace returns a workspace by ID
func (m *Manager) GetWorkspace(id string) *Workspace {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.workspaces[id]
}

// GetActiveWorkspace returns the currently active workspace
func (m *Manager) GetActiveWorkspace() *Workspace {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.workspaces[m.activeID]
}

// SwitchWorkspace switches to a workspace by ID
func (m *Manager) SwitchWorkspace(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.workspaces[id]; !ok {
		return ErrWorkspaceNotFound
	}

	oldWs := m.workspaces[m.activeID]
	m.activeID = id
	newWs := m.workspaces[id]

	// Notify listeners
	if m.onSwitch != nil {
		m.onSwitch(oldWs, newWs)
	}

	return nil
}

// SwitchWorkspaceByIndex switches to a workspace by index (1-based)
func (m *Manager) SwitchWorkspaceByIndex(index int) error {
	m.mu.RLock()
	workspaces := m.ListWorkspaces()
	m.mu.RUnlock()

	if index < 1 || index > len(workspaces) {
		return ErrWorkspaceNotFound
	}

	return m.SwitchWorkspace(workspaces[index-1].ID)
}

// DeleteWorkspace deletes a workspace by ID
func (m *Manager) DeleteWorkspace(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.workspaces) <= 1 {
		return ErrCannotDeleteLastWorkspace
	}

	if _, ok := m.workspaces[id]; !ok {
		return ErrWorkspaceNotFound
	}

	// If deleting active workspace, switch to another
	if m.activeID == id {
		for wsID := range m.workspaces {
			if wsID != id {
				m.activeID = wsID
				break
			}
		}
	}

	delete(m.workspaces, id)

	// Delete from storage if available
	if m.storage != nil {
		if err := m.storage.DeleteWorkspace(id); err != nil {
			return fmt.Errorf("failed to delete workspace from storage: %w", err)
		}
	}

	return nil
}

// ListWorkspaces returns all workspaces sorted by creation time
func (m *Manager) ListWorkspaces() []*Workspace {
	m.mu.RLock()
	defer m.mu.RUnlock()

	workspaces := make([]*Workspace, 0, len(m.workspaces))
	for _, w := range m.workspaces {
		workspaces = append(workspaces, w)
	}

	// Sort by creation time
	for i := 0; i < len(workspaces)-1; i++ {
		for j := i + 1; j < len(workspaces); j++ {
			if workspaces[i].CreatedAt.After(workspaces[j].CreatedAt) {
				workspaces[i], workspaces[j] = workspaces[j], workspaces[i]
			}
		}
	}

	return workspaces
}

// GetWorkspaceCount returns the number of workspaces
func (m *Manager) GetWorkspaceCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.workspaces)
}

// RenameWorkspace renames a workspace
func (m *Manager) RenameWorkspace(id, name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	w, ok := m.workspaces[id]
	if !ok {
		return ErrWorkspaceNotFound
	}

	w.Name = name
	w.UpdatedAt = time.Now()

	// Persist if storage available
	if m.storage != nil {
		if err := m.storage.SaveWorkspace(w); err != nil {
			return fmt.Errorf("failed to save workspace: %w", err)
		}
	}

	return nil
}

// SetLayout changes the layout of a workspace
func (m *Manager) SetLayout(id string, layout LayoutType) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	w, ok := m.workspaces[id]
	if !ok {
		return ErrWorkspaceNotFound
	}

	w.Layout = layout
	w.UpdatedAt = time.Now()

	// Recalculate positions
	if m.width > 0 && m.height > 0 {
		ApplyLayoutToWorkspace(w, m.width, m.height)
	}

	// Persist if storage available
	if m.storage != nil {
		if err := m.storage.SaveWorkspace(w); err != nil {
			return fmt.Errorf("failed to save workspace: %w", err)
		}
	}

	// Notify listeners
	if m.onChange != nil {
		m.onChange(w)
	}

	return nil
}

// SetSize updates the manager dimensions and recalculates all layouts
func (m *Manager) SetSize(width, height int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.width = width
	m.height = height

	// Recalculate all workspace layouts
	for _, w := range m.workspaces {
		ApplyLayoutToWorkspace(w, width, height)
	}
}

// SaveAll persists all workspaces
func (m *Manager) SaveAll() error {
	if m.storage == nil {
		return nil
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, w := range m.workspaces {
		if err := m.storage.SaveWorkspace(w); err != nil {
			return fmt.Errorf("failed to save workspace %s: %w", w.ID, err)
		}
	}

	return nil
}

// OnChange sets a callback for workspace changes
func (m *Manager) OnChange(fn func(*Workspace)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onChange = fn
}

// OnSwitch sets a callback for workspace switching
func (m *Manager) OnSwitch(fn func(old, new *Workspace)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onSwitch = fn
}

// DuplicateWorkspace creates a copy of an existing workspace
func (m *Manager) DuplicateWorkspace(id string) (*Workspace, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.workspaces) >= MaxWorkspaces {
		return nil, ErrMaxWorkspacesReached
	}

	original, ok := m.workspaces[id]
	if !ok {
		return nil, ErrWorkspaceNotFound
	}

	// Clone the workspace with a new ID
	newID := generateID()
	w := original.Clone()
	w.ID = newID
	w.Name = fmt.Sprintf("%s (copy)", original.Name)
	w.CreatedAt = time.Now()
	w.UpdatedAt = time.Now()

	// Generate new IDs for panes
	for i := range w.Panes {
		w.Panes[i].ID = fmt.Sprintf("pane-%d", i+1)
	}
	if len(w.Panes) > 0 {
		w.ActivePane = w.Panes[0].ID
	}

	m.workspaces[newID] = w

	// Persist if storage available
	if m.storage != nil {
		if err := m.storage.SaveWorkspace(w); err != nil {
			delete(m.workspaces, newID)
			return nil, fmt.Errorf("failed to save workspace: %w", err)
		}
	}

	return w, nil
}

// GetActiveIndex returns the 1-based index of the active workspace
func (m *Manager) GetActiveIndex() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	workspaces := m.ListWorkspaces()
	for i, w := range workspaces {
		if w.ID == m.activeID {
			return i + 1
		}
	}
	return 0
}

// AddPaneToWorkspace adds a pane to the specified workspace
func (m *Manager) AddPaneToWorkspace(wsID string, pane PaneState) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	w, ok := m.workspaces[wsID]
	if !ok {
		return ErrWorkspaceNotFound
	}

	w.AddPane(pane)

	// Recalculate layout
	if m.width > 0 && m.height > 0 {
		ApplyLayoutToWorkspace(w, m.width, m.height)
	}

	// Persist if storage available
	if m.storage != nil {
		if err := m.storage.SaveWorkspace(w); err != nil {
			return fmt.Errorf("failed to save workspace: %w", err)
		}
	}

	if m.onChange != nil {
		m.onChange(w)
	}

	return nil
}

// RemovePaneFromWorkspace removes a pane from the specified workspace
func (m *Manager) RemovePaneFromWorkspace(wsID, paneID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	w, ok := m.workspaces[wsID]
	if !ok {
		return ErrWorkspaceNotFound
	}

	if !w.RemovePane(paneID) {
		return ErrPaneNotFound
	}

	// Recalculate layout
	if m.width > 0 && m.height > 0 {
		ApplyLayoutToWorkspace(w, m.width, m.height)
	}

	// Update active pane if needed
	if w.ActivePane == paneID && len(w.Panes) > 0 {
		w.ActivePane = w.Panes[0].ID
	}

	// Persist if storage available
	if m.storage != nil {
		if err := m.storage.SaveWorkspace(w); err != nil {
			return fmt.Errorf("failed to save workspace: %w", err)
		}
	}

	if m.onChange != nil {
		m.onChange(w)
	}

	return nil
}

// generateID generates a unique workspace ID
func generateID() string {
	return fmt.Sprintf("ws-%s", uuid.New().String()[:8])
}
