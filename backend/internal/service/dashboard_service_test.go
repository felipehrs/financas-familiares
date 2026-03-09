package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/felipehrs/financas-familiares/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockRendaFixaRepositoryForDashboard implementa RendaFixaRepositoryForDashboard para testes.
type MockRendaFixaRepositoryForDashboard struct {
	rendas []*domain.RendaFixa
	err    error
}

func (m *MockRendaFixaRepositoryForDashboard) ListarVigentesPorMes(mes, ano int) ([]*domain.RendaFixa, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.rendas, nil
}

// MockDespesaCartaoRepositoryForDashboard implementa DespesaCartaoRepositoryForDashboard para testes.
type MockDespesaCartaoRepositoryForDashboard struct {
	despesas []*domain.DespesaCartao
	err      error
}

func (m *MockDespesaCartaoRepositoryForDashboard) ListarPorFaturaGlobal(mes, ano int) ([]*domain.DespesaCartao, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.despesas, nil
}

// MockRendaVariavelRepositoryForDashboard implementa RendaVariavelRepositoryForDashboard para testes.
type MockRendaVariavelRepositoryForDashboard struct {
	rendas []*domain.RendaVariavel
	err    error
}

func (m *MockRendaVariavelRepositoryForDashboard) ListarPorMes(mes, ano int) ([]*domain.RendaVariavel, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.rendas, nil
}

// MockRendaExtraRepositoryForDashboard implementa RendaExtraRepositoryForDashboard para testes.
type MockRendaExtraRepositoryForDashboard struct {
	rendas []*domain.RendaExtra
	err    error
}

func (m *MockRendaExtraRepositoryForDashboard) ListarPorMes(mes, ano int) ([]*domain.RendaExtra, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.rendas, nil
}

// MockRendimentoRepositoryForDashboard implementa RendimentoRepositoryForDashboard para testes.
type MockRendimentoRepositoryForDashboard struct {
	rendimentos []*domain.RendimentoInvestimento
	err         error
}

func (m *MockRendimentoRepositoryForDashboard) ListarPorMes(mes, ano int) ([]*domain.RendimentoInvestimento, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.rendimentos, nil
}

// MockAssinaturaRepositoryForDashboard implementa AssinaturaRepositoryForDashboard para testes.
type MockAssinaturaRepositoryForDashboard struct {
	assinaturas []*domain.Assinatura
	err         error
}

func (m *MockAssinaturaRepositoryForDashboard) ListarAtivas() ([]*domain.Assinatura, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.assinaturas, nil
}

// MockContaFixaRepositoryForDashboard implementa ContaFixaRepositoryForDashboard para testes.
type MockContaFixaRepositoryForDashboard struct {
	contas []*domain.ContaFixa
	err    error
}

func (m *MockContaFixaRepositoryForDashboard) ListarAtivas() ([]*domain.ContaFixa, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.contas, nil
}

// MockDespesaGeralRepositoryForDashboard implementa DespesaGeralRepositoryForDashboard para testes.
type MockDespesaGeralRepositoryForDashboard struct {
	despesas []*domain.DespesaGeral
	err      error
}

func (m *MockDespesaGeralRepositoryForDashboard) ListarPorMes(mes, ano int) ([]*domain.DespesaGeral, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.despesas, nil
}

// MockDashboardRepositoryForCategorias implementa DashboardRepositoryForCategorias para testes.
type MockDashboardRepositoryForCategorias struct {
	rows []domain.CategoriaTotalRaw
	err  error
}

func (m *MockDashboardRepositoryForCategorias) DespesasPorCategoria(mes, ano int) ([]domain.CategoriaTotalRaw, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.rows, nil
}

// dataInicioPadrao é uma data no passado para que os testes que não testam proporcionalidade
// recebam o valor cheio (mês intermediário).
var dataInicioPadrao = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// newSvcSimples cria um DashboardService com mocks vazios para os repositórios secundários,
// útil para testes que focam apenas em renda fixa e despesas de cartão.
func newSvcSimples(
	rendaRepo *MockRendaFixaRepositoryForDashboard,
	despesaRepo *MockDespesaCartaoRepositoryForDashboard,
) *service.DashboardService {
	return service.NewDashboardService(
		rendaRepo,
		despesaRepo,
		&MockRendaVariavelRepositoryForDashboard{},
		&MockRendaExtraRepositoryForDashboard{},
		&MockRendimentoRepositoryForDashboard{},
		&MockAssinaturaRepositoryForDashboard{},
		&MockContaFixaRepositoryForDashboard{},
		&MockDespesaGeralRepositoryForDashboard{},
		&MockDashboardRepositoryForCategorias{},
	)
}

