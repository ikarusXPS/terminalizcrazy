package plugins

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SafetyPlugin tests

func TestNewSafetyPlugin(t *testing.T) {
	plugin := NewSafetyPlugin()

	assert.NotNil(t, plugin)
	assert.True(t, plugin.enabled)
	assert.NotEmpty(t, plugin.dangerousPatterns)
}

func TestSafetyPlugin_Initialize(t *testing.T) {
	plugin := NewSafetyPlugin()

	err := plugin.Initialize(map[string]interface{}{
		"enabled": false,
	})

	assert.NoError(t, err)
	assert.False(t, plugin.enabled)
}

func TestSafetyPlugin_Execute_BlocksDangerous(t *testing.T) {
	plugin := NewSafetyPlugin()
	plugin.Initialize(nil)

	tests := []struct {
		name    string
		command string
		blocked bool
	}{
		{"rm -rf /", "rm -rf /", true},
		{"rm -rf *", "rm -rf *", true},
		{"mkfs.ext4", "mkfs.ext4 /dev/sda", true},
		{"dd if=", "dd if=/dev/zero of=/dev/sda", true},
		{"curl pipe to sh", "curl http://example.com | sh", true},
		{"wget pipe to sh", "wget http://example.com -O - | sh", true},
		{"safe command", "ls -la", false},
		{"git status", "git status", false},
		{"cat file", "cat /etc/passwd", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hookCtx := &HookContext{
				HookType: HookPreCommand,
				Command:  tt.command,
			}

			result, err := plugin.Execute(context.Background(), hookCtx)

			assert.NoError(t, err)
			assert.Equal(t, tt.blocked, result.Cancel)
		})
	}
}

func TestSafetyPlugin_Execute_DisabledDoesNotBlock(t *testing.T) {
	plugin := NewSafetyPlugin()
	plugin.Initialize(map[string]interface{}{"enabled": false})

	hookCtx := &HookContext{
		HookType: HookPreCommand,
		Command:  "rm -rf /",
	}

	result, err := plugin.Execute(context.Background(), hookCtx)

	assert.NoError(t, err)
	assert.False(t, result.Cancel)
}

func TestSafetyPlugin_Execute_WrongHookType(t *testing.T) {
	plugin := NewSafetyPlugin()
	plugin.Initialize(nil)

	hookCtx := &HookContext{
		HookType: HookPostCommand,
		Command:  "rm -rf /",
	}

	result, err := plugin.Execute(context.Background(), hookCtx)

	assert.NoError(t, err)
	assert.False(t, result.Cancel)
}

func TestSafetyPlugin_Shutdown(t *testing.T) {
	plugin := NewSafetyPlugin()

	err := plugin.Shutdown()

	assert.NoError(t, err)
}

func TestSafetyPlugin_GetInfo(t *testing.T) {
	plugin := NewSafetyPlugin()

	info := plugin.GetInfo()

	assert.Equal(t, "Safety Guard", info.Name)
	assert.Equal(t, "1.0.0", info.Version)
	assert.Equal(t, PluginTypeNative, info.Type)
	assert.Contains(t, info.Hooks, HookPreCommand)
	assert.Equal(t, 1, info.Config.Priority)
}

// TimestampPlugin tests

func TestNewTimestampPlugin(t *testing.T) {
	plugin := NewTimestampPlugin()

	assert.NotNil(t, plugin)
	assert.True(t, plugin.enabled)
	assert.Equal(t, "2006-01-02 15:04:05", plugin.format)
}

func TestTimestampPlugin_Initialize(t *testing.T) {
	plugin := NewTimestampPlugin()

	err := plugin.Initialize(map[string]interface{}{
		"enabled": false,
		"format":  "15:04:05",
	})

	assert.NoError(t, err)
	assert.False(t, plugin.enabled)
	assert.Equal(t, "15:04:05", plugin.format)
}

func TestTimestampPlugin_Execute(t *testing.T) {
	plugin := NewTimestampPlugin()
	plugin.Initialize(nil)

	hookCtx := &HookContext{
		HookType: HookPostCommand,
		Output:   "command output",
	}

	result, err := plugin.Execute(context.Background(), hookCtx)

	assert.NoError(t, err)
	assert.True(t, result.Modified)
	assert.Contains(t, result.NewMessage, "command output")
	assert.Contains(t, result.NewMessage, "[")
}

func TestTimestampPlugin_Execute_Disabled(t *testing.T) {
	plugin := NewTimestampPlugin()
	plugin.Initialize(map[string]interface{}{"enabled": false})

	hookCtx := &HookContext{
		HookType: HookPostCommand,
		Output:   "command output",
	}

	result, err := plugin.Execute(context.Background(), hookCtx)

	assert.NoError(t, err)
	assert.False(t, result.Modified)
}

