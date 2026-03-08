export interface RendaFixa {
  id: string
  descricao: string
  membro_id: string
  valor: number
  dia_recebimento: number
  ativa: boolean
  data_inicio: string       // "YYYY-MM-DD"
  data_fim: string | null   // "YYYY-MM-DD" or null
}

export interface CriarRendaFixaRequest {
  descricao: string
  membro_id: string
  valor: number
  dia_recebimento: number
  data_inicio: string
  data_fim?: string
}

export interface AtualizarRendaFixaRequest {
  descricao: string
  membro_id: string
  valor: number
  dia_recebimento: number
  ativa: boolean
  data_inicio: string
  data_fim?: string
}
