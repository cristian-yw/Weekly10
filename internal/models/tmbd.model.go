package models

type TMDBMovie struct {
	ID           int     `json:"id"`
	TMDBID       int     `json:"tmdb_id"`
	Title        string  `json:"title"`
	Overview     string  `json:"overview"`
	ReleaseDate  string  `json:"release_date"`
	GenreIDs     []int   `json:"genre_ids"`
	Runtime      int     `json:"runtime"`
	Popularity   float64 `json:"popularity"`
	VoteAverage  float64 `json:"vote_average"`
	VoteCount    int     `json:"vote_count"`
	PosterPath   string  `json:"poster_path"`
	BackdropPath string  `json:"backdrop_path"`
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
