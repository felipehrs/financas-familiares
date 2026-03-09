package service_test

import (
	"testing"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/felipehrs/financas-familiares/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockMembroRepository implementa MembroRepository para testes unitários.
// Usa campos de controle para simular respostas do repositório.
type MockMembroRepository struct {
	returnMembro  *domain.Membro
	returnMembros []*domain.Membro
	returnError   error

	// rastreamento de chamadas
	criarChamado     bool
	inativarChamado  bool
	atualizarChamado bool
	membroRecebido   *domain.Membro
}

func (m *MockMembroRepository) Criar(membro *domain.Membro) (*domain.Membro, error) {
	m.criarChamado = true
	m.membroRecebido = membro
	if m.returnError != nil {
		return nil, m.returnError
	}
	// simular ID gerado pelo banco
	criado := *membro
	criado.ID = "uuid-gerado-mock"
	return &criado, nil
}

func (m *MockMembroRepository) BuscarPorID(id string) (*domain.Membro, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnMembro, nil
}

func (m *MockMembroRepository) Listar() ([]*domain.Membro, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnMembros, nil
}

func (m *MockMembroRepository) Atualizar(membro *domain.Membro) (*domain.Membro, error) {
	m.atualizarChamado = true
	m.membroRecebido = membro
	if m.returnError != nil {
		return nil, m.returnError
	}
	atualizado := *membro
	return &atualizado, nil
}

func (m *MockMembroRepository) Inativar(id string) error {
	m.inativarChamado = true
	return m.returnError
}

// ---- Testes de Criar ----

func TestCriarMembro_Sucesso(t *testing.T) {
	mock := &MockMembroRepository{}
	svc := service.NewMembroService(mock)

	membro, err := svc.Criar("Ana Lima", "cônjuge")

	require.NoError(t, err)
	assert.Equal(t, "uuid-gerado-mock", membro.ID)
	assert.Equal(t, "Ana Lima", membro.Nome)
	assert.Equal(t, "cônjuge", membro.Relacionamento)
	assert.True(t, membro.Ativo, "membro criado deve estar ativo por padrão")
	assert.True(t, mock.criarChamado)
	assert.True(t, mock.membroRecebido.Ativo, "deve enviar Ativo=true ao repositório")
}

func TestCriarMembro_NomeVazio(t *testing.T) {
	mock := &MockMembroRepository{}
	svc := service.NewMembroService(mock)

	_, err := svc.Criar("", "cônjuge")

	assert.ErrorIs(t, err, domain.ErrNomeObrigatorio)
	assert.False(t, mock.criarChamado, "repositório não deve ser chamado com nome vazio")
}

func TestCriarMembro_NomeApenasEspacos(t *testing.T) {
	mock := &MockMembroRepository{}
	svc := service.NewMembroService(mock)

	_, err := svc.Criar("   ", "cônjuge")

	assert.ErrorIs(t, err, domain.ErrNomeObrigatorio)
	assert.False(t, mock.criarChamado)
}

// ---- Testes de BuscarPorID ----

func TestBuscarMembro_Sucesso(t *testing.T) {
	esperado := &domain.Membro{
		ID:             "uuid-1",
		Nome:           "Carlos Souza",
		Relacionamento: "filho",
		Ativo:          true,
	}
	mock := &MockMembroRepository{returnMembro: esperado}
	svc := service.NewMembroService(mock)

	membro, err := svc.BuscarPorID("uuid-1")

	require.NoError(t, err)
	assert.Equal(t, esperado.ID, membro.ID)
	assert.Equal(t, esperado.Nome, membro.Nome)
}

func TestBuscarMembro_NaoEncontrado(t *testing.T) {
	mock := &MockMembroRepository{returnError: domain.ErrMembroNaoEncontrado}
	svc := service.NewMembroService(mock)

	_, err := svc.BuscarPorID("uuid-inexistente")

	assert.ErrorIs(t, err, domain.ErrMembroNaoEncontrado)
}

// ---- Testes de Listar ----

func TestListarMembros_Sucesso(t *testing.T) {
	membros := []*domain.Membro{
		{ID: "uuid-1", Nome: "Ana Lima", Ativo: true},
		{ID: "uuid-2", Nome: "Carlos Souza", Ativo: false},
	}
	mock := &MockMembroRepository{returnMembros: membros}
	svc := service.NewMembroService(mock)

	resultado, err := svc.Listar()

	require.NoError(t, err)
	assert.Len(t, resultado, 2)
	assert.Equal(t, "uuid-1", resultado[0].ID)
	assert.Equal(t, "uuid-2", resultado[1].ID)
}

func TestListarMembros_ListaVazia(t *testing.T) {
	mock := &MockMembroRepository{returnMembros: []*domain.Membro{}}
	svc := service.NewMembroService(mock)

	resultado, err := svc.Listar()

	require.NoError(t, err)
	assert.Empty(t, resultado)
}

// ---- Testes de Atualizar ----

func TestAtualizarMembro_Sucesso(t *testing.T) {
	existente := &domain.Membro{
		ID:             "uuid-1",
		Nome:           "Ana Lima",
		Relacionamento: "cônjuge",
		Ativo:          true,
	}
	mock := &MockMembroRepository{returnMembro: existente}
	svc := service.NewMembroService(mock)

	atualizado, err := svc.Atualizar("uuid-1", "Ana Souza", "esposa", true)

	require.NoError(t, err)
	assert.Equal(t, "uuid-1", atualizado.ID)
	assert.Equal(t, "Ana Souza", atualizado.Nome)
	assert.Equal(t, "esposa", atualizado.Relacionamento)
	assert.True(t, mock.atualizarChamado)
}

func TestAtualizarMembro_NaoEncontrado(t *testing.T) {
	mock := &MockMembroRepository{returnError: domain.ErrMembroNaoEncontrado}
	svc := service.NewMembroService(mock)

	_, err := svc.Atualizar("uuid-inexistente", "Nome Qualquer", "", true)

	assert.ErrorIs(t, err, domain.ErrMembroNaoEncontrado)
}

func TestAtualizarMembro_NomeVazio(t *testing.T) {
	existente := &domain.Membro{ID: "uuid-1", Nome: "Ana Lima", Ativo: true}
	mock := &MockMembroRepository{returnMembro: existente}
	svc := service.NewMembroService(mock)

	_, err := svc.Atualizar("uuid-1", "", "", true)

	assert.ErrorIs(t, err, domain.ErrNomeObrigatorio)
	assert.False(t, mock.atualizarChamado)
}

// ---- Testes de Inativar ----

func TestInativarMembro_Sucesso(t *testing.T) {
	existente := &domain.Membro{ID: "uuid-1", Nome: "Ana Lima", Ativo: true}
	mock := &MockMembroRepository{returnMembro: existente}
	svc := service.NewMembroService(mock)

	err := svc.Inativar("uuid-1")

	require.NoError(t, err)
	assert.True(t, mock.inativarChamado)
}

func TestInativarMembro_NaoEncontrado(t *testing.T) {
	mock := &MockMembroRepository{returnError: domain.ErrMembroNaoEncontrado}
	svc := service.NewMembroService(mock)

	err := svc.Inativar("uuid-inexistente")

	assert.ErrorIs(t, err, domain.ErrMembroNaoEncontrado)
}
