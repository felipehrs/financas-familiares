export interface DespesaGeral {
  id: string
  membro_id: string
  categoria_id: string | null
  descricao: string
  data: string         // "YYYY-MM-DD"
  valor: number
  forma_pagamento: string
  observacoes: string | null
}

export interface CriarDespesaGeralRequest {
  membro_id: string
  categoria_id?: string
  descricao: string
  data: string
  valor: number
  forma_pagamento: string
  observacoes?: string
}

export type AtualizarDespesaGeralRequest = CriarDespesaGeralRequest
