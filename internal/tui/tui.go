package tui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/terminalizcrazy/terminalizcrazy/internal/ai"
	"github.com/terminalizcrazy/terminalizcrazy/internal/clipboard"
	"github.com/terminalizcrazy/terminalizcrazy/internal/collab"
	"github.com/terminalizcrazy/terminalizcrazy/internal/config"
	"github.com/terminalizcrazy/terminalizcrazy/internal/executor"
	"github.com/terminalizcrazy/terminalizcrazy/internal/project"
	"github.com/terminalizcrazy/terminalizcrazy/internal/secretguard"
	"github.com/terminalizcrazy/terminalizcrazy/internal/storage"
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4"))

	versionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	inputStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	statusConnectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#04B575"))

	statusDisconnectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFAA00"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	userMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	aiMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	systemMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Italic(true)

	commandStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#2D2D2D")).
			Foreground(lipgloss.Color("#E0E0E0")).
			Padding(0, 1)

	outputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	spinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4"))

	historyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666"))

	copyNoticeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	sessionItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC")).
			PaddingLeft(2)

	sessionSelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true).
			PaddingLeft(2)

	sessionHeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true).
			MarginBottom(1)

	collabUserStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4"))

	shareCodeStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#2D2D2D")).
			Foreground(lipgloss.Color("#FFEAA7")).
			Bold(true).
			Padding(0, 1)
)

// ViewMode represents the current view
type ViewMode int

const (
	ViewChat ViewMode = iota
	ViewSessionSelect
	ViewCollabJoin
)

// Message types for async operations
type aiResponseMsg struct {
	response *ai.Response
	err      error
}

type cmdResultMsg struct {
	result *executor.Result
}

type historyLoadedMsg struct {
	commands []string
}

type sessionsLoadedMsg struct {
	sessions []storage.Session
}

type sessionRestoredMsg struct {
	messages []storage.Message
	session  *storage.Session
}

type collabConnectedMsg struct {
	shareCode string
}

type collabJoinedMsg struct {
	roomID string
}

type collabMessageMsg struct {
	message *collab.Message
}

type collabErrorMsg struct {
	err error
}

type collabServerStartedMsg struct {
	port int
}

// ConfirmState represents the confirmation dialog state
type ConfirmState struct {
	Active    bool
	Command   string
	RiskLevel executor.RiskLevel
}

// Model represents the application state
type Model struct {
	config      *config.Config
	version     string
	input       textinput.Model
	viewport    viewport.Model
	spinner     spinner.Model
	aiService   *ai.Service
	secretGuard *secretguard.Guard
	executor    *executor.Executor
	storage     *storage.Storage
	clipboard   *clipboard.Manager

	sessionID       string
	messages        []ChatMessage
	commandHistory  []string
	historyIndex    int
	lastCommand     string
	lastCopied      string
	loading         bool
	executing       bool
	confirmState    ConfirmState
	width           int
	height          int
	ready           bool
	quitting        bool
	viewportReady   bool
	browsingHistory bool
	showCopyNotice  bool
	copyNoticeTime  time.Time

	// Session selection
	viewMode          ViewMode
	availableSessions []storage.Session
	sessionIndex      int

	// Smart sessions
	project          *project.Project
	suggestedSession *storage.Session

	// Collaboration
	collabClient   *collab.CollabClient
	collabServer   *collab.Server
	collabUsers    []*collab.User
	shareCode      string
	joinCodeInput  string
	isSharing      bool
	isCollaborating bool
}

// ChatMessage represents a chat message
type ChatMessage struct {
	Role      string
	Content   string
	Command   string
	Success   bool
	Timestamp time.Time
}

