package service

import (
	"strings"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
)

// CartaoCreditoRepository define as operações de persistência necessárias para cartões de crédito.
// Declarada aqui para evitar import circular entre service e repository.
type CartaoCreditoRepository interface {
	Criar(cartao *domain.CartaoCredito) (*domain.CartaoCredito, error)
	BuscarPorID(id string) (*domain.CartaoCredito, error)
	Listar() ([]*domain.CartaoCredito, error)
	Atualizar(cartao *domain.CartaoCredito) (*domain.CartaoCredito, error)
	Inativar(id string) error
}

// CartaoCreditoServiceInterface define os métodos públicos do serviço de cartões de crédito.
// Redeclarada nos handlers para desacoplamento.
type CartaoCreditoServiceInterface interface {
	Criar(nome, membroID string, diaFechamento, diaVencimento int, limite *float64) (*domain.CartaoCredito, error)
	BuscarPorID(id string) (*domain.CartaoCredito, error)
	Listar() ([]*domain.CartaoCredito, error)
	Atualizar(id, nome, membroID string, diaFechamento, diaVencimento int, limite *float64, ativo bool) (*domain.CartaoCredito, error)
	Inativar(id string) error
}

// CartaoCreditoService implementa a lógica de negócio para cartões de crédito.
type CartaoCreditoService struct {
	repo CartaoCreditoRepository
}

// NewCartaoCreditoService cria uma nova instância do CartaoCreditoService.
func NewCartaoCreditoService(repo CartaoCreditoRepository) *CartaoCreditoService {
	return &CartaoCreditoService{repo: repo}
}

func validarCartao(nome, membroID string, diaFechamento, diaVencimento int) error {
	if strings.TrimSpace(nome) == "" {
		return domain.ErrNomeCartaoObrigatorio
	}
	if strings.TrimSpace(membroID) == "" {
		return domain.ErrMembroIDObrigatorio
	}
	if diaFechamento < 1 || diaFechamento > 31 {
		return domain.ErrDiaFechamentoInvalido
	}
	if diaVencimento < 1 || diaVencimento > 31 {
		return domain.ErrDiaVencimentoInvalido
	}
	return nil
}

// Criar cria um novo cartão de crédito.
// Valida os campos obrigatórios e define Ativo=true por padrão.
func (s *CartaoCreditoService) Criar(nome, membroID string, diaFechamento, diaVencimento int, limite *float64) (*domain.CartaoCredito, error) {
	if err := validarCartao(nome, membroID, diaFechamento, diaVencimento); err != nil {
		return nil, err
	}

	cartao := &domain.CartaoCredito{
		Nome:          nome,
		MembroID:      membroID,
		DiaFechamento: diaFechamento,
		DiaVencimento: diaVencimento,
		Limite:        limite,
		Ativo:         true,
	}

	return s.repo.Criar(cartao)
}

// BuscarPorID retorna um cartão pelo seu ID.
// Retorna ErrCartaoNaoEncontrado se não existir.
func (s *CartaoCreditoService) BuscarPorID(id string) (*domain.CartaoCredito, error) {
	return s.repo.BuscarPorID(id)
}

// Listar retorna todos os cartões (ativos e inativos, sem deletados).
func (s *CartaoCreditoService) Listar() ([]*domain.CartaoCredito, error) {
	return s.repo.Listar()
}

// Atualizar atualiza os dados de um cartão existente.
// Valida os campos obrigatórios antes de buscar no repositório.
func (s *CartaoCreditoService) Atualizar(id, nome, membroID string, diaFechamento, diaVencimento int, limite *float64, ativo bool) (*domain.CartaoCredito, error) {
	if err := validarCartao(nome, membroID, diaFechamento, diaVencimento); err != nil {
		return nil, err
	}

	cartao, err := s.repo.BuscarPorID(id)
	if err != nil {
		return nil, err
	}

	cartao.Nome = nome
	cartao.MembroID = membroID
	cartao.DiaFechamento = diaFechamento
	cartao.DiaVencimento = diaVencimento
	cartao.Limite = limite
	cartao.Ativo = ativo

	return s.repo.Atualizar(cartao)
}

// Inativar marca um cartão como inativo (soft delete lógico).
// Não exclui o registro — apenas seta Ativo=false.
// Retorna ErrCartaoNaoEncontrado se o ID não existir.
func (s *CartaoCreditoService) Inativar(id string) error {
	_, err := s.repo.BuscarPorID(id)
	if err != nil {
		return err
	}
	return s.repo.Inativar(id)
}
