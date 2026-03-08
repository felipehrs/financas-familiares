import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useAuth } from '@/hooks/useAuth'
import { listarDespesasPorCartao, criarDespesa, excluirDespesa } from '@/api/despesas_cartao'
import { listarCartoes } from '@/api/cartoes_credito'
import { listarCategorias } from '@/api/categorias'
import type { DespesaCartao } from '@/types/despesa_cartao'
import type { CartaoCredito } from '@/types/cartao_credito'
import type { Categoria } from '@/types/categoria'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

// ─── Schema de validação ──────────────────────────────────────────────────────

const despesaSchema = z.object({
  descricao: z.string().min(2, 'Descrição é obrigatória e deve ter pelo menos 2 caracteres'),
  data_compra: z.string().min(1, 'Data da compra é obrigatória'),
  valor_total: z.coerce
    .number({ invalid_type_error: 'Informe o valor total' })
    .positive('Valor deve ser maior que zero'),
  categoria_id: z.string().optional(),
})

type DespesaFormValues = z.infer<typeof despesaSchema>

// ─── Helpers ─────────────────────────────────────────────────────────────────

function formatarMoeda(valor: number): string {
  return valor.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })
}

function agruparPorFatura(despesas: DespesaCartao[]): Map<string, DespesaCartao[]> {
  const grupos = new Map<string, DespesaCartao[]>()
  for (const despesa of despesas) {
    const chave = despesa.fatura
    if (!grupos.has(chave)) {
      grupos.set(chave, [])
    }
    grupos.get(chave)!.push(despesa)
  }
  return grupos
}

// ─── Componente ───────────────────────────────────────────────────────────────