// NewModel creates a new TUI model
func NewModel(cfg *config.Config, version string) Model {
	// Text input
	ti := textinput.New()
	ti.Placeholder = "Ask anything... (e.g., 'how to find large files')"
	ti.Focus()
	ti.CharLimit = 500

	// Spinner for loading state
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = spinnerStyle

	// Initialize AI service
	var aiService *ai.Service
	if cfg.HasAIKey() {
		provider := ai.ProviderAnthropic
		apiKey := cfg.AnthropicKey
		if cfg.AIProvider == "openai" {
			provider = ai.ProviderOpenAI
			apiKey = cfg.OpenAIKey
		}
		aiService, _ = ai.NewService(provider, apiKey)
	}

	// Initialize SecretGuard
	sg := secretguard.New(cfg.SecretGuardEnabled)

	// Initialize Executor
	exec := executor.New()

	// Initialize Storage
	home, _ := os.UserHomeDir()
	dataDir := filepath.Join(home, ".terminalizcrazy")
	store, err := storage.New(dataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialize storage: %v\n", err)
	}

	// Detect project
	workDir, _ := os.Getwd()
	detector := project.NewDetector(workDir)
	proj := detector.Detect()

	// Generate session ID
	sessionID := uuid.New().String()[:8]

	// Generate smart session name
	sessionName := proj.GenerateSessionName()

	// Create session in storage
	if store != nil {
		store.CreateSession(sessionID, sessionName, workDir)
	}

	// Initialize Clipboard
	cb, _ := clipboard.New()

	return Model{
		config:         cfg,
		version:        version,
		input:          ti,
		spinner:        sp,
		aiService:      aiService,
		secretGuard:    sg,
		executor:       exec,
		storage:        store,
		clipboard:      cb,
		sessionID:      sessionID,
		messages:       []ChatMessage{},
		commandHistory: []string{},
		historyIndex:   -1,
		project:        proj,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick, m.loadCommandHistory(), m.loadSessions())
}

// loadSessions loads existing sessions from storage
func (m *Model) loadSessions() tea.Cmd {
	return func() tea.Msg {
		if m.storage == nil {
			return sessionsLoadedMsg{sessions: []storage.Session{}}
		}

		sessions, err := m.storage.ListSessions(10)
		if err != nil {
			return sessionsLoadedMsg{sessions: []storage.Session{}}
		}

		return sessionsLoadedMsg{sessions: sessions}
	}
}

// restoreSession loads messages from a stored session
func (m *Model) restoreSession(session storage.Session) tea.Cmd {
	return func() tea.Msg {
		if m.storage == nil {
			return sessionRestoredMsg{messages: []storage.Message{}, session: &session}
		}

		messages, err := m.storage.GetSessionMessages(session.ID, 100)
		if err != nil {
			return sessionRestoredMsg{messages: []storage.Message{}, session: &session}
		}

		return sessionRestoredMsg{messages: messages, session: &session}
	}
}

// startSharing starts a collaboration server and creates a room
func (m *Model) startSharing() tea.Cmd {
	return func() tea.Msg {
		// Start local server
		port := 8765
		m.collabServer = collab.NewServer(port)

		go func() {
			m.collabServer.Start()
		}()

		// Give server time to start
		time.Sleep(100 * time.Millisecond)

		// Create client and room
		m.collabClient = collab.NewClient(fmt.Sprintf("http://localhost:%d", port))
		m.collabClient.SetMessageHandler(func(msg *collab.Message) {
			// This will be handled via tea.Cmd
		})

		userName := "Host"
		if m.project != nil {
			userName = m.project.Name
		}

		shareCode, err := m.collabClient.CreateRoom(m.sessionID, userName)
		if err != nil {
			return collabErrorMsg{err: err}
		}

		return collabConnectedMsg{shareCode: shareCode}
	}
}

// joinCollab joins an existing collaboration room
func (m *Model) joinCollab(shareCode string) tea.Cmd {
	return func() tea.Msg {
		// For now, assume server is on localhost:8765
		// In production, this would be configurable
		m.collabClient = collab.NewClient("http://localhost:8765")

		userName := "Guest"
		if m.project != nil {
			userName = m.project.Name
		}

		err := m.collabClient.JoinRoom(shareCode, m.sessionID, userName)
		if err != nil {
			return collabErrorMsg{err: err}
		}

		m.shareCode = shareCode
		return collabJoinedMsg{roomID: m.collabClient.GetRoomID()}
	}
}

