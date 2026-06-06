FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api

FROM alpine:3.20

RUN apk add --no-cache ca-certificates wget

WORKDIR /app

COPY --from=builder /api .

RUN mkdir -p /app/data

ENV HTTP_ADDR=:8080 \
    DB_PATH=/app/data/requests.db \
    JWT_SECRET=change-me-in-production

EXPOSE 8080

CMD ["./api"]