// newSvcCategorias cria um DashboardService com mocks vazios para todos os repositórios exceto categorias.
func newSvcCategorias(catRepo *MockDashboardRepositoryForCategorias) *service.DashboardService {
	return service.NewDashboardService(
		&MockRendaFixaRepositoryForDashboard{},
		&MockDespesaCartaoRepositoryForDashboard{},
		&MockRendaVariavelRepositoryForDashboard{},
		&MockRendaExtraRepositoryForDashboard{},
		&MockRendimentoRepositoryForDashboard{},
		&MockAssinaturaRepositoryForDashboard{},
		&MockContaFixaRepositoryForDashboard{},
		&MockDespesaGeralRepositoryForDashboard{},
		catRepo,
	)
}

// newSvcVazio cria um DashboardService com todos os mocks vazios, útil para testes de EvolucaoMensal.
func newSvcVazio() *service.DashboardService {
	return service.NewDashboardService(
		&MockRendaFixaRepositoryForDashboard{},
		&MockDespesaCartaoRepositoryForDashboard{},
		&MockRendaVariavelRepositoryForDashboard{},
		&MockRendaExtraRepositoryForDashboard{},
		&MockRendimentoRepositoryForDashboard{},
		&MockAssinaturaRepositoryForDashboard{},
		&MockContaFixaRepositoryForDashboard{},
		&MockDespesaGeralRepositoryForDashboard{},
		&MockDashboardRepositoryForCategorias{},
	)
}

// ---- Testes de EvolucaoMensal ----

func TestEvolucaoMensal_Retorna12Pontos(t *testing.T) {
	svc := newSvcVazio()
	pontos, err := svc.EvolucaoMensal(12)
	require.NoError(t, err)
	assert.Len(t, pontos, 12)
}

func TestEvolucaoMensal_OrdemCronologica(t *testing.T) {
	svc := newSvcVazio()
	pontos, err := svc.EvolucaoMensal(12)
	require.NoError(t, err)
	// Verificar que cada ponto é posterior ao anterior
	for i := 1; i < len(pontos); i++ {
		prev := time.Date(pontos[i-1].Ano, time.Month(pontos[i-1].Mes), 1, 0, 0, 0, 0, time.UTC)
		curr := time.Date(pontos[i].Ano, time.Month(pontos[i].Mes), 1, 0, 0, 0, 0, time.UTC)
		assert.True(t, curr.After(prev), "ponto %d deve ser posterior ao ponto %d", i, i-1)
	}
}

func TestEvolucaoMensal_UltimoMesEhAtual(t *testing.T) {
	svc := newSvcVazio()
	pontos, err := svc.EvolucaoMensal(12)
	require.NoError(t, err)
	agora := time.Now()
	ultimo := pontos[len(pontos)-1]
	assert.Equal(t, int(agora.Month()), ultimo.Mes)
	assert.Equal(t, agora.Year(), ultimo.Ano)
}

func TestEvolucaoMensal_ErroPropaganado(t *testing.T) {
	mockRendaFixa := &MockRendaFixaRepositoryForDashboard{err: errors.New("db error")}
	svc := service.NewDashboardService(
		mockRendaFixa,
		&MockDespesaCartaoRepositoryForDashboard{},
		&MockRendaVariavelRepositoryForDashboard{},
		&MockRendaExtraRepositoryForDashboard{},
		&MockRendimentoRepositoryForDashboard{},
		&MockAssinaturaRepositoryForDashboard{},
		&MockContaFixaRepositoryForDashboard{},
		&MockDespesaGeralRepositoryForDashboard{},
		&MockDashboardRepositoryForCategorias{},
	)
	_, err := svc.EvolucaoMensal(12)
	assert.Error(t, err)
}

// ---- Testes de ResumoMensal ----

