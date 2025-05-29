// integration/user_crud_test.go
package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"encontro/cmd/server"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var baseURL string

// TestMain настраивает подключение к БД, проверяет её доступность,
// поднимает тестовый HTTPS-сервер и проверяет, что эндпоинт /health работает.
func TestMain(m *testing.M) {
	// 1. Читаем DSN из окружения
	dsn := os.Getenv("TEST_DSN")
	if dsn == "" {
		log.Fatal("TEST_DSN must be set")
	}

	// 2. Открываем подключение к БД
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}

	// 3. Ждём, пока БД примет соединение
	for i := 0; i < 10; i++ {
		if err := db.Ping(); err == nil {
			log.Println("DB connection established")
			break
		}
		log.Printf("Waiting for DB... attempt %d/10", i+1)
		time.Sleep(2 * time.Second)
		if i == 9 {
			log.Fatalf("could not connect to DB: %v", err)
		}
	}

	// 4. Инициализируем роутер приложения
	router := server.SetupRouter(db)

	// 5. Запускаем тестовый TLS-сервер
	ts := httptest.NewTLSServer(router)
	defer ts.Close()

	// 6. Ждём, пока сервер ответит 200 на /health
	client := ts.Client()
	healthURL := ts.URL + "/health"
	for i := 0; i < 10; i++ {
		resp, err := client.Get(healthURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Println("Server health check passed")
			break
		}
		log.Printf("Waiting for server... attempt %d/10", i+1)
		time.Sleep(2 * time.Second)
		if i == 9 {
			log.Fatalf("server health check failed: %v", err)
		}
	}

	// 7. Устанавливаем базовый URL для CRUD тестов
	baseURL = ts.URL + "/api/v1/users"

	os.Exit(m.Run())
}

// TestUserCRUD проверяет полный цикл CRUD-операций для ресурса User
func TestUserCRUD(t *testing.T) {
	type User struct {
		ID    int64  `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	// CREATE
	payload := `{"name":"Test","email":"test@ex.com"}`
	res, err := http.Post(baseURL, "application/json", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatalf("Create request failed: %v", err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", res.StatusCode)
	}
	var u User
	if err := json.NewDecoder(res.Body).Decode(&u); err != nil {
		t.Fatalf("failed to decode create response: %v", err)
	}
	if u.ID == 0 {
		t.Fatal("empty ID after create")
	}

	// READ
	res, err = http.Get(fmt.Sprintf("%s/%d", baseURL, u.ID))
	if err != nil {
		t.Fatalf("GetByID request failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on get, got %d", res.StatusCode)
	}
	var got User
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode get response: %v", err)
	}
	if got.Email != u.Email {
		t.Errorf("expected email %s, got %s", u.Email, got.Email)
	}

	// UPDATE
	updatePayload := `{"name":"Updated","email":"upd@ex.com"}`
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%d", baseURL, u.ID), bytes.NewBufferString(updatePayload))
	req.Header.Set("Content-Type", "application/json")
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Update request failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on update, got %d", res.StatusCode)
	}

	// LIST
	res, err = http.Get(baseURL)
	if err != nil {
		t.Fatalf("List request failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on list, got %d", res.StatusCode)
	}
	var list []User
	if err := json.NewDecoder(res.Body).Decode(&list); err != nil {
		t.Fatalf("failed to decode list response: %v", err)
	}
	if len(list) == 0 {
		t.Error("expected non-empty user list")
	}

	// DELETE
	req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%d", baseURL, u.ID), nil)
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Delete request failed: %v", err)
	}
	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204 on delete, got %d", res.StatusCode)
	}

	// GET после DELETE -> 404
	res, err = http.Get(fmt.Sprintf("%s/%d", baseURL, u.ID))
	if err != nil {
		t.Fatalf("GetByID after delete failed: %v", err)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", res.StatusCode)
	}
}
