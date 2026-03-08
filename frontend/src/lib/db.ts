import Dexie, { type EntityTable } from 'dexie'

export interface LocalMembro {
  id: string
  nome: string
  relacionamento: string
  ativo: boolean
  updated_at: string
  deleted_at: string | null
}

export interface LocalCategoria {
  id: string
  nome: string
  updated_at: string
  deleted_at: string | null
}

export interface LocalCartaoCredito {
  id: string
  nome: string
  membro_id: string
  dia_fechamento: number
  dia_vencimento: number
  limite: number | null
  ativo: boolean
  updated_at: string
  deleted_at: string | null
}

export interface LocalDespesaCartao {
  id: string
  compra_id: string
  cartao_id: string
  categoria_id: string | null
  descricao: string
  data_compra: string
  valor_total: number
  numero_parcelas: number
  parcela_numero: number
  valor_parcela: number
  fatura_mes: number
  fatura_ano: number
  fatura: string
  updated_at: string
  deleted_at: string | null
}

export interface LocalAssinatura {
  id: string
  nome: string
  membro_id: string
  categoria_id: string | null
  valor: number
  dia_cobranca: number
  forma_pagamento: string
  status: 'ativa' | 'pausada' | 'cancelada'
  updated_at: string
  deleted_at: string | null
}

export interface LocalContaFixa {
  id: string
  descricao: string
  membro_id: string
  categoria_id: string | null
  valor: number
  dia_vencimento: number
  forma_pagamento: string
  ativa: boolean
  updated_at: string
  deleted_at: string | null
}

export interface LocalDespesaGeral {
  id: string
  membro_id: string
  categoria_id: string | null
  descricao: string
  data: string
  valor: number
  forma_pagamento: string
  observacoes: string | null
  updated_at: string
  deleted_at: string | null
}

export interface LocalRendaFixa {
  id: string
  descricao: string
  membro_id: string
  valor: number
  dia_recebimento: number
  ativa: boolean
  data_inicio: string
  data_fim: string | null
  updated_at: string
  deleted_at: string | null
}

export interface LocalRendaVariavel {
  id: string
  descricao: string
  membro_id: string
  mes_referencia: number
  ano_referencia: number
  valor: number
  data_recebimento: string
  updated_at: string
  deleted_at: string | null
}

export interface LocalRendaExtra {
  id: string
  descricao: string
  membro_id: string
  data_recebimento: string
  valor: number
  updated_at: string
  deleted_at: string | null
}

export interface LocalRendimentoInvestimento {
  id: string
  descricao: string
  membro_id: string
  data: string
  valor: number
  valor_distribuido: number
  updated_at: string
  deleted_at: string | null
}

export interface SyncQueueItem {
  id?: number
  method: 'POST' | 'PUT' | 'PATCH' | 'DELETE'
  endpoint: string
  body: unknown
  token: string
  created_at: string
  retries: number
}

export class FinancasDB extends Dexie {
  membros!: EntityTable<LocalMembro, 'id'>
  categorias!: EntityTable<LocalCategoria, 'id'>
  cartoes_credito!: EntityTable<LocalCartaoCredito, 'id'>
  despesas_cartao!: EntityTable<LocalDespesaCartao, 'id'>
  assinaturas!: EntityTable<LocalAssinatura, 'id'>
  contas_fixas!: EntityTable<LocalContaFixa, 'id'>
  despesas_gerais!: EntityTable<LocalDespesaGeral, 'id'>
  rendas_fixas!: EntityTable<LocalRendaFixa, 'id'>
  rendas_variaveis!: EntityTable<LocalRendaVariavel, 'id'>
  rendas_extras!: EntityTable<LocalRendaExtra, 'id'>
  rendimentos_investimento!: EntityTable<LocalRendimentoInvestimento, 'id'>
  sync_queue!: EntityTable<SyncQueueItem, 'id'>

  constructor() {
    super('financas-familiares')
    this.version(1).stores({
      membros: 'id, ativo, updated_at',
      categorias: 'id, updated_at',
      cartoes_credito: 'id, membro_id, updated_at',
      despesas_cartao: 'id, cartao_id, fatura_mes, fatura_ano, updated_at',
      assinaturas: 'id, membro_id, updated_at',
      contas_fixas: 'id, membro_id, updated_at',
      despesas_gerais: 'id, membro_id, data, updated_at',
      rendas_fixas: 'id, membro_id, updated_at',
      rendas_variaveis: 'id, membro_id, updated_at',
      rendas_extras: 'id, membro_id, updated_at',
      rendimentos_investimento: 'id, membro_id, updated_at',
      sync_queue: '++id, created_at',
    })
  }
}

export const db = new FinancasDB()
