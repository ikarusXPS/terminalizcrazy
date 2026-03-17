package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectGo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "project-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create go.mod
	goMod := `module github.com/example/myproject

go 1.21

require (
	github.com/charmbracelet/bubbletea v0.24.0
)
`
	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644)

	detector := NewDetector(tmpDir)
	project := detector.Detect()

	if project.Type != TypeGo {
		t.Errorf("Expected TypeGo, got %v", project.Type)
	}

	if project.Name != "myproject" {
		t.Errorf("Expected 'myproject', got '%s'", project.Name)
	}

	if project.Framework != "Bubble Tea" {
		t.Errorf("Expected 'Bubble Tea', got '%s'", project.Framework)
	}
}

func TestDetectNode(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "project-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create package.json
	packageJSON := `{
  "name": "my-react-app",
  "version": "1.0.0",
  "description": "A React application",
  "dependencies": {
    "react": "^18.0.0",
    "next": "^13.0.0"
  }
}
`
	os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(packageJSON), 0644)

	detector := NewDetector(tmpDir)
	project := detector.Detect()

	if project.Type != TypeNode {
		t.Errorf("Expected TypeNode, got %v", project.Type)
	}

	if project.Name != "my-react-app" {
		t.Errorf("Expected 'my-react-app', got '%s'", project.Name)
	}

	if project.Version != "1.0.0" {
		t.Errorf("Expected '1.0.0', got '%s'", project.Version)
	}

	if project.Framework != "Next.js" {
		t.Errorf("Expected 'Next.js', got '%s'", project.Framework)
	}
}

func TestDetectPython(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "project-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create requirements.txt
	requirements := `django>=4.0
djangorestframework
celery
`
	os.WriteFile(filepath.Join(tmpDir, "requirements.txt"), []byte(requirements), 0644)

	detector := NewDetector(tmpDir)
	project := detector.Detect()

	if project.Type != TypePython {
		t.Errorf("Expected TypePython, got %v", project.Type)
	}

	if project.Framework != "Django" {
		t.Errorf("Expected 'Django', got '%s'", project.Framework)
	}
}

func TestDetectRust(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "project-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create Cargo.toml
	cargoToml := `[package]
name = "my-rust-app"
version = "0.1.0"
description = "A Rust application"

[dependencies]
axum = "0.6"
tokio = { version = "1", features = ["full"] }
`
	os.WriteFile(filepath.Join(tmpDir, "Cargo.toml"), []byte(cargoToml), 0644)

	detector := NewDetector(tmpDir)
	project := detector.Detect()

	if project.Type != TypeRust {
		t.Errorf("Expected TypeRust, got %v", project.Type)
	}

	if project.Name != "my-rust-app" {
		t.Errorf("Expected 'my-rust-app', got '%s'", project.Name)
	}

	if project.Framework != "Axum" {
		t.Errorf("Expected 'Axum', got '%s'", project.Framework)
	}
}

func TestDetectDocker(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "project-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create Dockerfile
	os.WriteFile(filepath.Join(tmpDir, "Dockerfile"), []byte("FROM alpine"), 0644)

	detector := NewDetector(tmpDir)
	project := detector.Detect()

	if project.Type != TypeDocker {
		t.Errorf("Expected TypeDocker, got %v", project.Type)
	}
}

func TestDetectUnknown(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "project-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a random file
	os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("# Hello"), 0644)

	detector := NewDetector(tmpDir)
	project := detector.Detect()

	if project.Type != TypeUnknown {
		t.Errorf("Expected TypeUnknown, got %v", project.Type)
	}

	// Name should be directory name
	expected := filepath.Base(tmpDir)
	if project.Name != expected {
		t.Errorf("Expected '%s', got '%s'", expected, project.Name)
	}
}

func TestGetTypeIcon(t *testing.T) {
	tests := []struct {
		projectType ProjectType
		expected    string
	}{
		{TypeGo, "🐹"},
		{TypeNode, "📦"},
		{TypePython, "🐍"},
		{TypeRust, "🦀"},
		{TypeUnknown, "📁"},
	}

	for _, tt := range tests {
		t.Run(string(tt.projectType), func(t *testing.T) {
			icon := GetTypeIcon(tt.projectType)
			if icon != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, icon)
			}
		})
	}
}

func TestGenerateSessionName(t *testing.T) {
	tests := []struct {
		name     string
		project  Project
		expected string
	}{
		{
			name: "with framework",
			project: Project{
				Name:      "myapp",
				Type:      TypeNode,
				Framework: "Next.js",
			},
			expected: "📦 myapp (Next.js)",
		},
		{
			name: "without framework",
			project: Project{
				Name: "myapp",
				Type: TypeGo,
			},
			expected: "🐹 myapp",
		},
		{
			name: "unknown type",
			project: Project{
				Name: "myapp",
				Type: TypeUnknown,
			},
			expected: "myapp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.project.GenerateSessionName()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestGetAIContext(t *testing.T) {
	project := Project{
		Name:        "myapp",
		Type:        TypeGo,
		Framework:   "Gin",
		Description: "A web server",
	}

	context := project.GetAIContext()

	if context == "" {
		t.Error("Expected non-empty context")
	}

	// Check that it contains expected parts
	if !contains(context, "myapp") {
		t.Error("Context should contain project name")
	}

	if !contains(context, "Go") {
		t.Error("Context should contain project type")
	}

	if !contains(context, "Gin") {
		t.Error("Context should contain framework")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
