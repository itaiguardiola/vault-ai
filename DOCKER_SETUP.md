# Docker Setup Guide - Askara Local Edition

This guide explains how to run Askara completely locally using Docker. Everything runs in containers - no manual installation required!

## 🚀 Quick Start (One Command!)

```bash
docker-compose up
```

That's it! The first time will take a few minutes to:
1. Pull Docker images (Qdrant, Ollama, build the app)
2. Download required models (nomic-embed-text, llama3)
3. Start all services

**Access the application**: http://localhost:8100

## 📋 Prerequisites

- Docker (version 20.10+)
- Docker Compose (version 2.0+)
- 8GB+ RAM (16GB recommended for llama3)
- 10GB+ free disk space (for models)

### Install Docker

**Linux:**
```bash
curl -fsSL https://get.docker.com | sh
```

**Mac:**
```bash
brew install --cask docker
```

**Windows:**
Download from https://docs.docker.com/desktop/install/windows-install/

## 🏗️ Architecture

The docker-compose setup includes 4 services:

```
┌─────────────────┐
│   vault-app     │  ← Main application (port 8100)
│   (Go + React)  │
└────────┬────────┘
         │
    ┌────┴────┬──────────┐
    │         │          │
┌───▼───┐ ┌──▼────┐ ┌───▼────────┐
│Qdrant │ │Ollama │ │ollama-setup│
│:6333  │ │:11434 │ │(one-time)  │
└───────┘ └───────┘ └────────────┘
```

### Services

1. **vault-app**: The main Askara application
   - Ports: 8100 (web interface)
   - Depends on: Qdrant, Ollama

2. **qdrant**: Local vector database
   - Ports: 6333 (HTTP), 6334 (gRPC)
   - Storage: Persistent volume `qdrant_storage`

3. **ollama**: Local LLM service
   - Ports: 11434 (HTTP API)
   - Storage: Persistent volume `ollama_models`

4. **ollama-setup**: Model downloader (runs once)
   - Downloads required models automatically
   - Exits after completion

## ⚙️ Configuration

### Default Configuration

The application is pre-configured for Docker with these defaults:

```bash
OLLAMA_ENDPOINT=http://ollama:11434
OLLAMA_EMBED_MODEL=nomic-embed-text
OLLAMA_CHAT_MODEL=llama3
QDRANT_API_ENDPOINT=http://qdrant:6333
```

### Custom Configuration

You can change settings in two ways:

#### 1. Through the Web UI (Recommended)

1. Navigate to http://localhost:8100/config
2. Modify settings in the configuration page
3. Test connections
4. Save and restart

#### 2. Environment Variables

Create a `.env` file:

```bash
# .env file
OLLAMA_ENDPOINT=http://ollama:11434
OLLAMA_EMBED_MODEL=nomic-embed-text
OLLAMA_CHAT_MODEL=mistral  # or llama3, llama3:70b
QDRANT_API_ENDPOINT=http://qdrant:6333
CONFIG_PATH=/app/data/config.json
```

Then run:
```bash
docker-compose --env-file .env up
```

## 🎯 Common Commands

### Start Services

```bash
# Start all services (detached mode)
docker-compose up -d

# Start and watch logs
docker-compose up

# Start specific service
docker-compose up vault-app
```

### Stop Services

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (DELETES ALL DATA!)
docker-compose down -v
```

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f vault-app
docker-compose logs -f ollama
docker-compose logs -f qdrant
```

### Restart Services

```bash
# Restart all
docker-compose restart

# Restart specific service
docker-compose restart vault-app
```

### Rebuild After Code Changes

```bash
# Rebuild and restart
docker-compose up --build

# Force rebuild (no cache)
docker-compose build --no-cache
docker-compose up
```

## 📦 Managing Ollama Models

### List Downloaded Models

```bash
docker-compose exec ollama ollama list
```

### Pull Additional Models

```bash
# Mistral (lighter, faster)
docker-compose exec ollama ollama pull mistral

# Llama3 70B (better quality, requires 40GB+ RAM)
docker-compose exec ollama ollama pull llama3:70b

# Other embedding models
docker-compose exec ollama ollama pull bge-large
```

### Remove Models

