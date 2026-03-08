import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import { buscarResumoMensal } from '@/api/dashboard'
import type { ResumoMensal } from '@/types/dashboard'
import { Card, CardContent } from '@/components/ui/card'

const formatarMoeda = (valor: number) =>
  new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(valor)

export function DashboardPage() {
  const { accessToken, logout } = useAuth()

  const hoje = new Date()
  const [mes, setMes] = useState(hoje.getMonth() + 1)
  const [ano, setAno] = useState(hoje.getFullYear())
  const [resumo, setResumo] = useState<ResumoMensal | null>(null)
  const [loading, setLoading] = useState(true)
  const [erro, setErro] = useState<string | null>(null)

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

  useEffect(() => {
    void carregarResumo(mes, ano)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [mes, ano])

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
        <button onClick={logout} className="text-sm text-red-600 hover:underline">
          Sair
        </button>
      </div>

      <nav className="mb-6 flex gap-4">
        <Link to="/membros" className="text-blue-600 hover:underline">
          Membros da Família
        </Link>
        <Link to="/categorias" className="text-blue-600 hover:underline">
          Categorias
        </Link>
        <Link to="/cartoes" className="text-blue-600 hover:underline">
          Cartões de Crédito
        </Link>
        <Link to="/rendas-fixas" className="text-blue-600 hover:underline">
          Rendas Fixas
        </Link>
        <Link to="/assinaturas" className="text-blue-600 hover:underline">
          Assinaturas
        </Link>
        <Link to="/contas-fixas" className="text-blue-600 hover:underline">
          Contas Fixas
        </Link>
        <Link to="/despesas-gerais" className="text-blue-600 hover:underline">
          Despesas Gerais
        </Link>
        <Link to="/rendas-variaveis" className="text-blue-600 hover:underline">
          Rendas Variáveis
        </Link>
        <Link to="/rendas-extras" className="text-blue-600 hover:underline">
          Rendas Extras
        </Link>
        <Link to="/rendimentos-investimento" className="text-blue-600 hover:underline">
          Rendimentos de Investimento
        </Link>
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
        <p className="text-red-600" role="alert">
          {erro}
        </p>
      ) : resumo ? (
        <div className="space-y-4">
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {/* Seção 1 — Rendas Operacionais */}
            <Card>
              <CardContent className="pt-4">
                <p className="text-sm font-medium mb-1">Rendas Operacionais</p>
                <p className="text-xl font-semibold text-green-600 mb-3">
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
                <p className="text-xl font-semibold text-red-600 mb-3">
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
                className={`text-2xl font-semibold ${resumo.saldo >= 0 ? 'text-green-600' : 'text-red-600'}`}
              >
                {formatarMoeda(resumo.saldo)}
              </p>
            </CardContent>
          </Card>
        </div>
      ) : null}
    </div>
  )
}
