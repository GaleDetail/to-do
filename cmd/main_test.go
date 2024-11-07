package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"to-do/migrations"
	"to-do/models"
)

var db *gorm.DB

// TestMain виконує налаштування тестового середовища
func TestMain(m *testing.M) {
	// Завантаження .env файлу
	if err := godotenv.Load(); err != nil {
		log.Println("Could not load .env file for tests:", err)
	}

	// Підключення до тестової бази даних
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DBUSER"),
		os.Getenv("DBPASS"),
		os.Getenv("DBHOST"),
		os.Getenv("DBPORT"),
		os.Getenv("DBNAME"),
	)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// Виконання міграцій для тестової бази даних
	if err := migrations.Migrate(db); err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Запуск тестів
	code := m.Run()
	os.Exit(code)
}

// Тест для створення нової задачі
func TestCreateTask(t *testing.T) {
	task := models.Task{
		Title:       "Test Task",
		Description: "This is a test description",
		Status:      "pending",
	}

	if err := db.Create(&task).Error; err != nil {
		t.Fatalf("Could not create task: %v", err)
	}

	// Перевіряємо, чи було присвоєно ID
	if task.ID == 0 {
		t.Errorf("Expected task ID to be set, got 0")
	}
}

// Тест для отримання задачі
func TestGetTask(t *testing.T) {
	// Створюємо задачу для тесту
	task := models.Task{
		Title:       "Fetch Task",
		Description: "Task to fetch",
		Status:      "in_progress",
	}

	if err := db.Create(&task).Error; err != nil {
		t.Fatalf("Could not create task: %v", err)
	}

	var fetchedTask models.Task
	if err := db.First(&fetchedTask, task.ID).Error; err != nil {
		t.Fatalf("Could not fetch task: %v", err)
	}

	if fetchedTask.Title != task.Title {
		t.Errorf("Expected title %s, got %s", task.Title, fetchedTask.Title)
	}
}

// Тест для оновлення задачі
func TestUpdateTask(t *testing.T) {
	// Створюємо задачу для оновлення
	task := models.Task{
		Title:       "Update Task",
		Description: "Task to update",
		Status:      "pending",
	}

	if err := db.Create(&task).Error; err != nil {
		t.Fatalf("Could not create task: %v", err)
	}

	// Оновлюємо статус задачі
	if err := db.Model(&task).Update("Status", "completed").Error; err != nil {
		t.Fatalf("Could not update task: %v", err)
	}

	// Перевіряємо, чи оновлено статус
	var updatedTask models.Task
	if err := db.First(&updatedTask, task.ID).Error; err != nil {
		t.Fatalf("Could not fetch updated task: %v", err)
	}

	if updatedTask.Status != "completed" {
		t.Errorf("Expected status 'completed', got %s", updatedTask.Status)
	}
}

// Тест для видалення задачі
func TestDeleteTask(t *testing.T) {
	// Створюємо задачу для видалення
	task := models.Task{
		Title:       "Delete Task",
		Description: "Task to delete",
		Status:      "pending",
	}

	if err := db.Create(&task).Error; err != nil {
		t.Fatalf("Could not create task: %v", err)
	}

	// Видаляємо задачу
	if err := db.Delete(&task).Error; err != nil {
		t.Fatalf("Could not delete task: %v", err)
	}

	// Перевіряємо, чи задача видалена
	var deletedTask models.Task
	if err := db.First(&deletedTask, task.ID).Error; err == nil {
		t.Errorf("Expected task to be deleted, but it was found")
	}
}

// Тест для HTTP-обробника /ping
func TestPingHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(pingHandler)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":"Server is running"}`
	if rec.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rec.Body.String(), expected)
	}
}
