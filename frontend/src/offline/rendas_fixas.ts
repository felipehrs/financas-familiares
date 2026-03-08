import * as api from '@/api/rendas_fixas'
import { db } from '@/lib/db'
import type { RendaFixa, CriarRendaFixaRequest, AtualizarRendaFixaRequest } from '@/types/renda_fixa'

const API_BASE = import.meta.env.VITE_API_URL ?? ''

export async function listarRendasFixas(token: string): Promise<RendaFixa[]> {
  try {
    const rendas = await api.listarRendasFixas(token)
    const now = new Date().toISOString()
    await db.rendas_fixas.bulkPut(
      rendas.map(r => ({ ...r, updated_at: now, deleted_at: null })),
    )
    return rendas
  } catch {
    const cached = await db.rendas_fixas.filter(r => r.deleted_at === null).toArray()
    return cached.map(({ updated_at: _u, deleted_at: _d, ...r }) => r)
  }
}

export async function criarRendaFixa(token: string, data: CriarRendaFixaRequest): Promise<RendaFixa> {
  if (!navigator.onLine) {
    const id = crypto.randomUUID()
    const now = new Date().toISOString()
    const local: RendaFixa = {
      id,
      descricao: data.descricao,
      membro_id: data.membro_id,
      valor: data.valor,
      dia_recebimento: data.dia_recebimento,
      ativa: true,
      data_inicio: data.data_inicio,
      data_fim: data.data_fim ?? null,
    }
    await db.rendas_fixas.add({ ...local, updated_at: now, deleted_at: null })
    await db.sync_queue.add({
      method: 'POST',
      endpoint: `${API_BASE}/api/v1/rendas-fixas`,
      body: data,
      token,
      created_at: now,
      retries: 0,
    })
    return local
  }
  const renda = await api.criarRendaFixa(token, data)
  const now = new Date().toISOString()
  await db.rendas_fixas.put({ ...renda, updated_at: now, deleted_at: null })
  return renda
}

export async function atualizarRendaFixa(
  token: string,
  id: string,
  data: AtualizarRendaFixaRequest,
): Promise<RendaFixa> {
  if (!navigator.onLine) {
    const now = new Date().toISOString()
    const local: RendaFixa = {
      id,
      descricao: data.descricao,
      membro_id: data.membro_id,
      valor: data.valor,
      dia_recebimento: data.dia_recebimento,
      ativa: data.ativa,
      data_inicio: data.data_inicio,
      data_fim: data.data_fim ?? null,
    }
    await db.rendas_fixas.put({ ...local, updated_at: now, deleted_at: null })
    await db.sync_queue.add({
      method: 'PUT',
      endpoint: `${API_BASE}/api/v1/rendas-fixas/${id}`,
      body: data,
      token,
      created_at: now,
      retries: 0,
    })
    return local
  }
  const renda = await api.atualizarRendaFixa(token, id, data)
  const now = new Date().toISOString()
  await db.rendas_fixas.put({ ...renda, updated_at: now, deleted_at: null })
  return renda
}

export async function inativarRendaFixa(token: string, id: string): Promise<void> {
  if (!navigator.onLine) {
    const now = new Date().toISOString()
    await db.rendas_fixas.update(id, { ativa: false, updated_at: now })
    await db.sync_queue.add({
      method: 'PATCH',
      endpoint: `${API_BASE}/api/v1/rendas-fixas/${id}/inativar`,
      body: null,
      token,
      created_at: now,
      retries: 0,
    })
    return
  }
  await api.inativarRendaFixa(token, id)
  const now = new Date().toISOString()
  await db.rendas_fixas.update(id, { ativa: false, updated_at: now })
}

export async function listarRendasFixasVigentes(
  token: string,
  mes: number,
  ano: number,
): Promise<RendaFixa[]> {
  try {
    const rendas = await api.listarRendasFixasVigentes(token, mes, ano)
    const now = new Date().toISOString()
    await db.rendas_fixas.bulkPut(
      rendas.map(r => ({ ...r, updated_at: now, deleted_at: null })),
    )
    return rendas
  } catch {
    const refDate = new Date(ano, mes - 1, 1)
    const cached = await db.rendas_fixas
      .filter(r => {
        if (r.deleted_at !== null || !r.ativa) return false
        const inicio = new Date(r.data_inicio)
        if (inicio > refDate) return false
        if (r.data_fim) {
          const fim = new Date(r.data_fim)
          if (fim < refDate) return false
        }
        return true
      })
      .toArray()
    return cached.map(({ updated_at: _u, deleted_at: _d, ...r }) => r)
  }
}
