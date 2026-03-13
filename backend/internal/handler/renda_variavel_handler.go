package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// RendaVariavelServiceInterface define os métodos do service usados pelo handler.
// Redeclarada aqui para desacoplar o pacote handler do service sem importação circular.
type RendaVariavelServiceInterface interface {
	Criar(familiaID, descricao, membroID string, mesReferencia, anoReferencia int, valor float64, dataRecebimento time.Time) (*domain.RendaVariavel, error)
	BuscarPorID(familiaID, id string) (*domain.RendaVariavel, error)
	Listar(familiaID string) ([]*domain.RendaVariavel, error)
	ListarPorMes(familiaID string, mes, ano int) ([]*domain.RendaVariavel, error)
	Atualizar(familiaID, id, descricao, membroID string, mesReferencia, anoReferencia int, valor float64, dataRecebimento time.Time) (*domain.RendaVariavel, error)
	Excluir(familiaID, id string) error
}

// RendaVariavelHandler contém os handlers HTTP para rendas variáveis.
type RendaVariavelHandler struct {
	svc RendaVariavelServiceInterface
}

// NewRendaVariavelHandler cria um novo RendaVariavelHandler.
func NewRendaVariavelHandler(svc RendaVariavelServiceInterface) *RendaVariavelHandler {
	return &RendaVariavelHandler{svc: svc}
}

// rendaVariavelResponse é a estrutura de resposta JSON para uma renda variável.
type rendaVariavelResponse struct {
	ID              string  `json:"id"`
	Descricao       string  `json:"descricao"`
	MembroID        string  `json:"membro_id"`
	MesReferencia   int     `json:"mes_referencia"`
	AnoReferencia   int     `json:"ano_referencia"`
	Valor           float64 `json:"valor"`
	DataRecebimento string  `json:"data_recebimento"`
}

func toRendaVariavelResponse(r *domain.RendaVariavel) rendaVariavelResponse {
	return rendaVariavelResponse{
		ID:              r.ID,
		Descricao:       r.Descricao,
		MembroID:        r.MembroID,
		MesReferencia:   r.MesReferencia,
		AnoReferencia:   r.AnoReferencia,
		Valor:           r.Valor,
		DataRecebimento: r.DataRecebimento.Format("2006-01-02"),
	}
}

type criarRendaVariavelRequest struct {
	Descricao       string  `json:"descricao"`
	MembroID        string  `json:"membro_id"`
	MesReferencia   int     `json:"mes_referencia"`
	AnoReferencia   int     `json:"ano_referencia"`
	Valor           float64 `json:"valor"`
	DataRecebimento string  `json:"data_recebimento"`
}

type atualizarRendaVariavelRequest struct {
	Descricao       string  `json:"descricao"`
	MembroID        string  `json:"membro_id"`
	MesReferencia   int     `json:"mes_referencia"`
	AnoReferencia   int     `json:"ano_referencia"`
	Valor           float64 `json:"valor"`
	DataRecebimento string  `json:"data_recebimento"`
}

func rendaVariavelErroParaHTTP(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrRendaVariavelNaoEncontrada):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrDescricaoRendaVariavelObrigatoria),
		errors.Is(err, domain.ErrMembroIDRendaVariavelObrigatorio),
		errors.Is(err, domain.ErrMesReferenciaRendaVariavelInvalido),
		errors.Is(err, domain.ErrAnoReferenciaRendaVariavelInvalido),
		errors.Is(err, domain.ErrValorRendaVariavelInvalido),
		errors.Is(err, domain.ErrDataRecebimentoRendaVariavelObrigatoria):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
	}
}

// Listar retorna todas as rendas variáveis, com filtro opcional por mes e ano.
// GET /api/v1/rendas-variaveis
// GET /api/v1/rendas-variaveis?mes=3&ano=2026
func (h *RendaVariavelHandler) Listar(c *gin.Context) {
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

		resp := make([]rendaVariavelResponse, 0, len(rendas))
		for _, r := range rendas {
			resp = append(resp, toRendaVariavelResponse(r))
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	rendas, err := h.svc.Listar(familiaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	resp := make([]rendaVariavelResponse, 0, len(rendas))
	for _, r := range rendas {
		resp = append(resp, toRendaVariavelResponse(r))
	}
	c.JSON(http.StatusOK, resp)
}

// Criar cria uma nova renda variável.
// POST /api/v1/rendas-variaveis
func (h *RendaVariavelHandler) Criar(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	var req criarRendaVariavelRequest
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

	renda, err := h.svc.Criar(familiaID, req.Descricao, req.MembroID, req.MesReferencia, req.AnoReferencia, req.Valor, dataRecebimento)
	if err != nil {
		rendaVariavelErroParaHTTP(c, err)
		return
	}

	c.JSON(http.StatusCreated, toRendaVariavelResponse(renda))
}

// BuscarPorID retorna uma renda variável pelo seu ID.
// GET /api/v1/rendas-variaveis/:id
func (h *RendaVariavelHandler) BuscarPorID(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	renda, err := h.svc.BuscarPorID(familiaID, id)
	if err != nil {
		rendaVariavelErroParaHTTP(c, err)
		return
	}

	c.JSON(http.StatusOK, toRendaVariavelResponse(renda))
}

// Atualizar atualiza os dados de uma renda variável existente.
// PUT /api/v1/rendas-variaveis/:id
func (h *RendaVariavelHandler) Atualizar(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	var req atualizarRendaVariavelRequest
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

	renda, err := h.svc.Atualizar(familiaID, id, req.Descricao, req.MembroID, req.MesReferencia, req.AnoReferencia, req.Valor, dataRecebimento)
	if err != nil {
		rendaVariavelErroParaHTTP(c, err)
		return
	}

	c.JSON(http.StatusOK, toRendaVariavelResponse(renda))
}

// Excluir realiza o soft-delete de uma renda variável.
// DELETE /api/v1/rendas-variaveis/:id
func (h *RendaVariavelHandler) Excluir(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	if err := h.svc.Excluir(familiaID, id); err != nil {
		rendaVariavelErroParaHTTP(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
