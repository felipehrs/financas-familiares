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

// MockRendimentoInvestimentoService implementa RendimentoInvestimentoServiceInterface (declarada no handler) para testes.
type MockRendimentoInvestimentoService struct {
	CriarFn        func(descricao, membroID string, data time.Time, valor, valorDistribuido float64) (*domain.RendimentoInvestimento, error)
	BuscarFn       func(id string) (*domain.RendimentoInvestimento, error)
	ListarFn       func() ([]*domain.RendimentoInvestimento, error)
	ListarPorMesFn func(mes, ano int) ([]*domain.RendimentoInvestimento, error)
	AtualizarFn    func(id, descricao, membroID string, data time.Time, valor, valorDistribuido float64) (*domain.RendimentoInvestimento, error)
	ExcluirFn      func(id string) error
}

func (m *MockRendimentoInvestimentoService) Criar(descricao, membroID string, data time.Time, valor, valorDistribuido float64) (*domain.RendimentoInvestimento, error) {
	return m.CriarFn(descricao, membroID, data, valor, valorDistribuido)
}

func (m *MockRendimentoInvestimentoService) BuscarPorID(id string) (*domain.RendimentoInvestimento, error) {
	return m.BuscarFn(id)
}

func (m *MockRendimentoInvestimentoService) Listar() ([]*domain.RendimentoInvestimento, error) {
	return m.ListarFn()
}

func (m *MockRendimentoInvestimentoService) ListarPorMes(mes, ano int) ([]*domain.RendimentoInvestimento, error) {
	return m.ListarPorMesFn(mes, ano)
}

func (m *MockRendimentoInvestimentoService) Atualizar(id, descricao, membroID string, data time.Time, valor, valorDistribuido float64) (*domain.RendimentoInvestimento, error) {
	return m.AtualizarFn(id, descricao, membroID, data, valor, valorDistribuido)
}

func (m *MockRendimentoInvestimentoService) Excluir(id string) error {
	return m.ExcluirFn(id)
}

func setupRendimentoInvestimentoRouter(svc handler.RendimentoInvestimentoServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewRendimentoInvestimentoHandler(svc)
	v1 := r.Group("/api/v1")
	{
		v1.GET("/rendimentos-investimento", h.Listar)
		v1.POST("/rendimentos-investimento", h.Criar)
		v1.GET("/rendimentos-investimento/:id", h.BuscarPorID)
		v1.PUT("/rendimentos-investimento/:id", h.Atualizar)
		v1.DELETE("/rendimentos-investimento/:id", h.Excluir)
	}
	return r
}

var rendimentoInvestimentoExemplo = &domain.RendimentoInvestimento{
	ID:               "uuid-1",
	Descricao:        "Rendimento CDB",
	MembroID:         "membro-1",
	Data:             time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC),
	Valor:            1000.00,
	ValorDistribuido: 0,
}

// ---- POST /rendimentos-investimento ----

func TestCriarRendimentoInvestimentoHandler_Sucesso(t *testing.T) {
	svc := &MockRendimentoInvestimentoService{
		CriarFn: func(descricao, membroID string, data time.Time, valor, valorDistribuido float64) (*domain.RendimentoInvestimento, error) {
			return &domain.RendimentoInvestimento{
				ID:               "uuid-novo",
				Descricao:        descricao,
				MembroID:         membroID,
				Data:             data,
				Valor:            valor,
				ValorDistribuido: valorDistribuido,
			}, nil
		},
	}
	r := setupRendimentoInvestimentoRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"descricao": "Rendimento CDB",
		"membro_id": "membro-1",
		"data":      "2026-03-08",
		"valor":     1000.00,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/rendimentos-investimento", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-novo", resp["id"])
	assert.Equal(t, "Rendimento CDB", resp["descricao"])
	assert.Equal(t, "2026-03-08", resp["data"])
}

