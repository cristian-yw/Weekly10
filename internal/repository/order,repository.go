package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/cristian-yw/Weekly10/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	DB *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{DB: db}
}

// 1. Get Schedule
func (r *OrderRepository) GetSchedule(
	ctx context.Context,
	movieID int,
	cinemaName, locationName string,
	startTime, date string, // gunakan string/nullable sesuai kebutuhan
) ([]models.Schedule, error) {

	// Base query
	query := `
		SELECT s.id, m.title, c.name, l.location, t.start_time, s.date
		FROM schedules s
		JOIN movies m ON m.id = s.movie_id
		JOIN cinemas c ON c.id = s.cinema_id
		JOIN locations l ON l.id = s.location_id
		JOIN times t ON t.id = s.time_id
		WHERE m.id = $1`
	args := []interface{}{movieID}
	argPos := 2

	// Tambah filter bila diisi
	if cinemaName != "" {
		query += fmt.Sprintf(" AND c.name ILIKE $%d", argPos)
		args = append(args, "%"+cinemaName+"%")
		argPos++
	}
	if locationName != "" {
		query += fmt.Sprintf(" AND l.location ILIKE $%d", argPos)
		args = append(args, "%"+locationName+"%")
		argPos++
	}
	if startTime != "" {
		query += fmt.Sprintf(" AND t.start_time = $%d", argPos)
		args = append(args, startTime)
		argPos++
	}
	if date != "" {
		query += fmt.Sprintf(" AND s.date = $%d", argPos)
		args = append(args, date)
		argPos++
	}

	query += " ORDER BY s.date ASC"

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var s models.Schedule
		if err := rows.Scan(&s.ID, &s.MovieTitle, &s.Cinema, &s.Location, &s.StartTime, &s.Date); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return schedules, nil
}

// 2. Get Available Seat
func (r *OrderRepository) GetAvailableSeats(ctx context.Context, scheduleID int) ([]models.Seat, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT DISTINCT s.id, s.cinema_id, s.seat_code, true as is_booked
FROM seats s
JOIN order_seats os ON s.seat_code = os.seat_code
JOIN orders o ON o.id = os.order_id
JOIN schedules sch ON sch.id = o.schedule_id
WHERE o.schedule_id = $1
  AND o.status = 'paid'
  AND s.cinema_id = sch.cinema_id
ORDER BY s.seat_code;

	`, scheduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []models.Seat
	for rows.Next() {
		var seat models.Seat
		if err := rows.Scan(&seat.ID, &seat.CinemaID, &seat.SeatCode, &seat.IsBooked); err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}

	return seats, nil
}

// 3. Get Movie Detail
func (r *OrderRepository) GetMovieDetail(ctx context.Context, movieID int) (*models.MovieDetail, error) {
	row := r.DB.QueryRow(ctx, `
	SELECT
		m.id,
		m.title,
		m.overview,
		TO_CHAR(m.release_date, 'YYYY-MM-DD') AS release_date,
		COALESCE(m.runtime, 0) AS runtime,
		COALESCE(array_agg(DISTINCT g.name) FILTER (WHERE g.name IS NOT NULL), '{}') AS genres,
		COALESCE(array_agg(DISTINCT p.name) FILTER (WHERE c.role='Actor' AND p.name IS NOT NULL), '{}') AS casts,
		COALESCE(d.name, '') AS director,
		COALESCE(m.poster_path, '') AS poster_path,
		COALESCE(m.backdrop_path, '') AS backdrop_path
	FROM movies m
	LEFT JOIN movie_genres mg ON mg.movie_id = m.id
	LEFT JOIN genres g ON g.id = mg.genre_id
	LEFT JOIN movie_casts c ON c.movie_id = m.id
	LEFT JOIN persons p ON p.id = c.person_id
	LEFT JOIN movie_casts md ON md.movie_id = m.id AND md.role = 'Director'
	LEFT JOIN persons d ON d.id = md.person_id
	WHERE m.id = $1
	GROUP BY m.id, d.name
	`, movieID)

	var md models.MovieDetail

	// Pastikan struct MovieDetail memiliki:
	// ID int
	// Title string
	// Overview string
	// ReleaseDate string
	// Runtime int
	// Genres []string
	// Casts []string
	// Director string
	// PosterPath string
	// BackdropPath string
	if err := row.Scan(
		&md.ID,
		&md.Title,
		&md.Overview,
		&md.ReleaseDate,
		&md.Runtime,
		&md.Genres,
		&md.Casts,
		&md.Director,
		&md.PosterPath,
		&md.BackdropPath,
	); err != nil {
		return nil, err
	}

	return &md, nil
}

// 4. Create Order
func (r *OrderRepository) CreateOrder(ctx context.Context, userID, scheduleID, totalPrice int, seats []string) (*models.Order, error) {
	tx, err := r.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// 1. Insert ke orders
	var orderID int
	err = tx.QueryRow(ctx, `
		INSERT INTO orders (user_id, schedule_id, total_price, status, order_date)
		VALUES ($1, $2, $3, 'paid', NOW())
		RETURNING id
	`, userID, scheduleID, totalPrice).Scan(&orderID)
	if err != nil {
		return nil, err
	}

	// 2. Insert ke order_seats
	for _, seatCode := range seats {
		_, err := tx.Exec(ctx, `
			INSERT INTO order_seats (order_id, seat_code)
			VALUES ($1, $2)
		`, orderID, seatCode)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &models.Order{
		ID:         orderID,
		UserID:     userID,
		ScheduleID: scheduleID,
		TotalPrice: totalPrice,
		Status:     "paid",
		OrderDate:  time.Now(),
		Seats:      seats,
	}, nil
}
