package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// HookType defines when a plugin hook is triggered
type HookType string

const (
	// HookPreCommand is triggered before a command is executed
	HookPreCommand HookType = "pre_command"

	// HookPostCommand is triggered after a command is executed
	HookPostCommand HookType = "post_command"

	// HookPreAI is triggered before an AI request
	HookPreAI HookType = "pre_ai"

	// HookPostAI is triggered after an AI response
	HookPostAI HookType = "post_ai"

	// HookOnKeyPress is triggered on key press
	HookOnKeyPress HookType = "on_key_press"

	// HookOnMessage is triggered when a message is added
	HookOnMessage HookType = "on_message"

	// HookOnSessionStart is triggered when a session starts
	HookOnSessionStart HookType = "on_session_start"

	// HookOnSessionEnd is triggered when a session ends
	HookOnSessionEnd HookType = "on_session_end"
)

// PluginType defines the type of plugin
type PluginType string

const (
	// PluginTypeNative is a native Go plugin
	PluginTypeNative PluginType = "native"

	// PluginTypeWASM is a WebAssembly plugin
	PluginTypeWASM PluginType = "wasm"

	// PluginTypeScript is a script-based plugin (lua, python, etc.)
	PluginTypeScript PluginType = "script"
)

// PluginState represents the state of a plugin
type PluginState string

const (
	PluginStateUnloaded PluginState = "unloaded"
	PluginStateLoading  PluginState = "loading"
	PluginStateActive   PluginState = "active"
	PluginStateStopped  PluginState = "stopped"
	PluginStateError    PluginState = "error"
)

// Plugin represents a plugin
type Plugin struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Version     string      `json:"version"`
	Description string      `json:"description,omitempty"`
	Author      string      `json:"author,omitempty"`
	Type        PluginType  `json:"type"`
	Hooks       []HookType  `json:"hooks"`
	Config      PluginConfig `json:"config,omitempty"`
	State       PluginState `json:"state"`
	LoadedAt    *time.Time  `json:"loaded_at,omitempty"`
	Path        string      `json:"path,omitempty"`
	Error       string      `json:"error,omitempty"`
}

// PluginConfig holds plugin configuration
type PluginConfig struct {
	Enabled     bool                   `json:"enabled"`
	Priority    int                    `json:"priority"` // Lower = higher priority
	Settings    map[string]interface{} `json:"settings,omitempty"`
}

// PluginManifest represents the plugin.json manifest
type PluginManifest struct {
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Description string                 `json:"description,omitempty"`
	Author      string                 `json:"author,omitempty"`
	Type        PluginType             `json:"type"`
	Hooks       []HookType             `json:"hooks"`
	Main        string                 `json:"main"` // Entry point file
	Settings    map[string]SettingDef  `json:"settings,omitempty"`
}

// SettingDef defines a plugin setting
type SettingDef struct {
	Type        string      `json:"type"` // string, number, boolean
	Default     interface{} `json:"default,omitempty"`
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required,omitempty"`
}

// HookContext provides context to plugin hooks
type HookContext struct {
	HookType    HookType               `json:"hook_type"`
	Command     string                 `json:"command,omitempty"`
	Message     string                 `json:"message,omitempty"`
	Input       string                 `json:"input,omitempty"`
	Output      string                 `json:"output,omitempty"`
	KeyPressed  string                 `json:"key_pressed,omitempty"`
	SessionID   string                 `json:"session_id,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
}

// HookResult represents the result of a hook execution
type HookResult struct {
	Modified    bool                   `json:"modified"`
	Cancel      bool                   `json:"cancel"`      // Cancel the operation
	NewCommand  string                 `json:"new_command,omitempty"`
	NewMessage  string                 `json:"new_message,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Error       string                 `json:"error,omitempty"`
}

// PluginHandler is the interface that plugins must implement
type PluginHandler interface {
	// Initialize is called when the plugin is loaded
	Initialize(config map[string]interface{}) error

	// Execute is called when a hook is triggered
	Execute(ctx context.Context, hookCtx *HookContext) (*HookResult, error)

	// Shutdown is called when the plugin is unloaded
	Shutdown() error

	// GetInfo returns plugin information
	GetInfo() *Plugin
}

// PluginManager manages all plugins
type PluginManager struct {
	plugins     map[string]*Plugin
	handlers    map[string]PluginHandler
	pluginDir   string
	mu          sync.RWMutex

	// Hook chains
	hooks map[HookType][]string // Hook type -> ordered plugin IDs
}

