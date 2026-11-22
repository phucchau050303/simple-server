// @title Registration API
// @version 1.0
// @description API for the registration server.
// @host localhost:8080
// @BasePath /
package main

import (
	"crypto/rand"
	"encoding/hex"
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

func newRequestID() string {
	id := make([]byte, 8)
	if _, err := rand.Read(id); err != nil {
		return fmt.Sprintf("fallback-%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(id)
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := newRequestID()
		w.Header().Set("X-Request-ID", requestID)
		r.Header.Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r)
	})
}

func requestLogger(r *http.Request) *log.Logger {
	return log.New(log.Writer(), fmt.Sprintf("[request_id=%s] ", r.Header.Get("X-Request-ID")), log.Flags())
}

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
		requestLog := requestLogger(r)
		w.Header().Set("Content-Type", "application/json") //set content type in header

		//check if right method is sent
		requestLog.Println("checking if POST method is sent")
		if r.Method != http.MethodPost {
			//http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			w.WriteHeader(http.StatusMethodNotAllowed)
			response := map[string]string{
				"error": "method not allowed",
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		//initialise request body struct
		var reqBody SignUpRequest

		//decode request body from front end
		requestLog.Println("Decoding request body from front end")
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			//http.Error(w, "invalid request body", http.StatusBadRequest)
			requestLog.Println("Decoding failed, invalid request body")
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]string{
				"error": "invalid request body",
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		//log.Logger(hthtp server receive smth)
		//call Sign Up from credential service
		requestLog.Println("Calling SignUp service through gRPC to process and insert data")
		grpcResponse, err := credentialClient.SignUp(r.Context(), &credentialv1.SignUpRequest{
			Email:    reqBody.Email,
			Password: reqBody.Password,
		})

		//return error if grpc call fails
		if err != nil {
			requestLog.Printf("gRPC call failed: %v", err)
			//http.Error(w, fmt.Sprintf("gRPC call failed, internal service error: %v", err), http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			response := map[string]string{
				"error": err.Error(),
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		//return error if response is invalid
		if !grpcResponse.GetIsValid() {
			requestLog.Printf("Response from gRPC is invalid: %v", err)
			//http.Error(w, grpcResponse.GetMessage(), http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]string{
				"error": grpcResponse.GetMessage(),
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message":" user created succesfully"}`))
		requestLog.Println("User created succesfully")
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
		requestLog := requestLogger(r)
		w.Header().Set("Content-Type", "application/json")

		//check if the right method is sent
		if r.Method != http.MethodPost {
			//http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			w.WriteHeader(http.StatusMethodNotAllowed)
			response := map[string]string{
				"error": "method not allowed",
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		//declare and initialise request body struct
		var reqBody struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		//decode the request body from front end and check if it is valid
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			//http.Error(w, "invalid request body", http.StatusBadRequest)
			requestLog.Println("Decoding failed, invalid request body")
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]string{
				"error": "invalid request body",
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		//call SignIn funciton from credential service
		grpcResponse, err := credentialClient.SignIn(r.Context(), &credentialv1.SignInRequest{
			Email:    reqBody.Email,
			Password: reqBody.Password,
		})

		// if error is thrown, return error
		if err != nil {
			requestLog.Printf("gRPC call failed: %v", err)
			//http.Error(w, fmt.Sprintf("gRPC call failed, internal service error: %v", err), http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			response := map[string]string{
				"error": err.Error(),
			}
			json.NewEncoder(w).Encode(response)
		}

		//if response is invalid, return 401 unauthorised
		if !grpcResponse.GetIsValid() {
			requestLog.Printf("Response from gRPC is invalid: %v", err)
			//http.Error(w, grpcResponse.GetMessage(), http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]string{
				"error": grpcResponse.GetMessage(),
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		//set content type of the header

		//set 200 OK status
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"succesfully logged in"}`))
	})

	router.HandleFunc("POST /signup", signUpHandler(credentialClient))

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", requestIDMiddleware(router)))
}
