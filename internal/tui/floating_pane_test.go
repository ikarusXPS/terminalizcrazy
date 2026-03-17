package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFloatingPane(t *testing.T) {
	styles := DefaultStyles()
	pane := NewPane("test-1", PaneTypeChat, "Test", 100, 50, styles)
	fp := NewFloatingPane(pane, 10, 20, 1)

	assert.Equal(t, 10, fp.PosX)
	assert.Equal(t, 20, fp.PosY)
	assert.Equal(t, 1, fp.ZIndex)
	assert.False(t, fp.Dragging)
	assert.False(t, fp.Minimized)
}

func TestFloatingPaneSetPosition(t *testing.T) {
	styles := DefaultStyles()
	pane := NewPane("test-1", PaneTypeChat, "Test", 100, 50, styles)
	fp := NewFloatingPane(pane, 0, 0, 1)

	fp.SetPosition(30, 40)

	assert.Equal(t, 30, fp.PosX)
	assert.Equal(t, 40, fp.PosY)
}

func TestFloatingPaneDrag(t *testing.T) {
	styles := DefaultStyles()
	pane := NewPane("test-1", PaneTypeChat, "Test", 100, 50, styles)
	fp := NewFloatingPane(pane, 10, 20, 1)

	// Start drag
	fp.StartDrag(15, 25)
	assert.True(t, fp.Dragging)
	assert.Equal(t, 5, fp.DragOffX)
	assert.Equal(t, 5, fp.DragOffY)

	// Update drag
	fp.UpdateDrag(50, 60)
	assert.Equal(t, 45, fp.PosX)
	assert.Equal(t, 55, fp.PosY)

	// End drag
	fp.EndDrag()
	assert.False(t, fp.Dragging)
}

func TestFloatingPaneToggleMinimize(t *testing.T) {
	styles := DefaultStyles()
	pane := NewPane("test-1", PaneTypeChat, "Test", 100, 50, styles)
	fp := NewFloatingPane(pane, 0, 0, 1)

	assert.False(t, fp.Minimized)

	fp.ToggleMinimize()
	assert.True(t, fp.Minimized)

	fp.ToggleMinimize()
	assert.False(t, fp.Minimized)
}

func TestFloatingPaneContains(t *testing.T) {
	styles := DefaultStyles()
	pane := NewPane("test-1", PaneTypeChat, "Test", 100, 50, styles)
	fp := NewFloatingPane(pane, 10, 20, 1)

	// Inside
	assert.True(t, fp.Contains(10, 20))
	assert.True(t, fp.Contains(50, 40))
	assert.True(t, fp.Contains(109, 69))

	// Outside
	assert.False(t, fp.Contains(9, 20))
	assert.False(t, fp.Contains(10, 19))
	assert.False(t, fp.Contains(110, 20))
	assert.False(t, fp.Contains(10, 70))
}

func TestNewFloatingManager(t *testing.T) {
	fm := NewFloatingManager()

	assert.NotNil(t, fm)
	assert.Equal(t, 0, fm.GetCount())
}

func TestFloatingManagerAdd(t *testing.T) {
	fm := NewFloatingManager()
	styles := DefaultStyles()
	pane := NewPane("test-1", PaneTypeChat, "Test", 100, 50, styles)

	fp := fm.Add(pane, 10, 20)

	assert.NotNil(t, fp)
	assert.Equal(t, 1, fm.GetCount())
	assert.Equal(t, 1, fp.ZIndex)
}

func TestFloatingManagerRemove(t *testing.T) {
	fm := NewFloatingManager()
	styles := DefaultStyles()
	pane := NewPane("test-1", PaneTypeChat, "Test", 100, 50, styles)

	fm.Add(pane, 10, 20)
	assert.Equal(t, 1, fm.GetCount())

	result := fm.Remove("test-1")
	assert.True(t, result)
	assert.Equal(t, 0, fm.GetCount())

	// Remove non-existent
	result = fm.Remove("test-999")
	assert.False(t, result)
}

func TestFloatingManagerGet(t *testing.T) {
	fm := NewFloatingManager()
	styles := DefaultStyles()
	pane := NewPane("test-1", PaneTypeChat, "Test", 100, 50, styles)

	fm.Add(pane, 10, 20)

	fp := fm.Get("test-1")
	assert.NotNil(t, fp)
	assert.Equal(t, "test-1", fp.ID)

	// Get non-existent
	nilFp := fm.Get("test-999")
	assert.Nil(t, nilFp)
}

