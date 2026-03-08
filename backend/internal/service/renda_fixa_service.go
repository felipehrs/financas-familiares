package service

import (
	"strings"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
)

// RendaFixaRepository define as operações de persistência necessárias para rendas fixas.
// Declarada aqui para evitar import circular entre service e repository.
type RendaFixaRepository interface {
	Criar(r *domain.RendaFixa) (*domain.RendaFixa, error)
	BuscarPorID(id string) (*domain.RendaFixa, error)
	Listar() ([]*domain.RendaFixa, error)
	ListarAtivas() ([]*domain.RendaFixa, error)
	ListarVigentesPorMes(mes, ano int) ([]*domain.RendaFixa, error)
	Atualizar(r *domain.RendaFixa) (*domain.RendaFixa, error)
	Inativar(id string) error
}

// RendaFixaServiceInterface define os métodos públicos do serviço de rendas fixas.
// Redeclarada nos handlers para desacoplamento.
type RendaFixaServiceInterface interface {
	Criar(descricao, membroID string, valor float64, diaRecebimento int, dataInicio time.Time, dataFim *time.Time) (*domain.RendaFixa, error)
	BuscarPorID(id string) (*domain.RendaFixa, error)
	Listar() ([]*domain.RendaFixa, error)
	ListarAtivas() ([]*domain.RendaFixa, error)
	ListarVigentesPorMes(mes, ano int) ([]*domain.RendaFixa, error)
	Atualizar(id, descricao, membroID string, valor float64, diaRecebimento int, ativa bool, dataInicio time.Time, dataFim *time.Time) (*domain.RendaFixa, error)
	Inativar(id string) error
}

// RendaFixaService implementa a lógica de negócio para rendas fixas.
type RendaFixaService struct {
	repo RendaFixaRepository
}

// NewRendaFixaService cria uma nova instância do RendaFixaService.
func NewRendaFixaService(repo RendaFixaRepository) *RendaFixaService {
	return &RendaFixaService{repo: repo}
}

func validarRendaFixa(descricao, membroID string, valor float64, diaRecebimento int, dataInicio time.Time) error {
	if strings.TrimSpace(descricao) == "" {
		return domain.ErrDescricaoRendaObrigatoria
	}
	if strings.TrimSpace(membroID) == "" {
		return domain.ErrMembroIDRendaObrigatorio
	}
	if valor <= 0 {
		return domain.ErrValorRendaInvalido
	}
	if diaRecebimento < 1 || diaRecebimento > 31 {
		return domain.ErrDiaRecebimentoInvalido
	}
	if dataInicio.IsZero() {
		return domain.ErrDataInicioRendaObrigatoria
	}
	return nil
}

// Criar cria uma nova renda fixa.
// Valida os campos obrigatórios e define Ativa=true por padrão.
func (s *RendaFixaService) Criar(descricao, membroID string, valor float64, diaRecebimento int, dataInicio time.Time, dataFim *time.Time) (*domain.RendaFixa, error) {
	if err := validarRendaFixa(descricao, membroID, valor, diaRecebimento, dataInicio); err != nil {
		return nil, err
	}

	renda := &domain.RendaFixa{
		Descricao:      descricao,
		MembroID:       membroID,
		Valor:          valor,
		DiaRecebimento: diaRecebimento,
		Ativa:          true,
		DataInicio:     dataInicio,
		DataFim:        dataFim,
	}

	return s.repo.Criar(renda)
}

// BuscarPorID retorna uma renda fixa pelo seu ID.
// Retorna ErrRendaFixaNaoEncontrada se não existir.
func (s *RendaFixaService) BuscarPorID(id string) (*domain.RendaFixa, error) {
	return s.repo.BuscarPorID(id)
}

// Listar retorna todas as rendas fixas (ativas e inativas, sem deletadas).
func (s *RendaFixaService) Listar() ([]*domain.RendaFixa, error) {
	return s.repo.Listar()
}

// ListarAtivas retorna apenas as rendas fixas ativas (para cálculo de saldo e projeções).
func (s *RendaFixaService) ListarAtivas() ([]*domain.RendaFixa, error) {
	return s.repo.ListarAtivas()
}

// ListarVigentesPorMes retorna as rendas fixas vigentes em um dado mês/ano.
// Uma renda é vigente se: data_inicio <= mês/ano AND (data_fim IS NULL OR data_fim >= mês/ano).
func (s *RendaFixaService) ListarVigentesPorMes(mes, ano int) ([]*domain.RendaFixa, error) {
	return s.repo.ListarVigentesPorMes(mes, ano)
}

// Atualizar atualiza os dados de uma renda fixa existente.
// Valida os campos obrigatórios antes de buscar no repositório.
func (s *RendaFixaService) Atualizar(id, descricao, membroID string, valor float64, diaRecebimento int, ativa bool, dataInicio time.Time, dataFim *time.Time) (*domain.RendaFixa, error) {
	if err := validarRendaFixa(descricao, membroID, valor, diaRecebimento, dataInicio); err != nil {
		return nil, err
	}

	renda, err := s.repo.BuscarPorID(id)
	if err != nil {
		return nil, err
	}

	renda.Descricao = descricao
	renda.MembroID = membroID
	renda.Valor = valor
	renda.DiaRecebimento = diaRecebimento
	renda.Ativa = ativa
	renda.DataInicio = dataInicio
	renda.DataFim = dataFim

	return s.repo.Atualizar(renda)
}

// Inativar marca uma renda fixa como inativa (soft delete lógico).
// Não exclui o registro — apenas seta Ativa=false.
// Retorna ErrRendaFixaNaoEncontrada se o ID não existir.
func (s *RendaFixaService) Inativar(id string) error {
	_, err := s.repo.BuscarPorID(id)
	if err != nil {
		return err
	}
	return s.repo.Inativar(id)
}
