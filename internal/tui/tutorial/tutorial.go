package tutorial

import (
	tea "github.com/charmbracelet/bubbletea"
)

// TutorialState represents the current state of the tutorial
type TutorialState int

const (
	StateWelcome TutorialState = iota
	StateFirstQuestion
	StateExecuteCommand
	StateRiskConfirmation
	StateHistoryNavigation
	StateClipboard
	StateSessionSharing
	StateSessionRestore
	StateComplete
)

// Tutorial manages the interactive tutorial experience
type Tutorial struct {
	state     TutorialState
	steps     []Step
	current   int
	completed bool
	skipped   bool
}

// Step represents a single tutorial step
type Step struct {
	ID          string
	Title       string
	Description string
	Instruction string
	KeyHint     string
	Validator   func(msg tea.Msg) bool
}

// New creates a new Tutorial instance
func New() *Tutorial {
	return &Tutorial{
		state:     StateWelcome,
		steps:     GetSteps(),
		current:   0,
		completed: false,
		skipped:   false,
	}
}

// CurrentStep returns the current tutorial step
func (t *Tutorial) CurrentStep() *Step {
	if t.current >= 0 && t.current < len(t.steps) {
		return &t.steps[t.current]
	}
	return nil
}

// Advance moves to the next tutorial step
func (t *Tutorial) Advance() bool {
	if t.current < len(t.steps)-1 {
		t.current++
		t.state = TutorialState(t.current)
		return true
	}
	t.completed = true
	return false
}

// Skip skips the tutorial entirely
func (t *Tutorial) Skip() {
	t.skipped = true
	t.completed = true
}

// IsCompleted returns true if the tutorial is finished
func (t *Tutorial) IsCompleted() bool {
	return t.completed
}

// IsSkipped returns true if the tutorial was skipped
func (t *Tutorial) IsSkipped() bool {
	return t.skipped
}

// Reset restarts the tutorial from the beginning
func (t *Tutorial) Reset() {
	t.state = StateWelcome
	t.current = 0
	t.completed = false
	t.skipped = false
}

// GetProgress returns current step and total steps
func (t *Tutorial) GetProgress() (current, total int) {
	return t.current + 1, len(t.steps)
}

// ValidateInput checks if the user action matches the expected tutorial step
func (t *Tutorial) ValidateInput(msg tea.Msg) bool {
	step := t.CurrentStep()
	if step == nil || step.Validator == nil {
		return false
	}
	return step.Validator(msg)
}

// GetState returns the current tutorial state
func (t *Tutorial) GetState() TutorialState {
	return t.state
}

// GetWelcomeMessage returns the welcome message for the tutorial
func (t *Tutorial) GetWelcomeMessage() string {
	return WelcomeMessage
}

// GetCompletionMessage returns the completion message
func (t *Tutorial) GetCompletionMessage() string {
	return CompletionMessage
}

// ShouldShowTutorial returns true if tutorial should be displayed
func ShouldShowTutorial(tutorialCompleted bool) bool {
	return !tutorialCompleted
}
