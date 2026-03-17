package tutorial

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// GetSteps returns all tutorial steps
func GetSteps() []Step {
	return []Step{
		{
			ID:          "welcome",
			Title:       "Willkommen",
			Description: WelcomeDescription,
			Instruction: "Druecken Sie Enter um fortzufahren",
			KeyHint:     "Enter",
			Validator: func(msg tea.Msg) bool {
				if keyMsg, ok := msg.(tea.KeyMsg); ok {
					return keyMsg.Type == tea.KeyEnter
				}
				return false
			},
		},
		{
			ID:          "first_question",
			Title:       "Erste Frage",
			Description: FirstQuestionDescription,
			Instruction: "Tippen Sie: list files",
			KeyHint:     "Eingabe + Enter",
			Validator: func(msg tea.Msg) bool {
				if keyMsg, ok := msg.(tea.KeyMsg); ok {
					return keyMsg.Type == tea.KeyEnter
				}
				return false
			},
		},
		{
			ID:          "execute_command",
			Title:       "Befehl ausfuehren",
			Description: ExecuteCommandDescription,
			Instruction: "Druecken Sie Ctrl+E",
			KeyHint:     "Ctrl+E",
			Validator: func(msg tea.Msg) bool {
				if keyMsg, ok := msg.(tea.KeyMsg); ok {
					return keyMsg.Type == tea.KeyCtrlE
				}
				return false
			},
		},
		{
			ID:          "risk_confirmation",
			Title:       "Risiko verstehen",
			Description: RiskConfirmationDescription,
			Instruction: "Druecken Sie Y oder N",
			KeyHint:     "Y/N",
			Validator: func(msg tea.Msg) bool {
				if keyMsg, ok := msg.(tea.KeyMsg); ok {
					key := strings.ToLower(keyMsg.String())
					return key == "y" || key == "n"
				}
				return false
			},
		},
		{
			ID:          "history_navigation",
			Title:       "History nutzen",
			Description: HistoryNavigationDescription,
			Instruction: "Druecken Sie Pfeil hoch oder runter",
			KeyHint:     "Pfeiltasten",
			Validator: func(msg tea.Msg) bool {
				if keyMsg, ok := msg.(tea.KeyMsg); ok {
					return keyMsg.Type == tea.KeyUp || keyMsg.Type == tea.KeyDown
				}
				return false
			},
		},
		{
			ID:          "clipboard",
			Title:       "Clipboard verwenden",
			Description: ClipboardDescription,
			Instruction: "Druecken Sie Ctrl+Y",
			KeyHint:     "Ctrl+Y",
			Validator: func(msg tea.Msg) bool {
				if keyMsg, ok := msg.(tea.KeyMsg); ok {
					return keyMsg.Type == tea.KeyCtrlY
				}
				return false
			},
		},
		{
			ID:          "session_sharing",
			Title:       "Session teilen",
			Description: SessionSharingDescription,
			Instruction: "Druecken Sie Ctrl+S (optional: Esc zum Ueberspringen)",
			KeyHint:     "Ctrl+S / Esc",
			Validator: func(msg tea.Msg) bool {
				if keyMsg, ok := msg.(tea.KeyMsg); ok {
					return keyMsg.Type == tea.KeyCtrlS || keyMsg.Type == tea.KeyEsc
				}
				return false
			},
		},
		{
			ID:          "session_restore",
			Title:       "Session laden",
			Description: SessionRestoreDescription,
			Instruction: "Navigieren Sie mit Pfeiltasten und Enter",
			KeyHint:     "Pfeile + Enter",
			Validator: func(msg tea.Msg) bool {
				if keyMsg, ok := msg.(tea.KeyMsg); ok {
					return keyMsg.Type == tea.KeyEnter || keyMsg.Type == tea.KeyEsc
				}
				return false
			},
		},
		{
			ID:          "complete",
			Title:       "Fertig\!",
			Description: CompleteDescription,
			Instruction: "Druecken Sie Enter um zu beginnen",
			KeyHint:     "Enter",
			Validator: func(msg tea.Msg) bool {
				if keyMsg, ok := msg.(tea.KeyMsg); ok {
					return keyMsg.Type == tea.KeyEnter
				}
				return false
			},
		},
	}
}

// GetStepByID returns a step by its ID
func GetStepByID(id string) *Step {
	for _, step := range GetSteps() {
		if step.ID == id {
			return &step
		}
	}
	return nil
}

// GetTotalSteps returns the total number of tutorial steps
func GetTotalSteps() int {
	return len(GetSteps())
}
