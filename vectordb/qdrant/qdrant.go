package qdrant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/itaiguardiola/askara/chunk"
	"github.com/itaiguardiola/askara/vectordb"
	cache "github.com/patrickmn/go-cache"
)

const (
	VECTOR_SIZE     = 768 // nomic-embed-text (local Ollama model)
	VECTOR_DISTANCE = "Cosine"
	BATCH_SIZE      = 500
)

type Qdrant struct {
	Endpoint string
	cache    *cache.Cache
}

type Point struct {
	ID      int               `json:"id"`
	Vector  []float32         `json:"vector"`
	Payload map[string]string `json:"payload,omitempty"`
}

type Match struct {
	ID      int               `json:"id"`
	Score   float32           `json:"score"`
	Payload map[string]string `json:"payload"`
	Version int               `json:"version"`
}

type SearchResult struct {
	Result []Match `json:"result"`
	Status string  `json:"status"`
	Time   float64 `json:"time"`
}

type NamespaceConfig struct {
	Vectors struct {
		Size     int    `json:"size"`
		Distance string `json:"distance"`
	} `json:"vectors"`
}

func New(endpoint string) (*Qdrant, error) {
	return &Qdrant{
		Endpoint: endpoint,
		cache:    cache.New(5*time.Minute, 10*time.Minute),
	}, nil
}

func (q *Qdrant) NamespaceExists(uuid string) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/collections/%s", q.Endpoint, uuid), nil)
	if err != nil {
		return false, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	} else if resp.StatusCode != http.StatusOK {
		body := make([]byte, 1024)
		_, _ = resp.Body.Read(body)
		return false, fmt.Errorf("failed to check namespace, status code: %d, %s", resp.StatusCode, body)
	}

	return true, nil
}

func (q *Qdrant) CreateNamespace(uuid string) error {
	if _, found := q.cache.Get(uuid); found {
		return nil
	}

	if exists, err := q.NamespaceExists(uuid); err != nil {
		return err
	} else if exists {
		return nil
	}

	config := NamespaceConfig{}
	config.Vectors.Size = VECTOR_SIZE
	config.Vectors.Distance = VECTOR_DISTANCE

	jsonData, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/collections/%s", q.Endpoint, uuid), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create namespace, status code: %d", resp.StatusCode)
	}

	q.cache.Set(uuid, true, cache.DefaultExpiration)

	return nil
}

func (q *Qdrant) UpsertEmbeddings(embeddings [][]float32, chunks []chunk.Chunk, uuid string) error {
	if err := q.CreateNamespace(uuid); err != nil {
		return err
	}

	points := make([]Point, len(embeddings))

	for i, embedding := range embeddings {
		points[i].ID = i
		points[i].Vector = embedding
		if i < len(chunks) {
			points[i].Payload = map[string]string{
				"start": fmt.Sprintf("%d", chunks[i].Start),
				"end":   fmt.Sprintf("%d", chunks[i].End),
				"title": chunks[i].Title,
				"text":  chunks[i].Text,
			}
		}
	}

	for i := 0; i < len(points); i += BATCH_SIZE {
		end := i + BATCH_SIZE
		if end > len(points) {
			end = len(points)
		}

		data := map[string][]Point{"points": points[i:end]}
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/collections/%s/points", q.Endpoint, uuid), bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to upsert embeddings, status code: %d", resp.StatusCode)
		}
	}

	return nil
}

func (q *Qdrant) Retrieve(questionEmbedding []float32, topK int, uuid string) ([]vectordb.QueryMatch, error) {
	data := map[string]interface{}{
		"vector":       questionEmbedding,
		"top":          topK,
		"with_payload": true,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/collections/%s/points/search", q.Endpoint, uuid), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve embeddings, status code: %d", resp.StatusCode)
	}

	var searchResult SearchResult
	err = json.NewDecoder(resp.Body).Decode(&searchResult)
	if err != nil {
		return nil, err
	}

	// Convert qdrantMatch to QueryMatch
	queryMatches := make([]vectordb.QueryMatch, len(searchResult.Result))
	for i, result := range searchResult.Result {
		queryMatches[i].ID = fmt.Sprintf("%d", result.ID)
		queryMatches[i].Score = result.Score
		queryMatches[i].Metadata = result.Payload
	}

	return queryMatches, nil
}

// ListCollections lists all collections in Qdrant
func (q *Qdrant) ListCollections() ([]vectordb.CollectionInfo, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/collections", q.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list collections, status code: %d", resp.StatusCode)
	}

	var result struct {
		Result struct {
			Collections []struct {
				Name string `json:"name"`
			} `json:"collections"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Get detailed info for each collection
	collections := make([]vectordb.CollectionInfo, 0)
	for _, coll := range result.Result.Collections {
		info, err := q.GetCollectionInfo(coll.Name)
		if err != nil {
			// Skip collections we can't read
			continue
		}
		collections = append(collections, *info)
	}

	return collections, nil
}

// GetCollectionInfo gets information about a specific collection
func (q *Qdrant) GetCollectionInfo(uuid string) (*vectordb.CollectionInfo, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/collections/%s", q.Endpoint, uuid), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("collection not found: %s", uuid)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get collection info, status code: %d", resp.StatusCode)
	}

	var result struct {
		Result struct {
			PointsCount int `json:"points_count"`
			Config      struct {
				Params struct {
					Vectors struct {
						Size int `json:"size"`
					} `json:"vectors"`
				} `json:"params"`
			} `json:"config"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &vectordb.CollectionInfo{
		Name:       uuid,
		PointCount: result.Result.PointsCount,
		VectorSize: result.Result.Config.Params.Vectors.Size,
	}, nil
}

// GetDocuments gets all unique documents in a collection
func (q *Qdrant) GetDocuments(uuid string) ([]vectordb.DocumentInfo, error) {
	// Scroll through all points to get unique document names
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/collections/%s/points/scroll", q.Endpoint, uuid),
		bytes.NewBuffer([]byte(`{"limit":1000,"with_payload":true,"with_vector":false}`)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to scroll points, status code: %d", resp.StatusCode)
	}

	var result struct {
		Result struct {
			Points []struct {
				Payload map[string]interface{} `json:"payload"`
			} `json:"points"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Group by document name (title)
	docMap := make(map[string]*vectordb.DocumentInfo)
	for _, point := range result.Result.Points {
		if title, ok := point.Payload["title"].(string); ok && title != "" {
			if _, exists := docMap[title]; !exists {
				docMap[title] = &vectordb.DocumentInfo{
					Name:       title,
					ChunkCount: 0,
					Titles:     []string{title},
				}
			}
			docMap[title].ChunkCount++
		}
	}

	// Convert map to slice
	documents := make([]vectordb.DocumentInfo, 0, len(docMap))
	for _, doc := range docMap {
		documents = append(documents, *doc)
	}

	return documents, nil
}

// DeleteCollection deletes an entire collection
func (q *Qdrant) DeleteCollection(uuid string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/collections/%s", q.Endpoint, uuid), nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body := make([]byte, 1024)
		resp.Body.Read(body)
		return fmt.Errorf("failed to delete collection, status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Remove from cache
	q.cache.Delete(uuid)

	return nil
}
