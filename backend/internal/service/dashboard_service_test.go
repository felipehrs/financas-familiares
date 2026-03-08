package service_test

import (
	"errors"
	"testing"

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

func (m *MockRendaFixaRepositoryForDashboard) ListarAtivas() ([]*domain.RendaFixa, error) {
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

// ---- Testes de ResumoMensal ----

func TestResumoMensal_ComRendasEDespesas(t *testing.T) {
	rendaRepo := &MockRendaFixaRepositoryForDashboard{
		rendas: []*domain.RendaFixa{
			{ID: "r-1", Descricao: "Salário", Valor: 5000.00, Ativa: true},
			{ID: "r-2", Descricao: "Freelance", Valor: 1500.00, Ativa: true},
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
			{ID: "r-1", Descricao: "Salário", Valor: 3000.00, Ativa: true},
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
			{ID: "r-1", Descricao: "Salário", Valor: 3000.00, Ativa: true},
		},
	}
	despesaRepo := &MockDespesaCartaoRepositoryForDashboard{err: erroEsperado}
	svc := service.NewDashboardService(rendaRepo, despesaRepo)

	resumo, err := svc.ResumoMensal(3, 2026)

	assert.Nil(t, resumo)
	assert.ErrorIs(t, err, erroEsperado)
}
