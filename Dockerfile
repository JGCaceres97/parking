# 1. Web
FROM node:krypton-slim AS web

WORKDIR /app

ENV NODE_ENV=production

COPY web/package*.json ./

RUN npm ci --include=dev

COPY web .

RUN npm run build

# 2. API
FROM golang:1.25.3-alpine AS api

WORKDIR /app

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY go.mod go.sum ./
RUN go mod download

# Copiar compilado de web y handler
COPY --from=web /app/dist ./web/dist
COPY web/web.go ./web

COPY . .

# Pruebas y compilación
RUN CGO_ENABLED=0 GOOS=linux go test -failfast -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /parking-system ./cmd/main.go

# 3. Run
FROM alpine:latest

RUN apk update && apk --no-cache add ca-certificates busybox dos2unix

WORKDIR /app

# Binarios
COPY --from=api /parking-system .
COPY --from=api /go/bin/goose /usr/bin/goose

# Script de entrada
COPY docker-entrypoint.sh .
RUN dos2unix docker-entrypoint.sh && chmod +x docker-entrypoint.sh

# Archivos necesarios
COPY ./migrations ./migrations

EXPOSE 3000

ENTRYPOINT [ "/app/docker-entrypoint.sh" ]

CMD [ "/app/parking-system" ]
