package tui

import (
	"fmt"
	"sync"
)

// PaneManager manages all panes and their layout
type PaneManager struct {
	root      *PaneNode
	panes     map[string]*Pane
	focusedID string
	width     int
	height    int
	styles    *Styles
	paneCount int
	mu        sync.RWMutex

	// Enhanced features
	floating   *FloatingManager
	zoom       *ZoomState
	inputSync  *InputSynchronizer
}

// NewPaneManager creates a new pane manager
func NewPaneManager(width, height int, styles *Styles) *PaneManager {
	pm := &PaneManager{
		panes:     make(map[string]*Pane),
		width:     width,
		height:    height,
		styles:    styles,
		floating:  NewFloatingManager(),
		zoom:      NewZoomState(),
		inputSync: NewInputSynchronizer(),
	}

	// Create default pane
	mainPane := pm.createPane(PaneTypeChat, "Main")
	mainPane.Focus()
	pm.focusedID = mainPane.ID

	pm.root = NewPaneNode(mainPane)
	pm.root.SetSize(0, 0, width, height)

	return pm
}

// createPane creates a new pane with auto-generated ID
func (pm *PaneManager) createPane(paneType PaneType, title string) *Pane {
	pm.paneCount++
	id := fmt.Sprintf("pane-%d", pm.paneCount)
	pane := NewPane(id, paneType, title, pm.width, pm.height, pm.styles)
	pm.panes[id] = pane
	return pane
}

// SetSize updates the manager dimensions
func (pm *PaneManager) SetSize(width, height int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.width = width
	pm.height = height

	if pm.root != nil {
		pm.root.SetSize(0, 0, width, height)
	}
}

// SplitVertical splits the focused pane vertically (left/right)
func (pm *PaneManager) SplitVertical(paneType PaneType, title string) *Pane {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	focused := pm.root.GetFocusedPane()
	if focused == nil {
		return nil
	}

	newPane := pm.createPane(paneType, title)

	// Find the node containing the focused pane and split it
	pm.splitNodeContaining(pm.root, focused.ID, SplitVertical, newPane)

	return newPane
}

// SplitHorizontal splits the focused pane horizontally (top/bottom)
func (pm *PaneManager) SplitHorizontal(paneType PaneType, title string) *Pane {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	focused := pm.root.GetFocusedPane()
	if focused == nil {
		return nil
	}

	newPane := pm.createPane(paneType, title)

	// Find the node containing the focused pane and split it
	pm.splitNodeContaining(pm.root, focused.ID, SplitHorizontal, newPane)

	return newPane
}

// splitNodeContaining finds and splits the node containing the given pane ID
func (pm *PaneManager) splitNodeContaining(node *PaneNode, paneID string, direction SplitDirection, newPane *Pane) bool {
	if node.IsLeaf {
		if node.Pane != nil && node.Pane.ID == paneID {
			node.Split(direction, newPane)
			return true
		}
		return false
	}

	for _, child := range node.Children {
		if child != nil && pm.splitNodeContaining(child, paneID, direction, newPane) {
			return true
		}
	}

	return false
}

// ClosePane closes a pane by ID
func (pm *PaneManager) ClosePane(id string) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Don't close if it's the last pane
	if len(pm.panes) <= 1 {
		return false
	}

	// Remove from tree
	if pm.root.RemovePane(id) {
		delete(pm.panes, id)

		// If we closed the focused pane, focus another
		if pm.focusedID == id {
			panes := pm.root.GetPanes()
			if len(panes) > 0 {
				panes[0].Focus()
				pm.focusedID = panes[0].ID
			}
		}

		return true
	}

	return false
}

// CloseFocusedPane closes the currently focused pane
func (pm *PaneManager) CloseFocusedPane() bool {
	pm.mu.RLock()
	id := pm.focusedID
	pm.mu.RUnlock()

	return pm.ClosePane(id)
}

// FocusPane focuses a specific pane
func (pm *PaneManager) FocusPane(id string) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pane, ok := pm.panes[id]
	if !ok {
		return false
	}

	// Blur current focused pane
	if current := pm.root.GetFocusedPane(); current != nil {
		current.Blur()
	}

	pane.Focus()
	pm.focusedID = id
	return true
}

