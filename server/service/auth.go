package service

import (
	"errors"
	"server/controller"
	"server/repositories" // llamando a la base de datos
	"strings"
	"github.com/markbates/goth/gothic"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db * Database
}


func (c * AuthService) loginLocal( email, password string)(string, error){
	if email == "" || password == ""{
		return "" , errors.New("Email y password son requeridos")
	}
	if len(password) < 8{
		return "", errors.New("La contraseña debe contener al menos 8 caracteres ")
	}
	if !strings.Contains(email, "@"){
		return "", errors.New("El email debe contener '@' ")
	}

	user, err := s.repo.FindUserByEmail(email)
		if err != nil{
			return "", errors.New("Email no registrado")
		}
	
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.password),
		[]byte(password),
	)

	if err != nil{
		return "", errors.New("Contraseña incorrecta.")
	}

	token:= generateJWT(user.id) //UNCION QUE GENERA JWT EN /UTILS

	return token, nil
}

func (c *AuthService) LoginWithGoogle (var user goth.User)(string, error){

/*

BUSCAR MAIL SI ES Q ESTÁ REGISTRADO, SI NO LO ESTA LO INSERTAMOS EN LA BD.
SI ESTÁ PERO NO TIENE PROVIDER == GOOGLE, SE LO ASIGNAMOS Y LE ASIGNAMOS EL GOOGLEID
SI ESTA TODO REGISTRADO DESDE GOOGLE, DIRECTAMENTE PASA CON JWT

*/

}