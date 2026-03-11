package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// RendaExtraServiceInterface define os métodos do service usados pelo handler.
// Redeclarada aqui para desacoplar o pacote handler do service sem importação circular.
type RendaExtraServiceInterface interface {
	Criar(familiaID, descricao, membroID string, dataRecebimento time.Time, valor float64) (*domain.RendaExtra, error)
	BuscarPorID(familiaID, id string) (*domain.RendaExtra, error)
	Listar(familiaID string) ([]*domain.RendaExtra, error)
	ListarPorMes(familiaID string, mes, ano int) ([]*domain.RendaExtra, error)
	Atualizar(familiaID, id, descricao, membroID string, dataRecebimento time.Time, valor float64) (*domain.RendaExtra, error)
	Excluir(familiaID, id string) error
}

// RendaExtraHandler contém os handlers HTTP para rendas extras.
type RendaExtraHandler struct {
	svc RendaExtraServiceInterface
}

// NewRendaExtraHandler cria um novo RendaExtraHandler.
func NewRendaExtraHandler(svc RendaExtraServiceInterface) *RendaExtraHandler {
	return &RendaExtraHandler{svc: svc}
}

// rendaExtraResponse é a estrutura de resposta JSON para uma renda extra.
type rendaExtraResponse struct {
	ID              string  `json:"id"`
	Descricao       string  `json:"descricao"`
	MembroID        string  `json:"membro_id"`
	DataRecebimento string  `json:"data_recebimento"`
	Valor           float64 `json:"valor"`
}

func toRendaExtraResponse(r *domain.RendaExtra) rendaExtraResponse {
	return rendaExtraResponse{
		ID:              r.ID,
		Descricao:       r.Descricao,
		MembroID:        r.MembroID,
		DataRecebimento: r.DataRecebimento.Format("2006-01-02"),
		Valor:           r.Valor,
	}
}

type criarRendaExtraRequest struct {
	Descricao       string  `json:"descricao"`
	MembroID        string  `json:"membro_id"`
	DataRecebimento string  `json:"data_recebimento"`
	Valor           float64 `json:"valor"`
}

type atualizarRendaExtraRequest struct {
	Descricao       string  `json:"descricao"`
	MembroID        string  `json:"membro_id"`
	DataRecebimento string  `json:"data_recebimento"`
	Valor           float64 `json:"valor"`
}

func rendaExtraErroParaHTTP(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrRendaExtraNaoEncontrada):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrDescricaoRendaExtraObrigatoria),
		errors.Is(err, domain.ErrMembroIDRendaExtraObrigatorio),
		errors.Is(err, domain.ErrValorRendaExtraInvalido),
		errors.Is(err, domain.ErrDataRecebimentoRendaExtraObrigatoria):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
	}
}

// Listar retorna todas as rendas extras, com filtro opcional por mes e ano.
// GET /api/v1/rendas-extras
// GET /api/v1/rendas-extras?mes=3&ano=2026
func (h *RendaExtraHandler) Listar(c *gin.Context) {
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

		rendas, err := h.svc.ListarPorMes(familiaID, mes, ano)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
			return
		}

		resp := make([]rendaExtraResponse, 0, len(rendas))
		for _, r := range rendas {
			resp = append(resp, toRendaExtraResponse(r))
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	rendas, err := h.svc.Listar(familiaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	resp := make([]rendaExtraResponse, 0, len(rendas))
	for _, r := range rendas {
		resp = append(resp, toRendaExtraResponse(r))
	}
	c.JSON(http.StatusOK, resp)
}

// Criar cria uma nova renda extra.
// POST /api/v1/rendas-extras
func (h *RendaExtraHandler) Criar(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	var req criarRendaExtraRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	var dataRecebimento time.Time
	if req.DataRecebimento != "" {
		parsed, err := time.Parse("2006-01-02", req.DataRecebimento)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "formato de data_recebimento inválido, use YYYY-MM-DD"})
			return
		}
		dataRecebimento = parsed
	}

	renda, err := h.svc.Criar(familiaID, req.Descricao, req.MembroID, dataRecebimento, req.Valor)
	if err != nil {
		rendaExtraErroParaHTTP(c, err)
		return
	}

	c.JSON(http.StatusCreated, toRendaExtraResponse(renda))
}

// BuscarPorID retorna uma renda extra pelo seu ID.
// GET /api/v1/rendas-extras/:id
func (h *RendaExtraHandler) BuscarPorID(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	renda, err := h.svc.BuscarPorID(familiaID, id)
	if err != nil {
		rendaExtraErroParaHTTP(c, err)
		return
	}

	c.JSON(http.StatusOK, toRendaExtraResponse(renda))
}

// Atualizar atualiza os dados de uma renda extra existente.
// PUT /api/v1/rendas-extras/:id
func (h *RendaExtraHandler) Atualizar(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	var req atualizarRendaExtraRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	var dataRecebimento time.Time
	if req.DataRecebimento != "" {
		parsed, err := time.Parse("2006-01-02", req.DataRecebimento)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "formato de data_recebimento inválido, use YYYY-MM-DD"})
			return
		}
		dataRecebimento = parsed
	}

	renda, err := h.svc.Atualizar(familiaID, id, req.Descricao, req.MembroID, dataRecebimento, req.Valor)
	if err != nil {
		rendaExtraErroParaHTTP(c, err)
		return
	}

	c.JSON(http.StatusOK, toRendaExtraResponse(renda))
}

// Excluir realiza o soft-delete de uma renda extra.
// DELETE /api/v1/rendas-extras/:id
func (h *RendaExtraHandler) Excluir(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	if err := h.svc.Excluir(familiaID, id); err != nil {
		rendaExtraErroParaHTTP(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
