.PHONY: run test tidy build docker-up docker-down docker-logs

run:
	go run ./cmd/api

test:
	go test ./...

tidy:
	go mod tidy

build:
	go build -o bin/api ./cmd/api

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f
