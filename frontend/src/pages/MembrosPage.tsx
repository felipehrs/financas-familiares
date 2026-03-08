import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Link } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import {
  listarMembros,
  criarMembro,
  atualizarMembro,
  inativarMembro,
} from '@/offline/membros'
import type { Membro, AtualizarMembroRequest } from '@/types/membro'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

// ─── Schema de validação ──────────────────────────────────────────────────────

const membroSchema = z.object({
  nome: z
    .string()
    .min(2, 'Nome é obrigatório e deve ter pelo menos 2 caracteres'),
  relacionamento: z.string().optional(),
})

type MembroFormValues = z.infer<typeof membroSchema>

// ─── Componente ───────────────────────────────────────────────────────────────

type FormMode = 'hidden' | 'criar' | 'editar'

export function MembrosPage() {
  const { accessToken } = useAuth()

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
  } = useForm<MembroFormValues>({
    resolver: zodResolver(membroSchema),
    defaultValues: { nome: '', relacionamento: '' },
  })

  // ─── Carrega lista ──────────────────────────────────────────────────────────

  async function carregarMembros() {
    if (!accessToken) return
    try {
      const lista = await listarMembros(accessToken)
      setMembros(lista)
    } catch (err) {
      setApiError(err instanceof Error ? err.message : 'Erro ao carregar membros')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void carregarMembros()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  // ─── Handlers de UI ─────────────────────────────────────────────────────────

  function abrirFormCriacao() {
    setFormMode('criar')
    setEditandoId(null)
    setApiError(null)
    reset({ nome: '', relacionamento: '' })
  }

  function abrirFormEdicao(membro: Membro) {
    setFormMode('editar')
    setEditandoId(membro.id)
    setApiError(null)
    reset({ nome: membro.nome, relacionamento: membro.relacionamento ?? '' })
  }

  function fecharForm() {
    setFormMode('hidden')
    setEditandoId(null)
    setApiError(null)
    reset({ nome: '', relacionamento: '' })
  }

  // ─── Submit ─────────────────────────────────────────────────────────────────

  async function onSubmit(values: MembroFormValues) {
    if (!accessToken) return
    setApiError(null)

    try {
      if (formMode === 'criar') {
        await criarMembro(accessToken, {
          nome: values.nome,
          relacionamento: values.relacionamento || undefined,
        })
      } else if (formMode === 'editar' && editandoId) {
        const membro = membros.find((m) => m.id === editandoId)!
        const payload: AtualizarMembroRequest = {
          nome: values.nome,
          relacionamento: values.relacionamento || undefined,
          ativo: membro.ativo,
        }
        await atualizarMembro(accessToken, editandoId, payload)
      }

      fecharForm()
      setLoading(true)
      await carregarMembros()
    } catch (err) {
      setApiError(err instanceof Error ? err.message : 'Erro desconhecido')
    }
  }

  async function handleInativar(id: string) {
    if (!accessToken) return
    setApiError(null)
    try {
      await inativarMembro(accessToken, id)
      setLoading(true)
      await carregarMembros()
    } catch (err) {
      setApiError(err instanceof Error ? err.message : 'Erro ao inativar membro')
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
        <h1 className="text-2xl font-bold">Membros da Família</h1>
        {formMode === 'hidden' && (
          <Button onClick={abrirFormCriacao}>Adicionar membro</Button>
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
                <Label htmlFor="nome">Nome</Label>
                <Input
                  id="nome"
                  {...register('nome')}
                  placeholder="Nome do membro"
                  className="mt-1"
                />
                {errors.nome && (
                  <p className="mt-1 text-sm text-red-600">{errors.nome.message}</p>
                )}
              </div>

              <div className="mb-4">
                <Label htmlFor="relacionamento">Relacionamento</Label>
                <Input
                  id="relacionamento"
                  {...register('relacionamento')}
                  placeholder="Ex: Cônjuge, Filho(a), Pai..."
                  className="mt-1"
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

      {/* Lista de membros */}
      {loading ? (
        <p className="text-muted-foreground">Carregando...</p>
      ) : membros.length === 0 ? (
        <p className="text-muted-foreground">Nenhum membro cadastrado</p>
      ) : (
        <ul className="space-y-3">
          {membros.map((membro) => (
            <li key={membro.id}>
              <Card>
                <CardContent className="pt-4 flex items-center justify-between">
                  <div>
                    <p className="font-medium">{membro.nome}</p>
                    {membro.relacionamento && (
                      <p className="text-sm text-muted-foreground">{membro.relacionamento}</p>
                    )}
                    <span
                      className={`inline-block mt-1 text-xs font-semibold px-2 py-0.5 rounded-full ${
                        membro.ativo
                          ? 'bg-green-100 text-green-700'
                          : 'bg-gray-100 text-gray-500'
                      }`}
                    >
                      {membro.ativo ? 'Ativo' : 'Inativo'}
                    </span>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => abrirFormEdicao(membro)}
                    >
                      Editar
                    </Button>
                    {membro.ativo && (
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => handleInativar(membro.id)}
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
