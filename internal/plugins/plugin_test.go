package plugins

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPluginManager(t *testing.T) {
	pm := NewPluginManager("/tmp/plugins")

	assert.NotNil(t, pm)
	assert.Equal(t, "/tmp/plugins", pm.GetPluginDir())
	assert.Empty(t, pm.ListPlugins())
}

func TestPluginManager_GetPluginDir(t *testing.T) {
	pm := NewPluginManager("/custom/path")

	assert.Equal(t, "/custom/path", pm.GetPluginDir())
}

func TestPluginManager_LoadPlugins_EmptyDir(t *testing.T) {
	pm := NewPluginManager("")

	err := pm.LoadPlugins()

	assert.NoError(t, err)
}

func TestPluginManager_LoadPlugins_WithDir(t *testing.T) {
	tmpDir := t.TempDir()
	pm := NewPluginManager(tmpDir)

	err := pm.LoadPlugins()

	assert.NoError(t, err)
}

func TestPluginManager_LoadPlugins_WithManifest(t *testing.T) {
	tmpDir := t.TempDir()
	pluginDir := filepath.Join(tmpDir, "test-plugin")
	err := os.MkdirAll(pluginDir, 0755)
	require.NoError(t, err)

	manifest := `{
		"name": "Test Plugin",
		"version": "1.0.0",
		"description": "A test plugin",
		"author": "Test Author",
		"type": "native",
		"hooks": ["pre_command"],
		"main": "main.go"
	}`
	err = os.WriteFile(filepath.Join(pluginDir, "plugin.json"), []byte(manifest), 0644)
	require.NoError(t, err)

	pm := NewPluginManager(tmpDir)
	err = pm.LoadPlugins()

	require.NoError(t, err)
	plugins := pm.ListPlugins()
	assert.Len(t, plugins, 1)
	assert.Equal(t, "Test Plugin", plugins[0].Name)
	assert.Equal(t, "1.0.0", plugins[0].Version)
	assert.Equal(t, PluginStateUnloaded, plugins[0].State)
}

func TestPluginManager_EnableDisablePlugin(t *testing.T) {
	tmpDir := t.TempDir()
	pluginDir := filepath.Join(tmpDir, "test-plugin")
	err := os.MkdirAll(pluginDir, 0755)
	require.NoError(t, err)

	manifest := `{"name": "Test", "version": "1.0", "type": "native", "hooks": ["pre_command"]}`
	err = os.WriteFile(filepath.Join(pluginDir, "plugin.json"), []byte(manifest), 0644)
	require.NoError(t, err)

	pm := NewPluginManager(tmpDir)
	err = pm.LoadPlugins()
	require.NoError(t, err)

	// Enable plugin
	err = pm.EnablePlugin("test-plugin")
	assert.NoError(t, err)

	plugin := pm.GetPlugin("test-plugin")
	assert.Equal(t, PluginStateActive, plugin.State)
	assert.True(t, plugin.Config.Enabled)
	assert.NotNil(t, plugin.LoadedAt)

	// Disable plugin
	err = pm.DisablePlugin("test-plugin")
	assert.NoError(t, err)

	plugin = pm.GetPlugin("test-plugin")
	assert.Equal(t, PluginStateStopped, plugin.State)
	assert.False(t, plugin.Config.Enabled)
}

func TestPluginManager_EnablePlugin_NotFound(t *testing.T) {
	pm := NewPluginManager("")

	err := pm.EnablePlugin("nonexistent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plugin not found")
}

func TestPluginManager_DisablePlugin_NotFound(t *testing.T) {
	pm := NewPluginManager("")

	err := pm.DisablePlugin("nonexistent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plugin not found")
}

func TestPluginManager_RegisterHandler(t *testing.T) {
	pm := NewPluginManager("")
	handler := NewSafetyPlugin()

	err := pm.RegisterHandler("safety", handler)

	assert.NoError(t, err)
	plugin := pm.GetPlugin("safety")
	assert.NotNil(t, plugin)
	assert.Equal(t, "safety", plugin.ID)
	assert.Equal(t, PluginStateActive, plugin.State)
	assert.Equal(t, PluginTypeNative, plugin.Type)
}

func TestPluginManager_GetPlugin(t *testing.T) {
	pm := NewPluginManager("")
	handler := NewSafetyPlugin()
	pm.RegisterHandler("safety", handler)

	plugin := pm.GetPlugin("safety")
	assert.NotNil(t, plugin)
	assert.Equal(t, "safety", plugin.ID)

	nilPlugin := pm.GetPlugin("nonexistent")
	assert.Nil(t, nilPlugin)
}

