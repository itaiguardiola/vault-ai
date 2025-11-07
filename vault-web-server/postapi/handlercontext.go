package postapi

import (
	"github.com/pashpashpash/vault/llm"
	"github.com/pashpashpash/vault/vectordb"

	cache "github.com/patrickmn/go-cache"
)

type HandlerContext struct {
	llmClient llm.Client
	cache     *cache.Cache
	vectorDB  vectordb.VectorDB
}

func NewHandlerContext(llmClient llm.Client, vectorDB vectordb.VectorDB) *HandlerContext {
	return &HandlerContext{
		llmClient: llmClient,
		cache:     cache.New(cache.NoExpiration, cache.NoExpiration),
		vectorDB:  vectorDB,
	}
}
