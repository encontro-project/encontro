#!/bin/bash

# Цвета для вывода
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Базовый URL
BASE_URL="https://localhost:8443"

# Функция для вывода результата
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
    else
        echo -e "${RED}✗ $2${NC}"
        echo -e "${YELLOW}Ответ сервера:${NC}"
        echo "$3"
    fi
}

echo "Начинаем тестирование CRUD операций..."

# 1. Создание комнаты
echo -e "\n1. Создание комнаты"
ROOM_RESPONSE=$(curl -s -X POST "$BASE_URL/api/rooms" \
    -k \
    -H "Content-Type: application/json" \
    -d '{"name": "Test Room", "description": "Room for testing CRUD operations"}')
echo "Ответ: $ROOM_RESPONSE"
ROOM_ID=$(echo $ROOM_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
print_result $? "Создание комнаты" "$ROOM_RESPONSE"

# 2. Получение списка комнат
echo -e "\n2. Получение списка комнат"
ROOMS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/rooms?page=1&page_size=10" -k)
print_result $? "Получение списка комнат" "$ROOMS_RESPONSE"

# 3. Получение конкретной комнаты
echo -e "\n3. Получение конкретной комнаты"
ROOM_DETAIL_RESPONSE=$(curl -s -X GET "$BASE_URL/api/rooms/$ROOM_ID" -k)
print_result $? "Получение конкретной комнаты" "$ROOM_DETAIL_RESPONSE"

# 4. Создание сообщения
echo -e "\n4. Создание сообщения"
MESSAGE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/messages" \
    -k \
    -H "Content-Type: application/json" \
    -d "{\"room_id\": \"$ROOM_ID\", \"content\": \"Test message\", \"type\": \"text\", \"user_id\": \"test-user\"}")
echo "Ответ: $MESSAGE_RESPONSE"
MESSAGE_ID=$(echo $MESSAGE_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
print_result $? "Создание сообщения" "$MESSAGE_RESPONSE"

# 5. Получение сообщений комнаты
echo -e "\n5. Получение сообщений комнаты"
MESSAGES_RESPONSE=$(curl -s -X GET "$BASE_URL/api/messages?room_id=$ROOM_ID&page=1&page_size=10" -k)
print_result $? "Получение сообщений комнаты" "$MESSAGES_RESPONSE"

# 6. Получение конкретного сообщения
echo -e "\n6. Получение конкретного сообщения"
MESSAGE_DETAIL_RESPONSE=$(curl -s -X GET "$BASE_URL/api/messages/$MESSAGE_ID" -k)
print_result $? "Получение конкретного сообщения" "$MESSAGE_DETAIL_RESPONSE"

# 7. Удаление сообщения
echo -e "\n7. Удаление сообщения"
DELETE_MESSAGE_RESPONSE=$(curl -s -X DELETE "$BASE_URL/api/messages/$MESSAGE_ID" -k)
print_result $? "Удаление сообщения" "$DELETE_MESSAGE_RESPONSE"

# 8. Удаление комнаты
echo -e "\n8. Удаление комнаты"
DELETE_ROOM_RESPONSE=$(curl -s -X DELETE "$BASE_URL/api/rooms/$ROOM_ID" -k)
print_result $? "Удаление комнаты" "$DELETE_ROOM_RESPONSE"

echo -e "\nТестирование завершено!" 