// FocusNext focuses the next pane
func (pm *PaneManager) FocusNext() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.root.FocusNext()

	// Update focused ID
	if focused := pm.root.GetFocusedPane(); focused != nil {
		pm.focusedID = focused.ID
	}
}

// FocusPrevious focuses the previous pane
func (pm *PaneManager) FocusPrevious() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.root.FocusPrevious()

	// Update focused ID
	if focused := pm.root.GetFocusedPane(); focused != nil {
		pm.focusedID = focused.ID
	}
}

// FocusDirection focuses a pane in the given direction
func (pm *PaneManager) FocusDirection(direction string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.root.FocusDirection(direction)

	// Update focused ID
	if focused := pm.root.GetFocusedPane(); focused != nil {
		pm.focusedID = focused.ID
	}
}

// GetFocusedPane returns the currently focused pane
func (pm *PaneManager) GetFocusedPane() *Pane {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return pm.root.GetFocusedPane()
}

// GetPane returns a pane by ID
func (pm *PaneManager) GetPane(id string) *Pane {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return pm.panes[id]
}

// GetAllPanes returns all panes
func (pm *PaneManager) GetAllPanes() []*Pane {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return pm.root.GetPanes()
}

// GetPaneCount returns the number of panes
func (pm *PaneManager) GetPaneCount() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return len(pm.panes)
}

// View renders all panes
func (pm *PaneManager) View() string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return pm.root.View()
}

// SetPaneContent sets content for a specific pane
func (pm *PaneManager) SetPaneContent(id, content string) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pane, ok := pm.panes[id]
	if !ok {
		return false
	}

	pane.SetContent(content)
	return true
}

// AppendPaneContent appends content to a specific pane
func (pm *PaneManager) AppendPaneContent(id, content string) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pane, ok := pm.panes[id]
	if !ok {
		return false
	}

	pane.AppendContent(content)
	return true
}

// SetFocusedPaneContent sets content for the focused pane
func (pm *PaneManager) SetFocusedPaneContent(content string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if focused := pm.root.GetFocusedPane(); focused != nil {
		focused.SetContent(content)
	}
}

// GetPanesByType returns all panes of a specific type
func (pm *PaneManager) GetPanesByType(paneType PaneType) []*Pane {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var result []*Pane
	for _, pane := range pm.panes {
		if pane.Type == paneType {
			result = append(result, pane)
		}
	}
	return result
}

// ResizeFocused adjusts the split ratio for the focused pane
func (pm *PaneManager) ResizeFocused(delta float64) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Find the parent node of the focused pane and adjust its ratio
	focused := pm.root.GetFocusedPane()
	if focused == nil {
		return
	}

	pm.adjustNodeRatio(pm.root, focused.ID, delta)
}

// adjustNodeRatio finds and adjusts the ratio of the parent node
func (pm *PaneManager) adjustNodeRatio(node *PaneNode, paneID string, delta float64) bool {
	if node.IsLeaf {
		return node.Pane != nil && node.Pane.ID == paneID
	}

	for i, child := range node.Children {
		if child != nil && pm.adjustNodeRatio(child, paneID, delta) {
			// Adjust ratio based on which child contains the pane
			if i == 0 {
				node.Ratio += delta
			} else {
				node.Ratio -= delta
			}

			// Clamp ratio
			if node.Ratio < 0.2 {
				node.Ratio = 0.2
			}
			if node.Ratio > 0.8 {
				node.Ratio = 0.8
			}

			node.recalculateSizes()
			return true
		}
	}

	return false
}

// SwapPanes swaps two panes by ID
func (pm *PaneManager) SwapPanes(id1, id2 string) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pane1, ok1 := pm.panes[id1]
	pane2, ok2 := pm.panes[id2]

	if !ok1 || !ok2 {
		return false
	}

	// Swap content and type
	pane1.Content, pane2.Content = pane2.Content, pane1.Content
	pane1.Type, pane2.Type = pane2.Type, pane1.Type
	pane1.Title, pane2.Title = pane2.Title, pane1.Title

	// Update viewports
	pane1.Viewport.SetContent(pane1.Content)
	pane2.Viewport.SetContent(pane2.Content)

	return true
}

