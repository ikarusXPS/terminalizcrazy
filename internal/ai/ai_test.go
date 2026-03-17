package ai

import (
	"testing"
)

func TestDetectRequestType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected RequestType
	}{
		// Command requests
		{
			name:     "how to find files",
			input:    "how to find large files",
			expected: RequestTypeCommand,
		},
		{
			name:     "show me command",
			input:    "show me how to list all docker containers",
			expected: RequestTypeCommand,
		},
		{
			name:     "German command request",
			input:    "zeige mir wie ich dateien finde",
			expected: RequestTypeCommand,
		},
		{
			name:     "find command",
			input:    "find all .go files in this directory",
			expected: RequestTypeCommand,
		},

		// Explain requests
		{
			name:     "explain error",
			input:    "explain this error: permission denied",
			expected: RequestTypeExplain,
		},
		{
			name:     "what does mean",
			input:    "what does this error mean",
			expected: RequestTypeExplain,
		},
		{
			name:     "German explain",
			input:    "erkläre mir diesen fehler",
			expected: RequestTypeExplain,
		},

		// Chat requests (default)
		{
			name:     "general question",
			input:    "hello",
			expected: RequestTypeChat,
		},
		{
			name:     "opinion question",
			input:    "is rust better than go",
			expected: RequestTypeChat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectRequestType(tt.input)
			if result != tt.expected {
				t.Errorf("detectRequestType(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractCommand(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name: "shell code block",
			content: "Here's the command:\n```shell\nfind . -size +100M\n```\nThis finds large files.",
			expected: "find . -size +100M",
		},
		{
			name: "bash code block",
			content: "```bash\ngit status\n```",
			expected: "git status",
		},
		{
			name: "plain code block",
			content: "```\nls -la\n```",
			expected: "ls -la",
		},
		{
			name:     "no code block",
			content:  "Just some text without a command",
			expected: "",
		},
		{
			name: "multiline command",
			content: "```shell\ngit add .\ngit commit -m \"test\"\n```",
			expected: "git add .\ngit commit -m \"test\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractCommand(tt.content)
			if result != tt.expected {
				t.Errorf("extractCommand() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestExtractExplanation(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "with explanation after code block",
			content:  "```shell\nls -la\n```\nThis lists all files including hidden ones.",
			expected: "This lists all files including hidden ones.",
		},
		{
			name:     "no explanation",
			content:  "```shell\nls -la\n```",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractExplanation(tt.content)
			if result != tt.expected {
				t.Errorf("extractExplanation() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGetRequestContext(t *testing.T) {
	ctx := getRequestContext()

	if ctx == nil {
		t.Fatal("getRequestContext() returned nil")
	}

	if ctx.OS == "" {
		t.Error("OS should not be empty")
	}

	if ctx.Shell == "" {
		t.Error("Shell should not be empty")
	}
}

func TestNewService(t *testing.T) {
	// Test with empty API key
	_, err := NewService(ProviderAnthropic, "")
	if err == nil {
		t.Error("NewService should fail with empty API key")
	}

	// Test with invalid provider
	_, err = NewService("invalid", "some-key")
	if err == nil {
		t.Error("NewService should fail with invalid provider")
	}
}
