package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ProviderOllama is the Ollama provider constant
const ProviderOllama Provider = "ollama"

// OllamaClient implements the Client interface for Ollama
type OllamaClient struct {
	baseURL    string
	model      string
	httpClient *http.Client
}

// OllamaConfig holds configuration for Ollama
type OllamaConfig struct {
	BaseURL string // Default: http://localhost:11434
	Model   string // Default: codellama
	Timeout time.Duration
}

// DefaultOllamaConfig returns default Ollama configuration
func DefaultOllamaConfig() *OllamaConfig {
	return &OllamaConfig{
		BaseURL: "http://localhost:11434",
		Model:   "codellama",
		Timeout: 120 * time.Second,
	}
}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient(config *OllamaConfig) (*OllamaClient, error) {
	if config == nil {
		config = DefaultOllamaConfig()
	}

	if config.BaseURL == "" {
		config.BaseURL = "http://localhost:11434"
	}

	if config.Model == "" {
		config.Model = "codellama"
	}

	if config.Timeout == 0 {
		config.Timeout = 120 * time.Second
	}

	client := &OllamaClient{
		baseURL: strings.TrimSuffix(config.BaseURL, "/"),
		model:   config.Model,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}

	return client, nil
}

// ollamaRequest represents a request to Ollama API
type ollamaRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	System  string                 `json:"system,omitempty"`
	Stream  bool                   `json:"stream"`
	Options map[string]interface{} `json:"options,omitempty"`
}

// ollamaResponse represents a response from Ollama API
type ollamaResponse struct {
	Model     string `json:"model"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
	Context   []int  `json:"context,omitempty"`
	Error     string `json:"error,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

// ollamaModelsResponse represents the list models response
type ollamaModelsResponse struct {
	Models []ollamaModel `json:"models"`
}

// ollamaModel represents an Ollama model
type ollamaModel struct {
	Name       string `json:"name"`
	ModifiedAt string `json:"modified_at"`
	Size       int64  `json:"size"`
}

// Complete implements the Client interface
func (c *OllamaClient) Complete(ctx context.Context, req *Request) (*Response, error) {
	systemPrompt := c.buildSystemPrompt(req)
	userPrompt := c.buildUserPrompt(req)

	ollamaReq := ollamaRequest{
		Model:  c.model,
		Prompt: userPrompt,
		System: systemPrompt,
		Stream: false,
		Options: map[string]interface{}{
			"temperature": 0.7,
			"num_predict": 1024,
		},
	}

	body, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var ollamaResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if ollamaResp.Error != "" {
		return nil, fmt.Errorf("ollama error: %s", ollamaResp.Error)
	}

	return c.parseResponse(ollamaResp.Response, req.Type)
}

// Provider returns the provider name
func (c *OllamaClient) Provider() Provider {
	return ProviderOllama
}

// buildSystemPrompt creates the system prompt based on request type
func (c *OllamaClient) buildSystemPrompt(req *Request) string {
	var sb strings.Builder

	sb.WriteString("You are an expert command-line assistant. ")

	switch req.Type {
	case RequestTypeCommand:
		sb.WriteString(`Your task is to convert natural language descriptions into shell commands.

Rules:
1. Output ONLY the command, no explanations
2. Use proper flags and options
3. Consider the operating system and shell
4. Prefer safe, non-destructive commands
5. If multiple commands are needed, chain them with && or provide them on separate lines

Format your response as:
COMMAND: <the shell command>
EXPLANATION: <brief explanation of what it does>`)

	case RequestTypeExplain:
		sb.WriteString(`Your task is to explain commands, errors, or technical concepts.

Rules:
1. Be clear and concise
2. Use simple language
3. Provide examples when helpful
4. Suggest fixes for errors`)

	case RequestTypeChat:
		sb.WriteString(`Your task is to help with technical questions.

Rules:
1. Be helpful and informative
2. Provide code examples when appropriate
3. Consider the user's context`)
	}

	return sb.String()
}

// buildUserPrompt creates the user prompt with context
func (c *OllamaClient) buildUserPrompt(req *Request) string {
	var sb strings.Builder

	// Add context
	if req.Context != nil {
		sb.WriteString(fmt.Sprintf("Operating System: %s\n", req.Context.OS))
		sb.WriteString(fmt.Sprintf("Shell: %s\n", req.Context.Shell))
		sb.WriteString(fmt.Sprintf("Current Directory: %s\n", req.Context.CurrentDir))

		if req.Context.ProjectName != "" {
			sb.WriteString(fmt.Sprintf("Project: %s (%s)\n", req.Context.ProjectName, req.Context.ProjectType))
		}

		if len(req.Context.RecentHistory) > 0 {
			sb.WriteString("Recent commands:\n")
			for _, cmd := range req.Context.RecentHistory {
				sb.WriteString(fmt.Sprintf("  $ %s\n", cmd))
			}
		}

		sb.WriteString("\n")
	}

	// Add user message
	sb.WriteString(req.UserMessage)

	return sb.String()
}

// parseResponse parses the Ollama response
func (c *OllamaClient) parseResponse(content string, reqType RequestType) (*Response, error) {
	resp := &Response{
		Content:    content,
		Provider:   ProviderOllama,
		Confidence: 0.7, // Local models typically have lower confidence
	}

	if reqType == RequestTypeCommand {
		// Try to extract command from response
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)

			// Look for COMMAND: prefix
			if strings.HasPrefix(strings.ToUpper(line), "COMMAND:") {
				resp.Command = strings.TrimSpace(strings.TrimPrefix(line, "COMMAND:"))
				resp.Command = strings.TrimPrefix(resp.Command, "command:")
				resp.Command = strings.TrimSpace(resp.Command)
				continue
			}

			// Look for code blocks
			if strings.HasPrefix(line, "```") {
				continue
			}

			// Look for EXPLANATION: prefix
			if strings.HasPrefix(strings.ToUpper(line), "EXPLANATION:") {
				explanation := strings.TrimSpace(strings.TrimPrefix(line, "EXPLANATION:"))
				explanation = strings.TrimPrefix(explanation, "explanation:")
				resp.Explanation = strings.TrimSpace(explanation)
				continue
			}
		}

		// If no command was found, try to extract it from code blocks
		if resp.Command == "" {
			resp.Command = extractCommandFromContent(content)
		}
	}

	return resp, nil
}

