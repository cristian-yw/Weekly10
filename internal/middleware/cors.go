package middleware

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(ctx *gin.Context) {
	whilelist := []string{"http://localhost:3000"}
	origin := ctx.GetHeader("Origin")
	if slices.Contains(whilelist, origin) {
		ctx.Header("Access-Control-Allow-Origin", origin)
	}
	ctx.Header("Access-Control-Allow-Methods", "GET")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}
	ctx.Next()
}
