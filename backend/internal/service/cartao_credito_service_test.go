package service_test

import (
	"testing"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/felipehrs/financas-familiares/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockCartaoCreditoRepository implementa CartaoCreditoRepository para testes unitários.
type MockCartaoCreditoRepository struct {
	returnCartao  *domain.CartaoCredito
	returnCartoes []*domain.CartaoCredito
	returnError   error

	criarChamado     bool
	atualizarChamado bool
	inativarChamado  bool
	cartaoRecebido   *domain.CartaoCredito
}

func (m *MockCartaoCreditoRepository) Criar(cartao *domain.CartaoCredito) (*domain.CartaoCredito, error) {
	m.criarChamado = true
	m.cartaoRecebido = cartao
	if m.returnError != nil {
		return nil, m.returnError
	}
	criado := *cartao
	criado.ID = "uuid-gerado-mock"
	return &criado, nil
}

func (m *MockCartaoCreditoRepository) BuscarPorID(id string) (*domain.CartaoCredito, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnCartao, nil
}

func (m *MockCartaoCreditoRepository) Listar() ([]*domain.CartaoCredito, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnCartoes, nil
}

func (m *MockCartaoCreditoRepository) Atualizar(cartao *domain.CartaoCredito) (*domain.CartaoCredito, error) {
	m.atualizarChamado = true
	m.cartaoRecebido = cartao
	if m.returnError != nil {
		return nil, m.returnError
	}
	atualizado := *cartao
	return &atualizado, nil
}

func (m *MockCartaoCreditoRepository) Inativar(id string) error {
	m.inativarChamado = true
	return m.returnError
}

func limitePtr(v float64) *float64 { return &v }

// ---- Testes de Criar ----

func TestCriarCartao_Sucesso(t *testing.T) {
	mock := &MockCartaoCreditoRepository{}
	svc := service.NewCartaoCreditoService(mock)

	cartao, err := svc.Criar("Nubank", "membro-1", 10, 17, limitePtr(5000))

	require.NoError(t, err)
	assert.Equal(t, "uuid-gerado-mock", cartao.ID)
	assert.Equal(t, "Nubank", cartao.Nome)
	assert.Equal(t, "membro-1", cartao.MembroID)
	assert.Equal(t, 10, cartao.DiaFechamento)
	assert.Equal(t, 17, cartao.DiaVencimento)
	assert.True(t, cartao.Ativo, "cartão criado deve estar ativo por padrão")
	assert.True(t, mock.criarChamado)
}

func TestCriarCartao_SemLimite(t *testing.T) {
	mock := &MockCartaoCreditoRepository{}
	svc := service.NewCartaoCreditoService(mock)

	cartao, err := svc.Criar("Inter", "membro-1", 5, 12, nil)

	require.NoError(t, err)
	assert.Nil(t, cartao.Limite)
}

func TestCriarCartao_NomeVazio(t *testing.T) {
	mock := &MockCartaoCreditoRepository{}
	svc := service.NewCartaoCreditoService(mock)

	_, err := svc.Criar("", "membro-1", 10, 17, nil)

	assert.ErrorIs(t, err, domain.ErrNomeCartaoObrigatorio)
	assert.False(t, mock.criarChamado)
}

func TestCriarCartao_NomeApenasEspacos(t *testing.T) {
	mock := &MockCartaoCreditoRepository{}
	svc := service.NewCartaoCreditoService(mock)

	_, err := svc.Criar("   ", "membro-1", 10, 17, nil)

	assert.ErrorIs(t, err, domain.ErrNomeCartaoObrigatorio)
	assert.False(t, mock.criarChamado)
}

func TestCriarCartao_MembroIDVazio(t *testing.T) {
	mock := &MockCartaoCreditoRepository{}
	svc := service.NewCartaoCreditoService(mock)

	_, err := svc.Criar("Nubank", "", 10, 17, nil)

	assert.ErrorIs(t, err, domain.ErrMembroIDObrigatorio)
	assert.False(t, mock.criarChamado)
}

func TestCriarCartao_DiaFechamentoInvalido(t *testing.T) {
	mock := &MockCartaoCreditoRepository{}
	svc := service.NewCartaoCreditoService(mock)

	_, err := svc.Criar("Nubank", "membro-1", 0, 17, nil)

	assert.ErrorIs(t, err, domain.ErrDiaFechamentoInvalido)

	_, err = svc.Criar("Nubank", "membro-1", 32, 17, nil)

	assert.ErrorIs(t, err, domain.ErrDiaFechamentoInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarCartao_DiaVencimentoInvalido(t *testing.T) {
	mock := &MockCartaoCreditoRepository{}
	svc := service.NewCartaoCreditoService(mock)

	_, err := svc.Criar("Nubank", "membro-1", 10, 0, nil)

	assert.ErrorIs(t, err, domain.ErrDiaVencimentoInvalido)

	_, err = svc.Criar("Nubank", "membro-1", 10, 32, nil)

	assert.ErrorIs(t, err, domain.ErrDiaVencimentoInvalido)
	assert.False(t, mock.criarChamado)
}

// ---- Testes de BuscarPorID ----

func TestBuscarCartao_Sucesso(t *testing.T) {
	esperado := &domain.CartaoCredito{
		ID:            "uuid-1",
		Nome:          "Nubank",
		MembroID:      "membro-1",
		DiaFechamento: 10,
		DiaVencimento: 17,
		Ativo:         true,
	}
	mock := &MockCartaoCreditoRepository{returnCartao: esperado}
	svc := service.NewCartaoCreditoService(mock)

	cartao, err := svc.BuscarPorID("uuid-1")

	require.NoError(t, err)
	assert.Equal(t, esperado.ID, cartao.ID)
	assert.Equal(t, esperado.Nome, cartao.Nome)
}

func TestBuscarCartao_NaoEncontrado(t *testing.T) {
	mock := &MockCartaoCreditoRepository{returnError: domain.ErrCartaoNaoEncontrado}
	svc := service.NewCartaoCreditoService(mock)

	_, err := svc.BuscarPorID("uuid-inexistente")

	assert.ErrorIs(t, err, domain.ErrCartaoNaoEncontrado)
}

// ---- Testes de Listar ----

func TestListarCartoes_Sucesso(t *testing.T) {
	cartoes := []*domain.CartaoCredito{
		{ID: "uuid-1", Nome: "Nubank", Ativo: true},
		{ID: "uuid-2", Nome: "Inter", Ativo: false},
	}
	mock := &MockCartaoCreditoRepository{returnCartoes: cartoes}
	svc := service.NewCartaoCreditoService(mock)

	resultado, err := svc.Listar()

	require.NoError(t, err)
	assert.Len(t, resultado, 2)
}

func TestListarCartoes_ListaVazia(t *testing.T) {
	mock := &MockCartaoCreditoRepository{returnCartoes: []*domain.CartaoCredito{}}
	svc := service.NewCartaoCreditoService(mock)

	resultado, err := svc.Listar()

	require.NoError(t, err)
	assert.Empty(t, resultado)
}

// ---- Testes de Atualizar ----

func TestAtualizarCartao_Sucesso(t *testing.T) {
	existente := &domain.CartaoCredito{
		ID:            "uuid-1",
		Nome:          "Nubank",
		MembroID:      "membro-1",
		DiaFechamento: 10,
		DiaVencimento: 17,
		Ativo:         true,
	}
	mock := &MockCartaoCreditoRepository{returnCartao: existente}
	svc := service.NewCartaoCreditoService(mock)

	atualizado, err := svc.Atualizar("uuid-1", "Nubank Gold", "membro-1", 15, 20, limitePtr(8000), true)

	require.NoError(t, err)
	assert.Equal(t, "uuid-1", atualizado.ID)
	assert.Equal(t, "Nubank Gold", atualizado.Nome)
	assert.Equal(t, 15, atualizado.DiaFechamento)
	assert.True(t, mock.atualizarChamado)
}

func TestAtualizarCartao_NaoEncontrado(t *testing.T) {
	mock := &MockCartaoCreditoRepository{returnError: domain.ErrCartaoNaoEncontrado}
	svc := service.NewCartaoCreditoService(mock)

	_, err := svc.Atualizar("uuid-inexistente", "Nubank", "membro-1", 10, 17, nil, true)

	assert.ErrorIs(t, err, domain.ErrCartaoNaoEncontrado)
}

func TestAtualizarCartao_NomeVazio(t *testing.T) {
	existente := &domain.CartaoCredito{ID: "uuid-1", Nome: "Nubank", MembroID: "membro-1", DiaFechamento: 10, DiaVencimento: 17, Ativo: true}
	mock := &MockCartaoCreditoRepository{returnCartao: existente}
	svc := service.NewCartaoCreditoService(mock)

	_, err := svc.Atualizar("uuid-1", "", "membro-1", 10, 17, nil, true)

	assert.ErrorIs(t, err, domain.ErrNomeCartaoObrigatorio)
	assert.False(t, mock.atualizarChamado)
}

// ---- Testes de Inativar ----

func TestInativarCartao_Sucesso(t *testing.T) {
	existente := &domain.CartaoCredito{ID: "uuid-1", Nome: "Nubank", Ativo: true}
	mock := &MockCartaoCreditoRepository{returnCartao: existente}
	svc := service.NewCartaoCreditoService(mock)

	err := svc.Inativar("uuid-1")

	require.NoError(t, err)
	assert.True(t, mock.inativarChamado)
}

func TestInativarCartao_NaoEncontrado(t *testing.T) {
	mock := &MockCartaoCreditoRepository{returnError: domain.ErrCartaoNaoEncontrado}
	svc := service.NewCartaoCreditoService(mock)

	err := svc.Inativar("uuid-inexistente")

	assert.ErrorIs(t, err, domain.ErrCartaoNaoEncontrado)
}
