package domain

import (
	"errors"
	"time"
)

// DespesaGeral representa uma despesa geral (à vista, débito ou PIX) no sistema.
type DespesaGeral struct {
	ID             string
	MembroID       string
	CategoriaID    *string // opcional
	Descricao      string
	Data           time.Time
	Valor          float64
	FormaPagamento string  // "dinheiro", "debito", "pix"
	Observacoes    *string // opcional
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

var (
	// ErrDespesaGeralNaoEncontrada é retornado quando uma despesa geral não é encontrada pelo ID.
	ErrDespesaGeralNaoEncontrada = errors.New("despesa geral não encontrada")

	// ErrDescricaoDespesaGeralObrigatoria é retornado quando a descrição não é informada.
	ErrDescricaoDespesaGeralObrigatoria = errors.New("descrição da despesa geral é obrigatória")

	// ErrMembroIDDespesaGeralObrigatorio é retornado quando o membro responsável não é informado.
	ErrMembroIDDespesaGeralObrigatorio = errors.New("membro ID é obrigatório para despesa geral")

	// ErrValorDespesaGeralInvalido é retornado quando o valor da despesa é zero ou negativo.
	ErrValorDespesaGeralInvalido = errors.New("valor da despesa geral deve ser maior que zero")

	// ErrDataDespesaGeralObrigatoria é retornado quando a data não é informada.
	ErrDataDespesaGeralObrigatoria = errors.New("data da despesa geral é obrigatória")

	// ErrFormaPagamentoDespesaGeralObrigatoria é retornado quando a forma de pagamento não é informada.
	ErrFormaPagamentoDespesaGeralObrigatoria = errors.New("forma de pagamento da despesa geral é obrigatória")
)
