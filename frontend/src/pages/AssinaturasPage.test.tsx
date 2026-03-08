import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import { AssinaturasPage } from './AssinaturasPage'
import type { Assinatura } from '@/types/assinatura'
import type { Membro } from '@/types/membro'
import type { Categoria } from '@/types/categoria'

vi.mock('@/api/assinaturas', () => ({
  listarAssinaturas: vi.fn(),
  criarAssinatura: vi.fn(),
  atualizarAssinatura: vi.fn(),
  alterarStatusAssinatura: vi.fn(),
}))

vi.mock('@/api/membros', () => ({
  listarMembros: vi.fn(),
}))

vi.mock('@/api/categorias', () => ({
  listarCategorias: vi.fn(),
}))

vi.mock('@/hooks/useAuth', () => ({
  useAuth: vi.fn(() => ({ accessToken: 'fake-token' })),
}))

import { listarAssinaturas, criarAssinatura, alterarStatusAssinatura } from '@/api/assinaturas'
import { listarMembros } from '@/api/membros'
import { listarCategorias } from '@/api/categorias'

const membrosFixture: Membro[] = [
  { id: 'm-1', nome: 'Felipe', relacionamento: 'titular', ativo: true },
]

const categoriasFixture: Categoria[] = [
  { id: 'cat-1', nome: 'Entretenimento' },
]

const assinaturasFixture: Assinatura[] = [
  {
    id: 'a-1',
    nome: 'Netflix',
    membro_id: 'm-1',
    categoria_id: null,
    valor: 49.90,
    dia_cobranca: 15,
    forma_pagamento: 'Cartão de crédito',
    status: 'ativa',
  },
]

function renderPage() {
  return render(
    <MemoryRouter>
      <AssinaturasPage />
    </MemoryRouter>
  )
}

describe('AssinaturasPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(listarAssinaturas).mockResolvedValue(assinaturasFixture)
    vi.mocked(listarMembros).mockResolvedValue(membrosFixture)
    vi.mocked(listarCategorias).mockResolvedValue(categoriasFixture)
  })

  // ─── 1. Renderiza lista de assinaturas após carregar ──────────────────────
  it('renderiza lista de assinaturas após carregar', async () => {
    renderPage()

    await waitFor(() => {
      expect(screen.getByText('Netflix')).toBeInTheDocument()
    })
    expect(screen.getAllByText('Felipe').length).toBeGreaterThan(0)
  })

  // ─── 2. Exibe mensagem quando lista está vazia ─────────────────────────────
  it('exibe "Nenhuma assinatura" quando lista vazia', async () => {
    vi.mocked(listarAssinaturas).mockResolvedValue([])

    renderPage()

    await waitFor(() => {
      expect(screen.getByText(/nenhuma assinatura/i)).toBeInTheDocument()
    })
  })

  // ─── 3. Exibe formulário ao clicar em "Nova assinatura" ───────────────────
  it('exibe formulário ao clicar em "Nova assinatura"', async () => {
    vi.mocked(listarAssinaturas).mockResolvedValue([])
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /nova assinatura/i }))
    await user.click(screen.getByRole('button', { name: /nova assinatura/i }))

    expect(screen.getByLabelText(/nome da assinatura/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/dia de cobrança/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/forma de pagamento/i)).toBeInTheDocument()
  })

  // ─── 4. Chama criarAssinatura ao submeter formulário válido ───────────────
  it('chama criarAssinatura ao submeter formulário válido', async () => {
    vi.mocked(listarAssinaturas).mockResolvedValue([])
    vi.mocked(criarAssinatura).mockResolvedValue(assinaturasFixture[0])

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /nova assinatura/i }))
    await user.click(screen.getByRole('button', { name: /nova assinatura/i }))

    await waitFor(() => screen.getByLabelText(/nome da assinatura/i))

    await user.type(screen.getByLabelText(/nome da assinatura/i), 'Netflix')
    await user.selectOptions(screen.getByLabelText(/membro responsável/i), 'm-1')

    await user.clear(screen.getByLabelText(/valor/i))
    await user.type(screen.getByLabelText(/valor/i), '49.90')

    await user.clear(screen.getByLabelText(/dia de cobrança/i))
    await user.type(screen.getByLabelText(/dia de cobrança/i), '15')

    await user.selectOptions(screen.getByLabelText(/forma de pagamento/i), 'Cartão de crédito')

    await user.click(screen.getByRole('button', { name: /salvar/i }))

    await waitFor(() => {
      expect(criarAssinatura).toHaveBeenCalledWith(
        'fake-token',
        expect.objectContaining({ nome: 'Netflix', membro_id: 'm-1', valor: 49.90 }),
      )
    })
  })

  // ─── 5. Exibe erro de validação quando nome está vazio ────────────────────
  it('exibe erro de validação quando nome vazio', async () => {
    vi.mocked(listarAssinaturas).mockResolvedValue([])
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /nova assinatura/i }))
    await user.click(screen.getByRole('button', { name: /nova assinatura/i }))
    await user.click(screen.getByRole('button', { name: /salvar/i }))

    await waitFor(() => {
      expect(screen.getByText(/nome é obrigatório/i)).toBeInTheDocument()
    })
  })

  // ─── 6. Chama alterarStatusAssinatura ao pausar assinatura ativa ──────────
  it('chama alterarStatusAssinatura ao pausar assinatura ativa', async () => {
    vi.mocked(alterarStatusAssinatura).mockResolvedValue({ ...assinaturasFixture[0], status: 'pausada' })

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByText('Netflix'))

    const items = screen.getAllByRole('listitem')
    const netflixItem = items.find((el) => el.textContent?.includes('Netflix'))!
    await user.click(within(netflixItem).getByRole('button', { name: /pausar/i }))

    await waitFor(() => {
      expect(alterarStatusAssinatura).toHaveBeenCalledWith('fake-token', 'a-1', 'pausada')
    })
  })

  // ─── 7. Exibe membros no select do formulário ─────────────────────────────
  it('exibe membros no select do formulário', async () => {
    vi.mocked(listarAssinaturas).mockResolvedValue([])
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /nova assinatura/i }))
    await user.click(screen.getByRole('button', { name: /nova assinatura/i }))

    await waitFor(() => screen.getByLabelText(/membro responsável/i))

    const select = screen.getByLabelText(/membro responsável/i) as HTMLSelectElement
    const options = Array.from(select.options).map((o) => o.text)
    expect(options).toContain('Felipe')
  })

  // ─── 8. Abre formulário de edição com dados preenchidos ───────────────────
  it('abre formulário de edição preenchido ao clicar em Editar', async () => {
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByText('Netflix'))

    const items = screen.getAllByRole('listitem')
    const netflixItem = items.find((el) => el.textContent?.includes('Netflix'))!
    await user.click(within(netflixItem).getByRole('button', { name: /editar/i }))

    await waitFor(() => {
      const nomeInput = screen.getByLabelText(/nome da assinatura/i) as HTMLInputElement
      expect(nomeInput.value).toBe('Netflix')
    })
  })
})
