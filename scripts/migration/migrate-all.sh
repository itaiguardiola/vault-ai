#!/bin/bash
# migrate-all.sh
# Master migration script - runs all migration steps in order

set -e  # Exit on error

echo ""
echo "=========================================="
echo "  Askara Migration Script"
echo "  vault-ai → Askara"
echo "=========================================="
echo ""
echo "This script will rename the entire project from 'vault-ai'/'OP Vault' to 'Askara'"
echo ""
echo "What will be changed:"
echo "  - Go module path: github.com/pashpashpash/vault → github.com/pashpashpash/askara"
echo "  - Docker services: vault-* → askara-*"
echo "  - Binary name: vault-web-server → askara-web-server"
echo "  - All documentation and UI branding"
echo "  - Directory names"
echo ""

# Get the project root
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$PROJECT_ROOT"

echo "Project root: $PROJECT_ROOT"
echo ""

# Confirmation prompt
read -p "Do you want to proceed with the migration? (yes/no): " -r
echo
if [[ ! $REPLY =~ ^[Yy]es$ ]]; then
    echo "Migration cancelled."
    exit 0
fi

# Create backup
echo "Creating backup..."
BACKUP_NAME="vault-ai-backup-$(date +%Y%m%d-%H%M%S).tar.gz"
cd ..
tar -czf "$BACKUP_NAME" --exclude='vault-ai/node_modules' --exclude='vault-ai/.git' vault-ai/
echo "  ✓ Backup created: $(pwd)/$BACKUP_NAME"
cd "$PROJECT_ROOT"
echo ""

# Make all migration scripts executable
chmod +x scripts/migration/*.sh

echo "Starting migration..."
echo ""

# Phase 1: Go imports
./scripts/migration/1-migrate-go-imports.sh
if [ $? -ne 0 ]; then
    echo "Error in step 1. Aborting."
    exit 1
fi

# Phase 2: Docker
./scripts/migration/2-migrate-docker.sh
if [ $? -ne 0 ]; then
    echo "Error in step 2. Aborting."
    exit 1
fi

# Phase 3: Package config
./scripts/migration/3-migrate-package.sh
if [ $? -ne 0 ]; then
    echo "Error in step 3. Aborting."
    exit 1
fi

# Phase 4: Documentation
./scripts/migration/4-migrate-docs.sh
if [ $? -ne 0 ]; then
    echo "Error in step 4. Aborting."
    exit 1
fi

# Phase 5: Frontend
./scripts/migration/5-migrate-frontend.sh
if [ $? -ne 0 ]; then
    echo "Error in step 5. Aborting."
    exit 1
fi

# Phase 6: Rename directories
./scripts/migration/6-migrate-rename-dirs.sh
if [ $? -ne 0 ]; then
    echo "Error in step 6. Aborting."
    exit 1
fi

echo ""
echo "=========================================="
echo "  Migration Complete!"
echo "=========================================="
echo ""
echo "All migration steps completed successfully."
echo ""
echo "Next steps:"
echo ""
echo "1. Review changes:"
echo "   git status"
echo "   git diff"
echo ""
echo "2. Clean up backup files (.bak files created during migration):"
echo "   find . -name '*.bak' -type f -delete"
echo ""
echo "3. Test Go build:"
echo "   go mod tidy"
echo "   go build ./askara-web-server"
echo ""
echo "4. Test Docker build:"
echo "   docker-compose build"
echo ""
echo "5. Test application:"
echo "   docker-compose up"
echo ""
echo "6. Rename main directory (from parent directory):"
echo "   cd .. && mv vault-ai askara && cd askara"
echo ""
echo "7. If everything works, commit changes:"
echo "   git add -A"
echo "   git commit -m 'Migrate project from vault-ai to Askara'"
echo ""
echo "8. Update remote repository:"
echo "   - Rename repository on GitHub/GitLab"
echo "   - Update git remote: git remote set-url origin <new-url>"
echo "   - Push changes: git push -u origin <branch>"
echo ""
echo "Backup location: $(cd .. && pwd)/$BACKUP_NAME"
echo ""
