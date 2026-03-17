package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

func init() {
	// Set timeNowUnixNano to use actual time
	timeNowUnixNano = func() int64 {
		return time.Now().UnixNano()
	}
}

// Plan represents a multi-step task plan
type Plan struct {
	ID          string       `json:"id"`
	Goal        string       `json:"goal"`
	Tasks       []Task       `json:"tasks"`
	Status      PlanStatus   `json:"status"`
	CurrentTask int          `json:"current_task"`
	Context     *PlanContext `json:"context,omitempty"`
}

// PlanStatus represents the status of a plan
type PlanStatus string

const (
	PlanStatusPending   PlanStatus = "pending"
	PlanStatusApproved  PlanStatus = "approved"
	PlanStatusRunning   PlanStatus = "running"
	PlanStatusCompleted PlanStatus = "completed"
	PlanStatusFailed    PlanStatus = "failed"
	PlanStatusCancelled PlanStatus = "cancelled"
)

// Task represents a single task in a plan
type Task struct {
	ID           string          `json:"id"`
	Sequence     int             `json:"sequence"`
	Description  string          `json:"description"`
	Command      string          `json:"command"`
	Status       TaskStatus      `json:"status"`
	Output       string          `json:"output,omitempty"`
	Error        string          `json:"error,omitempty"`
	Verification *Verification   `json:"verification,omitempty"`
	Dependencies []string        `json:"dependencies,omitempty"`
}

// TaskStatus represents the status of a task
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusSkipped   TaskStatus = "skipped"
)

// Verification defines how to verify a task completed successfully
type Verification struct {
	Type         VerificationType `json:"type"`
	Command      string           `json:"command,omitempty"`
	ExpectedCode int              `json:"expected_code,omitempty"`
	Contains     string           `json:"contains,omitempty"`
}

// VerificationType defines verification methods
type VerificationType string

const (
	VerificationExitCode VerificationType = "exit_code"
	VerificationOutput   VerificationType = "output_contains"
	VerificationCommand  VerificationType = "run_command"
)

// PlanContext provides context for planning
type PlanContext struct {
	CurrentDir     string   `json:"current_dir"`
	OS             string   `json:"os"`
	Shell          string   `json:"shell"`
	ProjectType    string   `json:"project_type,omitempty"`
	ProjectName    string   `json:"project_name,omitempty"`
	AvailableTools []string `json:"available_tools,omitempty"`
}

// Planner creates execution plans from goals
type Planner struct {
	client Client
}

// NewPlanner creates a new Planner
func NewPlanner(client Client) *Planner {
	return &Planner{
		client: client,
	}
}

// CreatePlan generates a plan from a goal
func (p *Planner) CreatePlan(ctx context.Context, goal string, planCtx *PlanContext) (*Plan, error) {
	prompt := p.buildPlanningPrompt(goal, planCtx)

	req := &Request{
		UserMessage: prompt,
		Context: &RequestContext{
			CurrentDir: planCtx.CurrentDir,
			OS:         planCtx.OS,
			Shell:      planCtx.Shell,
		},
		Type: RequestTypeChat,
	}

	resp, err := p.client.Complete(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate plan: %w", err)
	}

	plan, err := p.parsePlanResponse(resp.Content, goal)
	if err != nil {
		return nil, fmt.Errorf("failed to parse plan: %w", err)
	}

	plan.Context = planCtx
	plan.Status = PlanStatusPending

	return plan, nil
}

