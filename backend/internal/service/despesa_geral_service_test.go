package service_test

import (
	"testing"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/felipehrs/financas-familiares/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockDespesaGeralRepository implementa DespesaGeralRepositoryInterface para testes unitários.
type MockDespesaGeralRepository struct {
	returnDespesa  *domain.DespesaGeral
	returnDespesas []*domain.DespesaGeral
	returnError    error

	// rastreamento de chamadas
	criarChamado     bool
	excluirChamado   bool
	atualizarChamado bool
}

func (m *MockDespesaGeralRepository) Criar(familiaID string, d *domain.DespesaGeral) (*domain.DespesaGeral, error) {
	m.criarChamado = true
	if m.returnError != nil {
		return nil, m.returnError
	}
	criado := *d
	criado.ID = "uuid-gerado-mock"
	return &criado, nil
}

func (m *MockDespesaGeralRepository) BuscarPorID(familiaID, id string) (*domain.DespesaGeral, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnDespesa, nil
}

func (m *MockDespesaGeralRepository) Listar(familiaID string) ([]*domain.DespesaGeral, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnDespesas, nil
}

func (m *MockDespesaGeralRepository) ListarPorMes(familiaID string, mes, ano int) ([]*domain.DespesaGeral, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnDespesas, nil
}

func (m *MockDespesaGeralRepository) Atualizar(familiaID string, d *domain.DespesaGeral) (*domain.DespesaGeral, error) {
	m.atualizarChamado = true
	if m.returnError != nil {
		return nil, m.returnError
	}
	atualizado := *d
	return &atualizado, nil
}

func (m *MockDespesaGeralRepository) Excluir(familiaID, id string) error {
	m.excluirChamado = true
	return m.returnError
}

// helpers
var dataValida = time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC)

// ---- Testes de Criar ----

func TestCriarDespesaGeral_Sucesso(t *testing.T) {
	mock := &MockDespesaGeralRepository{}
	svc := service.NewDespesaGeralService(mock)

	despesa, err := svc.Criar("familia-id", "membro-1", nil, "Mercado", dataValida, 250.00, "pix", nil)

	require.NoError(t, err)
	assert.Equal(t, "uuid-gerado-mock", despesa.ID)
	assert.Equal(t, "membro-1", despesa.MembroID)
	assert.Equal(t, "Mercado", despesa.Descricao)
	assert.Equal(t, 250.00, despesa.Valor)
	assert.Equal(t, "pix", despesa.FormaPagamento)
	assert.True(t, mock.criarChamado)
}

func TestCriarDespesaGeral_SemDescricao(t *testing.T) {
	mock := &MockDespesaGeralRepository{}
	svc := service.NewDespesaGeralService(mock)

	_, err := svc.Criar("familia-id", "membro-1", nil, "", dataValida, 250.00, "pix", nil)

	assert.ErrorIs(t, err, domain.ErrDescricaoDespesaGeralObrigatoria)
	assert.False(t, mock.criarChamado, "repositório não deve ser chamado com descrição vazia")
}

func TestCriarDespesaGeral_SemMembroID(t *testing.T) {
	mock := &MockDespesaGeralRepository{}
	svc := service.NewDespesaGeralService(mock)

	_, err := svc.Criar("familia-id", "", nil, "Mercado", dataValida, 250.00, "pix", nil)

	assert.ErrorIs(t, err, domain.ErrMembroIDDespesaGeralObrigatorio)
	assert.False(t, mock.criarChamado)
}

