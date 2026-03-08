import * as api from '@/api/membros'
import { db } from '@/lib/db'
import type { Membro, CriarMembroRequest, AtualizarMembroRequest } from '@/types/membro'

const API_BASE = import.meta.env.VITE_API_URL ?? ''

export async function listarMembros(token: string): Promise<Membro[]> {
  try {
    const membros = await api.listarMembros(token)
    const now = new Date().toISOString()
    await db.membros.bulkPut(
      membros.map(m => ({ ...m, updated_at: now, deleted_at: null })),
    )
    return membros
  } catch {
    const cached = await db.membros.filter(m => m.deleted_at === null).toArray()
    return cached.map(({ updated_at: _u, deleted_at: _d, ...m }) => m)
  }
}

export async function criarMembro(token: string, data: CriarMembroRequest): Promise<Membro> {
  if (!navigator.onLine) {
    const id = crypto.randomUUID()
    const now = new Date().toISOString()
    const local: Membro = {
      id,
      nome: data.nome,
      relacionamento: data.relacionamento ?? '',
      ativo: true,
    }
    await db.membros.add({ ...local, updated_at: now, deleted_at: null })
    await db.sync_queue.add({
      method: 'POST',
      endpoint: `${API_BASE}/api/v1/membros`,
      body: data,
      token,
      created_at: now,
      retries: 0,
    })
    return local
  }
  const membro = await api.criarMembro(token, data)
  const now = new Date().toISOString()
  await db.membros.put({ ...membro, updated_at: now, deleted_at: null })
  return membro
}

export async function atualizarMembro(
  token: string,
  id: string,
  data: AtualizarMembroRequest,
): Promise<Membro> {
  if (!navigator.onLine) {
    const now = new Date().toISOString()
    const existing = await db.membros.get(id)
    const local: Membro = {
      id,
      nome: data.nome,
      relacionamento: data.relacionamento ?? existing?.relacionamento ?? '',
      ativo: data.ativo,
    }
    await db.membros.put({ ...local, updated_at: now, deleted_at: null })
    await db.sync_queue.add({
      method: 'PUT',
      endpoint: `${API_BASE}/api/v1/membros/${id}`,
      body: data,
      token,
      created_at: now,
      retries: 0,
    })
    return local
  }
  const membro = await api.atualizarMembro(token, id, data)
  const now = new Date().toISOString()
  await db.membros.put({ ...membro, updated_at: now, deleted_at: null })
  return membro
}

export async function inativarMembro(token: string, id: string): Promise<void> {
  if (!navigator.onLine) {
    const now = new Date().toISOString()
    await db.membros.update(id, { ativo: false, updated_at: now })
    await db.sync_queue.add({
      method: 'PATCH',
      endpoint: `${API_BASE}/api/v1/membros/${id}/inativar`,
      body: null,
      token,
      created_at: now,
      retries: 0,
    })
    return
  }
  await api.inativarMembro(token, id)
  const now = new Date().toISOString()
  await db.membros.update(id, { ativo: false, updated_at: now })
}
