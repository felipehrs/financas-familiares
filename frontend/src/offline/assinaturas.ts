import * as api from '@/api/assinaturas'
import { db } from '@/lib/db'
import type { Assinatura, CriarAssinaturaRequest, AtualizarAssinaturaRequest } from '@/types/assinatura'

const API_BASE = import.meta.env.VITE_API_URL ?? ''

export async function listarAssinaturas(token: string): Promise<Assinatura[]> {
  try {
    const assinaturas = await api.listarAssinaturas(token)
    const now = new Date().toISOString()
    await db.assinaturas.bulkPut(
      assinaturas.map(a => ({ ...a, updated_at: now, deleted_at: null })),
    )
    return assinaturas
  } catch {
    const cached = await db.assinaturas.filter(a => a.deleted_at === null).toArray()
    return cached.map(({ updated_at: _u, deleted_at: _d, ...a }) => a)
  }
}

export async function criarAssinatura(token: string, data: CriarAssinaturaRequest): Promise<Assinatura> {
  if (!navigator.onLine) {
    const id = crypto.randomUUID()
    const now = new Date().toISOString()
    const local: Assinatura = {
      id,
      nome: data.nome,
      membro_id: data.membro_id,
      categoria_id: data.categoria_id ?? null,
      valor: data.valor,
      dia_cobranca: data.dia_cobranca,
      forma_pagamento: data.forma_pagamento,
      status: 'ativa',
    }
    await db.assinaturas.add({ ...local, updated_at: now, deleted_at: null })
    await db.sync_queue.add({
      method: 'POST',
      endpoint: `${API_BASE}/api/v1/assinaturas`,
      body: data,
      token,
      created_at: now,
      retries: 0,
    })
    return local
  }
  const assinatura = await api.criarAssinatura(token, data)
  const now = new Date().toISOString()
  await db.assinaturas.put({ ...assinatura, updated_at: now, deleted_at: null })
  return assinatura
}

export async function atualizarAssinatura(
  token: string,
  id: string,
  data: AtualizarAssinaturaRequest,
): Promise<Assinatura> {
  if (!navigator.onLine) {
    const now = new Date().toISOString()
    const local: Assinatura = {
      id,
      nome: data.nome,
      membro_id: data.membro_id,
      categoria_id: data.categoria_id ?? null,
      valor: data.valor,
      dia_cobranca: data.dia_cobranca,
      forma_pagamento: data.forma_pagamento,
      status: data.status,
    }
    await db.assinaturas.put({ ...local, updated_at: now, deleted_at: null })
    await db.sync_queue.add({
      method: 'PUT',
      endpoint: `${API_BASE}/api/v1/assinaturas/${id}`,
      body: data,
      token,
      created_at: now,
      retries: 0,
    })
    return local
  }
  const assinatura = await api.atualizarAssinatura(token, id, data)
  const now = new Date().toISOString()
  await db.assinaturas.put({ ...assinatura, updated_at: now, deleted_at: null })
  return assinatura
}

export async function alterarStatusAssinatura(
  token: string,
  id: string,
  status: string,
): Promise<Assinatura> {
  if (!navigator.onLine) {
    const now = new Date().toISOString()
    await db.assinaturas.update(id, { status: status as Assinatura['status'], updated_at: now })
    const updated = await db.assinaturas.get(id)
    const local: Assinatura = {
      id,
      nome: updated?.nome ?? '',
      membro_id: updated?.membro_id ?? '',
      categoria_id: updated?.categoria_id ?? null,
      valor: updated?.valor ?? 0,
      dia_cobranca: updated?.dia_cobranca ?? 1,
      forma_pagamento: updated?.forma_pagamento ?? '',
      status: status as Assinatura['status'],
    }
    await db.sync_queue.add({
      method: 'PATCH',
      endpoint: `${API_BASE}/api/v1/assinaturas/${id}/status`,
      body: { status },
      token,
      created_at: now,
      retries: 0,
    })
    return local
  }
  const assinatura = await api.alterarStatusAssinatura(token, id, status)
  const now = new Date().toISOString()
  await db.assinaturas.put({ ...assinatura, updated_at: now, deleted_at: null })
  return assinatura
}