func TestTimestampPlugin_Execute_WrongHookType(t *testing.T) {
	plugin := NewTimestampPlugin()
	plugin.Initialize(nil)

	hookCtx := &HookContext{
		HookType: HookPreCommand,
		Output:   "command output",
	}

	result, err := plugin.Execute(context.Background(), hookCtx)

	assert.NoError(t, err)
	assert.False(t, result.Modified)
}

func TestTimestampPlugin_Shutdown(t *testing.T) {
	plugin := NewTimestampPlugin()

	err := plugin.Shutdown()

	assert.NoError(t, err)
}

func TestTimestampPlugin_GetInfo(t *testing.T) {
	plugin := NewTimestampPlugin()

	info := plugin.GetInfo()

	assert.Equal(t, "Timestamp", info.Name)
	assert.Equal(t, "1.0.0", info.Version)
	assert.Equal(t, PluginTypeNative, info.Type)
	assert.Contains(t, info.Hooks, HookPostCommand)
}

// AliasPlugin tests

func TestNewAliasPlugin(t *testing.T) {
	plugin := NewAliasPlugin()

	assert.NotNil(t, plugin)
	assert.True(t, plugin.enabled)
	assert.NotEmpty(t, plugin.aliases)
}

func TestAliasPlugin_Initialize(t *testing.T) {
	plugin := NewAliasPlugin()

	err := plugin.Initialize(map[string]interface{}{
		"enabled": false,
		"aliases": map[string]interface{}{
			"custom": "custom command",
		},
	})

	assert.NoError(t, err)
	assert.False(t, plugin.enabled)
	assert.Equal(t, "custom command", plugin.aliases["custom"])
}

func TestAliasPlugin_Execute(t *testing.T) {
	plugin := NewAliasPlugin()
	plugin.Initialize(nil)

	tests := []struct {
		name       string
		command    string
		expected   string
		shouldMod  bool
	}{
		{"ll alias", "ll", "ls -la", true},
		{"gs alias", "gs", "git status", true},
		{"gd alias", "gd", "git diff", true},
		{"gp alias", "gp", "git push", true},
		{"k alias", "k get pods", "kubectl get pods", true},
		{"no alias", "echo hello", "", false},
		{"empty command", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hookCtx := &HookContext{
				HookType: HookPreCommand,
				Command:  tt.command,
			}

			result, err := plugin.Execute(context.Background(), hookCtx)

			assert.NoError(t, err)
			assert.Equal(t, tt.shouldMod, result.Modified)
			if tt.shouldMod {
				assert.Equal(t, tt.expected, result.NewCommand)
			}
		})
	}
}

func TestAliasPlugin_Execute_Disabled(t *testing.T) {
	plugin := NewAliasPlugin()
	plugin.Initialize(map[string]interface{}{"enabled": false})

	hookCtx := &HookContext{
		HookType: HookPreCommand,
		Command:  "ll",
	}

	result, err := plugin.Execute(context.Background(), hookCtx)

	assert.NoError(t, err)
	assert.False(t, result.Modified)
}

func TestAliasPlugin_Execute_WrongHookType(t *testing.T) {
	plugin := NewAliasPlugin()
	plugin.Initialize(nil)

	hookCtx := &HookContext{
		HookType: HookPostCommand,
		Command:  "ll",
	}

	result, err := plugin.Execute(context.Background(), hookCtx)

	assert.NoError(t, err)
	assert.False(t, result.Modified)
}

func TestAliasPlugin_AddRemoveAlias(t *testing.T) {
	plugin := NewAliasPlugin()

	plugin.AddAlias("custom", "custom command")
	assert.Equal(t, "custom command", plugin.aliases["custom"])

	plugin.RemoveAlias("custom")
	_, ok := plugin.aliases["custom"]
	assert.False(t, ok)
}

func TestAliasPlugin_GetAliases(t *testing.T) {
	plugin := NewAliasPlugin()

	aliases := plugin.GetAliases()

	assert.NotEmpty(t, aliases)
	assert.Equal(t, "ls -la", aliases["ll"])
}

func TestAliasPlugin_Shutdown(t *testing.T) {
	plugin := NewAliasPlugin()

	err := plugin.Shutdown()

	assert.NoError(t, err)
}

func TestAliasPlugin_GetInfo(t *testing.T) {
	plugin := NewAliasPlugin()

	info := plugin.GetInfo()

	assert.Equal(t, "Aliases", info.Name)
	assert.Equal(t, "1.0.0", info.Version)
	assert.Equal(t, PluginTypeNative, info.Type)
	assert.Contains(t, info.Hooks, HookPreCommand)
	assert.Equal(t, 10, info.Config.Priority)
}

// HistoryLoggerPlugin tests

func TestNewHistoryLoggerPlugin(t *testing.T) {
	plugin := NewHistoryLoggerPlugin()

	assert.NotNil(t, plugin)
	assert.True(t, plugin.enabled)
	assert.Equal(t, 1000, plugin.maxItems)
	assert.Empty(t, plugin.history)
}

