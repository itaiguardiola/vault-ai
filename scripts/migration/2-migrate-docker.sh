#!/bin/bash
# 2-migrate-docker.sh
# Updates Docker configuration files

set -e

echo "=========================================="
echo "  Step 2: Updating Docker Configuration"
echo "=========================================="
echo ""

# Get the project root
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$PROJECT_ROOT"

echo "Project root: $PROJECT_ROOT"
echo ""

# Update docker-compose.yml
echo "Updating docker-compose.yml..."
if [ -f "docker-compose.yml" ]; then
    cp docker-compose.yml docker-compose.yml.bak

    # Service and container names
    sed -i 's/vault-app/askara-app/g' docker-compose.yml
    sed -i 's/vault-qdrant/askara-qdrant/g' docker-compose.yml
    sed -i 's/vault-ollama/askara-ollama/g' docker-compose.yml

    # Network name
    sed -i 's/vault-network/askara-network/g' docker-compose.yml

    # Comments
    sed -i 's/OP Vault Application/Askara Application/g' docker-compose.yml

    echo "  ✓ docker-compose.yml updated"
else
    echo "  ✗ docker-compose.yml not found!"
    exit 1
fi

# Update Dockerfile
echo ""
echo "Updating Dockerfile..."
if [ -f "Dockerfile" ]; then
    cp Dockerfile Dockerfile.bak

    # Comments and branding
    sed -i 's/OP Vault Local Edition/Askara/g' Dockerfile
    sed -i 's/Multi-stage build for OP Vault/Multi-stage build for Askara/g' Dockerfile

    # Binary names
    sed -i 's/vault-web-server/askara-web-server/g' Dockerfile

    echo "  ✓ Dockerfile updated"
else
    echo "  ✗ Dockerfile not found!"
    exit 1
fi

echo ""
echo "Step 2 complete. Docker configuration updated."
echo ""
