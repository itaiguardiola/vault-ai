# Migration Scripts

This directory contains automated scripts for migrating the project from "vault-ai" to "Askara".

## Scripts Overview

| Script | Purpose | Changes |
|--------|---------|---------|
| **migrate-all.sh** | Master script - runs all steps | Runs scripts 1-6 in order |
| **1-migrate-go-imports.sh** | Update Go module path | go.mod + 16 Go files |
| **2-migrate-docker.sh** | Update Docker config | docker-compose.yml, Dockerfile |
| **3-migrate-package.sh** | Update Node config | package.json |
| **4-migrate-docs.sh** | Update documentation | 7 markdown files |
| **5-migrate-frontend.sh** | Update UI branding | React components, config |
| **6-migrate-rename-dirs.sh** | Rename directories | vault-web-server → askara-web-server |

## Usage

### Quick Start (Recommended)

Run the master script to execute all migration steps:

```bash
./scripts/migration/migrate-all.sh
```

This will:
1. Create an automatic backup
2. Run all migration steps in order
3. Show detailed progress
4. Provide next steps

### Individual Scripts

You can also run scripts individually for testing:

```bash
# Step 1: Update Go imports
./scripts/migration/1-migrate-go-imports.sh

# Step 2: Update Docker configuration
./scripts/migration/2-migrate-docker.sh

# Step 3: Update package.json
./scripts/migration/3-migrate-package.sh

# Step 4: Update documentation
./scripts/migration/4-migrate-docs.sh

# Step 5: Update frontend branding
./scripts/migration/5-migrate-frontend.sh

# Step 6: Rename directories
./scripts/migration/6-migrate-rename-dirs.sh
```

**Note:** Running individual scripts is useful for:
- Testing specific changes
- Re-running a failed step
- Understanding what each step does

## Script Details

### 1-migrate-go-imports.sh

Updates the Go module path and all import statements.

**Changes:**
- `go.mod`: Module path
- All `.go` files: Import statements

**Affected Files:** ~16 Go files

**After running:** Execute `go mod tidy` to verify

### 2-migrate-docker.sh

Updates Docker and docker-compose configuration.

**Changes:**
- Service names: vault-* → askara-*
- Container names
- Network name: vault-network → askara-network
- Comments and branding

**Affected Files:**
- `docker-compose.yml`
- `Dockerfile`

### 3-migrate-package.sh

Updates Node.js package configuration.

**Changes:**
- Package name
- Description
- Repository URLs
- Script paths

**Affected Files:**
- `package.json`

### 4-migrate-docs.sh

Updates all documentation files.

**Changes:**
- Branding: "OP Vault" / "The Vault" → "Askara"
- Binary names: vault-web-server → askara-web-server
- Repository URLs

**Affected Files:**
- `README.md`
- `SETUP_LOCAL.md`
- `DOCKER_SETUP.md`
- `CHANGES_SUMMARY.md`
- `DOCKER_AND_CONFIG_SUMMARY.md`
- `DOCUMENT_LIBRARY_FEATURE.md`
- `OCR_FEATURE.md`

**Note:** Original references to vault.pash.city are preserved

### 5-migrate-frontend.sh

Updates React components and UI text.

**Changes:**
- All JSX files: "OP Vault" / "The Vault" → "Askara"
- `main.go`: Page titles, favicon path
- `config/websites.json`: Site configuration

**Affected Files:** 5+ React components

### 6-migrate-rename-dirs.sh

Renames project directories.

**Changes:**
- `vault-web-server/` → `askara-web-server/`

**Note:** Main project directory (vault-ai → askara) must be renamed manually

## Backup & Safety

### Automatic Backup

The master script (`migrate-all.sh`) automatically creates a backup before making changes:

```
vault-ai-backup-YYYYMMDD-HHMMSS.tar.gz
```

Located in the parent directory of vault-ai.

### Manual Backup

Each individual script creates `.bak` files for changed files:

```
go.mod.bak
docker-compose.yml.bak
package.json.bak
etc.
```

### Cleanup

After successful migration, remove backup files:

```bash
find . -name '*.bak' -type f -delete
```

## Verification

After running migration scripts:

### 1. Check Git Status

```bash
git status
git diff
```

### 2. Test Go Build

```bash
go mod tidy
go build ./askara-web-server
```

### 3. Test Docker Build

```bash
docker-compose build
docker-compose up
```

### 4. Check Application

Visit http://localhost:8100 and verify:
- Branding shows "Askara"
- All features work
- No "vault" references in UI

## Rollback

If something goes wrong:

```bash
# Stop services
docker-compose down

# Restore from backup
cd ..
rm -rf vault-ai
tar -xzf vault-ai-backup-YYYYMMDD-HHMMSS.tar.gz
cd vault-ai
```

## Requirements

- Bash shell
- sed (GNU sed recommended)
- find
- tar
- chmod

These tools are standard on Linux and macOS.

## Notes

- All scripts use `set -e` to exit on error
- Backup files are created before modifications
- Scripts can be re-run safely (idempotent)
- Each script reports its progress and results

## Full Documentation

For complete migration documentation, see:
- `../../MIGRATION_TO_ASKARA.md` - Complete migration plan
- `../../MIGRATION_QUICK_START.md` - Quick start guide

## Support

If you encounter issues:
1. Check script output for error messages
2. Review backup files (.bak) to see what changed
3. Restore from backup and try again
4. Run individual scripts to isolate the problem
