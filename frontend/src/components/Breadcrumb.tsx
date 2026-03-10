import { Link, useLocation } from 'react-router-dom'

const ROUTE_LABELS: Record<string, string> = {
  '/dashboard': 'Dashboard',
  '/membros': 'Membros',
  '/categorias': 'Categorias',
  '/cartoes': 'Cartões',
  '/rendas-fixas': 'Rendas Fixas',
  '/assinaturas': 'Assinaturas',
  '/contas-fixas': 'Contas Fixas',
  '/despesas-gerais': 'Despesas Gerais',
  '/rendas-variaveis': 'Rendas Variáveis',
  '/rendas-extras': 'Rendas Extras',
  '/rendimentos-investimento': 'Rendimentos de Investimento',
  '/rendas/historico': 'Histórico de Rendas',
}

// UUID regex to detect dynamic segments
const UUID_REGEX = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i

interface BreadcrumbItem {
  label: string
  path: string | null
}

function buildBreadcrumbs(pathname: string): BreadcrumbItem[] {
  const items: BreadcrumbItem[] = [{ label: 'Dashboard', path: '/dashboard' }]

  if (pathname === '/dashboard') {
    return [{ label: 'Dashboard', path: null }]
  }

  // Handle special compound routes like /rendas/historico
  if (pathname === '/rendas/historico') {
    items.push({ label: 'Histórico de Rendas', path: null })
    return items
  }

  // Handle /cartoes/:cartaoId/despesas
  const despesasCartaoMatch = pathname.match(/^\/cartoes\/([^/]+)\/despesas$/)
  if (despesasCartaoMatch) {
    items.push({ label: 'Cartões', path: '/cartoes' })
    items.push({ label: 'Despesas', path: null })
    return items
  }

  // Generic: split path into segments and build breadcrumb
  const segments = pathname.split('/').filter(Boolean)
  let accumulatedPath = ''

  for (let i = 0; i < segments.length; i++) {
    const segment = segments[i]
    accumulatedPath += `/${segment}`
    const isLast = i === segments.length - 1

    // Skip UUID segments
    if (UUID_REGEX.test(segment)) {
      continue
    }

    const label = ROUTE_LABELS[accumulatedPath] ?? segment
    items.push({ label, path: isLast ? null : accumulatedPath })
  }

  return items
}

export function Breadcrumb() {
  const { pathname } = useLocation()
  const items = buildBreadcrumbs(pathname)

  if (items.length === 1 && items[0]?.path === null) {
    // On Dashboard itself, no breadcrumb needed
    return null
  }

  return (
    <nav aria-label="Breadcrumb" className="flex items-center gap-1 text-sm mb-4">
      {items.map((item, index) => {
        const isLast = index === items.length - 1
        return (
          <span key={index} className="flex items-center gap-1">
            {index > 0 && (
              <span className="text-muted-foreground select-none">›</span>
            )}
            {isLast || item.path === null ? (
              <span className="text-foreground font-medium">{item.label}</span>
            ) : (
              <Link
                to={item.path}
                className="text-muted-foreground hover:text-foreground transition-colors"
              >
                {item.label}
              </Link>
            )}
          </span>
        )
      })}
    </nav>
  )
}