func TestResumoMensal_ComRendasEDespesas(t *testing.T) {
	rendaRepo := &MockRendaFixaRepositoryForDashboard{
		rendas: []*domain.RendaFixa{
			{ID: "r-1", Descricao: "Salário", Valor: 5000.00, Ativa: true, DataInicio: dataInicioPadrao},
			{ID: "r-2", Descricao: "Freelance", Valor: 1500.00, Ativa: true, DataInicio: dataInicioPadrao},
		},
	}
	despesaRepo := &MockDespesaCartaoRepositoryForDashboard{
		despesas: []*domain.DespesaCartao{
			{ID: "d-1", ValorParcela: 800.00, FaturaMes: 3, FaturaAno: 2026},
			{ID: "d-2", ValorParcela: 200.00, FaturaMes: 3, FaturaAno: 2026},
		},
	}
	svc := newSvcSimples(rendaRepo, despesaRepo)

	resumo, err := svc.ResumoMensal(3, 2026)

	require.NoError(t, err)
	assert.Equal(t, 3, resumo.Mes)
	assert.Equal(t, 2026, resumo.Ano)
	assert.Equal(t, 6500.00, resumo.TotalRendaFixa)
	assert.Equal(t, 1000.00, resumo.TotalFaturaCartoes)
	assert.Equal(t, 6500.00, resumo.TotalRendasOperacionais)
	assert.Equal(t, 1000.00, resumo.TotalDespesas)
	assert.Equal(t, 5500.00, resumo.Saldo)
}

func TestResumoMensal_SemDespesas(t *testing.T) {
	rendaRepo := &MockRendaFixaRepositoryForDashboard{
		rendas: []*domain.RendaFixa{
			{ID: "r-1", Descricao: "Salário", Valor: 3000.00, Ativa: true, DataInicio: dataInicioPadrao},
		},
	}
	despesaRepo := &MockDespesaCartaoRepositoryForDashboard{
		despesas: []*domain.DespesaCartao{},
	}
	svc := newSvcSimples(rendaRepo, despesaRepo)

	resumo, err := svc.ResumoMensal(3, 2026)

	require.NoError(t, err)
	assert.Equal(t, 3000.00, resumo.TotalRendaFixa)
	assert.Equal(t, 0.0, resumo.TotalDespesas)
	assert.Equal(t, 3000.00, resumo.Saldo)
}

func TestResumoMensal_SemRendas(t *testing.T) {
	rendaRepo := &MockRendaFixaRepositoryForDashboard{
		rendas: []*domain.RendaFixa{},
	}
	despesaRepo := &MockDespesaCartaoRepositoryForDashboard{
		despesas: []*domain.DespesaCartao{
			{ID: "d-1", ValorParcela: 500.00, FaturaMes: 3, FaturaAno: 2026},
		},
	}
	svc := newSvcSimples(rendaRepo, despesaRepo)

	resumo, err := svc.ResumoMensal(3, 2026)

	require.NoError(t, err)
	assert.Equal(t, 0.0, resumo.TotalRendaFixa)
	assert.Equal(t, 500.00, resumo.TotalDespesas)
	assert.Equal(t, -500.00, resumo.Saldo)
}

func TestResumoMensal_ErroNoRepoDeRendas(t *testing.T) {
	erroEsperado := errors.New("falha no banco de dados")
	rendaRepo := &MockRendaFixaRepositoryForDashboard{err: erroEsperado}
	despesaRepo := &MockDespesaCartaoRepositoryForDashboard{}
	svc := newSvcSimples(rendaRepo, despesaRepo)

	resumo, err := svc.ResumoMensal(3, 2026)

	assert.Nil(t, resumo)
	assert.ErrorIs(t, err, erroEsperado)
}

func TestResumoMensal_ErroNoRepoDeDespesas(t *testing.T) {
	erroEsperado := errors.New("falha no banco de dados")
	rendaRepo := &MockRendaFixaRepositoryForDashboard{
		rendas: []*domain.RendaFixa{
			{ID: "r-1", Descricao: "Salário", Valor: 3000.00, Ativa: true, DataInicio: dataInicioPadrao},
		},
	}
	despesaRepo := &MockDespesaCartaoRepositoryForDashboard{err: erroEsperado}
	svc := newSvcSimples(rendaRepo, despesaRepo)

	resumo, err := svc.ResumoMensal(3, 2026)

	assert.Nil(t, resumo)
	assert.ErrorIs(t, err, erroEsperado)
}

// ---- Testes de ValorProporcionado (RN10) ----

