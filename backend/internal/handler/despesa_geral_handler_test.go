package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/felipehrs/financas-familiares/backend/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockDespesaGeralService implementa DespesaGeralServiceInterface (declarada no handler) para testes.
type MockDespesaGeralService struct {
	CriarFn        func(familiaID, membroID string, categoriaID *string, descricao string, data time.Time, valor float64, formaPagamento string, observacoes *string) (*domain.DespesaGeral, error)
	BuscarFn       func(familiaID, id string) (*domain.DespesaGeral, error)
	ListarFn       func(familiaID string) ([]*domain.DespesaGeral, error)
	ListarPorMesFn func(familiaID string, mes, ano int) ([]*domain.DespesaGeral, error)
	AtualizarFn    func(familiaID, id, membroID string, categoriaID *string, descricao string, data time.Time, valor float64, formaPagamento string, observacoes *string) (*domain.DespesaGeral, error)
	ExcluirFn      func(familiaID, id string) error
}

func (m *MockDespesaGeralService) Criar(familiaID, membroID string, categoriaID *string, descricao string, data time.Time, valor float64, formaPagamento string, observacoes *string) (*domain.DespesaGeral, error) {
	return m.CriarFn(familiaID, membroID, categoriaID, descricao, data, valor, formaPagamento, observacoes)
}

func (m *MockDespesaGeralService) BuscarPorID(familiaID, id string) (*domain.DespesaGeral, error) {
	return m.BuscarFn(familiaID, id)
}

func (m *MockDespesaGeralService) Listar(familiaID string) ([]*domain.DespesaGeral, error) {
	return m.ListarFn(familiaID)
}

func (m *MockDespesaGeralService) ListarPorMes(familiaID string, mes, ano int) ([]*domain.DespesaGeral, error) {
	return m.ListarPorMesFn(familiaID, mes, ano)
}

func (m *MockDespesaGeralService) Atualizar(familiaID, id, membroID string, categoriaID *string, descricao string, data time.Time, valor float64, formaPagamento string, observacoes *string) (*domain.DespesaGeral, error) {
	return m.AtualizarFn(familiaID, id, membroID, categoriaID, descricao, data, valor, formaPagamento, observacoes)
}

func (m *MockDespesaGeralService) Excluir(familiaID, id string) error {
	return m.ExcluirFn(familiaID, id)
}

func setupDespesaGeralRouter(svc handler.DespesaGeralServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", "usuario-teste-uuid")
		c.Set("familiaID", "familia-teste-uuid")
		c.Next()
	})
	h := handler.NewDespesaGeralHandler(svc)
	v1 := r.Group("/api/v1")
	{
		v1.GET("/despesas-gerais", h.Listar)
		v1.POST("/despesas-gerais", h.Criar)
		v1.GET("/despesas-gerais/:id", h.BuscarPorID)
		v1.PUT("/despesas-gerais/:id", h.Atualizar)
		v1.DELETE("/despesas-gerais/:id", h.Excluir)
	}
	return r
}

var despesaGeralExemplo = &domain.DespesaGeral{
	ID:             "uuid-1",
	MembroID:       "membro-1",
	Descricao:      "Mercado",
	Data:           time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC),
	Valor:          250.00,
	FormaPagamento: "pix",
}

// ---- POST /despesas-gerais ----

func TestCriarDespesaGeralHandler_Sucesso(t *testing.T) {
	svc := &MockDespesaGeralService{
		CriarFn: func(familiaID, membroID string, categoriaID *string, descricao string, data time.Time, valor float64, formaPagamento string, observacoes *string) (*domain.DespesaGeral, error) {
			return &domain.DespesaGeral{
				ID:             "uuid-novo",
				MembroID:       membroID,
				Descricao:      descricao,
				Data:           data,
				Valor:          valor,
				FormaPagamento: formaPagamento,
			}, nil
		},
	}
	r := setupDespesaGeralRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"membro_id":       "membro-1",
		"descricao":       "Mercado",
		"data":            "2026-03-08",
		"valor":           250.00,
		"forma_pagamento": "pix",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/despesas-gerais", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-novo", resp["id"])
	assert.Equal(t, "Mercado", resp["descricao"])
	assert.Equal(t, "2026-03-08", resp["data"])
}

