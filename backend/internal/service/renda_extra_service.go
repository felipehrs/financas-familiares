package service

import (
	"strings"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
)

// RendaExtraRepositoryInterface define as operações de persistência para rendas extras.
// Declarada aqui para evitar import circular entre service e repository.
type RendaExtraRepositoryInterface interface {
	Criar(familiaID string, r *domain.RendaExtra) (*domain.RendaExtra, error)
	BuscarPorID(familiaID, id string) (*domain.RendaExtra, error)
	Listar(familiaID string) ([]*domain.RendaExtra, error)
	ListarPorMes(familiaID string, mes, ano int) ([]*domain.RendaExtra, error)
	Atualizar(familiaID string, r *domain.RendaExtra) (*domain.RendaExtra, error)
	Excluir(familiaID, id string) error
}

// RendaExtraServiceInterface define os métodos públicos do serviço de rendas extras.
// Redeclarada nos handlers para desacoplamento.
type RendaExtraServiceInterface interface {
	Criar(familiaID, descricao, membroID string, dataRecebimento time.Time, valor float64) (*domain.RendaExtra, error)
	BuscarPorID(familiaID, id string) (*domain.RendaExtra, error)
	Listar(familiaID string) ([]*domain.RendaExtra, error)
	ListarPorMes(familiaID string, mes, ano int) ([]*domain.RendaExtra, error)
	Atualizar(familiaID, id, descricao, membroID string, dataRecebimento time.Time, valor float64) (*domain.RendaExtra, error)
	Excluir(familiaID, id string) error
}

// RendaExtraService implementa a lógica de negócio para rendas extras.
type RendaExtraService struct {
	repo RendaExtraRepositoryInterface
}

// NewRendaExtraService cria uma nova instância do RendaExtraService.
func NewRendaExtraService(repo RendaExtraRepositoryInterface) *RendaExtraService {
	return &RendaExtraService{repo: repo}
}

// validarCamposRendaExtra valida os campos obrigatórios de criação e atualização.
func validarCamposRendaExtra(descricao, membroID string, dataRecebimento time.Time, valor float64) error {
	if strings.TrimSpace(descricao) == "" {
		return domain.ErrDescricaoRendaExtraObrigatoria
	}
	if strings.TrimSpace(membroID) == "" {
		return domain.ErrMembroIDRendaExtraObrigatorio
	}
	if valor <= 0 {
		return domain.ErrValorRendaExtraInvalido
	}
	if dataRecebimento.IsZero() {
		return domain.ErrDataRecebimentoRendaExtraObrigatoria
	}
	return nil
}

// Criar cria uma nova renda extra.
// Valida campos obrigatórios antes de persistir.
func (s *RendaExtraService) Criar(familiaID, descricao, membroID string, dataRecebimento time.Time, valor float64) (*domain.RendaExtra, error) {
	if err := validarCamposRendaExtra(descricao, membroID, dataRecebimento, valor); err != nil {
		return nil, err
	}

	renda := &domain.RendaExtra{
		Descricao:       descricao,
		MembroID:        membroID,
		DataRecebimento: dataRecebimento,
		Valor:           valor,
	}

	return s.repo.Criar(familiaID, renda)
}

// BuscarPorID retorna uma renda extra pelo seu ID.
// Retorna ErrRendaExtraNaoEncontrada se não existir.
func (s *RendaExtraService) BuscarPorID(familiaID, id string) (*domain.RendaExtra, error) {
	return s.repo.BuscarPorID(familiaID, id)
}

// Listar retorna todas as rendas extras não excluídas.
func (s *RendaExtraService) Listar(familiaID string) ([]*domain.RendaExtra, error) {
	return s.repo.Listar(familiaID)
}

// ListarPorMes retorna as rendas extras do mês e ano especificados (pela data_recebimento).
func (s *RendaExtraService) ListarPorMes(familiaID string, mes, ano int) ([]*domain.RendaExtra, error) {
	return s.repo.ListarPorMes(familiaID, mes, ano)
}

// Atualizar atualiza os dados de uma renda extra existente.
// Valida campos obrigatórios antes de buscar no repositório.
func (s *RendaExtraService) Atualizar(familiaID, id, descricao, membroID string, dataRecebimento time.Time, valor float64) (*domain.RendaExtra, error) {
	if err := validarCamposRendaExtra(descricao, membroID, dataRecebimento, valor); err != nil {
		return nil, err
	}

	renda, err := s.repo.BuscarPorID(familiaID, id)
	if err != nil {
		return nil, err
	}

	renda.Descricao = descricao
	renda.MembroID = membroID
	renda.DataRecebimento = dataRecebimento
	renda.Valor = valor

	return s.repo.Atualizar(familiaID, renda)
}

// Excluir realiza o soft-delete de uma renda extra existente.
// Retorna ErrRendaExtraNaoEncontrada se não existir.
func (s *RendaExtraService) Excluir(familiaID, id string) error {
	_, err := s.repo.BuscarPorID(familiaID, id)
	if err != nil {
		return err
	}

	return s.repo.Excluir(familiaID, id)
}
