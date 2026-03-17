package workflows

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terminalizcrazy/terminalizcrazy/internal/executor"
)

func TestBuiltInWorkflows(t *testing.T) {
	workflows := BuiltInWorkflows()

	assert.Len(t, workflows, 6)

	names := make(map[string]bool)
	for _, wf := range workflows {
		names[wf.Name] = true
		assert.NotEmpty(t, wf.Name)
		assert.NotEmpty(t, wf.Description)
		assert.NotEmpty(t, wf.Steps)
	}

	assert.True(t, names["git-feature-branch"])
	assert.True(t, names["docker-build"])
	assert.True(t, names["go-test-build"])
	assert.True(t, names["node-test-build"])
	assert.True(t, names["git-cleanup"])
	assert.True(t, names["database-backup"])
}

func TestGitFeatureBranch(t *testing.T) {
	wf := GitFeatureBranch()

	assert.Equal(t, "git-feature-branch", wf.Name)
	assert.Contains(t, wf.Description, "feature branch")
	assert.Contains(t, wf.Tags, "git")

	// Check variables
	assert.Len(t, wf.Variables, 2)
	assert.Equal(t, "branch_name", wf.Variables[0].Name)
	assert.True(t, wf.Variables[0].Required)
	assert.Equal(t, "base_branch", wf.Variables[1].Name)
	assert.Equal(t, "main", wf.Variables[1].Default)

	// Check steps
	assert.Len(t, wf.Steps, 5)
	assert.Equal(t, "fetch-latest", wf.Steps[0].Name)
	assert.Equal(t, "checkout-base", wf.Steps[1].Name)
	assert.Equal(t, "pull-latest", wf.Steps[2].Name)
	assert.Equal(t, "create-branch", wf.Steps[3].Name)
	assert.Equal(t, "push-branch", wf.Steps[4].Name)
}

func TestDockerBuild(t *testing.T) {
	wf := DockerBuild()

	assert.Equal(t, "docker-build", wf.Name)
	assert.Contains(t, wf.Description, "Docker")
	assert.Contains(t, wf.Tags, "docker")

	// Check variables
	assert.Len(t, wf.Variables, 4)

	varNames := make(map[string]string)
	for _, v := range wf.Variables {
		varNames[v.Name] = v.Default
	}

	assert.True(t, wf.Variables[0].Required) // image_name
	assert.Equal(t, "latest", varNames["tag"])
	assert.Equal(t, "Dockerfile", varNames["dockerfile"])
	assert.Equal(t, "no", varNames["push"])

	// Check steps
	assert.Len(t, wf.Steps, 4)
	assert.Equal(t, "lint-dockerfile", wf.Steps[0].Name)
	assert.Equal(t, "build-image", wf.Steps[1].Name)
	assert.Equal(t, "test-image", wf.Steps[2].Name)
	assert.Equal(t, "push-image", wf.Steps[3].Name)

	// Check push condition
	assert.Equal(t, "${push}==yes", wf.Steps[3].Condition)
}

func TestGoTestAndBuild(t *testing.T) {
	wf := GoTestAndBuild()

	assert.Equal(t, "go-test-build", wf.Name)
	assert.Contains(t, wf.Description, "Go")
	assert.Contains(t, wf.Tags, "go")
	assert.Contains(t, wf.Tags, "ci")

	// Check variables
	assert.Len(t, wf.Variables, 2)
	assert.Equal(t, "output_name", wf.Variables[0].Name)
	assert.Equal(t, "app", wf.Variables[0].Default)

	// Check steps
	assert.Len(t, wf.Steps, 5)
	assert.Equal(t, "tidy-modules", wf.Steps[0].Name)
	assert.Equal(t, "verify-modules", wf.Steps[1].Name)
	assert.Equal(t, "run-vet", wf.Steps[2].Name)
	assert.Equal(t, "run-tests", wf.Steps[3].Name)
	assert.Equal(t, "build-binary", wf.Steps[4].Name)
}

