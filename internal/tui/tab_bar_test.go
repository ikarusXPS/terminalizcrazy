package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTabBar(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	assert.NotNil(t, tb)
	assert.Equal(t, 1, tb.GetTabCount()) // Initial tab
	assert.Equal(t, 0, tb.GetActiveIndex())
}

func TestNewTabBar_DefaultMaxTabs(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 0)

	// Should default to 10
	for i := 0; i < 9; i++ {
		tb.AddTab("Tab", "")
	}
	assert.Equal(t, 10, tb.GetTabCount())

	// Should not add more
	result := tb.AddTab("Tab", "")
	assert.Nil(t, result)
}

func TestTabBar_AddTab(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	tab := tb.AddTab("New Tab", "📁")

	assert.NotNil(t, tab)
	assert.Equal(t, "New Tab", tab.Title)
	assert.Equal(t, "📁", tab.Icon)
	assert.Equal(t, 2, tb.GetTabCount())
	assert.Equal(t, 1, tb.GetActiveIndex())
}

func TestTabBar_AddTab_MaxReached(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 2)

	tb.AddTab("Tab 2", "")

	result := tb.AddTab("Tab 3", "")

	assert.Nil(t, result)
	assert.Equal(t, 2, tb.GetTabCount())
}

func TestTabBar_AddTabAtIndex(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	tb.AddTab("Tab 2", "")
	tab := tb.AddTabAtIndex(1, "Middle Tab", "")

	assert.NotNil(t, tab)
	assert.Equal(t, 1, tb.GetActiveIndex())
	assert.Equal(t, "Middle Tab", tb.GetTab(1).Title)
}

func TestTabBar_AddTabAtIndex_BoundsCheck(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	tab := tb.AddTabAtIndex(-1, "Tab", "")
	assert.NotNil(t, tab)
	assert.Equal(t, 0, tb.GetActiveIndex())

	tab = tb.AddTabAtIndex(100, "Tab", "")
	assert.NotNil(t, tab)
	assert.Equal(t, tb.GetTabCount()-1, tb.GetActiveIndex())
}

func TestTabBar_CloseTab(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	tb.AddTab("Tab 2", "")
	assert.Equal(t, 2, tb.GetTabCount())

	result := tb.CloseTab(1)

	assert.True(t, result)
	assert.Equal(t, 1, tb.GetTabCount())
}

func TestTabBar_CloseTab_LastTab(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	result := tb.CloseTab(0)

	assert.False(t, result)
	assert.Equal(t, 1, tb.GetTabCount())
}

func TestTabBar_CloseTab_PinnedTab(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	tb.AddTab("Tab 2", "")
	tb.PinTab(0, true)

	result := tb.CloseTab(0)

	assert.False(t, result)
	assert.Equal(t, 2, tb.GetTabCount())
}

func TestTabBar_CloseTab_InvalidIndex(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	result := tb.CloseTab(-1)
	assert.False(t, result)

	result = tb.CloseTab(100)
	assert.False(t, result)
}

func TestTabBar_CloseActiveTab(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	tb.AddTab("Tab 2", "")
	tb.SetActiveTab(1)

	result := tb.CloseActiveTab()

	assert.True(t, result)
	assert.Equal(t, 1, tb.GetTabCount())
}

func TestTabBar_CloseTabByID(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	tab := tb.AddTab("Tab 2", "")

	result := tb.CloseTabByID(tab.ID)

	assert.True(t, result)
	assert.Equal(t, 1, tb.GetTabCount())
}

func TestTabBar_CloseTabByID_NotFound(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	result := tb.CloseTabByID("nonexistent")

	assert.False(t, result)
}

func TestTabBar_SetActiveTab(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	tb.AddTab("Tab 2", "")
	tb.AddTab("Tab 3", "")

	result := tb.SetActiveTab(0)

	assert.True(t, result)
	assert.Equal(t, 0, tb.GetActiveIndex())
}

func TestTabBar_SetActiveTab_InvalidIndex(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	result := tb.SetActiveTab(-1)
	assert.False(t, result)

	result = tb.SetActiveTab(100)
	assert.False(t, result)
}

func TestTabBar_SetActiveTabByID(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	tab := tb.AddTab("Tab 2", "")
	tb.SetActiveTab(0)

	result := tb.SetActiveTabByID(tab.ID)

	assert.True(t, result)
	assert.Equal(t, 1, tb.GetActiveIndex())
}