// handleCollabMessage processes incoming collaboration messages
func (m *Model) handleCollabMessage(msg *collab.Message) {
	switch msg.Type {
	case collab.MsgTypeChat:
		// Add chat message from other user
		if msg.UserID != m.sessionID {
			chatMsg := ChatMessage{
				Role:      "collab",
				Content:   fmt.Sprintf("[%s]: %s", msg.UserName, msg.Content),
				Timestamp: msg.Timestamp,
			}
			m.messages = append(m.messages, chatMsg)
			m.updateViewportContent()
		}

	case collab.MsgTypeCommand:
		// Show command from other user
		if msg.UserID != m.sessionID {
			m.addSystemMessage(fmt.Sprintf("[%s] suggested: %s", msg.UserName, msg.Command))
			m.lastCommand = msg.Command
		}

	case collab.MsgTypeOutput:
		// Show output from other user
		if msg.UserID != m.sessionID {
			m.addSystemMessage(fmt.Sprintf("[%s] output:\n%s", msg.UserName, msg.Content))
		}

	case collab.MsgTypeJoin:
		m.addSystemMessage(fmt.Sprintf("%s joined the session", msg.UserName))
		m.collabUsers = m.collabClient.GetUsers()

	case collab.MsgTypeLeave:
		m.addSystemMessage(fmt.Sprintf("%s left the session", msg.UserName))
		m.collabUsers = m.collabClient.GetUsers()

	case collab.MsgTypeUserList:
		m.collabUsers = m.collabClient.GetUsers()
	}
}

// disconnectCollab disconnects from collaboration
func (m *Model) disconnectCollab() {
	if m.collabClient != nil {
		m.collabClient.Disconnect()
		m.collabClient = nil
	}

	if m.collabServer != nil {
		m.collabServer.Stop()
		m.collabServer = nil
	}

	m.isSharing = false
	m.isCollaborating = false
	m.shareCode = ""
	m.collabUsers = nil

	m.addSystemMessage("Disconnected from collaboration")
}

// broadcastToCollab sends a message to collaboration if connected
func (m *Model) broadcastToCollab(msgType collab.MessageType, content, command string) {
	if m.collabClient != nil && m.collabClient.IsConnected() {
		msg := &collab.Message{
			Type:    msgType,
			Content: content,
			Command: command,
		}
		m.collabClient.SendMessage(msg)
	}
}

