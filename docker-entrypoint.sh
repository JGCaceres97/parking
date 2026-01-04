#!/bin/sh

# Salida del script automática en caso de error
set -e

MIGRATIONS_DIR="./migrations/$DB_DRIVER"

if [ ! -d "$MIGRATIONS_DIR" ]; then
  echo "❌ Directorio de migraciones no existe: $MIGRATIONS_DIR"
  exit 1
fi

# 1. Construir DSN dependiendo el driver
case "$DB_DRIVER" in
  sqlite)
    if [ -z "$SQLITE_DSN" ]; then
      echo "❌ SQLITE_DSN no definido"
      exit 1
    fi

    DB_CONN="$SQLITE_DSN"
    ;;

  mysql)
    DB_CONN="$MYSQL_USER:$MYSQL_PASSWORD@tcp($DB_HOST:$DB_PORT)/$MYSQL_DATABASE?parseTime=true&loc=UTC"

    # 1.1. Esperar a que la Base de Datos esté disponible
    echo "⏳ Esperando que MySQL esté disponible en $DB_HOST:$DB_PORT..."
    until nc -z "$DB_HOST" "$DB_PORT"; do
      sleep 1
    done
    echo "✅ MySQL listo."
    ;;

  *)
    echo "❌ DB_DRIVER no soportado: $DB_DRIVER"
    exit 1
    ;;
esac

# 2. Ejecutar las migraciones
echo "⬆️ Iniciando migraciones desde $MIGRATIONS_DIR..."
/usr/bin/goose "$DB_DRIVER" "$DB_CONN" up -dir "$MIGRATIONS_DIR"

# 3. Ejecutar el comando principal de la aplicación (definido en CMD)
echo "✅ Migraciones completadas. Iniciando la aplicación Go..."
exec "$@"
