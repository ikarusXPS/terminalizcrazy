package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Tab represents a single tab
type Tab struct {
	ID       string
	Title    string
	Icon     string
	Modified bool
	Pinned   bool
	Content  interface{} // Can hold tab-specific data
}

// TabBar manages multiple tabs
type TabBar struct {
	tabs      []*Tab
	activeIdx int
	styles    *Styles
	width     int
	maxTabs   int
	tabCount  int
}

// NewTabBar creates a new tab bar
func NewTabBar(styles *Styles, maxTabs int) *TabBar {
	if maxTabs <= 0 {
		maxTabs = 10
	}

	tb := &TabBar{
		tabs:     make([]*Tab, 0),
		styles:   styles,
		maxTabs:  maxTabs,
		tabCount: 0,
	}

	// Create initial tab
	tb.AddTab("Main", "")

	return tb
}

// AddTab adds a new tab
func (tb *TabBar) AddTab(title, icon string) *Tab {
	if len(tb.tabs) >= tb.maxTabs {
		return nil
	}

	tb.tabCount++
	tab := &Tab{
		ID:    fmt.Sprintf("tab-%d", tb.tabCount),
		Title: title,
		Icon:  icon,
	}

	tb.tabs = append(tb.tabs, tab)
	tb.activeIdx = len(tb.tabs) - 1

	return tab
}

// AddTabAtIndex adds a new tab at a specific index
func (tb *TabBar) AddTabAtIndex(index int, title, icon string) *Tab {
	if len(tb.tabs) >= tb.maxTabs {
		return nil
	}

	tb.tabCount++
	tab := &Tab{
		ID:    fmt.Sprintf("tab-%d", tb.tabCount),
		Title: title,
		Icon:  icon,
	}

	if index < 0 {
		index = 0
	}
	if index > len(tb.tabs) {
		index = len(tb.tabs)
	}

	// Insert at index
	tb.tabs = append(tb.tabs[:index], append([]*Tab{tab}, tb.tabs[index:]...)...)
	tb.activeIdx = index

	return tab
}

// CloseTab closes a tab by index
func (tb *TabBar) CloseTab(index int) bool {
	if index < 0 || index >= len(tb.tabs) {
		return false
	}

	// Don't close if it's the last tab
	if len(tb.tabs) <= 1 {
		return false
	}

	// Don't close pinned tabs
	if tb.tabs[index].Pinned {
		return false
	}

	// Remove the tab
	tb.tabs = append(tb.tabs[:index], tb.tabs[index+1:]...)

	// Adjust active index
	if tb.activeIdx >= len(tb.tabs) {
		tb.activeIdx = len(tb.tabs) - 1
	} else if tb.activeIdx > index {
		tb.activeIdx--
	}

	return true
}

// CloseActiveTab closes the currently active tab
func (tb *TabBar) CloseActiveTab() bool {
	return tb.CloseTab(tb.activeIdx)
}

// CloseTabByID closes a tab by ID
func (tb *TabBar) CloseTabByID(id string) bool {
	for i, tab := range tb.tabs {
		if tab.ID == id {
			return tb.CloseTab(i)
		}
	}
	return false
}

// SetActiveTab sets the active tab by index
func (tb *TabBar) SetActiveTab(index int) bool {
	if index < 0 || index >= len(tb.tabs) {
		return false
	}
	tb.activeIdx = index
	return true
}

// SetActiveTabByID sets the active tab by ID
func (tb *TabBar) SetActiveTabByID(id string) bool {
	for i, tab := range tb.tabs {
		if tab.ID == id {
			tb.activeIdx = i
			return true
		}
	}
	return false
}

// NextTab switches to the next tab
func (tb *TabBar) NextTab() {
	if len(tb.tabs) > 1 {
		tb.activeIdx = (tb.activeIdx + 1) % len(tb.tabs)
	}
}

// PrevTab switches to the previous tab
func (tb *TabBar) PrevTab() {
	if len(tb.tabs) > 1 {
		tb.activeIdx--
		if tb.activeIdx < 0 {
			tb.activeIdx = len(tb.tabs) - 1
		}
	}
}

