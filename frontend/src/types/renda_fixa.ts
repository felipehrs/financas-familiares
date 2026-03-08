export interface RendaFixa {
  id: string
  descricao: string
  membro_id: string
  valor: number
  dia_recebimento: number
  ativa: boolean
}

export interface CriarRendaFixaRequest {
  descricao: string
  membro_id: string
  valor: number
  dia_recebimento: number
}

export interface AtualizarRendaFixaRequest {
  descricao: string
  membro_id: string
  valor: number
  dia_recebimento: number
  ativa: boolean
}
