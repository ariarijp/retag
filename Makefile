.PHONY: all test lint

all: test lint

test:
	go test -v ./...

lint:
	golangci-lint run --out-format=github-actions