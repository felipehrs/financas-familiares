import type { Categoria, CriarCategoriaRequest, AtualizarCategoriaRequest } from '@/types/categoria'

const API_BASE = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'

async function handleResponse<T>(response: Response): Promise<T> {
  if (response.status === 401) {
    throw new Error('Não autorizado')
  }

  if (response.status === 409) {
    throw new Error('categoria possui registros vinculados e não pode ser excluída')
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

export async function listarCategorias(token: string): Promise<Categoria[]> {
  const response = await fetch(`${API_BASE}/api/v1/categorias`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<Categoria[]>(response)
}

export async function criarCategoria(token: string, data: CriarCategoriaRequest): Promise<Categoria> {
  const response = await fetch(`${API_BASE}/api/v1/categorias`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  return handleResponse<Categoria>(response)
}

export async function atualizarCategoria(
  token: string,
  id: string,
  data: AtualizarCategoriaRequest,
): Promise<Categoria> {
  const response = await fetch(`${API_BASE}/api/v1/categorias/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  return handleResponse<Categoria>(response)
}

export async function excluirCategoria(token: string, id: string): Promise<void> {
  const response = await fetch(`${API_BASE}/api/v1/categorias/${id}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${token}` },
  })
  return handleResponse<void>(response)
}
