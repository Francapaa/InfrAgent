package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	models "server/model"
	"strings"

	"google.golang.org/genai"
)

// GeminiClient es el cliente para interactuar con Gemini
type GeminiClient struct {
	client *genai.Client
	model  string
}

// ConnectionToGeminiLLM crea una nueva conexión con Gemini
func ConnectionToGeminiLLM(apikey, model string) *GeminiClient {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apikey,
	})
	if err != nil {
		log.Fatal("No pudimos conectarnos con GEMINI: ", err)
	}

	if model == "" {
		model = "gemini-2.0-flash-exp"
	}

	return &GeminiClient{
		client: client,
		model:  model,
	}
}

// Decide llama a Gemini, parsea la respuesta y valida la decisión
func (g *GeminiClient) Decide(ctx context.Context, agentCtx models.AgentRunContext) (*models.LLMDecision, error) {
	// 1. CREAR PROMPT
	prompt, err := g.CreatePrompt(agentCtx)
	if err != nil {
		return nil, fmt.Errorf("error creating prompt: %w", err)
	}

	log.Printf("[Gemini] Enviando prompt al modelo %s", g.model)

	// 2. LLAMAR A GEMINI
	result, err := g.client.Models.GenerateContent(
		ctx,
		g.model,
		genai.Text(prompt),
		&genai.GenerateContentConfig{
			Temperature: genai.Ptr(float32(0.7)), // Balance entre creatividad y precisión
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error llamando a Gemini: %w", err)
	}

	// 3. EXTRAER RESPUESTA
	responseText := result.Text()
	log.Printf("[Gemini] Respuesta recibida: %s", responseText)

	// 4. PARSEAR JSON
	decision, err := g.ParseResponse(responseText)
	if err != nil {
		return nil, fmt.Errorf("error parseando respuesta: %w", err)
	}

	// 5. VALIDAR DECISIÓN
	if err := g.ValidateDecision(decision, agentCtx); err != nil {
		log.Printf("[Gemini] Decisión inválida: %v", err)
		// Retornar decisión segura (wait) si la validación falla
		return &models.LLMDecision{
			Action:       "wait",
			Target:       "",
			Reasoning:    fmt.Sprintf("Decisión original rechazada: %v", err),
			Confidence:   0.0,
			ShouldNotify: true,
		}, nil
	}

	log.Printf("[Gemini] ✅ Decisión válida: %s (target=%s, confidence=%.2f)",
		decision.Action, decision.Target, decision.Confidence)

	return decision, nil
}

// CreatePrompt construye el prompt detallado para Gemini
func (g *GeminiClient) CreatePrompt(agentCtx models.AgentRunContext) (string, error) {
	var sb strings.Builder

	sb.WriteString("Eres un agente autónomo de monitoreo de infraestructura.\n")
	sb.WriteString("Tu trabajo es analizar incidentes y decidir la mejor acción.\n\n")

	// ============================================================
	// SITUACIÓN ACTUAL
	// ============================================================
	sb.WriteString("## SITUACIÓN ACTUAL\n")
	if len(agentCtx.CurrentEvents) == 0 {
		sb.WriteString("No hay eventos pendientes.\n")
	} else {
		for i, event := range agentCtx.CurrentEvents {
			sb.WriteString(fmt.Sprintf("%d. Tipo: %s\n", i+1, event.Type))
			sb.WriteString(fmt.Sprintf("   Servicio: %s\n", event.Service))
			sb.WriteString(fmt.Sprintf("   Severidad: %s\n", event.Severity))

			if len(event.Data) > 0 {
				dataJSON, _ := json.Marshal(event.Data)
				sb.WriteString(fmt.Sprintf("   Detalles: %s\n", string(dataJSON)))
			}
			sb.WriteString("\n")
		}
	}

	// ============================================================
	// HISTORIAL RECIENTE
	// ============================================================
	sb.WriteString("## HISTORIAL RECIENTE\n")
	sb.WriteString(fmt.Sprintf("Reinicios en la última hora: %d\n", agentCtx.RestartCountHour))

	if len(agentCtx.RecentActions) == 0 {
		sb.WriteString("No hay acciones previas.\n")
	} else {
		sb.WriteString("Últimas 5 acciones:\n")
		limit := 5
		if len(agentCtx.RecentActions) < limit {
			limit = len(agentCtx.RecentActions)
		}
		for i := 0; i < limit; i++ {
			action := agentCtx.RecentActions[i]
			sb.WriteString(fmt.Sprintf("- %s en %s (%s): %s\n",
				action.Type, action.Target, action.Status, action.Reasoning))
		}
	}
	sb.WriteString("\n")

	// ============================================================
	// ESTADO DE SERVICIOS
	// ============================================================
	sb.WriteString("## ESTADO DE SERVICIOS\n")
	if len(agentCtx.ServiceHealth) == 0 {
		sb.WriteString("Estado desconocido (sin datos)\n")
	} else {
		for service, status := range agentCtx.ServiceHealth {
			sb.WriteString(fmt.Sprintf("- %s: %s\n", service, status))
		}
	}
	sb.WriteString("\n")

	// ============================================================
	// REGLAS Y LÍMITES DEL CLIENTE
	// ============================================================
	sb.WriteString("## REGLAS Y LÍMITES\n")
	sb.WriteString(fmt.Sprintf("- Máximo de reinicios por hora: %d\n", agentCtx.ClientConfig.MaxRestartsPerHour))
	sb.WriteString(fmt.Sprintf("- Reinicios actuales: %d\n", agentCtx.RestartCountHour))
	sb.WriteString(fmt.Sprintf("- Acciones permitidas: %v\n", agentCtx.ClientConfig.AllowedActions))
	sb.WriteString(fmt.Sprintf("- Notificar al usuario en el reinicio #%d\n", agentCtx.ClientConfig.NotifyOnNthRestart))
	sb.WriteString("\n")

	// ============================================================
	// ACCIONES DISPONIBLES
	// ============================================================
	sb.WriteString("## ACCIONES DISPONIBLES\n")
	sb.WriteString("1. restart - Reinicia un servicio (úsalo con moderación)\n")
	sb.WriteString("2. notify - Alerta al dueño (para problemas críticos)\n")
	sb.WriteString("3. wait - No hacer nada y observar (cuando no estés seguro)\n")
	sb.WriteString("4. scale - Escalar réplicas hacia arriba (para alta carga)\n")
	sb.WriteString("5. rollback - Volver a versión anterior (si deploy reciente falló)\n")
	sb.WriteString("\n")

	// ============================================================
	// FORMATO DE RESPUESTA
	// ============================================================
	sb.WriteString("## TU TAREA\n")
	sb.WriteString("Analiza la situación y decide la mejor acción.\n")
	sb.WriteString("Responde SOLO con un objeto JSON (sin markdown, sin texto adicional):\n\n")

	sb.WriteString("{\n")
	sb.WriteString(`  "action": "restart",` + "\n")
	sb.WriteString(`  "target": "payments-api",` + "\n")
	sb.WriteString(`  "params": {},` + "\n")
	sb.WriteString(`  "reasoning": "Primera falla detectada, dependencias están saludables, seguro reiniciar",` + "\n")
	sb.WriteString(`  "confidence": 0.85,` + "\n")
	sb.WriteString(`  "should_notify": false` + "\n")
	sb.WriteString("}\n\n")

	// ============================================================
	// REGLAS IMPORTANTES
	// ============================================================
	sb.WriteString("REGLAS IMPORTANTES:\n")
	sb.WriteString("- NUNCA reinicies más del máximo permitido por hora\n")
	sb.WriteString("- Si el contador de reinicios está cerca del límite, prefiere 'notify' o 'wait'\n")
	sb.WriteString("- Siempre explica tu razonamiento claramente\n")
	sb.WriteString("- Usa 'wait' cuando la situación no esté clara\n")
	sb.WriteString("- La confianza debe estar entre 0.0 y 1.0\n")
	sb.WriteString("- Solo usa acciones de la lista permitida\n")
	sb.WriteString("- Si un servicio ya fue reiniciado varias veces, probablemente necesite intervención humana\n")

	return sb.String(), nil
}

// ParseResponse extrae y parsea el JSON de la respuesta de Gemini
func (g *GeminiClient) ParseResponse(responseText string) (*models.LLMDecision, error) {
	// Limpiar respuesta (remover markdown si existe)
	cleaned := strings.TrimSpace(responseText)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)

	// Parsear JSON
	var decision models.LLMDecision
	if err := json.Unmarshal([]byte(cleaned), &decision); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w (response: %s)", err, cleaned)
	}

	// Inicializar params si es nil
	if decision.Params == nil {
		decision.Params = make(map[string]interface{})
	}

	return &decision, nil
}

