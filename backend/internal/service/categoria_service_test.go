package service_test

import (
	"testing"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/felipehrs/financas-familiares/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockCategoriaRepository implementa CategoriaRepository para testes unitários.
// Usa campos de controle para simular respostas do repositório.
type MockCategoriaRepository struct {
	returnCategoria  *domain.Categoria
	returnCategorias []*domain.Categoria
	returnError      error

	// rastreamento de chamadas
	criarChamado      bool
	atualizarChamado  bool
	excluirChamado    bool
	categoriaRecebida *domain.Categoria
}

func (m *MockCategoriaRepository) Criar(familiaID string, categoria *domain.Categoria) (*domain.Categoria, error) {
	m.criarChamado = true
	m.categoriaRecebida = categoria
	if m.returnError != nil {
		return nil, m.returnError
	}
	criada := *categoria
	criada.ID = "uuid-gerado-mock"
	return &criada, nil
}

func (m *MockCategoriaRepository) Listar(familiaID string) ([]*domain.Categoria, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnCategorias, nil
}

func (m *MockCategoriaRepository) BuscarPorID(familiaID, id string) (*domain.Categoria, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnCategoria, nil
}

func (m *MockCategoriaRepository) Atualizar(familiaID string, categoria *domain.Categoria) (*domain.Categoria, error) {
	m.atualizarChamado = true
	m.categoriaRecebida = categoria
	if m.returnError != nil {
		return nil, m.returnError
	}
	atualizada := *categoria
	return &atualizada, nil
}

func (m *MockCategoriaRepository) Excluir(familiaID, id string) error {
	m.excluirChamado = true
	return m.returnError
}

// ---- Testes de Criar ----

func TestCriarCategoria_Sucesso(t *testing.T) {
	mock := &MockCategoriaRepository{}
	svc := service.NewCategoriaService(mock)

	cat, err := svc.Criar("familia-id", "Alimentação")

	require.NoError(t, err)
	assert.Equal(t, "uuid-gerado-mock", cat.ID)
	assert.Equal(t, "Alimentação", cat.Nome)
	assert.True(t, mock.criarChamado)
	assert.Equal(t, "Alimentação", mock.categoriaRecebida.Nome)
}

func TestCriarCategoria_NomeVazio(t *testing.T) {
	mock := &MockCategoriaRepository{}
	svc := service.NewCategoriaService(mock)

	_, err := svc.Criar("familia-id", "")

	assert.ErrorIs(t, err, domain.ErrNomeCategoriaObrigatorio)
	assert.False(t, mock.criarChamado, "repositório não deve ser chamado com nome vazio")
}

func TestCriarCategoria_NomeApenasEspacos(t *testing.T) {
	mock := &MockCategoriaRepository{}
	svc := service.NewCategoriaService(mock)

	_, err := svc.Criar("familia-id", "   ")

	assert.ErrorIs(t, err, domain.ErrNomeCategoriaObrigatorio)
	assert.False(t, mock.criarChamado, "repositório não deve ser chamado com nome apenas de espaços")
}

// ---- Testes de Listar ----

func TestListarCategorias_Sucesso(t *testing.T) {
	categorias := []*domain.Categoria{
		{ID: "uuid-1", Nome: "Alimentação"},
		{ID: "uuid-2", Nome: "Transporte"},
	}
	mock := &MockCategoriaRepository{returnCategorias: categorias}
	svc := service.NewCategoriaService(mock)

	resultado, err := svc.Listar("familia-id")

	require.NoError(t, err)
	assert.Len(t, resultado, 2)
	assert.Equal(t, "uuid-1", resultado[0].ID)
	assert.Equal(t, "uuid-2", resultado[1].ID)
}

func TestListarCategorias_ListaVazia(t *testing.T) {
	mock := &MockCategoriaRepository{returnCategorias: []*domain.Categoria{}}
	svc := service.NewCategoriaService(mock)

	resultado, err := svc.Listar("familia-id")

	require.NoError(t, err)
	assert.Empty(t, resultado)
}

// ---- Testes de Atualizar ----

func TestAtualizarCategoria_Sucesso(t *testing.T) {
	existente := &domain.Categoria{
		ID:   "uuid-1",
		Nome: "Alimentação",
	}
	mock := &MockCategoriaRepository{returnCategoria: existente}
	svc := service.NewCategoriaService(mock)

	atualizada, err := svc.Atualizar("familia-id", "uuid-1", "Alimentação e Bebidas")

	require.NoError(t, err)
	assert.Equal(t, "uuid-1", atualizada.ID)
	assert.Equal(t, "Alimentação e Bebidas", atualizada.Nome)
	assert.True(t, mock.atualizarChamado)
	assert.Equal(t, "Alimentação e Bebidas", mock.categoriaRecebida.Nome)
}

func TestAtualizarCategoria_NaoEncontrada(t *testing.T) {
	mock := &MockCategoriaRepository{returnError: domain.ErrCategoriaNaoEncontrada}
	svc := service.NewCategoriaService(mock)

	_, err := svc.Atualizar("familia-id", "uuid-inexistente", "Novo Nome")

	assert.ErrorIs(t, err, domain.ErrCategoriaNaoEncontrada)
	assert.False(t, mock.atualizarChamado, "Atualizar do repositório não deve ser chamado quando não encontrado")
}

func TestAtualizarCategoria_NomeVazio(t *testing.T) {
	existente := &domain.Categoria{ID: "uuid-1", Nome: "Alimentação"}
	mock := &MockCategoriaRepository{returnCategoria: existente}
	svc := service.NewCategoriaService(mock)

	_, err := svc.Atualizar("familia-id", "uuid-1", "")

	assert.ErrorIs(t, err, domain.ErrNomeCategoriaObrigatorio)
	assert.False(t, mock.atualizarChamado, "Atualizar do repositório não deve ser chamado com nome vazio")
}

// ---- Testes de Excluir ----

func TestExcluirCategoria_Sucesso(t *testing.T) {
	existente := &domain.Categoria{ID: "uuid-1", Nome: "Alimentação"}
	mock := &MockCategoriaRepository{returnCategoria: existente}
	svc := service.NewCategoriaService(mock)

	err := svc.Excluir("familia-id", "uuid-1")

	require.NoError(t, err)
	assert.True(t, mock.excluirChamado)
}

func TestExcluirCategoria_NaoEncontrada(t *testing.T) {
	mock := &MockCategoriaRepository{returnError: domain.ErrCategoriaNaoEncontrada}
	svc := service.NewCategoriaService(mock)

	err := svc.Excluir("familia-id", "uuid-inexistente")

	assert.ErrorIs(t, err, domain.ErrCategoriaNaoEncontrada)
	assert.False(t, mock.excluirChamado, "Excluir do repositório não deve ser chamado quando não encontrado")
}

func TestExcluirCategoria_ComVinculos(t *testing.T) {
	existente := &domain.Categoria{ID: "uuid-1", Nome: "Alimentação"}
	// BuscarPorID retorna sucesso, mas Excluir retorna ErrCategoriaComVinculos
	mock := &MockCategoriaRepository{
		returnCategoria: existente,
		returnError:     nil,
	}
	svc := service.NewCategoriaService(mock)

	// Substituir mock para que BuscarPorID funcione mas Excluir falhe
	mockComVinculos := &MockCategoriaRepositoryComVinculos{returnCategoria: existente}
	svcComVinculos := service.NewCategoriaService(mockComVinculos)

	err := svcComVinculos.Excluir("familia-id", "uuid-1")

	assert.ErrorIs(t, err, domain.ErrCategoriaComVinculos)
	assert.True(t, mockComVinculos.excluirChamado)

	// garantir que o mock simples não tem erro
	_ = svc
}

// MockCategoriaRepositoryComVinculos simula repositório onde Excluir retorna ErrCategoriaComVinculos.
type MockCategoriaRepositoryComVinculos struct {
	returnCategoria *domain.Categoria
	excluirChamado  bool
}

func (m *MockCategoriaRepositoryComVinculos) Criar(familiaID string, categoria *domain.Categoria) (*domain.Categoria, error) {
	return nil, nil
}

func (m *MockCategoriaRepositoryComVinculos) Listar(familiaID string) ([]*domain.Categoria, error) {
	return nil, nil
}

func (m *MockCategoriaRepositoryComVinculos) BuscarPorID(familiaID, id string) (*domain.Categoria, error) {
	return m.returnCategoria, nil
}

func (m *MockCategoriaRepositoryComVinculos) Atualizar(familiaID string, categoria *domain.Categoria) (*domain.Categoria, error) {
	return nil, nil
}

func (m *MockCategoriaRepositoryComVinculos) Excluir(familiaID, id string) error {
	m.excluirChamado = true
	return domain.ErrCategoriaComVinculos
}
