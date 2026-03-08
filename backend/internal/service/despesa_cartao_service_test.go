package service_test

import (
	"testing"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/felipehrs/financas-familiares/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockDespesaCartaoRepository implementa DespesaCartaoRepository para testes unitários.
type MockDespesaCartaoRepository struct {
	returnDespesa   *domain.DespesaCartao
	returnDespesas  []*domain.DespesaCartao
	returnError     error
	criarChamado    bool
	excluirChamado  bool
	despesaRecebida *domain.DespesaCartao
}

func (m *MockDespesaCartaoRepository) Criar(d *domain.DespesaCartao) (*domain.DespesaCartao, error) {
	m.criarChamado = true
	m.despesaRecebida = d
	if m.returnError != nil {
		return nil, m.returnError
	}
	criada := *d
	criada.ID = "uuid-despesa-mock"
	return &criada, nil
}

func (m *MockDespesaCartaoRepository) ListarPorCartao(cartaoID string) ([]*domain.DespesaCartao, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnDespesas, nil
}

func (m *MockDespesaCartaoRepository) ListarPorFatura(cartaoID string, mes, ano int) ([]*domain.DespesaCartao, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnDespesas, nil
}

func (m *MockDespesaCartaoRepository) BuscarPorID(id string) (*domain.DespesaCartao, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnDespesa, nil
}

func (m *MockDespesaCartaoRepository) Excluir(id string) error {
	m.excluirChamado = true
	return m.returnError
}

// MockCartaoRepositoryForDespesa implementa CartaoRepositoryForDespesa para testes unitários.
type MockCartaoRepositoryForDespesa struct {
	returnCartao *domain.CartaoCredito
	returnError  error
}

func (m *MockCartaoRepositoryForDespesa) BuscarPorID(id string) (*domain.CartaoCredito, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.returnCartao, nil
}

func cartaoComFechamento(dia int) *domain.CartaoCredito {
	return &domain.CartaoCredito{
		ID:            "cartao-1",
		Nome:          "Nubank",
		MembroID:      "membro-1",
		DiaFechamento: dia,
		DiaVencimento: 20,
		Ativo:         true,
	}
}

// ---- Testes de Criar ----

func TestCriarDespesaCartao_CompraDiaIgualFechamento_FaturaCorrentemesAtual(t *testing.T) {
	// Compra no dia 10, fechamento dia 10 → fatura mês corrente (março/2026)
	despesaRepo := &MockDespesaCartaoRepository{}
	cartaoRepo := &MockCartaoRepositoryForDespesa{returnCartao: cartaoComFechamento(10)}
	svc := service.NewDespesaCartaoService(despesaRepo, cartaoRepo)

	dataCompra := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
	despesa, err := svc.Criar("cartao-1", "Supermercado", nil, dataCompra, 150.00)

	require.NoError(t, err)
	assert.Equal(t, "uuid-despesa-mock", despesa.ID)
	assert.Equal(t, 3, despesa.FaturaMes)
	assert.Equal(t, 2026, despesa.FaturaAno)
	assert.Equal(t, 150.00, despesa.ValorParcela)
	assert.Equal(t, 1, despesa.NumeroParcelas)
	assert.True(t, despesaRepo.criarChamado)
}

func TestCriarDespesaCartao_CompraDiaMenorFechamento_FaturaAtual(t *testing.T) {
	// Compra no dia 5, fechamento dia 10 → fatura mês corrente (março/2026)
	despesaRepo := &MockDespesaCartaoRepository{}
	cartaoRepo := &MockCartaoRepositoryForDespesa{returnCartao: cartaoComFechamento(10)}
	svc := service.NewDespesaCartaoService(despesaRepo, cartaoRepo)

	dataCompra := time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC)
	despesa, err := svc.Criar("cartao-1", "Restaurante", nil, dataCompra, 80.00)

	require.NoError(t, err)
	assert.Equal(t, 3, despesa.FaturaMes)
	assert.Equal(t, 2026, despesa.FaturaAno)
}

