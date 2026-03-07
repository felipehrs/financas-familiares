package domain

import "errors"

// Membro representa um membro da família no sistema.
type Membro struct {
	ID             string
	Nome           string
	Relacionamento string
	Ativo          bool
}

var (
	// ErrMembroNaoEncontrado é retornado quando um membro não é encontrado pelo ID.
	ErrMembroNaoEncontrado = errors.New("membro não encontrado")

	// ErrMembroComVinculos é retornado quando se tenta excluir um membro com vínculos ativos.
	ErrMembroComVinculos = errors.New("membro possui vínculos ativos e não pode ser excluído")

	// ErrNomeObrigatorio é retornado quando o nome do membro não é informado.
	ErrNomeObrigatorio = errors.New("nome é obrigatório")
)
