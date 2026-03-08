import type { RendaExtra, CriarRendaExtraRequest, AtualizarRendaExtraRequest } from '@/types/renda_extra'

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

export async function listarRendasExtras(
  token: string,
  mes?: number,
  ano?: number,
): Promise<RendaExtra[]> {
  const params = new URLSearchParams()
  if (mes !== undefined) params.set('mes', String(mes))
  if (ano !== undefined) params.set('ano', String(ano))
  const query = params.toString() ? `?${params.toString()}` : ''
  const response = await fetch(`${API_BASE}/api/v1/rendas-extras${query}`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<RendaExtra[]>(response)
}

export async function criarRendaExtra(
  token: string,
  data: CriarRendaExtraRequest,
): Promise<RendaExtra> {
  const response = await fetch(`${API_BASE}/api/v1/rendas-extras`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  return handleResponse<RendaExtra>(response)
}

export async function atualizarRendaExtra(
  token: string,
  id: string,
  data: AtualizarRendaExtraRequest,
): Promise<RendaExtra> {
  const response = await fetch(`${API_BASE}/api/v1/rendas-extras/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  return handleResponse<RendaExtra>(response)
}

export async function excluirRendaExtra(token: string, id: string): Promise<void> {
  const response = await fetch(`${API_BASE}/api/v1/rendas-extras/${id}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<void>(response)
}
