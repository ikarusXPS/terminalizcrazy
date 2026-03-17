package project

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// ProjectType represents the detected project type
type ProjectType string

const (
	TypeUnknown    ProjectType = "unknown"
	TypeGo         ProjectType = "go"
	TypeNode       ProjectType = "node"
	TypePython     ProjectType = "python"
	TypeRust       ProjectType = "rust"
	TypeJava       ProjectType = "java"
	TypeDotNet     ProjectType = "dotnet"
	TypeRuby       ProjectType = "ruby"
	TypePHP        ProjectType = "php"
	TypeDocker     ProjectType = "docker"
	TypeTerraform  ProjectType = "terraform"
	TypeKubernetes ProjectType = "kubernetes"
)

// Project represents detected project information
type Project struct {
	Name        string      `json:"name"`
	Type        ProjectType `json:"type"`
	Path        string      `json:"path"`
	Description string      `json:"description,omitempty"`
	Version     string      `json:"version,omitempty"`
	Framework   string      `json:"framework,omitempty"`
}

// Detector handles project detection
type Detector struct {
	workDir string
}

// NewDetector creates a new project detector
func NewDetector(workDir string) *Detector {
	return &Detector{workDir: workDir}
}

// Detect analyzes the working directory and returns project info
func (d *Detector) Detect() *Project {
	project := &Project{
		Path: d.workDir,
		Type: TypeUnknown,
		Name: filepath.Base(d.workDir),
	}

	// Try to detect project type in order of specificity
	detectors := []func(*Project) bool{
		d.detectGo,
		d.detectNode,
		d.detectPython,
		d.detectRust,
		d.detectJava,
		d.detectDotNet,
		d.detectRuby,
		d.detectPHP,
		d.detectDocker,
		d.detectTerraform,
		d.detectKubernetes,
	}

	for _, detect := range detectors {
		if detect(project) {
			break
		}
	}

	return project
}

// detectGo checks for Go project
func (d *Detector) detectGo(p *Project) bool {
	goModPath := filepath.Join(d.workDir, "go.mod")
	if _, err := os.Stat(goModPath); err != nil {
		return false
	}

	p.Type = TypeGo

	// Parse go.mod for module name
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return true
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			moduleName := strings.TrimPrefix(line, "module ")
			// Extract last part of module path as name
			parts := strings.Split(moduleName, "/")
			p.Name = parts[len(parts)-1]
			p.Description = moduleName
			break
		}
	}

	// Detect framework
	if d.fileContains(goModPath, "github.com/gin-gonic/gin") {
		p.Framework = "Gin"
	} else if d.fileContains(goModPath, "github.com/labstack/echo") {
		p.Framework = "Echo"
	} else if d.fileContains(goModPath, "github.com/gofiber/fiber") {
		p.Framework = "Fiber"
	} else if d.fileContains(goModPath, "github.com/charmbracelet/bubbletea") {
		p.Framework = "Bubble Tea"
	}

	return true
}

// detectNode checks for Node.js project
func (d *Detector) detectNode(p *Project) bool {
	packagePath := filepath.Join(d.workDir, "package.json")
	if _, err := os.Stat(packagePath); err != nil {
		return false
	}

	p.Type = TypeNode

	// Parse package.json
	content, err := os.ReadFile(packagePath)
	if err != nil {
		return true
	}

	var pkg struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		Description string `json:"description"`
	}

	if err := json.Unmarshal(content, &pkg); err == nil {
		if pkg.Name != "" {
			p.Name = pkg.Name
		}
		p.Version = pkg.Version
		p.Description = pkg.Description
	}

	// Detect framework
	if d.fileContains(packagePath, "\"next\"") {
		p.Framework = "Next.js"
	} else if d.fileContains(packagePath, "\"react\"") {
		p.Framework = "React"
	} else if d.fileContains(packagePath, "\"vue\"") {
		p.Framework = "Vue"
	} else if d.fileContains(packagePath, "\"express\"") {
		p.Framework = "Express"
	} else if d.fileContains(packagePath, "\"nestjs\"") || d.fileContains(packagePath, "\"@nestjs/core\"") {
		p.Framework = "NestJS"
	} else if d.fileContains(packagePath, "\"svelte\"") {
		p.Framework = "Svelte"
	} else if d.fileContains(packagePath, "\"angular\"") || d.fileContains(packagePath, "\"@angular/core\"") {
		p.Framework = "Angular"
	}

	return true
}