func TestCriarDespesaGeralHandler_BodyInvalido(t *testing.T) {
	svc := &MockDespesaGeralService{}
	r := setupDespesaGeralRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/despesas-gerais", bytes.NewReader([]byte("nao-e-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCriarDespesaGeralHandler_ErroValidacao(t *testing.T) {
	svc := &MockDespesaGeralService{
		CriarFn: func(familiaID, membroID string, categoriaID *string, descricao string, data time.Time, valor float64, formaPagamento string, observacoes *string) (*domain.DespesaGeral, error) {
			return nil, domain.ErrDescricaoDespesaGeralObrigatoria
		},
	}
	r := setupDespesaGeralRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"membro_id":       "membro-1",
		"descricao":       "",
		"data":            "2026-03-08",
		"valor":           250.00,
		"forma_pagamento": "pix",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/despesas-gerais", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrDescricaoDespesaGeralObrigatoria.Error(), resp["error"])
}

// ---- GET /despesas-gerais ----

func TestListarDespesasGeraisHandler_Sucesso(t *testing.T) {
	despesas := []*domain.DespesaGeral{
		{ID: "uuid-1", MembroID: "membro-1", Descricao: "Mercado", Data: time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC), Valor: 250.00, FormaPagamento: "pix"},
		{ID: "uuid-2", MembroID: "membro-1", Descricao: "Farmácia", Data: time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC), Valor: 80.00, FormaPagamento: "dinheiro"},
	}
	svc := &MockDespesaGeralService{
		ListarFn: func(familiaID string) ([]*domain.DespesaGeral, error) {
			return despesas, nil
		},
	}
	r := setupDespesaGeralRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/despesas-gerais", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "uuid-1", resp[0]["id"])
}

func TestListarDespesasGeraisHandler_FiltradaPorMes(t *testing.T) {
	despesas := []*domain.DespesaGeral{
		{ID: "uuid-1", MembroID: "membro-1", Descricao: "Mercado", Data: time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC), Valor: 250.00, FormaPagamento: "pix"},
	}
	svc := &MockDespesaGeralService{
		ListarPorMesFn: func(familiaID string, mes, ano int) ([]*domain.DespesaGeral, error) {
			assert.Equal(t, 3, mes)
			assert.Equal(t, 2026, ano)
			return despesas, nil
		},
	}
	r := setupDespesaGeralRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/despesas-gerais?mes=3&ano=2026", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "uuid-1", resp[0]["id"])
}

// ---- GET /despesas-gerais/:id ----

func TestBuscarDespesaGeralHandler_Existente(t *testing.T) {
	svc := &MockDespesaGeralService{
		BuscarFn: func(familiaID, id string) (*domain.DespesaGeral, error) {
			return despesaGeralExemplo, nil
		},
	}
	r := setupDespesaGeralRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/despesas-gerais/uuid-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resp["id"])
	assert.Equal(t, "Mercado", resp["descricao"])
}

func TestBuscarDespesaGeralHandler_Inexistente(t *testing.T) {
	svc := &MockDespesaGeralService{
		BuscarFn: func(familiaID, id string) (*domain.DespesaGeral, error) {
			return nil, domain.ErrDespesaGeralNaoEncontrada
		},
	}
	r := setupDespesaGeralRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/despesas-gerais/uuid-inexistente", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrDespesaGeralNaoEncontrada.Error(), resp["error"])
}

// ---- PUT /despesas-gerais/:id ----

func TestAtualizarDespesaGeralHandler_Sucesso(t *testing.T) {
	svc := &MockDespesaGeralService{
		AtualizarFn: func(familiaID, id, membroID string, categoriaID *string, descricao string, data time.Time, valor float64, formaPagamento string, observacoes *string) (*domain.DespesaGeral, error) {
			return &domain.DespesaGeral{
				ID:             id,
				MembroID:       membroID,
				Descricao:      descricao,
				Data:           data,
				Valor:          valor,
				FormaPagamento: formaPagamento,
			}, nil
		},
	}
	r := setupDespesaGeralRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"membro_id":       "membro-1",
		"descricao":       "Supermercado",
		"data":            "2026-03-08",
		"valor":           300.00,
		"forma_pagamento": "debito",
	})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/despesas-gerais/uuid-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resp["id"])
	assert.Equal(t, "Supermercado", resp["descricao"])
}

// ---- DELETE /despesas-gerais/:id ----

func TestExcluirDespesaGeralHandler_Sucesso(t *testing.T) {
	svc := &MockDespesaGeralService{
		ExcluirFn: func(familiaID, id string) error {
			return nil
		},
	}
	r := setupDespesaGeralRouter(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/despesas-gerais/uuid-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestExcluirDespesaGeralHandler_Inexistente(t *testing.T) {
	svc := &MockDespesaGeralService{
		ExcluirFn: func(familiaID, id string) error {
			return domain.ErrDespesaGeralNaoEncontrada
		},
	}
	r := setupDespesaGeralRouter(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/despesas-gerais/uuid-inexistente", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrDespesaGeralNaoEncontrada.Error(), resp["error"])
}
