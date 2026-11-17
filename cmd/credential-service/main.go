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
	//query to check if email exist
	err := s.db.QueryRowContext(ctx, "SELECT id, passwordhash FROM users WHERE email = $1", email).Scan(&id, &storedHash)
	if err == sql.ErrNoRows {
		//if no user found, return error
		return &credentialv1.SignInResponse{
			IsValid: false,
			Message: "user not found",
		}, nil
	}
	// if different error, return that error
	if err != nil {
		return nil, err
	}

	//call function to verify password and hash
	match, err := credential.Verify(password, storedHash)
	if err != nil {
		return nil, err
	}

	//if match false, return password not matched
	if match == false {
		return &credentialv1.SignInResponse{
			IsValid: false,
			Message: "Password not match",
		}, nil
	} else {
		//else return success
		return &credentialv1.SignInResponse{
			IsValid: true,
			UserId:  fmt.Sprint(id),
			Message: "Sucessfully logged in",
		}, nil
	}
}

func (s *server) SignUp(ctx context.Context, req *credentialv1.SignUpRequest) (*credentialv1.SignUpResponse, error) {
	email := req.Email
	password := req.Password
	var id int

	//query to check if email exist
	err := s.db.QueryRowContext(ctx, "SELECT id FROM users WHERE email = $1", email).Scan(&id)
	if err == nil {
		//if yes, return user already registered
		return &credentialv1.SignUpResponse{
			IsValid: false,
			Message: "User already registered",
		}, err
	}
	//if different error, return that error
	if err != sql.ErrNoRows {
		return nil,
			fmt.Errorf("Unexpected error: %v", err)
	}

	//call function to hash password
	passwordHash, err := credential.Hash(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	//insert data into postgres db
	_, err = s.db.ExecContext(ctx, "INSERT INTO users (email, passwordhash, createdat) VALUES ($1, $2, NOW())", email, passwordHash)
	// if error  thrown, return error
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	//return sucess
	return &credentialv1.SignUpResponse{
		IsValid: true,
		Message: "user created sucessfully",
	}, nil
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
