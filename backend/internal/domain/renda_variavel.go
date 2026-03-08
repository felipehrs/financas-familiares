package domain

import (
	"errors"
	"time"
)

// RendaVariavel representa uma renda variável mensal no sistema (freelance, comissões, etc.).
type RendaVariavel struct {
	ID              string
	Descricao       string
	MembroID        string
	MesReferencia   int
	AnoReferencia   int
	Valor           float64
	DataRecebimento time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

var (
	// ErrRendaVariavelNaoEncontrada é retornado quando uma renda variável não é encontrada pelo ID.
	ErrRendaVariavelNaoEncontrada = errors.New("renda variável não encontrada")

	// ErrDescricaoRendaVariavelObrigatoria é retornado quando a descrição não é informada.
	ErrDescricaoRendaVariavelObrigatoria = errors.New("descrição da renda variável é obrigatória")

	// ErrMembroIDRendaVariavelObrigatorio é retornado quando o membro não é informado.
	ErrMembroIDRendaVariavelObrigatorio = errors.New("membro ID é obrigatório para renda variável")

	// ErrMesReferenciaRendaVariavelInvalido é retornado quando o mês de referência não está entre 1 e 12.
	ErrMesReferenciaRendaVariavelInvalido = errors.New("mês de referência da renda variável deve ser entre 1 e 12")

	// ErrAnoReferenciaRendaVariavelInvalido é retornado quando o ano de referência é inválido.
	ErrAnoReferenciaRendaVariavelInvalido = errors.New("ano de referência da renda variável deve ser maior que zero")

	// ErrValorRendaVariavelInvalido é retornado quando o valor é zero ou negativo.
	ErrValorRendaVariavelInvalido = errors.New("valor da renda variável deve ser maior que zero")

	// ErrDataRecebimentoRendaVariavelObrigatoria é retornado quando a data de recebimento não é informada.
	ErrDataRecebimentoRendaVariavelObrigatoria = errors.New("data de recebimento da renda variável é obrigatória")
)
