package workflows

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/terminalizcrazy/terminalizcrazy/internal/executor"
)

// Workflow represents a reusable workflow template
type Workflow struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Steps       []WorkflowStep  `json:"steps"`
	Variables   []Variable      `json:"variables,omitempty"`
	Tags        []string        `json:"tags,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// WorkflowStep represents a single step in a workflow
type WorkflowStep struct {
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Command     string           `json:"command"`
	OnFail      OnFailAction     `json:"on_fail,omitempty"`
	Condition   string           `json:"condition,omitempty"`
	Timeout     time.Duration    `json:"timeout,omitempty"`
	Retries     int              `json:"retries,omitempty"`
	CaptureAs   string           `json:"capture_as,omitempty"` // Variable name to capture output
}

// OnFailAction defines what to do when a step fails
type OnFailAction string

const (
	OnFailStop     OnFailAction = "stop"
	OnFailSkip     OnFailAction = "skip"
	OnFailContinue OnFailAction = "continue"
	OnFailRetry    OnFailAction = "retry"
)

// Variable represents a workflow variable
type Variable struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Default     string `json:"default,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Type        string `json:"type,omitempty"` // string, number, boolean, path
}

// WorkflowExecution represents a running workflow
type WorkflowExecution struct {
	ID           string                     `json:"id"`
	WorkflowID   string                     `json:"workflow_id"`
	WorkflowName string                     `json:"workflow_name"`
	Status       ExecutionStatus            `json:"status"`
	Variables    map[string]string          `json:"variables"`
	StepResults  []StepResult               `json:"step_results"`
	CurrentStep  int                        `json:"current_step"`
	StartedAt    time.Time                  `json:"started_at"`
	CompletedAt  *time.Time                 `json:"completed_at,omitempty"`
	Error        string                     `json:"error,omitempty"`
}

// ExecutionStatus represents the status of a workflow execution
type ExecutionStatus string

const (
	ExecutionPending   ExecutionStatus = "pending"
	ExecutionRunning   ExecutionStatus = "running"
	ExecutionCompleted ExecutionStatus = "completed"
	ExecutionFailed    ExecutionStatus = "failed"
	ExecutionCancelled ExecutionStatus = "cancelled"
)

// StepResult represents the result of a workflow step
type StepResult struct {
	StepName   string        `json:"step_name"`
	Command    string        `json:"command"`
	Output     string        `json:"output"`
	Error      string        `json:"error,omitempty"`
	ExitCode   int           `json:"exit_code"`
	Success    bool          `json:"success"`
	Duration   time.Duration `json:"duration"`
	Skipped    bool          `json:"skipped"`
	RetryCount int           `json:"retry_count,omitempty"`
}

// WorkflowEngine executes workflows
type WorkflowEngine struct {
	executor   *executor.Executor
	workflows  map[string]*Workflow
	executions map[string]*WorkflowExecution

	// Callbacks
	onStepStart    func(*WorkflowExecution, *WorkflowStep)
	onStepComplete func(*WorkflowExecution, *StepResult)
	onComplete     func(*WorkflowExecution)
}

// NewWorkflowEngine creates a new workflow engine
func NewWorkflowEngine(exec *executor.Executor) *WorkflowEngine {
	return &WorkflowEngine{
		executor:   exec,
		workflows:  make(map[string]*Workflow),
		executions: make(map[string]*WorkflowExecution),
	}
}

// SetCallbacks sets the callback functions
func (e *WorkflowEngine) SetCallbacks(
	onStepStart func(*WorkflowExecution, *WorkflowStep),
	onStepComplete func(*WorkflowExecution, *StepResult),
	onComplete func(*WorkflowExecution),
) {
	e.onStepStart = onStepStart
	e.onStepComplete = onStepComplete
	e.onComplete = onComplete
}

