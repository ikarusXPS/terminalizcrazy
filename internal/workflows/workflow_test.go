package workflows

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terminalizcrazy/terminalizcrazy/internal/executor"
)

func TestNewWorkflowEngine(t *testing.T) {
	exec := executor.New()
	engine := NewWorkflowEngine(exec)

	assert.NotNil(t, engine)
	assert.Empty(t, engine.ListWorkflows())
}

func TestWorkflowEngine_RegisterWorkflow(t *testing.T) {
	exec := executor.New()
	engine := NewWorkflowEngine(exec)

	workflow := &Workflow{
		Name:        "test-workflow",
		Description: "Test workflow",
		Steps: []WorkflowStep{
			{Name: "step1", Command: "echo hello"},
		},
	}

	err := engine.RegisterWorkflow(workflow)

	assert.NoError(t, err)
	assert.NotEmpty(t, workflow.ID)
	assert.Len(t, engine.ListWorkflows(), 1)
}

func TestWorkflowEngine_RegisterWorkflow_NoName(t *testing.T) {
	exec := executor.New()
	engine := NewWorkflowEngine(exec)

	workflow := &Workflow{
		Steps: []WorkflowStep{
			{Name: "step1", Command: "echo hello"},
		},
	}

	err := engine.RegisterWorkflow(workflow)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "workflow name is required")
}

func TestWorkflowEngine_RegisterWorkflow_NoSteps(t *testing.T) {
	exec := executor.New()
	engine := NewWorkflowEngine(exec)

	workflow := &Workflow{
		Name: "test-workflow",
	}

	err := engine.RegisterWorkflow(workflow)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must have at least one step")
}

func TestWorkflowEngine_GetWorkflow(t *testing.T) {
	exec := executor.New()
	engine := NewWorkflowEngine(exec)

	workflow := &Workflow{
		Name: "test-workflow",
		Steps: []WorkflowStep{
			{Name: "step1", Command: "echo hello"},
		},
	}
	engine.RegisterWorkflow(workflow)

	retrieved := engine.GetWorkflow("test-workflow")
	assert.NotNil(t, retrieved)
	assert.Equal(t, "test-workflow", retrieved.Name)

	notFound := engine.GetWorkflow("nonexistent")
	assert.Nil(t, notFound)
}

func TestWorkflowEngine_ListWorkflows(t *testing.T) {
	exec := executor.New()
	engine := NewWorkflowEngine(exec)

	wf1 := &Workflow{Name: "wf1", Steps: []WorkflowStep{{Name: "s", Command: "echo"}}}
	wf2 := &Workflow{Name: "wf2", Steps: []WorkflowStep{{Name: "s", Command: "echo"}}}
	engine.RegisterWorkflow(wf1)
	engine.RegisterWorkflow(wf2)

	workflows := engine.ListWorkflows()
	assert.Len(t, workflows, 2)
}

func TestWorkflowEngine_SetCallbacks(t *testing.T) {
	exec := executor.New()
	engine := NewWorkflowEngine(exec)

	startCalled := false
	completeCalled := false
	doneCalled := false

	engine.SetCallbacks(
		func(exec *WorkflowExecution, step *WorkflowStep) { startCalled = true },
		func(exec *WorkflowExecution, result *StepResult) { completeCalled = true },
		func(exec *WorkflowExecution) { doneCalled = true },
	)

	workflow := &Workflow{
		Name: "test",
		Steps: []WorkflowStep{
			{Name: "step1", Command: "echo hello"},
		},
	}
	engine.RegisterWorkflow(workflow)
	engine.Execute(context.Background(), "test", nil)

	assert.True(t, startCalled)
	assert.True(t, completeCalled)
	assert.True(t, doneCalled)
}

func TestWorkflowEngine_Execute_Success(t *testing.T) {
	exec := executor.New()
	engine := NewWorkflowEngine(exec)

	workflow := &Workflow{
		Name: "test",
		Steps: []WorkflowStep{
			{Name: "step1", Command: "echo hello"},
		},
	}
	engine.RegisterWorkflow(workflow)

	execution, err := engine.Execute(context.Background(), "test", nil)

	require.NoError(t, err)
	assert.Equal(t, ExecutionCompleted, execution.Status)
	assert.Len(t, execution.StepResults, 1)
	assert.True(t, execution.StepResults[0].Success)
}

func TestWorkflowEngine_Execute_NotFound(t *testing.T) {
	exec := executor.New()
	engine := NewWorkflowEngine(exec)

	_, err := engine.Execute(context.Background(), "nonexistent", nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "workflow not found")
}

