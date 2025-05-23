#!/bin/bash

# Установка переменных окружения для базы данных
export DB_HOST=77.221.159.137
export DB_PORT=5432
export DB_USER=kratoff
export DB_PASSWORD=Oleg253535
export DB_NAME=encontro
export DB_SSLMODE=require

# Загрузка тестовых данных
psql "host=$DB_HOST port=$DB_PORT user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME sslmode=$DB_SSLMODE" -f scripts/testdata.sql

# Запуск сервера
go run cmd/server/main.go 