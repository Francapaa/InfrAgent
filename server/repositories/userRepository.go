package repositories

import (
	"context"
	"database/sql"
	"fmt"
	models "server/model"
	"server/utils"
)

type ClientStorage interface {
	// Client operations
	CreateClient(ctx context.Context, user *models.Client) error
	GetClient(ctx context.Context, id string) (*models.Client, error)
	GetClientByAPIKey(ctx context.Context, apiKey string) (*models.Client, error)
	GetClientByEmail(ctx context.Context, email string) (*models.Client, error)
	UpdateClient(ctx context.Context, user *models.Client) error
}

type PostgresStorage struct {
	db *sql.DB
}

/*func NewPostgresStorage(connStr string) (*PostgresStorage, error){
	db, err := sql.Open
}*/

func (s *PostgresStorage) CreateClient(ctx context.Context, user *models.Client) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO clients (id, email, password, company_name, api_key_hash, web_hook_secret, web_hook_url, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`, user.ID, user.Email, user.Password, user.CompanyName, user.APIKeyHash, user.WebhookSecret, user.WebhookURL, user.CreatedAt, user.UpdatedAt)
	return err

}

func (s *PostgresStorage) GetClient(ctx context.Context, id string) (*models.Client, error) {

	var c models.Client

	err := s.db.QueryRowContext(ctx, `
	SELECT id, email, password, company_name, api_key_hash, web_hook_secret, web_hook_url, created_at, updated_at
	FROM clients 
	WHERE id = $1 
`, id).Scan(&c.ID, &c.Email, &c.CompanyName, &c.APIKeyHash, &c.WebhookSecret, &c.WebhookURL, &c.CreatedAt, &c.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("NO EXISTE ESE ID / DOESNT EXIST THESE ID")
	}
	return &c, nil
}

func (s *PostgresStorage) GetClientByAPIKey(ctx context.Context, APIKey string) (*models.Client, error) {

	rows, err := s.db.QueryContext(ctx, `
		SELECT id, email, password, company_name, api_key_hash, web_hook_secret, web_hook_url, created_at, updated_at
		FROM clients
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Client

		if err := rows.Scan(&c.ID, &c.Email, &c.CompanyName, &c.APIKeyHash, &c.WebhookSecret, &c.WebhookURL, &c.CreatedAt, &c.UpdatedAt); err != nil {
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
	SELECT id, email, password, company_name, api_key_hash, web_hook_secret, web_hook_url, created_at, updated_at
	FROM clients 
	WHERE email = $1 
`, email).Scan(&c.ID, &c.Email, &c.CompanyName, &c.APIKeyHash, &c.WebhookSecret, &c.WebhookURL, &c.CreatedAt, &c.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("NO EXISTE ESE EMAIL / DOESNT EXIST THESE EMAIL")
	}
	return &c, nil

}

func (s *PostgresStorage) UpdateClient(ctx context.Context, user *models.Client) error {

	_, err := s.db.ExecContext(ctx, `
		UPDATE clients
		SET metodo = $1,
		google_id = $2
		WHERE id = $3
	`, user.Metodo, user.GoogleID, user.ID)
	return err

}
