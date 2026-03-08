package service

import (
	"strings"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/google/uuid"
)

// DespesaCartaoRepository define as operações de persistência necessárias para despesas de cartão.
// Declarada aqui para evitar import circular entre service e repository.
type DespesaCartaoRepository interface {
	Criar(d *domain.DespesaCartao) (*domain.DespesaCartao, error)
	ListarPorCartao(cartaoID string) ([]*domain.DespesaCartao, error)
	ListarPorFatura(cartaoID string, mes, ano int) ([]*domain.DespesaCartao, error)
	BuscarPorID(id string) (*domain.DespesaCartao, error)
	Excluir(id string) error
}

// CartaoRepositoryForDespesa define os métodos de cartão que o service de despesas precisa.
type CartaoRepositoryForDespesa interface {
	BuscarPorID(id string) (*domain.CartaoCredito, error)
}

// DespesaCartaoServiceInterface define os métodos públicos do serviço de despesas de cartão.
// Redeclarada nos handlers para desacoplamento.
type DespesaCartaoServiceInterface interface {
	Criar(cartaoID, descricao string, categoriaID *string, dataCompra time.Time, valorTotal float64, numeroParcelas int) ([]*domain.DespesaCartao, error)
	ListarPorCartao(cartaoID string) ([]*domain.DespesaCartao, error)
	ListarPorFatura(cartaoID string, mes, ano int) ([]*domain.DespesaCartao, error)
	BuscarPorID(id string) (*domain.DespesaCartao, error)
	Excluir(id string) error
}

// DespesaCartaoService implementa a lógica de negócio para despesas de cartão de crédito.
type DespesaCartaoService struct {
	repo       DespesaCartaoRepository
	cartaoRepo CartaoRepositoryForDespesa
}

// NewDespesaCartaoService cria uma nova instância do DespesaCartaoService.
func NewDespesaCartaoService(repo DespesaCartaoRepository, cartaoRepo CartaoRepositoryForDespesa) *DespesaCartaoService {
	return &DespesaCartaoService{repo: repo, cartaoRepo: cartaoRepo}
}

// calcularFatura aplica a RN01: determina o mês e ano da fatura com base na data de compra e
// no dia de fechamento do cartão.
func calcularFatura(dataCompra time.Time, diaFechamento int) (mes, ano int) {
	if dataCompra.Day() <= diaFechamento {
		return int(dataCompra.Month()), dataCompra.Year()
	}
	// Compra após o fechamento: fatura do próximo mês
	proximo := dataCompra.AddDate(0, 1, 0)
	return int(proximo.Month()), proximo.Year()
}

// proximaFatura avança um mês a partir de mes/ano (RN03).
func proximaFatura(mes, ano int) (int, int) {
	if mes == 12 {
		return 1, ano + 1
	}
	return mes + 1, ano
}

// Criar cria uma ou mais despesas de cartão (uma por parcela — RN02, RN03).
// Valida os campos obrigatórios, busca o cartão para obter DiaFechamento,
// calcula a fatura inicial via RN01 e distribui as parcelas por mês.
// Retorna slice com todas as parcelas criadas.
func (s *DespesaCartaoService) Criar(cartaoID, descricao string, categoriaID *string, dataCompra time.Time, valorTotal float64, numeroParcelas int) ([]*domain.DespesaCartao, error) {
	if strings.TrimSpace(cartaoID) == "" {
		return nil, domain.ErrCartaoIDObrigatorio
	}
	if strings.TrimSpace(descricao) == "" {
		return nil, domain.ErrDescricaoObrigatoria
	}
	if valorTotal <= 0 {
		return nil, domain.ErrValorTotalInvalido
	}
	if numeroParcelas < 1 {
		return nil, domain.ErrNumeroParcelas
	}

	cartao, err := s.cartaoRepo.BuscarPorID(cartaoID)
	if err != nil {
		return nil, err
	}

	faturaMes, faturaAno := calcularFatura(dataCompra, cartao.DiaFechamento)
	valorParcela := valorTotal / float64(numeroParcelas)
	compraID := uuid.NewString()

	resultado := make([]*domain.DespesaCartao, 0, numeroParcelas)
	for i := 1; i <= numeroParcelas; i++ {
		despesa := &domain.DespesaCartao{
			CompraID:       compraID,
			CartaoID:       cartaoID,
			CategoriaID:    categoriaID,
			Descricao:      descricao,
			DataCompra:     dataCompra,
			ValorTotal:     valorTotal,
			NumeroParcelas: numeroParcelas,
			ParcelaNumero:  i,
			ValorParcela:   valorParcela,
			FaturaMes:      faturaMes,
			FaturaAno:      faturaAno,
		}

		criada, err := s.repo.Criar(despesa)
		if err != nil {
			return nil, err
		}
		resultado = append(resultado, criada)

		// Avança para a fatura do próximo mês (RN03)
		faturaMes, faturaAno = proximaFatura(faturaMes, faturaAno)
	}

	return resultado, nil
}

// ListarPorCartao retorna todas as despesas de um cartão (sem deletadas).
func (s *DespesaCartaoService) ListarPorCartao(cartaoID string) ([]*domain.DespesaCartao, error) {
	return s.repo.ListarPorCartao(cartaoID)
}

// ListarPorFatura retorna as despesas de um cartão filtradas por mês/ano de fatura.
func (s *DespesaCartaoService) ListarPorFatura(cartaoID string, mes, ano int) ([]*domain.DespesaCartao, error) {
	return s.repo.ListarPorFatura(cartaoID, mes, ano)
}

// BuscarPorID retorna uma despesa pelo seu ID.
// Retorna ErrDespesaCartaoNaoEncontrada se não existir.
func (s *DespesaCartaoService) BuscarPorID(id string) (*domain.DespesaCartao, error) {
	return s.repo.BuscarPorID(id)
}

// Excluir remove (soft delete) uma despesa de cartão pelo ID.
// Retorna ErrDespesaCartaoNaoEncontrada se não existir.
func (s *DespesaCartaoService) Excluir(id string) error {
	_, err := s.repo.BuscarPorID(id)
	if err != nil {
		return err
	}
	return s.repo.Excluir(id)
}
