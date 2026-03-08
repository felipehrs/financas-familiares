package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

// RendaExtraRepository implementa service.RendaExtraRepositoryInterface usando PostgreSQL via sqlx.
type RendaExtraRepository struct {
	db *sqlx.DB
}

// NewRendaExtraRepository cria uma nova instância do RendaExtraRepository.
func NewRendaExtraRepository(db *sqlx.DB) *RendaExtraRepository {
	return &RendaExtraRepository{db: db}
}

// rendaExtraRow é a estrutura de scan para as colunas do banco.
type rendaExtraRow struct {
	ID              string    `db:"id"`
	Descricao       string    `db:"descricao"`
	MembroID        string    `db:"membro_id"`
	DataRecebimento time.Time `db:"data_recebimento"`
	Valor           float64   `db:"valor"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

func (r rendaExtraRow) toDomain() *domain.RendaExtra {
	return &domain.RendaExtra{
		ID:              r.ID,
		Descricao:       r.Descricao,
		MembroID:        r.MembroID,
		DataRecebimento: r.DataRecebimento,
		Valor:           r.Valor,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}

// Criar persiste uma nova renda extra no banco e retorna o registro criado (com ID gerado).
func (r *RendaExtraRepository) Criar(re *domain.RendaExtra) (*domain.RendaExtra, error) {
	var row rendaExtraRow
	err := r.db.QueryRowx(`
		INSERT INTO rendas_extras (descricao, membro_id, data_recebimento, valor)
		VALUES ($1, $2, $3, $4)
		RETURNING id, descricao, membro_id, data_recebimento, valor, created_at, updated_at
	`, re.Descricao, re.MembroID, re.DataRecebimento, re.Valor).StructScan(&row)
	if err != nil {
		return nil, err
	}
	return row.toDomain(), nil
}

// BuscarPorID busca uma renda extra pelo seu UUID.
// Retorna domain.ErrRendaExtraNaoEncontrada se não existir ou estiver soft-deleted.
func (r *RendaExtraRepository) BuscarPorID(id string) (*domain.RendaExtra, error) {
	var row rendaExtraRow
	err := r.db.QueryRowx(`
		SELECT id, descricao, membro_id, data_recebimento, valor, created_at, updated_at
		FROM rendas_extras
		WHERE id = $1 AND deleted_at IS NULL
	`, id).StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRendaExtraNaoEncontrada
		}
		return nil, err
	}
	return row.toDomain(), nil
}

// Listar retorna todas as rendas extras não excluídas, ordenadas por data_recebimento decrescente.
func (r *RendaExtraRepository) Listar() ([]*domain.RendaExtra, error) {
	var rows []rendaExtraRow
	err := r.db.Select(&rows, `
		SELECT id, descricao, membro_id, data_recebimento, valor, created_at, updated_at
		FROM rendas_extras
		WHERE deleted_at IS NULL
		ORDER BY data_recebimento DESC
	`)
	if err != nil {
		return nil, err
	}

	rendas := make([]*domain.RendaExtra, 0, len(rows))
	for _, row := range rows {
		rendas = append(rendas, row.toDomain())
	}
	return rendas, nil
}

// ListarPorMes retorna as rendas extras cujo mês e ano de data_recebimento correspondam aos parâmetros.
func (r *RendaExtraRepository) ListarPorMes(mes, ano int) ([]*domain.RendaExtra, error) {
	var rows []rendaExtraRow
	err := r.db.Select(&rows, `
		SELECT id, descricao, membro_id, data_recebimento, valor, created_at, updated_at
		FROM rendas_extras
		WHERE deleted_at IS NULL
		  AND EXTRACT(MONTH FROM data_recebimento) = $1
		  AND EXTRACT(YEAR FROM data_recebimento) = $2
		ORDER BY data_recebimento DESC
	`, mes, ano)
	if err != nil {
		return nil, err
	}

	rendas := make([]*domain.RendaExtra, 0, len(rows))
	for _, row := range rows {
		rendas = append(rendas, row.toDomain())
	}
	return rendas, nil
}

// Atualizar atualiza os campos de uma renda extra existente.
// Retorna o registro atualizado ou ErrRendaExtraNaoEncontrada se não existir.
func (r *RendaExtraRepository) Atualizar(re *domain.RendaExtra) (*domain.RendaExtra, error) {
	var row rendaExtraRow
	err := r.db.QueryRowx(`
		UPDATE rendas_extras
		SET descricao=$2, membro_id=$3, data_recebimento=$4, valor=$5, updated_at=NOW()
		WHERE id=$1 AND deleted_at IS NULL
		RETURNING id, descricao, membro_id, data_recebimento, valor, created_at, updated_at
	`, re.ID, re.Descricao, re.MembroID, re.DataRecebimento, re.Valor).StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRendaExtraNaoEncontrada
		}
		return nil, err
	}
	return row.toDomain(), nil
}

// Excluir realiza o soft-delete de uma renda extra existente.
// Retorna ErrRendaExtraNaoEncontrada se não existir ou já estiver excluída.
func (r *RendaExtraRepository) Excluir(id string) error {
	result, err := r.db.Exec(`
		UPDATE rendas_extras
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
		return domain.ErrRendaExtraNaoEncontrada
	}
	return nil
}
