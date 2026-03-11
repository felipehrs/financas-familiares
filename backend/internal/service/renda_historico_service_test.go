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

// ---- Mocks para RendaHistoricoService ----

type MockRendaFixaRepoForHistorico struct {
	rendas []*domain.RendaFixa
	err    error
}

func (m *MockRendaFixaRepoForHistorico) Listar(familiaID string) ([]*domain.RendaFixa, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.rendas, nil
}

func (m *MockRendaFixaRepoForHistorico) ListarVigentesPorMes(familiaID string, mes, ano int) ([]*domain.RendaFixa, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.rendas, nil
}

type MockRendaVariavelRepoForHistorico struct {
	rendas []*domain.RendaVariavel
	err    error
}

func (m *MockRendaVariavelRepoForHistorico) Listar(familiaID string) ([]*domain.RendaVariavel, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.rendas, nil
}

func (m *MockRendaVariavelRepoForHistorico) ListarPorMes(familiaID string, mes, ano int) ([]*domain.RendaVariavel, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.rendas, nil
}

type MockRendaExtraRepoForHistorico struct {
	rendas []*domain.RendaExtra
	err    error
}

func (m *MockRendaExtraRepoForHistorico) Listar(familiaID string) ([]*domain.RendaExtra, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.rendas, nil
}

func (m *MockRendaExtraRepoForHistorico) ListarPorMes(familiaID string, mes, ano int) ([]*domain.RendaExtra, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.rendas, nil
}

type MockRendimentoRepoForHistorico struct {
	rendimentos []*domain.RendimentoInvestimento
	err         error
}

func (m *MockRendimentoRepoForHistorico) Listar(familiaID string) ([]*domain.RendimentoInvestimento, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.rendimentos, nil
}

func (m *MockRendimentoRepoForHistorico) ListarPorMes(familiaID string, mes, ano int) ([]*domain.RendimentoInvestimento, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.rendimentos, nil
}

// ---- Helpers ----

var dataInicioPadraoHistorico = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func newHistoricoSvc(
	fixaRepo *MockRendaFixaRepoForHistorico,
	variavelRepo *MockRendaVariavelRepoForHistorico,
	extraRepo *MockRendaExtraRepoForHistorico,
	rendimentoRepo *MockRendimentoRepoForHistorico,
) *service.RendaHistoricoService {
	return service.NewRendaHistoricoService(fixaRepo, variavelRepo, extraRepo, rendimentoRepo)
}

func newHistoricoSvcVazio() *service.RendaHistoricoService {
	return newHistoricoSvc(
		&MockRendaFixaRepoForHistorico{},
		&MockRendaVariavelRepoForHistorico{},
		&MockRendaExtraRepoForHistorico{},
		&MockRendimentoRepoForHistorico{},
	)
}

// ---- Testes ----

// TestHistorico_TodosTiposRetornadosSemFiltro verifica que todos os tipos de renda são retornados
// quando nenhum filtro é aplicado.
func TestHistorico_TodosTiposRetornadosSemFiltro(t *testing.T) {
	svc := newHistoricoSvc(
		&MockRendaFixaRepoForHistorico{rendas: []*domain.RendaFixa{
			{ID: "rf-1", Descricao: "Salário", MembroID: "m-1", Valor: 5000, Ativa: true, DataInicio: dataInicioPadraoHistorico},
		}},
		&MockRendaVariavelRepoForHistorico{rendas: []*domain.RendaVariavel{
			{ID: "rv-1", Descricao: "Freelance", MembroID: "m-1", Valor: 1200, MesReferencia: 3, AnoReferencia: 2026},
		}},
		&MockRendaExtraRepoForHistorico{rendas: []*domain.RendaExtra{
			{ID: "re-1", Descricao: "Bônus", MembroID: "m-2", Valor: 500, DataRecebimento: time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)},
		}},
		&MockRendimentoRepoForHistorico{rendimentos: []*domain.RendimentoInvestimento{
			{ID: "ri-1", Descricao: "CDB", MembroID: "m-1", Valor: 2000, ValorDistribuido: 400, Data: time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC)},
		}},
	)

	resultado, err := svc.BuscarHistorico("familia-id", domain.FiltroHistoricoRendas{})
	require.NoError(t, err)
	assert.Len(t, resultado.Itens, 4)
}

