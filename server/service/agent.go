package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	models "server/model"
	"server/repositories"
	"server/service/agent/llm"
	service "server/service/exec"
	"server/ws"
	"time"

	"github.com/google/uuid"
)

// AgentEngine es el motor principal del agente autónomo
type AgentEngine struct {
	Gemini  *llm.GeminiClient
	Events  repositories.EventStorage
	Actions repositories.ActionStorage
	Agents  repositories.AgentStorage
	Client  repositories.ClientStorage
}

func (e *AgentEngine) assembleContext(ctx context.Context, agent *models.Agent, events []models.Event) models.AgentRunContext {
	since := time.Now().Add(-1 * time.Hour)
	restartCount, _ := e.Actions.CountActionsSince(context.Background(), agent.ID, "restart", since)

	return models.AgentRunContext{
		CurrentEvents:    events,
		RestartCountHour: restartCount,
		ClientConfig:     models.ClientConfig{},
	}
}

func (e *AgentEngine) RunTick(ctx context.Context, agentId string) error {
	log.Printf("[AgentEngine] Iniciando tick para agente %s", agentId)

	agent, err := e.Agents.GetAgent(ctx, agentId)
	if err != nil {
		return fmt.Errorf("error getting agent: %w", err)
	}

	// Enviar estado inicial: analyzing
	e.broadcastState(agent, "analyzing", "Analizando eventos pendientes...", nil, nil)

	events, err := e.Events.GetPendingEvents(ctx, agentId)
	if err != nil {
		return fmt.Errorf("error getting pending events: %w", err)
	}

	clientIDUUID, err := uuid.Parse(agent.ClientID)
	if err != nil {
		return fmt.Errorf("invalid client ID format: %w", err)
	}

	client, err := e.Client.GetClient(ctx, clientIDUUID)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}

	// Si no hay eventos, solo actualizar estado a idle
	if len(events) == 0 {
		log.Printf("[AgentEngine] No hay eventos pendientes para agente %s", agentId)
		e.Agents.UpdateAgentState(ctx, agentId, "idle")
		e.broadcastState(agent, "idle", "Monitoreando servicios...", nil, nil)
		return nil
	}

	log.Printf("[AgentEngine] Encontrados %d eventos pendientes para agente %s", len(events), agentId)

	// Obtener acciones recientes para el contexto
	recentActions, _ := e.Actions.GetRecentActions(ctx, agentId, 5)

	runCtx := models.AgentRunContext{
		CurrentEvents:    events,
		RecentActions:    recentActions,
		RestartCountHour: 0,
		ServiceHealth:    make(map[string]string),
		ClientConfig:     models.ClientConfig{},
	}

	// Calcular reinicios en la última hora
	since := time.Now().Add(-1 * time.Hour)
	restartCount, _ := e.Actions.CountActionsSince(ctx, agentId, "restart", since)
	runCtx.RestartCountHour = restartCount

	// Actualizar estado a executing antes de consultar LLM
	e.Agents.UpdateAgentState(ctx, agentId, "executing")
	e.broadcastState(agent, "executing", fmt.Sprintf("Consultando LLM sobre %d eventos...", len(events)), events, recentActions)

	// Consultar al LLM
	decision, err := e.Gemini.Decide(ctx, runCtx)
	if err != nil {
		log.Printf("[AgentEngine] Error consultando LLM: %v", err)
		e.Agents.UpdateAgentState(ctx, agentId, "error")
		e.broadcastState(agent, "error", fmt.Sprintf("Error: %v", err), events, recentActions)
		return err
	}

	log.Printf("[AgentEngine] Decisión del LLM: %s (target=%s, confidence=%.2f)",
		decision.Action, decision.Target, decision.Confidence)

	// Ejecutar la decisión
	executor := &service.Executor{}
	action := executor.Execute(ctx, decision, agent, client)
	action.ClientID = agent.ClientID
	action.ID = uuid.New().String()
	action.CreatedAt = time.Now()
	action.ExecutedAt = &action.CreatedAt

	// Guardar la acción
	if err := e.Actions.SaveAction(ctx, action); err != nil {
		log.Printf("[AgentEngine] Error guardando acción: %v", err)
	}

	// Marcar eventos como procesados
	for _, ev := range events {
		if err := e.Events.MarkEventProcessed(ctx, ev.ID); err != nil {
			log.Printf("[AgentEngine] Error marcando evento %s como procesado: %v", ev.ID, err)
		}
	}

	// Obtener acciones actualizadas para el broadcast
	updatedActions, _ := e.Actions.GetRecentActions(ctx, agentId, 10)

	// Actualizar estado final
	finalStatus := "idle"
	finalTask := fmt.Sprintf("Acción completada: %s %s", decision.Action, decision.Target)
	if action.Status == "failed" {
		finalStatus = "error"
		finalTask = fmt.Sprintf("Error ejecutando %s: %v", decision.Action, action.Result)
	}

	e.Agents.UpdateAgentState(ctx, agentId, finalStatus)
	e.broadcastState(agent, finalStatus, finalTask, events, updatedActions)

	// Settear cooldown si es necesario
	if decision.Action == "restart" {
		e.Agents.SetAgentCooldown(ctx, agentId, 5*time.Minute)
		log.Printf("[AgentEngine] Agente %s en cooldown por 5 minutos", agentId)
	}

	log.Printf("[AgentEngine] Tick completado para agente %s", agentId)
	return nil
}

