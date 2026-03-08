export interface CartaoCredito {
  id: string
  nome: string
  membro_id: string
  dia_fechamento: number
  dia_vencimento: number
  limite: number | null
  ativo: boolean
}

export interface CriarCartaoCreditoRequest {
  nome: string
  membro_id: string
  dia_fechamento: number
  dia_vencimento: number
  limite?: number
}

export interface AtualizarCartaoCreditoRequest {
  nome: string
  membro_id: string
  dia_fechamento: number
  dia_vencimento: number
  limite?: number
  ativo: boolean
}