// GetActiveTab returns the currently active tab
func (tb *TabBar) GetActiveTab() *Tab {
	if tb.activeIdx >= 0 && tb.activeIdx < len(tb.tabs) {
		return tb.tabs[tb.activeIdx]
	}
	return nil
}

// GetActiveIndex returns the active tab index
func (tb *TabBar) GetActiveIndex() int {
	return tb.activeIdx
}

// GetTab returns a tab by index
func (tb *TabBar) GetTab(index int) *Tab {
	if index >= 0 && index < len(tb.tabs) {
		return tb.tabs[index]
	}
	return nil
}

// GetTabByID returns a tab by ID
func (tb *TabBar) GetTabByID(id string) *Tab {
	for _, tab := range tb.tabs {
		if tab.ID == id {
			return tab
		}
	}
	return nil
}

// GetTabs returns all tabs
func (tb *TabBar) GetTabs() []*Tab {
	return tb.tabs
}

// GetTabCount returns the number of tabs
func (tb *TabBar) GetTabCount() int {
	return len(tb.tabs)
}

// SetTabTitle sets the title of a tab
func (tb *TabBar) SetTabTitle(index int, title string) bool {
	if index >= 0 && index < len(tb.tabs) {
		tb.tabs[index].Title = title
		return true
	}
	return false
}

// SetActiveTabTitle sets the title of the active tab
func (tb *TabBar) SetActiveTabTitle(title string) bool {
	return tb.SetTabTitle(tb.activeIdx, title)
}

// MarkModified marks a tab as modified
func (tb *TabBar) MarkModified(index int, modified bool) bool {
	if index >= 0 && index < len(tb.tabs) {
		tb.tabs[index].Modified = modified
		return true
	}
	return false
}

// MarkActiveModified marks the active tab as modified
func (tb *TabBar) MarkActiveModified(modified bool) bool {
	return tb.MarkModified(tb.activeIdx, modified)
}

// PinTab pins or unpins a tab
func (tb *TabBar) PinTab(index int, pinned bool) bool {
	if index >= 0 && index < len(tb.tabs) {
		tb.tabs[index].Pinned = pinned
		return true
	}
	return false
}

// MoveTab moves a tab from one position to another
func (tb *TabBar) MoveTab(from, to int) bool {
	if from < 0 || from >= len(tb.tabs) || to < 0 || to >= len(tb.tabs) {
		return false
	}

	if from == to {
		return true
	}

	tab := tb.tabs[from]

	// Remove from old position
	tb.tabs = append(tb.tabs[:from], tb.tabs[from+1:]...)

	// Insert at new position
	if to > from {
		to-- // Adjust for removal
	}
	tb.tabs = append(tb.tabs[:to], append([]*Tab{tab}, tb.tabs[to:]...)...)

	// Update active index if needed
	if tb.activeIdx == from {
		tb.activeIdx = to
	} else if from < tb.activeIdx && to >= tb.activeIdx {
		tb.activeIdx--
	} else if from > tb.activeIdx && to <= tb.activeIdx {
		tb.activeIdx++
	}

	return true
}

// SetWidth sets the tab bar width
func (tb *TabBar) SetWidth(width int) {
	tb.width = width
}

// View renders the tab bar
func (tb *TabBar) View() string {
	if len(tb.tabs) == 0 {
		return ""
	}

	var tabs []string
	maxTabWidth := 20

	for i, tab := range tb.tabs {
		title := tab.Title
		if tab.Icon != "" {
			title = tab.Icon + " " + title
		}

		// Add modified indicator
		if tab.Modified {
			title = "● " + title
		}

		// Add pinned indicator
		if tab.Pinned {
			title = "📌 " + title
		}

		// Truncate if too long
		if len(title) > maxTabWidth {
			title = title[:maxTabWidth-1] + "…"
		}

		var style lipgloss.Style
		if i == tb.activeIdx {
			style = tb.styles.TabActive
		} else {
			style = tb.styles.TabInactive
		}

		tabs = append(tabs, style.Render(title))
	}

	// Join tabs horizontally
	tabContent := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

	// Add new tab button
	newTabBtn := tb.styles.TabInactive.Render("+")
	tabContent = lipgloss.JoinHorizontal(lipgloss.Top, tabContent, newTabBtn)

	// Apply tab bar style
	return tb.styles.TabBar.Width(tb.width).Render(tabContent)
}

