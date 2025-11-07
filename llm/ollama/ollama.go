package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/pashpashpash/vault/chunk"
)

const (
	DefaultEndpoint      = "http://localhost:11434"
	DefaultEmbedModel    = "nomic-embed-text"
	DefaultChatModel     = "llama3"
	EmbeddingDimension   = 768 // nomic-embed-text produces 768-dimensional vectors
	DefaultRetryAttempts = 3
	RetryDelay           = 5 * time.Second
)

// Client represents an Ollama API client
type Client struct {
	Endpoint   string
	EmbedModel string
	ChatModel  string
	HTTPClient *http.Client
}

// NewClient creates a new Ollama client
func NewClient(endpoint, embedModel, chatModel string) *Client {
	if endpoint == "" {
		endpoint = DefaultEndpoint
	}
	if embedModel == "" {
		embedModel = DefaultEmbedModel
	}
	if chatModel == "" {
		chatModel = DefaultChatModel
	}

	return &Client{
		Endpoint:   endpoint,
		EmbedModel: embedModel,
		ChatModel:  chatModel,
		HTTPClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// EmbedRequest represents a request to the Ollama embeddings API
type EmbedRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// EmbedResponse represents a response from the Ollama embeddings API
type EmbedResponse struct {
	Embedding []float32 `json:"embedding"`
}

// ChatRequest represents a request to the Ollama chat API
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
	Options  *ChatOptions  `json:"options,omitempty"`
}

// ChatMessage represents a message in the chat
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatOptions represents options for the chat request
type ChatOptions struct {
	Temperature      float32  `json:"temperature,omitempty"`
	TopP             float32  `json:"top_p,omitempty"`
	NumPredict       int      `json:"num_predict,omitempty"`
	Stop             []string `json:"stop,omitempty"`
	FrequencyPenalty float32  `json:"frequency_penalty,omitempty"`
	PresencePenalty  float32  `json:"presence_penalty,omitempty"`
}

// ChatResponse represents a response from the Ollama chat API
type ChatResponse struct {
	Model     string      `json:"model"`
	CreatedAt string      `json:"created_at"`
	Message   ChatMessage `json:"message"`
	Done      bool        `json:"done"`
}

// GetEmbedding gets an embedding for a single text
func (c *Client) GetEmbedding(ctx context.Context, text string) ([]float32, error) {
	req := EmbedRequest{
		Model:  c.EmbedModel,
		Prompt: text,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.Endpoint+"/api/embeddings", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var embedResp EmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return embedResp.Embedding, nil
}

// GetEmbeddingsWithRetry gets embeddings for multiple texts with retry logic
func (c *Client) GetEmbeddingsWithRetry(texts []string, maxRetries int) ([][]float32, error) {
	embeddings := make([][]float32, 0, len(texts))

	for _, text := range texts {
		var embedding []float32
		var err error

		for i := 0; i < maxRetries; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			embedding, err = c.GetEmbedding(ctx, text)
			cancel()

			if err == nil {
				break
			}

			log.Printf("[Ollama] Embedding request failed (attempt %d/%d): %v", i+1, maxRetries, err)
			if i < maxRetries-1 {
				time.Sleep(RetryDelay)
			}
		}

		if err != nil {
			return nil, fmt.Errorf("failed to get embedding after %d attempts: %w", maxRetries, err)
		}

		embeddings = append(embeddings, embedding)
	}

	return embeddings, nil
}

// GetEmbeddings gets embeddings for chunks with batch processing
func (c *Client) GetEmbeddings(chunks []chunk.Chunk, batchSize int) ([][]float32, error) {
	embeddings := make([][]float32, 0, len(chunks))

	for i := 0; i < len(chunks); i += batchSize {
		iEnd := min(len(chunks), i+batchSize)

		texts := make([]string, 0, iEnd-i)
		for _, chunk := range chunks[i:iEnd] {
			texts = append(texts, chunk.Text)
		}

		log.Printf("[Ollama] Getting embeddings for batch %d-%d of %d chunks", i, iEnd, len(chunks))

		batchEmbeddings, err := c.GetEmbeddingsWithRetry(texts, DefaultRetryAttempts)
		if err != nil {
			return nil, err
		}

		embeddings = append(embeddings, batchEmbeddings...)
	}

	return embeddings, nil
}

// CreateChatCompletion creates a chat completion
func (c *Client) CreateChatCompletion(ctx context.Context, messages []ChatMessage, options *ChatOptions) (string, error) {
	req := ChatRequest{
		Model:    c.ChatModel,
		Messages: messages,
		Stream:   false,
		Options:  options,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.Endpoint+"/api/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return chatResp.Message.Content, nil
}

// CreateChatCompletionSimple creates a chat completion with simplified parameters (implements llm.Client interface)
func (c *Client) CreateChatCompletionSimple(ctx context.Context, prompt string, systemPrompt string, maxTokens int, temperature float32) (string, error) {
	messages := []ChatMessage{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	options := &ChatOptions{
		Temperature: temperature,
		NumPredict:  maxTokens,
	}

	return c.CreateChatCompletion(ctx, messages, options)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
