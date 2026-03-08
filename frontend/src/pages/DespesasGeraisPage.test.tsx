import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import { DespesasGeraisPage } from './DespesasGeraisPage'
import type { DespesaGeral } from '@/types/despesa_geral'
import type { Membro } from '@/types/membro'

vi.mock('@/api/despesas_gerais', () => ({
  listarDespesasGerais: vi.fn(),
  criarDespesaGeral: vi.fn(),
  atualizarDespesaGeral: vi.fn(),
  excluirDespesaGeral: vi.fn(),
}))

vi.mock('@/api/membros', () => ({
  listarMembros: vi.fn(),
}))

vi.mock('@/api/categorias', () => ({
  listarCategorias: vi.fn(),
}))

vi.mock('@/hooks/useAuth', () => ({
  useAuth: () => ({ accessToken: 'token-fake' }),
}))

import { listarDespesasGerais, criarDespesaGeral, atualizarDespesaGeral, excluirDespesaGeral } from '@/api/despesas_gerais'
import { listarMembros } from '@/api/membros'
import { listarCategorias } from '@/api/categorias'

const despesasFixture: DespesaGeral[] = [
  {
    id: '1',
    membro_id: 'membro-1',
    categoria_id: null,
    descricao: 'Supermercado',
    data: '2026-03-05',
    valor: 350.00,
    forma_pagamento: 'pix',
    observacoes: null,
  },
]

const membroFixture: Membro = { id: 'membro-1', nome: 'Felipe', relacionamento: 'titular', ativo: true }

function renderPage() {
  return render(
    <MemoryRouter>
      <DespesasGeraisPage />
    </MemoryRouter>
  )
}

describe('DespesasGeraisPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(listarDespesasGerais).mockResolvedValue(despesasFixture)
    vi.mocked(listarMembros).mockResolvedValue([membroFixture])
    vi.mocked(listarCategorias).mockResolvedValue([])
  })

  // ─── 1. Exibe "Carregando..." enquanto carrega ────────────────────────────
  it('exibe "Carregando..." enquanto os dados são carregados', () => {
    vi.mocked(listarDespesasGerais).mockReturnValue(new Promise(() => {}))

    renderPage()

    expect(screen.getByText('Carregando...')).toBeInTheDocument()
  })

  // ─── 2. Exibe lista de despesas após carregar ─────────────────────────────
  it('exibe lista de despesas após carregar', async () => {
    renderPage()

    await waitFor(() => {
      expect(screen.getByText('Supermercado')).toBeInTheDocument()
    })
    expect(screen.getAllByText('Felipe').length).toBeGreaterThan(0)
    expect(screen.getByText(/pix/i)).toBeInTheDocument()
  })

  // ─── 3. Exibe "Nenhuma despesa cadastrada" quando lista vazia ─────────────
  it('exibe "Nenhuma despesa cadastrada" quando lista está vazia', async () => {
    vi.mocked(listarDespesasGerais).mockResolvedValue([])

    renderPage()

    await waitFor(() => {
      expect(screen.getByText(/nenhuma despesa cadastrada/i)).toBeInTheDocument()
    })
  })

  // ─── 4. Abre formulário ao clicar em "Nova despesa" ──────────────────────
  it('abre formulário ao clicar em "Nova despesa"', async () => {
    vi.mocked(listarDespesasGerais).mockResolvedValue([])
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /nova despesa/i }))
    await user.click(screen.getByRole('button', { name: /nova despesa/i }))

    expect(screen.getByLabelText(/data/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/membro responsável/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/descrição/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/forma de pagamento/i)).toBeInTheDocument()
  })

  // ─── 5. Submete formulário e recarrega lista ──────────────────────────────
  it('submete formulário de criação e recarrega lista', async () => {
    vi.mocked(listarDespesasGerais).mockResolvedValue([])
    vi.mocked(criarDespesaGeral).mockResolvedValue(despesasFixture[0])

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /nova despesa/i }))
    await user.click(screen.getByRole('button', { name: /nova despesa/i }))

    await waitFor(() => screen.getByLabelText(/descrição/i))

    await user.type(screen.getByLabelText(/data/i), '2026-03-05')
    await user.selectOptions(screen.getByLabelText(/membro responsável/i), 'membro-1')
    await user.type(screen.getByLabelText(/descrição/i), 'Supermercado')

    await user.clear(screen.getByLabelText(/valor/i))
    await user.type(screen.getByLabelText(/valor/i), '350')

    await user.selectOptions(screen.getByLabelText(/forma de pagamento/i), 'pix')

    await user.click(screen.getByRole('button', { name: /salvar/i }))

    await waitFor(() => {
      expect(criarDespesaGeral).toHaveBeenCalledWith(
        'token-fake',
        expect.objectContaining({
          descricao: 'Supermercado',
          membro_id: 'membro-1',
          valor: 350,
          forma_pagamento: 'pix',
        }),
      )
    })
  })

  // ─── 6. Exclui despesa e recarrega lista ─────────────────────────────────
  it('exclui despesa e recarrega lista', async () => {
    vi.mocked(excluirDespesaGeral).mockResolvedValue(undefined)

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByText('Supermercado'))

    const items = screen.getAllByRole('listitem')
    const supermercadoItem = items.find((el) => el.textContent?.includes('Supermercado'))!
    await user.click(within(supermercadoItem).getByRole('button', { name: /excluir/i }))

    await waitFor(() => {
      expect(excluirDespesaGeral).toHaveBeenCalledWith('token-fake', '1')
    })
  })

  // ─── 7. Exibe erro da API quando carregamento falha ───────────────────────
  it('exibe erro da API quando carregamento falha', async () => {
    vi.mocked(listarDespesasGerais).mockRejectedValue(new Error('Erro de rede'))

    renderPage()

    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument()
      expect(screen.getByText(/erro de rede/i)).toBeInTheDocument()
    })
  })

  // ─── 8. Formulário de edição preenche campos com valores da despesa ───────
  it('abre formulário de edição preenchido ao clicar em Editar', async () => {
    vi.mocked(atualizarDespesaGeral).mockResolvedValue(despesasFixture[0])
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByText('Supermercado'))

    const items = screen.getAllByRole('listitem')
    const supermercadoItem = items.find((el) => el.textContent?.includes('Supermercado'))!
    await user.click(within(supermercadoItem).getByRole('button', { name: /editar/i }))

    await waitFor(() => {
      const descricaoInput = screen.getByLabelText(/descrição/i) as HTMLInputElement
      expect(descricaoInput.value).toBe('Supermercado')
    })

    const valorInput = screen.getByLabelText(/valor/i) as HTMLInputElement
    expect(valorInput.value).toBe('350')

    const formaPagamentoSelect = screen.getByLabelText(/forma de pagamento/i) as HTMLSelectElement
    expect(formaPagamentoSelect.value).toBe('pix')
  })
})
