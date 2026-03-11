package domain

import (
	"errors"
	"time"
)

type Familia struct {
	ID        string
	Nome      string
	OwnerID   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrFamiliaNaoEncontrada = errors.New("família não encontrada")
)