// TestHistorico_FiltroTipoFixa verifica que apenas rendas fixas são retornadas.
func TestHistorico_FiltroTipoFixa(t *testing.T) {
	svc := newHistoricoSvc(
		&MockRendaFixaRepoForHistorico{rendas: []*domain.RendaFixa{
			{ID: "rf-1", Descricao: "Salário", MembroID: "m-1", Valor: 5000, Ativa: true, DataInicio: dataInicioPadraoHistorico},
		}},
		&MockRendaVariavelRepoForHistorico{rendas: []*domain.RendaVariavel{
			{ID: "rv-1", Descricao: "Freelance", MembroID: "m-1", Valor: 1200, MesReferencia: 3, AnoReferencia: 2026},
		}},
		&MockRendaExtraRepoForHistorico{},
		&MockRendimentoRepoForHistorico{},
	)

	resultado, err := svc.BuscarHistorico("familia-id", domain.FiltroHistoricoRendas{Tipo: domain.TipoRendaFixa})
	require.NoError(t, err)
	assert.Len(t, resultado.Itens, 1)
	assert.Equal(t, domain.TipoRendaFixa, resultado.Itens[0].Tipo)
	assert.Equal(t, "rf-1", resultado.Itens[0].ID)
}

// TestHistorico_FiltroTipoVariavel verifica que apenas rendas variáveis são retornadas.
func TestHistorico_FiltroTipoVariavel(t *testing.T) {
	svc := newHistoricoSvc(
		&MockRendaFixaRepoForHistorico{rendas: []*domain.RendaFixa{
			{ID: "rf-1", Descricao: "Salário", MembroID: "m-1", Valor: 5000, Ativa: true, DataInicio: dataInicioPadraoHistorico},
		}},
		&MockRendaVariavelRepoForHistorico{rendas: []*domain.RendaVariavel{
			{ID: "rv-1", Descricao: "Freelance", MembroID: "m-1", Valor: 1200, MesReferencia: 3, AnoReferencia: 2026},
			{ID: "rv-2", Descricao: "Comissão", MembroID: "m-2", Valor: 800, MesReferencia: 2, AnoReferencia: 2026},
		}},
		&MockRendaExtraRepoForHistorico{},
		&MockRendimentoRepoForHistorico{},
	)

	resultado, err := svc.BuscarHistorico("familia-id", domain.FiltroHistoricoRendas{Tipo: domain.TipoRendaVariavel})
	require.NoError(t, err)
	assert.Len(t, resultado.Itens, 2)
	for _, item := range resultado.Itens {
		assert.Equal(t, domain.TipoRendaVariavel, item.Tipo)
	}
}

// TestHistorico_FiltroTipoExtra verifica que apenas rendas extras são retornadas.
func TestHistorico_FiltroTipoExtra(t *testing.T) {
	svc := newHistoricoSvc(
		&MockRendaFixaRepoForHistorico{},
		&MockRendaVariavelRepoForHistorico{},
		&MockRendaExtraRepoForHistorico{rendas: []*domain.RendaExtra{
			{ID: "re-1", Descricao: "Bônus", MembroID: "m-1", Valor: 500, DataRecebimento: time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)},
		}},
		&MockRendimentoRepoForHistorico{rendimentos: []*domain.RendimentoInvestimento{
			{ID: "ri-1", Descricao: "CDB", MembroID: "m-1", Valor: 2000, ValorDistribuido: 400, Data: time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC)},
		}},
	)

	resultado, err := svc.BuscarHistorico("familia-id", domain.FiltroHistoricoRendas{Tipo: domain.TipoRendaExtra})
	require.NoError(t, err)
	assert.Len(t, resultado.Itens, 1)
	assert.Equal(t, domain.TipoRendaExtra, resultado.Itens[0].Tipo)
	assert.Equal(t, "re-1", resultado.Itens[0].ID)
}