func TestNodeTestAndBuild(t *testing.T) {
	wf := NodeTestAndBuild()

	assert.Equal(t, "node-test-build", wf.Name)
	assert.Contains(t, wf.Description, "Node")
	assert.Contains(t, wf.Tags, "node")
	assert.Contains(t, wf.Tags, "javascript")

	// Check variables
	assert.Len(t, wf.Variables, 1)
	assert.Equal(t, "package_manager", wf.Variables[0].Name)
	assert.Equal(t, "npm", wf.Variables[0].Default)

	// Check steps
	assert.Len(t, wf.Steps, 5)
	assert.Equal(t, "install-deps", wf.Steps[0].Name)
	assert.Equal(t, "run-lint", wf.Steps[1].Name)
	assert.Equal(t, "run-typecheck", wf.Steps[2].Name)
	assert.Equal(t, "run-tests", wf.Steps[3].Name)
	assert.Equal(t, "build", wf.Steps[4].Name)
}

func TestGitCleanup(t *testing.T) {
	wf := GitCleanup()

	assert.Equal(t, "git-cleanup", wf.Name)
	assert.Contains(t, wf.Description, "Clean")
	assert.Contains(t, wf.Tags, "git")
	assert.Contains(t, wf.Tags, "maintenance")

	// Check variables
	assert.Len(t, wf.Variables, 1)
	assert.Equal(t, "protected_branches", wf.Variables[0].Name)
	assert.Equal(t, "main,master,develop", wf.Variables[0].Default)

	// Check steps
	assert.Len(t, wf.Steps, 4)
	assert.Equal(t, "fetch-prune", wf.Steps[0].Name)
	assert.Equal(t, "list-merged", wf.Steps[1].Name)
	assert.Equal(t, "merged_branches", wf.Steps[1].CaptureAs)
}

func TestDatabaseBackup(t *testing.T) {
	wf := DatabaseBackup()

	assert.Equal(t, "database-backup", wf.Name)
	assert.Contains(t, wf.Description, "PostgreSQL")
	assert.Contains(t, wf.Tags, "database")
	assert.Contains(t, wf.Tags, "backup")

	// Check variables
	assert.Len(t, wf.Variables, 3)

	varMap := make(map[string]Variable)
	for _, v := range wf.Variables {
		varMap[v.Name] = v
	}

	assert.True(t, varMap["database_url"].Required)
	assert.Equal(t, "./backups", varMap["backup_dir"].Default)
	assert.Equal(t, "7", varMap["retention_days"].Default)

	// Check steps
	assert.Len(t, wf.Steps, 6)
	assert.Equal(t, "create-backup-dir", wf.Steps[0].Name)
	assert.Equal(t, "generate-timestamp", wf.Steps[1].Name)
	assert.Equal(t, "timestamp", wf.Steps[1].CaptureAs)
}

func TestRegisterBuiltInWorkflows(t *testing.T) {
	exec := executor.New()
	engine := NewWorkflowEngine(exec)

	err := RegisterBuiltInWorkflows(engine)

	require.NoError(t, err)
	workflows := engine.ListWorkflows()
	assert.Len(t, workflows, 6)
}

func TestBuiltInWorkflowsHaveValidStructure(t *testing.T) {
	for _, wf := range BuiltInWorkflows() {
		t.Run(wf.Name, func(t *testing.T) {
			assert.NotEmpty(t, wf.Name)
			assert.NotEmpty(t, wf.Description)
			assert.NotEmpty(t, wf.Steps)
			assert.NotEmpty(t, wf.Tags)

			for i, step := range wf.Steps {
				assert.NotEmpty(t, step.Name, "step %d has no name", i)
				assert.NotEmpty(t, step.Command, "step %d has no command", i)
			}
		})
	}
}

func TestBuiltInWorkflowsCanBeRegistered(t *testing.T) {
	exec := executor.New()
	engine := NewWorkflowEngine(exec)

	for _, wf := range BuiltInWorkflows() {
		t.Run(wf.Name, func(t *testing.T) {
			err := engine.RegisterWorkflow(wf)
			assert.NoError(t, err)

			retrieved := engine.GetWorkflow(wf.Name)
			assert.NotNil(t, retrieved)
			assert.Equal(t, wf.Name, retrieved.Name)
		})
	}
}
