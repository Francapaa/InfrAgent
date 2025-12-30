package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

type UserRegister struct {
	// El ID interno de Mongo. Usamos el tipo primitive.ObjectID
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Nombre string             `bson:"nombre" json:"nombre"`
	Email  string             `bson:"email" json:"email"`
	// El campo de Google que puede ser opcional
	GoogleID string `bson:"googleId,omitempty" json:"googleId,omitempty"`
	// El hash de la contraseña (no se envía en el JSON por seguridad)
	Password string `bson:"password,omitempty" json:"-"`
	Metodo   string `bson:"metodo" json:"metodo"` // "local" o "google"
}

type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
} // register no devuelve token, solo login x JWT
