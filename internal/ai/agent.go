package ai

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/terminalizcrazy/terminalizcrazy/internal/executor"
)

// AgentMode defines how the agent operates
type AgentMode string

const (
	// AgentModeOff - Agent is disabled, manual execution only
	AgentModeOff AgentMode = "off"

	// AgentModeSuggest - Agent creates plans but requires approval
	AgentModeSuggest AgentMode = "suggest"

	// AgentModeAuto - Agent executes plans automatically
	AgentModeAuto AgentMode = "auto"
)

// Agent orchestrates autonomous task completion
type Agent struct {
	client   Client
	executor *executor.Executor
	planner  *Planner
	mode     AgentMode

	currentPlan *Plan
	planHistory []*Plan

	// Callbacks for TUI integration
	onPlanCreated    func(*Plan)
	onTaskStarted    func(*Task)
	onTaskCompleted  func(*Task)
	onTaskFailed     func(*Task)
	onPlanCompleted  func(*Plan)
	onApprovalNeeded func(*Plan)

	mu sync.RWMutex
}

// AgentConfig holds configuration for the agent
type AgentConfig struct {
	Mode             AgentMode
	MaxTasksPerPlan  int
	TaskTimeout      time.Duration
	RequireApproval  bool
	AllowDangerous   bool
	MaxRetries       int
}

// DefaultAgentConfig returns sensible defaults
func DefaultAgentConfig() *AgentConfig {
	return &AgentConfig{
		Mode:             AgentModeSuggest,
		MaxTasksPerPlan:  10,
		TaskTimeout:      60 * time.Second,
		RequireApproval:  true,
		AllowDangerous:   false,
		MaxRetries:       2,
	}
}

// NewAgent creates a new Agent
func NewAgent(client Client, exec *executor.Executor, config *AgentConfig) *Agent {
	if config == nil {
		config = DefaultAgentConfig()
	}

	return &Agent{
		client:      client,
		executor:    exec,
		planner:     NewPlanner(client),
		mode:        config.Mode,
		planHistory: make([]*Plan, 0),
	}
}

// SetMode changes the agent mode
func (a *Agent) SetMode(mode AgentMode) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.mode = mode
}

// GetMode returns the current mode
func (a *Agent) GetMode() AgentMode {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.mode
}

// SetCallbacks sets the callback functions for TUI integration
func (a *Agent) SetCallbacks(
	onPlanCreated func(*Plan),
	onTaskStarted func(*Task),
	onTaskCompleted func(*Task),
	onTaskFailed func(*Task),
	onPlanCompleted func(*Plan),
	onApprovalNeeded func(*Plan),
) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.onPlanCreated = onPlanCreated
	a.onTaskStarted = onTaskStarted
	a.onTaskCompleted = onTaskCompleted
	a.onTaskFailed = onTaskFailed
	a.onPlanCompleted = onPlanCompleted
	a.onApprovalNeeded = onApprovalNeeded
}

// ProcessGoal takes a user goal and creates/executes a plan
func (a *Agent) ProcessGoal(ctx context.Context, goal string, planCtx *PlanContext) (*Plan, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Create the plan
	plan, err := a.planner.CreatePlan(ctx, goal, planCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to create plan: %w", err)
	}

	// Validate the plan
	if err := a.planner.ValidatePlan(plan); err != nil {
		return nil, fmt.Errorf("plan validation failed: %w", err)
	}

	a.currentPlan = plan
	a.planHistory = append(a.planHistory, plan)

	// Notify plan created
	if a.onPlanCreated != nil {
		a.onPlanCreated(plan)
	}

	// Handle based on mode
	switch a.mode {
	case AgentModeOff:
		return plan, nil

	case AgentModeSuggest:
		plan.Status = PlanStatusPending
		if a.onApprovalNeeded != nil {
			a.onApprovalNeeded(plan)
		}
		return plan, nil

	case AgentModeAuto:
		plan.Status = PlanStatusApproved
		go a.executePlanAsync(ctx, plan)
		return plan, nil
	}

	return plan, nil
}

// ApprovePlan approves a pending plan for execution
func (a *Agent) ApprovePlan(ctx context.Context, planID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.currentPlan == nil || a.currentPlan.ID != planID {
		return fmt.Errorf("plan not found: %s", planID)
	}

	if a.currentPlan.Status != PlanStatusPending {
		return fmt.Errorf("plan is not pending approval")
	}

	a.currentPlan.Status = PlanStatusApproved
	go a.executePlanAsync(ctx, a.currentPlan)

	return nil
}

// RejectPlan cancels a pending plan
func (a *Agent) RejectPlan(planID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.currentPlan == nil || a.currentPlan.ID != planID {
		return fmt.Errorf("plan not found: %s", planID)
	}

	a.currentPlan.Status = PlanStatusCancelled
	return nil
}

