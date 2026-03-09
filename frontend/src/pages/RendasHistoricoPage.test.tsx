import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import { RendasHistoricoPage } from './RendasHistoricoPage'
import type { HistoricoRendas } from '@/types/renda_historico'
import type { Membro } from '@/types/membro'

vi.mock('@/api/rendas_historico')
vi.mock('@/api/membros')
vi.mock('@/hooks/useAuth', () => ({
  useAuth: () => ({ accessToken: 'token-fake' }),
}))

import { buscarHistoricoRendas } from '@/api/rendas_historico'
import { listarMembros } from '@/api/membros'

const mockHistorico: HistoricoRendas = {
  itens: [
    { id: '1', tipo: 'fixa', descricao: 'Salário', membro_id: 'm1', valor: 5000, mes: 3, ano: 2026, ativa: true },
    { id: '2', tipo: 'variavel', descricao: 'Freelance', membro_id: 'm1', valor: 1000, mes: 3, ano: 2026 },
    { id: '3', tipo: 'extra', descricao: 'Bônus', membro_id: 'm2', valor: 500, data: '2026-03-10' },
    { id: '4', tipo: 'investimento', descricao: 'Dividendos', membro_id: 'm1', valor: 2000, valor_distribuido: 200, data: '2026-03-15' },
  ],
  resumo: {
    total_fixas: 5000,
    total_variaveis: 1000,
    total_extras: 500,
    total_investimentos_distribuidos: 200,
    total_geral: 6700,
  },
}

const mockMembros: Membro[] = [
  { id: 'm1', nome: 'Felipe', relacionamento: 'titular', ativo: true },
  { id: 'm2', nome: 'Maria', relacionamento: 'conjuge', ativo: true },
]

function renderPage() {
  return render(
    <MemoryRouter>
      <RendasHistoricoPage />
    </MemoryRouter>,
  )
}

describe('RendasHistoricoPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(buscarHistoricoRendas).mockResolvedValue(mockHistorico)
    vi.mocked(listarMembros).mockResolvedValue(mockMembros)
  })

  // ─── 1. Exibe "Carregando..." enquanto carrega ────────────────────────────
  it('exibe "Carregando..." enquanto os dados são carregados', () => {
    vi.mocked(buscarHistoricoRendas).mockReturnValue(new Promise(() => {}))

    renderPage()

    expect(screen.getByText('Carregando...')).toBeInTheDocument()
  })

  // ─── 2. Exibe itens após carregar ────────────────────────────────────────
  it('exibe itens do histórico após carregar', async () => {
    renderPage()

    await waitFor(() => {
      expect(screen.getByText('Salário')).toBeInTheDocument()
      expect(screen.getByText('Freelance')).toBeInTheDocument()
      expect(screen.getByText('Bônus')).toBeInTheDocument()
      expect(screen.getByText('Dividendos')).toBeInTheDocument()
    })
  })

  // ─── 3. Exibe totais do resumo formatados como moeda ─────────────────────
  it('exibe totais do resumo formatados como moeda', async () => {
    renderPage()

    await waitFor(() => {
      expect(screen.getByText('Resumo')).toBeInTheDocument()
    })

    // total_geral = 6700
    expect(screen.getByText(/6\.700,00/)).toBeInTheDocument()
    // total_variaveis = 1000
    expect(screen.getAllByText(/1\.000,00/).length).toBeGreaterThan(0)
    // total_extras = 500
    expect(screen.getAllByText(/500,00/).length).toBeGreaterThan(0)
    // total_investimentos_distribuidos = 200
    expect(screen.getAllByText(/200,00/).length).toBeGreaterThan(0)
  })

  // ─── 4. Exibe "Nenhuma renda encontrada" quando vazio ────────────────────
  it('exibe "Nenhuma renda encontrada" quando itens está vazio', async () => {
    vi.mocked(buscarHistoricoRendas).mockResolvedValue({
      itens: [],
      resumo: { total_fixas: 0, total_variaveis: 0, total_extras: 0, total_investimentos_distribuidos: 0, total_geral: 0 },
    })

    renderPage()

    await waitFor(() => {
      expect(screen.getByText('Nenhuma renda encontrada')).toBeInTheDocument()
    })
  })

  // ─── 5. Exibe erro quando API falha ──────────────────────────────────────
  it('exibe erro quando API falha', async () => {
    vi.mocked(buscarHistoricoRendas).mockRejectedValue(new Error('Erro de rede'))

    renderPage()

    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument()
      expect(screen.getByText(/erro de rede/i)).toBeInTheDocument()
    })
  })

  // ─── 6. Badge mostra o label correto para cada tipo ───────────────────────
  it('badge exibe label correto para cada tipo de renda', async () => {
    renderPage()

    await waitFor(() => {
      expect(screen.getByText('Renda Fixa')).toBeInTheDocument()
      expect(screen.getByText('Renda Variável')).toBeInTheDocument()
      expect(screen.getByText('Renda Extra')).toBeInTheDocument()
      expect(screen.getByText('Rendimento de Investimento')).toBeInTheDocument()
    })
  })

  // ─── 7. Re-busca quando filtro de tipo muda ───────────────────────────────
  it('re-busca quando o filtro de tipo muda', async () => {
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => {
      expect(screen.getByText('Salário')).toBeInTheDocument()
    })

    // buscarHistoricoRendas should have been called once initially
    expect(buscarHistoricoRendas).toHaveBeenCalledTimes(1)

    // Change tipo filter
    await user.selectOptions(screen.getByLabelText('Tipo'), 'fixa')

    await waitFor(() => {
      expect(buscarHistoricoRendas).toHaveBeenCalledTimes(2)
    })

    expect(buscarHistoricoRendas).toHaveBeenLastCalledWith(
      'token-fake',
      expect.objectContaining({ tipo: 'fixa' }),
    )
  })
})
