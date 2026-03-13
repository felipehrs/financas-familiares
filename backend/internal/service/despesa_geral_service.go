package service

import (
	"strings"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
)

// DespesaGeralRepositoryInterface define as operações de persistência para despesas gerais.
// Declarada aqui para evitar import circular entre service e repository.
type DespesaGeralRepositoryInterface interface {
	Criar(familiaID string, d *domain.DespesaGeral) (*domain.DespesaGeral, error)
	BuscarPorID(familiaID, id string) (*domain.DespesaGeral, error)
	Listar(familiaID string) ([]*domain.DespesaGeral, error)
	ListarPorMes(familiaID string, mes, ano int) ([]*domain.DespesaGeral, error)
	Atualizar(familiaID string, d *domain.DespesaGeral) (*domain.DespesaGeral, error)
	Excluir(familiaID, id string) error
}

// DespesaGeralServiceInterface define os métodos públicos do serviço de despesas gerais.
// Redeclarada nos handlers para desacoplamento.
type DespesaGeralServiceInterface interface {
	Criar(familiaID, membroID string, categoriaID *string, descricao string, data time.Time, valor float64, formaPagamento string, observacoes *string) (*domain.DespesaGeral, error)
	BuscarPorID(familiaID, id string) (*domain.DespesaGeral, error)
	Listar(familiaID string) ([]*domain.DespesaGeral, error)
	ListarPorMes(familiaID string, mes, ano int) ([]*domain.DespesaGeral, error)
	Atualizar(familiaID, id, membroID string, categoriaID *string, descricao string, data time.Time, valor float64, formaPagamento string, observacoes *string) (*domain.DespesaGeral, error)
	Excluir(familiaID, id string) error
}

// DespesaGeralService implementa a lógica de negócio para despesas gerais.
type DespesaGeralService struct {
	repo DespesaGeralRepositoryInterface
}

// NewDespesaGeralService cria uma nova instância do DespesaGeralService.
func NewDespesaGeralService(repo DespesaGeralRepositoryInterface) *DespesaGeralService {
	return &DespesaGeralService{repo: repo}
}

// validarCamposDespesaGeral valida os campos obrigatórios de criação e atualização.
func validarCamposDespesaGeral(membroID, descricao string, data time.Time, valor float64, formaPagamento string) error {
	if strings.TrimSpace(membroID) == "" {
		return domain.ErrMembroIDDespesaGeralObrigatorio
	}
	if strings.TrimSpace(descricao) == "" {
		return domain.ErrDescricaoDespesaGeralObrigatoria
	}
	if data.IsZero() {
		return domain.ErrDataDespesaGeralObrigatoria
	}
	if valor <= 0 {
		return domain.ErrValorDespesaGeralInvalido
	}
	if strings.TrimSpace(formaPagamento) == "" {
		return domain.ErrFormaPagamentoDespesaGeralObrigatoria
	}
	return nil
}

// Criar cria uma nova despesa geral.
// Valida campos obrigatórios antes de persistir.
func (s *DespesaGeralService) Criar(familiaID, membroID string, categoriaID *string, descricao string, data time.Time, valor float64, formaPagamento string, observacoes *string) (*domain.DespesaGeral, error) {
	if err := validarCamposDespesaGeral(membroID, descricao, data, valor, formaPagamento); err != nil {
		return nil, err
	}

	despesa := &domain.DespesaGeral{
		MembroID:       membroID,
		CategoriaID:    categoriaID,
		Descricao:      descricao,
		Data:           data,
		Valor:          valor,
		FormaPagamento: formaPagamento,
		Observacoes:    observacoes,
	}

	return s.repo.Criar(familiaID, despesa)
}

// BuscarPorID retorna uma despesa geral pelo seu ID.
// Retorna ErrDespesaGeralNaoEncontrada se não existir.
func (s *DespesaGeralService) BuscarPorID(familiaID, id string) (*domain.DespesaGeral, error) {
	return s.repo.BuscarPorID(familiaID, id)
}

// Listar retorna todas as despesas gerais não excluídas.
func (s *DespesaGeralService) Listar(familiaID string) ([]*domain.DespesaGeral, error) {
	return s.repo.Listar(familiaID)
}

// ListarPorMes retorna as despesas gerais do mês e ano especificados.
func (s *DespesaGeralService) ListarPorMes(familiaID string, mes, ano int) ([]*domain.DespesaGeral, error) {
	return s.repo.ListarPorMes(familiaID, mes, ano)
}

// Atualizar atualiza os dados de uma despesa geral existente.
// Valida campos obrigatórios antes de buscar no repositório.
func (s *DespesaGeralService) Atualizar(familiaID, id, membroID string, categoriaID *string, descricao string, data time.Time, valor float64, formaPagamento string, observacoes *string) (*domain.DespesaGeral, error) {
	if err := validarCamposDespesaGeral(membroID, descricao, data, valor, formaPagamento); err != nil {
		return nil, err
	}

	despesa, err := s.repo.BuscarPorID(familiaID, id)
	if err != nil {
		return nil, err
	}

	despesa.MembroID = membroID
	despesa.CategoriaID = categoriaID
	despesa.Descricao = descricao
	despesa.Data = data
	despesa.Valor = valor
	despesa.FormaPagamento = formaPagamento
	despesa.Observacoes = observacoes

	return s.repo.Atualizar(familiaID, despesa)
}

// Excluir realiza o soft-delete de uma despesa geral existente.
// Retorna ErrDespesaGeralNaoEncontrada se não existir.
func (s *DespesaGeralService) Excluir(familiaID, id string) error {
	_, err := s.repo.BuscarPorID(familiaID, id)
	if err != nil {
		return err
	}

	return s.repo.Excluir(familiaID, id)
}
