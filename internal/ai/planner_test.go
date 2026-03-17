package ai

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlanStatusConstants(t *testing.T) {
	assert.Equal(t, PlanStatus("pending"), PlanStatusPending)
	assert.Equal(t, PlanStatus("approved"), PlanStatusApproved)
	assert.Equal(t, PlanStatus("running"), PlanStatusRunning)
	assert.Equal(t, PlanStatus("completed"), PlanStatusCompleted)
	assert.Equal(t, PlanStatus("failed"), PlanStatusFailed)
	assert.Equal(t, PlanStatus("cancelled"), PlanStatusCancelled)
}

func TestTaskStatusConstants(t *testing.T) {
	assert.Equal(t, TaskStatus("pending"), TaskStatusPending)
	assert.Equal(t, TaskStatus("running"), TaskStatusRunning)
	assert.Equal(t, TaskStatus("completed"), TaskStatusCompleted)
	assert.Equal(t, TaskStatus("failed"), TaskStatusFailed)
	assert.Equal(t, TaskStatus("skipped"), TaskStatusSkipped)
}

func TestVerificationTypeConstants(t *testing.T) {
	assert.Equal(t, VerificationType("exit_code"), VerificationExitCode)
	assert.Equal(t, VerificationType("output_contains"), VerificationOutput)
	assert.Equal(t, VerificationType("run_command"), VerificationCommand)
}

func TestNewPlanner(t *testing.T) {
	client := NewMockClient()
	planner := NewPlanner(client)

	assert.NotNil(t, planner)
	assert.Equal(t, client, planner.client)
}

func TestPlanner_ValidatePlan(t *testing.T) {
	client := NewMockClient()
	planner := NewPlanner(client)

	t.Run("nil plan", func(t *testing.T) {
		err := planner.ValidatePlan(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nil")
	})

	t.Run("empty goal", func(t *testing.T) {
		plan := &Plan{
			ID:    "plan-1",
			Goal:  "",
			Tasks: []Task{{ID: "task-1"}},
		}
		err := planner.ValidatePlan(plan)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no goal")
	})

	t.Run("no tasks", func(t *testing.T) {
		plan := &Plan{
			ID:    "plan-1",
			Goal:  "Test goal",
			Tasks: []Task{},
		}
		err := planner.ValidatePlan(plan)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no tasks")
	})

	t.Run("dangerous command", func(t *testing.T) {
		plan := &Plan{
			ID:   "plan-1",
			Goal: "Test goal",
			Tasks: []Task{
				{ID: "task-1", Sequence: 1, Command: "rm -rf /"},
			},
		}
		err := planner.ValidatePlan(plan)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "dangerous")
	})

	t.Run("valid plan", func(t *testing.T) {
		plan := &Plan{
			ID:   "plan-1",
			Goal: "Test goal",
			Tasks: []Task{
				{ID: "task-1", Sequence: 1, Command: "echo hello"},
			},
		}
		err := planner.ValidatePlan(plan)
		assert.NoError(t, err)
	})
}