func TestCriarDespesaCartao_CompraDiaMaiorFechamento_FaturaProximoMes(t *testing.T) {
	// Compra no dia 15, fechamento dia 10 → fatura mês seguinte (abril/2026)
	despesaRepo := &MockDespesaCartaoRepository{}
	cartaoRepo := &MockCartaoRepositoryForDespesa{returnCartao: cartaoComFechamento(10)}
	svc := service.NewDespesaCartaoService(despesaRepo, cartaoRepo)

	dataCompra := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)
	despesa, err := svc.Criar("cartao-1", "Eletrônico", nil, dataCompra, 500.00)

	require.NoError(t, err)
	assert.Equal(t, 4, despesa.FaturaMes)
	assert.Equal(t, 2026, despesa.FaturaAno)
}

func TestCriarDespesaCartao_DezembroDiaMaiorFechamento_FaturaJaneiroProximoAno(t *testing.T) {
	// Compra em dezembro, dia 20, fechamento dia 10 → fatura janeiro do ano seguinte
	despesaRepo := &MockDespesaCartaoRepository{}
	cartaoRepo := &MockCartaoRepositoryForDespesa{returnCartao: cartaoComFechamento(10)}
	svc := service.NewDespesaCartaoService(despesaRepo, cartaoRepo)

	dataCompra := time.Date(2026, 12, 20, 0, 0, 0, 0, time.UTC)
	despesa, err := svc.Criar("cartao-1", "Natal", nil, dataCompra, 300.00)

	require.NoError(t, err)
	assert.Equal(t, 1, despesa.FaturaMes)
	assert.Equal(t, 2027, despesa.FaturaAno)
}

func TestCriarDespesaCartao_CartaoNaoEncontrado(t *testing.T) {
	despesaRepo := &MockDespesaCartaoRepository{}
	cartaoRepo := &MockCartaoRepositoryForDespesa{returnError: domain.ErrCartaoNaoEncontrado}
	svc := service.NewDespesaCartaoService(despesaRepo, cartaoRepo)

	dataCompra := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
	_, err := svc.Criar("cartao-inexistente", "Supermercado", nil, dataCompra, 150.00)

	assert.ErrorIs(t, err, domain.ErrCartaoNaoEncontrado)
	assert.False(t, despesaRepo.criarChamado)
}

func TestCriarDespesaCartao_CartaoIDVazio(t *testing.T) {
	despesaRepo := &MockDespesaCartaoRepository{}
	cartaoRepo := &MockCartaoRepositoryForDespesa{}
	svc := service.NewDespesaCartaoService(despesaRepo, cartaoRepo)

	dataCompra := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
	_, err := svc.Criar("", "Supermercado", nil, dataCompra, 150.00)

	assert.ErrorIs(t, err, domain.ErrCartaoIDObrigatorio)
	assert.False(t, despesaRepo.criarChamado)
}

func TestCriarDespesaCartao_DescricaoVazia(t *testing.T) {
	despesaRepo := &MockDespesaCartaoRepository{}
	cartaoRepo := &MockCartaoRepositoryForDespesa{}
	svc := service.NewDespesaCartaoService(despesaRepo, cartaoRepo)

	dataCompra := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
	_, err := svc.Criar("cartao-1", "", nil, dataCompra, 150.00)

	assert.ErrorIs(t, err, domain.ErrDescricaoObrigatoria)
	assert.False(t, despesaRepo.criarChamado)
}

func TestCriarDespesaCartao_DescricaoApenasEspacos(t *testing.T) {
	despesaRepo := &MockDespesaCartaoRepository{}
	cartaoRepo := &MockCartaoRepositoryForDespesa{}
	svc := service.NewDespesaCartaoService(despesaRepo, cartaoRepo)

	dataCompra := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
	_, err := svc.Criar("cartao-1", "   ", nil, dataCompra, 150.00)

	assert.ErrorIs(t, err, domain.ErrDescricaoObrigatoria)
}

func TestCriarDespesaCartao_ValorZero(t *testing.T) {
	despesaRepo := &MockDespesaCartaoRepository{}
	cartaoRepo := &MockCartaoRepositoryForDespesa{}
	svc := service.NewDespesaCartaoService(despesaRepo, cartaoRepo)

	dataCompra := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
	_, err := svc.Criar("cartao-1", "Supermercado", nil, dataCompra, 0)

	assert.ErrorIs(t, err, domain.ErrValorTotalInvalido)
	assert.False(t, despesaRepo.criarChamado)
}