// buildPlanningPrompt creates the prompt for plan generation
func (p *Planner) buildPlanningPrompt(goal string, planCtx *PlanContext) string {
	var sb strings.Builder

	sb.WriteString("You are a task planner. Create a step-by-step plan to accomplish the following goal.\n\n")
	sb.WriteString(fmt.Sprintf("Goal: %s\n\n", goal))

	if planCtx != nil {
		sb.WriteString("Context:\n")
		sb.WriteString(fmt.Sprintf("- Operating System: %s\n", planCtx.OS))
		sb.WriteString(fmt.Sprintf("- Shell: %s\n", planCtx.Shell))
		sb.WriteString(fmt.Sprintf("- Current Directory: %s\n", planCtx.CurrentDir))
		if planCtx.ProjectType != "" {
			sb.WriteString(fmt.Sprintf("- Project Type: %s\n", planCtx.ProjectType))
		}
		if planCtx.ProjectName != "" {
			sb.WriteString(fmt.Sprintf("- Project Name: %s\n", planCtx.ProjectName))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(`Respond with a JSON plan in this exact format:
{
  "tasks": [
    {
      "description": "Brief description of what this step does",
      "command": "the shell command to execute",
      "verification": {
        "type": "exit_code",
        "expected_code": 0
      }
    }
  ]
}

Rules:
1. Each task should be a single, atomic command
2. Tasks should be in logical order
3. Include verification for important tasks
4. Use safe, non-destructive commands when possible
5. Add appropriate flags for non-interactive execution
6. Maximum 10 tasks per plan

Verification types:
- "exit_code": Check exit code (expected_code: 0 for success)
- "output_contains": Check if output contains text (contains: "expected text")
- "run_command": Run a separate verification command (command: "verify command")

Respond ONLY with the JSON, no other text.`)

	return sb.String()
}

// parsePlanResponse parses the AI response into a Plan
func (p *Planner) parsePlanResponse(content string, goal string) (*Plan, error) {
	// Extract JSON from response
	jsonContent := extractJSON(content)
	if jsonContent == "" {
		return nil, fmt.Errorf("no JSON found in response")
	}

	var parsed struct {
		Tasks []struct {
			Description  string `json:"description"`
			Command      string `json:"command"`
			Verification *struct {
				Type         string `json:"type"`
				ExpectedCode int    `json:"expected_code"`
				Contains     string `json:"contains"`
				Command      string `json:"command"`
			} `json:"verification,omitempty"`
		} `json:"tasks"`
	}

	if err := json.Unmarshal([]byte(jsonContent), &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if len(parsed.Tasks) == 0 {
		return nil, fmt.Errorf("plan contains no tasks")
	}

	plan := &Plan{
		ID:          generatePlanID(),
		Goal:        goal,
		Tasks:       make([]Task, len(parsed.Tasks)),
		CurrentTask: 0,
	}

	for i, t := range parsed.Tasks {
		task := Task{
			ID:          fmt.Sprintf("task-%d", i+1),
			Sequence:    i + 1,
			Description: t.Description,
			Command:     t.Command,
			Status:      TaskStatusPending,
		}

		if t.Verification != nil {
			task.Verification = &Verification{
				Type:         VerificationType(t.Verification.Type),
				ExpectedCode: t.Verification.ExpectedCode,
				Contains:     t.Verification.Contains,
				Command:      t.Verification.Command,
			}
		}

		plan.Tasks[i] = task
	}

	return plan, nil
}

// extractJSON extracts JSON object from text
func extractJSON(text string) string {
	// Try to find JSON object
	start := strings.Index(text, "{")
	if start == -1 {
		return ""
	}

	// Find matching closing brace
	depth := 0
	for i := start; i < len(text); i++ {
		switch text[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return text[start : i+1]
			}
		}
	}

	return ""
}

// generatePlanID generates a unique plan ID
func generatePlanID() string {
	// Use UUID for proper unique IDs
	return fmt.Sprintf("plan-%s", generateShortID())
}

// generateShortID generates a short random ID
func generateShortID() string {
	// Simple but effective ID generation
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 8)
	for i := range result {
		// Use time-based seed mixed with position
		seed := int(timeNowUnixNano()) + i
		result[i] = chars[seed%len(chars)]
	}
	return string(result)
}

// timeNowUnixNano returns current time in nanoseconds
var timeNowUnixNano = func() int64 {
	return int64(1234567890) // Will be set properly by init
}

// ValidatePlan checks if a plan is valid
func (p *Planner) ValidatePlan(plan *Plan) error {
	if plan == nil {
		return fmt.Errorf("plan is nil")
	}

	if plan.Goal == "" {
		return fmt.Errorf("plan has no goal")
	}

	if len(plan.Tasks) == 0 {
		return fmt.Errorf("plan has no tasks")
	}

	// Check for dangerous commands
	for _, task := range plan.Tasks {
		if isDangerousCommand(task.Command) {
			return fmt.Errorf("task %d contains potentially dangerous command: %s",
				task.Sequence, task.Description)
		}
	}

	return nil
}

// isDangerousCommand checks if a command is dangerous
func isDangerousCommand(cmd string) bool {
	lower := strings.ToLower(cmd)

	dangerous := []string{
		"rm -rf /",
		"rm -rf /*",
		":(){ :|:& };:",
		"> /dev/sda",
		"dd if=/dev/zero",
		"mkfs.",
		"format c:",
		"del /f /s /q c:\\",
	}

	for _, d := range dangerous {
		if strings.Contains(lower, d) {
			return true
		}
	}

	// Check for command injection patterns
	injectionPatterns := []string{
		"$(", "`", "&&", "||", ";",
	}

	// Only flag injection if combined with dangerous keywords
	for _, pattern := range injectionPatterns {
		if strings.Contains(cmd, pattern) {
			for _, keyword := range []string{"rm ", "del ", "format", "mkfs"} {
				if strings.Contains(lower, keyword) {
					return true
				}
			}
		}
	}

	return false
}

// SummarizePlan creates a human-readable summary
func (p *Planner) SummarizePlan(plan *Plan) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Plan: %s\n", plan.Goal))
	sb.WriteString(fmt.Sprintf("Status: %s\n", plan.Status))
	sb.WriteString(fmt.Sprintf("Tasks: %d\n\n", len(plan.Tasks)))

	for _, task := range plan.Tasks {
		statusIcon := getStatusIcon(task.Status)
		sb.WriteString(fmt.Sprintf("%s %d. %s\n", statusIcon, task.Sequence, task.Description))
		sb.WriteString(fmt.Sprintf("   $ %s\n", task.Command))
		if task.Output != "" {
			lines := strings.Split(task.Output, "\n")
			if len(lines) > 3 {
				sb.WriteString(fmt.Sprintf("   Output: %s... (%d more lines)\n",
					strings.Join(lines[:3], "\n   "), len(lines)-3))
			} else {
				sb.WriteString(fmt.Sprintf("   Output: %s\n", task.Output))
			}
		}
		if task.Error != "" {
			sb.WriteString(fmt.Sprintf("   Error: %s\n", task.Error))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// getStatusIcon returns an icon for the status
func getStatusIcon(status TaskStatus) string {
	switch status {
	case TaskStatusPending:
		return "○"
	case TaskStatusRunning:
		return "◐"
	case TaskStatusCompleted:
		return "✓"
	case TaskStatusFailed:
		return "✗"
	case TaskStatusSkipped:
		return "⊘"
	default:
		return "?"
	}
}

// GetNextTask returns the next pending task
func (plan *Plan) GetNextTask() *Task {
	for i := range plan.Tasks {
		if plan.Tasks[i].Status == TaskStatusPending {
			return &plan.Tasks[i]
		}
	}
	return nil
}

// IsComplete returns true if all tasks are done
func (plan *Plan) IsComplete() bool {
	for _, task := range plan.Tasks {
		if task.Status == TaskStatusPending || task.Status == TaskStatusRunning {
			return false
		}
	}
	return true
}

// GetProgress returns completion percentage
func (plan *Plan) GetProgress() float64 {
	if len(plan.Tasks) == 0 {
		return 0
	}

	completed := 0
	for _, task := range plan.Tasks {
		if task.Status == TaskStatusCompleted || task.Status == TaskStatusSkipped {
			completed++
		}
	}

	return float64(completed) / float64(len(plan.Tasks)) * 100
}

// UpdateTaskStatus updates a task's status
func (plan *Plan) UpdateTaskStatus(taskID string, status TaskStatus, output, errMsg string) {
	for i := range plan.Tasks {
		if plan.Tasks[i].ID == taskID {
			plan.Tasks[i].Status = status
			plan.Tasks[i].Output = output
			plan.Tasks[i].Error = errMsg
			break
		}
	}
}

// ReplanOnFailure creates a recovery prompt for failed tasks
func (p *Planner) ReplanOnFailure(ctx context.Context, plan *Plan, failedTask *Task) (*Plan, error) {
	prompt := fmt.Sprintf(`The following task failed:
Task: %s
Command: %s
Error: %s

Original Goal: %s

Please create a revised plan to:
1. Fix the issue that caused the failure
2. Complete the remaining tasks

Consider alternative approaches if the original command is not working.`,
		failedTask.Description,
		failedTask.Command,
		failedTask.Error,
		plan.Goal,
	)

	return p.CreatePlan(ctx, prompt, plan.Context)
}

// CommandExtractor extracts commands from natural language
type CommandExtractor struct {
	patterns []*regexp.Regexp
}

// NewCommandExtractor creates a new extractor
func NewCommandExtractor() *CommandExtractor {
	patterns := []*regexp.Regexp{
		regexp.MustCompile("```(?:bash|sh|shell)?\n([^`]+)```"),
		regexp.MustCompile("`([^`]+)`"),
		regexp.MustCompile("\\$\\s+(.+)"),
	}

	return &CommandExtractor{patterns: patterns}
}

// Extract extracts commands from text
func (e *CommandExtractor) Extract(text string) []string {
	var commands []string
	seen := make(map[string]bool)

	for _, pattern := range e.patterns {
		matches := pattern.FindAllStringSubmatch(text, -1)
		for _, match := range matches {
			if len(match) > 1 {
				cmd := strings.TrimSpace(match[1])
				if cmd != "" && !seen[cmd] {
					seen[cmd] = true
					commands = append(commands, cmd)
				}
			}
		}
	}

	return commands
}
