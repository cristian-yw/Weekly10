package routers

import (
	"github.com/cristian-yw/Weekly10/internal/handlers"
	"github.com/cristian-yw/Weekly10/internal/middleware"
	"github.com/cristian-yw/Weekly10/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitAuthRouter(r *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	// buat repository dan handler
	authRepo := repository.NewAuthRepository(db)
	authHandler := handlers.NewAuthHandler(authRepo, rdb)

	api := r.Group("/auth")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
		// api.GET("/profile", middleware.AuthMiddleware(), authHandler.Profile)
		api.POST("/logout", middleware.AuthMiddleware(rdb), authHandler.Logout)
	}
}
