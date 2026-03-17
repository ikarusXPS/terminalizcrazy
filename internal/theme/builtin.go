package theme

// BuiltinThemes returns all built-in themes
func BuiltinThemes() map[string]*Theme {
	return map[string]*Theme{
		"dracula":           DraculaTheme(),
		"nord":              NordTheme(),
		"catppuccin-mocha":  CatppuccinMochaTheme(),
		"gruvbox-dark":      GruvboxDarkTheme(),
		"tokyo-night":       TokyoNightTheme(),
	}
}

// DraculaTheme returns the Dracula theme
func DraculaTheme() *Theme {
	return &Theme{
		Name:    "Dracula",
		Author:  "Zeno Rocha",
		Version: "1.0.0",
		Colors: ColorPalette{
			Background: "#282a36",
			Foreground: "#f8f8f2",
			Primary:    "#bd93f9",
			Secondary:  "#50fa7b",
			Warning:    "#ffb86c",
			Error:      "#ff5555",
			Success:    "#50fa7b",
			Muted:      "#6272a4",

			UserMessage:   "#ff79c6",
			AIMessage:     "#50fa7b",
			SystemMessage: "#6272a4",

			PaneBorderFocused:   "#bd93f9",
			PaneBorderUnfocused: "#6272a4",
			TabActive:           "#bd93f9",
			TabInactive:         "#6272a4",

			Selection: "#44475a",
			Comment:   "#6272a4",
			Cyan:      "#8be9fd",
			Green:     "#50fa7b",
			Orange:    "#ffb86c",
			Pink:      "#ff79c6",
			Purple:    "#bd93f9",
			Red:       "#ff5555",
			Yellow:    "#f1fa8c",
		},
	}
}

// NordTheme returns the Nord theme
func NordTheme() *Theme {
	return &Theme{
		Name:    "Nord",
		Author:  "Arctic Ice Studio",
		Version: "1.0.0",
		Colors: ColorPalette{
			Background: "#2e3440",
			Foreground: "#eceff4",
			Primary:    "#88c0d0",
			Secondary:  "#a3be8c",
			Warning:    "#ebcb8b",
			Error:      "#bf616a",
			Success:    "#a3be8c",
			Muted:      "#4c566a",

			UserMessage:   "#81a1c1",
			AIMessage:     "#a3be8c",
			SystemMessage: "#4c566a",

			PaneBorderFocused:   "#88c0d0",
			PaneBorderUnfocused: "#4c566a",
			TabActive:           "#88c0d0",
			TabInactive:         "#4c566a",

			Selection: "#434c5e",
			Comment:   "#616e88",
			Cyan:      "#8fbcbb",
			Green:     "#a3be8c",
			Orange:    "#d08770",
			Pink:      "#b48ead",
			Purple:    "#b48ead",
			Red:       "#bf616a",
			Yellow:    "#ebcb8b",
		},
	}
}

// CatppuccinMochaTheme returns the Catppuccin Mocha theme
func CatppuccinMochaTheme() *Theme {
	return &Theme{
		Name:    "Catppuccin Mocha",
		Author:  "Catppuccin",
		Version: "1.0.0",
		Colors: ColorPalette{
			Background: "#1e1e2e",
			Foreground: "#cdd6f4",
			Primary:    "#cba6f7",
			Secondary:  "#a6e3a1",
			Warning:    "#fab387",
			Error:      "#f38ba8",
			Success:    "#a6e3a1",
			Muted:      "#6c7086",

			UserMessage:   "#f5c2e7",
			AIMessage:     "#a6e3a1",
			SystemMessage: "#6c7086",

			PaneBorderFocused:   "#cba6f7",
			PaneBorderUnfocused: "#6c7086",
			TabActive:           "#cba6f7",
			TabInactive:         "#6c7086",

			Selection: "#313244",
			Comment:   "#6c7086",
			Cyan:      "#94e2d5",
			Green:     "#a6e3a1",
			Orange:    "#fab387",
			Pink:      "#f5c2e7",
			Purple:    "#cba6f7",
			Red:       "#f38ba8",
			Yellow:    "#f9e2af",
		},
	}
}

// GruvboxDarkTheme returns the Gruvbox Dark theme
func GruvboxDarkTheme() *Theme {
	return &Theme{
		Name:    "Gruvbox Dark",
		Author:  "morhetz",
		Version: "1.0.0",
		Colors: ColorPalette{
			Background: "#282828",
			Foreground: "#ebdbb2",
			Primary:    "#d3869b",
			Secondary:  "#b8bb26",
			Warning:    "#fe8019",
			Error:      "#fb4934",
			Success:    "#b8bb26",
			Muted:      "#928374",

			UserMessage:   "#83a598",
			AIMessage:     "#b8bb26",
			SystemMessage: "#928374",

			PaneBorderFocused:   "#d3869b",
			PaneBorderUnfocused: "#928374",
			TabActive:           "#d3869b",
			TabInactive:         "#928374",

			Selection: "#3c3836",
			Comment:   "#928374",
			Cyan:      "#8ec07c",
			Green:     "#b8bb26",
			Orange:    "#fe8019",
			Pink:      "#d3869b",
			Purple:    "#d3869b",
			Red:       "#fb4934",
			Yellow:    "#fabd2f",
		},
	}
}

// TokyoNightTheme returns the Tokyo Night theme
func TokyoNightTheme() *Theme {
	return &Theme{
		Name:    "Tokyo Night",
		Author:  "folke",
		Version: "1.0.0",
		Colors: ColorPalette{
			Background: "#1a1b26",
			Foreground: "#c0caf5",
			Primary:    "#7aa2f7",
			Secondary:  "#9ece6a",
			Warning:    "#e0af68",
			Error:      "#f7768e",
			Success:    "#9ece6a",
			Muted:      "#565f89",

			UserMessage:   "#bb9af7",
			AIMessage:     "#9ece6a",
			SystemMessage: "#565f89",

			PaneBorderFocused:   "#7aa2f7",
			PaneBorderUnfocused: "#565f89",
			TabActive:           "#7aa2f7",
			TabInactive:         "#565f89",

			Selection: "#283457",
			Comment:   "#565f89",
			Cyan:      "#7dcfff",
			Green:     "#9ece6a",
			Orange:    "#ff9e64",
			Pink:      "#bb9af7",
			Purple:    "#bb9af7",
			Red:       "#f7768e",
			Yellow:    "#e0af68",
		},
	}
}

// GetBuiltinTheme returns a built-in theme by name
func GetBuiltinTheme(name string) *Theme {
	themes := BuiltinThemes()
	if theme, ok := themes[name]; ok {
		return theme.Clone()
	}
	return nil
}

// ListBuiltinThemes returns a list of built-in theme names
func ListBuiltinThemes() []string {
	return []string{
		"dracula",
		"nord",
		"catppuccin-mocha",
		"gruvbox-dark",
		"tokyo-night",
	}
}
