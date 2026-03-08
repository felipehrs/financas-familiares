package service

import "github.com/felipehrs/financas-familiares/backend/internal/domain"

// RendaFixaRepositoryForDashboard define os métodos de renda fixa usados pelo DashboardService.
// Redeclarada aqui para evitar import circular e manter desacoplamento.
type RendaFixaRepositoryForDashboard interface {
	ListarAtivas() ([]*domain.RendaFixa, error)
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

// ResumoMensal calcula o resumo financeiro do mês especificado aplicando a RN06 (MVP).
// TOTAL_RENDAS = soma(valor) das rendas_fixas ativas.
// TOTAL_DESPESAS = soma(valor_parcela) das despesas_cartao da fatura do mês/ano.
// SALDO = TOTAL_RENDAS − TOTAL_DESPESAS.
func (s *DashboardService) ResumoMensal(mes, ano int) (*ResumoMensal, error) {
	rendas, err := s.rendaRepo.ListarAtivas()
	if err != nil {
		return nil, err
	}

	var totalRendas float64
	for _, r := range rendas {
		totalRendas += r.Valor
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
