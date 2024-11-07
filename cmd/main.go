package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"to-do/migrations"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const port = ":8080"

// Завантаження .env файлу
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func jsonResponse(key, message string) ([]byte, error) {
	jsonResponse := simplejson.New()
	jsonResponse.Set(key, message)
	return jsonResponse.MarshalJSON()
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := jsonResponse("status", "Server is running")
	if err != nil {
		log.Printf("Error generating JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(payload); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/ping", pingHandler).Methods("GET")
	return r
}

func main() {
	// Формування DSN для GORM
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DBUSER"),
		os.Getenv("DBPASS"),
		os.Getenv("DBHOST"),
		os.Getenv("DBPORT"),
		os.Getenv("DBNAME"),
	)

	// Підключення до бази даних за допомогою GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// Виконання міграцій
	if err := migrations.Migrate(db); err != nil {
		log.Fatal("Migration failed:", err)
	}

	fmt.Println("Connected to the database and migrations applied!")

	// Запуск сервера
	if err := http.ListenAndServe(port, router()); err != nil {
		log.Fatal(err)
	}
}
