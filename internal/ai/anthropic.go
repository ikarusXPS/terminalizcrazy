package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/liushuangls/go-anthropic/v2"
)

// AnthropicClient implements the Client interface for Anthropic
type AnthropicClient struct {
	client *anthropic.Client
}

// NewAnthropicClient creates a new Anthropic client
func NewAnthropicClient(apiKey string) (*AnthropicClient, error) {
	client := anthropic.NewClient(apiKey)
	return &AnthropicClient{client: client}, nil
}

// Provider returns the provider name
func (c *AnthropicClient) Provider() Provider {
	return ProviderAnthropic
}

// Complete sends a request to Anthropic and returns a response
func (c *AnthropicClient) Complete(ctx context.Context, req *Request) (*Response, error) {
	systemPrompt := buildSystemPrompt(req)
	userMessage := buildUserMessage(req)

	resp, err := c.client.CreateMessages(ctx, anthropic.MessagesRequest{
		Model:       anthropic.ModelClaude3Dot5Sonnet20241022,
		MultiSystem: anthropic.NewMultiSystemMessages(systemPrompt),
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage(userMessage),
		},
		MaxTokens: 1024,
	})

	if err != nil {
		return nil, fmt.Errorf("anthropic API error: %w", err)
	}

	// Extract text from response
	var content string
	for _, block := range resp.Content {
		if block.Type == "text" {
			content += block.GetText()
		}
	}

	return parseResponse(content, req.Type, ProviderAnthropic), nil
}

// buildSystemPrompt creates the system prompt based on request type
func buildSystemPrompt(req *Request) string {
	basePrompt := `You are TerminalizCrazy, an AI assistant integrated into a terminal application.
You help users with shell commands, explain errors, and answer technical questions.

Current environment:
- OS: %s
- Shell: %s
- Working Directory: %s
`

	ctx := req.Context
	if ctx == nil {
		ctx = &RequestContext{OS: "unknown", Shell: "bash", CurrentDir: "."}
	}

	prompt := fmt.Sprintf(basePrompt, ctx.OS, ctx.Shell, ctx.CurrentDir)

	// Add project context if available
	if ctx.ProjectName != "" || ctx.ProjectType != "" {
		prompt += "\nProject context:\n"
		if ctx.ProjectName != "" {
			prompt += fmt.Sprintf("- Project: %s\n", ctx.ProjectName)
		}
		if ctx.ProjectType != "" {
			prompt += fmt.Sprintf("- Type: %s\n", ctx.ProjectType)
		}
		if ctx.ProjectFramework != "" {
			prompt += fmt.Sprintf("- Framework: %s\n", ctx.ProjectFramework)
		}
		prompt += "\nUse project-appropriate commands and best practices.\n"
	}

	prompt += `
Guidelines:
- Be concise and direct
- For commands: provide the exact command first, then a brief explanation
- For errors: explain what went wrong and how to fix it
- Use the appropriate shell syntax for the user's OS
- If unsure, say so rather than guessing
`

	switch req.Type {
	case RequestTypeCommand:
		prompt += `
For command requests:
- Start your response with the command in a code block
- Use the format: ` + "```" + `shell
<command here>
` + "```" + `
- Then provide a brief explanation
- If multiple steps are needed, show them in order
`
	case RequestTypeExplain:
		prompt += `
For explanations:
- Be clear and educational
- Explain the root cause
- Provide a solution or workaround
- Use examples if helpful
`
	}

	return prompt
}

// buildUserMessage formats the user's message
func buildUserMessage(req *Request) string {
	return req.UserMessage
}

// parseResponse extracts structured data from the AI response
func parseResponse(content string, reqType RequestType, provider Provider) *Response {
	resp := &Response{
		Content:    content,
		Provider:   provider,
		Confidence: 0.8, // Default confidence
	}

	// Extract command from code block if present
	if reqType == RequestTypeCommand {
		resp.Command = extractCommand(content)
		if resp.Command != "" {
			resp.Explanation = extractExplanation(content)
		}
	}

	return resp
}

// extractCommand extracts a command from markdown code blocks
func extractCommand(content string) string {
	// Look for ```shell, ```bash, ```powershell, or just ```
	markers := []string{"```shell", "```bash", "```powershell", "```sh", "```"}

	for _, marker := range markers {
		if idx := strings.Index(content, marker); idx != -1 {
			start := idx + len(marker)
			// Skip newline after marker
			if start < len(content) && content[start] == '\n' {
				start++
			}

			// Find closing ```
			end := strings.Index(content[start:], "```")
			if end != -1 {
				cmd := strings.TrimSpace(content[start : start+end])
				return cmd
			}
		}
	}

	// Fallback: look for lines that look like commands
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and markdown
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "-") {
			continue
		}
		// If it looks like a command (starts with common command prefixes)
		cmdPrefixes := []string{"$", ">", "git ", "npm ", "go ", "docker ", "cd ", "ls ", "dir ", "cat ", "echo "}
		for _, prefix := range cmdPrefixes {
			if strings.HasPrefix(line, prefix) {
				return strings.TrimPrefix(line, "$")
			}
		}
	}

	return ""
}

// extractExplanation extracts the explanation after the command
func extractExplanation(content string) string {
	// Find the end of the code block
	if idx := strings.LastIndex(content, "```"); idx != -1 {
		explanation := strings.TrimSpace(content[idx+3:])
		if explanation != "" {
			return explanation
		}
	}

	return ""
}

// CompleteStream sends a streaming request to Anthropic
// Note: Falls back to non-streaming if stream API is unavailable
func (c *AnthropicClient) CompleteStream(ctx context.Context, req *Request, handler func(StreamingResponse)) error {
	// Use the standard Complete method and simulate streaming
	// This is a fallback approach; real streaming would require library support
	resp, err := c.Complete(ctx, req)
	if err != nil {
		return err
	}

	// Send the complete response as a single chunk
	handler(StreamingResponse{
		Delta:    resp.Content,
		Done:     false,
		FullText: resp.Content,
	})

	handler(StreamingResponse{
		Done:     true,
		Command:  resp.Command,
		FullText: resp.Content,
	})

	return nil
}
