package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

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

// TaskStatus represents the status of a task
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusSkipped   TaskStatus = "skipped"
)

// PlanTask represents a task in the plan view
type PlanTask struct {
	ID          string
	Sequence    int
	Description string
	Command     string
	Status      TaskStatus
	Output      string
	Error       string
}

// Plan represents an agent plan
type Plan struct {
	ID          string
	Goal        string
	Tasks       []PlanTask
	Status      PlanStatus
	CurrentTask int
}

// PlanViewStyles holds styles for the plan view
type PlanViewStyles struct {
	Title         lipgloss.Style
	Goal          lipgloss.Style
	TaskPending   lipgloss.Style
	TaskRunning   lipgloss.Style
	TaskCompleted lipgloss.Style
	TaskFailed    lipgloss.Style
	TaskSkipped   lipgloss.Style
	Command       lipgloss.Style
	Output        lipgloss.Style
	Error         lipgloss.Style
	Progress      lipgloss.Style
	Status        lipgloss.Style
	Help          lipgloss.Style
}

// PlanView handles plan display
type PlanView struct {
	viewport viewport.Model
	plan     *Plan
	styles   PlanViewStyles
	width    int
	height   int
	ready    bool
	expanded map[string]bool // Track expanded tasks
}

// NewPlanView creates a new plan view
func NewPlanView(width, height int, styles PlanViewStyles) *PlanView {
	vp := viewport.New(width, height)
	vp.HighPerformanceRendering = false

	return &PlanView{
		viewport: vp,
		styles:   styles,
		width:    width,
		height:   height,
		ready:    true,
		expanded: make(map[string]bool),
	}
}

// SetSize updates the view dimensions
func (p *PlanView) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.viewport.Width = width
	p.viewport.Height = height
	p.updateContent()
}

// SetPlan sets the current plan
func (p *PlanView) SetPlan(plan *Plan) {
	p.plan = plan
	p.expanded = make(map[string]bool)
	p.updateContent()
}

// Clear clears the current plan
func (p *PlanView) Clear() {
	p.plan = nil
	p.expanded = make(map[string]bool)
	p.updateContent()
}

// GetPlan returns the current plan
func (p *PlanView) GetPlan() *Plan {
	return p.plan
}

// ToggleTaskExpanded toggles task expansion
func (p *PlanView) ToggleTaskExpanded(taskID string) {
	p.expanded[taskID] = !p.expanded[taskID]
	p.updateContent()
}

// ExpandAll expands all tasks
func (p *PlanView) ExpandAll() {
	if p.plan == nil {
		return
	}
	for _, task := range p.plan.Tasks {
		p.expanded[task.ID] = true
	}
	p.updateContent()
}

// CollapseAll collapses all tasks
func (p *PlanView) CollapseAll() {
	p.expanded = make(map[string]bool)
	p.updateContent()
}

// updateContent updates the viewport content
func (p *PlanView) updateContent() {
	if !p.ready {
		return
	}

	if p.plan == nil {
		p.viewport.SetContent("No active plan")
		return
	}

	var content strings.Builder

	// Plan header
	content.WriteString(p.styles.Title.Render("📋 Agent Plan"))
	content.WriteString("\n\n")

	// Goal
	content.WriteString(p.styles.Goal.Render("Goal: " + p.plan.Goal))
	content.WriteString("\n\n")

	// Status and progress
	statusIcon := getStatusIcon(p.plan.Status)
	statusColor := getStatusColor(p.plan.Status)
	statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(statusColor)).Bold(true)
	content.WriteString(fmt.Sprintf("Status: %s %s", statusIcon, statusStyle.Render(string(p.plan.Status))))
	content.WriteString("\n")

	// Progress bar
	progress := p.getProgress()
	progressBar := p.renderProgressBar(progress, 30)
	content.WriteString(fmt.Sprintf("Progress: %s %.0f%%", progressBar, progress))
	content.WriteString("\n\n")

	// Tasks
	content.WriteString(p.styles.Title.Render("Tasks:"))
	content.WriteString("\n\n")

	for _, task := range p.plan.Tasks {
		taskIcon := getTaskIcon(task.Status)
		taskStyle := p.getTaskStyle(task.Status)

		// Task header
		content.WriteString(fmt.Sprintf("%s %d. %s\n",
			taskIcon,
			task.Sequence,
			taskStyle.Render(task.Description),
		))

		// Command
		content.WriteString(fmt.Sprintf("   %s\n", p.styles.Command.Render("$ "+task.Command)))

		// Expanded content
		if p.expanded[task.ID] {
			if task.Output != "" {
				content.WriteString(p.styles.Output.Render("   Output:\n"))
				lines := strings.Split(task.Output, "\n")
				for _, line := range lines {
					content.WriteString("   " + p.styles.Output.Render(line) + "\n")
				}
			}
			if task.Error != "" {
				content.WriteString(p.styles.Error.Render("   Error: "+task.Error) + "\n")
			}
		}

		content.WriteString("\n")
	}

	// Help
	if p.plan.Status == PlanStatusPending {
		content.WriteString("\n")
		content.WriteString(p.styles.Help.Render("Press Y to approve, N to reject"))
	} else if p.plan.Status == PlanStatusRunning {
		content.WriteString("\n")
		content.WriteString(p.styles.Help.Render("Press C to cancel"))
	}

	p.viewport.SetContent(content.String())
}

