package llm

import (
	"context"

	"github.com/pashpashpash/vault/chunk"
)

// Client is an interface for LLM providers (OpenAI, Ollama, etc.)
type Client interface {
	// GetEmbedding gets an embedding for a single text
	GetEmbedding(ctx context.Context, text string) ([]float32, error)

	// GetEmbeddings gets embeddings for multiple chunks
	GetEmbeddings(chunks []chunk.Chunk, batchSize int) ([][]float32, error)

	// CreateChatCompletionSimple creates a chat completion with the given messages and options
	CreateChatCompletionSimple(ctx context.Context, prompt string, systemPrompt string, maxTokens int, temperature float32) (string, error)
}
