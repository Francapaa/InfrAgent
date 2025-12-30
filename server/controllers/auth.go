package controllers

import (
	"net/http"
	models "server/model"
	"server/service"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

type authController struct {
	authService services.authService
}

func loginControllers(c *gin.Context) {

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Success": false,
			"Message": "Datos faltantes o invalidos" + err.Error(),
		})
		return
	}

	token, err := service.LoginUser(user.Email, user.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Success": false,
			"Message": "Falló en SERVICES",
		})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Success: true,
		Message: "Todo ok, todo perfecto",
		Token:   token,
	})
	return
}

func (c *authController) GoogleLogin(ctx *gin.Context) {
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

// el c * authController nos permite identificar a q pertenece la funcion
// es como si seria parte del objeto authController, se lo llama RECEPTOR
func (c *authController) getAuthCallBackFunction(ctx *gin.Context) {

	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		ctx.JSON(401, gin.H{"Error": "Unauthorized" + err.Error()})
		return
	}
	token, err := c.authService.LoginWithGoogle(user)
	if err != nil {
		ctx.JSON(500, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"Token": token})
}