// NewPluginManager creates a new plugin manager
func NewPluginManager(pluginDir string) *PluginManager {
	return &PluginManager{
		plugins:   make(map[string]*Plugin),
		handlers:  make(map[string]PluginHandler),
		pluginDir: pluginDir,
		hooks:     make(map[HookType][]string),
	}
}

// GetPluginDir returns the plugin directory
func (pm *PluginManager) GetPluginDir() string {
	return pm.pluginDir
}

// LoadPlugins loads all plugins from the plugin directory
func (pm *PluginManager) LoadPlugins() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if pm.pluginDir == "" {
		return nil
	}

	// Create plugin directory if it doesn't exist
	if err := os.MkdirAll(pm.pluginDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory: %w", err)
	}

	entries, err := os.ReadDir(pm.pluginDir)
	if err != nil {
		return fmt.Errorf("failed to read plugin directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pluginPath := filepath.Join(pm.pluginDir, entry.Name())
		manifestPath := filepath.Join(pluginPath, "plugin.json")

		// Check if manifest exists
		if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
			continue
		}

		// Load manifest
		manifest, err := pm.loadManifest(manifestPath)
		if err != nil {
			fmt.Printf("Warning: Failed to load plugin %s: %v\n", entry.Name(), err)
			continue
		}

		// Create plugin instance
		plugin := &Plugin{
			ID:          entry.Name(),
			Name:        manifest.Name,
			Version:     manifest.Version,
			Description: manifest.Description,
			Author:      manifest.Author,
			Type:        manifest.Type,
			Hooks:       manifest.Hooks,
			Path:        pluginPath,
			State:       PluginStateUnloaded,
			Config: PluginConfig{
				Enabled:  true,
				Priority: 100,
				Settings: make(map[string]interface{}),
			},
		}

		// Set default settings
		for name, def := range manifest.Settings {
			if def.Default != nil {
				plugin.Config.Settings[name] = def.Default
			}
		}

		pm.plugins[plugin.ID] = plugin
	}

	return nil
}

// loadManifest loads a plugin manifest
func (pm *PluginManager) loadManifest(path string) (*PluginManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest PluginManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

// EnablePlugin enables a plugin
func (pm *PluginManager) EnablePlugin(id string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	plugin, ok := pm.plugins[id]
	if !ok {
		return fmt.Errorf("plugin not found: %s", id)
	}

	plugin.Config.Enabled = true
	now := time.Now()
	plugin.LoadedAt = &now
	plugin.State = PluginStateActive

	// Register hooks
	for _, hook := range plugin.Hooks {
		pm.registerHook(hook, plugin.ID)
	}

	return nil
}

// DisablePlugin disables a plugin
func (pm *PluginManager) DisablePlugin(id string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	plugin, ok := pm.plugins[id]
	if !ok {
		return fmt.Errorf("plugin not found: %s", id)
	}

	plugin.Config.Enabled = false
	plugin.State = PluginStateStopped

	// Unregister hooks
	for _, hook := range plugin.Hooks {
		pm.unregisterHook(hook, plugin.ID)
	}

	// Shutdown handler if exists
	if handler, ok := pm.handlers[id]; ok {
		handler.Shutdown()
		delete(pm.handlers, id)
	}

	return nil
}

// registerHook registers a plugin for a hook
func (pm *PluginManager) registerHook(hookType HookType, pluginID string) {
	hooks := pm.hooks[hookType]

	// Check if already registered
	for _, id := range hooks {
		if id == pluginID {
			return
		}
	}

	pm.hooks[hookType] = append(hooks, pluginID)

	// Sort by priority
	pm.sortHooksByPriority(hookType)
}

// unregisterHook unregisters a plugin from a hook
func (pm *PluginManager) unregisterHook(hookType HookType, pluginID string) {
	hooks := pm.hooks[hookType]
	newHooks := make([]string, 0, len(hooks))

	for _, id := range hooks {
		if id != pluginID {
			newHooks = append(newHooks, id)
		}
	}

	pm.hooks[hookType] = newHooks
}

// sortHooksByPriority sorts hook registrations by plugin priority
func (pm *PluginManager) sortHooksByPriority(hookType HookType) {
	hooks := pm.hooks[hookType]

	// Simple bubble sort for small lists
	for i := 0; i < len(hooks)-1; i++ {
		for j := 0; j < len(hooks)-i-1; j++ {
			p1 := pm.plugins[hooks[j]]
			p2 := pm.plugins[hooks[j+1]]
			if p1 != nil && p2 != nil && p1.Config.Priority > p2.Config.Priority {
				hooks[j], hooks[j+1] = hooks[j+1], hooks[j]
			}
		}
	}
}

// ExecuteHook executes all plugins for a hook
func (pm *PluginManager) ExecuteHook(ctx context.Context, hookCtx *HookContext) (*HookResult, error) {
	pm.mu.RLock()
	pluginIDs := pm.hooks[hookCtx.HookType]
	pm.mu.RUnlock()

	result := &HookResult{}

	for _, pluginID := range pluginIDs {
		pm.mu.RLock()
		plugin := pm.plugins[pluginID]
		handler := pm.handlers[pluginID]
		pm.mu.RUnlock()

		if plugin == nil || !plugin.Config.Enabled {
			continue
		}

		// Execute handler if available
		if handler != nil {
			hookResult, err := handler.Execute(ctx, hookCtx)
			if err != nil {
				plugin.Error = err.Error()
				continue
			}

			// Merge results
			if hookResult.Cancel {
				result.Cancel = true
				return result, nil
			}
			if hookResult.Modified {
				result.Modified = true
				if hookResult.NewCommand != "" {
					result.NewCommand = hookResult.NewCommand
					hookCtx.Command = hookResult.NewCommand
				}
				if hookResult.NewMessage != "" {
					result.NewMessage = hookResult.NewMessage
					hookCtx.Message = hookResult.NewMessage
				}
			}
		}
	}

	return result, nil
}

// RegisterHandler registers a native plugin handler
func (pm *PluginManager) RegisterHandler(id string, handler PluginHandler) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	plugin := handler.GetInfo()
	plugin.ID = id
	plugin.Type = PluginTypeNative
	plugin.State = PluginStateActive

	pm.plugins[id] = plugin
	pm.handlers[id] = handler

	// Register hooks
	for _, hook := range plugin.Hooks {
		pm.registerHook(hook, id)
	}

	return nil
}

