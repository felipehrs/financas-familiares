import * as api from '@/api/categorias'
import { db } from '@/lib/db'
import type { Categoria, CriarCategoriaRequest, AtualizarCategoriaRequest } from '@/types/categoria'

const API_BASE = import.meta.env.VITE_API_URL ?? ''

export async function listarCategorias(token: string): Promise<Categoria[]> {
  try {
    const categorias = await api.listarCategorias(token)
    const now = new Date().toISOString()
    await db.categorias.bulkPut(
      categorias.map(c => ({ ...c, updated_at: now, deleted_at: null })),
    )
    return categorias
  } catch {
    const cached = await db.categorias.filter(c => c.deleted_at === null).toArray()
    return cached.map(({ updated_at: _u, deleted_at: _d, ...c }) => c)
  }
}

export async function criarCategoria(token: string, data: CriarCategoriaRequest): Promise<Categoria> {
  if (!navigator.onLine) {
    const id = crypto.randomUUID()
    const now = new Date().toISOString()
    const local: Categoria = { id, nome: data.nome }
    await db.categorias.add({ ...local, updated_at: now, deleted_at: null })
    await db.sync_queue.add({
      method: 'POST',
      endpoint: `${API_BASE}/api/v1/categorias`,
      body: data,
      token,
      created_at: now,
      retries: 0,
    })
    return local
  }
  const categoria = await api.criarCategoria(token, data)
  const now = new Date().toISOString()
  await db.categorias.put({ ...categoria, updated_at: now, deleted_at: null })
  return categoria
}

export async function atualizarCategoria(
  token: string,
  id: string,
  data: AtualizarCategoriaRequest,
): Promise<Categoria> {
  if (!navigator.onLine) {
    const now = new Date().toISOString()
    const local: Categoria = { id, nome: data.nome }
    await db.categorias.put({ ...local, updated_at: now, deleted_at: null })
    await db.sync_queue.add({
      method: 'PUT',
      endpoint: `${API_BASE}/api/v1/categorias/${id}`,
      body: data,
      token,
      created_at: now,
      retries: 0,
    })
    return local
  }
  const categoria = await api.atualizarCategoria(token, id, data)
  const now = new Date().toISOString()
  await db.categorias.put({ ...categoria, updated_at: now, deleted_at: null })
  return categoria
}

export async function excluirCategoria(token: string, id: string): Promise<void> {
  if (!navigator.onLine) {
    const now = new Date().toISOString()
    await db.categorias.update(id, { deleted_at: now, updated_at: now })
    await db.sync_queue.add({
      method: 'DELETE',
      endpoint: `${API_BASE}/api/v1/categorias/${id}`,
      body: null,
      token,
      created_at: now,
      retries: 0,
    })
    return
  }
  await api.excluirCategoria(token, id)
  const now = new Date().toISOString()
  await db.categorias.update(id, { deleted_at: now, updated_at: now })
}
