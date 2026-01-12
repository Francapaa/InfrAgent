package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	models "server/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

var eventRepo *EventRepository

func InitEventRepository(db *pgxpool.Pool) {
	eventRepo = &EventRepository{DB: db}
}

type EventRepository struct {
	DB *pgxpool.Pool
}

type EventStorage interface {
	CreateEvent(ctx context.Context, event *models.Event) error
	GetPendingEvents(ctx context.Context, agentId string) ([]models.Event, error)
	MarkEventProcessed(ctx context.Context, eventId string) error
}

func CreateEvent(ctx context.Context, event *models.Event) error {
	return eventRepo.CreateEvent(ctx, event)
}

func (r *EventRepository) CreateEvent(ctx context.Context, event *models.Event) error {
	dataJSON, err := json.Marshal(event.Data)
	if err != nil {
		return fmt.Errorf("marshal event data: %w", err)
	}

	_, err = r.DB.Exec(ctx, `
		INSERT INTO events (id, client_id, agent_id, type, service, severity, data, processed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, event.ID, event.ClientID, event.AgentID, event.Type, event.Service, event.Severity, dataJSON, event.ProcessedAt, event.CreatedAt)
	return err
}

func ReturnGetPendingEvents(ctx context.Context, agentId string) ([]models.Event, error) {
	return eventRepo.GetPendingEvents(ctx, agentId)
}

func (r *EventRepository) GetPendingEvents(ctx context.Context, agentId string) ([]models.Event, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT id, client_id, agent_id, type, service, severity, data, processed_at, created_at
		FROM events
		WHERE agent_id = $1 AND processed_at IS NULL
		ORDER BY created_at
		LIMIT 50
	`, agentId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		var e models.Event
		var dataJSON []byte

		err := rows.Scan(&e.ID, &e.ClientID, &e.AgentID, &e.Type, &e.Service, &e.Severity, &dataJSON, &e.ProcessedAt, &e.CreatedAt)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(dataJSON, &e.Data); err != nil {
			return nil, fmt.Errorf("unmarshal event data: %w", err)
		}

		events = append(events, e)
	}

	return events, nil
}

func ReturnMarkEventProcessed(ctx context.Context, eventId string) error {
	return eventRepo.MarkEventProcessed(ctx, eventId)
}

func (r *EventRepository) MarkEventProcessed(ctx context.Context, eventId string) error {
	_, err := r.DB.Exec(ctx, `
		UPDATE events
		SET processed_at = NOW()
		WHERE id = $1
	`, eventId)
	return err
}
