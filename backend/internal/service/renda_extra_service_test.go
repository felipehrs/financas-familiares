package service_test

import (
	"testing"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/felipehrs/financas-familiares/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockRendaExtraRepository implementa RendaExtraRepositoryInterface para testes unitários.
type MockRendaExtraRepository struct {
	returnRenda  *domain.RendaExtra
	returnRendas []*domain.RendaExtra
	returnError  error

	// rastreamento de chamadas
	criarChamado     bool
	excluirChamado   bool
	atualizarChamado bool
}

func (m *MockRendaExtraRepository) Criar(r *domain.RendaExtra) (*domain.RendaExtra, error) {
	m.criarChamado = true
	if m.returnError != nil {
		return nil, m.returnError
	}
	criado := *r
	criado.ID = "uuid-gerado-mock"
	return &criado, nil
}

func (m *MockRendaExtraRepository) BuscarPorID(id string) (*domain.RendaExtra, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnRenda, nil
}

func (m *MockRendaExtraRepository) Listar() ([]*domain.RendaExtra, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnRendas, nil
}

func (m *MockRendaExtraRepository) ListarPorMes(mes, ano int) ([]*domain.RendaExtra, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnRendas, nil
}

func (m *MockRendaExtraRepository) Atualizar(r *domain.RendaExtra) (*domain.RendaExtra, error) {
	m.atualizarChamado = true
	if m.returnError != nil {
		return nil, m.returnError
	}
	atualizado := *r
	return &atualizado, nil
}

func (m *MockRendaExtraRepository) Excluir(id string) error {
	m.excluirChamado = true
	return m.returnError
}

// helpers
var dataRecebimentoRendaExtraValida = time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC)

// ---- Testes de Criar ----

func TestCriarRendaExtra_Sucesso(t *testing.T) {
	mock := &MockRendaExtraRepository{}
	svc := service.NewRendaExtraService(mock)

	renda, err := svc.Criar("Bônus março", "membro-1", dataRecebimentoRendaExtraValida, 2000.00)

	require.NoError(t, err)
	assert.Equal(t, "uuid-gerado-mock", renda.ID)
	assert.Equal(t, "membro-1", renda.MembroID)
	assert.Equal(t, "Bônus março", renda.Descricao)
	assert.Equal(t, 2000.00, renda.Valor)
	assert.True(t, mock.criarChamado)
}

func TestCriarRendaExtra_SemDescricao(t *testing.T) {
	mock := &MockRendaExtraRepository{}
	svc := service.NewRendaExtraService(mock)

	_, err := svc.Criar("", "membro-1", dataRecebimentoRendaExtraValida, 2000.00)

	assert.ErrorIs(t, err, domain.ErrDescricaoRendaExtraObrigatoria)
	assert.False(t, mock.criarChamado, "repositório não deve ser chamado com descrição vazia")
}

