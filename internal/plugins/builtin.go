package plugins

import (
	"context"
	"regexp"
	"strings"
	"time"
)

// SafetyPlugin prevents dangerous commands
type SafetyPlugin struct {
	enabled         bool
	dangerousPatterns []string
	config          map[string]interface{}
}

// NewSafetyPlugin creates a new safety plugin
func NewSafetyPlugin() *SafetyPlugin {
	return &SafetyPlugin{
		enabled: true,
		dangerousPatterns: []string{
			`rm\s+-rf\s+/`,
			`rm\s+-rf\s+\*`,
			`mkfs\.`,
			`dd\s+if=`,
			`:()\{\s*:\|:&\s*\};:`,
			`>\s*/dev/sda`,
			`chmod\s+-R\s+777\s+/`,
			`curl.*\|\s*sh`,
			`wget.*\|\s*sh`,
		},
	}
}

// Initialize initializes the plugin
func (p *SafetyPlugin) Initialize(config map[string]interface{}) error {
	p.config = config
	if enabled, ok := config["enabled"].(bool); ok {
		p.enabled = enabled
	}
	return nil
}

// Execute executes the plugin hook
func (p *SafetyPlugin) Execute(ctx context.Context, hookCtx *HookContext) (*HookResult, error) {
	if !p.enabled || hookCtx.HookType != HookPreCommand {
		return &HookResult{}, nil
	}

	command := hookCtx.Command

	for _, pattern := range p.dangerousPatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}

		if re.MatchString(strings.ToLower(command)) {
			return &HookResult{
				Cancel:     true,
				Error:      "Blocked: This command is potentially dangerous",
				Metadata: map[string]interface{}{
					"blocked_pattern": pattern,
				},
			}, nil
		}
	}

	return &HookResult{}, nil
}

// Shutdown shuts down the plugin
func (p *SafetyPlugin) Shutdown() error {
	return nil
}

// GetInfo returns plugin information
func (p *SafetyPlugin) GetInfo() *Plugin {
	return &Plugin{
		Name:        "Safety Guard",
		Version:     "1.0.0",
		Description: "Prevents execution of dangerous commands",
		Author:      "TerminalizCrazy",
		Type:        PluginTypeNative,
		Hooks:       []HookType{HookPreCommand},
		Config: PluginConfig{
			Enabled:  p.enabled,
			Priority: 1, // High priority - run first
		},
	}
}

// TimestampPlugin adds timestamps to outputs
type TimestampPlugin struct {
	enabled bool
	format  string
	config  map[string]interface{}
}

// NewTimestampPlugin creates a new timestamp plugin
func NewTimestampPlugin() *TimestampPlugin {
	return &TimestampPlugin{
		enabled: true,
		format:  "2006-01-02 15:04:05",
	}
}

// Initialize initializes the plugin
func (p *TimestampPlugin) Initialize(config map[string]interface{}) error {
	p.config = config
	if enabled, ok := config["enabled"].(bool); ok {
		p.enabled = enabled
	}
	if format, ok := config["format"].(string); ok {
		p.format = format
	}
	return nil
}

// Execute executes the plugin hook
func (p *TimestampPlugin) Execute(ctx context.Context, hookCtx *HookContext) (*HookResult, error) {
	if !p.enabled || hookCtx.HookType != HookPostCommand {
		return &HookResult{}, nil
	}

	timestamp := time.Now().Format(p.format)
	newOutput := hookCtx.Output + "\n[" + timestamp + "]"

	return &HookResult{
		Modified:   true,
		NewMessage: newOutput,
	}, nil
}

// Shutdown shuts down the plugin
func (p *TimestampPlugin) Shutdown() error {
	return nil
}

// GetInfo returns plugin information
func (p *TimestampPlugin) GetInfo() *Plugin {
	return &Plugin{
		Name:        "Timestamp",
		Version:     "1.0.0",
		Description: "Adds timestamps to command outputs",
		Author:      "TerminalizCrazy",
		Type:        PluginTypeNative,
		Hooks:       []HookType{HookPostCommand},
		Config: PluginConfig{
			Enabled:  p.enabled,
			Priority: 100,
		},
	}
}

// AliasPlugin provides command aliases
type AliasPlugin struct {
	enabled bool
	aliases map[string]string
	config  map[string]interface{}
}

// NewAliasPlugin creates a new alias plugin
func NewAliasPlugin() *AliasPlugin {
	return &AliasPlugin{
		enabled: true,
		aliases: map[string]string{
			"ll":     "ls -la",
			"la":     "ls -la",
			"..":     "cd ..",
			"...":    "cd ../..",
			"gs":     "git status",
			"gd":     "git diff",
			"gp":     "git push",
			"gl":     "git log --oneline -10",
			"gc":     "git commit",
			"gco":    "git checkout",
			"gbr":    "git branch",
			"grh":    "git reset --hard HEAD",
			"cls":    "clear",
			"clr":    "clear",
			"k":      "kubectl",
			"d":      "docker",
			"dc":     "docker-compose",
			"tf":     "terraform",
			"py":     "python3",
			"python": "python3",
			"pip":    "pip3",
		},
	}
}

// Initialize initializes the plugin
func (p *AliasPlugin) Initialize(config map[string]interface{}) error {
	p.config = config
	if enabled, ok := config["enabled"].(bool); ok {
		p.enabled = enabled
	}
	if aliases, ok := config["aliases"].(map[string]interface{}); ok {
		for k, v := range aliases {
			if vs, ok := v.(string); ok {
				p.aliases[k] = vs
			}
		}
	}
	return nil
}

