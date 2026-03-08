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

// MockContaFixaService implementa ContaFixaServiceInterface (declarada no handler) para testes.
type MockContaFixaService struct {
	CriarFn       func(descricao, membroID string, categoriaID *string, valor float64, diaVencimento int, formaPagamento string) (*domain.ContaFixa, error)
	BuscarFn      func(id string) (*domain.ContaFixa, error)
	ListarFn      func() ([]*domain.ContaFixa, error)
	AtualizarFn   func(id, descricao, membroID string, categoriaID *string, valor float64, diaVencimento int, formaPagamento string, ativa bool) (*domain.ContaFixa, error)
	AlterarAtivoFn func(id string, ativa bool) (*domain.ContaFixa, error)
}

func (m *MockContaFixaService) Criar(descricao, membroID string, categoriaID *string, valor float64, diaVencimento int, formaPagamento string) (*domain.ContaFixa, error) {
	return m.CriarFn(descricao, membroID, categoriaID, valor, diaVencimento, formaPagamento)
}

func (m *MockContaFixaService) BuscarPorID(id string) (*domain.ContaFixa, error) {
	return m.BuscarFn(id)
}

func (m *MockContaFixaService) Listar() ([]*domain.ContaFixa, error) {
	return m.ListarFn()
}

func (m *MockContaFixaService) Atualizar(id, descricao, membroID string, categoriaID *string, valor float64, diaVencimento int, formaPagamento string, ativa bool) (*domain.ContaFixa, error) {
	return m.AtualizarFn(id, descricao, membroID, categoriaID, valor, diaVencimento, formaPagamento, ativa)
}

func (m *MockContaFixaService) AlterarAtivo(id string, ativa bool) (*domain.ContaFixa, error) {
	return m.AlterarAtivoFn(id, ativa)
}

func setupContaFixaRouter(svc handler.ContaFixaServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewContaFixaHandler(svc)
	v1 := r.Group("/api/v1")
	{
		v1.GET("/contas-fixas", h.Listar)
		v1.POST("/contas-fixas", h.Criar)
		v1.GET("/contas-fixas/:id", h.BuscarPorID)
		v1.PUT("/contas-fixas/:id", h.Atualizar)
		v1.PATCH("/contas-fixas/:id/ativo", h.AlterarAtivo)
	}
	return r
}

// ---- GET /contas-fixas ----

func TestListarContasFixasHandler_Sucesso(t *testing.T) {
	contas := []*domain.ContaFixa{
		{ID: "uuid-1", Descricao: "Internet", MembroID: "membro-1", Valor: 150.00, DiaVencimento: 10, FormaPagamento: "debito", Ativa: true},
		{ID: "uuid-2", Descricao: "Água", MembroID: "membro-1", Valor: 80.00, DiaVencimento: 15, FormaPagamento: "boleto", Ativa: true},
	}
	svc := &MockContaFixaService{
		ListarFn: func() ([]*domain.ContaFixa, error) {
			return contas, nil
		},
	}
	r := setupContaFixaRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/contas-fixas", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "uuid-1", resp[0]["id"])
	assert.Equal(t, "Internet", resp[0]["descricao"])
}

func TestListarContasFixasHandler_ListaVazia(t *testing.T) {
	svc := &MockContaFixaService{
		ListarFn: func() ([]*domain.ContaFixa, error) {
			return []*domain.ContaFixa{}, nil
		},
	}
	r := setupContaFixaRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/contas-fixas", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Empty(t, resp)
}

// ---- POST /contas-fixas ----

func TestCriarContaFixaHandler_Sucesso(t *testing.T) {
	svc := &MockContaFixaService{
		CriarFn: func(descricao, membroID string, categoriaID *string, valor float64, diaVencimento int, formaPagamento string) (*domain.ContaFixa, error) {
			return &domain.ContaFixa{
				ID:             "uuid-novo",
				Descricao:      descricao,
				MembroID:       membroID,
				CategoriaID:    categoriaID,
				Valor:          valor,
				DiaVencimento:  diaVencimento,
				FormaPagamento: formaPagamento,
				Ativa:          true,
			}, nil
		},
	}
	r := setupContaFixaRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"descricao":       "Internet",
		"membro_id":       "membro-1",
		"valor":           150.00,
		"dia_vencimento":  10,
		"forma_pagamento": "debito",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/contas-fixas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-novo", resp["id"])
	assert.Equal(t, "Internet", resp["descricao"])
	assert.Equal(t, true, resp["ativa"])
}

func TestCriarContaFixaHandler_DescricaoVazia(t *testing.T) {
	svc := &MockContaFixaService{
		CriarFn: func(descricao, membroID string, categoriaID *string, valor float64, diaVencimento int, formaPagamento string) (*domain.ContaFixa, error) {
			return nil, domain.ErrDescricaoContaFixaObrigatoria
		},
	}
	r := setupContaFixaRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"descricao":       "",
		"membro_id":       "membro-1",
		"valor":           150.00,
		"dia_vencimento":  10,
		"forma_pagamento": "debito",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/contas-fixas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrDescricaoContaFixaObrigatoria.Error(), resp["error"])
}

