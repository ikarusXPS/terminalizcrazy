package theme

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// ColorPalette defines the full color palette for a theme
type ColorPalette struct {
	// Base colors
	Background string `yaml:"background"`
	Foreground string `yaml:"foreground"`

	// Semantic colors
	Primary   string `yaml:"primary"`
	Secondary string `yaml:"secondary"`
	Warning   string `yaml:"warning"`
	Error     string `yaml:"error"`
	Success   string `yaml:"success"`
	Muted     string `yaml:"muted"`

	// Chat colors
	UserMessage   string `yaml:"user_message"`
	AIMessage     string `yaml:"ai_message"`
	SystemMessage string `yaml:"system_message"`

	// UI element colors
	PaneBorderFocused   string `yaml:"pane_border_focused"`
	PaneBorderUnfocused string `yaml:"pane_border_unfocused"`
	TabActive           string `yaml:"tab_active"`
	TabInactive         string `yaml:"tab_inactive"`

	// Additional colors
	Selection string `yaml:"selection"`
	Comment   string `yaml:"comment"`
	Cyan      string `yaml:"cyan"`
	Green     string `yaml:"green"`
	Orange    string `yaml:"orange"`
	Pink      string `yaml:"pink"`
	Purple    string `yaml:"purple"`
	Red       string `yaml:"red"`
	Yellow    string `yaml:"yellow"`
}

// Theme represents a color theme
type Theme struct {
	Name    string       `yaml:"name"`
	Author  string       `yaml:"author"`
	Version string       `yaml:"version"`
	Colors  ColorPalette `yaml:"colors"`
}

// Validate checks if the theme has all required colors
func (t *Theme) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("theme name is required")
	}
	if t.Colors.Background == "" {
		return fmt.Errorf("background color is required")
	}
	if t.Colors.Foreground == "" {
		return fmt.Errorf("foreground color is required")
	}
	if t.Colors.Primary == "" {
		return fmt.Errorf("primary color is required")
	}
	return nil
}

// ApplyDefaults fills in missing colors with sensible defaults
func (t *Theme) ApplyDefaults() {
	defaults := DefaultTheme()

	if t.Colors.Secondary == "" {
		t.Colors.Secondary = defaults.Colors.Secondary
	}
	if t.Colors.Warning == "" {
		t.Colors.Warning = defaults.Colors.Warning
	}
	if t.Colors.Error == "" {
		t.Colors.Error = defaults.Colors.Error
	}
	if t.Colors.Success == "" {
		t.Colors.Success = defaults.Colors.Success
	}
	if t.Colors.Muted == "" {
		t.Colors.Muted = defaults.Colors.Muted
	}
	if t.Colors.UserMessage == "" {
		t.Colors.UserMessage = t.Colors.Primary
	}
	if t.Colors.AIMessage == "" {
		t.Colors.AIMessage = t.Colors.Secondary
	}
	if t.Colors.SystemMessage == "" {
		t.Colors.SystemMessage = t.Colors.Muted
	}
	if t.Colors.PaneBorderFocused == "" {
		t.Colors.PaneBorderFocused = t.Colors.Primary
	}
	if t.Colors.PaneBorderUnfocused == "" {
		t.Colors.PaneBorderUnfocused = t.Colors.Muted
	}
	if t.Colors.TabActive == "" {
		t.Colors.TabActive = t.Colors.Primary
	}
	if t.Colors.TabInactive == "" {
		t.Colors.TabInactive = t.Colors.Muted
	}
	if t.Colors.Selection == "" {
		t.Colors.Selection = "#44475a"
	}
	if t.Colors.Comment == "" {
		t.Colors.Comment = t.Colors.Muted
	}
}

// Color returns a lipgloss.Color for the given hex color
func (t *Theme) Color(hex string) lipgloss.Color {
	return lipgloss.Color(hex)
}

// BackgroundColor returns the background color as lipgloss.Color
func (t *Theme) BackgroundColor() lipgloss.Color {
	return lipgloss.Color(t.Colors.Background)
}

// ForegroundColor returns the foreground color as lipgloss.Color
func (t *Theme) ForegroundColor() lipgloss.Color {
	return lipgloss.Color(t.Colors.Foreground)
}

// PrimaryColor returns the primary color as lipgloss.Color
func (t *Theme) PrimaryColor() lipgloss.Color {
	return lipgloss.Color(t.Colors.Primary)
}

// SecondaryColor returns the secondary color as lipgloss.Color
func (t *Theme) SecondaryColor() lipgloss.Color {
	return lipgloss.Color(t.Colors.Secondary)
}

// WarningColor returns the warning color as lipgloss.Color
func (t *Theme) WarningColor() lipgloss.Color {
	return lipgloss.Color(t.Colors.Warning)
}

// ErrorColor returns the error color as lipgloss.Color
func (t *Theme) ErrorColor() lipgloss.Color {
	return lipgloss.Color(t.Colors.Error)
}

// SuccessColor returns the success color as lipgloss.Color
func (t *Theme) SuccessColor() lipgloss.Color {
	return lipgloss.Color(t.Colors.Success)
}

// MutedColor returns the muted color as lipgloss.Color
func (t *Theme) MutedColor() lipgloss.Color {
	return lipgloss.Color(t.Colors.Muted)
}

// Clone creates a deep copy of the theme
func (t *Theme) Clone() *Theme {
	return &Theme{
		Name:    t.Name,
		Author:  t.Author,
		Version: t.Version,
		Colors:  t.Colors,
	}
}

// DefaultTheme returns the default built-in theme
func DefaultTheme() *Theme {
	return &Theme{
		Name:    "Default",
		Author:  "TerminalizCrazy",
		Version: "1.0.0",
		Colors: ColorPalette{
			Background: "#1e1e2e",
			Foreground: "#cdd6f4",
			Primary:    "#7D56F4",
			Secondary:  "#04B575",
			Warning:    "#FFAA00",
			Error:      "#FF6B6B",
			Success:    "#04B575",
			Muted:      "#888888",

			UserMessage:   "#7D56F4",
			AIMessage:     "#04B575",
			SystemMessage: "#888888",

			PaneBorderFocused:   "#7D56F4",
			PaneBorderUnfocused: "#888888",
			TabActive:           "#7D56F4",
			TabInactive:         "#888888",

			Selection: "#44475a",
			Comment:   "#888888",
		},
	}
}
