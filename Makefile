EXE = d8b
PACKAGE = github.com/DillonBarker/d8b/cmd/d8b
COMMAND_PACKAGE = github.com/DillonBarker/d8b/internal/main.go
GO_VERSION = $(shell go version | cut -d ' ' -f3)
COMMIT_HASH	= $(shell git rev-parse HEAD)

include version.make

.PHONY: all build lint fmt clean

all: build

build: test
	@go build -ldflags="-X $(COMMAND_PACKAGE).BuildVersion=$(VERSION) -X $(COMMAND_PACKAGE).BuildCommit=$(COMMIT_HASH) -X $(COMMAND_PACKAGE).BuildGoVersion=$(GO_VERSION)" -o $(GOPATH)/bin/$(EXE) $(PACKAGE)

test: lint
	@gotestsum -f short-verbose ./...

install-dependencies:
	@go get ./...

lint: fmt
	@golangci-lint run

fmt:
	@gofmt -w ./

clean:
	@rm -rf $(GOPATH)/bin/$(EXE)