func TestFloatingManagerBringToFront(t *testing.T) {
	fm := NewFloatingManager()
	styles := DefaultStyles()

	pane1 := NewPane("test-1", PaneTypeChat, "Test 1", 100, 50, styles)
	pane2 := NewPane("test-2", PaneTypeChat, "Test 2", 100, 50, styles)

	fp1 := fm.Add(pane1, 0, 0)
	fp2 := fm.Add(pane2, 10, 10)

	assert.Equal(t, 1, fp1.ZIndex)
	assert.Equal(t, 2, fp2.ZIndex)

	fm.BringToFront("test-1")

	fp1 = fm.Get("test-1")
	assert.Equal(t, 3, fp1.ZIndex)
}

func TestFloatingManagerGetAtPosition(t *testing.T) {
	fm := NewFloatingManager()
	styles := DefaultStyles()

	pane1 := NewPane("test-1", PaneTypeChat, "Test 1", 50, 30, styles)
	pane2 := NewPane("test-2", PaneTypeChat, "Test 2", 50, 30, styles)

	fm.Add(pane1, 0, 0)
	fm.Add(pane2, 25, 15)

	// Position only in pane1
	fp := fm.GetAtPosition(10, 10)
	assert.NotNil(t, fp)
	assert.Equal(t, "test-1", fp.ID)

	// Position only in pane2
	fp = fm.GetAtPosition(60, 35)
	assert.NotNil(t, fp)
	assert.Equal(t, "test-2", fp.ID)

	// Position in overlap (pane2 has higher z-index)
	fp = fm.GetAtPosition(30, 20)
	assert.NotNil(t, fp)
	assert.Equal(t, "test-2", fp.ID)

	// Position outside all panes
	fp = fm.GetAtPosition(200, 200)
	assert.Nil(t, fp)
}

func TestFloatingManagerGetAll(t *testing.T) {
	fm := NewFloatingManager()
	styles := DefaultStyles()

	pane1 := NewPane("test-1", PaneTypeChat, "Test 1", 50, 30, styles)
	pane2 := NewPane("test-2", PaneTypeChat, "Test 2", 50, 30, styles)
	pane3 := NewPane("test-3", PaneTypeChat, "Test 3", 50, 30, styles)

	fm.Add(pane1, 0, 0)
	fm.Add(pane2, 10, 10)
	fm.Add(pane3, 20, 20)

	all := fm.GetAll()
	assert.Len(t, all, 3)

	// Should be sorted by z-index (ascending)
	assert.Equal(t, 1, all[0].ZIndex)
	assert.Equal(t, 2, all[1].ZIndex)
	assert.Equal(t, 3, all[2].ZIndex)
}

func TestFloatingManagerClear(t *testing.T) {
	fm := NewFloatingManager()
	styles := DefaultStyles()

	pane1 := NewPane("test-1", PaneTypeChat, "Test 1", 50, 30, styles)
	pane2 := NewPane("test-2", PaneTypeChat, "Test 2", 50, 30, styles)

	fm.Add(pane1, 0, 0)
	fm.Add(pane2, 10, 10)

	assert.Equal(t, 2, fm.GetCount())

	fm.Clear()

	assert.Equal(t, 0, fm.GetCount())
}

func TestFloatingManagerIsFloating(t *testing.T) {
	fm := NewFloatingManager()
	styles := DefaultStyles()

	pane := NewPane("test-1", PaneTypeChat, "Test", 100, 50, styles)
	fm.Add(pane, 0, 0)

	assert.True(t, fm.IsFloating("test-1"))
	assert.False(t, fm.IsFloating("test-999"))
}

func TestFloatingManagerConstrainToScreen(t *testing.T) {
	fm := NewFloatingManager()
	styles := DefaultStyles()

	pane := NewPane("test-1", PaneTypeChat, "Test", 50, 30, styles)
	fp := fm.Add(pane, -10, -10)

	fm.ConstrainToScreen(100, 80)

	assert.Equal(t, 0, fp.PosX)
	assert.Equal(t, 0, fp.PosY)

	// Test overflow
	fp.SetPosition(80, 60)
	fm.ConstrainToScreen(100, 80)

	assert.Equal(t, 50, fp.PosX) // 100 - 50
	assert.Equal(t, 50, fp.PosY) // 80 - 30
}
