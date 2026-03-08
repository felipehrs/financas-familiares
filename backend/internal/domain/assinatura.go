package domain

import (
	"errors"
	"time"
)

// Assinatura representa uma assinatura recorrente no sistema.
type Assinatura struct {
	ID             string
	Nome           string
	MembroID       string
	CategoriaID    *string // optional
	Valor          float64
	DiaCobranca    int // 1-31
	FormaPagamento string
	Status         string // "ativa", "pausada", "cancelada"
	CreatedAt      time.Time
}

var (
	// ErrAssinaturaNaoEncontrada é retornado quando uma assinatura não é encontrada pelo ID.
	ErrAssinaturaNaoEncontrada = errors.New("assinatura não encontrada")

	// ErrNomeAssinaturaObrigatorio é retornado quando o nome da assinatura não é informado.
	ErrNomeAssinaturaObrigatorio = errors.New("nome da assinatura é obrigatório")

	// ErrMembroIDAssinaturaObrigatorio é retornado quando o membro responsável não é informado.
	ErrMembroIDAssinaturaObrigatorio = errors.New("membro responsável é obrigatório")

	// ErrValorAssinaturaInvalido é retornado quando o valor da assinatura é inválido.
	ErrValorAssinaturaInvalido = errors.New("valor deve ser maior que zero")

	// ErrDiaCobrancaInvalido é retornado quando o dia de cobrança está fora do intervalo 1-31.
	ErrDiaCobrancaInvalido = errors.New("dia de cobrança deve ser entre 1 e 31")

	// ErrFormaPagamentoObrigatoria é retornado quando a forma de pagamento não é informada.
	ErrFormaPagamentoObrigatoria = errors.New("forma de pagamento é obrigatória")

	// ErrStatusInvalido é retornado quando o status informado não é válido.
	ErrStatusInvalido = errors.New("status deve ser: ativa, pausada ou cancelada")
)

// StatusAssinaturaValidos contém os valores permitidos para o status de uma assinatura.
var StatusAssinaturaValidos = map[string]bool{
	"ativa":     true,
	"pausada":   true,
	"cancelada": true,
}
