package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestStorage(t *testing.T) (*Storage, func()) {
	tmpDir, err := os.MkdirTemp("", "terminalizcrazy-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	storage, err := New(tmpDir)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to create storage: %v", err)
	}

	cleanup := func() {
		storage.Close()
		os.RemoveAll(tmpDir)
	}

	return storage, cleanup
}

func TestNewStorage(t *testing.T) {
	storage, cleanup := setupTestStorage(t)
	defer cleanup()

	if storage == nil {
		t.Fatal("Storage should not be nil")
	}

	// Check that database file exists
	if _, err := os.Stat(storage.dbPath); os.IsNotExist(err) {
		t.Error("Database file should exist")
	}
}

func TestCreateAndGetSession(t *testing.T) {
	storage, cleanup := setupTestStorage(t)
	defer cleanup()

	// Create session
	session, err := storage.CreateSession("test-123", "Test Session", "/home/user")
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	if session.ID != "test-123" {
		t.Errorf("Expected ID 'test-123', got '%s'", session.ID)
	}

	// Get session
	retrieved, err := storage.GetSession("test-123")
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Retrieved session should not be nil")
	}

	if retrieved.Name != "Test Session" {
		t.Errorf("Expected name 'Test Session', got '%s'", retrieved.Name)
	}
}

func TestGetNonExistentSession(t *testing.T) {
	storage, cleanup := setupTestStorage(t)
	defer cleanup()

	session, err := storage.GetSession("non-existent")
	if err != nil {
		t.Fatalf("Should not error for non-existent session: %v", err)
	}

	if session != nil {
		t.Error("Session should be nil for non-existent ID")
	}
}

func TestListSessions(t *testing.T) {
	storage, cleanup := setupTestStorage(t)
	defer cleanup()

	// Create multiple sessions
	storage.CreateSession("session-1", "Session 1", "/dir1")
	storage.CreateSession("session-2", "Session 2", "/dir2")
	storage.CreateSession("session-3", "Session 3", "/dir3")

	sessions, err := storage.ListSessions(10)
	if err != nil {
		t.Fatalf("Failed to list sessions: %v", err)
	}

	if len(sessions) != 3 {
		t.Errorf("Expected 3 sessions, got %d", len(sessions))
	}
}

func TestSaveAndGetMessages(t *testing.T) {
	storage, cleanup := setupTestStorage(t)
	defer cleanup()

	// Create session first
	storage.CreateSession("test-session", "Test", "/home")

	// Save messages
	_, err := storage.SaveMessage("test-session", "user", "Hello AI", "", true)
	if err != nil {
		t.Fatalf("Failed to save message: %v", err)
	}

	_, err = storage.SaveMessage("test-session", "ai", "Hello! How can I help?", "ls -la", true)
	if err != nil {
		t.Fatalf("Failed to save AI message: %v", err)
	}

	// Get messages
	messages, err := storage.GetSessionMessages("test-session", 100)
	if err != nil {
		t.Fatalf("Failed to get messages: %v", err)
	}

	if len(messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(messages))
	}

	if messages[0].Role != "user" {
		t.Errorf("First message should be from user")
	}

	if messages[1].Command != "ls -la" {
		t.Errorf("Second message should have command 'ls -la'")
	}
}

func TestSaveAndGetCommandHistory(t *testing.T) {
	storage, cleanup := setupTestStorage(t)
	defer cleanup()

	// Save commands
	_, err := storage.SaveCommand("ls -la", "file1\nfile2", true, 50)
	if err != nil {
		t.Fatalf("Failed to save command: %v", err)
	}

	_, err = storage.SaveCommand("git status", "On branch main", true, 100)
	if err != nil {
		t.Fatalf("Failed to save command: %v", err)
	}

	// Get history
	history, err := storage.GetCommandHistory(10)
	if err != nil {
		t.Fatalf("Failed to get history: %v", err)
	}

	if len(history) != 2 {
		t.Errorf("Expected 2 commands, got %d", len(history))
	}

	// Most recent first
	if history[0].Command != "git status" {
		t.Errorf("Expected 'git status' first, got '%s'", history[0].Command)
	}
}

func TestSearchCommands(t *testing.T) {
	storage, cleanup := setupTestStorage(t)
	defer cleanup()

	// Save various commands
	storage.SaveCommand("ls -la", "", true, 10)
	storage.SaveCommand("git status", "", true, 20)
	storage.SaveCommand("git push", "", true, 30)
	storage.SaveCommand("npm install", "", true, 40)

	// Search for git commands
	results, err := storage.SearchCommands("git", 10)
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 git commands, got %d", len(results))
	}
}

func TestGetUniqueCommands(t *testing.T) {
	storage, cleanup := setupTestStorage(t)
	defer cleanup()

	// Save same command multiple times
	storage.SaveCommand("ls -la", "", true, 10)
	storage.SaveCommand("git status", "", true, 20)
	storage.SaveCommand("ls -la", "", true, 30) // duplicate

	unique, err := storage.GetUniqueCommands(10)
	if err != nil {
		t.Fatalf("Failed to get unique commands: %v", err)
	}

	if len(unique) != 2 {
		t.Errorf("Expected 2 unique commands, got %d", len(unique))
	}
}

func TestClearCommandHistory(t *testing.T) {
	storage, cleanup := setupTestStorage(t)
	defer cleanup()

	// Save commands
	storage.SaveCommand("cmd1", "", true, 10)
	storage.SaveCommand("cmd2", "", true, 20)

	// Clear
	err := storage.ClearCommandHistory()
	if err != nil {
		t.Fatalf("Failed to clear history: %v", err)
	}

	// Verify empty
	history, _ := storage.GetCommandHistory(10)
	if len(history) != 0 {
		t.Errorf("History should be empty, got %d items", len(history))
	}
}

func TestGetStats(t *testing.T) {
	storage, cleanup := setupTestStorage(t)
	defer cleanup()

	// Create some data
	storage.CreateSession("s1", "Session 1", "/")
	storage.SaveMessage("s1", "user", "Hello", "", true)
	storage.SaveMessage("s1", "ai", "Hi", "", true)
	storage.SaveCommand("ls", "", true, 10)

	stats, err := storage.GetStats()
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	if stats["sessions"] != 1 {
		t.Errorf("Expected 1 session, got %d", stats["sessions"])
	}

	if stats["messages"] != 2 {
		t.Errorf("Expected 2 messages, got %d", stats["messages"])
	}

	if stats["commands"] != 1 {
		t.Errorf("Expected 1 command, got %d", stats["commands"])
	}
}

func TestDeleteSession(t *testing.T) {
	storage, cleanup := setupTestStorage(t)
	defer cleanup()

	// Create and delete
	storage.CreateSession("to-delete", "Delete Me", "/")
	storage.SaveMessage("to-delete", "user", "Test", "", true)

	err := storage.DeleteSession("to-delete")
	if err != nil {
		t.Fatalf("Failed to delete session: %v", err)
	}

	// Verify deleted
	session, _ := storage.GetSession("to-delete")
	if session != nil {
		t.Error("Session should be deleted")
	}
}

func TestStorageInCustomDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "custom-storage-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	customDir := filepath.Join(tmpDir, "nested", "data")

	storage, err := New(customDir)
	if err != nil {
		t.Fatalf("Failed to create storage in nested dir: %v", err)
	}
	defer storage.Close()

	// Verify directory was created
	if _, err := os.Stat(customDir); os.IsNotExist(err) {
		t.Error("Custom directory should be created")
	}
}
