package routers

import (
	"net/http"

	docs "github.com/cristian-yw/Weekly10/docs"
	"github.com/cristian-yw/Weekly10/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(db *pgxpool.Pool, rdb *redis.Client) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.MyLogger)
	router.Use(middleware.CORSMiddleware())

	InitAuthRouter(router, db, rdb)
	InitMovieRouter(router, db, rdb)
	InitOrderRouter(router, db, rdb)
	InitUserRouter(router, db, rdb)
	InitAdminMovieRouter(router, db, rdb)
	Initschedule(router, db, rdb)

	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Static("/uploads", "./uploads")

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"Message": "Rute Salah",
			"Status":  "Rute tidak ditemukan",
		})
	})
	return router
}
