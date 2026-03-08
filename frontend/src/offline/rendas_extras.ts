import * as api from '@/api/rendas_extras'
import { db } from '@/lib/db'
import type { RendaExtra, CriarRendaExtraRequest, AtualizarRendaExtraRequest } from '@/types/renda_extra'

const API_BASE = import.meta.env.VITE_API_URL ?? ''

export async function listarRendasExtras(
  token: string,
  mes?: number,
  ano?: number,
): Promise<RendaExtra[]> {
  try {
    const rendas = await api.listarRendasExtras(token, mes, ano)
    const now = new Date().toISOString()
    await db.rendas_extras.bulkPut(
      rendas.map(r => ({ ...r, updated_at: now, deleted_at: null })),
    )
    return rendas
  } catch {
    let collection = db.rendas_extras.filter(r => r.deleted_at === null)
    if (mes !== undefined || ano !== undefined) {
      collection = db.rendas_extras.filter(r => {
        if (r.deleted_at !== null) return false
        const date = new Date(r.data_recebimento)
        if (mes !== undefined && date.getMonth() + 1 !== mes) return false
        if (ano !== undefined && date.getFullYear() !== ano) return false
        return true
      })
    }
    const cached = await collection.toArray()
    return cached.map(({ updated_at: _u, deleted_at: _d, ...r }) => r)
  }
}

export async function criarRendaExtra(
  token: string,
  data: CriarRendaExtraRequest,
): Promise<RendaExtra> {
  if (!navigator.onLine) {
    const id = crypto.randomUUID()
    const now = new Date().toISOString()
    const local: RendaExtra = {
      id,
      descricao: data.descricao,
      membro_id: data.membro_id,
      data_recebimento: data.data_recebimento,
      valor: data.valor,
    }
    await db.rendas_extras.add({ ...local, updated_at: now, deleted_at: null })
    await db.sync_queue.add({
      method: 'POST',
      endpoint: `${API_BASE}/api/v1/rendas-extras`,
      body: data,
      token,
      created_at: now,
      retries: 0,
    })
    return local
  }
  const renda = await api.criarRendaExtra(token, data)
  const now = new Date().toISOString()
  await db.rendas_extras.put({ ...renda, updated_at: now, deleted_at: null })
  return renda
}

export async function atualizarRendaExtra(
  token: string,
  id: string,
  data: AtualizarRendaExtraRequest,
): Promise<RendaExtra> {
  if (!navigator.onLine) {
    const now = new Date().toISOString()
    const local: RendaExtra = {
      id,
      descricao: data.descricao,
      membro_id: data.membro_id,
      data_recebimento: data.data_recebimento,
      valor: data.valor,
    }
    await db.rendas_extras.put({ ...local, updated_at: now, deleted_at: null })
    await db.sync_queue.add({
      method: 'PUT',
      endpoint: `${API_BASE}/api/v1/rendas-extras/${id}`,
      body: data,
      token,
      created_at: now,
      retries: 0,
    })
    return local
  }
  const renda = await api.atualizarRendaExtra(token, id, data)
  const now = new Date().toISOString()
  await db.rendas_extras.put({ ...renda, updated_at: now, deleted_at: null })
  return renda
}

export async function excluirRendaExtra(token: string, id: string): Promise<void> {
  if (!navigator.onLine) {
    const now = new Date().toISOString()
    await db.rendas_extras.update(id, { deleted_at: now, updated_at: now })
    await db.sync_queue.add({
      method: 'DELETE',
      endpoint: `${API_BASE}/api/v1/rendas-extras/${id}`,
      body: null,
      token,
      created_at: now,
      retries: 0,
    })
    return
  }
  await api.excluirRendaExtra(token, id)
  const now = new Date().toISOString()
  await db.rendas_extras.update(id, { deleted_at: now, updated_at: now })
}
