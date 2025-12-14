.PHONY: build test lint clean install run help

# Build the binary
build:
	@echo "Building growth..."
	@mkdir -p bin
	@go build -o bin/growth cmd/growth/main.go
	@echo "Build complete: bin/growth"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.txt ./...
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report: coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.txt coverage.html
	@echo "Clean complete"

# Install to $GOPATH/bin
install:
	@echo "Installing growth..."
	@go install cmd/growth/main.go

# Run the application
run:
	@go run cmd/growth/main.go

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p bin
	@GOOS=linux GOARCH=amd64 go build -o bin/growth-linux-amd64 cmd/growth/main.go
	@GOOS=darwin GOARCH=amd64 go build -o bin/growth-darwin-amd64 cmd/growth/main.go
	@GOOS=darwin GOARCH=arm64 go build -o bin/growth-darwin-arm64 cmd/growth/main.go
	@GOOS=windows GOARCH=amd64 go build -o bin/growth-windows-amd64.exe cmd/growth/main.go
	@echo "Multi-platform build complete"

# Show help
help:
	@echo "Available targets:"
	@echo "  build          - Build the growth binary"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  lint           - Run golangci-lint"
	@echo "  clean          - Remove build artifacts"
	@echo "  install        - Install to GOPATH/bin"
	@echo "  run            - Run the application"
	@echo "  build-all      - Build for multiple platforms"
	@echo "  help           - Show this help message"
