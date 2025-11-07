#!/bin/bash
# 3-migrate-package.sh
# Updates package.json configuration

set -e

echo "=========================================="
echo "  Step 3: Updating package.json"
echo "=========================================="
echo ""

# Get the project root
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$PROJECT_ROOT"

echo "Project root: $PROJECT_ROOT"
echo ""

echo "Updating package.json..."
if [ -f "package.json" ]; then
    cp package.json package.json.bak

    # Package name
    sed -i 's/"vault-web-server"/"askara-web-server"/g' package.json

    # Description
    sed -i 's/vault core website/askara core website/g' package.json

    # Repository URLs
    sed -i 's|pashpashpash/vault|pashpashpash/askara|g' package.json

    # Script paths
    sed -i 's|./vault-web-server|./askara-web-server|g' package.json
    sed -i 's|vault-web-server/|askara-web-server/|g' package.json
    sed -i 's|bin/vault-web-server|bin/askara-web-server|g' package.json

    echo "  ✓ package.json updated"
else
    echo "  ✗ package.json not found!"
    exit 1
fi

echo ""
echo "Step 3 complete. Package configuration updated."
echo ""
