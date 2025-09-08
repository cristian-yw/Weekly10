package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cristian-yw/Weekly10/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MovieRepository struct {
	DB *pgxpool.Pool
}

func NewMovieRepository(db *pgxpool.Pool) *MovieRepository {
	return &MovieRepository{DB: db}
}

// -------------------- DB OPS --------------------

// Insert / Update Movie
func (r *MovieRepository) upsertMovie(m models.TMDBMovie) (int, error) {
	var movieID int
	err := r.DB.QueryRow(
		context.Background(),
		`INSERT INTO movies (tmdb_id, title, overview, release_date, runtime, poster_path, backdrop_path, popularity, vote_average, vote_count, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,NOW(),NOW())
		ON CONFLICT (tmdb_id) DO UPDATE SET
			title=EXCLUDED.title,
			overview=EXCLUDED.overview,
			release_date=EXCLUDED.release_date,
			runtime=EXCLUDED.runtime,
			poster_path=EXCLUDED.poster_path,
			backdrop_path=EXCLUDED.backdrop_path,
			popularity=EXCLUDED.popularity,
			vote_average=EXCLUDED.vote_average,
			vote_count=EXCLUDED.vote_count,
			updated_at=NOW()
		RETURNING id`,
		m.ID, m.Title, m.Overview, m.ReleaseDate, m.Runtime,
		m.PosterPath, m.BackdropPath, m.Popularity, m.VoteAverage, m.VoteCount,
	).Scan(&movieID)

	return movieID, err
}

// Insert genre jika belum ada
func (r *MovieRepository) upsertGenre(tmdbID int, name string) (int, error) {
	var genreID int
	err := r.DB.QueryRow(context.Background(), `
		INSERT INTO genres (tmdb_id, name)
		VALUES ($1,$2)
		ON CONFLICT (tmdb_id) DO UPDATE SET name=EXCLUDED.name
		RETURNING id
	`, tmdbID, name).Scan(&genreID)

	return genreID, err
}

// Relasi movie_genres
func (r *MovieRepository) linkMovieGenre(movieID, genreID int) error {
	_, err := r.DB.Exec(context.Background(), `
		INSERT INTO movie_genres (movie_id, genre_id)
		VALUES ($1,$2)
		ON CONFLICT DO NOTHING
	`, movieID, genreID)
	return err
}

// Insert kategori
func (r *MovieRepository) upsertCategory(name string) (int, error) {
	var categoryID int
	err := r.DB.QueryRow(context.Background(), `
		INSERT INTO categories (name)
		VALUES ($1)
		ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name
		RETURNING id
	`, name).Scan(&categoryID)

	return categoryID, err
}

// Relasi movie_categories
func (r *MovieRepository) linkMovieCategory(movieID, categoryID int) error {
	_, err := r.DB.Exec(context.Background(), `
		INSERT INTO movie_categories (movie_id, category_id)
		VALUES ($1,$2)
		ON CONFLICT DO NOTHING
	`, movieID, categoryID)
	return err
}

func (r *MovieRepository) fetchGenres(apiKey string) (map[int]string, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/genre/movie/list?api_key=%s", apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var genreResp models.TMDBGenreResponse
	if err := json.NewDecoder(resp.Body).Decode(&genreResp); err != nil {
		return nil, err
	}

	genreMap := make(map[int]string)
	for _, g := range genreResp.Genres {
		genreMap[g.ID] = g.Name
	}
	return genreMap, nil
}

// Sinkronisasi popular movies
func (r *MovieRepository) SyncPopular(apiKey string) error {
	genreMap, err := r.fetchGenres(apiKey)
	if err != nil {
		return err
	}

	// Fetch popular movies
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/upcoming?api_key=%s", apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result models.TMDBResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	categoryID, err := r.upsertCategory("popular")
	if err != nil {
		return err
	}

	// Simpan semua movie
	for _, m := range result.Results {
		movieID, err := r.upsertMovie(m)
		if err != nil {
			return err
		}

		// Simpan genre
		for _, gid := range m.GenreIDs {
			name := genreMap[gid]
			genreID, err := r.upsertGenre(gid, name)
			if err != nil {
				return err
			}
			if err := r.linkMovieGenre(movieID, genreID); err != nil {
				return err
			}
		}

		if err := r.linkMovieCategory(movieID, categoryID); err != nil {
			return err
		}
	}

	return nil
}

