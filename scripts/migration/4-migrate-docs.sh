#!/bin/bash
# 4-migrate-docs.sh
# Updates all documentation files

set -e

echo "=========================================="
echo "  Step 4: Updating Documentation"
echo "=========================================="
echo ""

# Get the project root
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$PROJECT_ROOT"

echo "Project root: $PROJECT_ROOT"
echo ""

# List of documentation files
DOCS=(
    "README.md"
    "SETUP_LOCAL.md"
    "DOCKER_SETUP.md"
    "CHANGES_SUMMARY.md"
    "DOCKER_AND_CONFIG_SUMMARY.md"
    "DOCUMENT_LIBRARY_FEATURE.md"
    "OCR_FEATURE.md"
)

COUNT=0

for doc in "${DOCS[@]}"; do
    if [ -f "$doc" ]; then
        echo "Updating $doc..."
        cp "$doc" "$doc.bak"

        # Replace branding (but keep original vault.pash.city references)
        sed -i 's/OP Vault - Local Edition/Askara/g' "$doc"
        sed -i 's/OP Vault/Askara/g' "$doc"
        sed -i 's/The Vault/Askara/g' "$doc"

        # Replace paths and binary names
        sed -i 's/vault-web-server/askara-web-server/g' "$doc"

        # Update repository references (but not URLs to vault.pash.city)
        sed -i 's|github.com/pashpashpash/vault\([^.]\)|github.com/pashpashpash/askara\1|g' "$doc"
        sed -i 's|pashpashpash/vault#|pashpashpash/askara#|g' "$doc"
        sed -i 's|pashpashpash/vault/|pashpashpash/askara/|g' "$doc"

        # Update image references
        sed -i 's|vault_library.png|askara_library.png|g' "$doc"
        sed -i 's|vault-favicon.png|askara-favicon.png|g' "$doc"

        echo "  ✓ $doc updated"
        ((COUNT++))
    else
        echo "  ⚠ $doc not found (skipping)"
    fi
done

echo ""
echo "Updated $COUNT documentation files"
echo ""
echo "Note: Original references to vault.pash.city (the original project) have been preserved"
echo ""
echo "Step 4 complete. Documentation updated."
echo ""
