package service

import "errors"

// llamando a la base de datos

type authService struct {
	db *Database
}

func (c *authService) loginLocal(email, password string) (string, error) {
	if email == "" || password == "" {
		return "", errors.New("Email y password son requeridos")
	}
}
