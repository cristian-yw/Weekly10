package routers

import (
	"github.com/cristian-yw/Weekly10/internal/handlers"
	"github.com/cristian-yw/Weekly10/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitMovieRouter(r *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	movieRepo := repository.NewMovieRepository(db, rdb)
	movieHandler := handlers.NewMovieHandler(movieRepo)

	api := r.Group("/movies")
	{
		api.GET("/upcoming", movieHandler.GetUpcomingMovies)
		api.GET("/popular", movieHandler.GetPopularMovies)
		api.GET("/filter", movieHandler.GetMoviesWithFilter)
		api.GET("/all", movieHandler.GetAllMovies)
	}
}
