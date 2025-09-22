package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cristian-yw/Weekly10/internal/config"
	"github.com/cristian-yw/Weekly10/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

var jwtKey = []byte(getSecret())

func getSecret() string {
	if os.Getenv("JWT_SECRET") == "" {
		return "SECRET_JWT_KEY"
	}
	return os.Getenv("JWT_SECRET")
}

func GenerateJWT(userID int, role string) (string, error) {
	claims := &models.JWTClaim{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func AuthMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		if val, _ := rdb.Get(config.Ctx, "blacklist:"+tokenString).Result(); val == "true" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token sudah logout"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaim{}, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalid"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*models.JWTClaim); ok {
			c.Set("userID", claims.UserID)
			c.Set("role", claims.Role)
			c.Set("token", tokenString) // untuk logout handler
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role.(string) != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin only"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func UserOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role.(string) != "user" {
			c.JSON(http.StatusForbidden, gin.H{"error": "User only"})
			c.Abort()
			return
		}
		c.Next()
	}
}
