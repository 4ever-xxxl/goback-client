# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=goback-client.exe
BINARY_UNIX=$(BINARY_NAME)_unix

# Default target executed when no arguments are given to make.
all: test build

# Build the project
build: 
	$(GOBUILD) -o $(BINARY_NAME) -v

# Run the project
run: build
	./$(BINARY_NAME)

# Test the project
test: 
	$(GOTEST) -v ./...

# Clean the project
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Install dependencies
deps: 
	$(GOGET) -v ./...

# Cross compile for Linux
build-linux: 
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

.PHONY: all build clean run test deps build-linux