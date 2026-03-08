package service

import (
	"strings"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
)

// DespesaGeralRepositoryInterface define as operações de persistência para despesas gerais.
// Declarada aqui para evitar import circular entre service e repository.
type DespesaGeralRepositoryInterface interface {
	Criar(d *domain.DespesaGeral) (*domain.DespesaGeral, error)
	BuscarPorID(id string) (*domain.DespesaGeral, error)
	Listar() ([]*domain.DespesaGeral, error)
	ListarPorMes(mes, ano int) ([]*domain.DespesaGeral, error)
	Atualizar(d *domain.DespesaGeral) (*domain.DespesaGeral, error)
	Excluir(id string) error
}

// DespesaGeralServiceInterface define os métodos públicos do serviço de despesas gerais.
// Redeclarada nos handlers para desacoplamento.
type DespesaGeralServiceInterface interface {
	Criar(membroID string, categoriaID *string, descricao string, data time.Time, valor float64, formaPagamento string, observacoes *string) (*domain.DespesaGeral, error)
	BuscarPorID(id string) (*domain.DespesaGeral, error)
	Listar() ([]*domain.DespesaGeral, error)
	ListarPorMes(mes, ano int) ([]*domain.DespesaGeral, error)
	Atualizar(id, membroID string, categoriaID *string, descricao string, data time.Time, valor float64, formaPagamento string, observacoes *string) (*domain.DespesaGeral, error)
	Excluir(id string) error
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
func (s *DespesaGeralService) Criar(membroID string, categoriaID *string, descricao string, data time.Time, valor float64, formaPagamento string, observacoes *string) (*domain.DespesaGeral, error) {
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

	return s.repo.Criar(despesa)
}

// BuscarPorID retorna uma despesa geral pelo seu ID.
// Retorna ErrDespesaGeralNaoEncontrada se não existir.
func (s *DespesaGeralService) BuscarPorID(id string) (*domain.DespesaGeral, error) {
	return s.repo.BuscarPorID(id)
}

// Listar retorna todas as despesas gerais não excluídas.
func (s *DespesaGeralService) Listar() ([]*domain.DespesaGeral, error) {
	return s.repo.Listar()
}

// ListarPorMes retorna as despesas gerais do mês e ano especificados.
func (s *DespesaGeralService) ListarPorMes(mes, ano int) ([]*domain.DespesaGeral, error) {
	return s.repo.ListarPorMes(mes, ano)
}

// Atualizar atualiza os dados de uma despesa geral existente.
// Valida campos obrigatórios antes de buscar no repositório.
func (s *DespesaGeralService) Atualizar(id, membroID string, categoriaID *string, descricao string, data time.Time, valor float64, formaPagamento string, observacoes *string) (*domain.DespesaGeral, error) {
	if err := validarCamposDespesaGeral(membroID, descricao, data, valor, formaPagamento); err != nil {
		return nil, err
	}

	despesa, err := s.repo.BuscarPorID(id)
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

	return s.repo.Atualizar(despesa)
}

// Excluir realiza o soft-delete de uma despesa geral existente.
// Retorna ErrDespesaGeralNaoEncontrada se não existir.
func (s *DespesaGeralService) Excluir(id string) error {
	_, err := s.repo.BuscarPorID(id)
	if err != nil {
		return err
	}

	return s.repo.Excluir(id)
}
