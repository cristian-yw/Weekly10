package models

import "time"

type Schedule struct {
	ID         int    `json:"id"`
	MovieID    int    `json:"movie_id"`
	Cinema     string `json:"cinema"`
	MovieTitle string `json:"movie_title"`
	Location   string `json:"location"`
	StartTime  string `json:"start_time"`
	Date       string `json:"date"`
}

type Seat struct {
	ID       int    `json:"id"`
	CinemaID int    `json:"cinema_id"`
	SeatCode string `json:"seat_code"`
	IsBooked bool   `json:"is_booked"`
}

type MovieDetail struct {
	ID           int        `json:"id"`
	Title        string     `json:"title"`
	Overview     string     `json:"overview"`
	ReleaseDate  string     `json:"release_date"`
	Runtime      int        `json:"runtime"`
	Genres       []string   `json:"genres"`
	Casts        []string   `json:"casts"`
	Director     string     `json:"director"`
	PosterPath   string     `json:"poster_path"`
	BackdropPath string     `json:"backdrop_path"`
	Schedules    []Schedule `json:"schedules"`
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

// Schedule2 dipakai untuk input request Add Movie (jadwal baru)
type Schedule2 struct {
	CinemaID   int    `json:"cinema_id"`
	LocationID int    `json:"location_id"`
	TimeID     int    `json:"time_id"`
	Date       string `json:"date"`
}
