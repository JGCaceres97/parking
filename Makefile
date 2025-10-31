include .env
export $(shell sed 's/=.*//' .env)

# Variables de la Base de Datos para Goose
# Se componen en el Makefile para que Goose pueda usarlas directamente.
DB_DRIVER=mysql
DB_USER=$(MYSQL_USER)
DB_PASS=$(MYSQL_PASSWORD)
DB_HOST=$(DB_HOST)
DB_PORT=$(DB_PORT)
DB_NAME=$(MYSQL_DATABASE)
DB_CONN=$(DB_USER):$(DB_PASS)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)?parseTime=true&loc=UTC

GOOSE_DIR=./migrations

.PHONY: all db-up db-down migrate-up migrate-down migrate-status

## DB (Docker Compose)
# ------------------------------------------------------------
# Inicia los contenedores de Docker en segundo plano (incluyendo MySQL).
db-up:
	@echo "🛠️ Levantando contenedores Docker..."
	docker compose up -d

# Detiene y elimina los contenedores.
db-down:
	@echo "🛑 Deteniendo contenedores Docker..."
	docker compose down

## Migraciones (Goose)
# ------------------------------------------------------------
# Aplica todas las migraciones pendientes.
# go tool goose mysql "parkingUser:parkingUserPassword@tcp(localhost:3306)/parkingDb?parseTime=true&loc=UTC" up -dir ./migrations
migrate-up:
	@echo "⬆️ Aplicando migraciones pendientes..."
	go tool goose $(DB_DRIVER) "$(DB_CONN)" up -dir $(GOOSE_DIR)


# Revierte la última migración aplicada.
migrate-down:
	@echo "⬇️ Revertiendo la última migración..."
	go tool goose $(DB_DRIVER) "$(DB_CONN)" down -dir $(GOOSE_DIR)

# Revierte todas las migraciones aplicadas.
migrate-reset:
	@echo "🧹 Revirtiendo todas las migraciones..."
	go tool goose $(DB_DRIVER) "$(DB_CONN)" reset -dir $(GOOSE_DIR)

# Muestra el estado actual de las migraciones.
migrate-status:
	@echo "📊 Verificando estado de las migraciones..."
	go tool goose $(DB_DRIVER) "$(DB_CONN)" status -dir $(GOOSE_DIR)

# Comando por defecto
all: db-up migrate-up
