package main

import (
	"database/sql"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const port = ":8080"

func loadSqlFile(db *sql.DB, filePath string) error {
	query, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("Error reading file: %v", err)
	}

	if _, err := db.Exec(string(query)); err != nil {
		return fmt.Errorf("Error executing query: %v", err)
	}

	fmt.Printf("Successfully loaded file %s\n", filePath)
	return nil
}
func initDb(db *sql.DB) error {
	files := []string{"database/sql/create_users_table.sql",
		"database/sql/create_records_table.sql",
	}
	for _, file := range files {
		if err := loadSqlFile(db, file); err != nil {
			return err
		}
	}
	return nil
}
func init() {
	// Завантаження .env файлу

	err := godotenv.Load()
	if err != nil {
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
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%s", os.Getenv("DBHOST"), os.Getenv("DBPORT")),
		DBName: os.Getenv("DBNAME"),
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // Закриваємо з'єднання при завершенні програми

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	if err := initDb(db); err != nil {
		log.Fatal("Error initializing the database:", err)
	}
	if err := http.ListenAndServe(port, router()); err != nil {
		log.Fatal(err)
	}
}
