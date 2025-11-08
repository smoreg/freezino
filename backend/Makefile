.PHONY: run build dev clean test install

# Application name
APP_NAME=freezino-server

# Build directory
BUILD_DIR=./bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Main package path
MAIN_PATH=./cmd/server

# Run the application
run:
	@echo "🚀 Running application..."
	@$(GORUN) $(MAIN_PATH)/main.go

# Build the application
build:
	@echo "🔨 Building application..."
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)/main.go
	@echo "✅ Build complete: $(BUILD_DIR)/$(APP_NAME)"

# Development mode with auto-reload (requires air or similar)
dev:
	@echo "🔥 Starting development server..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "⚠️  'air' not installed. Installing..."; \
		go install github.com/air-verse/air@latest; \
		air; \
	fi

# Install dependencies
install:
	@echo "📦 Installing dependencies..."
	@$(GOGET) -v ./...
	@$(GOMOD) tidy
	@echo "✅ Dependencies installed"

# Run tests
test:
	@echo "🧪 Running tests..."
	@$(GOTEST) -v ./...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete"

# Tidy dependencies
tidy:
	@echo "📦 Tidying dependencies..."
	@$(GOMOD) tidy
	@echo "✅ Dependencies tidied"

# Format code
fmt:
	@echo "🎨 Formatting code..."
	@$(GOCMD) fmt ./...
	@echo "✅ Code formatted"

# Run linter (requires golangci-lint)
lint:
	@echo "🔍 Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "⚠️  'golangci-lint' not installed. Skipping..."; \
	fi

# Help
help:
	@echo "Available commands:"
	@echo "  make run     - Run the application"
	@echo "  make build   - Build the application"
	@echo "  make dev     - Run in development mode with auto-reload"
	@echo "  make install - Install dependencies"
	@echo "  make test    - Run tests"
	@echo "  make clean   - Clean build artifacts"
	@echo "  make tidy    - Tidy dependencies"
	@echo "  make fmt     - Format code"
	@echo "  make lint    - Run linter"
	@echo "  make help    - Show this help message"
