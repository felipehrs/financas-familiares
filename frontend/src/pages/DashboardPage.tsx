import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import { useTheme } from '@/hooks/useTheme'
import { buscarResumoMensal, buscarCategoriasDespesas, buscarEvolucaoMensal, buscarProjecao } from '@/api/dashboard'
import type { ResumoMensal, ResumoCategorias, PontoEvolucao, MesProjecao } from '@/types/dashboard'
import { Card, CardContent } from '@/components/ui/card'
import { GraficoCategoriaDespesas } from '@/components/GraficoCategoriaDespesas'
import { GraficoEvolucaoMensal } from '@/components/GraficoEvolucaoMensal'
import { TabelaProjecao } from '@/components/TabelaProjecao'
import {
  Users, Tag, CreditCard, TrendingUp, RefreshCcw, Building2,
  ShoppingBag, BarChart2, Gift, PiggyBank, Moon, Sun, History,
} from 'lucide-react'

const formatarMoeda = (valor: number) =>
  new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(valor)

export function DashboardPage() {
  const { accessToken, logout } = useAuth()
  const { theme, toggleTheme } = useTheme()

  const hoje = new Date()
  const [mes, setMes] = useState(hoje.getMonth() + 1)
  const [ano, setAno] = useState(hoje.getFullYear())
  const [resumo, setResumo] = useState<ResumoMensal | null>(null)
  const [loading, setLoading] = useState(true)
  const [erro, setErro] = useState<string | null>(null)
  const [resumoCategorias, setResumoCategorias] = useState<ResumoCategorias | null>(null)
  const [loadingCategorias, setLoadingCategorias] = useState(true)
  const [erroCategorias, setErroCategorias] = useState<string | null>(null)
  const [evolucao, setEvolucao] = useState<PontoEvolucao[]>([])
  const [loadingEvolucao, setLoadingEvolucao] = useState(true)
  const [erroEvolucao, setErroEvolucao] = useState<string | null>(null)
  const [projecao, setProjecao] = useState<MesProjecao[]>([])
  const [loadingProjecao, setLoadingProjecao] = useState(true)
  const [erroProjecao, setErroProjecao] = useState<string | null>(null)

  async function carregarResumo(mesSelecionado: number, anoSelecionado: number) {
    if (!accessToken) return
    setLoading(true)
    setErro(null)
    try {
      const dados = await buscarResumoMensal(accessToken, mesSelecionado, anoSelecionado)
      setResumo(dados)
    } catch (err) {
      setErro(err instanceof Error ? err.message : 'Erro ao carregar resumo')
    } finally {
      setLoading(false)
    }
  }

  async function carregarCategorias(m: number, a: number) {
    setLoadingCategorias(true)
    setErroCategorias(null)
    try {
      const data = await buscarCategoriasDespesas(accessToken!, m, a)
      setResumoCategorias(data)
    } catch (err) {
      setErroCategorias(err instanceof Error ? err.message : 'Erro ao carregar categorias')
    } finally {
      setLoadingCategorias(false)
    }
  }

  async function carregarEvolucao() {
    if (!accessToken) return
    setLoadingEvolucao(true)
    setErroEvolucao(null)
    try {
      const data = await buscarEvolucaoMensal(accessToken)
      setEvolucao(data)
    } catch (err) {
      setErroEvolucao(err instanceof Error ? err.message : 'Erro ao carregar evolução')
    } finally {
      setLoadingEvolucao(false)
    }
  }

  async function carregarProjecao() {
    if (!accessToken) return
    setLoadingProjecao(true)
    setErroProjecao(null)
    try {
      const data = await buscarProjecao(accessToken)
      setProjecao(data)
    } catch (err) {
      setErroProjecao(err instanceof Error ? err.message : 'Erro ao carregar projeção')
    } finally {
      setLoadingProjecao(false)
    }
  }

  useEffect(() => {
    void carregarResumo(mes, ano)
    void carregarCategorias(mes, ano)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [mes, ano])

  useEffect(() => {
    void carregarEvolucao()
    void carregarProjecao()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const meses = [
    { value: 1, label: 'Janeiro' },
    { value: 2, label: 'Fevereiro' },
    { value: 3, label: 'Março' },
    { value: 4, label: 'Abril' },
    { value: 5, label: 'Maio' },
    { value: 6, label: 'Junho' },
    { value: 7, label: 'Julho' },
    { value: 8, label: 'Agosto' },
    { value: 9, label: 'Setembro' },
    { value: 10, label: 'Outubro' },
    { value: 11, label: 'Novembro' },
    { value: 12, label: 'Dezembro' },
  ]

  const anosDisponiveis = Array.from({ length: 5 }, (_, i) => hoje.getFullYear() - 2 + i)

  return (
    <div className="p-6 max-w-4xl mx-auto">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Dashboard</h1>
        <div className="flex items-center gap-2">
          <button
            onClick={toggleTheme}
            aria-label="Alternar tema"
            className="rounded-md p-2 text-muted-foreground hover:bg-accent hover:text-foreground transition-colors"
          >
            {theme === 'dark' ? <Sun size={18} /> : <Moon size={18} />}
          </button>
          <button onClick={logout} className="text-sm text-red-600 dark:text-red-400 hover:underline">
            Sair
          </button>
        </div>
      </div>

      <nav className="mb-6">
        <p className="text-xs font-medium text-muted-foreground uppercase tracking-wide mb-3">Lançamentos</p>
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-2 mb-4">
          {[
            { to: '/cartoes', icon: CreditCard, label: 'Despesas Cartão', color: 'text-purple-600 dark:text-purple-400 bg-purple-50 dark:bg-purple-950/30' },
            { to: '/despesas-gerais', icon: ShoppingBag, label: 'Despesas Gerais', color: 'text-orange-600 dark:text-orange-400 bg-orange-50 dark:bg-orange-950/30' },
            { to: '/rendas-variaveis', icon: BarChart2, label: 'Rendas Variáveis', color: 'text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-950/30' },
            { to: '/rendas-extras', icon: Gift, label: 'Rendas Extras', color: 'text-pink-600 dark:text-pink-400 bg-pink-50 dark:bg-pink-950/30' },
            { to: '/rendas/historico', icon: History, label: 'Histórico de Rendas', color: 'text-teal-600 dark:text-teal-400 bg-teal-50 dark:bg-teal-950/30' },
          ].map(({ to, icon: Icon, label, color }) => (
            <Link
              key={to}
              to={to}
              className="flex items-center gap-2 rounded-lg border bg-card px-3 py-2.5 text-sm font-medium hover:bg-accent transition-colors"
            >
              <span className={`rounded-md p-1.5 ${color}`}>
                <Icon size={14} />
              </span>
              {label}
            </Link>
          ))}
        </div>

        <p className="text-xs font-medium text-muted-foreground uppercase tracking-wide mb-3">Cadastros</p>
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-2">
          {[
            { to: '/membros', icon: Users, label: 'Membros' },
            { to: '/categorias', icon: Tag, label: 'Categorias' },
            { to: '/cartoes', icon: CreditCard, label: 'Cartões' },
            { to: '/rendas-fixas', icon: TrendingUp, label: 'Rendas Fixas' },
            { to: '/assinaturas', icon: RefreshCcw, label: 'Assinaturas' },
            { to: '/contas-fixas', icon: Building2, label: 'Contas Fixas' },
            { to: '/rendimentos-investimento', icon: PiggyBank, label: 'Investimentos' },
          ].map(({ to, icon: Icon, label }) => (
            <Link
              key={to}
              to={to}
              className="flex items-center gap-2 rounded-lg border bg-card px-3 py-2 text-sm text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
            >
              <Icon size={14} />
              {label}
            </Link>
          ))}
        </div>
      </nav>

      {/* Seletor de mês/ano */}
      <div className="mb-6 flex gap-3 items-center">
        <select
          aria-label="Mês"
          value={mes}
          onChange={(e) => setMes(Number(e.target.value))}
          className="rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
        >
          {meses.map((m) => (
            <option key={m.value} value={m.value}>
              {m.label}
            </option>
          ))}
        </select>

        <select
          aria-label="Ano"
          value={ano}
          onChange={(e) => setAno(Number(e.target.value))}
          className="rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
        >
          {anosDisponiveis.map((a) => (
            <option key={a} value={a}>
              {a}
            </option>
          ))}
        </select>
      </div>

      {/* Conteúdo */}
      {loading ? (
        <p className="text-muted-foreground">Carregando...</p>
      ) : erro ? (
        <p className="text-red-600 dark:text-red-400" role="alert">
          {erro}
        </p>
      ) : resumo ? (
        <div className="space-y-4">
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {/* Seção 1 — Rendas Operacionais */}
            <Card>
              <CardContent className="pt-4">
                <p className="text-sm font-medium mb-1">Rendas Operacionais</p>
                <p className="text-xl font-semibold text-green-600 dark:text-green-400 mb-3">
                  {formatarMoeda(resumo.total_rendas_operacionais)}
                </p>
                <div className="space-y-1 text-sm">
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Renda Fixa</span>
                    <span>{formatarMoeda(resumo.total_renda_fixa)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Renda Variável</span>
                    <span>{formatarMoeda(resumo.total_renda_variavel)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Renda Extra</span>
                    <span>{formatarMoeda(resumo.total_renda_extra)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Rendimento Distribuído</span>
                    <span>{formatarMoeda(resumo.total_rendimento_distribuido)}</span>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Seção 2 — Rendimentos de Investimento */}
            <Card>
              <CardContent className="pt-4">
                <p className="text-sm font-medium mb-0.5">Rendimentos de Investimento</p>
                <p className="text-xs text-muted-foreground mb-1">(apenas informativo)</p>
                <p className="text-xl font-semibold mb-3">
                  {formatarMoeda(resumo.total_rendimento_investimento)}
                </p>
                <div className="space-y-1 text-sm">
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Valor distribuído ao orçamento</span>
                    <span>{formatarMoeda(resumo.total_rendimento_distribuido)}</span>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Seção 3 — Despesas */}
            <Card>
              <CardContent className="pt-4">
                <p className="text-sm font-medium mb-1">Despesas</p>
                <p className="text-xl font-semibold text-red-600 dark:text-red-400 mb-3">
                  {formatarMoeda(resumo.total_despesas)}
                </p>
                <div className="space-y-1 text-sm">
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Fatura Cartões</span>
                    <span>{formatarMoeda(resumo.total_fatura_cartoes)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Assinaturas</span>
                    <span>{formatarMoeda(resumo.total_assinaturas)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Contas Fixas</span>
                    <span>{formatarMoeda(resumo.total_contas_fixas)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Despesas Gerais</span>
                    <span>{formatarMoeda(resumo.total_despesas_gerais)}</span>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Seção 4 — Saldo */}
          <Card>
            <CardContent className="pt-4">
              <p className="text-sm font-medium mb-1">Saldo do Mês</p>
              <p
                className={`text-2xl font-semibold ${resumo.saldo >= 0 ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'}`}
              >
                {formatarMoeda(resumo.saldo)}
              </p>
            </CardContent>
          </Card>

          {/* Seção — Despesas por Categoria */}
          <Card>
            <CardContent className="pt-4">
              <p className="font-medium mb-4">Despesas por Categoria</p>
              {loadingCategorias ? (
                <p className="text-muted-foreground text-sm">Carregando...</p>
              ) : erroCategorias ? (
                <p className="text-red-600 dark:text-red-400 text-sm" role="alert">{erroCategorias}</p>
              ) : resumoCategorias ? (
                <GraficoCategoriaDespesas categorias={resumoCategorias.categorias} />
              ) : null}
            </CardContent>
          </Card>

          {/* Seção — Evolução Mensal */}
          <Card>
            <CardContent className="pt-4">
              <p className="font-medium mb-4">Evolução Mensal</p>
              {loadingEvolucao ? (
                <p className="text-muted-foreground text-sm">Carregando...</p>
              ) : erroEvolucao ? (
                <p className="text-red-600 dark:text-red-400 text-sm" role="alert">{erroEvolucao}</p>
              ) : (
                <GraficoEvolucaoMensal pontos={evolucao} />
              )}
            </CardContent>
          </Card>

          {/* Seção — Projeção dos Próximos 3 Meses */}
          <Card>
            <CardContent className="pt-4">
              <p className="font-medium mb-4">Projeção dos Próximos 3 Meses</p>
              {loadingProjecao ? (
                <p className="text-muted-foreground text-sm">Carregando...</p>
              ) : erroProjecao ? (
                <p className="text-red-600 dark:text-red-400 text-sm" role="alert">{erroProjecao}</p>
              ) : (
                <TabelaProjecao projecao={projecao} />
              )}
            </CardContent>
          </Card>
        </div>
      ) : null}
    </div>
  )
}