func TestCriarDespesaCartao_ValorNegativo(t *testing.T) {
	despesaRepo := &MockDespesaCartaoRepository{}
	cartaoRepo := &MockCartaoRepositoryForDespesa{}
	svc := service.NewDespesaCartaoService(despesaRepo, cartaoRepo)

	dataCompra := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
	_, err := svc.Criar("cartao-1", "Supermercado", nil, dataCompra, -50.00)

	assert.ErrorIs(t, err, domain.ErrValorTotalInvalido)
}

func TestCriarDespesaCartao_ComCategoriaID(t *testing.T) {
	despesaRepo := &MockDespesaCartaoRepository{}
	cartaoRepo := &MockCartaoRepositoryForDespesa{returnCartao: cartaoComFechamento(10)}
	svc := service.NewDespesaCartaoService(despesaRepo, cartaoRepo)

	catID := "cat-uuid-1"
	dataCompra := time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC)
	despesa, err := svc.Criar("cartao-1", "Mercado", &catID, dataCompra, 200.00)

	require.NoError(t, err)
	require.NotNil(t, despesa.CategoriaID)
	assert.Equal(t, "cat-uuid-1", *despesa.CategoriaID)
}

// ---- Testes de ListarPorCartao ----

func TestListarDespesasPorCartao_RetornaLista(t *testing.T) {
	despesas := []*domain.DespesaCartao{
		{ID: "d-1", CartaoID: "cartao-1", Descricao: "Supermercado", ValorTotal: 150.00},
		{ID: "d-2", CartaoID: "cartao-1", Descricao: "Farmácia", ValorTotal: 50.00},
	}
	despesaRepo := &MockDespesaCartaoRepository{returnDespesas: despesas}
	cartaoRepo := &MockCartaoRepositoryForDespesa{}
	svc := service.NewDespesaCartaoService(despesaRepo, cartaoRepo)

	resultado, err := svc.ListarPorCartao("cartao-1")

	require.NoError(t, err)
	assert.Len(t, resultado, 2)
}

func TestListarDespesasPorCartao_ListaVazia(t *testing.T) {
	despesaRepo := &MockDespesaCartaoRepository{returnDespesas: []*domain.DespesaCartao{}}
	cartaoRepo := &MockCartaoRepositoryForDespesa{}
	svc := service.NewDespesaCartaoService(despesaRepo, cartaoRepo)

	resultado, err := svc.ListarPorCartao("cartao-1")

	require.NoError(t, err)
	assert.Empty(t, resultado)
}

// ---- Testes de ListarPorFatura ----

func TestListarDespesasPorFatura_RetornaListaFiltrada(t *testing.T) {
	despesas := []*domain.DespesaCartao{
		{ID: "d-1", CartaoID: "cartao-1", FaturaMes: 3, FaturaAno: 2026, Descricao: "Compra A"},
	}
	despesaRepo := &MockDespesaCartaoRepository{returnDespesas: despesas}
	cartaoRepo := &MockCartaoRepositoryForDespesa{}
	svc := service.NewDespesaCartaoService(despesaRepo, cartaoRepo)

	resultado, err := svc.ListarPorFatura("cartao-1", 3, 2026)

	require.NoError(t, err)
	assert.Len(t, resultado, 1)
	assert.Equal(t, 3, resultado[0].FaturaMes)
	assert.Equal(t, 2026, resultado[0].FaturaAno)
}

// ---- Testes de Excluir ----

func TestExcluirDespesaCartao_Sucesso(t *testing.T) {
	existente := &domain.DespesaCartao{ID: "d-1", CartaoID: "cartao-1", Descricao: "Supermercado"}
	despesaRepo := &MockDespesaCartaoRepository{returnDespesa: existente}
	cartaoRepo := &MockCartaoRepositoryForDespesa{}
	svc := service.NewDespesaCartaoService(despesaRepo, cartaoRepo)

	err := svc.Excluir("d-1")

	require.NoError(t, err)
	assert.True(t, despesaRepo.excluirChamado)
}

func TestExcluirDespesaCartao_NaoEncontrada(t *testing.T) {
	despesaRepo := &MockDespesaCartaoRepository{returnError: domain.ErrDespesaCartaoNaoEncontrada}
	cartaoRepo := &MockCartaoRepositoryForDespesa{}
	svc := service.NewDespesaCartaoService(despesaRepo, cartaoRepo)

	err := svc.Excluir("d-inexistente")

	assert.ErrorIs(t, err, domain.ErrDespesaCartaoNaoEncontrada)
	assert.False(t, despesaRepo.excluirChamado)
}
