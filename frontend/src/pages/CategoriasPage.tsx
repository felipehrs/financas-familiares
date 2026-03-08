import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Link } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import {
  listarCategorias,
  criarCategoria,
  atualizarCategoria,
  excluirCategoria,
} from '@/api/categorias'
import type { Categoria } from '@/types/categoria'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

// ─── Schema de validação ──────────────────────────────────────────────────────

const categoriaSchema = z.object({
  nome: z
    .string()
    .min(2, 'Nome é obrigatório e deve ter pelo menos 2 caracteres'),
})

type CategoriaFormValues = z.infer<typeof categoriaSchema>

// ─── Componente ───────────────────────────────────────────────────────────────

type FormMode = 'hidden' | 'criar' | 'editar'

export function CategoriasPage() {
  const { accessToken } = useAuth()

  const [categorias, setCategorias] = useState<Categoria[]>([])
  const [loading, setLoading] = useState(true)
  const [apiError, setApiError] = useState<string | null>(null)
  const [formMode, setFormMode] = useState<FormMode>('hidden')
  const [editandoId, setEditandoId] = useState<string | null>(null)

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<CategoriaFormValues>({
    resolver: zodResolver(categoriaSchema),
    defaultValues: { nome: '' },
  })

  // ─── Carrega lista ──────────────────────────────────────────────────────────

  async function carregarCategorias() {
    if (!accessToken) return
    try {
      const lista = await listarCategorias(accessToken)
      setCategorias(lista ?? [])
    } catch (err) {
      setApiError(err instanceof Error ? err.message : 'Erro ao carregar categorias')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void carregarCategorias()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  // ─── Handlers de UI ─────────────────────────────────────────────────────────

  function abrirFormCriacao() {
    setFormMode('criar')
    setEditandoId(null)
    setApiError(null)
    reset({ nome: '' })
  }

  function abrirFormEdicao(categoria: Categoria) {
    setFormMode('editar')
    setEditandoId(categoria.id)
    setApiError(null)
    reset({ nome: categoria.nome })
  }

  function fecharForm() {
    setFormMode('hidden')
    setEditandoId(null)
    setApiError(null)
    reset({ nome: '' })
  }

  // ─── Submit ─────────────────────────────────────────────────────────────────

  async function onSubmit(values: CategoriaFormValues) {
    if (!accessToken) return
    setApiError(null)

    try {
      if (formMode === 'criar') {
        await criarCategoria(accessToken, { nome: values.nome })
      } else if (formMode === 'editar' && editandoId) {
        await atualizarCategoria(accessToken, editandoId, { nome: values.nome })
      }

      fecharForm()
      setLoading(true)
      await carregarCategorias()
    } catch (err) {
      setApiError(err instanceof Error ? err.message : 'Erro desconhecido')
    }
  }

  async function handleExcluir(categoria: Categoria) {
    if (!accessToken) return
    const confirmado = window.confirm(`Deseja excluir a categoria '${categoria.nome}'?`)
    if (!confirmado) return

    setApiError(null)
    try {
      await excluirCategoria(accessToken, categoria.id)
      setLoading(true)
      await carregarCategorias()
    } catch (err) {
      setApiError(err instanceof Error ? err.message : 'Erro ao excluir categoria')
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
        <h1 className="text-2xl font-bold">Categorias</h1>
        {formMode === 'hidden' && (
          <Button onClick={abrirFormCriacao}>Adicionar categoria</Button>
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
                  placeholder="Nome da categoria"
                  className="mt-1"
                />
                {errors.nome && (
                  <p className="mt-1 text-sm text-red-600">{errors.nome.message}</p>
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

      {/* Lista de categorias */}
      {loading ? (
        <p className="text-muted-foreground">Carregando...</p>
      ) : categorias.length === 0 ? (
        <p className="text-muted-foreground">Nenhuma categoria cadastrada</p>
      ) : (
        <ul className="space-y-3">
          {categorias.map((categoria) => (
            <li key={categoria.id}>
              <Card>
                <CardContent className="pt-4 flex items-center justify-between">
                  <p className="font-medium">{categoria.nome}</p>
                  <div className="flex gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => abrirFormEdicao(categoria)}
                    >
                      Editar
                    </Button>
                    <Button
                      variant="destructive"
                      size="sm"
                      onClick={() => handleExcluir(categoria)}
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
