package ai

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Provider represents an AI provider
type Provider string

const (
	ProviderAnthropic Provider = "anthropic"
	ProviderOpenAI    Provider = "openai"
	// ProviderGemini is defined in gemini.go
	// ProviderOllama is defined in ollama.go
)

// Request represents an AI request
type Request struct {
	// UserMessage is the user's input
	UserMessage string

	// Context provides additional context (current directory, OS, etc.)
	Context *RequestContext

	// Type specifies what kind of response we want
	Type RequestType
}

// RequestType specifies the type of AI request
type RequestType string

const (
	// RequestTypeCommand - Convert natural language to shell command
	RequestTypeCommand RequestType = "command"

	// RequestTypeExplain - Explain an error or command
	RequestTypeExplain RequestType = "explain"

	// RequestTypeChat - General chat/question
	RequestTypeChat RequestType = "chat"
)

// RequestContext provides context for the AI
type RequestContext struct {
	CurrentDir     string
	OS             string
	Shell          string
	RecentHistory  []string
	ProjectName    string
	ProjectType    string
	ProjectFramework string
}

// Response represents an AI response
type Response struct {
	// Content is the main response text
	Content string

	// Command is the suggested command (for RequestTypeCommand)
	Command string

	// Explanation provides additional context
	Explanation string

	// Confidence indicates how confident the AI is (0-1)
	Confidence float64

	// Provider indicates which AI was used
	Provider Provider
}

// Client is the interface for AI providers
type Client interface {
	// Complete sends a request and returns a response
	Complete(ctx context.Context, req *Request) (*Response, error)

	// Provider returns the provider name
	Provider() Provider
}

// StreamingResponse represents a chunk of a streaming response
type StreamingResponse struct {
	Delta    string // The new text chunk
	Done     bool   // Whether streaming is complete
	Command  string // Extracted command (only set when Done)
	FullText string // Full accumulated text (only set when Done)
	Err      error  // Error that occurred during streaming
}

// StreamingClient extends Client with streaming support
type StreamingClient interface {
	Client

	// CompleteStream sends a request and streams the response
	CompleteStream(ctx context.Context, req *Request, handler func(StreamingResponse)) error
}

// Service manages AI interactions
type Service struct {
	client   Client
	provider Provider
}

// NewService creates a new AI service
func NewService(provider Provider, apiKey string) (*Service, error) {
	// Ollama doesn't require an API key
	if provider != ProviderOllama && apiKey == "" {
		return nil, errors.New("API key is required")
	}

	var client Client
	var err error

	switch provider {
	case ProviderGemini:
		client, err = NewGeminiClient(apiKey)
	case ProviderAnthropic:
		client, err = NewAnthropicClient(apiKey)
	case ProviderOpenAI:
		client, err = NewOpenAIClient(apiKey)
	case ProviderOllama:
		client, err = NewOllamaClient(nil) // Uses default config
	default:
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}

	if err != nil {
		return nil, err
	}

	return &Service{
		client:   client,
		provider: provider,
	}, nil
}

// NewServiceWithGemini creates an AI service with custom Gemini model
func NewServiceWithGemini(apiKey, model string) (*Service, error) {
	client, err := NewGeminiClientWithModel(apiKey, model)
	if err != nil {
		return nil, err
	}

	return &Service{
		client:   client,
		provider: ProviderGemini,
	}, nil
}

// NewServiceWithOllama creates an AI service with custom Ollama config
func NewServiceWithOllama(config *OllamaConfig) (*Service, error) {
	client, err := NewOllamaClient(config)
	if err != nil {
		return nil, err
	}

	return &Service{
		client:   client,
		provider: ProviderOllama,
	}, nil
}

// ProcessInput processes user input and returns an AI response
func (s *Service) ProcessInput(ctx context.Context, input string) (*Response, error) {
	return s.ProcessInputWithContext(ctx, input, nil)
}

// ProcessInputWithContext processes user input with additional project context
func (s *Service) ProcessInputWithContext(ctx context.Context, input string, projectCtx *RequestContext) (*Response, error) {
	reqType := detectRequestType(input)

	reqContext := getRequestContext()
	// Merge project context if provided
	if projectCtx != nil {
		reqContext.ProjectName = projectCtx.ProjectName
		reqContext.ProjectType = projectCtx.ProjectType
		reqContext.ProjectFramework = projectCtx.ProjectFramework
		if projectCtx.CurrentDir != "" {
			reqContext.CurrentDir = projectCtx.CurrentDir
		}
	}

	req := &Request{
		UserMessage: input,
		Context:     reqContext,
		Type:        reqType,
	}

	return s.client.Complete(ctx, req)
}

// SuggestCommand converts natural language to a shell command
func (s *Service) SuggestCommand(ctx context.Context, description string) (*Response, error) {
	req := &Request{
		UserMessage: description,
		Context:     getRequestContext(),
		Type:        RequestTypeCommand,
	}

	return s.client.Complete(ctx, req)
}

// ExplainError explains an error message
func (s *Service) ExplainError(ctx context.Context, errorMsg string) (*Response, error) {
	req := &Request{
		UserMessage: fmt.Sprintf("Explain this error and suggest a fix:\n%s", errorMsg),
		Context:     getRequestContext(),
		Type:        RequestTypeExplain,
	}

	return s.client.Complete(ctx, req)
}

// detectRequestType determines what type of request the user is making
func detectRequestType(input string) RequestType {
	lower := strings.ToLower(input)

	// Command-related keywords
	commandKeywords := []string{
		"how to", "how do i", "command for", "show me how",
		"wie kann ich", "zeige mir", "befehl für",
		"run", "execute", "find", "list", "delete", "create",
		"finde", "lösche", "erstelle", "zeige",
	}

	for _, kw := range commandKeywords {
		if strings.Contains(lower, kw) {
			return RequestTypeCommand
		}
	}

	// Explain-related keywords
	explainKeywords := []string{
		"explain", "what does", "what is", "why",
		"erkläre", "was bedeutet", "was ist", "warum",
		"error", "fehler", "problem",
	}

	for _, kw := range explainKeywords {
		if strings.Contains(lower, kw) {
			return RequestTypeExplain
		}
	}

	return RequestTypeChat
}

// getRequestContext gathers context about the current environment
func getRequestContext() *RequestContext {
	shell := "bash"
	if runtime.GOOS == "windows" {
		shell = "powershell"
	}

	return &RequestContext{
		CurrentDir: ".", // Will be set by caller if needed
		OS:         runtime.GOOS,
		Shell:      shell,
	}
}

// GetProvider returns the current provider
func (s *Service) GetProvider() Provider {
	return s.provider
}

// GetClient returns the underlying AI client
func (s *Service) GetClient() Client {
	return s.client
}

// SupportsStreaming returns true if the current provider supports streaming
func (s *Service) SupportsStreaming() bool {
	_, ok := s.client.(StreamingClient)
	return ok
}

// ProcessInputStreaming processes user input with streaming response
func (s *Service) ProcessInputStreaming(ctx context.Context, input string, handler func(StreamingResponse)) error {
	streamClient, ok := s.client.(StreamingClient)
	if !ok {
		return fmt.Errorf("streaming not supported by provider %s", s.provider)
	}

	reqType := detectRequestType(input)
	reqContext := getRequestContext()

	req := &Request{
		UserMessage: input,
		Context:     reqContext,
		Type:        reqType,
	}

	return streamClient.CompleteStream(ctx, req, handler)
}
