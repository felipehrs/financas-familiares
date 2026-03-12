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

// MockAssinaturaService implementa AssinaturaServiceInterface (declarada no handler) para testes.
type MockAssinaturaService struct {
	CriarFn         func(familiaID, nome, membroID string, categoriaID *string, valor float64, diaCobranca int, formaPagamento string) (*domain.Assinatura, error)
	BuscarFn        func(familiaID, id string) (*domain.Assinatura, error)
	ListarFn        func(familiaID string) ([]*domain.Assinatura, error)
	AtualizarFn     func(familiaID, id, nome, membroID string, categoriaID *string, valor float64, diaCobranca int, formaPagamento, status string) (*domain.Assinatura, error)
	AlterarStatusFn func(familiaID, id, status string) (*domain.Assinatura, error)
}

func (m *MockAssinaturaService) Criar(familiaID, nome, membroID string, categoriaID *string, valor float64, diaCobranca int, formaPagamento string) (*domain.Assinatura, error) {
	return m.CriarFn(familiaID, nome, membroID, categoriaID, valor, diaCobranca, formaPagamento)
}

func (m *MockAssinaturaService) BuscarPorID(familiaID, id string) (*domain.Assinatura, error) {
	return m.BuscarFn(familiaID, id)
}

func (m *MockAssinaturaService) Listar(familiaID string) ([]*domain.Assinatura, error) {
	return m.ListarFn(familiaID)
}

func (m *MockAssinaturaService) Atualizar(familiaID, id, nome, membroID string, categoriaID *string, valor float64, diaCobranca int, formaPagamento, status string) (*domain.Assinatura, error) {
	return m.AtualizarFn(familiaID, id, nome, membroID, categoriaID, valor, diaCobranca, formaPagamento, status)
}

func (m *MockAssinaturaService) AlterarStatus(familiaID, id, status string) (*domain.Assinatura, error) {
	return m.AlterarStatusFn(familiaID, id, status)
}

func setupAssinaturaRouter(svc handler.AssinaturaServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", "usuario-teste-uuid")
		c.Set("familiaID", "familia-teste-uuid")
		c.Next()
	})
	h := handler.NewAssinaturaHandler(svc)
	v1 := r.Group("/api/v1")
	{
		v1.GET("/assinaturas", h.Listar)
		v1.POST("/assinaturas", h.Criar)
		v1.GET("/assinaturas/:id", h.BuscarPorID)
		v1.PUT("/assinaturas/:id", h.Atualizar)
		v1.PATCH("/assinaturas/:id/status", h.AlterarStatus)
	}
	return r
}

// ---- GET /assinaturas ----

func TestListarAssinaturasHandler_Sucesso(t *testing.T) {
	assinaturas := []*domain.Assinatura{
		{ID: "uuid-1", Nome: "Netflix", MembroID: "membro-1", Valor: 39.90, DiaCobranca: 15, FormaPagamento: "cartao_credito", Status: "ativa"},
		{ID: "uuid-2", Nome: "Spotify", MembroID: "membro-1", Valor: 19.90, DiaCobranca: 10, FormaPagamento: "debito", Status: "ativa"},
	}
	svc := &MockAssinaturaService{
		ListarFn: func(familiaID string) ([]*domain.Assinatura, error) {
			return assinaturas, nil
		},
	}
	r := setupAssinaturaRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/assinaturas", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "uuid-1", resp[0]["id"])
	assert.Equal(t, "Netflix", resp[0]["nome"])
}

func TestListarAssinaturasHandler_ListaVazia(t *testing.T) {
	svc := &MockAssinaturaService{
		ListarFn: func(familiaID string) ([]*domain.Assinatura, error) {
			return []*domain.Assinatura{}, nil
		},
	}
	r := setupAssinaturaRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/assinaturas", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Empty(t, resp)
}

// ---- POST /assinaturas ----

func TestCriarAssinaturaHandler_Sucesso(t *testing.T) {
	svc := &MockAssinaturaService{
		CriarFn: func(familiaID, nome, membroID string, categoriaID *string, valor float64, diaCobranca int, formaPagamento string) (*domain.Assinatura, error) {
			return &domain.Assinatura{
				ID:             "uuid-novo",
				Nome:           nome,
				MembroID:       membroID,
				CategoriaID:    categoriaID,
				Valor:          valor,
				DiaCobranca:    diaCobranca,
				FormaPagamento: formaPagamento,
				Status:         "ativa",
			}, nil
		},
	}
	r := setupAssinaturaRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"nome":            "Netflix",
		"membro_id":       "membro-1",
		"valor":           39.90,
		"dia_cobranca":    15,
		"forma_pagamento": "cartao_credito",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/assinaturas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-novo", resp["id"])
	assert.Equal(t, "Netflix", resp["nome"])
	assert.Equal(t, "ativa", resp["status"])
}

func TestCriarAssinaturaHandler_NomeVazio(t *testing.T) {
	svc := &MockAssinaturaService{
		CriarFn: func(familiaID, nome, membroID string, categoriaID *string, valor float64, diaCobranca int, formaPagamento string) (*domain.Assinatura, error) {
			return nil, domain.ErrNomeAssinaturaObrigatorio
		},
	}
	r := setupAssinaturaRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"nome":            "",
		"membro_id":       "membro-1",
		"valor":           39.90,
		"dia_cobranca":    15,
		"forma_pagamento": "cartao_credito",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/assinaturas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrNomeAssinaturaObrigatorio.Error(), resp["error"])
}

