package service

import (
	"strings"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
)

// MembroRepository define as operações de persistência necessárias para membros.
// Declarada aqui para evitar import circular entre service e repository.
type MembroRepository interface {
	Criar(membro *domain.Membro) (*domain.Membro, error)
	BuscarPorID(id string) (*domain.Membro, error)
	Listar() ([]*domain.Membro, error)
	Atualizar(membro *domain.Membro) (*domain.Membro, error)
	Inativar(id string) error
}

// MembroServiceInterface define os métodos públicos do serviço de membros.
// Redeclarada nos handlers para desacoplamento.
type MembroServiceInterface interface {
	Criar(nome, relacionamento string) (*domain.Membro, error)
	BuscarPorID(id string) (*domain.Membro, error)
	Listar() ([]*domain.Membro, error)
	Atualizar(id, nome, relacionamento string, ativo bool) (*domain.Membro, error)
	Inativar(id string) error
}

// MembroService implementa a lógica de negócio para membros da família.
type MembroService struct {
	repo MembroRepository
}

// NewMembroService cria uma nova instância do MembroService.
func NewMembroService(repo MembroRepository) *MembroService {
	return &MembroService{repo: repo}
}

// Criar cria um novo membro da família.
// Valida que o nome não está vazio e define Ativo=true por padrão.
func (s *MembroService) Criar(nome, relacionamento string) (*domain.Membro, error) {
	if strings.TrimSpace(nome) == "" {
		return nil, domain.ErrNomeObrigatorio
	}

	membro := &domain.Membro{
		Nome:           nome,
		Relacionamento: relacionamento,
		Ativo:          true,
	}

	return s.repo.Criar(membro)
}

// BuscarPorID retorna um membro pelo seu ID.
// Retorna ErrMembroNaoEncontrado se não existir.
func (s *MembroService) BuscarPorID(id string) (*domain.Membro, error) {
	return s.repo.BuscarPorID(id)
}

// Listar retorna todos os membros (ativos e inativos, sem deletados).
func (s *MembroService) Listar() ([]*domain.Membro, error) {
	return s.repo.Listar()
}

// Atualizar atualiza os dados de um membro existente.
// Valida que o nome não está vazio antes de buscar no repositório.
func (s *MembroService) Atualizar(id, nome, relacionamento string, ativo bool) (*domain.Membro, error) {
	if strings.TrimSpace(nome) == "" {
		return nil, domain.ErrNomeObrigatorio
	}

	membro, err := s.repo.BuscarPorID(id)
	if err != nil {
		return nil, err
	}

	membro.Nome = nome
	membro.Relacionamento = relacionamento
	membro.Ativo = ativo

	return s.repo.Atualizar(membro)
}

// Inativar marca um membro como inativo (soft delete lógico).
// Não exclui o registro — apenas seta Ativo=false.
// Retorna ErrMembroNaoEncontrado se o ID não existir.
func (s *MembroService) Inativar(id string) error {
	_, err := s.repo.BuscarPorID(id)
	if err != nil {
		return err
	}
	return s.repo.Inativar(id)
}
