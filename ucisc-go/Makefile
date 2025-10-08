.PHONY: build test bench lint clean run-simple run-fib help

# Default target
all: build test

# Build the emulator
build:
	@echo "Building ucisc-go..."
	go build -o bin/ucisc ./cmd/ucisc

# Build with optimizations for release
build-release:
	@echo "Building optimized release..."
	go build -ldflags="-s -w" -o bin/ucisc ./cmd/ucisc

# Run all tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	go test -race ./...

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Generate test coverage
coverage:
	@echo "Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Run simple example
run-simple: build
	@echo "Running simple example..."
	./bin/ucisc run examples/simple.hex

# Run fibonacci example
run-fib: build
	@echo "Running fibonacci example..."
	./bin/ucisc run examples/fibonacci.hex

# Install to GOPATH/bin
install:
	@echo "Installing to GOPATH/bin..."
	go install ./cmd/ucisc

# Help
help:
	@echo "ucisc-go Makefile targets:"
	@echo "  make build        - Build the emulator"
	@echo "  make test         - Run all tests"
	@echo "  make test-race    - Run tests with race detection"
	@echo "  make bench        - Run benchmarks"
	@echo "  make lint         - Run linter"
	@echo "  make coverage     - Generate coverage report"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make run-simple   - Run simple example"
	@echo "  make run-fib      - Run fibonacci example"
	@echo "  make install      - Install to GOPATH/bin"
