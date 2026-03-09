package service

import (
	"math"
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

// RendaVariavelRepositoryForDashboard define os métodos de renda variável usados pelo DashboardService.
type RendaVariavelRepositoryForDashboard interface {
	ListarPorMes(mes, ano int) ([]*domain.RendaVariavel, error)
}

// RendaExtraRepositoryForDashboard define os métodos de renda extra usados pelo DashboardService.
type RendaExtraRepositoryForDashboard interface {
	ListarPorMes(mes, ano int) ([]*domain.RendaExtra, error)
}

// RendimentoRepositoryForDashboard define os métodos de rendimento de investimento usados pelo DashboardService.
type RendimentoRepositoryForDashboard interface {
	ListarPorMes(mes, ano int) ([]*domain.RendimentoInvestimento, error)
}

// AssinaturaRepositoryForDashboard define os métodos de assinatura usados pelo DashboardService.
type AssinaturaRepositoryForDashboard interface {
	ListarAtivas() ([]*domain.Assinatura, error)
}

// ContaFixaRepositoryForDashboard define os métodos de conta fixa usados pelo DashboardService.
type ContaFixaRepositoryForDashboard interface {
	ListarAtivas() ([]*domain.ContaFixa, error)
}

// DespesaGeralRepositoryForDashboard define os métodos de despesa geral usados pelo DashboardService.
type DespesaGeralRepositoryForDashboard interface {
	ListarPorMes(mes, ano int) ([]*domain.DespesaGeral, error)
}

// CategoriaDespesa representa uma categoria com seu total e percentual de despesas do mês.
type CategoriaDespesa struct {
	Nome       string  `json:"nome"`
	Total      float64 `json:"total"`
	Percentual float64 `json:"percentual"`
}

// ResumoCategorias agrupa as categorias de despesas de um mês.
type ResumoCategorias struct {
	Mes           int                `json:"mes"`
	Ano           int                `json:"ano"`
	TotalDespesas float64            `json:"total_despesas"`
	Categorias    []CategoriaDespesa `json:"categorias"`
}

// DashboardRepositoryForCategorias define a query agregada de despesas por categoria.
type DashboardRepositoryForCategorias interface {
	DespesasPorCategoria(mes, ano int) ([]domain.CategoriaTotalRaw, error)
}

// ResumoMensal contém os totais calculados para um determinado mês/ano.
type ResumoMensal struct {
	Mes int `json:"mes"`
	Ano int `json:"ano"`
	// Rendas operacionais (entram no saldo — RN06)
	TotalRendaFixa             float64 `json:"total_renda_fixa"`
	TotalRendaVariavel         float64 `json:"total_renda_variavel"`
	TotalRendaExtra            float64 `json:"total_renda_extra"`
	TotalRendimentoDistribuido float64 `json:"total_rendimento_distribuido"`
	TotalRendasOperacionais    float64 `json:"total_rendas_operacionais"`
	// Rendimentos informativos (não entram no saldo — RN08)
	TotalRendimentoInvestimento float64 `json:"total_rendimento_investimento"`
	// Despesas
	TotalFaturaCartoes  float64 `json:"total_fatura_cartoes"`
	TotalAssinaturas    float64 `json:"total_assinaturas"`
	TotalContasFixas    float64 `json:"total_contas_fixas"`
	TotalDespesasGerais float64 `json:"total_despesas_gerais"`
	TotalDespesas       float64 `json:"total_despesas"`
	// Resultado
	Saldo float64 `json:"saldo"`
}

// DashboardService implementa a lógica de negócio para o dashboard financeiro.
type DashboardService struct {
	rendaRepo         RendaFixaRepositoryForDashboard
	despesaRepo       DespesaCartaoRepositoryForDashboard
	rendaVariavelRepo RendaVariavelRepositoryForDashboard
	rendaExtraRepo    RendaExtraRepositoryForDashboard
	rendimentoRepo    RendimentoRepositoryForDashboard
	assinaturaRepo    AssinaturaRepositoryForDashboard
	contaFixaRepo     ContaFixaRepositoryForDashboard
	despesaGeralRepo  DespesaGeralRepositoryForDashboard
	categoriasRepo    DashboardRepositoryForCategorias
}

// NewDashboardService cria uma nova instância do DashboardService.
func NewDashboardService(
	rendaRepo RendaFixaRepositoryForDashboard,
	despesaRepo DespesaCartaoRepositoryForDashboard,
	rendaVariavelRepo RendaVariavelRepositoryForDashboard,
	rendaExtraRepo RendaExtraRepositoryForDashboard,
	rendimentoRepo RendimentoRepositoryForDashboard,
	assinaturaRepo AssinaturaRepositoryForDashboard,
	contaFixaRepo ContaFixaRepositoryForDashboard,
	despesaGeralRepo DespesaGeralRepositoryForDashboard,
	categoriasRepo DashboardRepositoryForCategorias,
) *DashboardService {
	return &DashboardService{
		rendaRepo:         rendaRepo,
		despesaRepo:       despesaRepo,
		rendaVariavelRepo: rendaVariavelRepo,
		rendaExtraRepo:    rendaExtraRepo,
		rendimentoRepo:    rendimentoRepo,
		assinaturaRepo:    assinaturaRepo,
		contaFixaRepo:     contaFixaRepo,
		despesaGeralRepo:  despesaGeralRepo,
		categoriasRepo:    categoriasRepo,
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

// ResumoMensal calcula o resumo financeiro do mês especificado aplicando RN06 e RN08.
// TotalRendasOperacionais = RendaFixa + RendaVariavel + RendaExtra + ValorDistribuido (RN06).
// TotalRendimentoInvestimento = soma(Valor) dos rendimentos — informativo apenas (RN08).
// TotalDespesas = FaturaCartoes + Assinaturas + ContasFixas + DespesasGerais.
// Saldo = TotalRendasOperacionais − TotalDespesas.
func (s *DashboardService) ResumoMensal(mes, ano int) (*ResumoMensal, error) {
	// --- Renda Fixa ---
	rendas, err := s.rendaRepo.ListarVigentesPorMes(mes, ano)
	if err != nil {
		return nil, err
	}
	var totalRendaFixa float64
	for _, r := range rendas {
		totalRendaFixa += ValorProporcionado(r, mes, ano)
	}

	// --- Despesas de Cartão ---
	despesasCartao, err := s.despesaRepo.ListarPorFaturaGlobal(mes, ano)
	if err != nil {
		return nil, err
	}
	var totalFaturaCartoes float64
	for _, d := range despesasCartao {
		totalFaturaCartoes += d.ValorParcela
	}

	// --- Renda Variável ---
	rendasVariaveis, err := s.rendaVariavelRepo.ListarPorMes(mes, ano)
	if err != nil {
		return nil, err
	}
	var totalRendaVariavel float64
	for _, r := range rendasVariaveis {
		totalRendaVariavel += r.Valor
	}

	// --- Renda Extra ---
	rendasExtras, err := s.rendaExtraRepo.ListarPorMes(mes, ano)
	if err != nil {
		return nil, err
	}
	var totalRendaExtra float64
	for _, r := range rendasExtras {
		totalRendaExtra += r.Valor
	}

	// --- Rendimentos de Investimento ---
	rendimentos, err := s.rendimentoRepo.ListarPorMes(mes, ano)
	if err != nil {
		return nil, err
	}
	var totalRendimentoInvestimento float64
	var totalRendimentoDistribuido float64
	for _, r := range rendimentos {
		totalRendimentoInvestimento += r.Valor
		totalRendimentoDistribuido += r.ValorDistribuido
	}

	// --- Assinaturas Ativas ---
	assinaturas, err := s.assinaturaRepo.ListarAtivas()
	if err != nil {
		return nil, err
	}
	var totalAssinaturas float64
	for _, a := range assinaturas {
		totalAssinaturas += a.Valor
	}

	// --- Contas Fixas Ativas ---
	contasFixas, err := s.contaFixaRepo.ListarAtivas()
	if err != nil {
		return nil, err
	}
	var totalContasFixas float64
	for _, c := range contasFixas {
		totalContasFixas += c.Valor
	}

	// --- Despesas Gerais ---
	despesasGerais, err := s.despesaGeralRepo.ListarPorMes(mes, ano)
	if err != nil {
		return nil, err
	}
	var totalDespesasGerais float64
	for _, d := range despesasGerais {
		totalDespesasGerais += d.Valor
	}

	// --- Totalizadores ---
	totalRendasOperacionais := totalRendaFixa + totalRendaVariavel + totalRendaExtra + totalRendimentoDistribuido
	totalDespesas := totalFaturaCartoes + totalAssinaturas + totalContasFixas + totalDespesasGerais
	saldo := totalRendasOperacionais - totalDespesas

	return &ResumoMensal{
		Mes:                         mes,
		Ano:                         ano,
		TotalRendaFixa:              totalRendaFixa,
		TotalRendaVariavel:          totalRendaVariavel,
		TotalRendaExtra:             totalRendaExtra,
		TotalRendimentoDistribuido:  totalRendimentoDistribuido,
		TotalRendasOperacionais:     totalRendasOperacionais,
		TotalRendimentoInvestimento: totalRendimentoInvestimento,
		TotalFaturaCartoes:          totalFaturaCartoes,
		TotalAssinaturas:            totalAssinaturas,
		TotalContasFixas:            totalContasFixas,
		TotalDespesasGerais:         totalDespesasGerais,
		TotalDespesas:               totalDespesas,
		Saldo:                       saldo,
	}, nil
}

// PontoEvolucao representa os totais de um mês para o gráfico de evolução.
type PontoEvolucao struct {
	Mes           int     `json:"mes"`
	Ano           int     `json:"ano"`
	TotalRendas   float64 `json:"total_rendas"`
	TotalDespesas float64 `json:"total_despesas"`
	Saldo         float64 `json:"saldo"`
}

// EvolucaoMensal retorna os totais dos últimos qtdMeses meses em ordem cronológica crescente.
func (s *DashboardService) EvolucaoMensal(qtdMeses int) ([]PontoEvolucao, error) {
	agora := time.Now()
	inicio := agora.AddDate(0, -(qtdMeses - 1), 0)
	inicio = time.Date(inicio.Year(), inicio.Month(), 1, 0, 0, 0, 0, time.UTC)

	pontos := make([]PontoEvolucao, 0, qtdMeses)
	for i := 0; i < qtdMeses; i++ {
		mesAtual := inicio.AddDate(0, i, 0)
		resumo, err := s.ResumoMensal(int(mesAtual.Month()), mesAtual.Year())
		if err != nil {
			return nil, err
		}
		pontos = append(pontos, PontoEvolucao{
			Mes:           int(mesAtual.Month()),
			Ano:           mesAtual.Year(),
			TotalRendas:   resumo.TotalRendasOperacionais,
			TotalDespesas: resumo.TotalDespesas,
			Saldo:         resumo.Saldo,
		})
	}
	return pontos, nil
}

// DespesasPorCategoria retorna o resumo de despesas agrupadas por categoria para um mês/ano.
func (s *DashboardService) DespesasPorCategoria(mes, ano int) (*ResumoCategorias, error) {
	rows, err := s.categoriasRepo.DespesasPorCategoria(mes, ano)
	if err != nil {
		return nil, err
	}

	var totalDespesas float64
	for _, r := range rows {
		totalDespesas += r.Total
	}

	categorias := make([]CategoriaDespesa, 0, len(rows))
	for _, r := range rows {
		var percentual float64
		if totalDespesas > 0 {
			percentual = math.Round((r.Total/totalDespesas)*100*100) / 100
		}
		categorias = append(categorias, CategoriaDespesa{
			Nome:       r.Nome,
			Total:      r.Total,
			Percentual: percentual,
		})
	}

	return &ResumoCategorias{
		Mes:           mes,
		Ano:           ano,
		TotalDespesas: totalDespesas,
		Categorias:    categorias,
	}, nil
}
