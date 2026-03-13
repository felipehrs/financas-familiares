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

// MockCategoriaService implementa CategoriaServiceInterface (declarada no handler) para testes.
type MockCategoriaService struct {
	CriarFn     func(familiaID, nome string) (*domain.Categoria, error)
	ListarFn    func(familiaID string) ([]*domain.Categoria, error)
	AtualizarFn func(familiaID, id, nome string) (*domain.Categoria, error)
	ExcluirFn   func(familiaID, id string) error
}

func (m *MockCategoriaService) Criar(familiaID, nome string) (*domain.Categoria, error) {
	return m.CriarFn(familiaID, nome)
}

func (m *MockCategoriaService) Listar(familiaID string) ([]*domain.Categoria, error) {
	return m.ListarFn(familiaID)
}

func (m *MockCategoriaService) Atualizar(familiaID, id, nome string) (*domain.Categoria, error) {
	return m.AtualizarFn(familiaID, id, nome)
}

func (m *MockCategoriaService) Excluir(familiaID, id string) error {
	return m.ExcluirFn(familiaID, id)
}

func setupCategoriaRouter(svc handler.CategoriaServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", "usuario-teste-uuid")
		c.Set("familiaID", "familia-teste-uuid")
		c.Next()
	})
	h := handler.NewCategoriaHandler(svc)
	v1 := r.Group("/api/v1")
	{
		v1.GET("/categorias", h.Listar)
		v1.POST("/categorias", h.Criar)
		v1.PUT("/categorias/:id", h.Atualizar)
		v1.DELETE("/categorias/:id", h.Excluir)
	}
	return r
}

// ---- GET /categorias ----

func TestListarCategoriasHandler_Sucesso(t *testing.T) {
	categorias := []*domain.Categoria{
		{ID: "uuid-1", Nome: "Alimentação"},
		{ID: "uuid-2", Nome: "Transporte"},
	}
	svc := &MockCategoriaService{
		ListarFn: func(familiaID string) ([]*domain.Categoria, error) {
			return categorias, nil
		},
	}
	r := setupCategoriaRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/categorias", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "uuid-1", resp[0]["id"])
	assert.Equal(t, "Alimentação", resp[0]["nome"])
}

// ---- POST /categorias ----

func TestCriarCategoriaHandler_Sucesso(t *testing.T) {
	svc := &MockCategoriaService{
		CriarFn: func(familiaID, nome string) (*domain.Categoria, error) {
			return &domain.Categoria{
				ID:   "uuid-novo",
				Nome: nome,
			}, nil
		},
	}
	r := setupCategoriaRouter(svc)

	body, _ := json.Marshal(map[string]string{"nome": "Lazer"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/categorias", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-novo", resp["id"])
	assert.Equal(t, "Lazer", resp["nome"])
}

func TestCriarCategoriaHandler_NomeVazio(t *testing.T) {
	svc := &MockCategoriaService{
		CriarFn: func(familiaID, nome string) (*domain.Categoria, error) {
			return nil, domain.ErrNomeCategoriaObrigatorio
		},
	}
	r := setupCategoriaRouter(svc)

	body, _ := json.Marshal(map[string]string{"nome": ""})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/categorias", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrNomeCategoriaObrigatorio.Error(), resp["error"])
}

func TestCriarCategoriaHandler_BodyInvalido(t *testing.T) {
	svc := &MockCategoriaService{}
	r := setupCategoriaRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/categorias", bytes.NewReader([]byte("nao-e-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ---- PUT /categorias/:id ----

func TestAtualizarCategoriaHandler_Sucesso(t *testing.T) {
	svc := &MockCategoriaService{
		AtualizarFn: func(familiaID, id, nome string) (*domain.Categoria, error) {
			return &domain.Categoria{
				ID:   id,
				Nome: nome,
			}, nil
		},
	}
	r := setupCategoriaRouter(svc)

	body, _ := json.Marshal(map[string]string{"nome": "Alimentação e Bebidas"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/categorias/uuid-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resp["id"])
	assert.Equal(t, "Alimentação e Bebidas", resp["nome"])
}

func TestAtualizarCategoriaHandler_NaoEncontrada(t *testing.T) {
	svc := &MockCategoriaService{
		AtualizarFn: func(familiaID, id, nome string) (*domain.Categoria, error) {
			return nil, domain.ErrCategoriaNaoEncontrada
		},
	}
	r := setupCategoriaRouter(svc)

	body, _ := json.Marshal(map[string]string{"nome": "Qualquer"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/categorias/uuid-inexistente", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrCategoriaNaoEncontrada.Error(), resp["error"])
}

// ---- DELETE /categorias/:id ----

func TestExcluirCategoriaHandler_Sucesso(t *testing.T) {
	svc := &MockCategoriaService{
		ExcluirFn: func(familiaID, id string) error {
			return nil
		},
	}
	r := setupCategoriaRouter(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/categorias/uuid-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestExcluirCategoriaHandler_NaoEncontrada(t *testing.T) {
	svc := &MockCategoriaService{
		ExcluirFn: func(familiaID, id string) error {
			return domain.ErrCategoriaNaoEncontrada
		},
	}
	r := setupCategoriaRouter(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/categorias/uuid-inexistente", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrCategoriaNaoEncontrada.Error(), resp["error"])
}

func TestExcluirCategoriaHandler_ComVinculos(t *testing.T) {
	svc := &MockCategoriaService{
		ExcluirFn: func(familiaID, id string) error {
			return domain.ErrCategoriaComVinculos
		},
	}
	r := setupCategoriaRouter(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/categorias/uuid-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusConflict, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrCategoriaComVinculos.Error(), resp["error"])
}
