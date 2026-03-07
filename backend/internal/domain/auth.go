package domain

import "errors"

var (
	ErrCredenciaisInvalidas = errors.New("credenciais inválidas")
	ErrTokenInvalido        = errors.New("token inválido")
	ErrTokenExpirado        = errors.New("token expirado")
	ErrRefreshTokenInvalido = errors.New("refresh token inválido")
)
