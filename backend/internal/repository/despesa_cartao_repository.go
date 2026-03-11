package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/felipehrs/financas-familiares/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

// DespesaCartaoRepository implementa service.DespesaCartaoRepository usando PostgreSQL via sqlx.
type DespesaCartaoRepository struct {
	db *sqlx.DB
}

// NewDespesaCartaoRepository cria uma nova instância do DespesaCartaoRepository.
func NewDespesaCartaoRepository(db *sqlx.DB) *DespesaCartaoRepository {
	return &DespesaCartaoRepository{db: db}
}

// despesaCartaoRow é a estrutura de scan para as colunas do banco.
type despesaCartaoRow struct {
	ID             string    `db:"id"`
	CompraID       string    `db:"compra_id"`
	CartaoID       string    `db:"cartao_id"`
	CategoriaID    *string   `db:"categoria_id"`
	Descricao      string    `db:"descricao"`
	DataCompra     time.Time `db:"data_compra"`
	ValorTotal     float64   `db:"valor_total"`
	NumeroParcelas int       `db:"numero_parcelas"`
	ParcelaNumero  int       `db:"parcela_numero"`
	ValorParcela   float64   `db:"valor_parcela"`
	FaturaMes      int       `db:"fatura_mes"`
	FaturaAno      int       `db:"fatura_ano"`
}

func (r despesaCartaoRow) toDomain() *domain.DespesaCartao {
	return &domain.DespesaCartao{
		ID:             r.ID,
		CompraID:       r.CompraID,
		CartaoID:       r.CartaoID,
		CategoriaID:    r.CategoriaID,
		Descricao:      r.Descricao,
		DataCompra:     r.DataCompra,
		ValorTotal:     r.ValorTotal,
		NumeroParcelas: r.NumeroParcelas,
		ParcelaNumero:  r.ParcelaNumero,
		ValorParcela:   r.ValorParcela,
		FaturaMes:      r.FaturaMes,
		FaturaAno:      r.FaturaAno,
	}
}

// Criar persiste uma nova despesa de cartão no banco e retorna o registro criado (com ID gerado).
func (r *DespesaCartaoRepository) Criar(familiaID string, d *domain.DespesaCartao) (*domain.DespesaCartao, error) {
	var row despesaCartaoRow
	err := r.db.QueryRowx(`
		INSERT INTO despesas_cartao
			(familia_id, compra_id, cartao_id, categoria_id, descricao, data_compra, valor_total, numero_parcelas, parcela_numero, valor_parcela, fatura_mes, fatura_ano)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, compra_id, cartao_id, categoria_id, descricao, data_compra, valor_total, numero_parcelas, parcela_numero, valor_parcela, fatura_mes, fatura_ano
	`, familiaID, d.CompraID, d.CartaoID, d.CategoriaID, d.Descricao, d.DataCompra, d.ValorTotal, d.NumeroParcelas, d.ParcelaNumero, d.ValorParcela, d.FaturaMes, d.FaturaAno).StructScan(&row)
	if err != nil {
		return nil, err
	}
	return row.toDomain(), nil
}

// ListarPorCartao retorna todas as despesas não excluídas de um cartão, ordenadas por data de compra.
func (r *DespesaCartaoRepository) ListarPorCartao(familiaID, cartaoID string) ([]*domain.DespesaCartao, error) {
	var rows []despesaCartaoRow
	err := r.db.Select(&rows, `
		SELECT id, compra_id, cartao_id, categoria_id, descricao, data_compra, valor_total, numero_parcelas, parcela_numero, valor_parcela, fatura_mes, fatura_ano
		FROM despesas_cartao
		WHERE cartao_id = $1 AND familia_id = $2 AND deleted_at IS NULL
		ORDER BY data_compra DESC, parcela_numero ASC
	`, cartaoID, familiaID)
	if err != nil {
		return nil, err
	}

	despesas := make([]*domain.DespesaCartao, 0, len(rows))
	for _, row := range rows {
		despesas = append(despesas, row.toDomain())
	}
	return despesas, nil
}

// ListarPorFatura retorna as despesas de um cartão filtradas por mês e ano de fatura.
func (r *DespesaCartaoRepository) ListarPorFatura(familiaID, cartaoID string, mes, ano int) ([]*domain.DespesaCartao, error) {
	var rows []despesaCartaoRow
	err := r.db.Select(&rows, `
		SELECT id, compra_id, cartao_id, categoria_id, descricao, data_compra, valor_total, numero_parcelas, parcela_numero, valor_parcela, fatura_mes, fatura_ano
		FROM despesas_cartao
		WHERE cartao_id = $1 AND familia_id = $2 AND fatura_mes = $3 AND fatura_ano = $4 AND deleted_at IS NULL
		ORDER BY data_compra DESC, parcela_numero ASC
	`, cartaoID, familiaID, mes, ano)
	if err != nil {
		return nil, err
	}

	despesas := make([]*domain.DespesaCartao, 0, len(rows))
	for _, row := range rows {
		despesas = append(despesas, row.toDomain())
	}
	return despesas, nil
}

// BuscarPorID busca uma despesa pelo seu UUID.
// Retorna domain.ErrDespesaCartaoNaoEncontrada se não existir ou estiver soft-deleted.
func (r *DespesaCartaoRepository) BuscarPorID(familiaID, id string) (*domain.DespesaCartao, error) {
	var row despesaCartaoRow
	err := r.db.QueryRowx(`
		SELECT id, compra_id, cartao_id, categoria_id, descricao, data_compra, valor_total, numero_parcelas, parcela_numero, valor_parcela, fatura_mes, fatura_ano
		FROM despesas_cartao
		WHERE id = $1 AND familia_id = $2 AND deleted_at IS NULL
	`, id, familiaID).StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrDespesaCartaoNaoEncontrada
		}
		return nil, err
	}
	return row.toDomain(), nil
}

// ListarPorFaturaGlobal retorna todas as despesas de todos os cartões filtradas por mês e ano de fatura.
// Usado pelo DashboardService para calcular o total de despesas de cartão no mês.
func (r *DespesaCartaoRepository) ListarPorFaturaGlobal(familiaID string, mes, ano int) ([]*domain.DespesaCartao, error) {
	var rows []despesaCartaoRow
	err := r.db.Select(&rows, `
		SELECT id, compra_id, cartao_id, categoria_id, descricao, data_compra, valor_total, numero_parcelas, parcela_numero, valor_parcela, fatura_mes, fatura_ano
		FROM despesas_cartao
		WHERE familia_id = $1 AND fatura_mes = $2 AND fatura_ano = $3 AND deleted_at IS NULL
		ORDER BY data_compra DESC, parcela_numero ASC
	`, familiaID, mes, ano)
	if err != nil {
		return nil, err
	}

	despesas := make([]*domain.DespesaCartao, 0, len(rows))
	for _, row := range rows {
		despesas = append(despesas, row.toDomain())
	}
	return despesas, nil
}

// Excluir realiza soft delete de uma despesa setando deleted_at = NOW().
func (r *DespesaCartaoRepository) Excluir(familiaID, id string) error {
	result, err := r.db.Exec(`
		UPDATE despesas_cartao
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $2 AND familia_id = $1 AND deleted_at IS NULL
	`, familiaID, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrDespesaCartaoNaoEncontrada
	}
	return nil
}
