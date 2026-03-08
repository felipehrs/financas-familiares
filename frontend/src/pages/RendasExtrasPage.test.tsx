import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import { RendasExtrasPage } from './RendasExtrasPage'
import type { RendaExtra } from '@/types/renda_extra'
import type { Membro } from '@/types/membro'

vi.mock('@/offline/rendas_extras', () => ({
  listarRendasExtras: vi.fn(),
  criarRendaExtra: vi.fn(),
  atualizarRendaExtra: vi.fn(),
  excluirRendaExtra: vi.fn(),
}))

vi.mock('@/api/membros', () => ({
  listarMembros: vi.fn(),
}))

vi.mock('@/hooks/useAuth', () => ({
  useAuth: () => ({ accessToken: 'token-fake' }),
}))

import { listarRendasExtras, criarRendaExtra, excluirRendaExtra } from '@/offline/rendas_extras'
import { listarMembros } from '@/api/membros'

const rendasFixture: RendaExtra[] = [
  {
    id: '1',
    descricao: 'Bônus de projeto',
    membro_id: 'membro-1',
    data_recebimento: '2026-03-20',
    valor: 2000.00,
  },
]

const membroFixture: Membro = { id: 'membro-1', nome: 'Felipe', relacionamento: 'titular', ativo: true }

function renderPage() {
  return render(
    <MemoryRouter>
      <RendasExtrasPage />
    </MemoryRouter>
  )
}

describe('RendasExtrasPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(listarRendasExtras).mockResolvedValue(rendasFixture)
    vi.mocked(listarMembros).mockResolvedValue([membroFixture])
  })

  // ─── 1. Exibe "Carregando..." enquanto carrega ────────────────────────────
  it('exibe "Carregando..." enquanto os dados são carregados', () => {
    vi.mocked(listarRendasExtras).mockReturnValue(new Promise(() => {}))

    renderPage()

    expect(screen.getByText('Carregando...')).toBeInTheDocument()
  })

  // ─── 2. Exibe lista de rendas após carregar ───────────────────────────────
  it('exibe lista de rendas após carregar', async () => {
    renderPage()

    await waitFor(() => {
      expect(screen.getByText('Bônus de projeto')).toBeInTheDocument()
    })
    expect(screen.getAllByText('Felipe').length).toBeGreaterThan(0)
  })

  // ─── 3. Exibe "Nenhuma renda extra cadastrada" quando vazia ───────────────
  it('exibe "Nenhuma renda extra cadastrada" quando lista está vazia', async () => {
    vi.mocked(listarRendasExtras).mockResolvedValue([])

    renderPage()

    await waitFor(() => {
      expect(screen.getByText(/nenhuma renda extra cadastrada/i)).toBeInTheDocument()
    })
  })

  // ─── 4. Abre formulário ao clicar em "Nova renda extra" ──────────────────
  it('abre formulário ao clicar em "Nova renda extra"', async () => {
    vi.mocked(listarRendasExtras).mockResolvedValue([])
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /nova renda extra/i }))
    await user.click(screen.getByRole('button', { name: /nova renda extra/i }))

    expect(screen.getByLabelText(/descrição/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/membro/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/data de recebimento/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/valor/i)).toBeInTheDocument()
  })

  // ─── 5. Submete formulário e chama criarRendaExtra ────────────────────────
  it('submete formulário e chama criarRendaExtra', async () => {
    vi.mocked(listarRendasExtras).mockResolvedValue([])
    vi.mocked(criarRendaExtra).mockResolvedValue(rendasFixture[0])

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /nova renda extra/i }))
    await user.click(screen.getByRole('button', { name: /nova renda extra/i }))

    await waitFor(() => screen.getByLabelText(/descrição/i))

    await user.type(screen.getByLabelText(/descrição/i), 'Bônus de projeto')
    await user.selectOptions(screen.getByLabelText(/membro/i), 'membro-1')
    await user.type(screen.getByLabelText(/data de recebimento/i), '2026-03-20')
    await user.clear(screen.getByLabelText(/valor/i))
    await user.type(screen.getByLabelText(/valor/i), '2000')

    await user.click(screen.getByRole('button', { name: /salvar/i }))

    await waitFor(() => {
      expect(criarRendaExtra).toHaveBeenCalledWith(
        'token-fake',
        expect.objectContaining({
          descricao: 'Bônus de projeto',
          membro_id: 'membro-1',
          data_recebimento: '2026-03-20',
          valor: 2000,
        }),
      )
    })
  })

  // ─── 6. Exclui renda e chama excluirRendaExtra ────────────────────────────
  it('exclui renda e chama excluirRendaExtra', async () => {
    vi.mocked(excluirRendaExtra).mockResolvedValue(undefined)

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByText('Bônus de projeto'))

    const items = screen.getAllByRole('listitem')
    const rendaItem = items.find((el) => el.textContent?.includes('Bônus de projeto'))!
    await user.click(within(rendaItem).getByRole('button', { name: /excluir/i }))

    await waitFor(() => {
      expect(excluirRendaExtra).toHaveBeenCalledWith('token-fake', '1')
    })
  })

  // ─── 7. Exibe erro da API quando carregamento falha ───────────────────────
  it('exibe erro da API quando carregamento falha', async () => {
    vi.mocked(listarRendasExtras).mockRejectedValue(new Error('Erro de rede'))

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

    await waitFor(() => screen.getByText('Bônus de projeto'))

    const items = screen.getAllByRole('listitem')
    const rendaItem = items.find((el) => el.textContent?.includes('Bônus de projeto'))!
    await user.click(within(rendaItem).getByRole('button', { name: /editar/i }))

    await waitFor(() => {
      const descricaoInput = screen.getByLabelText(/descrição/i) as HTMLInputElement
      expect(descricaoInput.value).toBe('Bônus de projeto')
    })

    const form = screen.getByRole('button', { name: /salvar/i }).closest('form')!
    const valorInput = within(form).getByLabelText(/valor/i) as HTMLInputElement
    expect(valorInput.value).toBe('2000')
  })
})
