# Migration Plan: vault-ai → Askara

## Overview

This document outlines the complete migration plan to rename the project from **"vault-ai"** / **"OP Vault"** to **"Askara"**.

**Target Name:** Askara
**Current Names:**
- Repository: `vault-ai`
- Go Module: `github.com/pashpashpash/vault`
- Application: "OP Vault" / "The Vault"
- Binary: `vault-web-server`
- Docker Services: `vault-app`, `vault-qdrant`, `vault-ollama`, etc.

---

## Migration Checklist

### Phase 1: Go Module & Imports (16 files)

**Critical:** Go module path changes require updating all import statements.

#### 1.1 Update go.mod
- [ ] Change module path from `github.com/pashpashpash/vault` to `github.com/pashpashpash/askara`

#### 1.2 Update all Go import statements (16 files)
Replace `github.com/pashpashpash/vault` → `github.com/pashpashpash/askara`

- [ ] `vault-web-server/main.go`
- [ ] `vault-web-server/postapi/handlercontext.go`
- [ ] `vault-web-server/postapi/questions.go`
- [ ] `vault-web-server/postapi/fileupload.go`
- [ ] `vault-web-server/postapi/formparseverify.go`
- [ ] `vault-web-server/postapi/openai.go`
- [ ] `vault-web-server/getapi/getapi.go`
- [ ] `vectordb/pinecone/pinecone.go`
- [ ] `vectordb/qdrant/qdrant.go`
- [ ] `vectordb/vectordb.go`
- [ ] `llm/ollama/ollama.go`
- [ ] `llm/interface.go`
- [ ] `validator/validator.go`
- [ ] `validator/email.go`
- [ ] `form/callquestions.go`
- [ ] `form/form.go`

### Phase 2: Directory & File Renaming

#### 2.1 Rename main directory
- [ ] Rename `/home/user/vault-ai` → `/home/user/askara`

#### 2.2 Rename vault-web-server directory
- [ ] Rename `vault-web-server/` → `askara-web-server/`
- [ ] Update all references to this path in documentation and scripts

#### 2.3 Update binary names
- [ ] Change binary output from `vault-web-server` → `askara-web-server`
- [ ] Update in `Dockerfile` (line 36)
- [ ] Update in `Dockerfile` (line 51, 71)
- [ ] Update in `package.json` scripts

### Phase 3: Docker Configuration

#### 3.1 Update docker-compose.yml
- [ ] Service names:
  - `vault-app` → `askara-app`
  - `vault-qdrant` → `askara-qdrant`
  - `vault-ollama` → `askara-ollama`
  - `vault-ollama-setup` → `askara-ollama-setup`
- [ ] Container names (lines 7, 28, 57, 83)
- [ ] Network name: `vault-network` → `askara-network` (lines 23, 42, 78, 103, 113)
- [ ] Comment on line 52: "OP Vault Application" → "Askara Application"

#### 3.2 Update Dockerfile
- [ ] Line 1: Comment "OP Vault Local Edition" → "Askara"
- [ ] Binary paths: `vault-web-server` → `askara-web-server`

### Phase 4: Package Configuration

#### 4.1 Update package.json
- [ ] `name`: "vault-web-server" → "askara-web-server"
- [ ] `description`: "vault core website" → "askara core website"
- [ ] `repository.url`: Update GitHub URL
- [ ] `bugs.url`: Update GitHub URL
- [ ] `homepage`: Update GitHub URL
- [ ] Scripts: Update paths from `vault-web-server` → `askara-web-server`

### Phase 5: Documentation (11 files)

#### 5.1 README.md
- [ ] Title: "OP Vault - Local Edition" → "Askara"
- [ ] All mentions of "OP Vault" → "Askara"
- [ ] All mentions of "The Vault" → "Askara"
- [ ] Original version reference: vault.pash.city
- [ ] Image references: `/static/img/common/vault_library.png`
- [ ] Path references: `vault-web-server/`
- [ ] Binary references: `vault-web-server`

#### 5.2 Other Documentation Files
- [ ] `SETUP_LOCAL.md`
- [ ] `DOCKER_SETUP.md`
- [ ] `CHANGES_SUMMARY.md`
- [ ] `DOCKER_AND_CONFIG_SUMMARY.md`
- [ ] `DOCUMENT_LIBRARY_FEATURE.md`
- [ ] `OCR_FEATURE.md`

