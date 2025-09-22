package models

type TMDBMovie struct {
	ID           int      `json:"id"`
	TMDBID       *int     `json:"tmdb_id"`
	Title        string   `json:"title"`
	Overview     string   `json:"overview"`
	ReleaseDate  string   `json:"release_date"`
	GenreIDs     []int    `json:"-"`
	Genres       []string `json:"genres"`
	Runtime      int      `json:"runtime"`
	Popularity   float64  `json:"popularity"`
	VoteAverage  float64  `json:"vote_average"`
	VoteCount    int      `json:"vote_count"`
	PosterPath   string   `json:"poster_path"`
	BackdropPath string   `json:"backdrop_path"`
}

type TMDBGenre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TMDBGenreResponse struct {
	Genres []TMDBGenre `json:"genres"`
}

type TMDBResponse struct {
	Results []TMDBMovie `json:"results"`
}

// ============================
// REQUEST UNTUK ADD MOVIE
// ============================

type NewMovieRequest struct {
	TMDBID       *int              `json:"tmdb_id"` // boleh custom
	Title        string            `json:"title"`
	Overview     string            `json:"overview"`
	ReleaseDate  string            `json:"release_date"`
	Runtime      int               `json:"runtime"`
	PosterPath   string            `json:"poster_path"`
	BackdropPath string            `json:"backdrop_path"`
	Popularity   float64           `json:"popularity"`
	VoteAverage  float64           `json:"vote_average"`
	VoteCount    int               `json:"vote_count"`
	Genres       []string          `json:"genres"`
	Schedules    []ScheduleRequest `json:"schedules"`
	Director     *PersonRequest    `json:"director"`
	Casts        []PersonRequest   `json:"casts"`
}

type PersonRequest struct {
	TMDBID int    `json:"tmdb_id"`
	Name   string `json:"name"`
}

type ScheduleRequest struct {
	CinemaID   int    `json:"cinema_id"`
	LocationID int    `json:"location_id"`
	TimeID     int    `json:"time_id"`
	Date       string `json:"date"`
	Price      int    `json:"price"`
}