func TestPluginManager_ListPlugins(t *testing.T) {
	pm := NewPluginManager("")
	pm.RegisterHandler("safety", NewSafetyPlugin())
	pm.RegisterHandler("alias", NewAliasPlugin())

	plugins := pm.ListPlugins()

	assert.Len(t, plugins, 2)
}

func TestPluginManager_ListActivePlugins(t *testing.T) {
	pm := NewPluginManager("")
	pm.RegisterHandler("safety", NewSafetyPlugin())
	pm.RegisterHandler("alias", NewAliasPlugin())

	// Disable one
	pm.DisablePlugin("safety")

	active := pm.ListActivePlugins()
	assert.Len(t, active, 1)
	assert.Equal(t, "alias", active[0].ID)
}

func TestPluginManager_SetGetPluginSetting(t *testing.T) {
	pm := NewPluginManager("")
	pm.RegisterHandler("safety", NewSafetyPlugin())

	err := pm.SetPluginSetting("safety", "custom_key", "custom_value")
	assert.NoError(t, err)

	val, ok := pm.GetPluginSetting("safety", "custom_key")
	assert.True(t, ok)
	assert.Equal(t, "custom_value", val)

	// Non-existent key
	val, ok = pm.GetPluginSetting("safety", "nonexistent")
	assert.False(t, ok)
	assert.Nil(t, val)
}

func TestPluginManager_SetPluginSetting_NotFound(t *testing.T) {
	pm := NewPluginManager("")

	err := pm.SetPluginSetting("nonexistent", "key", "value")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plugin not found")
}

func TestPluginManager_GetPluginSetting_NotFound(t *testing.T) {
	pm := NewPluginManager("")

	val, ok := pm.GetPluginSetting("nonexistent", "key")

	assert.False(t, ok)
	assert.Nil(t, val)
}

func TestPluginManager_SetPluginPriority(t *testing.T) {
	pm := NewPluginManager("")
	pm.RegisterHandler("safety", NewSafetyPlugin())

	err := pm.SetPluginPriority("safety", 50)

	assert.NoError(t, err)
	plugin := pm.GetPlugin("safety")
	assert.Equal(t, 50, plugin.Config.Priority)
}

func TestPluginManager_SetPluginPriority_NotFound(t *testing.T) {
	pm := NewPluginManager("")

	err := pm.SetPluginPriority("nonexistent", 50)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plugin not found")
}

func TestPluginManager_UnloadAll(t *testing.T) {
	pm := NewPluginManager("")
	pm.RegisterHandler("safety", NewSafetyPlugin())
	pm.RegisterHandler("alias", NewAliasPlugin())

	pm.UnloadAll()

	plugins := pm.ListPlugins()
	for _, p := range plugins {
		assert.Equal(t, PluginStateUnloaded, p.State)
	}
}

func TestPluginManager_ExecuteHook(t *testing.T) {
	pm := NewPluginManager("")
	pm.RegisterHandler("alias", NewAliasPlugin())

	hookCtx := &HookContext{
		HookType:  HookPreCommand,
		Command:   "ll",
		Timestamp: time.Now(),
	}

	result, err := pm.ExecuteHook(context.Background(), hookCtx)

	assert.NoError(t, err)
	assert.True(t, result.Modified)
	assert.Equal(t, "ls -la", result.NewCommand)
}

func TestPluginManager_ExecuteHook_Cancel(t *testing.T) {
	pm := NewPluginManager("")
	pm.RegisterHandler("safety", NewSafetyPlugin())

	hookCtx := &HookContext{
		HookType:  HookPreCommand,
		Command:   "rm -rf /",
		Timestamp: time.Now(),
	}

	result, err := pm.ExecuteHook(context.Background(), hookCtx)

	assert.NoError(t, err)
	assert.True(t, result.Cancel)
}

func TestPluginManager_ExecuteHook_DisabledPlugin(t *testing.T) {
	pm := NewPluginManager("")
	pm.RegisterHandler("alias", NewAliasPlugin())
	pm.DisablePlugin("alias")

	hookCtx := &HookContext{
		HookType:  HookPreCommand,
		Command:   "ll",
		Timestamp: time.Now(),
	}

	result, err := pm.ExecuteHook(context.Background(), hookCtx)

	assert.NoError(t, err)
	assert.False(t, result.Modified)
}