// TestHistorico_FiltroTipoInvestimento verifica que apenas rendimentos de investimento são retornados.
func TestHistorico_FiltroTipoInvestimento(t *testing.T) {
	svc := newHistoricoSvc(
		&MockRendaFixaRepoForHistorico{rendas: []*domain.RendaFixa{
			{ID: "rf-1", Descricao: "Salário", MembroID: "m-1", Valor: 5000, Ativa: true, DataInicio: dataInicioPadraoHistorico},
		}},
		&MockRendaVariavelRepoForHistorico{},
		&MockRendaExtraRepoForHistorico{},
		&MockRendimentoRepoForHistorico{rendimentos: []*domain.RendimentoInvestimento{
			{ID: "ri-1", Descricao: "CDB", MembroID: "m-1", Valor: 2000, ValorDistribuido: 400, Data: time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC)},
			{ID: "ri-2", Descricao: "Tesouro", MembroID: "m-2", Valor: 1000, ValorDistribuido: 0, Data: time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC)},
		}},
	)

	resultado, err := svc.BuscarHistorico("familia-id", domain.FiltroHistoricoRendas{Tipo: domain.TipoRendaInvestimento})
	require.NoError(t, err)
	assert.Len(t, resultado.Itens, 2)
	for _, item := range resultado.Itens {
		assert.Equal(t, domain.TipoRendaInvestimento, item.Tipo)
		assert.NotNil(t, item.ValorDistribuido)
	}
}

// TestHistorico_FiltroMembroID verifica que apenas itens do membro especificado são retornados.
func TestHistorico_FiltroMembroID(t *testing.T) {
	svc := newHistoricoSvc(
		&MockRendaFixaRepoForHistorico{rendas: []*domain.RendaFixa{
			{ID: "rf-1", Descricao: "Salário M1", MembroID: "m-1", Valor: 5000, Ativa: true, DataInicio: dataInicioPadraoHistorico},
			{ID: "rf-2", Descricao: "Salário M2", MembroID: "m-2", Valor: 4000, Ativa: true, DataInicio: dataInicioPadraoHistorico},
		}},
		&MockRendaVariavelRepoForHistorico{rendas: []*domain.RendaVariavel{
			{ID: "rv-1", Descricao: "Freelance M1", MembroID: "m-1", Valor: 1200, MesReferencia: 3, AnoReferencia: 2026},
		}},
		&MockRendaExtraRepoForHistorico{rendas: []*domain.RendaExtra{
			{ID: "re-1", Descricao: "Bônus M2", MembroID: "m-2", Valor: 500, DataRecebimento: time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)},
		}},
		&MockRendimentoRepoForHistorico{},
	)

	resultado, err := svc.BuscarHistorico("familia-id", domain.FiltroHistoricoRendas{MembroID: "m-1"})
	require.NoError(t, err)
	// rf-1 + rv-1 = 2 itens
	assert.Len(t, resultado.Itens, 2)
	for _, item := range resultado.Itens {
		assert.Equal(t, "m-1", item.MembroID)
	}
}

