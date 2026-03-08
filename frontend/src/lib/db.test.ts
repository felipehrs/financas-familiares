import { describe, it, expect, beforeEach } from 'vitest'
import { db } from './db'

describe('FinancasDB', () => {
  beforeEach(async () => {
    await db.delete()
    await db.open()
  })

  it('cria todas as tabelas esperadas', () => {
    const tabelas = db.tables.map(t => t.name)
    expect(tabelas).toContain('despesas_gerais')
    expect(tabelas).toContain('membros')
    expect(tabelas).toContain('categorias')
    expect(tabelas).toContain('cartoes_credito')
    expect(tabelas).toContain('assinaturas')
    expect(tabelas).toContain('contas_fixas')
    expect(tabelas).toContain('rendas_fixas')
    expect(tabelas).toContain('sync_queue')
  })

  it('armazena e recupera um membro', async () => {
    await db.membros.add({
      id: 'uuid-1',
      nome: 'Felipe',
      relacionamento: 'titular',
      ativo: true,
      updated_at: '2026-01-01T00:00:00Z',
      deleted_at: null,
    })
    const result = await db.membros.get('uuid-1')
    expect(result?.nome).toBe('Felipe')
  })

  it('armazena e recupera item na sync_queue', async () => {
    const id = await db.sync_queue.add({
      method: 'POST',
      endpoint: '/api/v1/membros',
      body: { nome: 'Teste' },
      token: 'tok',
      created_at: '2026-01-01T00:00:00Z',
      retries: 0,
    })
    const item = await db.sync_queue.get(id)
    expect(item?.endpoint).toBe('/api/v1/membros')
  })
})
