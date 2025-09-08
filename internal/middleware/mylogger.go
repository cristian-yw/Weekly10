package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func MyLogger(ctx *gin.Context) {
	log.Println("start")
	start := time.Now()
	ctx.Next()
	duration := time.Since(start)
	log.Printf("Durasi Request: %dus\n", duration.Microseconds())
}
