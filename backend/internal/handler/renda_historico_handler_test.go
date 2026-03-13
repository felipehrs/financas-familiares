package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/felipehrs/financas-familiares/backend/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockRendaHistoricoService implementa RendaHistoricoServiceInterface para testes.
type MockRendaHistoricoService struct {
	BuscarHistoricoFn func(familiaID string, filtro domain.FiltroHistoricoRendas) (*domain.HistoricoRendas, error)
}

func (m *MockRendaHistoricoService) BuscarHistorico(familiaID string, filtro domain.FiltroHistoricoRendas) (*domain.HistoricoRendas, error) {
	return m.BuscarHistoricoFn(familiaID, filtro)
}

func setupRendaHistoricoRouter(svc handler.RendaHistoricoServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", "usuario-teste-uuid")
		c.Set("familiaID", "familia-teste-uuid")
		c.Next()
	})
	h := handler.NewRendaHistoricoHandler(svc)
	v1 := r.Group("/api/v1")
	v1.GET("/rendas/historico", h.Historico)
	return r
}

func historicoVazio() *domain.HistoricoRendas {
	return &domain.HistoricoRendas{
		Itens:  []domain.ItemRendaHistorico{},
		Resumo: domain.ResumoRendaHistorico{},
	}
}

func TestRendaHistoricoHandler_SucessoSemFiltros(t *testing.T) {
	svc := &MockRendaHistoricoService{
		BuscarHistoricoFn: func(familiaID string, filtro domain.FiltroHistoricoRendas) (*domain.HistoricoRendas, error) {
			assert.Equal(t, domain.TipoRenda(""), filtro.Tipo)
			assert.Equal(t, "", filtro.MembroID)
			assert.Equal(t, 0, filtro.Mes)
			assert.Equal(t, 0, filtro.Ano)
			return historicoVazio(), nil
		},
	}
	r := setupRendaHistoricoRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rendas/historico", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp, "itens")
	assert.Contains(t, resp, "resumo")
}

func TestRendaHistoricoHandler_SucessoComTodosOsFiltros(t *testing.T) {
	svc := &MockRendaHistoricoService{
		BuscarHistoricoFn: func(familiaID string, filtro domain.FiltroHistoricoRendas) (*domain.HistoricoRendas, error) {
			assert.Equal(t, domain.TipoRendaFixa, filtro.Tipo)
			assert.Equal(t, "membro-abc", filtro.MembroID)
			assert.Equal(t, 3, filtro.Mes)
			assert.Equal(t, 2026, filtro.Ano)
			mes := 3
			ano := 2026
			ativa := true
			return &domain.HistoricoRendas{
				Itens: []domain.ItemRendaHistorico{
					{ID: "rf-1", Tipo: domain.TipoRendaFixa, Descricao: "Salário", MembroID: "membro-abc", Valor: 5000, Mes: &mes, Ano: &ano, Ativa: &ativa},
				},
				Resumo: domain.ResumoRendaHistorico{
					TotalFixas: 5000,
					TotalGeral: 5000,
				},
			}, nil
		},
	}
	r := setupRendaHistoricoRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rendas/historico?tipo=fixa&membro_id=membro-abc&mes=3&ano=2026", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	itens, ok := resp["itens"].([]any)
	require.True(t, ok)
	assert.Len(t, itens, 1)

	resumo, ok := resp["resumo"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, 5000.0, resumo["total_fixas"])
	assert.Equal(t, 5000.0, resumo["total_geral"])
}

func TestRendaHistoricoHandler_TipoInvalido(t *testing.T) {
	svc := &MockRendaHistoricoService{}
	r := setupRendaHistoricoRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rendas/historico?tipo=invalido", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp["error"])
}

func TestRendaHistoricoHandler_MesInvalido(t *testing.T) {
	svc := &MockRendaHistoricoService{}
	r := setupRendaHistoricoRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rendas/historico?mes=13&ano=2026", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp["error"])
}

func TestRendaHistoricoHandler_MesNaoNumerico(t *testing.T) {
	svc := &MockRendaHistoricoService{}
	r := setupRendaHistoricoRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rendas/historico?mes=abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRendaHistoricoHandler_AnoInvalido(t *testing.T) {
	svc := &MockRendaHistoricoService{}
	r := setupRendaHistoricoRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rendas/historico?mes=3&ano=2000", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp["error"])
}

func TestRendaHistoricoHandler_AnoNaoNumerico(t *testing.T) {
	svc := &MockRendaHistoricoService{}
	r := setupRendaHistoricoRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rendas/historico?ano=abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRendaHistoricoHandler_ErroInterno(t *testing.T) {
	svc := &MockRendaHistoricoService{
		BuscarHistoricoFn: func(familiaID string, filtro domain.FiltroHistoricoRendas) (*domain.HistoricoRendas, error) {
			return nil, errors.New("falha no banco de dados")
		},
	}
	r := setupRendaHistoricoRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rendas/historico", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "erro interno", resp["error"])
}

func TestRendaHistoricoHandler_FiltroTiposValidos(t *testing.T) {
	tiposValidos := []string{"fixa", "variavel", "extra", "investimento"}

	for _, tipo := range tiposValidos {
		t.Run("tipo="+tipo, func(t *testing.T) {
			svc := &MockRendaHistoricoService{
				BuscarHistoricoFn: func(familiaID string, filtro domain.FiltroHistoricoRendas) (*domain.HistoricoRendas, error) {
					return historicoVazio(), nil
				},
			}
			r := setupRendaHistoricoRouter(svc)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/rendas/historico?tipo="+tipo, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "tipo=%s deve retornar 200", tipo)
		})
	}
}
