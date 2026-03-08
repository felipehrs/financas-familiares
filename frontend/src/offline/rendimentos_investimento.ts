import * as api from '@/api/rendimentos_investimento'
import { db } from '@/lib/db'
import type { RendimentoInvestimento, CriarRendimentoInvestimentoRequest, AtualizarRendimentoInvestimentoRequest } from '@/types/rendimento_investimento'

const API_BASE = import.meta.env.VITE_API_URL ?? ''

export async function listarRendimentos(
  token: string,
  mes?: number,
  ano?: number,
): Promise<RendimentoInvestimento[]> {
  try {
    const rendimentos = await api.listarRendimentos(token, mes, ano)
    const now = new Date().toISOString()
    await db.rendimentos_investimento.bulkPut(
      rendimentos.map(r => ({ ...r, updated_at: now, deleted_at: null })),
    )
    return rendimentos
  } catch {
    let collection = db.rendimentos_investimento.filter(r => r.deleted_at === null)
    if (mes !== undefined || ano !== undefined) {
      collection = db.rendimentos_investimento.filter(r => {
        if (r.deleted_at !== null) return false
        const date = new Date(r.data)
        if (mes !== undefined && date.getMonth() + 1 !== mes) return false
        if (ano !== undefined && date.getFullYear() !== ano) return false
        return true
      })
    }
    const cached = await collection.toArray()
    return cached.map(({ updated_at: _u, deleted_at: _d, ...r }) => r)
  }
}

export async function criarRendimento(
  token: string,
  data: CriarRendimentoInvestimentoRequest,
): Promise<RendimentoInvestimento> {
  if (!navigator.onLine) {
    const id = crypto.randomUUID()
    const now = new Date().toISOString()
    const local: RendimentoInvestimento = {
      id,
      descricao: data.descricao,
      membro_id: data.membro_id,
      data: data.data,
      valor: data.valor,
      valor_distribuido: data.valor_distribuido ?? 0,
    }
    await db.rendimentos_investimento.add({ ...local, updated_at: now, deleted_at: null })
    await db.sync_queue.add({
      method: 'POST',
      endpoint: `${API_BASE}/api/v1/rendimentos-investimento`,
      body: data,
      token,
      created_at: now,
      retries: 0,
    })
    return local
  }
  const rendimento = await api.criarRendimento(token, data)
  const now = new Date().toISOString()
  await db.rendimentos_investimento.put({ ...rendimento, updated_at: now, deleted_at: null })
  return rendimento
}

export async function atualizarRendimento(
  token: string,
  id: string,
  data: AtualizarRendimentoInvestimentoRequest,
): Promise<RendimentoInvestimento> {
  if (!navigator.onLine) {
    const now = new Date().toISOString()
    const local: RendimentoInvestimento = {
      id,
      descricao: data.descricao,
      membro_id: data.membro_id,
      data: data.data,
      valor: data.valor,
      valor_distribuido: data.valor_distribuido,
    }
    await db.rendimentos_investimento.put({ ...local, updated_at: now, deleted_at: null })
    await db.sync_queue.add({
      method: 'PUT',
      endpoint: `${API_BASE}/api/v1/rendimentos-investimento/${id}`,
      body: data,
      token,
      created_at: now,
      retries: 0,
    })
    return local
  }
  const rendimento = await api.atualizarRendimento(token, id, data)
  const now = new Date().toISOString()
  await db.rendimentos_investimento.put({ ...rendimento, updated_at: now, deleted_at: null })
  return rendimento
}

export async function excluirRendimento(token: string, id: string): Promise<void> {
  if (!navigator.onLine) {
    const now = new Date().toISOString()
    await db.rendimentos_investimento.update(id, { deleted_at: now, updated_at: now })
    await db.sync_queue.add({
      method: 'DELETE',
      endpoint: `${API_BASE}/api/v1/rendimentos-investimento/${id}`,
      body: null,
      token,
      created_at: now,
      retries: 0,
    })
    return
  }
  await api.excluirRendimento(token, id)
  const now = new Date().toISOString()
  await db.rendimentos_investimento.update(id, { deleted_at: now, updated_at: now })
}
