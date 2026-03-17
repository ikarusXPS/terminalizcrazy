package workflows

import "time"

// BuiltInWorkflows returns all built-in workflow templates
func BuiltInWorkflows() []*Workflow {
	return []*Workflow{
		GitFeatureBranch(),
		DockerBuild(),
		GoTestAndBuild(),
		NodeTestAndBuild(),
		GitCleanup(),
		DatabaseBackup(),
	}
}

// GitFeatureBranch returns a workflow for creating a git feature branch
func GitFeatureBranch() *Workflow {
	return &Workflow{
		Name:        "git-feature-branch",
		Description: "Create a new git feature branch with proper setup",
		Tags:        []string{"git", "development"},
		Variables: []Variable{
			{
				Name:        "branch_name",
				Description: "Name of the feature branch",
				Required:    true,
			},
			{
				Name:        "base_branch",
				Description: "Base branch to create from",
				Default:     "main",
			},
		},
		Steps: []WorkflowStep{
			{
				Name:        "fetch-latest",
				Description: "Fetch latest changes from remote",
				Command:     "git fetch origin",
				OnFail:      OnFailStop,
			},
			{
				Name:        "checkout-base",
				Description: "Checkout base branch",
				Command:     "git checkout ${base_branch}",
				OnFail:      OnFailStop,
			},
			{
				Name:        "pull-latest",
				Description: "Pull latest changes",
				Command:     "git pull origin ${base_branch}",
				OnFail:      OnFailStop,
			},
			{
				Name:        "create-branch",
				Description: "Create and checkout new branch",
				Command:     "git checkout -b feature/${branch_name}",
				OnFail:      OnFailStop,
			},
			{
				Name:        "push-branch",
				Description: "Push branch to remote",
				Command:     "git push -u origin feature/${branch_name}",
				OnFail:      OnFailSkip,
			},
		},
	}
}

// DockerBuild returns a workflow for building and pushing a Docker image
func DockerBuild() *Workflow {
	return &Workflow{
		Name:        "docker-build",
		Description: "Build and optionally push a Docker image",
		Tags:        []string{"docker", "deployment"},
		Variables: []Variable{
			{
				Name:        "image_name",
				Description: "Docker image name",
				Required:    true,
			},
			{
				Name:        "tag",
				Description: "Image tag",
				Default:     "latest",
			},
			{
				Name:        "dockerfile",
				Description: "Path to Dockerfile",
				Default:     "Dockerfile",
			},
			{
				Name:        "push",
				Description: "Push to registry (yes/no)",
				Default:     "no",
			},
		},
		Steps: []WorkflowStep{
			{
				Name:        "lint-dockerfile",
				Description: "Lint Dockerfile",
				Command:     "docker run --rm -i hadolint/hadolint < ${dockerfile}",
				OnFail:      OnFailSkip,
				Timeout:     30 * time.Second,
			},
			{
				Name:        "build-image",
				Description: "Build Docker image",
				Command:     "docker build -t ${image_name}:${tag} -f ${dockerfile} .",
				OnFail:      OnFailStop,
				Timeout:     10 * time.Minute,
			},
			{
				Name:        "test-image",
				Description: "Run basic image tests",
				Command:     "docker run --rm ${image_name}:${tag} --version || echo 'No version command'",
				OnFail:      OnFailSkip,
				Timeout:     30 * time.Second,
			},
			{
				Name:        "push-image",
				Description: "Push image to registry",
				Command:     "docker push ${image_name}:${tag}",
				Condition:   "${push}==yes",
				OnFail:      OnFailStop,
				Timeout:     5 * time.Minute,
			},
		},
	}
}

// GoTestAndBuild returns a workflow for Go testing and building
func GoTestAndBuild() *Workflow {
	return &Workflow{
		Name:        "go-test-build",
		Description: "Run Go tests and build binary",
		Tags:        []string{"go", "development", "ci"},
		Variables: []Variable{
			{
				Name:        "output_name",
				Description: "Output binary name",
				Default:     "app",
			},
			{
				Name:        "build_flags",
				Description: "Additional build flags",
				Default:     "",
			},
		},
		Steps: []WorkflowStep{
			{
				Name:        "tidy-modules",
				Description: "Tidy Go modules",
				Command:     "go mod tidy",
				OnFail:      OnFailStop,
			},
			{
				Name:        "verify-modules",
				Description: "Verify module checksums",
				Command:     "go mod verify",
				OnFail:      OnFailSkip,
			},
			{
				Name:        "run-vet",
				Description: "Run go vet",
				Command:     "go vet ./...",
				OnFail:      OnFailSkip,
			},
			{
				Name:        "run-tests",
				Description: "Run tests with coverage",
				Command:     "go test -v -cover ./...",
				OnFail:      OnFailStop,
				Timeout:     5 * time.Minute,
			},
			{
				Name:        "build-binary",
				Description: "Build binary",
				Command:     "go build ${build_flags} -o ${output_name} ./...",
				OnFail:      OnFailStop,
				Timeout:     2 * time.Minute,
			},
		},
	}
}

