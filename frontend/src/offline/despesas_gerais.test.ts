import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { db } from '@/lib/db'
import { listarDespesasGerais, criarDespesaGeral, excluirDespesaGeral } from './despesas_gerais'

vi.mock('@/api/despesas_gerais', () => ({
  listarDespesasGerais: vi.fn(),
  criarDespesaGeral: vi.fn(),
  atualizarDespesaGeral: vi.fn(),
  excluirDespesaGeral: vi.fn(),
}))

import * as api from '@/api/despesas_gerais'

const despesaBase = {
  id: 'uuid-1',
  membro_id: 'm1',
  categoria_id: null,
  descricao: 'Mercado',
  data: '2026-03-15',
  valor: 100,
  forma_pagamento: 'pix',
  observacoes: null,
}

describe('offline/despesas_gerais', () => {
  beforeEach(async () => {
    await db.delete()
    await db.open()
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  describe('listarDespesasGerais', () => {
    it('retorna dados da API e atualiza cache', async () => {
      vi.mocked(api.listarDespesasGerais).mockResolvedValue([despesaBase])
      const result = await listarDespesasGerais('token')
      expect(result).toEqual([despesaBase])
      expect(await db.despesas_gerais.count()).toBe(1)
    })

    it('retorna do cache quando API falha', async () => {
      await db.despesas_gerais.add({ ...despesaBase, updated_at: '2026-03-15T00:00:00Z', deleted_at: null })
      vi.mocked(api.listarDespesasGerais).mockRejectedValue(new Error('offline'))
      const result = await listarDespesasGerais('token')
      expect(result.length).toBe(1)
      expect(result[0].descricao).toBe('Mercado')
    })

    it('não retorna registros com deleted_at no fallback', async () => {
      await db.despesas_gerais.add({ ...despesaBase, updated_at: '2026-03-15T00:00:00Z', deleted_at: '2026-03-16T00:00:00Z' })
      vi.mocked(api.listarDespesasGerais).mockRejectedValue(new Error('offline'))
      const result = await listarDespesasGerais('token')
      expect(result).toHaveLength(0)
    })
  })

  describe('criarDespesaGeral', () => {
    it('salva localmente e enfileira quando offline', async () => {
      vi.stubGlobal('navigator', { onLine: false })
      const result = await criarDespesaGeral('token', {
        membro_id: 'm1',
        descricao: 'Lanche',
        data: '2026-03-08',
        valor: 25,
        forma_pagamento: 'dinheiro',
      })
      expect(result.descricao).toBe('Lanche')
      expect(await db.despesas_gerais.count()).toBe(1)
      expect(await db.sync_queue.count()).toBe(1)
      const qi = await db.sync_queue.toCollection().first()
      expect(qi?.method).toBe('POST')
    })

    it('chama API e não enfileira quando online', async () => {
      vi.stubGlobal('navigator', { onLine: true })
      vi.mocked(api.criarDespesaGeral).mockResolvedValue({ ...despesaBase, id: 'srv-1' })
      await criarDespesaGeral('token', { membro_id: 'm1', descricao: 'Online', data: '2026-03-08', valor: 200, forma_pagamento: 'pix' })
      expect(await db.sync_queue.count()).toBe(0)
      expect(api.criarDespesaGeral).toHaveBeenCalled()
    })
  })

  describe('excluirDespesaGeral', () => {
    it('soft delete e enfileira quando offline', async () => {
      vi.stubGlobal('navigator', { onLine: false })
      await db.despesas_gerais.add({ ...despesaBase, updated_at: '2026-03-01T00:00:00Z', deleted_at: null })
      await excluirDespesaGeral('token', 'uuid-1')
      const reg = await db.despesas_gerais.get('uuid-1')
      expect(reg?.deleted_at).not.toBeNull()
      expect(await db.sync_queue.count()).toBe(1)
    })
  })
})
