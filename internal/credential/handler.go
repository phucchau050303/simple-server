package credential

import (
	"context"

	credentialv1 "simple-webser/proto/credential/v1"
)

type Server struct {
	credentialv1.UnimplementedCredentialServiceServer

	service Service
}

func NewServer(service Service) *Server {
	return &Server{service: service}
}

func (s *Server) SignIn(ctx context.Context, req *credentialv1.SignInRequest) (*credentialv1.SignInResponse, error) {
	return s.service.SignIn(ctx, req.Email, req.Password)
}

func (s *Server) SignUp(ctx context.Context, req *credentialv1.SignUpRequest) (*credentialv1.SignUpResponse, error) {
	return s.service.SignUp(ctx, req.Email, req.Password)
}
