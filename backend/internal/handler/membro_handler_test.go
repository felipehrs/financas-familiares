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

// MockMembroService implementa MembroServiceInterface (declarada no handler) para testes.
type MockMembroService struct {
	CriarFn     func(familiaID, nome, relacionamento string) (*domain.Membro, error)
	BuscarFn    func(familiaID, id string) (*domain.Membro, error)
	ListarFn    func(familiaID string) ([]*domain.Membro, error)
	AtualizarFn func(familiaID, id, nome, relacionamento string, ativo bool) (*domain.Membro, error)
	InativarFn  func(familiaID, id string) error
}

func (m *MockMembroService) Criar(familiaID, nome, relacionamento string) (*domain.Membro, error) {
	return m.CriarFn(familiaID, nome, relacionamento)
}

func (m *MockMembroService) BuscarPorID(familiaID, id string) (*domain.Membro, error) {
	return m.BuscarFn(familiaID, id)
}

func (m *MockMembroService) Listar(familiaID string) ([]*domain.Membro, error) {
	return m.ListarFn(familiaID)
}

func (m *MockMembroService) Atualizar(familiaID, id, nome, relacionamento string, ativo bool) (*domain.Membro, error) {
	return m.AtualizarFn(familiaID, id, nome, relacionamento, ativo)
}

func (m *MockMembroService) Inativar(familiaID, id string) error {
	return m.InativarFn(familiaID, id)
}

func setupMembroRouter(svc handler.MembroServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", "usuario-teste-uuid")
		c.Set("familiaID", "familia-teste-uuid")
		c.Next()
	})
	h := handler.NewMembroHandler(svc)
	v1 := r.Group("/api/v1")
	{
		v1.GET("/membros", h.Listar)
		v1.POST("/membros", h.Criar)
		v1.GET("/membros/:id", h.BuscarPorID)
		v1.PUT("/membros/:id", h.Atualizar)
		v1.PATCH("/membros/:id/inativar", h.Inativar)
	}
	return r
}

// ---- GET /membros ----

func TestListarMembrosHandler_Sucesso(t *testing.T) {
	membros := []*domain.Membro{
		{ID: "uuid-1", Nome: "Ana Lima", Relacionamento: "cônjuge", Ativo: true},
		{ID: "uuid-2", Nome: "Carlos Souza", Relacionamento: "filho", Ativo: true},
	}
	svc := &MockMembroService{
		ListarFn: func(familiaID string) ([]*domain.Membro, error) {
			return membros, nil
		},
	}
	r := setupMembroRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/membros", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "uuid-1", resp[0]["id"])
	assert.Equal(t, "Ana Lima", resp[0]["nome"])
}

// ---- POST /membros ----

func TestCriarMembroHandler_Sucesso(t *testing.T) {
	svc := &MockMembroService{
		CriarFn: func(familiaID, nome, relacionamento string) (*domain.Membro, error) {
			return &domain.Membro{
				ID:             "uuid-novo",
				Nome:           nome,
				Relacionamento: relacionamento,
				Ativo:          true,
			}, nil
		},
	}
	r := setupMembroRouter(svc)

	body, _ := json.Marshal(map[string]string{
		"nome":           "Beatriz Costa",
		"relacionamento": "filha",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/membros", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-novo", resp["id"])
	assert.Equal(t, "Beatriz Costa", resp["nome"])
	assert.Equal(t, "filha", resp["relacionamento"])
	assert.Equal(t, true, resp["ativo"])
}

func TestCriarMembroHandler_NomeVazio(t *testing.T) {
	svc := &MockMembroService{
		CriarFn: func(familiaID, nome, relacionamento string) (*domain.Membro, error) {
			return nil, domain.ErrNomeObrigatorio
		},
	}
	r := setupMembroRouter(svc)

	body, _ := json.Marshal(map[string]string{
		"nome":           "",
		"relacionamento": "cônjuge",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/membros", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.ErrNomeObrigatorio.Error(), resp["error"])
}

func TestCriarMembroHandler_BodyInvalido(t *testing.T) {
	svc := &MockMembroService{}
	r := setupMembroRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/membros", bytes.NewReader([]byte("nao-e-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ---- GET /membros/:id ----

func TestBuscarMembroHandler_Sucesso(t *testing.T) {
	svc := &MockMembroService{
		BuscarFn: func(familiaID, id string) (*domain.Membro, error) {
			return &domain.Membro{
				ID:             id,
				Nome:           "Pedro Alves",
				Relacionamento: "pai",
				Ativo:          true,
			}, nil
		},
	}
	r := setupMembroRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/membros/uuid-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resp["id"])
	assert.Equal(t, "Pedro Alves", resp["nome"])
}

func TestBuscarMembroHandler_NaoEncontrado(t *testing.T) {
	svc := &MockMembroService{
		BuscarFn: func(familiaID, id string) (*domain.Membro, error) {
			return nil, domain.ErrMembroNaoEncontrado
		},
	}
	r := setupMembroRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/membros/uuid-inexistente", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "membro não encontrado", resp["error"])
}

// ---- PUT /membros/:id ----

func TestAtualizarMembroHandler_Sucesso(t *testing.T) {
	svc := &MockMembroService{
		AtualizarFn: func(familiaID, id, nome, relacionamento string, ativo bool) (*domain.Membro, error) {
			return &domain.Membro{
				ID:             id,
				Nome:           nome,
				Relacionamento: relacionamento,
				Ativo:          ativo,
			}, nil
		},
	}
	r := setupMembroRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"nome":           "Ana Souza",
		"relacionamento": "esposa",
		"ativo":          true,
	})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/membros/uuid-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resp["id"])
	assert.Equal(t, "Ana Souza", resp["nome"])
	assert.Equal(t, "esposa", resp["relacionamento"])
	assert.Equal(t, true, resp["ativo"])
}

func TestAtualizarMembroHandler_NaoEncontrado(t *testing.T) {
	svc := &MockMembroService{
		AtualizarFn: func(familiaID, id, nome, relacionamento string, ativo bool) (*domain.Membro, error) {
			return nil, domain.ErrMembroNaoEncontrado
		},
	}
	r := setupMembroRouter(svc)

	body, _ := json.Marshal(map[string]any{
		"nome":           "Nome Qualquer",
		"relacionamento": "",
		"ativo":          true,
	})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/membros/uuid-inexistente", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "membro não encontrado", resp["error"])
}

func TestAtualizarMembroHandler_BodyInvalido(t *testing.T) {
	svc := &MockMembroService{}
	r := setupMembroRouter(svc)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/membros/uuid-1", bytes.NewReader([]byte("nao-e-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ---- PATCH /membros/:id/inativar ----

func TestInativarMembroHandler_Sucesso(t *testing.T) {
	svc := &MockMembroService{
		InativarFn: func(familiaID, id string) error {
			return nil
		},
	}
	r := setupMembroRouter(svc)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/membros/uuid-1/inativar", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "membro inativado com sucesso", resp["message"])
}

func TestInativarMembroHandler_NaoEncontrado(t *testing.T) {
	svc := &MockMembroService{
		InativarFn: func(familiaID, id string) error {
			return domain.ErrMembroNaoEncontrado
		},
	}
	r := setupMembroRouter(svc)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/membros/uuid-inexistente/inativar", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "membro não encontrado", resp["error"])
}
