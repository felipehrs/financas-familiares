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

// MockDespesaCartaoService implementa DespesaCartaoServiceInterface para testes.
type MockDespesaCartaoService struct {
	CriarFn           func(cartaoID, descricao string, categoriaID *string, dataCompra time.Time, valorTotal float64, numeroParcelas int) ([]*domain.DespesaCartao, error)
	ListarPorCartaoFn func(cartaoID string) ([]*domain.DespesaCartao, error)
	ListarPorFaturaFn func(cartaoID string, mes, ano int) ([]*domain.DespesaCartao, error)
	BuscarPorIDFn     func(id string) (*domain.DespesaCartao, error)
	ExcluirFn         func(id string) error
}

func (m *MockDespesaCartaoService) Criar(cartaoID, descricao string, categoriaID *string, dataCompra time.Time, valorTotal float64, numeroParcelas int) ([]*domain.DespesaCartao, error) {
	return m.CriarFn(cartaoID, descricao, categoriaID, dataCompra, valorTotal, numeroParcelas)
}

func (m *MockDespesaCartaoService) ListarPorCartao(cartaoID string) ([]*domain.DespesaCartao, error) {
	return m.ListarPorCartaoFn(cartaoID)
}

func (m *MockDespesaCartaoService) ListarPorFatura(cartaoID string, mes, ano int) ([]*domain.DespesaCartao, error) {
	return m.ListarPorFaturaFn(cartaoID, mes, ano)
}

func (m *MockDespesaCartaoService) BuscarPorID(id string) (*domain.DespesaCartao, error) {
	return m.BuscarPorIDFn(id)
}

func (m *MockDespesaCartaoService) Excluir(id string) error {
	return m.ExcluirFn(id)
}

func setupDespesaRouter(svc handler.DespesaCartaoServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewDespesaCartaoHandler(svc)
	v1 := r.Group("/api/v1")
	{
		v1.GET("/cartoes/:id/despesas", h.ListarPorCartao)
		v1.GET("/cartoes/:id/despesas/fatura", h.ListarPorFatura)
		v1.POST("/cartoes/:id/despesas", h.Criar)
		v1.DELETE("/despesas/:id", h.Excluir)
	}
	return r
}

func despesaFixa() *domain.DespesaCartao {
	return &domain.DespesaCartao{
		ID:             "d-uuid-1",
		CompraID:       "compra-uuid-1",
		CartaoID:       "cartao-1",
		Descricao:      "Supermercado",
		DataCompra:     time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC),
		ValorTotal:     150.00,
		NumeroParcelas: 1,
		ParcelaNumero:  1,
		ValorParcela:   150.00,
		FaturaMes:      3,
		FaturaAno:      2026,
	}
}

// ---- GET /cartoes/:id/despesas ----

func TestListarDespesasPorCartaoHandler_Sucesso(t *testing.T) {
	despesas := []*domain.DespesaCartao{despesaFixa()}
	svc := &MockDespesaCartaoService{
		ListarPorCartaoFn: func(cartaoID string) ([]*domain.DespesaCartao, error) {
			return despesas, nil
		},
	}
	r := setupDespesaRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cartoes/cartao-1/despesas", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "d-uuid-1", resp[0]["id"])
	assert.Equal(t, "Supermercado", resp[0]["descricao"])
	assert.Equal(t, "MAR/26", resp[0]["fatura"])
	assert.Equal(t, "compra-uuid-1", resp[0]["compra_id"])
	assert.Equal(t, float64(1), resp[0]["parcela_numero"])
}

func TestListarDespesasPorCartaoHandler_ListaVazia(t *testing.T) {
	svc := &MockDespesaCartaoService{
		ListarPorCartaoFn: func(cartaoID string) ([]*domain.DespesaCartao, error) {
			return []*domain.DespesaCartao{}, nil
		},
	}
	r := setupDespesaRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cartoes/cartao-1/despesas", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Empty(t, resp)
}

// ---- GET /cartoes/:id/despesas/fatura ----

func TestListarDespesasPorFaturaHandler_Sucesso(t *testing.T) {
	despesas := []*domain.DespesaCartao{despesaFixa()}
	svc := &MockDespesaCartaoService{
		ListarPorFaturaFn: func(cartaoID string, mes, ano int) ([]*domain.DespesaCartao, error) {
			return despesas, nil
		},
	}
	r := setupDespesaRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cartoes/cartao-1/despesas/fatura?mes=3&ano=2026", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "MAR/26", resp[0]["fatura"])
}

