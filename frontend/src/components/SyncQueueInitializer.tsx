import { useSyncQueue } from '@/hooks/useSyncQueue'

interface Props {
  children: React.ReactNode
}

export function SyncQueueInitializer({ children }: Props) {
  useSyncQueue()
  return <>{children}</>
}
