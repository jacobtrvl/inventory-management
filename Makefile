# Copyright 2025 Jacob Philip. All rights reserved.

# Variables
BINARY_NAME=inventory
MAIN_PATH=./cmd/inventory
BUILD_DIR=./bin

# Default target
.DEFAULT_GOAL := build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Test the application
test:
	@echo "Running tests..."
	go test ./...

ratelimiter:
	@echo "Running rate limiter..."
	go run ./cmd/ratelimiter

# Show help
help:
	@echo "Available targets:"
	@echo "  build  - Build the application"
	@echo "  run    - Build and run the application"
	@echo "  clean  - Remove build artifacts"
	@echo "  test   - Run tests"
	@echo "  ratelimiter - Run the rate limiter example"

.PHONY: build run clean test fmt vet help