// GetPlugin returns a plugin by ID
func (pm *PluginManager) GetPlugin(id string) *Plugin {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.plugins[id]
}

// ListPlugins returns all plugins
func (pm *PluginManager) ListPlugins() []*Plugin {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	plugins := make([]*Plugin, 0, len(pm.plugins))
	for _, p := range pm.plugins {
		plugins = append(plugins, p)
	}
	return plugins
}

// ListActivePlugins returns all active plugins
func (pm *PluginManager) ListActivePlugins() []*Plugin {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var active []*Plugin
	for _, p := range pm.plugins {
		if p.State == PluginStateActive {
			active = append(active, p)
		}
	}
	return active
}

// SetPluginSetting sets a plugin setting
func (pm *PluginManager) SetPluginSetting(pluginID, key string, value interface{}) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	plugin, ok := pm.plugins[pluginID]
	if !ok {
		return fmt.Errorf("plugin not found: %s", pluginID)
	}

	if plugin.Config.Settings == nil {
		plugin.Config.Settings = make(map[string]interface{})
	}
	plugin.Config.Settings[key] = value

	return nil
}

// GetPluginSetting gets a plugin setting
func (pm *PluginManager) GetPluginSetting(pluginID, key string) (interface{}, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	plugin, ok := pm.plugins[pluginID]
	if !ok || plugin.Config.Settings == nil {
		return nil, false
	}

	val, ok := plugin.Config.Settings[key]
	return val, ok
}

// SetPluginPriority sets a plugin's priority
func (pm *PluginManager) SetPluginPriority(pluginID string, priority int) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	plugin, ok := pm.plugins[pluginID]
	if !ok {
		return fmt.Errorf("plugin not found: %s", pluginID)
	}

	plugin.Config.Priority = priority

	// Re-sort hooks
	for _, hook := range plugin.Hooks {
		pm.sortHooksByPriority(hook)
	}

	return nil
}

// UnloadAll unloads all plugins
func (pm *PluginManager) UnloadAll() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for id, handler := range pm.handlers {
		handler.Shutdown()
		delete(pm.handlers, id)
	}

	for _, plugin := range pm.plugins {
		plugin.State = PluginStateUnloaded
	}

	pm.hooks = make(map[HookType][]string)
}
