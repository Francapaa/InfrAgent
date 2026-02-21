package middleware

import (
	"log"
	"net/http"
	"server/repositories"
	"server/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserIDKey es la clave usada para almacenar el ID del usuario en el contexto de Gin
const UserIDKey = "userID"

type Middleware struct {
	client repositories.ClientStorage
}

func AuthMiddlewareApiKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey, err := c.Cookie("auth_token")

		if err != nil {
			log.Printf("[AuthMiddleware] Error al extraer cookie: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		utils.IsValidAPIKeyMiddleware(apiKey)

		c.Next()
	}
}

// JWTMiddleware valida tokens JWT para autenticación de usuarios
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader, err := c.Cookie("auth_token")
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

		c.Set(UserIDKey, claims.UserID)
		c.Next()
	}
}

// ProfileCompleteMiddleware verifica que el usuario haya completado su perfil
// antes de permitir el acceso a rutas protegidas como el dashboard
func ProfileCompleteMiddleware(client repositories.ClientStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener userID del contexto (seteado por JWTMiddleware)
		userID, exists := c.Get(UserIDKey)
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
		missingFields := []string{}
		hasErrors := false

		// Validar WebhookURL
		if user.WebhookURL == "" || strings.TrimSpace(user.WebhookURL) == "" {
			missingFields = append(missingFields, "webhook_url")
			hasErrors = true
		} else if !strings.HasPrefix(user.WebhookURL, "https://") {
			missingFields = append(missingFields, "webhook_url_must_start_with_https")
			hasErrors = true
		}

		// Validar CompanyName
		if user.CompanyName == "" || strings.TrimSpace(user.CompanyName) == "" {
			missingFields = append(missingFields, "company_name")
			hasErrors = true
		}

		if hasErrors {
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
