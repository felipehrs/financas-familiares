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

// MockRendaVariavelService implementa RendaVariavelServiceInterface (declarada no handler) para testes.
type MockRendaVariavelService struct {
	CriarFn        func(descricao, membroID string, mesReferencia, anoReferencia int, valor float64, dataRecebimento time.Time) (*domain.RendaVariavel, error)
	BuscarFn       func(id string) (*domain.RendaVariavel, error)
	ListarFn       func() ([]*domain.RendaVariavel, error)
	ListarPorMesFn func(mes, ano int) ([]*domain.RendaVariavel, error)
	AtualizarFn    func(id, descricao, membroID string, mesReferencia, anoReferencia int, valor float64, dataRecebimento time.Time) (*domain.RendaVariavel, error)
	ExcluirFn      func(id string) error
}

func (m *MockRendaVariavelService) Criar(descricao, membroID string, mesReferencia, anoReferencia int, valor float64, dataRecebimento time.Time) (*domain.RendaVariavel, error) {
	return m.CriarFn(descricao, membroID, mesReferencia, anoReferencia, valor, dataRecebimento)
}

func (m *MockRendaVariavelService) BuscarPorID(id string) (*domain.RendaVariavel, error) {
	return m.BuscarFn(id)
}

func (m *MockRendaVariavelService) Listar() ([]*domain.RendaVariavel, error) {
	return m.ListarFn()
}

func (m *MockRendaVariavelService) ListarPorMes(mes, ano int) ([]*domain.RendaVariavel, error) {
	return m.ListarPorMesFn(mes, ano)
}

func (m *MockRendaVariavelService) Atualizar(id, descricao, membroID string, mesReferencia, anoReferencia int, valor float64, dataRecebimento time.Time) (*domain.RendaVariavel, error) {
	return m.AtualizarFn(id, descricao, membroID, mesReferencia, anoReferencia, valor, dataRecebimento)
}

func (m *MockRendaVariavelService) Excluir(id string) error {
	return m.ExcluirFn(id)
}

func setupRendaVariavelRouter(svc handler.RendaVariavelServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewRendaVariavelHandler(svc)
	v1 := r.Group("/api/v1")
	{
		v1.GET("/rendas-variaveis", h.Listar)
		v1.POST("/rendas-variaveis", h.Criar)
		v1.GET("/rendas-variaveis/:id", h.BuscarPorID)
		v1.PUT("/rendas-variaveis/:id", h.Atualizar)
		v1.DELETE("/rendas-variaveis/:id", h.Excluir)
	}
	return r
}

var rendaVariavelExemplo = &domain.RendaVariavel{
	ID:              "uuid-1",
	Descricao:       "Freelance março",
	MembroID:        "membro-1",
	MesReferencia:   3,
	AnoReferencia:   2026,
	Valor:           1500.00,
	DataRecebimento: time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC),
}

// ---- POST /rendas-variaveis ----

func TestCriarRendaVariavelHandler_Sucesso(t *testing.T) {
	svc := &MockRendaVariavelService{
		CriarFn: func(descricao, membroID string, mesReferencia, anoReferencia int, valor float64, dataRecebimento time.Time) (*domain.RendaVariavel, error) {
			return &domain.RendaVariavel{
				ID:              "uuid-novo",
				Descricao:       descricao,
				MembroID:        membroID,
				MesReferencia:   mesReferencia,
				AnoReferencia:   anoReferencia,
				Valor:           valor,
				DataRecebimento: dataRecebimento,
			}, nil
		},
	}
	r := setupRendaVariavelRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"descricao":        "Freelance março",
		"membro_id":        "membro-1",
		"mes_referencia":   3,
		"ano_referencia":   2026,
		"valor":            1500.00,
		"data_recebimento": "2026-03-08",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/rendas-variaveis", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-novo", resp["id"])
	assert.Equal(t, "Freelance março", resp["descricao"])
	assert.Equal(t, "2026-03-08", resp["data_recebimento"])
}

