// @title Registration API
// @version 1.0
// @description API for the registration server.
// @host localhost:8080
// @BasePath /
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "simple-webser/docs"
	_ "simple-webser/internal/storage/postgres"
	credentialv1 "simple-webser/proto/credential/v1"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// healthHandler godoc
// @Summary Health check
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
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
	conn, err := grpc.NewClient("credential:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to credential service: %v", err)
	}
	defer conn.Close()

	credentialClient := credentialv1.NewCredentialServiceClient(conn)

	router.Handle("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	router.HandleFunc("GET /health", healthHandler)
	// SignInHandler godoc
	// @Summary      User Sign In
	// @Description  Authenticate user with email and password via gRPC
	// @Tags         auth
	// @Accept       json
	// @Produce      json
	// @Param        request body reqBody true "Login Credentials"
	// @Success      200 {object} map[string]string
	// @Failure      400,401,500 {string} string
	// @Router /signin [post]
	router.HandleFunc("POST /signin", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var reqBody struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		grpcResponse, err := credentialClient.SignIn(r.Context(), &credentialv1.SignInRequest{
			Email:    reqBody.Email,
			Password: reqBody.Password,
		})

		if err != nil {
			log.Printf("❌ gRPC call failed: %v", err)
			http.Error(w, fmt.Sprintf("internal service error: %v", err),
				http.StatusInternalServerError)
			return
		}

		if !grpcResponse.GetIsValid() {
			http.Error(w, grpcResponse.GetMessage(), http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "message":"succesfully logged in"}`))
	})
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
