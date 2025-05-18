#!/bin/bash

# Базовый URL API
BASE_URL="http://localhost:8080"

# Цвета для вывода
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

echo "🧪 Тестирование API..."

# 1. Создание комнаты
echo -e "\n1. Создание комнаты"
ROOM_RESPONSE=$(curl -s -X POST "$BASE_URL/rooms" \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Room"}' | jq '.')
ROOM_ID=$(echo $ROOM_RESPONSE | jq -r '.id')
echo "Создана комната: $ROOM_RESPONSE"

# 2. Получение списка комнат
echo -e "\n2. Получение списка комнат"
curl -s "$BASE_URL/rooms?page=1&page_size=20" | jq '.'

# 3. Создание сообщения
echo -e "\n3. Создание сообщения"
MESSAGE_RESPONSE=$(curl -s -X POST "$BASE_URL/rooms/$ROOM_ID/messages" \
  -H "Content-Type: application/json" \
  -d '{"content": "Test message"}' | jq '.')
MESSAGE_ID=$(echo $MESSAGE_RESPONSE | jq -r '.id')
echo "Создано сообщение: $MESSAGE_RESPONSE"

# 4. Получение списка сообщений
echo -e "\n4. Получение списка сообщений"
curl -s "$BASE_URL/rooms/$ROOM_ID/messages?page=1&page_size=20" | jq '.'

# 5. Получение сообщения по ID
echo -e "\n5. Получение сообщения по ID"
curl -s "$BASE_URL/rooms/$ROOM_ID/messages/$MESSAGE_ID" | jq '.'

# 6. Удаление сообщения
echo -e "\n6. Удаление сообщения"
curl -s -X DELETE "$BASE_URL/rooms/$ROOM_ID/messages/$MESSAGE_ID"
echo "Сообщение удалено"

# 7. Удаление комнаты
echo -e "\n7. Удаление комнаты"
curl -s -X DELETE "$BASE_URL/rooms/$ROOM_ID"
echo "Комната удалена"

echo -e "\n✅ Тестирование завершено" 