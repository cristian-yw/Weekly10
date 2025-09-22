package models

type MovieRequest struct {
	Title        string      `json:"title"`
	Overview     string      `json:"overview"`
	ReleaseDate  string      `json:"release_date"`
	Runtime      int         `json:"runtime"`
	PosterPath   string      `json:"poster_path"`
	BackdropPath string      `json:"backdrop_path"`
	Popularity   float64     `json:"popularity"`
	VoteAverage  float64     `json:"vote_average"`
	VoteCount    int         `json:"vote_count"`
	Genres       []string    `json:"genres"`
	Schedules    []Schedule2 `json:"schedules"`
}
type SuccessMessage struct {
	Message string `json:"message" example:"Movie updated successfully"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"something went wrong"`
}
