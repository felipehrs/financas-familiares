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
	ResumoMensal(familiaID string, mes, ano int) (*service.ResumoMensal, error)
	DespesasPorCategoria(familiaID string, mes, ano int) (*service.ResumoCategorias, error)
	EvolucaoMensal(familiaID string, qtdMeses int) ([]service.PontoEvolucao, error)
	ProjecaoProximosMeses(familiaID string, qtdMeses int) ([]service.MesProjecao, error)
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
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

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

	resumo, err := h.svc.ResumoMensal(familiaID, mes, ano)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	c.JSON(http.StatusOK, resumo)
}

// EvolucaoMensal retorna os totais dos últimos 12 meses para o gráfico de evolução.
// GET /api/v1/dashboard/evolucao
func (h *DashboardHandler) EvolucaoMensal(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	pontos, err := h.svc.EvolucaoMensal(familiaID, 12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}
	c.JSON(http.StatusOK, pontos)
}

// Projecao retorna a projeção dos próximos 3 meses.
// GET /api/v1/dashboard/projecao
func (h *DashboardHandler) Projecao(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	projecao, err := h.svc.ProjecaoProximosMeses(familiaID, 3)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}
	c.JSON(http.StatusOK, projecao)
}

// DespesasPorCategoria retorna as despesas agrupadas por categoria para um mês/ano.
// GET /api/v1/dashboard/categorias?mes=3&ano=2026
func (h *DashboardHandler) DespesasPorCategoria(c *gin.Context) {
	familiaID, ok := getFamiliaID(c)
	if !ok {
		return
	}

	mesStr := c.Query("mes")
	anoStr := c.Query("ano")

	var mes, ano int
	if mesStr != "" && anoStr != "" {
		var err error
		mes, err = strconv.Atoi(mesStr)
		if err != nil || mes < 1 || mes > 12 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetro mes inválido"})
			return
		}
		ano, err = strconv.Atoi(anoStr)
		if err != nil || ano < 2000 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetro ano inválido"})
			return
		}
	} else {
		now := time.Now()
		mes = int(now.Month())
		ano = now.Year()
	}

	resumo, err := h.svc.DespesasPorCategoria(familiaID, mes, ano)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	c.JSON(http.StatusOK, resumo)
}
