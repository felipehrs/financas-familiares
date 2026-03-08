import type { Assinatura, CriarAssinaturaRequest, AtualizarAssinaturaRequest } from '@/types/assinatura'

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

export async function listarAssinaturas(token: string): Promise<Assinatura[]> {
  const response = await fetch(`${API_BASE}/api/v1/assinaturas`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<Assinatura[]>(response)
}

export async function criarAssinatura(token: string, data: CriarAssinaturaRequest): Promise<Assinatura> {
  const response = await fetch(`${API_BASE}/api/v1/assinaturas`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  return handleResponse<Assinatura>(response)
}

export async function atualizarAssinatura(
  token: string,
  id: string,
  data: AtualizarAssinaturaRequest,
): Promise<Assinatura> {
  const response = await fetch(`${API_BASE}/api/v1/assinaturas/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  return handleResponse<Assinatura>(response)
}

export async function alterarStatusAssinatura(token: string, id: string, status: string): Promise<Assinatura> {
  const response = await fetch(`${API_BASE}/api/v1/assinaturas/${id}/status`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ status }),
  })
  return handleResponse<Assinatura>(response)
}
