package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/handler"
	"github.com/felipehrs/financas-familiares/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockDashboardService implementa DashboardServiceInterface para testes.
type MockDashboardService struct {
	ResumoMensalFn         func(mes, ano int) (*service.ResumoMensal, error)
	DespesasPorCategoriaFn func(mes, ano int) (*service.ResumoCategorias, error)
	EvolucaoMensalFn       func(qtdMeses int) ([]service.PontoEvolucao, error)
}

func (m *MockDashboardService) ResumoMensal(mes, ano int) (*service.ResumoMensal, error) {
	return m.ResumoMensalFn(mes, ano)
}

func (m *MockDashboardService) DespesasPorCategoria(mes, ano int) (*service.ResumoCategorias, error) {
	if m.DespesasPorCategoriaFn != nil {
		return m.DespesasPorCategoriaFn(mes, ano)
	}
	return nil, nil
}

func (m *MockDashboardService) EvolucaoMensal(qtdMeses int) ([]service.PontoEvolucao, error) {
	if m.EvolucaoMensalFn != nil {
		return m.EvolucaoMensalFn(qtdMeses)
	}
	return []service.PontoEvolucao{}, nil
}

func setupDashboardRouter(svc handler.DashboardServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewDashboardHandler(svc)
	v1 := r.Group("/api/v1")
	{
		v1.GET("/dashboard/resumo", h.ResumoMensal)
		v1.GET("/dashboard/categorias", h.DespesasPorCategoria)
		v1.GET("/dashboard/evolucao", h.EvolucaoMensal)
	}
	return r
}

// ---- Testes de GET /dashboard/resumo ----

func TestResumoMensalHandler_SucessoComParamsExplicitos(t *testing.T) {
	svc := &MockDashboardService{
		ResumoMensalFn: func(mes, ano int) (*service.ResumoMensal, error) {
			assert.Equal(t, 3, mes)
			assert.Equal(t, 2026, ano)
			return &service.ResumoMensal{
				Mes:                     mes,
				Ano:                     ano,
				TotalRendasOperacionais: 5000.00,
				TotalDespesas:           1000.00,
				Saldo:                   4000.00,
			}, nil
		},
	}
	r := setupDashboardRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/resumo?mes=3&ano=2026", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(3), resp["mes"])
	assert.Equal(t, float64(2026), resp["ano"])
	assert.Equal(t, 5000.00, resp["total_rendas_operacionais"])
	assert.Equal(t, 1000.00, resp["total_despesas"])
	assert.Equal(t, 4000.00, resp["saldo"])
}

func TestResumoMensalHandler_SucessoSemParams(t *testing.T) {
	agora := time.Now()
	mesAtual := int(agora.Month())
	anoAtual := agora.Year()

	svc := &MockDashboardService{
		ResumoMensalFn: func(mes, ano int) (*service.ResumoMensal, error) {
			assert.Equal(t, mesAtual, mes)
			assert.Equal(t, anoAtual, ano)
			return &service.ResumoMensal{
				Mes:                     mes,
				Ano:                     ano,
				TotalRendasOperacionais: 3000.00,
				TotalDespesas:           500.00,
				Saldo:                   2500.00,
			}, nil
		},
	}
	r := setupDashboardRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/resumo", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(mesAtual), resp["mes"])
	assert.Equal(t, float64(anoAtual), resp["ano"])
}

func TestResumoMensalHandler_MesInvalido(t *testing.T) {
	svc := &MockDashboardService{}
	r := setupDashboardRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/resumo?mes=13&ano=2026", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp["error"])
}

func TestResumoMensalHandler_MesNaoNumerico(t *testing.T) {
	svc := &MockDashboardService{}
	r := setupDashboardRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/resumo?mes=abc&ano=2026", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestResumoMensalHandler_AnoInvalido(t *testing.T) {
	svc := &MockDashboardService{}
	r := setupDashboardRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/resumo?mes=3&ano=2000", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp["error"])
}

