package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaneTypes(t *testing.T) {
	assert.Equal(t, PaneType("chat"), PaneTypeChat)
	assert.Equal(t, PaneType("terminal"), PaneTypeTerminal)
	assert.Equal(t, PaneType("plan"), PaneTypePlan)
	assert.Equal(t, PaneType("output"), PaneTypeOutput)
}

func TestNewPane(t *testing.T) {
	styles := DefaultStyles()
	pane := NewPane("test-1", PaneTypeChat, "Test Title", 100, 50, styles)

	assert.Equal(t, "test-1", pane.ID)
	assert.Equal(t, PaneTypeChat, pane.Type)
	assert.Equal(t, "Test Title", pane.Title)
	assert.Equal(t, 100, pane.Width)
	assert.Equal(t, 50, pane.Height)
	assert.False(t, pane.Focused)
	assert.NotNil(t, pane.styles)
}

func TestPane_SetSize(t *testing.T) {
	styles := DefaultStyles()
	pane := NewPane("test-1", PaneTypeChat, "Test", 100, 50, styles)

	pane.SetSize(200, 100)

	assert.Equal(t, 200, pane.Width)
	assert.Equal(t, 100, pane.Height)
}

func TestPane_SetContent(t *testing.T) {
	styles := DefaultStyles()
	pane := NewPane("test-1", PaneTypeChat, "Test", 100, 50, styles)

	pane.SetContent("Hello World")

	assert.Equal(t, "Hello World", pane.Content)
}

func TestPane_AppendContent(t *testing.T) {
	styles := DefaultStyles()
	pane := NewPane("test-1", PaneTypeChat, "Test", 100, 50, styles)

	pane.SetContent("Hello")
	pane.AppendContent(" World")

	assert.Equal(t, "Hello World", pane.Content)
}

func TestPane_FocusBlur(t *testing.T) {
	styles := DefaultStyles()
	pane := NewPane("test-1", PaneTypeChat, "Test", 100, 50, styles)

	assert.False(t, pane.Focused)

	pane.Focus()
	assert.True(t, pane.Focused)

	pane.Blur()
	assert.False(t, pane.Focused)
}

func TestPane_View(t *testing.T) {
	styles := DefaultStyles()
	pane := NewPane("test-1", PaneTypeChat, "Test", 100, 50, styles)
	pane.SetContent("Content here")

	view := pane.View()

	assert.NotEmpty(t, view)
	assert.Contains(t, view, "Test")
}

func TestPane_View_Focused(t *testing.T) {
	styles := DefaultStyles()
	pane := NewPane("test-1", PaneTypeChat, "Test", 100, 50, styles)
	pane.Focus()

	view := pane.View()

	assert.NotEmpty(t, view)
	assert.Contains(t, view, "●") // Focus indicator
}

func TestSplitDirections(t *testing.T) {
	assert.Equal(t, SplitDirection(0), SplitHorizontal)
	assert.Equal(t, SplitDirection(1), SplitVertical)
}

func TestNewPaneNode(t *testing.T) {
	styles := DefaultStyles()
	pane := NewPane("test-1", PaneTypeChat, "Test", 100, 50, styles)
	node := NewPaneNode(pane)

	assert.True(t, node.IsLeaf)
	assert.Equal(t, pane, node.Pane)
	assert.Equal(t, 0.5, node.Ratio)
}

func TestPaneNode_Split(t *testing.T) {
	styles := DefaultStyles()
	pane1 := NewPane("pane-1", PaneTypeChat, "Pane 1", 100, 50, styles)
	pane2 := NewPane("pane-2", PaneTypeChat, "Pane 2", 100, 50, styles)

	node := NewPaneNode(pane1)
	node.SetSize(0, 0, 200, 100)

	node.Split(SplitVertical, pane2)

	assert.False(t, node.IsLeaf)
	assert.Equal(t, SplitVertical, node.SplitDir)
	assert.NotNil(t, node.Children[0])
	assert.NotNil(t, node.Children[1])
}