// broadcastState envía el estado actual del agente a todos los clientes WebSocket conectados
func (e *AgentEngine) broadcastState(agent *models.Agent, status string, currentTask string, events []models.Event, actions []models.Action) {
	// Convertir eventos
	var eventInfos []models.EventInfo

	for _, ev := range events {
		eventInfos = append(eventInfos, models.EventInfo{
			ID:          ev.ID,
			Type:        ev.Type,
			Service:     ev.Service,
			Severity:    ev.Severity,
			Data:        ev.Data,
			ProcessedAt: ev.ProcessedAt,
			CreatedAt:   ev.CreatedAt,
		})
	}

	// Convertir acciones
	var actionInfos []models.ActionInfo

	for _, act := range actions {
		description := act.Reasoning
		if len(description) > 100 {
			description = description[:100] + "..."
		}
		actionInfos = append(actionInfos, models.ActionInfo{
			ID:          act.ID,
			Type:        act.Type,
			Target:      act.Target,
			Status:      act.Status,
			Reasoning:   act.Reasoning,
			Confidence:  act.Confidence,
			ExecutedAt:  act.ExecutedAt,
			CreatedAt:   act.CreatedAt,
			Description: description,
		})
	}

	// Generar métricas simuladas basadas en eventos
	metrics := e.generateSimulatedMetrics(events)

	msg := models.WebSocketMessage{
		Agents: []models.AgentInfo{{
			ID:            agent.ID,
			ClientID:      agent.ClientID,
			State:         agent.State,
			LastTickAt:    agent.LastTickAt,
			CooldownUntil: agent.CooldownUntil,
		}},
		Events:      eventInfos,
		Actions:     actionInfos,
		Status:      status,
		CurrentTask: currentTask,
		Metrics:     metrics,
		Timestamp:   time.Now().Format(time.RFC3339),
		ClientID:    agent.ClientID,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[AgentEngine] Error serializando mensaje WebSocket: %v", err)
		return
	}

	ws.BroadcastMessage(jsonData)
	log.Printf("[AgentEngine] Estado broadcasteado: %s - %s", status, currentTask)
}

// generateSimulatedMetrics genera métricas simuladas basadas en los eventos
func (e *AgentEngine) generateSimulatedMetrics(events []models.Event) models.MetricsInfo {
	// Contar errores en eventos
	errorsCount := 0
	for _, ev := range events {
		if ev.Severity == "critical" || ev.Severity == "error" || ev.Type == "error_spike" {
			errorsCount++
		}
	}

	// Generar valores simulados basados en la cantidad de eventos
	cpuUsage := 20.0 + float64(len(events))*5.0
	if cpuUsage > 95.0 {
		cpuUsage = 95.0
	}

	memoryUsage := 30.0 + float64(errorsCount)*10.0
	if memoryUsage > 90.0 {
		memoryUsage = 90.0
	}

	return models.MetricsInfo{
		CpuUsage:          cpuUsage,
		MemoryUsage:       memoryUsage,
		ActiveConnections: 100 + len(events)*50,
		ErrorsDetected:    errorsCount,
	}
}

// GetAgentState obtiene el estado completo de un agente para mostrar en el dashboard
func (e *AgentEngine) GetAgentState(ctx context.Context, agentID string) (*models.WebSocketMessage, error) {
	agent, err := e.Agents.GetAgent(ctx, agentID)
	if err != nil {
		return nil, err
	}

	events, err := e.Events.GetPendingEvents(ctx, agentID)
	if err != nil {
		return nil, err
	}

	actions, err := e.Actions.GetRecentActions(ctx, agentID, 10)
	if err != nil {
		return nil, err
	}

	// Convertir a formato de respuesta
	var eventInfos []models.EventInfo
	for _, ev := range events {
		eventInfos = append(eventInfos, models.EventInfo{
			ID:          ev.ID,
			Type:        ev.Type,
			Service:     ev.Service,
			Severity:    ev.Severity,
			Data:        ev.Data,
			ProcessedAt: ev.ProcessedAt,
			CreatedAt:   ev.CreatedAt,
		})
	}

	var actionInfos []models.ActionInfo
	for _, act := range actions {
		description := act.Reasoning
		if len(description) > 100 {
			description = description[:100] + "..."
		}
		actionInfos = append(actionInfos, models.ActionInfo{
			ID:          act.ID,
			Type:        act.Type,
			Target:      act.Target,
			Status:      act.Status,
			Reasoning:   act.Reasoning,
			Confidence:  act.Confidence,
			ExecutedAt:  act.ExecutedAt,
			CreatedAt:   act.CreatedAt,
			Description: description,
		})
	}

	metrics := e.generateSimulatedMetrics(events)

	return &models.WebSocketMessage{
		Agents: []models.AgentInfo{{
			ID:            agent.ID,
			ClientID:      agent.ClientID,
			State:         agent.State,
			LastTickAt:    agent.LastTickAt,
			CooldownUntil: agent.CooldownUntil,
		}},
		Events:      eventInfos,
		Actions:     actionInfos,
		Status:      agent.State,
		CurrentTask: "Agente activo",
		Metrics:     metrics,
		Timestamp:   time.Now().Format(time.RFC3339),
		ClientID:    agent.ClientID,
	}, nil
}
