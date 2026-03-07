package handler

import (
	"errors"
	"net/http"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// AuthServiceInterface define os métodos do service usados pelo handler.
// Redeclarada aqui para desacoplar o pacote handler do service sem importação circular.
type AuthServiceInterface interface {
	Login(email, senha string) (string, string, error)
	RefreshToken(refreshToken string) (string, error)
	ValidateAccessToken(token string) (string, error)
}

// AuthHandler contém os handlers HTTP para autenticação.
type AuthHandler struct {
	svc AuthServiceInterface
}

// NewAuthHandler cria um novo AuthHandler.
func NewAuthHandler(svc AuthServiceInterface) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type loginRequest struct {
	Email string `json:"email" binding:"required,email"`
	Senha string `json:"senha" binding:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Login autentica o usuário e retorna os tokens.
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	accessToken, refreshToken, err := h.svc.Login(req.Email, req.Senha)
	if err != nil {
		if errors.Is(err, domain.ErrCredenciaisInvalidas) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciais inválidas"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    900, // 15 minutos em segundos
	})
}

// Refresh emite um novo access token a partir de um refresh token válido.
// POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	newAccessToken, err := h.svc.RefreshToken(req.RefreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrRefreshTokenInvalido) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token inválido"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
		"token_type":   "Bearer",
		"expires_in":   900,
	})
}
