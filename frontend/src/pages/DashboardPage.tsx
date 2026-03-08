import { Link } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'

export function DashboardPage() {
  const { logout } = useAuth()
  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold">Dashboard</h1>
      <p>Em construção...</p>
      <nav className="mt-4 flex gap-4">
        <Link to="/membros" className="text-blue-600 hover:underline">
          Membros da Família
        </Link>
        <Link to="/categorias" className="text-blue-600 hover:underline">
          Categorias
        </Link>
        <Link to="/cartoes" className="text-blue-600 hover:underline">
          Cartões de Crédito
        </Link>
      </nav>
      <button className="mt-4" onClick={logout}>Sair</button>
    </div>
  )
}
