package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

// AssinaturaRepository implementa service.AssinaturaRepository usando PostgreSQL via sqlx.
type AssinaturaRepository struct {
	db *sqlx.DB
}

// NewAssinaturaRepository cria uma nova instância do AssinaturaRepository.
func NewAssinaturaRepository(db *sqlx.DB) *AssinaturaRepository {
	return &AssinaturaRepository{db: db}
}

// assinaturaRow é a estrutura de scan para as colunas do banco.
type assinaturaRow struct {
	ID             string    `db:"id"`
	Nome           string    `db:"nome"`
	MembroID       string    `db:"membro_id"`
	CategoriaID    *string   `db:"categoria_id"`
	Valor          float64   `db:"valor"`
	DiaCobranca    int       `db:"dia_cobranca"`
	FormaPagamento string    `db:"forma_pagamento"`
	Status         string    `db:"status"`
	CreatedAt      time.Time `db:"created_at"`
}

func (r assinaturaRow) toDomain() *domain.Assinatura {
	return &domain.Assinatura{
		ID:             r.ID,
		Nome:           r.Nome,
		MembroID:       r.MembroID,
		CategoriaID:    r.CategoriaID,
		Valor:          r.Valor,
		DiaCobranca:    r.DiaCobranca,
		FormaPagamento: r.FormaPagamento,
		Status:         r.Status,
		CreatedAt:      r.CreatedAt,
	}
}

// Criar persiste uma nova assinatura no banco e retorna o registro criado (com ID gerado).
func (r *AssinaturaRepository) Criar(familiaID string, a *domain.Assinatura) (*domain.Assinatura, error) {
	var row assinaturaRow
	err := r.db.QueryRowx(`
		INSERT INTO assinaturas (familia_id, nome, membro_id, categoria_id, valor, dia_cobranca, forma_pagamento, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, nome, membro_id, categoria_id, valor, dia_cobranca, forma_pagamento, status, created_at
	`, familiaID, a.Nome, a.MembroID, a.CategoriaID, a.Valor, a.DiaCobranca, a.FormaPagamento, a.Status).StructScan(&row)
	if err != nil {
		return nil, err
	}
	return row.toDomain(), nil
}

// BuscarPorID busca uma assinatura pelo seu UUID.
// Retorna domain.ErrAssinaturaNaoEncontrada se não existir ou estiver soft-deleted.
func (r *AssinaturaRepository) BuscarPorID(familiaID, id string) (*domain.Assinatura, error) {
	var row assinaturaRow
	err := r.db.QueryRowx(`
		SELECT id, nome, membro_id, categoria_id, valor, dia_cobranca, forma_pagamento, status, created_at
		FROM assinaturas
		WHERE id = $1 AND familia_id = $2 AND deleted_at IS NULL
	`, id, familiaID).StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrAssinaturaNaoEncontrada
		}
		return nil, err
	}
	return row.toDomain(), nil
}

// Listar retorna todas as assinaturas não excluídas, ordenadas por nome.
func (r *AssinaturaRepository) Listar(familiaID string) ([]*domain.Assinatura, error) {
	var rows []assinaturaRow
	err := r.db.Select(&rows, `
		SELECT id, nome, membro_id, categoria_id, valor, dia_cobranca, forma_pagamento, status, created_at
		FROM assinaturas
		WHERE familia_id = $1 AND deleted_at IS NULL
		ORDER BY nome ASC
	`, familiaID)
	if err != nil {
		return nil, err
	}

	assinaturas := make([]*domain.Assinatura, 0, len(rows))
	for _, row := range rows {
		assinaturas = append(assinaturas, row.toDomain())
	}
	return assinaturas, nil
}

// Atualizar atualiza os campos de uma assinatura existente.
// Retorna o registro atualizado ou ErrAssinaturaNaoEncontrada se não existir.
func (r *AssinaturaRepository) Atualizar(familiaID string, a *domain.Assinatura) (*domain.Assinatura, error) {
	var row assinaturaRow
	err := r.db.QueryRowx(`
		UPDATE assinaturas
		SET nome = $3, membro_id = $4, categoria_id = $5, valor = $6, dia_cobranca = $7, forma_pagamento = $8, status = $9, updated_at = NOW()
		WHERE id = $2 AND familia_id = $1 AND deleted_at IS NULL
		RETURNING id, nome, membro_id, categoria_id, valor, dia_cobranca, forma_pagamento, status, created_at
	`, familiaID, a.ID, a.Nome, a.MembroID, a.CategoriaID, a.Valor, a.DiaCobranca, a.FormaPagamento, a.Status).StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrAssinaturaNaoEncontrada
		}
		return nil, err
	}
	return row.toDomain(), nil
}

// ListarAtivas retorna todas as assinaturas com status "ativa" e não excluídas.
func (r *AssinaturaRepository) ListarAtivas(familiaID string) ([]*domain.Assinatura, error) {
	var rows []assinaturaRow
	err := r.db.Select(&rows, `
		SELECT id, nome, membro_id, categoria_id, valor, dia_cobranca, forma_pagamento, status, created_at
		FROM assinaturas
		WHERE familia_id = $1 AND deleted_at IS NULL
		  AND status = 'ativa'
		ORDER BY nome ASC
	`, familiaID)
	if err != nil {
		return nil, err
	}
	assinaturas := make([]*domain.Assinatura, 0, len(rows))
	for _, row := range rows {
		assinaturas = append(assinaturas, row.toDomain())
	}
	return assinaturas, nil
}

// AlterarStatus atualiza apenas o status de uma assinatura existente.
// Retorna o registro atualizado ou ErrAssinaturaNaoEncontrada se não existir.
func (r *AssinaturaRepository) AlterarStatus(familiaID, id, status string) (*domain.Assinatura, error) {
	var row assinaturaRow
	err := r.db.QueryRowx(`
		UPDATE assinaturas
		SET status = $3, updated_at = NOW()
		WHERE id = $2 AND familia_id = $1 AND deleted_at IS NULL
		RETURNING id, nome, membro_id, categoria_id, valor, dia_cobranca, forma_pagamento, status, created_at
	`, familiaID, id, status).StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrAssinaturaNaoEncontrada
		}
		return nil, err
	}
	return row.toDomain(), nil
}