// extractCommandFromContent extracts a command from content
func extractCommandFromContent(content string) string {
	// Try to find command in code blocks
	lines := strings.Split(content, "\n")
	inCodeBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}

		if inCodeBlock && trimmed != "" {
			return trimmed
		}
	}

	// Try to find lines starting with $
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "$ ") {
			return strings.TrimPrefix(trimmed, "$ ")
		}
	}

	return ""
}

// IsAvailable checks if Ollama is running and accessible
func (c *OllamaClient) IsAvailable(ctx context.Context) bool {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/tags", nil)
	if err != nil {
		return false
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// ListModels returns available Ollama models
func (c *OllamaClient) ListModels(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/tags", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list models: status %d", resp.StatusCode)
	}

	var modelsResp ollamaModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	models := make([]string, len(modelsResp.Models))
	for i, m := range modelsResp.Models {
		models[i] = m.Name
	}

	return models, nil
}

// SetModel changes the model being used
func (c *OllamaClient) SetModel(model string) {
	c.model = model
}

// GetModel returns the current model
func (c *OllamaClient) GetModel() string {
	return c.model
}

// PullModel pulls a model from Ollama library
func (c *OllamaClient) PullModel(ctx context.Context, model string) error {
	reqBody := map[string]string{"name": model}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/pull", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to pull model: %s", string(bodyBytes))
	}

	return nil
}

// StreamingResponse represents a streaming response chunk
type StreamingResponse struct {
	Content string
	Done    bool
	Error   error
}

// CompleteStream streams the response
func (c *OllamaClient) CompleteStream(ctx context.Context, req *Request, handler func(StreamingResponse)) error {
	systemPrompt := c.buildSystemPrompt(req)
	userPrompt := c.buildUserPrompt(req)

	ollamaReq := ollamaRequest{
		Model:  c.model,
		Prompt: userPrompt,
		System: systemPrompt,
		Stream: true,
	}

	body, err := json.Marshal(ollamaReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ollama API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	decoder := json.NewDecoder(resp.Body)
	for {
		var chunk ollamaResponse
		if err := decoder.Decode(&chunk); err != nil {
			if err == io.EOF {
				break
			}
			handler(StreamingResponse{Error: err})
			return err
		}

		handler(StreamingResponse{
			Content: chunk.Response,
			Done:    chunk.Done,
		})

		if chunk.Done {
			break
		}
	}

	return nil
}
