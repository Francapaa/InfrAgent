package models

import (
	"time"

	"github.com/google/uuid"
)

// info que llega al endpoint cuando un usuario se registra localmente
type ClientRegister struct {
	Nombre      string `json:"nombre"`
	Email       string `json:"email"`
	CompanyName string `json:"company_name"`
	Password    string `json:"password"`
	WebhookURL  string `json:"webhook_url"`
}

// Client represents a registered user/company
type Client struct {
	ID            uuid.UUID `json:"id"`
	Nombre        string    `json:"nombre"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	CompanyName   string    `json:"company_name"`
	Metodo        string    `json:"metodo"` // este "metodo" nos dice si inicio con google id o LOCAL
	GoogleID      string    `json:"google_id,omitempty"`
	APIKeyHash    string    `json:"-"`
	WebhookSecret string    `json:"-"`
	WebhookURL    string    `json:"webhook_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type LoginResponse struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	Token         string `json:"token"`
	WebHookSecret string `json:"-"` //estos dos valores (webHooksecret y apiKey) se lo damos en txt plano
	ApiKey        string `json:"-"`
}

// Agent represents an autonomous agent managing a client's infrastructure
type Agent struct {
	ID            string     `json:"id"`
	ClientID      string     `json:"client_id"`
	State         string     `json:"state"` // "idle", "analyzing", "acting", "cooldown"
	LastTickAt    *time.Time `json:"last_tick_at"`
	CooldownUntil time.Time  `json:"cooldown_until"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// Event represents an incident or metric reported by the client's SDK
type Event struct {
	ID          string                 `json:"id"`
	ClientID    string                 `json:"client_id"`
	AgentID     string                 `json:"agent_id"`
	Type        string                 `json:"type"`         // "app_down", "high_cpu", "error_spike", etc
	Service     string                 `json:"service"`      // "api", "db", "worker", etc
	Severity    string                 `json:"severity"`     // "info", "warning", "critical"
	Data        map[string]interface{} `json:"data"`         // Flexible metadata
	ProcessedAt *time.Time             `json:"processed_at"` // nil = pending
	CreatedAt   time.Time              `json:"created_at"`
}

// Action represents a decision made and executed by the agent
type Action struct {
	ID          string                 `json:"id"`
	AgentID     string                 `json:"agent_id"`
	ClientID    string                 `json:"client_id"`
	Type        string                 `json:"type"`   // "restart", "notify", "wait", "scale", etc
	Target      string                 `json:"target"` // "api", "db", etc
	Params      map[string]interface{} `json:"params"`
	Reasoning   string                 `json:"reasoning"`   // Why the agent chose this
	Description string                 `json:"description"` //brief description of what the agent did
	Confidence  float64                `json:"confidence"`  // LLM confidence score
	Status      string                 `json:"status"`      // "pending", "success", "failed"
	Result      map[string]interface{} `json:"result"`      // Execution result
	ExecutedAt  *time.Time             `json:"executed_at"`
	CreatedAt   time.Time              `json:"created_at"`
}

// Notification represents an alert sent to the client
type Notification struct {
	ID        string     `json:"id"`
	ClientID  string     `json:"client_id"`
	ActionID  *string    `json:"action_id,omitempty"`
	Type      string     `json:"type"` // "email", "slack", "webhook"
	Recipient string     `json:"recipient"`
	Subject   string     `json:"subject,omitempty"`
	Body      string     `json:"body"`
	SentAt    *time.Time `json:"sent_at"`
	Status    string     `json:"status"` // "pending", "sent", "failed"
	CreatedAt time.Time  `json:"created_at"`
}

// AgentContext is the full context passed to the LLM for decision making
type AgentContext struct {
	CurrentEvents      []Event           `json:"current_events"`
	RecentActions      []Action          `json:"recent_actions"` // Last 10 actions
	RestartCountHour   int               `json:"restart_count_hour"`
	ServiceHealth      map[string]string `json:"service_health"` // service -> "up"/"down"
	ClientConfig       ClientConfig      `json:"client_config"`
	HistoricalPatterns []string          `json:"historical_patterns,omitempty"` // Future: ML insights
}

type AgentRunContext struct {
	CurrentEvents    []Event
	RecentActions    []Action
	RestartCountHour int
	ServiceHealth    map[string]string
	ClientConfig     ClientConfig
}

// ClientConfig represents the rules and limits for this client
type ClientConfig struct {
	MaxRestartsPerHour int      `json:"max_restarts_per_hour"`
	AllowedActions     []string `json:"allowed_actions"` // Whitelist
	NotifyOnNthRestart int      `json:"notify_on_nth_restart"`
	CooldownMinutes    int      `json:"cooldown_minutes"`
}

// LLMDecision represents the decision made by the LLM
// PENSAR: LLM DECISION PERTENECE A UN AGENTE? O COMO IDENTIFICAMOS ESA DECISION
type LLMDecision struct {
	Action       string                 `json:"action"`
	Target       string                 `json:"target"`
	Params       map[string]interface{} `json:"params,omitempty"`
	Reasoning    string                 `json:"reasoning"`
	Confidence   float64                `json:"confidence"`
	Alternative  string                 `json:"alternative,omitempty"` // Fallback plan
	ShouldNotify bool                   `json:"should_notify"`
}

// CompleteRegistrationRequest represents the request to complete registration after Google login
type CompleteRegistrationRequest struct {
	CompanyName string `json:"company_name" binding:"required"`
	WebhookURL  string `json:"webhook_url" binding:"required"`
}

// CompleteRegistrationResponse represents the response after completing registration
type CompleteRegistrationResponse struct {
	ClientID      uuid.UUID `json:"client_id"`
	APIKey        string    `json:"api_key"`
	WebhookSecret string    `json:"webhook_secret"`
}

type WebSocketMessage struct {
	Agents      []Agent     `json:"agents"`
	Events      []Event     `json:"events"`
	Actions     []Action    `json:"actions"`
	Status      string      `json:"status"`
	CurrentTask string      `json:"currentTask"`
	Metrics     MetricsInfo `json:"metrics"`
	Timestamp   string      `json:"timestamp"`
	ClientID    string      `json:"clientId,omitempty"`
}

type MetricsInfo struct {
	CpuUsage          float64 `json:"cpuUsage"`
	MemoryUsage       float64 `json:"memoryUsage"`
	ActiveConnections int     `json:"activeConnections"`
	ErrorsDetected    int     `json:"errorsDetected"`
}
