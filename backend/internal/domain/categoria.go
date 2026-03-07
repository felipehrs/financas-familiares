package domain

import "errors"

// Categoria representa uma categoria de despesa ou receita no sistema.
type Categoria struct {
	ID   string
	Nome string
}

var (
	// ErrCategoriaNaoEncontrada é retornado quando uma categoria não é encontrada pelo ID.
	ErrCategoriaNaoEncontrada = errors.New("categoria não encontrada")

	// ErrCategoriaComVinculos é retornado quando se tenta excluir uma categoria com registros vinculados.
	ErrCategoriaComVinculos = errors.New("categoria possui registros vinculados e não pode ser excluída")

	// ErrNomeCategoriaObrigatorio é retornado quando o nome da categoria não é informado.
	ErrNomeCategoriaObrigatorio = errors.New("nome da categoria é obrigatório")
)
