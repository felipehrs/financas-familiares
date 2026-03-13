package handler

import (
	"errors"
	"net/http"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// CartaoCreditoServiceInterface define os métodos do service usados pelo handler.
// Redeclarada aqui para desacoplar o pacote handler do service sem importação circular.
type CartaoCreditoServiceInterface interface {
	Criar(familiaID, nome, membroID string, diaFechamento, diaVencimento int, limite *float64) (*domain.CartaoCredito, error)
	BuscarPorID(familiaID, id string) (*domain.CartaoCredito, error)
	Listar(familiaID string) ([]*domain.CartaoCredito, error)
	Atualizar(familiaID, id, nome, membroID string, diaFechamento, diaVencimento int, limite *float64, ativo bool) (*domain.CartaoCredito, error)
	Inativar(familiaID, id string) error
}

// CartaoCreditoHandler contém os handlers HTTP para cartões de crédito.
type CartaoCreditoHandler struct {
	svc CartaoCreditoServiceInterface
}

// NewCartaoCreditoHandler cria um novo CartaoCreditoHandler.
func NewCartaoCreditoHandler(svc CartaoCreditoServiceInterface) *CartaoCreditoHandler {
	return &CartaoCreditoHandler{svc: svc}
}

// cartaoCreditoResponse é a estrutura de resposta JSON para um cartão de crédito.
type cartaoCreditoResponse struct {
	ID            string   `json:"id"`
	Nome          string   `json:"nome"`
	MembroID      string   `json:"membro_id"`
	DiaFechamento int      `json:"dia_fechamento"`
	DiaVencimento int      `json:"dia_vencimento"`
	Limite        *float64 `json:"limite"`
	Ativo         bool     `json:"ativo"`
}

func toCartaoCreditoResponse(c *domain.CartaoCredito) cartaoCreditoResponse {
	return cartaoCreditoResponse{
		ID:            c.ID,
		Nome:          c.Nome,
		MembroID:      c.MembroID,
		DiaFechamento: c.DiaFechamento,
		DiaVencimento: c.DiaVencimento,
		Limite:        c.Limite,
		Ativo:         c.Ativo,
	}
}

type criarCartaoCreditoRequest struct {
	Nome          string   `json:"nome"`
	MembroID      string   `json:"membro_id"`
	DiaFechamento int      `json:"dia_fechamento"`
	DiaVencimento int      `json:"dia_vencimento"`
	Limite        *float64 `json:"limite"`
}

type atualizarCartaoCreditoRequest struct {
	Nome          string   `json:"nome"           binding:"required"`
	MembroID      string   `json:"membro_id"      binding:"required"`
	DiaFechamento int      `json:"dia_fechamento" binding:"required"`
	DiaVencimento int      `json:"dia_vencimento" binding:"required"`
	Limite        *float64 `json:"limite"`
	Ativo         bool     `json:"ativo"`
}

func (h *CartaoCreditoHandler) erroDominio(c *gin.Context, err error) bool {
	switch {
	case errors.Is(err, domain.ErrNomeCartaoObrigatorio),
		errors.Is(err, domain.ErrMembroIDObrigatorio),
		errors.Is(err, domain.ErrDiaFechamentoInvalido),
		errors.Is(err, domain.ErrDiaVencimentoInvalido):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return true
	case errors.Is(err, domain.ErrCartaoNaoEncontrado):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return true
	}
	return false
}

// Listar retorna todos os cartões de crédito.
// GET /api/v1/cartoes
func (h *CartaoCreditoHandler) Listar(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	cartoes, err := h.svc.Listar(familiaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	resp := make([]cartaoCreditoResponse, 0, len(cartoes))
	for _, cartao := range cartoes {
		resp = append(resp, toCartaoCreditoResponse(cartao))
	}

	c.JSON(http.StatusOK, resp)
}

// Criar cria um novo cartão de crédito.
// POST /api/v1/cartoes
func (h *CartaoCreditoHandler) Criar(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	var req criarCartaoCreditoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	cartao, err := h.svc.Criar(familiaID, req.Nome, req.MembroID, req.DiaFechamento, req.DiaVencimento, req.Limite)
	if err != nil {
		if h.erroDominio(c, err) {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	c.JSON(http.StatusCreated, toCartaoCreditoResponse(cartao))
}

// BuscarPorID retorna um cartão pelo seu ID.
// GET /api/v1/cartoes/:id
func (h *CartaoCreditoHandler) BuscarPorID(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	cartao, err := h.svc.BuscarPorID(familiaID, id)
	if err != nil {
		if h.erroDominio(c, err) {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	c.JSON(http.StatusOK, toCartaoCreditoResponse(cartao))
}

// Atualizar atualiza os dados de um cartão existente.
// PUT /api/v1/cartoes/:id
func (h *CartaoCreditoHandler) Atualizar(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	var req atualizarCartaoCreditoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	cartao, err := h.svc.Atualizar(familiaID, id, req.Nome, req.MembroID, req.DiaFechamento, req.DiaVencimento, req.Limite, req.Ativo)
	if err != nil {
		if h.erroDominio(c, err) {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	c.JSON(http.StatusOK, toCartaoCreditoResponse(cartao))
}

// Inativar marca um cartão como inativo.
// PATCH /api/v1/cartoes/:id/inativar
func (h *CartaoCreditoHandler) Inativar(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	if err := h.svc.Inativar(familiaID, id); err != nil {
		if h.erroDominio(c, err) {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cartão inativado com sucesso"})
}
