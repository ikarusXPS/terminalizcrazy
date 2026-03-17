package workspace

import (
	"encoding/json"
	"time"
)

// LayoutType defines the type of workspace layout
type LayoutType string

const (
	// LayoutQuad is a 2x2 grid layout (default)
	LayoutQuad LayoutType = "quad"
	// LayoutTall is 1 main pane (60%) + 2 side panes stacked
	LayoutTall LayoutType = "tall"
	// LayoutWide is 1 top pane (60%) + 2 bottom panes
	LayoutWide LayoutType = "wide"
	// LayoutStack is 4 vertical panes stacked
	LayoutStack LayoutType = "stack"
	// LayoutSingle is a single pane
	LayoutSingle LayoutType = "single"
	// LayoutCustom is a user-defined layout
	LayoutCustom LayoutType = "custom"
)

// PaneType defines the type of content in a pane
type PaneType string

const (
	PaneTypeChat     PaneType = "chat"
	PaneTypeTerminal PaneType = "terminal"
	PaneTypePlan     PaneType = "plan"
	PaneTypeOutput   PaneType = "output"
	PaneTypeHistory  PaneType = "history"
	PaneTypeFiles    PaneType = "files"
)

// PaneState represents the saved state of a pane
type PaneState struct {
	ID        string   `json:"id"`
	Type      PaneType `json:"type"`
	Title     string   `json:"title"`
	X         int      `json:"x"`
	Y         int      `json:"y"`
	Width     int      `json:"width"`
	Height    int      `json:"height"`
	Content   string   `json:"content,omitempty"`
	Floating  bool     `json:"floating,omitempty"`
	Minimized bool     `json:"minimized,omitempty"`
	ZIndex    int      `json:"z_index,omitempty"`
}

// Workspace represents a collection of panes with a specific layout
type Workspace struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Layout      LayoutType   `json:"layout"`
	Panes       []PaneState  `json:"panes"`
	ActivePane  string       `json:"active_pane"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// NewWorkspace creates a new workspace with the given layout
func NewWorkspace(id, name string, layout LayoutType) *Workspace {
	now := time.Now()
	return &Workspace{
		ID:        id,
		Name:      name,
		Layout:    layout,
		Panes:     make([]PaneState, 0),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// AddPane adds a pane to the workspace
func (w *Workspace) AddPane(pane PaneState) {
	w.Panes = append(w.Panes, pane)
	w.UpdatedAt = time.Now()
}

// RemovePane removes a pane from the workspace by ID
func (w *Workspace) RemovePane(id string) bool {
	for i, pane := range w.Panes {
		if pane.ID == id {
			w.Panes = append(w.Panes[:i], w.Panes[i+1:]...)
			w.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

// GetPane returns a pane by ID
func (w *Workspace) GetPane(id string) *PaneState {
	for i := range w.Panes {
		if w.Panes[i].ID == id {
			return &w.Panes[i]
		}
	}
	return nil
}

// UpdatePane updates a pane in the workspace
func (w *Workspace) UpdatePane(pane PaneState) bool {
	for i := range w.Panes {
		if w.Panes[i].ID == pane.ID {
			w.Panes[i] = pane
			w.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

// SetActivePane sets the active pane
func (w *Workspace) SetActivePane(id string) {
	w.ActivePane = id
	w.UpdatedAt = time.Now()
}

// Clone creates a deep copy of the workspace
func (w *Workspace) Clone() *Workspace {
	panes := make([]PaneState, len(w.Panes))
	copy(panes, w.Panes)

	return &Workspace{
		ID:          w.ID,
		Name:        w.Name,
		Description: w.Description,
		Layout:      w.Layout,
		Panes:       panes,
		ActivePane:  w.ActivePane,
		CreatedAt:   w.CreatedAt,
		UpdatedAt:   w.UpdatedAt,
	}
}

// ToJSON serializes the workspace to JSON
func (w *Workspace) ToJSON() (string, error) {
	data, err := json.Marshal(w)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON deserializes a workspace from JSON
func FromJSON(data string) (*Workspace, error) {
	var w Workspace
	if err := json.Unmarshal([]byte(data), &w); err != nil {
		return nil, err
	}
	return &w, nil
}

// GetPaneCount returns the number of panes
func (w *Workspace) GetPaneCount() int {
	return len(w.Panes)
}

// GetFloatingPanes returns all floating panes
func (w *Workspace) GetFloatingPanes() []PaneState {
	var floating []PaneState
	for _, pane := range w.Panes {
		if pane.Floating {
			floating = append(floating, pane)
		}
	}
	return floating
}

// GetDockedPanes returns all non-floating panes
func (w *Workspace) GetDockedPanes() []PaneState {
	var docked []PaneState
	for _, pane := range w.Panes {
		if !pane.Floating {
			docked = append(docked, pane)
		}
	}
	return docked
}

// Validate checks if the workspace is valid
func (w *Workspace) Validate() error {
	if w.ID == "" {
		return ErrInvalidWorkspaceID
	}
	if w.Name == "" {
		return ErrInvalidWorkspaceName
	}
	return nil
}

// LayoutDescription returns a human-readable description of the layout
func (l LayoutType) Description() string {
	switch l {
	case LayoutQuad:
		return "2x2 grid layout"
	case LayoutTall:
		return "Main pane with side stack"
	case LayoutWide:
		return "Top pane with bottom row"
	case LayoutStack:
		return "Vertical stack"
	case LayoutSingle:
		return "Single pane"
	case LayoutCustom:
		return "Custom layout"
	default:
		return "Unknown layout"
	}
}

// AvailableLayouts returns all available layout types
func AvailableLayouts() []LayoutType {
	return []LayoutType{
		LayoutQuad,
		LayoutTall,
		LayoutWide,
		LayoutStack,
		LayoutSingle,
	}
}
