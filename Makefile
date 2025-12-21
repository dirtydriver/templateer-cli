# Makefile for the project-generator

# Variables
BINARY_NAME := templateer
BUILD_DIR   := bin
SRC         := ./main.go

# Ensure these targets are always executed even if a file with the same name exists
.PHONY: all build run test fmt lint vet deps clean

# Default target
all: build

# Build the binary into the specified build directory.
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC)

# Run the application (builds first if needed).
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run unit tests.
test:
	@echo "Running tests..."
	go test -v ./...

# Run unit tests.
cover:
	@echo "Running Coverage..."
	go test -cover ./...

# Format the Go source files.
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint the Go source files.
# Note: Ensure that golint is installed (`go install golang.org/x/lint/golint@latest`)
lint:
	@echo "Linting code..."
	golint ./...

# Run go vet for static analysis.
vet:
	@echo "Running go vet..."
	go vet ./...

# Download dependencies.
deps:
	@echo "Downloading dependencies..."
	go mod download

# Clean up the build artifacts.
clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)