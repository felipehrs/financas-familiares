export type StatusAssinatura = 'ativa' | 'pausada' | 'cancelada'

export interface Assinatura {
  id: string
  nome: string
  membro_id: string
  categoria_id: string | null
  valor: number
  dia_cobranca: number
  forma_pagamento: string
  status: StatusAssinatura
}

export interface CriarAssinaturaRequest {
  nome: string
  membro_id: string
  categoria_id?: string
  valor: number
  dia_cobranca: number
  forma_pagamento: string
}

export interface AtualizarAssinaturaRequest {
  nome: string
  membro_id: string
  categoria_id?: string
  valor: number
  dia_cobranca: number
  forma_pagamento: string
  status: StatusAssinatura
}
