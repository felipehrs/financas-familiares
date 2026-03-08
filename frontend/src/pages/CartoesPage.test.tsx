import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import { CartoesPage } from './CartoesPage'
import type { CartaoCredito } from '@/types/cartao_credito'
import type { Membro } from '@/types/membro'

vi.mock('@/api/cartoes_credito', () => ({
  listarCartoes: vi.fn(),
  criarCartao: vi.fn(),
  atualizarCartao: vi.fn(),
  inativarCartao: vi.fn(),
}))

vi.mock('@/api/membros', () => ({
  listarMembros: vi.fn(),
}))

vi.mock('@/hooks/useAuth', () => ({
  useAuth: vi.fn(() => ({ accessToken: 'fake-token' })),
}))

import { listarCartoes, criarCartao, inativarCartao } from '@/api/cartoes_credito'
import { listarMembros } from '@/api/membros'

const membrosFixture: Membro[] = [
  { id: 'm-1', nome: 'Ana Silva', relacionamento: 'Cônjuge', ativo: true },
  { id: 'm-2', nome: 'Carlos Silva', relacionamento: 'Filho', ativo: false },
]

const cartoesFixture: CartaoCredito[] = [
  { id: '1', nome: 'Nubank', membro_id: 'm-1', dia_fechamento: 10, dia_vencimento: 17, limite: 5000, ativo: true },
  { id: '2', nome: 'Inter', membro_id: 'm-1', dia_fechamento: 5, dia_vencimento: 12, limite: null, ativo: false },
]

function renderPage() {
  return render(
    <MemoryRouter>
      <CartoesPage />
    </MemoryRouter>
  )
}

describe('CartoesPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(listarCartoes).mockResolvedValue(cartoesFixture)
    vi.mocked(listarMembros).mockResolvedValue(membrosFixture)
  })

  // ─── 1. Renderiza lista de cartões após carregar ───────────────────────────
  it('renderiza lista de cartões após carregar', async () => {
    renderPage()

    await waitFor(() => {
      expect(screen.getByText('Nubank')).toBeInTheDocument()
      expect(screen.getByText('Inter')).toBeInTheDocument()
    })
    expect(screen.getAllByText('Ana Silva').length).toBeGreaterThan(0)
  })

  // ─── 2. Exibe mensagem quando lista está vazia ─────────────────────────────
  it('exibe "Nenhum cartão cadastrado" quando lista vazia', async () => {
    vi.mocked(listarCartoes).mockResolvedValue([])

    renderPage()

    await waitFor(() => {
      expect(screen.getByText(/nenhum cartão cadastrado/i)).toBeInTheDocument()
    })
  })

  // ─── 3. Exibe formulário ao clicar em "Adicionar cartão" ──────────────────
  it('exibe formulário ao clicar em "Adicionar cartão"', async () => {
    vi.mocked(listarCartoes).mockResolvedValue([])
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /adicionar cartão/i }))
    await user.click(screen.getByRole('button', { name: /adicionar cartão/i }))

    expect(screen.getByLabelText(/nome do cartão/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/dia de fechamento/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/dia de vencimento/i)).toBeInTheDocument()
  })

  // ─── 4. Chama criarCartao ao submeter formulário válido ───────────────────
  it('chama criarCartao ao submeter formulário válido', async () => {
    vi.mocked(listarCartoes).mockResolvedValue([])
    vi.mocked(criarCartao).mockResolvedValue(cartoesFixture[0])

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /adicionar cartão/i }))
    await user.click(screen.getByRole('button', { name: /adicionar cartão/i }))

    // aguardar formulário estar visível
    await waitFor(() => screen.getByLabelText(/nome do cartão/i))

    await user.type(screen.getByLabelText(/nome do cartão/i), 'Nubank')
    await user.selectOptions(screen.getByLabelText(/membro responsável/i), 'm-1')

    // campos type=number: usar userEvent.type para acionar corretamente o react-hook-form
    await user.clear(screen.getByLabelText(/dia de fechamento/i))
    await user.type(screen.getByLabelText(/dia de fechamento/i), '10')
    await user.clear(screen.getByLabelText(/dia de vencimento/i))
    await user.type(screen.getByLabelText(/dia de vencimento/i), '17')

    // limite é opcional mas o campo vazio coerce para 0 (falha em .positive()),
    // então informamos um valor válido
    await user.clear(screen.getByLabelText(/limite/i))
    await user.type(screen.getByLabelText(/limite/i), '5000')

    await user.click(screen.getByRole('button', { name: /salvar/i }))

    await waitFor(() => {
      expect(criarCartao).toHaveBeenCalledWith(
        'fake-token',
        expect.objectContaining({ nome: 'Nubank', membro_id: 'm-1' }),
      )
    })
  })

  // ─── 5. Exibe erro de validação quando nome está vazio ────────────────────
  it('exibe erro de validação quando nome está vazio', async () => {
    vi.mocked(listarCartoes).mockResolvedValue([])
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /adicionar cartão/i }))
    await user.click(screen.getByRole('button', { name: /adicionar cartão/i }))
    await user.click(screen.getByRole('button', { name: /salvar/i }))

    await waitFor(() => {
      expect(screen.getByText(/nome é obrigatório/i)).toBeInTheDocument()
    })
  })

  // ─── 6. Botão "Inativar" aparece apenas para cartões ativos ──────────────
  it('botão "Inativar" aparece apenas para cartões ativos', async () => {
    renderPage()

    await waitFor(() => screen.getByText('Nubank'))

    const items = screen.getAllByRole('listitem')
    const nubankItem = items.find((el) => el.textContent?.includes('Nubank'))!
    const interItem = items.find((el) => el.textContent?.includes('Inter'))!

    expect(within(nubankItem).getByRole('button', { name: /inativar/i })).toBeInTheDocument()
    expect(within(interItem).queryByRole('button', { name: /inativar/i })).not.toBeInTheDocument()
  })

  // ─── 7. Chama inativarCartao ao clicar em "Inativar" ─────────────────────
  it('chama inativarCartao ao clicar em "Inativar"', async () => {
    vi.mocked(inativarCartao).mockResolvedValue(undefined)

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByText('Nubank'))

    const items = screen.getAllByRole('listitem')
    const nubankItem = items.find((el) => el.textContent?.includes('Nubank'))!
    await user.click(within(nubankItem).getByRole('button', { name: /inativar/i }))

    await waitFor(() => {
      expect(inativarCartao).toHaveBeenCalledWith('fake-token', '1')
    })
  })

  // ─── 8. Abre formulário de edição com dados preenchidos ──────────────────
  it('abre formulário de edição ao clicar em "Editar" com dados preenchidos', async () => {
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByText('Nubank'))

    const items = screen.getAllByRole('listitem')
    const nubankItem = items.find((el) => el.textContent?.includes('Nubank'))!
    await user.click(within(nubankItem).getByRole('button', { name: /editar/i }))

    await waitFor(() => {
      const nomeInput = screen.getByLabelText(/nome do cartão/i) as HTMLInputElement
      expect(nomeInput.value).toBe('Nubank')
    })
  })
})
