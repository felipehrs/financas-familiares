package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// DespesaCartaoServiceInterface define os métodos do service usados pelo handler.
// Redeclarada aqui para desacoplar o pacote handler do service sem importação circular.
type DespesaCartaoServiceInterface interface {
	Criar(cartaoID, descricao string, categoriaID *string, dataCompra time.Time, valorTotal float64, numeroParcelas int) ([]*domain.DespesaCartao, error)
	ListarPorCartao(cartaoID string) ([]*domain.DespesaCartao, error)
	ListarPorFatura(cartaoID string, mes, ano int) ([]*domain.DespesaCartao, error)
	BuscarPorID(id string) (*domain.DespesaCartao, error)
	Excluir(id string) error
}

// DespesaCartaoHandler contém os handlers HTTP para despesas de cartão de crédito.
type DespesaCartaoHandler struct {
	svc DespesaCartaoServiceInterface
}

// NewDespesaCartaoHandler cria um novo DespesaCartaoHandler.
func NewDespesaCartaoHandler(svc DespesaCartaoServiceInterface) *DespesaCartaoHandler {
	return &DespesaCartaoHandler{svc: svc}
}

// mesesPTBR mapeia número do mês para abreviação em português.
var mesesPTBR = map[int]string{
	1:  "JAN",
	2:  "FEV",
	3:  "MAR",
	4:  "ABR",
	5:  "MAI",
	6:  "JUN",
	7:  "JUL",
	8:  "AGO",
	9:  "SET",
	10: "OUT",
	11: "NOV",
	12: "DEZ",
}

// formatarFatura retorna a fatura no formato "MAR/26".
func formatarFatura(mes, ano int) string {
	abrev := mesesPTBR[mes]
	anoAbrev := ano % 100
	return fmt.Sprintf("%s/%02d", abrev, anoAbrev)
}

// despesaCartaoResponse é a estrutura de resposta JSON para uma despesa de cartão.
type despesaCartaoResponse struct {
	ID             string  `json:"id"`
	CompraID       string  `json:"compra_id"`
	CartaoID       string  `json:"cartao_id"`
	CategoriaID    *string `json:"categoria_id"`
	Descricao      string  `json:"descricao"`
	DataCompra     string  `json:"data_compra"`
	ValorTotal     float64 `json:"valor_total"`
	NumeroParcelas int     `json:"numero_parcelas"`
	ParcelaNumero  int     `json:"parcela_numero"`
	ValorParcela   float64 `json:"valor_parcela"`
	FaturaMes      int     `json:"fatura_mes"`
	FaturaAno      int     `json:"fatura_ano"`
	Fatura         string  `json:"fatura"`
}

func toDespesaCartaoResponse(d *domain.DespesaCartao) despesaCartaoResponse {
	return despesaCartaoResponse{
		ID:             d.ID,
		CompraID:       d.CompraID,
		CartaoID:       d.CartaoID,
		CategoriaID:    d.CategoriaID,
		Descricao:      d.Descricao,
		DataCompra:     d.DataCompra.Format("2006-01-02"),
		ValorTotal:     d.ValorTotal,
		NumeroParcelas: d.NumeroParcelas,
		ParcelaNumero:  d.ParcelaNumero,
		ValorParcela:   d.ValorParcela,
		FaturaMes:      d.FaturaMes,
		FaturaAno:      d.FaturaAno,
		Fatura:         formatarFatura(d.FaturaMes, d.FaturaAno),
	}
}

type criarDespesaCartaoRequest struct {
	Descricao      string  `json:"descricao"`
	CategoriaID    *string `json:"categoria_id"`
	DataCompra     string  `json:"data_compra"`
	ValorTotal     float64 `json:"valor_total"`
	NumeroParcelas int     `json:"numero_parcelas"`
}

func (h *DespesaCartaoHandler) erroDominio(c *gin.Context, err error) bool {
	switch {
	case errors.Is(err, domain.ErrCartaoIDObrigatorio),
		errors.Is(err, domain.ErrDescricaoObrigatoria),
		errors.Is(err, domain.ErrValorTotalInvalido),
		errors.Is(err, domain.ErrNumeroParcelas):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return true
	case errors.Is(err, domain.ErrDespesaCartaoNaoEncontrada):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return true
	case errors.Is(err, domain.ErrCartaoNaoEncontrado):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return true
	}
	return false
}

// ListarPorCartao retorna todas as despesas de um cartão.
// GET /api/v1/cartoes/:id/despesas
func (h *DespesaCartaoHandler) ListarPorCartao(c *gin.Context) {
	cartaoID := c.Param("id")

	despesas, err := h.svc.ListarPorCartao(cartaoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	resp := make([]despesaCartaoResponse, 0, len(despesas))
	for _, d := range despesas {
		resp = append(resp, toDespesaCartaoResponse(d))
	}

	c.JSON(http.StatusOK, resp)
}

// ListarPorFatura retorna as despesas de um cartão filtradas por mês/ano de fatura.
// GET /api/v1/cartoes/:id/despesas/fatura?mes=3&ano=2026
func (h *DespesaCartaoHandler) ListarPorFatura(c *gin.Context) {
	cartaoID := c.Param("id")

	mesStr := c.Query("mes")
	anoStr := c.Query("ano")

	if mesStr == "" || anoStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetros 'mes' e 'ano' são obrigatórios"})
		return
	}

	mes, err := strconv.Atoi(mesStr)
	if err != nil || mes < 1 || mes > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetro 'mes' inválido"})
		return
	}

	ano, err := strconv.Atoi(anoStr)
	if err != nil || ano < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetro 'ano' inválido"})
		return
	}

	despesas, err := h.svc.ListarPorFatura(cartaoID, mes, ano)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	resp := make([]despesaCartaoResponse, 0, len(despesas))
	for _, d := range despesas {
		resp = append(resp, toDespesaCartaoResponse(d))
	}

	c.JSON(http.StatusOK, resp)
}

// Criar cria uma ou mais despesas em um cartão (uma por parcela — RN03).
// POST /api/v1/cartoes/:id/despesas
// Responde com array JSON de todas as parcelas criadas.
func (h *DespesaCartaoHandler) Criar(c *gin.Context) {
	cartaoID := c.Param("id")

	var req criarDespesaCartaoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	dataCompra, err := time.Parse("2006-01-02", req.DataCompra)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data_compra inválida, use o formato YYYY-MM-DD"})
		return
	}

	// Mantém compatibilidade: se numero_parcelas não informado (0), assume 1
	if req.NumeroParcelas < 1 {
		req.NumeroParcelas = 1
	}

	despesas, err := h.svc.Criar(cartaoID, req.Descricao, req.CategoriaID, dataCompra, req.ValorTotal, req.NumeroParcelas)
	if err != nil {
		if h.erroDominio(c, err) {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	resp := make([]despesaCartaoResponse, 0, len(despesas))
	for _, d := range despesas {
		resp = append(resp, toDespesaCartaoResponse(d))
	}

	c.JSON(http.StatusCreated, resp)
}

// Excluir remove (soft delete) uma despesa de cartão.
// DELETE /api/v1/despesas/:id
func (h *DespesaCartaoHandler) Excluir(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.Excluir(id); err != nil {
		if h.erroDominio(c, err) {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "despesa excluída com sucesso"})
}
