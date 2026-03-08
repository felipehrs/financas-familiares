package domain

import "errors"

// CartaoCredito representa um cartão de crédito vinculado a um membro da família.
type CartaoCredito struct {
	ID            string
	Nome          string
	MembroID      string
	DiaFechamento int
	DiaVencimento int
	Limite        *float64 // opcional
	Ativo         bool
}

var (
	// ErrCartaoNaoEncontrado é retornado quando um cartão não é encontrado pelo ID.
	ErrCartaoNaoEncontrado = errors.New("cartão não encontrado")

	// ErrCartaoComVinculos é retornado quando se tenta excluir um cartão com despesas vinculadas.
	ErrCartaoComVinculos = errors.New("cartão possui despesas vinculadas e não pode ser excluído")

	// ErrNomeCartaoObrigatorio é retornado quando o nome do cartão não é informado.
	ErrNomeCartaoObrigatorio = errors.New("nome do cartão é obrigatório")

	// ErrDiaFechamentoInvalido é retornado quando o dia de fechamento está fora do intervalo válido.
	ErrDiaFechamentoInvalido = errors.New("dia de fechamento deve ser entre 1 e 31")

	// ErrDiaVencimentoInvalido é retornado quando o dia de vencimento está fora do intervalo válido.
	ErrDiaVencimentoInvalido = errors.New("dia de vencimento deve ser entre 1 e 31")

	// ErrMembroIDObrigatorio é retornado quando o membro responsável não é informado.
	ErrMembroIDObrigatorio = errors.New("membro responsável é obrigatório")
)
