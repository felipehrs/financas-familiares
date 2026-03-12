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

// MockRendaExtraService implementa RendaExtraServiceInterface (declarada no handler) para testes.
type MockRendaExtraService struct {
	CriarFn        func(familiaID, descricao, membroID string, dataRecebimento time.Time, valor float64) (*domain.RendaExtra, error)
	BuscarFn       func(familiaID, id string) (*domain.RendaExtra, error)
	ListarFn       func(familiaID string) ([]*domain.RendaExtra, error)
	ListarPorMesFn func(familiaID string, mes, ano int) ([]*domain.RendaExtra, error)
	AtualizarFn    func(familiaID, id, descricao, membroID string, dataRecebimento time.Time, valor float64) (*domain.RendaExtra, error)
	ExcluirFn      func(familiaID, id string) error
}

func (m *MockRendaExtraService) Criar(familiaID, descricao, membroID string, dataRecebimento time.Time, valor float64) (*domain.RendaExtra, error) {
	return m.CriarFn(familiaID, descricao, membroID, dataRecebimento, valor)
}

func (m *MockRendaExtraService) BuscarPorID(familiaID, id string) (*domain.RendaExtra, error) {
	return m.BuscarFn(familiaID, id)
}

func (m *MockRendaExtraService) Listar(familiaID string) ([]*domain.RendaExtra, error) {
	return m.ListarFn(familiaID)
}

func (m *MockRendaExtraService) ListarPorMes(familiaID string, mes, ano int) ([]*domain.RendaExtra, error) {
	return m.ListarPorMesFn(familiaID, mes, ano)
}

func (m *MockRendaExtraService) Atualizar(familiaID, id, descricao, membroID string, dataRecebimento time.Time, valor float64) (*domain.RendaExtra, error) {
	return m.AtualizarFn(familiaID, id, descricao, membroID, dataRecebimento, valor)
}

func (m *MockRendaExtraService) Excluir(familiaID, id string) error {
	return m.ExcluirFn(familiaID, id)
}

func setupRendaExtraRouter(svc handler.RendaExtraServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", "usuario-teste-uuid")
		c.Set("familiaID", "familia-teste-uuid")
		c.Next()
	})
	h := handler.NewRendaExtraHandler(svc)
	v1 := r.Group("/api/v1")
	{
		v1.GET("/rendas-extras", h.Listar)
		v1.POST("/rendas-extras", h.Criar)
		v1.GET("/rendas-extras/:id", h.BuscarPorID)
		v1.PUT("/rendas-extras/:id", h.Atualizar)
		v1.DELETE("/rendas-extras/:id", h.Excluir)
	}
	return r
}

var rendaExtraExemplo = &domain.RendaExtra{
	ID:              "uuid-1",
	Descricao:       "Bônus março",
	MembroID:        "membro-1",
	DataRecebimento: time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC),
	Valor:           2000.00,
}

// ---- POST /rendas-extras ----

func TestCriarRendaExtraHandler_Sucesso(t *testing.T) {
	svc := &MockRendaExtraService{
		CriarFn: func(familiaID, descricao, membroID string, dataRecebimento time.Time, valor float64) (*domain.RendaExtra, error) {
			return &domain.RendaExtra{
				ID:              "uuid-novo",
				Descricao:       descricao,
				MembroID:        membroID,
				DataRecebimento: dataRecebimento,
				Valor:           valor,
			}, nil
		},
	}
	r := setupRendaExtraRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"descricao":        "Bônus março",
		"membro_id":        "membro-1",
		"data_recebimento": "2026-03-08",
		"valor":            2000.00,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/rendas-extras", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-novo", resp["id"])
	assert.Equal(t, "Bônus março", resp["descricao"])
	assert.Equal(t, "2026-03-08", resp["data_recebimento"])
}

