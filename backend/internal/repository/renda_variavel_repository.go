package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

// RendaVariavelRepository implementa service.RendaVariavelRepositoryInterface usando PostgreSQL via sqlx.
type RendaVariavelRepository struct {
	db *sqlx.DB
}

// NewRendaVariavelRepository cria uma nova instância do RendaVariavelRepository.
func NewRendaVariavelRepository(db *sqlx.DB) *RendaVariavelRepository {
	return &RendaVariavelRepository{db: db}
}

// rendaVariavelRow é a estrutura de scan para as colunas do banco.
type rendaVariavelRow struct {
	ID              string    `db:"id"`
	Descricao       string    `db:"descricao"`
	MembroID        string    `db:"membro_id"`
	MesReferencia   int       `db:"mes_referencia"`
	AnoReferencia   int       `db:"ano_referencia"`
	Valor           float64   `db:"valor"`
	DataRecebimento time.Time `db:"data_recebimento"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

func (r rendaVariavelRow) toDomain() *domain.RendaVariavel {
	return &domain.RendaVariavel{
		ID:              r.ID,
		Descricao:       r.Descricao,
		MembroID:        r.MembroID,
		MesReferencia:   r.MesReferencia,
		AnoReferencia:   r.AnoReferencia,
		Valor:           r.Valor,
		DataRecebimento: r.DataRecebimento,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}

// Criar persiste uma nova renda variável no banco e retorna o registro criado (com ID gerado).
func (r *RendaVariavelRepository) Criar(rv *domain.RendaVariavel) (*domain.RendaVariavel, error) {
	var row rendaVariavelRow
	err := r.db.QueryRowx(`
		INSERT INTO rendas_variaveis (descricao, membro_id, mes_referencia, ano_referencia, valor, data_recebimento)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, descricao, membro_id, mes_referencia, ano_referencia, valor, data_recebimento, created_at, updated_at
	`, rv.Descricao, rv.MembroID, rv.MesReferencia, rv.AnoReferencia, rv.Valor, rv.DataRecebimento).StructScan(&row)
	if err != nil {
		return nil, err
	}
	return row.toDomain(), nil
}

// BuscarPorID busca uma renda variável pelo seu UUID.
// Retorna domain.ErrRendaVariavelNaoEncontrada se não existir ou estiver soft-deleted.
func (r *RendaVariavelRepository) BuscarPorID(id string) (*domain.RendaVariavel, error) {
	var row rendaVariavelRow
	err := r.db.QueryRowx(`
		SELECT id, descricao, membro_id, mes_referencia, ano_referencia, valor, data_recebimento, created_at, updated_at
		FROM rendas_variaveis
		WHERE id = $1 AND deleted_at IS NULL
	`, id).StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRendaVariavelNaoEncontrada
		}
		return nil, err
	}
	return row.toDomain(), nil
}

// Listar retorna todas as rendas variáveis não excluídas, ordenadas por ano e mês de referência decrescentes.
func (r *RendaVariavelRepository) Listar() ([]*domain.RendaVariavel, error) {
	var rows []rendaVariavelRow
	err := r.db.Select(&rows, `
		SELECT id, descricao, membro_id, mes_referencia, ano_referencia, valor, data_recebimento, created_at, updated_at
		FROM rendas_variaveis
		WHERE deleted_at IS NULL
		ORDER BY ano_referencia DESC, mes_referencia DESC
	`)
	if err != nil {
		return nil, err
	}

	rendas := make([]*domain.RendaVariavel, 0, len(rows))
	for _, row := range rows {
		rendas = append(rendas, row.toDomain())
	}
	return rendas, nil
}

// ListarPorMes retorna as rendas variáveis do mês e ano de referência especificados.
func (r *RendaVariavelRepository) ListarPorMes(mes, ano int) ([]*domain.RendaVariavel, error) {
	var rows []rendaVariavelRow
	err := r.db.Select(&rows, `
		SELECT id, descricao, membro_id, mes_referencia, ano_referencia, valor, data_recebimento, created_at, updated_at
		FROM rendas_variaveis
		WHERE deleted_at IS NULL
		  AND mes_referencia = $1
		  AND ano_referencia = $2
		ORDER BY data_recebimento DESC
	`, mes, ano)
	if err != nil {
		return nil, err
	}

	rendas := make([]*domain.RendaVariavel, 0, len(rows))
	for _, row := range rows {
		rendas = append(rendas, row.toDomain())
	}
	return rendas, nil
}

// Atualizar atualiza os campos de uma renda variável existente.
// Retorna o registro atualizado ou ErrRendaVariavelNaoEncontrada se não existir.
func (r *RendaVariavelRepository) Atualizar(rv *domain.RendaVariavel) (*domain.RendaVariavel, error) {
	var row rendaVariavelRow
	err := r.db.QueryRowx(`
		UPDATE rendas_variaveis
		SET descricao=$2, membro_id=$3, mes_referencia=$4, ano_referencia=$5, valor=$6, data_recebimento=$7, updated_at=NOW()
		WHERE id=$1 AND deleted_at IS NULL
		RETURNING id, descricao, membro_id, mes_referencia, ano_referencia, valor, data_recebimento, created_at, updated_at
	`, rv.ID, rv.Descricao, rv.MembroID, rv.MesReferencia, rv.AnoReferencia, rv.Valor, rv.DataRecebimento).StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRendaVariavelNaoEncontrada
		}
		return nil, err
	}
	return row.toDomain(), nil
}

// Excluir realiza o soft-delete de uma renda variável existente.
// Retorna ErrRendaVariavelNaoEncontrada se não existir ou já estiver excluída.
func (r *RendaVariavelRepository) Excluir(id string) error {
	result, err := r.db.Exec(`
		UPDATE rendas_variaveis
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
		return domain.ErrRendaVariavelNaoEncontrada
	}
	return nil
}
