#!/bin/bash
# 6-migrate-rename-dirs.sh
# Renames directories (should be run last)

set -e

echo "=========================================="
echo "  Step 6: Renaming Directories"
echo "=========================================="
echo ""

# This script should be run from the parent directory of vault-ai
CURRENT_DIR="$(basename "$(pwd)")"

if [ "$CURRENT_DIR" != "vault-ai" ]; then
    echo "Error: This script must be run from within the vault-ai directory"
    echo "Current directory: $CURRENT_DIR"
    exit 1
fi

echo "Current directory: $(pwd)"
echo ""

# Rename vault-web-server directory
echo "Renaming vault-web-server → askara-web-server..."
if [ -d "vault-web-server" ]; then
    mv vault-web-server askara-web-server
    echo "  ✓ vault-web-server renamed to askara-web-server"

    # Update any remaining references in package.json
    if [ -f "package.json" ]; then
        sed -i 's|vault-web-server|askara-web-server|g' package.json
    fi

    # Update references in go.mod if any
    if [ -f "go.mod" ]; then
        sed -i 's|vault-web-server|askara-web-server|g' go.mod
    fi
else
    echo "  ⚠ vault-web-server directory not found (may have been renamed already)"
fi

echo ""
echo "Step 6 complete. Directories renamed."
echo ""
echo "IMPORTANT: To rename the main project directory from 'vault-ai' to 'askara',"
echo "you must run the following command from the PARENT directory:"
echo ""
echo "  cd .. && mv vault-ai askara && cd askara"
echo ""