func TestHistoryLoggerPlugin_Initialize(t *testing.T) {
	plugin := NewHistoryLoggerPlugin()

	err := plugin.Initialize(map[string]interface{}{
		"enabled":   false,
		"max_items": 500,
	})

	assert.NoError(t, err)
	assert.False(t, plugin.enabled)
	assert.Equal(t, 500, plugin.maxItems)
}

func TestHistoryLoggerPlugin_Execute(t *testing.T) {
	plugin := NewHistoryLoggerPlugin()
	plugin.Initialize(nil)

	hookCtx := &HookContext{
		HookType:  HookPostCommand,
		Command:   "ls -la",
		Output:    "file1\nfile2",
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"exit_code": 0,
			"duration":  time.Millisecond * 100,
		},
	}

	result, err := plugin.Execute(context.Background(), hookCtx)

	assert.NoError(t, err)
	assert.False(t, result.Modified)

	history := plugin.GetHistory()
	assert.Len(t, history, 1)
	assert.Equal(t, "ls -la", history[0].Command)
	assert.Equal(t, "file1\nfile2", history[0].Output)
	assert.Equal(t, 0, history[0].ExitCode)
}

func TestHistoryLoggerPlugin_Execute_Disabled(t *testing.T) {
	plugin := NewHistoryLoggerPlugin()
	plugin.Initialize(map[string]interface{}{"enabled": false})

	hookCtx := &HookContext{
		HookType: HookPostCommand,
		Command:  "ls -la",
	}

	result, err := plugin.Execute(context.Background(), hookCtx)

	assert.NoError(t, err)
	assert.Empty(t, plugin.GetHistory())
	assert.False(t, result.Modified)
}

func TestHistoryLoggerPlugin_TrimHistory(t *testing.T) {
	plugin := NewHistoryLoggerPlugin()
	plugin.Initialize(map[string]interface{}{"max_items": 5})

	for i := 0; i < 10; i++ {
		hookCtx := &HookContext{
			HookType:  HookPostCommand,
			Command:   "cmd",
			Timestamp: time.Now(),
		}
		plugin.Execute(context.Background(), hookCtx)
	}

	history := plugin.GetHistory()
	assert.Len(t, history, 5)
}

func TestHistoryLoggerPlugin_SearchHistory(t *testing.T) {
	plugin := NewHistoryLoggerPlugin()
	plugin.Initialize(nil)

	commands := []string{"git status", "git commit", "ls -la", "git push"}
	for _, cmd := range commands {
		hookCtx := &HookContext{
			HookType:  HookPostCommand,
			Command:   cmd,
			Timestamp: time.Now(),
		}
		plugin.Execute(context.Background(), hookCtx)
	}

	results := plugin.SearchHistory("git")
	assert.Len(t, results, 3)

	results = plugin.SearchHistory("ls")
	assert.Len(t, results, 1)

	results = plugin.SearchHistory("nonexistent")
	assert.Empty(t, results)
}

func TestHistoryLoggerPlugin_SearchHistory_InvalidRegex(t *testing.T) {
	plugin := NewHistoryLoggerPlugin()

	results := plugin.SearchHistory("[invalid")

	assert.Empty(t, results)
}

func TestHistoryLoggerPlugin_Shutdown(t *testing.T) {
	plugin := NewHistoryLoggerPlugin()

	err := plugin.Shutdown()

	assert.NoError(t, err)
}

func TestHistoryLoggerPlugin_GetInfo(t *testing.T) {
	plugin := NewHistoryLoggerPlugin()

	info := plugin.GetInfo()

	assert.Equal(t, "History Logger", info.Name)
	assert.Equal(t, "1.0.0", info.Version)
	assert.Equal(t, PluginTypeNative, info.Type)
	assert.Contains(t, info.Hooks, HookPostCommand)
	assert.Equal(t, 200, info.Config.Priority)
}

func TestCommandRecord(t *testing.T) {
	now := time.Now()
	record := CommandRecord{
		Command:   "ls -la",
		Output:    "output",
		ExitCode:  0,
		Timestamp: now,
		Duration:  time.Second,
	}

	assert.Equal(t, "ls -la", record.Command)
	assert.Equal(t, "output", record.Output)
	assert.Equal(t, 0, record.ExitCode)
	assert.Equal(t, now, record.Timestamp)
	assert.Equal(t, time.Second, record.Duration)
}

// RegisterBuiltInPlugins tests

func TestRegisterBuiltInPlugins(t *testing.T) {
	pm := NewPluginManager("")

	err := RegisterBuiltInPlugins(pm)

	require.NoError(t, err)

	plugins := pm.ListPlugins()
	assert.Len(t, plugins, 4)

	// Check all plugins are active
	for _, p := range plugins {
		assert.Equal(t, PluginStateActive, p.State)
	}
}
