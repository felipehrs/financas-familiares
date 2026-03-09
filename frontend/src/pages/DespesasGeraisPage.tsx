import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { BackToDashboard } from '@/components/BackToDashboard'
import { useAuth } from '@/hooks/useAuth'
import { listarDespesasGerais, criarDespesaGeral, atualizarDespesaGeral, excluirDespesaGeral } from '@/offline/despesas_gerais'
import { listarMembros } from '@/api/membros'
import { listarCategorias } from '@/api/categorias'
import type { DespesaGeral } from '@/types/despesa_geral'
import type { Membro } from '@/types/membro'
import type { Categoria } from '@/types/categoria'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

// ─── Constantes ───────────────────────────────────────────────────────────────

const FORMAS_PAGAMENTO = [
  { value: 'dinheiro', label: 'Dinheiro' },
  { value: 'debito', label: 'Débito' },
  { value: 'pix', label: 'PIX' },
]

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

const SELECT_CLASS = 'mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2'

// ─── Schema de validação ──────────────────────────────────────────────────────

const despesaGeralSchema = z.object({
  data: z.string().min(1, 'Data é obrigatória'),
  membro_id: z.string().min(1, 'Membro é obrigatório'),
  descricao: z.string().min(2, 'Descrição é obrigatória e deve ter pelo menos 2 caracteres'),
  categoria_id: z.string().optional(),
  valor: z.coerce.number({ invalid_type_error: 'Informe o valor' }).positive('Valor deve ser maior que zero'),
  forma_pagamento: z.string().min(1, 'Forma de pagamento é obrigatória'),
  observacoes: z.string().optional(),
})

type DespesaGeralFormValues = z.infer<typeof despesaGeralSchema>

// ─── Helpers ─────────────────────────────────────────────────────────────────

function formatarMoeda(valor: number): string {
  return valor.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })
}

// ─── Componente ───────────────────────────────────────────────────────────────

type FormMode = 'hidden' | 'criar' | 'editar'

