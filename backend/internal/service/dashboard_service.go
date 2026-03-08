package service

import (
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
)

// RendaFixaRepositoryForDashboard define os métodos de renda fixa usados pelo DashboardService.
// Redeclarada aqui para evitar import circular e manter desacoplamento.
type RendaFixaRepositoryForDashboard interface {
	ListarVigentesPorMes(mes, ano int) ([]*domain.RendaFixa, error)
}

// DespesaCartaoRepositoryForDashboard define os métodos de despesa de cartão usados pelo DashboardService.
// Redeclarada aqui para evitar import circular e manter desacoplamento.
type DespesaCartaoRepositoryForDashboard interface {
	ListarPorFaturaGlobal(mes, ano int) ([]*domain.DespesaCartao, error)
}

// ResumoMensal contém os totais calculados para um determinado mês/ano.
type ResumoMensal struct {
	Mes           int     `json:"mes"`
	Ano           int     `json:"ano"`
	TotalRendas   float64 `json:"total_rendas"`
	TotalDespesas float64 `json:"total_despesas"`
	Saldo         float64 `json:"saldo"`
}

// DashboardService implementa a lógica de negócio para o dashboard financeiro.
type DashboardService struct {
	rendaRepo   RendaFixaRepositoryForDashboard
	despesaRepo DespesaCartaoRepositoryForDashboard
}

// NewDashboardService cria uma nova instância do DashboardService.
func NewDashboardService(rendaRepo RendaFixaRepositoryForDashboard, despesaRepo DespesaCartaoRepositoryForDashboard) *DashboardService {
	return &DashboardService{
		rendaRepo:   rendaRepo,
		despesaRepo: despesaRepo,
	}
}

// ValorProporcionado calcula o valor efetivo de uma renda fixa para o mês/ano informado,
// aplicando a RN10: proporcionalidade no mês de início e/ou fim.
func ValorProporcionado(r *domain.RendaFixa, mes, ano int) float64 {
	// time.Date(ano, time.Month(mes+1), 0, ...) retorna o último dia do mês (dia 0 do mês seguinte).
	diasNoMes := time.Date(ano, time.Month(mes+1), 0, 0, 0, 0, 0, time.UTC).Day()

	inicioMes := r.DataInicio.Month() == time.Month(mes) && r.DataInicio.Year() == ano
	fimMes := r.DataFim != nil && r.DataFim.Month() == time.Month(mes) && r.DataFim.Year() == ano

	switch {
	case inicioMes && fimMes:
		dias := r.DataFim.Day() - r.DataInicio.Day() + 1
		return r.Valor * float64(dias) / float64(diasNoMes)
	case inicioMes:
		dias := diasNoMes - r.DataInicio.Day() + 1
		return r.Valor * float64(dias) / float64(diasNoMes)
	case fimMes:
		return r.Valor * float64(r.DataFim.Day()) / float64(diasNoMes)
	default:
		return r.Valor
	}
}

// ResumoMensal calcula o resumo financeiro do mês especificado aplicando a RN06 e RN10.
// TOTAL_RENDAS = soma proporcional das rendas_fixas vigentes no mês/ano.
// TOTAL_DESPESAS = soma(valor_parcela) das despesas_cartao da fatura do mês/ano.
// SALDO = TOTAL_RENDAS − TOTAL_DESPESAS.
func (s *DashboardService) ResumoMensal(mes, ano int) (*ResumoMensal, error) {
	rendas, err := s.rendaRepo.ListarVigentesPorMes(mes, ano)
	if err != nil {
		return nil, err
	}

	var totalRendas float64
	for _, r := range rendas {
		totalRendas += ValorProporcionado(r, mes, ano)
	}

	despesas, err := s.despesaRepo.ListarPorFaturaGlobal(mes, ano)
	if err != nil {
		return nil, err
	}

	var totalDespesas float64
	for _, d := range despesas {
		totalDespesas += d.ValorParcela
	}

	return &ResumoMensal{
		Mes:           mes,
		Ano:           ano,
		TotalRendas:   totalRendas,
		TotalDespesas: totalDespesas,
		Saldo:         totalRendas - totalDespesas,
	}, nil
}
