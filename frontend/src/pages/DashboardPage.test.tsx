import { render, screen, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import { DashboardPage } from './DashboardPage'
import type { ResumoMensal } from '@/types/dashboard'

vi.mock('@/api/dashboard', () => ({
  buscarResumoMensal: vi.fn(),
  buscarCategoriasDespesas: vi.fn(),
}))
vi.mock('@/hooks/useAuth', () => ({ useAuth: vi.fn(() => ({ accessToken: 'fake-token', logout: vi.fn() })) }))

import { buscarResumoMensal, buscarCategoriasDespesas } from '@/api/dashboard'

const resumoFixture: ResumoMensal = {
  mes: 3,
  ano: 2026,
  total_renda_fixa: 5000,
  total_renda_variavel: 1500,
  total_renda_extra: 500,
  total_rendimento_distribuido: 200,
  total_rendas_operacionais: 7200,
  total_rendimento_investimento: 1000,
  total_fatura_cartoes: 800,
  total_assinaturas: 150,
  total_contas_fixas: 300,
  total_despesas_gerais: 450,
  total_despesas: 1700,
  saldo: 5500,
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
    vi.mocked(buscarCategoriasDespesas).mockResolvedValue({
      mes: 3,
      ano: 2026,
      total_despesas: 1700,
      categorias: [
        { nome: 'Alimentação', total: 999, percentual: 58.82 },
        { nome: 'Sem categoria', total: 701, percentual: 41.18 },
      ],
    })
  })

  // ─── 1. Exibe os cards com valores formatados após carregar ───────────────
  it('exibe os cards com valores formatados após carregar', async () => {
    vi.mocked(buscarResumoMensal).mockResolvedValue(resumoFixture)

    renderPage()

    await waitFor(() => {
      expect(screen.getByText('Rendas Operacionais')).toBeInTheDocument()
      expect(screen.getByText('Despesas')).toBeInTheDocument()
      expect(screen.getByText('Saldo do Mês')).toBeInTheDocument()
    })

    expect(screen.getByText(/7\.200,00/)).toBeInTheDocument()
    expect(screen.getByText(/1\.700,00/)).toBeInTheDocument()
    expect(screen.getByText(/5\.500,00/)).toBeInTheDocument()
  })

  // ─── 2. Exibe subtotais de rendas operacionais ────────────────────────────
  it('exibe subtotais de rendas operacionais', async () => {
    vi.mocked(buscarResumoMensal).mockResolvedValue(resumoFixture)

    renderPage()

    await waitFor(() => screen.getByText('Rendas Operacionais'))

    expect(screen.getByText('Renda Fixa')).toBeInTheDocument()
    expect(screen.getByText('Renda Variável')).toBeInTheDocument()
    expect(screen.getByText('Renda Extra')).toBeInTheDocument()
    expect(screen.getByText('Rendimento Distribuído')).toBeInTheDocument()
    expect(screen.getByText('R$ 5.000,00')).toBeInTheDocument()
    expect(screen.getByText('R$ 1.500,00')).toBeInTheDocument()
    expect(screen.getByText('R$ 500,00')).toBeInTheDocument()
  })

  // ─── 3. Exibe seção de rendimentos de investimento ────────────────────────
  it('exibe seção de rendimentos de investimento com texto informativo', async () => {
    vi.mocked(buscarResumoMensal).mockResolvedValue(resumoFixture)

    renderPage()

    await waitFor(() => screen.getByText('Rendimentos de Investimento'))

    expect(screen.getByText('(apenas informativo)')).toBeInTheDocument()
    expect(screen.getByText(/1\.000,00/)).toBeInTheDocument()
  })

  // ─── 4. Exibe subtotais de despesas ───────────────────────────────────────
  it('exibe subtotais de despesas', async () => {
    vi.mocked(buscarResumoMensal).mockResolvedValue(resumoFixture)

    renderPage()

    await waitFor(() => screen.getByText('Despesas'))

    expect(screen.getByText('Fatura Cartões')).toBeInTheDocument()
    expect(screen.getAllByText('Assinaturas').length).toBeGreaterThanOrEqual(1)
    expect(screen.getAllByText('Contas Fixas').length).toBeGreaterThanOrEqual(1)
    expect(screen.getAllByText('Despesas Gerais').length).toBeGreaterThanOrEqual(1)
    expect(screen.getByText('R$ 800,00')).toBeInTheDocument()
    expect(screen.getByText('R$ 150,00')).toBeInTheDocument()
    expect(screen.getByText('R$ 300,00')).toBeInTheDocument()
    expect(screen.getByText('R$ 450,00')).toBeInTheDocument()
  })

  // ─── 5. Saldo positivo tem classe/cor verde ───────────────────────────────
  it('saldo positivo exibe valor com classe text-green-600', async () => {
    vi.mocked(buscarResumoMensal).mockResolvedValue(resumoFixture)

    renderPage()

    await waitFor(() => screen.getByText('Saldo do Mês'))

    const saldoValor = screen.getByText(/5\.500,00/)
    expect(saldoValor.className).toContain('text-green-600')
  })

  // ─── 6. Saldo negativo tem classe/cor vermelha ────────────────────────────
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

  // ─── 7. Exibe "Carregando..." enquanto aguarda resposta ───────────────────
  it('exibe "Carregando..." enquanto aguarda resposta', () => {
    vi.mocked(buscarResumoMensal).mockReturnValue(new Promise(() => {}))

    renderPage()

    expect(screen.getByText('Carregando...')).toBeInTheDocument()
  })

  // ─── 8. Exibe mensagem de erro quando API falha ───────────────────────────
  it('exibe mensagem de erro quando API falha', async () => {
    vi.mocked(buscarResumoMensal).mockRejectedValue(new Error('Falha na conexão'))

    renderPage()

    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument()
      expect(screen.getByText('Falha na conexão')).toBeInTheDocument()
    })
  })

  // ─── 9. Exibe seção de despesas por categoria ─────────────────────────────
  it('exibe seção de despesas por categoria', async () => {
    vi.mocked(buscarResumoMensal).mockResolvedValue(resumoFixture)
    renderPage()
    await waitFor(() => {
      expect(screen.getByText('Despesas por Categoria')).toBeInTheDocument()
      expect(screen.getByText('Alimentação')).toBeInTheDocument()
      expect(screen.getByText('Sem categoria')).toBeInTheDocument()
    })
  })

  // ─── 10. Exibe erro de categorias quando API falha ────────────────────────
  it('exibe erro de categorias quando API falha', async () => {
    vi.mocked(buscarResumoMensal).mockResolvedValue(resumoFixture)
    vi.mocked(buscarCategoriasDespesas).mockRejectedValue(new Error('Erro de categorias'))
    renderPage()
    await waitFor(() => {
      expect(screen.getAllByRole('alert').length).toBeGreaterThan(0)
    })
  })
})