func TestPaneNode_Split_OnlyLeaf(t *testing.T) {
	styles := DefaultStyles()
	pane1 := NewPane("pane-1", PaneTypeChat, "Pane 1", 100, 50, styles)
	pane2 := NewPane("pane-2", PaneTypeChat, "Pane 2", 100, 50, styles)
	pane3 := NewPane("pane-3", PaneTypeChat, "Pane 3", 100, 50, styles)

	node := NewPaneNode(pane1)
	node.SetSize(0, 0, 200, 100)
	node.Split(SplitVertical, pane2)

	// Try to split a non-leaf node
	node.Split(SplitHorizontal, pane3)

	// Should still be a split node with 2 children
	assert.False(t, node.IsLeaf)
	assert.NotNil(t, node.Children[0])
	assert.NotNil(t, node.Children[1])
}

func TestPaneNode_SetSize(t *testing.T) {
	styles := DefaultStyles()
	pane := NewPane("test-1", PaneTypeChat, "Test", 100, 50, styles)
	node := NewPaneNode(pane)

	node.SetSize(10, 20, 200, 100)

	assert.Equal(t, 10, node.X)
	assert.Equal(t, 20, node.Y)
	assert.Equal(t, 200, node.Width)
	assert.Equal(t, 100, node.Height)
}

func TestPaneNode_GetPanes(t *testing.T) {
	styles := DefaultStyles()
	pane1 := NewPane("pane-1", PaneTypeChat, "Pane 1", 100, 50, styles)
	pane2 := NewPane("pane-2", PaneTypeChat, "Pane 2", 100, 50, styles)

	node := NewPaneNode(pane1)
	node.SetSize(0, 0, 200, 100)
	node.Split(SplitVertical, pane2)

	panes := node.GetPanes()

	assert.Len(t, panes, 2)
}

func TestPaneNode_GetPanes_Empty(t *testing.T) {
	node := &PaneNode{IsLeaf: true, Pane: nil}

	panes := node.GetPanes()

	assert.Empty(t, panes)
}

func TestPaneNode_FindPane(t *testing.T) {
	styles := DefaultStyles()
	pane1 := NewPane("pane-1", PaneTypeChat, "Pane 1", 100, 50, styles)
	pane2 := NewPane("pane-2", PaneTypeChat, "Pane 2", 100, 50, styles)

	node := NewPaneNode(pane1)
	node.SetSize(0, 0, 200, 100)
	node.Split(SplitVertical, pane2)

	found := node.FindPane("pane-1")
	assert.NotNil(t, found)
	assert.Equal(t, "pane-1", found.ID)

	found = node.FindPane("pane-2")
	assert.NotNil(t, found)
	assert.Equal(t, "pane-2", found.ID)

	notFound := node.FindPane("pane-999")
	assert.Nil(t, notFound)
}

func TestPaneNode_RemovePane(t *testing.T) {
	styles := DefaultStyles()
	pane1 := NewPane("pane-1", PaneTypeChat, "Pane 1", 100, 50, styles)
	pane2 := NewPane("pane-2", PaneTypeChat, "Pane 2", 100, 50, styles)

	node := NewPaneNode(pane1)
	node.SetSize(0, 0, 200, 100)
	node.Split(SplitVertical, pane2)

	removed := node.RemovePane("pane-2")

	assert.True(t, removed)
	panes := node.GetPanes()
	assert.Len(t, panes, 1)
	assert.Equal(t, "pane-1", panes[0].ID)
}

func TestPaneNode_RemovePane_NotFound(t *testing.T) {
	styles := DefaultStyles()
	pane1 := NewPane("pane-1", PaneTypeChat, "Pane 1", 100, 50, styles)
	pane2 := NewPane("pane-2", PaneTypeChat, "Pane 2", 100, 50, styles)

	node := NewPaneNode(pane1)
	node.SetSize(0, 0, 200, 100)
	node.Split(SplitVertical, pane2)

	removed := node.RemovePane("pane-999")

	assert.False(t, removed)
}

func TestPaneNode_View(t *testing.T) {
	styles := DefaultStyles()
	pane1 := NewPane("pane-1", PaneTypeChat, "Pane 1", 100, 50, styles)

	node := NewPaneNode(pane1)
	node.SetSize(0, 0, 100, 50)

	view := node.View()

	assert.NotEmpty(t, view)
}

