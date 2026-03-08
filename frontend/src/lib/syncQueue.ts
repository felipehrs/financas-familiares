import { db } from './db'

export async function processarFila(token: string): Promise<void> {
  const pendentes = await db.sync_queue.orderBy('created_at').toArray()

  for (const item of pendentes) {
    try {
      const response = await fetch(item.endpoint, {
        method: item.method,
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: item.method !== 'DELETE' ? JSON.stringify(item.body) : undefined,
      })

      if (response.ok || response.status === 409 || response.status === 404) {
        await db.sync_queue.delete(item.id!)
      } else {
        await db.sync_queue.update(item.id!, { retries: item.retries + 1 })
      }
    } catch {
      await db.sync_queue.update(item.id!, { retries: item.retries + 1 })
    }
  }
}

export function configurarListenerOnline(getToken: () => string | null): void {
  window.addEventListener('online', () => {
    const token = getToken()
    if (token) {
      void processarFila(token)
    }
  })
}
