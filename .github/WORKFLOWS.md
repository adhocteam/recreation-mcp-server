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
   - Runs golangci-lint with multiple linters
   - Checks for common issues and code smells
   - Ensures consistent code style

3. **Build** - Verifies the project builds
   - Compiles the MCP server binary
   - Uploads build artifact for 7 days
   - Ensures no build errors

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
[![codecov](https://codecov.io/gh/adhocteam/recreation-mcp-server/branch/main/graph/badge.svg)](https://codecov.io/gh/adhocteam/recreation-mcp-server)
[![Go Report Card](https://goreportcard.com/badge/github.com/adhocteam/recreation-mcp-server)](https://goreportcard.com/report/github.com/adhocteam/recreation-mcp-server)
```

## Troubleshooting

### Tests fail locally but pass in CI (or vice versa)

- Ensure you're using the same Go version (1.21+)
- Run `go mod tidy` to sync dependencies
- Check for environment-specific issues

### Linting failures

- Run `make lint` locally to see the same errors
- Install golangci-lint: `brew install golangci-lint` (macOS)
- Configure your editor to use `.golangci.yml`

### Build failures

- Verify all dependencies are available: `go mod verify`
- Ensure no uncommitted changes to `go.mod` or `go.sum`
- Check for syntax errors: `go build ./...`