// RegisterWorkflow registers a workflow
func (e *WorkflowEngine) RegisterWorkflow(workflow *Workflow) error {
	if workflow.Name == "" {
		return fmt.Errorf("workflow name is required")
	}
	if len(workflow.Steps) == 0 {
		return fmt.Errorf("workflow must have at least one step")
	}

	workflow.ID = generateWorkflowID(workflow.Name)
	e.workflows[workflow.Name] = workflow
	return nil
}

// GetWorkflow retrieves a workflow by name
func (e *WorkflowEngine) GetWorkflow(name string) *Workflow {
	return e.workflows[name]
}

// ListWorkflows returns all registered workflows
func (e *WorkflowEngine) ListWorkflows() []*Workflow {
	workflows := make([]*Workflow, 0, len(e.workflows))
	for _, wf := range e.workflows {
		workflows = append(workflows, wf)
	}
	return workflows
}

// Execute runs a workflow
func (e *WorkflowEngine) Execute(ctx context.Context, workflowName string, variables map[string]string) (*WorkflowExecution, error) {
	workflow, ok := e.workflows[workflowName]
	if !ok {
		return nil, fmt.Errorf("workflow not found: %s", workflowName)
	}

	// Validate and fill in default variables
	vars, err := e.prepareVariables(workflow, variables)
	if err != nil {
		return nil, err
	}

	// Create execution
	execution := &WorkflowExecution{
		ID:           generateExecutionID(),
		WorkflowID:   workflow.ID,
		WorkflowName: workflow.Name,
		Status:       ExecutionRunning,
		Variables:    vars,
		StepResults:  make([]StepResult, 0),
		CurrentStep:  0,
		StartedAt:    time.Now(),
	}

	e.executions[execution.ID] = execution

	// Execute steps
	for i, step := range workflow.Steps {
		execution.CurrentStep = i

		// Check condition
		if step.Condition != "" && !e.evaluateCondition(step.Condition, vars) {
			execution.StepResults = append(execution.StepResults, StepResult{
				StepName: step.Name,
				Skipped:  true,
			})
			continue
		}

		if e.onStepStart != nil {
			e.onStepStart(execution, &step)
		}

		result := e.executeStep(ctx, &step, vars)
		execution.StepResults = append(execution.StepResults, result)

		// Capture output as variable if specified
		if step.CaptureAs != "" && result.Success {
			vars[step.CaptureAs] = strings.TrimSpace(result.Output)
		}

		if e.onStepComplete != nil {
			e.onStepComplete(execution, &result)
		}

		// Handle failure
		if !result.Success && !result.Skipped {
			switch step.OnFail {
			case OnFailStop, "":
				execution.Status = ExecutionFailed
				execution.Error = fmt.Sprintf("Step '%s' failed: %s", step.Name, result.Error)
				now := time.Now()
				execution.CompletedAt = &now
				if e.onComplete != nil {
					e.onComplete(execution)
				}
				return execution, fmt.Errorf("workflow failed at step '%s': %s", step.Name, result.Error)

			case OnFailSkip, OnFailContinue:
				continue

			case OnFailRetry:
				// Already handled in executeStep
			}
		}
	}

	execution.Status = ExecutionCompleted
	now := time.Now()
	execution.CompletedAt = &now

	if e.onComplete != nil {
		e.onComplete(execution)
	}

	return execution, nil
}

// executeStep executes a single workflow step
func (e *WorkflowEngine) executeStep(ctx context.Context, step *WorkflowStep, vars map[string]string) StepResult {
	// Substitute variables in command
	command := substituteVariables(step.Command, vars)

	result := StepResult{
		StepName: step.Name,
		Command:  command,
	}

	// Set timeout
	timeout := step.Timeout
	if timeout == 0 {
		timeout = 60 * time.Second
	}

	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Execute with retries
	maxRetries := step.Retries
	if maxRetries < 0 {
		maxRetries = 0
	}

	start := time.Now()

	for attempt := 0; attempt <= maxRetries; attempt++ {
		result.RetryCount = attempt

		execResult := e.executor.Execute(execCtx, command)

		result.Output = execResult.Output
		result.Error = execResult.Error
		result.ExitCode = execResult.ExitCode
		result.Success = execResult.Success
		result.Duration = time.Since(start)

		if result.Success {
			break
		}

		if attempt < maxRetries {
			// Wait before retry
			time.Sleep(time.Second * time.Duration(attempt+1))
		}
	}

	return result
}

