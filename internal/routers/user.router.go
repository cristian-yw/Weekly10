package routers

import (
	"github.com/cristian-yw/Weekly10/internal/handlers"
	"github.com/cristian-yw/Weekly10/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitAuthRouter(r *gin.Engine, db *pgxpool.Pool) {
	// buat repository dan handler
	authRepo := repository.NewAuthRepository(db)
	authHandler := handlers.NewAuthHandler(authRepo)

	api := r.Group("/auth")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
		// api.GET("/profile", middleware.AuthMiddleware(), authHandler.Profile)
	}
}
