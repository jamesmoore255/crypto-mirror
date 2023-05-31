# Variables
GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
BINARY_NAME = server
BUILD_DIR = ./build

# Default target
all: build

# Build target
build:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

# Test target
test:
	$(GOTEST) ./...

# Clean target
clean:
	rm -rf $(BUILD_DIR)

.PHONY: all build test clean

.PHONY: run
run:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server
	$(BUILD_DIR)/$(BINARY_NAME)