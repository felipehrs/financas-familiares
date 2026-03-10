import type { LoginRequest, LoginResponse, RefreshResponse } from '@/types/auth'

const API_BASE = import.meta.env.VITE_API_URL ?? ''

export async function login(data: LoginRequest): Promise<LoginResponse> {
  const response = await fetch(`${API_BASE}/api/v1/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })

  if (response.status === 401) {
    throw new Error('Credenciais inválidas')
  }

  if (!response.ok) {
    throw new Error(`Erro ao fazer login: ${response.statusText}`)
  }

  return response.json() as Promise<LoginResponse>
}

export async function refreshToken(token: string): Promise<RefreshResponse> {
  const response = await fetch(`${API_BASE}/api/v1/auth/refresh`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refresh_token: token }),
  })

  if (!response.ok) {
    throw response
  }

  return response.json() as Promise<RefreshResponse>
}
