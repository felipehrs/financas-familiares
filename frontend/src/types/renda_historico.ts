export type TipoRenda = 'fixa' | 'variavel' | 'extra' | 'investimento'

export interface ItemRendaHistorico {
  id: string
  tipo: TipoRenda
  descricao: string
  membro_id: string
  valor: number
  mes?: number
  ano?: number
  data?: string         // "YYYY-MM-DD"
  ativa?: boolean
  valor_distribuido?: number
}

export interface ResumoRendaHistorico {
  total_fixas: number
  total_variaveis: number
  total_extras: number
  total_investimentos_distribuidos: number
  total_geral: number
}

export interface HistoricoRendas {
  itens: ItemRendaHistorico[]
  resumo: ResumoRendaHistorico
}

export interface FiltroHistoricoRendas {
  tipo?: TipoRenda | ''
  membro_id?: string
  mes?: number
  ano?: number
}
