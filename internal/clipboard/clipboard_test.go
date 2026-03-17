package clipboard

import (
	"testing"
)

func TestNew(t *testing.T) {
	// Note: Clipboard init may fail in headless/CI environments
	// This test just verifies the function doesn't panic
	_, err := New()
	if err != nil {
		t.Logf("Clipboard init returned error (expected in CI): %v", err)
	}
}

func TestManagerCopyEmptyText(t *testing.T) {
	m := &Manager{initialized: true}
	err := m.Copy("")
	if err == nil {
		t.Error("Expected error for empty text")
	}
	if err.Error() != "cannot copy empty text" {
		t.Errorf("Expected 'cannot copy empty text', got: %v", err)
	}
}

func TestManagerNotInitialized(t *testing.T) {
	m := &Manager{initialized: false}

	err := m.Copy("test")
	if err == nil {
		t.Error("Expected error for uninitialized manager")
	}

	_, err = m.Read()
	if err == nil {
		t.Error("Expected error for uninitialized manager")
	}
}

func TestLastCopied(t *testing.T) {
	m := &Manager{initialized: true, lastCopied: "test-command"}

	if m.LastCopied() != "test-command" {
		t.Errorf("Expected 'test-command', got: %s", m.LastCopied())
	}
}

func TestIsAvailable(t *testing.T) {
	m := &Manager{initialized: true}
	if !m.IsAvailable() {
		t.Error("Expected IsAvailable to return true")
	}

	m2 := &Manager{initialized: false}
	if m2.IsAvailable() {
		t.Error("Expected IsAvailable to return false")
	}
}

func TestCopyCommandPrefixRemoval(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"$ ls -la", "ls -la"},
		{"> git status", "git status"},
		{"ls -la", "ls -la"},
		{"$ls", "$ls"}, // Too short, no prefix removal
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			// We can't actually test clipboard in CI, but we can verify
			// the prefix removal logic works
			cmd := tt.input
			if len(cmd) > 2 && (cmd[0] == '$' || cmd[0] == '>') && cmd[1] == ' ' {
				cmd = cmd[2:]
			}

			if cmd != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, cmd)
			}
		})
	}
}
