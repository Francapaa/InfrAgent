package service

import (
	"context"
	"errors"
	models "server/model"
	"server/repositories"
	"server/utils"
	"strings"

	"github.com/google/uuid"
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
	if gothUser.UserID == "" {
		return "", errors.New("google user ID is required")
	}
	if gothUser.Email == "" {
		return "", errors.New("google email is required")
	}

	existingByGoogleID, err := l.client.GetClientByGoogleID(ctx, gothUser.UserID)
	if err == nil && existingByGoogleID != nil {
		token, err := utils.GenerateJWT(existingByGoogleID.ID)
		if err != nil {
			return "", errors.New("error generating token: " + err.Error())
		}
		return token, nil
	}

	existingByEmail, err := l.client.GetClientByEmail(ctx, gothUser.Email)
	if err == nil && existingByEmail != nil {
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

	newUser := &models.Client{
		CompanyName:   "",
		Email:         gothUser.Email,
		GoogleID:      gothUser.UserID,
		Metodo:        "google",
		APIKeyHash:    "",
		WebhookSecret: "",
		WebhookURL:    "",
	}

	err = l.client.CreateClient(ctx, newUser)
	if err != nil {
		return "", errors.New("error creating user: " + err.Error())
	}

	token, err := utils.GenerateJWT(newUser.ID)
	if err != nil {
		return "", errors.New("error generating token: " + err.Error())
	}

	return token, nil
}

func (l *Login) CompleteRegistration(ctx context.Context, userID string, companyName string, webhookURL string) (*models.CompleteRegistrationResponse, error) {
	if companyName == "" {
		return nil, errors.New("company_name is required")
	}

	if webhookURL == "" {
		return nil, errors.New("webhook_url is required")
	}

	if !strings.HasPrefix(webhookURL, "https://") {
		return nil, errors.New("webhook_url must use HTTPS")
	}

	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := l.client.GetClient(ctx, userIDUUID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.Metodo != "google" {
		return nil, errors.New("only google users can complete registration this way")
	}

	if user.WebhookURL != "" {
		return nil, errors.New("registration already completed")
	}

	apiKey, err := utils.GenerateAPIKey()
	if err != nil {
		return nil, errors.New("error generating API key")
	}

	apiKeyHashed := utils.HashAPIKey(apiKey)
	webhookSecret, err := utils.WebHookSecret()
	if err != nil {
		return nil, errors.New("error generating webhook secret")
	}

	user.CompanyName = companyName
	user.WebhookURL = webhookURL
	user.APIKeyHash = apiKeyHashed
	user.WebhookSecret = webhookSecret

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
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := l.client.GetClient(ctx, userIDUUID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (l *Login) GetUserByEmail(ctx context.Context, email string) (*models.Client, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	user, err := l.client.GetClientByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
