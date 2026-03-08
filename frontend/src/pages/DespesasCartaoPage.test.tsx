import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { DespesasCartaoPage } from './DespesasCartaoPage'
import type { DespesaCartao } from '@/types/despesa_cartao'
import type { CartaoCredito } from '@/types/cartao_credito'
import type { Categoria } from '@/types/categoria'

vi.mock('@/api/despesas_cartao', () => ({
  listarDespesasPorCartao: vi.fn(),
  criarDespesa: vi.fn(),
  excluirDespesa: vi.fn(),
}))

vi.mock('@/api/categorias', () => ({
  listarCategorias: vi.fn(),
}))

vi.mock('@/api/cartoes_credito', () => ({
  listarCartoes: vi.fn(),
}))

vi.mock('@/hooks/useAuth', () => ({
  useAuth: vi.fn(() => ({ accessToken: 'fake-token' })),
}))

import { listarDespesasPorCartao, criarDespesa, excluirDespesa } from '@/api/despesas_cartao'
import { listarCategorias } from '@/api/categorias'
import { listarCartoes } from '@/api/cartoes_credito'

const cartoesFixture: CartaoCredito[] = [
  { id: '1', nome: 'Nubank', membro_id: 'm-1', dia_fechamento: 10, dia_vencimento: 17, limite: 5000, ativo: true },
]

const categoriasFixture: Categoria[] = [
  { id: 'cat-1', nome: 'Alimentação' },
  { id: 'cat-2', nome: 'Saúde' },
]

const despesasFixture: DespesaCartao[] = [
  {
    id: 'd1',
    cartao_id: '1',
    descricao: 'Supermercado',
    valor_total: 150,
    valor_parcela: 150,
    numero_parcelas: 1,
    fatura_mes: 3,
    fatura_ano: 2026,
    fatura: 'MAR/26',
    data_compra: '2026-03-10',
    categoria_id: null,
  },
]

function renderPage() {
  return render(
    <MemoryRouter initialEntries={['/cartoes/1/despesas']}>
      <Routes>
        <Route path="/cartoes/:cartaoId/despesas" element={<DespesasCartaoPage />} />
      </Routes>
    </MemoryRouter>
  )
}

describe('DespesasCartaoPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(listarDespesasPorCartao).mockResolvedValue(despesasFixture)
    vi.mocked(listarCartoes).mockResolvedValue(cartoesFixture)
    vi.mocked(listarCategorias).mockResolvedValue(categoriasFixture)
  })

  // ─── 1. Renderiza lista de despesas após carregar ─────────────────────────
  it('renderiza lista de despesas após carregar', async () => {
    renderPage()

    await waitFor(() => {
      expect(screen.getByText('Supermercado')).toBeInTheDocument()
    })
    expect(screen.getByText('MAR/26')).toBeInTheDocument()
  })

  // ─── 2. Exibe mensagem quando lista vazia ─────────────────────────────────
  it('exibe "Nenhuma despesa" quando lista vazia', async () => {
    vi.mocked(listarDespesasPorCartao).mockResolvedValue([])

    renderPage()

    await waitFor(() => {
      expect(screen.getByText(/nenhuma despesa/i)).toBeInTheDocument()
    })
  })

  // ─── 3. Exibe formulário ao clicar em "Nova despesa" ─────────────────────
  it('exibe formulário ao clicar em "Nova despesa"', async () => {
    vi.mocked(listarDespesasPorCartao).mockResolvedValue([])
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /nova despesa/i }))
    await user.click(screen.getByRole('button', { name: /nova despesa/i }))

    expect(screen.getByLabelText(/descrição/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/data da compra/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/valor total/i)).toBeInTheDocument()
  })

  // ─── 4. Chama criarDespesa ao submeter formulário válido ─────────────────
  it('chama criarDespesa ao submeter formulário válido', async () => {
    vi.mocked(listarDespesasPorCartao).mockResolvedValue([])
    vi.mocked(criarDespesa).mockResolvedValue(despesasFixture[0])

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /nova despesa/i }))
    await user.click(screen.getByRole('button', { name: /nova despesa/i }))

    await waitFor(() => screen.getByLabelText(/descrição/i))

    await user.type(screen.getByLabelText(/descrição/i), 'Supermercado')
    await user.type(screen.getByLabelText(/data da compra/i), '2026-03-10')
    await user.clear(screen.getByLabelText(/valor total/i))
    await user.type(screen.getByLabelText(/valor total/i), '150')

    await user.click(screen.getByRole('button', { name: /salvar/i }))

    await waitFor(() => {
      expect(criarDespesa).toHaveBeenCalledWith(
        'fake-token',
        '1',
        expect.objectContaining({ descricao: 'Supermercado', valor_total: 150 }),
      )
    })
  })

  // ─── 5. Exibe erro de validação quando campos obrigatórios vazios ─────────
  it('exibe erro de validação quando campos obrigatórios vazios', async () => {
    vi.mocked(listarDespesasPorCartao).mockResolvedValue([])
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /nova despesa/i }))
    await user.click(screen.getByRole('button', { name: /nova despesa/i }))
    await user.click(screen.getByRole('button', { name: /salvar/i }))

    await waitFor(() => {
      expect(screen.getByText(/descrição é obrigatória/i)).toBeInTheDocument()
    })
  })

  // ─── 6. Chama excluirDespesa ao clicar em "Excluir" ──────────────────────
  it('chama excluirDespesa ao clicar em "Excluir"', async () => {
    vi.mocked(excluirDespesa).mockResolvedValue(undefined)

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByText('Supermercado'))

    await user.click(screen.getByRole('button', { name: /excluir/i }))

    await waitFor(() => {
      expect(excluirDespesa).toHaveBeenCalledWith('fake-token', 'd1')
    })
  })
})
