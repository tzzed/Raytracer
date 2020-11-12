# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=rt
BINARY_UNIX=$(BINARY_NAME)_unix

.PHONY: all $(NAME) test testrace build gen

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

gen: $(NAME)
	go generate ./...

test:
	go test -cover -timeout=1m ./...

testtinygo:
	go test -tags=tinygo -cover -timeout=1m ./...

bench:
	go test -v -run=^\$$ -benchmem -bench=. ./...
	cd cmd/genji && go test -v -run=^\$$ -benchmem -bench=. ./...
