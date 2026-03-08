package service

import (
	"strings"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
)

// ContaFixaRepository define as operações de persistência necessárias para contas fixas.
// Declarada aqui para evitar import circular entre service e repository.
type ContaFixaRepository interface {
	Criar(c *domain.ContaFixa) (*domain.ContaFixa, error)
	BuscarPorID(id string) (*domain.ContaFixa, error)
	Listar() ([]*domain.ContaFixa, error)
	Atualizar(c *domain.ContaFixa) (*domain.ContaFixa, error)
	AlterarAtivo(id string, ativa bool) (*domain.ContaFixa, error)
}

// ContaFixaServiceInterface define os métodos públicos do serviço de contas fixas.
// Redeclarada nos handlers para desacoplamento.
type ContaFixaServiceInterface interface {
	Criar(descricao, membroID string, categoriaID *string, valor float64, diaVencimento int, formaPagamento string) (*domain.ContaFixa, error)
	BuscarPorID(id string) (*domain.ContaFixa, error)
	Listar() ([]*domain.ContaFixa, error)
	Atualizar(id, descricao, membroID string, categoriaID *string, valor float64, diaVencimento int, formaPagamento string, ativa bool) (*domain.ContaFixa, error)
	AlterarAtivo(id string, ativa bool) (*domain.ContaFixa, error)
}

// ContaFixaService implementa a lógica de negócio para contas fixas mensais.
type ContaFixaService struct {
	repo ContaFixaRepository
}

// NewContaFixaService cria uma nova instância do ContaFixaService.
func NewContaFixaService(repo ContaFixaRepository) *ContaFixaService {
	return &ContaFixaService{repo: repo}
}

// validarCampos valida os campos comuns de criação e atualização.
func validarCamposContaFixa(descricao, membroID string, valor float64, diaVencimento int, formaPagamento string) error {
	if strings.TrimSpace(descricao) == "" {
		return domain.ErrDescricaoContaFixaObrigatoria
	}
	if strings.TrimSpace(membroID) == "" {
		return domain.ErrMembroIDContaFixaObrigatorio
	}
	if valor <= 0 {
		return domain.ErrValorContaFixaInvalido
	}
	if diaVencimento < 1 || diaVencimento > 28 {
		return domain.ErrDiaVencimentoContaFixaInvalido
	}
	if strings.TrimSpace(formaPagamento) == "" {
		return domain.ErrFormaPagamentoContaFixaObrigatoria
	}
	return nil
}

// Criar cria uma nova conta fixa mensal.
// Valida campos obrigatórios e define Ativa=true por padrão.
func (s *ContaFixaService) Criar(descricao, membroID string, categoriaID *string, valor float64, diaVencimento int, formaPagamento string) (*domain.ContaFixa, error) {
	if err := validarCamposContaFixa(descricao, membroID, valor, diaVencimento, formaPagamento); err != nil {
		return nil, err
	}

	conta := &domain.ContaFixa{
		Descricao:      descricao,
		MembroID:       membroID,
		CategoriaID:    categoriaID,
		Valor:          valor,
		DiaVencimento:  diaVencimento,
		FormaPagamento: formaPagamento,
		Ativa:          true,
	}

	return s.repo.Criar(conta)
}

// BuscarPorID retorna uma conta fixa pelo seu ID.
// Retorna ErrContaFixaNaoEncontrada se não existir.
func (s *ContaFixaService) BuscarPorID(id string) (*domain.ContaFixa, error) {
	return s.repo.BuscarPorID(id)
}

// Listar retorna todas as contas fixas não excluídas.
func (s *ContaFixaService) Listar() ([]*domain.ContaFixa, error) {
	return s.repo.Listar()
}

// Atualizar atualiza os dados de uma conta fixa existente.
// Valida campos obrigatórios antes de buscar no repositório.
func (s *ContaFixaService) Atualizar(id, descricao, membroID string, categoriaID *string, valor float64, diaVencimento int, formaPagamento string, ativa bool) (*domain.ContaFixa, error) {
	if err := validarCamposContaFixa(descricao, membroID, valor, diaVencimento, formaPagamento); err != nil {
		return nil, err
	}

	conta, err := s.repo.BuscarPorID(id)
	if err != nil {
		return nil, err
	}

	conta.Descricao = descricao
	conta.MembroID = membroID
	conta.CategoriaID = categoriaID
	conta.Valor = valor
	conta.DiaVencimento = diaVencimento
	conta.FormaPagamento = formaPagamento
	conta.Ativa = ativa

	return s.repo.Atualizar(conta)
}

// AlterarAtivo altera o estado ativo/inativo de uma conta fixa existente.
// Confirma a existência do registro antes de delegar ao repositório.
func (s *ContaFixaService) AlterarAtivo(id string, ativa bool) (*domain.ContaFixa, error) {
	_, err := s.repo.BuscarPorID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.AlterarAtivo(id, ativa)
}
