package credential

import (
	"context"
	"database/sql"
	"fmt"

	"simple-webser/internal/user"
	credentialv1 "simple-webser/proto/credential/v1"
)

type Service interface {
	SignIn(ctx context.Context, email string, password string) (*credentialv1.SignInResponse, error)
	SignUp(ctx context.Context, email string, password string) (*credentialv1.SignUpResponse, error)
}

type service struct {
	users user.Repository
}

func NewService(users user.Repository) Service {
	return &service{users: users}
}

func (s *service) SignIn(ctx context.Context, email string, password string) (*credentialv1.SignInResponse, error) {
	id, storedHash, err := s.users.FindCredentials(ctx, email)
	if err == sql.ErrNoRows {
		return &credentialv1.SignInResponse{
			IsValid: false,
			Message: "user not found",
		}, nil
	}
	if err != nil {
		return nil, err
	}

	match, err := Verify(password, storedHash)
	if err != nil {
		return nil, err
	}
	if !match {
		return &credentialv1.SignInResponse{
			IsValid: false,
			Message: "Password not match",
		}, nil
	}

	return &credentialv1.SignInResponse{
		IsValid: true,
		UserId:  fmt.Sprint(id),
		Message: "Sucessfully logged in",
	}, nil
}

func (s *service) SignUp(ctx context.Context, email string, password string) (*credentialv1.SignUpResponse, error) {
	if email == "" || password == "" {
		return &credentialv1.SignUpResponse{
			IsValid: false,
			Message: "Email and password cannot be null!",
		}, fmt.Errorf("email and password cannot be null")
	}
	if FormatCheck(password) == false {
		return &credentialv1.SignUpResponse{
			IsValid: false,
			Message: "Password is not in correct format!",
		}, nil
	}

	exists, err := s.users.Exists(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user: %w", err)
	}
	if exists {
		return &credentialv1.SignUpResponse{
			IsValid: false,
			Message: "User already registered",
		}, nil
	}

	passwordHash, err := Hash(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	if _, err := s.users.Create(ctx, email, passwordHash); err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return &credentialv1.SignUpResponse{
		IsValid: true,
		Message: "user created sucessfully",
	}, nil
}
