#!/bin/sh
set -e

# Ожидаем доступность БД и применяем миграции
until goose -dir /app/migrations postgres "$DB_DSN" up; do
  echo "Waiting for DB to be ready..."
  sleep 2
done

# Запуск приложения
exec /usr/local/bin/app
