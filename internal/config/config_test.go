package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetViper() {
	viper.Reset()
}

func TestLoad(t *testing.T) {
	resetViper()

	cfg, err := Load()

	require.NoError(t, err)
	assert.NotNil(t, cfg)

	// Check defaults
	assert.Equal(t, "anthropic", cfg.AIProvider)
	assert.Equal(t, "default", cfg.Theme)
	assert.True(t, cfg.ShowWelcome)
	assert.Equal(t, 1000, cfg.HistoryLimit)
	assert.True(t, cfg.SecretGuardEnabled)
	assert.False(t, cfg.Debug)
	assert.Equal(t, "info", cfg.LogLevel)
}

func TestLoadDefaults_Ollama(t *testing.T) {
	resetViper()

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, "http://localhost:11434", cfg.OllamaURL)
	assert.Equal(t, "codellama", cfg.OllamaModel)
	assert.False(t, cfg.OllamaEnabled)
}

func TestLoadDefaults_Agent(t *testing.T) {
	resetViper()

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, "suggest", cfg.AgentMode)
	assert.Equal(t, 10, cfg.AgentMaxTasks)
}

func TestLoadDefaults_Appearance(t *testing.T) {
	resetViper()

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, "default", cfg.Appearance.Theme)
	assert.Equal(t, "", cfg.Appearance.BackgroundColor)
	assert.Equal(t, 1.0, cfg.Appearance.Transparency)
	assert.True(t, cfg.Appearance.EnableAnimations)
	assert.Equal(t, 14, cfg.Appearance.FontSize)
	assert.True(t, cfg.Appearance.ThemeHotReload)
}

func TestLoadDefaults_Pane(t *testing.T) {
	resetViper()

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, "rounded", cfg.Pane.BorderStyle)
	assert.Equal(t, 1, cfg.Pane.ActiveBorderWidth)
	assert.Equal(t, 0.8, cfg.Pane.InactiveOpacity)
	assert.True(t, cfg.Pane.ShowPaneTitles)
	assert.Equal(t, 20, cfg.Pane.MinPaneWidth)
	assert.Equal(t, 5, cfg.Pane.MinPaneHeight)
}

func TestLoadDefaults_Workspace(t *testing.T) {
	resetViper()

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, "quad", cfg.Workspace.DefaultLayout)
	assert.True(t, cfg.Workspace.AutoSave)
	assert.Equal(t, 60, cfg.Workspace.AutoSaveInterval)
	assert.True(t, cfg.Workspace.RestoreOnStartup)
	assert.Equal(t, 10, cfg.Workspace.MaxWorkspaces)
}

func TestConfig_HasAIKey(t *testing.T) {
	tests := []struct {
		name         string
		anthropicKey string
		openAIKey    string
		want         bool
	}{
		{
			name:         "no keys",
			anthropicKey: "",
			openAIKey:    "",
			want:         false,
		},
		{
			name:         "anthropic key only",
			anthropicKey: "sk-ant-123",
			openAIKey:    "",
			want:         true,
		},
		{
			name:         "openai key only",
			anthropicKey: "",
			openAIKey:    "sk-123",
			want:         true,
		},
		{
			name:         "both keys",
			anthropicKey: "sk-ant-123",
			openAIKey:    "sk-123",
			want:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				AnthropicKey: tt.anthropicKey,
				OpenAIKey:    tt.openAIKey,
			}
			assert.Equal(t, tt.want, cfg.HasAIKey())
		})
	}
}

func TestConfig_GetActiveAIKey(t *testing.T) {
	tests := []struct {
		name         string
		provider     string
		anthropicKey string
		openAIKey    string
		want         string
	}{
		{
			name:         "anthropic provider",
			provider:     "anthropic",
			anthropicKey: "sk-ant-123",
			openAIKey:    "sk-456",
			want:         "sk-ant-123",
		},
		{
			name:         "openai provider",
			provider:     "openai",
			anthropicKey: "sk-ant-123",
			openAIKey:    "sk-456",
			want:         "sk-456",
		},
		{
			name:         "ollama provider (no key needed)",
			provider:     "ollama",
			anthropicKey: "sk-ant-123",
			openAIKey:    "sk-456",
			want:         "",
		},
		{
			name:         "default provider (anthropic)",
			provider:     "",
			anthropicKey: "sk-ant-123",
			openAIKey:    "sk-456",
			want:         "sk-ant-123",
		},
		{
			name:         "unknown provider defaults to anthropic",
			provider:     "unknown",
			anthropicKey: "sk-ant-123",
			openAIKey:    "sk-456",
			want:         "sk-ant-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				AIProvider:   tt.provider,
				AnthropicKey: tt.anthropicKey,
				OpenAIKey:    tt.openAIKey,
			}
			assert.Equal(t, tt.want, cfg.GetActiveAIKey())
		})
	}
}

