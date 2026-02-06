package service

import (
	"context"
	"errors"
	"fmt"
	models "server/model"
	"server/repositories"
	"server/utils"
	"strings"

	"github.com/markbates/goth"
)

type Login struct {
	client repositories.ClientStorage
}

func NewLogin(client repositories.ClientStorage) *Login {
	return &Login{
		client: client,
	}
}

func (l *Login) LoginWithGoogle(gothUser goth.User) (string, error) {
	ctx := context.Background()
	fmt.Println("ENTRE AL LOGIN WITH GOOGLE")
	// Validar que tenemos datos de Google
	if gothUser.UserID == "" {
		return "", errors.New("google user ID is required")
	}
	if gothUser.Email == "" {
		return "", errors.New("google email is required")
	}

	// PASO 1: Buscar si el GoogleID ya existe
	fmt.Println(gothUser.UserID)
	existingByGoogleID, err := l.client.GetClientByGoogleID(ctx, gothUser.UserID)
	fmt.Println(existingByGoogleID)
	if err == nil && existingByGoogleID != nil {
		fmt.Println("TE ENCONTRE POR GOOGLE ID")
		// Usuario ya registrado con Google - login normal
		token, err := utils.GenerateJWT(existingByGoogleID.ID)
		if err != nil {
			return "", errors.New("error generating token: " + err.Error())
		}

		return token, nil
	}

	// PASO 2: Buscar por Email
	existingByEmail, err := l.client.GetClientByEmail(ctx, gothUser.Email)
	if err == nil && existingByEmail != nil {
		fmt.Println("TE ENCONTRE POR MAIL")
		fmt.Println("SOS ESTE: ", existingByEmail)
		fmt.Println("HASTA ACA SOS")
		// El email existe pero no tiene GoogleID - hacer UPDATE
		existingByEmail.Metodo = "google"
		existingByEmail.GoogleID = gothUser.UserID

		err = l.client.UpdateClient(ctx, existingByEmail)
		if err != nil {
			return "", errors.New("error updating user with google ID: " + err.Error())
		}

		token, err := utils.GenerateJWT(existingByEmail.ID)
		if err != nil {
			return "", errors.New("error generating token: " + err.Error())
		}
		return token, nil
	}

	// PASO 3: No existe ni GoogleID ni Email - INSERT nuevo usuario
	// El ID lo genera PostgreSQL (gen_random_uuid())
	fmt.Println("TE VOY A INSERTAR PORQUE NO TE CONOZCO")
	newUser := &models.Client{
		// ID se deja vacío para que PostgreSQL genere el UUID
		CompanyName:   "",
		Email:         gothUser.Email,
		GoogleID:      gothUser.UserID,
		Metodo:        "google",
		APIKeyHash:    "",
		WebhookSecret: "",
		WebhookURL:    "",
		// WebhookURL queda vacío temporalmente, se completará luego
	}

	/*
		LOS DATOS MAS IMPORTANTES QUEDAN VACIOS YA QUE SE VAN A INGRESAR POSTERIORMENTE AL REGISTRO
	*/

	// Crear usuario en BD
	err = l.client.CreateClient(ctx, newUser)
	if err != nil {
		return "", errors.New("error creating user: " + err.Error())
	}

	// IMPORTANTE: No devolvemos apiKey ni webhookSecret todavía
	// El usuario debe completar el registro con webhook_url primero
	token, err := utils.GenerateJWT(newUser.ID)
	if err != nil {
		return "", errors.New("error generating token: " + err.Error())
	}

	return token, nil
}

func (l *Login) CompleteRegistration(ctx context.Context, userID string, companyName string, webhookURL string) (*models.CompleteRegistrationResponse, error) {
	// Validar company_name
	if companyName == "" {
		return nil, errors.New("company_name is required")
	}

	// Validar webhook_url
	if webhookURL == "" {
		return nil, errors.New("webhook_url is required")
	}

	// Validar que webhook_url sea HTTPS
	if !strings.HasPrefix(webhookURL, "https://") {
		return nil, errors.New("webhook_url must use HTTPS")
	}

	// Buscar usuario por ID
	user, err := l.client.GetClient(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Verificar que sea usuario de Google
	if user.Metodo != "google" {
		return nil, errors.New("only google users can complete registration this way")
	}

	// Verificar que no haya completado el registro antes
	if user.WebhookURL != "" {
		return nil, errors.New("registration already completed")
	}

	// Generar nueva API Key y Webhook Secret (por seguridad, diferentes a los temporales)
	apiKey, err := utils.GenerateAPIKey()
	if err != nil {
		return nil, errors.New("error generating API key")
	}

	apiKeyHashed := utils.HashAPIKey(apiKey)
	webhookSecret, err := utils.WebHookSecret()
	if err != nil {
		return nil, errors.New("error generating webhook secret")
	}

	// Actualizar usuario con company_name, webhook_url y nuevas credenciales
	user.CompanyName = companyName
	user.WebhookURL = webhookURL
	user.APIKeyHash = apiKeyHashed
	user.WebhookSecret = webhookSecret

	// Usar el método específico para completar registro
	err = l.client.UpdateClientComplete(ctx, user)
	if err != nil {
		return nil, errors.New("error updating user: " + err.Error())
	}

	return &models.CompleteRegistrationResponse{
		ClientID:      user.ID,
		APIKey:        apiKey,
		WebhookSecret: webhookSecret,
	}, nil
}

func (l *Login) GetUserByID(ctx context.Context, userID string) (*models.Client, error) {
	user, err := l.client.GetClient(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
