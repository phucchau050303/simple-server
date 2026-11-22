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

	router.Handle("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	router.HandleFunc("GET /health", registration.HealthHandler)
	router.HandleFunc("POST /signin", registration.SignInHandler(credentialClient))
	router.HandleFunc("POST /signup", registration.SignUpHandler(credentialClient))

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", registration.RequestIDMiddleware(router)))
}
