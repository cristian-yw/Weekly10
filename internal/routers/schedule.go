package routers

import (
	"github.com/cristian-yw/Weekly10/internal/handlers"
	"github.com/cristian-yw/Weekly10/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func Initschedule(r *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	schedulesRepo := repository.NewScheduleRepository(db)
	schedulesHandler := handlers.NewScheduleHandler(schedulesRepo)

	r.GET("/genres", schedulesHandler.GetGenresHandler)
	r.GET("/cinemas", schedulesHandler.GetCinemasHandler)
	r.GET("/locations", schedulesHandler.GetLocationsHandler)

}