// getTaskStyle returns the style for a task status
func (p *PlanView) getTaskStyle(status TaskStatus) lipgloss.Style {
	switch status {
	case TaskStatusPending:
		return p.styles.TaskPending
	case TaskStatusRunning:
		return p.styles.TaskRunning
	case TaskStatusCompleted:
		return p.styles.TaskCompleted
	case TaskStatusFailed:
		return p.styles.TaskFailed
	case TaskStatusSkipped:
		return p.styles.TaskSkipped
	default:
		return p.styles.TaskPending
	}
}

// getProgress calculates completion percentage
func (p *PlanView) getProgress() float64 {
	if p.plan == nil || len(p.plan.Tasks) == 0 {
		return 0
	}

	completed := 0
	for _, task := range p.plan.Tasks {
		if task.Status == TaskStatusCompleted || task.Status == TaskStatusSkipped {
			completed++
		}
	}

	return float64(completed) / float64(len(p.plan.Tasks)) * 100
}

// renderProgressBar renders a text progress bar
func (p *PlanView) renderProgressBar(percentage float64, width int) string {
	filled := int(percentage / 100 * float64(width))
	empty := width - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	return p.styles.Progress.Render("[" + bar + "]")
}

// getStatusIcon returns an icon for plan status
func getStatusIcon(status PlanStatus) string {
	switch status {
	case PlanStatusPending:
		return "⏸"
	case PlanStatusApproved:
		return "✓"
	case PlanStatusRunning:
		return "▶"
	case PlanStatusCompleted:
		return "✅"
	case PlanStatusFailed:
		return "❌"
	case PlanStatusCancelled:
		return "⊘"
	default:
		return "?"
	}
}

// getStatusColor returns a color for plan status
func getStatusColor(status PlanStatus) string {
	switch status {
	case PlanStatusPending:
		return "#FFAA00"
	case PlanStatusApproved:
		return "#04B575"
	case PlanStatusRunning:
		return "#7D56F4"
	case PlanStatusCompleted:
		return "#04B575"
	case PlanStatusFailed:
		return "#FF6B6B"
	case PlanStatusCancelled:
		return "#888888"
	default:
		return "#888888"
	}
}

// getTaskIcon returns an icon for task status
func getTaskIcon(status TaskStatus) string {
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

// ScrollUp scrolls the viewport up
func (p *PlanView) ScrollUp(lines int) {
	p.viewport.LineUp(lines)
}

// ScrollDown scrolls the viewport down
func (p *PlanView) ScrollDown(lines int) {
	p.viewport.LineDown(lines)
}

// View renders the plan view
func (p *PlanView) View() string {
	return p.viewport.View()
}

// UpdateTaskStatus updates a task's status in the view
func (p *PlanView) UpdateTaskStatus(taskID string, status TaskStatus, output, errMsg string) {
	if p.plan == nil {
		return
	}

	for i := range p.plan.Tasks {
		if p.plan.Tasks[i].ID == taskID {
			p.plan.Tasks[i].Status = status
			p.plan.Tasks[i].Output = output
			p.plan.Tasks[i].Error = errMsg
			break
		}
	}

	p.updateContent()
}

// UpdatePlanStatus updates the plan status
func (p *PlanView) UpdatePlanStatus(status PlanStatus) {
	if p.plan == nil {
		return
	}

	p.plan.Status = status
	p.updateContent()
}

// GetCurrentTask returns the currently running task
func (p *PlanView) GetCurrentTask() *PlanTask {
	if p.plan == nil {
		return nil
	}

	for i := range p.plan.Tasks {
		if p.plan.Tasks[i].Status == TaskStatusRunning {
			return &p.plan.Tasks[i]
		}
	}

	return nil
}

// GetPendingTasks returns all pending tasks
func (p *PlanView) GetPendingTasks() []PlanTask {
	if p.plan == nil {
		return nil
	}

	var pending []PlanTask
	for _, task := range p.plan.Tasks {
		if task.Status == TaskStatusPending {
			pending = append(pending, task)
		}
	}

	return pending
}

// GetFailedTasks returns all failed tasks
func (p *PlanView) GetFailedTasks() []PlanTask {
	if p.plan == nil {
		return nil
	}

	var failed []PlanTask
	for _, task := range p.plan.Tasks {
		if task.Status == TaskStatusFailed {
			failed = append(failed, task)
		}
	}

	return failed
}

// IsComplete returns true if the plan is complete
func (p *PlanView) IsComplete() bool {
	return p.plan != nil && (p.plan.Status == PlanStatusCompleted ||
		p.plan.Status == PlanStatusFailed ||
		p.plan.Status == PlanStatusCancelled)
}

// IsPending returns true if the plan is pending approval
func (p *PlanView) IsPending() bool {
	return p.plan != nil && p.plan.Status == PlanStatusPending
}

// IsRunning returns true if the plan is running
func (p *PlanView) IsRunning() bool {
	return p.plan != nil && p.plan.Status == PlanStatusRunning
}
