package service

import (
	"strings"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
)

// RendaVariavelRepositoryInterface define as operações de persistência para rendas variáveis.
// Declarada aqui para evitar import circular entre service e repository.
type RendaVariavelRepositoryInterface interface {
	Criar(r *domain.RendaVariavel) (*domain.RendaVariavel, error)
	BuscarPorID(id string) (*domain.RendaVariavel, error)
	Listar() ([]*domain.RendaVariavel, error)
	ListarPorMes(mes, ano int) ([]*domain.RendaVariavel, error)
	Atualizar(r *domain.RendaVariavel) (*domain.RendaVariavel, error)
	Excluir(id string) error
}

// RendaVariavelServiceInterface define os métodos públicos do serviço de rendas variáveis.
// Redeclarada nos handlers para desacoplamento.
type RendaVariavelServiceInterface interface {
	Criar(descricao, membroID string, mesReferencia, anoReferencia int, valor float64, dataRecebimento time.Time) (*domain.RendaVariavel, error)
	BuscarPorID(id string) (*domain.RendaVariavel, error)
	Listar() ([]*domain.RendaVariavel, error)
	ListarPorMes(mes, ano int) ([]*domain.RendaVariavel, error)
	Atualizar(id, descricao, membroID string, mesReferencia, anoReferencia int, valor float64, dataRecebimento time.Time) (*domain.RendaVariavel, error)
	Excluir(id string) error
}

// RendaVariavelService implementa a lógica de negócio para rendas variáveis.
type RendaVariavelService struct {
	repo RendaVariavelRepositoryInterface
}

// NewRendaVariavelService cria uma nova instância do RendaVariavelService.
func NewRendaVariavelService(repo RendaVariavelRepositoryInterface) *RendaVariavelService {
	return &RendaVariavelService{repo: repo}
}

// validarCamposRendaVariavel valida os campos obrigatórios de criação e atualização.
func validarCamposRendaVariavel(descricao, membroID string, mesReferencia, anoReferencia int, valor float64, dataRecebimento time.Time) error {
	if strings.TrimSpace(descricao) == "" {
		return domain.ErrDescricaoRendaVariavelObrigatoria
	}
	if strings.TrimSpace(membroID) == "" {
		return domain.ErrMembroIDRendaVariavelObrigatorio
	}
	if mesReferencia < 1 || mesReferencia > 12 {
		return domain.ErrMesReferenciaRendaVariavelInvalido
	}
	if anoReferencia <= 0 {
		return domain.ErrAnoReferenciaRendaVariavelInvalido
	}
	if valor <= 0 {
		return domain.ErrValorRendaVariavelInvalido
	}
	if dataRecebimento.IsZero() {
		return domain.ErrDataRecebimentoRendaVariavelObrigatoria
	}
	return nil
}

// Criar cria uma nova renda variável.
// Valida campos obrigatórios antes de persistir.
func (s *RendaVariavelService) Criar(descricao, membroID string, mesReferencia, anoReferencia int, valor float64, dataRecebimento time.Time) (*domain.RendaVariavel, error) {
	if err := validarCamposRendaVariavel(descricao, membroID, mesReferencia, anoReferencia, valor, dataRecebimento); err != nil {
		return nil, err
	}

	renda := &domain.RendaVariavel{
		Descricao:       descricao,
		MembroID:        membroID,
		MesReferencia:   mesReferencia,
		AnoReferencia:   anoReferencia,
		Valor:           valor,
		DataRecebimento: dataRecebimento,
	}

	return s.repo.Criar(renda)
}

// BuscarPorID retorna uma renda variável pelo seu ID.
// Retorna ErrRendaVariavelNaoEncontrada se não existir.
func (s *RendaVariavelService) BuscarPorID(id string) (*domain.RendaVariavel, error) {
	return s.repo.BuscarPorID(id)
}

// Listar retorna todas as rendas variáveis não excluídas.
func (s *RendaVariavelService) Listar() ([]*domain.RendaVariavel, error) {
	return s.repo.Listar()
}

// ListarPorMes retorna as rendas variáveis do mês e ano especificados.
func (s *RendaVariavelService) ListarPorMes(mes, ano int) ([]*domain.RendaVariavel, error) {
	return s.repo.ListarPorMes(mes, ano)
}

// Atualizar atualiza os dados de uma renda variável existente.
// Valida campos obrigatórios antes de buscar no repositório.
func (s *RendaVariavelService) Atualizar(id, descricao, membroID string, mesReferencia, anoReferencia int, valor float64, dataRecebimento time.Time) (*domain.RendaVariavel, error) {
	if err := validarCamposRendaVariavel(descricao, membroID, mesReferencia, anoReferencia, valor, dataRecebimento); err != nil {
		return nil, err
	}

	renda, err := s.repo.BuscarPorID(id)
	if err != nil {
		return nil, err
	}

	renda.Descricao = descricao
	renda.MembroID = membroID
	renda.MesReferencia = mesReferencia
	renda.AnoReferencia = anoReferencia
	renda.Valor = valor
	renda.DataRecebimento = dataRecebimento

	return s.repo.Atualizar(renda)
}

// Excluir realiza o soft-delete de uma renda variável existente.
// Retorna ErrRendaVariavelNaoEncontrada se não existir.
func (s *RendaVariavelService) Excluir(id string) error {
	_, err := s.repo.BuscarPorID(id)
	if err != nil {
		return err
	}

	return s.repo.Excluir(id)
}
