package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

// ContaFixaRepository implementa service.ContaFixaRepository usando PostgreSQL via sqlx.
type ContaFixaRepository struct {
	db *sqlx.DB
}

// NewContaFixaRepository cria uma nova instância do ContaFixaRepository.
func NewContaFixaRepository(db *sqlx.DB) *ContaFixaRepository {
	return &ContaFixaRepository{db: db}
}

// contaFixaRow é a estrutura de scan para as colunas do banco.
type contaFixaRow struct {
	ID             string    `db:"id"`
	Descricao      string    `db:"descricao"`
	MembroID       string    `db:"membro_id"`
	CategoriaID    *string   `db:"categoria_id"`
	Valor          float64   `db:"valor"`
	DiaVencimento  int       `db:"dia_vencimento"`
	FormaPagamento string    `db:"forma_pagamento"`
	Ativa          bool      `db:"ativa"`
	CreatedAt      time.Time `db:"created_at"`
}

func (r contaFixaRow) toDomain() *domain.ContaFixa {
	return &domain.ContaFixa{
		ID:             r.ID,
		Descricao:      r.Descricao,
		MembroID:       r.MembroID,
		CategoriaID:    r.CategoriaID,
		Valor:          r.Valor,
		DiaVencimento:  r.DiaVencimento,
		FormaPagamento: r.FormaPagamento,
		Ativa:          r.Ativa,
		CreatedAt:      r.CreatedAt,
	}
}

// Criar persiste uma nova conta fixa no banco e retorna o registro criado (com ID gerado).
func (r *ContaFixaRepository) Criar(c *domain.ContaFixa) (*domain.ContaFixa, error) {
	var row contaFixaRow
	err := r.db.QueryRowx(`
		INSERT INTO contas_fixas (descricao, membro_id, categoria_id, valor, dia_vencimento, forma_pagamento, ativa)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, descricao, membro_id, categoria_id, valor, dia_vencimento, forma_pagamento, ativa, created_at
	`, c.Descricao, c.MembroID, c.CategoriaID, c.Valor, c.DiaVencimento, c.FormaPagamento, c.Ativa).StructScan(&row)
	if err != nil {
		return nil, err
	}
	return row.toDomain(), nil
}

// BuscarPorID busca uma conta fixa pelo seu UUID.
// Retorna domain.ErrContaFixaNaoEncontrada se não existir ou estiver soft-deleted.
func (r *ContaFixaRepository) BuscarPorID(id string) (*domain.ContaFixa, error) {
	var row contaFixaRow
	err := r.db.QueryRowx(`
		SELECT id, descricao, membro_id, categoria_id, valor, dia_vencimento, forma_pagamento, ativa, created_at
		FROM contas_fixas
		WHERE id = $1 AND deleted_at IS NULL
	`, id).StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrContaFixaNaoEncontrada
		}
		return nil, err
	}
	return row.toDomain(), nil
}

// Listar retorna todas as contas fixas não excluídas, ordenadas por descrição.
func (r *ContaFixaRepository) Listar() ([]*domain.ContaFixa, error) {
	var rows []contaFixaRow
	err := r.db.Select(&rows, `
		SELECT id, descricao, membro_id, categoria_id, valor, dia_vencimento, forma_pagamento, ativa, created_at
		FROM contas_fixas
		WHERE deleted_at IS NULL
		ORDER BY descricao ASC
	`)
	if err != nil {
		return nil, err
	}

	contas := make([]*domain.ContaFixa, 0, len(rows))
	for _, row := range rows {
		contas = append(contas, row.toDomain())
	}
	return contas, nil
}

// Atualizar atualiza os campos de uma conta fixa existente.
// Retorna o registro atualizado ou ErrContaFixaNaoEncontrada se não existir.
func (r *ContaFixaRepository) Atualizar(c *domain.ContaFixa) (*domain.ContaFixa, error) {
	var row contaFixaRow
	err := r.db.QueryRowx(`
		UPDATE contas_fixas
		SET descricao=$2, membro_id=$3, categoria_id=$4, valor=$5, dia_vencimento=$6, forma_pagamento=$7, ativa=$8, updated_at=NOW()
		WHERE id=$1 AND deleted_at IS NULL
		RETURNING id, descricao, membro_id, categoria_id, valor, dia_vencimento, forma_pagamento, ativa, created_at
	`, c.ID, c.Descricao, c.MembroID, c.CategoriaID, c.Valor, c.DiaVencimento, c.FormaPagamento, c.Ativa).StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrContaFixaNaoEncontrada
		}
		return nil, err
	}
	return row.toDomain(), nil
}

// ListarAtivas retorna todas as contas fixas com ativa = true e não excluídas.
func (r *ContaFixaRepository) ListarAtivas() ([]*domain.ContaFixa, error) {
	var rows []contaFixaRow
	err := r.db.Select(&rows, `
		SELECT id, descricao, membro_id, categoria_id, valor, dia_vencimento, forma_pagamento, ativa, created_at
		FROM contas_fixas
		WHERE deleted_at IS NULL
		  AND ativa = true
		ORDER BY descricao ASC
	`)
	if err != nil {
		return nil, err
	}
	contas := make([]*domain.ContaFixa, 0, len(rows))
	for _, row := range rows {
		contas = append(contas, row.toDomain())
	}
	return contas, nil
}

// AlterarAtivo atualiza apenas o estado ativa de uma conta fixa existente.
// Retorna o registro atualizado ou ErrContaFixaNaoEncontrada se não existir.
func (r *ContaFixaRepository) AlterarAtivo(id string, ativa bool) (*domain.ContaFixa, error) {
	var row contaFixaRow
	err := r.db.QueryRowx(`
		UPDATE contas_fixas
		SET ativa=$2, updated_at=NOW()
		WHERE id=$1 AND deleted_at IS NULL
		RETURNING id, descricao, membro_id, categoria_id, valor, dia_vencimento, forma_pagamento, ativa, created_at
	`, id, ativa).StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrContaFixaNaoEncontrada
		}
		return nil, err
	}
	return row.toDomain(), nil
}
