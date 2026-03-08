import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useAuth } from '@/hooks/useAuth'
import { listarCartoes, criarCartao, atualizarCartao, inativarCartao } from '@/offline/cartoes_credito'
import { listarMembros } from '@/api/membros'
import type { CartaoCredito, AtualizarCartaoCreditoRequest } from '@/types/cartao_credito'
import type { Membro } from '@/types/membro'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

// ─── Schema de validação ──────────────────────────────────────────────────────

const cartaoSchema = z.object({
  nome: z.string().min(2, 'Nome é obrigatório e deve ter pelo menos 2 caracteres'),
  membro_id: z.string().min(1, 'Membro responsável é obrigatório'),
  dia_fechamento: z.coerce
    .number({ invalid_type_error: 'Informe o dia de fechamento' })
    .int()
    .min(1, 'Deve ser entre 1 e 31')
    .max(31, 'Deve ser entre 1 e 31'),
  dia_vencimento: z.coerce
    .number({ invalid_type_error: 'Informe o dia de vencimento' })
    .int()
    .min(1, 'Deve ser entre 1 e 31')
    .max(31, 'Deve ser entre 1 e 31'),
  limite: z.coerce.number().positive('Limite deve ser positivo').optional().or(z.literal(undefined)),
})

type CartaoFormValues = z.infer<typeof cartaoSchema>

// ─── Componente ───────────────────────────────────────────────────────────────

type FormMode = 'hidden' | 'criar' | 'editar'