func TestCriarRendaVariavelHandler_BodyInvalido(t *testing.T) {
	svc := &MockRendaVariavelService{}
	r := setupRendaVariavelRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rendas-variaveis", bytes.NewReader([]byte("nao-e-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCriarRendaVariavelHandler_ErroValidacao(t *testing.T) {
	svc := &MockRendaVariavelService{
		CriarFn: func(descricao, membroID string, mesReferencia, anoReferencia int, valor float64, dataRecebimento time.Time) (*domain.RendaVariavel, error) {
			return nil, domain.ErrDescricaoRendaVariavelObrigatoria
		},
	}
	r := setupRendaVariavelRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"descricao":        "",
		"membro_id":        "membro-1",
		"mes_referencia":   3,
		"ano_referencia":   2026,
		"valor":            1500.00,
		"data_recebimento": "2026-03-08",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/rendas-variaveis", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrDescricaoRendaVariavelObrigatoria.Error(), resp["error"])
}

// ---- GET /rendas-variaveis ----

func TestListarRendasVariaveisHandler_Sucesso(t *testing.T) {
	rendas := []*domain.RendaVariavel{
		{ID: "uuid-1", Descricao: "Freelance março", MembroID: "membro-1", MesReferencia: 3, AnoReferencia: 2026, Valor: 1500.00, DataRecebimento: time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC)},
		{ID: "uuid-2", Descricao: "Comissão venda", MembroID: "membro-1", MesReferencia: 3, AnoReferencia: 2026, Valor: 800.00, DataRecebimento: time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC)},
	}
	svc := &MockRendaVariavelService{
		ListarFn: func() ([]*domain.RendaVariavel, error) {
			return rendas, nil
		},
	}
	r := setupRendaVariavelRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rendas-variaveis", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "uuid-1", resp[0]["id"])
}

func TestListarRendasVariaveisHandler_FiltradaPorMes(t *testing.T) {
	rendas := []*domain.RendaVariavel{
		{ID: "uuid-1", Descricao: "Freelance março", MembroID: "membro-1", MesReferencia: 3, AnoReferencia: 2026, Valor: 1500.00, DataRecebimento: time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC)},
	}
	svc := &MockRendaVariavelService{
		ListarPorMesFn: func(mes, ano int) ([]*domain.RendaVariavel, error) {
			assert.Equal(t, 3, mes)
			assert.Equal(t, 2026, ano)
			return rendas, nil
		},
	}
	r := setupRendaVariavelRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rendas-variaveis?mes=3&ano=2026", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "uuid-1", resp[0]["id"])
}

// ---- GET /rendas-variaveis/:id ----

func TestBuscarRendaVariavelHandler_Existente(t *testing.T) {
	svc := &MockRendaVariavelService{
		BuscarFn: func(id string) (*domain.RendaVariavel, error) {
			return rendaVariavelExemplo, nil
		},
	}
	r := setupRendaVariavelRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rendas-variaveis/uuid-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resp["id"])
	assert.Equal(t, "Freelance março", resp["descricao"])
}

func TestBuscarRendaVariavelHandler_Inexistente(t *testing.T) {
	svc := &MockRendaVariavelService{
		BuscarFn: func(id string) (*domain.RendaVariavel, error) {
			return nil, domain.ErrRendaVariavelNaoEncontrada
		},
	}
	r := setupRendaVariavelRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rendas-variaveis/uuid-inexistente", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrRendaVariavelNaoEncontrada.Error(), resp["error"])
}

// ---- PUT /rendas-variaveis/:id ----

func TestAtualizarRendaVariavelHandler_Sucesso(t *testing.T) {
	svc := &MockRendaVariavelService{
		AtualizarFn: func(id, descricao, membroID string, mesReferencia, anoReferencia int, valor float64, dataRecebimento time.Time) (*domain.RendaVariavel, error) {
			return &domain.RendaVariavel{
				ID:              id,
				Descricao:       descricao,
				MembroID:        membroID,
				MesReferencia:   mesReferencia,
				AnoReferencia:   anoReferencia,
				Valor:           valor,
				DataRecebimento: dataRecebimento,
			}, nil
		},
	}
	r := setupRendaVariavelRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"descricao":        "Freelance atualizado",
		"membro_id":        "membro-1",
		"mes_referencia":   3,
		"ano_referencia":   2026,
		"valor":            2000.00,
		"data_recebimento": "2026-03-10",
	})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/rendas-variaveis/uuid-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resp["id"])
	assert.Equal(t, "Freelance atualizado", resp["descricao"])
}

// ---- DELETE /rendas-variaveis/:id ----

func TestExcluirRendaVariavelHandler_Sucesso(t *testing.T) {
	svc := &MockRendaVariavelService{
		ExcluirFn: func(id string) error {
			return nil
		},
	}
	r := setupRendaVariavelRouter(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/rendas-variaveis/uuid-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestExcluirRendaVariavelHandler_Inexistente(t *testing.T) {
	svc := &MockRendaVariavelService{
		ExcluirFn: func(id string) error {
			return domain.ErrRendaVariavelNaoEncontrada
		},
	}
	r := setupRendaVariavelRouter(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/rendas-variaveis/uuid-inexistente", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrRendaVariavelNaoEncontrada.Error(), resp["error"])
}
