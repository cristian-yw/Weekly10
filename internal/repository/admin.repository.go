package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cristian-yw/Weekly10/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepository struct {
	DB *pgxpool.Pool
}

func NewAdminRepository(db *pgxpool.Pool) *AdminRepository {
	return &AdminRepository{DB: db}
}

func (r *AdminRepository) CreateMovie(ctx context.Context, req models.NewMovieRequest) (int, error) {
	tx, err := r.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	// 1. Insert movie
	var movieID int
	err = tx.QueryRow(ctx, `
        INSERT INTO movies (tmdb_id, title, overview, release_date, runtime,
                            poster_path, backdrop_path, popularity, vote_average, vote_count, created_at, updated_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,NOW(),NOW())
        RETURNING id
    `, req.TMDBID, req.Title, req.Overview, req.ReleaseDate, req.Runtime,
		req.PosterPath, req.BackdropPath, req.Popularity, req.VoteAverage, req.VoteCount).
		Scan(&movieID)
	if err != nil {
		return 0, err
	}

	// 2. Insert genres
	for _, g := range req.Genres {
		_, err = tx.Exec(ctx, `
            INSERT INTO genres (name)
            VALUES ($1)
            ON CONFLICT (name) DO NOTHING
        `, g)
		if err != nil {
			return 0, err
		}

		_, err = tx.Exec(ctx, `
            INSERT INTO movie_genres (movie_id, genre_id)
            SELECT $1, id FROM genres WHERE name = $2
            ON CONFLICT DO NOTHING
        `, movieID, g)
		if err != nil {
			return 0, err
		}
	}

	// 3. Insert schedules
	for _, s := range req.Schedules {
		_, err = tx.Exec(ctx, `
            INSERT INTO schedules (movie_id, cinema_id, location_id, time_id, date, price)
            VALUES ($1, $2, $3, $4, $5, $6)
        `, movieID, s.CinemaID, s.LocationID, s.TimeID, s.Date, s.Price)
		if err != nil {
			return 0, err
		}
	}

	// 4. Insert director
	_, err = tx.Exec(ctx, `
        INSERT INTO persons (tmdb_id, name)
        VALUES ($1, $2)
        ON CONFLICT (tmdb_id) DO NOTHING
    `, req.Director.TMDBID, req.Director.Name)
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(ctx, `
        INSERT INTO movie_casts (movie_id, person_id, role)
        SELECT $1, id, 'Director' FROM persons WHERE tmdb_id = $2
        ON CONFLICT DO NOTHING
    `, movieID, req.Director.TMDBID)
	if err != nil {
		return 0, err
	}

	// 5. Insert casts
	for _, cast := range req.Casts {
		_, err = tx.Exec(ctx, `
            INSERT INTO persons (tmdb_id, name)
            VALUES ($1, $2)
            ON CONFLICT (tmdb_id) DO NOTHING
        `, cast.TMDBID, cast.Name)
		if err != nil {
			return 0, err
		}

		_, err = tx.Exec(ctx, `
            INSERT INTO movie_casts (movie_id, person_id, role)
            SELECT $1, id, 'Actor' FROM persons WHERE tmdb_id = $2
            ON CONFLICT DO NOTHING
        `, movieID, cast.TMDBID)
		if err != nil {
			return 0, err
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return movieID, nil
}

// GetMovieByID
func (r *AdminRepository) GetMovieByID(ctx context.Context, id int) (*models.TMDBMovie, error) {
	row := r.DB.QueryRow(ctx, `
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
		WHERE m.id = $1
		GROUP BY m.id
	`, id)

	var m models.TMDBMovie
	var genreNames []string
	if err := row.Scan(
		&m.ID,
		&m.TMDBID,
		&m.Title,
		&m.Overview,
		&m.ReleaseDate,
		&m.Popularity,
		&m.VoteAverage,
		&m.VoteCount,
		&genreNames,
		&m.PosterPath,
		&m.BackdropPath,
		&m.Runtime,
	); err != nil {
		return nil, err
	}
	m.Genres = genreNames
	return &m, nil
}

// UpdateMovie
func (r *AdminRepository) PatchMovie(ctx context.Context, id int, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}

	setClauses := []string{}
	args := []interface{}{}
	i := 1

	for k, v := range fields {
		setClauses = append(setClauses, fmt.Sprintf("%s=$%d", k, i))
		args = append(args, v)
		i++
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at=NOW()"))

	query := fmt.Sprintf(`
        UPDATE movies
        SET %s
        WHERE id=$%d
    `, strings.Join(setClauses, ", "), i)

	args = append(args, id)

	_, err := r.DB.Exec(ctx, query, args...)
	return err
}

// DeleteMovie
func (r *AdminRepository) DeleteMovie(ctx context.Context, id int) error {
	_, err := r.DB.Exec(ctx, `DELETE FROM movies WHERE id=$1`, id)
	return err
}

//
// -------------------- UPSERT / LINK HELPERS --------------------
//

func (r *AdminRepository) UpsertMovie(m models.TMDBMovie) (int, error) {
	var movieID int
	err := r.DB.QueryRow(
		context.Background(),
		`INSERT INTO movies (
			tmdb_id, title, overview, release_date, runtime,
			poster_path, backdrop_path, popularity, vote_average, vote_count, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,NOW(),NOW())
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
		m.TMDBID, m.Title, m.Overview, m.ReleaseDate, m.Runtime,
		m.PosterPath, m.BackdropPath, m.Popularity, m.VoteAverage, m.VoteCount,
	).Scan(&movieID)

	return movieID, err
}

// UpsertGenre: insert/update genre by tmdb_id
func (r *AdminRepository) UpsertGenre(tmdbID int, name string) (int, error) {
	var genreID int
	err := r.DB.QueryRow(context.Background(), `
		INSERT INTO genres (tmdb_id, name)
		VALUES ($1,$2)
		ON CONFLICT (tmdb_id) DO UPDATE SET name=EXCLUDED.name
		RETURNING id
	`, tmdbID, name).Scan(&genreID)

	return genreID, err
}

// LinkMovieGenre: relasi movie - genre
func (r *AdminRepository) LinkMovieGenre(movieID, genreID int) error {
	_, err := r.DB.Exec(context.Background(), `
		INSERT INTO movie_genres (movie_id, genre_id)
		VALUES ($1,$2)
		ON CONFLICT DO NOTHING
	`, movieID, genreID)
	return err
}

// UpsertCategory: insert/update category by name
func (r *AdminRepository) UpsertCategory(name string) (int, error) {
	var categoryID int
	err := r.DB.QueryRow(context.Background(), `
		INSERT INTO categories (name)
		VALUES ($1)
		ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name
		RETURNING id
	`, name).Scan(&categoryID)

	return categoryID, err
}

// LinkMovieCategory: relasi movie - category
func (r *AdminRepository) LinkMovieCategory(movieID, categoryID int) error {
	_, err := r.DB.Exec(context.Background(), `
		INSERT INTO movie_categories (movie_id, category_id)
		VALUES ($1,$2)
		ON CONFLICT DO NOTHING
	`, movieID, categoryID)
	return err
}

//
// -------------------- SYNC & FETCH (existing code) --------------------
//

func (r *AdminRepository) fetchGenres(apiKey string) (map[int]string, error) {
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

// SyncPopular (keputusan: menggunakan upcoming/popular sesuai kebutuhan)
func (r *AdminRepository) SyncPopular(apiKey string) error {
	genreMap, err := r.fetchGenres(apiKey)
	if err != nil {
		return err
	}

	// NOTE: jika ingin sync popular ganti /movie/upcoming -> /movie/popular
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

	categoryID, err := r.UpsertCategory("popular")
	if err != nil {
		return err
	}

	for _, m := range result.Results {
		movieID, err := r.UpsertMovie(m)
		if err != nil {
			return err
		}

		// Convert genre IDs to names using genreMap
		var genreNames []string
		for _, gid := range m.GenreIDs {
			if name, ok := genreMap[gid]; ok {
				genreNames = append(genreNames, name)
			}
		}

		// Use the same genre handling as CreateMovie
		for _, gname := range genreNames {
			var gid int
			if err := r.DB.QueryRow(context.Background(),
				`INSERT INTO genres (name)
				 VALUES ($1)
				 ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name
				 RETURNING id`,
				gname).Scan(&gid); err != nil {
				return fmt.Errorf("insert genre %s: %w", gname, err)
			}

			if _, err := r.DB.Exec(context.Background(),
				`INSERT INTO movie_genres (movie_id, genre_id)
				 VALUES ($1,$2)
				 ON CONFLICT DO NOTHING`,
				movieID, gid); err != nil {
				return fmt.Errorf("link genre %s: %w", gname, err)
			}
		}

		if err := r.LinkMovieCategory(movieID, categoryID); err != nil {
			return err
		}
	}

	return nil
}
