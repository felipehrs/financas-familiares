import { useEffect, useState } from 'react'
import { useAuth } from '@/hooks/useAuth'
import { buscarHistoricoRendas } from '@/api/rendas_historico'
import { listarMembros } from '@/api/membros'
import type { HistoricoRendas, FiltroHistoricoRendas, TipoRenda } from '@/types/renda_historico'
import type { Membro } from '@/types/membro'
import { Card, CardContent } from '@/components/ui/card'

// ─── Constantes ───────────────────────────────────────────────────────────────

const MESES = [
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

const SELECT_CLASS =
  'mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2'

const TIPO_LABELS: Record<TipoRenda, string> = {
  fixa: 'Renda Fixa',
  variavel: 'Renda Variável',
  extra: 'Renda Extra',
  investimento: 'Rendimento de Investimento',
}

const TIPO_BADGE_CLASSES: Record<TipoRenda, string> = {
  fixa: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400',
  variavel: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400',
  extra: 'bg-pink-100 text-pink-800 dark:bg-pink-900/30 dark:text-pink-400',
  investimento: 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400',
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

function formatarMoeda(valor: number): string {
  return new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(valor)
}

function formatarPeriodo(
  tipo: TipoRenda,
  mes?: number,
  ano?: number,
  data?: string,
  filtroMes?: number,
  filtroAno?: number,
): string {
  if (tipo === 'fixa') {
    if (filtroMes && filtroAno) {
      const mesStr = String(filtroMes).padStart(2, '0')
      return `${mesStr}/${filtroAno}`
    }
    return 'Recorrente'
  }
  if (tipo === 'variavel' && mes && ano) {
    const mesStr = String(mes).padStart(2, '0')
    return `${mesStr}/${ano}`
  }
  if ((tipo === 'extra' || tipo === 'investimento') && data) {
    // data is "YYYY-MM-DD"
    const [year, month, day] = data.split('-')
    return `${day}/${month}/${year}`
  }
  return '—'
}

// ─── Componente ───────────────────────────────────────────────────────────────

export function RendasHistoricoPage() {
  const { accessToken } = useAuth()

  const hoje = new Date()

  const [historico, setHistorico] = useState<HistoricoRendas | null>(null)
  const [membros, setMembros] = useState<Membro[]>([])
  const [loading, setLoading] = useState(true)
  const [apiError, setApiError] = useState<string | null>(null)

  const [filtroTipo, setFiltroTipo] = useState<TipoRenda | ''>('')
  const [filtroMembroId, setFiltroMembroId] = useState('')
  const [filtroMes, setFiltroMes] = useState(0)
  const [filtroAno, setFiltroAno] = useState(hoje.getFullYear())

  // ─── Carrega dados ──────────────────────────────────────────────────────────

  useEffect(() => {
    if (!accessToken) return

    const filtro: FiltroHistoricoRendas = {}
    if (filtroTipo) filtro.tipo = filtroTipo
    if (filtroMembroId) filtro.membro_id = filtroMembroId
    if (filtroMes !== 0) filtro.mes = filtroMes
    if (filtroAno !== 0) filtro.ano = filtroAno

    setLoading(true)
    setApiError(null)

    void Promise.all([
      buscarHistoricoRendas(accessToken, filtro),
      listarMembros(accessToken),
    ])
      .then(([hist, memb]) => {
        setHistorico(hist)
        setMembros(memb)
      })
      .catch((err: unknown) => {
        setApiError(err instanceof Error ? err.message : 'Erro ao carregar dados')
      })
      .finally(() => {
        setLoading(false)
      })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [filtroTipo, filtroMembroId, filtroMes, filtroAno])

  // ─── Helpers ────────────────────────────────────────────────────────────────

  function nomeMembroById(membroId: string): string {
    return membros.find((m) => m.id === membroId)?.nome ?? membroId
  }

  // ─── Render ──────────────────────────────────────────────────────────────────

  return (
    <div className="p-6 max-w-4xl mx-auto">

      <h1 className="text-2xl font-bold mb-6">Histórico de Rendas</h1>

      {/* Erro global de API */}
      {apiError && (
        <p className="mb-4 text-sm text-red-600 dark:text-red-400" role="alert">
          {apiError}
        </p>
      )}

      {/* Filtros */}
      <div className="mb-6 flex flex-wrap gap-3 items-end">
        <div>
          <label htmlFor="filtro-tipo" className="block text-sm font-medium mb-1">
            Tipo
          </label>
          <select
            id="filtro-tipo"
            aria-label="Tipo"
            value={filtroTipo}
            onChange={(e) => setFiltroTipo(e.target.value as TipoRenda | '')}
            className={SELECT_CLASS}
          >
            <option value="">Todos os tipos</option>
            <option value="fixa">Renda Fixa</option>
            <option value="variavel">Renda Variável</option>
            <option value="extra">Renda Extra</option>
            <option value="investimento">Rendimento de Investimento</option>
          </select>
        </div>

        <div>
          <label htmlFor="filtro-membro" className="block text-sm font-medium mb-1">
            Membro
          </label>
          <select
            id="filtro-membro"
            aria-label="Membro"
            value={filtroMembroId}
            onChange={(e) => setFiltroMembroId(e.target.value)}
            className={SELECT_CLASS}
          >
            <option value="">Todos os membros</option>
            {membros.map((m) => (
              <option key={m.id} value={m.id}>
                {m.nome}
              </option>
            ))}
          </select>
        </div>

        <div>
          <label htmlFor="filtro-mes" className="block text-sm font-medium mb-1">
            Mês
          </label>
          <select
            id="filtro-mes"
            aria-label="Mês"
            value={filtroMes}
            onChange={(e) => setFiltroMes(Number(e.target.value))}
            className={SELECT_CLASS}
          >
            <option value={0}>Todos os meses</option>
            {MESES.map((m) => (
              <option key={m.value} value={m.value}>
                {m.label}
              </option>
            ))}
          </select>
        </div>

        <div>
          <label htmlFor="filtro-ano" className="block text-sm font-medium mb-1">
            Ano
          </label>
          <input
            id="filtro-ano"
            aria-label="Ano"
            type="number"
            value={filtroAno}
            onChange={(e) => setFiltroAno(Number(e.target.value))}
            className="mt-1 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 w-28"
          />
        </div>
      </div>

      {/* Conteúdo */}
      {loading ? (
        <p className="text-muted-foreground">Carregando...</p>
      ) : !historico || historico.itens.length === 0 ? (
        <p className="text-muted-foreground">Nenhuma renda encontrada</p>
      ) : (
        <>
          {/* Tabela */}
          <div className="overflow-x-auto mb-6">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b">
                  <th className="text-left py-2 pr-4 font-medium">Tipo</th>
                  <th className="text-left py-2 pr-4 font-medium">Descrição</th>
                  <th className="text-left py-2 pr-4 font-medium">Membro</th>
                  <th className="text-right py-2 pr-4 font-medium">Valor</th>
                  <th className="text-left py-2 font-medium">Período</th>
                </tr>
              </thead>
              <tbody>
                {historico.itens.map((item) => (
                  <tr key={item.id} className="border-b last:border-0">
                    <td className="py-2 pr-4">
                      <span
                        className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${TIPO_BADGE_CLASSES[item.tipo]}`}
                      >
                        {TIPO_LABELS[item.tipo]}
                      </span>
                    </td>
                    <td className="py-2 pr-4">{item.descricao}</td>
                    <td className="py-2 pr-4 text-muted-foreground">{nomeMembroById(item.membro_id)}</td>
                    <td className="py-2 pr-4 text-right font-medium">{formatarMoeda(item.valor)}</td>
                    <td className="py-2 text-muted-foreground">
                      {formatarPeriodo(item.tipo, item.mes, item.ano, item.data, filtroMes || undefined, filtroAno || undefined)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Resumo */}
          <Card>
            <CardContent className="pt-4">
              <p className="text-sm font-medium mb-3">Resumo</p>
              <div className="space-y-1 text-sm">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Rendas Fixas</span>
                  <span>{formatarMoeda(historico.resumo.total_fixas)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Rendas Variáveis</span>
                  <span>{formatarMoeda(historico.resumo.total_variaveis)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Rendas Extras</span>
                  <span>{formatarMoeda(historico.resumo.total_extras)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Investimentos Distribuídos</span>
                  <span>{formatarMoeda(historico.resumo.total_investimentos_distribuidos)}</span>
                </div>
                <div className="flex justify-between font-semibold border-t pt-1 mt-1">
                  <span>Total Geral</span>
                  <span>{formatarMoeda(historico.resumo.total_geral)}</span>
                </div>
              </div>
            </CardContent>
          </Card>
        </>
      )}
    </div>
  )
}
