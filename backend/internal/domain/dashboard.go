package domain

// CategoriaTotalRaw é o resultado bruto da query de despesas por categoria.
type CategoriaTotalRaw struct {
	Nome  string
	Total float64
}
