package service_test

import (
	"testing"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/felipehrs/financas-familiares/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockRendaVariavelRepository implementa RendaVariavelRepositoryInterface para testes unitários.
type MockRendaVariavelRepository struct {
	returnRenda  *domain.RendaVariavel
	returnRendas []*domain.RendaVariavel
	returnError  error

	// rastreamento de chamadas
	criarChamado     bool
	excluirChamado   bool
	atualizarChamado bool
}

func (m *MockRendaVariavelRepository) Criar(r *domain.RendaVariavel) (*domain.RendaVariavel, error) {
	m.criarChamado = true
	if m.returnError != nil {
		return nil, m.returnError
	}
	criado := *r
	criado.ID = "uuid-gerado-mock"
	return &criado, nil
}

func (m *MockRendaVariavelRepository) BuscarPorID(id string) (*domain.RendaVariavel, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnRenda, nil
}

func (m *MockRendaVariavelRepository) Listar() ([]*domain.RendaVariavel, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnRendas, nil
}

func (m *MockRendaVariavelRepository) ListarPorMes(mes, ano int) ([]*domain.RendaVariavel, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnRendas, nil
}

func (m *MockRendaVariavelRepository) Atualizar(r *domain.RendaVariavel) (*domain.RendaVariavel, error) {
	m.atualizarChamado = true
	if m.returnError != nil {
		return nil, m.returnError
	}
	atualizado := *r
	return &atualizado, nil
}

func (m *MockRendaVariavelRepository) Excluir(id string) error {
	m.excluirChamado = true
	return m.returnError
}

// helpers
var dataRecebimentoValida = time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC)

// ---- Testes de Criar ----

func TestCriarRendaVariavel_Sucesso(t *testing.T) {
	mock := &MockRendaVariavelRepository{}
	svc := service.NewRendaVariavelService(mock)

	renda, err := svc.Criar("Freelance março", "membro-1", 3, 2026, 1500.00, dataRecebimentoValida)

	require.NoError(t, err)
	assert.Equal(t, "uuid-gerado-mock", renda.ID)
	assert.Equal(t, "membro-1", renda.MembroID)
	assert.Equal(t, "Freelance março", renda.Descricao)
	assert.Equal(t, 3, renda.MesReferencia)
	assert.Equal(t, 2026, renda.AnoReferencia)
	assert.Equal(t, 1500.00, renda.Valor)
	assert.True(t, mock.criarChamado)
}

func TestCriarRendaVariavel_SemDescricao(t *testing.T) {
	mock := &MockRendaVariavelRepository{}
	svc := service.NewRendaVariavelService(mock)

	_, err := svc.Criar("", "membro-1", 3, 2026, 1500.00, dataRecebimentoValida)

	assert.ErrorIs(t, err, domain.ErrDescricaoRendaVariavelObrigatoria)
	assert.False(t, mock.criarChamado, "repositório não deve ser chamado com descrição vazia")
}

