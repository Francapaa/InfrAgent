package unit

import (
	"server/utils"
	"testing"
)

/*
 TESTING DE APIKEYHASH
1) PETICION SIN API KEY => 401 UNAUTHORIZED
2) PETICION CON API KEY DIFERENTE => 401 UNAUTHORIZED
3) PETICION CON API KEY BIEN PERO HASH DIFERENTE => 401 UNAUTHORIZED
4) PETICION CON API KEY CORRECTA Y HASH CORRECTO => 200 OK (ACEPTA EVENTO)

*/

func TestIsValidApiKey(t *testing.T) {

	testScenaries := []struct {
		name     string
		inputKey string
		expected bool
	}{
		{"Key correcta", "agent_key_1235444123", true},
		{"Key incorrecta", "agent_bearer_21i39124", false},
		{"Key incorrecta formato", "agent_key 123564", false},
		{"Key vacia", "", false},
		{"Key sin prefijo", "12u3123", false},
	}

	for _, tt := range testScenaries {
		t.Run(tt.name, func(t *testing.T) {
			resultado := utils.IsValidAPIKeyMiddleware(tt.inputKey)

			if resultado != tt.expected {
				t.Errorf("Se esperaba %s y se obtuvo %v", tt.inputKey, resultado)
			}
		})
	}

}
