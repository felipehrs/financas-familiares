import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import { RendimentosInvestimentoPage } from './RendimentosInvestimentoPage'
import type { RendimentoInvestimento } from '@/types/rendimento_investimento'
import type { Membro } from '@/types/membro'

vi.mock('@/api/rendimentos_investimento', () => ({
  listarRendimentos: vi.fn(),
  criarRendimento: vi.fn(),
  atualizarRendimento: vi.fn(),
  excluirRendimento: vi.fn(),
}))

vi.mock('@/api/membros', () => ({
  listarMembros: vi.fn(),
}))

vi.mock('@/hooks/useAuth', () => ({
  useAuth: () => ({ accessToken: 'token-fake' }),
}))

import { listarRendimentos, criarRendimento, excluirRendimento } from '@/api/rendimentos_investimento'
import { listarMembros } from '@/api/membros'

const rendimentosFixture: RendimentoInvestimento[] = [
  {
    id: '1',
    descricao: 'Dividendos FII XPML11',
    membro_id: 'membro-1',
    data: '2026-03-15',
    valor: 500.00,
    valor_distribuido: 200.00,
  },
]

const membroFixture: Membro = { id: 'membro-1', nome: 'Felipe', relacionamento: 'titular', ativo: true }

function renderPage() {
  return render(
    <MemoryRouter>
      <RendimentosInvestimentoPage />
    </MemoryRouter>
  )
}

describe('RendimentosInvestimentoPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(listarRendimentos).mockResolvedValue(rendimentosFixture)
    vi.mocked(listarMembros).mockResolvedValue([membroFixture])
  })

  // ─── 1. Exibe "Carregando..." enquanto carrega ────────────────────────────
  it('exibe "Carregando..." enquanto os dados são carregados', () => {
    vi.mocked(listarRendimentos).mockReturnValue(new Promise(() => {}))

    renderPage()

    expect(screen.getByText('Carregando...')).toBeInTheDocument()
  })

  // ─── 2. Exibe lista após carregar ─────────────────────────────────────────
  it('exibe lista de rendimentos após carregar', async () => {
    renderPage()

    await waitFor(() => {
      expect(screen.getByText('Dividendos FII XPML11')).toBeInTheDocument()
    })
    expect(screen.getAllByText('Felipe').length).toBeGreaterThan(0)
  })

  // ─── 3. Exibe "Nenhum rendimento cadastrado" quando vazia ─────────────────
  it('exibe "Nenhum rendimento cadastrado" quando lista está vazia', async () => {
    vi.mocked(listarRendimentos).mockResolvedValue([])

    renderPage()

    await waitFor(() => {
      expect(screen.getByText(/nenhum rendimento cadastrado/i)).toBeInTheDocument()
    })
  })

  // ─── 4. Abre formulário ao clicar em "Novo rendimento" ───────────────────
  it('abre formulário ao clicar em "Novo rendimento"', async () => {
    vi.mocked(listarRendimentos).mockResolvedValue([])
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /novo rendimento/i }))
    await user.click(screen.getByRole('button', { name: /novo rendimento/i }))

    expect(screen.getByLabelText(/descrição/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/membro/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/data/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/valor total recebido/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/valor distribuído/i)).toBeInTheDocument()
  })

  // ─── 5. Submete formulário e chama criarRendimento ────────────────────────
  it('submete formulário e chama criarRendimento', async () => {
    vi.mocked(listarRendimentos).mockResolvedValue([])
    vi.mocked(criarRendimento).mockResolvedValue(rendimentosFixture[0])

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /novo rendimento/i }))
    await user.click(screen.getByRole('button', { name: /novo rendimento/i }))

    await waitFor(() => screen.getByLabelText(/descrição/i))

    await user.type(screen.getByLabelText(/descrição/i), 'Dividendos FII XPML11')
    await user.selectOptions(screen.getByLabelText(/membro/i), 'membro-1')
    await user.type(screen.getByLabelText(/data/i), '2026-03-15')
    await user.clear(screen.getByLabelText(/valor total recebido/i))
    await user.type(screen.getByLabelText(/valor total recebido/i), '500')
    await user.clear(screen.getByLabelText(/valor distribuído/i))
    await user.type(screen.getByLabelText(/valor distribuído/i), '200')

    await user.click(screen.getByRole('button', { name: /salvar/i }))

    await waitFor(() => {
      expect(criarRendimento).toHaveBeenCalledWith(
        'token-fake',
        expect.objectContaining({
          descricao: 'Dividendos FII XPML11',
          membro_id: 'membro-1',
          data: '2026-03-15',
          valor: 500,
          valor_distribuido: 200,
        }),
      )
    })
  })

  // ─── 6. Exclui rendimento e chama excluirRendimento ───────────────────────
  it('exclui rendimento e chama excluirRendimento', async () => {
    vi.mocked(excluirRendimento).mockResolvedValue(undefined)

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByText('Dividendos FII XPML11'))

    const items = screen.getAllByRole('listitem')
    const rendimentoItem = items.find((el) => el.textContent?.includes('Dividendos FII XPML11'))!
    await user.click(within(rendimentoItem).getByRole('button', { name: /excluir/i }))

    await waitFor(() => {
      expect(excluirRendimento).toHaveBeenCalledWith('token-fake', '1')
    })
  })

  // ─── 7. Exibe erro da API quando falha ───────────────────────────────────
  it('exibe erro da API quando carregamento falha', async () => {
    vi.mocked(listarRendimentos).mockRejectedValue(new Error('Erro de rede'))

    renderPage()

    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument()
      expect(screen.getByText(/erro de rede/i)).toBeInTheDocument()
    })
  })

  // ─── 8. Abre formulário de edição preenchido ao clicar em Editar ──────────
  it('abre formulário de edição preenchido ao clicar em Editar', async () => {
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByText('Dividendos FII XPML11'))

    const items = screen.getAllByRole('listitem')
    const rendimentoItem = items.find((el) => el.textContent?.includes('Dividendos FII XPML11'))!
    await user.click(within(rendimentoItem).getByRole('button', { name: /editar/i }))

    await waitFor(() => {
      const descricaoInput = screen.getByLabelText(/descrição/i) as HTMLInputElement
      expect(descricaoInput.value).toBe('Dividendos FII XPML11')
    })

    const form = screen.getByRole('button', { name: /salvar/i }).closest('form')!
    const valorInput = within(form).getByLabelText(/valor total recebido/i) as HTMLInputElement
    expect(valorInput.value).toBe('500')
    const valorDistribuidoInput = within(form).getByLabelText(/valor distribuído/i) as HTMLInputElement
    expect(valorDistribuidoInput.value).toBe('200')
  })
})
