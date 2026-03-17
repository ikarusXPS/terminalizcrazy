package ai

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/terminalizcrazy/terminalizcrazy/internal/executor"
)

// MockClient implements the Client interface for testing
type MockClient struct {
	responses    []*Response
	responseIdx  int
	lastRequest  *Request
	shouldError  bool
	errorMessage string
}

func NewMockClient() *MockClient {
	return &MockClient{
		responses: make([]*Response, 0),
	}
}

func (m *MockClient) Complete(ctx context.Context, req *Request) (*Response, error) {
	m.lastRequest = req
	if m.shouldError {
		return nil, assert.AnError
	}
	if m.responseIdx < len(m.responses) {
		resp := m.responses[m.responseIdx]
		m.responseIdx++
		return resp, nil
	}
	return &Response{Content: "default response"}, nil
}

func (m *MockClient) Provider() Provider {
	return "mock"
}

func (m *MockClient) AddResponse(content string) {
	m.responses = append(m.responses, &Response{Content: content})
}

func (m *MockClient) SetError(msg string) {
	m.shouldError = true
	m.errorMessage = msg
}

func TestAgentModeConstants(t *testing.T) {
	assert.Equal(t, AgentMode("off"), AgentModeOff)
	assert.Equal(t, AgentMode("suggest"), AgentModeSuggest)
	assert.Equal(t, AgentMode("auto"), AgentModeAuto)
}

func TestDefaultAgentConfig(t *testing.T) {
	config := DefaultAgentConfig()

	assert.NotNil(t, config)
	assert.Equal(t, AgentModeSuggest, config.Mode)
	assert.Equal(t, 10, config.MaxTasksPerPlan)
	assert.Equal(t, 60*1000000000, int(config.TaskTimeout.Nanoseconds()))
	assert.True(t, config.RequireApproval)
	assert.False(t, config.AllowDangerous)
	assert.Equal(t, 2, config.MaxRetries)
}

func TestNewAgent(t *testing.T) {
	client := NewMockClient()
	exec := executor.New()

	t.Run("with default config", func(t *testing.T) {
		agent := NewAgent(client, exec, nil)

		assert.NotNil(t, agent)
		assert.Equal(t, AgentModeSuggest, agent.GetMode())
		assert.Nil(t, agent.GetCurrentPlan())
	})

	t.Run("with custom config", func(t *testing.T) {
		config := &AgentConfig{
			Mode: AgentModeAuto,
		}
		agent := NewAgent(client, exec, config)

		assert.NotNil(t, agent)
		assert.Equal(t, AgentModeAuto, agent.GetMode())
	})
}

func TestAgent_SetGetMode(t *testing.T) {
	client := NewMockClient()
	exec := executor.New()
	agent := NewAgent(client, exec, nil)

	agent.SetMode(AgentModeOff)
	assert.Equal(t, AgentModeOff, agent.GetMode())

	agent.SetMode(AgentModeAuto)
	assert.Equal(t, AgentModeAuto, agent.GetMode())

	agent.SetMode(AgentModeSuggest)
	assert.Equal(t, AgentModeSuggest, agent.GetMode())
}

func TestAgent_SetCallbacks(t *testing.T) {
	client := NewMockClient()
	exec := executor.New()
	agent := NewAgent(client, exec, nil)

	planCreatedCalled := false
	taskStartedCalled := false
	taskCompletedCalled := false
	taskFailedCalled := false
	planCompletedCalled := false
	approvalNeededCalled := false

	agent.SetCallbacks(
		func(p *Plan) { planCreatedCalled = true },
		func(t *Task) { taskStartedCalled = true },
		func(t *Task) { taskCompletedCalled = true },
		func(t *Task) { taskFailedCalled = true },
		func(p *Plan) { planCompletedCalled = true },
		func(p *Plan) { approvalNeededCalled = true },
	)

	assert.NotNil(t, agent.onPlanCreated)
	assert.NotNil(t, agent.onTaskStarted)
	assert.NotNil(t, agent.onTaskCompleted)
	assert.NotNil(t, agent.onTaskFailed)
	assert.NotNil(t, agent.onPlanCompleted)
	assert.NotNil(t, agent.onApprovalNeeded)

	// Test callbacks work
	agent.onPlanCreated(&Plan{})
	agent.onTaskStarted(&Task{})
	agent.onTaskCompleted(&Task{})
	agent.onTaskFailed(&Task{})
	agent.onPlanCompleted(&Plan{})
	agent.onApprovalNeeded(&Plan{})

	assert.True(t, planCreatedCalled)
	assert.True(t, taskStartedCalled)
	assert.True(t, taskCompletedCalled)
	assert.True(t, taskFailedCalled)
	assert.True(t, planCompletedCalled)
	assert.True(t, approvalNeededCalled)
}

