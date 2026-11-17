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

type SignUpRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123!"`
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

// SignUpHandler godoc
// @Summary      User Sign Up
// @Description  Register new user via gRPC
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body SignUpRequest true "Registration Credentials"
// @Success      201 {object} map[string]string
// @Failure      400 {string} string "Bad Request"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /signup [post]
func signUpHandler(credentialClient credentialv1.CredentialServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json") //set content type in header

		//check if right method is sent
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		//initialise request body struct
		var reqBody SignUpRequest

		//decode request body from front end
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		//call Sign Up from credential service
		grpcResponse, err := credentialClient.SignUp(r.Context(), &credentialv1.SignUpRequest{
			Email:    reqBody.Email,
			Password: reqBody.Password,
		})

		//return error if grpc call fails
		if err != nil {
			http.Error(w, fmt.Sprintf("gRPC call failed, internal service error: %v", err), http.StatusInternalServerError)
			return
		}

		//return error if response is invalid
		if !grpcResponse.GetIsValid() {
			http.Error(w, grpcResponse.GetMessage(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"status": "created", "message":" user created succesfully"}`))
	}
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

	router.HandleFunc("POST /signin", func(w http.ResponseWriter, r *http.Request) {
		//check if the right method is sent
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		//declare and initialise request body struct
		var reqBody struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		//decode the request body from front end and check if it is valid
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		//call SignIn funciton from credential service
		grpcResponse, err := credentialClient.SignIn(r.Context(), &credentialv1.SignInRequest{
			Email:    reqBody.Email,
			Password: reqBody.Password,
		})

		// if error is thrown, return error
		if err != nil {
			log.Printf("gRPC call failed: %v", err)
			http.Error(w, fmt.Sprintf("internal service error: %v", err),
				http.StatusInternalServerError)
			return
		}

		//if response is invalid, return 401 unauthorised
		if !grpcResponse.GetIsValid() {
			http.Error(w, grpcResponse.GetMessage(), http.StatusUnauthorized)
			return
		}

		//set content type of the header
		w.Header().Set("Content-Type", "application/json")
		//set 200 OK status
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "message":"succesfully logged in"}`))
	})

	router.HandleFunc("POST /signup", signUpHandler(credentialClient))

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
