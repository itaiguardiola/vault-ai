# Askara Migration - Quick Start Guide

## TL;DR - Fast Migration

```bash
# 1. Run the master migration script
./scripts/migration/migrate-all.sh

# 2. Clean up backup files
find . -name '*.bak' -type f -delete

# 3. Test the build
go mod tidy
go build ./askara-web-server

# 4. Rename the main directory (from parent)
cd .. && mv vault-ai askara && cd askara

# 5. Commit changes
git add -A
git commit -m "Migrate project from vault-ai to Askara"
```

## What Gets Changed

| Category | Old | New |
|----------|-----|-----|
| **Go Module** | `github.com/pashpashpash/vault` | `github.com/pashpashpash/askara` |
| **Binary** | `vault-web-server` | `askara-web-server` |
| **Docker Service** | `vault-app` | `askara-app` |
| **Docker Network** | `vault-network` | `askara-network` |
| **Branding** | "OP Vault" / "The Vault" | "Askara" |
| **Directory** | `vault-ai/vault-web-server/` | `askara/askara-web-server/` |

## Pre-Migration Checklist

- [ ] Commit or stash any uncommitted changes
- [ ] Stop running Docker containers: `docker-compose down`
- [ ] Note: Backup is created automatically by the script
- [ ] Estimated time: 30-60 minutes

## Migration Steps

### Step 1: Run Migration Script

```bash
cd /path/to/vault-ai
./scripts/migration/migrate-all.sh
```

The script will:
1. Create a timestamped backup
2. Update Go module and imports (16 files)
3. Update Docker configuration
4. Update package.json
5. Update documentation (7 files)
6. Update frontend branding (5+ files)
7. Rename `vault-web-server` → `askara-web-server`

### Step 2: Review Changes

```bash
# See what changed
git status

# Review specific changes
git diff

# Check for any remaining "vault" references
grep -r "vault" --include="*.go" --include="*.jsx" --include="*.md" . | grep -v ".git" | grep -v "node_modules" | grep -v "vault.pash.city"
```

### Step 3: Test Build

```bash
# Tidy Go modules
go mod tidy

# Build the binary
go build ./askara-web-server

# Should create: ./bin/askara-web-server
```

### Step 4: Test Docker

```bash
# Build containers
docker-compose build

# Start services
docker-compose up -d

# Check logs
docker-compose logs -f askara-app

# Test endpoint
curl http://localhost:8100
```

### Step 5: Verify Functionality

Visit http://localhost:8100 and check:
- [ ] Page title shows "Askara"
- [ ] No "OP Vault" or "vault" in UI
- [ ] Can upload a document
- [ ] Can ask a question
- [ ] Settings page works
- [ ] Documents page works

### Step 6: Rename Main Directory

```bash
# From parent directory
cd /path/to
mv vault-ai askara
cd askara

# Verify it worked
pwd  # Should show /path/to/askara
```

### Step 7: Commit Changes

```bash
# Add all changes
git add -A

# Commit
git commit -m "Migrate project from vault-ai to Askara

- Update Go module path to github.com/pashpashpash/askara
- Rename binary to askara-web-server
- Update Docker services and network names
- Update all documentation and branding
- Rename directories for consistency"

# Optional: Tag the migration
git tag -a v1.0.0-askara -m "First release as Askara"
```

### Step 8: Update Remote Repository

```bash
# Rename repository on GitHub/GitLab first, then:

# Update remote URL
git remote set-url origin https://github.com/YOUR_USERNAME/askara.git

# Push changes
git push -u origin main

# Push tags
git push --tags
```

## Verification Checklist

After migration, verify:

- [ ] ✅ Go build succeeds without errors
- [ ] ✅ Docker build succeeds
- [ ] ✅ All containers start (askara-app, askara-qdrant, askara-ollama)
- [ ] ✅ Web interface loads at http://localhost:8100
- [ ] ✅ Page title and branding show "Askara"
- [ ] ✅ Can upload documents successfully
- [ ] ✅ Can query documents successfully
- [ ] ✅ Configuration page works
- [ ] ✅ Documents page works
- [ ] ✅ No console errors in browser
- [ ] ✅ No "vault" references in UI (except vault.pash.city links to original)

## Troubleshooting

### Build Errors

```bash
# If Go build fails:
go mod tidy
go clean -cache
go build ./askara-web-server
```

### Docker Errors

```bash
# If containers fail to start:
docker-compose down
docker-compose build --no-cache
docker-compose up
```

### Partial Migration

If you need to restart:

```bash
# Restore from backup
cd /path/to
rm -rf vault-ai  # or askara
tar -xzf vault-ai-backup-YYYYMMDD-HHMMSS.tar.gz
cd vault-ai

# Run migration again
./scripts/migration/migrate-all.sh
```

## Rollback

If something goes wrong:

```bash
# Stop services
docker-compose down

# Restore from backup
cd /path/to
rm -rf askara  # or vault-ai if not renamed yet
tar -xzf vault-ai-backup-YYYYMMDD-HHMMSS.tar.gz
cd vault-ai

# Resume with original version
docker-compose up
```

## Manual Steps (Optional)

### Update Favicon and Images

```bash
# Rename or replace image files
cd static/img/logos/
mv vault-favicon.png askara-favicon.png
# (or create new askara-favicon.png)

cd ../common/
mv vault_library.png askara_library.png
# (or create new askara_library.png)
```

### Clean Up Backup Files

```bash
# Remove .bak files created during migration
find . -name '*.bak' -type f -delete
```

### Update GitHub Repository

1. Go to repository Settings
2. General → Repository name
3. Rename to "askara"
4. Update description: "Askara - AI-powered document Q&A system with local LLM support"
5. Update topics/tags

## Need Help?

- Full documentation: See `MIGRATION_TO_ASKARA.md`
- Individual scripts: See `scripts/migration/` directory
- Each script can be run independently for testing

## Summary

The migration script automates ~95% of the renaming process. The main manual steps are:
1. Running the script
2. Testing the build
3. Renaming the main directory
4. Updating the remote repository

Total time: **30-60 minutes** (including testing)

**Backup is created automatically**, so you can always roll back if needed.
