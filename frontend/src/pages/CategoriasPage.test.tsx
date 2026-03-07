import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import { CategoriasPage } from './CategoriasPage'
import type { Categoria } from '@/types/categoria'

// Mock do módulo de API
vi.mock('@/api/categorias', () => ({
  listarCategorias: vi.fn(),
  criarCategoria: vi.fn(),
  atualizarCategoria: vi.fn(),
  excluirCategoria: vi.fn(),
}))

// Mock do useAuth
vi.mock('@/hooks/useAuth', () => ({
  useAuth: vi.fn(() => ({ accessToken: 'fake-token' })),
}))

import {
  listarCategorias,
  criarCategoria,
  atualizarCategoria,
  excluirCategoria,
} from '@/api/categorias'

const categoriasFixture: Categoria[] = [
  { id: '1', nome: 'Alimentação' },
  { id: '2', nome: 'Transporte' },
]

function renderPage() {
  return render(
    <MemoryRouter>
      <CategoriasPage />
    </MemoryRouter>
  )
}

describe('CategoriasPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  // ─── 1. Renderiza lista de categorias após carregar ───────────────────────
  it('renderiza lista de categorias após carregar', async () => {
    vi.mocked(listarCategorias).mockResolvedValueOnce(categoriasFixture)

    renderPage()

    await waitFor(() => {
      expect(screen.getByText('Alimentação')).toBeInTheDocument()
      expect(screen.getByText('Transporte')).toBeInTheDocument()
    })
  })

  // ─── 2. Exibe mensagem quando lista está vazia ────────────────────────────
  it('exibe "Nenhuma categoria cadastrada" quando lista vazia', async () => {
    vi.mocked(listarCategorias).mockResolvedValueOnce([])

    renderPage()

    await waitFor(() => {
      expect(screen.getByText(/nenhuma categoria cadastrada/i)).toBeInTheDocument()
    })
  })

  // ─── 3. Exibe formulário ao clicar em "Adicionar categoria" ──────────────
  it('exibe formulário ao clicar em "Adicionar categoria"', async () => {
    vi.mocked(listarCategorias).mockResolvedValueOnce([])
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /adicionar categoria/i }))
    await user.click(screen.getByRole('button', { name: /adicionar categoria/i }))

    expect(screen.getByLabelText(/nome/i)).toBeInTheDocument()
  })

  // ─── 4. Chama criarCategoria ao submeter formulário válido ───────────────
  it('chama criarCategoria ao submeter formulário válido e atualiza a lista', async () => {
    vi.mocked(listarCategorias)
      .mockResolvedValueOnce([])
      .mockResolvedValueOnce([{ id: '3', nome: 'Lazer' }])
    vi.mocked(criarCategoria).mockResolvedValueOnce({ id: '3', nome: 'Lazer' })

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /adicionar categoria/i }))
    await user.click(screen.getByRole('button', { name: /adicionar categoria/i }))

    await user.type(screen.getByLabelText(/nome/i), 'Lazer')
    await user.click(screen.getByRole('button', { name: /salvar/i }))

    await waitFor(() => {
      expect(criarCategoria).toHaveBeenCalledWith('fake-token', { nome: 'Lazer' })
    })

    await waitFor(() => {
      expect(screen.getByText('Lazer')).toBeInTheDocument()
    })
  })

  // ─── 5. Exibe erro de validação quando nome está vazio ────────────────────
  it('exibe erro de validação quando nome está vazio', async () => {
    vi.mocked(listarCategorias).mockResolvedValueOnce([])
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /adicionar categoria/i }))
    await user.click(screen.getByRole('button', { name: /adicionar categoria/i }))
    await user.click(screen.getByRole('button', { name: /salvar/i }))

    await waitFor(() => {
      expect(screen.getByText(/nome é obrigatório/i)).toBeInTheDocument()
    })
  })

  // ─── 6. Exibe erro de API quando criarCategoria falha ────────────────────
  it('exibe erro de API quando criarCategoria falha', async () => {
    vi.mocked(listarCategorias).mockResolvedValueOnce([])
    vi.mocked(criarCategoria).mockRejectedValueOnce(new Error('Erro interno do servidor'))

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByRole('button', { name: /adicionar categoria/i }))
    await user.click(screen.getByRole('button', { name: /adicionar categoria/i }))

    await user.type(screen.getByLabelText(/nome/i), 'Lazer')
    await user.click(screen.getByRole('button', { name: /salvar/i }))

    await waitFor(() => {
      expect(screen.getByText(/erro interno do servidor/i)).toBeInTheDocument()
    })
  })

  // ─── 7. Chama window.confirm antes de excluir ────────────────────────────
  it('chama window.confirm antes de excluir', async () => {
    vi.mocked(listarCategorias)
      .mockResolvedValueOnce(categoriasFixture)
      .mockResolvedValueOnce(categoriasFixture)
    vi.mocked(excluirCategoria).mockResolvedValueOnce(undefined)
    const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(true)

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByText('Alimentação'))

    const items = screen.getAllByRole('listitem')
    const alimentacaoItem = items.find((el) => el.textContent?.includes('Alimentação'))!
    await user.click(within(alimentacaoItem).getByRole('button', { name: /excluir/i }))

    expect(confirmSpy).toHaveBeenCalledWith("Deseja excluir a categoria 'Alimentação'?")

    confirmSpy.mockRestore()
  })

  // ─── 8. Chama excluirCategoria após confirmação e remove da lista ─────────
  it('chama excluirCategoria após confirmação e remove da lista', async () => {
    vi.mocked(listarCategorias)
      .mockResolvedValueOnce(categoriasFixture)
      .mockResolvedValueOnce([categoriasFixture[1]])
    vi.mocked(excluirCategoria).mockResolvedValueOnce(undefined)
    vi.spyOn(window, 'confirm').mockReturnValue(true)

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByText('Alimentação'))

    const items = screen.getAllByRole('listitem')
    const alimentacaoItem = items.find((el) => el.textContent?.includes('Alimentação'))!
    await user.click(within(alimentacaoItem).getByRole('button', { name: /excluir/i }))

    await waitFor(() => {
      expect(excluirCategoria).toHaveBeenCalledWith('fake-token', '1')
    })

    await waitFor(() => {
      expect(screen.queryByText('Alimentação')).not.toBeInTheDocument()
      expect(screen.getByText('Transporte')).toBeInTheDocument()
    })

    vi.restoreAllMocks()
  })

  // ─── 9. Exibe erro 409 quando excluirCategoria falha com vínculo ──────────
  it('exibe erro 409 quando excluirCategoria falha com vínculo', async () => {
    vi.mocked(listarCategorias).mockResolvedValueOnce(categoriasFixture)
    vi.mocked(excluirCategoria).mockRejectedValueOnce(
      new Error('categoria possui registros vinculados e não pode ser excluída'),
    )
    vi.spyOn(window, 'confirm').mockReturnValue(true)

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByText('Alimentação'))

    const items = screen.getAllByRole('listitem')
    const alimentacaoItem = items.find((el) => el.textContent?.includes('Alimentação'))!
    await user.click(within(alimentacaoItem).getByRole('button', { name: /excluir/i }))

    await waitFor(() => {
      expect(
        screen.getByText(/categoria possui registros vinculados e não pode ser excluída/i),
      ).toBeInTheDocument()
    })

    vi.restoreAllMocks()
  })

  // ─── 10. Abre formulário de edição com nome preenchido ───────────────────
  it('abre formulário de edição ao clicar em "Editar" com nome preenchido', async () => {
    vi.mocked(listarCategorias).mockResolvedValueOnce(categoriasFixture)
    const user = userEvent.setup()

    renderPage()

    await waitFor(() => screen.getByText('Alimentação'))

    const items = screen.getAllByRole('listitem')
    const alimentacaoItem = items.find((el) => el.textContent?.includes('Alimentação'))!
    await user.click(within(alimentacaoItem).getByRole('button', { name: /editar/i }))

    await waitFor(() => {
      const nomeInput = screen.getByLabelText(/nome/i) as HTMLInputElement
      expect(nomeInput.value).toBe('Alimentação')
    })
  })

  // ─── 11. Chama atualizarCategoria ao submeter edição válida ──────────────
  it('chama atualizarCategoria ao submeter edição válida', async () => {
    vi.mocked(listarCategorias)
      .mockResolvedValueOnce(categoriasFixture)
      .mockResolvedValueOnce([
        { id: '1', nome: 'Alimentação e Bebidas' },
        categoriasFixture[1],
      ])
    vi.mocked(atualizarCategoria).mockResolvedValueOnce({
      id: '1',
      nome: 'Alimentação e Bebidas',
    })

    const user = userEvent.setup()
    renderPage()

    await waitFor(() => screen.getByText('Alimentação'))

    const items = screen.getAllByRole('listitem')
    const alimentacaoItem = items.find((el) => el.textContent?.includes('Alimentação'))!
    await user.click(within(alimentacaoItem).getByRole('button', { name: /editar/i }))

    const nomeInput = screen.getByLabelText(/nome/i)
    await user.clear(nomeInput)
    await user.type(nomeInput, 'Alimentação e Bebidas')
    await user.click(screen.getByRole('button', { name: /salvar/i }))

    await waitFor(() => {
      expect(atualizarCategoria).toHaveBeenCalledWith('fake-token', '1', {
        nome: 'Alimentação e Bebidas',
      })
    })
  })
})