func TestAgent_GetCurrentPlan(t *testing.T) {
	client := NewMockClient()
	exec := executor.New()
	agent := NewAgent(client, exec, nil)

	// Initially no plan
	assert.Nil(t, agent.GetCurrentPlan())

	// After setting a plan
	plan := &Plan{ID: "test-plan", Goal: "test goal"}
	agent.currentPlan = plan

	current := agent.GetCurrentPlan()
	assert.NotNil(t, current)
	assert.Equal(t, "test-plan", current.ID)
}

func TestAgent_GetPlanHistory(t *testing.T) {
	client := NewMockClient()
	exec := executor.New()
	agent := NewAgent(client, exec, nil)

	// Initially empty
	history := agent.GetPlanHistory()
	assert.Empty(t, history)

	// After adding plans
	agent.planHistory = append(agent.planHistory, &Plan{ID: "plan-1"})
	agent.planHistory = append(agent.planHistory, &Plan{ID: "plan-2"})

	history = agent.GetPlanHistory()
	assert.Len(t, history, 2)
}

func TestAgent_ApprovePlan(t *testing.T) {
	client := NewMockClient()
	exec := executor.New()
	agent := NewAgent(client, exec, nil)

	t.Run("no current plan", func(t *testing.T) {
		err := agent.ApprovePlan(context.Background(), "nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "plan not found")
	})

	t.Run("wrong plan ID", func(t *testing.T) {
		agent.currentPlan = &Plan{ID: "plan-1", Status: PlanStatusPending}
		err := agent.ApprovePlan(context.Background(), "plan-2")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "plan not found")
	})

	t.Run("plan not pending", func(t *testing.T) {
		agent.currentPlan = &Plan{ID: "plan-1", Status: PlanStatusRunning}
		err := agent.ApprovePlan(context.Background(), "plan-1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not pending")
	})

	t.Run("successful approval", func(t *testing.T) {
		agent.currentPlan = &Plan{
			ID:     "plan-1",
			Status: PlanStatusPending,
			Tasks:  []Task{}, // Empty tasks so execution completes quickly
		}
		err := agent.ApprovePlan(context.Background(), "plan-1")
		assert.NoError(t, err)
		assert.Equal(t, PlanStatusApproved, agent.currentPlan.Status)
	})
}

func TestAgent_RejectPlan(t *testing.T) {
	client := NewMockClient()
	exec := executor.New()
	agent := NewAgent(client, exec, nil)

	t.Run("no current plan", func(t *testing.T) {
		err := agent.RejectPlan("nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "plan not found")
	})

	t.Run("wrong plan ID", func(t *testing.T) {
		agent.currentPlan = &Plan{ID: "plan-1"}
		err := agent.RejectPlan("plan-2")
		assert.Error(t, err)
	})

	t.Run("successful rejection", func(t *testing.T) {
		agent.currentPlan = &Plan{ID: "plan-1", Status: PlanStatusPending}
		err := agent.RejectPlan("plan-1")
		assert.NoError(t, err)
		assert.Equal(t, PlanStatusCancelled, agent.currentPlan.Status)
	})
}

func TestAgent_CancelCurrentPlan(t *testing.T) {
	client := NewMockClient()
	exec := executor.New()
	agent := NewAgent(client, exec, nil)

	t.Run("no active plan", func(t *testing.T) {
		err := agent.CancelCurrentPlan()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no active plan")
	})

	t.Run("plan not running", func(t *testing.T) {
		agent.currentPlan = &Plan{ID: "plan-1", Status: PlanStatusPending}
		err := agent.CancelCurrentPlan()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not running")
	})

	t.Run("successful cancel", func(t *testing.T) {
		agent.currentPlan = &Plan{ID: "plan-1", Status: PlanStatusRunning}
		err := agent.CancelCurrentPlan()
		assert.NoError(t, err)
		assert.Equal(t, PlanStatusCancelled, agent.currentPlan.Status)
	})
}

func TestAgent_SkipTask(t *testing.T) {
	client := NewMockClient()
	exec := executor.New()
	agent := NewAgent(client, exec, nil)

	t.Run("no active plan", func(t *testing.T) {
		err := agent.SkipTask("task-1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no active plan")
	})

	t.Run("task not found", func(t *testing.T) {
		agent.currentPlan = &Plan{
			ID:    "plan-1",
			Tasks: []Task{{ID: "task-1", Status: TaskStatusPending}},
		}
		err := agent.SkipTask("task-99")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "task not found")
	})

	t.Run("task not pending", func(t *testing.T) {
		agent.currentPlan = &Plan{
			ID:    "plan-1",
			Tasks: []Task{{ID: "task-1", Status: TaskStatusCompleted}},
		}
		err := agent.SkipTask("task-1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only skip pending")
	})

	t.Run("successful skip", func(t *testing.T) {
		agent.currentPlan = &Plan{
			ID:    "plan-1",
			Tasks: []Task{{ID: "task-1", Status: TaskStatusPending}},
		}
		err := agent.SkipTask("task-1")
		assert.NoError(t, err)
		assert.Equal(t, TaskStatusSkipped, agent.currentPlan.Tasks[0].Status)
	})
}

