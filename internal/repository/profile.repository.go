package repository

import (
	"context"

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
func (r *UserRepository) EditProfile(ctx context.Context, userID int, firstName, lastName, phone, avatarURL string) (models.UserProfile, error) {
	_, err := r.DB.Exec(ctx, `
		UPDATE profiles
		SET first_name=$1, last_name=$2, phone=$3, avatar_url=$4, updated_at=NOW()
		WHERE user_id=$5
	`, firstName, lastName, phone, avatarURL, userID)
	if err != nil {
		return models.UserProfile{}, err
	}
	return r.GetProfile(ctx, userID)
}