// Get Upcoming Movies
func (r *MovieRepository) GetUpcomingMovies(ctx context.Context, limit, offset int) ([]models.TMDBMovie, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT 
			m.id,
			m.tmdb_id,
			m.title,
			m.overview,
			COALESCE(TO_CHAR(m.release_date, 'YYYY-MM-DD'), '') AS release_date,
			m.popularity,
			m.vote_average,
			m.vote_count,
			COALESCE(array_agg(DISTINCT g.id) FILTER (WHERE g.id IS NOT NULL), '{}') AS genre_ids,
			COALESCE(m.poster_path, '') AS poster_path,
			COALESCE(m.backdrop_path, '') AS backdrop_path,
			COALESCE(m.runtime, 0) AS runtime
		FROM movies m
		LEFT JOIN movie_genres mg ON m.id = mg.movie_id
		LEFT JOIN genres g ON mg.genre_id = g.id
		WHERE m.release_date > NOW()
		GROUP BY m.id
		ORDER BY m.release_date ASC
		LIMIT $1 OFFSET $2;
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.TMDBMovie
	for rows.Next() {
		var m models.TMDBMovie
		var genreIDs []int32

		if err := rows.Scan(
			&m.ID,
			&m.TMDBID,
			&m.Title,
			&m.Overview,
			&m.ReleaseDate,
			&m.Popularity,
			&m.VoteAverage,
			&m.VoteCount,
			&genreIDs,
			&m.PosterPath,
			&m.BackdropPath,
			&m.Runtime,
		); err != nil {
			return nil, err
		}
		m.GenreIDs = make([]int, len(genreIDs))
		for i, g := range genreIDs {
			m.GenreIDs[i] = int(g)
		}

		movies = append(movies, m)
	}
	return movies, nil
}

func (r *MovieRepository) GetPopularMovies(ctx context.Context, limit, offset int) ([]models.TMDBMovie, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT 
			m.id,
			m.tmdb_id,
			m.title,
			m.overview,
			COALESCE(TO_CHAR(m.release_date, 'YYYY-MM-DD'), '') AS release_date,
			m.popularity,
			m.vote_average,
			m.vote_count,
			COALESCE(array_agg(DISTINCT g.id) FILTER (WHERE g.id IS NOT NULL), '{}') AS genre_ids,
			COALESCE(m.poster_path, '') AS poster_path,
			COALESCE(m.backdrop_path, '') AS backdrop_path,
			COALESCE(m.runtime, 0) AS runtime
		FROM movies m
		LEFT JOIN movie_genres mg ON m.id = mg.movie_id
		LEFT JOIN genres g ON mg.genre_id = g.id
		GROUP BY m.id
		ORDER BY m.popularity DESC
		LIMIT $1 OFFSET $2;
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.TMDBMovie
	for rows.Next() {
		var m models.TMDBMovie
		var genreIDs []int32

		if err := rows.Scan(
			&m.ID,
			&m.TMDBID,
			&m.Title,
			&m.Overview,
			&m.ReleaseDate,
			&m.Popularity,
			&m.VoteAverage,
			&m.VoteCount,
			&genreIDs,
			&m.PosterPath,
			&m.BackdropPath,
			&m.Runtime,
		); err != nil {
			return nil, err
		}

		m.GenreIDs = make([]int, len(genreIDs))
		for i, g := range genreIDs {
			m.GenreIDs[i] = int(g)
		}

		movies = append(movies, m)
	}

	return movies, nil
}

// Get Movie with Filter
func (r *MovieRepository) GetMoviesWithFilter(ctx context.Context, name string, genreID int, limit, offset int) ([]models.TMDBMovie, error) {
	query := `
		SELECT 
			m.id, 
			m.tmdb_id, 
			m.title, 
			m.overview, 
			COALESCE(TO_CHAR(m.release_date, 'YYYY-MM-DD'), '') AS release_date, 
			m.popularity, 
			m.vote_average, 
			m.vote_count, 
			COALESCE(array_agg(DISTINCT g.id) FILTER (WHERE g.id IS NOT NULL), '{}') AS genre_ids
		FROM movies m
		LEFT JOIN movie_genres mg ON mg.movie_id = m.id
		LEFT JOIN genres g ON g.id = mg.genre_id
		WHERE ($1 = '' OR m.title ILIKE '%' || $1 || '%')
		  AND ($2 = 0 OR g.tmdb_id = $2)
		GROUP BY m.id
		ORDER BY m.release_date DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.DB.Query(ctx, query, name, genreID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.TMDBMovie
	for rows.Next() {
		var m models.TMDBMovie
		var genreIDs []int32
		err := rows.Scan(
			&m.ID,
			&m.TMDBID,
			&m.Title,
			&m.Overview,
			&m.ReleaseDate,
			&m.Popularity,
			&m.VoteAverage,
			&m.VoteCount,
			&genreIDs,
		)
		if err != nil {
			return nil, err
		}
		m.GenreIDs = make([]int, len(genreIDs))
		for i, g := range genreIDs {
			m.GenreIDs[i] = int(g)
		}
		movies = append(movies, m)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return movies, nil
}
