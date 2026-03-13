package handler

import (
	"errors"
	"net/http"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// ContaFixaServiceInterface define os métodos do service usados pelo handler.
// Redeclarada aqui para desacoplar o pacote handler do service sem importação circular.
type ContaFixaServiceInterface interface {
	Criar(familiaID, descricao, membroID string, categoriaID *string, valor float64, diaVencimento int, formaPagamento string) (*domain.ContaFixa, error)
	BuscarPorID(familiaID, id string) (*domain.ContaFixa, error)
	Listar(familiaID string) ([]*domain.ContaFixa, error)
	Atualizar(familiaID, id, descricao, membroID string, categoriaID *string, valor float64, diaVencimento int, formaPagamento string, ativa bool) (*domain.ContaFixa, error)
	AlterarAtivo(familiaID, id string, ativa bool) (*domain.ContaFixa, error)
}

// ContaFixaHandler contém os handlers HTTP para contas fixas mensais.
type ContaFixaHandler struct {
	svc ContaFixaServiceInterface
}

// NewContaFixaHandler cria um novo ContaFixaHandler.
func NewContaFixaHandler(svc ContaFixaServiceInterface) *ContaFixaHandler {
	return &ContaFixaHandler{svc: svc}
}

// contaFixaResponse é a estrutura de resposta JSON para uma conta fixa.
type contaFixaResponse struct {
	ID             string  `json:"id"`
	Descricao      string  `json:"descricao"`
	MembroID       string  `json:"membro_id"`
	CategoriaID    *string `json:"categoria_id"`
	Valor          float64 `json:"valor"`
	DiaVencimento  int     `json:"dia_vencimento"`
	FormaPagamento string  `json:"forma_pagamento"`
	Ativa          bool    `json:"ativa"`
}

func toContaFixaResponse(c *domain.ContaFixa) contaFixaResponse {
	return contaFixaResponse{
		ID:             c.ID,
		Descricao:      c.Descricao,
		MembroID:       c.MembroID,
		CategoriaID:    c.CategoriaID,
		Valor:          c.Valor,
		DiaVencimento:  c.DiaVencimento,
		FormaPagamento: c.FormaPagamento,
		Ativa:          c.Ativa,
	}
}

type criarContaFixaRequest struct {
	Descricao      string  `json:"descricao"`
	MembroID       string  `json:"membro_id"`
	CategoriaID    *string `json:"categoria_id"`
	Valor          float64 `json:"valor"`
	DiaVencimento  int     `json:"dia_vencimento"`
	FormaPagamento string  `json:"forma_pagamento"`
}

type atualizarContaFixaRequest struct {
	Descricao      string  `json:"descricao"`
	MembroID       string  `json:"membro_id"`
	CategoriaID    *string `json:"categoria_id"`
	Valor          float64 `json:"valor"`
	DiaVencimento  int     `json:"dia_vencimento"`
	FormaPagamento string  `json:"forma_pagamento"`
	Ativa          bool    `json:"ativa"`
}

type alterarAtivoRequest struct {
	Ativa bool `json:"ativa"`
}

func contaFixaErroParaHTTP(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrContaFixaNaoEncontrada):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrDescricaoContaFixaObrigatoria),
		errors.Is(err, domain.ErrMembroIDContaFixaObrigatorio),
		errors.Is(err, domain.ErrValorContaFixaInvalido),
		errors.Is(err, domain.ErrDiaVencimentoContaFixaInvalido),
		errors.Is(err, domain.ErrFormaPagamentoContaFixaObrigatoria):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
	}
}

// Listar retorna todas as contas fixas mensais.
// GET /api/v1/contas-fixas
func (h *ContaFixaHandler) Listar(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	contas, err := h.svc.Listar(familiaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	resp := make([]contaFixaResponse, 0, len(contas))
	for _, conta := range contas {
		resp = append(resp, toContaFixaResponse(conta))
	}

	c.JSON(http.StatusOK, resp)
}

// Criar cria uma nova conta fixa mensal.
// POST /api/v1/contas-fixas
func (h *ContaFixaHandler) Criar(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	var req criarContaFixaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	conta, err := h.svc.Criar(familiaID, req.Descricao, req.MembroID, req.CategoriaID, req.Valor, req.DiaVencimento, req.FormaPagamento)
	if err != nil {
		contaFixaErroParaHTTP(c, err)
		return
	}

	c.JSON(http.StatusCreated, toContaFixaResponse(conta))
}

// BuscarPorID retorna uma conta fixa pelo seu ID.
// GET /api/v1/contas-fixas/:id
func (h *ContaFixaHandler) BuscarPorID(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	conta, err := h.svc.BuscarPorID(familiaID, id)
	if err != nil {
		contaFixaErroParaHTTP(c, err)
		return
	}

	c.JSON(http.StatusOK, toContaFixaResponse(conta))
}

// Atualizar atualiza os dados de uma conta fixa existente.
// PUT /api/v1/contas-fixas/:id
func (h *ContaFixaHandler) Atualizar(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	var req atualizarContaFixaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	conta, err := h.svc.Atualizar(familiaID, id, req.Descricao, req.MembroID, req.CategoriaID, req.Valor, req.DiaVencimento, req.FormaPagamento, req.Ativa)
	if err != nil {
		contaFixaErroParaHTTP(c, err)
		return
	}

	c.JSON(http.StatusOK, toContaFixaResponse(conta))
}

// AlterarAtivo altera o estado ativo/inativo de uma conta fixa existente.
// PATCH /api/v1/contas-fixas/:id/ativo
func (h *ContaFixaHandler) AlterarAtivo(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	var req alterarAtivoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	conta, err := h.svc.AlterarAtivo(familiaID, id, req.Ativa)
	if err != nil {
		contaFixaErroParaHTTP(c, err)
		return
	}

	c.JSON(http.StatusOK, toContaFixaResponse(conta))
}