func TestPlanner_SummarizePlan(t *testing.T) {
	client := NewMockClient()
	planner := NewPlanner(client)

	t.Run("basic plan", func(t *testing.T) {
		plan := &Plan{
			ID:     "plan-1",
			Goal:   "Build project",
			Status: PlanStatusPending,
			Tasks: []Task{
				{ID: "task-1", Sequence: 1, Description: "Install deps", Command: "npm install", Status: TaskStatusPending},
				{ID: "task-2", Sequence: 2, Description: "Build", Command: "npm run build", Status: TaskStatusPending},
			},
		}

		summary := planner.SummarizePlan(plan)

		assert.Contains(t, summary, "Build project")
		assert.Contains(t, summary, "pending")
		assert.Contains(t, summary, "Tasks: 2")
		assert.Contains(t, summary, "Install deps")
		assert.Contains(t, summary, "npm install")
		assert.Contains(t, summary, "Build")
		assert.Contains(t, summary, "npm run build")
	})

	t.Run("plan with output", func(t *testing.T) {
		plan := &Plan{
			ID:     "plan-1",
			Goal:   "Test",
			Status: PlanStatusRunning,
			Tasks: []Task{
				{
					ID:          "task-1",
					Sequence:    1,
					Description: "Echo",
					Command:     "echo hello",
					Status:      TaskStatusCompleted,
					Output:      "hello",
				},
			},
		}

		summary := planner.SummarizePlan(plan)
		assert.Contains(t, summary, "Output: hello")
	})

	t.Run("plan with error", func(t *testing.T) {
		plan := &Plan{
			ID:     "plan-1",
			Goal:   "Test",
			Status: PlanStatusFailed,
			Tasks: []Task{
				{
					ID:          "task-1",
					Sequence:    1,
					Description: "Fail",
					Command:     "exit 1",
					Status:      TaskStatusFailed,
					Error:       "command failed",
				},
			},
		}

		summary := planner.SummarizePlan(plan)
		assert.Contains(t, summary, "Error: command failed")
	})

	t.Run("plan with long output", func(t *testing.T) {
		plan := &Plan{
			ID:     "plan-1",
			Goal:   "Test",
			Status: PlanStatusCompleted,
			Tasks: []Task{
				{
					ID:          "task-1",
					Sequence:    1,
					Description: "Multi-line",
					Command:     "echo",
					Status:      TaskStatusCompleted,
					Output:      "line1\nline2\nline3\nline4\nline5",
				},
			},
		}

		summary := planner.SummarizePlan(plan)
		assert.Contains(t, summary, "more lines")
	})
}

func TestExtractJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple JSON",
			input:    `{"tasks": []}`,
			expected: `{"tasks": []}`,
		},
		{
			name:     "JSON with text before",
			input:    `Here is the plan: {"tasks": [{"description": "test"}]}`,
			expected: `{"tasks": [{"description": "test"}]}`,
		},
		{
			name:     "JSON with text after",
			input:    `{"tasks": []} Let me know if you need more.`,
			expected: `{"tasks": []}`,
		},
		{
			name:     "nested JSON",
			input:    `{"outer": {"inner": {"deep": true}}}`,
			expected: `{"outer": {"inner": {"deep": true}}}`,
		},
		{
			name:     "no JSON",
			input:    "This is just plain text",
			expected: "",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "unclosed brace",
			input:    "{incomplete",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractJSON(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsDangerousCommand(t *testing.T) {
	tests := []struct {
		name      string
		command   string
		dangerous bool
	}{
		// Dangerous commands
		{"rm -rf /", "rm -rf /", true},
		{"rm -rf /*", "rm -rf /*", true},
		{"fork bomb", ":(){ :|:& };:", true},
		{"dd zero", "dd if=/dev/zero of=/dev/sda", true},
		{"mkfs", "mkfs.ext4 /dev/sda1", true},
		{"format C", "format c:", true},
		{"del /f /s", "del /f /s /q c:\\", true},
		{"rm with injection", "$(rm -rf /home)", true},
		{"rm with backtick", "`rm -rf /tmp`", true},
		{"rm with semicolon", "ls; rm -rf /", true},

		// Safe commands
		{"echo", "echo hello", false},
		{"ls", "ls -la", false},
		{"git status", "git status", false},
		{"npm install", "npm install", false},
		{"cd", "cd /home", false},
		{"mkdir", "mkdir new_folder", false},
		{"cat", "cat file.txt", false},
		{"grep", "grep pattern file", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isDangerousCommand(tt.command)
			assert.Equal(t, tt.dangerous, result)
		})
	}
}

func TestGetStatusIcon(t *testing.T) {
	tests := []struct {
		status   TaskStatus
		expected string
	}{
		{TaskStatusPending, "○"},
		{TaskStatusRunning, "◐"},
		{TaskStatusCompleted, "✓"},
		{TaskStatusFailed, "✗"},
		{TaskStatusSkipped, "⊘"},
		{TaskStatus("unknown"), "?"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			result := getStatusIcon(tt.status)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPlan_GetNextTask(t *testing.T) {
	t.Run("with pending tasks", func(t *testing.T) {
		plan := &Plan{
			Tasks: []Task{
				{ID: "task-1", Status: TaskStatusCompleted},
				{ID: "task-2", Status: TaskStatusPending},
				{ID: "task-3", Status: TaskStatusPending},
			},
		}

		task := plan.GetNextTask()
		require.NotNil(t, task)
		assert.Equal(t, "task-2", task.ID)
	})

	t.Run("all completed", func(t *testing.T) {
		plan := &Plan{
			Tasks: []Task{
				{ID: "task-1", Status: TaskStatusCompleted},
				{ID: "task-2", Status: TaskStatusCompleted},
			},
		}

		task := plan.GetNextTask()
		assert.Nil(t, task)
	})

	t.Run("empty tasks", func(t *testing.T) {
		plan := &Plan{Tasks: []Task{}}
		task := plan.GetNextTask()
		assert.Nil(t, task)
	})
}

func TestPlan_IsComplete(t *testing.T) {
	tests := []struct {
		name     string
		tasks    []Task
		complete bool
	}{
		{
			name:     "empty tasks",
			tasks:    []Task{},
			complete: true,
		},
		{
			name: "all completed",
			tasks: []Task{
				{Status: TaskStatusCompleted},
				{Status: TaskStatusCompleted},
			},
			complete: true,
		},
		{
			name: "all skipped",
			tasks: []Task{
				{Status: TaskStatusSkipped},
			},
			complete: true,
		},
		{
			name: "mixed completed and skipped",
			tasks: []Task{
				{Status: TaskStatusCompleted},
				{Status: TaskStatusSkipped},
			},
			complete: true,
		},
		{
			name: "has pending",
			tasks: []Task{
				{Status: TaskStatusCompleted},
				{Status: TaskStatusPending},
			},
			complete: false,
		},
		{
			name: "has running",
			tasks: []Task{
				{Status: TaskStatusRunning},
			},
			complete: false,
		},
		{
			name: "has failed",
			tasks: []Task{
				{Status: TaskStatusFailed},
			},
			complete: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := &Plan{Tasks: tt.tasks}
			assert.Equal(t, tt.complete, plan.IsComplete())
		})
	}
}

func TestPlan_GetProgress(t *testing.T) {
	tests := []struct {
		name     string
		tasks    []Task
		expected float64
	}{
		{
			name:     "empty tasks",
			tasks:    []Task{},
			expected: 0,
		},
		{
			name: "all completed",
			tasks: []Task{
				{Status: TaskStatusCompleted},
				{Status: TaskStatusCompleted},
			},
			expected: 100,
		},
		{
			name: "half completed",
			tasks: []Task{
				{Status: TaskStatusCompleted},
				{Status: TaskStatusPending},
			},
			expected: 50,
		},
		{
			name: "one of three",
			tasks: []Task{
				{Status: TaskStatusCompleted},
				{Status: TaskStatusPending},
				{Status: TaskStatusPending},
			},
			expected: 100.0 / 3,
		},
		{
			name: "skipped counts as completed",
			tasks: []Task{
				{Status: TaskStatusSkipped},
				{Status: TaskStatusPending},
			},
			expected: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := &Plan{Tasks: tt.tasks}
			result := plan.GetProgress()
			assert.InDelta(t, tt.expected, result, 0.01)
		})
	}
}

func TestPlan_UpdateTaskStatus(t *testing.T) {
	plan := &Plan{
		Tasks: []Task{
			{ID: "task-1", Status: TaskStatusPending},
			{ID: "task-2", Status: TaskStatusPending},
		},
	}

	plan.UpdateTaskStatus("task-1", TaskStatusCompleted, "output here", "")

	assert.Equal(t, TaskStatusCompleted, plan.Tasks[0].Status)
	assert.Equal(t, "output here", plan.Tasks[0].Output)
	assert.Equal(t, "", plan.Tasks[0].Error)

	// Update with error
	plan.UpdateTaskStatus("task-2", TaskStatusFailed, "", "error occurred")

	assert.Equal(t, TaskStatusFailed, plan.Tasks[1].Status)
	assert.Equal(t, "error occurred", plan.Tasks[1].Error)

	// Update nonexistent task (should not panic)
	plan.UpdateTaskStatus("task-99", TaskStatusCompleted, "", "")
	assert.Len(t, plan.Tasks, 2)
}

func TestNewCommandExtractor(t *testing.T) {
	extractor := NewCommandExtractor()

	assert.NotNil(t, extractor)
	assert.NotEmpty(t, extractor.patterns)
}

func TestCommandExtractor_Extract(t *testing.T) {
	extractor := NewCommandExtractor()

	tests := []struct {
		name        string
		text        string
		minExpected int  // Minimum number of commands expected
		contains    []string // Strings that should be in the results
	}{
		{
			name:        "bash code block",
			text:        "```bash\nls -la\n```",
			minExpected: 1,
			contains:    []string{"ls -la"},
		},
		{
			name:        "shell code block",
			text:        "```shell\ngit status\n```",
			minExpected: 1,
			contains:    []string{"git status"},
		},
		{
			name:        "plain code block",
			text:        "```\nnpm install\n```",
			minExpected: 1,
			contains:    []string{"npm install"},
		},
		{
			name:        "inline code",
			text:        "Run `echo hello` to print",
			minExpected: 1,
			contains:    []string{"echo hello"},
		},
		{
			name:        "dollar sign prefix",
			text:        "$ git commit -m 'test'",
			minExpected: 1,
			contains:    []string{"git commit -m 'test'"},
		},
		{
			name:        "no commands",
			text:        "This is just plain text without commands.",
			minExpected: 0,
			contains:    nil,
		},
		{
			name:        "empty string",
			text:        "",
			minExpected: 0,
			contains:    nil,
		},
		{
			name:        "duplicate commands",
			text:        "`ls` and then `ls` again",
			minExpected: 1,
			contains:    []string{"ls"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractor.Extract(tt.text)
			assert.GreaterOrEqual(t, len(result), tt.minExpected)
			for _, expected := range tt.contains {
				found := false
				for _, cmd := range result {
					if cmd == expected {
						found = true
						break
					}
				}
				assert.True(t, found, "Expected to find %q in results %v", expected, result)
			}
		})
	}
}

func TestTask_Verification(t *testing.T) {
	t.Run("exit code verification", func(t *testing.T) {
		v := &Verification{
			Type:         VerificationExitCode,
			ExpectedCode: 0,
		}

		assert.Equal(t, VerificationExitCode, v.Type)
		assert.Equal(t, 0, v.ExpectedCode)
	})

	t.Run("output verification", func(t *testing.T) {
		v := &Verification{
			Type:     VerificationOutput,
			Contains: "success",
		}

		assert.Equal(t, VerificationOutput, v.Type)
		assert.Equal(t, "success", v.Contains)
	})

	t.Run("command verification", func(t *testing.T) {
		v := &Verification{
			Type:    VerificationCommand,
			Command: "test -f output.txt",
		}

		assert.Equal(t, VerificationCommand, v.Type)
		assert.Equal(t, "test -f output.txt", v.Command)
	})
}

func TestPlanContext(t *testing.T) {
	ctx := &PlanContext{
		CurrentDir:     "/home/user/project",
		OS:             "linux",
		Shell:          "bash",
		ProjectType:    "go",
		ProjectName:    "myapp",
		AvailableTools: []string{"go", "git", "make"},
	}

	assert.Equal(t, "/home/user/project", ctx.CurrentDir)
	assert.Equal(t, "linux", ctx.OS)
	assert.Equal(t, "bash", ctx.Shell)
	assert.Equal(t, "go", ctx.ProjectType)
	assert.Equal(t, "myapp", ctx.ProjectName)
	assert.Equal(t, []string{"go", "git", "make"}, ctx.AvailableTools)
}

func TestGeneratePlanID(t *testing.T) {
	id1 := generatePlanID()
	id2 := generatePlanID()

	assert.True(t, len(id1) > 0)
	assert.True(t, len(id2) > 0)
	assert.True(t, len(id1) > 5) // plan- prefix + short ID
}

func TestGenerateShortID(t *testing.T) {
	id := generateShortID()

	assert.Equal(t, 8, len(id))
	// All chars should be alphanumeric lowercase
	for _, c := range id {
		valid := (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
		assert.True(t, valid, "Invalid character in ID: %c", c)
	}
}
