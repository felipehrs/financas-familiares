export interface RendaVariavel {
  id: string
  descricao: string
  membro_id: string
  mes_referencia: number   // 1-12
  ano_referencia: number
  valor: number
  data_recebimento: string // "YYYY-MM-DD"
}

export interface CriarRendaVariavelRequest {
  descricao: string
  membro_id: string
  mes_referencia: number
  ano_referencia: number
  valor: number
  data_recebimento: string
}

export interface AtualizarRendaVariavelRequest extends CriarRendaVariavelRequest {}
