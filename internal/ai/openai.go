package ai

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// OpenAIClient implements the Client interface for OpenAI
type OpenAIClient struct {
	client *openai.Client
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient(apiKey string) (*OpenAIClient, error) {
	client := openai.NewClient(apiKey)
	return &OpenAIClient{client: client}, nil
}

// Provider returns the provider name
func (c *OpenAIClient) Provider() Provider {
	return ProviderOpenAI
}

// Complete sends a request to OpenAI and returns a response
func (c *OpenAIClient) Complete(ctx context.Context, req *Request) (*Response, error) {
	systemPrompt := buildSystemPrompt(req)
	userMessage := buildUserMessage(req)

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userMessage,
			},
		},
		MaxTokens:   1024,
		Temperature: 0.7,
	})

	if err != nil {
		return nil, fmt.Errorf("openai API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	content := resp.Choices[0].Message.Content
	return parseResponse(content, req.Type, ProviderOpenAI), nil
}

// CompleteStream sends a streaming request to OpenAI
func (c *OpenAIClient) CompleteStream(ctx context.Context, req *Request, handler func(StreamingResponse)) error {
	systemPrompt := buildSystemPrompt(req)
	userMessage := buildUserMessage(req)

	stream, err := c.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userMessage,
			},
		},
		MaxTokens:   1024,
		Temperature: 0.7,
		Stream:      true,
	})

	if err != nil {
		return fmt.Errorf("openai stream error: %w", err)
	}
	defer stream.Close()

	var fullText strings.Builder

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			// Stream completed
			finalText := fullText.String()
			handler(StreamingResponse{
				Done:     true,
				Command:  extractCommand(finalText),
				FullText: finalText,
			})
			return nil
		}

		if err != nil {
			return fmt.Errorf("stream recv error: %w", err)
		}

		if len(response.Choices) > 0 {
			delta := response.Choices[0].Delta.Content
			if delta != "" {
				fullText.WriteString(delta)
				handler(StreamingResponse{
					Delta: delta,
					Done:  false,
				})
			}
		}
	}
}
