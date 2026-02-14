package sdk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sdk/models"
	"time"

	"github.com/gin-gonic/gin"
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

func (a *AgentSDK) Run(port string) error {
	r := gin.Default()

	r.POST("/webhook/agent", a.handleWebhook)
	fmt.Printf("[SDK] escuchando en el puerto %s \n", port)
	return r.Run(":" + port)
}

func (a *AgentSDK) handleWebhook(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot read body"})
		return
	}

	signature := c.GetHeader("X-Agent-Signature")
	if signature == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing signature"})
		return
	}

	if !verifySignature(body, signature, a.webHookSecret) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
		return
	}

	var decision models.LLMDecision
	if err := json.Unmarshal(body, &decision); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	if decision.Confidence < 0.9 || decision.Action == "restart" {
		fmt.Printf("[AGENTE] verificando si de verdad la %s esta caido", decision.Target)

		isHealthy := a.checkLocalHealth()

		if isHealthy {
			c.JSON(http.StatusOK, gin.H{"status": "aborted", "reason": "local health passed, works normally"})
			return
		}
	}

	handler, exists := a.actions[decision.Action]

	if !exists { // una accion que no existe
		c.JSON(http.StatusNotImplemented, gin.H{"status": "aborted", "reason": "this action doesnt exists"})
		return
	}

	if err := handler(decision.Target, decision.Params); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (a *AgentSDK) checkLocalHealth() bool {

	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(a.healthCheck)

	if err != nil || resp.StatusCode != 200 {
		return false
	}
	return true
}

func verifySignature(body []byte, signature, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
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
