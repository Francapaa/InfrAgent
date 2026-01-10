package llm

import (
	"context"
	"log"

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

//esta funcion llama a la API de gemini (generando el prompt parseandola y validandolo)

/*
	COMO VALIDA? 
	=> SI LA DECISION QUE TOMA EL LLM ESTA DISPONIBLE DENTRO DE LAS DECISIONES QUE 
	   PUEDE TOMAR EL AGENTE (ESPERAR, REINICIAR BD O SERVIDOR, NOTIFICAR AL USER)

*/
func (g * GeminiClient) Decide (ctx context.Context, agentCtx storage.AgentContext)(*storage.LLMDecision, error){

	/*
	TO-DO: 
	1) PARSEAR EL PROMPT
	2) PASARSELO AL LLM
	3) PARSEAR RESPUESTA
	4) VER SI LA ACTION ESTA DENTRO DE LAS ACTIONS QUE PUEDE TOMAR EL AGENTE
	5) RETORNAR LA DECISION TOMADA, LUEGO EL AGENTE VE QUE ACCION PUEDE TOMAR. DONDE? &STORAGE.LLMDECISION.ACTION (LO DEVUELVE EL LLM)
	
	*/


}