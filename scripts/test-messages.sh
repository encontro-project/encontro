#!/bin/bash

# Цвета для вывода
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Базовый URL
BASE_URL="https://localhost:8443/api"
USER_ID="test-user-123"

# Функция для вывода результата
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
    else
        echo -e "${RED}✗ $2${NC}"
        echo "Response: $3"
    fi
}

# Функция для проверки JSON ответа
check_json_response() {
    if ! echo "$1" | jq -e . >/dev/null 2>&1; then
        return 1
    fi
    return 0
}

# Создаем комнату
echo "Создаем тестовую комнату..."
ROOM_RESPONSE=$(curl -s -k -X POST "$BASE_URL/rooms" \
    -H "Content-Type: application/json" \
    -d '{"name": "Test Room", "description": "Room for testing messages"}')
if ! check_json_response "$ROOM_RESPONSE"; then
    print_result 1 "Создание комнаты" "Invalid JSON response"
    exit 1
fi
ROOM_ID=$(echo "$ROOM_RESPONSE" | jq -r '.id')
if [ -z "$ROOM_ID" ] || [ "$ROOM_ID" = "null" ]; then
    print_result 1 "Создание комнаты" "Failed to get room ID"
    exit 1
fi
print_result 0 "Создание комнаты" "$ROOM_RESPONSE"

# Создаем сообщение
echo -e "\nСоздаем тестовое сообщение..."
MESSAGE_RESPONSE=$(curl -s -k -X POST "$BASE_URL/messages" \
    -H "Content-Type: application/json" \
    -H "X-User-ID: $USER_ID" \
    -d "{\"content\": \"Test message\", \"room_id\": \"$ROOM_ID\"}")
if ! check_json_response "$MESSAGE_RESPONSE"; then
    print_result 1 "Создание сообщения" "Invalid JSON response"
    exit 1
fi
MESSAGE_ID=$(echo "$MESSAGE_RESPONSE" | jq -r '.id')
if [ -z "$MESSAGE_ID" ] || [ "$MESSAGE_ID" = "null" ]; then
    print_result 1 "Создание сообщения" "Failed to get message ID"
    exit 1
fi
print_result 0 "Создание сообщения" "$MESSAGE_RESPONSE"

# Сохраняем ID сообщения для последующих проверок
SAVED_MESSAGE_ID="$MESSAGE_ID"

# Получаем список сообщений
echo -e "\nПолучаем список сообщений..."
MESSAGES_RESPONSE=$(curl -s -k -X GET "$BASE_URL/messages?room_id=$ROOM_ID&page=1&page_size=10" \
    -H "X-User-ID: $USER_ID")
print_result $? "Получение списка сообщений" "$MESSAGES_RESPONSE"

# Получаем сообщение по ID
echo -e "\nПолучаем сообщение по ID..."
GET_MESSAGE_RESPONSE=$(curl -s -k -X GET "$BASE_URL/messages/$SAVED_MESSAGE_ID" \
    -H "X-User-ID: $USER_ID")
print_result $? "Получение сообщения по ID" "$GET_MESSAGE_RESPONSE"

# Обновляем сообщение
echo -e "\nОбновляем сообщение..."
UPDATE_RESPONSE=$(curl -s -k -X PUT "$BASE_URL/messages/$SAVED_MESSAGE_ID" \
    -H "Content-Type: application/json" \
    -H "X-User-ID: $USER_ID" \
    -d '{"content": "Updated test message"}')
print_result $? "Обновление сообщения" "$UPDATE_RESPONSE"

# Проверяем обновленное сообщение
echo -e "\nПроверяем обновленное сообщение..."
UPDATED_MESSAGE_RESPONSE=$(curl -s -k -X GET "$BASE_URL/messages/$SAVED_MESSAGE_ID" \
    -H "X-User-ID: $USER_ID")
print_result $? "Проверка обновленного сообщения" "$UPDATED_MESSAGE_RESPONSE"

# Удаляем сообщение
echo -e "\nУдаляем сообщение..."
DELETE_RESPONSE=$(curl -s -k -X DELETE "$BASE_URL/messages/$SAVED_MESSAGE_ID" \
    -H "X-User-ID: $USER_ID")
print_result $? "Удаление сообщения" "$DELETE_RESPONSE"

# Проверяем, что сообщение удалено
echo -e "\nПроверяем, что сообщение удалено..."
CHECK_DELETED_RESPONSE=$(curl -s -k -X GET "$BASE_URL/messages/$SAVED_MESSAGE_ID" \
    -H "X-User-ID: $USER_ID")
if [[ $CHECK_DELETED_RESPONSE == *"message not found"* ]]; then
    print_result 0 "Проверка удаления сообщения" "$CHECK_DELETED_RESPONSE"
else
    print_result 1 "Проверка удаления сообщения" "$CHECK_DELETED_RESPONSE"
fi

# Удаляем тестовую комнату
echo -e "\nУдаляем тестовую комнату..."
DELETE_ROOM_RESPONSE=$(curl -s -k -X DELETE "$BASE_URL/rooms/$ROOM_ID" \
    -H "X-User-ID: $USER_ID")
print_result $? "Удаление комнаты" "$DELETE_ROOM_RESPONSE"

echo -e "\nТестирование завершено!" 