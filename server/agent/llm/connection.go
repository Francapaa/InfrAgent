package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"server/model"

	"google.golang.org/genai"
)

type GeminiClient struct {
	client *genai.Client
	model  string
}

func ConnectionToGeminiLLM(apikey, model string) *GeminiClient {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apikey,
	})
	if err != nil {
		log.Fatal("No pudimos conextarnos con GEMINI ", err)
	}

	if model == "" {
		model = "gemini-2.0-flash-exp"
	}

	return &GeminiClient{
		client: client,
		model:  model,
	}

}

func (g *GeminiClient) Decide(ctx context.Context, agentCtx *model.AgentContext) (*model.LLMDecision, error) {
	prompt, err := g.CreatePrompt(ctx, agentCtx)
	if err != nil {
		return nil, fmt.Errorf("create prompt: %w", err)
	}

	// Simplified placeholder response for now
	decision := &model.LLMDecision{
		Action:     "notify",
		Target:     "",
		Reasoning:  "Placeholder decision",
		Confidence: 0.8,
	}

	return decision, nil
}

func (g *GeminiClient) CreatePrompt(ctx context.Context, agentCtx *model.AgentContext) (string, error) {
	var sb strings.Builder

	sb.WriteString("You are an infrastructure monitoring agent. Analyze the following context and decide on an action.\n\n")

	sb.WriteString("Current Events:\n")
	for _, e := range agentCtx.CurrentEvents {
		sb.WriteString(fmt.Sprintf("- %s on %s (%s): %v\n", e.Type, e.Service, e.Severity, e.Data))
	}

	sb.WriteString("\nRecent Actions:\n")
	for _, a := range agentCtx.RecentActions {
		sb.WriteString(fmt.Sprintf("- %s on %s: %s (confidence: %.2f)\n", a.Type, a.Target, a.Status, a.Confidence))
	}

	sb.WriteString(fmt.Sprintf("\nRestart Count Last Hour: %d\n", agentCtx.RestartCountHour))
	sb.WriteString(fmt.Sprintf("Client Config: Max restarts/hour: %d, Allowed actions: %v\n", agentCtx.ClientConfig.MaxRestartsPerHour, agentCtx.ClientConfig.AllowedActions))

	sb.WriteString("\nRespond with JSON: {\"action\": \"restart|notify|wait|scale\", \"target\": \"api|db|etc\", \"params\": {}, \"reasoning\": \"why\", \"confidence\": 0.0-1.0}\n")

	return sb.String(), nil
}
