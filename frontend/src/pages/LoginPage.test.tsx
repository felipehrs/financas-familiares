import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import { LoginPage } from './LoginPage'
import { AuthProvider } from '@/store/authStore'

// Mock de react-router-dom navigate
const mockNavigate = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  }
})

// Mock da API de autenticação
vi.mock('@/api/auth', () => ({
  login: vi.fn(),
  refreshToken: vi.fn(),
}))

import { login as mockLogin } from '@/api/auth'

function renderLoginPage() {
  return render(
    <MemoryRouter>
      <AuthProvider>
        <LoginPage />
      </AuthProvider>
    </MemoryRouter>
  )
}

describe('LoginPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
  })

  it('renderiza campos de email e senha', () => {
    renderLoginPage()
    expect(screen.getByLabelText(/e-mail/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/senha/i)).toBeInTheDocument()
  })

  it('renderiza botão de submit', () => {
    renderLoginPage()
    expect(screen.getByRole('button', { name: /entrar/i })).toBeInTheDocument()
  })

  it('renderiza título e subtítulo', () => {
    renderLoginPage()
    expect(screen.getByText('Finanças Familiares')).toBeInTheDocument()
    expect(screen.getByText('Acesse sua conta')).toBeInTheDocument()
  })

  it('exibe erro de validação quando email inválido é submetido', async () => {
    const user = userEvent.setup()
    renderLoginPage()

    await user.type(screen.getByLabelText(/e-mail/i), 'email-invalido')
    await user.type(screen.getByLabelText(/senha/i), 'senha123')
    await user.click(screen.getByRole('button', { name: /entrar/i }))

    await waitFor(() => {
      expect(screen.getByText(/e-mail inválido/i)).toBeInTheDocument()
    })
  })

  it('exibe erro de validação quando senha tem menos de 6 caracteres', async () => {
    const user = userEvent.setup()
    renderLoginPage()

    await user.type(screen.getByLabelText(/e-mail/i), 'user@example.com')
    await user.type(screen.getByLabelText(/senha/i), '123')
    await user.click(screen.getByRole('button', { name: /entrar/i }))

    await waitFor(() => {
      expect(screen.getByText(/mínimo 6 caracteres/i)).toBeInTheDocument()
    })
  })

  it('exibe mensagem de erro quando API retorna 401', async () => {
    vi.mocked(mockLogin).mockRejectedValueOnce(new Error('Credenciais inválidas'))

    const user = userEvent.setup()
    renderLoginPage()

    await user.type(screen.getByLabelText(/e-mail/i), 'user@example.com')
    await user.type(screen.getByLabelText(/senha/i), 'senha123')
    await user.click(screen.getByRole('button', { name: /entrar/i }))

    await waitFor(() => {
      expect(
        screen.getByText(/credenciais inválidas\. verifique seu e-mail e senha\./i)
      ).toBeInTheDocument()
    })
  })

  it('chama login e redireciona quando credenciais válidas', async () => {
    vi.mocked(mockLogin).mockResolvedValueOnce({
      access_token: 'access-token-123',
      refresh_token: 'refresh-token-456',
      token_type: 'Bearer',
      expires_in: 3600,
    })

    const user = userEvent.setup()
    renderLoginPage()

    await user.type(screen.getByLabelText(/e-mail/i), 'user@example.com')
    await user.type(screen.getByLabelText(/senha/i), 'senha123')
    await user.click(screen.getByRole('button', { name: /entrar/i }))

    await waitFor(() => {
      expect(mockLogin).toHaveBeenCalledWith({
        email: 'user@example.com',
        senha: 'senha123',
      })
      expect(mockNavigate).toHaveBeenCalledWith('/dashboard')
    })
  })

  it('exibe estado de loading durante submit', async () => {
    vi.mocked(mockLogin).mockImplementation(
      () => new Promise((resolve) => setTimeout(() => resolve({
        access_token: 'token',
        refresh_token: 'refresh',
        token_type: 'Bearer',
        expires_in: 3600,
      }), 100))
    )

    const user = userEvent.setup()
    renderLoginPage()

    await user.type(screen.getByLabelText(/e-mail/i), 'user@example.com')
    await user.type(screen.getByLabelText(/senha/i), 'senha123')
    await user.click(screen.getByRole('button', { name: /entrar/i }))

    expect(screen.getByRole('button', { name: /entrando/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /entrando/i })).toBeDisabled()

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith('/dashboard')
    })
  })
})
