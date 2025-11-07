# Local Setup Guide - OP Vault

This guide will help you set up OP Vault to run **completely locally** without internet.

## Initial Setup (Requires Internet - One Time Only)

### 1. Download and Cache All Dependencies

```bash
# Download Go module dependencies
go mod download

# Install Node.js dependencies
npm install
```

This downloads all dependencies to your local cache. After this, you can work offline.

### 2. Install Local Services

**Install Ollama:**
```bash
# Linux
curl -fsSL https://ollama.com/install.sh | sh

# macOS
brew install ollama
```

**Download Ollama Models:**
```bash
# Start Ollama
ollama serve &

# Download required models
ollama pull nomic-embed-text  # 768-dim embeddings
ollama pull llama3            # 8B chat model (recommended)

# Optional: Download alternative models
ollama pull mistral           # 7B alternative
ollama pull llama3:70b        # 70B high-quality (requires 40GB+ RAM)
```

**Install Docker (for Qdrant):**
```bash
# Follow: https://docs.docker.com/get-docker/

# Pull Qdrant image
docker pull qdrant/qdrant
```

## Running Offline

Once the initial setup is complete, you can disconnect from the internet and run everything locally:

### Start Services

```bash
# Terminal 1: Start Qdrant (runs in background)
docker run -d -p 6333:6333 \
  -v $(pwd)/qdrant_storage:/qdrant/storage \
  --name vault-qdrant \
  qdrant/qdrant

# Terminal 2: Start Ollama (if not running as service)
ollama serve

# Terminal 3: Start the backend
npm start

# Terminal 4: Start frontend dev server
npm run dev
```

### Access Application

Open http://localhost:8100 in your browser

## Environment Variables

The application uses these defaults (no configuration needed):

```bash
OLLAMA_ENDPOINT="http://localhost:11434"
OLLAMA_EMBED_MODEL="nomic-embed-text"
OLLAMA_CHAT_MODEL="llama3"
QDRANT_API_ENDPOINT="http://localhost:6333"
```

Override by exporting before running:
```bash
export OLLAMA_CHAT_MODEL="mistral"
npm start
```

## Architecture Changes from Original

| Component | Original (Cloud) | New (Local) |
|-----------|-----------------|-------------|
| LLM | OpenAI GPT-3.5/4 | Ollama (llama3/mistral) |
| Embeddings | OpenAI ada-002 (1536-dim) | nomic-embed-text (768-dim) |
| Vector DB | Pinecone (cloud) | Qdrant (local) |
| Cost | ~$0.002/1K tokens | Free |
| Privacy | Data sent to cloud | Data stays local |
| Internet | Required | Not required |

## File Changes

### New Files Created:
- `llm/interface.go` - LLM client interface
- `llm/ollama/ollama.go` - Ollama client implementation

### Modified Files:
- `vault-web-server/main.go` - Initialize Ollama instead of OpenAI
- `vault-web-server/postapi/handlercontext.go` - Use LLM interface
- `vault-web-server/postapi/fileupload.go` - Use Ollama embeddings
- `vault-web-server/postapi/questions.go` - Use Ollama chat
- `vectordb/qdrant/qdrant.go` - Update vector dimensions to 768

### Removed Dependencies:
- OpenAI API (go-openai) - still in go.mod but not used
- Pinecone - still in code but not initialized

## Troubleshooting

### "Cannot connect to Ollama"
```bash
# Check if running
curl http://localhost:11434/api/tags

# Start Ollama
ollama serve
```

### "Cannot connect to Qdrant"
```bash
# Check if running
docker ps | grep qdrant

# Start Qdrant
docker start vault-qdrant
# OR
docker run -d -p 6333:6333 -v $(pwd)/qdrant_storage:/qdrant/storage qdrant/qdrant
```

### "Model not found"
```bash
# List installed models
ollama list

# Pull missing model
ollama pull nomic-embed-text
ollama pull llama3
```

### Slow Performance
- Use `mistral` instead of `llama3` for faster responses
- Reduce batch size in `llm/ollama/ollama.go` (line 184)
- Use quantized models: `ollama pull llama3:7b-q4_0`

### Out of Memory
- Use smaller model: `mistral` or `llama3:7b-q4_0`
- Reduce `NumPredict` in chat options (default: 512 tokens)
- Close other applications

## Testing the Setup

### Test Ollama Embeddings:
```bash
curl http://localhost:11434/api/embeddings -d '{
  "model": "nomic-embed-text",
  "prompt": "Hello world"
}'
```

Should return a 768-dimensional vector.

### Test Ollama Chat:
```bash
curl http://localhost:11434/api/chat -d '{
  "model": "llama3",
  "messages": [{"role": "user", "content": "Hello!"}],
  "stream": false
}'
```

### Test Qdrant:
```bash
curl http://localhost:6333/collections
```

Should return `{"result":{"collections":[]},"status":"ok","time":...}`

## Performance Benchmarks

Typical performance on modern hardware:

| Task | llama3 (8B) | mistral (7B) | llama3:70b |
|------|-------------|--------------|------------|
| Embedding (per chunk) | ~50ms | ~40ms | N/A |
| Chat response (512 tokens) | ~2-5s | ~1-3s | ~10-30s |
| RAM usage | ~8GB | ~4GB | ~40GB |

## Advantages of Local Setup

✅ **Complete Privacy**: Documents never leave your machine
✅ **No API Costs**: No per-token charges
✅ **Offline Capable**: Works without internet
✅ **Unlimited Usage**: Query as much as you want
✅ **Full Control**: Customize models and parameters
✅ **Data Sovereignty**: Compliance-friendly
✅ **Low Latency**: No network round trips

## Next Steps

1. Upload some test documents (PDFs, text files)
2. Ask questions about your documents
3. Experiment with different models for your use case
4. Adjust parameters for your hardware

Happy local AI querying! 🚀
