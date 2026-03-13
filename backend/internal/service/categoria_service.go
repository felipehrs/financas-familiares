package service

import (
	"strings"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
)

// CategoriaRepository define as operações de persistência necessárias para categorias.
// Declarada aqui para evitar import circular entre service e repository.
type CategoriaRepository interface {
	Criar(familiaID string, categoria *domain.Categoria) (*domain.Categoria, error)
	Listar(familiaID string) ([]*domain.Categoria, error)
	BuscarPorID(familiaID, id string) (*domain.Categoria, error)
	Atualizar(familiaID string, categoria *domain.Categoria) (*domain.Categoria, error)
	Excluir(familiaID, id string) error // retorna ErrCategoriaComVinculos se tiver vínculos
}

// CategoriaServiceInterface define os métodos públicos do serviço de categorias.
// Redeclarada nos handlers para desacoplamento.
type CategoriaServiceInterface interface {
	Criar(familiaID, nome string) (*domain.Categoria, error)
	Listar(familiaID string) ([]*domain.Categoria, error)
	Atualizar(familiaID, id, nome string) (*domain.Categoria, error)
	Excluir(familiaID, id string) error
}

// CategoriaService implementa a lógica de negócio para categorias.
type CategoriaService struct {
	repo CategoriaRepository
}

// NewCategoriaService cria uma nova instância do CategoriaService.
func NewCategoriaService(repo CategoriaRepository) *CategoriaService {
	return &CategoriaService{repo: repo}
}

// Criar cria uma nova categoria.
// Valida que o nome não está vazio.
func (s *CategoriaService) Criar(familiaID, nome string) (*domain.Categoria, error) {
	if strings.TrimSpace(nome) == "" {
		return nil, domain.ErrNomeCategoriaObrigatorio
	}

	categoria := &domain.Categoria{
		Nome: nome,
	}

	return s.repo.Criar(familiaID, categoria)
}

// Listar retorna todas as categorias não excluídas.
func (s *CategoriaService) Listar(familiaID string) ([]*domain.Categoria, error) {
	return s.repo.Listar(familiaID)
}

// Atualizar atualiza o nome de uma categoria existente.
// Valida que o nome não está vazio antes de buscar no repositório.
func (s *CategoriaService) Atualizar(familiaID, id, nome string) (*domain.Categoria, error) {
	if strings.TrimSpace(nome) == "" {
		return nil, domain.ErrNomeCategoriaObrigatorio
	}

	categoria, err := s.repo.BuscarPorID(familiaID, id)
	if err != nil {
		return nil, err
	}

	categoria.Nome = nome

	return s.repo.Atualizar(familiaID, categoria)
}

// Excluir remove uma categoria (soft delete).
// Retorna ErrCategoriaNaoEncontrada se o ID não existir.
// Retorna ErrCategoriaComVinculos se houver registros vinculados.
func (s *CategoriaService) Excluir(familiaID, id string) error {
	_, err := s.repo.BuscarPorID(familiaID, id)
	if err != nil {
		return err
	}
	return s.repo.Excluir(familiaID, id)
}
