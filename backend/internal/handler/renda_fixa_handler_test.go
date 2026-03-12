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

// MockRendaFixaService implementa RendaFixaServiceInterface para testes.
type MockRendaFixaService struct {
	CriarFn                func(familiaID, descricao, membroID string, valor float64, diaRecebimento int, dataInicio time.Time, dataFim *time.Time) (*domain.RendaFixa, error)
	BuscarFn               func(familiaID, id string) (*domain.RendaFixa, error)
	ListarFn               func(familiaID string) ([]*domain.RendaFixa, error)
	ListarAtivasFn         func(familiaID string) ([]*domain.RendaFixa, error)
	ListarVigentesPorMesFn func(familiaID string, mes, ano int) ([]*domain.RendaFixa, error)
	AtualizarFn            func(familiaID, id, descricao, membroID string, valor float64, diaRecebimento int, ativa bool, dataInicio time.Time, dataFim *time.Time) (*domain.RendaFixa, error)
	InativarFn             func(familiaID, id string) error
}

func (m *MockRendaFixaService) Criar(familiaID, descricao, membroID string, valor float64, diaRecebimento int, dataInicio time.Time, dataFim *time.Time) (*domain.RendaFixa, error) {
	return m.CriarFn(familiaID, descricao, membroID, valor, diaRecebimento, dataInicio, dataFim)
}

func (m *MockRendaFixaService) BuscarPorID(familiaID, id string) (*domain.RendaFixa, error) {
	return m.BuscarFn(familiaID, id)
}

func (m *MockRendaFixaService) Listar(familiaID string) ([]*domain.RendaFixa, error) {
	return m.ListarFn(familiaID)
}

func (m *MockRendaFixaService) ListarAtivas(familiaID string) ([]*domain.RendaFixa, error) {
	return m.ListarAtivasFn(familiaID)
}

func (m *MockRendaFixaService) ListarVigentesPorMes(familiaID string, mes, ano int) ([]*domain.RendaFixa, error) {
	return m.ListarVigentesPorMesFn(familiaID, mes, ano)
}

func (m *MockRendaFixaService) Atualizar(familiaID, id, descricao, membroID string, valor float64, diaRecebimento int, ativa bool, dataInicio time.Time, dataFim *time.Time) (*domain.RendaFixa, error) {
	return m.AtualizarFn(familiaID, id, descricao, membroID, valor, diaRecebimento, ativa, dataInicio, dataFim)
}

func (m *MockRendaFixaService) Inativar(familiaID, id string) error {
	return m.InativarFn(familiaID, id)
}

func setupRendaFixaRouter(svc handler.RendaFixaServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", "usuario-teste-uuid")
		c.Set("familiaID", "familia-teste-uuid")
		c.Next()
	})
	h := handler.NewRendaFixaHandler(svc)
	v1 := r.Group("/api/v1")
	{
		v1.GET("/rendas-fixas", h.Listar)
		v1.POST("/rendas-fixas", h.Criar)
		v1.GET("/rendas-fixas/vigentes", h.ListarVigentesPorMes)
		v1.GET("/rendas-fixas/:id", h.BuscarPorID)
		v1.PUT("/rendas-fixas/:id", h.Atualizar)
		v1.PATCH("/rendas-fixas/:id/inativar", h.Inativar)
	}
	return r
}

// ---- GET /rendas-fixas ----

func TestListarRendasFixasHandler_Sucesso(t *testing.T) {
	rendas := []*domain.RendaFixa{
		{ID: "uuid-1", Descricao: "Salário", MembroID: "m-1", Valor: 5000.00, DiaRecebimento: 5, Ativa: true},
		{ID: "uuid-2", Descricao: "Aposentadoria", MembroID: "m-2", Valor: 2000.00, DiaRecebimento: 10, Ativa: true},
	}
	svc := &MockRendaFixaService{
		ListarFn: func(familiaID string) ([]*domain.RendaFixa, error) { return rendas, nil },
	}
	r := setupRendaFixaRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rendas-fixas", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "uuid-1", resp[0]["id"])
	assert.Equal(t, "Salário", resp[0]["descricao"])
}

// ---- POST /rendas-fixas ----

func TestCriarRendaFixaHandler_Sucesso(t *testing.T) {
	svc := &MockRendaFixaService{
		CriarFn: func(familiaID, descricao, membroID string, valor float64, diaRecebimento int, dataInicio time.Time, dataFim *time.Time) (*domain.RendaFixa, error) {
			return &domain.RendaFixa{
				ID:             "uuid-novo",
				Descricao:      descricao,
				MembroID:       membroID,
				Valor:          valor,
				DiaRecebimento: diaRecebimento,
				Ativa:          true,
				DataInicio:     dataInicio,
			}, nil
		},
	}
	r := setupRendaFixaRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"descricao":       "Salário",
		"membro_id":       "membro-1",
		"valor":           5000.00,
		"dia_recebimento": 5,
		"data_inicio":     "2026-01-01",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/rendas-fixas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-novo", resp["id"])
	assert.Equal(t, "Salário", resp["descricao"])
	assert.Equal(t, true, resp["ativa"])
	assert.Equal(t, "2026-01-01", resp["data_inicio"])
}

func TestCriarRendaFixaHandler_DataInicioAusente(t *testing.T) {
	svc := &MockRendaFixaService{}
	r := setupRendaFixaRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"descricao":       "Salário",
		"membro_id":       "membro-1",
		"valor":           5000.00,
		"dia_recebimento": 5,
		// data_inicio ausente
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/rendas-fixas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp["error"])
}

