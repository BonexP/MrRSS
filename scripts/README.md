# Scripts

This directory contains cross-platform automation scripts for development, testing, and release processes.

## Available Scripts

### check.sh / check.ps1

Runs all code quality checks including:

- Go vet (static analysis)
- Go formatting checks
- Frontend linting
- Frontend tests
- Build verification

**Usage:**

```bash
# Linux/macOS
./scripts/check.sh

# Windows
.\scripts\check.ps1
```

### pre-release.sh / pre-release.ps1

Performs pre-release checks including:

- All checks from `check.sh/ps1`
- Go module cleanliness
- Frontend dependency audit
- Version consistency across all files

**Usage:**

```bash
# Linux/macOS
./scripts/pre-release.sh

# Windows
.\scripts\pre-release.ps1
```

### bump-version.sh / bump-version.ps1

Updates version numbers across all relevant files:

- `frontend/package.json`
- `wails.json` (version and productVersion)
- `internal/version/version.go`
- `frontend/src/components/modals/settings/AboutTab.vue`
- README badges

**Usage:**

```bash
# Linux/macOS
./scripts/bump-version.sh 1.2.4

# Windows
.\scripts\bump-version.ps1 -NewVersion 1.2.4
```

## Integration with CI/CD

These scripts are used in GitHub Actions workflows:

- `check.sh/ps1` is called during the test workflow
- `pre-release.sh/ps1` can be used in release workflows

## Adding New Scripts

When adding new scripts:

1. Create both `.sh` (bash) and `.ps1` (PowerShell) versions
2. Add documentation to this README
3. Update the Makefile if needed
4. Test on both Linux and Windows environments
