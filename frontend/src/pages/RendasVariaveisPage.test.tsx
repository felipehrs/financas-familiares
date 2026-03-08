import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import { RendasVariaveisPage } from './RendasVariaveisPage'
import type { RendaVariavel } from '@/types/renda_variavel'
import type { Membro } from '@/types/membro'

vi.mock('@/offline/rendas_variaveis', () => ({
  listarRendasVariaveis: vi.fn(),
  criarRendaVariavel: vi.fn(),
  atualizarRendaVariavel: vi.fn(),
  excluirRendaVariavel: vi.fn(),
}))

vi.mock('@/api/membros', () => ({
  listarMembros: vi.fn(),
}))

vi.mock('@/hooks/useAuth', () => ({
  useAuth: () => ({ accessToken: 'token-fake' }),
}))

import { listarRendasVariaveis, criarRendaVariavel, excluirRendaVariavel } from '@/offline/rendas_variaveis'
import { listarMembros } from '@/api/membros'

const rendasFixture: RendaVariavel[] = [
  {
    id: '1',
    descricao: 'Freelance de março',
    membro_id: 'membro-1',
    mes_referencia: 3,
    ano_referencia: 2026,
    valor: 1500.00,
    data_recebimento: '2026-03-15',
  },
]

const membroFixture: Membro = { id: 'membro-1', nome: 'Felipe', relacionamento: 'titular', ativo: true }

function renderPage() {
  return render(
    <MemoryRouter>
      <RendasVariaveisPage />
    </MemoryRouter>
  )
}

describe('RendasVariaveisPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(listarRendasVariaveis).mockResolvedValue(rendasFixture)
    vi.mocked(listarMembros).mockResolvedValue([membroFixture])
  })

  // ─── 1. Exibe "Carregando..." enquanto carrega ────────────────────────────
  it('exibe "Carregando..." enquanto os dados são carregados', () => {
    vi.mocked(listarRendasVariaveis).mockReturnValue(new Promise(() => {}))

    renderPage()

    expect(screen.getByText('Carregando...')).toBeInTheDocument()
  })

  // ─── 2. Exibe lista de rendas após carregar ───────────────────────────────
  it('exibe lista de rendas após carregar', async () => {
    renderPage()

    await waitFor(() => {
      expect(screen.getByText('Freelance de março')).toBeInTheDocument()
    })
    expect(screen.getAllByText('Felipe').length).toBeGreaterThan(0)
  })

  // ─── 3. Exibe "Nenhuma renda variável cadastrada" quando vazia ────────────
  it('exibe "Nenhuma renda variável cadastrada" quando lista está vazia', async () => {
    vi.mocked(listarRendasVariaveis).mockResolvedValue([])

    renderPage()

    await waitFor(() => {
      expect(screen.getByText(/nenhuma renda variável cadastrada/i)).toBeInTheDocument()
    })
  })

  // ─── 4. Abre formulário ao clicar em "Nova renda variável" ───────────────
  it('abre formulário ao clicar em "Nova renda variável"', async () => {
    vi.mocked(listarRendasVariaveis).mockResolvedValue([])
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /nova renda variável/i }))
    await user.click(screen.getByRole('button', { name: /nova renda variável/i }))

    expect(screen.getByLabelText(/descrição/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/membro/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/mês de referência/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/valor/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/data de recebimento/i)).toBeInTheDocument()
  })

  // ─── 5. Submete formulário e chama criarRendaVariavel ────────────────────
  it('submete formulário e chama criarRendaVariavel', async () => {
    vi.mocked(listarRendasVariaveis).mockResolvedValue([])
    vi.mocked(criarRendaVariavel).mockResolvedValue(rendasFixture[0])

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /nova renda variável/i }))
    await user.click(screen.getByRole('button', { name: /nova renda variável/i }))

    await waitFor(() => screen.getByLabelText(/descrição/i))

    await user.type(screen.getByLabelText(/descrição/i), 'Freelance de março')
    await user.selectOptions(screen.getByLabelText(/membro/i), 'membro-1')
    await user.selectOptions(screen.getByLabelText(/mês de referência/i), '3')
    await user.clear(screen.getByLabelText(/ano de referência/i))
    await user.type(screen.getByLabelText(/ano de referência/i), '2026')
    await user.clear(screen.getByLabelText(/valor/i))
    await user.type(screen.getByLabelText(/valor/i), '1500')
    await user.type(screen.getByLabelText(/data de recebimento/i), '2026-03-15')

    await user.click(screen.getByRole('button', { name: /salvar/i }))

    await waitFor(() => {
      expect(criarRendaVariavel).toHaveBeenCalledWith(
        'token-fake',
        expect.objectContaining({
          descricao: 'Freelance de março',
          membro_id: 'membro-1',
          mes_referencia: 3,
          ano_referencia: 2026,
          valor: 1500,
        }),
      )
    })
  })

  // ─── 6. Exclui renda e chama excluirRendaVariavel ─────────────────────────
  it('exclui renda e chama excluirRendaVariavel', async () => {
    vi.mocked(excluirRendaVariavel).mockResolvedValue(undefined)

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByText('Freelance de março'))

    const items = screen.getAllByRole('listitem')
    const rendaItem = items.find((el) => el.textContent?.includes('Freelance de março'))!
    await user.click(within(rendaItem).getByRole('button', { name: /excluir/i }))

    await waitFor(() => {
      expect(excluirRendaVariavel).toHaveBeenCalledWith('token-fake', '1')
    })
  })

  // ─── 7. Exibe erro da API quando carregamento falha ───────────────────────
  it('exibe erro da API quando carregamento falha', async () => {
    vi.mocked(listarRendasVariaveis).mockRejectedValue(new Error('Erro de rede'))

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

    await waitFor(() => screen.getByText('Freelance de março'))

    const items = screen.getAllByRole('listitem')
    const rendaItem = items.find((el) => el.textContent?.includes('Freelance de março'))!
    await user.click(within(rendaItem).getByRole('button', { name: /editar/i }))

    await waitFor(() => {
      const descricaoInput = screen.getByLabelText(/descrição/i) as HTMLInputElement
      expect(descricaoInput.value).toBe('Freelance de março')
    })

    const form = screen.getByRole('button', { name: /salvar/i }).closest('form')!
    const valorInput = within(form).getByLabelText(/valor/i) as HTMLInputElement
    expect(valorInput.value).toBe('1500')
  })
})
