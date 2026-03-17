package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	// AI settings
	AIProvider    string `mapstructure:"ai_provider"`
	AnthropicKey  string `mapstructure:"anthropic_api_key"`
	OpenAIKey     string `mapstructure:"openai_api_key"`

	// Ollama settings
	OllamaURL     string `mapstructure:"ollama_url"`
	OllamaModel   string `mapstructure:"ollama_model"`
	OllamaEnabled bool   `mapstructure:"ollama_enabled"`

	// Agent settings
	AgentMode     string `mapstructure:"agent_mode"` // off, suggest, auto
	AgentMaxTasks int    `mapstructure:"agent_max_tasks"`

	// UI settings
	Theme         string `mapstructure:"theme"`
	ShowWelcome   bool   `mapstructure:"show_welcome"`

	// Appearance settings
	Appearance AppearanceConfig `mapstructure:"appearance"`

	// Pane settings
	Pane PaneConfig `mapstructure:"pane"`

	// Workspace settings
	Workspace WorkspaceConfig `mapstructure:"workspace"`

	// Session settings
	SessionDir    string `mapstructure:"session_dir"`
	HistoryLimit  int    `mapstructure:"history_limit"`

	// SecretGuard settings
	SecretGuardEnabled bool `mapstructure:"secret_guard_enabled"`

	// Debug settings
	Debug    bool   `mapstructure:"debug"`
	LogLevel string `mapstructure:"log_level"`
}

// AppearanceConfig holds appearance-related settings
type AppearanceConfig struct {
	Theme           string  `mapstructure:"theme"`
	BackgroundColor string  `mapstructure:"background_color"`
	Transparency    float64 `mapstructure:"transparency"`
	EnableAnimations bool   `mapstructure:"enable_animations"`
	FontSize        int     `mapstructure:"font_size"`
	ThemeHotReload  bool    `mapstructure:"theme_hot_reload"`
}

// PaneConfig holds pane-related settings
type PaneConfig struct {
	BorderStyle       string  `mapstructure:"border_style"` // rounded, normal, double, hidden
	ActiveBorderWidth int     `mapstructure:"active_border_width"`
	InactiveOpacity   float64 `mapstructure:"inactive_opacity"`
	ShowPaneTitles    bool    `mapstructure:"show_pane_titles"`
	MinPaneWidth      int     `mapstructure:"min_pane_width"`
	MinPaneHeight     int     `mapstructure:"min_pane_height"`
}

// WorkspaceConfig holds workspace-related settings
type WorkspaceConfig struct {
	DefaultLayout     string `mapstructure:"default_layout"` // quad, tall, wide, stack, single
	AutoSave          bool   `mapstructure:"auto_save"`
	AutoSaveInterval  int    `mapstructure:"auto_save_interval"` // seconds
	RestoreOnStartup  bool   `mapstructure:"restore_on_startup"`
	MaxWorkspaces     int    `mapstructure:"max_workspaces"`
}

