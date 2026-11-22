// @title Registration API
// @version 1.0
// @description API for the registration server.
// @host localhost:8080
// @BasePath /
package main

import (
	"log"
	"net/http"
	_ "simple-webser/docs"
	"simple-webser/internal/registration"
	_ "simple-webser/internal/storage/postgres"
	credentialv1 "simple-webser/proto/credential/v1"

	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func requestLogger(r *http.Request) *log.Logger {
	return registration.RequestLogger(r)
}

func main() {
	router := http.NewServeMux()
	conn, err := grpc.NewClient("credential:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to credential service: %v", err)
	}
	defer conn.Close()

	credentialClient := credentialv1.NewCredentialServiceClient(conn)

	registrationService := registration.NewService(credentialClient)

	router.Handle("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// healthHandler godoc
	// @Summary Health check
	// @Produce json
	// @Success 200 {object} HealthResponse
	// @Router /health [get]
	router.HandleFunc("GET /health", registration.HealthHandler)

	router.HandleFunc("POST /signin", registration.SignInHandler(registrationService))

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
	router.HandleFunc("POST /signup", registration.SignUpHandler(registrationService))

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", registration.RequestIDMiddleware(router)))
}