func TestAgent_ModifyTask(t *testing.T) {
	client := NewMockClient()
	exec := executor.New()
	agent := NewAgent(client, exec, nil)

	t.Run("no active plan", func(t *testing.T) {
		err := agent.ModifyTask("task-1", "new command")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no active plan")
	})

	t.Run("task not found", func(t *testing.T) {
		agent.currentPlan = &Plan{
			ID:    "plan-1",
			Tasks: []Task{{ID: "task-1", Command: "old", Status: TaskStatusPending}},
		}
		err := agent.ModifyTask("task-99", "new")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "task not found")
	})

	t.Run("task not pending", func(t *testing.T) {
		agent.currentPlan = &Plan{
			ID:    "plan-1",
			Tasks: []Task{{ID: "task-1", Command: "old", Status: TaskStatusRunning}},
		}
		err := agent.ModifyTask("task-1", "new")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only modify pending")
	})

	t.Run("successful modify", func(t *testing.T) {
		agent.currentPlan = &Plan{
			ID:    "plan-1",
			Tasks: []Task{{ID: "task-1", Command: "old command", Status: TaskStatusPending}},
		}
		err := agent.ModifyTask("task-1", "new command")
		assert.NoError(t, err)
		assert.Equal(t, "new command", agent.currentPlan.Tasks[0].Command)
	})
}

func TestAgent_AddTask(t *testing.T) {
	client := NewMockClient()
	exec := executor.New()
	agent := NewAgent(client, exec, nil)

	t.Run("no active plan", func(t *testing.T) {
		err := agent.AddTask("desc", "cmd", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no active plan")
	})

	t.Run("add to end", func(t *testing.T) {
		agent.currentPlan = &Plan{
			ID: "plan-1",
			Tasks: []Task{
				{ID: "task-1", Sequence: 1, Description: "Task 1"},
			},
		}
		err := agent.AddTask("New Task", "new command", "")
		assert.NoError(t, err)
		assert.Len(t, agent.currentPlan.Tasks, 2)
		assert.Equal(t, "New Task", agent.currentPlan.Tasks[1].Description)
	})

	t.Run("add after specific task", func(t *testing.T) {
		agent.currentPlan = &Plan{
			ID: "plan-1",
			Tasks: []Task{
				{ID: "task-1", Sequence: 1, Description: "Task 1"},
				{ID: "task-2", Sequence: 2, Description: "Task 2"},
			},
		}
		err := agent.AddTask("Middle Task", "middle command", "task-1")
		assert.NoError(t, err)
		assert.Len(t, agent.currentPlan.Tasks, 3)
		assert.Equal(t, "Middle Task", agent.currentPlan.Tasks[1].Description)
		assert.Equal(t, 2, agent.currentPlan.Tasks[1].Sequence)
	})

	t.Run("add after nonexistent task", func(t *testing.T) {
		agent.currentPlan = &Plan{
			ID:    "plan-1",
			Tasks: []Task{{ID: "task-1"}},
		}
		err := agent.AddTask("New", "cmd", "task-99")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "task not found")
	})
}

