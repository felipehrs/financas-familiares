import type { CartaoCredito, CriarCartaoCreditoRequest, AtualizarCartaoCreditoRequest } from '@/types/cartao_credito'

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

export async function listarCartoes(token: string): Promise<CartaoCredito[]> {
  const response = await fetch(`${API_BASE}/api/v1/cartoes`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<CartaoCredito[]>(response)
}

export async function criarCartao(token: string, data: CriarCartaoCreditoRequest): Promise<CartaoCredito> {
  const response = await fetch(`${API_BASE}/api/v1/cartoes`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  return handleResponse<CartaoCredito>(response)
}

export async function atualizarCartao(
  token: string,
  id: string,
  data: AtualizarCartaoCreditoRequest,
): Promise<CartaoCredito> {
  const response = await fetch(`${API_BASE}/api/v1/cartoes/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  return handleResponse<CartaoCredito>(response)
}

export async function inativarCartao(token: string, id: string): Promise<void> {
  const response = await fetch(`${API_BASE}/api/v1/cartoes/${id}/inativar`, {
    method: 'PATCH',
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<void>(response)
}