func TestConfig_IsOllamaConfigured(t *testing.T) {
	tests := []struct {
		name    string
		enabled bool
		url     string
		want    bool
	}{
		{
			name:    "enabled with URL",
			enabled: true,
			url:     "http://localhost:11434",
			want:    true,
		},
		{
			name:    "enabled without URL",
			enabled: true,
			url:     "",
			want:    false,
		},
		{
			name:    "disabled with URL",
			enabled: false,
			url:     "http://localhost:11434",
			want:    false,
		},
		{
			name:    "disabled without URL",
			enabled: false,
			url:     "",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				OllamaEnabled: tt.enabled,
				OllamaURL:     tt.url,
			}
			assert.Equal(t, tt.want, cfg.IsOllamaConfigured())
		})
	}
}

func TestConfig_GetOllamaURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "custom URL",
			url:  "http://remote:11434",
			want: "http://remote:11434",
		},
		{
			name: "empty URL returns default",
			url:  "",
			want: "http://localhost:11434",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{OllamaURL: tt.url}
			assert.Equal(t, tt.want, cfg.GetOllamaURL())
		})
	}
}

func TestConfig_GetOllamaModel(t *testing.T) {
	tests := []struct {
		name  string
		model string
		want  string
	}{
		{
			name:  "custom model",
			model: "llama2",
			want:  "llama2",
		},
		{
			name:  "empty model returns default",
			model: "",
			want:  "codellama",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{OllamaModel: tt.model}
			assert.Equal(t, tt.want, cfg.GetOllamaModel())
		})
	}
}

func TestConfig_GetAgentMode(t *testing.T) {
	tests := []struct {
		name string
		mode string
		want string
	}{
		{
			name: "off mode",
			mode: "off",
			want: "off",
		},
		{
			name: "suggest mode",
			mode: "suggest",
			want: "suggest",
		},
		{
			name: "auto mode",
			mode: "auto",
			want: "auto",
		},
		{
			name: "empty returns default",
			mode: "",
			want: "suggest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{AgentMode: tt.mode}
			assert.Equal(t, tt.want, cfg.GetAgentMode())
		})
	}
}

func TestConfig_GetTheme(t *testing.T) {
	tests := []struct {
		name            string
		appearanceTheme string
		legacyTheme     string
		want            string
	}{
		{
			name:            "appearance theme takes precedence",
			appearanceTheme: "dracula",
			legacyTheme:     "nord",
			want:            "dracula",
		},
		{
			name:            "legacy theme fallback",
			appearanceTheme: "",
			legacyTheme:     "nord",
			want:            "nord",
		},
		{
			name:            "default when both empty",
			appearanceTheme: "",
			legacyTheme:     "",
			want:            "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Theme: tt.legacyTheme,
				Appearance: AppearanceConfig{
					Theme: tt.appearanceTheme,
				},
			}
			assert.Equal(t, tt.want, cfg.GetTheme())
		})
	}
}

func TestConfig_GetThemesDir(t *testing.T) {
	cfg := &Config{}
	dir := cfg.GetThemesDir()

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	expected := filepath.Join(home, ".terminalizcrazy", "themes")
	assert.Equal(t, expected, dir)
}

func TestConfig_GetConfigDir(t *testing.T) {
	cfg := &Config{}
	dir := cfg.GetConfigDir()

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	expected := filepath.Join(home, ".terminalizcrazy")
	assert.Equal(t, expected, dir)
}

func TestConfig_GetDataDir(t *testing.T) {
	cfg := &Config{}
	dir := cfg.GetDataDir()

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	expected := filepath.Join(home, ".terminalizcrazy")
	assert.Equal(t, expected, dir)
}

func TestConfig_GetBorderStyle(t *testing.T) {
	tests := []struct {
		name  string
		style string
		want  string
	}{
		{name: "rounded", style: "rounded", want: "rounded"},
		{name: "normal", style: "normal", want: "normal"},
		{name: "double", style: "double", want: "double"},
		{name: "hidden", style: "hidden", want: "hidden"},
		{name: "invalid returns default", style: "invalid", want: "rounded"},
		{name: "empty returns default", style: "", want: "rounded"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Pane: PaneConfig{BorderStyle: tt.style},
			}
			assert.Equal(t, tt.want, cfg.GetBorderStyle())
		})
	}
}

func TestConfig_GetDefaultLayout(t *testing.T) {
	tests := []struct {
		name   string
		layout string
		want   string
	}{
		{name: "quad", layout: "quad", want: "quad"},
		{name: "tall", layout: "tall", want: "tall"},
		{name: "wide", layout: "wide", want: "wide"},
		{name: "stack", layout: "stack", want: "stack"},
		{name: "single", layout: "single", want: "single"},
		{name: "invalid returns default", layout: "invalid", want: "quad"},
		{name: "empty returns default", layout: "", want: "quad"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Workspace: WorkspaceConfig{DefaultLayout: tt.layout},
			}
			assert.Equal(t, tt.want, cfg.GetDefaultLayout())
		})
	}
}