func TestTabBar_SetActiveTabByID_NotFound(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	result := tb.SetActiveTabByID("nonexistent")

	assert.False(t, result)
}

func TestTabBar_NextPrevTab(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	tb.AddTab("Tab 2", "")
	tb.AddTab("Tab 3", "")
	tb.SetActiveTab(0)

	tb.NextTab()
	assert.Equal(t, 1, tb.GetActiveIndex())

	tb.NextTab()
	assert.Equal(t, 2, tb.GetActiveIndex())

	tb.NextTab()
	assert.Equal(t, 0, tb.GetActiveIndex()) // Wrap

	tb.PrevTab()
	assert.Equal(t, 2, tb.GetActiveIndex()) // Wrap back

	tb.PrevTab()
	assert.Equal(t, 1, tb.GetActiveIndex())
}

func TestTabBar_GetActiveTab(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	tb.AddTab("Tab 2", "")

	active := tb.GetActiveTab()

	assert.NotNil(t, active)
	assert.Equal(t, "Tab 2", active.Title)
}

func TestTabBar_GetTab(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	tab := tb.GetTab(0)
	assert.NotNil(t, tab)

	tab = tb.GetTab(-1)
	assert.Nil(t, tab)

	tab = tb.GetTab(100)
	assert.Nil(t, tab)
}

func TestTabBar_GetTabByID(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	addedTab := tb.AddTab("Test", "")

	tab := tb.GetTabByID(addedTab.ID)
	assert.NotNil(t, tab)
	assert.Equal(t, addedTab.ID, tab.ID)

	tab = tb.GetTabByID("nonexistent")
	assert.Nil(t, tab)
}

func TestTabBar_GetTabs(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	tb.AddTab("Tab 2", "")

	tabs := tb.GetTabs()

	assert.Len(t, tabs, 2)
}

func TestTabBar_SetTabTitle(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	result := tb.SetTabTitle(0, "New Title")

	assert.True(t, result)
	assert.Equal(t, "New Title", tb.GetTab(0).Title)
}

func TestTabBar_SetTabTitle_InvalidIndex(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	result := tb.SetTabTitle(-1, "Title")
	assert.False(t, result)

	result = tb.SetTabTitle(100, "Title")
	assert.False(t, result)
}

func TestTabBar_SetActiveTabTitle(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	result := tb.SetActiveTabTitle("New Title")

	assert.True(t, result)
	assert.Equal(t, "New Title", tb.GetActiveTab().Title)
}

func TestTabBar_MarkModified(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	result := tb.MarkModified(0, true)

	assert.True(t, result)
	assert.True(t, tb.GetTab(0).Modified)

	tb.MarkModified(0, false)
	assert.False(t, tb.GetTab(0).Modified)
}

func TestTabBar_MarkModified_InvalidIndex(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	result := tb.MarkModified(-1, true)
	assert.False(t, result)
}

func TestTabBar_MarkActiveModified(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	result := tb.MarkActiveModified(true)

	assert.True(t, result)
	assert.True(t, tb.GetActiveTab().Modified)
}

func TestTabBar_PinTab(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	result := tb.PinTab(0, true)

	assert.True(t, result)
	assert.True(t, tb.GetTab(0).Pinned)

	tb.PinTab(0, false)
	assert.False(t, tb.GetTab(0).Pinned)
}

func TestTabBar_PinTab_InvalidIndex(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	result := tb.PinTab(-1, true)
	assert.False(t, result)
}

func TestTabBar_MoveTab(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	tb.AddTab("Tab 2", "")
	tb.AddTab("Tab 3", "")
	tb.SetActiveTab(0)

	result := tb.MoveTab(0, 2)

	assert.True(t, result)
	// Active index should follow the moved tab
}

func TestTabBar_MoveTab_SamePosition(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	tb.AddTab("Tab 2", "")

	result := tb.MoveTab(1, 1)

	assert.True(t, result)
}

func TestTabBar_MoveTab_InvalidIndex(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	result := tb.MoveTab(-1, 0)
	assert.False(t, result)

	result = tb.MoveTab(0, 100)
	assert.False(t, result)
}

func TestTabBar_SetWidth(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	tb.SetWidth(200)

	assert.Equal(t, 200, tb.width)
}

