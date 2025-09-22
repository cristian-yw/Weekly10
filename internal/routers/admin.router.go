package routers

import (
	"github.com/cristian-yw/Weekly10/internal/handlers"
	"github.com/cristian-yw/Weekly10/internal/middleware"
	"github.com/cristian-yw/Weekly10/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitAdminMovieRouter(r *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	movieRepo := repository.NewAdminRepository(db)
	movieHandler := handlers.NewAdminHandler(movieRepo)

	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware(rdb), middleware.AdminOnly())
	{
		admin.POST("/sync/popular", movieHandler.SyncPopular)
		admin.POST("/movies", movieHandler.CreateMovie)       // Create Movie
		admin.GET("/movies/:id", movieHandler.GetMovieByID)   // Get Movie by ID
		admin.PATCH("/movies/:id", movieHandler.PatchMovie)   // Update Movie
		admin.DELETE("/movies/:id", movieHandler.DeleteMovie) // Delete Movie
	}
}
