package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

// PaneType defines the type of content in a pane
type PaneType string

const (
	PaneTypeChat     PaneType = "chat"
	PaneTypeTerminal PaneType = "terminal"
	PaneTypePlan     PaneType = "plan"
	PaneTypeOutput   PaneType = "output"
)

// Pane represents a single content pane
type Pane struct {
	ID       string
	Type     PaneType
	Title    string
	Focused  bool
	X        int
	Y        int
	Width    int
	Height   int
	Viewport viewport.Model
	Content  string
	styles   *Styles
}

// NewPane creates a new pane
func NewPane(id string, paneType PaneType, title string, width, height int, styles *Styles) *Pane {
	vp := viewport.New(width-2, height-2) // Account for borders

	return &Pane{
		ID:       id,
		Type:     paneType,
		Title:    title,
		Width:    width,
		Height:   height,
		Viewport: vp,
		styles:   styles,
	}
}

// SetSize updates the pane dimensions
func (p *Pane) SetSize(width, height int) {
	p.Width = width
	p.Height = height
	p.Viewport.Width = width - 2
	p.Viewport.Height = height - 2
}

// SetContent updates the pane content
func (p *Pane) SetContent(content string) {
	p.Content = content
	p.Viewport.SetContent(content)
}

// AppendContent appends content to the pane
func (p *Pane) AppendContent(content string) {
	p.Content += content
	p.Viewport.SetContent(p.Content)
	p.Viewport.GotoBottom()
}

// Focus focuses the pane
func (p *Pane) Focus() {
	p.Focused = true
}

// Blur removes focus from the pane
func (p *Pane) Blur() {
	p.Focused = false
}

// View renders the pane
func (p *Pane) View() string {
	// Choose border style based on focus
	borderStyle := p.styles.PaneUnfocused
	if p.Focused {
		borderStyle = p.styles.PaneFocused
	}

	// Build the title bar
	titleText := p.styles.PaneTitle.Render(p.Title)
	if p.Focused {
		titleText = "● " + titleText
	}

	// Create content area with viewport
	content := p.Viewport.View()

	// Combine title and content
	return borderStyle.
		Width(p.Width).
		Height(p.Height).
		Render(fmt.Sprintf("%s\n%s", titleText, content))
}

// SplitDirection defines the direction of a split
type SplitDirection int

const (
	SplitHorizontal SplitDirection = iota
	SplitVertical
)

// PaneNode represents a node in the pane layout tree
type PaneNode struct {
	IsLeaf    bool
	Pane      *Pane
	Split     SplitDirection
	Ratio     float64 // 0.0 to 1.0
	Children  [2]*PaneNode
	X, Y      int
	Width     int
	Height    int
}

// NewPaneNode creates a leaf node with a pane
func NewPaneNode(pane *Pane) *PaneNode {
	return &PaneNode{
		IsLeaf: true,
		Pane:   pane,
		Ratio:  0.5,
	}
}

// Split splits the node into two children
func (n *PaneNode) Split(direction SplitDirection, newPane *Pane) {
	if !n.IsLeaf {
		return // Can only split leaf nodes
	}

	n.IsLeaf = false
	n.Split = direction
	n.Ratio = 0.5

	// Move current pane to first child
	n.Children[0] = NewPaneNode(n.Pane)
	n.Pane = nil

	// Create new pane as second child
	n.Children[1] = NewPaneNode(newPane)

	// Recalculate sizes
	n.recalculateSizes()
}

// recalculateSizes updates child sizes based on ratio
func (n *PaneNode) recalculateSizes() {
	if n.IsLeaf {
		if n.Pane != nil {
			n.Pane.SetSize(n.Width, n.Height)
			n.Pane.X = n.X
			n.Pane.Y = n.Y
		}
		return
	}

	switch n.Split {
	case SplitHorizontal:
		// Split top/bottom
		topHeight := int(float64(n.Height) * n.Ratio)
		bottomHeight := n.Height - topHeight

		n.Children[0].X = n.X
		n.Children[0].Y = n.Y
		n.Children[0].Width = n.Width
		n.Children[0].Height = topHeight

		n.Children[1].X = n.X
		n.Children[1].Y = n.Y + topHeight
		n.Children[1].Width = n.Width
		n.Children[1].Height = bottomHeight

	case SplitVertical:
		// Split left/right
		leftWidth := int(float64(n.Width) * n.Ratio)
		rightWidth := n.Width - leftWidth

		n.Children[0].X = n.X
		n.Children[0].Y = n.Y
		n.Children[0].Width = leftWidth
		n.Children[0].Height = n.Height

		n.Children[1].X = n.X + leftWidth
		n.Children[1].Y = n.Y
		n.Children[1].Width = rightWidth
		n.Children[1].Height = n.Height
	}

	// Recurse to children
	n.Children[0].recalculateSizes()
	n.Children[1].recalculateSizes()
}

