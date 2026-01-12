package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	models "server/model"
	"time"
)

type EventStorage interface {
	SaveAction(ctx context.Context, action *models.Action) error
	GetRecentActions(ctx context.Context, agentID string, limit int) ([]models.Action, error)
	CountActionsSince(ctx context.Context, agentID, actionType string, since time.Time) (int, error)
}

func (s *PostgresStorage) SaveAction(ctx context.Context, action *models.Action) error {

	paramsJSON, err := json.Marshal(action.Params) // convierte  Go data structures (structs, maps, slices) into a JSON-formatted byte slice (a string representation)

	if err != nil {
		return fmt.Errorf("marshal result: %w", err)
	}

	resultJSON, err := json.Marshal(action.Result)
	if err != nil {
		return fmt.Errorf("marshal result: %w", err)
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO actions (id, agent_id, client_id, type, target, params, reasoning, confidence, status, result, executed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`, action.ID, action.AgentID, action.ClientID, action.Type, action.Target, paramsJSON,
		action.Reasoning, action.Confidence, action.Status, resultJSON, action.ExecutedAt, action.CreatedAt)

	return err
}

// return the recent actions of the AGENT
func (s *PostgresStorage) GetRecentActions(ctx context.Context, agentId string, limit int) ([]models.Action, err) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, agent_id, client_id, type, target, params, reasoning, confidence, status, result, executed_at, created_at
		FROM actions
		WHERE agent_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, agentId, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var actions []models.Action

	for rows.Next() {
		var a models.Action
		var paramsJSON, resultJSON []byte

		err := rows.Scan(&a.ID, &a.AgentID, &a.ClientID, &a.Type, &a.Target, &paramsJSON,
			&a.Reasoning, &a.Confidence, &a.Status, &resultJSON, &a.ExecutedAt, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(paramsJSON, &a.Params)
		json.Unmarshal(resultJSON, &a.Result)

		actions = append(actions, a)
	}

	return actions, nil
}

func (s *PostgresStorage) CountActionsSince(ctx context.Context, agentID, actionType string, since time.Time) (int, error) {

}
