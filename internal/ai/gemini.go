package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// ProviderGemini is the Gemini provider constant
const ProviderGemini Provider = "gemini"

// GeminiClient implements the Client interface for Google Gemini
type GeminiClient struct {
	client *genai.Client
	model  string
}

// NewGeminiClient creates a new Gemini client
func NewGeminiClient(apiKey string) (*GeminiClient, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &GeminiClient{
		client: client,
		model:  "gemini-1.5-flash", // Default model
	}, nil
}

// NewGeminiClientWithModel creates a new Gemini client with a specific model
func NewGeminiClientWithModel(apiKey, model string) (*GeminiClient, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	if model == "" {
		model = "gemini-1.5-flash"
	}

	return &GeminiClient{
		client: client,
		model:  model,
	}, nil
}

// Provider returns the provider name
func (c *GeminiClient) Provider() Provider {
	return ProviderGemini
}

// Complete sends a request to Gemini and returns a response
func (c *GeminiClient) Complete(ctx context.Context, req *Request) (*Response, error) {
	model := c.client.GenerativeModel(c.model)

	// Set up the system instruction
	systemPrompt := buildSystemPrompt(req)
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemPrompt)},
	}

	// Configure safety settings to be more permissive for technical content
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockOnlyHigh,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockOnlyHigh,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockOnlyHigh,
		},
	}

	// Generate content
	resp, err := model.GenerateContent(ctx, genai.Text(req.UserMessage))
	if err != nil {
		return nil, fmt.Errorf("gemini API error: %w", err)
	}

	// Extract text from response
	var content strings.Builder
	for _, candidate := range resp.Candidates {
		if candidate.Content != nil {
			for _, part := range candidate.Content.Parts {
				if text, ok := part.(genai.Text); ok {
					content.WriteString(string(text))
				}
			}
		}
	}

	return parseResponse(content.String(), req.Type, ProviderGemini), nil
}

// CompleteStream sends a streaming request to Gemini
func (c *GeminiClient) CompleteStream(ctx context.Context, req *Request, handler func(StreamingResponse)) error {
	model := c.client.GenerativeModel(c.model)

	// Set up the system instruction
	systemPrompt := buildSystemPrompt(req)
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemPrompt)},
	}

	// Configure safety settings
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockOnlyHigh,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockOnlyHigh,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockOnlyHigh,
		},
	}

	// Stream content
	iter := model.GenerateContentStream(ctx, genai.Text(req.UserMessage))

	var fullText strings.Builder

	for {
		resp, err := iter.Next()
		if err != nil {
			// Check if we've reached the end of the stream
			if err.Error() == "iterator done" || strings.Contains(err.Error(), "iterator") {
				finalText := fullText.String()
				handler(StreamingResponse{
					Done:     true,
					Command:  extractCommand(finalText),
					FullText: finalText,
				})
				return nil
			}
			return fmt.Errorf("stream error: %w", err)
		}

		// Extract text from response chunk
		for _, candidate := range resp.Candidates {
			if candidate.Content != nil {
				for _, part := range candidate.Content.Parts {
					if text, ok := part.(genai.Text); ok {
						chunk := string(text)
						fullText.WriteString(chunk)
						handler(StreamingResponse{
							Delta: chunk,
							Done:  false,
						})
					}
				}
			}
		}
	}
}

// Close closes the Gemini client
func (c *GeminiClient) Close() error {
	return c.client.Close()
}

// SetModel changes the model used by the client
func (c *GeminiClient) SetModel(model string) {
	c.model = model
}

// GetModel returns the current model
func (c *GeminiClient) GetModel() string {
	return c.model
}

// ListModels returns available Gemini models
func (c *GeminiClient) ListModels() []string {
	return []string{
		"gemini-1.5-flash",
		"gemini-1.5-flash-8b",
		"gemini-1.5-pro",
		"gemini-1.0-pro",
		"gemini-2.0-flash-exp",
	}
}