// prepareVariables validates and fills in default variables
func (e *WorkflowEngine) prepareVariables(workflow *Workflow, provided map[string]string) (map[string]string, error) {
	vars := make(map[string]string)

	// Copy provided variables
	for k, v := range provided {
		vars[k] = v
	}

	// Check required and fill defaults
	for _, v := range workflow.Variables {
		if _, ok := vars[v.Name]; !ok {
			if v.Default != "" {
				vars[v.Name] = v.Default
			} else if v.Required {
				return nil, fmt.Errorf("required variable '%s' not provided", v.Name)
			}
		}
	}

	return vars, nil
}

// evaluateCondition evaluates a simple condition
func (e *WorkflowEngine) evaluateCondition(condition string, vars map[string]string) bool {
	// Simple variable existence check: ${VAR}
	if strings.HasPrefix(condition, "${") && strings.HasSuffix(condition, "}") {
		varName := condition[2 : len(condition)-1]
		val, ok := vars[varName]
		return ok && val != ""
	}

	// Simple equality check: ${VAR}==value
	if strings.Contains(condition, "==") {
		parts := strings.SplitN(condition, "==", 2)
		if len(parts) == 2 {
			left := substituteVariables(strings.TrimSpace(parts[0]), vars)
			right := strings.TrimSpace(parts[1])
			return left == right
		}
	}

	// Simple inequality check: ${VAR}!=value
	if strings.Contains(condition, "!=") {
		parts := strings.SplitN(condition, "!=", 2)
		if len(parts) == 2 {
			left := substituteVariables(strings.TrimSpace(parts[0]), vars)
			right := strings.TrimSpace(parts[1])
			return left != right
		}
	}

	return true
}

// substituteVariables replaces ${var} patterns with values
func substituteVariables(text string, vars map[string]string) string {
	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	return re.ReplaceAllStringFunc(text, func(match string) string {
		varName := match[2 : len(match)-1]
		if val, ok := vars[varName]; ok {
			return val
		}
		return match
	})
}

// generateWorkflowID generates a workflow ID
func generateWorkflowID(name string) string {
	return fmt.Sprintf("wf-%s-%d", strings.ToLower(strings.ReplaceAll(name, " ", "-")), time.Now().UnixNano()%10000)
}

// generateExecutionID generates an execution ID
func generateExecutionID() string {
	return fmt.Sprintf("exec-%d", time.Now().UnixNano()%1000000)
}

// GetExecution retrieves an execution by ID
func (e *WorkflowEngine) GetExecution(id string) *WorkflowExecution {
	return e.executions[id]
}

// ListExecutions returns all executions
func (e *WorkflowEngine) ListExecutions() []*WorkflowExecution {
	execs := make([]*WorkflowExecution, 0, len(e.executions))
	for _, exec := range e.executions {
		execs = append(execs, exec)
	}
	return execs
}

// ToJSON serializes a workflow to JSON
func (w *Workflow) ToJSON() (string, error) {
	data, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON deserializes a workflow from JSON
func WorkflowFromJSON(data string) (*Workflow, error) {
	var workflow Workflow
	if err := json.Unmarshal([]byte(data), &workflow); err != nil {
		return nil, err
	}
	return &workflow, nil
}

// Summary returns a brief summary of the workflow
func (w *Workflow) Summary() string {
	return fmt.Sprintf("%s (%d steps)", w.Name, len(w.Steps))
}

// ExecutionSummary returns a brief summary of an execution
func (e *WorkflowExecution) ExecutionSummary() string {
	completed := 0
	failed := 0
	for _, r := range e.StepResults {
		if r.Success {
			completed++
		} else if !r.Skipped {
			failed++
		}
	}
	return fmt.Sprintf("%s: %d completed, %d failed", e.Status, completed, failed)
}
