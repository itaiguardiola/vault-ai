#!/bin/bash
# 5-migrate-frontend.sh
# Updates frontend branding and UI text

set -e

echo "=========================================="
echo "  Step 5: Updating Frontend Branding"
echo "=========================================="
echo ""

# Get the project root
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$PROJECT_ROOT"

echo "Project root: $PROJECT_ROOT"
echo ""

# Update all JSX files in components
echo "Updating React components..."
if [ -d "components" ]; then
    JSX_FILES=$(find components -type f -name "*.jsx")
    COUNT=0

    for file in $JSX_FILES; do
        if grep -q -E "(OP Vault|The Vault)" "$file"; then
            cp "$file" "$file.bak"
            sed -i 's/OP Vault/Askara/g' "$file"
            sed -i 's/The Vault/Askara/g' "$file"
            echo "  ✓ Updated: $file"
            ((COUNT++))
        fi
    done

    echo ""
    echo "Updated $COUNT React component files"
else
    echo "  ⚠ components directory not found"
fi

# Update main.go page configuration
echo ""
echo "Updating server configuration..."
if [ -f "vault-web-server/main.go" ]; then
    cp vault-web-server/main.go vault-web-server/main.go.bak

    # Page titles
    sed -i 's/"The Vault | OP Question-Answer Stack"/"Askara | AI Question-Answer System"/g' vault-web-server/main.go

    # Favicon
    sed -i 's|vault-favicon.png|askara-favicon.png|g' vault-web-server/main.go

    # Site name
    sed -i 's/currentSite = "vault"/currentSite = "askara"/g' vault-web-server/main.go

    echo "  ✓ vault-web-server/main.go updated"
else
    echo "  ✗ vault-web-server/main.go not found!"
    exit 1
fi

# Update config/websites.json if it exists
echo ""
echo "Updating website configuration..."
if [ -f "config/websites.json" ]; then
    cp config/websites.json config/websites.json.bak

    sed -i 's/"The Vault"/"Askara"/g' config/websites.json
    sed -i 's|vault-favicon.png|askara-favicon.png|g' config/websites.json

    echo "  ✓ config/websites.json updated"
else
    echo "  ⚠ config/websites.json not found (skipping)"
fi

echo ""
echo "Step 5 complete. Frontend branding updated."
echo ""
echo "Note: You may need to rename/replace favicon and logo image files manually:"
echo "  - /static/img/logos/vault-favicon.png → askara-favicon.png"
echo "  - /static/img/common/vault_library.png → askara_library.png"
echo ""
