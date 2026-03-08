package handler

import (
	"errors"
	"net/http"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// RendaFixaServiceInterface define os métodos do service usados pelo handler.
// Redeclarada aqui para desacoplar o pacote handler do service sem importação circular.
type RendaFixaServiceInterface interface {
	Criar(descricao, membroID string, valor float64, diaRecebimento int) (*domain.RendaFixa, error)
	BuscarPorID(id string) (*domain.RendaFixa, error)
	Listar() ([]*domain.RendaFixa, error)
	ListarAtivas() ([]*domain.RendaFixa, error)
	Atualizar(id, descricao, membroID string, valor float64, diaRecebimento int, ativa bool) (*domain.RendaFixa, error)
	Inativar(id string) error
}

// RendaFixaHandler contém os handlers HTTP para rendas fixas.
type RendaFixaHandler struct {
	svc RendaFixaServiceInterface
}

// NewRendaFixaHandler cria um novo RendaFixaHandler.
func NewRendaFixaHandler(svc RendaFixaServiceInterface) *RendaFixaHandler {
	return &RendaFixaHandler{svc: svc}
}

// rendaFixaResponse é a estrutura de resposta JSON para uma renda fixa.
type rendaFixaResponse struct {
	ID             string  `json:"id"`
	Descricao      string  `json:"descricao"`
	MembroID       string  `json:"membro_id"`
	Valor          float64 `json:"valor"`
	DiaRecebimento int     `json:"dia_recebimento"`
	Ativa          bool    `json:"ativa"`
}

func toRendaFixaResponse(r *domain.RendaFixa) rendaFixaResponse {
	return rendaFixaResponse{
		ID:             r.ID,
		Descricao:      r.Descricao,
		MembroID:       r.MembroID,
		Valor:          r.Valor,
		DiaRecebimento: r.DiaRecebimento,
		Ativa:          r.Ativa,
	}
}

type criarRendaFixaRequest struct {
	Descricao      string  `json:"descricao"`
	MembroID       string  `json:"membro_id"`
	Valor          float64 `json:"valor"`
	DiaRecebimento int     `json:"dia_recebimento"`
}

type atualizarRendaFixaRequest struct {
	Descricao      string  `json:"descricao"       binding:"required"`
	MembroID       string  `json:"membro_id"       binding:"required"`
	Valor          float64 `json:"valor"           binding:"required"`
	DiaRecebimento int     `json:"dia_recebimento" binding:"required"`
	Ativa          bool    `json:"ativa"`
}

func (h *RendaFixaHandler) erroDominio(c *gin.Context, err error) bool {
	switch {
	case errors.Is(err, domain.ErrDescricaoRendaObrigatoria),
		errors.Is(err, domain.ErrMembroIDRendaObrigatorio),
		errors.Is(err, domain.ErrValorRendaInvalido),
		errors.Is(err, domain.ErrDiaRecebimentoInvalido):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return true
	case errors.Is(err, domain.ErrRendaFixaNaoEncontrada):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return true
	}
	return false
}

// Listar retorna todas as rendas fixas.
// GET /api/v1/rendas-fixas
func (h *RendaFixaHandler) Listar(c *gin.Context) {
	rendas, err := h.svc.Listar()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	resp := make([]rendaFixaResponse, 0, len(rendas))
	for _, r := range rendas {
		resp = append(resp, toRendaFixaResponse(r))
	}

	c.JSON(http.StatusOK, resp)
}

// Criar cria uma nova renda fixa.
// POST /api/v1/rendas-fixas
func (h *RendaFixaHandler) Criar(c *gin.Context) {
	var req criarRendaFixaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	renda, err := h.svc.Criar(req.Descricao, req.MembroID, req.Valor, req.DiaRecebimento)
	if err != nil {
		if h.erroDominio(c, err) {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	c.JSON(http.StatusCreated, toRendaFixaResponse(renda))
}

// BuscarPorID retorna uma renda fixa pelo seu ID.
// GET /api/v1/rendas-fixas/:id
func (h *RendaFixaHandler) BuscarPorID(c *gin.Context) {
	id := c.Param("id")

	renda, err := h.svc.BuscarPorID(id)
	if err != nil {
		if h.erroDominio(c, err) {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	c.JSON(http.StatusOK, toRendaFixaResponse(renda))
}

// Atualizar atualiza os dados de uma renda fixa existente.
// PUT /api/v1/rendas-fixas/:id
func (h *RendaFixaHandler) Atualizar(c *gin.Context) {
	id := c.Param("id")

	var req atualizarRendaFixaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	renda, err := h.svc.Atualizar(id, req.Descricao, req.MembroID, req.Valor, req.DiaRecebimento, req.Ativa)
	if err != nil {
		if h.erroDominio(c, err) {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	c.JSON(http.StatusOK, toRendaFixaResponse(renda))
}

// Inativar marca uma renda fixa como inativa.
// PATCH /api/v1/rendas-fixas/:id/inativar
func (h *RendaFixaHandler) Inativar(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.Inativar(id); err != nil {
		if h.erroDominio(c, err) {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "renda fixa inativada com sucesso"})
}
