package models

type LLMDecision struct {
	Action       string                 `json:"action"`
	Target       string                 `json:"target"`
	Params       map[string]interface{} `json:"params,omitempty"`
	Reasoning    string                 `json:"reasoning"`
	Confidence   float64                `json:"confidence"`
	Alternative  string                 `json:"alternative,omitempty"` // Fallback plan
	ShouldNotify bool                   `json:"should_notify"`
}
