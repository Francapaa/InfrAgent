package service

import (
	"context"
	"errors"
	"fmt"
	models "server/model"
	"server/repositories" // llamando a la base de datos
	"server/utils"
	"strings"

	"github.com/markbates/goth"
	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	client repositories.ClientStorage
}

func (l *Login) LoginLocal(email, password string) (string, error) {
	if email == "" || password == "" {
		return "", errors.New("Email y password son requeridos")
	}
	if len(password) < 8 {
		return "", errors.New("La contraseña debe contener al menos 8 caracteres ")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("El email debe contener '@' ")
	}

	ctx := context.Background()

	user, err := l.client.GetClientByEmail(ctx, email)
	if err != nil {
		return "", errors.New("Email no registrado")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)

	if err != nil {
		return "", errors.New("Contraseña incorrecta.")
	}

	token, err := utils.GenerateJWT(user.ID) //FUNCION QUE GENERA JWT EN /UTILS
	if err != nil {
		return "", errors.New("Error al generar el token" + err.Error())
	}

	return token, nil
}

func (l *Login) LoginWithGoogle(gothUser goth.User) (string, error) {

	ctx := context.Background()

	fmt.Println("Entro al login with google SERVICE")
	existingUser, err := l.client.GetClientByEmail(ctx, gothUser.Email)

	if err != nil {
		newUser := &models.Client{
			Nombre:   gothUser.Name,
			Email:    gothUser.Email,
			GoogleID: gothUser.UserID,
			Metodo:   "google",
			Password: "",
		}
		apiKeyDeUsuario, err := utils.GenerateAPIKey()
		apiKeyHashed := utils.HashAPIKey(apiKeyDeUsuario)
		// METODOS PARA PODER CREAR LA API KEY PARA CADA USUARIO
		newUser.APIKeyHash = apiKeyHashed
		newUser.WebhookSecret, err = utils.WebHookSecret()
		err = l.client.CreateClient(ctx, newUser)
		if err != nil {
			return "", errors.New("Error al crear el usuario")
		}
		token, err := utils.GenerateJWT(newUser.ID)
		if err != nil {
			return "", errors.New("Error al generar el token" + err.Error())
		}
		//aca deberiamos devolver el usuario
		return token, nil
	}

	if existingUser.Metodo == "" || existingUser.Metodo != "google" {
		existingUser.Metodo = "google"
		existingUser.GoogleID = gothUser.UserID

		err = l.client.UpdateClient(ctx, existingUser)
		if err != nil {
			return "", errors.New("Error al actualizar el usuario" + err.Error())
		}
	}

	token, err := utils.GenerateJWT(existingUser.ID)
	if err != nil {
		return "", errors.New("Error al generar el token" + err.Error())
	}
	return token, nil

	/*

	   BUSCAR MAIL SI ES Q ESTÁ REGISTRADO, SI NO LO ESTA LO INSERTAMOS EN LA BD.
	   SI ESTÁ PERO NO TIENE PROVIDER == GOOGLE, SE LO ASIGNAMOS Y LE ASIGNAMOS EL GOOGLEID
	   SI ESTA TODO REGISTRADO DESDE GOOGLE, DIRECTAMENTE PASA CON JWT

	*/

}

func (l *Login) Register(userRegister models.ClientRegister) (models.LoginResponse, error) {

	ctx := context.Background()

	_, err := l.client.GetClientByEmail(ctx, userRegister.Email)

	if err == nil {
		return models.LoginResponse{}, errors.New("email already registered")
	}
	if !errors.Is(err, repositories.ErrUserNotFound) {
		return models.LoginResponse{}, err
	}

	passwordHashed, err := utils.HashPassword(userRegister.Password)

	if err != nil {
		return models.LoginResponse{}, errors.New("we cant hash your password, try again")
	}

	newUser := &models.Client{
		Nombre:      userRegister.Nombre,
		Email:       userRegister.Email,
		Password:    passwordHashed,
		CompanyName: userRegister.CompanyName,
		WebhookURL:  userRegister.WebhookURL,
	}

	newUser.Metodo = "local"
	newUser.WebhookSecret, err = utils.WebHookSecret()
	apiKey, err := utils.GenerateAPIKey()
	newUser.APIKeyHash = utils.HashAPIKey(apiKey)
	newUser.WebhookSecret, err = utils.WebHookSecret()
	tokenReturned, err := utils.GenerateJWT(newUser.ID)

	if err != nil {
		return models.LoginResponse{}, errors.New("ha ocurrido un error creando el webhooksecret")
	}

	l.client.CreateClient(ctx, newUser)

	return models.LoginResponse{
		Success:       true,
		Message:       "we have been created your profile successfully",
		Token:         tokenReturned,
		WebHookSecret: newUser.WebhookSecret,
		ApiKey:        apiKey,
	}, nil

}
