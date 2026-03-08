package domain

import (
	"errors"
	"time"
)

// RendaExtra representa uma renda extra (pontual) no sistema (bônus, venda, etc.).
type RendaExtra struct {
	ID              string
	Descricao       string
	MembroID        string
	DataRecebimento time.Time
	Valor           float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

var (
	// ErrRendaExtraNaoEncontrada é retornado quando uma renda extra não é encontrada pelo ID.
	ErrRendaExtraNaoEncontrada = errors.New("renda extra não encontrada")

	// ErrDescricaoRendaExtraObrigatoria é retornado quando a descrição não é informada.
	ErrDescricaoRendaExtraObrigatoria = errors.New("descrição da renda extra é obrigatória")

	// ErrMembroIDRendaExtraObrigatorio é retornado quando o membro não é informado.
	ErrMembroIDRendaExtraObrigatorio = errors.New("membro ID é obrigatório para renda extra")

	// ErrValorRendaExtraInvalido é retornado quando o valor é zero ou negativo.
	ErrValorRendaExtraInvalido = errors.New("valor da renda extra deve ser maior que zero")

	// ErrDataRecebimentoRendaExtraObrigatoria é retornado quando a data de recebimento não é informada.
	ErrDataRecebimentoRendaExtraObrigatoria = errors.New("data de recebimento da renda extra é obrigatória")
)
