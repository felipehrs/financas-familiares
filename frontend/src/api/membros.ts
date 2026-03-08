import type { Membro, CriarMembroRequest, AtualizarMembroRequest } from '@/types/membro'

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

  // 204 No Content — sem corpo
  if (response.status === 204) {
    return undefined as unknown as T
  }

  return response.json() as Promise<T>
}

export async function listarMembros(token: string): Promise<Membro[]> {
  const response = await fetch(`${API_BASE}/api/v1/membros`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<Membro[]>(response)
}

export async function criarMembro(token: string, data: CriarMembroRequest): Promise<Membro> {
  const response = await fetch(`${API_BASE}/api/v1/membros`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  return handleResponse<Membro>(response)
}

export async function atualizarMembro(
  token: string,
  id: string,
  data: AtualizarMembroRequest,
): Promise<Membro> {
  const response = await fetch(`${API_BASE}/api/v1/membros/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  return handleResponse<Membro>(response)
}

export async function inativarMembro(token: string, id: string): Promise<void> {
  const response = await fetch(`${API_BASE}/api/v1/membros/${id}/inativar`, {
    method: 'PATCH',
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<void>(response)
}
