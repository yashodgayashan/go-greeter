# Makefile for Go Greeter Service

# Variables
BINARY_NAME=go-greeter
MAIN_FILE=main.go
TEST_TIMEOUT=30s

# Default target
.PHONY: help
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Build commands
.PHONY: build
build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) $(MAIN_FILE)

.PHONY: build-linux
build-linux: ## Build for Linux
	@echo "Building $(BINARY_NAME) for Linux..."
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux $(MAIN_FILE)

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME) $(BINARY_NAME)-linux
	rm -f coverage.out coverage.html

# Test commands
.PHONY: test
test: ## Run all tests
	@echo "Running tests..."
	go test -v -timeout $(TEST_TIMEOUT)

.PHONY: test-short
test-short: ## Run tests without verbose output
	@echo "Running tests (short)..."
	go test -timeout $(TEST_TIMEOUT)

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-coverage-func
test-coverage-func: ## Show function-level coverage
	@echo "Running tests with function coverage..."
	go test -cover -coverprofile=coverage.out
	go tool cover -func=coverage.out

.PHONY: benchmark
benchmark: ## Run benchmark tests
	@echo "Running benchmarks..."
	go test -bench=. -benchmem

# Linting commands
.PHONY: lint
lint: ## Run golangci-lint
	@echo "Running golangci-lint..."
	golangci-lint run

.PHONY: lint-clean
lint-clean: ## Run golangci-lint excluding common test file issues
	@echo "Running golangci-lint (excluding test file issues)..."
	golangci-lint run --exclude="Error return value of.*io\\.ReadAll.*is not checked" --exclude="Error return value of.*json\\.Marshal.*is not checked"

.PHONY: lint-comprehensive
lint-comprehensive: ## Run comprehensive linting with nice output
	@echo "Running comprehensive linting..."
	./lint.sh

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with auto-fix
	@echo "Running golangci-lint with auto-fix..."
	golangci-lint run --fix

.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting Go code..."
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

.PHONY: mod-tidy
mod-tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	go mod tidy

# Quality checks (combination of multiple checks)
.PHONY: check
check: fmt vet lint test ## Run all quality checks (format, vet, lint, test)

.PHONY: check-all
check-all: fmt vet lint test-coverage benchmark ## Run all checks including coverage and benchmarks

# Run commands
.PHONY: run
run: ## Run the application
	@echo "Starting $(BINARY_NAME)..."
	go run $(MAIN_FILE)

.PHONY: run-build
run-build: build ## Build and run the application
	@echo "Running built $(BINARY_NAME)..."
	./$(BINARY_NAME)

# Development commands
.PHONY: deps
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download

.PHONY: deps-update
deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# Docker commands (if Dockerfile exists)
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) .

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 9090:9090 $(BINARY_NAME)

# API testing commands
.PHONY: test-api
test-api: ## Test API endpoints manually (requires server to be running)
	@echo "Testing API endpoints..."
	@echo "Testing /greeter/greet:"
	curl -s "http://localhost:9090/greeter/greet?name=World" || echo "Server not running?"
	@echo ""
	@echo "Testing /greeter/health:"
	curl -s "http://localhost:9090/greeter/health" | jq '.' 2>/dev/null || curl -s "http://localhost:9090/greeter/health"
	@echo ""
	@echo "Testing /greeter/farewell:"
	curl -s "http://localhost:9090/greeter/farewell?name=World"
	@echo ""

# Install linting tools
.PHONY: install-tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@echo "Tools installed successfully!"

# Git hooks
.PHONY: install-hooks
install-hooks: ## Install git pre-commit hooks
	@echo "Installing git pre-commit hooks..."
	@echo '#!/bin/sh\nmake check' > .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "Pre-commit hook installed!"

# Show project info
.PHONY: info
info: ## Show project information
	@echo "Project: $(BINARY_NAME)"
	@echo "Go version: $$(go version)"
	@echo "Git branch: $$(git branch --show-current 2>/dev/null || echo 'unknown')"
	@echo "Git commit: $$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
	@echo "Files:"
	@find . -name "*.go" -not -path "./vendor/*" | wc -l | xargs echo "  Go files:"
	@wc -l *.go | tail -1 