package user

import "context"

type Repository interface {
	FindCredentials(ctx context.Context, email string) (int, string, error)
	Exists(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, email string, passwordHash string) (int, error)
}
