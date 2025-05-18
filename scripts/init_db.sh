#!/bin/bash

# Проверяем наличие psql
if ! command -v psql &> /dev/null; then
    echo "❌ psql не установлен. Пожалуйста, установите PostgreSQL client tools."
    exit 1
fi

# Проверяем обязательные переменные окружения
if [ -z "$DB_USER" ]; then
    echo "❌ Переменная окружения DB_USER не задана!"
    exit 1
fi
if [ -z "$DB_PASSWORD" ]; then
    echo "❌ Переменная окружения DB_PASSWORD не задана!"
    exit 1
fi

# Параметры подключения к БД из переменных окружения или значения по умолчанию
DB_HOST="${DB_HOST:-77.221.159.137}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-encontro}"
DB_SSLMODE="${DB_SSLMODE:-require}"

echo "🔄 Создание базы данных..."

# Создаем базу данных
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "CREATE DATABASE $DB_NAME;"

if [ $? -ne 0 ]; then
    echo "⚠️ База данных уже существует или не удалось её создать"
fi

echo "🔄 Применение схемы базы данных..."

# Применяем схему базы данных
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -v sslmode=$DB_SSLMODE -f internal/infrastructure/database/schema.sql

if [ $? -eq 0 ]; then
    echo "✅ База данных успешно инициализирована"
else
    echo "❌ Ошибка при инициализации базы данных"
    exit 1
fi 