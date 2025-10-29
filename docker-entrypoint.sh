#!/bin/sh

DB_DRIVER="mysql"
DB_CONN="$MYSQL_USER:$MYSQL_PASSWORD@tcp($DB_HOST:$DB_PORT)/$MYSQL_DATABASE?parseTime=true"
GOOSE_DIR="./migrations"

# 1. Esperar a que la Base de Datos esté disponible
echo "⏳ Esperando que MySQL esté disponible en $DB_HOST:$DB_PORT..."
until nc -z "$DB_HOST" "$DB_PORT"; do
  sleep 1
done
echo "✅ MySQL listo. Iniciando migraciones..."

# 2. Ejecutar las migraciones con Goose
/usr/bin/goose "$DB_DRIVER" "$DB_CONN" up -dir "$GOOSE_DIR"

if [ $? -ne 0 ]; then
  echo "❌ ERROR: Las migraciones de Goose fallaron. Saliendo del contenedor."
  exit 1
fi

echo "✅ Migraciones completadas. Iniciando la aplicación Go..."

# 4. Ejecutar el comando principal de la aplicación (definido en CMD)
exec "$@"