func TestCriarAssinaturaHandler_BodyInvalido(t *testing.T) {
	svc := &MockAssinaturaService{}
	r := setupAssinaturaRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/assinaturas", bytes.NewReader([]byte("nao-e-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ---- GET /assinaturas/:id ----

func TestBuscarAssinaturaHandler_Sucesso(t *testing.T) {
	svc := &MockAssinaturaService{
		BuscarFn: func(familiaID, id string) (*domain.Assinatura, error) {
			return &domain.Assinatura{
				ID:             id,
				Nome:           "Netflix",
				MembroID:       "membro-1",
				Valor:          39.90,
				DiaCobranca:    15,
				FormaPagamento: "cartao_credito",
				Status:         "ativa",
			}, nil
		},
	}
	r := setupAssinaturaRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/assinaturas/uuid-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resp["id"])
	assert.Equal(t, "Netflix", resp["nome"])
}

func TestBuscarAssinaturaHandler_NaoEncontrada(t *testing.T) {
	svc := &MockAssinaturaService{
		BuscarFn: func(familiaID, id string) (*domain.Assinatura, error) {
			return nil, domain.ErrAssinaturaNaoEncontrada
		},
	}
	r := setupAssinaturaRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/assinaturas/uuid-inexistente", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrAssinaturaNaoEncontrada.Error(), resp["error"])
}

// ---- PUT /assinaturas/:id ----

func TestAtualizarAssinaturaHandler_Sucesso(t *testing.T) {
	svc := &MockAssinaturaService{
		AtualizarFn: func(familiaID, id, nome, membroID string, categoriaID *string, valor float64, diaCobranca int, formaPagamento, status string) (*domain.Assinatura, error) {
			return &domain.Assinatura{
				ID:             id,
				Nome:           nome,
				MembroID:       membroID,
				CategoriaID:    categoriaID,
				Valor:          valor,
				DiaCobranca:    diaCobranca,
				FormaPagamento: formaPagamento,
				Status:         status,
			}, nil
		},
	}
	r := setupAssinaturaRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"nome":            "Netflix Premium",
		"membro_id":       "membro-1",
		"valor":           55.90,
		"dia_cobranca":    15,
		"forma_pagamento": "cartao_credito",
		"status":          "ativa",
	})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/assinaturas/uuid-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resp["id"])
	assert.Equal(t, "Netflix Premium", resp["nome"])
}

func TestAtualizarAssinaturaHandler_NaoEncontrada(t *testing.T) {
	svc := &MockAssinaturaService{
		AtualizarFn: func(familiaID, id, nome, membroID string, categoriaID *string, valor float64, diaCobranca int, formaPagamento, status string) (*domain.Assinatura, error) {
			return nil, domain.ErrAssinaturaNaoEncontrada
		},
	}
	r := setupAssinaturaRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"nome":            "Netflix",
		"membro_id":       "membro-1",
		"valor":           39.90,
		"dia_cobranca":    15,
		"forma_pagamento": "cartao_credito",
		"status":          "ativa",
	})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/assinaturas/uuid-inexistente", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrAssinaturaNaoEncontrada.Error(), resp["error"])
}

func TestAtualizarAssinaturaHandler_BodyInvalido(t *testing.T) {
	svc := &MockAssinaturaService{}
	r := setupAssinaturaRouter(svc)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/assinaturas/uuid-1", bytes.NewReader([]byte("nao-e-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ---- PATCH /assinaturas/:id/status ----

func TestAlterarStatusHandler_Sucesso(t *testing.T) {
	svc := &MockAssinaturaService{
		AlterarStatusFn: func(familiaID, id, status string) (*domain.Assinatura, error) {
			return &domain.Assinatura{
				ID:             id,
				Nome:           "Netflix",
				MembroID:       "membro-1",
				Valor:          39.90,
				DiaCobranca:    15,
				FormaPagamento: "cartao_credito",
				Status:         status,
			}, nil
		},
	}
	r := setupAssinaturaRouter(svc)

	body, _ := json.Marshal(map[string]string{
		"status": "pausada",
	})
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/assinaturas/uuid-1/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resp["id"])
	assert.Equal(t, "pausada", resp["status"])
}

func TestAlterarStatusHandler_StatusInvalido(t *testing.T) {
	svc := &MockAssinaturaService{
		AlterarStatusFn: func(familiaID, id, status string) (*domain.Assinatura, error) {
			return nil, domain.ErrStatusInvalido
		},
	}
	r := setupAssinaturaRouter(svc)

	body, _ := json.Marshal(map[string]string{
		"status": "invalido",
	})
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/assinaturas/uuid-1/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrStatusInvalido.Error(), resp["error"])
}

func TestAlterarStatusHandler_NaoEncontrada(t *testing.T) {
	svc := &MockAssinaturaService{
		AlterarStatusFn: func(familiaID, id, status string) (*domain.Assinatura, error) {
			return nil, domain.ErrAssinaturaNaoEncontrada
		},
	}
	r := setupAssinaturaRouter(svc)

	body, _ := json.Marshal(map[string]string{
		"status": "cancelada",
	})
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/assinaturas/uuid-inexistente/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrAssinaturaNaoEncontrada.Error(), resp["error"])
}
