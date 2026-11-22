package registration

import (
	"context"
	credentialv1 "simple-webser/proto/credential/v1"
)

type CredentialService interface {
	SignUp(ctx context.Context, email string, password string) (*credentialv1.SignUpResponse, error)
	SignIn(ctx context.Context, email string, password string) (*credentialv1.SignInResponse, error)
}

type service struct {
	credentialClient credentialv1.CredentialServiceClient
}

func NewService(client credentialv1.CredentialServiceClient) CredentialService {
	return &service{
		credentialClient: client,
	}
}

func (s *service) SignUp(
	ctx context.Context,
	email string,
	password string,
) (*credentialv1.SignUpResponse, error) {
	return s.credentialClient.SignUp(
		ctx,
		&credentialv1.SignUpRequest{
			Email:    email,
			Password: password,
		},
	)
}

func (s *service) SignIn(
	ctx context.Context,
	email string,
	password string,
) (*credentialv1.SignInResponse, error) {
	return s.credentialClient.SignIn(
		ctx,
		&credentialv1.SignInRequest{
			Email:    email,
			Password: password,
		},
	)
}
