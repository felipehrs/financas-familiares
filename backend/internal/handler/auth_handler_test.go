package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/felipehrs/financas-familiares/backend/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockAuthService implementa AuthServiceInterface para testes do handler.
type MockAuthService struct {
	LoginFn               func(email, senha string) (string, string, error)
	RefreshTokenFn        func(refreshToken string) (string, error)
	ValidateAccessTokenFn func(token string) (string, string, error)
}

func (m *MockAuthService) Login(email, senha string) (string, string, error) {
	return m.LoginFn(email, senha)
}

func (m *MockAuthService) RefreshToken(refreshToken string) (string, error) {
	return m.RefreshTokenFn(refreshToken)
}

func (m *MockAuthService) ValidateAccessToken(token string) (string, string, error) {
	return m.ValidateAccessTokenFn(token)
}

func setupRouter(svc handler.AuthServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewAuthHandler(svc)
	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.Refresh)
	}
	return r
}

func TestLoginHandler_Sucesso(t *testing.T) {
	mockSvc := &MockAuthService{
		LoginFn: func(email, senha string) (string, string, error) {
			return "access-token-mock", "refresh-token-mock", nil
		},
	}
	r := setupRouter(mockSvc)

	body, _ := json.Marshal(map[string]string{
		"email": "user@example.com",
		"senha": "senha123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "access-token-mock", resp["access_token"])
	assert.Equal(t, "refresh-token-mock", resp["refresh_token"])
	assert.Equal(t, "Bearer", resp["token_type"])
	assert.Equal(t, float64(900), resp["expires_in"])
}

func TestLoginHandler_CredenciaisInvalidas(t *testing.T) {
	mockSvc := &MockAuthService{
		LoginFn: func(email, senha string) (string, string, error) {
			return "", "", domain.ErrCredenciaisInvalidas
		},
	}
	r := setupRouter(mockSvc)

	body, _ := json.Marshal(map[string]string{
		"email": "user@example.com",
		"senha": "senha-errada",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "credenciais inválidas", resp["error"])
}

func TestLoginHandler_BodyInvalido(t *testing.T) {
	mockSvc := &MockAuthService{}
	r := setupRouter(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader([]byte("nao-e-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRefreshHandler_Sucesso(t *testing.T) {
	mockSvc := &MockAuthService{
		RefreshTokenFn: func(refreshToken string) (string, error) {
			return "novo-access-token-mock", nil
		},
	}
	r := setupRouter(mockSvc)

	body, _ := json.Marshal(map[string]string{
		"refresh_token": "refresh-token-valido",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "novo-access-token-mock", resp["access_token"])
	assert.Equal(t, "Bearer", resp["token_type"])
	assert.Equal(t, float64(900), resp["expires_in"])
}

func TestRefreshHandler_Invalido(t *testing.T) {
	mockSvc := &MockAuthService{
		RefreshTokenFn: func(refreshToken string) (string, error) {
			return "", domain.ErrRefreshTokenInvalido
		},
	}
	r := setupRouter(mockSvc)

	body, _ := json.Marshal(map[string]string{
		"refresh_token": "token-invalido",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "refresh token inválido", resp["error"])
}

func TestRefreshHandler_BodyInvalido(t *testing.T) {
	mockSvc := &MockAuthService{}
	r := setupRouter(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader([]byte("nao-e-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