### Phase 6: Frontend Branding (5+ files)

#### 6.1 React Components
Search and replace in all JSX files:
- [ ] "OP Vault" → "Askara"
- [ ] "The Vault" → "Askara"
- [ ] "vault" references in UI text

Files to update:
- [ ] `components/Pages/DocumentsPage/index.jsx`
- [ ] `components/Pages/ConfigPage/index.jsx`
- [ ] `components/Pages/LandingPage/index.jsx`
- [ ] `components/Header/index.jsx`
- [ ] `components/Footer/index.jsx`

#### 6.2 Page Titles and Metadata
- [ ] `vault-web-server/main.go` (lines 96-97):
  - PageTitle: "The Vault | OP Question-Answer Stack" → "Askara | AI Question-Answer System"
  - PageIcon: `/img/logos/vault-favicon.png` → `/img/logos/askara-favicon.png`

#### 6.3 Configuration Files
- [ ] `config/websites.json`:
  - All "vault" references
  - URLs, titles, meta tags

### Phase 7: Source Code Comments & Strings

#### 7.1 Go Files
Review and update comments containing:
- [ ] "vault" references
- [ ] "OP Vault" references
- [ ] Product descriptions

#### 7.2 Main Configuration
- [ ] `vault-web-server/main.go`:
  - Line 96: currentSite = "vault" → "askara"
  - Line 80: certFile and keyFile paths (if applicable)
  - Line 81: currentHost = "vault.pash.city" (keep as reference to original)
  - Line 86: GitHub redirect URL

### Phase 8: Assets & Static Files

#### 8.1 Images
- [ ] Rename/replace: `/static/img/logos/vault-favicon.png` → `askara-favicon.png`
- [ ] Update: `/static/img/common/vault_library.png` (if branded)
- [ ] Check for other vault-branded images

#### 8.2 CSS/Less Files
- [ ] Search for "vault" in class names or comments
- [ ] Update any vault-specific styling

### Phase 9: Scripts & Build Files

#### 9.1 Shell Scripts
- [ ] `scripts/go-compile.sh` - Check for "vault" references
- [ ] `scripts/source-me.sh` - Check for "vault" references
- [ ] Any other scripts in `/scripts` directory

### Phase 10: Git & Repository

#### 10.1 Repository Settings
- [ ] Repository name on GitHub/GitLab
- [ ] Repository description
- [ ] Repository topics/tags

#### 10.2 Git Configuration
- [ ] Update `.git/config` remote URLs (if needed)
- [ ] Update any git hooks

#### 10.3 Create Migration Tag
- [ ] Tag current commit: `git tag -a v-pre-askara-migration -m "Last commit before Askara migration"`

---

## Detailed File Changes

### Critical: Go Module Path Migration

**Old:** `github.com/pashpashpash/vault`
**New:** `github.com/pashpashpash/askara`

This affects **every Go file** that imports internal packages:

```go
// Before
import "github.com/pashpashpash/vault/chunk"

// After
import "github.com/pashpashpash/askara/chunk"
```

### Docker Service Naming Convention

**Pattern:** `askara-<service>`

```yaml
# Before
services:
  vault-app:
    container_name: vault-app
  qdrant:
    container_name: vault-qdrant

# After
services:
  askara-app:
    container_name: askara-app
  qdrant:
    container_name: askara-qdrant
```

### Binary Naming

**Old Binary:** `vault-web-server`
**New Binary:** `askara-web-server`

Update in:
- `Dockerfile`: Build output path
- `package.json`: Scripts
- `vault-web-server/main.go`: Any self-references
- Documentation: All command examples

### Frontend Branding Guidelines

**Replacements:**
- "OP Vault" → "Askara"
- "The Vault" → "Askara"
- "vault-ai" → "Askara"

**Keep Original References:**
- Links to original project: vault.pash.city
- Git history references
- Attribution comments

---

## Migration Scripts

### Script 1: Update Go Imports

```bash
#!/bin/bash
# migrate-go-imports.sh

echo "Updating Go module path and imports..."

# Update go.mod
sed -i 's|module github.com/pashpashpash/vault|module github.com/pashpashpash/askara|g' go.mod

# Update all Go files
find . -type f -name "*.go" -not -path "*/vendor/*" -not -path "*/.git/*" -exec \
  sed -i 's|github.com/pashpashpash/vault|github.com/pashpashpash/askara|g' {} +

echo "Go imports updated. Run 'go mod tidy' to verify."
```