// TestValorProporcionado_MesInicio — data_inicio = 2026-03-08, consulta MAR/2026
// Dias restantes: 31 - 8 + 1 = 24. Esperado: valor × 24/31
func TestValorProporcionado_MesInicio(t *testing.T) {
	dataInicio := time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC)
	r := &domain.RendaFixa{
		Valor:      31000.00,
		DataInicio: dataInicio,
	}

	resultado := service.ValorProporcionado(r, 3, 2026)

	esperado := 31000.00 * 24.0 / 31.0
	assert.InDelta(t, esperado, resultado, 0.001)
}

// TestValorProporcionado_MesFim — data_fim = 2026-03-15, consulta MAR/2026
// Esperado: valor × 15/31
func TestValorProporcionado_MesFim(t *testing.T) {
	dataInicio := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	dataFim := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)
	r := &domain.RendaFixa{
		Valor:      31000.00,
		DataInicio: dataInicio,
		DataFim:    &dataFim,
	}

	resultado := service.ValorProporcionado(r, 3, 2026)

	esperado := 31000.00 * 15.0 / 31.0
	assert.InDelta(t, esperado, resultado, 0.001)
}

// TestValorProporcionado_MesInicioEFim — data_inicio = 2026-03-05, data_fim = 2026-03-20, consulta MAR/2026
// Dias: 20 - 5 + 1 = 16. Esperado: valor × 16/31
func TestValorProporcionado_MesInicioEFim(t *testing.T) {
	dataInicio := time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC)
	dataFim := time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)
	r := &domain.RendaFixa{
		Valor:      31000.00,
		DataInicio: dataInicio,
		DataFim:    &dataFim,
	}

	resultado := service.ValorProporcionado(r, 3, 2026)

	esperado := 31000.00 * 16.0 / 31.0
	assert.InDelta(t, esperado, resultado, 0.001)
}

// TestValorProporcionado_MesIntermediario — data_inicio = 2026-01-01, consulta MAR/2026
// Mês intermediário: deve retornar o valor cheio.
func TestValorProporcionado_MesIntermediario(t *testing.T) {
	dataInicio := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	r := &domain.RendaFixa{
		Valor:      5000.00,
		DataInicio: dataInicio,
	}

	resultado := service.ValorProporcionado(r, 3, 2026)

	assert.Equal(t, 5000.00, resultado)
}

// TestResumoMensal_ComProporcionalidade — renda de R$ 15000 com data_inicio = 2026-03-08, MAR/2026
// Dias restantes: 24. Esperado TotalRendaFixa ≈ 15000 × 24/31
func TestResumoMensal_ComProporcionalidade(t *testing.T) {
	dataInicio := time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC)
	rendaRepo := &MockRendaFixaRepositoryForDashboard{
		rendas: []*domain.RendaFixa{
			{ID: "r-1", Descricao: "Salário", Valor: 15000.00, Ativa: true, DataInicio: dataInicio},
		},
	}
	despesaRepo := &MockDespesaCartaoRepositoryForDashboard{
		despesas: []*domain.DespesaCartao{},
	}
	svc := newSvcSimples(rendaRepo, despesaRepo)

	resumo, err := svc.ResumoMensal(3, 2026)

	require.NoError(t, err)
	esperado := 15000.00 * 24.0 / 31.0
	assert.InDelta(t, esperado, resumo.TotalRendaFixa, 0.001)
	assert.InDelta(t, esperado, resumo.Saldo, 0.001)
}

// ---- Novos testes RN06/RN08 ----