export function CartoesPage() {
  const { accessToken } = useAuth()

  const [cartoes, setCartoes] = useState<CartaoCredito[]>([])
  const [membros, setMembros] = useState<Membro[]>([])
  const [loading, setLoading] = useState(true)
  const [apiError, setApiError] = useState<string | null>(null)
  const [formMode, setFormMode] = useState<FormMode>('hidden')
  const [editandoId, setEditandoId] = useState<string | null>(null)

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<CartaoFormValues>({
    resolver: zodResolver(cartaoSchema),
    defaultValues: { nome: '', membro_id: '', dia_fechamento: '' as unknown as number, dia_vencimento: '' as unknown as number, limite: '' as unknown as number | undefined },
  })

  // ─── Carrega dados ──────────────────────────────────────────────────────────

  async function carregarDados() {
    if (!accessToken) return
    try {
      const [listaCartoes, listaMembros] = await Promise.all([
        listarCartoes(accessToken),
        listarMembros(accessToken),
      ])
      setCartoes(listaCartoes)
      setMembros(listaMembros)
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

  // ─── Helpers ────────────────────────────────────────────────────────────────

  function nomeMembro(membroId: string): string {
    return membros.find((m) => m.id === membroId)?.nome ?? membroId
  }

  // ─── Handlers de UI ─────────────────────────────────────────────────────────

  function abrirFormCriacao() {
    setFormMode('criar')
    setEditandoId(null)
    setApiError(null)
    reset({ nome: '', membro_id: '', dia_fechamento: undefined, dia_vencimento: undefined, limite: undefined })
  }

  function abrirFormEdicao(cartao: CartaoCredito) {
    setFormMode('editar')
    setEditandoId(cartao.id)
    setApiError(null)
    reset({
      nome: cartao.nome,
      membro_id: cartao.membro_id,
      dia_fechamento: cartao.dia_fechamento,
      dia_vencimento: cartao.dia_vencimento,
      limite: cartao.limite ?? undefined,
    })
  }

  function fecharForm() {
    setFormMode('hidden')
    setEditandoId(null)
    setApiError(null)
    reset({ nome: '', membro_id: '', dia_fechamento: undefined, dia_vencimento: undefined, limite: undefined })
  }

  // ─── Submit ─────────────────────────────────────────────────────────────────

  async function onSubmit(values: CartaoFormValues) {
    if (!accessToken) return
    setApiError(null)

    try {
      if (formMode === 'criar') {
        await criarCartao(accessToken, {
          nome: values.nome,
          membro_id: values.membro_id,
          dia_fechamento: values.dia_fechamento,
          dia_vencimento: values.dia_vencimento,
          limite: values.limite,
        })
      } else if (formMode === 'editar' && editandoId) {
        const cartao = cartoes.find((c) => c.id === editandoId)!
        const payload: AtualizarCartaoCreditoRequest = {
          nome: values.nome,
          membro_id: values.membro_id,
          dia_fechamento: values.dia_fechamento,
          dia_vencimento: values.dia_vencimento,
          limite: values.limite,
          ativo: cartao.ativo,
        }
        await atualizarCartao(accessToken, editandoId, payload)
      }

      fecharForm()
      setLoading(true)
      await carregarDados()
    } catch (err) {
      setApiError(err instanceof Error ? err.message : 'Erro desconhecido')
    }
  }

  async function handleInativar(id: string) {
    if (!accessToken) return
    setApiError(null)
    try {
      await inativarCartao(accessToken, id)
      setLoading(true)
      await carregarDados()
    } catch (err) {
      setApiError(err instanceof Error ? err.message : 'Erro ao inativar cartão')
    }
  }

  // ─── Render ──────────────────────────────────────────────────────────────────

  return (
    <div className="p-6 max-w-2xl mx-auto">
      <div className="mb-4">
        <Link to="/" className="text-blue-600 hover:underline text-sm">
          ← Voltar para o dashboard
        </Link>
      </div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Cartões de Crédito</h1>
        {formMode === 'hidden' && (
          <Button onClick={abrirFormCriacao}>Adicionar cartão</Button>
        )}
      </div>

      {/* Erro global de API */}
      {apiError && (
        <p className="mb-4 text-sm text-red-600" role="alert">
          {apiError}
        </p>
      )}

      {/* Formulário inline de criação / edição */}
      {formMode !== 'hidden' && (
        <Card className="mb-6">
          <CardContent className="pt-4">
            <form onSubmit={handleSubmit(onSubmit)} noValidate>
              <div className="mb-4">
                <Label htmlFor="nome">Nome do cartão</Label>
                <Input
                  id="nome"
                  {...register('nome')}
                  placeholder="Ex: Nubank, Inter, Itaú..."
                  className="mt-1"
                />
                {errors.nome && (
                  <p className="mt-1 text-sm text-red-600">{errors.nome.message}</p>
                )}
              </div>

              <div className="mb-4">
                <Label htmlFor="membro_id">Membro responsável</Label>
                <select
                  id="membro_id"
                  {...register('membro_id')}
                  className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                >
                  <option value="">Selecione um membro</option>
                  {membros.filter((m) => m.ativo).map((membro) => (
                    <option key={membro.id} value={membro.id}>
                      {membro.nome}
                    </option>
                  ))}
                </select>
                {errors.membro_id && (
                  <p className="mt-1 text-sm text-red-600">{errors.membro_id.message}</p>
                )}
              </div>

              <div className="mb-4 grid grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="dia_fechamento">Dia de fechamento</Label>
                  <Input
                    id="dia_fechamento"
                    type="number"
                    min={1}
                    max={31}
                    {...register('dia_fechamento')}
                    placeholder="Ex: 10"
                    className="mt-1"
                  />
                  {errors.dia_fechamento && (
                    <p className="mt-1 text-sm text-red-600">{errors.dia_fechamento.message}</p>
                  )}
                </div>
                <div>
                  <Label htmlFor="dia_vencimento">Dia de vencimento</Label>
                  <Input
                    id="dia_vencimento"
                    type="number"
                    min={1}
                    max={31}
                    {...register('dia_vencimento')}
                    placeholder="Ex: 17"
                    className="mt-1"
                  />
                  {errors.dia_vencimento && (
                    <p className="mt-1 text-sm text-red-600">{errors.dia_vencimento.message}</p>
                  )}
                </div>
              </div>

              <div className="mb-4">
                <Label htmlFor="limite">Limite (opcional)</Label>
                <Input
                  id="limite"
                  type="number"
                  min={0}
                  step="0.01"
                  {...register('limite')}
                  placeholder="Ex: 5000.00"
                  className="mt-1"
                />
                {errors.limite && (
                  <p className="mt-1 text-sm text-red-600">{errors.limite.message}</p>
                )}
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

      {/* Lista de cartões */}
      {loading ? (
        <p className="text-muted-foreground">Carregando...</p>
      ) : cartoes.length === 0 ? (
        <p className="text-muted-foreground">Nenhum cartão cadastrado</p>
      ) : (
        <ul className="space-y-3">
          {cartoes.map((cartao) => (
            <li key={cartao.id}>
              <Card>
                <CardContent className="pt-4 flex items-center justify-between">
                  <div>
                    <p className="font-medium">{cartao.nome}</p>
                    <p className="text-sm text-muted-foreground">{nomeMembro(cartao.membro_id)}</p>
                    <p className="text-sm text-muted-foreground">
                      Fecha dia {cartao.dia_fechamento} · Vence dia {cartao.dia_vencimento}
                      {cartao.limite != null && ` · Limite R$ ${cartao.limite.toFixed(2)}`}
                    </p>
                    <span
                      className={`inline-block mt-1 text-xs font-semibold px-2 py-0.5 rounded-full ${
                        cartao.ativo
                          ? 'bg-green-100 text-green-700'
                          : 'bg-gray-100 text-gray-500'
                      }`}
                    >
                      {cartao.ativo ? 'Ativo' : 'Inativo'}
                    </span>
                  </div>
                  <div className="flex gap-2">
                    <Link to={`/cartoes/${cartao.id}/despesas`}>
                      <Button variant="outline" size="sm">
                        Ver despesas
                      </Button>
                    </Link>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => abrirFormEdicao(cartao)}
                    >
                      Editar
                    </Button>
                    {cartao.ativo && (
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => handleInativar(cartao.id)}
                      >
                        Inativar
                      </Button>
                    )}
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
