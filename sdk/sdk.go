package sdk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Francapaa/InfrAgent/sdk/models"
)

//this is the folder which the users can install

type ActionFunc func(target string, params map[string]interface{}) error

type AgentSDK struct {
	apiKey        string
	webHookSecret string
	backendURL    string
	actions       (map[string]ActionFunc)
	healthCheck   string // url para verificar el /health
}

func NewSDK(apiKey, backendURL string, webHookSecret string) *AgentSDK {
	return &AgentSDK{
		apiKey:        apiKey,
		backendURL:    backendURL,
		webHookSecret: webHookSecret,
		actions:       make(map[string]ActionFunc),
	}
}

func (a *AgentSDK) On(action string, fn ActionFunc) {

	a.actions[action] = fn

}

// esta funcion es la que se encarga de verificar la firma para que no hayan problemas de seguridad.
func (a *AgentSDK) verifyHMAC(payload []byte, signature string) bool {
	h := hmac.New(sha256.New, []byte(a.webHookSecret))
	h.Write(payload)
	expected := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

func (a *AgentSDK) Run(port string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook/agent", a.handleWebhook)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	return server.ListenAndServe()
}

func (a *AgentSDK) handleWebhook(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "metodo NO permitido", http.StatusMethodNotAllowed)
		return
	}

	signature := r.Header.Get("X-InfrAgent-signature")
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}
	if !a.verifyHMAC(body, signature) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid signature"})
		return
	}

	// 3. Parsear el JSON
	var decision models.LLMDecision
	if err := json.Unmarshal(body, &decision); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid payload"})
		return
	}

	// 4. Ejecutar la acción
	handler, exists := a.actions[decision.Action]
	if !exists {
		w.WriteHeader(http.StatusNotImplemented)
		json.NewEncoder(w).Encode(map[string]string{"status": "aborted", "reason": "action not implemented"})
		return
	}

	if err := handler(decision.Target, decision.Params); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// 5. Éxito
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})

}

func (a *AgentSDK) checkLocalHealth() bool {

	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(a.healthCheck)

	if err != nil || resp.StatusCode != 200 {
		return false
	}
	return true
}

/*

SDK :=  agent.NewSDK(
"sajdasidasndsa",
"www.mibackend.com"
"daskdjsdasdasdásdp929hu5ejumkg,lbl j8iras"
)

NewSdk.Run() ESTO EJECUTA TODO

*/

/*
1) HOW THE SDK CAN EXECUTE THE ACTIONS
2) BEFORE TO EXECUTE ANY ACTION WE HAVE TO ASK TO THE OWNER OF THE PROJECT
 OUR AGENT EXECUTE SOMETHING WITH SYSMTECTL (VERY POWERFUL)
3) BAD PROMPT TO LLM, LLM SAID THAT WE HAVE TO RESTART THE API, BUT THE API IS WORKING WELL

IF CONFIDENCE < 70% THE AGENT CAN REQUEST TO THE /HEALTH ENDPOINT
IF /HEALTH === BAD REQUEST, SO WE CAN RESTART THE API
IF /HEALT === GOOD REQUEST, SO WE HAVE TO IMPROVE "SOMETHING" ABOUT THE LLM

*/
