package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// DespesaGeralServiceInterface define os métodos do service usados pelo handler.
// Redeclarada aqui para desacoplar o pacote handler do service sem importação circular.
type DespesaGeralServiceInterface interface {
	Criar(membroID string, categoriaID *string, descricao string, data time.Time, valor float64, formaPagamento string, observacoes *string) (*domain.DespesaGeral, error)
	BuscarPorID(id string) (*domain.DespesaGeral, error)
	Listar() ([]*domain.DespesaGeral, error)
	ListarPorMes(mes, ano int) ([]*domain.DespesaGeral, error)
	Atualizar(id, membroID string, categoriaID *string, descricao string, data time.Time, valor float64, formaPagamento string, observacoes *string) (*domain.DespesaGeral, error)
	Excluir(id string) error
}

// DespesaGeralHandler contém os handlers HTTP para despesas gerais.
type DespesaGeralHandler struct {
	svc DespesaGeralServiceInterface
}

// NewDespesaGeralHandler cria um novo DespesaGeralHandler.
func NewDespesaGeralHandler(svc DespesaGeralServiceInterface) *DespesaGeralHandler {
	return &DespesaGeralHandler{svc: svc}
}

// despesaGeralResponse é a estrutura de resposta JSON para uma despesa geral.
type despesaGeralResponse struct {
	ID             string  `json:"id"`
	MembroID       string  `json:"membro_id"`
	CategoriaID    *string `json:"categoria_id"`
	Descricao      string  `json:"descricao"`
	Data           string  `json:"data"`
	Valor          float64 `json:"valor"`
	FormaPagamento string  `json:"forma_pagamento"`
	Observacoes    *string `json:"observacoes"`
}

func toDespesaGeralResponse(d *domain.DespesaGeral) despesaGeralResponse {
	return despesaGeralResponse{
		ID:             d.ID,
		MembroID:       d.MembroID,
		CategoriaID:    d.CategoriaID,
		Descricao:      d.Descricao,
		Data:           d.Data.Format("2006-01-02"),
		Valor:          d.Valor,
		FormaPagamento: d.FormaPagamento,
		Observacoes:    d.Observacoes,
	}
}

type criarDespesaGeralRequest struct {
	MembroID       string  `json:"membro_id"`
	CategoriaID    *string `json:"categoria_id"`
	Descricao      string  `json:"descricao"`
	Data           string  `json:"data"`
	Valor          float64 `json:"valor"`
	FormaPagamento string  `json:"forma_pagamento"`
	Observacoes    *string `json:"observacoes"`
}

type atualizarDespesaGeralRequest struct {
	MembroID       string  `json:"membro_id"`
	CategoriaID    *string `json:"categoria_id"`
	Descricao      string  `json:"descricao"`
	Data           string  `json:"data"`
	Valor          float64 `json:"valor"`
	FormaPagamento string  `json:"forma_pagamento"`
	Observacoes    *string `json:"observacoes"`
}

func despesaGeralErroParaHTTP(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrDespesaGeralNaoEncontrada):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrDescricaoDespesaGeralObrigatoria),
		errors.Is(err, domain.ErrMembroIDDespesaGeralObrigatorio),
		errors.Is(err, domain.ErrValorDespesaGeralInvalido),
		errors.Is(err, domain.ErrDataDespesaGeralObrigatoria),
		errors.Is(err, domain.ErrFormaPagamentoDespesaGeralObrigatoria):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
	}
}

// Listar retorna todas as despesas gerais, com filtro opcional por mes e ano.
// GET /api/v1/despesas-gerais
// GET /api/v1/despesas-gerais?mes=3&ano=2026
func (h *DespesaGeralHandler) Listar(c *gin.Context) {
	mesStr := c.Query("mes")
	anoStr := c.Query("ano")

	if mesStr != "" && anoStr != "" {
		mes, err := strconv.Atoi(mesStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetro mes inválido"})
			return
		}
		ano, err := strconv.Atoi(anoStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetro ano inválido"})
			return
		}

		despesas, err := h.svc.ListarPorMes(mes, ano)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
			return
		}

		resp := make([]despesaGeralResponse, 0, len(despesas))
		for _, d := range despesas {
			resp = append(resp, toDespesaGeralResponse(d))
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	despesas, err := h.svc.Listar()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	resp := make([]despesaGeralResponse, 0, len(despesas))
	for _, d := range despesas {
		resp = append(resp, toDespesaGeralResponse(d))
	}
	c.JSON(http.StatusOK, resp)
}

// Criar cria uma nova despesa geral.
// POST /api/v1/despesas-gerais
func (h *DespesaGeralHandler) Criar(c *gin.Context) {
	var req criarDespesaGeralRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	var data time.Time
	if req.Data != "" {
		parsed, err := time.Parse("2006-01-02", req.Data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "formato de data inválido, use YYYY-MM-DD"})
			return
		}
		data = parsed
	}

	despesa, err := h.svc.Criar(req.MembroID, req.CategoriaID, req.Descricao, data, req.Valor, req.FormaPagamento, req.Observacoes)
	if err != nil {
		despesaGeralErroParaHTTP(c, err)
		return
	}

	c.JSON(http.StatusCreated, toDespesaGeralResponse(despesa))
}

// BuscarPorID retorna uma despesa geral pelo seu ID.
// GET /api/v1/despesas-gerais/:id
func (h *DespesaGeralHandler) BuscarPorID(c *gin.Context) {
	id := c.Param("id")

	despesa, err := h.svc.BuscarPorID(id)
	if err != nil {
		despesaGeralErroParaHTTP(c, err)
		return
	}

	c.JSON(http.StatusOK, toDespesaGeralResponse(despesa))
}

// Atualizar atualiza os dados de uma despesa geral existente.
// PUT /api/v1/despesas-gerais/:id
func (h *DespesaGeralHandler) Atualizar(c *gin.Context) {
	id := c.Param("id")

	var req atualizarDespesaGeralRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	var data time.Time
	if req.Data != "" {
		parsed, err := time.Parse("2006-01-02", req.Data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "formato de data inválido, use YYYY-MM-DD"})
			return
		}
		data = parsed
	}

	despesa, err := h.svc.Atualizar(id, req.MembroID, req.CategoriaID, req.Descricao, data, req.Valor, req.FormaPagamento, req.Observacoes)
	if err != nil {
		despesaGeralErroParaHTTP(c, err)
		return
	}

	c.JSON(http.StatusOK, toDespesaGeralResponse(despesa))
}

// Excluir realiza o soft-delete de uma despesa geral.
// DELETE /api/v1/despesas-gerais/:id
func (h *DespesaGeralHandler) Excluir(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.Excluir(id); err != nil {
		despesaGeralErroParaHTTP(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
