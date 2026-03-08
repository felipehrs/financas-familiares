package service_test

import (
	"testing"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/felipehrs/financas-familiares/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockRendaFixaRepository implementa RendaFixaRepository para testes unitários.
type MockRendaFixaRepository struct {
	returnRenda   *domain.RendaFixa
	returnRendas  []*domain.RendaFixa
	returnError   error

	criarChamado     bool
	atualizarChamado bool
	inativarChamado  bool
	rendaRecebida    *domain.RendaFixa
}

func (m *MockRendaFixaRepository) Criar(r *domain.RendaFixa) (*domain.RendaFixa, error) {
	m.criarChamado = true
	m.rendaRecebida = r
	if m.returnError != nil {
		return nil, m.returnError
	}
	criado := *r
	criado.ID = "uuid-gerado-mock"
	return &criado, nil
}

func (m *MockRendaFixaRepository) BuscarPorID(id string) (*domain.RendaFixa, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnRenda, nil
}

func (m *MockRendaFixaRepository) Listar() ([]*domain.RendaFixa, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnRendas, nil
}

func (m *MockRendaFixaRepository) ListarAtivas() ([]*domain.RendaFixa, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnRendas, nil
}

func (m *MockRendaFixaRepository) Atualizar(r *domain.RendaFixa) (*domain.RendaFixa, error) {
	m.atualizarChamado = true
	m.rendaRecebida = r
	if m.returnError != nil {
		return nil, m.returnError
	}
	atualizado := *r
	return &atualizado, nil
}

func (m *MockRendaFixaRepository) Inativar(id string) error {
	m.inativarChamado = true
	return m.returnError
}

// ---- Testes de Criar ----

func TestCriarRendaFixa_Sucesso(t *testing.T) {
	mock := &MockRendaFixaRepository{}
	svc := service.NewRendaFixaService(mock)

	renda, err := svc.Criar("Salário", "membro-1", 5000.00, 5)

	require.NoError(t, err)
	assert.Equal(t, "uuid-gerado-mock", renda.ID)
	assert.Equal(t, "Salário", renda.Descricao)
	assert.Equal(t, "membro-1", renda.MembroID)
	assert.Equal(t, 5000.00, renda.Valor)
	assert.Equal(t, 5, renda.DiaRecebimento)
	assert.True(t, renda.Ativa, "renda criada deve estar ativa por padrão")
	assert.True(t, mock.criarChamado)
}

func TestCriarRendaFixa_DescricaoVazia(t *testing.T) {
	mock := &MockRendaFixaRepository{}
	svc := service.NewRendaFixaService(mock)

	_, err := svc.Criar("", "membro-1", 5000.00, 5)

	assert.ErrorIs(t, err, domain.ErrDescricaoRendaObrigatoria)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendaFixa_DescricaoApenasEspacos(t *testing.T) {
	mock := &MockRendaFixaRepository{}
	svc := service.NewRendaFixaService(mock)

	_, err := svc.Criar("   ", "membro-1", 5000.00, 5)

	assert.ErrorIs(t, err, domain.ErrDescricaoRendaObrigatoria)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendaFixa_MembroIDVazio(t *testing.T) {
	mock := &MockRendaFixaRepository{}
	svc := service.NewRendaFixaService(mock)

	_, err := svc.Criar("Salário", "", 5000.00, 5)

	assert.ErrorIs(t, err, domain.ErrMembroIDRendaObrigatorio)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendaFixa_ValorZero(t *testing.T) {
	mock := &MockRendaFixaRepository{}
	svc := service.NewRendaFixaService(mock)

	_, err := svc.Criar("Salário", "membro-1", 0, 5)

	assert.ErrorIs(t, err, domain.ErrValorRendaInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendaFixa_ValorNegativo(t *testing.T) {
	mock := &MockRendaFixaRepository{}
	svc := service.NewRendaFixaService(mock)

	_, err := svc.Criar("Salário", "membro-1", -100, 5)

	assert.ErrorIs(t, err, domain.ErrValorRendaInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendaFixa_DiaRecebimentoZero(t *testing.T) {
	mock := &MockRendaFixaRepository{}
	svc := service.NewRendaFixaService(mock)

	_, err := svc.Criar("Salário", "membro-1", 5000.00, 0)

	assert.ErrorIs(t, err, domain.ErrDiaRecebimentoInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendaFixa_DiaRecebimentoAcimaDe31(t *testing.T) {
	mock := &MockRendaFixaRepository{}
	svc := service.NewRendaFixaService(mock)

	_, err := svc.Criar("Salário", "membro-1", 5000.00, 32)

	assert.ErrorIs(t, err, domain.ErrDiaRecebimentoInvalido)
	assert.False(t, mock.criarChamado)
}

// ---- Testes de BuscarPorID ----

func TestBuscarRendaFixa_Sucesso(t *testing.T) {
	esperado := &domain.RendaFixa{
		ID:             "uuid-1",
		Descricao:      "Salário",
		MembroID:       "membro-1",
		Valor:          5000.00,
		DiaRecebimento: 5,
		Ativa:          true,
	}
	mock := &MockRendaFixaRepository{returnRenda: esperado}
	svc := service.NewRendaFixaService(mock)

	renda, err := svc.BuscarPorID("uuid-1")

	require.NoError(t, err)
	assert.Equal(t, esperado.ID, renda.ID)
	assert.Equal(t, esperado.Descricao, renda.Descricao)
}

func TestBuscarRendaFixa_NaoEncontrada(t *testing.T) {
	mock := &MockRendaFixaRepository{returnError: domain.ErrRendaFixaNaoEncontrada}
	svc := service.NewRendaFixaService(mock)

	_, err := svc.BuscarPorID("uuid-inexistente")

	assert.ErrorIs(t, err, domain.ErrRendaFixaNaoEncontrada)
}

// ---- Testes de Listar ----

func TestListarRendasFixas_Sucesso(t *testing.T) {
	rendas := []*domain.RendaFixa{
		{ID: "uuid-1", Descricao: "Salário", Ativa: true},
		{ID: "uuid-2", Descricao: "Aposentadoria", Ativa: false},
	}
	mock := &MockRendaFixaRepository{returnRendas: rendas}
	svc := service.NewRendaFixaService(mock)

	resultado, err := svc.Listar()

	require.NoError(t, err)
	assert.Len(t, resultado, 2)
}

// ---- Testes de ListarAtivas ----

func TestListarRendasFixasAtivas_Sucesso(t *testing.T) {
	rendas := []*domain.RendaFixa{
		{ID: "uuid-1", Descricao: "Salário", Ativa: true},
	}
	mock := &MockRendaFixaRepository{returnRendas: rendas}
	svc := service.NewRendaFixaService(mock)

	resultado, err := svc.ListarAtivas()

	require.NoError(t, err)
	assert.Len(t, resultado, 1)
	assert.True(t, resultado[0].Ativa)
}

// ---- Testes de Atualizar ----

func TestAtualizarRendaFixa_Sucesso(t *testing.T) {
	existente := &domain.RendaFixa{
		ID:             "uuid-1",
		Descricao:      "Salário",
		MembroID:       "membro-1",
		Valor:          5000.00,
		DiaRecebimento: 5,
		Ativa:          true,
	}
	mock := &MockRendaFixaRepository{returnRenda: existente}
	svc := service.NewRendaFixaService(mock)

	atualizado, err := svc.Atualizar("uuid-1", "Salário Atualizado", "membro-1", 6000.00, 10, true)

	require.NoError(t, err)
	assert.Equal(t, "uuid-1", atualizado.ID)
	assert.Equal(t, "Salário Atualizado", atualizado.Descricao)
	assert.Equal(t, 6000.00, atualizado.Valor)
	assert.Equal(t, 10, atualizado.DiaRecebimento)
	assert.True(t, mock.atualizarChamado)
}

func TestAtualizarRendaFixa_NaoEncontrada(t *testing.T) {
	mock := &MockRendaFixaRepository{returnError: domain.ErrRendaFixaNaoEncontrada}
	svc := service.NewRendaFixaService(mock)

	_, err := svc.Atualizar("uuid-inexistente", "Salário", "membro-1", 5000.00, 5, true)

	assert.ErrorIs(t, err, domain.ErrRendaFixaNaoEncontrada)
}

func TestAtualizarRendaFixa_DescricaoVazia(t *testing.T) {
	existente := &domain.RendaFixa{ID: "uuid-1", Descricao: "Salário", MembroID: "membro-1", Valor: 5000.00, DiaRecebimento: 5, Ativa: true}
	mock := &MockRendaFixaRepository{returnRenda: existente}
	svc := service.NewRendaFixaService(mock)

	_, err := svc.Atualizar("uuid-1", "", "membro-1", 5000.00, 5, true)

	assert.ErrorIs(t, err, domain.ErrDescricaoRendaObrigatoria)
	assert.False(t, mock.atualizarChamado)
}

// ---- Testes de Inativar ----

func TestInativarRendaFixa_Sucesso(t *testing.T) {
	existente := &domain.RendaFixa{ID: "uuid-1", Descricao: "Salário", Ativa: true}
	mock := &MockRendaFixaRepository{returnRenda: existente}
	svc := service.NewRendaFixaService(mock)

	err := svc.Inativar("uuid-1")

	require.NoError(t, err)
	assert.True(t, mock.inativarChamado)
}

func TestInativarRendaFixa_NaoEncontrada(t *testing.T) {
	mock := &MockRendaFixaRepository{returnError: domain.ErrRendaFixaNaoEncontrada}
	svc := service.NewRendaFixaService(mock)

	err := svc.Inativar("uuid-inexistente")

	assert.ErrorIs(t, err, domain.ErrRendaFixaNaoEncontrada)
}
