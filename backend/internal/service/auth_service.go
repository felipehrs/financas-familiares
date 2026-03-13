package service

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// RefreshTokenRecord representa um registro de refresh token armazenado.
type RefreshTokenRecord struct {
	UsuarioID string
	TokenHash string
	ExpiresAt time.Time
	Revogado  bool
}

// AuthRepository define as operações de persistência necessárias para autenticação.
type AuthRepository interface {
	FindByEmail(email string) (*domain.Usuario, error)
	SaveRefreshToken(usuarioID, tokenHash string, expiresAt time.Time) error
	FindRefreshToken(tokenHash string) (*RefreshTokenRecord, error)
	RevokeRefreshToken(tokenHash string) error
}

// FamiliaRepository define as operações de persistência necessárias para famílias.
type FamiliaRepository interface {
	BuscarFamiliaPorUsuario(usuarioID string) (string, error)
}

// AuthServiceInterface define os métodos públicos do serviço de autenticação.
type AuthServiceInterface interface {
	Login(email, senha string) (accessToken string, refreshToken string, err error)
	RefreshToken(refreshToken string) (newAccessToken string, err error)
	ValidateAccessToken(token string) (usuarioID, familiaID string, err error)
}

// AuthService implementa a lógica de autenticação.
type AuthService struct {
	repo        AuthRepository
	familiaRepo FamiliaRepository
	jwtSecret   string
}

// NewAuthService cria uma nova instância do AuthService.
func NewAuthService(repo AuthRepository, familiaRepo FamiliaRepository, jwtSecret string) *AuthService {
	return &AuthService{
		repo:        repo,
		familiaRepo: familiaRepo,
		jwtSecret:   jwtSecret,
	}
}

const (
	accessTokenDuration  = 15 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour
)

// Login autentica o usuário com email e senha.
// Retorna access token (JWT 15min) e refresh token (UUID, 7 dias).
// Em caso de falha (email ou senha inválidos), retorna sempre ErrCredenciaisInvalidas.
func (s *AuthService) Login(email, senha string) (string, string, error) {
	usuario, err := s.repo.FindByEmail(email)
	if err != nil {
		// não revelar se o email não existe ou a senha está errada
		return "", "", domain.ErrCredenciaisInvalidas
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usuario.SenhaHash), []byte(senha)); err != nil {
		return "", "", domain.ErrCredenciaisInvalidas
	}

	familiaID, err := s.familiaRepo.BuscarFamiliaPorUsuario(usuario.ID)
	if err != nil {
		return "", "", fmt.Errorf("erro ao buscar família do usuário: %w", err)
	}

	accessToken, err := criarAccessToken(usuario.ID, familiaID, s.jwtSecret)
	if err != nil {
		return "", "", fmt.Errorf("erro ao gerar access token: %w", err)
	}

	refreshTokenRaw := uuid.New().String()
	tokenHash := hashSHA256(refreshTokenRaw)
	expiresAt := time.Now().Add(refreshTokenDuration)

	if err := s.repo.SaveRefreshToken(usuario.ID, tokenHash, expiresAt); err != nil {
		return "", "", fmt.Errorf("erro ao salvar refresh token: %w", err)
	}

	return accessToken, refreshTokenRaw, nil
}

// RefreshToken valida um refresh token e emite um novo access token.
// O refresh token é revogado após o uso (rotation).
func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	tokenHash := hashSHA256(refreshToken)

	record, err := s.repo.FindRefreshToken(tokenHash)
	if err != nil {
		return "", domain.ErrRefreshTokenInvalido
	}
	if record == nil {
		return "", domain.ErrRefreshTokenInvalido
	}
	if record.Revogado {
		return "", domain.ErrRefreshTokenInvalido
	}
	if time.Now().After(record.ExpiresAt) {
		return "", domain.ErrRefreshTokenInvalido
	}

	// revogar o token usado (refresh token rotation)
	_ = s.repo.RevokeRefreshToken(tokenHash)

	familiaID, err := s.familiaRepo.BuscarFamiliaPorUsuario(record.UsuarioID)
	if err != nil {
		return "", fmt.Errorf("erro ao buscar família do usuário: %w", err)
	}

	newAccessToken, err := criarAccessToken(record.UsuarioID, familiaID, s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("erro ao gerar access token: %w", err)
	}

	return newAccessToken, nil
}

// ValidateAccessToken valida um JWT e retorna o usuarioID e familiaID extraídos dos claims.
func (s *AuthService) ValidateAccessToken(token string) (string, string, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", t.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		if isExpiredError(err) {
			return "", "", domain.ErrTokenExpirado
		}
		return "", "", domain.ErrTokenInvalido
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok || !parsed.Valid {
		return "", "", domain.ErrTokenInvalido
	}

	sub, err := claims.GetSubject()
	if err != nil || sub == "" {
		return "", "", domain.ErrTokenInvalido
	}

	familiaID, _ := claims["familia_id"].(string)

	return sub, familiaID, nil
}

// criarAccessToken gera um JWT com expiração de 15 minutos.
func criarAccessToken(usuarioID, familiaID, secret string) (string, error) {
	return CriarAccessTokenComExpiracao(usuarioID, familiaID, secret, time.Now().Add(accessTokenDuration))
}

// CriarAccessTokenComExpiracao é exportada para uso nos testes (permite simular tokens expirados).
func CriarAccessTokenComExpiracao(usuarioID, familiaID, secret string, expiresAt time.Time) (string, error) {
	claims := jwt.MapClaims{
		"sub":        usuarioID,
		"familia_id": familiaID,
		"exp":        expiresAt.Unix(),
		"iat":        time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// hashSHA256 calcula o hash SHA-256 de uma string e retorna em hexadecimal.
func hashSHA256(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}

// isExpiredError verifica se o erro é de token expirado.
func isExpiredError(err error) bool {
	return err != nil && (err.Error() == "token has expired" ||
		containsString(err.Error(), "expired"))
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
