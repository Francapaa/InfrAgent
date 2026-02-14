package service

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
	req.Header.Set("X-Agent-Signature", signPayload(payload, client.WebhookSecret))

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

func signPayload(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}
