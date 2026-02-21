package controllers

import (
	"net/http"
	models "server/model"
	"server/service"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

type LoginController struct {
	service *service.Login
}

func NewLoginController(loginService *service.Login) *LoginController {
	return &LoginController{
		service: loginService,
	}
}

func NewWebSocketController() *WebSocketController {
	return &WebSocketController{}
}

func (lc *LoginController) GoogleLogin(ctx *gin.Context) {
	provider := ctx.Param("provider")
	ctx.Request.URL.RawQuery = "provider=" + provider
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

// el c * authController nos permite identificar a q pertenece la funcion
// es como si seria parte del objeto authController, se lo llama RECEPTOR
func (lc *LoginController) GetAuthCallBackFunction(ctx *gin.Context) {

	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		ctx.JSON(401, gin.H{"Error": "Unauthorized" + err.Error()})
		return
	}
	token, err := lc.service.LoginWithGoogle(user)
	if err != nil {
		ctx.JSON(500, gin.H{"Error": err.Error()})
		return
	}

	// Buscar el usuario para verificar si tiene perfil completo
	existingUser, err := lc.service.GetUserByEmail(ctx, user.Email)
	if err != nil {
		ctx.JSON(500, gin.H{"Error": "Error retrieving user"})
		return
	}

	const sessionDuration = 8 * 3600

	ctx.SetCookie(
		"auth_token",
		token,
		sessionDuration,
		"/",
		"",
		true,
		true,
	)

	// Verificar si el usuario ya completó el perfil
	// Si tiene webhook_url y company_name, ir directo al dashboard
	if existingUser.WebhookURL != "" && existingUser.CompanyName != "" {
		ctx.Redirect(http.StatusFound, "http://localhost:3000/dashboard")
		return
	}

	// Si no tiene perfil completo, ir a onboarding
	ctx.Redirect(http.StatusFound, "http://localhost:3000/onboarding")
}

func (lc *LoginController) CompleteRegistration(ctx *gin.Context) {
	var req models.CompleteRegistrationRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "webhook_url is required"})
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	response, err := lc.service.CompleteRegistration(ctx, userIDStr, req.CompanyName, req.WebhookURL)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (lc *LoginController) GetCurrentUser(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := lc.service.GetUserByID(ctx, userIDStr)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"email":        user.Email,
		"company_name": user.CompanyName,
		"webhook_url":  user.WebhookURL,
		"metodo":       user.Metodo,
		"created_at":   user.CreatedAt,
		"updated_at":   user.UpdatedAt,
	})
}
