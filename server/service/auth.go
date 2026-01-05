package service

import (
	"errors"
	"fmt"
	models "server/model"
	"server/repositories" // llamando a la base de datos
	"server/utils"
	"strings"

	"github.com/markbates/goth"
	"golang.org/x/crypto/bcrypt"
)

var repo *repositories.UserRepository

func initService(userRepo *repositories.UserRepository) {
	repo = userRepo
}

func LoginLocal(email, password string) (string, error) {
	if email == "" || password == "" {
		return "", errors.New("Email y password son requeridos")
	}
	if len(password) < 8 {
		return "", errors.New("La contraseña debe contener al menos 8 caracteres ")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("El email debe contener '@' ")
	}

	user, err := repo.FindUserByEmail(email)
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

	token, err := utils.GenerateJWT(user.ID.Hex()) //UNCION QUE GENERA JWT EN /UTILS
	if err != nil {
		return "", errors.New("Error al generar el token" + err.Error())
	}

	return token, nil
}

func LoginWithGoogle(gothUser goth.User) (string, error) {

	fmt.Println("Entro al login with google SERVICE")
	existingUser, err := repo.FindUserByEmail(gothUser.Email)

	if err != nil {
		newUser := models.UserRegister{
			Nombre:   gothUser.Name,
			Email:    gothUser.Email,
			GoogleID: gothUser.UserID,
			Metodo:   "google",
			Password: "",
		}
		err = repo.CreateUser(newUser)
		if err != nil {
			return "", errors.New("Error al crear el usuario")
		}
		token, err := utils.GenerateJWT(newUser.ID.Hex())
		if err != nil {
			return "", errors.New("Error al generar el token" + err.Error())
		}
		return token, nil
	}

	if existingUser.Metodo == "" || existingUser.Metodo != "google" {
		existingUser.Metodo = "google"
		existingUser.GoogleID = gothUser.UserID

		err = repo.UpdateUser(*existingUser)
		if err != nil {
			return "", errors.New("Error al actualizar el usuario" + err.Error())
		}
	}

	token, err := utils.GenerateJWT(existingUser.ID.Hex())
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
