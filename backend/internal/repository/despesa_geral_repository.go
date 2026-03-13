package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

// DespesaGeralRepository implementa service.DespesaGeralRepositoryInterface usando PostgreSQL via sqlx.
type DespesaGeralRepository struct {
	db *sqlx.DB
}

// NewDespesaGeralRepository cria uma nova instância do DespesaGeralRepository.
func NewDespesaGeralRepository(db *sqlx.DB) *DespesaGeralRepository {
	return &DespesaGeralRepository{db: db}
}

// despesaGeralRow é a estrutura de scan para as colunas do banco.
type despesaGeralRow struct {
	ID             string    `db:"id"`
	MembroID       string    `db:"membro_id"`
	CategoriaID    *string   `db:"categoria_id"`
	Descricao      string    `db:"descricao"`
	Data           time.Time `db:"data"`
	Valor          float64   `db:"valor"`
	FormaPagamento string    `db:"forma_pagamento"`
	Observacoes    *string   `db:"observacoes"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

func (r despesaGeralRow) toDomain() *domain.DespesaGeral {
	return &domain.DespesaGeral{
		ID:             r.ID,
		MembroID:       r.MembroID,
		CategoriaID:    r.CategoriaID,
		Descricao:      r.Descricao,
		Data:           r.Data,
		Valor:          r.Valor,
		FormaPagamento: r.FormaPagamento,
		Observacoes:    r.Observacoes,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}

// Criar persiste uma nova despesa geral no banco e retorna o registro criado (com ID gerado).
func (r *DespesaGeralRepository) Criar(familiaID string, d *domain.DespesaGeral) (*domain.DespesaGeral, error) {
	var row despesaGeralRow
	err := r.db.QueryRowx(`
		INSERT INTO despesas_gerais (familia_id, membro_id, categoria_id, descricao, data, valor, forma_pagamento, observacoes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, membro_id, categoria_id, descricao, data, valor, forma_pagamento, observacoes, created_at, updated_at
	`, familiaID, d.MembroID, d.CategoriaID, d.Descricao, d.Data, d.Valor, d.FormaPagamento, d.Observacoes).StructScan(&row)
	if err != nil {
		return nil, err
	}
	return row.toDomain(), nil
}

// BuscarPorID busca uma despesa geral pelo seu UUID.
// Retorna domain.ErrDespesaGeralNaoEncontrada se não existir ou estiver soft-deleted.
func (r *DespesaGeralRepository) BuscarPorID(familiaID, id string) (*domain.DespesaGeral, error) {
	var row despesaGeralRow
	err := r.db.QueryRowx(`
		SELECT id, membro_id, categoria_id, descricao, data, valor, forma_pagamento, observacoes, created_at, updated_at
		FROM despesas_gerais
		WHERE id = $1 AND familia_id = $2 AND deleted_at IS NULL
	`, id, familiaID).StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrDespesaGeralNaoEncontrada
		}
		return nil, err
	}
	return row.toDomain(), nil
}

// Listar retorna todas as despesas gerais não excluídas, ordenadas por data decrescente.
func (r *DespesaGeralRepository) Listar(familiaID string) ([]*domain.DespesaGeral, error) {
	var rows []despesaGeralRow
	err := r.db.Select(&rows, `
		SELECT id, membro_id, categoria_id, descricao, data, valor, forma_pagamento, observacoes, created_at, updated_at
		FROM despesas_gerais
		WHERE familia_id = $1 AND deleted_at IS NULL
		ORDER BY data DESC
	`, familiaID)
	if err != nil {
		return nil, err
	}

	despesas := make([]*domain.DespesaGeral, 0, len(rows))
	for _, row := range rows {
		despesas = append(despesas, row.toDomain())
	}
	return despesas, nil
}

// ListarPorMes retorna as despesas gerais do mês e ano especificados, ordenadas por data decrescente.
func (r *DespesaGeralRepository) ListarPorMes(familiaID string, mes, ano int) ([]*domain.DespesaGeral, error) {
	var rows []despesaGeralRow
	err := r.db.Select(&rows, `
		SELECT id, membro_id, categoria_id, descricao, data, valor, forma_pagamento, observacoes, created_at, updated_at
		FROM despesas_gerais
		WHERE familia_id = $1
		  AND deleted_at IS NULL
		  AND EXTRACT(MONTH FROM data) = $2
		  AND EXTRACT(YEAR FROM data) = $3
		ORDER BY data DESC
	`, familiaID, mes, ano)
	if err != nil {
		return nil, err
	}

	despesas := make([]*domain.DespesaGeral, 0, len(rows))
	for _, row := range rows {
		despesas = append(despesas, row.toDomain())
	}
	return despesas, nil
}

// Atualizar atualiza os campos de uma despesa geral existente.
// Retorna o registro atualizado ou ErrDespesaGeralNaoEncontrada se não existir.
func (r *DespesaGeralRepository) Atualizar(familiaID string, d *domain.DespesaGeral) (*domain.DespesaGeral, error) {
	var row despesaGeralRow
	err := r.db.QueryRowx(`
		UPDATE despesas_gerais
		SET membro_id=$3, categoria_id=$4, descricao=$5, data=$6, valor=$7, forma_pagamento=$8, observacoes=$9, updated_at=NOW()
		WHERE id=$2 AND familia_id=$1 AND deleted_at IS NULL
		RETURNING id, membro_id, categoria_id, descricao, data, valor, forma_pagamento, observacoes, created_at, updated_at
	`, familiaID, d.ID, d.MembroID, d.CategoriaID, d.Descricao, d.Data, d.Valor, d.FormaPagamento, d.Observacoes).StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrDespesaGeralNaoEncontrada
		}
		return nil, err
	}
	return row.toDomain(), nil
}

// Excluir realiza o soft-delete de uma despesa geral existente.
// Retorna ErrDespesaGeralNaoEncontrada se não existir ou já estiver excluída.
func (r *DespesaGeralRepository) Excluir(familiaID, id string) error {
	result, err := r.db.Exec(`
		UPDATE despesas_gerais
		SET deleted_at = NOW()
		WHERE id = $2 AND familia_id = $1 AND deleted_at IS NULL
	`, familiaID, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrDespesaGeralNaoEncontrada
	}
	return nil
}
