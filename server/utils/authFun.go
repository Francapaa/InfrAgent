package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"golang.org/x/crypto/bcrypt"
)

const (
	key    = "z£=AO#ylytj^j}Un0rI8q;D<7G8e1op"
	maxAge = 86400 * 30
	isProd = false
)

func NewAuth() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando variables de entorno")
	}

	GoogleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	GoogleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd

	gothic.Store = store

	goth.UseProviders(
		google.New(GoogleClientId, GoogleClientSecret, "http://localhost:8080/auth/google/callback"),
	)
}

func GenerateJWT(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(key))
}

func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "agent_key_" + base64.URLEncoding.EncodeToString(bytes), nil
}

func WebHookSecret() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "web_hook_secret_" + base64.URLEncoding.EncodeToString(bytes), nil
}

func HashAPIKey(apiKey string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash api key: %w", err)
	}

	fmt.Printf("Hashed API KEY: %s", string(hash))
	return string(hash), nil
}

func IsValidaAPIKey(apiKey, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(apiKey))

	return err == nil
}

func CompareHashedAPIKey(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

/*

ACA ESTAN CONVIVIENDO TODAS LAS FUNCIONES QUE SERAN DE UTILIDAD A LA HORA DEL LOGIN DE LOS USUARIOS
IN THIS FILE WE HAVE ALL FUNCTION WHICH ARE USED TO CREATE THE SDK IN THE LOGIN

*/
