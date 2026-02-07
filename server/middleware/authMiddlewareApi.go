package middleware

import (
	"fmt"
	"net/http"
	"server/repositories"
	"server/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	client repositories.ClientStorage
}

func (m *Middleware) authMiddlewareApi() gin.HandlerFunc {

	return func(c *gin.Context) {
		authHeader := c.GetHeader("auth_token")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			c.Abort()
			return
		}
		apiKey := parts[1]

		client, err := m.client.GetClientByAPIKey(c.Request.Context(), apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		c.Set("client", client)
		c.Next()
	}

}

// JWTMiddleware valida tokens JWT para autenticación de usuarios
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader, err := c.Cookie("auth_token")
		fmt.Println("header de autenticacion: ", authHeader)
		if authHeader == "" || err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateJWT(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