// NodeTestAndBuild returns a workflow for Node.js testing and building
func NodeTestAndBuild() *Workflow {
	return &Workflow{
		Name:        "node-test-build",
		Description: "Run Node.js tests and build",
		Tags:        []string{"node", "javascript", "development", "ci"},
		Variables: []Variable{
			{
				Name:        "package_manager",
				Description: "Package manager (npm/yarn/pnpm)",
				Default:     "npm",
			},
		},
		Steps: []WorkflowStep{
			{
				Name:        "install-deps",
				Description: "Install dependencies",
				Command:     "${package_manager} install",
				OnFail:      OnFailStop,
				Timeout:     5 * time.Minute,
			},
			{
				Name:        "run-lint",
				Description: "Run linter",
				Command:     "${package_manager} run lint",
				OnFail:      OnFailSkip,
			},
			{
				Name:        "run-typecheck",
				Description: "Run type checking",
				Command:     "${package_manager} run typecheck",
				OnFail:      OnFailSkip,
			},
			{
				Name:        "run-tests",
				Description: "Run tests",
				Command:     "${package_manager} test",
				OnFail:      OnFailStop,
				Timeout:     5 * time.Minute,
			},
			{
				Name:        "build",
				Description: "Build project",
				Command:     "${package_manager} run build",
				OnFail:      OnFailStop,
				Timeout:     5 * time.Minute,
			},
		},
	}
}

// GitCleanup returns a workflow for cleaning up git branches
func GitCleanup() *Workflow {
	return &Workflow{
		Name:        "git-cleanup",
		Description: "Clean up merged and stale git branches",
		Tags:        []string{"git", "maintenance"},
		Variables: []Variable{
			{
				Name:        "protected_branches",
				Description: "Branches to protect (comma-separated)",
				Default:     "main,master,develop",
			},
		},
		Steps: []WorkflowStep{
			{
				Name:        "fetch-prune",
				Description: "Fetch and prune remote tracking branches",
				Command:     "git fetch --prune",
				OnFail:      OnFailStop,
			},
			{
				Name:        "list-merged",
				Description: "List merged branches",
				Command:     "git branch --merged main | grep -v 'main\\|master\\|develop\\|\\*'",
				OnFail:      OnFailSkip,
				CaptureAs:   "merged_branches",
			},
			{
				Name:        "delete-merged-local",
				Description: "Delete merged local branches",
				Command:     "git branch --merged main | grep -v 'main\\|master\\|develop\\|\\*' | xargs -r git branch -d",
				OnFail:      OnFailSkip,
			},
			{
				Name:        "show-status",
				Description: "Show current branch status",
				Command:     "git branch -vv",
				OnFail:      OnFailSkip,
			},
		},
	}
}

// DatabaseBackup returns a workflow for database backup
func DatabaseBackup() *Workflow {
	return &Workflow{
		Name:        "database-backup",
		Description: "Backup PostgreSQL database",
		Tags:        []string{"database", "backup", "postgres"},
		Variables: []Variable{
			{
				Name:        "database_url",
				Description: "Database connection URL",
				Required:    true,
			},
			{
				Name:        "backup_dir",
				Description: "Backup directory",
				Default:     "./backups",
			},
			{
				Name:        "retention_days",
				Description: "Days to keep backups",
				Default:     "7",
			},
		},
		Steps: []WorkflowStep{
			{
				Name:        "create-backup-dir",
				Description: "Create backup directory",
				Command:     "mkdir -p ${backup_dir}",
				OnFail:      OnFailStop,
			},
			{
				Name:        "generate-timestamp",
				Description: "Generate backup filename",
				Command:     "date +%Y%m%d_%H%M%S",
				CaptureAs:   "timestamp",
				OnFail:      OnFailStop,
			},
			{
				Name:        "backup-database",
				Description: "Create database backup",
				Command:     "pg_dump ${database_url} > ${backup_dir}/backup_${timestamp}.sql",
				OnFail:      OnFailStop,
				Timeout:     30 * time.Minute,
			},
			{
				Name:        "compress-backup",
				Description: "Compress backup file",
				Command:     "gzip ${backup_dir}/backup_${timestamp}.sql",
				OnFail:      OnFailSkip,
			},
			{
				Name:        "cleanup-old",
				Description: "Remove old backups",
				Command:     "find ${backup_dir} -name '*.gz' -mtime +${retention_days} -delete",
				OnFail:      OnFailSkip,
			},
			{
				Name:        "list-backups",
				Description: "List current backups",
				Command:     "ls -la ${backup_dir}",
				OnFail:      OnFailSkip,
			},
		},
	}
}

// RegisterBuiltInWorkflows registers all built-in workflows with an engine
func RegisterBuiltInWorkflows(engine *WorkflowEngine) error {
	for _, wf := range BuiltInWorkflows() {
		if err := engine.RegisterWorkflow(wf); err != nil {
			return err
		}
	}
	return nil
}
