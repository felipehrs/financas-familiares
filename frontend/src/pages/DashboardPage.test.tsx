import { render, screen, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import { DashboardPage } from './DashboardPage'
import type { ResumoMensal } from '@/types/dashboard'

vi.mock('@/api/dashboard', () => ({ buscarResumoMensal: vi.fn() }))
vi.mock('@/hooks/useAuth', () => ({ useAuth: vi.fn(() => ({ accessToken: 'fake-token', logout: vi.fn() })) }))

import { buscarResumoMensal } from '@/api/dashboard'

const resumoFixture: ResumoMensal = {
  mes: 3,
  ano: 2026,
  total_rendas: 5000,
  total_despesas: 1200,
  saldo: 3800,
}

function renderPage() {
  return render(
    <MemoryRouter>
      <DashboardPage />
    </MemoryRouter>,
  )
}

describe('DashboardPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  // ─── 1. Exibe os 3 cards com valores formatados após carregar ─────────────
  it('exibe os 3 cards com valores formatados após carregar', async () => {
    vi.mocked(buscarResumoMensal).mockResolvedValue(resumoFixture)

    renderPage()

    await waitFor(() => {
      expect(screen.getByText('Total Rendas')).toBeInTheDocument()
      expect(screen.getByText('Total Despesas')).toBeInTheDocument()
      expect(screen.getByText('Saldo do Mês')).toBeInTheDocument()
    })

    expect(screen.getByText(/5\.000,00/)).toBeInTheDocument()
    expect(screen.getByText(/1\.200,00/)).toBeInTheDocument()
    expect(screen.getByText(/3\.800,00/)).toBeInTheDocument()
  })

  // ─── 2. Saldo positivo tem classe/cor verde ───────────────────────────────
  it('saldo positivo exibe valor com classe text-green-600', async () => {
    vi.mocked(buscarResumoMensal).mockResolvedValue(resumoFixture)

    renderPage()

    await waitFor(() => screen.getByText('Saldo do Mês'))

    const saldoValor = screen.getByText(/3\.800,00/)
    expect(saldoValor.className).toContain('text-green-600')
  })

  // ─── 3. Saldo negativo tem classe/cor vermelha ────────────────────────────
  it('saldo negativo exibe valor com classe text-red-600', async () => {
    vi.mocked(buscarResumoMensal).mockResolvedValue({
      ...resumoFixture,
      saldo: -500,
    })

    renderPage()

    await waitFor(() => screen.getByText('Saldo do Mês'))

    const saldoValor = screen.getByText(/-R\$/)
    expect(saldoValor.className).toContain('text-red-600')
  })

  // ─── 4. Exibe "Carregando..." enquanto aguarda resposta ───────────────────
  it('exibe "Carregando..." enquanto aguarda resposta', () => {
    vi.mocked(buscarResumoMensal).mockReturnValue(new Promise(() => {}))

    renderPage()

    expect(screen.getByText('Carregando...')).toBeInTheDocument()
  })

  // ─── 5. Exibe mensagem de erro quando API falha ───────────────────────────
  it('exibe mensagem de erro quando API falha', async () => {
    vi.mocked(buscarResumoMensal).mockRejectedValue(new Error('Falha na conexão'))

    renderPage()

    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument()
      expect(screen.getByText('Falha na conexão')).toBeInTheDocument()
    })
  })
})
