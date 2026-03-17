package tui

import (
	"fmt"
	"sync"

	"github.com/charmbracelet/lipgloss"
)

// FloatingPane represents a pane that can float above other panes
type FloatingPane struct {
	*Pane
	ZIndex    int
	PosX      int // Screen position X
	PosY      int // Screen position Y
	Dragging  bool
	DragOffX  int
	DragOffY  int
	Minimized bool
}

// NewFloatingPane creates a new floating pane from an existing pane
func NewFloatingPane(pane *Pane, x, y, zIndex int) *FloatingPane {
	return &FloatingPane{
		Pane:   pane,
		PosX:   x,
		PosY:   y,
		ZIndex: zIndex,
	}
}

// SetPosition sets the floating pane position
func (fp *FloatingPane) SetPosition(x, y int) {
	fp.PosX = x
	fp.PosY = y
}

// StartDrag starts dragging the pane
func (fp *FloatingPane) StartDrag(mouseX, mouseY int) {
	fp.Dragging = true
	fp.DragOffX = mouseX - fp.PosX
	fp.DragOffY = mouseY - fp.PosY
}

// UpdateDrag updates the pane position while dragging
func (fp *FloatingPane) UpdateDrag(mouseX, mouseY int) {
	if fp.Dragging {
		fp.PosX = mouseX - fp.DragOffX
		fp.PosY = mouseY - fp.DragOffY
	}
}

// EndDrag ends dragging the pane
func (fp *FloatingPane) EndDrag() {
	fp.Dragging = false
}

// ToggleMinimize toggles the minimized state
func (fp *FloatingPane) ToggleMinimize() {
	fp.Minimized = !fp.Minimized
}

// Contains checks if a point is inside the floating pane
func (fp *FloatingPane) Contains(x, y int) bool {
	return x >= fp.PosX && x < fp.PosX+fp.Width &&
		y >= fp.PosY && y < fp.PosY+fp.Height
}

// View renders the floating pane with floating-specific styling
func (fp *FloatingPane) View() string {
	if fp.Minimized {
		// Show only title bar when minimized
		return fp.styles.PaneTitle.
			Width(fp.Width).
			Render(fmt.Sprintf("▼ %s", fp.Title))
	}

	// Add shadow effect for floating appearance
	content := fp.Pane.View()

	// Add floating indicator
	if fp.Focused {
		content = lipgloss.JoinVertical(lipgloss.Left,
			fp.styles.PaneTitle.Render(fmt.Sprintf("◇ %s [floating]", fp.Title)),
			content,
		)
	}

	return content
}

// FloatingManager manages floating panes
type FloatingManager struct {
	panes     []*FloatingPane
	maxZIndex int
	mu        sync.RWMutex
}

// NewFloatingManager creates a new floating pane manager
func NewFloatingManager() *FloatingManager {
	return &FloatingManager{
		panes:     make([]*FloatingPane, 0),
		maxZIndex: 0,
	}
}

// Add adds a pane to the floating manager
func (fm *FloatingManager) Add(pane *Pane, x, y int) *FloatingPane {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fm.maxZIndex++
	fp := NewFloatingPane(pane, x, y, fm.maxZIndex)
	fm.panes = append(fm.panes, fp)
	return fp
}

// Remove removes a floating pane by pane ID
func (fm *FloatingManager) Remove(paneID string) bool {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	for i, fp := range fm.panes {
		if fp.ID == paneID {
			fm.panes = append(fm.panes[:i], fm.panes[i+1:]...)
			return true
		}
	}
	return false
}

// Get returns a floating pane by pane ID
func (fm *FloatingManager) Get(paneID string) *FloatingPane {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	for _, fp := range fm.panes {
		if fp.ID == paneID {
			return fp
		}
	}
	return nil
}

// BringToFront brings a floating pane to the front
func (fm *FloatingManager) BringToFront(paneID string) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	for _, fp := range fm.panes {
		if fp.ID == paneID {
			fm.maxZIndex++
			fp.ZIndex = fm.maxZIndex
			return
		}
	}
}

// GetAtPosition returns the topmost floating pane at the given position
func (fm *FloatingManager) GetAtPosition(x, y int) *FloatingPane {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	var topPane *FloatingPane
	topZ := -1

	for _, fp := range fm.panes {
		if fp.Contains(x, y) && fp.ZIndex > topZ {
			topPane = fp
			topZ = fp.ZIndex
		}
	}

	return topPane
}

// GetAll returns all floating panes sorted by z-index
func (fm *FloatingManager) GetAll() []*FloatingPane {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	// Copy and sort by z-index
	result := make([]*FloatingPane, len(fm.panes))
	copy(result, fm.panes)

	// Simple bubble sort by z-index
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].ZIndex > result[j].ZIndex {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result
}

// GetCount returns the number of floating panes
func (fm *FloatingManager) GetCount() int {
	fm.mu.RLock()
	defer fm.mu.RUnlock()
	return len(fm.panes)
}

// Clear removes all floating panes
func (fm *FloatingManager) Clear() {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.panes = make([]*FloatingPane, 0)
	fm.maxZIndex = 0
}

// IsFloating checks if a pane is currently floating
func (fm *FloatingManager) IsFloating(paneID string) bool {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	for _, fp := range fm.panes {
		if fp.ID == paneID {
			return true
		}
	}
	return false
}

// ConstrainToScreen ensures all floating panes are within screen bounds
func (fm *FloatingManager) ConstrainToScreen(screenWidth, screenHeight int) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	for _, fp := range fm.panes {
		// Constrain X
		if fp.PosX < 0 {
			fp.PosX = 0
		}
		if fp.PosX+fp.Width > screenWidth {
			fp.PosX = screenWidth - fp.Width
		}

		// Constrain Y
		if fp.PosY < 0 {
			fp.PosY = 0
		}
		if fp.PosY+fp.Height > screenHeight {
			fp.PosY = screenHeight - fp.Height
		}
	}
}
