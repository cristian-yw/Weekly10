package repository

import (
	"context"
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
func (r *OrderRepository) GetSchedule(ctx context.Context, movieID int) ([]models.Schedule, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT s.id, m.title, c.name, l.location, t.start_time, s.date
		FROM schedules s
		JOIN movies m ON m.id = s.movie_id
		JOIN cinemas c ON c.id = s.cinema_id
		JOIN locations l ON l.id = s.location_id
		JOIN times t ON t.id = s.time_id
		WHERE m.id = $1
		ORDER BY s.date ASC`, movieID)
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
	return schedules, nil
}

// 2. Get Available Seat
func (r *OrderRepository) GetAvailableSeats(ctx context.Context, scheduleID int) ([]models.Seat, error) {
	// Ambil semua seat untuk cinema jadwal ini
	rows, err := r.DB.Query(ctx, `
		SELECT s.id, s.cinema_id, s.seat_code,
		       CASE 
		         WHEN os.seat_code IS NULL THEN false
		         ELSE true
		       END AS is_booked
		FROM seats s
		LEFT JOIN (
		    SELECT os.seat_code
		    FROM order_seats os
		    JOIN orders o ON o.id = os.order_id
		    WHERE o.schedule_id = $1 AND o.status = 'paid'
		) os ON s.seat_code = os.seat_code
		ORDER BY s.seat_code
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
		SELECT m.id, m.title, m.overview, TO_CHAR(m.release_date, 'YYYY-MM-DD'),
		       m.runtime,
		       array_agg(DISTINCT g.name) AS genres,
		       array_agg(DISTINCT p.name) AS casts
		FROM movies m
		LEFT JOIN movie_genres mg ON mg.movie_id = m.id
		LEFT JOIN genres g ON g.id = mg.genre_id
		LEFT JOIN movie_casts mc ON mc.movie_id = m.id
		LEFT JOIN persons p ON p.id = mc.person_id
		WHERE m.id = $1
		GROUP BY m.id`, movieID)

	var md models.MovieDetail
	if err := row.Scan(&md.ID, &md.Title, &md.Overview, &md.ReleaseDate, &md.Runtime, &md.Genres, &md.Casts); err != nil {
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
