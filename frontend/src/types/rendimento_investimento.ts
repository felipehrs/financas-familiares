export interface RendimentoInvestimento {
  id: string
  descricao: string
  membro_id: string
  data: string             // "YYYY-MM-DD"
  valor: number
  valor_distribuido: number // >= 0, <= valor, default 0
}

export interface CriarRendimentoInvestimentoRequest {
  descricao: string
  membro_id: string
  data: string
  valor: number
  valor_distribuido?: number  // opcional, default 0
}

export interface AtualizarRendimentoInvestimentoRequest extends CriarRendimentoInvestimentoRequest {
  valor_distribuido: number   // obrigatório no update
}
