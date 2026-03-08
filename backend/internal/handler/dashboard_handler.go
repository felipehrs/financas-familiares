package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// DashboardServiceInterface define os métodos do service usados pelo handler de dashboard.
// Redeclarada aqui para desacoplar o pacote handler do service sem importação circular.
type DashboardServiceInterface interface {
	ResumoMensal(mes, ano int) (*service.ResumoMensal, error)
}

// DashboardHandler contém os handlers HTTP para o dashboard.
type DashboardHandler struct {
	svc DashboardServiceInterface
}

// NewDashboardHandler cria um novo DashboardHandler.
func NewDashboardHandler(svc DashboardServiceInterface) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

// ResumoMensal retorna o resumo financeiro do mês especificado.
// GET /api/v1/dashboard/resumo?mes=3&ano=2026
// Usa mês e ano atuais se os parâmetros não forem fornecidos.
func (h *DashboardHandler) ResumoMensal(c *gin.Context) {
	agora := time.Now()
	mes := int(agora.Month())
	ano := agora.Year()

	if mesParam := c.Query("mes"); mesParam != "" {
		v, err := strconv.Atoi(mesParam)
		if err != nil || v < 1 || v > 12 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "mes deve ser um número entre 1 e 12"})
			return
		}
		mes = v
	}

	if anoParam := c.Query("ano"); anoParam != "" {
		v, err := strconv.Atoi(anoParam)
		if err != nil || v <= 2000 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ano deve ser um número maior que 2000"})
			return
		}
		ano = v
	}

	resumo, err := h.svc.ResumoMensal(mes, ano)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	c.JSON(http.StatusOK, resumo)
}
