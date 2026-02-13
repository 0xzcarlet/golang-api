.PHONY: help run build test clean tidy

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application in development mode
	@echo "🚀 Running application..."
	@go run cmd/api/main.go

build: ## Build the application binary
	@echo "🔨 Building application..."
	@go build -o bin/api cmd/api/main.go
	@echo "✅ Build complete: bin/api"

test: ## Run tests
	@echo "🧪 Running tests..."
	@go test -v ./...

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	@rm -rf bin/
	@echo "✅ Clean complete"

tidy: ## Tidy dependencies
	@echo "📦 Tidying dependencies..."
	@go mod tidy
	@echo "✅ Dependencies tidied"

deps: ## Download dependencies
	@echo "📥 Downloading dependencies..."
	@go mod download
	@echo "✅ Dependencies downloaded"

dev: ## Run with hot-reload using Air
	@echo "🔥 Running with hot-reload..."
	@air

air-init: ## Initialize Air configuration
	@air init
