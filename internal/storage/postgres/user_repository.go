package postgres

import (
	"context"
	"database/sql"

	"simple-webser/internal/user"
)

var _ user.Repository = (*UserRepository)(nil)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindCredentials(ctx context.Context, email string) (int, string, error) {
	var (
		id           int
		passwordHash string
	)

	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, passwordhash FROM users WHERE email = $1",
		email,
	).Scan(&id, &passwordHash)
	if err != nil {
		return 0, "", err
	}

	return id, passwordHash, nil
}

func (r *UserRepository) Exists(ctx context.Context, email string) (bool, error) {
	var id int

	err := r.db.QueryRowContext(
		ctx,
		"SELECT id FROM users WHERE email = $1",
		email,
	).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *UserRepository) Create(ctx context.Context, email string, passwordHash string) (int, error) {
	var id int

	err := r.db.QueryRowContext(
		ctx,
		"INSERT INTO users (email, passwordhash, createdat) VALUES ($1, $2, NOW()) RETURNING id",
		email,
		passwordHash,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
