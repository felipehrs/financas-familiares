package domain

import "errors"

// RendaFixa representa uma renda fixa mensal de um membro da família.
type RendaFixa struct {
	ID             string
	Descricao      string
	MembroID       string
	Valor          float64
	DiaRecebimento int
	Ativa          bool
}

var (
	// ErrRendaFixaNaoEncontrada é retornado quando uma renda fixa não é encontrada pelo ID.
	ErrRendaFixaNaoEncontrada = errors.New("renda fixa não encontrada")

	// ErrDescricaoRendaObrigatoria é retornado quando a descrição da renda não é informada.
	ErrDescricaoRendaObrigatoria = errors.New("descrição da renda é obrigatória")

	// ErrMembroIDRendaObrigatorio é retornado quando o membro responsável não é informado.
	ErrMembroIDRendaObrigatorio = errors.New("membro responsável é obrigatório")

	// ErrValorRendaInvalido é retornado quando o valor da renda não é maior que zero.
	ErrValorRendaInvalido = errors.New("valor da renda deve ser maior que zero")

	// ErrDiaRecebimentoInvalido é retornado quando o dia de recebimento está fora do intervalo válido.
	ErrDiaRecebimentoInvalido = errors.New("dia de recebimento deve ser entre 1 e 31")
)
