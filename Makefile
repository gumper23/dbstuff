# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=test-dbhelper

all: test tidy build
build: 
		$(GOBUILD) -o $(BINARY_NAME) -v
tidy:
		$(GOMOD) tidy
test: 
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_LINUX)
		rm -f main
run: build
		./$(BINARY_NAME)
