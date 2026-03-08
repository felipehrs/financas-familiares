import * as api from '@/api/cartoes_credito'
import { db } from '@/lib/db'
import type { CartaoCredito, CriarCartaoCreditoRequest, AtualizarCartaoCreditoRequest } from '@/types/cartao_credito'

const API_BASE = import.meta.env.VITE_API_URL ?? ''

export async function listarCartoes(token: string): Promise<CartaoCredito[]> {
  try {
    const cartoes = await api.listarCartoes(token)
    const now = new Date().toISOString()
    await db.cartoes_credito.bulkPut(
      cartoes.map(c => ({ ...c, updated_at: now, deleted_at: null })),
    )
    return cartoes
  } catch {
    const cached = await db.cartoes_credito.filter(c => c.deleted_at === null).toArray()
    return cached.map(({ updated_at: _u, deleted_at: _d, ...c }) => c)
  }
}

export async function criarCartao(token: string, data: CriarCartaoCreditoRequest): Promise<CartaoCredito> {
  if (!navigator.onLine) {
    const id = crypto.randomUUID()
    const now = new Date().toISOString()
    const local: CartaoCredito = {
      id,
      nome: data.nome,
      membro_id: data.membro_id,
      dia_fechamento: data.dia_fechamento,
      dia_vencimento: data.dia_vencimento,
      limite: data.limite ?? null,
      ativo: true,
    }
    await db.cartoes_credito.add({ ...local, updated_at: now, deleted_at: null })
    await db.sync_queue.add({
      method: 'POST',
      endpoint: `${API_BASE}/api/v1/cartoes`,
      body: data,
      token,
      created_at: now,
      retries: 0,
    })
    return local
  }
  const cartao = await api.criarCartao(token, data)
  const now = new Date().toISOString()
  await db.cartoes_credito.put({ ...cartao, updated_at: now, deleted_at: null })
  return cartao
}

export async function atualizarCartao(
  token: string,
  id: string,
  data: AtualizarCartaoCreditoRequest,
): Promise<CartaoCredito> {
  if (!navigator.onLine) {
    const now = new Date().toISOString()
    const local: CartaoCredito = {
      id,
      nome: data.nome,
      membro_id: data.membro_id,
      dia_fechamento: data.dia_fechamento,
      dia_vencimento: data.dia_vencimento,
      limite: data.limite ?? null,
      ativo: data.ativo,
    }
    await db.cartoes_credito.put({ ...local, updated_at: now, deleted_at: null })
    await db.sync_queue.add({
      method: 'PUT',
      endpoint: `${API_BASE}/api/v1/cartoes/${id}`,
      body: data,
      token,
      created_at: now,
      retries: 0,
    })
    return local
  }
  const cartao = await api.atualizarCartao(token, id, data)
  const now = new Date().toISOString()
  await db.cartoes_credito.put({ ...cartao, updated_at: now, deleted_at: null })
  return cartao
}

export async function inativarCartao(token: string, id: string): Promise<void> {
  if (!navigator.onLine) {
    const now = new Date().toISOString()
    await db.cartoes_credito.update(id, { ativo: false, updated_at: now })
    await db.sync_queue.add({
      method: 'PATCH',
      endpoint: `${API_BASE}/api/v1/cartoes/${id}/inativar`,
      body: null,
      token,
      created_at: now,
      retries: 0,
    })
    return
  }
  await api.inativarCartao(token, id)
  const now = new Date().toISOString()
  await db.cartoes_credito.update(id, { ativo: false, updated_at: now })
}
