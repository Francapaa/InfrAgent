package unit

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"server/middleware"
	models "server/model"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClientStorage struct {
	mock.Mock
}

func (m *MockClientStorage) CreateClient(ctx context.Context, user *models.Client) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockClientStorage) GetClient(ctx context.Context, id uuid.UUID) (*models.Client, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Client), args.Error(1)
}

func (m *MockClientStorage) GetClientByAPIKey(ctx context.Context, apiKey string) (*models.Client, error) {
	args := m.Called(ctx, apiKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Client), args.Error(1)
}

func (m *MockClientStorage) GetClientByEmail(ctx context.Context, email string) (*models.Client, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Client), args.Error(1)
}

func (m *MockClientStorage) GetClientByGoogleID(ctx context.Context, googleID string) (*models.Client, error) {
	args := m.Called(ctx, googleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Client), args.Error(1)
}

func (m *MockClientStorage) UpdateClient(ctx context.Context, user *models.Client) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockClientStorage) UpdateClientComplete(ctx context.Context, user *models.Client) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockClientStorage) FixClientID(ctx context.Context, email string, newID string) error {
	args := m.Called(ctx, email, newID)
	return args.Error(0)
}

func TestProfileCompleteMiddleware_NoUserIDInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockStorage := new(MockClientStorage)
	router := gin.New()
	router.Use(middleware.ProfileCompleteMiddleware(mockStorage))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "unauthorized")
}

func TestProfileCompleteMiddleware_InvalidUserIDType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockStorage := new(MockClientStorage)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.UserIDKey, 12345)
		c.Next()
	})
	router.Use(middleware.ProfileCompleteMiddleware(mockStorage))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid user ID")
}

func TestProfileCompleteMiddleware_InvalidUUIDFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockStorage := new(MockClientStorage)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.UserIDKey, "invalid-uuid")
		c.Next()
	})
	router.Use(middleware.ProfileCompleteMiddleware(mockStorage))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid user ID format")
}

func TestProfileCompleteMiddleware_UserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	mockStorage := new(MockClientStorage)
	mockStorage.On("GetClient", mock.Anything, userID).Return(nil, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.UserIDKey, userID.String())
		c.Next()
	})
	router.Use(middleware.ProfileCompleteMiddleware(mockStorage))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "user not found")
	mockStorage.AssertExpectations(t)
}

func TestProfileCompleteMiddleware_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	mockStorage := new(MockClientStorage)
	mockStorage.On("GetClient", mock.Anything, userID).Return(nil, errors.New("database error"))

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.UserIDKey, userID.String())
		c.Next()
	})
	router.Use(middleware.ProfileCompleteMiddleware(mockStorage))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error retrieving user")
	mockStorage.AssertExpectations(t)
}

func TestProfileCompleteMiddleware_MissingWebhookURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	client := &models.Client{
		ID:          userID,
		Email:       "test@test.com",
		CompanyName: "Test Company",
		WebhookURL:  "",
	}
	mockStorage := new(MockClientStorage)
	mockStorage.On("GetClient", mock.Anything, userID).Return(client, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.UserIDKey, userID.String())
		c.Next()
	})
	router.Use(middleware.ProfileCompleteMiddleware(mockStorage))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "profile incomplete")
	assert.Contains(t, w.Body.String(), "webhook_url")
	mockStorage.AssertExpectations(t)
}

func TestProfileCompleteMiddleware_MissingCompanyName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	client := &models.Client{
		ID:          userID,
		Email:       "test@test.com",
		CompanyName: "",
		WebhookURL:  "https://webhook.test.com",
	}
	mockStorage := new(MockClientStorage)
	mockStorage.On("GetClient", mock.Anything, userID).Return(client, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.UserIDKey, userID.String())
		c.Next()
	})
	router.Use(middleware.ProfileCompleteMiddleware(mockStorage))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "profile incomplete")
	assert.Contains(t, w.Body.String(), "company_name")
	mockStorage.AssertExpectations(t)
}

func TestProfileCompleteMiddleware_BothFieldsMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	client := &models.Client{
		ID:          userID,
		Email:       "test@test.com",
		CompanyName: "",
		WebhookURL:  "",
	}
	mockStorage := new(MockClientStorage)
	mockStorage.On("GetClient", mock.Anything, userID).Return(client, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.UserIDKey, userID.String())
		c.Next()
	})
	router.Use(middleware.ProfileCompleteMiddleware(mockStorage))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "profile incomplete")
	assert.Contains(t, w.Body.String(), "webhook_url")
	assert.Contains(t, w.Body.String(), "company_name")
	mockStorage.AssertExpectations(t)
}

func TestProfileCompleteMiddleware_WebhookURLOnlySpaces(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	client := &models.Client{
		ID:          userID,
		Email:       "test@test.com",
		CompanyName: "Test Company",
		WebhookURL:  "     ",
	}
	mockStorage := new(MockClientStorage)
	mockStorage.On("GetClient", mock.Anything, userID).Return(client, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.UserIDKey, userID.String())
		c.Next()
	})
	router.Use(middleware.ProfileCompleteMiddleware(mockStorage))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "profile incomplete")
	assert.Contains(t, w.Body.String(), "webhook_url")
	mockStorage.AssertExpectations(t)
}

func TestProfileCompleteMiddleware_CompanyNameOnlySpaces(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	client := &models.Client{
		ID:          userID,
		Email:       "test@test.com",
		CompanyName: "    ",
		WebhookURL:  "https://webhook.test.com",
	}
	mockStorage := new(MockClientStorage)
	mockStorage.On("GetClient", mock.Anything, userID).Return(client, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.UserIDKey, userID.String())
		c.Next()
	})
	router.Use(middleware.ProfileCompleteMiddleware(mockStorage))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "profile incomplete")
	assert.Contains(t, w.Body.String(), "company_name")
	mockStorage.AssertExpectations(t)
}

func TestProfileCompleteMiddleware_WebhookURLWithoutHTTPS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	client := &models.Client{
		ID:          userID,
		Email:       "test@test.com",
		CompanyName: "Test Company",
		WebhookURL:  "http://webhook.test.com",
	}
	mockStorage := new(MockClientStorage)
	mockStorage.On("GetClient", mock.Anything, userID).Return(client, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.UserIDKey, userID.String())
		c.Next()
	})
	router.Use(middleware.ProfileCompleteMiddleware(mockStorage))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "profile incomplete")
	assert.Contains(t, w.Body.String(), "webhook_url_must_start_with_https")
	mockStorage.AssertExpectations(t)
}

func TestProfileCompleteMiddleware_ProfileComplete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	client := &models.Client{
		ID:          userID,
		Email:       "test@test.com",
		CompanyName: "Test Company",
		WebhookURL:  "https://webhook.test.com",
	}
	mockStorage := new(MockClientStorage)
	mockStorage.On("GetClient", mock.Anything, userID).Return(client, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.UserIDKey, userID.String())
		c.Next()
	})
	router.Use(middleware.ProfileCompleteMiddleware(mockStorage))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
	mockStorage.AssertExpectations(t)
}
