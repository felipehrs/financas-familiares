package service

import (
	"strings"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
)

// RendimentoInvestimentoRepositoryInterface define as operações de persistência para rendimentos de investimento.
// Declarada aqui para evitar import circular entre service e repository.
type RendimentoInvestimentoRepositoryInterface interface {
	Criar(familiaID string, r *domain.RendimentoInvestimento) (*domain.RendimentoInvestimento, error)
	BuscarPorID(familiaID, id string) (*domain.RendimentoInvestimento, error)
	Listar(familiaID string) ([]*domain.RendimentoInvestimento, error)
	ListarPorMes(familiaID string, mes, ano int) ([]*domain.RendimentoInvestimento, error)
	Atualizar(familiaID string, r *domain.RendimentoInvestimento) (*domain.RendimentoInvestimento, error)
	Excluir(familiaID, id string) error
}

// RendimentoInvestimentoServiceInterface define os métodos públicos do serviço de rendimentos de investimento.
// Redeclarada nos handlers para desacoplamento.
type RendimentoInvestimentoServiceInterface interface {
	Criar(familiaID, descricao, membroID string, data time.Time, valor, valorDistribuido float64) (*domain.RendimentoInvestimento, error)
	BuscarPorID(familiaID, id string) (*domain.RendimentoInvestimento, error)
	Listar(familiaID string) ([]*domain.RendimentoInvestimento, error)
	ListarPorMes(familiaID string, mes, ano int) ([]*domain.RendimentoInvestimento, error)
	Atualizar(familiaID, id, descricao, membroID string, data time.Time, valor, valorDistribuido float64) (*domain.RendimentoInvestimento, error)
	Excluir(familiaID, id string) error
}

// RendimentoInvestimentoService implementa a lógica de negócio para rendimentos de investimento.
type RendimentoInvestimentoService struct {
	repo RendimentoInvestimentoRepositoryInterface
}

// NewRendimentoInvestimentoService cria uma nova instância do RendimentoInvestimentoService.
func NewRendimentoInvestimentoService(repo RendimentoInvestimentoRepositoryInterface) *RendimentoInvestimentoService {
	return &RendimentoInvestimentoService{repo: repo}
}

// validarCamposRendimento valida os campos obrigatórios de criação e atualização.
func validarCamposRendimento(descricao, membroID string, data time.Time, valor, valorDistribuido float64) error {
	if strings.TrimSpace(descricao) == "" {
		return domain.ErrDescricaoRendimentoObrigatoria
	}
	if strings.TrimSpace(membroID) == "" {
		return domain.ErrMembroIDRendimentoObrigatorio
	}
	if valor <= 0 {
		return domain.ErrValorRendimentoInvalido
	}
	if data.IsZero() {
		return domain.ErrDataRendimentoObrigatoria
	}
	if valorDistribuido < 0 || valorDistribuido > valor {
		return domain.ErrValorDistribuidoInvalido
	}
	return nil
}

// Criar cria um novo rendimento de investimento.
// Valida campos obrigatórios antes de persistir.
func (s *RendimentoInvestimentoService) Criar(familiaID, descricao, membroID string, data time.Time, valor, valorDistribuido float64) (*domain.RendimentoInvestimento, error) {
	if err := validarCamposRendimento(descricao, membroID, data, valor, valorDistribuido); err != nil {
		return nil, err
	}

	rendimento := &domain.RendimentoInvestimento{
		Descricao:        descricao,
		MembroID:         membroID,
		Data:             data,
		Valor:            valor,
		ValorDistribuido: valorDistribuido,
	}

	return s.repo.Criar(familiaID, rendimento)
}

// BuscarPorID retorna um rendimento de investimento pelo seu ID.
// Retorna ErrRendimentoInvestimentoNaoEncontrado se não existir.
func (s *RendimentoInvestimentoService) BuscarPorID(familiaID, id string) (*domain.RendimentoInvestimento, error) {
	return s.repo.BuscarPorID(familiaID, id)
}

// Listar retorna todos os rendimentos de investimento não excluídos.
func (s *RendimentoInvestimentoService) Listar(familiaID string) ([]*domain.RendimentoInvestimento, error) {
	return s.repo.Listar(familiaID)
}

// ListarPorMes retorna os rendimentos do mês e ano especificados (pela data).
func (s *RendimentoInvestimentoService) ListarPorMes(familiaID string, mes, ano int) ([]*domain.RendimentoInvestimento, error) {
	return s.repo.ListarPorMes(familiaID, mes, ano)
}

// Atualizar atualiza os dados de um rendimento de investimento existente.
// Valida campos obrigatórios antes de buscar no repositório.
func (s *RendimentoInvestimentoService) Atualizar(familiaID, id, descricao, membroID string, data time.Time, valor, valorDistribuido float64) (*domain.RendimentoInvestimento, error) {
	if err := validarCamposRendimento(descricao, membroID, data, valor, valorDistribuido); err != nil {
		return nil, err
	}

	rendimento, err := s.repo.BuscarPorID(familiaID, id)
	if err != nil {
		return nil, err
	}

	rendimento.Descricao = descricao
	rendimento.MembroID = membroID
	rendimento.Data = data
	rendimento.Valor = valor
	rendimento.ValorDistribuido = valorDistribuido

	return s.repo.Atualizar(familiaID, rendimento)
}

// Excluir realiza o soft-delete de um rendimento de investimento existente.
// Retorna ErrRendimentoInvestimentoNaoEncontrado se não existir.
func (s *RendimentoInvestimentoService) Excluir(familiaID, id string) error {
	_, err := s.repo.BuscarPorID(familiaID, id)
	if err != nil {
		return err
	}

	return s.repo.Excluir(familiaID, id)
}
