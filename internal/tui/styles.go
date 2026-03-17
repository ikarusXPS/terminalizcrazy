package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/terminalizcrazy/terminalizcrazy/internal/theme"
)

// Color constants for theming
const (
	ColorPrimary     = "#7D56F4"
	ColorSecondary   = "#04B575"
	ColorWarning     = "#FFAA00"
	ColorError       = "#FF6B6B"
	ColorSuccess     = "#04B575"
	ColorMuted       = "#888888"
	ColorDark        = "#2D2D2D"
	ColorLight       = "#E0E0E0"
	ColorLighter     = "#AAAAAA"
	ColorCollabUser  = "#4ECDC4"
	ColorShareCode   = "#FFEAA7"
)

// Styles contains all the application styles
type Styles struct {
	Title              lipgloss.Style
	Version            lipgloss.Style
	Input              lipgloss.Style
	StatusConnected    lipgloss.Style
	StatusDisconnected lipgloss.Style
	Help               lipgloss.Style
	UserMsg            lipgloss.Style
	AIMsg              lipgloss.Style
	SystemMsg          lipgloss.Style
	Command            lipgloss.Style
	Output             lipgloss.Style
	Error              lipgloss.Style
	Success            lipgloss.Style
	Spinner            lipgloss.Style
	History            lipgloss.Style
	CopyNotice         lipgloss.Style
	SessionItem        lipgloss.Style
	SessionSelected    lipgloss.Style
	SessionHeader      lipgloss.Style
	CollabUser         lipgloss.Style
	ShareCode          lipgloss.Style

	// Tab styles
	Tab           lipgloss.Style
	TabActive     lipgloss.Style
	TabInactive   lipgloss.Style
	TabBar        lipgloss.Style

	// Pane styles
	Pane          lipgloss.Style
	PaneFocused   lipgloss.Style
	PaneUnfocused lipgloss.Style
	PaneBorder    lipgloss.Style
	PaneTitle     lipgloss.Style

	// Agent styles
	AgentStatus   lipgloss.Style
	AgentPlan     lipgloss.Style
	AgentTask     lipgloss.Style
	AgentProgress lipgloss.Style
}

// DefaultStyles returns the default application styles
func DefaultStyles() *Styles {
	return &Styles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(ColorPrimary)),

		Version: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMuted)),

		Input: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPrimary)).
			Padding(0, 1),

		StatusConnected: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSuccess)),

		StatusDisconnected: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorWarning)),

		Help: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")),

		UserMsg: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorPrimary)).
			Bold(true),

		AIMsg: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSecondary)),

		SystemMsg: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMuted)).
			Italic(true),

		Command: lipgloss.NewStyle().
			Background(lipgloss.Color(ColorDark)).
			Foreground(lipgloss.Color(ColorLight)).
			Padding(0, 1),

		Output: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorLighter)),

		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorError)),

		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSuccess)),

		Spinner: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorPrimary)),

		History: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")),

		CopyNotice: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSuccess)).
			Bold(true),

		SessionItem: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC")).
			PaddingLeft(2),

		SessionSelected: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorPrimary)).
			Bold(true).
			PaddingLeft(2),

		SessionHeader: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorPrimary)).
			Bold(true).
			MarginBottom(1),

		CollabUser: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorCollabUser)),

		ShareCode: lipgloss.NewStyle().
			Background(lipgloss.Color(ColorDark)).
			Foreground(lipgloss.Color(ColorShareCode)).
			Bold(true).
			Padding(0, 1),

		// Tab styles
		Tab: lipgloss.NewStyle().
			Padding(0, 2),

		TabActive: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorPrimary)).
			Background(lipgloss.Color(ColorDark)).
			Bold(true).
			Padding(0, 2),

		TabInactive: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMuted)).
			Padding(0, 2),

		TabBar: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(lipgloss.Color(ColorMuted)),

		// Pane styles
		Pane: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorMuted)),

		PaneFocused: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPrimary)),

		PaneUnfocused: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorMuted)),

		PaneBorder: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorMuted)),

		PaneTitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorPrimary)).
			Bold(true).
			Padding(0, 1),

		// Agent styles
		AgentStatus: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSecondary)).
			Bold(true),

		AgentPlan: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorSecondary)).
			Padding(1),

		AgentTask: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorLight)),

		AgentProgress: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSecondary)),
	}
}

// Theme represents a color theme
type Theme struct {
	Name    string
	Primary string
	Success string
	Warning string
	Error   string
	Muted   string
}

// GetTheme returns a theme by name
func GetTheme(name string) *Theme {
	themes := map[string]*Theme{
		"default": {
			Name:    "default",
			Primary: ColorPrimary,
			Success: ColorSuccess,
			Warning: ColorWarning,
			Error:   ColorError,
			Muted:   ColorMuted,
		},
		"dark": {
			Name:    "dark",
			Primary: "#BB86FC",
			Success: "#03DAC6",
			Warning: "#FFAB00",
			Error:   "#CF6679",
			Muted:   "#606060",
		},
		"light": {
			Name:    "light",
			Primary: "#6200EE",
			Success: "#00C853",
			Warning: "#FF6D00",
			Error:   "#B00020",
			Muted:   "#9E9E9E",
		},
	}

	if theme, ok := themes[name]; ok {
		return theme
	}
	return themes["default"]
}

