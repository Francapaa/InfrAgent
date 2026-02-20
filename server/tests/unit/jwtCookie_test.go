package unit

import (
	"server/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJwtCookie(t *testing.T) {

	userId := uuid.New()
	token, err := utils.GenerateJWT(userId)
	assert.NoError(t, err)
	if err != nil {
		t.Fatalf("Error generando token: %v", err)
	}

	claims, err := utils.ValidateJWT(token)

	if err != nil && claims.UserID != "" {
		t.Errorf("Se esperaba que devuelva un struct y se obtuvo %v", err)
	}

}
