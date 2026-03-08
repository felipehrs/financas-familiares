package service_test

import (
	"testing"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/felipehrs/financas-familiares/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockAssinaturaRepository implementa AssinaturaRepository para testes unitários.
// Usa campos de controle para simular respostas do repositório.
type MockAssinaturaRepository struct {
	returnAssinatura  *domain.Assinatura
	returnAssinaturas []*domain.Assinatura
	returnError       error

	// rastreamento de chamadas
	criarChamado         bool
	alterarStatusChamado bool
	atualizarChamado     bool
}

func (m *MockAssinaturaRepository) Criar(a *domain.Assinatura) (*domain.Assinatura, error) {
	m.criarChamado = true
	if m.returnError != nil {
		return nil, m.returnError
	}
	criado := *a
	criado.ID = "uuid-gerado-mock"
	return &criado, nil
}

func (m *MockAssinaturaRepository) BuscarPorID(id string) (*domain.Assinatura, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnAssinatura, nil
}

func (m *MockAssinaturaRepository) Listar() ([]*domain.Assinatura, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnAssinaturas, nil
}

func (m *MockAssinaturaRepository) Atualizar(a *domain.Assinatura) (*domain.Assinatura, error) {
	m.atualizarChamado = true
	if m.returnError != nil {
		return nil, m.returnError
	}
	atualizado := *a
	return &atualizado, nil
}

func (m *MockAssinaturaRepository) AlterarStatus(id, status string) (*domain.Assinatura, error) {
	m.alterarStatusChamado = true
	if m.returnError != nil {
		return nil, m.returnError
	}
	a := *m.returnAssinatura
	a.Status = status
	return &a, nil
}

// ---- Testes de Criar ----

func TestCriarAssinatura_Sucesso(t *testing.T) {
	mock := &MockAssinaturaRepository{}
	svc := service.NewAssinaturaService(mock)

	assinatura, err := svc.Criar("Netflix", "membro-1", nil, 39.90, 15, "cartao_credito")

	require.NoError(t, err)
	assert.Equal(t, "uuid-gerado-mock", assinatura.ID)
	assert.Equal(t, "Netflix", assinatura.Nome)
	assert.Equal(t, "membro-1", assinatura.MembroID)
	assert.Equal(t, 39.90, assinatura.Valor)
	assert.Equal(t, 15, assinatura.DiaCobranca)
	assert.Equal(t, "cartao_credito", assinatura.FormaPagamento)
	assert.Equal(t, "ativa", assinatura.Status)
	assert.True(t, mock.criarChamado)
}

func TestCriarAssinatura_NomeVazio(t *testing.T) {
	mock := &MockAssinaturaRepository{}
	svc := service.NewAssinaturaService(mock)

	_, err := svc.Criar("", "membro-1", nil, 39.90, 15, "cartao_credito")

	assert.ErrorIs(t, err, domain.ErrNomeAssinaturaObrigatorio)
	assert.False(t, mock.criarChamado, "repositório não deve ser chamado com nome vazio")
}

func TestCriarAssinatura_MembroIDVazio(t *testing.T) {
	mock := &MockAssinaturaRepository{}
	svc := service.NewAssinaturaService(mock)

	_, err := svc.Criar("Netflix", "", nil, 39.90, 15, "cartao_credito")

	assert.ErrorIs(t, err, domain.ErrMembroIDAssinaturaObrigatorio)
	assert.False(t, mock.criarChamado)
}

func TestCriarAssinatura_ValorZero(t *testing.T) {
	mock := &MockAssinaturaRepository{}
	svc := service.NewAssinaturaService(mock)

	_, err := svc.Criar("Netflix", "membro-1", nil, 0, 15, "cartao_credito")

	assert.ErrorIs(t, err, domain.ErrValorAssinaturaInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarAssinatura_ValorNegativo(t *testing.T) {
	mock := &MockAssinaturaRepository{}
	svc := service.NewAssinaturaService(mock)

	_, err := svc.Criar("Netflix", "membro-1", nil, -10, 15, "cartao_credito")

	assert.ErrorIs(t, err, domain.ErrValorAssinaturaInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarAssinatura_DiaCobranca0(t *testing.T) {
	mock := &MockAssinaturaRepository{}
	svc := service.NewAssinaturaService(mock)

	_, err := svc.Criar("Netflix", "membro-1", nil, 39.90, 0, "cartao_credito")

	assert.ErrorIs(t, err, domain.ErrDiaCobrancaInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarAssinatura_DiaCobranca32(t *testing.T) {
	mock := &MockAssinaturaRepository{}
	svc := service.NewAssinaturaService(mock)

	_, err := svc.Criar("Netflix", "membro-1", nil, 39.90, 32, "cartao_credito")

	assert.ErrorIs(t, err, domain.ErrDiaCobrancaInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarAssinatura_FormaPagamentoVazia(t *testing.T) {
	mock := &MockAssinaturaRepository{}
	svc := service.NewAssinaturaService(mock)

	_, err := svc.Criar("Netflix", "membro-1", nil, 39.90, 15, "")

	assert.ErrorIs(t, err, domain.ErrFormaPagamentoObrigatoria)
	assert.False(t, mock.criarChamado)
}

func TestCriarAssinatura_StatusDefaultAtiva(t *testing.T) {
	mock := &MockAssinaturaRepository{}
	svc := service.NewAssinaturaService(mock)

	assinatura, err := svc.Criar("Spotify", "membro-1", nil, 19.90, 10, "debito")

	require.NoError(t, err)
	assert.Equal(t, "ativa", assinatura.Status)
}

// ---- Testes de Listar ----

func TestListarAssinaturas_RetornaLista(t *testing.T) {
	assinaturas := []*domain.Assinatura{
		{ID: "uuid-1", Nome: "Netflix", Status: "ativa"},
		{ID: "uuid-2", Nome: "Spotify", Status: "ativa"},
	}
	mock := &MockAssinaturaRepository{returnAssinaturas: assinaturas}
	svc := service.NewAssinaturaService(mock)

	resultado, err := svc.Listar()

	require.NoError(t, err)
	assert.Len(t, resultado, 2)
	assert.Equal(t, "uuid-1", resultado[0].ID)
	assert.Equal(t, "uuid-2", resultado[1].ID)
}

func TestListarAssinaturas_ListaVazia(t *testing.T) {
	mock := &MockAssinaturaRepository{returnAssinaturas: []*domain.Assinatura{}}
	svc := service.NewAssinaturaService(mock)

	resultado, err := svc.Listar()

	require.NoError(t, err)
	assert.Empty(t, resultado)
}

// ---- Testes de AlterarStatus ----

func TestAlterarStatus_Sucesso(t *testing.T) {
	existente := &domain.Assinatura{ID: "uuid-1", Nome: "Netflix", Status: "ativa"}
	mock := &MockAssinaturaRepository{returnAssinatura: existente}
	svc := service.NewAssinaturaService(mock)

	resultado, err := svc.AlterarStatus("uuid-1", "pausada")

	require.NoError(t, err)
	assert.Equal(t, "pausada", resultado.Status)
	assert.True(t, mock.alterarStatusChamado)
}

func TestAlterarStatus_StatusInvalido(t *testing.T) {
	mock := &MockAssinaturaRepository{}
	svc := service.NewAssinaturaService(mock)

	_, err := svc.AlterarStatus("uuid-1", "invalido")

	assert.ErrorIs(t, err, domain.ErrStatusInvalido)
	assert.False(t, mock.alterarStatusChamado)
}

func TestAlterarStatus_NaoEncontrada(t *testing.T) {
	mock := &MockAssinaturaRepository{returnError: domain.ErrAssinaturaNaoEncontrada}
	svc := service.NewAssinaturaService(mock)

	_, err := svc.AlterarStatus("uuid-inexistente", "pausada")

	assert.ErrorIs(t, err, domain.ErrAssinaturaNaoEncontrada)
	assert.False(t, mock.alterarStatusChamado)
}

// ---- Testes de Atualizar ----

func TestAtualizarAssinatura_Sucesso(t *testing.T) {
	existente := &domain.Assinatura{
		ID:             "uuid-1",
		Nome:           "Netflix",
		MembroID:       "membro-1",
		Valor:          39.90,
		DiaCobranca:    15,
		FormaPagamento: "cartao_credito",
		Status:         "ativa",
	}
	mock := &MockAssinaturaRepository{returnAssinatura: existente}
	svc := service.NewAssinaturaService(mock)

	resultado, err := svc.Atualizar("uuid-1", "Netflix Premium", "membro-1", nil, 55.90, 15, "cartao_credito", "ativa")

	require.NoError(t, err)
	assert.Equal(t, "Netflix Premium", resultado.Nome)
	assert.Equal(t, 55.90, resultado.Valor)
	assert.True(t, mock.atualizarChamado)
}

func TestAtualizarAssinatura_NaoEncontrada(t *testing.T) {
	mock := &MockAssinaturaRepository{returnError: domain.ErrAssinaturaNaoEncontrada}
	svc := service.NewAssinaturaService(mock)

	_, err := svc.Atualizar("uuid-inexistente", "Netflix", "membro-1", nil, 39.90, 15, "cartao_credito", "ativa")

	assert.ErrorIs(t, err, domain.ErrAssinaturaNaoEncontrada)
	assert.False(t, mock.atualizarChamado)
}
