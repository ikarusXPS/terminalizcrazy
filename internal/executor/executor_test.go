package executor

import (
	"context"
	"runtime"
	"testing"
)

func TestExecuteSimpleCommand(t *testing.T) {
	e := New()
	ctx := context.Background()

	var cmd string
	if runtime.GOOS == "windows" {
		cmd = "echo hello"
	} else {
		cmd = "echo hello"
	}

	result := e.Execute(ctx, cmd)

	if !result.Success {
		t.Errorf("Expected success, got error: %s", result.Error)
	}

	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}
}

func TestExecuteFailingCommand(t *testing.T) {
	e := New()
	ctx := context.Background()

	// Command that doesn't exist
	result := e.Execute(ctx, "this_command_does_not_exist_12345")

	if result.Success {
		t.Error("Expected failure for non-existent command")
	}

	if result.ExitCode == 0 {
		t.Error("Expected non-zero exit code")
	}
}

func TestAssessRisk(t *testing.T) {
	e := New()

	tests := []struct {
		command  string
		expected RiskLevel
	}{
		// Low risk
		{"ls -la", RiskLow},
		{"cat file.txt", RiskLow},
		{"echo hello", RiskLow},
		{"pwd", RiskLow},
		{"git status", RiskLow},

		// Medium risk
		{"mv file1 file2", RiskMedium},
		{"cp -r dir1 dir2", RiskMedium},
		{"git push origin main", RiskMedium},
		{"npm install express", RiskMedium},
		{"mkdir new_folder", RiskMedium},

		// High risk
		{"rm -rf node_modules", RiskHigh},
		{"git reset --hard HEAD~1", RiskHigh},
		{"DROP TABLE users", RiskHigh},
		{"rmdir /s /q folder", RiskHigh},

		// Critical
		{"sudo rm -rf /", RiskCritical},
		{"chmod 777 /etc", RiskCritical},
		{"dd if=/dev/zero of=/dev/sda", RiskCritical},
		{"shutdown -h now", RiskCritical},
	}

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			result := e.AssessRisk(tt.command)
			if result != tt.expected {
				t.Errorf("AssessRisk(%q) = %v, want %v", tt.command, result, tt.expected)
			}
		})
	}
}

func TestShouldConfirm(t *testing.T) {
	e := New()

	tests := []struct {
		command  string
		expected bool
	}{
		{"ls -la", false},
		{"echo hello", false},
		{"rm -rf node_modules", true},
		{"git push", true},
		{"sudo apt update", true},
	}

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			result := e.ShouldConfirm(tt.command)
			if result != tt.expected {
				t.Errorf("ShouldConfirm(%q) = %v, want %v", tt.command, result, tt.expected)
			}
		})
	}
}

func TestGetRiskDescription(t *testing.T) {
	tests := []struct {
		level    RiskLevel
		contains string
	}{
		{RiskLow, "safe"},
		{RiskMedium, "modify"},
		{RiskHigh, "delete"},
		{RiskCritical, "damage"},
	}

	for _, tt := range tests {
		t.Run(tt.contains, func(t *testing.T) {
			desc := GetRiskDescription(tt.level)
			if desc == "" {
				t.Error("Description should not be empty")
			}
		})
	}
}

func TestResultFormatOutput(t *testing.T) {
	result := &Result{
		Command:  "echo test",
		Output:   "test\n",
		ExitCode: 0,
		Success:  true,
	}

	output := result.FormatOutput()

	if output == "" {
		t.Error("Output should not be empty")
	}

	if !contains(output, "echo test") {
		t.Error("Output should contain the command")
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