func TestCriarRendaVariavel_SemMembroID(t *testing.T) {
	mock := &MockRendaVariavelRepository{}
	svc := service.NewRendaVariavelService(mock)

	_, err := svc.Criar("Freelance março", "", 3, 2026, 1500.00, dataRecebimentoValida)

	assert.ErrorIs(t, err, domain.ErrMembroIDRendaVariavelObrigatorio)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendaVariavel_MesReferenciaZero(t *testing.T) {
	mock := &MockRendaVariavelRepository{}
	svc := service.NewRendaVariavelService(mock)

	_, err := svc.Criar("Freelance março", "membro-1", 0, 2026, 1500.00, dataRecebimentoValida)

	assert.ErrorIs(t, err, domain.ErrMesReferenciaRendaVariavelInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendaVariavel_MesReferenciaTreze(t *testing.T) {
	mock := &MockRendaVariavelRepository{}
	svc := service.NewRendaVariavelService(mock)

	_, err := svc.Criar("Freelance março", "membro-1", 13, 2026, 1500.00, dataRecebimentoValida)

	assert.ErrorIs(t, err, domain.ErrMesReferenciaRendaVariavelInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendaVariavel_AnoReferenciaZero(t *testing.T) {
	mock := &MockRendaVariavelRepository{}
	svc := service.NewRendaVariavelService(mock)

	_, err := svc.Criar("Freelance março", "membro-1", 3, 0, 1500.00, dataRecebimentoValida)

	assert.ErrorIs(t, err, domain.ErrAnoReferenciaRendaVariavelInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendaVariavel_ValorZero(t *testing.T) {
	mock := &MockRendaVariavelRepository{}
	svc := service.NewRendaVariavelService(mock)

	_, err := svc.Criar("Freelance março", "membro-1", 3, 2026, 0, dataRecebimentoValida)

	assert.ErrorIs(t, err, domain.ErrValorRendaVariavelInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendaVariavel_ValorNegativo(t *testing.T) {
	mock := &MockRendaVariavelRepository{}
	svc := service.NewRendaVariavelService(mock)

	_, err := svc.Criar("Freelance março", "membro-1", 3, 2026, -100.00, dataRecebimentoValida)

	assert.ErrorIs(t, err, domain.ErrValorRendaVariavelInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendaVariavel_DataRecebimentoZero(t *testing.T) {
	mock := &MockRendaVariavelRepository{}
	svc := service.NewRendaVariavelService(mock)

	_, err := svc.Criar("Freelance março", "membro-1", 3, 2026, 1500.00, time.Time{})

	assert.ErrorIs(t, err, domain.ErrDataRecebimentoRendaVariavelObrigatoria)
	assert.False(t, mock.criarChamado)
}

// ---- Testes de BuscarPorID ----

func TestBuscarPorIDRendaVariavel_Existente(t *testing.T) {
	existente := &domain.RendaVariavel{
		ID:              "uuid-1",
		Descricao:       "Freelance março",
		MembroID:        "membro-1",
		MesReferencia:   3,
		AnoReferencia:   2026,
		Valor:           1500.00,
		DataRecebimento: dataRecebimentoValida,
	}
	mock := &MockRendaVariavelRepository{returnRenda: existente}
	svc := service.NewRendaVariavelService(mock)

	resultado, err := svc.BuscarPorID("uuid-1")

	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resultado.ID)
	assert.Equal(t, "Freelance março", resultado.Descricao)
}

func TestBuscarPorIDRendaVariavel_Inexistente(t *testing.T) {
	mock := &MockRendaVariavelRepository{returnError: domain.ErrRendaVariavelNaoEncontrada}
	svc := service.NewRendaVariavelService(mock)

	_, err := svc.BuscarPorID("uuid-inexistente")

	assert.ErrorIs(t, err, domain.ErrRendaVariavelNaoEncontrada)
}

// ---- Testes de Listar ----

func TestListarRendasVariaveis_RetornaLista(t *testing.T) {
	rendas := []*domain.RendaVariavel{
		{ID: "uuid-1", Descricao: "Freelance março", Valor: 1500.00},
		{ID: "uuid-2", Descricao: "Comissão venda", Valor: 800.00},
	}
	mock := &MockRendaVariavelRepository{returnRendas: rendas}
	svc := service.NewRendaVariavelService(mock)

	resultado, err := svc.Listar()

	require.NoError(t, err)
	assert.Len(t, resultado, 2)
	assert.Equal(t, "uuid-1", resultado[0].ID)
	assert.Equal(t, "uuid-2", resultado[1].ID)
}

// ---- Testes de ListarPorMes ----

func TestListarRendasVariaveisPorMes_RetornaListaFiltrada(t *testing.T) {
	rendas := []*domain.RendaVariavel{
		{ID: "uuid-1", Descricao: "Freelance março", MesReferencia: 3, AnoReferencia: 2026},
		{ID: "uuid-2", Descricao: "Comissão março", MesReferencia: 3, AnoReferencia: 2026},
	}
	mock := &MockRendaVariavelRepository{returnRendas: rendas}
	svc := service.NewRendaVariavelService(mock)

	resultado, err := svc.ListarPorMes(3, 2026)

	require.NoError(t, err)
	assert.Len(t, resultado, 2)
}

// ---- Testes de Excluir ----

func TestExcluirRendaVariavel_Existente(t *testing.T) {
	existente := &domain.RendaVariavel{ID: "uuid-1", Descricao: "Freelance março"}
	mock := &MockRendaVariavelRepository{returnRenda: existente}
	svc := service.NewRendaVariavelService(mock)

	err := svc.Excluir("uuid-1")

	require.NoError(t, err)
	assert.True(t, mock.excluirChamado)
}

func TestExcluirRendaVariavel_Inexistente(t *testing.T) {
	mock := &MockRendaVariavelRepository{returnError: domain.ErrRendaVariavelNaoEncontrada}
	svc := service.NewRendaVariavelService(mock)

	err := svc.Excluir("uuid-inexistente")

	assert.ErrorIs(t, err, domain.ErrRendaVariavelNaoEncontrada)
	assert.False(t, mock.excluirChamado)
}
