package repository

import (
	"database/sql"
	"errors"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

// CartaoCreditoRepository implementa service.CartaoCreditoRepository usando PostgreSQL via sqlx.
type CartaoCreditoRepository struct {
	db *sqlx.DB
}

// NewCartaoCreditoRepository cria uma nova instância do CartaoCreditoRepository.
func NewCartaoCreditoRepository(db *sqlx.DB) *CartaoCreditoRepository {
	return &CartaoCreditoRepository{db: db}
}

// cartaoCreditoRow é a estrutura de scan para as colunas do banco.
type cartaoCreditoRow struct {
	ID            string   `db:"id"`
	Nome          string   `db:"nome"`
	MembroID      string   `db:"membro_id"`
	DiaFechamento int      `db:"dia_fechamento"`
	DiaVencimento int      `db:"dia_vencimento"`
	Limite        *float64 `db:"limite"`
	Ativo         bool     `db:"ativo"`
}

func (r cartaoCreditoRow) toDomain() *domain.CartaoCredito {
	return &domain.CartaoCredito{
		ID:            r.ID,
		Nome:          r.Nome,
		MembroID:      r.MembroID,
		DiaFechamento: r.DiaFechamento,
		DiaVencimento: r.DiaVencimento,
		Limite:        r.Limite,
		Ativo:         r.Ativo,
	}
}

// Criar persiste um novo cartão de crédito no banco e retorna o registro criado (com ID gerado).
func (r *CartaoCreditoRepository) Criar(cartao *domain.CartaoCredito) (*domain.CartaoCredito, error) {
	var row cartaoCreditoRow
	err := r.db.QueryRowx(`
		INSERT INTO cartoes_credito (nome, membro_id, dia_fechamento, dia_vencimento, limite, ativo)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, nome, membro_id, dia_fechamento, dia_vencimento, limite, ativo
	`, cartao.Nome, cartao.MembroID, cartao.DiaFechamento, cartao.DiaVencimento, cartao.Limite, cartao.Ativo).StructScan(&row)
	if err != nil {
		return nil, err
	}
	return row.toDomain(), nil
}

// BuscarPorID busca um cartão pelo seu UUID.
// Retorna domain.ErrCartaoNaoEncontrado se não existir ou estiver soft-deleted.
func (r *CartaoCreditoRepository) BuscarPorID(id string) (*domain.CartaoCredito, error) {
	var row cartaoCreditoRow
	err := r.db.QueryRowx(`
		SELECT id, nome, membro_id, dia_fechamento, dia_vencimento, limite, ativo
		FROM cartoes_credito
		WHERE id = $1 AND deleted_at IS NULL
	`, id).StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCartaoNaoEncontrado
		}
		return nil, err
	}
	return row.toDomain(), nil
}

// Listar retorna todos os cartões não excluídos (ativos e inativos).
func (r *CartaoCreditoRepository) Listar() ([]*domain.CartaoCredito, error) {
	var rows []cartaoCreditoRow
	err := r.db.Select(&rows, `
		SELECT id, nome, membro_id, dia_fechamento, dia_vencimento, limite, ativo
		FROM cartoes_credito
		WHERE deleted_at IS NULL
		ORDER BY nome ASC
	`)
	if err != nil {
		return nil, err
	}

	cartoes := make([]*domain.CartaoCredito, 0, len(rows))
	for _, row := range rows {
		cartoes = append(cartoes, row.toDomain())
	}
	return cartoes, nil
}

// Atualizar atualiza os dados de um cartão existente e retorna o registro atualizado.
func (r *CartaoCreditoRepository) Atualizar(cartao *domain.CartaoCredito) (*domain.CartaoCredito, error) {
	var row cartaoCreditoRow
	err := r.db.QueryRowx(`
		UPDATE cartoes_credito
		SET nome = $2, membro_id = $3, dia_fechamento = $4, dia_vencimento = $5, limite = $6, ativo = $7, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, nome, membro_id, dia_fechamento, dia_vencimento, limite, ativo
	`, cartao.ID, cartao.Nome, cartao.MembroID, cartao.DiaFechamento, cartao.DiaVencimento, cartao.Limite, cartao.Ativo).StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCartaoNaoEncontrado
		}
		return nil, err
	}
	return row.toDomain(), nil
}

// Inativar marca um cartão como inativo sem excluir o registro.
func (r *CartaoCreditoRepository) Inativar(id string) error {
	result, err := r.db.Exec(`
		UPDATE cartoes_credito
		SET ativo = FALSE, updated_at = NOW()
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
		return domain.ErrCartaoNaoEncontrado
	}
	return nil
}
