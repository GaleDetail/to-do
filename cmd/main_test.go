package main

import (
	"database/sql"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupTestDB() (*sql.DB, error) {
	// Завантаження поточного робочого каталогу для перевірки
	path, err := os.Getwd()
	fmt.Println("Current working directory:", path)

	// Формування DSN для підключення до тестової БД
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DBUSER"),
		os.Getenv("DBPASS"),
		os.Getenv("DBHOST"),
		os.Getenv("DBPORT"),
		os.Getenv("DBNAME"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to test database: %v", err)
	}

	// Перевірка з'єднання
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Error pinging test database: %v", err)
	}

	return db, nil
}

// TestMain забезпечує підготовку середовища перед усіма тестами
func TestMain(m *testing.M) {
	// Спроба завантажити .env файл перед запуском тестів
	if err := godotenv.Load(); err != nil {
		log.Println("Could not load .env file for tests")
	}

	// Запуск усіх тестів
	code := m.Run()
	os.Exit(code)
}

func TestJsonResponse(t *testing.T) {
	expectedKey := "status"
	expectedMessage := "Server is running"
	payload, err := jsonResponse(expectedKey, expectedMessage)
	if err != nil {
		t.Fatalf("Error generating JSON response: %v", err)
	}

	json, err := simplejson.NewJson(payload)
	if err != nil {
		t.Fatalf("Error parsing JSON response: %v", err)
	}

	actualMessage, err := json.Get(expectedKey).String()
	if err != nil {
		t.Fatalf("Error retrieving key '%s' from JSON: %v", expectedKey, err)
	}

	if actualMessage != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, actualMessage)
	}
}

func TestPingHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(pingHandler)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "application/json")
	}

	expected := `{"status":"Server is running"}`
	if actualMessage := rec.Body.String(); actualMessage != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actualMessage, expected)
	}
}

func TestLoadSqlFile(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Close()

	err = loadSqlFile(db, "database/sql/create_users_table.sql")
	if err != nil {
		t.Errorf("loadSqlFile failed: %v", err)
	}
}

func TestInitDb(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Close()

	err = initDb(db)
	if err != nil {
		t.Errorf("initDb failed: %v", err)
	}
}

func TestDatabaseConnection(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Errorf("Database connection ping failed: %v", err)
	}
}