// MaximizeFocused toggles maximizing the focused pane
func (pm *PaneManager) MaximizeFocused() {
	// This would temporarily hide other panes
	// Implementation depends on how we want to handle this state
	// For now, we'll just focus the pane
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	// The actual maximize logic would need additional state tracking
	// This is a placeholder for the API
}

// GetLayoutInfo returns information about the current layout
func (pm *PaneManager) GetLayoutInfo() string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return pm.root.PaneInfo()
}

// ResetLayout resets to a single pane
func (pm *PaneManager) ResetLayout() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Get current focused content
	var content string
	var paneType PaneType
	var title string

	if focused := pm.root.GetFocusedPane(); focused != nil {
		content = focused.Content
		paneType = focused.Type
		title = focused.Title
	} else {
		paneType = PaneTypeChat
		title = "Main"
	}

	// Clear existing panes
	pm.panes = make(map[string]*Pane)
	pm.paneCount = 0

	// Create new single pane
	mainPane := pm.createPane(paneType, title)
	mainPane.SetContent(content)
	mainPane.Focus()
	pm.focusedID = mainPane.ID

	pm.root = NewPaneNode(mainPane)
	pm.root.SetSize(0, 0, pm.width, pm.height)
}

// DuplicateFocused creates a duplicate of the focused pane
func (pm *PaneManager) DuplicateFocused() *Pane {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	focused := pm.root.GetFocusedPane()
	if focused == nil {
		return nil
	}

	// Create new pane with same type and content
	newPane := pm.createPane(focused.Type, focused.Title+" (copy)")
	newPane.SetContent(focused.Content)

	// Split vertically by default
	pm.splitNodeContaining(pm.root, focused.ID, SplitVertical, newPane)

	return newPane
}

// --- Zoom Methods ---

// ToggleZoom toggles zoom on the focused pane
func (pm *PaneManager) ToggleZoom() bool {
	if pm.zoom.IsZoomed() {
		return pm.RestoreZoom()
	}
	return pm.ZoomFocused()
}

// ZoomFocused zooms the currently focused pane to fill the entire area
func (pm *PaneManager) ZoomFocused() bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	focused := pm.root.GetFocusedPane()
	if focused == nil {
		return false
	}

	// Save current layout
	pm.zoom.SaveLayout(pm.root.GetPanes())
	pm.zoom.SetZoomed(true, focused.ID)

	// Set the focused pane to full size (will be handled in View)
	return true
}

// RestoreZoom restores the layout from before zoom
func (pm *PaneManager) RestoreZoom() bool {
	if !pm.zoom.IsZoomed() {
		return false
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	snapshot := pm.zoom.GetSavedLayout()
	if snapshot != nil {
		// Restore pane sizes from snapshot
		for _, ps := range snapshot.Panes {
			if pane, ok := pm.panes[ps.ID]; ok {
				pane.X = ps.X
				pane.Y = ps.Y
				pane.SetSize(ps.Width, ps.Height)
			}
		}
	}

	pm.zoom.ClearZoom()

	// Recalculate layout
	if pm.root != nil {
		pm.root.SetSize(0, 0, pm.width, pm.height)
	}

	return true
}

// IsZoomed returns whether any pane is currently zoomed
func (pm *PaneManager) IsZoomed() bool {
	return pm.zoom.IsZoomed()
}

// GetZoomedPaneID returns the ID of the zoomed pane
func (pm *PaneManager) GetZoomedPaneID() string {
	return pm.zoom.GetZoomedPane()
}

// --- Floating Pane Methods ---

// ToggleFloating toggles floating mode for the focused pane
func (pm *PaneManager) ToggleFloating() bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	focused := pm.root.GetFocusedPane()
	if focused == nil {
		return false
	}

	if pm.floating.IsFloating(focused.ID) {
		// Remove from floating
		return pm.floating.Remove(focused.ID)
	}

	// Add to floating at center of screen
	x := (pm.width - focused.Width) / 2
	y := (pm.height - focused.Height) / 2
	pm.floating.Add(focused, x, y)
	return true
}

