package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/terminalizcrazy/terminalizcrazy/internal/crypto"
)

// Config holds all application configuration
type Config struct {
	// AI settings
	AIProvider    string `mapstructure:"ai_provider"`
	AnthropicKey  string `mapstructure:"anthropic_api_key"`
	OpenAIKey     string `mapstructure:"openai_api_key"`
	GeminiKey     string `mapstructure:"gemini_api_key"`
	GeminiModel   string `mapstructure:"gemini_model"` // Default: gemini-1.5-flash

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

	// Data retention settings (GDPR compliance)
	Retention RetentionConfig `mapstructure:"retention"`

	// Debug settings
	Debug    bool   `mapstructure:"debug"`
	LogLevel string `mapstructure:"log_level"`
}

// RetentionConfig holds data retention policy settings
type RetentionConfig struct {
	// MessageRetentionDays is how long to keep chat messages (0 = forever)
	MessageRetentionDays int `mapstructure:"message_retention_days"`
	// CommandHistoryRetentionDays is how long to keep command history (0 = forever)
	CommandHistoryRetentionDays int `mapstructure:"command_history_retention_days"`
	// AgentPlanRetentionDays is how long to keep agent plans (0 = forever)
	AgentPlanRetentionDays int `mapstructure:"agent_plan_retention_days"`
	// AutoCleanupEnabled enables automatic cleanup on startup
	AutoCleanupEnabled bool `mapstructure:"auto_cleanup_enabled"`
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
	viper.SetDefault("ai_provider", "ollama") // Default to Ollama (local)
	viper.SetDefault("gemini_model", "gemini-1.5-flash") // Gemini model if used
	viper.SetDefault("theme", "default")
	viper.SetDefault("show_welcome", true)
	viper.SetDefault("history_limit", 1000)
	viper.SetDefault("secret_guard_enabled", true)
	viper.SetDefault("debug", false)
	viper.SetDefault("log_level", "info")

	// Ollama defaults (now the primary provider)
	viper.SetDefault("ollama_url", "http://localhost:11434")
	viper.SetDefault("ollama_model", "gemma4")
	viper.SetDefault("ollama_enabled", true)

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

	// Retention defaults (GDPR compliance)
	viper.SetDefault("retention.message_retention_days", 90)          // 90 days default
	viper.SetDefault("retention.command_history_retention_days", 90)  // 90 days default
	viper.SetDefault("retention.agent_plan_retention_days", 30)       // 30 days default
	viper.SetDefault("retention.auto_cleanup_enabled", true)          // Enabled by default

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
	viper.BindEnv("gemini_api_key", "GEMINI_API_KEY")
	viper.BindEnv("gemini_model", "GEMINI_MODEL")
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

	// Decrypt API keys if they are encrypted
	if err := cfg.decryptAPIKeys(); err != nil {
		return nil, fmt.Errorf("failed to decrypt API keys: %w", err)
	}

	return &cfg, nil
}

// decryptAPIKeys decrypts any encrypted API keys in the config
func (c *Config) decryptAPIKeys() error {
	km, err := crypto.NewKeyManager()
	if err != nil {
		return err
	}

	// Decrypt each API key if it's encrypted
	if c.GeminiKey != "" {
		decrypted, err := km.Decrypt(c.GeminiKey)
		if err != nil {
			return fmt.Errorf("gemini key: %w", err)
		}
		c.GeminiKey = decrypted
	}

	if c.AnthropicKey != "" {
		decrypted, err := km.Decrypt(c.AnthropicKey)
		if err != nil {
			return fmt.Errorf("anthropic key: %w", err)
		}
		c.AnthropicKey = decrypted
	}

	if c.OpenAIKey != "" {
		decrypted, err := km.Decrypt(c.OpenAIKey)
		if err != nil {
			return fmt.Errorf("openai key: %w", err)
		}
		c.OpenAIKey = decrypted
	}

	return nil
}

// HasAIKey returns true if at least one AI API key is configured
func (c *Config) HasAIKey() bool {
	return c.GeminiKey != "" || c.AnthropicKey != "" || c.OpenAIKey != ""
}

// GetActiveAIKey returns the API key for the configured provider
func (c *Config) GetActiveAIKey() string {
	switch c.AIProvider {
	case "gemini":
		return c.GeminiKey
	case "openai":
		return c.OpenAIKey
	case "anthropic":
		return c.AnthropicKey
	case "ollama":
		return "" // Ollama doesn't need an API key
	default:
		return c.GeminiKey // Default to Gemini
	}
}

// GetGeminiModel returns the Gemini model with default fallback
func (c *Config) GetGeminiModel() string {
	if c.GeminiModel == "" {
		return "gemini-1.5-flash"
	}
	return c.GeminiModel
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

// GetMessageRetentionDays returns message retention period (0 = forever)
func (c *Config) GetMessageRetentionDays() int {
	if c.Retention.MessageRetentionDays < 0 {
		return 0
	}
	return c.Retention.MessageRetentionDays
}

// GetCommandHistoryRetentionDays returns command history retention period (0 = forever)
func (c *Config) GetCommandHistoryRetentionDays() int {
	if c.Retention.CommandHistoryRetentionDays < 0 {
		return 0
	}
	return c.Retention.CommandHistoryRetentionDays
}

// GetAgentPlanRetentionDays returns agent plan retention period (0 = forever)
func (c *Config) GetAgentPlanRetentionDays() int {
	if c.Retention.AgentPlanRetentionDays < 0 {
		return 0
	}
	return c.Retention.AgentPlanRetentionDays
}

// IsAutoCleanupEnabled returns whether automatic data cleanup is enabled
func (c *Config) IsAutoCleanupEnabled() bool {
	return c.Retention.AutoCleanupEnabled
}

// EncryptAPIKey encrypts an API key for storage in config file
func EncryptAPIKey(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	km, err := crypto.NewKeyManager()
	if err != nil {
		return "", err
	}

	return km.Encrypt(plaintext)
}

// IsAPIKeyEncrypted checks if an API key value is encrypted
func IsAPIKeyEncrypted(value string) bool {
	return crypto.IsEncrypted(value)
}