func TestCriarRendaFixaHandler_DescricaoVazia(t *testing.T) {
	svc := &MockRendaFixaService{
		CriarFn: func(familiaID, descricao, membroID string, valor float64, diaRecebimento int, dataInicio time.Time, dataFim *time.Time) (*domain.RendaFixa, error) {
			return nil, domain.ErrDescricaoRendaObrigatoria
		},
	}
	r := setupRendaFixaRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"descricao": "", "membro_id": "membro-1", "valor": 5000.00, "dia_recebimento": 5, "data_inicio": "2026-01-01",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/rendas-fixas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrDescricaoRendaObrigatoria.Error(), resp["error"])
}

func TestCriarRendaFixaHandler_BodyInvalido(t *testing.T) {
	svc := &MockRendaFixaService{}
	r := setupRendaFixaRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rendas-fixas", bytes.NewReader([]byte("nao-e-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ---- GET /rendas-fixas/vigentes ----

func TestListarVigentesPorMesHandler_Sucesso(t *testing.T) {
	rendas := []*domain.RendaFixa{
		{ID: "uuid-1", Descricao: "Salário", MembroID: "m-1", Valor: 5000.00, DiaRecebimento: 5, Ativa: true, DataInicio: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
		{ID: "uuid-2", Descricao: "Aposentadoria", MembroID: "m-2", Valor: 2000.00, DiaRecebimento: 10, Ativa: true, DataInicio: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
	svc := &MockRendaFixaService{
		ListarVigentesPorMesFn: func(familiaID string, mes, ano int) ([]*domain.RendaFixa, error) {
			return rendas, nil
		},
	}
	r := setupRendaFixaRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rendas-fixas/vigentes?mes=3&ano=2026", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "uuid-1", resp[0]["id"])
	assert.Equal(t, "2026-01-01", resp[0]["data_inicio"])
}

// ---- GET /rendas-fixas/:id ----

func TestBuscarRendaFixaHandler_Sucesso(t *testing.T) {
	svc := &MockRendaFixaService{
		BuscarFn: func(familiaID, id string) (*domain.RendaFixa, error) {
			return &domain.RendaFixa{
				ID:             id,
				Descricao:      "Salário",
				MembroID:       "membro-1",
				Valor:          5000.00,
				DiaRecebimento: 5,
				Ativa:          true,
			}, nil
		},
	}
	r := setupRendaFixaRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rendas-fixas/uuid-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resp["id"])
}

func TestBuscarRendaFixaHandler_NaoEncontrada(t *testing.T) {
	svc := &MockRendaFixaService{
		BuscarFn: func(familiaID, id string) (*domain.RendaFixa, error) {
			return nil, domain.ErrRendaFixaNaoEncontrada
		},
	}
	r := setupRendaFixaRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rendas-fixas/uuid-inexistente", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

// ---- PUT /rendas-fixas/:id ----

func TestAtualizarRendaFixaHandler_Sucesso(t *testing.T) {
	svc := &MockRendaFixaService{
		AtualizarFn: func(familiaID, id, descricao, membroID string, valor float64, diaRecebimento int, ativa bool, dataInicio time.Time, dataFim *time.Time) (*domain.RendaFixa, error) {
			return &domain.RendaFixa{
				ID:             id,
				Descricao:      descricao,
				MembroID:       membroID,
				Valor:          valor,
				DiaRecebimento: diaRecebimento,
				Ativa:          ativa,
				DataInicio:     dataInicio,
			}, nil
		},
	}
	r := setupRendaFixaRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"descricao": "Salário Atualizado", "membro_id": "membro-1",
		"valor": 6000.00, "dia_recebimento": 10, "ativa": true, "data_inicio": "2026-01-01",
	})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/rendas-fixas/uuid-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resp["id"])
	assert.Equal(t, "Salário Atualizado", resp["descricao"])
}

func TestAtualizarRendaFixaHandler_NaoEncontrada(t *testing.T) {
	svc := &MockRendaFixaService{
		AtualizarFn: func(familiaID, id, descricao, membroID string, valor float64, diaRecebimento int, ativa bool, dataInicio time.Time, dataFim *time.Time) (*domain.RendaFixa, error) {
			return nil, domain.ErrRendaFixaNaoEncontrada
		},
	}
	r := setupRendaFixaRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"descricao": "Salário", "membro_id": "membro-1", "valor": 5000.00, "dia_recebimento": 5, "ativa": true, "data_inicio": "2026-01-01",
	})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/rendas-fixas/uuid-inexistente", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

// ---- PATCH /rendas-fixas/:id/inativar ----

func TestInativarRendaFixaHandler_Sucesso(t *testing.T) {
	svc := &MockRendaFixaService{
		InativarFn: func(familiaID, id string) error { return nil },
	}
	r := setupRendaFixaRouter(svc)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/rendas-fixas/uuid-1/inativar", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "renda fixa inativada com sucesso", resp["message"])
}

func TestInativarRendaFixaHandler_NaoEncontrada(t *testing.T) {
	svc := &MockRendaFixaService{
		InativarFn: func(familiaID, id string) error { return domain.ErrRendaFixaNaoEncontrada },
	}
	r := setupRendaFixaRouter(svc)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/rendas-fixas/uuid-inexistente/inativar", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}