// loadCommandHistory loads command history from storage
func (m *Model) loadCommandHistory() tea.Cmd {
	return func() tea.Msg {
		if m.storage == nil {
			return historyLoadedMsg{commands: []string{}}
		}

		commands, err := m.storage.GetUniqueCommands(100)
		if err != nil {
			return historyLoadedMsg{commands: []string{}}
		}

		return historyLoadedMsg{commands: commands}
	}
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle session select mode
		if m.viewMode == ViewSessionSelect {
			switch msg.String() {
			case "ctrl+c":
				m.quitting = true
				if m.storage != nil {
					m.storage.Close()
				}
				return m, tea.Quit

			case "up", "k":
				if m.sessionIndex > 0 {
					m.sessionIndex--
				}
				return m, nil

			case "down", "j":
				// +1 for "New Session" option
				if m.sessionIndex < len(m.availableSessions) {
					m.sessionIndex++
				}
				return m, nil

			case "enter":
				if m.sessionIndex == 0 {
					// New session selected
					m.viewMode = ViewChat
					return m, nil
				}
				// Restore existing session (index - 1 because 0 is "New Session")
				session := m.availableSessions[m.sessionIndex-1]
				return m, m.restoreSession(session)

			case "n":
				// Quick key for new session
				m.viewMode = ViewChat
				return m, nil

			case "esc":
				// Start new session on escape
				m.viewMode = ViewChat
				return m, nil
			}
			return m, nil
		}

		// Handle collab join view
		if m.viewMode == ViewCollabJoin {
			switch msg.String() {
			case "ctrl+c":
				m.quitting = true
				return m, tea.Quit

			case "esc":
				m.viewMode = ViewChat
				m.joinCodeInput = ""
				return m, nil

			case "enter":
				if len(m.joinCodeInput) >= 9 { // XXXX-XXXX
					return m, m.joinCollab(m.joinCodeInput)
				}
				return m, nil

			case "backspace":
				if len(m.joinCodeInput) > 0 {
					m.joinCodeInput = m.joinCodeInput[:len(m.joinCodeInput)-1]
				}
				return m, nil

			default:
				// Accept alphanumeric and dash
				if len(msg.String()) == 1 && len(m.joinCodeInput) < 9 {
					char := msg.String()[0]
					if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
						m.joinCodeInput += msg.String()
					}
				}
				return m, nil
			}
		}

		// Handle confirmation dialog
		if m.confirmState.Active {
			switch msg.String() {
			case "y", "Y":
				m.confirmState.Active = false
				m.executing = true
				cmd := m.confirmState.Command
				m.confirmState.Command = ""
				return m, tea.Batch(m.spinner.Tick, m.executeCommand(cmd))

			case "n", "N", "esc":
				m.confirmState.Active = false
				m.confirmState.Command = ""
				m.addSystemMessage("Command execution cancelled.")
				return m, nil
			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			if m.storage != nil {
				m.storage.Close()
			}
			return m, tea.Quit

		case "esc":
			if m.loading || m.executing {
				m.loading = false
				m.executing = false
			} else if m.browsingHistory {
				m.browsingHistory = false
				m.historyIndex = -1
				m.input.SetValue("")
			} else {
				m.quitting = true
				if m.storage != nil {
					m.storage.Close()
				}
				return m, tea.Quit
			}

		case "enter":
			if m.input.Value() != "" && !m.loading && !m.executing {
				userInput := m.input.Value()
				m.input.SetValue("")
				m.browsingHistory = false
				m.historyIndex = -1

				// Add to messages
				m.addUserMessage(userInput)

				// Check if AI is available
				if m.aiService != nil {
					m.loading = true
					return m, tea.Batch(m.spinner.Tick, m.sendAIRequest(userInput))
				} else {
					m.addAIMessage("AI not configured. Please set ANTHROPIC_API_KEY or OPENAI_API_KEY.", "")
				}
			}

		case "ctrl+e":
			if m.lastCommand != "" && !m.loading && !m.executing {
				return m, m.tryExecuteCommand(m.lastCommand)
			}

		case "ctrl+r":
			if m.lastCommand != "" {
				m.addSystemMessage(fmt.Sprintf("Last command: %s\nPress Ctrl+E to execute", m.lastCommand))
			}

		case "up":
			if !m.loading && !m.executing && len(m.commandHistory) > 0 {
				if !m.browsingHistory {
					m.browsingHistory = true
					m.historyIndex = 0
				} else if m.historyIndex < len(m.commandHistory)-1 {
					m.historyIndex++
				}
				if m.historyIndex < len(m.commandHistory) {
					m.input.SetValue(m.commandHistory[m.historyIndex])
					m.input.CursorEnd()
				}
			} else if m.viewportReady {
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
				cmds = append(cmds, cmd)
			}

		case "down":
			if m.browsingHistory {
				if m.historyIndex > 0 {
					m.historyIndex--
					m.input.SetValue(m.commandHistory[m.historyIndex])
					m.input.CursorEnd()
				} else {
					m.browsingHistory = false
					m.historyIndex = -1
					m.input.SetValue("")
				}
			} else if m.viewportReady {
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
				cmds = append(cmds, cmd)
			}

		case "pgup", "pgdown":
			if m.viewportReady {
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
				cmds = append(cmds, cmd)
			}

		case "ctrl+l":
			// Clear screen
			m.messages = []ChatMessage{}
			m.updateViewportContent()

		case "ctrl+y":
			// Copy last command to clipboard
			if m.lastCommand != "" && m.clipboard != nil {
				if err := m.clipboard.CopyCommand(m.lastCommand); err == nil {
					m.lastCopied = m.lastCommand
					m.showCopyNotice = true
					m.copyNoticeTime = time.Now()
					m.addSystemMessage(fmt.Sprintf("Copied to clipboard: %s", m.lastCommand))
				}
			}

		case "ctrl+s":
			// Start sharing session
			if !m.isSharing && !m.isCollaborating {
				return m, m.startSharing()
			}

		case "ctrl+j":
			// Join collaboration
			if !m.isSharing && !m.isCollaborating {
				m.viewMode = ViewCollabJoin
				m.joinCodeInput = ""
			}

		case "ctrl+d":
			// Disconnect from collaboration
			if m.isSharing || m.isCollaborating {
				m.disconnectCollab()
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		headerHeight := 3
		inputHeight := 3
		helpHeight := 2
		confirmHeight := 0
		if m.confirmState.Active {
			confirmHeight = 6
		}
		vpHeight := m.height - headerHeight - inputHeight - helpHeight - confirmHeight

		if !m.viewportReady {
			m.viewport = viewport.New(m.width-2, vpHeight)
			m.viewport.HighPerformanceRendering = false
			m.viewportReady = true
		} else {
			m.viewport.Width = m.width - 2
			m.viewport.Height = vpHeight
		}

		m.input.Width = m.width - 4
		m.ready = true
		m.updateViewportContent()

	case aiResponseMsg:
		m.loading = false

		if msg.err != nil {
			m.addAIMessage(fmt.Sprintf("Error: %v", msg.err), "")
		} else {
			content := m.secretGuard.Mask(msg.response.Content)
			command := msg.response.Command

			if command != "" {
				m.lastCommand = command
			}

			m.addAIMessage(content, command)
		}

	case cmdResultMsg:
		m.executing = false

		output := m.secretGuard.Mask(msg.result.Output)
		errorOutput := m.secretGuard.Mask(msg.result.Error)

		var content string
		if msg.result.Success {
			content = fmt.Sprintf("$ %s\n%s", msg.result.Command, output)
		} else {
			content = fmt.Sprintf("$ %s\n%s\n%s", msg.result.Command, output, errorOutput)
		}

		m.addOutputMessage(content, msg.result.Success)

		// Save to command history
		if m.storage != nil {
			m.storage.SaveCommand(
				msg.result.Command,
				output+errorOutput,
				msg.result.Success,
				msg.result.Duration.Milliseconds(),
			)
			// Add to local history
			m.commandHistory = append([]string{msg.result.Command}, m.commandHistory...)
			if len(m.commandHistory) > 100 {
				m.commandHistory = m.commandHistory[:100]
			}
		}

	case historyLoadedMsg:
		m.commandHistory = msg.commands

	case sessionsLoadedMsg:
		m.availableSessions = msg.sessions
		// If there are existing sessions, show session select menu
		if len(msg.sessions) > 0 {
			m.viewMode = ViewSessionSelect
			m.sessionIndex = 0

			// Find suggested session based on working directory
			workDir, _ := os.Getwd()
			for i, session := range msg.sessions {
				if session.WorkDir == workDir {
					m.suggestedSession = &msg.sessions[i]
					// Pre-select the suggested session (+1 because 0 is "New Session")
					m.sessionIndex = i + 1
					break
				}
			}
		}

	case sessionRestoredMsg:
		m.viewMode = ViewChat
		m.sessionID = msg.session.ID

		// Convert stored messages to chat messages
		m.messages = []ChatMessage{}
		for _, storedMsg := range msg.messages {
			chatMsg := ChatMessage{
				Role:      storedMsg.Role,
				Content:   storedMsg.Content,
				Command:   storedMsg.Command,
				Success:   storedMsg.Success,
				Timestamp: storedMsg.CreatedAt,
			}
			m.messages = append(m.messages, chatMsg)

			// Restore last command from AI messages
			if storedMsg.Role == "ai" && storedMsg.Command != "" {
				m.lastCommand = storedMsg.Command
			}
		}

		m.updateViewportContent()
		m.addSystemMessage(fmt.Sprintf("Session restored: %s (%d messages)", msg.session.Name, len(msg.messages)))

	case collabConnectedMsg:
		m.isSharing = true
		m.shareCode = msg.shareCode
		m.addSystemMessage(fmt.Sprintf("Sharing started! Share code: %s", msg.shareCode))

	case collabJoinedMsg:
		m.isCollaborating = true
		m.viewMode = ViewChat
		m.addSystemMessage(fmt.Sprintf("Joined collaboration room: %s", msg.roomID))

	case collabMessageMsg:
		m.handleCollabMessage(msg.message)

	case collabErrorMsg:
		m.addSystemMessage(fmt.Sprintf("Collaboration error: %v", msg.err))
		m.isSharing = false
		m.isCollaborating = false

	case collabServerStartedMsg:
		m.addSystemMessage(fmt.Sprintf("Collaboration server started on port %d", msg.port))

	case spinner.TickMsg:
		if m.loading || m.executing {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	// Update text input
	var inputCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)
	cmds = append(cmds, inputCmd)

	return m, tea.Batch(cmds...)
}

// Helper methods to add messages
func (m *Model) addUserMessage(content string) {
	msg := ChatMessage{
		Role:      "user",
		Content:   content,
		Timestamp: time.Now(),
	}
	m.messages = append(m.messages, msg)
	m.updateViewportContent()

	// Persist to storage
	if m.storage != nil {
		m.storage.SaveMessage(m.sessionID, "user", content, "", true)
	}
}

func (m *Model) addAIMessage(content, command string) {
	msg := ChatMessage{
		Role:      "ai",
		Content:   content,
		Command:   command,
		Timestamp: time.Now(),
	}
	m.messages = append(m.messages, msg)
	m.updateViewportContent()

	// Persist to storage
	if m.storage != nil {
		m.storage.SaveMessage(m.sessionID, "ai", content, command, true)
	}
}

func (m *Model) addSystemMessage(content string) {
	msg := ChatMessage{
		Role:      "system",
		Content:   content,
		Timestamp: time.Now(),
	}
	m.messages = append(m.messages, msg)
	m.updateViewportContent()
}

func (m *Model) addOutputMessage(content string, success bool) {
	msg := ChatMessage{
		Role:      "output",
		Content:   content,
		Success:   success,
		Timestamp: time.Now(),
	}
	m.messages = append(m.messages, msg)
	m.updateViewportContent()

	// Persist to storage
	if m.storage != nil {
		m.storage.SaveMessage(m.sessionID, "output", content, "", success)
	}
}

// tryExecuteCommand checks risk and either executes or asks for confirmation
func (m *Model) tryExecuteCommand(command string) tea.Cmd {
	risk := m.executor.AssessRisk(command)

	if risk >= executor.RiskMedium {
		m.confirmState = ConfirmState{
			Active:    true,
			Command:   command,
			RiskLevel: risk,
		}
		return nil
	}

	m.executing = true
	return tea.Batch(m.spinner.Tick, m.executeCommand(command))
}

// executeCommand creates a command to execute a shell command
func (m *Model) executeCommand(command string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		result := m.executor.Execute(ctx, command)
		return cmdResultMsg{result: result}
	}
}

// sendAIRequest creates a command that sends a request to the AI service
func (m *Model) sendAIRequest(input string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Build project context
		var projectCtx *ai.RequestContext
		if m.project != nil {
			workDir, _ := os.Getwd()
			projectCtx = &ai.RequestContext{
				CurrentDir:       workDir,
				ProjectName:      m.project.Name,
				ProjectType:      project.GetTypeLabel(m.project.Type),
				ProjectFramework: m.project.Framework,
			}
		}

		resp, err := m.aiService.ProcessInputWithContext(ctx, input, projectCtx)
		return aiResponseMsg{response: resp, err: err}
	}
}

