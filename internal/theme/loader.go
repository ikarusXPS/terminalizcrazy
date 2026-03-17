package theme

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

// Manager handles theme loading, switching, and hot-reload
type Manager struct {
	currentTheme *Theme
	themes       map[string]*Theme
	themesDir    string
	watcher      *fsnotify.Watcher
	onChange     func(*Theme)
	mu           sync.RWMutex
	watching     bool
}

// NewManager creates a new theme manager
func NewManager(themesDir string) (*Manager, error) {
	m := &Manager{
		themes:    make(map[string]*Theme),
		themesDir: themesDir,
	}

	// Load built-in themes
	for name, theme := range BuiltinThemes() {
		m.themes[name] = theme
	}

	// Set default theme
	m.currentTheme = DefaultTheme()
	m.themes["default"] = m.currentTheme

	// Create themes directory if it doesn't exist
	if themesDir != "" {
		if err := os.MkdirAll(themesDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create themes directory: %w", err)
		}

		// Load custom themes from directory
		if err := m.loadThemesFromDir(themesDir); err != nil {
			// Log but don't fail - custom themes are optional
			fmt.Printf("Warning: failed to load custom themes: %v\n", err)
		}
	}

	return m, nil
}

// loadThemesFromDir loads all YAML theme files from a directory
func (m *Manager) loadThemesFromDir(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return err
	}

	ymlFiles, err := filepath.Glob(filepath.Join(dir, "*.yml"))
	if err != nil {
		return err
	}
	files = append(files, ymlFiles...)

	for _, file := range files {
		theme, err := m.loadThemeFromFile(file)
		if err != nil {
			fmt.Printf("Warning: failed to load theme %s: %v\n", file, err)
			continue
		}
		m.themes[normalizeThemeName(theme.Name)] = theme
	}

	return nil
}

// loadThemeFromFile loads a theme from a YAML file
func (m *Manager) loadThemeFromFile(path string) (*Theme, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read theme file: %w", err)
	}

	var theme Theme
	if err := yaml.Unmarshal(data, &theme); err != nil {
		return nil, fmt.Errorf("failed to parse theme file: %w", err)
	}

	// Apply defaults for missing colors
	theme.ApplyDefaults()

	// Validate the theme
	if err := theme.Validate(); err != nil {
		return nil, fmt.Errorf("invalid theme: %w", err)
	}

	return &theme, nil
}

// StartWatching starts watching the themes directory for changes
func (m *Manager) StartWatching() error {
	if m.themesDir == "" {
		return fmt.Errorf("no themes directory configured")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}

	m.mu.Lock()
	m.watcher = watcher
	m.watching = true
	m.mu.Unlock()

	go m.watchLoop()

	if err := watcher.Add(m.themesDir); err != nil {
		return fmt.Errorf("failed to watch directory: %w", err)
	}

	return nil
}

// watchLoop handles file system events
func (m *Manager) watchLoop() {
	for {
		m.mu.RLock()
		watcher := m.watcher
		watching := m.watching
		m.mu.RUnlock()

		if !watching || watcher == nil {
			return
		}

		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Only handle write and create events for YAML files
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				ext := filepath.Ext(event.Name)
				if ext == ".yaml" || ext == ".yml" {
					m.handleThemeFileChange(event.Name)
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("Theme watcher error: %v\n", err)
		}
	}
}

// handleThemeFileChange reloads a theme file when it changes
func (m *Manager) handleThemeFileChange(path string) {
	theme, err := m.loadThemeFromFile(path)
	if err != nil {
		fmt.Printf("Warning: failed to reload theme %s: %v\n", path, err)
		return
	}

	m.mu.Lock()
	name := normalizeThemeName(theme.Name)
	m.themes[name] = theme

	// If this is the current theme, update it
	currentName := normalizeThemeName(m.currentTheme.Name)
	isCurrentTheme := currentName == name
	m.mu.Unlock()

	if isCurrentTheme {
		m.SetTheme(name)
	}
}

// StopWatching stops watching the themes directory
func (m *Manager) StopWatching() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.watching = false
	if m.watcher != nil {
		err := m.watcher.Close()
		m.watcher = nil
		return err
	}
	return nil
}

// SetTheme sets the current theme by name
func (m *Manager) SetTheme(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	normalized := normalizeThemeName(name)
	theme, ok := m.themes[normalized]
	if !ok {
		return fmt.Errorf("theme not found: %s", name)
	}

	m.currentTheme = theme

	// Notify listeners
	if m.onChange != nil {
		m.onChange(theme)
	}

	return nil
}

// CurrentTheme returns the current theme
func (m *Manager) CurrentTheme() *Theme {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentTheme
}

// GetTheme returns a theme by name
func (m *Manager) GetTheme(name string) *Theme {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.themes[normalizeThemeName(name)]
}

// ListThemes returns all available theme names
func (m *Manager) ListThemes() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.themes))
	for name := range m.themes {
		names = append(names, name)
	}
	return names
}

// OnChange sets a callback to be called when the theme changes
func (m *Manager) OnChange(fn func(*Theme)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onChange = fn
}

// SaveTheme saves a theme to a YAML file
func (m *Manager) SaveTheme(theme *Theme) error {
	if m.themesDir == "" {
		return fmt.Errorf("no themes directory configured")
	}

	filename := normalizeThemeName(theme.Name) + ".yaml"
	path := filepath.Join(m.themesDir, filename)

	data, err := yaml.Marshal(theme)
	if err != nil {
		return fmt.Errorf("failed to marshal theme: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write theme file: %w", err)
	}

	// Add to themes map
	m.mu.Lock()
	m.themes[normalizeThemeName(theme.Name)] = theme
	m.mu.Unlock()

	return nil
}

// RegisterTheme adds a theme to the manager
func (m *Manager) RegisterTheme(theme *Theme) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.themes[normalizeThemeName(theme.Name)] = theme
}

// normalizeThemeName converts a theme name to lowercase for consistent lookup
func normalizeThemeName(name string) string {
	// Convert to lowercase and replace spaces with hyphens
	result := ""
	for _, r := range name {
		if r == ' ' {
			result += "-"
		} else if r >= 'A' && r <= 'Z' {
			result += string(r + 32)
		} else {
			result += string(r)
		}
	}
	return result
}

// Close cleans up the manager
func (m *Manager) Close() error {
	return m.StopWatching()
}
