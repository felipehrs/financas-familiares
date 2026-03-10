import { useState } from "react"
import { Link, Outlet, useLocation } from "react-router-dom"
import {
  CreditCard,
  Receipt,
  TrendingUp,
  Plus,
  History,
  Users,
  Tag,
  Wallet,
  RefreshCw,
  FileText,
  BarChart2,
  Sun,
  Moon,
  LogOut,
  Menu,
  X,
} from "lucide-react"
import { useAuth } from "@/hooks/useAuth"
import { useTheme } from "@/hooks/useTheme"
import { SyncQueueInitializer } from "@/components/SyncQueueInitializer"
import { Breadcrumb } from "@/components/Breadcrumb"

interface NavItem {
  to: string
  label: string
  icon: React.ElementType
}

const LANCAMENTOS: NavItem[] = [
  { to: "/cartoes", label: "Despesas de Cartão", icon: CreditCard },
  { to: "/despesas-gerais", label: "Despesas Gerais", icon: Receipt },
  { to: "/rendas-variaveis", label: "Rendas Variáveis", icon: TrendingUp },
  { to: "/rendas-extras", label: "Rendas Extras", icon: Plus },
  { to: "/rendas/historico", label: "Histórico de Rendas", icon: History },
]

const CADASTROS: NavItem[] = [
  { to: "/membros", label: "Membros", icon: Users },
  { to: "/categorias", label: "Categorias", icon: Tag },
  { to: "/cartoes", label: "Cartões", icon: CreditCard },
  { to: "/rendas-fixas", label: "Rendas Fixas", icon: Wallet },
  { to: "/assinaturas", label: "Assinaturas", icon: RefreshCw },
  { to: "/contas-fixas", label: "Contas Fixas", icon: FileText },
  { to: "/rendimentos-investimento", label: "Rendimentos de Investimento", icon: BarChart2 },
]

interface NavSectionProps {
  title: string
  items: NavItem[]
  currentPath: string
  onItemClick?: () => void
}

function NavSection({ title, items, currentPath, onItemClick }: NavSectionProps) {
  return (
    <div className="mb-4">
      <p className="px-3 text-xs font-medium text-muted-foreground uppercase tracking-wide mb-1">
        {title}
      </p>
      <ul>
        {items.map((item) => {
          const Icon = item.icon
          const isActive = currentPath === item.to
          return (
            <li key={`${title}-${item.to}-${item.label}`}>
              <Link
                to={item.to}
                onClick={onItemClick}
                className={`flex items-center gap-2 px-3 py-2 rounded-md text-sm transition-colors ${
                  isActive
                    ? "bg-accent text-accent-foreground font-medium"
                    : "text-muted-foreground hover:text-foreground hover:bg-accent"
                }`}
              >
                <Icon size={16} />
                {item.label}
              </Link>
            </li>
          )
        })}
      </ul>
    </div>
  )
}

interface SidebarContentProps {
  currentPath: string
  onItemClick?: () => void
}

function SidebarContent({ currentPath, onItemClick }: SidebarContentProps) {
  const { logout } = useAuth()
  const { theme, toggleTheme } = useTheme()

  return (
    <div className="flex flex-col h-full">
      <div className="px-3 py-4 border-b border-border mb-4">
        <Link
          to="/dashboard"
          onClick={onItemClick}
          className="text-base font-semibold text-foreground hover:text-foreground/80 transition-colors"
        >
          Finanças Familiares
        </Link>
      </div>

      <nav className="flex-1 overflow-y-auto px-2">
        <NavSection
          title="Lançamentos"
          items={LANCAMENTOS}
          currentPath={currentPath}
          onItemClick={onItemClick}
        />
        <div className="border-t border-border my-2" />
        <NavSection
          title="Cadastros"
          items={CADASTROS}
          currentPath={currentPath}
          onItemClick={onItemClick}
        />
      </nav>

      <div className="border-t border-border px-2 py-3 space-y-1">
        <button
          onClick={toggleTheme}
          aria-label="Alternar tema"
          className="flex items-center gap-2 w-full px-3 py-2 rounded-md text-sm text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
        >
          {theme === "dark" ? <Sun size={16} /> : <Moon size={16} />}
          {theme === "dark" ? "Tema claro" : "Tema escuro"}
        </button>
        <button
          onClick={logout}
          className="flex items-center gap-2 w-full px-3 py-2 rounded-md text-sm text-muted-foreground hover:text-red-600 dark:hover:text-red-400 hover:bg-accent transition-colors"
        >
          <LogOut size={16} />
          Sair
        </button>
      </div>
    </div>
  )
}

export function AppLayout() {
  const { pathname } = useLocation()
  const [drawerOpen, setDrawerOpen] = useState(false)
  const { theme, toggleTheme } = useTheme()

  function closeDrawer() {
    setDrawerOpen(false)
  }

  return (
    <SyncQueueInitializer>
      <div className="flex h-screen bg-background">
        <aside className="hidden md:flex md:flex-col w-64 border-r border-border bg-background flex-shrink-0">
          <SidebarContent currentPath={pathname} />
        </aside>

        {drawerOpen && (
          <div
            className="fixed inset-0 z-40 bg-black/50 md:hidden"
            onClick={closeDrawer}
            aria-hidden="true"
          />
        )}

        <div
          className={`fixed inset-y-0 left-0 z-50 w-72 bg-background border-r border-border flex flex-col transform transition-transform duration-200 ease-in-out md:hidden ${
            drawerOpen ? "translate-x-0" : "-translate-x-full"
          }`}
        >
          <div className="flex items-center justify-between px-4 py-3 border-b border-border">
            <span className="font-semibold text-foreground">Finanças Familiares</span>
            <button
              onClick={closeDrawer}
              aria-label="Fechar menu"
              className="p-1 rounded-md text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
            >
              <X size={18} />
            </button>
          </div>
          <div className="flex-1 overflow-y-auto py-2">
            <SidebarContent currentPath={pathname} onItemClick={closeDrawer} />
          </div>
        </div>

        <div className="flex flex-col flex-1 min-w-0 overflow-hidden">
          <header className="md:hidden flex items-center justify-between px-4 py-3 border-b border-border bg-background flex-shrink-0">
            <button
              onClick={() => setDrawerOpen(true)}
              aria-label="Abrir menu"
              className="p-1 rounded-md text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
            >
              <Menu size={20} />
            </button>
            <span className="font-semibold text-foreground text-sm">Finanças Familiares</span>
            <button
              onClick={toggleTheme}
              aria-label="Alternar tema"
              className="p-1 rounded-md text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
            >
              {theme === "dark" ? <Sun size={18} /> : <Moon size={18} />}
            </button>
          </header>

          <main className="flex-1 overflow-y-auto">
            <div className="px-6 pt-4">
              <Breadcrumb />
            </div>
            <Outlet />
          </main>
        </div>
      </div>
    </SyncQueueInitializer>
  )
}
