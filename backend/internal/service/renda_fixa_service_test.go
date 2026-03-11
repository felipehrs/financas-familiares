package service_test

import (
	"testing"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/felipehrs/financas-familiares/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockRendaFixaRepository implementa RendaFixaRepository para testes unitários.
type MockRendaFixaRepository struct {
	returnRenda  *domain.RendaFixa
	returnRendas []*domain.RendaFixa
	returnError  error

	criarChamado     bool
	atualizarChamado bool
	inativarChamado  bool
	rendaRecebida    *domain.RendaFixa
}

func (m *MockRendaFixaRepository) Criar(familiaID string, r *domain.RendaFixa) (*domain.RendaFixa, error) {
	m.criarChamado = true
	m.rendaRecebida = r
	if m.returnError != nil {
		return nil, m.returnError
	}
	criado := *r
	criado.ID = "uuid-gerado-mock"
	return &criado, nil
}

func (m *MockRendaFixaRepository) BuscarPorID(familiaID, id string) (*domain.RendaFixa, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnRenda, nil
}

func (m *MockRendaFixaRepository) Listar(familiaID string) ([]*domain.RendaFixa, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnRendas, nil
}

func (m *MockRendaFixaRepository) ListarAtivas(familiaID string) ([]*domain.RendaFixa, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnRendas, nil
}

func (m *MockRendaFixaRepository) ListarVigentesPorMes(familiaID string, mes, ano int) ([]*domain.RendaFixa, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnRendas, nil
}

func (m *MockRendaFixaRepository) Atualizar(familiaID string, r *domain.RendaFixa) (*domain.RendaFixa, error) {
	m.atualizarChamado = true
	m.rendaRecebida = r
	if m.returnError != nil {
		return nil, m.returnError
	}
	atualizado := *r
	return &atualizado, nil
}

func (m *MockRendaFixaRepository) Inativar(familiaID, id string) error {
	m.inativarChamado = true
	return m.returnError
}

// ---- Testes de Criar ----

func TestCriarRendaFixa_Sucesso(t *testing.T) {
	mock := &MockRendaFixaRepository{}
	svc := service.NewRendaFixaService(mock)

	dataInicio := time.Now()
	renda, err := svc.Criar("familia-id", "Salário", "membro-1", 5000.00, 5, dataInicio, nil)

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

	_, err := svc.Criar("familia-id", "", "membro-1", 5000.00, 5, time.Now(), nil)

	assert.ErrorIs(t, err, domain.ErrDescricaoRendaObrigatoria)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendaFixa_DescricaoApenasEspacos(t *testing.T) {
	mock := &MockRendaFixaRepository{}
	svc := service.NewRendaFixaService(mock)

	_, err := svc.Criar("familia-id", "   ", "membro-1", 5000.00, 5, time.Now(), nil)

	assert.ErrorIs(t, err, domain.ErrDescricaoRendaObrigatoria)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendaFixa_MembroIDVazio(t *testing.T) {
	mock := &MockRendaFixaRepository{}
	svc := service.NewRendaFixaService(mock)

	_, err := svc.Criar("familia-id", "Salário", "", 5000.00, 5, time.Now(), nil)

	assert.ErrorIs(t, err, domain.ErrMembroIDRendaObrigatorio)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendaFixa_ValorZero(t *testing.T) {
	mock := &MockRendaFixaRepository{}
	svc := service.NewRendaFixaService(mock)

	_, err := svc.Criar("familia-id", "Salário", "membro-1", 0, 5, time.Now(), nil)

	assert.ErrorIs(t, err, domain.ErrValorRendaInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendaFixa_ValorNegativo(t *testing.T) {
	mock := &MockRendaFixaRepository{}
	svc := service.NewRendaFixaService(mock)

	_, err := svc.Criar("familia-id", "Salário", "membro-1", -100, 5, time.Now(), nil)

	assert.ErrorIs(t, err, domain.ErrValorRendaInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendaFixa_DiaRecebimentoZero(t *testing.T) {
	mock := &MockRendaFixaRepository{}
	svc := service.NewRendaFixaService(mock)

	_, err := svc.Criar("familia-id", "Salário", "membro-1", 5000.00, 0, time.Now(), nil)

	assert.ErrorIs(t, err, domain.ErrDiaRecebimentoInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendaFixa_DiaRecebimentoAcimaDe31(t *testing.T) {
	mock := &MockRendaFixaRepository{}
	svc := service.NewRendaFixaService(mock)

	_, err := svc.Criar("familia-id", "Salário", "membro-1", 5000.00, 32, time.Now(), nil)

	assert.ErrorIs(t, err, domain.ErrDiaRecebimentoInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendaFixa_DataInicioZero(t *testing.T) {
	mock := &MockRendaFixaRepository{}
	svc := service.NewRendaFixaService(mock)

	_, err := svc.Criar("familia-id", "Salário", "membro-1", 5000.00, 5, time.Time{}, nil)

	assert.ErrorIs(t, err, domain.ErrDataInicioRendaObrigatoria)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendaFixa_DataFimNilPermitido(t *testing.T) {
	mock := &MockRendaFixaRepository{}
	svc := service.NewRendaFixaService(mock)

	renda, err := svc.Criar("familia-id", "Salário", "membro-1", 5000.00, 5, time.Now(), nil)

	require.NoError(t, err)
	assert.Nil(t, renda.DataFim)
	assert.True(t, mock.criarChamado)
}

func TestCriarRendaFixa_ComDataFim(t *testing.T) {
	mock := &MockRendaFixaRepository{}
	svc := service.NewRendaFixaService(mock)

	dataInicio := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	dataFim := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)

	renda, err := svc.Criar("familia-id", "Salário", "membro-1", 5000.00, 5, dataInicio, &dataFim)

	require.NoError(t, err)
	assert.True(t, mock.criarChamado)
	require.NotNil(t, renda.DataFim)
	assert.Equal(t, dataFim.Year(), renda.DataFim.Year())
	assert.Equal(t, dataFim.Month(), renda.DataFim.Month())
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

	renda, err := svc.BuscarPorID("familia-id", "uuid-1")

	require.NoError(t, err)
	assert.Equal(t, esperado.ID, renda.ID)
	assert.Equal(t, esperado.Descricao, renda.Descricao)
}

func TestBuscarRendaFixa_NaoEncontrada(t *testing.T) {
	mock := &MockRendaFixaRepository{returnError: domain.ErrRendaFixaNaoEncontrada}
	svc := service.NewRendaFixaService(mock)

	_, err := svc.BuscarPorID("familia-id", "uuid-inexistente")

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

	resultado, err := svc.Listar("familia-id")

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

	resultado, err := svc.ListarAtivas("familia-id")

	require.NoError(t, err)
	assert.Len(t, resultado, 1)
	assert.True(t, resultado[0].Ativa)
}

// ---- Testes de ListarVigentesPorMes ----

func TestListarVigentesPorMes_RetornaVigentes(t *testing.T) {
	rendas := []*domain.RendaFixa{
		{ID: "uuid-1", Descricao: "Salário", Ativa: true},
		{ID: "uuid-2", Descricao: "Freelance", Ativa: true},
	}
	mock := &MockRendaFixaRepository{returnRendas: rendas}
	svc := service.NewRendaFixaService(mock)

	resultado, err := svc.ListarVigentesPorMes("familia-id", 3, 2026)

	require.NoError(t, err)
	assert.Len(t, resultado, 2)
}

func TestListarVigentesPorMes_ListaVazia(t *testing.T) {
	mock := &MockRendaFixaRepository{returnRendas: []*domain.RendaFixa{}}
	svc := service.NewRendaFixaService(mock)

	resultado, err := svc.ListarVigentesPorMes("familia-id", 3, 2026)

	require.NoError(t, err)
	assert.Empty(t, resultado)
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

	atualizado, err := svc.Atualizar("familia-id", "uuid-1", "Salário Atualizado", "membro-1", 6000.00, 10, true, time.Now(), nil)

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

	_, err := svc.Atualizar("familia-id", "uuid-inexistente", "Salário", "membro-1", 5000.00, 5, true, time.Now(), nil)

	assert.ErrorIs(t, err, domain.ErrRendaFixaNaoEncontrada)
}

func TestAtualizarRendaFixa_DescricaoVazia(t *testing.T) {
	existente := &domain.RendaFixa{ID: "uuid-1", Descricao: "Salário", MembroID: "membro-1", Valor: 5000.00, DiaRecebimento: 5, Ativa: true}
	mock := &MockRendaFixaRepository{returnRenda: existente}
	svc := service.NewRendaFixaService(mock)

	_, err := svc.Atualizar("familia-id", "uuid-1", "", "membro-1", 5000.00, 5, true, time.Now(), nil)

	assert.ErrorIs(t, err, domain.ErrDescricaoRendaObrigatoria)
	assert.False(t, mock.atualizarChamado)
}

// ---- Testes de Inativar ----

func TestInativarRendaFixa_Sucesso(t *testing.T) {
	existente := &domain.RendaFixa{ID: "uuid-1", Descricao: "Salário", Ativa: true}
	mock := &MockRendaFixaRepository{returnRenda: existente}
	svc := service.NewRendaFixaService(mock)

	err := svc.Inativar("familia-id", "uuid-1")

	require.NoError(t, err)
	assert.True(t, mock.inativarChamado)
}

func TestInativarRendaFixa_NaoEncontrada(t *testing.T) {
	mock := &MockRendaFixaRepository{returnError: domain.ErrRendaFixaNaoEncontrada}
	svc := service.NewRendaFixaService(mock)

	err := svc.Inativar("familia-id", "uuid-inexistente")

	assert.ErrorIs(t, err, domain.ErrRendaFixaNaoEncontrada)
}