```bash
docker-compose exec ollama ollama rm mistral
```

## 🔍 Troubleshooting

### Service Won't Start

**Check logs:**
```bash
docker-compose logs vault-app
docker-compose logs ollama
docker-compose logs qdrant
```

**Common issues:**
- Port already in use: Change ports in docker-compose.yml
- Out of memory: Reduce model size or increase Docker memory limit
- Models not downloaded: Check ollama-setup logs

### Cannot Connect to Ollama

```bash
# Check if Ollama is running
docker-compose ps ollama

# Test Ollama directly
curl http://localhost:11434/api/tags

# Restart Ollama
docker-compose restart ollama
```

### Cannot Connect to Qdrant

```bash
# Check if Qdrant is running
docker-compose ps qdrant

# Test Qdrant directly
curl http://localhost:6333/collections

# Restart Qdrant
docker-compose restart qdrant
```

### Models Not Downloading

```bash
# Check setup logs
docker-compose logs ollama-setup

# Manually trigger model download
docker-compose exec ollama ollama pull nomic-embed-text
docker-compose exec ollama ollama pull llama3
```

### Out of Disk Space

```bash
# Check Docker disk usage
docker system df

# Clean up unused images/containers
docker system prune -a

# Remove old volumes (WARNING: deletes data!)
docker volume prune
```

### Application Crashes or Errors

```bash
# Check application logs
docker-compose logs -f vault-app

# Restart with fresh build
docker-compose down
docker-compose up --build
```

## 🎛️ Advanced Configuration

### GPU Support (NVIDIA)

Uncomment the GPU configuration in `docker-compose.yml`:

```yaml
ollama:
  # ...
  deploy:
    resources:
      reservations:
        devices:
          - driver: nvidia
            count: 1
            capabilities: [gpu]
```

### Custom Models

Edit `docker-compose.yml` to pull different models:

```yaml
ollama-setup:
  command: >
    -c "
    echo 'Pulling custom models...';
    ollama pull mistral;
    ollama pull codellama;
    ollama pull nomic-embed-text;
    "
```

### Port Conflicts

Change ports in `docker-compose.yml`:

```yaml
services:
  vault-app:
    ports:
      - "8200:8100"  # Change 8200 to your preferred port

  qdrant:
    ports:
      - "6444:6333"  # Change 6444 to your preferred port
```

### Persistent Configuration

Mount a config file:

```yaml
vault-app:
  volumes:
    - ./config.json:/app/config.json
```

## 🔐 Security Notes

- All services run on localhost by default
- No data is sent to external services
- Volumes persist data between restarts
- To expose to network, modify `docker-compose.yml` ports

## 📊 Resource Requirements

### Minimal Setup (mistral)
- RAM: 8GB
- Disk: 10GB
- CPU: 4 cores

### Recommended Setup (llama3)
- RAM: 16GB
- Disk: 20GB
- CPU: 8 cores

### High-Quality Setup (llama3:70b)
- RAM: 64GB
- Disk: 40GB
- CPU: 16 cores

## 🚀 Performance Tips

1. **Use SSD for Docker volumes**: Much faster than HDD
2. **Increase Docker memory**: Docker Desktop → Settings → Resources
3. **Use smaller models for faster responses**: mistral instead of llama3
4. **Enable GPU if available**: 5-10x faster inference
5. **Close other applications**: Free up RAM for models

## 📝 Docker Compose Commands Cheat Sheet

```bash
# Start
docker-compose up -d

# Stop
docker-compose down

# Rebuild
docker-compose up --build

# View logs
docker-compose logs -f

# Execute command in container
docker-compose exec vault-app /bin/sh
docker-compose exec ollama ollama list

# Check status
docker-compose ps

# Remove everything
docker-compose down -v
```

## 🎉 Next Steps

1. **Visit** http://localhost:8100
2. **Configure** settings at http://localhost:8100/config
3. **Upload** documents
4. **Ask** questions!

## 🆘 Need Help?

- Check logs: `docker-compose logs -f`
- Visit config page: http://localhost:8100/config
- Test connections in the config page
- Verify models: `docker-compose exec ollama ollama list`

Enjoy your fully local, privacy-focused document Q&A system! 🎊
