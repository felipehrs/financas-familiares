import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import { ContasFixasPage } from './ContasFixasPage'
import type { ContaFixa } from '@/types/conta_fixa'
import type { Membro } from '@/types/membro'

vi.mock('@/offline/contas_fixas', () => ({
  listarContasFixas: vi.fn(),
  criarContaFixa: vi.fn(),
  atualizarContaFixa: vi.fn(),
  alterarAtivoContaFixa: vi.fn(),
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

import { listarContasFixas, criarContaFixa, alterarAtivoContaFixa } from '@/offline/contas_fixas'
import { listarMembros } from '@/api/membros'
import { listarCategorias } from '@/api/categorias'

const contaFixaFixture: ContaFixa = {
  id: 'cf-1',
  descricao: 'Condomínio',
  membro_id: 'm-1',
  categoria_id: null,
  valor: 800,
  dia_vencimento: 10,
  forma_pagamento: 'PIX',
  ativa: true,
}

const membroFixture: Membro = { id: 'm-1', nome: 'Felipe', relacionamento: 'titular', ativo: true }

function renderPage() {
  return render(
    <MemoryRouter>
      <ContasFixasPage />
    </MemoryRouter>
  )
}

describe('ContasFixasPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(listarContasFixas).mockResolvedValue([contaFixaFixture])
    vi.mocked(listarMembros).mockResolvedValue([membroFixture])
    vi.mocked(listarCategorias).mockResolvedValue([])
  })

  // ─── 1. Renderiza lista de contas fixas após carregar ─────────────────────
  it('renderiza lista de contas fixas após carregar', async () => {
    renderPage()

    await waitFor(() => {
      expect(screen.getByText('Condomínio')).toBeInTheDocument()
    })
    expect(screen.getAllByText('Felipe').length).toBeGreaterThan(0)
  })

  // ─── 2. Exibe mensagem quando lista está vazia ─────────────────────────────
  it('exibe mensagem quando lista está vazia', async () => {
    vi.mocked(listarContasFixas).mockResolvedValue([])

    renderPage()

    await waitFor(() => {
      expect(screen.getByText(/nenhuma conta/i)).toBeInTheDocument()
    })
  })

  // ─── 3. Exibe formulário ao clicar em "Nova conta fixa" ───────────────────
  it('exibe formulário ao clicar em "Nova conta fixa"', async () => {
    vi.mocked(listarContasFixas).mockResolvedValue([])
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /nova conta fixa/i }))
    await user.click(screen.getByRole('button', { name: /nova conta fixa/i }))

    expect(screen.getByLabelText(/descrição/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/dia de vencimento/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/forma de pagamento/i)).toBeInTheDocument()
  })

  // ─── 4. Chama criarContaFixa ao submeter formulário válido ───────────────
  it('chama criarContaFixa ao submeter formulário válido', async () => {
    vi.mocked(listarContasFixas).mockResolvedValue([])
    vi.mocked(criarContaFixa).mockResolvedValue(contaFixaFixture)

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /nova conta fixa/i }))
    await user.click(screen.getByRole('button', { name: /nova conta fixa/i }))

    await waitFor(() => screen.getByLabelText(/descrição/i))

    await user.type(screen.getByLabelText(/descrição/i), 'Condomínio')
    await user.selectOptions(screen.getByLabelText(/membro responsável/i), 'm-1')

    await user.clear(screen.getByLabelText(/valor/i))
    await user.type(screen.getByLabelText(/valor/i), '800')

    await user.clear(screen.getByLabelText(/dia de vencimento/i))
    await user.type(screen.getByLabelText(/dia de vencimento/i), '10')

    await user.selectOptions(screen.getByLabelText(/forma de pagamento/i), 'PIX')

    await user.click(screen.getByRole('button', { name: /salvar/i }))

    await waitFor(() => {
      expect(criarContaFixa).toHaveBeenCalledWith(
        'fake-token',
        expect.objectContaining({ descricao: 'Condomínio', membro_id: 'm-1', valor: 800 }),
      )
    })
  })

  // ─── 5. Exibe erro de validação quando descrição está vazia ───────────────
  it('exibe erro de validação quando descrição está vazia', async () => {
    vi.mocked(listarContasFixas).mockResolvedValue([])
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /nova conta fixa/i }))
    await user.click(screen.getByRole('button', { name: /nova conta fixa/i }))
    await user.click(screen.getByRole('button', { name: /salvar/i }))

    await waitFor(() => {
      expect(screen.getByText(/descrição é obrigatória/i)).toBeInTheDocument()
    })
  })

  // ─── 6. Chama alterarAtivoContaFixa ao clicar em Inativar ─────────────────
  it('chama alterarAtivoContaFixa ao clicar em Inativar', async () => {
    vi.mocked(alterarAtivoContaFixa).mockResolvedValue({ ...contaFixaFixture, ativa: false })

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByText('Condomínio'))

    const items = screen.getAllByRole('listitem')
    const condominioItem = items.find((el) => el.textContent?.includes('Condomínio'))!
    await user.click(within(condominioItem).getByRole('button', { name: /inativar/i }))

    await waitFor(() => {
      expect(alterarAtivoContaFixa).toHaveBeenCalledWith('fake-token', 'cf-1', false)
    })
  })

  // ─── 7. Exibe membros no select do formulário ─────────────────────────────
  it('exibe membros no select do formulário', async () => {
    vi.mocked(listarContasFixas).mockResolvedValue([])
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /nova conta fixa/i }))
    await user.click(screen.getByRole('button', { name: /nova conta fixa/i }))

    await waitFor(() => screen.getByLabelText(/membro responsável/i))

    const select = screen.getByLabelText(/membro responsável/i) as HTMLSelectElement
    const options = Array.from(select.options).map((o) => o.text)
    expect(options).toContain('Felipe')
  })

  // ─── 8. Abre formulário de edição com dados preenchidos ───────────────────
  it('abre formulário de edição preenchido ao clicar em Editar', async () => {
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByText('Condomínio'))

    const items = screen.getAllByRole('listitem')
    const condominioItem = items.find((el) => el.textContent?.includes('Condomínio'))!
    await user.click(within(condominioItem).getByRole('button', { name: /editar/i }))

    await waitFor(() => {
      const descricaoInput = screen.getByLabelText(/descrição/i) as HTMLInputElement
      expect(descricaoInput.value).toBe('Condomínio')
    })
  })
})