// updateViewportContent updates the viewport with current messages
func (m *Model) updateViewportContent() {
	if !m.viewportReady {
		return
	}

	var content strings.Builder

	for _, msg := range m.messages {
		timestamp := msg.Timestamp.Format("15:04")

		switch msg.Role {
		case "user":
			content.WriteString(fmt.Sprintf("%s %s\n",
				userMsgStyle.Render("You:"),
				versionStyle.Render(timestamp),
			))
			content.WriteString(msg.Content + "\n\n")

		case "ai":
			content.WriteString(fmt.Sprintf("%s %s\n",
				aiMsgStyle.Render("AI:"),
				versionStyle.Render(timestamp),
			))

			if msg.Command != "" {
				content.WriteString(commandStyle.Render(msg.Command) + "\n")
				content.WriteString(helpStyle.Render("  Press Ctrl+E to execute") + "\n\n")

				explanation := extractExplanation(msg.Content)
				if explanation != "" {
					content.WriteString(explanation + "\n")
				}
			} else {
				content.WriteString(msg.Content + "\n")
			}
			content.WriteString("\n")

		case "system":
			content.WriteString(fmt.Sprintf("%s %s\n",
				systemMsgStyle.Render("System:"),
				versionStyle.Render(timestamp),
			))
			content.WriteString(systemMsgStyle.Render(msg.Content) + "\n\n")

		case "output":
			if msg.Success {
				content.WriteString(successStyle.Render("✓ Command executed") + " " + versionStyle.Render(timestamp) + "\n")
			} else {
				content.WriteString(errorStyle.Render("✗ Command failed") + " " + versionStyle.Render(timestamp) + "\n")
			}
			content.WriteString(outputStyle.Render(msg.Content) + "\n\n")
		}
	}

	m.viewport.SetContent(content.String())
	m.viewport.GotoBottom()
}