func TestWorkflowEngine_Execute_WithVariables(t *testing.T) {
	exec := executor.New()
	engine := NewWorkflowEngine(exec)

	workflow := &Workflow{
		Name: "test",
		Variables: []Variable{
			{Name: "greeting", Default: "hello"},
		},
		Steps: []WorkflowStep{
			{Name: "step1", Command: "echo ${greeting}"},
		},
	}
	engine.RegisterWorkflow(workflow)

	execution, err := engine.Execute(context.Background(), "test", map[string]string{
		"greeting": "world",
	})

	require.NoError(t, err)
	assert.Equal(t, "world", execution.Variables["greeting"])
}

func TestWorkflowEngine_Execute_RequiredVariable(t *testing.T) {
	exec := executor.New()
	engine := NewWorkflowEngine(exec)

	workflow := &Workflow{
		Name: "test",
		Variables: []Variable{
			{Name: "required_var", Required: true},
		},
		Steps: []WorkflowStep{
			{Name: "step1", Command: "echo ${required_var}"},
		},
	}
	engine.RegisterWorkflow(workflow)

	_, err := engine.Execute(context.Background(), "test", nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required variable")
}

func TestWorkflowEngine_Execute_DefaultVariable(t *testing.T) {
	exec := executor.New()
	engine := NewWorkflowEngine(exec)

	workflow := &Workflow{
		Name: "test",
		Variables: []Variable{
			{Name: "name", Default: "default-value"},
		},
		Steps: []WorkflowStep{
			{Name: "step1", Command: "echo ${name}"},
		},
	}
	engine.RegisterWorkflow(workflow)

	execution, err := engine.Execute(context.Background(), "test", nil)

	require.NoError(t, err)
	assert.Equal(t, "default-value", execution.Variables["name"])
}

func TestWorkflowEngine_Execute_SkippedCondition(t *testing.T) {
	exec := executor.New()
	engine := NewWorkflowEngine(exec)

	workflow := &Workflow{
		Name: "test",
		Steps: []WorkflowStep{
			{Name: "step1", Command: "echo step1"},
			{Name: "step2", Command: "echo step2", Condition: "${nonexistent}"},
		},
	}
	engine.RegisterWorkflow(workflow)

	execution, err := engine.Execute(context.Background(), "test", nil)

	require.NoError(t, err)
	assert.Len(t, execution.StepResults, 2)
	assert.False(t, execution.StepResults[0].Skipped)
	assert.True(t, execution.StepResults[1].Skipped)
}

func TestWorkflowEngine_Execute_ConditionEquality(t *testing.T) {
	exec := executor.New()
	engine := NewWorkflowEngine(exec)

	workflow := &Workflow{
		Name: "test",
		Variables: []Variable{
			{Name: "push", Default: "no"},
		},
		Steps: []WorkflowStep{
			{Name: "step1", Command: "echo always"},
			{Name: "step2", Command: "echo push", Condition: "${push}==yes"},
		},
	}
	engine.RegisterWorkflow(workflow)

	execution, err := engine.Execute(context.Background(), "test", nil)

	require.NoError(t, err)
	assert.False(t, execution.StepResults[0].Skipped)
	assert.True(t, execution.StepResults[1].Skipped)

	// With push=yes
	execution2, err := engine.Execute(context.Background(), "test", map[string]string{"push": "yes"})

	require.NoError(t, err)
	assert.False(t, execution2.StepResults[0].Skipped)
	assert.False(t, execution2.StepResults[1].Skipped)
}

func TestWorkflowEngine_Execute_ConditionInequality(t *testing.T) {
	exec := executor.New()
	engine := NewWorkflowEngine(exec)

	workflow := &Workflow{
		Name: "test",
		Variables: []Variable{
			{Name: "env", Default: "prod"},
		},
		Steps: []WorkflowStep{
			{Name: "step1", Command: "echo debug", Condition: "${env}!=prod"},
		},
	}
	engine.RegisterWorkflow(workflow)

	execution, err := engine.Execute(context.Background(), "test", nil)

	require.NoError(t, err)
	assert.True(t, execution.StepResults[0].Skipped)

	execution2, err := engine.Execute(context.Background(), "test", map[string]string{"env": "dev"})

	require.NoError(t, err)
	assert.False(t, execution2.StepResults[0].Skipped)
}

func TestWorkflowEngine_Execute_CaptureOutput(t *testing.T) {
	exec := executor.New()
	engine := NewWorkflowEngine(exec)

	workflow := &Workflow{
		Name: "test",
		Steps: []WorkflowStep{
			{Name: "capture", Command: "echo captured-value", CaptureAs: "output"},
			{Name: "use", Command: "echo ${output}"},
		},
	}
	engine.RegisterWorkflow(workflow)

	execution, err := engine.Execute(context.Background(), "test", nil)

	require.NoError(t, err)
	assert.Contains(t, execution.Variables["output"], "captured-value")
}

func TestWorkflowEngine_Execute_OnFailStop(t *testing.T) {
	exec := executor.New()
	engine := NewWorkflowEngine(exec)

	workflow := &Workflow{
		Name: "test",
		Steps: []WorkflowStep{
			{Name: "fail", Command: "exit 1", OnFail: OnFailStop},
			{Name: "after", Command: "echo after"},
		},
	}
	engine.RegisterWorkflow(workflow)

	execution, err := engine.Execute(context.Background(), "test", nil)

	assert.Error(t, err)
	assert.Equal(t, ExecutionFailed, execution.Status)
	assert.Len(t, execution.StepResults, 1)
}

func TestWorkflowEngine_Execute_OnFailContinue(t *testing.T) {
	exec := executor.New()
	engine := NewWorkflowEngine(exec)

	workflow := &Workflow{
		Name: "test",
		Steps: []WorkflowStep{
			{Name: "fail", Command: "exit 1", OnFail: OnFailContinue},
			{Name: "after", Command: "echo after"},
		},
	}
	engine.RegisterWorkflow(workflow)

	execution, err := engine.Execute(context.Background(), "test", nil)

	assert.NoError(t, err)
	assert.Equal(t, ExecutionCompleted, execution.Status)
	assert.Len(t, execution.StepResults, 2)
}

func TestWorkflowEngine_Execute_OnFailSkip(t *testing.T) {
	exec := executor.New()
	engine := NewWorkflowEngine(exec)

	workflow := &Workflow{
		Name: "test",
		Steps: []WorkflowStep{
			{Name: "fail", Command: "exit 1", OnFail: OnFailSkip},
			{Name: "after", Command: "echo after"},
		},
	}
	engine.RegisterWorkflow(workflow)

	execution, err := engine.Execute(context.Background(), "test", nil)

	assert.NoError(t, err)
	assert.Equal(t, ExecutionCompleted, execution.Status)
	assert.Len(t, execution.StepResults, 2)
}

func TestWorkflowEngine_GetExecution(t *testing.T) {
	exec := executor.New()
	engine := NewWorkflowEngine(exec)

	workflow := &Workflow{
		Name:  "test",
		Steps: []WorkflowStep{{Name: "s", Command: "echo"}},
	}
	engine.RegisterWorkflow(workflow)

	execution, _ := engine.Execute(context.Background(), "test", nil)

	retrieved := engine.GetExecution(execution.ID)
	assert.NotNil(t, retrieved)
	assert.Equal(t, execution.ID, retrieved.ID)
}

func TestWorkflowEngine_ListExecutions(t *testing.T) {
	exec := executor.New()
	engine := NewWorkflowEngine(exec)

	workflow := &Workflow{
		Name:  "test",
		Steps: []WorkflowStep{{Name: "s", Command: "echo"}},
	}
	engine.RegisterWorkflow(workflow)

	engine.Execute(context.Background(), "test", nil)
	engine.Execute(context.Background(), "test", nil)

	executions := engine.ListExecutions()
	assert.Len(t, executions, 2)
}

func TestWorkflow_ToJSON(t *testing.T) {
	workflow := &Workflow{
		Name:        "test",
		Description: "Test workflow",
		Steps: []WorkflowStep{
			{Name: "step1", Command: "echo hello"},
		},
	}

	json, err := workflow.ToJSON()

	assert.NoError(t, err)
	assert.Contains(t, json, "test")
	assert.Contains(t, json, "step1")
}

func TestWorkflowFromJSON(t *testing.T) {
	json := `{
		"name": "test",
		"description": "Test workflow",
		"steps": [
			{"name": "step1", "command": "echo hello"}
		]
	}`

	workflow, err := WorkflowFromJSON(json)

	assert.NoError(t, err)
	assert.Equal(t, "test", workflow.Name)
	assert.Len(t, workflow.Steps, 1)
}

func TestWorkflowFromJSON_Invalid(t *testing.T) {
	json := `invalid json`

	_, err := WorkflowFromJSON(json)

	assert.Error(t, err)
}

func TestWorkflow_Summary(t *testing.T) {
	workflow := &Workflow{
		Name: "test-workflow",
		Steps: []WorkflowStep{
			{Name: "step1", Command: "echo"},
			{Name: "step2", Command: "echo"},
			{Name: "step3", Command: "echo"},
		},
	}

	summary := workflow.Summary()

	assert.Equal(t, "test-workflow (3 steps)", summary)
}

func TestWorkflowExecution_ExecutionSummary(t *testing.T) {
	execution := &WorkflowExecution{
		Status: ExecutionCompleted,
		StepResults: []StepResult{
			{Success: true},
			{Success: true},
			{Success: false, Skipped: false},
			{Success: false, Skipped: true},
		},
	}

	summary := execution.ExecutionSummary()

	assert.Contains(t, summary, "completed")
	assert.Contains(t, summary, "2 completed")
	assert.Contains(t, summary, "1 failed")
}

func TestOnFailActions(t *testing.T) {
	assert.Equal(t, OnFailAction("stop"), OnFailStop)
	assert.Equal(t, OnFailAction("skip"), OnFailSkip)
	assert.Equal(t, OnFailAction("continue"), OnFailContinue)
	assert.Equal(t, OnFailAction("retry"), OnFailRetry)
}

func TestExecutionStatuses(t *testing.T) {
	assert.Equal(t, ExecutionStatus("pending"), ExecutionPending)
	assert.Equal(t, ExecutionStatus("running"), ExecutionRunning)
	assert.Equal(t, ExecutionStatus("completed"), ExecutionCompleted)
	assert.Equal(t, ExecutionStatus("failed"), ExecutionFailed)
	assert.Equal(t, ExecutionStatus("cancelled"), ExecutionCancelled)
}

func TestSubstituteVariables(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		vars   map[string]string
		expect string
	}{
		{
			name:   "single variable",
			text:   "echo ${name}",
			vars:   map[string]string{"name": "world"},
			expect: "echo world",
		},
		{
			name:   "multiple variables",
			text:   "${greeting} ${name}!",
			vars:   map[string]string{"greeting": "Hello", "name": "World"},
			expect: "Hello World!",
		},
		{
			name:   "missing variable",
			text:   "echo ${missing}",
			vars:   map[string]string{},
			expect: "echo ${missing}",
		},
		{
			name:   "no variables",
			text:   "echo hello",
			vars:   map[string]string{},
			expect: "echo hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := substituteVariables(tt.text, tt.vars)
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestWorkflowStep(t *testing.T) {
	step := WorkflowStep{
		Name:        "test-step",
		Description: "A test step",
		Command:     "echo hello",
		OnFail:      OnFailStop,
		Condition:   "${var}==yes",
		Timeout:     30 * time.Second,
		Retries:     3,
		CaptureAs:   "output",
	}

	assert.Equal(t, "test-step", step.Name)
	assert.Equal(t, "A test step", step.Description)
	assert.Equal(t, "echo hello", step.Command)
	assert.Equal(t, OnFailStop, step.OnFail)
	assert.Equal(t, "${var}==yes", step.Condition)
	assert.Equal(t, 30*time.Second, step.Timeout)
	assert.Equal(t, 3, step.Retries)
	assert.Equal(t, "output", step.CaptureAs)
}

func TestVariable(t *testing.T) {
	v := Variable{
		Name:        "test_var",
		Description: "A test variable",
		Default:     "default_value",
		Required:    true,
		Type:        "string",
	}

	assert.Equal(t, "test_var", v.Name)
	assert.Equal(t, "A test variable", v.Description)
	assert.Equal(t, "default_value", v.Default)
	assert.True(t, v.Required)
	assert.Equal(t, "string", v.Type)
}

func TestStepResult(t *testing.T) {
	result := StepResult{
		StepName:   "step1",
		Command:    "echo hello",
		Output:     "hello",
		Error:      "",
		ExitCode:   0,
		Success:    true,
		Duration:   100 * time.Millisecond,
		Skipped:    false,
		RetryCount: 0,
	}

	assert.Equal(t, "step1", result.StepName)
	assert.Equal(t, "echo hello", result.Command)
	assert.Equal(t, "hello", result.Output)
	assert.Empty(t, result.Error)
	assert.Equal(t, 0, result.ExitCode)
	assert.True(t, result.Success)
	assert.Equal(t, 100*time.Millisecond, result.Duration)
	assert.False(t, result.Skipped)
	assert.Equal(t, 0, result.RetryCount)
}