func TestConfig_GetTransparency(t *testing.T) {
	tests := []struct {
		name         string
		transparency float64
		want         float64
	}{
		{name: "valid 0.5", transparency: 0.5, want: 0.5},
		{name: "valid 1.0", transparency: 1.0, want: 1.0},
		{name: "valid 0.1", transparency: 0.1, want: 0.1},
		{name: "zero returns default", transparency: 0.0, want: 1.0},
		{name: "negative returns default", transparency: -0.5, want: 1.0},
		{name: "above 1 returns default", transparency: 1.5, want: 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Appearance: AppearanceConfig{Transparency: tt.transparency},
			}
			assert.Equal(t, tt.want, cfg.GetTransparency())
		})
	}
}

func TestConfig_GetInactiveOpacity(t *testing.T) {
	tests := []struct {
		name    string
		opacity float64
		want    float64
	}{
		{name: "valid 0.5", opacity: 0.5, want: 0.5},
		{name: "valid 1.0", opacity: 1.0, want: 1.0},
		{name: "valid 0.1", opacity: 0.1, want: 0.1},
		{name: "zero returns default", opacity: 0.0, want: 0.8},
		{name: "negative returns default", opacity: -0.5, want: 0.8},
		{name: "above 1 clamps to 1", opacity: 1.5, want: 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Pane: PaneConfig{InactiveOpacity: tt.opacity},
			}
			assert.Equal(t, tt.want, cfg.GetInactiveOpacity())
		})
	}
}

func TestLoadWithEnvVars(t *testing.T) {
	resetViper()

	// Set environment variables
	os.Setenv("ANTHROPIC_API_KEY", "test-anthropic-key")
	os.Setenv("OPENAI_API_KEY", "test-openai-key")
	os.Setenv("AI_PROVIDER", "openai")
	defer func() {
		os.Unsetenv("ANTHROPIC_API_KEY")
		os.Unsetenv("OPENAI_API_KEY")
		os.Unsetenv("AI_PROVIDER")
	}()

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, "test-anthropic-key", cfg.AnthropicKey)
	assert.Equal(t, "test-openai-key", cfg.OpenAIKey)
	assert.Equal(t, "openai", cfg.AIProvider)
}

func TestLoadWithOllamaEnvVars(t *testing.T) {
	resetViper()

	os.Setenv("OLLAMA_URL", "http://custom:11434")
	os.Setenv("OLLAMA_MODEL", "custom-model")
	os.Setenv("OLLAMA_ENABLED", "true")
	defer func() {
		os.Unsetenv("OLLAMA_URL")
		os.Unsetenv("OLLAMA_MODEL")
		os.Unsetenv("OLLAMA_ENABLED")
	}()

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, "http://custom:11434", cfg.OllamaURL)
	assert.Equal(t, "custom-model", cfg.OllamaModel)
	assert.True(t, cfg.OllamaEnabled)
}

func TestAppearanceConfig(t *testing.T) {
	cfg := AppearanceConfig{
		Theme:            "dracula",
		BackgroundColor:  "#282a36",
		Transparency:     0.95,
		EnableAnimations: false,
		FontSize:         16,
		ThemeHotReload:   false,
	}

	assert.Equal(t, "dracula", cfg.Theme)
	assert.Equal(t, "#282a36", cfg.BackgroundColor)
	assert.Equal(t, 0.95, cfg.Transparency)
	assert.False(t, cfg.EnableAnimations)
	assert.Equal(t, 16, cfg.FontSize)
	assert.False(t, cfg.ThemeHotReload)
}

func TestPaneConfig(t *testing.T) {
	cfg := PaneConfig{
		BorderStyle:       "double",
		ActiveBorderWidth: 2,
		InactiveOpacity:   0.6,
		ShowPaneTitles:    false,
		MinPaneWidth:      30,
		MinPaneHeight:     10,
	}

	assert.Equal(t, "double", cfg.BorderStyle)
	assert.Equal(t, 2, cfg.ActiveBorderWidth)
	assert.Equal(t, 0.6, cfg.InactiveOpacity)
	assert.False(t, cfg.ShowPaneTitles)
	assert.Equal(t, 30, cfg.MinPaneWidth)
	assert.Equal(t, 10, cfg.MinPaneHeight)
}

func TestWorkspaceConfig(t *testing.T) {
	cfg := WorkspaceConfig{
		DefaultLayout:    "tall",
		AutoSave:         false,
		AutoSaveInterval: 120,
		RestoreOnStartup: false,
		MaxWorkspaces:    20,
	}

	assert.Equal(t, "tall", cfg.DefaultLayout)
	assert.False(t, cfg.AutoSave)
	assert.Equal(t, 120, cfg.AutoSaveInterval)
	assert.False(t, cfg.RestoreOnStartup)
	assert.Equal(t, 20, cfg.MaxWorkspaces)
}
