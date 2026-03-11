package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// RendimentoInvestimentoServiceInterface define os métodos do service usados pelo handler.
// Redeclarada aqui para desacoplar o pacote handler do service sem importação circular.
type RendimentoInvestimentoServiceInterface interface {
	Criar(familiaID, descricao, membroID string, data time.Time, valor, valorDistribuido float64) (*domain.RendimentoInvestimento, error)
	BuscarPorID(familiaID, id string) (*domain.RendimentoInvestimento, error)
	Listar(familiaID string) ([]*domain.RendimentoInvestimento, error)
	ListarPorMes(familiaID string, mes, ano int) ([]*domain.RendimentoInvestimento, error)
	Atualizar(familiaID, id, descricao, membroID string, data time.Time, valor, valorDistribuido float64) (*domain.RendimentoInvestimento, error)
	Excluir(familiaID, id string) error
}

// RendimentoInvestimentoHandler contém os handlers HTTP para rendimentos de investimento.
type RendimentoInvestimentoHandler struct {
	svc RendimentoInvestimentoServiceInterface
}

// NewRendimentoInvestimentoHandler cria um novo RendimentoInvestimentoHandler.
func NewRendimentoInvestimentoHandler(svc RendimentoInvestimentoServiceInterface) *RendimentoInvestimentoHandler {
	return &RendimentoInvestimentoHandler{svc: svc}
}

// rendimentoInvestimentoResponse é a estrutura de resposta JSON para um rendimento de investimento.
type rendimentoInvestimentoResponse struct {
	ID               string  `json:"id"`
	Descricao        string  `json:"descricao"`
	MembroID         string  `json:"membro_id"`
	Data             string  `json:"data"`
	Valor            float64 `json:"valor"`
	ValorDistribuido float64 `json:"valor_distribuido"`
}

func toRendimentoInvestimentoResponse(r *domain.RendimentoInvestimento) rendimentoInvestimentoResponse {
	return rendimentoInvestimentoResponse{
		ID:               r.ID,
		Descricao:        r.Descricao,
		MembroID:         r.MembroID,
		Data:             r.Data.Format("2006-01-02"),
		Valor:            r.Valor,
		ValorDistribuido: r.ValorDistribuido,
	}
}

type criarRendimentoInvestimentoRequest struct {
	Descricao        string   `json:"descricao"`
	MembroID         string   `json:"membro_id"`
	Data             string   `json:"data"`
	Valor            float64  `json:"valor"`
	ValorDistribuido *float64 `json:"valor_distribuido"`
}

type atualizarRendimentoInvestimentoRequest struct {
	Descricao        string   `json:"descricao"`
	MembroID         string   `json:"membro_id"`
	Data             string   `json:"data"`
	Valor            float64  `json:"valor"`
	ValorDistribuido *float64 `json:"valor_distribuido"`
}

func rendimentoInvestimentoErroParaHTTP(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrRendimentoInvestimentoNaoEncontrado):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrDescricaoRendimentoObrigatoria),
		errors.Is(err, domain.ErrMembroIDRendimentoObrigatorio),
		errors.Is(err, domain.ErrValorRendimentoInvalido),
		errors.Is(err, domain.ErrDataRendimentoObrigatoria),
		errors.Is(err, domain.ErrValorDistribuidoInvalido):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
	}
}

// Listar retorna todos os rendimentos de investimento, com filtro opcional por mes e ano.
// GET /api/v1/rendimentos-investimento
// GET /api/v1/rendimentos-investimento?mes=3&ano=2026
func (h *RendimentoInvestimentoHandler) Listar(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

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

		rendimentos, err := h.svc.ListarPorMes(familiaID, mes, ano)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
			return
		}

		resp := make([]rendimentoInvestimentoResponse, 0, len(rendimentos))
		for _, r := range rendimentos {
			resp = append(resp, toRendimentoInvestimentoResponse(r))
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	rendimentos, err := h.svc.Listar(familiaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	resp := make([]rendimentoInvestimentoResponse, 0, len(rendimentos))
	for _, r := range rendimentos {
		resp = append(resp, toRendimentoInvestimentoResponse(r))
	}
	c.JSON(http.StatusOK, resp)
}

// Criar cria um novo rendimento de investimento.
// POST /api/v1/rendimentos-investimento
func (h *RendimentoInvestimentoHandler) Criar(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	var req criarRendimentoInvestimentoRequest
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

	valorDistribuido := 0.0
	if req.ValorDistribuido != nil {
		valorDistribuido = *req.ValorDistribuido
	}

	rendimento, err := h.svc.Criar(familiaID, req.Descricao, req.MembroID, data, req.Valor, valorDistribuido)
	if err != nil {
		rendimentoInvestimentoErroParaHTTP(c, err)
		return
	}

	c.JSON(http.StatusCreated, toRendimentoInvestimentoResponse(rendimento))
}

// BuscarPorID retorna um rendimento de investimento pelo seu ID.
// GET /api/v1/rendimentos-investimento/:id
func (h *RendimentoInvestimentoHandler) BuscarPorID(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	rendimento, err := h.svc.BuscarPorID(familiaID, id)
	if err != nil {
		rendimentoInvestimentoErroParaHTTP(c, err)
		return
	}

	c.JSON(http.StatusOK, toRendimentoInvestimentoResponse(rendimento))
}

// Atualizar atualiza os dados de um rendimento de investimento existente.
// PUT /api/v1/rendimentos-investimento/:id
func (h *RendimentoInvestimentoHandler) Atualizar(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	var req atualizarRendimentoInvestimentoRequest
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

	valorDistribuido := 0.0
	if req.ValorDistribuido != nil {
		valorDistribuido = *req.ValorDistribuido
	}

	rendimento, err := h.svc.Atualizar(familiaID, id, req.Descricao, req.MembroID, data, req.Valor, valorDistribuido)
	if err != nil {
		rendimentoInvestimentoErroParaHTTP(c, err)
		return
	}

	c.JSON(http.StatusOK, toRendimentoInvestimentoResponse(rendimento))
}

// Excluir realiza o soft-delete de um rendimento de investimento.
// DELETE /api/v1/rendimentos-investimento/:id
func (h *RendimentoInvestimentoHandler) Excluir(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	if err := h.svc.Excluir(familiaID, id); err != nil {
		rendimentoInvestimentoErroParaHTTP(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