func TestHookTypes(t *testing.T) {
	tests := []struct {
		hook HookType
		want string
	}{
		{HookPreCommand, "pre_command"},
		{HookPostCommand, "post_command"},
		{HookPreAI, "pre_ai"},
		{HookPostAI, "post_ai"},
		{HookOnKeyPress, "on_key_press"},
		{HookOnMessage, "on_message"},
		{HookOnSessionStart, "on_session_start"},
		{HookOnSessionEnd, "on_session_end"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			assert.Equal(t, HookType(tt.want), tt.hook)
		})
	}
}

func TestPluginTypes(t *testing.T) {
	assert.Equal(t, PluginType("native"), PluginTypeNative)
	assert.Equal(t, PluginType("wasm"), PluginTypeWASM)
	assert.Equal(t, PluginType("script"), PluginTypeScript)
}

func TestPluginStates(t *testing.T) {
	assert.Equal(t, PluginState("unloaded"), PluginStateUnloaded)
	assert.Equal(t, PluginState("loading"), PluginStateLoading)
	assert.Equal(t, PluginState("active"), PluginStateActive)
	assert.Equal(t, PluginState("stopped"), PluginStateStopped)
	assert.Equal(t, PluginState("error"), PluginStateError)
}

func TestPluginConfig(t *testing.T) {
	cfg := PluginConfig{
		Enabled:  true,
		Priority: 50,
		Settings: map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		},
	}

	assert.True(t, cfg.Enabled)
	assert.Equal(t, 50, cfg.Priority)
	assert.Equal(t, "value1", cfg.Settings["key1"])
	assert.Equal(t, 42, cfg.Settings["key2"])
}

func TestHookContext(t *testing.T) {
	now := time.Now()
	ctx := HookContext{
		HookType:   HookPreCommand,
		Command:    "ls -la",
		Message:    "message",
		Input:      "input",
		Output:     "output",
		KeyPressed: "ctrl+c",
		SessionID:  "session-1",
		Metadata:   map[string]interface{}{"key": "value"},
		Timestamp:  now,
	}

	assert.Equal(t, HookPreCommand, ctx.HookType)
	assert.Equal(t, "ls -la", ctx.Command)
	assert.Equal(t, "message", ctx.Message)
	assert.Equal(t, "input", ctx.Input)
	assert.Equal(t, "output", ctx.Output)
	assert.Equal(t, "ctrl+c", ctx.KeyPressed)
	assert.Equal(t, "session-1", ctx.SessionID)
	assert.Equal(t, "value", ctx.Metadata["key"])
	assert.Equal(t, now, ctx.Timestamp)
}

func TestHookResult(t *testing.T) {
	result := HookResult{
		Modified:   true,
		Cancel:     false,
		NewCommand: "new command",
		NewMessage: "new message",
		Metadata:   map[string]interface{}{"key": "value"},
		Error:      "",
	}

	assert.True(t, result.Modified)
	assert.False(t, result.Cancel)
	assert.Equal(t, "new command", result.NewCommand)
	assert.Equal(t, "new message", result.NewMessage)
	assert.Equal(t, "value", result.Metadata["key"])
}

func TestPluginManifest(t *testing.T) {
	manifest := PluginManifest{
		Name:        "Test Plugin",
		Version:     "1.0.0",
		Description: "A test plugin",
		Author:      "Test Author",
		Type:        PluginTypeNative,
		Hooks:       []HookType{HookPreCommand, HookPostCommand},
		Main:        "main.go",
		Settings: map[string]SettingDef{
			"enabled": {
				Type:        "boolean",
				Default:     true,
				Description: "Enable the plugin",
				Required:    false,
			},
		},
	}

	assert.Equal(t, "Test Plugin", manifest.Name)
	assert.Equal(t, "1.0.0", manifest.Version)
	assert.Equal(t, PluginTypeNative, manifest.Type)
	assert.Len(t, manifest.Hooks, 2)
	assert.Equal(t, "main.go", manifest.Main)
	assert.Equal(t, "boolean", manifest.Settings["enabled"].Type)
}

func TestSettingDef(t *testing.T) {
	def := SettingDef{
		Type:        "string",
		Default:     "default_value",
		Description: "A string setting",
		Required:    true,
	}

	assert.Equal(t, "string", def.Type)
	assert.Equal(t, "default_value", def.Default)
	assert.Equal(t, "A string setting", def.Description)
	assert.True(t, def.Required)
}
