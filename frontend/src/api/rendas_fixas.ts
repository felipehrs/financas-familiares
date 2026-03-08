import type { RendaFixa, CriarRendaFixaRequest, AtualizarRendaFixaRequest } from '@/types/renda_fixa'

const API_BASE = import.meta.env.VITE_API_URL ?? ''

async function handleResponse<T>(response: Response): Promise<T> {
  if (response.status === 401) {
    throw new Error('Não autorizado')
  }

  if (!response.ok) {
    let message = response.statusText
    try {
      const body = await response.json() as { error?: string }
      if (body.error) {
        message = body.error
      }
    } catch {
      // mantém statusText se body não for JSON válido
    }
    throw new Error(message)
  }

  if (response.status === 204) {
    return undefined as unknown as T
  }

  return response.json() as Promise<T>
}

export async function listarRendasFixas(token: string): Promise<RendaFixa[]> {
  const response = await fetch(`${API_BASE}/api/v1/rendas-fixas`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<RendaFixa[]>(response)
}

export async function criarRendaFixa(token: string, data: CriarRendaFixaRequest): Promise<RendaFixa> {
  const response = await fetch(`${API_BASE}/api/v1/rendas-fixas`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  return handleResponse<RendaFixa>(response)
}

export async function atualizarRendaFixa(
  token: string,
  id: string,
  data: AtualizarRendaFixaRequest,
): Promise<RendaFixa> {
  const response = await fetch(`${API_BASE}/api/v1/rendas-fixas/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  return handleResponse<RendaFixa>(response)
}

export async function inativarRendaFixa(token: string, id: string): Promise<void> {
  const response = await fetch(`${API_BASE}/api/v1/rendas-fixas/${id}/inativar`, {
    method: 'PATCH',
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<void>(response)
}

export async function listarRendasFixasVigentes(
  token: string,
  mes: number,
  ano: number,
): Promise<RendaFixa[]> {
  const response = await fetch(`${API_BASE}/api/v1/rendas-fixas/vigentes?mes=${mes}&ano=${ano}`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<RendaFixa[]>(response)
}