func TestCriarDespesaGeral_ValorZero(t *testing.T) {
	mock := &MockDespesaGeralRepository{}
	svc := service.NewDespesaGeralService(mock)

	_, err := svc.Criar("familia-id", "membro-1", nil, "Mercado", dataValida, 0, "pix", nil)

	assert.ErrorIs(t, err, domain.ErrValorDespesaGeralInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarDespesaGeral_ValorNegativo(t *testing.T) {
	mock := &MockDespesaGeralRepository{}
	svc := service.NewDespesaGeralService(mock)

	_, err := svc.Criar("familia-id", "membro-1", nil, "Mercado", dataValida, -10.00, "pix", nil)

	assert.ErrorIs(t, err, domain.ErrValorDespesaGeralInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarDespesaGeral_SemData(t *testing.T) {
	mock := &MockDespesaGeralRepository{}
	svc := service.NewDespesaGeralService(mock)

	_, err := svc.Criar("familia-id", "membro-1", nil, "Mercado", time.Time{}, 250.00, "pix", nil)

	assert.ErrorIs(t, err, domain.ErrDataDespesaGeralObrigatoria)
	assert.False(t, mock.criarChamado)
}

func TestCriarDespesaGeral_SemFormaPagamento(t *testing.T) {
	mock := &MockDespesaGeralRepository{}
	svc := service.NewDespesaGeralService(mock)

	_, err := svc.Criar("familia-id", "membro-1", nil, "Mercado", dataValida, 250.00, "", nil)

	assert.ErrorIs(t, err, domain.ErrFormaPagamentoDespesaGeralObrigatoria)
	assert.False(t, mock.criarChamado)
}

// ---- Testes de BuscarPorID ----

func TestBuscarPorIDDespesaGeral_Existente(t *testing.T) {
	existente := &domain.DespesaGeral{
		ID:             "uuid-1",
		MembroID:       "membro-1",
		Descricao:      "Mercado",
		Data:           dataValida,
		Valor:          250.00,
		FormaPagamento: "pix",
	}
	mock := &MockDespesaGeralRepository{returnDespesa: existente}
	svc := service.NewDespesaGeralService(mock)

	resultado, err := svc.BuscarPorID("familia-id", "uuid-1")

	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resultado.ID)
	assert.Equal(t, "Mercado", resultado.Descricao)
}

func TestBuscarPorIDDespesaGeral_Inexistente(t *testing.T) {
	mock := &MockDespesaGeralRepository{returnError: domain.ErrDespesaGeralNaoEncontrada}
	svc := service.NewDespesaGeralService(mock)

	_, err := svc.BuscarPorID("familia-id", "uuid-inexistente")

	assert.ErrorIs(t, err, domain.ErrDespesaGeralNaoEncontrada)
}

// ---- Testes de Listar ----

func TestListarDespesasGerais_RetornaLista(t *testing.T) {
	despesas := []*domain.DespesaGeral{
		{ID: "uuid-1", Descricao: "Mercado", Valor: 250.00},
		{ID: "uuid-2", Descricao: "Farmácia", Valor: 80.00},
	}
	mock := &MockDespesaGeralRepository{returnDespesas: despesas}
	svc := service.NewDespesaGeralService(mock)

	resultado, err := svc.Listar("familia-id")

	require.NoError(t, err)
	assert.Len(t, resultado, 2)
	assert.Equal(t, "uuid-1", resultado[0].ID)
	assert.Equal(t, "uuid-2", resultado[1].ID)
}

// ---- Testes de ListarPorMes ----

func TestListarDespesasGeraisPorMes_RetornaApenasMesAno(t *testing.T) {
	despesas := []*domain.DespesaGeral{
		{ID: "uuid-1", Descricao: "Mercado", Data: time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC)},
		{ID: "uuid-2", Descricao: "Farmácia", Data: time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)},
	}
	mock := &MockDespesaGeralRepository{returnDespesas: despesas}
	svc := service.NewDespesaGeralService(mock)

	resultado, err := svc.ListarPorMes("familia-id", 3, 2026)

	require.NoError(t, err)
	assert.Len(t, resultado, 2)
}

// ---- Testes de Excluir ----

func TestExcluirDespesaGeral_Existente(t *testing.T) {
	existente := &domain.DespesaGeral{ID: "uuid-1", Descricao: "Mercado"}
	mock := &MockDespesaGeralRepository{returnDespesa: existente}
	svc := service.NewDespesaGeralService(mock)

	err := svc.Excluir("familia-id", "uuid-1")

	require.NoError(t, err)
	assert.True(t, mock.excluirChamado)
}

func TestExcluirDespesaGeral_Inexistente(t *testing.T) {
	mock := &MockDespesaGeralRepository{returnError: domain.ErrDespesaGeralNaoEncontrada}
	svc := service.NewDespesaGeralService(mock)

	err := svc.Excluir("familia-id", "uuid-inexistente")

	assert.ErrorIs(t, err, domain.ErrDespesaGeralNaoEncontrada)
	assert.False(t, mock.excluirChamado)
}
