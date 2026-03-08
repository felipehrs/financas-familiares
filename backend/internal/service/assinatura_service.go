package service

import (
	"strings"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
)

// AssinaturaRepository define as operações de persistência necessárias para assinaturas.
// Declarada aqui para evitar import circular entre service e repository.
type AssinaturaRepository interface {
	Criar(a *domain.Assinatura) (*domain.Assinatura, error)
	BuscarPorID(id string) (*domain.Assinatura, error)
	Listar() ([]*domain.Assinatura, error)
	Atualizar(a *domain.Assinatura) (*domain.Assinatura, error)
	AlterarStatus(id, status string) (*domain.Assinatura, error)
}

// AssinaturaServiceInterface define os métodos públicos do serviço de assinaturas.
// Redeclarada nos handlers para desacoplamento.
type AssinaturaServiceInterface interface {
	Criar(nome, membroID string, categoriaID *string, valor float64, diaCobranca int, formaPagamento string) (*domain.Assinatura, error)
	BuscarPorID(id string) (*domain.Assinatura, error)
	Listar() ([]*domain.Assinatura, error)
	Atualizar(id, nome, membroID string, categoriaID *string, valor float64, diaCobranca int, formaPagamento, status string) (*domain.Assinatura, error)
	AlterarStatus(id, status string) (*domain.Assinatura, error)
}

// AssinaturaService implementa a lógica de negócio para assinaturas recorrentes.
type AssinaturaService struct {
	repo AssinaturaRepository
}

// NewAssinaturaService cria uma nova instância do AssinaturaService.
func NewAssinaturaService(repo AssinaturaRepository) *AssinaturaService {
	return &AssinaturaService{repo: repo}
}

// Criar cria uma nova assinatura recorrente.
// Valida campos obrigatórios e define Status="ativa" por padrão.
func (s *AssinaturaService) Criar(nome, membroID string, categoriaID *string, valor float64, diaCobranca int, formaPagamento string) (*domain.Assinatura, error) {
	if strings.TrimSpace(nome) == "" {
		return nil, domain.ErrNomeAssinaturaObrigatorio
	}
	if strings.TrimSpace(membroID) == "" {
		return nil, domain.ErrMembroIDAssinaturaObrigatorio
	}
	if valor <= 0 {
		return nil, domain.ErrValorAssinaturaInvalido
	}
	if diaCobranca < 1 || diaCobranca > 31 {
		return nil, domain.ErrDiaCobrancaInvalido
	}
	if strings.TrimSpace(formaPagamento) == "" {
		return nil, domain.ErrFormaPagamentoObrigatoria
	}

	assinatura := &domain.Assinatura{
		Nome:           nome,
		MembroID:       membroID,
		CategoriaID:    categoriaID,
		Valor:          valor,
		DiaCobranca:    diaCobranca,
		FormaPagamento: formaPagamento,
		Status:         "ativa",
	}

	return s.repo.Criar(assinatura)
}

// BuscarPorID retorna uma assinatura pelo seu ID.
// Retorna ErrAssinaturaNaoEncontrada se não existir.
func (s *AssinaturaService) BuscarPorID(id string) (*domain.Assinatura, error) {
	return s.repo.BuscarPorID(id)
}

// Listar retorna todas as assinaturas não excluídas.
func (s *AssinaturaService) Listar() ([]*domain.Assinatura, error) {
	return s.repo.Listar()
}

// Atualizar atualiza os dados de uma assinatura existente.
// Valida campos obrigatórios antes de buscar no repositório.
func (s *AssinaturaService) Atualizar(id, nome, membroID string, categoriaID *string, valor float64, diaCobranca int, formaPagamento, status string) (*domain.Assinatura, error) {
	if strings.TrimSpace(nome) == "" {
		return nil, domain.ErrNomeAssinaturaObrigatorio
	}
	if strings.TrimSpace(membroID) == "" {
		return nil, domain.ErrMembroIDAssinaturaObrigatorio
	}
	if valor <= 0 {
		return nil, domain.ErrValorAssinaturaInvalido
	}
	if diaCobranca < 1 || diaCobranca > 31 {
		return nil, domain.ErrDiaCobrancaInvalido
	}
	if strings.TrimSpace(formaPagamento) == "" {
		return nil, domain.ErrFormaPagamentoObrigatoria
	}
	if !domain.StatusAssinaturaValidos[status] {
		return nil, domain.ErrStatusInvalido
	}

	assinatura, err := s.repo.BuscarPorID(id)
	if err != nil {
		return nil, err
	}

	assinatura.Nome = nome
	assinatura.MembroID = membroID
	assinatura.CategoriaID = categoriaID
	assinatura.Valor = valor
	assinatura.DiaCobranca = diaCobranca
	assinatura.FormaPagamento = formaPagamento
	assinatura.Status = status

	return s.repo.Atualizar(assinatura)
}

// AlterarStatus altera o status de uma assinatura existente.
// Valida o status antes de buscar no repositório.
func (s *AssinaturaService) AlterarStatus(id, status string) (*domain.Assinatura, error) {
	if !domain.StatusAssinaturaValidos[status] {
		return nil, domain.ErrStatusInvalido
	}

	_, err := s.repo.BuscarPorID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.AlterarStatus(id, status)
}