export function DespesasCartaoPage() {
  const { cartaoId } = useParams<{ cartaoId: string }>()
  const { accessToken } = useAuth()

  const [despesas, setDespesas] = useState<DespesaCartao[]>([])
  const [cartao, setCartao] = useState<CartaoCredito | null>(null)
  const [categorias, setCategorias] = useState<Categoria[]>([])
  const [loading, setLoading] = useState(true)
  const [apiError, setApiError] = useState<string | null>(null)
  const [mostrarForm, setMostrarForm] = useState(false)

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<DespesaFormValues>({
    resolver: zodResolver(despesaSchema),
    defaultValues: {
      descricao: '',
      data_compra: '',
      valor_total: '' as unknown as number,
      categoria_id: '',
    },
  })

  // ─── Carrega dados ──────────────────────────────────────────────────────────

  async function carregarDados() {
    if (!accessToken || !cartaoId) return
    try {
      const [listaDespesas, listaCartoes, listaCategorias] = await Promise.all([
        listarDespesasPorCartao(accessToken, cartaoId),
        listarCartoes(accessToken),
        listarCategorias(accessToken),
      ])
      setDespesas(listaDespesas)
      setCartao(listaCartoes.find((c) => c.id === cartaoId) ?? null)
      setCategorias(listaCategorias)
    } catch (err) {
      setApiError(err instanceof Error ? err.message : 'Erro ao carregar dados')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void carregarDados()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  // ─── Handlers de UI ─────────────────────────────────────────────────────────

  function abrirForm() {
    setMostrarForm(true)
    setApiError(null)
    reset({ descricao: '', data_compra: '', valor_total: undefined, categoria_id: '' })
  }

  function fecharForm() {
    setMostrarForm(false)
    setApiError(null)
    reset({ descricao: '', data_compra: '', valor_total: undefined, categoria_id: '' })
  }

  // ─── Submit ─────────────────────────────────────────────────────────────────

  async function onSubmit(values: DespesaFormValues) {
    if (!accessToken || !cartaoId) return
    setApiError(null)

    try {
      await criarDespesa(accessToken, cartaoId, {
        descricao: values.descricao,
        data_compra: values.data_compra,
        valor_total: values.valor_total,
        categoria_id: values.categoria_id || undefined,
      })
      fecharForm()
      setLoading(true)
      await carregarDados()
    } catch (err) {
      setApiError(err instanceof Error ? err.message : 'Erro desconhecido')
    }
  }

  async function handleExcluir(id: string) {
    if (!accessToken) return
    setApiError(null)
    try {
      await excluirDespesa(accessToken, id)
      setLoading(true)
      await carregarDados()
    } catch (err) {
      setApiError(err instanceof Error ? err.message : 'Erro ao excluir despesa')
    }
  }

  // ─── Render ──────────────────────────────────────────────────────────────────

  const gruposFatura = agruparPorFatura(despesas)

  return (
    <div className="p-6 max-w-2xl mx-auto">
      <div className="flex items-center justify-between mb-6">
        <div>
          <Link to="/cartoes" className="text-sm text-muted-foreground hover:underline">
            ← Voltar para cartões
          </Link>
          <h1 className="text-2xl font-bold mt-1">
            Despesas{cartao ? ` — ${cartao.nome}` : ''}
          </h1>
        </div>
        {!mostrarForm && (
          <Button onClick={abrirForm}>Nova despesa</Button>
        )}
      </div>

      {/* Erro global de API */}
      {apiError && (
        <p className="mb-4 text-sm text-red-600" role="alert">
          {apiError}
        </p>
      )}

      {/* Formulário inline de criação */}
      {mostrarForm && (
        <Card className="mb-6">
          <CardContent className="pt-4">
            <form onSubmit={handleSubmit(onSubmit)} noValidate>
              <div className="mb-4">
                <Label htmlFor="descricao">Descrição</Label>
                <Input
                  id="descricao"
                  {...register('descricao')}
                  placeholder="Ex: Supermercado, Farmácia..."
                  className="mt-1"
                />
                {errors.descricao && (
                  <p className="mt-1 text-sm text-red-600">{errors.descricao.message}</p>
                )}
              </div>

              <div className="mb-4">
                <Label htmlFor="data_compra">Data da compra</Label>
                <Input
                  id="data_compra"
                  type="date"
                  {...register('data_compra')}
                  className="mt-1"
                />
                {errors.data_compra && (
                  <p className="mt-1 text-sm text-red-600">{errors.data_compra.message}</p>
                )}
              </div>

              <div className="mb-4">
                <Label htmlFor="valor_total">Valor total</Label>
                <Input
                  id="valor_total"
                  type="number"
                  min={0}
                  step="0.01"
                  {...register('valor_total')}
                  placeholder="Ex: 150.00"
                  className="mt-1"
                />
                {errors.valor_total && (
                  <p className="mt-1 text-sm text-red-600">{errors.valor_total.message}</p>
                )}
              </div>

              <div className="mb-4">
                <Label htmlFor="categoria_id">Categoria (opcional)</Label>
                <select
                  id="categoria_id"
                  {...register('categoria_id')}
                  className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                >
                  <option value="">Sem categoria</option>
                  {categorias.map((categoria) => (
                    <option key={categoria.id} value={categoria.id}>
                      {categoria.nome}
                    </option>
                  ))}
                </select>
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

      {/* Lista de despesas */}
      {loading ? (
        <p className="text-muted-foreground">Carregando...</p>
      ) : despesas.length === 0 ? (
        <p className="text-muted-foreground">Nenhuma despesa cadastrada</p>
      ) : (
        <div className="space-y-6">
          {Array.from(gruposFatura.entries()).map(([fatura, itens]) => (
            <div key={fatura}>
              <h2 className="text-lg font-semibold mb-3">{fatura}</h2>
              <ul className="space-y-3">
                {itens.map((despesa) => (
                  <li key={despesa.id}>
                    <Card>
                      <CardContent className="pt-4 flex items-center justify-between">
                        <div>
                          <p className="font-medium">{despesa.descricao}</p>
                          <p className="text-sm text-muted-foreground">
                            {formatarMoeda(despesa.valor_parcela)}
                            {despesa.numero_parcelas > 1 &&
                              ` (${despesa.numero_parcelas}x de ${formatarMoeda(despesa.valor_parcela)})`}
                          </p>
                          <p className="text-sm text-muted-foreground">
                            Compra em {despesa.data_compra} · Fatura {despesa.fatura}
                          </p>
                        </div>
                        <Button
                          variant="destructive"
                          size="sm"
                          onClick={() => handleExcluir(despesa.id)}
                        >
                          Excluir
                        </Button>
                      </CardContent>
                    </Card>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
