package main

import (
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const port = ":8080"

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
	if err := http.ListenAndServe(port, router()); err != nil {
		log.Fatal(err)
	}
}
