package domain

import (
	"errors"
	"time"
)

// ContaFixa representa uma conta fixa mensal no sistema.
type ContaFixa struct {
	ID             string
	Descricao      string
	MembroID       string
	CategoriaID    *string // optional
	Valor          float64
	DiaVencimento  int // 1-28
	FormaPagamento string
	Ativa          bool
	CreatedAt      time.Time
}

var (
	// ErrContaFixaNaoEncontrada é retornado quando uma conta fixa não é encontrada pelo ID.
	ErrContaFixaNaoEncontrada = errors.New("conta fixa não encontrada")

	// ErrDescricaoContaFixaObrigatoria é retornado quando a descrição da conta fixa não é informada.
	ErrDescricaoContaFixaObrigatoria = errors.New("descrição da conta fixa é obrigatória")

	// ErrMembroIDContaFixaObrigatorio é retornado quando o membro responsável não é informado.
	ErrMembroIDContaFixaObrigatorio = errors.New("membro ID é obrigatório para conta fixa")

	// ErrValorContaFixaInvalido é retornado quando o valor da conta fixa é inválido.
	ErrValorContaFixaInvalido = errors.New("valor da conta fixa deve ser maior que zero")

	// ErrDiaVencimentoContaFixaInvalido é retornado quando o dia de vencimento está fora do intervalo 1-28.
	ErrDiaVencimentoContaFixaInvalido = errors.New("dia de vencimento deve ser entre 1 e 28")

	// ErrFormaPagamentoContaFixaObrigatoria é retornado quando a forma de pagamento não é informada.
	ErrFormaPagamentoContaFixaObrigatoria = errors.New("forma de pagamento da conta fixa é obrigatória")
)
