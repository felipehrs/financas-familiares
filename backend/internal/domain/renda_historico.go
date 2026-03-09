package domain

// TipoRenda identifica a origem de uma renda no histórico.
type TipoRenda string

const (
	TipoRendaFixa         TipoRenda = "fixa"
	TipoRendaVariavel     TipoRenda = "variavel"
	TipoRendaExtra        TipoRenda = "extra"
	TipoRendaInvestimento TipoRenda = "investimento"
)

// ItemRendaHistorico representa um registro de renda no histórico filtrado.
type ItemRendaHistorico struct {
	ID               string    `json:"id"`
	Tipo             TipoRenda `json:"tipo"`
	Descricao        string    `json:"descricao"`
	MembroID         string    `json:"membro_id"`
	Valor            float64   `json:"valor"`
	Mes              *int      `json:"mes,omitempty"`
	Ano              *int      `json:"ano,omitempty"`
	Data             *string   `json:"data,omitempty"`
	Ativa            *bool     `json:"ativa,omitempty"`
	ValorDistribuido *float64  `json:"valor_distribuido,omitempty"`
}

// ResumoRendaHistorico agrega os totais por tipo de renda para o filtro aplicado.
type ResumoRendaHistorico struct {
	TotalFixas                     float64 `json:"total_fixas"`
	TotalVariaveis                 float64 `json:"total_variaveis"`
	TotalExtras                    float64 `json:"total_extras"`
	TotalInvestimentosDistribuidos float64 `json:"total_investimentos_distribuidos"`
	TotalGeral                     float64 `json:"total_geral"`
}

// HistoricoRendas contém a lista de itens e o resumo do histórico de rendas.
type HistoricoRendas struct {
	Itens  []ItemRendaHistorico `json:"itens"`
	Resumo ResumoRendaHistorico `json:"resumo"`
}

// FiltroHistoricoRendas encapsula os parâmetros de filtragem do histórico de rendas.
type FiltroHistoricoRendas struct {
	Tipo     TipoRenda
	MembroID string
	Mes      int // 0 = sem filtro
	Ano      int // 0 = sem filtro
}
