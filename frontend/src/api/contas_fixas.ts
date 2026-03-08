import type { ContaFixa, CriarContaFixaRequest, AtualizarContaFixaRequest } from '@/types/conta_fixa'

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

export async function listarContasFixas(token: string): Promise<ContaFixa[]> {
  const response = await fetch(`${API_BASE}/api/v1/contas-fixas`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<ContaFixa[]>(response)
}

export async function criarContaFixa(token: string, data: CriarContaFixaRequest): Promise<ContaFixa> {
  const response = await fetch(`${API_BASE}/api/v1/contas-fixas`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  return handleResponse<ContaFixa>(response)
}

export async function atualizarContaFixa(
  token: string,
  id: string,
  data: AtualizarContaFixaRequest,
): Promise<ContaFixa> {
  const response = await fetch(`${API_BASE}/api/v1/contas-fixas/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  return handleResponse<ContaFixa>(response)
}

export async function alterarAtivoContaFixa(token: string, id: string, ativa: boolean): Promise<ContaFixa> {
  const response = await fetch(`${API_BASE}/api/v1/contas-fixas/${id}/ativo`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ ativa }),
  })
  return handleResponse<ContaFixa>(response)
}
