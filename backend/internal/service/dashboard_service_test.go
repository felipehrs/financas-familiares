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

// dataInicioPadrao é uma data no passado para que os testes que não testam proporcionalidade
// recebam o valor cheio (mês intermediário).
var dataInicioPadrao = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

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
	svc := service.NewDashboardService(rendaRepo, despesaRepo)

	resumo, err := svc.ResumoMensal(3, 2026)

	require.NoError(t, err)
	assert.Equal(t, 3, resumo.Mes)
	assert.Equal(t, 2026, resumo.Ano)
	assert.Equal(t, 6500.00, resumo.TotalRendas)
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
	svc := service.NewDashboardService(rendaRepo, despesaRepo)

	resumo, err := svc.ResumoMensal(3, 2026)

	require.NoError(t, err)
	assert.Equal(t, 3000.00, resumo.TotalRendas)
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
	svc := service.NewDashboardService(rendaRepo, despesaRepo)

	resumo, err := svc.ResumoMensal(3, 2026)

	require.NoError(t, err)
	assert.Equal(t, 0.0, resumo.TotalRendas)
	assert.Equal(t, 500.00, resumo.TotalDespesas)
	assert.Equal(t, -500.00, resumo.Saldo)
}

func TestResumoMensal_ErroNoRepoDeRendas(t *testing.T) {
	erroEsperado := errors.New("falha no banco de dados")
	rendaRepo := &MockRendaFixaRepositoryForDashboard{err: erroEsperado}
	despesaRepo := &MockDespesaCartaoRepositoryForDashboard{}
	svc := service.NewDashboardService(rendaRepo, despesaRepo)

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
	svc := service.NewDashboardService(rendaRepo, despesaRepo)

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
// Dias restantes: 24. Esperado TotalRendas ≈ 15000 × 24/31
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
	svc := service.NewDashboardService(rendaRepo, despesaRepo)

	resumo, err := svc.ResumoMensal(3, 2026)

	require.NoError(t, err)
	esperado := 15000.00 * 24.0 / 31.0
	assert.InDelta(t, esperado, resumo.TotalRendas, 0.001)
	assert.InDelta(t, esperado, resumo.Saldo, 0.001)
}
