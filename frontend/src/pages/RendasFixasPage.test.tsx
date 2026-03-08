import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import { RendasFixasPage } from './RendasFixasPage'
import type { RendaFixa } from '@/types/renda_fixa'
import type { Membro } from '@/types/membro'

vi.mock('@/api/rendas_fixas', () => ({
  listarRendasFixas: vi.fn(),
  criarRendaFixa: vi.fn(),
  atualizarRendaFixa: vi.fn(),
  inativarRendaFixa: vi.fn(),
}))

vi.mock('@/api/membros', () => ({ listarMembros: vi.fn() }))

vi.mock('@/hooks/useAuth', () => ({ useAuth: vi.fn(() => ({ accessToken: 'fake-token' })) }))

import { listarRendasFixas, criarRendaFixa, inativarRendaFixa } from '@/api/rendas_fixas'
import { listarMembros } from '@/api/membros'

const membrosFixture: Membro[] = [
  { id: 'm-1', nome: 'Ana Silva', relacionamento: 'Cônjuge', ativo: true },
]

const rendasFixture: RendaFixa[] = [
  { id: 'r-1', descricao: 'Salário Ana', membro_id: 'm-1', valor: 5000, dia_recebimento: 5, ativa: true, data_inicio: '2026-01-01', data_fim: null },
  { id: 'r-2', descricao: 'Salário Carlos', membro_id: 'm-1', valor: 3000, dia_recebimento: 10, ativa: false, data_inicio: '2025-03-01', data_fim: '2026-02-28' },
]

function renderPage() {
  return render(
    <MemoryRouter>
      <RendasFixasPage />
    </MemoryRouter>
  )
}

