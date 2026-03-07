package repository

import (
	"database/sql"
	"errors"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

// CategoriaRepository implementa service.CategoriaRepository usando PostgreSQL via sqlx.
type CategoriaRepository struct {
	db *sqlx.DB
}

// NewCategoriaRepository cria uma nova instância do CategoriaRepository.
func NewCategoriaRepository(db *sqlx.DB) *CategoriaRepository {
	return &CategoriaRepository{db: db}
}

// categoriaRow é a estrutura de scan para as colunas do banco.
type categoriaRow struct {
	ID   string `db:"id"`
	Nome string `db:"nome"`
}

func (r categoriaRow) toDomain() *domain.Categoria {
	return &domain.Categoria{
		ID:   r.ID,
		Nome: r.Nome,
	}
}

// Criar persiste uma nova categoria no banco e retorna o registro criado (com ID gerado).
func (r *CategoriaRepository) Criar(categoria *domain.Categoria) (*domain.Categoria, error) {
	var row categoriaRow
	err := r.db.QueryRowx(`
		INSERT INTO categorias (nome)
		VALUES ($1)
		RETURNING id, nome
	`, categoria.Nome).StructScan(&row)
	if err != nil {
		return nil, err
	}
	return row.toDomain(), nil
}

// Listar retorna todas as categorias não excluídas, ordenadas por nome.
func (r *CategoriaRepository) Listar() ([]*domain.Categoria, error) {
	var rows []categoriaRow
	err := r.db.Select(&rows, `
		SELECT id, nome
		FROM categorias
		WHERE deleted_at IS NULL
		ORDER BY nome ASC
	`)
	if err != nil {
		return nil, err
	}

	categorias := make([]*domain.Categoria, 0, len(rows))
	for _, row := range rows {
		categorias = append(categorias, row.toDomain())
	}
	return categorias, nil
}

// BuscarPorID busca uma categoria pelo seu UUID.
// Retorna domain.ErrCategoriaNaoEncontrada se não existir ou estiver soft-deleted.
func (r *CategoriaRepository) BuscarPorID(id string) (*domain.Categoria, error) {
	var row categoriaRow
	err := r.db.QueryRowx(`
		SELECT id, nome
		FROM categorias
		WHERE id = $1 AND deleted_at IS NULL
	`, id).StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCategoriaNaoEncontrada
		}
		return nil, err
	}
	return row.toDomain(), nil
}

// Atualizar atualiza o nome de uma categoria existente e retorna o registro atualizado.
func (r *CategoriaRepository) Atualizar(categoria *domain.Categoria) (*domain.Categoria, error) {
	var row categoriaRow
	err := r.db.QueryRowx(`
		UPDATE categorias
		SET nome = $2, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, nome
	`, categoria.ID, categoria.Nome).StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCategoriaNaoEncontrada
		}
		return nil, err
	}
	return row.toDomain(), nil
}

// Excluir realiza o soft delete de uma categoria.
// Antes verifica se há registros vinculados em despesas_cartao, assinaturas,
// contas_fixas e despesas_gerais. Se houver, retorna ErrCategoriaComVinculos.
func (r *CategoriaRepository) Excluir(id string) error {
	var count int
	err := r.db.QueryRowx(`
		SELECT COUNT(*) FROM (
			SELECT id FROM despesas_cartao WHERE categoria_id = $1 AND deleted_at IS NULL
			UNION ALL
			SELECT id FROM assinaturas WHERE categoria_id = $1 AND deleted_at IS NULL
			UNION ALL
			SELECT id FROM contas_fixas WHERE categoria_id = $1 AND deleted_at IS NULL
			UNION ALL
			SELECT id FROM despesas_gerais WHERE categoria_id = $1 AND deleted_at IS NULL
		) AS vinculos
	`, id).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return domain.ErrCategoriaComVinculos
	}

	result, err := r.db.Exec(`
		UPDATE categorias
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
		return domain.ErrCategoriaNaoEncontrada
	}
	return nil
}
