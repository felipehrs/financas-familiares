import * as api from '@/api/contas_fixas'
import { db } from '@/lib/db'
import type { ContaFixa, CriarContaFixaRequest, AtualizarContaFixaRequest } from '@/types/conta_fixa'

const API_BASE = import.meta.env.VITE_API_URL ?? ''

export async function listarContasFixas(token: string): Promise<ContaFixa[]> {
  try {
    const contas = await api.listarContasFixas(token)
    const now = new Date().toISOString()
    await db.contas_fixas.bulkPut(
      contas.map(c => ({ ...c, updated_at: now, deleted_at: null })),
    )
    return contas
  } catch {
    const cached = await db.contas_fixas.filter(c => c.deleted_at === null).toArray()
    return cached.map(({ updated_at: _u, deleted_at: _d, ...c }) => c)
  }
}

export async function criarContaFixa(token: string, data: CriarContaFixaRequest): Promise<ContaFixa> {
  if (!navigator.onLine) {
    const id = crypto.randomUUID()
    const now = new Date().toISOString()
    const local: ContaFixa = {
      id,
      descricao: data.descricao,
      membro_id: data.membro_id,
      categoria_id: data.categoria_id ?? null,
      valor: data.valor,
      dia_vencimento: data.dia_vencimento,
      forma_pagamento: data.forma_pagamento,
      ativa: true,
    }
    await db.contas_fixas.add({ ...local, updated_at: now, deleted_at: null })
    await db.sync_queue.add({
      method: 'POST',
      endpoint: `${API_BASE}/api/v1/contas-fixas`,
      body: data,
      token,
      created_at: now,
      retries: 0,
    })
    return local
  }
  const conta = await api.criarContaFixa(token, data)
  const now = new Date().toISOString()
  await db.contas_fixas.put({ ...conta, updated_at: now, deleted_at: null })
  return conta
}

export async function atualizarContaFixa(
  token: string,
  id: string,
  data: AtualizarContaFixaRequest,
): Promise<ContaFixa> {
  if (!navigator.onLine) {
    const now = new Date().toISOString()
    const local: ContaFixa = {
      id,
      descricao: data.descricao,
      membro_id: data.membro_id,
      categoria_id: data.categoria_id ?? null,
      valor: data.valor,
      dia_vencimento: data.dia_vencimento,
      forma_pagamento: data.forma_pagamento,
      ativa: data.ativa,
    }
    await db.contas_fixas.put({ ...local, updated_at: now, deleted_at: null })
    await db.sync_queue.add({
      method: 'PUT',
      endpoint: `${API_BASE}/api/v1/contas-fixas/${id}`,
      body: data,
      token,
      created_at: now,
      retries: 0,
    })
    return local
  }
  const conta = await api.atualizarContaFixa(token, id, data)
  const now = new Date().toISOString()
  await db.contas_fixas.put({ ...conta, updated_at: now, deleted_at: null })
  return conta
}

export async function alterarAtivoContaFixa(
  token: string,
  id: string,
  ativa: boolean,
): Promise<ContaFixa> {
  if (!navigator.onLine) {
    const now = new Date().toISOString()
    await db.contas_fixas.update(id, { ativa, updated_at: now })
    const updated = await db.contas_fixas.get(id)
    const local: ContaFixa = {
      id,
      descricao: updated?.descricao ?? '',
      membro_id: updated?.membro_id ?? '',
      categoria_id: updated?.categoria_id ?? null,
      valor: updated?.valor ?? 0,
      dia_vencimento: updated?.dia_vencimento ?? 1,
      forma_pagamento: updated?.forma_pagamento ?? '',
      ativa,
    }
    await db.sync_queue.add({
      method: 'PATCH',
      endpoint: `${API_BASE}/api/v1/contas-fixas/${id}/ativo`,
      body: { ativa },
      token,
      created_at: now,
      retries: 0,
    })
    return local
  }
  const conta = await api.alterarAtivoContaFixa(token, id, ativa)
  const now = new Date().toISOString()
  await db.contas_fixas.put({ ...conta, updated_at: now, deleted_at: null })
  return conta
}
