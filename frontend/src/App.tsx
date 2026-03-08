import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from '@/store/authStore'
import { LoginPage } from '@/pages/LoginPage'
import { DashboardPage } from '@/pages/DashboardPage'
import { MembrosPage } from '@/pages/MembrosPage'
import { CategoriasPage } from '@/pages/CategoriasPage'
import { CartoesPage } from '@/pages/CartoesPage'
import { DespesasCartaoPage } from '@/pages/DespesasCartaoPage'
import { RendasFixasPage } from '@/pages/RendasFixasPage'
import { AssinaturasPage } from '@/pages/AssinaturasPage'
import { ProtectedRoute } from '@/components/ProtectedRoute'

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
          <Route path="/login" element={<LoginPage />} />
          <Route element={<ProtectedRoute />}>
            <Route path="/dashboard" element={<DashboardPage />} />
            <Route path="/membros" element={<MembrosPage />} />
            <Route path="/categorias" element={<CategoriasPage />} />
            <Route path="/cartoes" element={<CartoesPage />} />
            <Route path="/cartoes/:cartaoId/despesas" element={<DespesasCartaoPage />} />
            <Route path="/rendas-fixas" element={<RendasFixasPage />} />
            <Route path="/assinaturas" element={<AssinaturasPage />} />
          </Route>
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  )
}

export default App
