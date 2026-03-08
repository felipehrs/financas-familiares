package domain

import (
	"errors"
	"time"
)

// RendimentoInvestimento representa um rendimento de investimento no sistema.
type RendimentoInvestimento struct {
	ID               string
	Descricao        string
	MembroID         string
	Data             time.Time
	Valor            float64
	ValorDistribuido float64 // default 0, 0 <= ValorDistribuido <= Valor
	CreatedAt        time.Time
}

var (
	// ErrRendimentoInvestimentoNaoEncontrado é retornado quando um rendimento não é encontrado pelo ID.
	ErrRendimentoInvestimentoNaoEncontrado = errors.New("rendimento de investimento não encontrado")

	// ErrDescricaoRendimentoObrigatoria é retornado quando a descrição não é informada.
	ErrDescricaoRendimentoObrigatoria = errors.New("descrição do rendimento é obrigatória")

	// ErrMembroIDRendimentoObrigatorio é retornado quando o membro não é informado.
	ErrMembroIDRendimentoObrigatorio = errors.New("membro ID é obrigatório para rendimento de investimento")

	// ErrValorRendimentoInvalido é retornado quando o valor é zero ou negativo.
	ErrValorRendimentoInvalido = errors.New("valor do rendimento deve ser maior que zero")

	// ErrDataRendimentoObrigatoria é retornado quando a data não é informada.
	ErrDataRendimentoObrigatoria = errors.New("data do rendimento é obrigatória")

	// ErrValorDistribuidoInvalido é retornado quando valor_distribuido < 0 ou > valor.
	ErrValorDistribuidoInvalido = errors.New("valor distribuído deve ser entre 0 e o valor total do rendimento")
)
