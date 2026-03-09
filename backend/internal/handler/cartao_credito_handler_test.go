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

// MockCartaoCreditoService implementa CartaoCreditoServiceInterface para testes.
type MockCartaoCreditoService struct {
	CriarFn     func(nome, membroID string, diaFechamento, diaVencimento int, limite *float64) (*domain.CartaoCredito, error)
	BuscarFn    func(id string) (*domain.CartaoCredito, error)
	ListarFn    func() ([]*domain.CartaoCredito, error)
	AtualizarFn func(id, nome, membroID string, diaFechamento, diaVencimento int, limite *float64, ativo bool) (*domain.CartaoCredito, error)
	InativarFn  func(id string) error
}

func (m *MockCartaoCreditoService) Criar(nome, membroID string, diaFechamento, diaVencimento int, limite *float64) (*domain.CartaoCredito, error) {
	return m.CriarFn(nome, membroID, diaFechamento, diaVencimento, limite)
}

func (m *MockCartaoCreditoService) BuscarPorID(id string) (*domain.CartaoCredito, error) {
	return m.BuscarFn(id)
}

func (m *MockCartaoCreditoService) Listar() ([]*domain.CartaoCredito, error) {
	return m.ListarFn()
}

func (m *MockCartaoCreditoService) Atualizar(id, nome, membroID string, diaFechamento, diaVencimento int, limite *float64, ativo bool) (*domain.CartaoCredito, error) {
	return m.AtualizarFn(id, nome, membroID, diaFechamento, diaVencimento, limite, ativo)
}

func (m *MockCartaoCreditoService) Inativar(id string) error {
	return m.InativarFn(id)
}

func setupCartaoRouter(svc handler.CartaoCreditoServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewCartaoCreditoHandler(svc)
	v1 := r.Group("/api/v1")
	{
		v1.GET("/cartoes", h.Listar)
		v1.POST("/cartoes", h.Criar)
		v1.GET("/cartoes/:id", h.BuscarPorID)
		v1.PUT("/cartoes/:id", h.Atualizar)
		v1.PATCH("/cartoes/:id/inativar", h.Inativar)
	}
	return r
}

// ---- GET /cartoes ----

func TestListarCartoesHandler_Sucesso(t *testing.T) {
	cartoes := []*domain.CartaoCredito{
		{ID: "uuid-1", Nome: "Nubank", MembroID: "m-1", DiaFechamento: 10, DiaVencimento: 17, Ativo: true},
		{ID: "uuid-2", Nome: "Inter", MembroID: "m-1", DiaFechamento: 5, DiaVencimento: 12, Ativo: true},
	}
	svc := &MockCartaoCreditoService{
		ListarFn: func() ([]*domain.CartaoCredito, error) { return cartoes, nil },
	}
	r := setupCartaoRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cartoes", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "uuid-1", resp[0]["id"])
	assert.Equal(t, "Nubank", resp[0]["nome"])
}

// ---- POST /cartoes ----

func TestCriarCartaoHandler_Sucesso(t *testing.T) {
	svc := &MockCartaoCreditoService{
		CriarFn: func(nome, membroID string, diaFechamento, diaVencimento int, limite *float64) (*domain.CartaoCredito, error) {
			return &domain.CartaoCredito{
				ID:            "uuid-novo",
				Nome:          nome,
				MembroID:      membroID,
				DiaFechamento: diaFechamento,
				DiaVencimento: diaVencimento,
				Limite:        limite,
				Ativo:         true,
			}, nil
		},
	}
	r := setupCartaoRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"nome":           "Nubank",
		"membro_id":      "membro-1",
		"dia_fechamento": 10,
		"dia_vencimento": 17,
		"limite":         5000.0,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cartoes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-novo", resp["id"])
	assert.Equal(t, "Nubank", resp["nome"])
	assert.Equal(t, true, resp["ativo"])
}

func TestCriarCartaoHandler_NomeVazio(t *testing.T) {
	svc := &MockCartaoCreditoService{
		CriarFn: func(nome, membroID string, diaFechamento, diaVencimento int, limite *float64) (*domain.CartaoCredito, error) {
			return nil, domain.ErrNomeCartaoObrigatorio
		},
	}
	r := setupCartaoRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"nome": "", "membro_id": "membro-1", "dia_fechamento": 10, "dia_vencimento": 17,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cartoes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrNomeCartaoObrigatorio.Error(), resp["error"])
}

func TestCriarCartaoHandler_BodyInvalido(t *testing.T) {
	svc := &MockCartaoCreditoService{}
	r := setupCartaoRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cartoes", bytes.NewReader([]byte("nao-e-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ---- GET /cartoes/:id ----

func TestBuscarCartaoHandler_Sucesso(t *testing.T) {
	svc := &MockCartaoCreditoService{
		BuscarFn: func(id string) (*domain.CartaoCredito, error) {
			return &domain.CartaoCredito{
				ID:            id,
				Nome:          "Nubank",
				MembroID:      "membro-1",
				DiaFechamento: 10,
				DiaVencimento: 17,
				Ativo:         true,
			}, nil
		},
	}
	r := setupCartaoRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cartoes/uuid-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resp["id"])
}

func TestBuscarCartaoHandler_NaoEncontrado(t *testing.T) {
	svc := &MockCartaoCreditoService{
		BuscarFn: func(id string) (*domain.CartaoCredito, error) {
			return nil, domain.ErrCartaoNaoEncontrado
		},
	}
	r := setupCartaoRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cartoes/uuid-inexistente", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

// ---- PUT /cartoes/:id ----

func TestAtualizarCartaoHandler_Sucesso(t *testing.T) {
	svc := &MockCartaoCreditoService{
		AtualizarFn: func(id, nome, membroID string, diaFechamento, diaVencimento int, limite *float64, ativo bool) (*domain.CartaoCredito, error) {
			return &domain.CartaoCredito{
				ID:            id,
				Nome:          nome,
				MembroID:      membroID,
				DiaFechamento: diaFechamento,
				DiaVencimento: diaVencimento,
				Limite:        limite,
				Ativo:         ativo,
			}, nil
		},
	}
	r := setupCartaoRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"nome": "Nubank Gold", "membro_id": "membro-1",
		"dia_fechamento": 15, "dia_vencimento": 20, "ativo": true,
	})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/cartoes/uuid-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resp["id"])
	assert.Equal(t, "Nubank Gold", resp["nome"])
}

func TestAtualizarCartaoHandler_NaoEncontrado(t *testing.T) {
	svc := &MockCartaoCreditoService{
		AtualizarFn: func(id, nome, membroID string, diaFechamento, diaVencimento int, limite *float64, ativo bool) (*domain.CartaoCredito, error) {
			return nil, domain.ErrCartaoNaoEncontrado
		},
	}
	r := setupCartaoRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"nome": "Nubank", "membro_id": "membro-1", "dia_fechamento": 10, "dia_vencimento": 17, "ativo": true,
	})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/cartoes/uuid-inexistente", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

// ---- PATCH /cartoes/:id/inativar ----

func TestInativarCartaoHandler_Sucesso(t *testing.T) {
	svc := &MockCartaoCreditoService{
		InativarFn: func(id string) error { return nil },
	}
	r := setupCartaoRouter(svc)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/cartoes/uuid-1/inativar", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "cartão inativado com sucesso", resp["message"])
}

func TestInativarCartaoHandler_NaoEncontrado(t *testing.T) {
	svc := &MockCartaoCreditoService{
		InativarFn: func(id string) error { return domain.ErrCartaoNaoEncontrado },
	}
	r := setupCartaoRouter(svc)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/cartoes/uuid-inexistente/inativar", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}
