package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	models "server/model"
)

type clientStorage interface {
	GetClientConfig(ctx context.Context, agentID string) (models.ClientConfig, error)

	// Notification operations
	CreateNotification(ctx context.Context, notification *models.Notification) error
}

func (s *PostgresStorage) GetClientConfig(ctx context.Context, agentId string) (models.ClientConfig, error) {

	var cfg models.ClientConfig
	var allowedActionsJSON []byte

	err := s.db.QueryRowContext(ctx, `
	SELECT c.max_restarts_per_hour, c.allowed_actions, c.notify_on_nth_restart, c.cooldown_minutes
	FROM clients_config c 
	JOIN agents a on a.client_id = c.client_id
	where a.id = $1
`, agentId).Scan(&cfg.MaxRestartsPerHour, &cfg.AllowedActions, &cfg.NotifyOnNthRestart, &cfg.CooldownMinutes)

	if err == sql.ErrNoRows {
		return models.ClientConfig{
			MaxRestartsPerHour: 3,
			AllowedActions:     []string{"restart", "notify", "wait"},
			NotifyOnNthRestart: 3,
			CooldownMinutes:    5,
		}, nil
	}

	if err != nil {
		return cfg, err
	}

	json.Unmarshal(allowedActionsJSON, &cfg.AllowedActions)

	return cfg, nil
}

func (s *PostgresStorage) CreateNotification(ctx context.Context, notification *models.Notification) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO notifications (id, client_id, action_id, type, recipient, subject, body, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, notification.ID, notification.ClientID, notification.ActionID, notification.Type,
		notification.Recipient, notification.Subject, notification.Body, notification.Status, notification.CreatedAt)
	return err
}

// Close closes the database connection
func (s *PostgresStorage) Close() error {
	return s.db.Close()
}
