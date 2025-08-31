package middlewares

import (
	"fleet-pulse-users-service/internal/errors"
	"fleet-pulse-users-service/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header"})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := authService.ParseJWT(token)
		if err != nil {
			errors.HandleAuthErrors(c, err)
			c.Abort()
			return
		}
		c.Set("current_user_id", claims.UserID)
		c.Next()
	}
}
