package routers

import (
	"github.com/cristian-yw/Weekly10/internal/handlers"
	"github.com/cristian-yw/Weekly10/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitMovieRouter(r *gin.Engine, db *pgxpool.Pool) {
	movieHandler := handlers.NewMovieHandler(repository.NewMovieRepository(db))

	api := r.Group("/movies")
	{
		api.POST("/sync/popular", movieHandler.SyncPopular)
		api.GET("/upcoming", movieHandler.GetUpcomingMovies)
		api.GET("/popular", movieHandler.GetPopularMovies)
		api.GET("/filter", movieHandler.GetMoviesWithFilter)
	}
}