// detectPython checks for Python project
func (d *Detector) detectPython(p *Project) bool {
	// Check various Python project files
	pyprojectPath := filepath.Join(d.workDir, "pyproject.toml")
	setupPath := filepath.Join(d.workDir, "setup.py")
	requirementsPath := filepath.Join(d.workDir, "requirements.txt")

	hasPyproject := fileExists(pyprojectPath)
	hasSetup := fileExists(setupPath)
	hasRequirements := fileExists(requirementsPath)

	if !hasPyproject && !hasSetup && !hasRequirements {
		return false
	}

	p.Type = TypePython

	// Try to parse pyproject.toml
	if hasPyproject {
		content, err := os.ReadFile(pyprojectPath)
		if err == nil {
			var pyproject map[string]interface{}
			if _, err := toml.Decode(string(content), &pyproject); err == nil {
				if project, ok := pyproject["project"].(map[string]interface{}); ok {
					if name, ok := project["name"].(string); ok {
						p.Name = name
					}
					if version, ok := project["version"].(string); ok {
						p.Version = version
					}
					if desc, ok := project["description"].(string); ok {
						p.Description = desc
					}
				}
			}
		}

		// Detect framework
		if d.fileContains(pyprojectPath, "django") {
			p.Framework = "Django"
		} else if d.fileContains(pyprojectPath, "fastapi") {
			p.Framework = "FastAPI"
		} else if d.fileContains(pyprojectPath, "flask") {
			p.Framework = "Flask"
		}
	}

	// Check requirements.txt for framework
	if hasRequirements && p.Framework == "" {
		if d.fileContains(requirementsPath, "django") {
			p.Framework = "Django"
		} else if d.fileContains(requirementsPath, "fastapi") {
			p.Framework = "FastAPI"
		} else if d.fileContains(requirementsPath, "flask") {
			p.Framework = "Flask"
		}
	}

	return true
}

// detectRust checks for Rust project
func (d *Detector) detectRust(p *Project) bool {
	cargoPath := filepath.Join(d.workDir, "Cargo.toml")
	if _, err := os.Stat(cargoPath); err != nil {
		return false
	}

	p.Type = TypeRust

	// Parse Cargo.toml
	content, err := os.ReadFile(cargoPath)
	if err != nil {
		return true
	}

	var cargo map[string]interface{}
	if _, err := toml.Decode(string(content), &cargo); err == nil {
		if pkg, ok := cargo["package"].(map[string]interface{}); ok {
			if name, ok := pkg["name"].(string); ok {
				p.Name = name
			}
			if version, ok := pkg["version"].(string); ok {
				p.Version = version
			}
			if desc, ok := pkg["description"].(string); ok {
				p.Description = desc
			}
		}
	}

	// Detect framework
	if d.fileContains(cargoPath, "actix-web") {
		p.Framework = "Actix"
	} else if d.fileContains(cargoPath, "axum") {
		p.Framework = "Axum"
	} else if d.fileContains(cargoPath, "rocket") {
		p.Framework = "Rocket"
	} else if d.fileContains(cargoPath, "tauri") {
		p.Framework = "Tauri"
	}

	return true
}

// detectJava checks for Java project
func (d *Detector) detectJava(p *Project) bool {
	pomPath := filepath.Join(d.workDir, "pom.xml")
	gradlePath := filepath.Join(d.workDir, "build.gradle")
	gradleKtsPath := filepath.Join(d.workDir, "build.gradle.kts")

	if !fileExists(pomPath) && !fileExists(gradlePath) && !fileExists(gradleKtsPath) {
		return false
	}

	p.Type = TypeJava

	// Check for Spring
	if fileExists(pomPath) && d.fileContains(pomPath, "spring-boot") {
		p.Framework = "Spring Boot"
	} else if fileExists(gradlePath) && d.fileContains(gradlePath, "spring-boot") {
		p.Framework = "Spring Boot"
	}

	return true
}

// detectDotNet checks for .NET project
func (d *Detector) detectDotNet(p *Project) bool {
	// Look for .csproj or .sln files
	matches, _ := filepath.Glob(filepath.Join(d.workDir, "*.csproj"))
	if len(matches) == 0 {
		matches, _ = filepath.Glob(filepath.Join(d.workDir, "*.sln"))
	}

	if len(matches) == 0 {
		return false
	}

	p.Type = TypeDotNet

	// Extract name from first project file
	if len(matches) > 0 {
		base := filepath.Base(matches[0])
		p.Name = strings.TrimSuffix(base, filepath.Ext(base))
	}

	return true
}

// detectRuby checks for Ruby project
func (d *Detector) detectRuby(p *Project) bool {
	gemfilePath := filepath.Join(d.workDir, "Gemfile")
	if _, err := os.Stat(gemfilePath); err != nil {
		return false
	}

	p.Type = TypeRuby

	// Check for Rails
	if d.fileContains(gemfilePath, "rails") {
		p.Framework = "Rails"
	} else if d.fileContains(gemfilePath, "sinatra") {
		p.Framework = "Sinatra"
	}

	return true
}

