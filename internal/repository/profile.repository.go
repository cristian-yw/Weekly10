package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/cristian-yw/Weekly10/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetProfile(ctx context.Context, userID int) (models.UserProfile, error) {
	var profile models.UserProfile
	err := r.DB.QueryRow(ctx, `
		SELECT u.id, u.email, p.first_name, p.last_name, p.phone, p.avatar_url
		FROM users u
		JOIN profiles p ON p.user_id = u.id
		WHERE u.id=$1
	`, userID).Scan(&profile.ID, &profile.Email, &profile.FirstName, &profile.LastName, &profile.Phone, &profile.AvatarURL)
	return profile, err
}

func (r *UserRepository) GetHistory(ctx context.Context, userID int) ([]models.OrderHistory, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT o.id, m.title, o.total_price, o.status, o.order_date
		FROM orders o
		JOIN schedules s ON s.id = o.schedule_id
		JOIN movies m ON m.id = s.movie_id
		WHERE o.user_id=$1
		ORDER BY o.order_date DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.OrderHistory
	for rows.Next() {
		var h models.OrderHistory
		if err := rows.Scan(&h.OrderID, &h.MovieTitle, &h.TotalPrice, &h.Status, &h.Date); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, nil
}
func (r *UserRepository) UpdateProfile(
	ctx context.Context,
	id int,
	firstName, lastName, phone *string,
	avatarURL *string,
) error {
	setParts := []string{}
	args := []interface{}{}
	i := 1

	if firstName != nil {
		setParts = append(setParts, fmt.Sprintf("first_name = $%d", i))
		args = append(args, *firstName)
		i++
	}
	if lastName != nil {
		setParts = append(setParts, fmt.Sprintf("last_name = $%d", i))
		args = append(args, *lastName)
		i++
	}
	if phone != nil {
		setParts = append(setParts, fmt.Sprintf("phone = $%d", i))
		args = append(args, *phone)
		i++
	}
	if avatarURL != nil {
		setParts = append(setParts, fmt.Sprintf("avatar_url = $%d", i))
		args = append(args, *avatarURL)
		i++
	}

	// Jika tidak ada field yang dikirim, jangan update apa pun
	if len(setParts) == 0 {
		return nil
	}

	// selalu update updated_at
	setParts = append(setParts, fmt.Sprintf("updated_at = NOW()"))

	query := fmt.Sprintf(`UPDATE users SET %s WHERE id = $%d`,
		strings.Join(setParts, ", "), i)
	args = append(args, id)

	_, err := r.DB.Exec(ctx, query, args...)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.UserProfile, error) {
	query := `
		SELECT id, email, first_name, last_name, phone, avatar_url
		FROM users
		WHERE id = $1
	`
	u := &models.UserProfile{}
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.Email,
		&u.FirstName,
		&u.LastName,
		&u.Phone,
		&u.AvatarURL,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Ambil hash password asli
func (r *UserRepository) GetPasswordHash(ctx context.Context, userID int) (string, error) {
	var hash string
	err := r.DB.QueryRow(ctx,
		"SELECT password_hash FROM users WHERE id = $1",
		userID,
	).Scan(&hash)
	return hash, err
}

// Update password
func (r *UserRepository) UpdatePassword(ctx context.Context, userID int, newHashed string) error {
	_, err := r.DB.Exec(ctx,
		"UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2",
		newHashed, userID,
	)
	return err
}