func TestCriarRendimentoInvestimentoHandler_BodyInvalido(t *testing.T) {
	svc := &MockRendimentoInvestimentoService{}
	r := setupRendimentoInvestimentoRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rendimentos-investimento", bytes.NewReader([]byte("nao-e-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCriarRendimentoInvestimentoHandler_ErroValidacao(t *testing.T) {
	svc := &MockRendimentoInvestimentoService{
		CriarFn: func(descricao, membroID string, data time.Time, valor, valorDistribuido float64) (*domain.RendimentoInvestimento, error) {
			return nil, domain.ErrDescricaoRendimentoObrigatoria
		},
	}
	r := setupRendimentoInvestimentoRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"descricao": "",
		"membro_id": "membro-1",
		"data":      "2026-03-08",
		"valor":     1000.00,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/rendimentos-investimento", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrDescricaoRendimentoObrigatoria.Error(), resp["error"])
}

// ---- GET /rendimentos-investimento ----

func TestListarRendimentosInvestimentoHandler_Sucesso(t *testing.T) {
	rendimentos := []*domain.RendimentoInvestimento{
		{ID: "uuid-1", Descricao: "Rendimento CDB", MembroID: "membro-1", Data: time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC), Valor: 1000.00},
		{ID: "uuid-2", Descricao: "Rendimento LCI", MembroID: "membro-1", Data: time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC), Valor: 2000.00},
	}
	svc := &MockRendimentoInvestimentoService{
		ListarFn: func() ([]*domain.RendimentoInvestimento, error) {
			return rendimentos, nil
		},
	}
	r := setupRendimentoInvestimentoRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rendimentos-investimento", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "uuid-1", resp[0]["id"])
}

func TestListarRendimentosInvestimentoHandler_FiltradaPorMes(t *testing.T) {
	rendimentos := []*domain.RendimentoInvestimento{
		{ID: "uuid-1", Descricao: "Rendimento CDB", MembroID: "membro-1", Data: time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC), Valor: 1000.00},
	}
	svc := &MockRendimentoInvestimentoService{
		ListarPorMesFn: func(mes, ano int) ([]*domain.RendimentoInvestimento, error) {
			assert.Equal(t, 3, mes)
			assert.Equal(t, 2026, ano)
			return rendimentos, nil
		},
	}
	r := setupRendimentoInvestimentoRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rendimentos-investimento?mes=3&ano=2026", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "uuid-1", resp[0]["id"])
}

// ---- GET /rendimentos-investimento/:id ----

func TestBuscarRendimentoInvestimentoHandler_Existente(t *testing.T) {
	svc := &MockRendimentoInvestimentoService{
		BuscarFn: func(id string) (*domain.RendimentoInvestimento, error) {
			return rendimentoInvestimentoExemplo, nil
		},
	}
	r := setupRendimentoInvestimentoRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rendimentos-investimento/uuid-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resp["id"])
	assert.Equal(t, "Rendimento CDB", resp["descricao"])
}

func TestBuscarRendimentoInvestimentoHandler_Inexistente(t *testing.T) {
	svc := &MockRendimentoInvestimentoService{
		BuscarFn: func(id string) (*domain.RendimentoInvestimento, error) {
			return nil, domain.ErrRendimentoInvestimentoNaoEncontrado
		},
	}
	r := setupRendimentoInvestimentoRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rendimentos-investimento/uuid-inexistente", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrRendimentoInvestimentoNaoEncontrado.Error(), resp["error"])
}

// ---- PUT /rendimentos-investimento/:id ----

func TestAtualizarRendimentoInvestimentoHandler_Sucesso(t *testing.T) {
	svc := &MockRendimentoInvestimentoService{
		AtualizarFn: func(id, descricao, membroID string, data time.Time, valor, valorDistribuido float64) (*domain.RendimentoInvestimento, error) {
			return &domain.RendimentoInvestimento{
				ID:               id,
				Descricao:        descricao,
				MembroID:         membroID,
				Data:             data,
				Valor:            valor,
				ValorDistribuido: valorDistribuido,
			}, nil
		},
	}
	r := setupRendimentoInvestimentoRouter(svc)

	valorDistribuido := 300.0
	body, _ := json.Marshal(map[string]any{
		"descricao":         "Rendimento CDB atualizado",
		"membro_id":         "membro-1",
		"data":              "2026-03-10",
		"valor":             1500.00,
		"valor_distribuido": valorDistribuido,
	})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/rendimentos-investimento/uuid-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resp["id"])
	assert.Equal(t, "Rendimento CDB atualizado", resp["descricao"])
}

// ---- DELETE /rendimentos-investimento/:id ----

func TestExcluirRendimentoInvestimentoHandler_Sucesso(t *testing.T) {
	svc := &MockRendimentoInvestimentoService{
		ExcluirFn: func(id string) error {
			return nil
		},
	}
	r := setupRendimentoInvestimentoRouter(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/rendimentos-investimento/uuid-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestExcluirRendimentoInvestimentoHandler_Inexistente(t *testing.T) {
	svc := &MockRendimentoInvestimentoService{
		ExcluirFn: func(id string) error {
			return domain.ErrRendimentoInvestimentoNaoEncontrado
		},
	}
	r := setupRendimentoInvestimentoRouter(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/rendimentos-investimento/uuid-inexistente", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrRendimentoInvestimentoNaoEncontrado.Error(), resp["error"])
}