// detectPHP checks for PHP project
func (d *Detector) detectPHP(p *Project) bool {
	composerPath := filepath.Join(d.workDir, "composer.json")
	if _, err := os.Stat(composerPath); err != nil {
		return false
	}

	p.Type = TypePHP

	// Parse composer.json
	content, err := os.ReadFile(composerPath)
	if err != nil {
		return true
	}

	var composer struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.Unmarshal(content, &composer); err == nil {
		if composer.Name != "" {
			parts := strings.Split(composer.Name, "/")
			p.Name = parts[len(parts)-1]
		}
		p.Description = composer.Description
	}

	// Detect framework
	if d.fileContains(composerPath, "laravel") {
		p.Framework = "Laravel"
	} else if d.fileContains(composerPath, "symfony") {
		p.Framework = "Symfony"
	}

	return true
}

// detectDocker checks for Docker project
func (d *Detector) detectDocker(p *Project) bool {
	dockerfilePath := filepath.Join(d.workDir, "Dockerfile")
	composePath := filepath.Join(d.workDir, "docker-compose.yml")
	composeAltPath := filepath.Join(d.workDir, "docker-compose.yaml")

	if !fileExists(dockerfilePath) && !fileExists(composePath) && !fileExists(composeAltPath) {
		return false
	}

	p.Type = TypeDocker
	return true
}

// detectTerraform checks for Terraform project
func (d *Detector) detectTerraform(p *Project) bool {
	matches, _ := filepath.Glob(filepath.Join(d.workDir, "*.tf"))
	if len(matches) == 0 {
		return false
	}

	p.Type = TypeTerraform
	return true
}

// detectKubernetes checks for Kubernetes manifests
func (d *Detector) detectKubernetes(p *Project) bool {
	// Check for common k8s directories
	k8sDirs := []string{"k8s", "kubernetes", "manifests", "deploy"}
	for _, dir := range k8sDirs {
		if fileExists(filepath.Join(d.workDir, dir)) {
			p.Type = TypeKubernetes
			return true
		}
	}

	// Check for kustomization.yaml
	if fileExists(filepath.Join(d.workDir, "kustomization.yaml")) {
		p.Type = TypeKubernetes
		return true
	}

	return false
}

// fileContains checks if a file contains a string (case-insensitive)
func (d *Detector) fileContains(path, search string) bool {
	content, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(content)), strings.ToLower(search))
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// GetTypeIcon returns an icon for the project type
func GetTypeIcon(t ProjectType) string {
	icons := map[ProjectType]string{
		TypeGo:         "🐹",
		TypeNode:       "📦",
		TypePython:     "🐍",
		TypeRust:       "🦀",
		TypeJava:       "☕",
		TypeDotNet:     "🔷",
		TypeRuby:       "💎",
		TypePHP:        "🐘",
		TypeDocker:     "🐳",
		TypeTerraform:  "🏗️",
		TypeKubernetes: "☸️",
		TypeUnknown:    "📁",
	}

	if icon, ok := icons[t]; ok {
		return icon
	}
	return "📁"
}

// GetTypeLabel returns a human-readable label for the project type
func GetTypeLabel(t ProjectType) string {
	labels := map[ProjectType]string{
		TypeGo:         "Go",
		TypeNode:       "Node.js",
		TypePython:     "Python",
		TypeRust:       "Rust",
		TypeJava:       "Java",
		TypeDotNet:     ".NET",
		TypeRuby:       "Ruby",
		TypePHP:        "PHP",
		TypeDocker:     "Docker",
		TypeTerraform:  "Terraform",
		TypeKubernetes: "Kubernetes",
		TypeUnknown:    "Project",
	}

	if label, ok := labels[t]; ok {
		return label
	}
	return "Project"
}

// GenerateSessionName creates a smart session name from project info
func (p *Project) GenerateSessionName() string {
	icon := GetTypeIcon(p.Type)

	if p.Framework != "" {
		return fmt.Sprintf("%s %s (%s)", icon, p.Name, p.Framework)
	}

	if p.Type != TypeUnknown {
		return fmt.Sprintf("%s %s", icon, p.Name)
	}

	return p.Name
}

// GetAIContext returns context information for AI prompts
func (p *Project) GetAIContext() string {
	var parts []string

	parts = append(parts, fmt.Sprintf("Project: %s", p.Name))
	parts = append(parts, fmt.Sprintf("Type: %s", GetTypeLabel(p.Type)))

	if p.Framework != "" {
		parts = append(parts, fmt.Sprintf("Framework: %s", p.Framework))
	}

	if p.Description != "" {
		parts = append(parts, fmt.Sprintf("Description: %s", p.Description))
	}

	return strings.Join(parts, "\n")
}
