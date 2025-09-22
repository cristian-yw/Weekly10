package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cristian-yw/Weekly10/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type MovieRepository struct {
	DB  *pgxpool.Pool
	rdb *redis.Client
}

func NewMovieRepository(db *pgxpool.Pool, rdb *redis.Client) *MovieRepository {
	return &MovieRepository{DB: db, rdb: rdb}
}

const cacheTTL = 5 * time.Minute

// GetUpcomingMovies with Redis
func (r *MovieRepository) GetUpcomingMovies(ctx context.Context, limit, offset int) ([]models.TMDBMovie, error) {
	key := fmt.Sprintf("movies:upcoming:%d:%d", limit, offset)
	if data, err := r.rdb.Get(ctx, key).Bytes(); err == nil {
		var cached []models.TMDBMovie
		if json.Unmarshal(data, &cached) == nil {
			return cached, nil
		}
	}

	movies, err := r.getUpcomingMoviesDB(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	_ = r.cache(ctx, key, movies)
	return movies, nil
}

// GetPopularMovies with Redis
func (r *MovieRepository) GetPopularMovies(ctx context.Context, limit, offset int) ([]models.TMDBMovie, error) {
	key := fmt.Sprintf("movies:popular:%d:%d", limit, offset)
	if data, err := r.rdb.Get(ctx, key).Bytes(); err == nil {
		var cached []models.TMDBMovie
		if json.Unmarshal(data, &cached) == nil {
			return cached, nil
		}
	}

	movies, err := r.getPopularMoviesDB(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	_ = r.cache(ctx, key, movies)
	return movies, nil
}

// GetMoviesWithFilter with Redis
func (r *MovieRepository) GetMoviesWithFilter(ctx context.Context, name string, genreID int, limit, offset int) ([]models.TMDBMovie, error) {
	key := fmt.Sprintf("movies:filter:%s:%d:%d:%d", name, genreID, limit, offset)
	if data, err := r.rdb.Get(ctx, key).Bytes(); err == nil {
		var cached []models.TMDBMovie
		if json.Unmarshal(data, &cached) == nil {
			return cached, nil
		}
	}

	movies, err := r.getMoviesWithFilterDB(ctx, name, genreID, limit, offset)
	if err != nil {
		return nil, err
	}
	_ = r.cache(ctx, key, movies)
	return movies, nil
}

// GetAllMovies with Redis
func (r *MovieRepository) GetAllMovies(ctx context.Context) ([]models.TMDBMovie, error) {
	key := "movies:all"
	if data, err := r.rdb.Get(ctx, key).Bytes(); err == nil {
		var cached []models.TMDBMovie
		if json.Unmarshal(data, &cached) == nil {
			return cached, nil
		}
	}

	movies, err := r.getAllMoviesDB(ctx)
	if err != nil {
		return nil, err
	}
	_ = r.cache(ctx, key, movies)
	return movies, nil
}

func (r *MovieRepository) cache(ctx context.Context, key string, movies []models.TMDBMovie) error {
	data, err := json.Marshal(movies)
	if err != nil {
		return err
	}
	if err := r.rdb.Set(ctx, key, data, cacheTTL).Err(); err != nil {
		log.Println("Redis SET error:", err)
	}
	return nil
}

// -------------------- Raw DB queries --------------------

func (r *MovieRepository) getUpcomingMoviesDB(ctx context.Context, limit, offset int) ([]models.TMDBMovie, error) {
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

func (r *MovieRepository) getPopularMoviesDB(ctx context.Context, limit, offset int) ([]models.TMDBMovie, error) {
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

func (r *MovieRepository) getMoviesWithFilterDB(
	ctx context.Context,
	name string,
	genreID int,
	limit, offset int,
) ([]models.TMDBMovie, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT 
			m.id, 
			m.tmdb_id, 
			m.title, 
			m.overview, 
			COALESCE(TO_CHAR(m.release_date, 'YYYY-MM-DD'), '') AS release_date,
			COALESCE(m.runtime, 0) AS runtime,
			COALESCE(m.poster_path, '') AS poster_path,
			COALESCE(m.backdrop_path, '') AS backdrop_path,
			m.popularity, 
			m.vote_average, 
			m.vote_count, 
			COALESCE(array_agg(DISTINCT g.name) FILTER (WHERE g.name IS NOT NULL), '{}') AS genres
		FROM movies m
		LEFT JOIN movie_genres mg ON mg.movie_id = m.id
		LEFT JOIN genres g ON g.id = mg.genre_id
		WHERE ($1 = '' OR m.title ILIKE '%' || $1 || '%')
		  AND ($2 = 0 OR g.tmdb_id = $2)
		GROUP BY m.id
		ORDER BY m.release_date DESC
		LIMIT $3 OFFSET $4;
	`, name, genreID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.TMDBMovie
	for rows.Next() {
		var m models.TMDBMovie
		var genres []string
		if err := rows.Scan(
			&m.ID,
			&m.TMDBID,
			&m.Title,
			&m.Overview,
			&m.ReleaseDate,
			&m.Runtime,
			&m.PosterPath,
			&m.BackdropPath,
			&m.Popularity,
			&m.VoteAverage,
			&m.VoteCount,
			&genres,
		); err != nil {
			return nil, err
		}
		m.Genres = genres
		movies = append(movies, m)
	}

	// biar tidak null di JSON
	if movies == nil {
		movies = []models.TMDBMovie{}
	}

	return movies, nil
}

func (r *MovieRepository) getAllMoviesDB(ctx context.Context) ([]models.TMDBMovie, error) {
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
			COALESCE(array_agg(DISTINCT g.name) FILTER (WHERE g.name IS NOT NULL), '{}') AS genres,
			COALESCE(m.poster_path, '') AS poster_path,
			COALESCE(m.backdrop_path, '') AS backdrop_path,
			COALESCE(m.runtime, 0) AS runtime
		FROM movies m
		LEFT JOIN movie_genres mg ON m.id = mg.movie_id
		LEFT JOIN genres g ON mg.genre_id = g.id
		GROUP BY m.id
		ORDER BY m.release_date DESC;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.TMDBMovie
	for rows.Next() {
		var m models.TMDBMovie
		var genres []string
		if err := rows.Scan(
			&m.ID,
			&m.TMDBID,
			&m.Title,
			&m.Overview,
			&m.ReleaseDate,
			&m.Popularity,
			&m.VoteAverage,
			&m.VoteCount,
			&genres,
			&m.PosterPath,
			&m.BackdropPath,
			&m.Runtime,
		); err != nil {
			log.Println("Scan error:", err)
			return nil, err
		}
		m.Genres = genres
		movies = append(movies, m)
	}
	return movies, nil
}

func (r *MovieRepository) CountMoviesWithFilter(ctx context.Context, name string, genreID int) (int, error) {
	var total int

	query := `
        SELECT COUNT(DISTINCT m.id)
        FROM movies m
        LEFT JOIN movie_genres mg ON m.id = mg.movie_id
        WHERE 1=1
    `
	args := []interface{}{}
	argIdx := 1

	if name != "" {
		query += fmt.Sprintf(" AND m.title ILIKE $%d", argIdx)
		args = append(args, "%"+name+"%")
		argIdx++
	}

	if genreID != 0 {
		query += fmt.Sprintf(" AND mg.genre_id = $%d", argIdx)
		args = append(args, genreID)
		argIdx++
	}

	err := r.DB.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
