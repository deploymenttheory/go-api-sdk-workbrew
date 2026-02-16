.PHONY: help test test-unit test-acceptance test-acceptance-verbose test-acceptance-service test-all clean

# Default target
help:
	@echo "Workbrew SDK Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  make test                    - Run all unit tests"
	@echo "  make test-unit               - Run unit tests only"
	@echo "  make test-acceptance         - Run acceptance tests (requires API credentials)"
	@echo "  make test-acceptance-verbose - Run acceptance tests with verbose output"
	@echo "  make test-acceptance-service - Run acceptance tests for specific service (usage: make test-acceptance-service SERVICE=Devices)"
	@echo "  make test-all                - Run all tests (unit + acceptance)"
	@echo "  make clean                   - Clean test cache and artifacts"
	@echo ""
	@echo "Environment variables for acceptance tests:"
	@echo "  WORKBREW_API_KEY           - Your Workbrew API key (required)"
	@echo "  WORKBREW_WORKSPACE_NAME    - Your workspace name (required)"
	@echo "  WORKBREW_VERBOSE=true      - Enable verbose test output"
	@echo "  WORKBREW_SKIP_CLEANUP=true - Skip cleanup after tests"
	@echo ""
	@echo "Example usage:"
	@echo "  make test-unit"
	@echo "  WORKBREW_VERBOSE=true make test-acceptance"
	@echo "  make test-acceptance-service SERVICE=Devices"

# Run all unit tests
test: test-unit

# Run unit tests
test-unit:
	@echo "Running unit tests..."
	@go test -v -race -coverprofile=coverage.txt -covermode=atomic ./workbrew/client/... ./workbrew/services/...

# Run acceptance tests
test-acceptance:
	@echo "Running acceptance tests..."
	@echo "Note: This requires WORKBREW_API_KEY and WORKBREW_WORKSPACE_NAME to be set"
	@go test -v -timeout 30m ./workbrew/acceptance/...

# Run acceptance tests with verbose output
test-acceptance-verbose:
	@echo "Running acceptance tests (verbose)..."
	@WORKBREW_VERBOSE=true go test -v -timeout 30m ./workbrew/acceptance/...

# Run acceptance tests for a specific service
# Usage: make test-acceptance-service SERVICE=Devices
test-acceptance-service:
	@echo "Running acceptance tests for service: $(SERVICE)..."
	@go test -v -timeout 10m -run TestAcceptance_$(SERVICE) ./workbrew/acceptance/

# Run all tests (unit + acceptance)
test-all: test-unit test-acceptance
	@echo "All tests completed!"

# Clean test cache and build artifacts
clean:
	@echo "Cleaning test cache and artifacts..."
	@go clean -testcache
	@rm -f coverage.txt
	@rm -f workbrew/acceptance/*.log
	@echo "Clean complete!"

# Build the project
build:
	@echo "Building project..."
	@go build ./...

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

# Display coverage
coverage:
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.txt -covermode=atomic ./workbrew/...
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated: coverage.html"
