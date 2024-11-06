package main

import (
	"github.com/bitly/go-simplejson"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testJsonResponse(t *testing.T) {
	expectedKey := "key"
	expectedMessage := "Server is running"
	payload, err := jsonResponse(expectedKey, expectedMessage)
	if err != nil {
		log.Printf("Error generating JSON response: %v", err)
	}
	json, err := simplejson.NewJson(payload)
	if err != nil {
		t.Fatalf("Error getting key '%s' from JSON: %v", expectedKey, err)
	}
	actualMessage, err := json.Get(expectedKey).String()
	if err != nil {
		t.Fatalf("Error getting key '%s' from JSON: %v", expectedKey, err)
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
