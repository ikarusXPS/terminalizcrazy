package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultOllamaConfig(t *testing.T) {
	config := DefaultOllamaConfig()

	assert.NotNil(t, config)
	assert.Equal(t, "http://localhost:11434", config.BaseURL)
	assert.Equal(t, "codellama", config.Model)
	assert.Equal(t, 120*time.Second, config.Timeout)
}

func TestNewOllamaClient(t *testing.T) {
	t.Run("with nil config", func(t *testing.T) {
		client, err := NewOllamaClient(nil)

		require.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, "http://localhost:11434", client.baseURL)
		assert.Equal(t, "codellama", client.model)
	})

	t.Run("with custom config", func(t *testing.T) {
		config := &OllamaConfig{
			BaseURL: "http://custom:8080",
			Model:   "llama2",
			Timeout: 60 * time.Second,
		}

		client, err := NewOllamaClient(config)

		require.NoError(t, err)
		assert.Equal(t, "http://custom:8080", client.baseURL)
		assert.Equal(t, "llama2", client.model)
	})

	t.Run("with empty config values", func(t *testing.T) {
		config := &OllamaConfig{
			BaseURL: "",
			Model:   "",
			Timeout: 0,
		}

		client, err := NewOllamaClient(config)

		require.NoError(t, err)
		assert.Equal(t, "http://localhost:11434", client.baseURL)
		assert.Equal(t, "codellama", client.model)
	})

	t.Run("strips trailing slash", func(t *testing.T) {
		config := &OllamaConfig{
			BaseURL: "http://localhost:11434/",
			Model:   "test",
		}

		client, err := NewOllamaClient(config)

		require.NoError(t, err)
		assert.Equal(t, "http://localhost:11434", client.baseURL)
	})
}

func TestOllamaClient_Provider(t *testing.T) {
	client, _ := NewOllamaClient(nil)

	assert.Equal(t, ProviderOllama, client.Provider())
}

func TestOllamaClient_SetGetModel(t *testing.T) {
	client, _ := NewOllamaClient(nil)

	assert.Equal(t, "codellama", client.GetModel())

	client.SetModel("llama2")
	assert.Equal(t, "llama2", client.GetModel())

	client.SetModel("mistral")
	assert.Equal(t, "mistral", client.GetModel())
}

func TestOllamaClient_buildSystemPrompt(t *testing.T) {
	client, _ := NewOllamaClient(nil)

	tests := []struct {
		name     string
		reqType  RequestType
		contains []string
	}{
		{
			name:     "command request",
			reqType:  RequestTypeCommand,
			contains: []string{"command-line", "COMMAND:", "EXPLANATION:"},
		},
		{
			name:     "explain request",
			reqType:  RequestTypeExplain,
			contains: []string{"explain", "clear", "concise"},
		},
		{
			name:     "chat request",
			reqType:  RequestTypeChat,
			contains: []string{"technical", "helpful"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{Type: tt.reqType}
			prompt := client.buildSystemPrompt(req)

			for _, expected := range tt.contains {
				assert.Contains(t, prompt, expected)
			}
		})
	}
}

func TestOllamaClient_buildUserPrompt(t *testing.T) {
	client, _ := NewOllamaClient(nil)

	t.Run("without context", func(t *testing.T) {
		req := &Request{UserMessage: "Hello world"}
		prompt := client.buildUserPrompt(req)

		assert.Equal(t, "Hello world", prompt)
	})

	t.Run("with context", func(t *testing.T) {
		req := &Request{
			UserMessage: "list files",
			Context: &RequestContext{
				OS:         "linux",
				Shell:      "bash",
				CurrentDir: "/home/user",
			},
		}
		prompt := client.buildUserPrompt(req)

		assert.Contains(t, prompt, "linux")
		assert.Contains(t, prompt, "bash")
		assert.Contains(t, prompt, "/home/user")
		assert.Contains(t, prompt, "list files")
	})

	t.Run("with project context", func(t *testing.T) {
		req := &Request{
			UserMessage: "build",
			Context: &RequestContext{
				OS:          "darwin",
				Shell:       "zsh",
				CurrentDir:  "/projects/myapp",
				ProjectName: "MyApp",
				ProjectType: "node",
			},
		}
		prompt := client.buildUserPrompt(req)

		assert.Contains(t, prompt, "MyApp")
		assert.Contains(t, prompt, "node")
	})

	t.Run("with recent history", func(t *testing.T) {
		req := &Request{
			UserMessage: "run tests",
			Context: &RequestContext{
				OS:            "linux",
				Shell:         "bash",
				CurrentDir:    "/app",
				RecentHistory: []string{"npm install", "npm run build"},
			},
		}
		prompt := client.buildUserPrompt(req)

		assert.Contains(t, prompt, "Recent commands")
		assert.Contains(t, prompt, "npm install")
		assert.Contains(t, prompt, "npm run build")
	})
}

