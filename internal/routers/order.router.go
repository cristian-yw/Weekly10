package routers

import (
	"github.com/cristian-yw/Weekly10/internal/handlers"
	"github.com/cristian-yw/Weekly10/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitOrderRouter(r *gin.Engine, db *pgxpool.Pool) {
	// buat repository dan handler
	orderRepo := repository.NewOrderRepository(db)
	orderHandler := handlers.NewOrderHandler(orderRepo)

	api := r.Group("/orders")
	{
		api.GET("/:movieId/schedules", orderHandler.GetSchedule)
		api.GET("/seats/:scheduleId", orderHandler.GetAvailableSeats)
		api.GET("/:movieId", orderHandler.GetMovieDetail)
		api.POST("/", orderHandler.CreateOrder)
	}
}
