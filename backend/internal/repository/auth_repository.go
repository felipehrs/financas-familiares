package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/felipehrs/financas-familiares/backend/internal/service"
	"github.com/jmoiron/sqlx"
)

// AuthRepository implementa service.AuthRepository usando PostgreSQL via sqlx.
type AuthRepository struct {
	db *sqlx.DB
}

// NewAuthRepository cria uma nova instância do AuthRepository.
func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// FindByEmail busca um usuário pelo email.
// Retorna domain.ErrCredenciaisInvalidas se não encontrado.
func (r *AuthRepository) FindByEmail(email string) (*domain.Usuario, error) {
	var row struct {
		ID        string `db:"id"`
		Nome      string `db:"nome"`
		Email     string `db:"email"`
		SenhaHash string `db:"senha_hash"`
	}

	err := r.db.Get(&row, `
		SELECT id, nome, email, senha_hash
		FROM usuarios
		WHERE email = $1 AND deleted_at IS NULL
	`, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCredenciaisInvalidas
		}
		return nil, err
	}

	return &domain.Usuario{
		ID:        row.ID,
		Nome:      row.Nome,
		Email:     row.Email,
		SenhaHash: row.SenhaHash,
	}, nil
}

// SaveRefreshToken persiste um novo refresh token (já hasheado) no banco.
func (r *AuthRepository) SaveRefreshToken(usuarioID, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.Exec(`
		INSERT INTO refresh_tokens (usuario_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, usuarioID, tokenHash, expiresAt)
	return err
}

// FindRefreshToken busca um refresh token pelo hash.
// Retorna nil, nil se não encontrado.
func (r *AuthRepository) FindRefreshToken(tokenHash string) (*service.RefreshTokenRecord, error) {
	var row struct {
		UsuarioID string    `db:"usuario_id"`
		TokenHash string    `db:"token_hash"`
		ExpiresAt time.Time `db:"expires_at"`
		Revogado  bool      `db:"revogado"`
	}

	err := r.db.Get(&row, `
		SELECT usuario_id, token_hash, expires_at, revogado
		FROM refresh_tokens
		WHERE token_hash = $1 AND deleted_at IS NULL
	`, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &service.RefreshTokenRecord{
		UsuarioID: row.UsuarioID,
		TokenHash: row.TokenHash,
		ExpiresAt: row.ExpiresAt,
		Revogado:  row.Revogado,
	}, nil
}

// RevokeRefreshToken marca um refresh token como revogado.
func (r *AuthRepository) RevokeRefreshToken(tokenHash string) error {
	_, err := r.db.Exec(`
		UPDATE refresh_tokens
		SET revogado = TRUE, updated_at = NOW()
		WHERE token_hash = $1
	`, tokenHash)
	return err
}
