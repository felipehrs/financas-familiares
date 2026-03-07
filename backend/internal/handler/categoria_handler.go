package handler

import (
	"errors"
	"net/http"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// CategoriaServiceInterface define os métodos do service usados pelo handler.
// Redeclarada aqui para desacoplar o pacote handler do service sem importação circular.
type CategoriaServiceInterface interface {
	Criar(nome string) (*domain.Categoria, error)
	Listar() ([]*domain.Categoria, error)
	Atualizar(id, nome string) (*domain.Categoria, error)
	Excluir(id string) error
}

// CategoriaHandler contém os handlers HTTP para categorias.
type CategoriaHandler struct {
	svc CategoriaServiceInterface
}

// NewCategoriaHandler cria um novo CategoriaHandler.
func NewCategoriaHandler(svc CategoriaServiceInterface) *CategoriaHandler {
	return &CategoriaHandler{svc: svc}
}

// categoriaResponse é a estrutura de resposta JSON para uma categoria.
type categoriaResponse struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
}

func toCategoriaResponse(c *domain.Categoria) categoriaResponse {
	return categoriaResponse{
		ID:   c.ID,
		Nome: c.Nome,
	}
}

type criarCategoriaRequest struct {
	Nome string `json:"nome"`
}

type atualizarCategoriaRequest struct {
	Nome string `json:"nome"`
}

// Listar retorna todas as categorias.
// GET /api/v1/categorias
func (h *CategoriaHandler) Listar(c *gin.Context) {
	categorias, err := h.svc.Listar()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	resp := make([]categoriaResponse, 0, len(categorias))
	for _, cat := range categorias {
		resp = append(resp, toCategoriaResponse(cat))
	}

	c.JSON(http.StatusOK, resp)
}

// Criar cria uma nova categoria.
// POST /api/v1/categorias
func (h *CategoriaHandler) Criar(c *gin.Context) {
	var req criarCategoriaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	cat, err := h.svc.Criar(req.Nome)
	if err != nil {
		if errors.Is(err, domain.ErrNomeCategoriaObrigatorio) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	c.JSON(http.StatusCreated, toCategoriaResponse(cat))
}

// Atualizar atualiza o nome de uma categoria existente.
// PUT /api/v1/categorias/:id
func (h *CategoriaHandler) Atualizar(c *gin.Context) {
	id := c.Param("id")

	var req atualizarCategoriaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	cat, err := h.svc.Atualizar(id, req.Nome)
	if err != nil {
		if errors.Is(err, domain.ErrCategoriaNaoEncontrada) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domain.ErrNomeCategoriaObrigatorio) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	c.JSON(http.StatusOK, toCategoriaResponse(cat))
}

// Excluir remove uma categoria (soft delete).
// DELETE /api/v1/categorias/:id
// Retorna 409 Conflict se a categoria tiver registros vinculados.
func (h *CategoriaHandler) Excluir(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.Excluir(id); err != nil {
		if errors.Is(err, domain.ErrCategoriaNaoEncontrada) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domain.ErrCategoriaComVinculos) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	c.Status(http.StatusNoContent)
}