func TestAgent_GetPlanSummary(t *testing.T) {
	client := NewMockClient()
	exec := executor.New()
	agent := NewAgent(client, exec, nil)

	t.Run("no active plan", func(t *testing.T) {
		summary := agent.GetPlanSummary()
		assert.Equal(t, "No active plan", summary)
	})

	t.Run("with active plan", func(t *testing.T) {
		agent.currentPlan = &Plan{
			ID:     "plan-1",
			Goal:   "Test goal",
			Status: PlanStatusPending,
			Tasks: []Task{
				{ID: "task-1", Sequence: 1, Description: "First task", Command: "echo hello", Status: TaskStatusPending},
			},
		}
		summary := agent.GetPlanSummary()
		assert.Contains(t, summary, "Test goal")
		assert.Contains(t, summary, "First task")
	})
}

func TestAgent_GetStatus(t *testing.T) {
	client := NewMockClient()
	exec := executor.New()
	agent := NewAgent(client, exec, nil)

	t.Run("no plan", func(t *testing.T) {
		status := agent.GetStatus()
		assert.Equal(t, AgentModeSuggest, status.Mode)
		assert.False(t, status.HasPlan)
	})

	t.Run("with plan", func(t *testing.T) {
		agent.currentPlan = &Plan{
			ID:     "plan-1",
			Status: PlanStatusRunning,
			Tasks: []Task{
				{ID: "task-1", Description: "First", Status: TaskStatusCompleted},
				{ID: "task-2", Description: "Second", Status: TaskStatusRunning},
				{ID: "task-3", Description: "Third", Status: TaskStatusPending},
			},
		}

		status := agent.GetStatus()
		assert.True(t, status.HasPlan)
		assert.Equal(t, PlanStatusRunning, status.PlanStatus)
		assert.Equal(t, 3, status.TotalTasks)
		assert.Equal(t, "Second", status.CurrentTask)
	})
}

func TestAgent_QuickExecute(t *testing.T) {
	client := NewMockClient()
	exec := executor.New()
	agent := NewAgent(client, exec, nil)

	t.Run("safe command", func(t *testing.T) {
		result, err := agent.QuickExecute(context.Background(), "echo hello")
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("high risk command", func(t *testing.T) {
		result, err := agent.QuickExecute(context.Background(), "rm -rf /")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "high-risk")
	})
}

func TestAgent_RetryTask(t *testing.T) {
	client := NewMockClient()
	exec := executor.New()
	agent := NewAgent(client, exec, nil)

	t.Run("no active plan", func(t *testing.T) {
		err := agent.RetryTask(context.Background(), "task-1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no active plan")
	})

	t.Run("task not found", func(t *testing.T) {
		agent.currentPlan = &Plan{
			ID:    "plan-1",
			Tasks: []Task{{ID: "task-1", Status: TaskStatusFailed}},
		}
		err := agent.RetryTask(context.Background(), "task-99")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "task not found")
	})

	t.Run("task not failed", func(t *testing.T) {
		agent.currentPlan = &Plan{
			ID:    "plan-1",
			Tasks: []Task{{ID: "task-1", Status: TaskStatusCompleted}},
		}
		err := agent.RetryTask(context.Background(), "task-1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only retry failed")
	})
}

func TestAgent_ExecutePlan(t *testing.T) {
	client := NewMockClient()
	exec := executor.New()
	agent := NewAgent(client, exec, nil)

	t.Run("empty plan", func(t *testing.T) {
		plan := &Plan{
			ID:     "plan-1",
			Goal:   "Test",
			Status: PlanStatusApproved,
			Tasks:  []Task{},
		}

		err := agent.ExecutePlan(context.Background(), plan)
		assert.NoError(t, err)
		assert.Equal(t, PlanStatusCompleted, plan.Status)
	})

	t.Run("plan with simple tasks", func(t *testing.T) {
		plan := &Plan{
			ID:     "plan-1",
			Goal:   "Test",
			Status: PlanStatusApproved,
			Tasks: []Task{
				{ID: "task-1", Command: "echo hello", Status: TaskStatusPending},
			},
		}

		err := agent.ExecutePlan(context.Background(), plan)
		assert.NoError(t, err)
		assert.Equal(t, PlanStatusCompleted, plan.Status)
		assert.Equal(t, TaskStatusCompleted, plan.Tasks[0].Status)
	})
}