### Script 2: Update Docker Configuration

```bash
#!/bin/bash
# migrate-docker.sh

echo "Updating Docker configuration..."

# Update docker-compose.yml
sed -i 's/vault-app/askara-app/g' docker-compose.yml
sed -i 's/vault-qdrant/askara-qdrant/g' docker-compose.yml
sed -i 's/vault-ollama/askara-ollama/g' docker-compose.yml
sed -i 's/vault-network/askara-network/g' docker-compose.yml
sed -i 's/OP Vault Application/Askara Application/g' docker-compose.yml

# Update Dockerfile
sed -i 's/OP Vault Local Edition/Askara/g' Dockerfile
sed -i 's/vault-web-server/askara-web-server/g' Dockerfile

echo "Docker configuration updated."
```

### Script 3: Update Package Configuration

```bash
#!/bin/bash
# migrate-package.sh

echo "Updating package.json..."

sed -i 's/"vault-web-server"/"askara-web-server"/g' package.json
sed -i 's/vault core website/askara core website/g' package.json
sed -i 's|pashpashpash/vault|pashpashpash/askara|g' package.json

echo "package.json updated."
```

### Script 4: Update Documentation

```bash
#!/bin/bash
# migrate-docs.sh

echo "Updating documentation files..."

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

for doc in "${DOCS[@]}"; do
  if [ -f "$doc" ]; then
    echo "Updating $doc..."
    # Replace branding
    sed -i 's/OP Vault/Askara/g' "$doc"
    sed -i 's/The Vault/Askara/g' "$doc"
    # Replace paths
    sed -i 's/vault-web-server/askara-web-server/g' "$doc"
    # Update repository references (but keep original vault.pash.city links)
    sed -i 's|pashpashpash/vault[^.]|pashpashpash/askara|g' "$doc"
  fi
done

echo "Documentation updated."
```

### Script 5: Update Frontend

```bash
#!/bin/bash
# migrate-frontend.sh

echo "Updating frontend components..."

# Update all JSX files
find components -type f -name "*.jsx" -exec \
  sed -i 's/OP Vault/Askara/g; s/The Vault/Askara/g' {} +

# Update main.go page configuration
sed -i 's/"The Vault | OP Question-Answer Stack"/"Askara | AI Question-Answer System"/g' vault-web-server/main.go
sed -i 's|vault-favicon.png|askara-favicon.png|g' vault-web-server/main.go
sed -i 's/currentSite = "vault"/currentSite = "askara"/g' vault-web-server/main.go

echo "Frontend branding updated."
```

### Script 6: Rename Directories

```bash
#!/bin/bash
# migrate-rename-dirs.sh

echo "Renaming directories..."

# Must be run from parent directory of vault-ai
if [ -d "vault-ai" ]; then
  mv vault-ai askara
  echo "Renamed vault-ai → askara"
fi

cd askara

# Rename vault-web-server directory
if [ -d "vault-web-server" ]; then
  mv vault-web-server askara-web-server
  echo "Renamed vault-web-server → askara-web-server"

  # Update package.json references
  sed -i 's|./vault-web-server|./askara-web-server|g' package.json
  sed -i 's|vault-web-server/|askara-web-server/|g' package.json
fi

echo "Directory renaming complete."
```

### Master Migration Script

```bash
#!/bin/bash
# migrate-all.sh
# Master script to run all migration steps

set -e  # Exit on error

echo "=========================================="
echo "  Askara Migration Script"
echo "  vault-ai → Askara"
echo "=========================================="
echo ""

# Create backup
echo "Creating backup..."
cd ..
tar -czf vault-ai-backup-$(date +%Y%m%d-%H%M%S).tar.gz vault-ai/
cd vault-ai

echo "Starting migration..."
echo ""

# Phase 1: Go imports
echo "[1/6] Updating Go module and imports..."
./scripts/migrate-go-imports.sh

# Phase 2: Docker
echo "[2/6] Updating Docker configuration..."
./scripts/migrate-docker.sh

# Phase 3: Package config
echo "[3/6] Updating package.json..."
./scripts/migrate-package.sh

# Phase 4: Documentation
echo "[4/6] Updating documentation..."
./scripts/migrate-docs.sh

# Phase 5: Frontend
echo "[5/6] Updating frontend branding..."
./scripts/migrate-frontend.sh

# Phase 6: Rename directories (must be last)
echo "[6/6] Renaming directories..."
cd ..
./vault-ai/scripts/migrate-rename-dirs.sh

echo ""
echo "=========================================="
echo "  Migration Complete!"
echo "=========================================="
echo ""
echo "Next steps:"
echo "  1. Review changes: cd askara && git diff"
echo "  2. Test build: go mod tidy && go build ./askara-web-server"
echo "  3. Test Docker: docker-compose build"
echo "  4. Run application: docker-compose up"
echo "  5. If successful, commit: git add -A && git commit -m 'Migrate to Askara'"
echo ""
```

