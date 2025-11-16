.PHONY: build test clean install run fmt vet deps help

# Build variables
BINARY_NAME=rtc-scheduler
BUILD_DIR=bin
MAIN_PATH=cmd/rtc-scheduler/main.go
INSTALL_PATH=/usr/local/bin

# Version information
VERSION?=1.0.0
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

test:
	go test -v -race -cover ./...

clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)/
	rm -f coverage.out coverage.html

install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/
	sudo chmod +x $(INSTALL_PATH)/$(BINARY_NAME)

run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

fmt:
	go fmt ./...

vet:
	go vet ./...

deps:
	go mod download
	go mod tidy

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install it from https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

help:
	@echo "Available targets:"
	@echo "  build      - Build the binary"
	@echo "  test       - Run tests with race detection"
	@echo "  clean      - Clean build artifacts"
	@echo "  install    - Install binary to $(INSTALL_PATH)"
	@echo "  run        - Build and run the application"
	@echo "  fmt        - Format Go code"
	@echo "  vet        - Run go vet"
	@echo "  deps       - Download and tidy dependencies"
	@echo "  lint       - Run golangci-lint (if installed)"
	@echo "  coverage   - Generate coverage report"
	@echo "  help       - Show this help message"