func TestPaneNode_View_Split(t *testing.T) {
	styles := DefaultStyles()
	pane1 := NewPane("pane-1", PaneTypeChat, "Pane 1", 100, 50, styles)
	pane2 := NewPane("pane-2", PaneTypeChat, "Pane 2", 100, 50, styles)

	node := NewPaneNode(pane1)
	node.SetSize(0, 0, 200, 100)
	node.Split(SplitVertical, pane2)

	view := node.View()

	assert.NotEmpty(t, view)
}

func TestPaneNode_GetFocusedPane(t *testing.T) {
	styles := DefaultStyles()
	pane1 := NewPane("pane-1", PaneTypeChat, "Pane 1", 100, 50, styles)
	pane2 := NewPane("pane-2", PaneTypeChat, "Pane 2", 100, 50, styles)

	node := NewPaneNode(pane1)
	node.SetSize(0, 0, 200, 100)
	node.Split(SplitVertical, pane2)

	// No focus initially
	focused := node.GetFocusedPane()
	assert.Nil(t, focused)

	// Focus pane 2
	pane2.Focus()
	focused = node.GetFocusedPane()
	assert.NotNil(t, focused)
	assert.Equal(t, "pane-2", focused.ID)
}

func TestPaneNode_FocusNext(t *testing.T) {
	styles := DefaultStyles()
	pane1 := NewPane("pane-1", PaneTypeChat, "Pane 1", 100, 50, styles)
	pane2 := NewPane("pane-2", PaneTypeChat, "Pane 2", 100, 50, styles)

	node := NewPaneNode(pane1)
	node.SetSize(0, 0, 200, 100)
	node.Split(SplitVertical, pane2)

	pane1.Focus()

	node.FocusNext()

	assert.False(t, pane1.Focused)
	assert.True(t, pane2.Focused)
}

func TestPaneNode_FocusNext_Wrap(t *testing.T) {
	styles := DefaultStyles()
	pane1 := NewPane("pane-1", PaneTypeChat, "Pane 1", 100, 50, styles)
	pane2 := NewPane("pane-2", PaneTypeChat, "Pane 2", 100, 50, styles)

	node := NewPaneNode(pane1)
	node.SetSize(0, 0, 200, 100)
	node.Split(SplitVertical, pane2)

	pane2.Focus()

	node.FocusNext()

	panes := node.GetPanes()
	assert.True(t, panes[0].Focused)
}

func TestPaneNode_FocusPrevious(t *testing.T) {
	styles := DefaultStyles()
	pane1 := NewPane("pane-1", PaneTypeChat, "Pane 1", 100, 50, styles)
	pane2 := NewPane("pane-2", PaneTypeChat, "Pane 2", 100, 50, styles)

	node := NewPaneNode(pane1)
	node.SetSize(0, 0, 200, 100)
	node.Split(SplitVertical, pane2)

	pane2.Focus()

	node.FocusPrevious()

	panes := node.GetPanes()
	assert.True(t, panes[0].Focused)
	assert.False(t, panes[1].Focused)
}

func TestPaneNode_FocusDirection(t *testing.T) {
	styles := DefaultStyles()
	pane1 := NewPane("pane-1", PaneTypeChat, "Pane 1", 100, 50, styles)
	pane2 := NewPane("pane-2", PaneTypeChat, "Pane 2", 100, 50, styles)

	node := NewPaneNode(pane1)
	node.SetSize(0, 0, 200, 100)
	node.Split(SplitVertical, pane2)

	panes := node.GetPanes()
	panes[0].Focus()

	// Focus right
	node.FocusDirection("right")
	assert.False(t, panes[0].Focused)
	assert.True(t, panes[1].Focused)

	// Focus left
	node.FocusDirection("left")
	assert.True(t, panes[0].Focused)
	assert.False(t, panes[1].Focused)
}

func TestPaneNode_PaneInfo(t *testing.T) {
	styles := DefaultStyles()
	pane1 := NewPane("pane-1", PaneTypeChat, "Pane 1", 100, 50, styles)
	pane1.Focus()

	node := NewPaneNode(pane1)
	node.SetSize(0, 0, 100, 50)

	info := node.PaneInfo()

	assert.Contains(t, info, "pane-1")
	assert.Contains(t, info, "Pane 1")
	assert.Contains(t, info, "*") // Focus indicator
}