// extractExplanation extracts explanation text after command block
func extractExplanation(content string) string {
	if idx := strings.LastIndex(content, "```"); idx != -1 {
		explanation := strings.TrimSpace(content[idx+3:])
		if explanation != "" {
			return explanation
		}
	}
	return ""
}

// View renders the UI
func (m Model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	if !m.ready {
		return "Loading...\n"
	}

	// Session selection view
	if m.viewMode == ViewSessionSelect {
		return m.renderSessionSelect()
	}

	// Collab join view
	if m.viewMode == ViewCollabJoin {
		return m.renderCollabJoin()
	}

	var s strings.Builder

	// Header
	title := titleStyle.Render("⚡ TerminalizCrazy")
	version := versionStyle.Render(fmt.Sprintf("v%s", m.version))
	s.WriteString(fmt.Sprintf("%s %s  ", title, version))

	// AI status
	if m.aiService != nil {
		provider := string(m.aiService.GetProvider())
		s.WriteString(statusConnectedStyle.Render(fmt.Sprintf("● %s", provider)))
	} else {
		s.WriteString(statusDisconnectedStyle.Render("○ No AI"))
	}

	// Session info
	s.WriteString(versionStyle.Render(fmt.Sprintf("  [%s]", m.sessionID)))

	// Collaboration status
	if m.isSharing {
		s.WriteString("  ")
		s.WriteString(collabUserStyle.Render(fmt.Sprintf("📡 Sharing: %s", m.shareCode)))
		if len(m.collabUsers) > 1 {
			s.WriteString(collabUserStyle.Render(fmt.Sprintf(" (%d users)", len(m.collabUsers))))
		}
	} else if m.isCollaborating {
		s.WriteString("  ")
		s.WriteString(collabUserStyle.Render(fmt.Sprintf("🤝 Collaborating (%d users)", len(m.collabUsers))))
	}

	s.WriteString("\n\n")

	// Viewport (chat history)
	if m.viewportReady {
		s.WriteString(m.viewport.View())
		s.WriteString("\n")
	}

	// Confirmation dialog
	if m.confirmState.Active {
		riskColor := executor.GetRiskColor(m.confirmState.RiskLevel)
		riskDesc := executor.GetRiskDescription(m.confirmState.RiskLevel)

		confirmBox := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(riskColor)).
			Padding(0, 1).
			Render(fmt.Sprintf(
				"%s\n\n%s\n\n%s",
				lipgloss.NewStyle().Foreground(lipgloss.Color(riskColor)).Bold(true).Render(riskDesc),
				commandStyle.Render(m.confirmState.Command),
				"Execute? [Y]es / [N]o",
			))

		s.WriteString(confirmBox + "\n")
	}

	// Loading/Executing indicator
	if m.loading {
		s.WriteString(fmt.Sprintf("%s Thinking...\n", m.spinner.View()))
	} else if m.executing {
		s.WriteString(fmt.Sprintf("%s Executing...\n", m.spinner.View()))
	}

	// History indicator
	if m.browsingHistory && len(m.commandHistory) > 0 {
		s.WriteString(historyStyle.Render(fmt.Sprintf("  History [%d/%d]", m.historyIndex+1, len(m.commandHistory))) + "\n")
	}

	// Copy notice (show for 3 seconds)
	if m.showCopyNotice && time.Since(m.copyNoticeTime) < 3*time.Second {
		s.WriteString(copyNoticeStyle.Render("  ✓ Copied to clipboard") + "\n")
	} else {
		m.showCopyNotice = false
	}

	// Input
	s.WriteString(inputStyle.Render(m.input.View()))
	s.WriteString("\n")

	// Help
	var helpParts []string
	helpParts = append(helpParts, "Enter: Send", "↑↓: History")

	if m.lastCommand != "" {
		helpParts = append(helpParts, "Ctrl+E: Execute", "Ctrl+Y: Copy")
	}

	// Collaboration shortcuts
	if m.isSharing || m.isCollaborating {
		helpParts = append(helpParts, "Ctrl+D: Disconnect")
	} else {
		helpParts = append(helpParts, "Ctrl+S: Share", "Ctrl+J: Join")
	}

	helpParts = append(helpParts, "Esc: Quit")

	help := strings.Join(helpParts, " • ")
	s.WriteString(helpStyle.Render(help))

	return s.String()
}

