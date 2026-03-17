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

	// UI settings
	Theme         string `mapstructure:"theme"`
	ShowWelcome   bool   `mapstructure:"show_welcome"`

	// Session settings
	SessionDir    string `mapstructure:"session_dir"`
	HistoryLimit  int    `mapstructure:"history_limit"`

	// SecretGuard settings
	SecretGuardEnabled bool `mapstructure:"secret_guard_enabled"`

	// Debug settings
	Debug    bool   `mapstructure:"debug"`
	LogLevel string `mapstructure:"log_level"`
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
	default:
		return c.AnthropicKey
	}
}
