include .env
export $(shell sed 's/=.*//' .env)

GOOSE_DIR=./migrations

#-------------------------
# DB Según driver
#-------------------------
ifeq ($(DB_DRIVER),sqlite)
	DB_CONN := $(SQLITE_DSN)
	GOOSE_MIGRATIONS_DIR := $(GOOSE_DIR)/sqlite
endif

ifeq ($(DB_DRIVER),mysql)
	DB_CONN := $(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(MYSQL_DATABASE)?parseTime=true&loc=UTC
	GOOSE_MIGRATIONS_DIR := $(GOOSE_DIR)/mysql
endif

#-------------------------
# Targets
#-------------------------
.PHONY: all db-up db-down migrate-up migrate-down migrate-status

# Comando por defecto
all: db-up migrate-up

## DB (Docker Compose)
# ------------------------------------------------------------
# Inicia los contenedores de Docker en segundo plano.
db-up:
	ifeq ($(DB_DRIVER),mysql)
		@echo "🛠️ Levantando contenedores Docker..."
		docker compose up -d
	else
		@echo "ℹ️ SQLite no requiere Docker"
	endif

# Detiene y elimina los contenedores.
db-down:
	ifeq ($(DB_DRIVER),mysql)
		@echo "🛑 Deteniendo contenedores Docker..."
		docker compose down
	else
			@echo "ℹ️ SQLite no requiere Docker"
	endif

## Migraciones (Goose)
# ------------------------------------------------------------
# Aplica todas las migraciones pendientes.
# go tool goose sqlite "file:parking.db" up -dir ./migrations/sqlite/
# go tool goose mysql "parkingUser:parkingUserPassword@tcp(localhost:3306)/parkingDb?parseTime=true&loc=UTC" up -dir ./migrations/mysql
migrate-up:
	@echo "⬆️ Aplicando migraciones ($(DB_DRIVER))..."
	go tool goose $(DB_DRIVER) "$(DB_CONN)" up -dir $(GOOSE_MIGRATIONS_DIR)

# Revierte la última migración aplicada.
migrate-down:
	@echo "⬇️ Revertiendo la última migración ($(DB_DRIVER))..."
	go tool goose $(DB_DRIVER) "$(DB_CONN)" down -dir $(GOOSE_MIGRATIONS_DIR)

# Revierte todas las migraciones aplicadas.
migrate-reset:
	@echo "🧹 Revirtiendo todas las migraciones ($(DB_DRIVER))..."
	go tool goose $(DB_DRIVER) "$(DB_CONN)" reset -dir $(GOOSE_MIGRATIONS_DIR)

# Muestra el estado actual de las migraciones.
migrate-status:
	@echo "📊 Verificando estado de las migraciones..."
	go tool goose $(DB_DRIVER) "$(DB_CONN)" status -dir $(GOOSE_MIGRATIONS_DIR)