func TestCriarRendaExtra_SemMembroID(t *testing.T) {
	mock := &MockRendaExtraRepository{}
	svc := service.NewRendaExtraService(mock)

	_, err := svc.Criar("Bônus março", "", dataRecebimentoRendaExtraValida, 2000.00)

	assert.ErrorIs(t, err, domain.ErrMembroIDRendaExtraObrigatorio)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendaExtra_ValorZero(t *testing.T) {
	mock := &MockRendaExtraRepository{}
	svc := service.NewRendaExtraService(mock)

	_, err := svc.Criar("Bônus março", "membro-1", dataRecebimentoRendaExtraValida, 0)

	assert.ErrorIs(t, err, domain.ErrValorRendaExtraInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendaExtra_ValorNegativo(t *testing.T) {
	mock := &MockRendaExtraRepository{}
	svc := service.NewRendaExtraService(mock)

	_, err := svc.Criar("Bônus março", "membro-1", dataRecebimentoRendaExtraValida, -100.00)

	assert.ErrorIs(t, err, domain.ErrValorRendaExtraInvalido)
	assert.False(t, mock.criarChamado)
}

func TestCriarRendaExtra_DataRecebimentoZero(t *testing.T) {
	mock := &MockRendaExtraRepository{}
	svc := service.NewRendaExtraService(mock)

	_, err := svc.Criar("Bônus março", "membro-1", time.Time{}, 2000.00)

	assert.ErrorIs(t, err, domain.ErrDataRecebimentoRendaExtraObrigatoria)
	assert.False(t, mock.criarChamado)
}

// ---- Testes de BuscarPorID ----

func TestBuscarPorIDRendaExtra_Existente(t *testing.T) {
	existente := &domain.RendaExtra{
		ID:              "uuid-1",
		Descricao:       "Bônus março",
		MembroID:        "membro-1",
		DataRecebimento: dataRecebimentoRendaExtraValida,
		Valor:           2000.00,
	}
	mock := &MockRendaExtraRepository{returnRenda: existente}
	svc := service.NewRendaExtraService(mock)

	resultado, err := svc.BuscarPorID("uuid-1")

	require.NoError(t, err)
	assert.Equal(t, "uuid-1", resultado.ID)
	assert.Equal(t, "Bônus março", resultado.Descricao)
}

func TestBuscarPorIDRendaExtra_Inexistente(t *testing.T) {
	mock := &MockRendaExtraRepository{returnError: domain.ErrRendaExtraNaoEncontrada}
	svc := service.NewRendaExtraService(mock)

	_, err := svc.BuscarPorID("uuid-inexistente")

	assert.ErrorIs(t, err, domain.ErrRendaExtraNaoEncontrada)
}

// ---- Testes de Listar ----

func TestListarRendasExtras_RetornaLista(t *testing.T) {
	rendas := []*domain.RendaExtra{
		{ID: "uuid-1", Descricao: "Bônus março", Valor: 2000.00},
		{ID: "uuid-2", Descricao: "Venda equipamento", Valor: 500.00},
	}
	mock := &MockRendaExtraRepository{returnRendas: rendas}
	svc := service.NewRendaExtraService(mock)

	resultado, err := svc.Listar()

	require.NoError(t, err)
	assert.Len(t, resultado, 2)
	assert.Equal(t, "uuid-1", resultado[0].ID)
	assert.Equal(t, "uuid-2", resultado[1].ID)
}

// ---- Testes de ListarPorMes ----

func TestListarRendasExtrasPorMes_RetornaListaFiltrada(t *testing.T) {
	rendas := []*domain.RendaExtra{
		{ID: "uuid-1", Descricao: "Bônus março", DataRecebimento: time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC)},
		{ID: "uuid-2", Descricao: "Venda março", DataRecebimento: time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)},
	}
	mock := &MockRendaExtraRepository{returnRendas: rendas}
	svc := service.NewRendaExtraService(mock)

	resultado, err := svc.ListarPorMes(3, 2026)

	require.NoError(t, err)
	assert.Len(t, resultado, 2)
}

// ---- Testes de Excluir ----

func TestExcluirRendaExtra_Existente(t *testing.T) {
	existente := &domain.RendaExtra{ID: "uuid-1", Descricao: "Bônus março"}
	mock := &MockRendaExtraRepository{returnRenda: existente}
	svc := service.NewRendaExtraService(mock)

	err := svc.Excluir("uuid-1")

	require.NoError(t, err)
	assert.True(t, mock.excluirChamado)
}

func TestExcluirRendaExtra_Inexistente(t *testing.T) {
	mock := &MockRendaExtraRepository{returnError: domain.ErrRendaExtraNaoEncontrada}
	svc := service.NewRendaExtraService(mock)

	err := svc.Excluir("uuid-inexistente")

	assert.ErrorIs(t, err, domain.ErrRendaExtraNaoEncontrada)
	assert.False(t, mock.excluirChamado)
}
