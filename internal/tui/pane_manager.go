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
}

// NewPaneManager creates a new pane manager
func NewPaneManager(width, height int, styles *Styles) *PaneManager {
	pm := &PaneManager{
		panes:  make(map[string]*Pane),
		width:  width,
		height: height,
		styles: styles,
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
