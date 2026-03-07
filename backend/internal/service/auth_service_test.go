package service_test

import (
	"crypto/sha256"
	"fmt"
	"testing"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/felipehrs/financas-familiares/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// MockAuthRepository implementa AuthRepository para testes unitários.
type MockAuthRepository struct {
	usuarios      map[string]*domain.Usuario
	refreshTokens map[string]*service.RefreshTokenRecord
}

func newMockRepo() *MockAuthRepository {
	return &MockAuthRepository{
		usuarios:      make(map[string]*domain.Usuario),
		refreshTokens: make(map[string]*service.RefreshTokenRecord),
	}
}

func (m *MockAuthRepository) FindByEmail(email string) (*domain.Usuario, error) {
	u, ok := m.usuarios[email]
	if !ok {
		return nil, domain.ErrCredenciaisInvalidas
	}
	return u, nil
}

func (m *MockAuthRepository) SaveRefreshToken(usuarioID, tokenHash string, expiresAt time.Time) error {
	m.refreshTokens[tokenHash] = &service.RefreshTokenRecord{
		UsuarioID: usuarioID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		Revogado:  false,
	}
	return nil
}

func (m *MockAuthRepository) FindRefreshToken(tokenHash string) (*service.RefreshTokenRecord, error) {
	r, ok := m.refreshTokens[tokenHash]
	if !ok {
		return nil, nil
	}
	return r, nil
}

func (m *MockAuthRepository) RevokeRefreshToken(tokenHash string) error {
	r, ok := m.refreshTokens[tokenHash]
	if !ok {
		return fmt.Errorf("token não encontrado")
	}
	r.Revogado = true
	return nil
}

// helpers

func hashSHA256(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}

func criarUsuarioComSenha(t *testing.T, repo *MockAuthRepository, email, senha string) *domain.Usuario {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(senha), 12)
	require.NoError(t, err)
	u := &domain.Usuario{
		ID:        "usuario-id-1",
		Nome:      "Usuário Teste",
		Email:     email,
		SenhaHash: string(hash),
	}
	repo.usuarios[email] = u
	return u
}

// ---- Testes ----

func TestLogin_Sucesso(t *testing.T) {
	repo := newMockRepo()
	criarUsuarioComSenha(t, repo, "user@example.com", "senha123")
	svc := service.NewAuthService(repo, "segredo-jwt-teste")

	accessToken, refreshToken, err := svc.Login("user@example.com", "senha123")

	require.NoError(t, err)
	assert.NotEmpty(t, accessToken, "access token deve ser retornado")
	assert.NotEmpty(t, refreshToken, "refresh token deve ser retornado")
}

func TestLogin_EmailInvalido(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewAuthService(repo, "segredo-jwt-teste")

	_, _, err := svc.Login("naoexiste@example.com", "qualquersenha")

	assert.ErrorIs(t, err, domain.ErrCredenciaisInvalidas)
}

func TestLogin_SenhaInvalida(t *testing.T) {
	repo := newMockRepo()
	criarUsuarioComSenha(t, repo, "user@example.com", "senha-correta")
	svc := service.NewAuthService(repo, "segredo-jwt-teste")

	_, _, err := svc.Login("user@example.com", "senha-errada")

	// deve retornar o mesmo erro sem revelar qual campo é inválido
	assert.ErrorIs(t, err, domain.ErrCredenciaisInvalidas)
}

func TestRefreshToken_Sucesso(t *testing.T) {
	repo := newMockRepo()
	criarUsuarioComSenha(t, repo, "user@example.com", "senha123")
	svc := service.NewAuthService(repo, "segredo-jwt-teste")

	_, refreshToken, err := svc.Login("user@example.com", "senha123")
	require.NoError(t, err)

	newAccessToken, err := svc.RefreshToken(refreshToken)

	require.NoError(t, err)
	assert.NotEmpty(t, newAccessToken, "novo access token deve ser retornado")
}

func TestRefreshToken_Invalido(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewAuthService(repo, "segredo-jwt-teste")

	_, err := svc.RefreshToken("token-que-nao-existe")

	assert.ErrorIs(t, err, domain.ErrRefreshTokenInvalido)
}

func TestRefreshToken_Expirado(t *testing.T) {
	repo := newMockRepo()
	criarUsuarioComSenha(t, repo, "user@example.com", "senha123")
	svc := service.NewAuthService(repo, "segredo-jwt-teste")

	// inserir manualmente um refresh token já expirado no mock
	tokenRaw := "refresh-token-expirado"
	tokenHash := hashSHA256(tokenRaw)
	repo.refreshTokens[tokenHash] = &service.RefreshTokenRecord{
		UsuarioID: "usuario-id-1",
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(-24 * time.Hour), // expirado ontem
		Revogado:  false,
	}

	_, err := svc.RefreshToken(tokenRaw)

	assert.ErrorIs(t, err, domain.ErrRefreshTokenInvalido)
}

func TestRefreshToken_Revogado(t *testing.T) {
	repo := newMockRepo()
	criarUsuarioComSenha(t, repo, "user@example.com", "senha123")
	svc := service.NewAuthService(repo, "segredo-jwt-teste")

	_, refreshToken, err := svc.Login("user@example.com", "senha123")
	require.NoError(t, err)

	// revogar manualmente
	tokenHash := hashSHA256(refreshToken)
	err = repo.RevokeRefreshToken(tokenHash)
	require.NoError(t, err)

	_, err = svc.RefreshToken(refreshToken)

	assert.ErrorIs(t, err, domain.ErrRefreshTokenInvalido)
}

func TestValidateAccessToken_Sucesso(t *testing.T) {
	repo := newMockRepo()
	criarUsuarioComSenha(t, repo, "user@example.com", "senha123")
	svc := service.NewAuthService(repo, "segredo-jwt-teste")

	accessToken, _, err := svc.Login("user@example.com", "senha123")
	require.NoError(t, err)

	usuarioID, err := svc.ValidateAccessToken(accessToken)

	require.NoError(t, err)
	assert.Equal(t, "usuario-id-1", usuarioID)
}

func TestValidateAccessToken_Expirado(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewAuthService(repo, "segredo-jwt-teste")

	// criar token já expirado
	expiredToken, err := service.CriarAccessTokenComExpiracao("usuario-id-1", "segredo-jwt-teste", time.Now().Add(-1*time.Minute))
	require.NoError(t, err)

	_, err = svc.ValidateAccessToken(expiredToken)

	assert.ErrorIs(t, err, domain.ErrTokenExpirado)
}

func TestValidateAccessToken_Invalido(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewAuthService(repo, "segredo-jwt-teste")

	_, err := svc.ValidateAccessToken("token.invalido.qualquer")

	assert.ErrorIs(t, err, domain.ErrTokenInvalido)
}