func TestOllamaClient_parseResponse(t *testing.T) {
	client, _ := NewOllamaClient(nil)

	t.Run("chat response", func(t *testing.T) {
		resp, err := client.parseResponse("Hello, how can I help?", RequestTypeChat)

		require.NoError(t, err)
		assert.Equal(t, "Hello, how can I help?", resp.Content)
		assert.Equal(t, ProviderOllama, resp.Provider)
		assert.Equal(t, 0.7, resp.Confidence)
	})

	t.Run("command response with COMMAND prefix", func(t *testing.T) {
		content := "COMMAND: ls -la\nEXPLANATION: Lists all files"
		resp, err := client.parseResponse(content, RequestTypeCommand)

		require.NoError(t, err)
		assert.Equal(t, "ls -la", resp.Command)
		assert.Equal(t, "Lists all files", resp.Explanation)
	})

	t.Run("command response with code block", func(t *testing.T) {
		content := "Here's the command:\n```bash\ngit status\n```"
		resp, err := client.parseResponse(content, RequestTypeCommand)

		require.NoError(t, err)
		assert.Equal(t, "git status", resp.Command)
	})

	t.Run("command response with $ prefix", func(t *testing.T) {
		content := "Run this:\n$ echo hello"
		resp, err := client.parseResponse(content, RequestTypeCommand)

		require.NoError(t, err)
		assert.Equal(t, "echo hello", resp.Command)
	})
}

func TestExtractCommandFromContent(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "code block",
			content:  "```bash\nls -la\n```",
			expected: "ls -la",
		},
		{
			name:     "dollar prefix",
			content:  "$ git status",
			expected: "git status",
		},
		{
			name:     "no command",
			content:  "Just plain text",
			expected: "",
		},
		{
			name:     "empty code block",
			content:  "```\n```",
			expected: "",
		},
		{
			name:     "multiple lines in code block",
			content:  "```\nfirst line\nsecond line\n```",
			expected: "first line",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractCommandFromContent(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestOllamaClient_Complete(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/generate", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			var reqBody ollamaRequest
			json.NewDecoder(r.Body).Decode(&reqBody)
			assert.Equal(t, "codellama", reqBody.Model)
			assert.False(t, reqBody.Stream)

			resp := ollamaResponse{
				Model:    "codellama",
				Response: "COMMAND: ls -la\nEXPLANATION: List files",
				Done:     true,
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client, _ := NewOllamaClient(&OllamaConfig{BaseURL: server.URL})
		req := &Request{UserMessage: "list files", Type: RequestTypeCommand}

		resp, err := client.Complete(context.Background(), req)

		require.NoError(t, err)
		assert.Equal(t, "ls -la", resp.Command)
	})

	t.Run("server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}))
		defer server.Close()

		client, _ := NewOllamaClient(&OllamaConfig{BaseURL: server.URL})
		req := &Request{UserMessage: "test", Type: RequestTypeChat}

		resp, err := client.Complete(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "500")
	})

	t.Run("ollama error in response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := ollamaResponse{
				Error: "model not found",
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client, _ := NewOllamaClient(&OllamaConfig{BaseURL: server.URL})
		req := &Request{UserMessage: "test", Type: RequestTypeChat}

		resp, err := client.Complete(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "model not found")
	})
}

