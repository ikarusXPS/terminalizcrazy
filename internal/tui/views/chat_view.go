package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

// ChatMessage represents a chat message
type ChatMessage struct {
	Role      string
	Content   string
	Command   string
	Success   bool
	Timestamp time.Time
}

// ChatViewStyles holds styles for the chat view
type ChatViewStyles struct {
	UserMsg   lipgloss.Style
	AIMsg     lipgloss.Style
	SystemMsg lipgloss.Style
	Command   lipgloss.Style
	Output    lipgloss.Style
	Error     lipgloss.Style
	Success   lipgloss.Style
	Timestamp lipgloss.Style
	Help      lipgloss.Style
}

// ChatView handles chat display
type ChatView struct {
	viewport viewport.Model
	messages []ChatMessage
	styles   ChatViewStyles
	width    int
	height   int
	ready    bool
}

// NewChatView creates a new chat view
func NewChatView(width, height int, styles ChatViewStyles) *ChatView {
	vp := viewport.New(width, height)
	vp.HighPerformanceRendering = false

	return &ChatView{
		viewport: vp,
		messages: []ChatMessage{},
		styles:   styles,
		width:    width,
		height:   height,
		ready:    true,
	}
}

// SetSize updates the view dimensions
func (c *ChatView) SetSize(width, height int) {
	c.width = width
	c.height = height
	c.viewport.Width = width
	c.viewport.Height = height
	c.updateContent()
}

// AddMessage adds a message to the chat
func (c *ChatView) AddMessage(msg ChatMessage) {
	c.messages = append(c.messages, msg)
	c.updateContent()
}

// AddUserMessage adds a user message
func (c *ChatView) AddUserMessage(content string) {
	c.AddMessage(ChatMessage{
		Role:      "user",
		Content:   content,
		Timestamp: time.Now(),
	})
}

// AddAIMessage adds an AI message
func (c *ChatView) AddAIMessage(content, command string) {
	c.AddMessage(ChatMessage{
		Role:      "ai",
		Content:   content,
		Command:   command,
		Timestamp: time.Now(),
	})
}

// AddSystemMessage adds a system message
func (c *ChatView) AddSystemMessage(content string) {
	c.AddMessage(ChatMessage{
		Role:      "system",
		Content:   content,
		Timestamp: time.Now(),
	})
}

// AddOutputMessage adds an output message
func (c *ChatView) AddOutputMessage(content string, success bool) {
	c.AddMessage(ChatMessage{
		Role:      "output",
		Content:   content,
		Success:   success,
		Timestamp: time.Now(),
	})
}

// Clear clears all messages
func (c *ChatView) Clear() {
	c.messages = []ChatMessage{}
	c.updateContent()
}

// GetMessages returns all messages
func (c *ChatView) GetMessages() []ChatMessage {
	return c.messages
}

// GetLastCommand returns the last command from AI messages
func (c *ChatView) GetLastCommand() string {
	for i := len(c.messages) - 1; i >= 0; i-- {
		if c.messages[i].Role == "ai" && c.messages[i].Command != "" {
			return c.messages[i].Command
		}
	}
	return ""
}

// updateContent updates the viewport content
func (c *ChatView) updateContent() {
	if !c.ready {
		return
	}

	var content strings.Builder

	for _, msg := range c.messages {
		timestamp := msg.Timestamp.Format("15:04")

		switch msg.Role {
		case "user":
			content.WriteString(fmt.Sprintf("%s %s\n",
				c.styles.UserMsg.Render("You:"),
				c.styles.Timestamp.Render(timestamp),
			))
			content.WriteString(msg.Content + "\n\n")

		case "ai":
			content.WriteString(fmt.Sprintf("%s %s\n",
				c.styles.AIMsg.Render("AI:"),
				c.styles.Timestamp.Render(timestamp),
			))

			if msg.Command != "" {
				content.WriteString(c.styles.Command.Render(msg.Command) + "\n")
				content.WriteString(c.styles.Help.Render("  Press Ctrl+E to execute") + "\n\n")

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
				c.styles.SystemMsg.Render("System:"),
				c.styles.Timestamp.Render(timestamp),
			))
			content.WriteString(c.styles.SystemMsg.Render(msg.Content) + "\n\n")

		case "output":
			if msg.Success {
				content.WriteString(c.styles.Success.Render("✓ Command executed") + " " + c.styles.Timestamp.Render(timestamp) + "\n")
			} else {
				content.WriteString(c.styles.Error.Render("✗ Command failed") + " " + c.styles.Timestamp.Render(timestamp) + "\n")
			}
			content.WriteString(c.styles.Output.Render(msg.Content) + "\n\n")

		case "collab":
			content.WriteString(fmt.Sprintf("%s %s\n",
				c.styles.AIMsg.Render("Collaborator:"),
				c.styles.Timestamp.Render(timestamp),
			))
			content.WriteString(msg.Content + "\n\n")
		}
	}

	c.viewport.SetContent(content.String())
	c.viewport.GotoBottom()
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

// ScrollUp scrolls the viewport up
func (c *ChatView) ScrollUp(lines int) {
	c.viewport.LineUp(lines)
}

// ScrollDown scrolls the viewport down
func (c *ChatView) ScrollDown(lines int) {
	c.viewport.LineDown(lines)
}

// PageUp scrolls up one page
func (c *ChatView) PageUp() {
	c.viewport.ViewUp()
}

// PageDown scrolls down one page
func (c *ChatView) PageDown() {
	c.viewport.ViewDown()
}

// GotoTop scrolls to the top
func (c *ChatView) GotoTop() {
	c.viewport.GotoTop()
}

// GotoBottom scrolls to the bottom
func (c *ChatView) GotoBottom() {
	c.viewport.GotoBottom()
}

// View renders the chat view
func (c *ChatView) View() string {
	return c.viewport.View()
}

// LoadMessages loads messages from storage
func (c *ChatView) LoadMessages(messages []ChatMessage) {
	c.messages = messages
	c.updateContent()
}

// AppendMessages appends messages without clearing
func (c *ChatView) AppendMessages(messages []ChatMessage) {
	c.messages = append(c.messages, messages...)
	c.updateContent()
}

// GetMessageCount returns the number of messages
func (c *ChatView) GetMessageCount() int {
	return len(c.messages)
}

// SearchMessages searches for messages containing text
func (c *ChatView) SearchMessages(query string) []ChatMessage {
	var results []ChatMessage
	query = strings.ToLower(query)

	for _, msg := range c.messages {
		if strings.Contains(strings.ToLower(msg.Content), query) ||
			strings.Contains(strings.ToLower(msg.Command), query) {
			results = append(results, msg)
		}
	}

	return results
}

// ExportMessages exports messages as plain text
func (c *ChatView) ExportMessages() string {
	var sb strings.Builder

	for _, msg := range c.messages {
		timestamp := msg.Timestamp.Format("2006-01-02 15:04:05")
		sb.WriteString(fmt.Sprintf("[%s] %s:\n", timestamp, msg.Role))
		sb.WriteString(msg.Content)
		if msg.Command != "" {
			sb.WriteString(fmt.Sprintf("\nCommand: %s", msg.Command))
		}
		sb.WriteString("\n\n")
	}

	return sb.String()
}