// TestResumoMensal_RN06_TodosOsTipos — verifica que todos os tipos de renda e despesa
// são somados corretamente nos totalizadores.
func TestResumoMensal_RN06_TodosOsTipos(t *testing.T) {
	svc := service.NewDashboardService(
		&MockRendaFixaRepositoryForDashboard{
			rendas: []*domain.RendaFixa{
				{ID: "rf-1", Valor: 5000.00, DataInicio: dataInicioPadrao},
			},
		},
		&MockDespesaCartaoRepositoryForDashboard{
			despesas: []*domain.DespesaCartao{
				{ID: "dc-1", ValorParcela: 300.00},
			},
		},
		&MockRendaVariavelRepositoryForDashboard{
			rendas: []*domain.RendaVariavel{
				{ID: "rv-1", Valor: 1200.00},
				{ID: "rv-2", Valor: 800.00},
			},
		},
		&MockRendaExtraRepositoryForDashboard{
			rendas: []*domain.RendaExtra{
				{ID: "re-1", Valor: 500.00},
			},
		},
		&MockRendimentoRepositoryForDashboard{
			rendimentos: []*domain.RendimentoInvestimento{
				{ID: "ri-1", Valor: 1000.00, ValorDistribuido: 400.00},
			},
		},
		&MockAssinaturaRepositoryForDashboard{
			assinaturas: []*domain.Assinatura{
				{ID: "as-1", Valor: 50.00},
				{ID: "as-2", Valor: 30.00},
			},
		},
		&MockContaFixaRepositoryForDashboard{
			contas: []*domain.ContaFixa{
				{ID: "cf-1", Valor: 200.00},
			},
		},
		&MockDespesaGeralRepositoryForDashboard{
			despesas: []*domain.DespesaGeral{
				{ID: "dg-1", Valor: 150.00},
				{ID: "dg-2", Valor: 100.00},
			},
		},
		&MockDashboardRepositoryForCategorias{},
	)

	resumo, err := svc.ResumoMensal(3, 2026)

	require.NoError(t, err)
	assert.Equal(t, 3, resumo.Mes)
	assert.Equal(t, 2026, resumo.Ano)

	// Rendas
	assert.Equal(t, 5000.00, resumo.TotalRendaFixa)
	assert.Equal(t, 2000.00, resumo.TotalRendaVariavel) // 1200 + 800
	assert.Equal(t, 500.00, resumo.TotalRendaExtra)
	assert.Equal(t, 400.00, resumo.TotalRendimentoDistribuido)
	assert.Equal(t, 7900.00, resumo.TotalRendasOperacionais) // 5000 + 2000 + 500 + 400

	// Rendimento informativo
	assert.Equal(t, 1000.00, resumo.TotalRendimentoInvestimento)

	// Despesas
	assert.Equal(t, 300.00, resumo.TotalFaturaCartoes)
	assert.Equal(t, 80.00, resumo.TotalAssinaturas) // 50 + 30
	assert.Equal(t, 200.00, resumo.TotalContasFixas)
	assert.Equal(t, 250.00, resumo.TotalDespesasGerais) // 150 + 100
	assert.Equal(t, 830.00, resumo.TotalDespesas)       // 300 + 80 + 200 + 250

	// Saldo
	assert.Equal(t, 7070.00, resumo.Saldo) // 7900 - 830
}

// TestResumoMensal_RN08_RendimentoNaoDistribuidoNaoEntreNoSaldo — rendimento com ValorDistribuido = 0
// não deve entrar no saldo, mas deve aparecer em TotalRendimentoInvestimento.
func TestResumoMensal_RN08_RendimentoNaoDistribuidoNaoEntreNoSaldo(t *testing.T) {
	svc := service.NewDashboardService(
		&MockRendaFixaRepositoryForDashboard{
			rendas: []*domain.RendaFixa{
				{ID: "rf-1", Valor: 3000.00, DataInicio: dataInicioPadrao},
			},
		},
		&MockDespesaCartaoRepositoryForDashboard{},
		&MockRendaVariavelRepositoryForDashboard{},
		&MockRendaExtraRepositoryForDashboard{},
		&MockRendimentoRepositoryForDashboard{
			rendimentos: []*domain.RendimentoInvestimento{
				{ID: "ri-1", Valor: 2000.00, ValorDistribuido: 0.00},
			},
		},
		&MockAssinaturaRepositoryForDashboard{},
		&MockContaFixaRepositoryForDashboard{},
		&MockDespesaGeralRepositoryForDashboard{},
		&MockDashboardRepositoryForCategorias{},
	)

	resumo, err := svc.ResumoMensal(3, 2026)

	require.NoError(t, err)
	// Rendimento total é informativo
	assert.Equal(t, 2000.00, resumo.TotalRendimentoInvestimento)
	// Distribuído é zero — não entra no saldo
	assert.Equal(t, 0.00, resumo.TotalRendimentoDistribuido)
	// Rendas operacionais = apenas renda fixa
	assert.Equal(t, 3000.00, resumo.TotalRendasOperacionais)
	assert.Equal(t, 3000.00, resumo.Saldo)
}

// ---- Testes de DespesasPorCategoria ----

