import * as api from '@/api/rendas_variaveis'
import { db } from '@/lib/db'
import type { RendaVariavel, CriarRendaVariavelRequest, AtualizarRendaVariavelRequest } from '@/types/renda_variavel'

const API_BASE = import.meta.env.VITE_API_URL ?? ''

export async function listarRendasVariaveis(
  token: string,
  mes?: number,
  ano?: number,
): Promise<RendaVariavel[]> {
  try {
    const rendas = await api.listarRendasVariaveis(token, mes, ano)
    const now = new Date().toISOString()
    await db.rendas_variaveis.bulkPut(
      rendas.map(r => ({ ...r, updated_at: now, deleted_at: null })),
    )
    return rendas
  } catch {
    let collection = db.rendas_variaveis.filter(r => r.deleted_at === null)
    if (mes !== undefined || ano !== undefined) {
      collection = db.rendas_variaveis.filter(r => {
        if (r.deleted_at !== null) return false
        if (mes !== undefined && r.mes_referencia !== mes) return false
        if (ano !== undefined && r.ano_referencia !== ano) return false
        return true
      })
    }
    const cached = await collection.toArray()
    return cached.map(({ updated_at: _u, deleted_at: _d, ...r }) => r)
  }
}

export async function criarRendaVariavel(
  token: string,
  data: CriarRendaVariavelRequest,
): Promise<RendaVariavel> {
  if (!navigator.onLine) {
    const id = crypto.randomUUID()
    const now = new Date().toISOString()
    const local: RendaVariavel = {
      id,
      descricao: data.descricao,
      membro_id: data.membro_id,
      mes_referencia: data.mes_referencia,
      ano_referencia: data.ano_referencia,
      valor: data.valor,
      data_recebimento: data.data_recebimento,
    }
    await db.rendas_variaveis.add({ ...local, updated_at: now, deleted_at: null })
    await db.sync_queue.add({
      method: 'POST',
      endpoint: `${API_BASE}/api/v1/rendas-variaveis`,
      body: data,
      token,
      created_at: now,
      retries: 0,
    })
    return local
  }
  const renda = await api.criarRendaVariavel(token, data)
  const now = new Date().toISOString()
  await db.rendas_variaveis.put({ ...renda, updated_at: now, deleted_at: null })
  return renda
}

export async function atualizarRendaVariavel(
  token: string,
  id: string,
  data: AtualizarRendaVariavelRequest,
): Promise<RendaVariavel> {
  if (!navigator.onLine) {
    const now = new Date().toISOString()
    const local: RendaVariavel = {
      id,
      descricao: data.descricao,
      membro_id: data.membro_id,
      mes_referencia: data.mes_referencia,
      ano_referencia: data.ano_referencia,
      valor: data.valor,
      data_recebimento: data.data_recebimento,
    }
    await db.rendas_variaveis.put({ ...local, updated_at: now, deleted_at: null })
    await db.sync_queue.add({
      method: 'PUT',
      endpoint: `${API_BASE}/api/v1/rendas-variaveis/${id}`,
      body: data,
      token,
      created_at: now,
      retries: 0,
    })
    return local
  }
  const renda = await api.atualizarRendaVariavel(token, id, data)
  const now = new Date().toISOString()
  await db.rendas_variaveis.put({ ...renda, updated_at: now, deleted_at: null })
  return renda
}

export async function excluirRendaVariavel(token: string, id: string): Promise<void> {
  if (!navigator.onLine) {
    const now = new Date().toISOString()
    await db.rendas_variaveis.update(id, { deleted_at: now, updated_at: now })
    await db.sync_queue.add({
      method: 'DELETE',
      endpoint: `${API_BASE}/api/v1/rendas-variaveis/${id}`,
      body: null,
      token,
      created_at: now,
      retries: 0,
    })
    return
  }
  await api.excluirRendaVariavel(token, id)
  const now = new Date().toISOString()
  await db.rendas_variaveis.update(id, { deleted_at: now, updated_at: now })
}