// ValidateDecision valida que la decisión sea segura y siga las reglas
func (g *GeminiClient) ValidateDecision(decision *models.LLMDecision, ctx models.AgentRunContext) error {
	// 1. Validar que la acción esté en la lista permitida
	allowed := false
	for _, allowedAction := range ctx.ClientConfig.AllowedActions {
		if decision.Action == allowedAction {
			allowed = true
			break
		}
	}
	if !allowed {
		return fmt.Errorf("acción '%s' no está en la lista permitida: %v",
			decision.Action, ctx.ClientConfig.AllowedActions)
	}

	// 2. Validar límite de reinicios
	if decision.Action == "restart" {
		if ctx.RestartCountHour >= ctx.ClientConfig.MaxRestartsPerHour {
			return fmt.Errorf("límite de reinicios excedido (%d/%d)",
				ctx.RestartCountHour, ctx.ClientConfig.MaxRestartsPerHour)
		}
	}

	// 3. Validar confianza
	if decision.Confidence < 0 || decision.Confidence > 1 {
		return fmt.Errorf("confianza fuera de rango: %.2f (debe ser 0.0-1.0)", decision.Confidence)
	}

	// 4. Validar que target esté especificado para acciones que lo requieren
	if decision.Action != "wait" && decision.Action != "notify" {
		if decision.Target == "" {
			return fmt.Errorf("target requerido para acción '%s'", decision.Action)
		}
	}

	// 5. Validar que reasoning no esté vacío
	if decision.Reasoning == "" {
		return fmt.Errorf("reasoning es obligatorio")
	}

	return nil
}
