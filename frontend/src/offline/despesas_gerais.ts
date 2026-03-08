import * as api from '@/api/despesas_gerais'
import { db } from '@/lib/db'
import type { DespesaGeral, CriarDespesaGeralRequest, AtualizarDespesaGeralRequest } from '@/types/despesa_geral'

const API_BASE = import.meta.env.VITE_API_URL ?? ''

export async function listarDespesasGerais(
  token: string,
  mes?: number,
  ano?: number,
): Promise<DespesaGeral[]> {
  try {
    const despesas = await api.listarDespesasGerais(token, mes, ano)
    const now = new Date().toISOString()
    await db.despesas_gerais.bulkPut(
      despesas.map(d => ({ ...d, updated_at: now, deleted_at: null })),
    )
    return despesas
  } catch {
    let collection = db.despesas_gerais.filter(d => d.deleted_at === null)
    if (mes !== undefined || ano !== undefined) {
      collection = db.despesas_gerais.filter(d => {
        if (d.deleted_at !== null) return false
        const date = new Date(d.data)
        if (mes !== undefined && date.getMonth() + 1 !== mes) return false
        if (ano !== undefined && date.getFullYear() !== ano) return false
        return true
      })
    }
    const cached = await collection.toArray()
    return cached.map(({ updated_at: _u, deleted_at: _d, ...d }) => d)
  }
}

export async function criarDespesaGeral(
  token: string,
  data: CriarDespesaGeralRequest,
): Promise<DespesaGeral> {
  if (!navigator.onLine) {
    const id = crypto.randomUUID()
    const now = new Date().toISOString()
    const local: DespesaGeral = {
      id,
      membro_id: data.membro_id,
      categoria_id: data.categoria_id ?? null,
      descricao: data.descricao,
      data: data.data,
      valor: data.valor,
      forma_pagamento: data.forma_pagamento,
      observacoes: data.observacoes ?? null,
    }
    await db.despesas_gerais.add({ ...local, updated_at: now, deleted_at: null })
    await db.sync_queue.add({
      method: 'POST',
      endpoint: `${API_BASE}/api/v1/despesas-gerais`,
      body: data,
      token,
      created_at: now,
      retries: 0,
    })
    return local
  }
  const despesa = await api.criarDespesaGeral(token, data)
  const now = new Date().toISOString()
  await db.despesas_gerais.put({ ...despesa, updated_at: now, deleted_at: null })
  return despesa
}

export async function atualizarDespesaGeral(
  token: string,
  id: string,
  data: AtualizarDespesaGeralRequest,
): Promise<DespesaGeral> {
  if (!navigator.onLine) {
    const now = new Date().toISOString()
    const local: DespesaGeral = {
      id,
      membro_id: data.membro_id,
      categoria_id: data.categoria_id ?? null,
      descricao: data.descricao,
      data: data.data,
      valor: data.valor,
      forma_pagamento: data.forma_pagamento,
      observacoes: data.observacoes ?? null,
    }
    await db.despesas_gerais.put({ ...local, updated_at: now, deleted_at: null })
    await db.sync_queue.add({
      method: 'PUT',
      endpoint: `${API_BASE}/api/v1/despesas-gerais/${id}`,
      body: data,
      token,
      created_at: now,
      retries: 0,
    })
    return local
  }
  const despesa = await api.atualizarDespesaGeral(token, id, data)
  const now = new Date().toISOString()
  await db.despesas_gerais.put({ ...despesa, updated_at: now, deleted_at: null })
  return despesa
}

export async function excluirDespesaGeral(token: string, id: string): Promise<void> {
  if (!navigator.onLine) {
    const now = new Date().toISOString()
    await db.despesas_gerais.update(id, { deleted_at: now, updated_at: now })
    await db.sync_queue.add({
      method: 'DELETE',
      endpoint: `${API_BASE}/api/v1/despesas-gerais/${id}`,
      body: null,
      token,
      created_at: now,
      retries: 0,
    })
    return
  }
  await api.excluirDespesaGeral(token, id)
  const now = new Date().toISOString()
  await db.despesas_gerais.update(id, { deleted_at: now, updated_at: now })
}