// Execute executes the plugin hook
func (p *AliasPlugin) Execute(ctx context.Context, hookCtx *HookContext) (*HookResult, error) {
	if !p.enabled || hookCtx.HookType != HookPreCommand {
		return &HookResult{}, nil
	}

	command := hookCtx.Command
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return &HookResult{}, nil
	}

	// Check if first word is an alias
	if replacement, ok := p.aliases[parts[0]]; ok {
		newCommand := replacement
		if len(parts) > 1 {
			newCommand += " " + strings.Join(parts[1:], " ")
		}

		return &HookResult{
			Modified:   true,
			NewCommand: newCommand,
		}, nil
	}

	return &HookResult{}, nil
}

// Shutdown shuts down the plugin
func (p *AliasPlugin) Shutdown() error {
	return nil
}

// GetInfo returns plugin information
func (p *AliasPlugin) GetInfo() *Plugin {
	return &Plugin{
		Name:        "Aliases",
		Version:     "1.0.0",
		Description: "Provides command aliases and shortcuts",
		Author:      "TerminalizCrazy",
		Type:        PluginTypeNative,
		Hooks:       []HookType{HookPreCommand},
		Config: PluginConfig{
			Enabled:  p.enabled,
			Priority: 10, // Run early but after safety
		},
	}
}

// AddAlias adds a custom alias
func (p *AliasPlugin) AddAlias(name, command string) {
	p.aliases[name] = command
}

// RemoveAlias removes an alias
func (p *AliasPlugin) RemoveAlias(name string) {
	delete(p.aliases, name)
}

// GetAliases returns all aliases
func (p *AliasPlugin) GetAliases() map[string]string {
	result := make(map[string]string)
	for k, v := range p.aliases {
		result[k] = v
	}
	return result
}

// HistoryLoggerPlugin logs command history
type HistoryLoggerPlugin struct {
	enabled  bool
	history  []CommandRecord
	maxItems int
	config   map[string]interface{}
}

// CommandRecord represents a command in history
type CommandRecord struct {
	Command   string
	Output    string
	ExitCode  int
	Timestamp time.Time
	Duration  time.Duration
}

// NewHistoryLoggerPlugin creates a new history logger plugin
func NewHistoryLoggerPlugin() *HistoryLoggerPlugin {
	return &HistoryLoggerPlugin{
		enabled:  true,
		history:  make([]CommandRecord, 0),
		maxItems: 1000,
	}
}

// Initialize initializes the plugin
func (p *HistoryLoggerPlugin) Initialize(config map[string]interface{}) error {
	p.config = config
	if enabled, ok := config["enabled"].(bool); ok {
		p.enabled = enabled
	}
	if maxItems, ok := config["max_items"].(int); ok {
		p.maxItems = maxItems
	}
	return nil
}

// Execute executes the plugin hook
func (p *HistoryLoggerPlugin) Execute(ctx context.Context, hookCtx *HookContext) (*HookResult, error) {
	if !p.enabled {
		return &HookResult{}, nil
	}

	switch hookCtx.HookType {
	case HookPostCommand:
		record := CommandRecord{
			Command:   hookCtx.Command,
			Output:    hookCtx.Output,
			Timestamp: hookCtx.Timestamp,
		}
		if exitCode, ok := hookCtx.Metadata["exit_code"].(int); ok {
			record.ExitCode = exitCode
		}
		if duration, ok := hookCtx.Metadata["duration"].(time.Duration); ok {
			record.Duration = duration
		}

		p.history = append(p.history, record)

		// Trim history
		if len(p.history) > p.maxItems {
			p.history = p.history[len(p.history)-p.maxItems:]
		}
	}

	return &HookResult{}, nil
}

// Shutdown shuts down the plugin
func (p *HistoryLoggerPlugin) Shutdown() error {
	return nil
}

// GetInfo returns plugin information
func (p *HistoryLoggerPlugin) GetInfo() *Plugin {
	return &Plugin{
		Name:        "History Logger",
		Version:     "1.0.0",
		Description: "Logs command history for analysis",
		Author:      "TerminalizCrazy",
		Type:        PluginTypeNative,
		Hooks:       []HookType{HookPostCommand},
		Config: PluginConfig{
			Enabled:  p.enabled,
			Priority: 200, // Run later
		},
	}
}

// GetHistory returns the command history
func (p *HistoryLoggerPlugin) GetHistory() []CommandRecord {
	return p.history
}

// SearchHistory searches the history for commands matching a pattern
func (p *HistoryLoggerPlugin) SearchHistory(pattern string) []CommandRecord {
	var results []CommandRecord
	re, err := regexp.Compile(pattern)
	if err != nil {
		return results
	}

	for _, record := range p.history {
		if re.MatchString(record.Command) {
			results = append(results, record)
		}
	}

	return results
}

// RegisterBuiltInPlugins registers all built-in plugins with a manager
func RegisterBuiltInPlugins(pm *PluginManager) error {
	plugins := []PluginHandler{
		NewSafetyPlugin(),
		NewAliasPlugin(),
		NewTimestampPlugin(),
		NewHistoryLoggerPlugin(),
	}

	for i, handler := range plugins {
		id := handler.GetInfo().Name
		id = strings.ToLower(strings.ReplaceAll(id, " ", "-"))
		if err := pm.RegisterHandler(id+"-"+string(rune('0'+i)), handler); err != nil {
			return err
		}
	}

	return nil
}
