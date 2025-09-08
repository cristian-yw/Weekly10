package routers

import (
	"github.com/cristian-yw/Weekly10/internal/handlers"
	"github.com/cristian-yw/Weekly10/internal/middleware"
	"github.com/cristian-yw/Weekly10/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitAdminMovieRouter(r *gin.Engine, db *pgxpool.Pool) {
	movieRepo := repository.NewAdminRepository(db)
	movieHandler := handlers.NewAdminHandler(movieRepo)

	admin := r.Group("/admin/movies")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
	{
		admin.GET("", movieHandler.GetAllMovies)
		admin.PUT("/:id", movieHandler.UpdateMovie)
		admin.DELETE("/:id", movieHandler.DeleteMovie)
	}
}
