import { Link } from 'react-router-dom'
import { ChevronLeft } from 'lucide-react'

export function BackToDashboard() {
  return (
    <Link
      to="/dashboard"
      className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground transition-colors mb-6"
    >
      <ChevronLeft size={16} />
      Dashboard
    </Link>
  )
}
