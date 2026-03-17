package executor

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// Result represents the result of a command execution
type Result struct {
	Command    string
	Output     string
	Error      string
	ExitCode   int
	Duration   time.Duration
	Success    bool
}

// RiskLevel indicates how dangerous a command might be
type RiskLevel int

const (
	RiskLow    RiskLevel = iota // Safe commands (ls, cat, echo)
	RiskMedium                  // Commands that modify files (mv, cp)
	RiskHigh                    // Destructive commands (rm, drop, delete)
	RiskCritical                // System-altering commands (sudo, chmod 777)
)

// Executor handles command execution
type Executor struct {
	shell     string
	shellFlag string
	timeout   time.Duration
}

// New creates a new Executor
func New() *Executor {
	shell := "bash"
	shellFlag := "-c"

	if runtime.GOOS == "windows" {
		shell = "cmd"
		shellFlag = "/C"
	}

	return &Executor{
		shell:     shell,
		shellFlag: shellFlag,
		timeout:   60 * time.Second,
	}
}

// Execute runs a command and returns the result
func (e *Executor) Execute(ctx context.Context, command string) *Result {
	start := time.Now()

	// Create context with timeout
	execCtx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	// Create command
	cmd := exec.CommandContext(execCtx, e.shell, e.shellFlag, command)

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run command
	err := cmd.Run()

	result := &Result{
		Command:  command,
		Output:   stdout.String(),
		Error:    stderr.String(),
		Duration: time.Since(start),
		Success:  err == nil,
	}

	// Get exit code
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
			if result.Error == "" {
				result.Error = err.Error()
			}
		}
	}

	return result
}

// AssessRisk evaluates the risk level of a command
func (e *Executor) AssessRisk(command string) RiskLevel {
	lower := strings.ToLower(command)

	// Critical commands - system altering
	criticalPatterns := []string{
		"sudo", "su ", "chmod 777", "chown", "mkfs",
		"dd if=", ":(){ :|:& };:", "rm -rf /",
		"format", "> /dev/", "shutdown", "reboot",
		"reg delete", "del /f /s /q c:\\",
	}
	for _, pattern := range criticalPatterns {
		if strings.Contains(lower, pattern) {
			return RiskCritical
		}
	}

	// High risk - destructive
	highRiskPatterns := []string{
		"rm -rf", "rm -r", "rmdir", "del ", "erase",
		"drop table", "drop database", "truncate",
		"git reset --hard", "git clean -fd",
		"npm uninstall", "pip uninstall",
	}
	for _, pattern := range highRiskPatterns {
		if strings.Contains(lower, pattern) {
			return RiskHigh
		}
	}

	// Medium risk - modifying
	mediumRiskPatterns := []string{
		"mv ", "cp ", "rename", "move",
		"git push", "git commit", "git checkout",
		"npm install", "pip install", "go install",
		"mkdir", "touch", "wget", "curl -o",
	}
	for _, pattern := range mediumRiskPatterns {
		if strings.Contains(lower, pattern) {
			return RiskMedium
		}
	}

	return RiskLow
}

// GetRiskDescription returns a human-readable risk description
func GetRiskDescription(level RiskLevel) string {
	switch level {
	case RiskCritical:
		return "CRITICAL: This command could damage your system"
	case RiskHigh:
		return "HIGH: This command will delete or destroy data"
	case RiskMedium:
		return "MEDIUM: This command will modify files or system state"
	default:
		return "LOW: This command is safe to run"
	}
}

// GetRiskColor returns a color code for the risk level
func GetRiskColor(level RiskLevel) string {
	switch level {
	case RiskCritical:
		return "#FF0000" // Red
	case RiskHigh:
		return "#FF6B6B" // Light red
	case RiskMedium:
		return "#FFAA00" // Orange
	default:
		return "#04B575" // Green
	}
}

// ShouldConfirm returns true if the command requires confirmation
func (e *Executor) ShouldConfirm(command string) bool {
	risk := e.AssessRisk(command)
	return risk >= RiskMedium
}

// FormatOutput formats the command result for display
func (r *Result) FormatOutput() string {
	var sb strings.Builder

	// Command
	sb.WriteString(fmt.Sprintf("$ %s\n", r.Command))
	sb.WriteString(fmt.Sprintf("Duration: %v\n\n", r.Duration.Round(time.Millisecond)))

	// Output
	if r.Output != "" {
		sb.WriteString(r.Output)
		if !strings.HasSuffix(r.Output, "\n") {
			sb.WriteString("\n")
		}
	}

	// Error
	if r.Error != "" {
		sb.WriteString(fmt.Sprintf("\nError: %s", r.Error))
	}

	// Exit code
	if !r.Success {
		sb.WriteString(fmt.Sprintf("\nExit code: %d", r.ExitCode))
	}

	return sb.String()
}