// SetSize sets the node size
func (n *PaneNode) SetSize(x, y, width, height int) {
	n.X = x
	n.Y = y
	n.Width = width
	n.Height = height
	n.recalculateSizes()
}

// GetPanes returns all panes in the tree
func (n *PaneNode) GetPanes() []*Pane {
	if n.IsLeaf {
		if n.Pane != nil {
			return []*Pane{n.Pane}
		}
		return nil
	}

	var panes []*Pane
	if n.Children[0] != nil {
		panes = append(panes, n.Children[0].GetPanes()...)
	}
	if n.Children[1] != nil {
		panes = append(panes, n.Children[1].GetPanes()...)
	}
	return panes
}

// FindPane finds a pane by ID
func (n *PaneNode) FindPane(id string) *Pane {
	if n.IsLeaf {
		if n.Pane != nil && n.Pane.ID == id {
			return n.Pane
		}
		return nil
	}

	if n.Children[0] != nil {
		if pane := n.Children[0].FindPane(id); pane != nil {
			return pane
		}
	}
	if n.Children[1] != nil {
		if pane := n.Children[1].FindPane(id); pane != nil {
			return pane
		}
	}
	return nil
}

// RemovePane removes a pane from the tree
func (n *PaneNode) RemovePane(id string) bool {
	if n.IsLeaf {
		return false
	}

	// Check if either child contains the pane to remove
	for i, child := range n.Children {
		if child != nil && child.IsLeaf && child.Pane != nil && child.Pane.ID == id {
			// Replace this node with the other child
			other := n.Children[1-i]
			*n = *other
			n.recalculateSizes()
			return true
		}
	}

	// Recurse into children
	for _, child := range n.Children {
		if child != nil && child.RemovePane(id) {
			return true
		}
	}

	return false
}

// View renders the entire pane tree
func (n *PaneNode) View() string {
	if n.IsLeaf {
		if n.Pane != nil {
			return n.Pane.View()
		}
		return ""
	}

	child0View := n.Children[0].View()
	child1View := n.Children[1].View()

	switch n.Split {
	case SplitHorizontal:
		return lipgloss.JoinVertical(lipgloss.Left, child0View, child1View)
	case SplitVertical:
		return lipgloss.JoinHorizontal(lipgloss.Top, child0View, child1View)
	}

	return ""
}

// GetFocusedPane returns the currently focused pane
func (n *PaneNode) GetFocusedPane() *Pane {
	panes := n.GetPanes()
	for _, pane := range panes {
		if pane.Focused {
			return pane
		}
	}
	return nil
}

// FocusNext focuses the next pane in order
func (n *PaneNode) FocusNext() {
	panes := n.GetPanes()
	if len(panes) == 0 {
		return
	}

	// Find current focused index
	currentIdx := -1
	for i, pane := range panes {
		if pane.Focused {
			currentIdx = i
			pane.Blur()
			break
		}
	}

	// Focus next pane (or first if none focused)
	nextIdx := (currentIdx + 1) % len(panes)
	panes[nextIdx].Focus()
}

// FocusPrevious focuses the previous pane
func (n *PaneNode) FocusPrevious() {
	panes := n.GetPanes()
	if len(panes) == 0 {
		return
	}

	// Find current focused index
	currentIdx := -1
	for i, pane := range panes {
		if pane.Focused {
			currentIdx = i
			pane.Blur()
			break
		}
	}

	// Focus previous pane
	prevIdx := currentIdx - 1
	if prevIdx < 0 {
		prevIdx = len(panes) - 1
	}
	panes[prevIdx].Focus()
}

// FocusDirection focuses a pane in the given direction
func (n *PaneNode) FocusDirection(direction string) {
	panes := n.GetPanes()
	focused := n.GetFocusedPane()
	if focused == nil || len(panes) <= 1 {
		return
	}

	var bestPane *Pane
	bestDistance := -1

	for _, pane := range panes {
		if pane.ID == focused.ID {
			continue
		}

		dx := pane.X - focused.X
		dy := pane.Y - focused.Y

		var valid bool
		var distance int

		switch direction {
		case "left":
			valid = dx < 0
			distance = -dx
		case "right":
			valid = dx > 0
			distance = dx
		case "up":
			valid = dy < 0
			distance = -dy
		case "down":
			valid = dy > 0
			distance = dy
		}

		if valid && (bestDistance < 0 || distance < bestDistance) {
			bestPane = pane
			bestDistance = distance
		}
	}

	if bestPane != nil {
		focused.Blur()
		bestPane.Focus()
	}
}

// PaneInfo returns info about panes for debugging
func (n *PaneNode) PaneInfo() string {
	var sb strings.Builder
	panes := n.GetPanes()

	for _, pane := range panes {
		focus := " "
		if pane.Focused {
			focus = "*"
		}
		sb.WriteString(fmt.Sprintf("%s[%s] %s (%dx%d at %d,%d)\n",
			focus, pane.ID, pane.Title, pane.Width, pane.Height, pane.X, pane.Y))
	}

	return sb.String()
}
