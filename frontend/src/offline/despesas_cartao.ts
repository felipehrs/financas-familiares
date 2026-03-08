import * as api from '@/api/despesas_cartao'
import { db } from '@/lib/db'
import type { DespesaCartao, CriarDespesaCartaoRequest } from '@/types/despesa_cartao'

const API_BASE = import.meta.env.VITE_API_URL ?? ''

export async function listarDespesasPorCartao(token: string, cartaoId: string): Promise<DespesaCartao[]> {
  try {
    const despesas = await api.listarDespesasPorCartao(token, cartaoId)
    const now = new Date().toISOString()
    await db.despesas_cartao.bulkPut(
      despesas.map(d => ({ ...d, updated_at: now, deleted_at: null })),
    )
    return despesas
  } catch {
    const cached = await db.despesas_cartao
      .where('cartao_id')
      .equals(cartaoId)
      .filter(d => d.deleted_at === null)
      .toArray()
    return cached.map(({ updated_at: _u, deleted_at: _d, ...d }) => d)
  }
}

export async function listarDespesasPorFatura(
  token: string,
  cartaoId: string,
  mes: number,
  ano: number,
): Promise<DespesaCartao[]> {
  try {
    const despesas = await api.listarDespesasPorFatura(token, cartaoId, mes, ano)
    const now = new Date().toISOString()
    await db.despesas_cartao.bulkPut(
      despesas.map(d => ({ ...d, updated_at: now, deleted_at: null })),
    )
    return despesas
  } catch {
    const cached = await db.despesas_cartao
      .where('cartao_id')
      .equals(cartaoId)
      .filter(d => d.fatura_mes === mes && d.fatura_ano === ano && d.deleted_at === null)
      .toArray()
    return cached.map(({ updated_at: _u, deleted_at: _d, ...d }) => d)
  }
}

export async function criarDespesa(
  token: string,
  cartaoId: string,
  data: CriarDespesaCartaoRequest,
): Promise<DespesaCartao[]> {
  if (!navigator.onLine) {
    const now = new Date().toISOString()
    const compraId = crypto.randomUUID()
    const numParcelas = data.numero_parcelas || 1
    const valorParcela = data.valor_total / numParcelas
    const localItems: DespesaCartao[] = []
    for (let i = 0; i < numParcelas; i++) {
      const id = crypto.randomUUID()
      // Simple fatura calculation (approximate for offline, server will recalculate)
      const dataCompra = new Date(data.data_compra)
      const faturaDate = new Date(dataCompra)
      faturaDate.setMonth(faturaDate.getMonth() + i)
      const faturaMes = faturaDate.getMonth() + 1
      const faturaAno = faturaDate.getFullYear()
      const item: DespesaCartao = {
        id,
        compra_id: compraId,
        cartao_id: cartaoId,
        categoria_id: data.categoria_id ?? null,
        descricao: data.descricao,
        data_compra: data.data_compra,
        valor_total: data.valor_total,
        numero_parcelas: numParcelas,
        parcela_numero: i + 1,
        valor_parcela: valorParcela,
        fatura_mes: faturaMes,
        fatura_ano: faturaAno,
        fatura: `${String(faturaMes).padStart(2, '0')}/${String(faturaAno).slice(-2)}`,
      }
      localItems.push(item)
    }
    await db.despesas_cartao.bulkAdd(
      localItems.map(item => ({ ...item, updated_at: now, deleted_at: null })),
    )
    await db.sync_queue.add({
      method: 'POST',
      endpoint: `${API_BASE}/api/v1/cartoes/${cartaoId}/despesas`,
      body: data,
      token,
      created_at: now,
      retries: 0,
    })
    return localItems
  }
  const despesas = await api.criarDespesa(token, cartaoId, data)
  const now = new Date().toISOString()
  await db.despesas_cartao.bulkPut(
    despesas.map(d => ({ ...d, updated_at: now, deleted_at: null })),
  )
  return despesas
}

export async function excluirDespesa(token: string, id: string): Promise<void> {
  if (!navigator.onLine) {
    const now = new Date().toISOString()
    await db.despesas_cartao.update(id, { deleted_at: now, updated_at: now })
    await db.sync_queue.add({
      method: 'DELETE',
      endpoint: `${API_BASE}/api/v1/despesas/${id}`,
      body: null,
      token,
      created_at: now,
      retries: 0,
    })
    return
  }
  await api.excluirDespesa(token, id)
  const now = new Date().toISOString()
  await db.despesas_cartao.update(id, { deleted_at: now, updated_at: now })
}
