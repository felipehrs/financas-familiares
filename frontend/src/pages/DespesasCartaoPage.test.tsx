import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { DespesasCartaoPage } from './DespesasCartaoPage'
import type { DespesaCartao } from '@/types/despesa_cartao'
import type { CartaoCredito } from '@/types/cartao_credito'
import type { Categoria } from '@/types/categoria'

vi.mock('@/offline/despesas_cartao', () => ({
  listarDespesasPorCartao: vi.fn(),
  listarDespesasPorFatura: vi.fn(),
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

import { listarDespesasPorCartao, listarDespesasPorFatura, criarDespesa, excluirDespesa } from '@/offline/despesas_cartao'
import { listarCategorias } from '@/api/categorias'
import { listarCartoes } from '@/api/cartoes_credito'

const cartoesFixture: CartaoCredito[] = [
  { id: '1', nome: 'Nubank', membro_id: 'm-1', dia_fechamento: 10, dia_vencimento: 17, limite: 5000, ativo: true },
]

const cartaoSemLimiteFixture: CartaoCredito[] = [
  { id: '1', nome: 'Nubank', membro_id: 'm-1', dia_fechamento: 10, dia_vencimento: 17, limite: null, ativo: true },
]

const categoriasFixture: Categoria[] = [
  { id: 'cat-1', nome: 'Alimentação' },
  { id: 'cat-2', nome: 'Saúde' },
]

const despesasFixture: DespesaCartao[] = [
  {
    id: 'd1',
    compra_id: 'c1',
    cartao_id: '1',
    descricao: 'Supermercado',
    valor_total: 150,
    valor_parcela: 150,
    numero_parcelas: 1,
    parcela_numero: 1,
    fatura_mes: 3,
    fatura_ano: 2026,
    fatura: 'MAR/26',
    data_compra: '2026-03-10',
    categoria_id: null,
  },
]

// fixture para testes de total/saldo/alerta (100 + 200 = 300)
const despesasFaturaFixture: DespesaCartao[] = [
  {
    id: 'd1',
    compra_id: 'c1',
    cartao_id: '1',
    descricao: 'Supermercado',
    valor_total: 100,
    valor_parcela: 100,
    numero_parcelas: 1,
    parcela_numero: 1,
    fatura_mes: 3,
    fatura_ano: 2026,
    fatura: 'MAR/26',
    data_compra: '2026-03-10',
    categoria_id: null,
  },
  {
    id: 'd2',
    compra_id: 'c2',
    cartao_id: '1',
    descricao: 'Farmácia',
    valor_total: 200,
    valor_parcela: 200,
    numero_parcelas: 1,
    parcela_numero: 1,
    fatura_mes: 3,
    fatura_ano: 2026,
    fatura: 'MAR/26',
    data_compra: '2026-03-08',
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
    vi.mocked(listarDespesasPorFatura).mockResolvedValue([])
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
    expect(screen.getByLabelText(/número de parcelas/i)).toBeInTheDocument()
  })

  // ─── 4. Chama criarDespesa ao submeter formulário válido ─────────────────
  it('chama criarDespesa ao submeter formulário válido', async () => {
    vi.mocked(listarDespesasPorCartao).mockResolvedValue([])
    vi.mocked(criarDespesa).mockResolvedValue(despesasFixture)

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
        expect.objectContaining({ descricao: 'Supermercado', valor_total: 150, numero_parcelas: 1 }),
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

  // ─── 6. Exibe categorias no select ao abrir o formulário ─────────────────
  it('exibe categorias no select ao abrir o formulário', async () => {
    vi.mocked(listarDespesasPorCartao).mockResolvedValue([])
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /nova despesa/i }))
    await user.click(screen.getByRole('button', { name: /nova despesa/i }))

    await waitFor(() => {
      expect(screen.getByRole('option', { name: 'Alimentação' })).toBeInTheDocument()
      expect(screen.getByRole('option', { name: 'Saúde' })).toBeInTheDocument()
    })
  })

  // ─── 7. Exibe categorias mesmo quando listarDespesasPorCartao falha ───────
  it('exibe categorias no select mesmo quando o carregamento de despesas falha', async () => {
    vi.mocked(listarDespesasPorCartao).mockRejectedValue(new Error('Erro de rede'))
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /nova despesa/i }))
    await user.click(screen.getByRole('button', { name: /nova despesa/i }))

    await waitFor(() => {
      expect(screen.getByRole('option', { name: 'Alimentação' })).toBeInTheDocument()
      expect(screen.getByRole('option', { name: 'Saúde' })).toBeInTheDocument()
    })
  })

  // ─── 8. Chama excluirDespesa ao clicar em "Excluir" ──────────────────────
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

  // ─── 9. Exibe indicador "Parcela X/Y" para compras parceladas ────────────
  it('exibe indicador de parcela X/Y para compras parceladas', async () => {
    const parcelada: DespesaCartao = {
      id: 'd2',
      compra_id: 'c2',
      cartao_id: '1',
      descricao: 'Notebook',
      valor_total: 3000,
      valor_parcela: 1000,
      numero_parcelas: 3,
      parcela_numero: 1,
      fatura_mes: 3,
      fatura_ano: 2026,
      fatura: 'MAR/26',
      data_compra: '2026-03-05',
      categoria_id: null,
    }
    vi.mocked(listarDespesasPorCartao).mockResolvedValue([parcelada])

    renderPage()

    await waitFor(() => {
      expect(screen.getByText(/parcela 1\/3/i)).toBeInTheDocument()
    })
  })

  // ─── 10. Exibe erro de validação quando numero_parcelas é apagado ─────────
  it('exibe erro de validação quando numero_parcelas é menor que 1', async () => {
    vi.mocked(listarDespesasPorCartao).mockResolvedValue([])
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /nova despesa/i }))
    await user.click(screen.getByRole('button', { name: /nova despesa/i }))

    await user.clear(screen.getByLabelText(/número de parcelas/i))
    await user.click(screen.getByRole('button', { name: /salvar/i }))

    await waitFor(() => {
      expect(screen.getByText(/mínimo de 1 parcela/i)).toBeInTheDocument()
    })
  })

  // ─── 11. Seletor de mês/ano é renderizado ────────────────────────────────
  it('renderiza seletor de mês e ano para filtro de fatura', async () => {
    renderPage()

    await waitFor(() => {
      expect(screen.getByLabelText(/mês da fatura/i)).toBeInTheDocument()
    })
    expect(screen.getByLabelText(/ano da fatura/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /filtrar fatura/i })).toBeInTheDocument()
  })

  // ─── 12. Filtro chama listarDespesasPorFatura com mes/ano corretos ────────
  it('chama listarDespesasPorFatura com mês e ano selecionados ao filtrar', async () => {
    vi.mocked(listarDespesasPorFatura).mockResolvedValue(despesasFaturaFixture)
    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByLabelText(/mês da fatura/i))
    await user.selectOptions(screen.getByLabelText(/mês da fatura/i), '3')
    await user.clear(screen.getByLabelText(/ano da fatura/i))
    await user.type(screen.getByLabelText(/ano da fatura/i), '2026')
    await user.click(screen.getByRole('button', { name: /filtrar fatura/i }))

    await waitFor(() => {
      expect(listarDespesasPorFatura).toHaveBeenCalledWith('fake-token', '1', 3, 2026)
    })
  })

  // ─── 13. Total da fatura exibido quando filtro ativo ─────────────────────
  it('exibe total da fatura quando filtro de fatura está ativo', async () => {
    vi.mocked(listarDespesasPorFatura).mockResolvedValue(despesasFaturaFixture) // 100 + 200 = 300
    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /filtrar fatura/i }))
    await user.click(screen.getByRole('button', { name: /filtrar fatura/i }))

    await waitFor(() => {
      expect(screen.getByText(/total da fatura/i)).toBeInTheDocument()
      expect(screen.getByText(/300,00/)).toBeInTheDocument()
    })
  })

  // ─── 14. Saldo disponível exibido quando cartão tem limite ───────────────
  it('exibe saldo disponível quando cartão tem limite configurado', async () => {
    // cartoesFixture tem limite: 5000, total: 300 → saldo: 4.700
    vi.mocked(listarDespesasPorFatura).mockResolvedValue(despesasFaturaFixture)
    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /filtrar fatura/i }))
    await user.click(screen.getByRole('button', { name: /filtrar fatura/i }))

    await waitFor(() => {
      expect(screen.getByText(/saldo disponível/i)).toBeInTheDocument()
      expect(screen.getByText(/4\.700,00/)).toBeInTheDocument()
    })
  })

  // ─── 15. Saldo disponível NÃO exibido quando cartão sem limite ───────────
  it('não exibe saldo disponível quando cartão não tem limite', async () => {
    vi.mocked(listarCartoes).mockResolvedValue(cartaoSemLimiteFixture)
    vi.mocked(listarDespesasPorFatura).mockResolvedValue(despesasFaturaFixture)
    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /filtrar fatura/i }))
    await user.click(screen.getByRole('button', { name: /filtrar fatura/i }))

    await waitFor(() => screen.getByText(/total da fatura/i))
    expect(screen.queryByText(/saldo disponível/i)).not.toBeInTheDocument()
  })

  // ─── 16. Alerta exibido quando uso >= 80% do limite ──────────────────────
  it('exibe alerta quando uso da fatura ultrapassa 80% do limite', async () => {
    // limite 1000, despesas 900 = 90%
    const cartaoLimiteBaixo: CartaoCredito[] = [
      { id: '1', nome: 'Nubank', membro_id: 'm-1', dia_fechamento: 10, dia_vencimento: 17, limite: 1000, ativo: true },
    ]
    const despesas900: DespesaCartao[] = [
      { ...despesasFaturaFixture[0], valor_parcela: 900, valor_total: 900 },
    ]
    vi.mocked(listarCartoes).mockResolvedValue(cartaoLimiteBaixo)
    vi.mocked(listarDespesasPorFatura).mockResolvedValue(despesas900)
    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /filtrar fatura/i }))
    await user.click(screen.getByRole('button', { name: /filtrar fatura/i }))

    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument()
      expect(screen.getByText(/atenção/i)).toBeInTheDocument()
    })
  })

  // ─── 17. Sem alerta quando uso < 80% ─────────────────────────────────────
  it('não exibe alerta quando uso da fatura é menor que 80% do limite', async () => {
    // limite 5000, despesas 300 = 6%
    vi.mocked(listarDespesasPorFatura).mockResolvedValue(despesasFaturaFixture)
    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /filtrar fatura/i }))
    await user.click(screen.getByRole('button', { name: /filtrar fatura/i }))

    await waitFor(() => screen.getByText(/total da fatura/i))
    expect(screen.queryByRole('alert')).not.toBeInTheDocument()
  })

  // ─── 18. "Limpar filtro" volta a chamar listarDespesasPorCartao ───────────
  it('limpar filtro volta a exibir todas as despesas', async () => {
    vi.mocked(listarDespesasPorFatura).mockResolvedValue(despesasFaturaFixture)
    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /filtrar fatura/i }))
    await user.click(screen.getByRole('button', { name: /filtrar fatura/i }))
    await waitFor(() => screen.getByText(/total da fatura/i))

    await user.click(screen.getByRole('button', { name: /limpar filtro/i }))

    await waitFor(() => {
      // listarDespesasPorCartao foi chamado na carga inicial + após limpar filtro
      expect(listarDespesasPorCartao).toHaveBeenCalledTimes(2)
    })
    expect(screen.queryByText(/total da fatura/i)).not.toBeInTheDocument()
  })
})