func TestCriarContaFixaHandler_BodyInvalido(t *testing.T) {
	svc := &MockContaFixaService{}
	r := setupContaFixaRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/contas-fixas", bytes.NewReader([]byte("nao-e-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ---- GET /contas-fixas/:id ----

func TestBuscarContaFixaHandler_Sucesso(t *testing.T) {
	svc := &MockContaFixaService{
		BuscarFn: func(id string) (*domain.ContaFixa, error) {
			return &domain.ContaFixa{
				ID:             id,
				Descricao:      "Internet",
				MembroID:       "membro-1",
				Valor:          150.00,
				DiaVencimento:  10,
				FormaPagamento: "debito",
				Ativa:          true,
			}, nil
		},
	}
	r := setupContaFixaRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/contas-fixas/uuid-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resp["id"])
	assert.Equal(t, "Internet", resp["descricao"])
}

func TestBuscarContaFixaHandler_NaoEncontrada(t *testing.T) {
	svc := &MockContaFixaService{
		BuscarFn: func(id string) (*domain.ContaFixa, error) {
			return nil, domain.ErrContaFixaNaoEncontrada
		},
	}
	r := setupContaFixaRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/contas-fixas/uuid-inexistente", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrContaFixaNaoEncontrada.Error(), resp["error"])
}

// ---- PUT /contas-fixas/:id ----

func TestAtualizarContaFixaHandler_Sucesso(t *testing.T) {
	svc := &MockContaFixaService{
		AtualizarFn: func(id, descricao, membroID string, categoriaID *string, valor float64, diaVencimento int, formaPagamento string, ativa bool) (*domain.ContaFixa, error) {
			return &domain.ContaFixa{
				ID:             id,
				Descricao:      descricao,
				MembroID:       membroID,
				CategoriaID:    categoriaID,
				Valor:          valor,
				DiaVencimento:  diaVencimento,
				FormaPagamento: formaPagamento,
				Ativa:          ativa,
			}, nil
		},
	}
	r := setupContaFixaRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"descricao":       "Internet Fibra",
		"membro_id":       "membro-1",
		"valor":           200.00,
		"dia_vencimento":  10,
		"forma_pagamento": "debito",
		"ativa":           true,
	})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/contas-fixas/uuid-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resp["id"])
	assert.Equal(t, "Internet Fibra", resp["descricao"])
}

func TestAtualizarContaFixaHandler_NaoEncontrada(t *testing.T) {
	svc := &MockContaFixaService{
		AtualizarFn: func(id, descricao, membroID string, categoriaID *string, valor float64, diaVencimento int, formaPagamento string, ativa bool) (*domain.ContaFixa, error) {
			return nil, domain.ErrContaFixaNaoEncontrada
		},
	}
	r := setupContaFixaRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"descricao":       "Internet",
		"membro_id":       "membro-1",
		"valor":           150.00,
		"dia_vencimento":  10,
		"forma_pagamento": "debito",
		"ativa":           true,
	})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/contas-fixas/uuid-inexistente", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrContaFixaNaoEncontrada.Error(), resp["error"])
}

func TestAtualizarContaFixaHandler_BodyInvalido(t *testing.T) {
	svc := &MockContaFixaService{}
	r := setupContaFixaRouter(svc)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/contas-fixas/uuid-1", bytes.NewReader([]byte("nao-e-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ---- PATCH /contas-fixas/:id/ativo ----

func TestAlterarAtivoHandler_Sucesso(t *testing.T) {
	svc := &MockContaFixaService{
		AlterarAtivoFn: func(id string, ativa bool) (*domain.ContaFixa, error) {
			return &domain.ContaFixa{
				ID:             id,
				Descricao:      "Internet",
				MembroID:       "membro-1",
				Valor:          150.00,
				DiaVencimento:  10,
				FormaPagamento: "debito",
				Ativa:          ativa,
			}, nil
		},
	}
	r := setupContaFixaRouter(svc)

	body, _ := json.Marshal(map[string]bool{
		"ativa": false,
	})
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/contas-fixas/uuid-1/ativo", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resp["id"])
	assert.Equal(t, false, resp["ativa"])
}

func TestAlterarAtivoHandler_NaoEncontrada(t *testing.T) {
	svc := &MockContaFixaService{
		AlterarAtivoFn: func(id string, ativa bool) (*domain.ContaFixa, error) {
			return nil, domain.ErrContaFixaNaoEncontrada
		},
	}
	r := setupContaFixaRouter(svc)

	body, _ := json.Marshal(map[string]bool{
		"ativa": false,
	})
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/contas-fixas/uuid-inexistente/ativo", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrContaFixaNaoEncontrada.Error(), resp["error"])
}

func TestAlterarAtivoHandler_BodyInvalido(t *testing.T) {
	svc := &MockContaFixaService{}
	r := setupContaFixaRouter(svc)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/contas-fixas/uuid-1/ativo", bytes.NewReader([]byte("nao-e-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
