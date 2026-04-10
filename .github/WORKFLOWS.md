# GitHub Actions CI/CD

This project uses GitHub Actions for continuous integration and deployment.

## Workflows

### CI Workflow (`.github/workflows/ci.yml`)

Runs automatically on:
- Push to `main` or `develop` branches
- Pull requests targeting `main` or `develop`

**Jobs:**

1. **Test** - Runs all unit tests
   - Verifies Go module dependencies
   - Runs `go vet` for static analysis
   - Executes tests with race detection
   - Generates code coverage report
   - Uploads coverage to Codecov (optional)
   - Current coverage: ~25%

2. **Lint** - Code quality checks
   - Runs golangci-lint v2.5.0 with multiple linters
   - Checks for common issues and code smells
   - Ensures consistent code style
   - Uses configuration from `.golangci.yml`

3. **Build** - Verifies the project builds
   - Compiles the MCP server binary
   - Uploads build artifact for 7 days
   - Ensures no build errors

### Security Workflow (`.github/workflows/security.yml`)

Runs automatically on:
- Push to `main` or `develop` branches
- Pull requests targeting `main` or `develop`
- Weekly schedule (Monday at 09:00 UTC)
- Manual dispatch from GitHub Actions UI

**Checks:**

1. **Module Verification**
   - Runs `go mod verify` to validate module integrity

2. **Vulnerability Scan**
   - Installs and runs `govulncheck ./...`
   - Fails if reachable vulnerabilities are detected

3. **Dependency Update Visibility**
   - Runs `go list -m -u all`
   - Shows available module updates in workflow logs

## Running Locally

You can run the same checks locally before committing:

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run linting
make lint

# Run all pre-commit checks
make pre-commit

# Run vulnerability scan
govulncheck ./...

# Show available module updates
go list -m -u all
```

## Code Coverage

The workflow tracks code coverage and can optionally enforce minimum coverage thresholds. To enable coverage enforcement, uncomment the coverage check in the workflow:

```yaml
# Uncomment to fail CI if coverage is below 70%
if (( $(echo "$COVERAGE < 70.0" | bc -l) )); then
  echo "Coverage ${COVERAGE}% is below minimum 70%"
  exit 1
fi
```

## Codecov Integration (Optional)

To enable Codecov integration:

1. Sign up at [codecov.io](https://codecov.io)
2. Add your repository
3. Add the `CODECOV_TOKEN` secret to your GitHub repository
4. The workflow will automatically upload coverage reports

## Status Badges

Add these badges to your main README.md:

```markdown
[![CI](https://github.com/adhocteam/recreation-mcp-server/actions/workflows/ci.yml/badge.svg)](https://github.com/adhocteam/recreation-mcp-server/actions/workflows/ci.yml)
[![Security](https://github.com/adhocteam/recreation-mcp-server/actions/workflows/security.yml/badge.svg)](https://github.com/adhocteam/recreation-mcp-server/actions/workflows/security.yml)
[![codecov](https://codecov.io/gh/adhocteam/recreation-mcp-server/branch/main/graph/badge.svg)](https://codecov.io/gh/adhocteam/recreation-mcp-server)
[![Go Report Card](https://goreportcard.com/badge/github.com/adhocteam/recreation-mcp-server)](https://goreportcard.com/report/github.com/adhocteam/recreation-mcp-server)
```

## Troubleshooting

### Tests fail locally but pass in CI (or vice versa)

- Ensure you're using the same Go version (1.25+)
- Run `go mod tidy` to sync dependencies
- Check for environment-specific issues

### Linting failures

- Run `make lint` locally to see the same errors
- Install golangci-lint v2.5.0: `brew install golangci-lint` (macOS)
- Ensure you're using golangci-lint v2.x (check with `golangci-lint version`)
- Configure your editor to use `.golangci.yml`

### Build failures

- Verify all dependencies are available: `go mod verify`
- Ensure no uncommitted changes to `go.mod` or `go.sum`
- Check for syntax errors: `go build ./...`