// ApplyTheme applies a theme to the styles
func (s *Styles) ApplyTheme(t *Theme) {
	s.Title = s.Title.Foreground(lipgloss.Color(t.Primary))
	s.Input = s.Input.BorderForeground(lipgloss.Color(t.Primary))
	s.StatusConnected = s.StatusConnected.Foreground(lipgloss.Color(t.Success))
	s.StatusDisconnected = s.StatusDisconnected.Foreground(lipgloss.Color(t.Warning))
	s.UserMsg = s.UserMsg.Foreground(lipgloss.Color(t.Primary))
	s.AIMsg = s.AIMsg.Foreground(lipgloss.Color(t.Success))
	s.SystemMsg = s.SystemMsg.Foreground(lipgloss.Color(t.Muted))
	s.Error = s.Error.Foreground(lipgloss.Color(t.Error))
	s.Success = s.Success.Foreground(lipgloss.Color(t.Success))
	s.Spinner = s.Spinner.Foreground(lipgloss.Color(t.Primary))
	s.SessionSelected = s.SessionSelected.Foreground(lipgloss.Color(t.Primary))
	s.SessionHeader = s.SessionHeader.Foreground(lipgloss.Color(t.Primary))
	s.TabActive = s.TabActive.Foreground(lipgloss.Color(t.Primary))
	s.PaneFocused = s.PaneFocused.BorderForeground(lipgloss.Color(t.Primary))
	s.PaneTitle = s.PaneTitle.Foreground(lipgloss.Color(t.Primary))
	s.AgentStatus = s.AgentStatus.Foreground(lipgloss.Color(t.Success))
	s.AgentPlan = s.AgentPlan.BorderForeground(lipgloss.Color(t.Success))
	s.AgentProgress = s.AgentProgress.Foreground(lipgloss.Color(t.Success))
}

// ApplyAdvancedTheme applies a full theme from the theme package
func (s *Styles) ApplyAdvancedTheme(t *theme.Theme) {
	if t == nil {
		return
	}

	c := t.Colors

	// Base styles
	s.Title = s.Title.Foreground(lipgloss.Color(c.Primary))
	s.Version = s.Version.Foreground(lipgloss.Color(c.Muted))
	s.Input = s.Input.BorderForeground(lipgloss.Color(c.Primary))

	// Status styles
	s.StatusConnected = s.StatusConnected.Foreground(lipgloss.Color(c.Success))
	s.StatusDisconnected = s.StatusDisconnected.Foreground(lipgloss.Color(c.Warning))

	// Help
	s.Help = s.Help.Foreground(lipgloss.Color(c.Muted))

	// Message styles
	s.UserMsg = s.UserMsg.Foreground(lipgloss.Color(c.UserMessage))
	s.AIMsg = s.AIMsg.Foreground(lipgloss.Color(c.AIMessage))
	s.SystemMsg = s.SystemMsg.Foreground(lipgloss.Color(c.SystemMessage))

	// Command styles
	s.Command = s.Command.
		Background(lipgloss.Color(c.Background)).
		Foreground(lipgloss.Color(c.Foreground))

	s.Output = s.Output.Foreground(lipgloss.Color(c.Foreground))

	// Feedback styles
	s.Error = s.Error.Foreground(lipgloss.Color(c.Error))
	s.Success = s.Success.Foreground(lipgloss.Color(c.Success))
	s.Spinner = s.Spinner.Foreground(lipgloss.Color(c.Primary))

	// History
	s.History = s.History.Foreground(lipgloss.Color(c.Muted))
	s.CopyNotice = s.CopyNotice.Foreground(lipgloss.Color(c.Success))

	// Session styles
	s.SessionItem = s.SessionItem.Foreground(lipgloss.Color(c.Foreground))
	s.SessionSelected = s.SessionSelected.Foreground(lipgloss.Color(c.Primary))
	s.SessionHeader = s.SessionHeader.Foreground(lipgloss.Color(c.Primary))

	// Tab styles
	s.Tab = s.Tab.Foreground(lipgloss.Color(c.Foreground))
	s.TabActive = s.TabActive.
		Foreground(lipgloss.Color(c.TabActive)).
		Background(lipgloss.Color(c.Background))
	s.TabInactive = s.TabInactive.Foreground(lipgloss.Color(c.TabInactive))
	s.TabBar = s.TabBar.BorderForeground(lipgloss.Color(c.Muted))

	// Pane styles
	s.Pane = s.Pane.BorderForeground(lipgloss.Color(c.PaneBorderUnfocused))
	s.PaneFocused = s.PaneFocused.BorderForeground(lipgloss.Color(c.PaneBorderFocused))
	s.PaneUnfocused = s.PaneUnfocused.BorderForeground(lipgloss.Color(c.PaneBorderUnfocused))
	s.PaneBorder = s.PaneBorder.BorderForeground(lipgloss.Color(c.PaneBorderUnfocused))
	s.PaneTitle = s.PaneTitle.Foreground(lipgloss.Color(c.Primary))

	// Agent styles
	s.AgentStatus = s.AgentStatus.Foreground(lipgloss.Color(c.Success))
	s.AgentPlan = s.AgentPlan.BorderForeground(lipgloss.Color(c.Secondary))
	s.AgentTask = s.AgentTask.Foreground(lipgloss.Color(c.Foreground))
	s.AgentProgress = s.AgentProgress.Foreground(lipgloss.Color(c.Secondary))
}

// NewStylesWithTheme creates new styles and applies a theme
func NewStylesWithTheme(t *theme.Theme) *Styles {
	s := DefaultStyles()
	if t != nil {
		s.ApplyAdvancedTheme(t)
	}
	return s
}
