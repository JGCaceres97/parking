# 1. Build
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go test -failfast -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /parking-system ./cmd/main.go

# 2. Run
FROM alpine:latest

RUN apk update && apk --no-cache add ca-certificates busybox dos2unix

WORKDIR /app

# Binarios
COPY --from=builder /parking-system .
COPY --from=builder /go/bin/goose /usr/bin/goose

# Script de entrada
COPY docker-entrypoint.sh .
RUN dos2unix docker-entrypoint.sh && chmod +x docker-entrypoint.sh

# Archivos necesarios
COPY .env .
COPY ./migrations ./migrations

EXPOSE 3000

ENTRYPOINT [ "/app/docker-entrypoint.sh" ]

CMD [ "/app/parking-system" ]
