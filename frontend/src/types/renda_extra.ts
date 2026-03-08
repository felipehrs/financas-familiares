export interface RendaExtra {
  id: string
  descricao: string
  membro_id: string
  data_recebimento: string  // "YYYY-MM-DD"
  valor: number
}

export interface CriarRendaExtraRequest {
  descricao: string
  membro_id: string
  data_recebimento: string
  valor: number
}

export interface AtualizarRendaExtraRequest extends CriarRendaExtraRequest {}