func TestDespesasPorCategoria_MultiplasCategorias(t *testing.T) {
	mockCat := &MockDashboardRepositoryForCategorias{
		rows: []domain.CategoriaTotalRaw{
			{Nome: "Alimentação", Total: 650.0},
			{Nome: "Transporte", Total: 350.0},
		},
	}
	svc := newSvcCategorias(mockCat)

	resumo, err := svc.DespesasPorCategoria(3, 2026)

	require.NoError(t, err)
	assert.Equal(t, 1000.0, resumo.TotalDespesas)
	assert.Len(t, resumo.Categorias, 2)
	assert.Equal(t, "Alimentação", resumo.Categorias[0].Nome)
	assert.Equal(t, 65.0, resumo.Categorias[0].Percentual)
	assert.Equal(t, "Transporte", resumo.Categorias[1].Nome)
	assert.Equal(t, 35.0, resumo.Categorias[1].Percentual)
}

func TestDespesasPorCategoria_SemCategoria(t *testing.T) {
	mockCat := &MockDashboardRepositoryForCategorias{
		rows: []domain.CategoriaTotalRaw{
			{Nome: "Sem categoria", Total: 200.0},
		},
	}
	svc := newSvcCategorias(mockCat)

	resumo, err := svc.DespesasPorCategoria(3, 2026)

	require.NoError(t, err)
	assert.Equal(t, "Sem categoria", resumo.Categorias[0].Nome)
	assert.Equal(t, 100.0, resumo.Categorias[0].Percentual)
}

func TestDespesasPorCategoria_SemDespesas(t *testing.T) {
	mockCat := &MockDashboardRepositoryForCategorias{rows: []domain.CategoriaTotalRaw{}}
	svc := newSvcCategorias(mockCat)

	resumo, err := svc.DespesasPorCategoria(3, 2026)

	require.NoError(t, err)
	assert.Equal(t, 0.0, resumo.TotalDespesas)
	assert.Empty(t, resumo.Categorias)
}

func TestDespesasPorCategoria_ErroNoRepositorio(t *testing.T) {
	mockCat := &MockDashboardRepositoryForCategorias{err: errors.New("db error")}
	svc := newSvcCategorias(mockCat)

	_, err := svc.DespesasPorCategoria(3, 2026)

	assert.Error(t, err)
}

// ---- Testes de ProjecaoProximosMeses ----

func TestProjecaoProximosMeses_Retorna3Pontos(t *testing.T) {
	svc := newSvcVazio()
	projecoes, err := svc.ProjecaoProximosMeses(3)
	require.NoError(t, err)
	assert.Len(t, projecoes, 3)
}

func TestProjecaoProximosMeses_MesesSaoFuturos(t *testing.T) {
	svc := newSvcVazio()
	projecoes, err := svc.ProjecaoProximosMeses(3)
	require.NoError(t, err)
	agora := time.Now()
	for _, p := range projecoes {
		mesP := time.Date(p.Ano, time.Month(p.Mes), 1, 0, 0, 0, 0, time.UTC)
		mesAtual := time.Date(agora.Year(), agora.Month(), 1, 0, 0, 0, 0, time.UTC)
		assert.True(t, mesP.After(mesAtual))
	}
}

func TestProjecaoProximosMeses_OrdemCronologica(t *testing.T) {
	svc := newSvcVazio()
	projecoes, err := svc.ProjecaoProximosMeses(3)
	require.NoError(t, err)
	for i := 1; i < len(projecoes); i++ {
		prev := time.Date(projecoes[i-1].Ano, time.Month(projecoes[i-1].Mes), 1, 0, 0, 0, 0, time.UTC)
		curr := time.Date(projecoes[i].Ano, time.Month(projecoes[i].Mes), 1, 0, 0, 0, 0, time.UTC)
		assert.True(t, curr.After(prev))
	}
}

func TestProjecaoProximosMeses_SomaRendaFixa(t *testing.T) {
	dataInicio := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	mockRenda := &MockRendaFixaRepositoryForDashboard{
		rendas: []*domain.RendaFixa{
			{ID: "r1", Valor: 5000, DataInicio: dataInicio, DataFim: nil},
		},
	}
	svc := service.NewDashboardService(
		mockRenda,
		&MockDespesaCartaoRepositoryForDashboard{},
		&MockRendaVariavelRepositoryForDashboard{},
		&MockRendaExtraRepositoryForDashboard{},
		&MockRendimentoRepositoryForDashboard{},
		&MockAssinaturaRepositoryForDashboard{},
		&MockContaFixaRepositoryForDashboard{},
		&MockDespesaGeralRepositoryForDashboard{},
		&MockDashboardRepositoryForCategorias{},
	)
	projecoes, err := svc.ProjecaoProximosMeses(3)
	require.NoError(t, err)
	for _, p := range projecoes {
		assert.Equal(t, 5000.0, p.TotalRendas)
	}
}