// ViewCompact renders a compact tab bar (numbers only)
func (tb *TabBar) ViewCompact() string {
	if len(tb.tabs) == 0 {
		return ""
	}

	var tabs []string

	for i := range tb.tabs {
		num := fmt.Sprintf("%d", i+1)
		var style lipgloss.Style
		if i == tb.activeIdx {
			style = tb.styles.TabActive
		} else {
			style = tb.styles.TabInactive
		}
		tabs = append(tabs, style.Render(num))
	}

	return strings.Join(tabs, " ")
}

// HandleKey processes keyboard input for tab navigation
func (tb *TabBar) HandleKey(key string) bool {
	switch key {
	case "ctrl+t":
		tb.AddTab("New Tab", "")
		return true

	case "ctrl+w":
		tb.CloseActiveTab()
		return true

	case "ctrl+tab", "ctrl+pgdown":
		tb.NextTab()
		return true

	case "ctrl+shift+tab", "ctrl+pgup":
		tb.PrevTab()
		return true

	case "alt+1":
		tb.SetActiveTab(0)
		return true
	case "alt+2":
		tb.SetActiveTab(1)
		return true
	case "alt+3":
		tb.SetActiveTab(2)
		return true
	case "alt+4":
		tb.SetActiveTab(3)
		return true
	case "alt+5":
		tb.SetActiveTab(4)
		return true
	case "alt+6":
		tb.SetActiveTab(5)
		return true
	case "alt+7":
		tb.SetActiveTab(6)
		return true
	case "alt+8":
		tb.SetActiveTab(7)
		return true
	case "alt+9":
		// Jump to last tab
		if len(tb.tabs) > 0 {
			tb.SetActiveTab(len(tb.tabs) - 1)
		}
		return true
	}

	return false
}

// DuplicateActiveTab duplicates the active tab
func (tb *TabBar) DuplicateActiveTab() *Tab {
	if tb.activeIdx < 0 || tb.activeIdx >= len(tb.tabs) {
		return nil
	}

	active := tb.tabs[tb.activeIdx]
	newTab := tb.AddTab(active.Title+" (copy)", active.Icon)
	if newTab != nil {
		newTab.Content = active.Content
	}

	return newTab
}

// GetModifiedTabs returns all modified tabs
func (tb *TabBar) GetModifiedTabs() []*Tab {
	var modified []*Tab
	for _, tab := range tb.tabs {
		if tab.Modified {
			modified = append(modified, tab)
		}
	}
	return modified
}

// CloseAllUnpinned closes all unpinned tabs
func (tb *TabBar) CloseAllUnpinned() int {
	closed := 0
	for i := len(tb.tabs) - 1; i >= 0; i-- {
		if !tb.tabs[i].Pinned && len(tb.tabs) > 1 {
			if tb.CloseTab(i) {
				closed++
			}
		}
	}
	return closed
}

// CloseAllToRight closes all tabs to the right of the active tab
func (tb *TabBar) CloseAllToRight() int {
	closed := 0
	for i := len(tb.tabs) - 1; i > tb.activeIdx; i-- {
		if tb.CloseTab(i) {
			closed++
		}
	}
	return closed
}

// CloseAllToLeft closes all tabs to the left of the active tab
func (tb *TabBar) CloseAllToLeft() int {
	closed := 0
	for i := tb.activeIdx - 1; i >= 0; i-- {
		if !tb.tabs[i].Pinned && len(tb.tabs) > 1 {
			if tb.CloseTab(i) {
				closed++
			}
		}
	}
	return closed
}
