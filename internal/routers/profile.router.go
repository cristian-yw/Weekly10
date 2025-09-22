package routers

import (
	"github.com/cristian-yw/Weekly10/internal/handlers"
	"github.com/cristian-yw/Weekly10/internal/middleware"
	"github.com/cristian-yw/Weekly10/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitUserRouter(r *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	userRepo := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)

	api := r.Group("/user")
	api.Use(middleware.AuthMiddleware(rdb))
	{

		api.GET("/profile", userHandler.GetProfile)
		api.GET("/history", userHandler.GetHistory)
		api.PATCH("/profile", userHandler.UpdateProfile)
		api.PATCH("/password", userHandler.ChangePassword)
	}
}