// Load reads configuration from file and environment
func Load() (*Config, error) {
	// Set defaults
	viper.SetDefault("ai_provider", "anthropic")
	viper.SetDefault("theme", "default")
	viper.SetDefault("show_welcome", true)
	viper.SetDefault("history_limit", 1000)
	viper.SetDefault("secret_guard_enabled", true)
	viper.SetDefault("debug", false)
	viper.SetDefault("log_level", "info")

	// Ollama defaults
	viper.SetDefault("ollama_url", "http://localhost:11434")
	viper.SetDefault("ollama_model", "codellama")
	viper.SetDefault("ollama_enabled", false)

	// Agent defaults
	viper.SetDefault("agent_mode", "suggest")
	viper.SetDefault("agent_max_tasks", 10)

	// Appearance defaults
	viper.SetDefault("appearance.theme", "default")
	viper.SetDefault("appearance.background_color", "")
	viper.SetDefault("appearance.transparency", 1.0)
	viper.SetDefault("appearance.enable_animations", true)
	viper.SetDefault("appearance.font_size", 14)
	viper.SetDefault("appearance.theme_hot_reload", true)

	// Pane defaults
	viper.SetDefault("pane.border_style", "rounded")
	viper.SetDefault("pane.active_border_width", 1)
	viper.SetDefault("pane.inactive_opacity", 0.8)
	viper.SetDefault("pane.show_pane_titles", true)
	viper.SetDefault("pane.min_pane_width", 20)
	viper.SetDefault("pane.min_pane_height", 5)

	// Workspace defaults
	viper.SetDefault("workspace.default_layout", "quad")
	viper.SetDefault("workspace.auto_save", true)
	viper.SetDefault("workspace.auto_save_interval", 60)
	viper.SetDefault("workspace.restore_on_startup", true)
	viper.SetDefault("workspace.max_workspaces", 10)

	// Config file locations
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(home, ".terminalizcrazy")
	viper.SetDefault("session_dir", filepath.Join(configDir, "sessions"))

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return nil, err
	}

	// Config file settings
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(configDir)
	viper.AddConfigPath(".")

	// Read config file (ignore if not found)
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, err
		}
	}

	// Environment variables override config file
	viper.SetEnvPrefix("TERMINALIZCRAZY")
	viper.AutomaticEnv()

	// Also read from .env file style (without prefix)
	viper.BindEnv("anthropic_api_key", "ANTHROPIC_API_KEY")
	viper.BindEnv("openai_api_key", "OPENAI_API_KEY")
	viper.BindEnv("ai_provider", "AI_PROVIDER")
	viper.BindEnv("debug", "DEBUG")
	viper.BindEnv("log_level", "LOG_LEVEL")

	// Ollama environment variables
	viper.BindEnv("ollama_url", "OLLAMA_URL")
	viper.BindEnv("ollama_model", "OLLAMA_MODEL")
	viper.BindEnv("ollama_enabled", "OLLAMA_ENABLED")

	// Agent environment variables
	viper.BindEnv("agent_mode", "AGENT_MODE")
	viper.BindEnv("agent_max_tasks", "AGENT_MAX_TASKS")

	// Unmarshal into config struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// HasAIKey returns true if at least one AI API key is configured
func (c *Config) HasAIKey() bool {
	return c.AnthropicKey != "" || c.OpenAIKey != ""
}

// GetActiveAIKey returns the API key for the configured provider
func (c *Config) GetActiveAIKey() string {
	switch c.AIProvider {
	case "openai":
		return c.OpenAIKey
	case "ollama":
		return "" // Ollama doesn't need an API key
	default:
		return c.AnthropicKey
	}
}

// IsOllamaConfigured returns true if Ollama is properly configured
func (c *Config) IsOllamaConfigured() bool {
	return c.OllamaEnabled && c.OllamaURL != ""
}

// GetOllamaURL returns the Ollama URL with default fallback
func (c *Config) GetOllamaURL() string {
	if c.OllamaURL == "" {
		return "http://localhost:11434"
	}
	return c.OllamaURL
}

// GetOllamaModel returns the Ollama model with default fallback
func (c *Config) GetOllamaModel() string {
	if c.OllamaModel == "" {
		return "codellama"
	}
	return c.OllamaModel
}

// GetAgentMode returns the agent mode with default fallback
func (c *Config) GetAgentMode() string {
	if c.AgentMode == "" {
		return "suggest"
	}
	return c.AgentMode
}

// GetTheme returns the theme name (from appearance or legacy theme field)
func (c *Config) GetTheme() string {
	if c.Appearance.Theme != "" {
		return c.Appearance.Theme
	}
	if c.Theme != "" {
		return c.Theme
	}
	return "default"
}

// GetThemesDir returns the themes directory path
func (c *Config) GetThemesDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".terminalizcrazy", "themes")
}

// GetConfigDir returns the config directory path
func (c *Config) GetConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".terminalizcrazy")
}

// GetDataDir returns the data directory path (for database, etc.)
func (c *Config) GetDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".terminalizcrazy")
}

// GetBorderStyle returns the pane border style with validation
func (c *Config) GetBorderStyle() string {
	style := c.Pane.BorderStyle
	switch style {
	case "rounded", "normal", "double", "hidden":
		return style
	default:
		return "rounded"
	}
}

// GetDefaultLayout returns the default workspace layout with validation
func (c *Config) GetDefaultLayout() string {
	layout := c.Workspace.DefaultLayout
	switch layout {
	case "quad", "tall", "wide", "stack", "single":
		return layout
	default:
		return "quad"
	}
}

// GetTransparency returns the transparency value clamped to valid range
func (c *Config) GetTransparency() float64 {
	t := c.Appearance.Transparency
	if t <= 0 {
		return 1.0
	}
	if t > 1 {
		return 1.0
	}
	return t
}

// GetInactiveOpacity returns the inactive pane opacity clamped to valid range
func (c *Config) GetInactiveOpacity() float64 {
	o := c.Pane.InactiveOpacity
	if o <= 0 {
		return 0.8
	}
	if o > 1 {
		return 1.0
	}
	return o
}
