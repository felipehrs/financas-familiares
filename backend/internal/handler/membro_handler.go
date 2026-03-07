package handler

import (
	"errors"
	"net/http"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// MembroServiceInterface define os métodos do service usados pelo handler.
// Redeclarada aqui para desacoplar o pacote handler do service sem importação circular.
type MembroServiceInterface interface {
	Criar(nome, relacionamento string) (*domain.Membro, error)
	BuscarPorID(id string) (*domain.Membro, error)
	Listar() ([]*domain.Membro, error)
	Atualizar(id, nome, relacionamento string, ativo bool) (*domain.Membro, error)
	Inativar(id string) error
}

// MembroHandler contém os handlers HTTP para membros da família.
type MembroHandler struct {
	svc MembroServiceInterface
}

// NewMembroHandler cria um novo MembroHandler.
func NewMembroHandler(svc MembroServiceInterface) *MembroHandler {
	return &MembroHandler{svc: svc}
}

// membroResponse é a estrutura de resposta JSON para um membro.
type membroResponse struct {
	ID             string `json:"id"`
	Nome           string `json:"nome"`
	Relacionamento string `json:"relacionamento"`
	Ativo          bool   `json:"ativo"`
}

func toMembroResponse(m *domain.Membro) membroResponse {
	return membroResponse{
		ID:             m.ID,
		Nome:           m.Nome,
		Relacionamento: m.Relacionamento,
		Ativo:          m.Ativo,
	}
}

type criarMembroRequest struct {
	Nome           string `json:"nome"`
	Relacionamento string `json:"relacionamento"`
}

type atualizarMembroRequest struct {
	Nome           string `json:"nome"         binding:"required"`
	Relacionamento string `json:"relacionamento"`
	Ativo          bool   `json:"ativo"`
}

// Listar retorna todos os membros da família.
// GET /api/v1/membros
func (h *MembroHandler) Listar(c *gin.Context) {
	membros, err := h.svc.Listar()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	resp := make([]membroResponse, 0, len(membros))
	for _, m := range membros {
		resp = append(resp, toMembroResponse(m))
	}

	c.JSON(http.StatusOK, resp)
}

// Criar cria um novo membro da família.
// POST /api/v1/membros
func (h *MembroHandler) Criar(c *gin.Context) {
	var req criarMembroRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	membro, err := h.svc.Criar(req.Nome, req.Relacionamento)
	if err != nil {
		if errors.Is(err, domain.ErrNomeObrigatorio) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	c.JSON(http.StatusCreated, toMembroResponse(membro))
}

// BuscarPorID retorna um membro pelo seu ID.
// GET /api/v1/membros/:id
func (h *MembroHandler) BuscarPorID(c *gin.Context) {
	id := c.Param("id")

	membro, err := h.svc.BuscarPorID(id)
	if err != nil {
		if errors.Is(err, domain.ErrMembroNaoEncontrado) {
			c.JSON(http.StatusNotFound, gin.H{"error": "membro não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	c.JSON(http.StatusOK, toMembroResponse(membro))
}

// Atualizar atualiza os dados de um membro existente.
// PUT /api/v1/membros/:id
func (h *MembroHandler) Atualizar(c *gin.Context) {
	id := c.Param("id")

	var req atualizarMembroRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	membro, err := h.svc.Atualizar(id, req.Nome, req.Relacionamento, req.Ativo)
	if err != nil {
		if errors.Is(err, domain.ErrMembroNaoEncontrado) {
			c.JSON(http.StatusNotFound, gin.H{"error": "membro não encontrado"})
			return
		}
		if errors.Is(err, domain.ErrNomeObrigatorio) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	c.JSON(http.StatusOK, toMembroResponse(membro))
}

// Inativar marca um membro como inativo.
// PATCH /api/v1/membros/:id/inativar
func (h *MembroHandler) Inativar(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.Inativar(id); err != nil {
		if errors.Is(err, domain.ErrMembroNaoEncontrado) {
			c.JSON(http.StatusNotFound, gin.H{"error": "membro não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "membro inativado com sucesso"})
}