describe('RendasFixasPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(listarRendasFixas).mockResolvedValue(rendasFixture)
    vi.mocked(listarMembros).mockResolvedValue(membrosFixture)
  })

  // ─── 1. Renderiza lista de rendas fixas após carregar ─────────────────────
  it('renderiza lista de rendas fixas após carregar', async () => {
    renderPage()

    await waitFor(() => {
      expect(screen.getByText('Salário Ana')).toBeInTheDocument()
      expect(screen.getByText('Salário Carlos')).toBeInTheDocument()
    })
    expect(screen.getAllByText('Ana Silva').length).toBeGreaterThan(0)
  })

  // ─── 2. Exibe mensagem quando lista está vazia ─────────────────────────────
  it('exibe "Nenhuma renda fixa cadastrada" quando lista vazia', async () => {
    vi.mocked(listarRendasFixas).mockResolvedValue([])

    renderPage()

    await waitFor(() => {
      expect(screen.getByText(/nenhuma renda fixa cadastrada/i)).toBeInTheDocument()
    })
  })

  // ─── 3. Exibe formulário ao clicar em "Adicionar renda fixa" ──────────────
  it('exibe formulário ao clicar em "Adicionar renda fixa"', async () => {
    vi.mocked(listarRendasFixas).mockResolvedValue([])
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /adicionar renda fixa/i }))
    await user.click(screen.getByRole('button', { name: /adicionar renda fixa/i }))

    expect(screen.getByLabelText(/descrição/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/valor/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/dia de recebimento/i)).toBeInTheDocument()
  })

  // ─── 4. Chama criarRendaFixa ao submeter formulário válido ────────────────
  it('chama criarRendaFixa ao submeter formulário válido', async () => {
    vi.mocked(listarRendasFixas).mockResolvedValue([])
    vi.mocked(criarRendaFixa).mockResolvedValue(rendasFixture[0])

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /adicionar renda fixa/i }))
    await user.click(screen.getByRole('button', { name: /adicionar renda fixa/i }))

    await waitFor(() => screen.getByLabelText(/descrição/i))

    await user.type(screen.getByLabelText(/descrição/i), 'Salário Ana')
    await user.selectOptions(screen.getByLabelText(/membro responsável/i), 'm-1')

    await user.clear(screen.getByLabelText(/valor/i))
    await user.type(screen.getByLabelText(/valor/i), '5000')
    await user.clear(screen.getByLabelText(/dia de recebimento/i))
    await user.type(screen.getByLabelText(/dia de recebimento/i), '5')
    await user.type(screen.getByLabelText(/data de início/i), '2026-01-01')

    await user.click(screen.getByRole('button', { name: /salvar/i }))

    await waitFor(() => {
      expect(criarRendaFixa).toHaveBeenCalledWith(
        'fake-token',
        expect.objectContaining({ descricao: 'Salário Ana', membro_id: 'm-1', data_inicio: '2026-01-01' }),
      )
    })
  })

  // ─── 5. Exibe erro de validação quando descrição está vazia ───────────────
  it('exibe erro de validação quando descrição está vazia', async () => {
    vi.mocked(listarRendasFixas).mockResolvedValue([])
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /adicionar renda fixa/i }))
    await user.click(screen.getByRole('button', { name: /adicionar renda fixa/i }))
    await user.click(screen.getByRole('button', { name: /salvar/i }))

    await waitFor(() => {
      expect(screen.getByText(/descrição é obrigatória/i)).toBeInTheDocument()
    })
  })

  // ─── 6. Botão "Inativar" aparece apenas para rendas ativas ───────────────
  it('botão "Inativar" aparece apenas para rendas ativas', async () => {
    renderPage()

    await waitFor(() => screen.getByText('Salário Ana'))

    const items = screen.getAllByRole('listitem')
    const rendaAtivaItem = items.find((el) => el.textContent?.includes('Salário Ana'))!
    const rendaInativaItem = items.find((el) => el.textContent?.includes('Salário Carlos'))!

    expect(within(rendaAtivaItem).getByRole('button', { name: /inativar/i })).toBeInTheDocument()
    expect(within(rendaInativaItem).queryByRole('button', { name: /inativar/i })).not.toBeInTheDocument()
  })

  // ─── 7. Chama inativarRendaFixa ao clicar em "Inativar" ──────────────────
  it('chama inativarRendaFixa ao clicar em "Inativar"', async () => {
    vi.mocked(inativarRendaFixa).mockResolvedValue(undefined)

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByText('Salário Ana'))

    const items = screen.getAllByRole('listitem')
    const rendaAtivaItem = items.find((el) => el.textContent?.includes('Salário Ana'))!
    await user.click(within(rendaAtivaItem).getByRole('button', { name: /inativar/i }))

    await waitFor(() => {
      expect(inativarRendaFixa).toHaveBeenCalledWith('fake-token', 'r-1')
    })
  })

  // ─── 8. Abre formulário de edição com dados preenchidos ──────────────────
  it('abre formulário de edição ao clicar em "Editar" com dados preenchidos', async () => {
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByText('Salário Ana'))

    const items = screen.getAllByRole('listitem')
    const rendaAtivaItem = items.find((el) => el.textContent?.includes('Salário Ana'))!
    await user.click(within(rendaAtivaItem).getByRole('button', { name: /editar/i }))

    await waitFor(() => {
      const descricaoInput = screen.getByLabelText(/descrição/i) as HTMLInputElement
      expect(descricaoInput.value).toBe('Salário Ana')
    })
  })

  // ─── 9. Exibe data de início no formulário de edição preenchido ───────────
  it('exibe data de início no formulário de edição preenchido', async () => {
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByText('Salário Ana'))

    const items = screen.getAllByRole('listitem')
    const rendaAtivaItem = items.find((el) => el.textContent?.includes('Salário Ana'))!
    await user.click(within(rendaAtivaItem).getByRole('button', { name: /editar/i }))

    await waitFor(() => {
      const dataInicioInput = screen.getByLabelText(/data de início/i) as HTMLInputElement
      expect(dataInicioInput.value).toBe('2026-01-01')
    })
  })

  // ─── 10. Exibe data de fim quando preenchida ──────────────────────────────
  it('exibe data de fim quando preenchida', async () => {
    renderPage()

    await waitFor(() => {
      expect(screen.getByText(/Fim: 2026-02-28/)).toBeInTheDocument()
    })
  })

  // ─── 11. Exibe erro de validação quando data de início está vazia ─────────
  it('exibe erro de validação quando data de início está vazia', async () => {
    vi.mocked(listarRendasFixas).mockResolvedValue([])
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /adicionar renda fixa/i }))
    await user.click(screen.getByRole('button', { name: /adicionar renda fixa/i }))

    await user.type(screen.getByLabelText(/descrição/i), 'Salário Ana')
    await user.selectOptions(screen.getByLabelText(/membro responsável/i), 'm-1')
    await user.type(screen.getByLabelText(/valor/i), '5000')
    await user.type(screen.getByLabelText(/dia de recebimento/i), '5')
    // data_inicio is intentionally left empty

    await user.click(screen.getByRole('button', { name: /salvar/i }))

    await waitFor(() => {
      expect(screen.getByText(/data de início é obrigatória/i)).toBeInTheDocument()
    })
  })
})
