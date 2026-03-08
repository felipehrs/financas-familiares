export interface ContaFixa {
  id: string
  descricao: string
  membro_id: string
  categoria_id: string | null
  valor: number
  dia_vencimento: number
  forma_pagamento: string
  ativa: boolean
}

export interface CriarContaFixaRequest {
  descricao: string
  membro_id: string
  categoria_id?: string
  valor: number
  dia_vencimento: number
  forma_pagamento: string
}

export interface AtualizarContaFixaRequest {
  descricao: string
  membro_id: string
  categoria_id?: string
  valor: number
  dia_vencimento: number
  forma_pagamento: string
  ativa: boolean
}
