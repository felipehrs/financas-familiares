package service_test

import (
	"testing"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/felipehrs/financas-familiares/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockRendimentoInvestimentoRepository implementa RendimentoInvestimentoRepositoryInterface para testes unitários.
type MockRendimentoInvestimentoRepository struct {
	returnRendimento  *domain.RendimentoInvestimento
	returnRendimentos []*domain.RendimentoInvestimento
	returnError       error

	// rastreamento de chamadas
	criarChamado     bool
	excluirChamado   bool
	atualizarChamado bool
}

func (m *MockRendimentoInvestimentoRepository) Criar(familiaID string, r *domain.RendimentoInvestimento) (*domain.RendimentoInvestimento, error) {
	m.criarChamado = true
	if m.returnError != nil {
		return nil, m.returnError
	}
	criado := *r
	criado.ID = "uuid-gerado-mock"
	return &criado, nil
}

func (m *MockRendimentoInvestimentoRepository) BuscarPorID(familiaID, id string) (*domain.RendimentoInvestimento, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnRendimento, nil
}

func (m *MockRendimentoInvestimentoRepository) Listar(familiaID string) ([]*domain.RendimentoInvestimento, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnRendimentos, nil
}

func (m *MockRendimentoInvestimentoRepository) ListarPorMes(familiaID string, mes, ano int) ([]*domain.RendimentoInvestimento, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnRendimentos, nil
}

func (m *MockRendimentoInvestimentoRepository) Atualizar(familiaID string, r *domain.RendimentoInvestimento) (*domain.RendimentoInvestimento, error) {
	m.atualizarChamado = true
	if m.returnError != nil {
		return nil, m.returnError
	}
	atualizado := *r
	return &atualizado, nil
}

func (m *MockRendimentoInvestimentoRepository) Excluir(familiaID, id string) error {
	m.excluirChamado = true
	return m.returnError
}

// helpers
var dataRendimentoValida = time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC)

// ---- Testes de Criar ----

func TestCriarRendimentoInvestimento_SucessoSemValorDistribuido(t *testing.T) {
	mock := &MockRendimentoInvestimentoRepository{}
	svc := service.NewRendimentoInvestimentoService(mock)

	rendimento, err := svc.Criar("familia-id", "Rendimento CDB", "membro-1", dataRendimentoValida, 1000.00, 0)

	require.NoError(t, err)
	assert.Equal(t, "uuid-gerado-mock", rendimento.ID)
	assert.Equal(t, "membro-1", rendimento.MembroID)
	assert.Equal(t, "Rendimento CDB", rendimento.Descricao)
	assert.Equal(t, 1000.00, rendimento.Valor)
	assert.Equal(t, 0.0, rendimento.ValorDistribuido)
	assert.True(t, mock.criarChamado)
}

func TestCriarRendimentoInvestimento_SucessoComValorDistribuidoParcial(t *testing.T) {
	mock := &MockRendimentoInvestimentoRepository{}
	svc := service.NewRendimentoInvestimentoService(mock)

	rendimento, err := svc.Criar("familia-id", "Rendimento LCI", "membro-1", dataRendimentoValida, 2000.00, 500.00)

	require.NoError(t, err)
	assert.Equal(t, "uuid-gerado-mock", rendimento.ID)
	assert.Equal(t, 2000.00, rendimento.Valor)
	assert.Equal(t, 500.00, rendimento.ValorDistribuido)
	assert.True(t, mock.criarChamado)
}

func TestCriarRendimentoInvestimento_SemDescricao(t *testing.T) {
	mock := &MockRendimentoInvestimentoRepository{}
	svc := service.NewRendimentoInvestimentoService(mock)

	_, err := svc.Criar("familia-id", "", "membro-1", dataRendimentoValida, 1000.00, 0)

	assert.ErrorIs(t, err, domain.ErrDescricaoRendimentoObrigatoria)
	assert.False(t, mock.criarChamado, "repositório não deve ser chamado com descrição vazia")
}

