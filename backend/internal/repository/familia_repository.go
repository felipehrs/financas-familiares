package repository

import (
	"database/sql"
	"errors"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type FamiliaRepository struct {
	db *sqlx.DB
}

func NewFamiliaRepository(db *sqlx.DB) *FamiliaRepository {
	return &FamiliaRepository{db: db}
}

func (r *FamiliaRepository) BuscarFamiliaPorUsuario(usuarioID string) (string, error) {
	var familiaID string
	err := r.db.QueryRowx(`
		SELECT familia_id FROM familia_usuarios
		WHERE usuario_id = $1
		LIMIT 1
	`, usuarioID).Scan(&familiaID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.ErrFamiliaNaoEncontrada
		}
		return "", err
	}
	return familiaID, nil
}
