package service_test

import (
	"testing"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/felipehrs/financas-familiares/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockContaFixaRepository implementa ContaFixaRepository para testes unitários.
// Usa campos de controle para simular respostas do repositório.
type MockContaFixaRepository struct {
	returnContaFixa  *domain.ContaFixa
	returnContasFixas []*domain.ContaFixa
	returnError       error

	// rastreamento de chamadas
	criarChamado        bool
	alterarAtivoChamado bool
	atualizarChamado    bool
}

func (m *MockContaFixaRepository) Criar(c *domain.ContaFixa) (*domain.ContaFixa, error) {
	m.criarChamado = true
	if m.returnError != nil {
		return nil, m.returnError
	}
	criado := *c
	criado.ID = "uuid-gerado-mock"
	return &criado, nil
}

func (m *MockContaFixaRepository) BuscarPorID(id string) (*domain.ContaFixa, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnContaFixa, nil
}

func (m *MockContaFixaRepository) Listar() ([]*domain.ContaFixa, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnContasFixas, nil
}

func (m *MockContaFixaRepository) Atualizar(c *domain.ContaFixa) (*domain.ContaFixa, error) {
	m.atualizarChamado = true
	if m.returnError != nil {
		return nil, m.returnError
	}
	atualizado := *c
	return &atualizado, nil
}

func (m *MockContaFixaRepository) AlterarAtivo(id string, ativa bool) (*domain.ContaFixa, error) {
	m.alterarAtivoChamado = true
	if m.returnError != nil {
		return nil, m.returnError
	}
	c := *m.returnContaFixa
	c.Ativa = ativa
	return &c, nil
}

// ---- Testes de Criar ----

func TestCriarContaFixa_Sucesso(t *testing.T) {
	mock := &MockContaFixaRepository{}
	svc := service.NewContaFixaService(mock)

	conta, err := svc.Criar("Internet", "membro-1", nil, 150.00, 10, "debito")

	require.NoError(t, err)
	assert.Equal(t, "uuid-gerado-mock", conta.ID)
	assert.Equal(t, "Internet", conta.Descricao)
	assert.Equal(t, "membro-1", conta.MembroID)
	assert.Equal(t, 150.00, conta.Valor)
	assert.Equal(t, 10, conta.DiaVencimento)
	assert.Equal(t, "debito", conta.FormaPagamento)
	assert.True(t, conta.Ativa)
	assert.True(t, mock.criarChamado)
}

func TestCriarContaFixa_DescricaoVazia(t *testing.T) {
	mock := &MockContaFixaRepository{}
	svc := service.NewContaFixaService(mock)

	_, err := svc.Criar("", "membro-1", nil, 150.00, 10, "debito")

	assert.ErrorIs(t, err, domain.ErrDescricaoContaFixaObrigatoria)
	assert.False(t, mock.criarChamado, "repositório não deve ser chamado com descrição vazia")
}

func TestCriarContaFixa_MembroIDVazio(t *testing.T) {
	mock := &MockContaFixaRepository{}
	svc := service.NewContaFixaService(mock)

	_, err := svc.Criar("Internet", "", nil, 150.00, 10, "debito")

	assert.ErrorIs(t, err, domain.ErrMembroIDContaFixaObrigatorio)
	assert.False(t, mock.criarChamado)
}

func TestCriarContaFixa_ValorZero(t *testing.T) {
	mock := &MockContaFixaRepository{}
	svc := service.NewContaFixaService(mock)

	_, err := svc.Criar("Internet", "membro-1", nil, 0, 10, "debito")

	assert.ErrorIs(t, err, domain.ErrValorContaFixaInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarContaFixa_ValorNegativo(t *testing.T) {
	mock := &MockContaFixaRepository{}
	svc := service.NewContaFixaService(mock)

	_, err := svc.Criar("Internet", "membro-1", nil, -50, 10, "debito")

	assert.ErrorIs(t, err, domain.ErrValorContaFixaInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarContaFixa_DiaVencimento0(t *testing.T) {
	mock := &MockContaFixaRepository{}
	svc := service.NewContaFixaService(mock)

	_, err := svc.Criar("Internet", "membro-1", nil, 150.00, 0, "debito")

	assert.ErrorIs(t, err, domain.ErrDiaVencimentoContaFixaInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarContaFixa_DiaVencimento29(t *testing.T) {
	mock := &MockContaFixaRepository{}
	svc := service.NewContaFixaService(mock)

	_, err := svc.Criar("Internet", "membro-1", nil, 150.00, 29, "debito")

	assert.ErrorIs(t, err, domain.ErrDiaVencimentoContaFixaInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarContaFixa_FormaPagamentoVazia(t *testing.T) {
	mock := &MockContaFixaRepository{}
	svc := service.NewContaFixaService(mock)

	_, err := svc.Criar("Internet", "membro-1", nil, 150.00, 10, "")

	assert.ErrorIs(t, err, domain.ErrFormaPagamentoContaFixaObrigatoria)
	assert.False(t, mock.criarChamado)
}

func TestCriarContaFixa_AtivaDefaultTrue(t *testing.T) {
	mock := &MockContaFixaRepository{}
	svc := service.NewContaFixaService(mock)

	conta, err := svc.Criar("Água", "membro-1", nil, 80.00, 15, "boleto")

	require.NoError(t, err)
	assert.True(t, conta.Ativa)
}

// ---- Testes de Listar ----

func TestListarContasFixas_RetornaLista(t *testing.T) {
	contas := []*domain.ContaFixa{
		{ID: "uuid-1", Descricao: "Internet", Ativa: true},
		{ID: "uuid-2", Descricao: "Água", Ativa: true},
	}
	mock := &MockContaFixaRepository{returnContasFixas: contas}
	svc := service.NewContaFixaService(mock)

	resultado, err := svc.Listar()

	require.NoError(t, err)
	assert.Len(t, resultado, 2)
	assert.Equal(t, "uuid-1", resultado[0].ID)
	assert.Equal(t, "uuid-2", resultado[1].ID)
}

func TestListarContasFixas_ListaVazia(t *testing.T) {
	mock := &MockContaFixaRepository{returnContasFixas: []*domain.ContaFixa{}}
	svc := service.NewContaFixaService(mock)

	resultado, err := svc.Listar()

	require.NoError(t, err)
	assert.Empty(t, resultado)
}

// ---- Testes de AlterarAtivo ----

func TestAlterarAtivo_Sucesso(t *testing.T) {
	existente := &domain.ContaFixa{ID: "uuid-1", Descricao: "Internet", Ativa: true}
	mock := &MockContaFixaRepository{returnContaFixa: existente}
	svc := service.NewContaFixaService(mock)

	resultado, err := svc.AlterarAtivo("uuid-1", false)

	require.NoError(t, err)
	assert.False(t, resultado.Ativa)
	assert.True(t, mock.alterarAtivoChamado)
}

func TestAlterarAtivo_NaoEncontrada(t *testing.T) {
	mock := &MockContaFixaRepository{returnError: domain.ErrContaFixaNaoEncontrada}
	svc := service.NewContaFixaService(mock)

	_, err := svc.AlterarAtivo("uuid-inexistente", false)

	assert.ErrorIs(t, err, domain.ErrContaFixaNaoEncontrada)
	assert.False(t, mock.alterarAtivoChamado)
}

// ---- Testes de Atualizar ----

func TestAtualizarContaFixa_Sucesso(t *testing.T) {
	existente := &domain.ContaFixa{
		ID:             "uuid-1",
		Descricao:      "Internet",
		MembroID:       "membro-1",
		Valor:          150.00,
		DiaVencimento:  10,
		FormaPagamento: "debito",
		Ativa:          true,
	}
	mock := &MockContaFixaRepository{returnContaFixa: existente}
	svc := service.NewContaFixaService(mock)

	resultado, err := svc.Atualizar("uuid-1", "Internet Fibra", "membro-1", nil, 200.00, 10, "debito", true)

	require.NoError(t, err)
	assert.Equal(t, "Internet Fibra", resultado.Descricao)
	assert.Equal(t, 200.00, resultado.Valor)
	assert.True(t, mock.atualizarChamado)
}

func TestAtualizarContaFixa_NaoEncontrada(t *testing.T) {
	mock := &MockContaFixaRepository{returnError: domain.ErrContaFixaNaoEncontrada}
	svc := service.NewContaFixaService(mock)

	_, err := svc.Atualizar("uuid-inexistente", "Internet", "membro-1", nil, 150.00, 10, "debito", true)

	assert.ErrorIs(t, err, domain.ErrContaFixaNaoEncontrada)
	assert.False(t, mock.atualizarChamado)
}
