package repository

import (
	"context"

	"github.com/cristian-yw/Weekly10/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ScheduleRepository struct {
	DB *pgxpool.Pool
}

func NewScheduleRepository(db *pgxpool.Pool) *ScheduleRepository {
	return &ScheduleRepository{DB: db}
}

func (r *ScheduleRepository) GetGenres(ctx context.Context) ([]models.Genre, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, name FROM genres")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var g models.Genre
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, err
		}
		genres = append(genres, g)
	}
	return genres, nil
}

func (r *ScheduleRepository) GetCinemas(ctx context.Context) ([]models.Cinema, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, name FROM cinemas")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cinemas []models.Cinema
	for rows.Next() {
		var c models.Cinema
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		cinemas = append(cinemas, c)
	}
	return cinemas, nil
}

func (r *ScheduleRepository) GetLocations(ctx context.Context) ([]models.Location, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, location FROM locations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []models.Location
	for rows.Next() {
		var l models.Location
		if err := rows.Scan(&l.ID, &l.Name); err != nil {
			return nil, err
		}
		locations = append(locations, l)
	}
	return locations, nil
}
