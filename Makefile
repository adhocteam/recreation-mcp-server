.PHONY: help build test test-unit test-coverage test-integration clean run lint

# Default target
help:
	@echo "Recreation MCP Server - Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  make build           - Build the server binary"
	@echo "  make test            - Run all unit tests"
	@echo "  make test-unit       - Run unit tests only"
	@echo "  make test-coverage   - Run tests with coverage report"
	@echo "  make test-integration- Run integration tests (requires .env file)"
	@echo "  make clean           - Remove build artifacts"
	@echo "  make run             - Run the server"
	@echo "  make lint            - Run code quality checks"
	@echo "  make pre-commit      - Run checks before committing"
	@echo "  make ci-check        - Simulate GitHub Actions CI checks locally"
	@echo ""

# Build the server
build:
	@echo "Building MCP server..."
	@go build -o mcp-server ./cmd/server
	@echo "✓ Build complete: ./mcp-server"

# Run all unit tests
test:
	@echo "Running unit tests..."
	@go test ./... -v -timeout 30s

# Run tests without verbose output
test-unit:
	@echo "Running unit tests..."
	@go test ./... -timeout 30s -count=1

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out -covermode=atomic -timeout 30s
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out | grep total
	@echo "✓ Coverage report generated: coverage.html"

# Run integration tests (requires API keys in .env)
test-integration:
	@echo "Running integration tests..."
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found. Copy .env.example and add your API keys."; \
		exit 1; \
	fi
	@./scripts/test-server.sh

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f mcp-server
	@rm -f coverage.out coverage.html
	@echo "✓ Clean complete"

# Run the server
run: build
	@echo "Starting MCP server..."
	@./mcp-server

# Lint and format code
lint:
	@echo "Running code quality checks..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, running basic checks..."; \
		go vet ./...; \
		gofmt -l .; \
	fi
	@echo "✓ Lint complete"

# Run quick checks before committing
pre-commit: lint test-unit
	@echo "✓ Pre-commit checks passed"

# Simulate CI checks locally
ci-check: lint test-coverage
	@echo "Running CI simulation..."
	@go vet ./...
	@go build -v -o mcp-server ./cmd/server
	@rm -f mcp-server
	@echo "✓ CI checks passed - ready to push!"
