package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func GenerateJWT(userId uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId.String(),
		"exp": time.Now().Add(time.Hour * 8).Unix(),
	})
	return token.SignedString([]byte(key))
}

// Claims representa los claims del JWT
type Claims struct {
	UserID string `json:"sub"`
	jwt.RegisteredClaims
}

// ValidateJWT valida un token JWT y retorna los claims
func ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
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
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func HashAPIKey(apiKey string) string {
	h := sha256.New()
	h.Write([]byte(apiKey))
	return hex.EncodeToString(h.Sum(nil))
}

func IsValidaAPIKey(apiKey, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(apiKey))

	return err == nil
}

func CompareHashedAPIKey(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func HashPassword(password string) (string, error) {
	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(password), 256)
	if err != nil {
		return "", err
	}
	return string(passwordHashed), nil
}

// GenerateUUID genera un UUID v4 sin dependencias externas

/*

ACA ESTAN CONVIVIENDO TODAS LAS FUNCIONES QUE SERAN DE UTILIDAD A LA HORA DEL LOGIN DE LOS USUARIOS
IN THIS FILE WE HAVE ALL FUNCTION WHICH ARE USED TO CREATE THE SDK IN THE LOGIN

*/
