package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// RendaFixaServiceInterface define os métodos do service usados pelo handler.
// Redeclarada aqui para desacoplar o pacote handler do service sem importação circular.
type RendaFixaServiceInterface interface {
	Criar(familiaID, descricao, membroID string, valor float64, diaRecebimento int, dataInicio time.Time, dataFim *time.Time) (*domain.RendaFixa, error)
	BuscarPorID(familiaID, id string) (*domain.RendaFixa, error)
	Listar(familiaID string) ([]*domain.RendaFixa, error)
	ListarAtivas(familiaID string) ([]*domain.RendaFixa, error)
	ListarVigentesPorMes(familiaID string, mes, ano int) ([]*domain.RendaFixa, error)
	Atualizar(familiaID, id, descricao, membroID string, valor float64, diaRecebimento int, ativa bool, dataInicio time.Time, dataFim *time.Time) (*domain.RendaFixa, error)
	Inativar(familiaID, id string) error
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
	DataInicio     string  `json:"data_inicio"` // "YYYY-MM-DD"
	DataFim        *string `json:"data_fim"`    // "YYYY-MM-DD" or null
}

func toRendaFixaResponse(r *domain.RendaFixa) rendaFixaResponse {
	resp := rendaFixaResponse{
		ID:             r.ID,
		Descricao:      r.Descricao,
		MembroID:       r.MembroID,
		Valor:          r.Valor,
		DiaRecebimento: r.DiaRecebimento,
		Ativa:          r.Ativa,
		DataInicio:     r.DataInicio.Format("2006-01-02"),
	}
	if r.DataFim != nil {
		s := r.DataFim.Format("2006-01-02")
		resp.DataFim = &s
	}
	return resp
}

type criarRendaFixaRequest struct {
	Descricao      string  `json:"descricao"`
	MembroID       string  `json:"membro_id"`
	Valor          float64 `json:"valor"`
	DiaRecebimento int     `json:"dia_recebimento"`
	DataInicio     string  `json:"data_inicio"`
	DataFim        *string `json:"data_fim"`
}

type atualizarRendaFixaRequest struct {
	Descricao      string  `json:"descricao"       binding:"required"`
	MembroID       string  `json:"membro_id"       binding:"required"`
	Valor          float64 `json:"valor"           binding:"required"`
	DiaRecebimento int     `json:"dia_recebimento" binding:"required"`
	Ativa          bool    `json:"ativa"`
	DataInicio     string  `json:"data_inicio"`
	DataFim        *string `json:"data_fim"`
}

func (h *RendaFixaHandler) erroDominio(c *gin.Context, err error) bool {
	switch {
	case errors.Is(err, domain.ErrDescricaoRendaObrigatoria),
		errors.Is(err, domain.ErrMembroIDRendaObrigatorio),
		errors.Is(err, domain.ErrValorRendaInvalido),
		errors.Is(err, domain.ErrDiaRecebimentoInvalido),
		errors.Is(err, domain.ErrDataInicioRendaObrigatoria):
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
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	rendas, err := h.svc.Listar(familiaID)
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

// ListarVigentesPorMes retorna as rendas fixas vigentes no mês/ano informado.
// GET /api/v1/rendas-fixas/vigentes?mes=3&ano=2026
func (h *RendaFixaHandler) ListarVigentesPorMes(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	mesStr := c.Query("mes")
	anoStr := c.Query("ano")

	if mesStr == "" || anoStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetros mes e ano são obrigatórios"})
		return
	}

	mes, err := strconv.Atoi(mesStr)
	if err != nil || mes < 1 || mes > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetro mes inválido"})
		return
	}

	ano, err := strconv.Atoi(anoStr)
	if err != nil || ano < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetro ano inválido"})
		return
	}

	rendas, err := h.svc.ListarVigentesPorMes(familiaID, mes, ano)
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
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	var req criarRendaFixaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	if req.DataInicio == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrDataInicioRendaObrigatoria.Error()})
		return
	}

	dataInicio, err := time.Parse("2006-01-02", req.DataInicio)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data_inicio inválida, use o formato YYYY-MM-DD"})
		return
	}

	var dataFim *time.Time
	if req.DataFim != nil && *req.DataFim != "" {
		df, err := time.Parse("2006-01-02", *req.DataFim)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "data_fim inválida, use o formato YYYY-MM-DD"})
			return
		}
		dataFim = &df
	}

	renda, err := h.svc.Criar(familiaID, req.Descricao, req.MembroID, req.Valor, req.DiaRecebimento, dataInicio, dataFim)
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
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	renda, err := h.svc.BuscarPorID(familiaID, id)
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
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	var req atualizarRendaFixaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	if req.DataInicio == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrDataInicioRendaObrigatoria.Error()})
		return
	}

	dataInicio, err := time.Parse("2006-01-02", req.DataInicio)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data_inicio inválida, use o formato YYYY-MM-DD"})
		return
	}

	var dataFim *time.Time
	if req.DataFim != nil && *req.DataFim != "" {
		df, err := time.Parse("2006-01-02", *req.DataFim)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "data_fim inválida, use o formato YYYY-MM-DD"})
			return
		}
		dataFim = &df
	}

	renda, err := h.svc.Atualizar(familiaID, id, req.Descricao, req.MembroID, req.Valor, req.DiaRecebimento, req.Ativa, dataInicio, dataFim)
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

	c.JSON(http.StatusOK, gin.H{"message": "renda fixa inativada com sucesso"})
}
