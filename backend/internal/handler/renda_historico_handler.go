package handler

import (
	"net/http"
	"strconv"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// RendaHistoricoServiceInterface define os métodos do service usados pelo handler de histórico de rendas.
type RendaHistoricoServiceInterface interface {
	BuscarHistorico(familiaID string, filtro domain.FiltroHistoricoRendas) (*domain.HistoricoRendas, error)
}

// RendaHistoricoHandler contém os handlers HTTP para o histórico de rendas.
type RendaHistoricoHandler struct {
	svc RendaHistoricoServiceInterface
}

// NewRendaHistoricoHandler cria um novo RendaHistoricoHandler.
func NewRendaHistoricoHandler(svc RendaHistoricoServiceInterface) *RendaHistoricoHandler {
	return &RendaHistoricoHandler{svc: svc}
}

// Historico retorna o histórico e resumo de rendas com filtros opcionais.
// GET /api/v1/rendas/historico?tipo=fixa|variavel|extra|investimento&membro_id=...&mes=1-12&ano=2001+
func (h *RendaHistoricoHandler) Historico(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	var filtro domain.FiltroHistoricoRendas

	// Filtro: tipo
	if tipoStr := c.Query("tipo"); tipoStr != "" {
		tipo := domain.TipoRenda(tipoStr)
		switch tipo {
		case domain.TipoRendaFixa, domain.TipoRendaVariavel, domain.TipoRendaExtra, domain.TipoRendaInvestimento:
			filtro.Tipo = tipo
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "tipo inválido: deve ser fixa, variavel, extra ou investimento"})
			return
		}
	}

	// Filtro: membro_id
	filtro.MembroID = c.Query("membro_id")

	// Filtro: mes
	if mesStr := c.Query("mes"); mesStr != "" {
		v, err := strconv.Atoi(mesStr)
		if err != nil || v < 1 || v > 12 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "mes deve ser um número entre 1 e 12"})
			return
		}
		filtro.Mes = v
	}

	// Filtro: ano
	if anoStr := c.Query("ano"); anoStr != "" {
		v, err := strconv.Atoi(anoStr)
		if err != nil || v <= 2000 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ano deve ser um número maior que 2000"})
			return
		}
		filtro.Ano = v
	}

	historico, err := h.svc.BuscarHistorico(familiaID, filtro)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	c.JSON(http.StatusOK, historico)
}
