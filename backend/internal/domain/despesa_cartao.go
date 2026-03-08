package domain

import (
	"errors"
	"time"
)

// DespesaCartao representa uma despesa lançada em um cartão de crédito.
type DespesaCartao struct {
	ID             string
	CartaoID       string
	CategoriaID    *string // opcional
	Descricao      string
	DataCompra     time.Time
	ValorTotal     float64
	NumeroParcelas int
	ValorParcela   float64 // calculado
	FaturaMes      int     // calculado (1-12)
	FaturaAno      int     // calculado
}

var (
	// ErrDespesaCartaoNaoEncontrada é retornado quando uma despesa não é encontrada pelo ID.
	ErrDespesaCartaoNaoEncontrada = errors.New("despesa de cartão não encontrada")

	// ErrCartaoIDObrigatorio é retornado quando o cartão não é informado.
	ErrCartaoIDObrigatorio = errors.New("cartão é obrigatório")

	// ErrDescricaoObrigatoria é retornado quando a descrição não é informada.
	ErrDescricaoObrigatoria = errors.New("descrição é obrigatória")

	// ErrValorTotalInvalido é retornado quando o valor total é menor ou igual a zero.
	ErrValorTotalInvalido = errors.New("valor total deve ser maior que zero")

	// ErrNumeroParcelas é retornado quando o número de parcelas é menor que 1.
	ErrNumeroParcelas = errors.New("número de parcelas deve ser maior ou igual a 1")
)
