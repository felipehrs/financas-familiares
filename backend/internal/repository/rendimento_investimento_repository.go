package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

// RendimentoInvestimentoRepository implementa service.RendimentoInvestimentoRepositoryInterface usando PostgreSQL via sqlx.
type RendimentoInvestimentoRepository struct {
	db *sqlx.DB
}

// NewRendimentoInvestimentoRepository cria uma nova instância do RendimentoInvestimentoRepository.
func NewRendimentoInvestimentoRepository(db *sqlx.DB) *RendimentoInvestimentoRepository {
	return &RendimentoInvestimentoRepository{db: db}
}

// rendimentoInvestimentoRow é a estrutura de scan para as colunas do banco.
type rendimentoInvestimentoRow struct {
	ID               string    `db:"id"`
	Descricao        string    `db:"descricao"`
	MembroID         string    `db:"membro_id"`
	Data             time.Time `db:"data"`
	Valor            float64   `db:"valor"`
	ValorDistribuido float64   `db:"valor_distribuido"`
	CreatedAt        time.Time `db:"created_at"`
}

func (r rendimentoInvestimentoRow) toDomain() *domain.RendimentoInvestimento {
	return &domain.RendimentoInvestimento{
		ID:               r.ID,
		Descricao:        r.Descricao,
		MembroID:         r.MembroID,
		Data:             r.Data,
		Valor:            r.Valor,
		ValorDistribuido: r.ValorDistribuido,
		CreatedAt:        r.CreatedAt,
	}
}

// Criar persiste um novo rendimento de investimento no banco e retorna o registro criado (com ID gerado).
func (r *RendimentoInvestimentoRepository) Criar(ri *domain.RendimentoInvestimento) (*domain.RendimentoInvestimento, error) {
	var row rendimentoInvestimentoRow
	err := r.db.QueryRowx(`
		INSERT INTO rendimentos_investimento (descricao, membro_id, data, valor, valor_distribuido)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, descricao, membro_id, data, valor, valor_distribuido, created_at
	`, ri.Descricao, ri.MembroID, ri.Data, ri.Valor, ri.ValorDistribuido).StructScan(&row)
	if err != nil {
		return nil, err
	}
	return row.toDomain(), nil
}

// BuscarPorID busca um rendimento de investimento pelo seu UUID.
// Retorna domain.ErrRendimentoInvestimentoNaoEncontrado se não existir ou estiver soft-deleted.
func (r *RendimentoInvestimentoRepository) BuscarPorID(id string) (*domain.RendimentoInvestimento, error) {
	var row rendimentoInvestimentoRow
	err := r.db.QueryRowx(`
		SELECT id, descricao, membro_id, data, valor, valor_distribuido, created_at
		FROM rendimentos_investimento
		WHERE id = $1 AND deleted_at IS NULL
	`, id).StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRendimentoInvestimentoNaoEncontrado
		}
		return nil, err
	}
	return row.toDomain(), nil
}

// Listar retorna todos os rendimentos de investimento não excluídos, ordenados por data decrescente.
func (r *RendimentoInvestimentoRepository) Listar() ([]*domain.RendimentoInvestimento, error) {
	var rows []rendimentoInvestimentoRow
	err := r.db.Select(&rows, `
		SELECT id, descricao, membro_id, data, valor, valor_distribuido, created_at
		FROM rendimentos_investimento
		WHERE deleted_at IS NULL
		ORDER BY data DESC
	`)
	if err != nil {
		return nil, err
	}

	rendimentos := make([]*domain.RendimentoInvestimento, 0, len(rows))
	for _, row := range rows {
		rendimentos = append(rendimentos, row.toDomain())
	}
	return rendimentos, nil
}

// ListarPorMes retorna os rendimentos cujo mês e ano de data correspondam aos parâmetros.
func (r *RendimentoInvestimentoRepository) ListarPorMes(mes, ano int) ([]*domain.RendimentoInvestimento, error) {
	var rows []rendimentoInvestimentoRow
	err := r.db.Select(&rows, `
		SELECT id, descricao, membro_id, data, valor, valor_distribuido, created_at
		FROM rendimentos_investimento
		WHERE deleted_at IS NULL
		  AND EXTRACT(MONTH FROM data) = $1
		  AND EXTRACT(YEAR FROM data) = $2
		ORDER BY data DESC
	`, mes, ano)
	if err != nil {
		return nil, err
	}

	rendimentos := make([]*domain.RendimentoInvestimento, 0, len(rows))
	for _, row := range rows {
		rendimentos = append(rendimentos, row.toDomain())
	}
	return rendimentos, nil
}

// Atualizar atualiza os campos de um rendimento de investimento existente.
// Retorna o registro atualizado ou ErrRendimentoInvestimentoNaoEncontrado se não existir.
func (r *RendimentoInvestimentoRepository) Atualizar(ri *domain.RendimentoInvestimento) (*domain.RendimentoInvestimento, error) {
	var row rendimentoInvestimentoRow
	err := r.db.QueryRowx(`
		UPDATE rendimentos_investimento
		SET descricao=$2, membro_id=$3, data=$4, valor=$5, valor_distribuido=$6, updated_at=NOW()
		WHERE id=$1 AND deleted_at IS NULL
		RETURNING id, descricao, membro_id, data, valor, valor_distribuido, created_at
	`, ri.ID, ri.Descricao, ri.MembroID, ri.Data, ri.Valor, ri.ValorDistribuido).StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRendimentoInvestimentoNaoEncontrado
		}
		return nil, err
	}
	return row.toDomain(), nil
}

// Excluir realiza o soft-delete de um rendimento de investimento existente.
// Retorna ErrRendimentoInvestimentoNaoEncontrado se não existir ou já estiver excluído.
func (r *RendimentoInvestimentoRepository) Excluir(id string) error {
	result, err := r.db.Exec(`
		UPDATE rendimentos_investimento
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrRendimentoInvestimentoNaoEncontrado
	}
	return nil
}
