package repository

import (
	"context"
	"time"

	"github.com/cristian-yw/Weekly10/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepository struct {
	DB *pgxpool.Pool
}

func NewAdminRepository(db *pgxpool.Pool) *AdminRepository {
	return &AdminRepository{DB: db}
}

// Get all movies
func (r *AdminRepository) GetAllMovies(ctx context.Context) ([]models.Movie, error) {
	rows, err := r.DB.Query(ctx, `SELECT id, title, release_date, runtime, created_at, updated_at FROM movies ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.ReleaseDate, &m.Runtime, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, nil
}

// Edit movie
func (r *AdminRepository) UpdateMovie(ctx context.Context, id int, title string, runtime int) error {
	_, err := r.DB.Exec(ctx, `UPDATE movies SET title=$1, runtime=$2, updated_at=$3 WHERE id=$4`, title, runtime, time.Now(), id)
	return err
}

// Delete movie
func (r *AdminRepository) DeleteMovie(ctx context.Context, id int) error {
	_, err := r.DB.Exec(ctx, `DELETE FROM movies WHERE id=$1`, id)
	return err
}
