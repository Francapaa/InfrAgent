package middleware

import (
	"fmt"
	"net/http"
	"server/repositories"
	"server/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// ProfileCompleteMiddleware verifica que el usuario haya completado su perfil
// antes de permitir el acceso a rutas protegidas como el dashboard
func ProfileCompleteMiddleware(client repositories.ClientStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener userID del contexto (seteado por JWTMiddleware)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
			c.Abort()
			return
		}

		// Parsear UUID
		userUUID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID format"})
			c.Abort()
			return
		}

		// Obtener usuario de la base de datos
		user, err := client.GetClient(c.Request.Context(), userUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error retrieving user"})
			c.Abort()
			return
		}

		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			c.Abort()
			return
		}

		// Verificar que el perfil esté completo
		if user.WebhookURL == "" || user.CompanyName == "" {
			missingFields := []string{}

			if user.WebhookURL == "" {
				missingFields = append(missingFields, "webhook_url")
			}

			if user.CompanyName == "" {
				missingFields = append(missingFields, "company_name")
			}

			c.JSON(http.StatusForbidden, gin.H{
				"error":          "profile incomplete",
				"message":        "User must complete registration before accessing this resource",
				"redirect_to":    "/onboarding",
				"missing_fields": missingFields,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
