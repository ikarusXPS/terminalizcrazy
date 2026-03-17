package secretguard

import (
	"strings"
	"testing"
)

func TestScanDetectsSecrets(t *testing.T) {
	guard := New(true)

	tests := []struct {
		name     string
		input    string
		expected SecretType
	}{
		{
			name:     "AWS Access Key",
			input:    "AKIAIOSFODNN7EXAMPLE",
			expected: SecretTypeAWS,
		},
		{
			name:     "GitHub Token",
			input:    "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			expected: SecretTypeGitHub,
		},
		{
			name:     "Anthropic Key",
			input:    "sk-ant-api03-xxxxxxxxxxxxxxxxxxxx",
			expected: SecretTypeAnthropic,
		},
		{
			name:     "OpenAI Key",
			input:    "sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			expected: SecretTypeOpenAI,
		},
		{
			name:     "JWT Token",
			input:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U",
			expected: SecretTypeJWT,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detections := guard.Scan(tt.input)
			if len(detections) == 0 {
				t.Errorf("Expected to detect %s, got no detections", tt.expected)
				return
			}
			if detections[0].Type != tt.expected {
				t.Errorf("Expected type %s, got %s", tt.expected, detections[0].Type)
			}
		})
	}
}

func TestMaskReplacesSecrets(t *testing.T) {
	guard := New(true)

	input := "My API key is sk-ant-api03-xxxxxxxxxxxxxxxxxxxx please don't share"
	masked := guard.Mask(input)

	if strings.Contains(masked, "sk-ant-api03-xxxxxxxxxxxxxxxxxxxx") {
		t.Error("Masked output still contains the original secret")
	}

	if !strings.Contains(masked, "****") {
		t.Error("Masked output should contain asterisks")
	}
}

func TestDisabledGuardDoesNotMask(t *testing.T) {
	guard := New(false)

	input := "sk-ant-api03-xxxxxxxxxxxxxxxxxxxx"
	masked := guard.Mask(input)

	if masked != input {
		t.Error("Disabled guard should not modify input")
	}
}

func TestHasSecrets(t *testing.T) {
	guard := New(true)

	if !guard.HasSecrets("sk-ant-api03-xxxxxxxxxxxxxxxxxxxx") {
		t.Error("Should detect secret")
	}

	if guard.HasSecrets("just some normal text") {
		t.Error("Should not detect secret in normal text")
	}
}
