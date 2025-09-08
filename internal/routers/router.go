package routers

import (
	"net/http"

	docs "github.com/cristian-yw/Weekly10/docs"
	"github.com/cristian-yw/Weekly10/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(db *pgxpool.Pool) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.MyLogger)
	router.Use(middleware.CORSMiddleware)
	InitAuthRouter(router, db)
	InitMovieRouter(router, db)
	InitOrderRouter(router, db)
	InitUserRouter(router, db)
	InitAdminMovieRouter(router, db)

	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"Message": "Rute Salah",
			"Status":  "Rute tidak ditemukan",
		})
	})
	return router
}
