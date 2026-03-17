package tui

import (
	"sync"
)

// ZoomState tracks the zoom state of panes
type ZoomState struct {
	zoomed       bool
	zoomedPaneID string
	savedLayout  *LayoutSnapshot
	mu           sync.RWMutex
}

// LayoutSnapshot stores the layout state before zoom
type LayoutSnapshot struct {
	Panes []PaneSnapshot
}

// PaneSnapshot stores a single pane's state
type PaneSnapshot struct {
	ID      string
	X       int
	Y       int
	Width   int
	Height  int
	Visible bool
}

// NewZoomState creates a new zoom state manager
func NewZoomState() *ZoomState {
	return &ZoomState{
		zoomed: false,
	}
}

// IsZoomed returns whether a pane is currently zoomed
func (zs *ZoomState) IsZoomed() bool {
	zs.mu.RLock()
	defer zs.mu.RUnlock()
	return zs.zoomed
}

// GetZoomedPane returns the ID of the zoomed pane, or empty string if none
func (zs *ZoomState) GetZoomedPane() string {
	zs.mu.RLock()
	defer zs.mu.RUnlock()
	return zs.zoomedPaneID
}

// SaveLayout saves the current layout before zooming
func (zs *ZoomState) SaveLayout(panes []*Pane) {
	zs.mu.Lock()
	defer zs.mu.Unlock()

	snapshot := &LayoutSnapshot{
		Panes: make([]PaneSnapshot, len(panes)),
	}

	for i, p := range panes {
		snapshot.Panes[i] = PaneSnapshot{
			ID:      p.ID,
			X:       p.X,
			Y:       p.Y,
			Width:   p.Width,
			Height:  p.Height,
			Visible: true,
		}
	}

	zs.savedLayout = snapshot
}

// GetSavedLayout returns the saved layout snapshot
func (zs *ZoomState) GetSavedLayout() *LayoutSnapshot {
	zs.mu.RLock()
	defer zs.mu.RUnlock()
	return zs.savedLayout
}

// SetZoomed sets the zoom state
func (zs *ZoomState) SetZoomed(zoomed bool, paneID string) {
	zs.mu.Lock()
	defer zs.mu.Unlock()
	zs.zoomed = zoomed
	zs.zoomedPaneID = paneID
}

// ClearZoom clears the zoom state
func (zs *ZoomState) ClearZoom() {
	zs.mu.Lock()
	defer zs.mu.Unlock()
	zs.zoomed = false
	zs.zoomedPaneID = ""
	zs.savedLayout = nil
}

// ZoomManager integrates with PaneManager for zoom functionality
type ZoomManager struct {
	state       *ZoomState
	paneManager *PaneManager
}

// NewZoomManager creates a new zoom manager
func NewZoomManager(pm *PaneManager) *ZoomManager {
	return &ZoomManager{
		state:       NewZoomState(),
		paneManager: pm,
	}
}

// ToggleZoom toggles zoom on the focused pane
func (zm *ZoomManager) ToggleZoom() bool {
	if zm.state.IsZoomed() {
		return zm.Restore()
	}
	return zm.ZoomFocused()
}

// ZoomFocused zooms the currently focused pane
func (zm *ZoomManager) ZoomFocused() bool {
	focused := zm.paneManager.GetFocusedPane()
	if focused == nil {
		return false
	}

	return zm.ZoomPane(focused.ID)
}

// ZoomPane zooms a specific pane
func (zm *ZoomManager) ZoomPane(paneID string) bool {
	pane := zm.paneManager.GetPane(paneID)
	if pane == nil {
		return false
	}

	// Save current layout
	zm.state.SaveLayout(zm.paneManager.GetAllPanes())

	// Set zoom state
	zm.state.SetZoomed(true, paneID)

	// Resize the pane to full size
	pm := zm.paneManager
	pm.mu.Lock()
	pane.SetSize(pm.width, pm.height)
	pane.X = 0
	pane.Y = 0
	pm.mu.Unlock()

	return true
}

// Restore restores the layout from before zoom
func (zm *ZoomManager) Restore() bool {
	if !zm.state.IsZoomed() {
		return false
	}

	snapshot := zm.state.GetSavedLayout()
	if snapshot == nil {
		zm.state.ClearZoom()
		return false
	}

	// Restore each pane's position and size
	pm := zm.paneManager
	pm.mu.Lock()
	for _, ps := range snapshot.Panes {
		if pane := pm.panes[ps.ID]; pane != nil {
			pane.X = ps.X
			pane.Y = ps.Y
			pane.SetSize(ps.Width, ps.Height)
		}
	}
	pm.mu.Unlock()

	// Clear zoom state
	zm.state.ClearZoom()

	// Recalculate layout
	pm.SetSize(pm.width, pm.height)

	return true
}

// IsZoomed returns whether any pane is zoomed
func (zm *ZoomManager) IsZoomed() bool {
	return zm.state.IsZoomed()
}

// GetZoomedPaneID returns the ID of the zoomed pane
func (zm *ZoomManager) GetZoomedPaneID() string {
	return zm.state.GetZoomedPane()
}

// GetZoomStatus returns a status string for display
func (zm *ZoomManager) GetZoomStatus() string {
	if zm.state.IsZoomed() {
		return "[ZOOMED]"
	}
	return ""
}
