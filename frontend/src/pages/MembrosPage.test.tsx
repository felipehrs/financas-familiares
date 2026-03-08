import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import { MembrosPage } from './MembrosPage'
import type { Membro } from '@/types/membro'

// Mock do módulo de API
vi.mock('@/offline/membros', () => ({
  listarMembros: vi.fn(),
  criarMembro: vi.fn(),
  atualizarMembro: vi.fn(),
  inativarMembro: vi.fn(),
}))

// Mock do useAuth
vi.mock('@/hooks/useAuth', () => ({
  useAuth: vi.fn(() => ({ accessToken: 'fake-token' })),
}))

import {
  listarMembros,
  criarMembro,
  inativarMembro,
  atualizarMembro,
} from '@/offline/membros'

const membrosFixture: Membro[] = [
  { id: '1', nome: 'Ana Silva', relacionamento: 'Cônjuge', ativo: true },
  { id: '2', nome: 'Carlos Silva', relacionamento: 'Filho', ativo: false },
]

function renderPage() {
  return render(
    <MemoryRouter>
      <MembrosPage />
    </MemoryRouter>
  )
}

describe('MembrosPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  // ─── 1. Renderiza lista de membros após carregar ───────────────────────────
  it('renderiza lista de membros após carregar', async () => {
    vi.mocked(listarMembros).mockResolvedValueOnce(membrosFixture)

    renderPage()

    await waitFor(() => {
      expect(screen.getByText('Ana Silva')).toBeInTheDocument()
      expect(screen.getByText('Cônjuge')).toBeInTheDocument()
      expect(screen.getByText('Carlos Silva')).toBeInTheDocument()
      expect(screen.getByText('Filho')).toBeInTheDocument()
    })
  })

  // ─── 2. Exibe mensagem quando lista está vazia ─────────────────────────────
  it('exibe "Nenhum membro cadastrado" quando lista vazia', async () => {
    vi.mocked(listarMembros).mockResolvedValueOnce([])

    renderPage()

    await waitFor(() => {
      expect(screen.getByText(/nenhum membro cadastrado/i)).toBeInTheDocument()
    })
  })

  // ─── 3. Exibe formulário ao clicar em "Adicionar membro" ──────────────────
  it('exibe formulário ao clicar em "Adicionar membro"', async () => {
    vi.mocked(listarMembros).mockResolvedValueOnce([])
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /adicionar membro/i }))
    await user.click(screen.getByRole('button', { name: /adicionar membro/i }))

    expect(screen.getByLabelText(/nome/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/relacionamento/i)).toBeInTheDocument()
  })

  // ─── 4. Chama criarMembro ao submeter formulário válido ───────────────────
  it('chama criarMembro ao submeter formulário válido', async () => {
    vi.mocked(listarMembros)
      .mockResolvedValueOnce([])
      .mockResolvedValueOnce([
        { id: '3', nome: 'Pedro Costa', relacionamento: 'Pai', ativo: true },
      ])
    vi.mocked(criarMembro).mockResolvedValueOnce({
      id: '3',
      nome: 'Pedro Costa',
      relacionamento: 'Pai',
      ativo: true,
    })

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /adicionar membro/i }))
    await user.click(screen.getByRole('button', { name: /adicionar membro/i }))

    await user.type(screen.getByLabelText(/nome/i), 'Pedro Costa')
    await user.type(screen.getByLabelText(/relacionamento/i), 'Pai')
    await user.click(screen.getByRole('button', { name: /salvar/i }))

    await waitFor(() => {
      expect(criarMembro).toHaveBeenCalledWith('fake-token', {
        nome: 'Pedro Costa',
        relacionamento: 'Pai',
      })
    })
  })

  // ─── 5. Exibe erro de validação quando nome está vazio ────────────────────
  it('exibe erro de validação quando nome está vazio', async () => {
    vi.mocked(listarMembros).mockResolvedValueOnce([])
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /adicionar membro/i }))
    await user.click(screen.getByRole('button', { name: /adicionar membro/i }))
    await user.click(screen.getByRole('button', { name: /salvar/i }))

    await waitFor(() => {
      expect(screen.getByText(/nome é obrigatório/i)).toBeInTheDocument()
    })
  })

  // ─── 6. Exibe erro de API quando criarMembro falha ───────────────────────
  it('exibe erro de API quando criarMembro falha', async () => {
    vi.mocked(listarMembros).mockResolvedValueOnce([])
    vi.mocked(criarMembro).mockRejectedValueOnce(new Error('Erro interno do servidor'))

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /adicionar membro/i }))
    await user.click(screen.getByRole('button', { name: /adicionar membro/i }))

    await user.type(screen.getByLabelText(/nome/i), 'Pedro Costa')
    await user.click(screen.getByRole('button', { name: /salvar/i }))

    await waitFor(() => {
      expect(screen.getByText(/erro interno do servidor/i)).toBeInTheDocument()
    })
  })

  // ─── 7. Botão "Inativar" aparece apenas para membros ativos ──────────────
  it('botão "Inativar" aparece apenas para membros ativos', async () => {
    vi.mocked(listarMembros).mockResolvedValueOnce(membrosFixture)

    renderPage()

    await waitFor(() => screen.getByText('Ana Silva'))

    const items = screen.getAllByRole('listitem')
    const anaItem = items.find((el) => el.textContent?.includes('Ana Silva'))!
    const carlosItem = items.find((el) => el.textContent?.includes('Carlos Silva'))!

    expect(within(anaItem).getByRole('button', { name: /inativar/i })).toBeInTheDocument()
    expect(within(carlosItem).queryByRole('button', { name: /inativar/i })).not.toBeInTheDocument()
  })

  // ─── 8. Chama inativarMembro ao clicar em "Inativar" ─────────────────────
  it('chama inativarMembro ao clicar em "Inativar"', async () => {
    vi.mocked(listarMembros)
      .mockResolvedValueOnce(membrosFixture)
      .mockResolvedValueOnce([membrosFixture[1]])
    vi.mocked(inativarMembro).mockResolvedValueOnce(undefined)

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByText('Ana Silva'))

    const items = screen.getAllByRole('listitem')
    const anaItem = items.find((el) => el.textContent?.includes('Ana Silva'))!
    await user.click(within(anaItem).getByRole('button', { name: /inativar/i }))

    await waitFor(() => {
      expect(inativarMembro).toHaveBeenCalledWith('fake-token', '1')
    })
  })

  // ─── 9. Abre formulário de edição com dados preenchidos ──────────────────
  it('abre formulário de edição ao clicar em "Editar" com dados preenchidos', async () => {
    vi.mocked(listarMembros).mockResolvedValueOnce(membrosFixture)
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByText('Ana Silva'))

    const items = screen.getAllByRole('listitem')
    const anaItem = items.find((el) => el.textContent?.includes('Ana Silva'))!
    await user.click(within(anaItem).getByRole('button', { name: /editar/i }))

    await waitFor(() => {
      const nomeInput = screen.getByLabelText(/nome/i) as HTMLInputElement
      const relInput = screen.getByLabelText(/relacionamento/i) as HTMLInputElement
      expect(nomeInput.value).toBe('Ana Silva')
      expect(relInput.value).toBe('Cônjuge')
    })
  })
})
