package middleware

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(ctx *gin.Context) {
	whilelist := []string{"http://localhost:5173"}
	origin := ctx.GetHeader("Origin")
	if slices.Contains(whilelist, origin) {
		ctx.Header("Access-Control-Allow-Origin", origin)
	}
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}
	ctx.Next()
}
