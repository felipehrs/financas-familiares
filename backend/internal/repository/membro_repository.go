package repository

import (
	"database/sql"
	"errors"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

// MembroRepository implementa service.MembroRepository usando PostgreSQL via sqlx.
type MembroRepository struct {
	db *sqlx.DB
}

// NewMembroRepository cria uma nova instância do MembroRepository.
func NewMembroRepository(db *sqlx.DB) *MembroRepository {
	return &MembroRepository{db: db}
}

// membroRow é a estrutura de scan para as colunas do banco.
type membroRow struct {
	ID             string `db:"id"`
	Nome           string `db:"nome"`
	Relacionamento string `db:"relacionamento"`
	Ativo          bool   `db:"ativo"`
}

func (r membroRow) toDomain() *domain.Membro {
	return &domain.Membro{
		ID:             r.ID,
		Nome:           r.Nome,
		Relacionamento: r.Relacionamento,
		Ativo:          r.Ativo,
	}
}

// Criar persiste um novo membro no banco e retorna o registro criado (com ID gerado).
func (r *MembroRepository) Criar(familiaID string, membro *domain.Membro) (*domain.Membro, error) {
	var row membroRow
	err := r.db.QueryRowx(`
		INSERT INTO membros (familia_id, nome, relacionamento, ativo)
		VALUES ($1, $2, $3, $4)
		RETURNING id, nome, relacionamento, ativo
	`, familiaID, membro.Nome, membro.Relacionamento, membro.Ativo).StructScan(&row)
	if err != nil {
		return nil, err
	}
	return row.toDomain(), nil
}

// BuscarPorID busca um membro pelo seu UUID.
// Retorna domain.ErrMembroNaoEncontrado se não existir ou estiver soft-deleted.
func (r *MembroRepository) BuscarPorID(familiaID, id string) (*domain.Membro, error) {
	var row membroRow
	err := r.db.QueryRowx(`
		SELECT id, nome, relacionamento, ativo
		FROM membros
		WHERE id = $1 AND familia_id = $2 AND deleted_at IS NULL
	`, id, familiaID).StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrMembroNaoEncontrado
		}
		return nil, err
	}
	return row.toDomain(), nil
}

// Listar retorna todos os membros não excluídos (ativos e inativos).
func (r *MembroRepository) Listar(familiaID string) ([]*domain.Membro, error) {
	var rows []membroRow
	err := r.db.Select(&rows, `
		SELECT id, nome, relacionamento, ativo
		FROM membros
		WHERE familia_id = $1 AND deleted_at IS NULL
		ORDER BY nome ASC
	`, familiaID)
	if err != nil {
		return nil, err
	}

	membros := make([]*domain.Membro, 0, len(rows))
	for _, row := range rows {
		membros = append(membros, row.toDomain())
	}
	return membros, nil
}

// Atualizar atualiza nome, relacionamento e ativo de um membro existente.
// Retorna o registro atualizado.
func (r *MembroRepository) Atualizar(familiaID string, membro *domain.Membro) (*domain.Membro, error) {
	var row membroRow
	err := r.db.QueryRowx(`
		UPDATE membros
		SET nome = $3, relacionamento = $4, ativo = $5, updated_at = NOW()
		WHERE id = $2 AND familia_id = $1 AND deleted_at IS NULL
		RETURNING id, nome, relacionamento, ativo
	`, familiaID, membro.ID, membro.Nome, membro.Relacionamento, membro.Ativo).StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrMembroNaoEncontrado
		}
		return nil, err
	}
	return row.toDomain(), nil
}

// Inativar marca um membro como inativo sem excluir o registro.
// Atualiza ativo=false e updated_at=NOW().
func (r *MembroRepository) Inativar(familiaID, id string) error {
	result, err := r.db.Exec(`
		UPDATE membros
		SET ativo = FALSE, updated_at = NOW()
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
		return domain.ErrMembroNaoEncontrado
	}
	return nil
}
