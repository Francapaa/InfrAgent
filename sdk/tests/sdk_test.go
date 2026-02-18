package tests

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Francapaa/InfrAgent/sdk"
	"github.com/Francapaa/InfrAgent/sdk/models"
)

func TestHandleWebhook_MissingSignature(t *testing.T) {
	// Arrange
	sdk := sdk.NewSDK("test-api-key", "http://backend", "webhook-secret")

	payload := []byte(`{"action": "restart", "target": "api"}`)
	req := httptest.NewRequest(http.MethodPost, "/webhook/agent", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	// NO se establece el header X-InfrAgent-signature

	rr := httptest.NewRecorder()

	// Act
	sdk.HandleWebhook(rr, req)

	// Assert
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Se esperaba status 401, se obtuvo %d", rr.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("No se pudo parsear la respuesta JSON: %v", err)
	}

	if response["error"] != "invalid signature" {
		t.Errorf("Se esperaba error 'invalid signature', se obtuvo '%s'", response["error"])
	}
}

func TestHandleWebhook_TamperedSignature(t *testing.T) {
	// Arrange
	sdk := sdk.NewSDK("test-api-key", "http://backend", "webhook-secret")

	payload := []byte(`{"action": "restart", "target": "api"}`)
	req := httptest.NewRequest(http.MethodPost, "/webhook/agent", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	// Firma incorrecta/adulterada
	req.Header.Set("X-InfrAgent-signature", "firma-adulterada-incorrecta")

	rr := httptest.NewRecorder()

	// Act
	sdk.HandleWebhook(rr, req)

	// Assert
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Se esperaba status 401, se obtuvo %d", rr.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("No se pudo parsear la respuesta JSON: %v", err)
	}

	if response["error"] != "invalid signature" {
		t.Errorf("Se esperaba error 'invalid signature', se obtuvo '%s'", response["error"])
	}
}

func TestHandleWebhook_ValidSignature(t *testing.T) {
	// Arrange
	webhookSecret := "mi-secreto-super-seguro"
	sdk := sdk.NewSDK("test-api-key", "http://backend", webhookSecret)

	// Registrar una acción de prueba
	actionExecuted := false
	sdk.On("restart", func(target string, params map[string]interface{}) error {
		actionExecuted = true
		if target != "api" {
			t.Errorf("Se esperaba target 'api', se obtuvo '%s'", target)
		}
		return nil
	})

	decision := models.LLMDecision{
		Action: "restart",
		Target: "api",
		Params: map[string]interface{}{},
	}

	payload, _ := json.Marshal(decision)

	// Generar firma HMAC válida
	h := hmac.New(sha256.New, []byte(webhookSecret))
	h.Write(payload)
	validSignature := hex.EncodeToString(h.Sum(nil))

	req := httptest.NewRequest(http.MethodPost, "/webhook/agent", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-InfrAgent-signature", validSignature)

	rr := httptest.NewRecorder()

	// Act
	sdk.HandleWebhook(rr, req)

	// Assert
	if rr.Code != http.StatusOK {
		t.Errorf("Se esperaba status 200, se obtuvo %d", rr.Code)
	}

	if !actionExecuted {
		t.Error("La acción no fue ejecutada")
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("No se pudo parsear la respuesta JSON: %v", err)
	}

	if response["status"] != "success" {
		t.Errorf("Se esperaba status 'success', se obtuvo '%s'", response["status"])
	}
}

func TestHandleWebhook_WrongMethod(t *testing.T) {
	// Arrange
	sdk := sdk.NewSDK("test-api-key", "http://backend", "webhook-secret")

	req := httptest.NewRequest(http.MethodGet, "/webhook/agent", nil)
	rr := httptest.NewRecorder()

	// Act
	sdk.HandleWebhook(rr, req)

	// Assert
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Se esperaba status 405, se obtuvo %d", rr.Code)
	}
}