func TestProjecaoProximosMeses_SomaDespesas(t *testing.T) {
	mockDespesa := &MockDespesaCartaoRepositoryForDashboard{
		despesas: []*domain.DespesaCartao{
			{ID: "d1", ValorParcela: 300},
		},
	}
	mockAssinatura := &MockAssinaturaRepositoryForDashboard{
		assinaturas: []*domain.Assinatura{
			{ID: "a1", Valor: 50},
		},
	}
	mockConta := &MockContaFixaRepositoryForDashboard{
		contas: []*domain.ContaFixa{
			{ID: "c1", Valor: 200},
		},
	}
	svc := service.NewDashboardService(
		&MockRendaFixaRepositoryForDashboard{},
		mockDespesa,
		&MockRendaVariavelRepositoryForDashboard{},
		&MockRendaExtraRepositoryForDashboard{},
		&MockRendimentoRepositoryForDashboard{},
		mockAssinatura,
		mockConta,
		&MockDespesaGeralRepositoryForDashboard{},
		&MockDashboardRepositoryForCategorias{},
	)
	projecoes, err := svc.ProjecaoProximosMeses(3)
	require.NoError(t, err)
	for _, p := range projecoes {
		assert.Equal(t, 300.0, p.TotalCartoes)
		assert.Equal(t, 50.0, p.TotalAssinaturas)
		assert.Equal(t, 200.0, p.TotalContasFixas)
		assert.Equal(t, 550.0, p.TotalDespesas)
	}
}

func TestProjecaoProximosMeses_ErroAssinaturaRepo(t *testing.T) {
	mockAssinatura := &MockAssinaturaRepositoryForDashboard{err: errors.New("db error")}
	svc := service.NewDashboardService(
		&MockRendaFixaRepositoryForDashboard{},
		&MockDespesaCartaoRepositoryForDashboard{},
		&MockRendaVariavelRepositoryForDashboard{},
		&MockRendaExtraRepositoryForDashboard{},
		&MockRendimentoRepositoryForDashboard{},
		mockAssinatura,
		&MockContaFixaRepositoryForDashboard{},
		&MockDespesaGeralRepositoryForDashboard{},
		&MockDashboardRepositoryForCategorias{},
	)
	_, err := svc.ProjecaoProximosMeses(3)
	assert.Error(t, err)
}

// TestResumoMensal_RN08_RendimentoParcialmenteDistribuido — apenas ValorDistribuido entra no saldo,
// não o valor total do rendimento.
func TestResumoMensal_RN08_RendimentoParcialmenteDistribuido(t *testing.T) {
	svc := service.NewDashboardService(
		&MockRendaFixaRepositoryForDashboard{
			rendas: []*domain.RendaFixa{
				{ID: "rf-1", Valor: 3000.00, DataInicio: dataInicioPadrao},
			},
		},
		&MockDespesaCartaoRepositoryForDashboard{},
		&MockRendaVariavelRepositoryForDashboard{},
		&MockRendaExtraRepositoryForDashboard{},
		&MockRendimentoRepositoryForDashboard{
			rendimentos: []*domain.RendimentoInvestimento{
				{ID: "ri-1", Valor: 5000.00, ValorDistribuido: 1500.00},
			},
		},
		&MockAssinaturaRepositoryForDashboard{},
		&MockContaFixaRepositoryForDashboard{},
		&MockDespesaGeralRepositoryForDashboard{},
		&MockDashboardRepositoryForCategorias{},
	)

	resumo, err := svc.ResumoMensal(3, 2026)

	require.NoError(t, err)
	// Rendimento total é informativo
	assert.Equal(t, 5000.00, resumo.TotalRendimentoInvestimento)
	// Apenas o valor distribuído entra no saldo
	assert.Equal(t, 1500.00, resumo.TotalRendimentoDistribuido)
	// Rendas operacionais = renda fixa + valor distribuído
	assert.Equal(t, 4500.00, resumo.TotalRendasOperacionais) // 3000 + 1500
	assert.Equal(t, 4500.00, resumo.Saldo)
}
