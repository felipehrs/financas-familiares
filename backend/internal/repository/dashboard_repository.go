package repository

import (
	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

// DashboardRepository realiza queries agregadas para o dashboard.
type DashboardRepository struct {
	db *sqlx.DB
}

// NewDashboardRepository cria um novo DashboardRepository.
func NewDashboardRepository(db *sqlx.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

type categoriaTotalRow struct {
	Nome  string  `db:"categoria_nome"`
	Total float64 `db:"total"`
}

// DespesasPorCategoria retorna o total de despesas agrupado por categoria para um mês/ano.
// Inclui despesas de cartão (parcelas do mês), assinaturas ativas, contas fixas ativas e despesas gerais do mês.
// Itens sem categoria aparecem como "Sem categoria".
func (r *DashboardRepository) DespesasPorCategoria(familiaID string, mes, ano int) ([]domain.CategoriaTotalRaw, error) {
	var rows []categoriaTotalRow
	err := r.db.Select(&rows, `
		SELECT
			COALESCE(c.nome, 'Sem categoria') AS categoria_nome,
			SUM(d.valor) AS total
		FROM (
			SELECT categoria_id, valor_parcela AS valor
			FROM despesas_cartao
			WHERE fatura_mes = $2 AND fatura_ano = $3 AND familia_id = $1 AND deleted_at IS NULL

			UNION ALL

			SELECT categoria_id, valor
			FROM assinaturas
			WHERE status = 'ativa' AND familia_id = $1 AND deleted_at IS NULL

			UNION ALL

			SELECT categoria_id, valor
			FROM contas_fixas
			WHERE ativa = true AND familia_id = $1 AND deleted_at IS NULL

			UNION ALL

			SELECT categoria_id, valor
			FROM despesas_gerais
			WHERE EXTRACT(MONTH FROM data) = $2
			  AND EXTRACT(YEAR FROM data) = $3
			  AND familia_id = $1
			  AND deleted_at IS NULL
		) AS d
		LEFT JOIN categorias c ON c.id = d.categoria_id AND c.deleted_at IS NULL
		GROUP BY COALESCE(c.nome, 'Sem categoria')
		ORDER BY total DESC
	`, familiaID, mes, ano)
	if err != nil {
		return nil, err
	}

	result := make([]domain.CategoriaTotalRaw, 0, len(rows))
	for _, row := range rows {
		result = append(result, domain.CategoriaTotalRaw{
			Nome:  row.Nome,
			Total: row.Total,
		})
	}
	return result, nil
}