func TestTabBar_View(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	tb.SetWidth(200)
	view := tb.View()

	assert.NotEmpty(t, view)
	assert.Contains(t, view, "Main")
	assert.Contains(t, view, "+") // New tab button
}

func TestTabBar_View_Empty(t *testing.T) {
	styles := DefaultStyles()
	tb := &TabBar{styles: styles, tabs: []*Tab{}}

	view := tb.View()

	assert.Empty(t, view)
}

func TestTabBar_ViewCompact(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	tb.AddTab("Tab 2", "")

	view := tb.ViewCompact()

	assert.Contains(t, view, "1")
	assert.Contains(t, view, "2")
}

func TestTabBar_HandleKey(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	// Add new tab
	result := tb.HandleKey("ctrl+t")
	assert.True(t, result)
	assert.Equal(t, 2, tb.GetTabCount())

	// Close tab
	result = tb.HandleKey("ctrl+w")
	assert.True(t, result)
	assert.Equal(t, 1, tb.GetTabCount())

	// Next/Prev tab
	tb.AddTab("Tab 2", "")
	tb.SetActiveTab(0)

	result = tb.HandleKey("ctrl+tab")
	assert.True(t, result)
	assert.Equal(t, 1, tb.GetActiveIndex())

	result = tb.HandleKey("ctrl+shift+tab")
	assert.True(t, result)
	assert.Equal(t, 0, tb.GetActiveIndex())

	// Alt+number
	result = tb.HandleKey("alt+2")
	assert.True(t, result)
	assert.Equal(t, 1, tb.GetActiveIndex())

	result = tb.HandleKey("alt+1")
	assert.True(t, result)
	assert.Equal(t, 0, tb.GetActiveIndex())

	result = tb.HandleKey("alt+9")
	assert.True(t, result)
	assert.Equal(t, 1, tb.GetActiveIndex()) // Last tab
}

func TestTabBar_HandleKey_Unhandled(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	result := tb.HandleKey("unhandled")

	assert.False(t, result)
}

func TestTabBar_DuplicateActiveTab(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	tb.GetActiveTab().Title = "Original"

	dup := tb.DuplicateActiveTab()

	assert.NotNil(t, dup)
	assert.Contains(t, dup.Title, "copy")
	assert.Equal(t, 2, tb.GetTabCount())
}

func TestTabBar_GetModifiedTabs(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	tb.AddTab("Tab 2", "")
	tb.MarkModified(0, true)

	modified := tb.GetModifiedTabs()

	assert.Len(t, modified, 1)
	assert.Equal(t, "Main", modified[0].Title)
}

func TestTabBar_CloseAllUnpinned(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	tb.AddTab("Tab 2", "")
	tb.AddTab("Tab 3", "")
	tb.PinTab(0, true)

	closed := tb.CloseAllUnpinned()

	assert.Equal(t, 2, closed)
	assert.Equal(t, 1, tb.GetTabCount())
	assert.True(t, tb.GetTab(0).Pinned)
}

func TestTabBar_CloseAllToRight(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	tb.AddTab("Tab 2", "")
	tb.AddTab("Tab 3", "")
	tb.AddTab("Tab 4", "")
	tb.SetActiveTab(1)

	closed := tb.CloseAllToRight()

	assert.Equal(t, 2, closed)
	assert.Equal(t, 2, tb.GetTabCount())
}

func TestTabBar_CloseAllToLeft(t *testing.T) {
	styles := DefaultStyles()
	tb := NewTabBar(styles, 10)

	tb.AddTab("Tab 2", "")
	tb.AddTab("Tab 3", "")
	tb.AddTab("Tab 4", "")
	tb.SetActiveTab(2)

	closed := tb.CloseAllToLeft()

	assert.Equal(t, 2, closed)
	assert.Equal(t, 2, tb.GetTabCount())
}

func TestTab_Properties(t *testing.T) {
	tab := &Tab{
		ID:       "tab-1",
		Title:    "Test Tab",
		Icon:     "📁",
		Modified: true,
		Pinned:   true,
		Content:  "some content",
	}

	assert.Equal(t, "tab-1", tab.ID)
	assert.Equal(t, "Test Tab", tab.Title)
	assert.Equal(t, "📁", tab.Icon)
	assert.True(t, tab.Modified)
	assert.True(t, tab.Pinned)
	assert.Equal(t, "some content", tab.Content)
}
