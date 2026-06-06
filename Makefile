.PHONY: run test tidy build

run:
	go run ./cmd/api

test:
	go test ./...

tidy:
	go mod tidy

build:
	go build -o bin/api ./cmd/api
