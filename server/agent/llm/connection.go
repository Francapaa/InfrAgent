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

	if err != nil {
		log.Fatal(err)
	}

	return &GeminiClient{
		client: client,
		model:  model,
	}

}
