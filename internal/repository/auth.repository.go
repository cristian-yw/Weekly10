package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	DB *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{DB: db}
}

// Register user baru
func (r *AuthRepository) RegisterUser(email, passwordHash string) error {
	now := time.Now()
	_, err := r.DB.Exec(
		context.Background(),
		`INSERT INTO users (email, password_hash, role, created_at, updated_at) 
		 VALUES ($1, $2, 'user', $3, $3)`,
		email, passwordHash, now,
	)
	return err
}

// Cari user berdasarkan email
func (r *AuthRepository) GetUserByEmail(email string) (int, string, string, error) {
	var id int
	var passwordHash, role string

	err := r.DB.QueryRow(
		context.Background(),
		"SELECT id, password_hash, role FROM users WHERE email=$1",
		email,
	).Scan(&id, &passwordHash, &role)

	return id, passwordHash, role, err
}