func TestCriarRendimentoInvestimento_SemMembroID(t *testing.T) {
	mock := &MockRendimentoInvestimentoRepository{}
	svc := service.NewRendimentoInvestimentoService(mock)

	_, err := svc.Criar("familia-id", "Rendimento CDB", "", dataRendimentoValida, 1000.00, 0)

	assert.ErrorIs(t, err, domain.ErrMembroIDRendimentoObrigatorio)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendimentoInvestimento_ValorZero(t *testing.T) {
	mock := &MockRendimentoInvestimentoRepository{}
	svc := service.NewRendimentoInvestimentoService(mock)

	_, err := svc.Criar("familia-id", "Rendimento CDB", "membro-1", dataRendimentoValida, 0, 0)

	assert.ErrorIs(t, err, domain.ErrValorRendimentoInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendimentoInvestimento_ValorNegativo(t *testing.T) {
	mock := &MockRendimentoInvestimentoRepository{}
	svc := service.NewRendimentoInvestimentoService(mock)

	_, err := svc.Criar("familia-id", "Rendimento CDB", "membro-1", dataRendimentoValida, -100.00, 0)

	assert.ErrorIs(t, err, domain.ErrValorRendimentoInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendimentoInvestimento_DataZero(t *testing.T) {
	mock := &MockRendimentoInvestimentoRepository{}
	svc := service.NewRendimentoInvestimentoService(mock)

	_, err := svc.Criar("familia-id", "Rendimento CDB", "membro-1", time.Time{}, 1000.00, 0)

	assert.ErrorIs(t, err, domain.ErrDataRendimentoObrigatoria)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendimentoInvestimento_ValorDistribuidoMaiorQueValor(t *testing.T) {
	mock := &MockRendimentoInvestimentoRepository{}
	svc := service.NewRendimentoInvestimentoService(mock)

	_, err := svc.Criar("familia-id", "Rendimento CDB", "membro-1", dataRendimentoValida, 1000.00, 1500.00)

	assert.ErrorIs(t, err, domain.ErrValorDistribuidoInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendimentoInvestimento_ValorDistribuidoNegativo(t *testing.T) {
	mock := &MockRendimentoInvestimentoRepository{}
	svc := service.NewRendimentoInvestimentoService(mock)

	_, err := svc.Criar("familia-id", "Rendimento CDB", "membro-1", dataRendimentoValida, 1000.00, -50.00)

	assert.ErrorIs(t, err, domain.ErrValorDistribuidoInvalido)
	assert.False(t, mock.criarChamado)
}

// ---- Testes de BuscarPorID ----

func TestBuscarPorIDRendimentoInvestimento_Existente(t *testing.T) {
	existente := &domain.RendimentoInvestimento{
		ID:               "uuid-1",
		Descricao:        "Rendimento CDB",
		MembroID:         "membro-1",
		Data:             dataRendimentoValida,
		Valor:            1000.00,
		ValorDistribuido: 0,
	}
	mock := &MockRendimentoInvestimentoRepository{returnRendimento: existente}
	svc := service.NewRendimentoInvestimentoService(mock)

	resultado, err := svc.BuscarPorID("familia-id", "uuid-1")

	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resultado.ID)
	assert.Equal(t, "Rendimento CDB", resultado.Descricao)
}

func TestBuscarPorIDRendimentoInvestimento_Inexistente(t *testing.T) {
	mock := &MockRendimentoInvestimentoRepository{returnError: domain.ErrRendimentoInvestimentoNaoEncontrado}
	svc := service.NewRendimentoInvestimentoService(mock)

	_, err := svc.BuscarPorID("familia-id", "uuid-inexistente")

	assert.ErrorIs(t, err, domain.ErrRendimentoInvestimentoNaoEncontrado)
}

// ---- Testes de Listar ----

func TestListarRendimentosInvestimento_RetornaLista(t *testing.T) {
	rendimentos := []*domain.RendimentoInvestimento{
		{ID: "uuid-1", Descricao: "Rendimento CDB", Valor: 1000.00},
		{ID: "uuid-2", Descricao: "Rendimento LCI", Valor: 2000.00},
	}
	mock := &MockRendimentoInvestimentoRepository{returnRendimentos: rendimentos}
	svc := service.NewRendimentoInvestimentoService(mock)

	resultado, err := svc.Listar("familia-id")

	require.NoError(t, err)
	assert.Len(t, resultado, 2)
	assert.Equal(t, "uuid-1", resultado[0].ID)
	assert.Equal(t, "uuid-2", resultado[1].ID)
}

// ---- Testes de ListarPorMes ----

func TestListarRendimentosInvestimentoPorMes_RetornaListaFiltrada(t *testing.T) {
	rendimentos := []*domain.RendimentoInvestimento{
		{ID: "uuid-1", Descricao: "Rendimento CDB", Data: time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC)},
		{ID: "uuid-2", Descricao: "Rendimento LCI", Data: time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)},
	}
	mock := &MockRendimentoInvestimentoRepository{returnRendimentos: rendimentos}
	svc := service.NewRendimentoInvestimentoService(mock)

	resultado, err := svc.ListarPorMes("familia-id", 3, 2026)

	require.NoError(t, err)
	assert.Len(t, resultado, 2)
}

// ---- Testes de Excluir ----

func TestExcluirRendimentoInvestimento_Existente(t *testing.T) {
	existente := &domain.RendimentoInvestimento{ID: "uuid-1", Descricao: "Rendimento CDB"}
	mock := &MockRendimentoInvestimentoRepository{returnRendimento: existente}
	svc := service.NewRendimentoInvestimentoService(mock)

	err := svc.Excluir("familia-id", "uuid-1")

	require.NoError(t, err)
	assert.True(t, mock.excluirChamado)
}

func TestExcluirRendimentoInvestimento_Inexistente(t *testing.T) {
	mock := &MockRendimentoInvestimentoRepository{returnError: domain.ErrRendimentoInvestimentoNaoEncontrado}
	svc := service.NewRendimentoInvestimentoService(mock)

	err := svc.Excluir("familia-id", "uuid-inexistente")

	assert.ErrorIs(t, err, domain.ErrRendimentoInvestimentoNaoEncontrado)
	assert.False(t, mock.excluirChamado)
}
