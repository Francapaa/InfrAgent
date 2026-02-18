package service

import (
	"context"
	"fmt"
	models "server/model"
	"server/repositories"
)

// AgentDataService maneja la lógica de negocio para obtener datos del agente
type AgentDataService struct {
	agents repositories.AgentStorage
	engine *AgentEngine
}

// NewAgentDataService crea un nuevo servicio de datos del agente
func NewAgentDataService(agents repositories.AgentStorage, engine *AgentEngine) *AgentDataService {
	return &AgentDataService{
		agents: agents,
		engine: engine,
	}
}

// GetAgentStateForClient obtiene el estado completo del agente para un cliente específico
// Retorna el estado del agente o un error si no se encuentra
func (s *AgentDataService) GetAgentStateForClient(ctx context.Context, clientID string) (*models.WebSocketMessage, error) {
	// Obtener el agente asociado al cliente
	agent, err := s.agents.GetAgentByClientId(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("agente no encontrado para el cliente: %w", err)
	}

	// Obtener el estado completo del agente usando el engine
	state, err := s.engine.GetAgentState(ctx, agent.ID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener estado del agente: %w", err)
	}

	return state, nil
}

// GetAgentByClientId obtiene el agente directamente por ID del cliente
// Útil para operaciones que solo necesitan el agente sin su estado completo
func (s *AgentDataService) GetAgentByClientId(ctx context.Context, clientID string) (*models.Agent, error) {
	agent, err := s.agents.GetAgentByClientId(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("agente no encontrado: %w", err)
	}

	return &models.Agent{
		ID:            agent.ID,
		ClientID:      agent.ClientID,
		State:         agent.State,
		LastTickAt:    agent.LastTickAt,
		CooldownUntil: agent.CooldownUntil,
	}, nil
}

func (s *AgentDataService) GetLast30ActionsByAgent(ctx context.Context, clientID string) ([]models.Action, error) {

	actions, err := s.agents.GetLast30ActionsByAgent(ctx, clientID)

	if err != nil {
		return nil, fmt.Errorf("acciones no econtradas: %w", err)
	}

	return actions, nil
}
