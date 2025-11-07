#!/bin/bash
# 1-migrate-go-imports.sh
# Updates Go module path and all import statements

set -e

echo "=========================================="
echo "  Step 1: Updating Go Imports"
echo "=========================================="
echo ""

# Get the project root (parent of scripts directory)
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$PROJECT_ROOT"

echo "Project root: $PROJECT_ROOT"
echo ""

# Update go.mod
echo "Updating go.mod module path..."
if [ -f "go.mod" ]; then
    sed -i.bak 's|module github.com/pashpashpash/vault|module github.com/itaiguardiola/askara|g' go.mod
    echo "  ✓ go.mod updated"
else
    echo "  ✗ go.mod not found!"
    exit 1
fi

# Update all Go files
echo ""
echo "Updating import statements in Go files..."
GO_FILES=$(find . -type f -name "*.go" -not -path "*/vendor/*" -not -path "*/.git/*" -not -path "*/node_modules/*")
COUNT=0

for file in $GO_FILES; do
    if grep -q "github.com/pashpashpash/vault" "$file"; then
        sed -i.bak 's|github.com/pashpashpash/vault|github.com/itaiguardiola/askara|g' "$file"
        echo "  ✓ Updated: $file"
        ((COUNT++))
    fi
done

echo ""
echo "Updated $COUNT Go files"
echo ""
echo "Step 1 complete. Run 'go mod tidy' to verify."
echo ""
