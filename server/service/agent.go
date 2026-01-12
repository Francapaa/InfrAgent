package service

import (
	"context"
	"server/agent/llm"
	"server/repositories"
)

// ACA VA A ESTAR TODA LA LOGICA RELACIONADA AL AGENTE, EL WORKFLOW PRINCIPAL VA A ESTAR ALMACENADO EN ESTE
// ARCHIVO

type AgentEngine struct {
	gemini  *llm.GeminiClient
	events  *repositories.EventRepository
	actions repositories.ActionStorage
	agents  *repositories.AgentStorage
}

func (e *AgentEngine) RunTick(ctx context.Context, agentId string) error {
	agent, _ := e.agents.GetAgent(ctx, agentId)
	events, _ := e.events.ReturnGetPendingEvents(ctx, agentId)

	if len(events) == 0 {
		return nil
	}

	runCtx := e.assembleContext(agent, events)

	decision, err := e.gemin.Decide(ctx, runCtx)

	if err != nil {
		return nil
	}

	result := e.executor.Execute(decision)

	e.actions.SaveAction(ctx, result)

	for _, ev := range events {
		e.events.MarkEventProcessed(ctx, ev.id)
	}

	return nil
}
