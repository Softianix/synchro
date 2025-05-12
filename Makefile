BINARY := synchro

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@echo "→ Building $(BINARY)..."
	go build -o bin/$(BINARY) main.go

# Run tests
.PHONY: test
test:
	@echo "→ Running tests..."
	go test ./tests/...

# Build then run
.PHONY: run
run: build
	@echo "→ Running $(BINARY)..."
	./bin/$(BINARY)

# Clean up build artifacts
.PHONY: clean
clean:
	@echo "→ Cleaning..."
	rm -rf bin/

# Install binary to $GOPATH/bin (optional)
.PHONY: install
install:
	@echo "→ Installing $(BINARY)..."
	go install ./...

# Display help
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make           # builds the binary"
	@echo "  make build     # same as 'make'"
	@echo "  make test      # runs all tests"
	@echo "  make run       # builds then runs the binary"
	@echo "  make clean     # removes bin/ directory"
	@echo "  make install   # installs binary to GOPATH/bin"