func TestHandleWebhook_EmptySignature(t *testing.T) {
	// Arrange
	sdk := sdk.NewSDK("test-api-key", "http://backend", "webhook-secret")

	payload := []byte(`{"action": "restart", "target": "api"}`)
	req := httptest.NewRequest(http.MethodPost, "/webhook/agent", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-InfrAgent-signature", "") // Firma vacía

	rr := httptest.NewRecorder()

	// Act
	sdk.HandleWebhook(rr, req)

	// Assert
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Se esperaba status 401, se obtuvo %d", rr.Code)
	}
}

func TestHandleWebhook_ModifiedPayload(t *testing.T) {
	// Arrange
	webhookSecret := "mi-secreto-super-seguro"
	sdk := sdk.NewSDK("test-api-key", "http://backend", webhookSecret)

	// Payload original
	originalPayload := []byte(`{"action": "restart", "target": "api"}`)

	// Generar firma HMAC válida para el payload original
	h := hmac.New(sha256.New, []byte(webhookSecret))
	h.Write(originalPayload)
	validSignature := hex.EncodeToString(h.Sum(nil))

	// Payload modificado (diferente al que se firmó)
	modifiedPayload := []byte(`{"action": "restart", "target": "database"}`)

	req := httptest.NewRequest(http.MethodPost, "/webhook/agent", bytes.NewBuffer(modifiedPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-InfrAgent-signature", validSignature) // Firma del payload original

	rr := httptest.NewRecorder()

	// Act
	sdk.HandleWebhook(rr, req)

	// Assert - debe rechazar porque el payload fue modificado
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Se esperaba status 401 para payload modificado, se obtuvo %d", rr.Code)
	}
}

func TestHandleWebhook_UnimplementedAction(t *testing.T) {
	// Arrange
	webhookSecret := "mi-secreto-super-seguro"
	sdk := sdk.NewSDK("test-api-key", "http://backend", webhookSecret)

	// NO registramos ninguna acción

	decision := models.LLMDecision{
		Action: "unknown-action",
		Target: "api",
		Params: map[string]interface{}{},
	}

	payload, _ := json.Marshal(decision)

	// Generar firma HMAC válida
	h := hmac.New(sha256.New, []byte(webhookSecret))
	h.Write(payload)
	validSignature := hex.EncodeToString(h.Sum(nil))

	req := httptest.NewRequest(http.MethodPost, "/webhook/agent", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-InfrAgent-signature", validSignature)

	rr := httptest.NewRecorder()

	// Act
	sdk.HandleWebhook(rr, req)

	// Assert
	if rr.Code != http.StatusNotImplemented {
		t.Errorf("Se esperaba status 501, se obtuvo %d", rr.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("No se pudo parsear la respuesta JSON: %v", err)
	}

	if response["reason"] != "action not implemented" {
		t.Errorf("Se esperaba reason 'action not implemented', se obtuvo '%s'", response["reason"])
	}
}

func TestHandleWebhook_InvalidJSON(t *testing.T) {
	// Arrange
	webhookSecret := "mi-secreto-super-seguro"
	sdk := sdk.NewSDK("test-api-key", "http://backend", webhookSecret)

	// Payload que no es JSON válido
	invalidPayload := []byte(`{action: restart, target: api}`)

	// Generar firma HMAC válida
	h := hmac.New(sha256.New, []byte(webhookSecret))
	h.Write(invalidPayload)
	validSignature := hex.EncodeToString(h.Sum(nil))

	req := httptest.NewRequest(http.MethodPost, "/webhook/agent", bytes.NewBuffer(invalidPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-InfrAgent-signature", validSignature)

	rr := httptest.NewRecorder()

	// Act
	sdk.HandleWebhook(rr, req)

	// Assert
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Se esperaba status 400, se obtuvo %d", rr.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("No se pudo parsear la respuesta JSON: %v", err)
	}

	if response["error"] != "invalid payload" {
		t.Errorf("Se esperaba error 'invalid payload', se obtuvo '%s'", response["error"])
	}
}