// executePlanAsync executes the plan asynchronously
func (a *Agent) executePlanAsync(ctx context.Context, plan *Plan) {
	a.executePlan(ctx, plan)
}

// ExecutePlan executes a plan synchronously
func (a *Agent) ExecutePlan(ctx context.Context, plan *Plan) error {
	a.mu.Lock()
	a.currentPlan = plan
	a.mu.Unlock()

	return a.executePlan(ctx, plan)
}

// executePlan runs through all tasks in the plan
func (a *Agent) executePlan(ctx context.Context, plan *Plan) error {
	plan.Status = PlanStatusRunning

	for {
		task := plan.GetNextTask()
		if task == nil {
			break
		}

		// Execute the task
		err := a.executeTask(ctx, plan, task)
		if err != nil {
			plan.Status = PlanStatusFailed
			if a.onTaskFailed != nil {
				a.onTaskFailed(task)
			}
			return err
		}
	}

	plan.Status = PlanStatusCompleted
	if a.onPlanCompleted != nil {
		a.onPlanCompleted(plan)
	}

	return nil
}

// executeTask executes a single task
func (a *Agent) executeTask(ctx context.Context, plan *Plan, task *Task) error {
	task.Status = TaskStatusRunning

	if a.onTaskStarted != nil {
		a.onTaskStarted(task)
	}

	// Check risk level
	risk := a.executor.AssessRisk(task.Command)
	if risk >= executor.RiskHigh && a.mode != AgentModeAuto {
		task.Status = TaskStatusPending
		task.Error = "High-risk command requires manual approval"
		return fmt.Errorf("high-risk command: %s", task.Command)
	}

	// Execute the command
	result := a.executor.Execute(ctx, task.Command)

	task.Output = result.Output
	if result.Error != "" {
		task.Error = result.Error
	}

	// Verify if needed
	if task.Verification != nil {
		verified := a.verifyTask(ctx, task, result)
		if !verified {
			task.Status = TaskStatusFailed
			if task.Error == "" {
				task.Error = "Verification failed"
			}
			return fmt.Errorf("task verification failed: %s", task.Description)
		}
	} else {
		// Default verification: check exit code
		if !result.Success {
			task.Status = TaskStatusFailed
			return fmt.Errorf("task failed: %s", task.Description)
		}
	}

	task.Status = TaskStatusCompleted
	plan.UpdateTaskStatus(task.ID, TaskStatusCompleted, task.Output, task.Error)

	if a.onTaskCompleted != nil {
		a.onTaskCompleted(task)
	}

	return nil
}

// verifyTask verifies a task completed successfully
func (a *Agent) verifyTask(ctx context.Context, task *Task, result *executor.Result) bool {
	v := task.Verification

	switch v.Type {
	case VerificationExitCode:
		return result.ExitCode == v.ExpectedCode

	case VerificationOutput:
		return strings.Contains(result.Output, v.Contains)

	case VerificationCommand:
		// Run the verification command
		verifyResult := a.executor.Execute(ctx, v.Command)
		return verifyResult.Success
	}

	return result.Success
}

// GetCurrentPlan returns the current plan
func (a *Agent) GetCurrentPlan() *Plan {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.currentPlan
}

// GetPlanHistory returns all plans
func (a *Agent) GetPlanHistory() []*Plan {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.planHistory
}

// CancelCurrentPlan cancels the running plan
func (a *Agent) CancelCurrentPlan() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.currentPlan == nil {
		return fmt.Errorf("no active plan")
	}

	if a.currentPlan.Status != PlanStatusRunning {
		return fmt.Errorf("plan is not running")
	}

	a.currentPlan.Status = PlanStatusCancelled
	return nil
}

// GetPlanSummary returns a formatted summary of the current plan
func (a *Agent) GetPlanSummary() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.currentPlan == nil {
		return "No active plan"
	}

	return a.planner.SummarizePlan(a.currentPlan)
}

// SkipTask skips a task in the current plan
func (a *Agent) SkipTask(taskID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.currentPlan == nil {
		return fmt.Errorf("no active plan")
	}

	for i := range a.currentPlan.Tasks {
		if a.currentPlan.Tasks[i].ID == taskID {
			if a.currentPlan.Tasks[i].Status != TaskStatusPending {
				return fmt.Errorf("can only skip pending tasks")
			}
			a.currentPlan.Tasks[i].Status = TaskStatusSkipped
			return nil
		}
	}

	return fmt.Errorf("task not found: %s", taskID)
}