export function DespesasGeraisPage() {
  const { accessToken } = useAuth()

  const hoje = new Date()

  const [despesas, setDespesas] = useState<DespesaGeral[]>([])
  const [membros, setMembros] = useState<Membro[]>([])
  const [categorias, setCategorias] = useState<Categoria[]>([])
  const [loading, setLoading] = useState(true)
  const [apiError, setApiError] = useState<string | null>(null)
  const [formMode, setFormMode] = useState<FormMode>('hidden')
  const [editandoId, setEditandoId] = useState<string | null>(null)

  // ─── Estado do filtro de mês/ano ──────────────────────────────────────────
  const [filtroMes, setFiltroMes] = useState(hoje.getMonth() + 1)
  const [filtroAno, setFiltroAno] = useState(hoje.getFullYear())

  // ─── Estado dos filtros adicionais ────────────────────────────────────────
  const [filtroMembroId, setFiltroMembroId] = useState('')
  const [filtroCategoriaId, setFiltroCategoriaId] = useState('')
  const [filtroForma, setFiltroForma] = useState('')

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<DespesaGeralFormValues>({
    resolver: zodResolver(despesaGeralSchema),
    defaultValues: {
      data: '',
      membro_id: '',
      descricao: '',
      categoria_id: '',
      valor: '' as unknown as number,
      forma_pagamento: '',
      observacoes: '',
    },
  })

  // ─── Carrega dados ──────────────────────────────────────────────────────────

  async function carregarDados(mes?: number, ano?: number) {
    if (!accessToken) return
    try {
      const [listaDespesas, listaMembros, listaCategorias] = await Promise.all([
        listarDespesasGerais(accessToken, mes, ano),
        listarMembros(accessToken),
        listarCategorias(accessToken),
      ])
      setDespesas(listaDespesas)
      setMembros(listaMembros)
      setCategorias(listaCategorias)
    } catch (err) {
      setApiError(err instanceof Error ? err.message : 'Erro ao carregar dados')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void carregarDados(filtroMes, filtroAno)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [filtroMes, filtroAno])

  // ─── Helpers ────────────────────────────────────────────────────────────────

  function nomeMembro(membroId: string): string {
    return membros.find((m) => m.id === membroId)?.nome ?? membroId
  }

  // ─── Handlers de UI ─────────────────────────────────────────────────────────

  function abrirFormCriacao() {
    setFormMode('criar')
    setEditandoId(null)
    setApiError(null)
    reset({
      data: '',
      membro_id: '',
      descricao: '',
      categoria_id: '',
      valor: undefined,
      forma_pagamento: '',
      observacoes: '',
    })
  }

  function abrirFormEdicao(despesa: DespesaGeral) {
    setFormMode('editar')
    setEditandoId(despesa.id)
    setApiError(null)
    reset({
      data: despesa.data,
      membro_id: despesa.membro_id,
      descricao: despesa.descricao,
      categoria_id: despesa.categoria_id ?? '',
      valor: despesa.valor,
      forma_pagamento: despesa.forma_pagamento,
      observacoes: despesa.observacoes ?? '',
    })
  }

  function fecharForm() {
    setFormMode('hidden')
    setEditandoId(null)
    setApiError(null)
    reset({
      data: '',
      membro_id: '',
      descricao: '',
      categoria_id: '',
      valor: undefined,
      forma_pagamento: '',
      observacoes: '',
    })
  }

  // ─── Submit ─────────────────────────────────────────────────────────────────

  async function onSubmit(values: DespesaGeralFormValues) {
    if (!accessToken) return
    setApiError(null)

    try {
      if (formMode === 'criar') {
        await criarDespesaGeral(accessToken, {
          data: values.data,
          membro_id: values.membro_id,
          descricao: values.descricao,
          categoria_id: values.categoria_id || undefined,
          valor: values.valor,
          forma_pagamento: values.forma_pagamento,
          observacoes: values.observacoes || undefined,
        })
      } else if (formMode === 'editar' && editandoId) {
        await atualizarDespesaGeral(accessToken, editandoId, {
          data: values.data,
          membro_id: values.membro_id,
          descricao: values.descricao,
          categoria_id: values.categoria_id || undefined,
          valor: values.valor,
          forma_pagamento: values.forma_pagamento,
          observacoes: values.observacoes || undefined,
        })
      }

      fecharForm()
      setLoading(true)
      await carregarDados(filtroMes, filtroAno)
    } catch (err) {
      setApiError(err instanceof Error ? err.message : 'Erro desconhecido')
    }
  }

  async function handleExcluir(id: string) {
    if (!accessToken) return
    setApiError(null)
    try {
      await excluirDespesaGeral(accessToken, id)
      setLoading(true)
      await carregarDados(filtroMes, filtroAno)
    } catch (err) {
      setApiError(err instanceof Error ? err.message : 'Erro ao excluir despesa')
    }
  }

  // ─── Filtros adicionais (aplicados localmente) ────────────────────────────

  function despesasFiltradas(): DespesaGeral[] {
    return despesas.filter((d) => {
      if (filtroMembroId && d.membro_id !== filtroMembroId) return false
      if (filtroCategoriaId && d.categoria_id !== filtroCategoriaId) return false
      if (filtroForma && d.forma_pagamento !== filtroForma) return false
      return true
    })
  }

  // ─── Exportação CSV ───────────────────────────────────────────────────────

  function exportarCSV() {
    const mes = String(filtroMes).padStart(2, '0')
    const ano = String(filtroAno)
    const nomeArquivo = `despesas-gerais-${mes}-${ano}.csv`

    const cabecalho = ['Data', 'Descrição', 'Membro', 'Categoria', 'Valor', 'Forma de Pagamento', 'Observações']
    const linhas = despesasFiltradas().map((d) => {
      const nomeMem = membros.find((m) => m.id === d.membro_id)?.nome ?? d.membro_id
      const nomeCat = d.categoria_id ? (categorias.find((c) => c.id === d.categoria_id)?.nome ?? '') : ''
      return [
        d.data,
        d.descricao,
        nomeMem,
        nomeCat,
        d.valor.toFixed(2),
        d.forma_pagamento,
        d.observacoes ?? '',
      ]
    })

    const escapar = (v: string) => `"${v.replace(/"/g, '""')}"`
    const conteudo = [cabecalho, ...linhas].map((row) => row.map(escapar).join(',')).join('\n')

    const blob = new Blob(['\uFEFF' + conteudo], { type: 'text/csv;charset=utf-8;' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', nomeArquivo)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
  }

  // ─── Render ──────────────────────────────────────────────────────────────────

  const filtradas = despesasFiltradas()
  const despesasOrdenadas = [...filtradas].sort((a, b) => b.data.localeCompare(a.data))
  const totalFiltrado = filtradas.reduce((acc, d) => acc + d.valor, 0)

  return (
    <div className="p-6 max-w-2xl mx-auto">
      <BackToDashboard />

      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Despesas Gerais</h1>
        {formMode === 'hidden' && (
          <div className="flex gap-2">
            <Button variant="outline" onClick={exportarCSV}>Exportar CSV</Button>
            <Button onClick={abrirFormCriacao}>Nova despesa</Button>
          </div>
        )}
      </div>

      {/* Erro global de API */}
      {apiError && (
        <p className="mb-4 text-sm text-red-600 dark:text-red-400" role="alert">
          {apiError}
        </p>
      )}

      {/* Filtro de mês/ano */}
      <div className="mb-6 flex flex-wrap gap-3 items-end">
        <div>
          <label htmlFor="filtro-mes" className="block text-sm font-medium mb-1">
            Mês
          </label>
          <select
            id="filtro-mes"
            aria-label="Mês"
            value={filtroMes}
            onChange={(e) => setFiltroMes(Number(e.target.value))}
            className="mt-1 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
          >
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
        <div>
          <label htmlFor="filtro-membro" className="block text-sm font-medium mb-1">
            Membro
          </label>
          <select
            id="filtro-membro"
            aria-label="Filtrar por membro"
            value={filtroMembroId}
            onChange={(e) => setFiltroMembroId(e.target.value)}
            className="mt-1 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
          >
            <option value="">Todos</option>
            {membros.map((m) => (
              <option key={m.id} value={m.id}>
                {m.nome}
              </option>
            ))}
          </select>
        </div>
        <div>
          <label htmlFor="filtro-categoria" className="block text-sm font-medium mb-1">
            Categoria
          </label>
          <select
            id="filtro-categoria"
            aria-label="Filtrar por categoria"
            value={filtroCategoriaId}
            onChange={(e) => setFiltroCategoriaId(e.target.value)}
            className="mt-1 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
          >
            <option value="">Todas</option>
            {categorias.map((cat) => (
              <option key={cat.id} value={cat.id}>
                {cat.nome}
              </option>
            ))}
          </select>
        </div>
        <div>
          <label htmlFor="filtro-forma" className="block text-sm font-medium mb-1">
            Forma de pagamento
          </label>
          <select
            id="filtro-forma"
            aria-label="Filtrar por forma de pagamento"
            value={filtroForma}
            onChange={(e) => setFiltroForma(e.target.value)}
            className="mt-1 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
          >
            <option value="">Todas</option>
            {FORMAS_PAGAMENTO.map((f) => (
              <option key={f.value} value={f.value}>
                {f.label}
              </option>
            ))}
          </select>
        </div>
      </div>

      {/* Total das despesas filtradas */}
      <p className="mb-4 text-sm font-medium">
        Total: {formatarMoeda(totalFiltrado)} ({filtradas.length} {filtradas.length === 1 ? 'despesa' : 'despesas'})
      </p>

      {/* Formulário inline de criação / edição */}
      {formMode !== 'hidden' && (
        <Card className="mb-6">
          <CardContent className="pt-4">
            <form onSubmit={handleSubmit(onSubmit)} noValidate>
              <div className="mb-4">
                <Label htmlFor="data">Data</Label>
                <Input
                  id="data"
                  type="date"
                  {...register('data')}
                  className="mt-1"
                />
                {errors.data && (
                  <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.data.message}</p>
                )}
              </div>

              <div className="mb-4">
                <Label htmlFor="membro_id">Membro responsável</Label>
                <select
                  id="membro_id"
                  {...register('membro_id')}
                  className={SELECT_CLASS}
                >
                  <option value="">Selecione um membro</option>
                  {membros.filter((m) => m.ativo).map((membro) => (
                    <option key={membro.id} value={membro.id}>
                      {membro.nome}
                    </option>
                  ))}
                </select>
                {errors.membro_id && (
                  <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.membro_id.message}</p>
                )}
              </div>

              <div className="mb-4">
                <Label htmlFor="descricao">Descrição</Label>
                <Input
                  id="descricao"
                  {...register('descricao')}
                  placeholder="Ex: Supermercado, Farmácia..."
                  className="mt-1"
                />
                {errors.descricao && (
                  <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.descricao.message}</p>
                )}
              </div>

              <div className="mb-4">
                <Label htmlFor="categoria_id">Categoria (opcional)</Label>
                <select
                  id="categoria_id"
                  {...register('categoria_id')}
                  className={SELECT_CLASS}
                >
                  <option value="">Sem categoria</option>
                  {categorias.map((cat) => (
                    <option key={cat.id} value={cat.id}>
                      {cat.nome}
                    </option>
                  ))}
                </select>
              </div>

              <div className="mb-4">
                <Label htmlFor="valor">Valor (R$)</Label>
                <Input
                  id="valor"
                  type="number"
                  min={0}
                  step="0.01"
                  {...register('valor')}
                  placeholder="Ex: 150.00"
                  className="mt-1"
                />
                {errors.valor && (
                  <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.valor.message}</p>
                )}
              </div>

              <div className="mb-4">
                <Label htmlFor="forma_pagamento">Forma de pagamento</Label>
                <select
                  id="forma_pagamento"
                  {...register('forma_pagamento')}
                  className={SELECT_CLASS}
                >
                  <option value="">Selecione uma forma de pagamento</option>
                  {FORMAS_PAGAMENTO.map((forma) => (
                    <option key={forma.value} value={forma.value}>
                      {forma.label}
                    </option>
                  ))}
                </select>
                {errors.forma_pagamento && (
                  <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.forma_pagamento.message}</p>
                )}
              </div>

              <div className="mb-4">
                <Label htmlFor="observacoes">Observações (opcional)</Label>
                <textarea
                  id="observacoes"
                  {...register('observacoes')}
                  placeholder="Observações adicionais..."
                  rows={3}
                  className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 resize-none"
                />
              </div>

              <div className="flex gap-2">
                <Button type="submit" disabled={isSubmitting}>
                  {isSubmitting ? 'Salvando...' : 'Salvar'}
                </Button>
                <Button type="button" variant="outline" onClick={fecharForm}>
                  Cancelar
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      )}

      {/* Lista de despesas gerais */}
      {loading ? (
        <p className="text-muted-foreground">Carregando...</p>
      ) : despesasOrdenadas.length === 0 ? (
        <p className="text-muted-foreground">Nenhuma despesa cadastrada</p>
      ) : (
        <ul className="space-y-3">
          {despesasOrdenadas.map((despesa) => (
            <li key={despesa.id}>
              <Card>
                <CardContent className="pt-4 flex items-center justify-between">
                  <div>
                    <p className="font-medium">{despesa.descricao}</p>
                    <p className="text-sm text-muted-foreground">{nomeMembro(despesa.membro_id)}</p>
                    <p className="text-sm text-muted-foreground">
                      {despesa.data}
                      {' · '}
                      {formatarMoeda(despesa.valor)}
                      {' · '}
                      {despesa.forma_pagamento}
                    </p>
                  </div>
                  <div className="flex gap-2 flex-wrap justify-end">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => abrirFormEdicao(despesa)}
                    >
                      Editar
                    </Button>
                    <Button
                      variant="destructive"
                      size="sm"
                      onClick={() => handleExcluir(despesa.id)}
                    >
                      Excluir
                    </Button>
                  </div>
                </CardContent>
              </Card>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