func TestResumoMensalHandler_AnoNaoNumerico(t *testing.T) {
	svc := &MockDashboardService{}
	r := setupDashboardRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/resumo?mes=3&ano=abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestResumoMensalHandler_ErroInterno(t *testing.T) {
	svc := &MockDashboardService{
		ResumoMensalFn: func(mes, ano int) (*service.ResumoMensal, error) {
			return nil, errors.New("falha no banco de dados")
		},
	}
	r := setupDashboardRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/resumo?mes=3&ano=2026", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "erro interno", resp["error"])
}

// ---- Testes de GET /dashboard/categorias ----

func TestDespesasPorCategoriaHandler_SucessoComParams(t *testing.T) {
	svc := &MockDashboardService{
		DespesasPorCategoriaFn: func(mes, ano int) (*service.ResumoCategorias, error) {
			assert.Equal(t, 3, mes)
			assert.Equal(t, 2026, ano)
			return &service.ResumoCategorias{
				Mes:           mes,
				Ano:           ano,
				TotalDespesas: 1000.0,
				Categorias: []service.CategoriaDespesa{
					{Nome: "Alimentação", Total: 650.0, Percentual: 65.0},
					{Nome: "Transporte", Total: 350.0, Percentual: 35.0},
				},
			}, nil
		},
	}
	r := setupDashboardRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/categorias?mes=3&ano=2026", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(3), resp["mes"])
	assert.Equal(t, float64(2026), resp["ano"])
	assert.Equal(t, 1000.0, resp["total_despesas"])
	categorias, ok := resp["categorias"].([]any)
	require.True(t, ok)
	assert.Len(t, categorias, 2)
}

func TestDespesasPorCategoriaHandler_SucessoSemParams(t *testing.T) {
	agora := time.Now()
	mesAtual := int(agora.Month())
	anoAtual := agora.Year()

	svc := &MockDashboardService{
		DespesasPorCategoriaFn: func(mes, ano int) (*service.ResumoCategorias, error) {
			assert.Equal(t, mesAtual, mes)
			assert.Equal(t, anoAtual, ano)
			return &service.ResumoCategorias{
				Mes:           mes,
				Ano:           ano,
				TotalDespesas: 0.0,
				Categorias:    []service.CategoriaDespesa{},
			}, nil
		},
	}
	r := setupDashboardRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/categorias", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestDespesasPorCategoriaHandler_MesInvalido(t *testing.T) {
	svc := &MockDashboardService{}
	r := setupDashboardRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/categorias?mes=13&ano=2026", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp["error"])
}

func TestDespesasPorCategoriaHandler_AnoInvalido(t *testing.T) {
	svc := &MockDashboardService{}
	r := setupDashboardRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/categorias?mes=3&ano=1999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp["error"])
}

func TestDespesasPorCategoriaHandler_ErroInterno(t *testing.T) {
	svc := &MockDashboardService{
		DespesasPorCategoriaFn: func(mes, ano int) (*service.ResumoCategorias, error) {
			return nil, errors.New("falha no banco de dados")
		},
	}
	r := setupDashboardRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/categorias?mes=3&ano=2026", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "erro interno", resp["error"])
}

// ---- Testes de GET /dashboard/evolucao ----

func TestEvolucaoMensalHandler_Sucesso(t *testing.T) {
	svc := &MockDashboardService{
		EvolucaoMensalFn: func(qtdMeses int) ([]service.PontoEvolucao, error) {
			assert.Equal(t, 12, qtdMeses)
			return []service.PontoEvolucao{
				{Mes: 4, Ano: 2025, TotalRendas: 5000.00, TotalDespesas: 1000.00, Saldo: 4000.00},
				{Mes: 3, Ano: 2026, TotalRendas: 6000.00, TotalDespesas: 1500.00, Saldo: 4500.00},
			}, nil
		},
	}
	r := setupDashboardRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/evolucao", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, float64(4), resp[0]["mes"])
	assert.Equal(t, float64(2025), resp[0]["ano"])
	assert.Equal(t, 5000.00, resp[0]["total_rendas"])
}

func TestEvolucaoMensalHandler_ErroInterno(t *testing.T) {
	svc := &MockDashboardService{
		EvolucaoMensalFn: func(qtdMeses int) ([]service.PontoEvolucao, error) {
			return nil, errors.New("falha no banco de dados")
		},
	}
	r := setupDashboardRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/evolucao", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "erro interno", resp["error"])
}
