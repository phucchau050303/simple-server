package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) { //function must have w http.ResponseWrite and r *http.Request
	w.Header().Set("Content-Type", "application/json") // Set the Header - Content type to application/json

	// w.WriteHeader(http.StatusOK) // Set the status code

	response := HealthResponse{ // Initialise response
		Status:    "OK",
		Timestamp: time.Now(),
	}

	json.NewEncoder(w).Encode(response) // Encode the response to w json

}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /health", healthHandler)
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