// renderSessionSelect renders the session selection view
func (m Model) renderSessionSelect() string {
	var s strings.Builder

	// Header
	title := titleStyle.Render("⚡ TerminalizCrazy")
	version := versionStyle.Render(fmt.Sprintf("v%s", m.version))
	s.WriteString(fmt.Sprintf("%s %s\n\n", title, version))

	// Project info
	if m.project != nil && m.project.Type != project.TypeUnknown {
		icon := project.GetTypeIcon(m.project.Type)
		label := project.GetTypeLabel(m.project.Type)
		projectInfo := fmt.Sprintf("%s %s", icon, m.project.Name)
		if m.project.Framework != "" {
			projectInfo += fmt.Sprintf(" (%s)", m.project.Framework)
		}
		s.WriteString(aiMsgStyle.Render(fmt.Sprintf("Detected: %s %s", label, projectInfo)))
		s.WriteString("\n\n")
	}

	// Session selection header
	s.WriteString(sessionHeaderStyle.Render("Select Session:"))
	s.WriteString("\n\n")

	// New session option (always first)
	newSessionLabel := "New Session"
	if m.project != nil && m.project.Type != project.TypeUnknown {
		newSessionLabel = fmt.Sprintf("New Session: %s", m.project.GenerateSessionName())
	}
	if m.sessionIndex == 0 {
		s.WriteString(sessionSelectedStyle.Render(fmt.Sprintf("▶ %s", newSessionLabel)))
	} else {
		s.WriteString(sessionItemStyle.Render(fmt.Sprintf("  %s", newSessionLabel)))
	}
	s.WriteString("\n")

	// Existing sessions
	workDir, _ := os.Getwd()
	for i, session := range m.availableSessions {
		// Format timestamp
		timeAgo := formatTimeAgo(session.UpdatedAt)
		sessionLabel := fmt.Sprintf("%s • %s", session.Name, timeAgo)

		// Mark sessions from same directory
		isSameDir := session.WorkDir == workDir
		if isSameDir {
			sessionLabel += " ★"
		}

		if m.sessionIndex == i+1 {
			s.WriteString(sessionSelectedStyle.Render(fmt.Sprintf("▶ %s", sessionLabel)))
		} else {
			s.WriteString(sessionItemStyle.Render(fmt.Sprintf("  %s", sessionLabel)))
		}
		s.WriteString("\n")
	}

	// Suggestion hint
	if m.suggestedSession != nil {
		s.WriteString("\n")
		s.WriteString(helpStyle.Render("★ = same directory (recommended)"))
	}

	s.WriteString("\n")
	s.WriteString(helpStyle.Render("↑↓: Navigate • Enter: Select • N: New Session • Esc: New Session"))

	return s.String()
}

