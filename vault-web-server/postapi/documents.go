package postapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// DocumentsListHandler lists all collections and their documents
func (ctx *HandlerContext) DocumentsListHandler(w http.ResponseWriter, r *http.Request) {
	collections, err := ctx.vectorDB.ListCollections()
	if err != nil {
		log.Println("[DocumentsListHandler ERR] Error listing collections:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get documents for each collection
	type CollectionWithDocs struct {
		Collection string                    `json:"collection"`
		PointCount int                       `json:"point_count"`
		Documents  []map[string]interface{}  `json:"documents"`
	}

	result := make([]CollectionWithDocs, 0)
	for _, coll := range collections {
		docs, err := ctx.vectorDB.GetDocuments(coll.Name)
		if err != nil {
			log.Printf("[DocumentsListHandler WARN] Error getting documents for collection %s: %v", coll.Name, err)
			continue
		}

		// Convert to map for JSON response
		docMaps := make([]map[string]interface{}, len(docs))
		for i, doc := range docs {
			docMaps[i] = map[string]interface{}{
				"name":        doc.Name,
				"chunk_count": doc.ChunkCount,
			}
		}

		result = append(result, CollectionWithDocs{
			Collection: coll.Name,
			PointCount: coll.PointCount,
			Documents:  docMaps,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// DocumentsGetHandler gets documents for a specific collection (UUID)
func (ctx *HandlerContext) DocumentsGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	if uuid == "" {
		http.Error(w, "UUID parameter is required", http.StatusBadRequest)
		return
	}

	documents, err := ctx.vectorDB.GetDocuments(uuid)
	if err != nil {
		log.Printf("[DocumentsGetHandler ERR] Error getting documents for UUID %s: %v", uuid, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(documents)
}

// DocumentsDeleteHandler deletes a collection (UUID)
func (ctx *HandlerContext) DocumentsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	if uuid == "" {
		http.Error(w, "UUID parameter is required", http.StatusBadRequest)
		return
	}

	// Delete the collection
	err := ctx.vectorDB.DeleteCollection(uuid)
	if err != nil {
		log.Printf("[DocumentsDeleteHandler ERR] Error deleting collection %s: %v", uuid, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("[DocumentsDeleteHandler] Successfully deleted collection %s", uuid)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Collection deleted successfully",
		"collection": uuid,
	})
}

// CollectionStatsHandler gets statistics for a specific collection
func (ctx *HandlerContext) CollectionStatsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	if uuid == "" {
		http.Error(w, "UUID parameter is required", http.StatusBadRequest)
		return
	}

	// Get collection info
	collInfo, err := ctx.vectorDB.GetCollectionInfo(uuid)
	if err != nil {
		log.Printf("[CollectionStatsHandler ERR] Error getting collection info for %s: %v", uuid, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get documents
	documents, err := ctx.vectorDB.GetDocuments(uuid)
	if err != nil {
		log.Printf("[CollectionStatsHandler ERR] Error getting documents for %s: %v", uuid, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate statistics
	response := map[string]interface{}{
		"collection":     uuid,
		"point_count":    collInfo.PointCount,
		"vector_size":    collInfo.VectorSize,
		"document_count": len(documents),
		"documents":      documents,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
