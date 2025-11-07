package postapi

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
)

// Config represents the application configuration
type Config struct {
	OllamaEndpoint   string `json:"ollama_endpoint"`
	OllamaEmbedModel string `json:"ollama_embed_model"`
	OllamaChatModel  string `json:"ollama_chat_model"`
	QdrantEndpoint   string `json:"qdrant_endpoint"`
	MaxFileSize      int64  `json:"max_file_size_mb"`
	MaxTotalSize     int64  `json:"max_total_size_mb"`
}

// ConfigManager handles configuration persistence
type ConfigManager struct {
	config     *Config
	configPath string
	mu         sync.RWMutex
}

// NewConfigManager creates a new config manager
func NewConfigManager(configPath string) *ConfigManager {
	cm := &ConfigManager{
		configPath: configPath,
		config:     getDefaultConfig(),
	}
	cm.LoadConfig()
	return cm
}

// getDefaultConfig returns the default configuration
func getDefaultConfig() *Config {
	return &Config{
		OllamaEndpoint:   getEnvOrDefault("OLLAMA_ENDPOINT", "http://ollama:11434"),
		OllamaEmbedModel: getEnvOrDefault("OLLAMA_EMBED_MODEL", "nomic-embed-text"),
		OllamaChatModel:  getEnvOrDefault("OLLAMA_CHAT_MODEL", "llama3"),
		QdrantEndpoint:   getEnvOrDefault("QDRANT_API_ENDPOINT", "http://qdrant:6333"),
		MaxFileSize:      25,  // MB
		MaxTotalSize:     50,  // MB
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// LoadConfig loads configuration from file
func (cm *ConfigManager) LoadConfig() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Config file doesn't exist, use defaults
			log.Println("[ConfigManager] Config file not found, using defaults")
			return cm.SaveConfigUnlocked()
		}
		return err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		log.Printf("[ConfigManager] Error parsing config file, using defaults: %v", err)
		return nil
	}

	cm.config = &config
	log.Println("[ConfigManager] Configuration loaded successfully")
	return nil
}

// SaveConfigUnlocked saves configuration to file (must be called with lock held)
func (cm *ConfigManager) SaveConfigUnlocked() error {
	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return err
	}

	log.Println("[ConfigManager] Configuration saved successfully")
	return nil
}

// SaveConfig saves configuration to file
func (cm *ConfigManager) SaveConfig() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.SaveConfigUnlocked()
}

// GetConfig returns a copy of the current configuration
func (cm *ConfigManager) GetConfig() *Config {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// Return a copy to prevent external modification
	configCopy := *cm.config
	return &configCopy
}

// UpdateConfig updates the configuration
func (cm *ConfigManager) UpdateConfig(newConfig *Config) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.config = newConfig
	return cm.SaveConfigUnlocked()
}

// ConfigGetHandler handles GET requests for configuration
func (ctx *HandlerContext) ConfigGetHandler(w http.ResponseWriter, r *http.Request) {
	config := ctx.configManager.GetConfig()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(config); err != nil {
		log.Println("[ConfigGetHandler ERR] Error encoding config:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ConfigUpdateHandler handles POST requests to update configuration
func (ctx *HandlerContext) ConfigUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var newConfig Config
	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		log.Println("[ConfigUpdateHandler ERR] Error decoding request:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate configuration
	if newConfig.OllamaEndpoint == "" {
		http.Error(w, "ollama_endpoint is required", http.StatusBadRequest)
		return
	}
	if newConfig.QdrantEndpoint == "" {
		http.Error(w, "qdrant_endpoint is required", http.StatusBadRequest)
		return
	}
	if newConfig.OllamaEmbedModel == "" {
		newConfig.OllamaEmbedModel = "nomic-embed-text"
	}
	if newConfig.OllamaChatModel == "" {
		newConfig.OllamaChatModel = "llama3"
	}
	if newConfig.MaxFileSize <= 0 {
		newConfig.MaxFileSize = 25
	}
	if newConfig.MaxTotalSize <= 0 {
		newConfig.MaxTotalSize = 50
	}

	// Update configuration
	if err := ctx.configManager.UpdateConfig(&newConfig); err != nil {
		log.Println("[ConfigUpdateHandler ERR] Error updating config:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("[ConfigUpdateHandler] Configuration updated successfully")

	// Return updated config
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Configuration updated successfully. Please restart the application for changes to take effect.",
		"config":  newConfig,
	})
}

// ConfigTestHandler tests connectivity to Ollama and Qdrant
func (ctx *HandlerContext) ConfigTestHandler(w http.ResponseWriter, r *http.Request) {
	config := ctx.configManager.GetConfig()
	results := make(map[string]interface{})

	// Test Ollama connectivity
	ollamaResp, err := http.Get(config.OllamaEndpoint + "/api/tags")
	if err != nil {
		results["ollama"] = map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		}
	} else {
		defer ollamaResp.Body.Close()
		results["ollama"] = map[string]interface{}{
			"status":  "ok",
			"message": "Connected successfully",
		}
	}

	// Test Qdrant connectivity
	qdrantResp, err := http.Get(config.QdrantEndpoint + "/collections")
	if err != nil {
		results["qdrant"] = map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		}
	} else {
		defer qdrantResp.Body.Close()
		results["qdrant"] = map[string]interface{}{
			"status":  "ok",
			"message": "Connected successfully",
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
