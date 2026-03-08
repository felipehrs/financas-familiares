import type { RendaVariavel, CriarRendaVariavelRequest, AtualizarRendaVariavelRequest } from '@/types/renda_variavel'

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

export async function listarRendasVariaveis(
  token: string,
  mes?: number,
  ano?: number,
): Promise<RendaVariavel[]> {
  const params = new URLSearchParams()
  if (mes !== undefined) params.set('mes', String(mes))
  if (ano !== undefined) params.set('ano', String(ano))
  const query = params.toString() ? `?${params.toString()}` : ''
  const response = await fetch(`${API_BASE}/api/v1/rendas-variaveis${query}`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<RendaVariavel[]>(response)
}

export async function criarRendaVariavel(
  token: string,
  data: CriarRendaVariavelRequest,
): Promise<RendaVariavel> {
  const response = await fetch(`${API_BASE}/api/v1/rendas-variaveis`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  return handleResponse<RendaVariavel>(response)
}

export async function atualizarRendaVariavel(
  token: string,
  id: string,
  data: AtualizarRendaVariavelRequest,
): Promise<RendaVariavel> {
  const response = await fetch(`${API_BASE}/api/v1/rendas-variaveis/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  return handleResponse<RendaVariavel>(response)
}

export async function excluirRendaVariavel(token: string, id: string): Promise<void> {
  const response = await fetch(`${API_BASE}/api/v1/rendas-variaveis/${id}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<void>(response)
}
