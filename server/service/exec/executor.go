package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	models "server/model"
	"time"
)

type Executor struct{}

func (e *Executor) Execute(ctx context.Context, decision *models.LLMDecision, agent *models.Agent, client *models.Client) *models.Action {

	action := &models.Action{
		AgentID:   agent.ID,
		Type:      decision.Action,
		Target:    decision.Target,
		Reasoning: decision.Reasoning,
		Status:    "pending",
	}

	payload, _ := json.Marshal(decision)

	req, _ := http.NewRequestWithContext(ctx, "POST", client.WebhookURL, bytes.NewBuffer(payload))
	req.Header.Set("Content-type", "application/json")

	cliente := &http.Client{Timeout: 10 * time.Second}
	response, err := cliente.Do(req)

	if err != nil || response.StatusCode != 200 {
		action.Status = "failed"
		action.Result = map[string]interface{}{"error": "client_webhook_unreacheable"}
	} else {
		action.Status = "success"
	}

	return action
}
