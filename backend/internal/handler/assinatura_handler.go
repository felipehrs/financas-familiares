package handler

import (
	"errors"
	"net/http"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// AssinaturaServiceInterface define os métodos do service usados pelo handler.
// Redeclarada aqui para desacoplar o pacote handler do service sem importação circular.
type AssinaturaServiceInterface interface {
	Criar(familiaID, nome, membroID string, categoriaID *string, valor float64, diaCobranca int, formaPagamento string) (*domain.Assinatura, error)
	BuscarPorID(familiaID, id string) (*domain.Assinatura, error)
	Listar(familiaID string) ([]*domain.Assinatura, error)
	Atualizar(familiaID, id, nome, membroID string, categoriaID *string, valor float64, diaCobranca int, formaPagamento, status string) (*domain.Assinatura, error)
	AlterarStatus(familiaID, id, status string) (*domain.Assinatura, error)
}

// AssinaturaHandler contém os handlers HTTP para assinaturas recorrentes.
type AssinaturaHandler struct {
	svc AssinaturaServiceInterface
}

// NewAssinaturaHandler cria um novo AssinaturaHandler.
func NewAssinaturaHandler(svc AssinaturaServiceInterface) *AssinaturaHandler {
	return &AssinaturaHandler{svc: svc}
}

// assinaturaResponse é a estrutura de resposta JSON para uma assinatura.
type assinaturaResponse struct {
	ID             string  `json:"id"`
	Nome           string  `json:"nome"`
	MembroID       string  `json:"membro_id"`
	CategoriaID    *string `json:"categoria_id"`
	Valor          float64 `json:"valor"`
	DiaCobranca    int     `json:"dia_cobranca"`
	FormaPagamento string  `json:"forma_pagamento"`
	Status         string  `json:"status"`
}

func toAssinaturaResponse(a *domain.Assinatura) assinaturaResponse {
	return assinaturaResponse{
		ID:             a.ID,
		Nome:           a.Nome,
		MembroID:       a.MembroID,
		CategoriaID:    a.CategoriaID,
		Valor:          a.Valor,
		DiaCobranca:    a.DiaCobranca,
		FormaPagamento: a.FormaPagamento,
		Status:         a.Status,
	}
}

type criarAssinaturaRequest struct {
	Nome           string  `json:"nome"`
	MembroID       string  `json:"membro_id"`
	CategoriaID    *string `json:"categoria_id"`
	Valor          float64 `json:"valor"`
	DiaCobranca    int     `json:"dia_cobranca"`
	FormaPagamento string  `json:"forma_pagamento"`
}

type atualizarAssinaturaRequest struct {
	Nome           string  `json:"nome"`
	MembroID       string  `json:"membro_id"`
	CategoriaID    *string `json:"categoria_id"`
	Valor          float64 `json:"valor"`
	DiaCobranca    int     `json:"dia_cobranca"`
	FormaPagamento string  `json:"forma_pagamento"`
	Status         string  `json:"status"`
}

type alterarStatusRequest struct {
	Status string `json:"status"`
}

func assinaturaErroParaHTTP(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrAssinaturaNaoEncontrada):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrNomeAssinaturaObrigatorio),
		errors.Is(err, domain.ErrMembroIDAssinaturaObrigatorio),
		errors.Is(err, domain.ErrValorAssinaturaInvalido),
		errors.Is(err, domain.ErrDiaCobrancaInvalido),
		errors.Is(err, domain.ErrFormaPagamentoObrigatoria),
		errors.Is(err, domain.ErrStatusInvalido):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
	}
}

// Listar retorna todas as assinaturas recorrentes.
// GET /api/v1/assinaturas
func (h *AssinaturaHandler) Listar(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	assinaturas, err := h.svc.Listar(familiaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	resp := make([]assinaturaResponse, 0, len(assinaturas))
	for _, a := range assinaturas {
		resp = append(resp, toAssinaturaResponse(a))
	}

	c.JSON(http.StatusOK, resp)
}

// Criar cria uma nova assinatura recorrente.
// POST /api/v1/assinaturas
func (h *AssinaturaHandler) Criar(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	var req criarAssinaturaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	assinatura, err := h.svc.Criar(familiaID, req.Nome, req.MembroID, req.CategoriaID, req.Valor, req.DiaCobranca, req.FormaPagamento)
	if err != nil {
		assinaturaErroParaHTTP(c, err)
		return
	}

	c.JSON(http.StatusCreated, toAssinaturaResponse(assinatura))
}

// BuscarPorID retorna uma assinatura pelo seu ID.
// GET /api/v1/assinaturas/:id
func (h *AssinaturaHandler) BuscarPorID(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	assinatura, err := h.svc.BuscarPorID(familiaID, id)
	if err != nil {
		assinaturaErroParaHTTP(c, err)
		return
	}

	c.JSON(http.StatusOK, toAssinaturaResponse(assinatura))
}

// Atualizar atualiza os dados de uma assinatura existente.
// PUT /api/v1/assinaturas/:id
func (h *AssinaturaHandler) Atualizar(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	var req atualizarAssinaturaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	assinatura, err := h.svc.Atualizar(familiaID, id, req.Nome, req.MembroID, req.CategoriaID, req.Valor, req.DiaCobranca, req.FormaPagamento, req.Status)
	if err != nil {
		assinaturaErroParaHTTP(c, err)
		return
	}

	c.JSON(http.StatusOK, toAssinaturaResponse(assinatura))
}

// AlterarStatus altera o status de uma assinatura existente.
// PATCH /api/v1/assinaturas/:id/status
func (h *AssinaturaHandler) AlterarStatus(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	var req alterarStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	assinatura, err := h.svc.AlterarStatus(familiaID, id, req.Status)
	if err != nil {
		assinaturaErroParaHTTP(c, err)
		return
	}

	c.JSON(http.StatusOK, toAssinaturaResponse(assinatura))
}
