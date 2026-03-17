package clipboard

import (
	"errors"
	"sync"

	"golang.design/x/clipboard"
)

var (
	initOnce sync.Once
	initErr  error
)

// Manager handles clipboard operations
type Manager struct {
	initialized bool
	lastCopied  string
}

// New creates a new clipboard manager
func New() (*Manager, error) {
	// Initialize clipboard (only once)
	initOnce.Do(func() {
		initErr = clipboard.Init()
	})

	if initErr != nil {
		return nil, initErr
	}

	return &Manager{
		initialized: true,
	}, nil
}

// Copy copies text to the clipboard
func (m *Manager) Copy(text string) error {
	if !m.initialized {
		return errors.New("clipboard not initialized")
	}

	if text == "" {
		return errors.New("cannot copy empty text")
	}

	clipboard.Write(clipboard.FmtText, []byte(text))
	m.lastCopied = text

	return nil
}

// Read reads text from the clipboard
func (m *Manager) Read() (string, error) {
	if !m.initialized {
		return "", errors.New("clipboard not initialized")
	}

	data := clipboard.Read(clipboard.FmtText)
	if data == nil {
		return "", nil
	}

	return string(data), nil
}

// LastCopied returns the last copied text
func (m *Manager) LastCopied() string {
	return m.lastCopied
}

// IsAvailable checks if clipboard is available
func (m *Manager) IsAvailable() bool {
	return m.initialized
}

// CopyCommand copies a command with optional prefix removal
func (m *Manager) CopyCommand(command string) error {
	// Remove common prefixes
	cmd := command

	// Remove $ or > prefixes
	if len(cmd) > 2 && (cmd[0] == '$' || cmd[0] == '>') && cmd[1] == ' ' {
		cmd = cmd[2:]
	}

	return m.Copy(cmd)
}