// TestHistorico_FiltroMesEAno verifica que o filtro de mês+ano usa os métodos corretos.
func TestHistorico_FiltroMesEAno(t *testing.T) {
	svc := newHistoricoSvc(
		&MockRendaFixaRepoForHistorico{rendas: []*domain.RendaFixa{
			{ID: "rf-1", Descricao: "Salário", MembroID: "m-1", Valor: 5000, Ativa: true, DataInicio: dataInicioPadraoHistorico},
		}},
		&MockRendaVariavelRepoForHistorico{rendas: []*domain.RendaVariavel{
			{ID: "rv-1", Descricao: "Freelance", MembroID: "m-1", Valor: 1200, MesReferencia: 3, AnoReferencia: 2026},
		}},
		&MockRendaExtraRepoForHistorico{rendas: []*domain.RendaExtra{
			{ID: "re-1", Descricao: "Bônus", MembroID: "m-1", Valor: 500, DataRecebimento: time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)},
		}},
		&MockRendimentoRepoForHistorico{rendimentos: []*domain.RendimentoInvestimento{
			{ID: "ri-1", Descricao: "CDB", MembroID: "m-1", Valor: 2000, ValorDistribuido: 400, Data: time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC)},
		}},
	)

	resultado, err := svc.BuscarHistorico("familia-id", domain.FiltroHistoricoRendas{Mes: 3, Ano: 2026})
	require.NoError(t, err)
	// Todos os 4 itens devem ser retornados (mocks retornam todos independente do mês)
	assert.Len(t, resultado.Itens, 4)

	// A renda fixa deve ter mes e ano preenchidos
	var fixaItem *domain.ItemRendaHistorico
	for i, item := range resultado.Itens {
		if item.Tipo == domain.TipoRendaFixa {
			fixaItem = &resultado.Itens[i]
			break
		}
	}
	require.NotNil(t, fixaItem)
	require.NotNil(t, fixaItem.Mes)
	require.NotNil(t, fixaItem.Ano)
	assert.Equal(t, 3, *fixaItem.Mes)
	assert.Equal(t, 2026, *fixaItem.Ano)
}

// TestHistorico_ResumoTotaisCorretos verifica que os totais do resumo são calculados corretamente,
// especialmente que investimento usa ValorDistribuido, não Valor.
func TestHistorico_ResumoTotaisCorretos(t *testing.T) {
	svc := newHistoricoSvc(
		&MockRendaFixaRepoForHistorico{rendas: []*domain.RendaFixa{
			{ID: "rf-1", Descricao: "Salário", MembroID: "m-1", Valor: 5000, Ativa: true, DataInicio: dataInicioPadraoHistorico},
		}},
		&MockRendaVariavelRepoForHistorico{rendas: []*domain.RendaVariavel{
			{ID: "rv-1", Descricao: "Freelance", MembroID: "m-1", Valor: 1200, MesReferencia: 3, AnoReferencia: 2026},
		}},
		&MockRendaExtraRepoForHistorico{rendas: []*domain.RendaExtra{
			{ID: "re-1", Descricao: "Bônus", MembroID: "m-1", Valor: 500, DataRecebimento: time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)},
		}},
		&MockRendimentoRepoForHistorico{rendimentos: []*domain.RendimentoInvestimento{
			// Valor total = 2000, mas apenas 400 é distribuído
			{ID: "ri-1", Descricao: "CDB", MembroID: "m-1", Valor: 2000, ValorDistribuido: 400, Data: time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC)},
		}},
	)

	resultado, err := svc.BuscarHistorico("familia-id", domain.FiltroHistoricoRendas{})
	require.NoError(t, err)

	assert.Equal(t, 5000.0, resultado.Resumo.TotalFixas)
	assert.Equal(t, 1200.0, resultado.Resumo.TotalVariaveis)
	assert.Equal(t, 500.0, resultado.Resumo.TotalExtras)
	// Deve usar ValorDistribuido (400), não Valor (2000)
	assert.Equal(t, 400.0, resultado.Resumo.TotalInvestimentosDistribuidos)
	// TotalGeral = 5000 + 1200 + 500 + 400 = 7100
	assert.Equal(t, 7100.0, resultado.Resumo.TotalGeral)
}

// TestHistorico_InvestimentoNaoDistribuidoNaoEntraNaTotalGeral verifica que rendimento com
// ValorDistribuido = 0 não afeta o TotalGeral.
func TestHistorico_InvestimentoNaoDistribuidoNaoEntraNaTotalGeral(t *testing.T) {
	svc := newHistoricoSvc(
		&MockRendaFixaRepoForHistorico{},
		&MockRendaVariavelRepoForHistorico{},
		&MockRendaExtraRepoForHistorico{},
		&MockRendimentoRepoForHistorico{rendimentos: []*domain.RendimentoInvestimento{
			{ID: "ri-1", Descricao: "CDB", MembroID: "m-1", Valor: 5000, ValorDistribuido: 0, Data: time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC)},
		}},
	)

	resultado, err := svc.BuscarHistorico("familia-id", domain.FiltroHistoricoRendas{})
	require.NoError(t, err)

	assert.Len(t, resultado.Itens, 1)
	assert.Equal(t, 5000.0, resultado.Itens[0].Valor) // Valor informativo preservado
	assert.Equal(t, 0.0, resultado.Resumo.TotalInvestimentosDistribuidos)
	assert.Equal(t, 0.0, resultado.Resumo.TotalGeral)
}