// MakeFloating makes a pane floating
func (pm *PaneManager) MakeFloating(paneID string) *FloatingPane {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pane, ok := pm.panes[paneID]
	if !ok {
		return nil
	}

	// Center the floating pane
	x := (pm.width - pane.Width) / 2
	y := (pm.height - pane.Height) / 2

	return pm.floating.Add(pane, x, y)
}

// DockFloating docks a floating pane back to the layout
func (pm *PaneManager) DockFloating(paneID string) bool {
	return pm.floating.Remove(paneID)
}

// IsFloating checks if a pane is floating
func (pm *PaneManager) IsFloating(paneID string) bool {
	return pm.floating.IsFloating(paneID)
}

// GetFloatingPanes returns all floating panes
func (pm *PaneManager) GetFloatingPanes() []*FloatingPane {
	return pm.floating.GetAll()
}

// BringFloatingToFront brings a floating pane to the front
func (pm *PaneManager) BringFloatingToFront(paneID string) {
	pm.floating.BringToFront(paneID)
}

// --- Input Synchronization Methods ---

// ToggleBroadcast toggles broadcast mode (send input to all panes)
func (pm *PaneManager) ToggleBroadcast() bool {
	return pm.inputSync.ToggleBroadcast()
}

// IsBroadcastEnabled returns whether broadcast mode is enabled
func (pm *PaneManager) IsBroadcastEnabled() bool {
	return pm.inputSync.IsBroadcastEnabled()
}

// CreateSyncGroup creates a new input sync group
func (pm *PaneManager) CreateSyncGroup(id, name string) *SyncGroup {
	return pm.inputSync.CreateGroup(id, name)
}

// AddPaneToSyncGroup adds a pane to a sync group
func (pm *PaneManager) AddPaneToSyncGroup(groupID, paneID string) bool {
	return pm.inputSync.AddPaneToGroup(groupID, paneID)
}

// RemovePaneFromSyncGroup removes a pane from its sync group
func (pm *PaneManager) RemovePaneFromSyncGroup(paneID string) bool {
	return pm.inputSync.RemovePaneFromGroup(paneID)
}

// GetSyncTargets returns panes that should receive synced input from a source pane
func (pm *PaneManager) GetSyncTargets(sourcePaneID string) []string {
	pm.mu.RLock()
	allPaneIDs := make([]string, 0, len(pm.panes))
	for id := range pm.panes {
		allPaneIDs = append(allPaneIDs, id)
	}
	pm.mu.RUnlock()

	return pm.inputSync.GetSyncTargets(sourcePaneID, allPaneIDs)
}

// GetInputSyncStatus returns a status string for the input synchronizer
func (pm *PaneManager) GetInputSyncStatus() string {
	return pm.inputSync.GetStatus()
}

// --- View with Enhanced Features ---

// ViewWithEnhancements renders all panes including floating and zoom states
func (pm *PaneManager) ViewWithEnhancements() string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	// If zoomed, show only the zoomed pane
	if pm.zoom.IsZoomed() {
		zoomedID := pm.zoom.GetZoomedPane()
		if pane, ok := pm.panes[zoomedID]; ok {
			// Temporarily set pane to full size for rendering
			originalW, originalH := pane.Width, pane.Height
			pane.SetSize(pm.width, pm.height)
			view := pane.View()
			pane.SetSize(originalW, originalH)
			return view
		}
	}

	// Normal view
	baseView := pm.root.View()

	// Overlay floating panes (they render on top)
	// Note: In a real TUI, this would require proper compositing
	// For now, floating panes are managed separately
	return baseView
}

// GetStatusLine returns a status line showing current states
func (pm *PaneManager) GetStatusLine() string {
	status := ""

	if pm.zoom.IsZoomed() {
		status += "[ZOOMED] "
	}

	if pm.floating.GetCount() > 0 {
		status += fmt.Sprintf("[%d FLOATING] ", pm.floating.GetCount())
	}

	if syncStatus := pm.inputSync.GetStatus(); syncStatus != "" {
		status += syncStatus + " "
	}

	return status
}