func TestOllamaClient_IsAvailable(t *testing.T) {
	t.Run("available", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/tags", r.URL.Path)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client, _ := NewOllamaClient(&OllamaConfig{BaseURL: server.URL})

		assert.True(t, client.IsAvailable(context.Background()))
	})

	t.Run("not available", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		defer server.Close()

		client, _ := NewOllamaClient(&OllamaConfig{BaseURL: server.URL})

		assert.False(t, client.IsAvailable(context.Background()))
	})

	t.Run("connection refused", func(t *testing.T) {
		client, _ := NewOllamaClient(&OllamaConfig{BaseURL: "http://localhost:99999"})

		assert.False(t, client.IsAvailable(context.Background()))
	})
}

func TestOllamaClient_ListModels(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/tags", r.URL.Path)

			resp := ollamaModelsResponse{
				Models: []ollamaModel{
					{Name: "codellama:7b"},
					{Name: "llama2:13b"},
					{Name: "mistral:latest"},
				},
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client, _ := NewOllamaClient(&OllamaConfig{BaseURL: server.URL})

		models, err := client.ListModels(context.Background())

		require.NoError(t, err)
		assert.Len(t, models, 3)
		assert.Contains(t, models, "codellama:7b")
		assert.Contains(t, models, "llama2:13b")
		assert.Contains(t, models, "mistral:latest")
	})

	t.Run("server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client, _ := NewOllamaClient(&OllamaConfig{BaseURL: server.URL})

		models, err := client.ListModels(context.Background())

		assert.Error(t, err)
		assert.Nil(t, models)
	})
}

func TestOllamaClient_PullModel(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/pull", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			var reqBody map[string]string
			json.NewDecoder(r.Body).Decode(&reqBody)
			assert.Equal(t, "llama2", reqBody["name"])

			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client, _ := NewOllamaClient(&OllamaConfig{BaseURL: server.URL})

		err := client.PullModel(context.Background(), "llama2")

		assert.NoError(t, err)
	})

	t.Run("server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("model not found"))
		}))
		defer server.Close()

		client, _ := NewOllamaClient(&OllamaConfig{BaseURL: server.URL})

		err := client.PullModel(context.Background(), "nonexistent")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "model not found")
	})
}

func TestOllamaClient_CompleteStream(t *testing.T) {
	t.Run("successful stream", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/generate", r.URL.Path)

			var reqBody ollamaRequest
			json.NewDecoder(r.Body).Decode(&reqBody)
			assert.True(t, reqBody.Stream)

			// Send streaming response
			chunks := []ollamaResponse{
				{Response: "Hello", Done: false},
				{Response: " world", Done: false},
				{Response: "!", Done: true},
			}

			for _, chunk := range chunks {
				json.NewEncoder(w).Encode(chunk)
			}
		}))
		defer server.Close()

		client, _ := NewOllamaClient(&OllamaConfig{BaseURL: server.URL})
		req := &Request{UserMessage: "greet", Type: RequestTypeChat}

		var responses []StreamingResponse
		err := client.CompleteStream(context.Background(), req, func(resp StreamingResponse) {
			responses = append(responses, resp)
		})

		require.NoError(t, err)
		assert.Len(t, responses, 3)
		assert.Equal(t, "Hello", responses[0].Content)
		assert.Equal(t, " world", responses[1].Content)
		assert.Equal(t, "!", responses[2].Content)
		assert.True(t, responses[2].Done)
	})

	t.Run("server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client, _ := NewOllamaClient(&OllamaConfig{BaseURL: server.URL})
		req := &Request{UserMessage: "test", Type: RequestTypeChat}

		err := client.CompleteStream(context.Background(), req, func(resp StreamingResponse) {})

		assert.Error(t, err)
	})
}

func TestStreamingResponse(t *testing.T) {
	resp := StreamingResponse{
		Content: "test content",
		Done:    true,
		Error:   nil,
	}

	assert.Equal(t, "test content", resp.Content)
	assert.True(t, resp.Done)
	assert.Nil(t, resp.Error)
}