func TestListarDespesasPorFaturaHandler_ParametrosFaltando(t *testing.T) {
	svc := &MockDespesaCartaoService{}
	r := setupDespesaRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cartoes/cartao-1/despesas/fatura?mes=3", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListarDespesasPorFaturaHandler_MesInvalido(t *testing.T) {
	svc := &MockDespesaCartaoService{}
	r := setupDespesaRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cartoes/cartao-1/despesas/fatura?mes=13&ano=2026", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ---- POST /cartoes/:id/despesas ----

func TestCriarDespesaCartaoHandler_Sucesso_AVista(t *testing.T) {
	svc := &MockDespesaCartaoService{
		CriarFn: func(cartaoID, descricao string, categoriaID *string, dataCompra time.Time, valorTotal float64, numeroParcelas int) ([]*domain.DespesaCartao, error) {
			return []*domain.DespesaCartao{
				{
					ID:             "d-novo",
					CompraID:       "compra-novo",
					CartaoID:       cartaoID,
					Descricao:      descricao,
					DataCompra:     dataCompra,
					ValorTotal:     valorTotal,
					NumeroParcelas: 1,
					ParcelaNumero:  1,
					ValorParcela:   valorTotal,
					FaturaMes:      3,
					FaturaAno:      2026,
				},
			}, nil
		},
	}
	r := setupDespesaRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"descricao":   "Supermercado",
		"data_compra": "2026-03-05",
		"valor_total": 150.00,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cartoes/cartao-1/despesas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.Len(t, resp, 1)
	assert.Equal(t, "d-novo", resp[0]["id"])
	assert.Equal(t, "Supermercado", resp[0]["descricao"])
	assert.Equal(t, "MAR/26", resp[0]["fatura"])
	assert.Equal(t, float64(1), resp[0]["numero_parcelas"])
	assert.Equal(t, float64(1), resp[0]["parcela_numero"])
}

func TestCriarDespesaCartaoHandler_3Parcelas_RetornaArray3Elementos(t *testing.T) {
	var numeroParcelas int
	svc := &MockDespesaCartaoService{
		CriarFn: func(cartaoID, descricao string, categoriaID *string, dataCompra time.Time, valorTotal float64, np int) ([]*domain.DespesaCartao, error) {
			numeroParcelas = np
			return []*domain.DespesaCartao{
				{ID: "d-1", CompraID: "c-1", CartaoID: cartaoID, Descricao: descricao, DataCompra: dataCompra, ValorTotal: valorTotal, NumeroParcelas: 3, ParcelaNumero: 1, ValorParcela: 100.00, FaturaMes: 3, FaturaAno: 2026},
				{ID: "d-2", CompraID: "c-1", CartaoID: cartaoID, Descricao: descricao, DataCompra: dataCompra, ValorTotal: valorTotal, NumeroParcelas: 3, ParcelaNumero: 2, ValorParcela: 100.00, FaturaMes: 4, FaturaAno: 2026},
				{ID: "d-3", CompraID: "c-1", CartaoID: cartaoID, Descricao: descricao, DataCompra: dataCompra, ValorTotal: valorTotal, NumeroParcelas: 3, ParcelaNumero: 3, ValorParcela: 100.00, FaturaMes: 5, FaturaAno: 2026},
			}, nil
		},
	}
	r := setupDespesaRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"descricao":        "Notebook",
		"data_compra":      "2026-03-05",
		"valor_total":      300.00,
		"numero_parcelas":  3,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cartoes/cartao-1/despesas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.Len(t, resp, 3)
	assert.Equal(t, 3, numeroParcelas)
	assert.Equal(t, float64(1), resp[0]["parcela_numero"])
	assert.Equal(t, float64(2), resp[1]["parcela_numero"])
	assert.Equal(t, float64(3), resp[2]["parcela_numero"])
	assert.Equal(t, "MAR/26", resp[0]["fatura"])
	assert.Equal(t, "ABR/26", resp[1]["fatura"])
	assert.Equal(t, "MAI/26", resp[2]["fatura"])
}

func TestCriarDespesaCartaoHandler_NumeroParcelas0_DefaultsPara1(t *testing.T) {
	var numeroParcelas int
	svc := &MockDespesaCartaoService{
		CriarFn: func(cartaoID, descricao string, categoriaID *string, dataCompra time.Time, valorTotal float64, np int) ([]*domain.DespesaCartao, error) {
			numeroParcelas = np
			return []*domain.DespesaCartao{
				{ID: "d-1", CompraID: "c-1", CartaoID: cartaoID, Descricao: descricao, DataCompra: dataCompra, ValorTotal: valorTotal, NumeroParcelas: 1, ParcelaNumero: 1, ValorParcela: valorTotal, FaturaMes: 3, FaturaAno: 2026},
			}, nil
		},
	}
	r := setupDespesaRouter(svc)

	// Envia sem numero_parcelas (omitido → zero value em Go)
	body, _ := json.Marshal(map[string]any{
		"descricao":   "Supermercado",
		"data_compra": "2026-03-05",
		"valor_total": 150.00,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cartoes/cartao-1/despesas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, 1, numeroParcelas)
}

func TestCriarDespesaCartaoHandler_ErrNumeroParcelas_Returns400(t *testing.T) {
	svc := &MockDespesaCartaoService{
		CriarFn: func(cartaoID, descricao string, categoriaID *string, dataCompra time.Time, valorTotal float64, numeroParcelas int) ([]*domain.DespesaCartao, error) {
			return nil, domain.ErrNumeroParcelas
		},
	}
	r := setupDespesaRouter(svc)

	// Envia numero_parcelas=-1 para forçar erro no service
	body, _ := json.Marshal(map[string]any{
		"descricao":       "Supermercado",
		"data_compra":     "2026-03-05",
		"valor_total":     150.00,
		"numero_parcelas": -1,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cartoes/cartao-1/despesas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCriarDespesaCartaoHandler_DataInvalida(t *testing.T) {
	svc := &MockDespesaCartaoService{}
	r := setupDespesaRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"descricao":   "Supermercado",
		"data_compra": "data-invalida",
		"valor_total": 150.00,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cartoes/cartao-1/despesas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCriarDespesaCartaoHandler_BodyInvalido(t *testing.T) {
	svc := &MockDespesaCartaoService{}
	r := setupDespesaRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cartoes/cartao-1/despesas", bytes.NewReader([]byte("nao-e-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCriarDespesaCartaoHandler_DescricaoVazia(t *testing.T) {
	svc := &MockDespesaCartaoService{
		CriarFn: func(cartaoID, descricao string, categoriaID *string, dataCompra time.Time, valorTotal float64, numeroParcelas int) ([]*domain.DespesaCartao, error) {
			return nil, domain.ErrDescricaoObrigatoria
		},
	}
	r := setupDespesaRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"descricao":   "",
		"data_compra": "2026-03-05",
		"valor_total": 150.00,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cartoes/cartao-1/despesas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrDescricaoObrigatoria.Error(), resp["error"])
}

func TestCriarDespesaCartaoHandler_ValorInvalido(t *testing.T) {
	svc := &MockDespesaCartaoService{
		CriarFn: func(cartaoID, descricao string, categoriaID *string, dataCompra time.Time, valorTotal float64, numeroParcelas int) ([]*domain.DespesaCartao, error) {
			return nil, domain.ErrValorTotalInvalido
		},
	}
	r := setupDespesaRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"descricao":   "Supermercado",
		"data_compra": "2026-03-05",
		"valor_total": 0,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cartoes/cartao-1/despesas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCriarDespesaCartaoHandler_CartaoNaoEncontrado(t *testing.T) {
	svc := &MockDespesaCartaoService{
		CriarFn: func(cartaoID, descricao string, categoriaID *string, dataCompra time.Time, valorTotal float64, numeroParcelas int) ([]*domain.DespesaCartao, error) {
			return nil, domain.ErrCartaoNaoEncontrado
		},
	}
	r := setupDespesaRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"descricao":   "Supermercado",
		"data_compra": "2026-03-05",
		"valor_total": 150.00,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cartoes/cartao-inexistente/despesas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

// ---- DELETE /despesas/:id ----

func TestExcluirDespesaCartaoHandler_Sucesso(t *testing.T) {
	svc := &MockDespesaCartaoService{
		ExcluirFn: func(id string) error { return nil },
	}
	r := setupDespesaRouter(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/despesas/d-uuid-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "despesa excluída com sucesso", resp["message"])
}

func TestExcluirDespesaCartaoHandler_NaoEncontrada(t *testing.T) {
	svc := &MockDespesaCartaoService{
		ExcluirFn: func(id string) error { return domain.ErrDespesaCartaoNaoEncontrada },
	}
	r := setupDespesaRouter(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/despesas/d-inexistente", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}
