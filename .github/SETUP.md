# GitHub Actions Quick Setup Guide

## What's Included

✅ **Automated Testing** - Runs on every push and PR  
✅ **Code Coverage** - Tracks test coverage over time  
✅ **Linting** - Ensures code quality and consistency  
✅ **Build Verification** - Confirms the server compiles successfully  

## Files Created

```
.github/
  workflows/
    ci.yml              # Main CI workflow
  WORKFLOWS.md          # Detailed documentation
.golangci.yml           # Linting configuration
Makefile                # Updated with ci-check target
```

## Quick Start

### 1. Local Testing (Before Push)

Run the same checks that GitHub Actions will run:

```bash
# Quick pre-commit checks
make pre-commit

# Full CI simulation (recommended)
make ci-check
```

### 2. Push to GitHub

Once you push to GitHub, the workflow will automatically run. You can watch it at:
```
https://github.com/adhocteam/recreation-mcp-server/actions
```

### 3. Add Status Badges (Optional)

Add to your main README.md to show CI status:

```markdown
[![CI](https://github.com/{your-org-name}}/recreation-mcp-server/actions/workflows/ci.yml/badge.svg)](https://github.com/{your-org-name}}/recreation-mcp-server/actions/workflows/ci.yml)
```

## What Gets Checked

### Test Job
- ✓ Go module dependencies verified
- ✓ Static analysis with `go vet`
- ✓ All unit tests with race detection
- ✓ Code coverage report generated
- ✓ Coverage uploaded to Codecov (optional)

### Lint Job  
- ✓ golangci-lint v2.5.0 with 10+ linters
- ✓ Code formatting checks (gofmt, goimports)
- ✓ Import statement organization
- ✓ Common Go issues detected (errcheck, govet, staticcheck)
- ✓ Configuration: `.golangci.yml` (v2 format)

### Build Job
- ✓ Server binary compiles successfully
- ✓ Build artifact uploaded (kept for 7 days)

## Current Status

- **Test Coverage:** 25.3%
  - config: 94.3% ✓
  - cache: 49.4%
  - api: 27.0%
  
- **Test Runtime:** ~1 second ✓

- **All Tests:** Passing ✓

## Next Steps

1. **Enable Codecov (Optional)**
   - Sign up at codecov.io
   - Add `CODECOV_TOKEN` to GitHub secrets
   - Get coverage trend graphs and PR comments

2. **Increase Coverage**
   - Add handler tests → ~50-60% coverage
   - Add util tests → ~60-70% coverage
   - Uncomment coverage threshold in workflow

3. **Add More Linters (Optional)**
   - Install golangci-lint locally: `brew install golangci-lint`
   - Run `make lint` to see more detailed checks
   - Configure additional linters in `.golangci.yml`

## Troubleshooting

**Tests pass locally but fail in CI?**
- Check Go version matches (1.21+)
- Run `go mod tidy`
- Test with race detector: `go test -race ./...`

**Want to skip CI for a commit?**
```bash
git commit -m "docs: update README [skip ci]"
```

**Need to debug a failing workflow?**
- Check the Actions tab in GitHub
- Download build artifacts to inspect
- Re-run failed jobs with debug logging enabled

## Integration with Pull Requests

The workflow automatically:
- ✓ Runs on every PR
- ✓ Shows status in PR checks
- ✓ Blocks merge if tests fail (optional)
- ✓ Comments with coverage changes (with Codecov)

To require passing tests before merge:
1. Go to Settings → Branches
2. Add branch protection rule for `main`
3. Require status checks to pass: `Test`, `Lint`, `Build`
