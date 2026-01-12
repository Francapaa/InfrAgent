package repositories

import (
	"context"
	"database/sql"
	"fmt"
	models "server/model"
	"time"
)

type AgentStorage interface {
	GetAgent(ctx context.Context, id string) (*models.Agent, error)
	GetAgentByClientId(ctx context.Context, clientId string) (*models.Agent, error)
	UpdateAgentState(ctx context.Context, id string) error
	SetAgentCooldown(ctx context.Context, id string, duration time.Duration) error
}

// QUERY PARA OBTENER EL AGENTE EN ESPECIFICO PARA NUESTRO CLIENTE (LUEGO OPTIMIZAMOS)
func (s *PostgresStorage) GetAgent(ctx context.Context, id string) (*models.Agent, error) {
	var a models.Agent

	err := s.db.QueryRowContext(ctx, `
		SELECT id , client_id, state, last_tick_at, cooldown_until, created_at, updated_at
		FROM agents 
		WHERE id = $1
	`, id).Scan(&a.ID, &a.ClientID, &a.State, &a.LastTickAt, &a.CooldownUntil, &a.CreatedAt, &a.UpdatedAt)

	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("agent not found with id: " + id)
	}
	return &a, nil
}

// ACA LO OBTENEMOS MEDIANTE EL ID DEL CLIENTE
func (s *PostgresStorage) GetAgentByClientId(ctx context.Context, clientId string) (*models.Agent, error) {

	var a models.Agent

	err := s.db.QueryRowContext(ctx, `
		SELECT id , client_id, state, last_tick_at, cooldown_until, created_at, updated_at
		FROM agents 
		WHERE client_id = $1
	`, clientId).Scan(&a.ID, &a.ClientID, &a.State, &a.LastTickAt, &a.CooldownUntil, &a.CreatedAt, &a.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("agent not found with client_id: " + clientId)
	}

	return &a, nil
}

// updateamos el state del agent cada vez que el agent realice acciones
func (s *PostgresStorage) UpdateAgentState(ctx context.Context, id string, stateAgent string) error {

	_, err := s.db.ExecContext(ctx, `
		UPDATE agents
		SET state = $1,
		last_tick_at = NOW(),
		updated_at = NOW()
		WHERE id = $2
	`, stateAgent, id)

	return err
}

// seteamos el cooldown, luego del cooldown se vuelve a ejecutar el tick
func (s *PostgresStorage) SetAgentCooldown(ctx context.Context, id string, duration time.Duration) error {
	coolDownUntil := time.Now().Add(duration)

	_, err := s.db.ExecContext(ctx, `
		UPDATE agents
		SET cooldown_until = $1,
		updated_at = NOW()
		WHERE id = $2
	`, coolDownUntil, id)

	return err
}
