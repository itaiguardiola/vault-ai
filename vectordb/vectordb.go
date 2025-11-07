package vectordb

import (
	"github.com/pashpashpash/vault/chunk"
)

type QueryMatch struct {
	ID       string            `json:"id"`
	Score    float32           `json:"score"` // Use "score" instead of "distance"
	Metadata map[string]string `json:"metadata"`
}

type CollectionInfo struct {
	Name       string `json:"name"`        // UUID/namespace
	PointCount int    `json:"point_count"` // Number of chunks
	VectorSize int    `json:"vector_size"`
}

type DocumentInfo struct {
	Name       string   `json:"name"`        // Document filename
	ChunkCount int      `json:"chunk_count"` // Number of chunks
	Titles     []string `json:"titles"`      // All unique titles (usually just one per file)
}

type VectorDB interface {
	UpsertEmbeddings(embeddings [][]float32, chunks []chunk.Chunk, uuid string) error
	Retrieve(questionEmbedding []float32, topK int, uuid string) ([]QueryMatch, error)

	// Document management methods
	ListCollections() ([]CollectionInfo, error)
	GetCollectionInfo(uuid string) (*CollectionInfo, error)
	GetDocuments(uuid string) ([]DocumentInfo, error)
	DeleteCollection(uuid string) error
}
