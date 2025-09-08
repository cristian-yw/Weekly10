package models

import "time"

type Schedule struct {
	ID         int       `json:"id"`
	MovieTitle string    `json:"movie_title"`
	Cinema     string    `json:"cinema"`
	Location   string    `json:"location"`
	StartTime  time.Time `json:"start_time"`
	Date       time.Time `json:"date"`
}

type Seatcode struct {
	SeatCode string `json:"seat_code"`
}

type MovieDetail struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Overview    string   `json:"overview"`
	ReleaseDate string   `json:"release_date"`
	Runtime     int      `json:"runtime"`
	Genres      []string `json:"genres"`
	Casts       []string `json:"casts"`
}

type Order struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	ScheduleID int       `json:"schedule_id"`
	TotalPrice int       `json:"total_price"`
	Status     string    `json:"status"`
	OrderDate  time.Time `json:"order_date"`
	Seats      []string  `json:"seats"`
}
type Schedule2 struct {
	ID       int    `json:"id" example:"1"`
	MovieID  int    `json:"movie_id" example:"101"`
	Time     string `json:"time" example:"2025-09-08T14:00:00Z"`
	CinemaID int    `json:"cinema_id" example:"5"`
}

type Seat struct {
	ID       int    `json:"id"`
	CinemaID int    `json:"cinema_id"`
	SeatCode string `json:"seat_code"`
	IsBooked bool   `json:"is_booked"`
}
