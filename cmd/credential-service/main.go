package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"simple-webser/internal/credential"
	"simple-webser/internal/storage/postgres"
	credentialv1 "simple-webser/proto/credential/v1"

	"google.golang.org/grpc"
)

type server struct {
	// Embed grpc base struct
	credentialv1.UnimplementedCredentialServiceServer

	// database connection pool
	db *sql.DB
}

func (s *server) SignIn(ctx context.Context, req *credentialv1.SignInRequest) (*credentialv1.SignInResponse, error) {
	email := req.Email
	password := req.Password
	var (
		id         int
		storedHash string
	)
	err := s.db.QueryRowContext(ctx, "SELECT id, passwordhash FROM users WHERE email = $1", email).Scan(&id, &storedHash)
	if err == sql.ErrNoRows {
		return &credentialv1.SignInResponse{
			IsValid: false,
			Message: "user not found",
		}, nil
	}

	if err != nil {
		return nil, err
	}

	match, err := credential.Verify(password, storedHash)
	if err != nil {
		return nil, err
	}

	if match == false {
		return &credentialv1.SignInResponse{
			IsValid: false,
			Message: "Password not match",
		}, nil
	} else {
		return &credentialv1.SignInResponse{
			IsValid: true,
			UserId:  fmt.Sprint(id),
			Message: "Sucessfully logged in",
		}, nil
	}
}

func main() {
	db, err := postgres.Connect(os.Getenv("CONNECTION_STRING"))
	//db, err := postgres.Connect("postgres://peter:phcchu88@db:5432/registration?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()
	credentialv1.RegisterCredentialServiceServer(grpcServer, &server{db: db})

	log.Println("Credential Service listening on :50051")

	grpcServer.Serve(listener)
	// for {
	// 	conn, err := listener.Accept()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	go func() {
	// 		defer conn.Close()
	// 		// Handle the connection here.
	// 	}()
	// }

}
