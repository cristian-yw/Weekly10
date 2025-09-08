package routers

import (
	"github.com/cristian-yw/Weekly10/internal/handlers"
	"github.com/cristian-yw/Weekly10/internal/middleware"
	"github.com/cristian-yw/Weekly10/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitUserRouter(r *gin.Engine, db *pgxpool.Pool) {
	// Inisialisasi repository dan handler
	userRepo := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)

	api := r.Group("/auth")
	api.Use(middleware.AuthMiddleware())
	{

		api.GET("/profile", userHandler.GetProfile)
		api.GET("/history", userHandler.GetHistory)
		api.PUT("/profile", userHandler.EditProfile)
	}
}
