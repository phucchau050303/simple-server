package main

import (
	"log"
	"net"
	"os"
	"simple-webser/internal/credential"
	"simple-webser/internal/storage/postgres"
	credentialv1 "simple-webser/proto/credential/v1"

	"google.golang.org/grpc"
)

func main() {
	db, err := postgres.Connect(os.Getenv("CONNECTION_STRING"))

	if err != nil {
		log.Fatal(err)
	}
	userRepository := postgres.NewUserRepository(db)
	credentialService := credential.NewService(userRepository)
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(credential.RequestIDInterceptor))
	credentialv1.RegisterCredentialServiceServer(grpcServer, credential.NewServer(credentialService))

	log.Println("Credential Service listening on :50051")

	grpcServer.Serve(listener)

}
