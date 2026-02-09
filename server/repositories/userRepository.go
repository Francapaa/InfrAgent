package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	models "server/model"
	"server/utils"
	"time"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")

type ClientStorage interface {
	// Client operations
	CreateClient(ctx context.Context, user *models.Client) error
	GetClient(ctx context.Context, id uuid.UUID) (*models.Client, error)
	GetClientByAPIKey(ctx context.Context, apiKey string) (*models.Client, error)
	GetClientByEmail(ctx context.Context, email string) (*models.Client, error)
	GetClientByGoogleID(ctx context.Context, googleID string) (*models.Client, error)
	UpdateClient(ctx context.Context, user *models.Client) error
	UpdateClientComplete(ctx context.Context, user *models.Client) error
	FixClientID(ctx context.Context, email string, newID string) error
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (s *PostgresStorage) CreateClient(ctx context.Context, user *models.Client) error {
	var err error

	// Si el ID es nil (uuid.Nil), dejamos que PostgreSQL genere el UUID con DEFAULT
	if user.ID == uuid.Nil {
		err = s.db.QueryRowContext(ctx, `
			INSERT INTO clients (email, company_name, api_key_hash, webhook_secret, webhook_url, google_id, metodo, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id
		`, user.Email, user.CompanyName, user.APIKeyHash, user.WebhookSecret, user.WebhookURL, user.GoogleID, user.Metodo, time.Now(), time.Now()).Scan(&user.ID)
	} else {
		err = s.db.QueryRowContext(ctx, `
			INSERT INTO clients (id, nombre, email, company_name, api_key_hash, webhook_secret, webhook_url, google_id, metodo, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			RETURNING id
		`, user.ID, user.Nombre, user.Email, user.CompanyName, user.APIKeyHash, user.WebhookSecret, user.WebhookURL, user.GoogleID, user.Metodo, time.Now(), time.Now()).Scan(&user.ID)
	}

	return err
}

func (s *PostgresStorage) GetClient(ctx context.Context, id uuid.UUID) (*models.Client, error) {

	var c models.Client
	var idStr string

	idString := id.String()
	fmt.Printf("[Repository] Buscando cliente con ID: '%s'\n", idString)
	fmt.Printf("[Repository] Longitud del ID: %d\n", len(idString))
	fmt.Printf("[Repository] UUID objeto: %+v\n", id)

	err := s.db.QueryRowContext(ctx, `
	SELECT id::text, email, company_name ,metodo, google_id, api_key_hash, webhook_secret, webhook_url, created_at, updated_at
	FROM clients 
	WHERE id = $1::uuid
`, idString).Scan(&idStr, &c.Email, &c.CompanyName, &c.Metodo, &c.GoogleID, &c.APIKeyHash, &c.WebhookSecret, &c.WebhookURL, &c.CreatedAt, &c.UpdatedAt)

	fmt.Printf("[Repository] Error de query: %v\n", err)
	fmt.Printf("[Repository] ID encontrado en BD: '%s'\n", idStr)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("[Repository] No se encontraron filas\n")
			return nil, nil
		}
		return nil, err
	}

	// Parsear el ID string a uuid.UUID
	c.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing client ID: %w", err)
	}

	return &c, nil
}

func (s *PostgresStorage) GetClientByAPIKey(ctx context.Context, APIKey string) (*models.Client, error) {

	rows, err := s.db.QueryContext(ctx, `
		SELECT id, nombre, email, password, company_name, metodo, google_id, api_key_hash, web_hook_secret, web_hook_url, created_at, updated_at
		FROM clients
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Client

		if err := rows.Scan(&c.ID, &c.Nombre, &c.Email, &c.Password, &c.CompanyName, &c.Metodo, &c.GoogleID, &c.APIKeyHash, &c.WebhookSecret, &c.WebhookURL, &c.CreatedAt, &c.UpdatedAt); err != nil {
			continue
		}

		if utils.IsValidaAPIKey(APIKey, c.APIKeyHash) {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("invalid api key")
}

func (s *PostgresStorage) GetClientByEmail(ctx context.Context, email string) (*models.Client, error) {

	var c models.Client

	err := s.db.QueryRowContext(ctx, `
		SELECT id, nombre, email, password, company_name, metodo, google_id, api_key_hash, web_hook_secret, web_hook_url, created_at, updated_at
		FROM clients 
		WHERE email = $1 
	`, email).Scan(&c.ID, &c.Nombre, &c.Email, &c.Password, &c.CompanyName, &c.Metodo, &c.GoogleID, &c.APIKeyHash, &c.WebhookSecret, &c.WebhookURL, &c.CreatedAt, &c.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil

}

func (s *PostgresStorage) UpdateClient(ctx context.Context, user *models.Client) error {
	// Si el ID está vacío, no podemos actualizar (no sabemos qué registro actualizar)

	_, err := s.db.ExecContext(ctx, `
		UPDATE clients
		SET metodo = $1,
		google_id = $2
		WHERE id = $3
	`, user.Metodo, user.GoogleID, user.ID)
	return err

}

// FixClientID actualiza el ID de un cliente que fue creado sin ID (para migración de datos antiguos)
func (s *PostgresStorage) FixClientID(ctx context.Context, email string, newID string) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE clients
		SET id = $1
		WHERE email = $2 AND (id = '' OR id IS NULL)
	`, newID, email)
	return err
}

func (s *PostgresStorage) GetClientByGoogleID(ctx context.Context, googleID string) (*models.Client, error) {
	var c models.Client

	err := s.db.QueryRowContext(ctx, `
	SELECT id, email, company_name, metodo, google_id, api_key_hash, webhook_secret, webhook_url, created_at, updated_at
	FROM clients 
	WHERE google_id = $1 
`, googleID).Scan(&c.ID, &c.Email, &c.CompanyName, &c.Metodo, &c.GoogleID, &c.APIKeyHash, &c.WebhookSecret, &c.WebhookURL, &c.CreatedAt, &c.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (s *PostgresStorage) UpdateClientComplete(ctx context.Context, user *models.Client) error {

	_, err := s.db.ExecContext(ctx, `
		UPDATE clients
		SET company_name = $1,
		    webhook_url = $2,
		    api_key_hash = $3,
		    webhook_secret = $4,
		    updated_at = $5
		WHERE id = $6
	`, user.CompanyName, user.WebhookURL, user.APIKeyHash, user.WebhookSecret, time.Now(), user.ID)

	return err
}
