import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { db } from './db'
import { processarFila } from './syncQueue'

describe('processarFila', () => {
  beforeEach(async () => {
    await db.delete()
    await db.open()
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('remove item da fila após resposta 200', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: true, status: 200 }))
    await db.sync_queue.add({
      method: 'POST',
      endpoint: '/api/v1/despesas-gerais',
      body: { descricao: 'Mercado' },
      token: 'token-teste',
      created_at: '2026-01-01T10:00:00Z',
      retries: 0,
    })
    await processarFila('token-teste')
    expect(await db.sync_queue.count()).toBe(0)
  })

  it('remove item da fila após resposta 409 (last-write-wins)', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: false, status: 409 }))
    await db.sync_queue.add({
      method: 'PUT',
      endpoint: '/api/v1/despesas-gerais/uuid-1',
      body: { descricao: 'Atualizada' },
      token: 'token-teste',
      created_at: '2026-01-01T10:00:00Z',
      retries: 0,
    })
    await processarFila('token-teste')
    expect(await db.sync_queue.count()).toBe(0)
  })

  it('remove item da fila após resposta 404', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: false, status: 404 }))
    await db.sync_queue.add({
      method: 'PUT',
      endpoint: '/api/v1/despesas-gerais/uuid-deletado',
      body: {},
      token: 'token-teste',
      created_at: '2026-01-01T10:00:00Z',
      retries: 0,
    })
    await processarFila('token-teste')
    expect(await db.sync_queue.count()).toBe(0)
  })

  it('incrementa retries quando servidor retorna 500', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: false, status: 500 }))
    await db.sync_queue.add({
      method: 'POST',
      endpoint: '/api/v1/despesas-gerais',
      body: {},
      token: 'token-teste',
      created_at: '2026-01-01T10:00:00Z',
      retries: 0,
    })
    await processarFila('token-teste')
    const item = await db.sync_queue.toCollection().first()
    expect(item?.retries).toBe(1)
  })

  it('incrementa retries quando fetch lança erro de rede', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new Error('Network error')))
    await db.sync_queue.add({
      method: 'DELETE',
      endpoint: '/api/v1/despesas-gerais/uuid-1',
      body: null,
      token: 'token-teste',
      created_at: '2026-01-01T10:00:00Z',
      retries: 2,
    })
    await processarFila('token-teste')
    const item = await db.sync_queue.toCollection().first()
    expect(item?.retries).toBe(3)
  })

  it('processa múltiplos itens na ordem de criação', async () => {
    const fetchMock = vi.fn().mockResolvedValue({ ok: true, status: 200 })
    vi.stubGlobal('fetch', fetchMock)
    await db.sync_queue.bulkAdd([
      { method: 'POST', endpoint: '/api/v1/membros', body: { nome: 'A' }, token: 't', created_at: '2026-01-01T10:00:00Z', retries: 0 },
      { method: 'POST', endpoint: '/api/v1/membros', body: { nome: 'B' }, token: 't', created_at: '2026-01-01T10:00:01Z', retries: 0 },
    ])
    await processarFila('token-teste')
    expect(fetchMock).toHaveBeenCalledTimes(2)
    expect(await db.sync_queue.count()).toBe(0)
  })
})