---

## Post-Migration Verification

### 1. Go Build Test
```bash
cd askara
go mod tidy
go build ./askara-web-server
```

### 2. Docker Build Test
```bash
docker-compose build
```

### 3. Docker Run Test
```bash
docker-compose up -d
```

### 4. Check Services
```bash
# Application
curl http://localhost:8100

# Qdrant
curl http://localhost:6333

# Ollama
curl http://localhost:11434/api/tags
```

### 5. Frontend Check
Visit http://localhost:8100 and verify:
- [ ] Page title shows "Askara"
- [ ] No "OP Vault" or "The Vault" references
- [ ] All pages load correctly
- [ ] Settings page works
- [ ] Documents page works

### 6. Functionality Test
- [ ] Upload a test document
- [ ] Ask a question
- [ ] Verify answer retrieval
- [ ] Check document management
- [ ] Test configuration changes

---

## Rollback Plan

If migration fails:

```bash
# Stop services
docker-compose down

# Restore from backup
cd ..
rm -rf askara
tar -xzf vault-ai-backup-YYYYMMDD-HHMMSS.tar.gz

# Resume with old version
cd vault-ai
docker-compose up
```

---

## Breaking Changes for Users

### For Developers

1. **Go Import Paths Changed**
   - Old: `github.com/pashpashpash/vault/*`
   - New: `github.com/pashpashpash/askara/*`
   - Action: Update imports in custom code

2. **Binary Name Changed**
   - Old: `vault-web-server`
   - New: `askara-web-server`
   - Action: Update scripts and commands

3. **Directory Structure Changed**
   - Old: `vault-ai/vault-web-server/`
   - New: `askara/askara-web-server/`
   - Action: Update paths in scripts

### For Docker Users

1. **Container Names Changed**
   - Old: `vault-app`, `vault-qdrant`, `vault-ollama`
   - New: `askara-app`, `askara-qdrant`, `askara-ollama`
   - Action: Update any scripts referencing containers

2. **Network Name Changed**
   - Old: `vault-network`
   - New: `askara-network`
   - Action: Update external network configurations

3. **Data Persistence**
   - Volume names remain the same
   - No data loss expected
   - Backup recommended before migration

---

## Timeline

**Estimated Duration:** 2-3 hours

1. **Preparation** (30 min)
   - Backup current state
   - Review checklist
   - Test scripts

2. **Execution** (60 min)
   - Run migration scripts
   - Rename directories
   - Update assets

3. **Verification** (60 min)
   - Build tests
   - Docker tests
   - Functionality tests
   - Documentation review

---

## Success Criteria

Migration is successful when:

- ✅ All Go files compile without errors
- ✅ Docker containers build successfully
- ✅ Application starts and serves on port 8100
- ✅ All features work (upload, search, config, documents)
- ✅ No "vault" references in user-facing UI
- ✅ Documentation is accurate and consistent
- ✅ No broken links or missing assets

---

## Additional Considerations

### GitHub Repository

After migration:
1. Rename repository on GitHub: Settings → General → Repository name
2. Update repository description
3. Update About section
4. Update topics/tags
5. Consider archiving old "vault-ai" repository with redirect

### Domain & Deployment

If deploying to production:
- Update domain configuration
- Update SSL certificates
- Update DNS records
- Update environment variables on hosting platform

### Dependencies

No external dependency changes required. The project will continue to use:
- Ollama (same endpoints)
- Qdrant (same endpoints)
- All Go packages (no changes)
- All NPM packages (no changes)

---

## Contact & Support

For questions about this migration:
- Review this document thoroughly
- Check migration scripts for comments
- Test in development environment first
- Keep backups until migration is verified

**Remember:** This is a cosmetic/branding migration. Core functionality remains unchanged.
