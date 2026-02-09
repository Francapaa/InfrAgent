package service

import (
	"context"
	"fmt"
	models "server/model"
	"server/repositories"
	"server/service/agent/llm" // go importa por path del modulo + carpetas
	service "server/service/exec"
	"time"

	"github.com/google/uuid"
)

// ACA VA A ESTAR TODA LA LOGICA RELACIONADA AL AGENTE, EL WORKFLOW PRINCIPAL VA A ESTAR ALMACENADO EN ESTE
// ARCHIVO

type AgentEngine struct {
	gemini  *llm.GeminiClient
	events  repositories.EventRepository
	actions repositories.ActionStorage
	agents  repositories.AgentStorage
	client  repositories.ClientStorage
}

func (e *AgentEngine) assembleContext(ctx context.Context, agent *models.Agent, events []models.Event) models.AgentRunContext {

	since := time.Now().Add(-1 * time.Hour)
	restartCount, _ := e.actions.CountActionsSince(context.Background(), agent.ID, "restart", since)

	return models.AgentRunContext{
		CurrentEvents:    events,
		RestartCountHour: restartCount,
		ClientConfig:     models.ClientConfig{},
	}

}

func (e *AgentEngine) RunTick(ctx context.Context, agentId string) error {
	agent, err := e.agents.GetAgent(ctx, agentId) // aca le damos el estado al agente
	if err != nil {
		return fmt.Errorf("error getting agent: %w", err)
	}
	events, err := e.events.GetPendingEvents(ctx, agentId) // aca cargamos los eventos pendientes que tenga
	if err != nil {
		return fmt.Errorf("error getting pending events: %w", err)
	}
	clientIDUUID, err := uuid.Parse(agent.ClientID)
	if err != nil {
		return fmt.Errorf("invalid client ID format: %w", err)
	}

	client, err := e.client.GetClient(ctx, clientIDUUID)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}

	if len(events) == 0 {
		return nil
	}

	runCtx := models.AgentRunContext{
		CurrentEvents:    events,
		RestartCountHour: 0,
		ClientConfig:     models.ClientConfig{},
	}

	decision, err := e.gemini.Decide(ctx, runCtx)

	if err != nil {
		return err
	}

	executor := &service.Executor{}
	result := executor.Execute(ctx, decision, agent, client) // en execute tiene que ir algun switch con las opciones
	e.actions.SaveAction(ctx, result)                        // guardar esa accion en la base de datos

	for _, ev := range events {
		e.events.MarkEventProcessed(ctx, ev.ID) // marcamos el evento como procesado
	}

	//settear nuevo CD

	return nil
}