// renderCollabJoin renders the collaboration join view
func (m Model) renderCollabJoin() string {
	var s strings.Builder

	// Header
	title := titleStyle.Render("⚡ TerminalizCrazy")
	version := versionStyle.Render(fmt.Sprintf("v%s", m.version))
	s.WriteString(fmt.Sprintf("%s %s\n\n", title, version))

	// Join header
	s.WriteString(sessionHeaderStyle.Render("🤝 Join Collaboration"))
	s.WriteString("\n\n")

	s.WriteString("Enter share code:\n\n")

	// Code input display
	displayCode := m.joinCodeInput
	if len(displayCode) < 9 {
		displayCode += strings.Repeat("_", 9-len(displayCode))
	}
	// Format as XXXX-XXXX
	if len(displayCode) >= 4 {
		displayCode = displayCode[:4] + "-" + displayCode[4:]
	}

	s.WriteString(shareCodeStyle.Render(displayCode))
	s.WriteString("\n\n")

	if len(m.joinCodeInput) >= 9 {
		s.WriteString(successStyle.Render("Press Enter to join"))
	} else {
		s.WriteString(helpStyle.Render("Type the 8-character code"))
	}

	s.WriteString("\n\n")
	s.WriteString(helpStyle.Render("Esc: Cancel"))

	return s.String()
}

// formatTimeAgo formats a time as a relative string
func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return "just now"
	case duration < time.Hour:
		mins := int(duration.Minutes())
		if mins == 1 {
			return "1 min ago"
		}
		return fmt.Sprintf("%d mins ago", mins)
	case duration < 24*time.Hour:
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case duration < 7*24*time.Hour:
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "yesterday"
		}
		return fmt.Sprintf("%d days ago", days)
	default:
		return t.Format("2006-01-02")
	}
}

// Run starts the TUI application
func Run(cfg *config.Config, version string) error {
	p := tea.NewProgram(
		NewModel(cfg, version),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := p.Run()
	return err
}