// TestHistorico_ErroPropagadoFixa verifica que erros do repositório de renda fixa são propagados.
func TestHistorico_ErroPropagadoFixa(t *testing.T) {
	erroEsperado := errors.New("falha no banco de dados")
	svc := newHistoricoSvc(
		&MockRendaFixaRepoForHistorico{err: erroEsperado},
		&MockRendaVariavelRepoForHistorico{},
		&MockRendaExtraRepoForHistorico{},
		&MockRendimentoRepoForHistorico{},
	)

	_, err := svc.BuscarHistorico("familia-id", domain.FiltroHistoricoRendas{})
	assert.ErrorIs(t, err, erroEsperado)
}

// TestHistorico_ErroPropagadoVariavel verifica que erros do repositório de renda variável são propagados.
func TestHistorico_ErroPropagadoVariavel(t *testing.T) {
	erroEsperado := errors.New("falha no banco de dados")
	svc := newHistoricoSvc(
		&MockRendaFixaRepoForHistorico{},
		&MockRendaVariavelRepoForHistorico{err: erroEsperado},
		&MockRendaExtraRepoForHistorico{},
		&MockRendimentoRepoForHistorico{},
	)

	_, err := svc.BuscarHistorico("familia-id", domain.FiltroHistoricoRendas{})
	assert.ErrorIs(t, err, erroEsperado)
}

// TestHistorico_ErroPropagadoExtra verifica que erros do repositório de renda extra são propagados.
func TestHistorico_ErroPropagadoExtra(t *testing.T) {
	erroEsperado := errors.New("falha no banco de dados")
	svc := newHistoricoSvc(
		&MockRendaFixaRepoForHistorico{},
		&MockRendaVariavelRepoForHistorico{},
		&MockRendaExtraRepoForHistorico{err: erroEsperado},
		&MockRendimentoRepoForHistorico{},
	)

	_, err := svc.BuscarHistorico("familia-id", domain.FiltroHistoricoRendas{})
	assert.ErrorIs(t, err, erroEsperado)
}

// TestHistorico_ErroPropagadoRendimento verifica que erros do repositório de rendimento são propagados.
func TestHistorico_ErroPropagadoRendimento(t *testing.T) {
	erroEsperado := errors.New("falha no banco de dados")
	svc := newHistoricoSvc(
		&MockRendaFixaRepoForHistorico{},
		&MockRendaVariavelRepoForHistorico{},
		&MockRendaExtraRepoForHistorico{},
		&MockRendimentoRepoForHistorico{err: erroEsperado},
	)

	_, err := svc.BuscarHistorico("familia-id", domain.FiltroHistoricoRendas{})
	assert.ErrorIs(t, err, erroEsperado)
}

// TestHistorico_SemResultados verifica que o historico vazio retorna itens e resumo zerados.
func TestHistorico_SemResultados(t *testing.T) {
	svc := newHistoricoSvcVazio()

	resultado, err := svc.BuscarHistorico("familia-id", domain.FiltroHistoricoRendas{})
	require.NoError(t, err)
	assert.Empty(t, resultado.Itens)
	assert.Equal(t, 0.0, resultado.Resumo.TotalFixas)
	assert.Equal(t, 0.0, resultado.Resumo.TotalVariaveis)
	assert.Equal(t, 0.0, resultado.Resumo.TotalExtras)
	assert.Equal(t, 0.0, resultado.Resumo.TotalInvestimentosDistribuidos)
	assert.Equal(t, 0.0, resultado.Resumo.TotalGeral)
}