func TestCriarRendaExtraHandler_BodyInvalido(t *testing.T) {
	svc := &MockRendaExtraService{}
	r := setupRendaExtraRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rendas-extras", bytes.NewReader([]byte("nao-e-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCriarRendaExtraHandler_ErroValidacao(t *testing.T) {
	svc := &MockRendaExtraService{
		CriarFn: func(familiaID, descricao, membroID string, dataRecebimento time.Time, valor float64) (*domain.RendaExtra, error) {
			return nil, domain.ErrDescricaoRendaExtraObrigatoria
		},
	}
	r := setupRendaExtraRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"descricao":        "",
		"membro_id":        "membro-1",
		"data_recebimento": "2026-03-08",
		"valor":            2000.00,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/rendas-extras", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrDescricaoRendaExtraObrigatoria.Error(), resp["error"])
}

// ---- GET /rendas-extras ----

func TestListarRendasExtrasHandler_Sucesso(t *testing.T) {
	rendas := []*domain.RendaExtra{
		{ID: "uuid-1", Descricao: "Bônus março", MembroID: "membro-1", DataRecebimento: time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC), Valor: 2000.00},
		{ID: "uuid-2", Descricao: "Venda equipamento", MembroID: "membro-1", DataRecebimento: time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC), Valor: 500.00},
	}
	svc := &MockRendaExtraService{
		ListarFn: func(familiaID string) ([]*domain.RendaExtra, error) {
			return rendas, nil
		},
	}
	r := setupRendaExtraRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rendas-extras", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "uuid-1", resp[0]["id"])
}

func TestListarRendasExtrasHandler_FiltradaPorMes(t *testing.T) {
	rendas := []*domain.RendaExtra{
		{ID: "uuid-1", Descricao: "Bônus março", MembroID: "membro-1", DataRecebimento: time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC), Valor: 2000.00},
	}
	svc := &MockRendaExtraService{
		ListarPorMesFn: func(familiaID string, mes, ano int) ([]*domain.RendaExtra, error) {
			assert.Equal(t, 3, mes)
			assert.Equal(t, 2026, ano)
			return rendas, nil
		},
	}
	r := setupRendaExtraRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rendas-extras?mes=3&ano=2026", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "uuid-1", resp[0]["id"])
}

// ---- GET /rendas-extras/:id ----

func TestBuscarRendaExtraHandler_Existente(t *testing.T) {
	svc := &MockRendaExtraService{
		BuscarFn: func(familiaID, id string) (*domain.RendaExtra, error) {
			return rendaExtraExemplo, nil
		},
	}
	r := setupRendaExtraRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rendas-extras/uuid-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resp["id"])
	assert.Equal(t, "Bônus março", resp["descricao"])
}

func TestBuscarRendaExtraHandler_Inexistente(t *testing.T) {
	svc := &MockRendaExtraService{
		BuscarFn: func(familiaID, id string) (*domain.RendaExtra, error) {
			return nil, domain.ErrRendaExtraNaoEncontrada
		},
	}
	r := setupRendaExtraRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rendas-extras/uuid-inexistente", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrRendaExtraNaoEncontrada.Error(), resp["error"])
}

// ---- PUT /rendas-extras/:id ----

func TestAtualizarRendaExtraHandler_Sucesso(t *testing.T) {
	svc := &MockRendaExtraService{
		AtualizarFn: func(familiaID, id, descricao, membroID string, dataRecebimento time.Time, valor float64) (*domain.RendaExtra, error) {
			return &domain.RendaExtra{
				ID:              id,
				Descricao:       descricao,
				MembroID:        membroID,
				DataRecebimento: dataRecebimento,
				Valor:           valor,
			}, nil
		},
	}
	r := setupRendaExtraRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"descricao":        "Bônus atualizado",
		"membro_id":        "membro-1",
		"data_recebimento": "2026-03-10",
		"valor":            3000.00,
	})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/rendas-extras/uuid-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resp["id"])
	assert.Equal(t, "Bônus atualizado", resp["descricao"])
}

// ---- DELETE /rendas-extras/:id ----

func TestExcluirRendaExtraHandler_Sucesso(t *testing.T) {
	svc := &MockRendaExtraService{
		ExcluirFn: func(familiaID, id string) error {
			return nil
		},
	}
	r := setupRendaExtraRouter(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/rendas-extras/uuid-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestExcluirRendaExtraHandler_Inexistente(t *testing.T) {
	svc := &MockRendaExtraService{
		ExcluirFn: func(familiaID, id string) error {
			return domain.ErrRendaExtraNaoEncontrada
		},
	}
	r := setupRendaExtraRouter(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/rendas-extras/uuid-inexistente", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrRendaExtraNaoEncontrada.Error(), resp["error"])
}