// RetryTask retries a failed task
func (a *Agent) RetryTask(ctx context.Context, taskID string) error {
	a.mu.Lock()
	plan := a.currentPlan
	a.mu.Unlock()

	if plan == nil {
		return fmt.Errorf("no active plan")
	}

	for i := range plan.Tasks {
		if plan.Tasks[i].ID == taskID {
			if plan.Tasks[i].Status != TaskStatusFailed {
				return fmt.Errorf("can only retry failed tasks")
			}
			plan.Tasks[i].Status = TaskStatusPending
			plan.Tasks[i].Output = ""
			plan.Tasks[i].Error = ""

			// Re-execute the task
			return a.executeTask(ctx, plan, &plan.Tasks[i])
		}
	}

	return fmt.Errorf("task not found: %s", taskID)
}

// ModifyTask modifies a task's command before execution
func (a *Agent) ModifyTask(taskID, newCommand string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.currentPlan == nil {
		return fmt.Errorf("no active plan")
	}

	for i := range a.currentPlan.Tasks {
		if a.currentPlan.Tasks[i].ID == taskID {
			if a.currentPlan.Tasks[i].Status != TaskStatusPending {
				return fmt.Errorf("can only modify pending tasks")
			}
			a.currentPlan.Tasks[i].Command = newCommand
			return nil
		}
	}

	return fmt.Errorf("task not found: %s", taskID)
}

// AddTask adds a new task to the current plan
func (a *Agent) AddTask(description, command string, afterTaskID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.currentPlan == nil {
		return fmt.Errorf("no active plan")
	}

	newTask := Task{
		ID:          fmt.Sprintf("task-%d", len(a.currentPlan.Tasks)+1),
		Sequence:    len(a.currentPlan.Tasks) + 1,
		Description: description,
		Command:     command,
		Status:      TaskStatusPending,
	}

	if afterTaskID == "" {
		// Add to end
		a.currentPlan.Tasks = append(a.currentPlan.Tasks, newTask)
	} else {
		// Find position and insert
		for i := range a.currentPlan.Tasks {
			if a.currentPlan.Tasks[i].ID == afterTaskID {
				// Insert after this task
				newTasks := make([]Task, 0, len(a.currentPlan.Tasks)+1)
				newTasks = append(newTasks, a.currentPlan.Tasks[:i+1]...)
				newTasks = append(newTasks, newTask)
				newTasks = append(newTasks, a.currentPlan.Tasks[i+1:]...)
				a.currentPlan.Tasks = newTasks

				// Re-sequence
				for j := range a.currentPlan.Tasks {
					a.currentPlan.Tasks[j].Sequence = j + 1
				}
				return nil
			}
		}
		return fmt.Errorf("task not found: %s", afterTaskID)
	}

	return nil
}

// QuickExecute executes a single command without planning
func (a *Agent) QuickExecute(ctx context.Context, command string) (*executor.Result, error) {
	// Check risk level
	risk := a.executor.AssessRisk(command)
	if risk >= executor.RiskHigh {
		return nil, fmt.Errorf("high-risk command requires explicit approval")
	}

	return a.executor.Execute(ctx, command), nil
}

// ExplainPlan asks the AI to explain the plan
func (a *Agent) ExplainPlan(ctx context.Context) (string, error) {
	a.mu.RLock()
	plan := a.currentPlan
	a.mu.RUnlock()

	if plan == nil {
		return "", fmt.Errorf("no active plan")
	}

	prompt := fmt.Sprintf(`Explain this plan in simple terms:

Goal: %s

Tasks:
`, plan.Goal)

	for _, task := range plan.Tasks {
		prompt += fmt.Sprintf("- %s: %s\n", task.Description, task.Command)
	}

	prompt += "\nExplain what each command does and why it's needed."

	req := &Request{
		UserMessage: prompt,
		Type:        RequestTypeExplain,
	}

	resp, err := a.client.Complete(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.Content, nil
}

// AgentStatus represents the agent's current status
type AgentStatus struct {
	Mode         AgentMode  `json:"mode"`
	HasPlan      bool       `json:"has_plan"`
	PlanStatus   PlanStatus `json:"plan_status,omitempty"`
	PlanProgress float64    `json:"plan_progress,omitempty"`
	CurrentTask  string     `json:"current_task,omitempty"`
	TotalTasks   int        `json:"total_tasks,omitempty"`
}

// GetStatus returns the agent's current status
func (a *Agent) GetStatus() *AgentStatus {
	a.mu.RLock()
	defer a.mu.RUnlock()

	status := &AgentStatus{
		Mode:    a.mode,
		HasPlan: a.currentPlan != nil,
	}

	if a.currentPlan != nil {
		status.PlanStatus = a.currentPlan.Status
		status.PlanProgress = a.currentPlan.GetProgress()
		status.TotalTasks = len(a.currentPlan.Tasks)

		for _, task := range a.currentPlan.Tasks {
			if task.Status == TaskStatusRunning {
				status.CurrentTask = task.Description
				break
			}
		}
	}

	return status
}
