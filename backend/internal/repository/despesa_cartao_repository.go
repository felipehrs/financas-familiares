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
	CartaoID       string    `db:"cartao_id"`
	CategoriaID    *string   `db:"categoria_id"`
	Descricao      string    `db:"descricao"`
	DataCompra     time.Time `db:"data_compra"`
	ValorTotal     float64   `db:"valor_total"`
	NumeroParcelas int       `db:"numero_parcelas"`
	ValorParcela   float64   `db:"valor_parcela"`
	FaturaMes      int       `db:"fatura_mes"`
	FaturaAno      int       `db:"fatura_ano"`
}

func (r despesaCartaoRow) toDomain() *domain.DespesaCartao {
	return &domain.DespesaCartao{
		ID:             r.ID,
		CartaoID:       r.CartaoID,
		CategoriaID:    r.CategoriaID,
		Descricao:      r.Descricao,
		DataCompra:     r.DataCompra,
		ValorTotal:     r.ValorTotal,
		NumeroParcelas: r.NumeroParcelas,
		ValorParcela:   r.ValorParcela,
		FaturaMes:      r.FaturaMes,
		FaturaAno:      r.FaturaAno,
	}
}

// Criar persiste uma nova despesa de cartão no banco e retorna o registro criado (com ID gerado).
func (r *DespesaCartaoRepository) Criar(d *domain.DespesaCartao) (*domain.DespesaCartao, error) {
	var row despesaCartaoRow
	err := r.db.QueryRowx(`
		INSERT INTO despesas_cartao
			(cartao_id, categoria_id, descricao, data_compra, valor_total, numero_parcelas, valor_parcela, fatura_mes, fatura_ano)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, cartao_id, categoria_id, descricao, data_compra, valor_total, numero_parcelas, valor_parcela, fatura_mes, fatura_ano
	`, d.CartaoID, d.CategoriaID, d.Descricao, d.DataCompra, d.ValorTotal, d.NumeroParcelas, d.ValorParcela, d.FaturaMes, d.FaturaAno).StructScan(&row)
	if err != nil {
		return nil, err
	}
	return row.toDomain(), nil
}

// ListarPorCartao retorna todas as despesas não excluídas de um cartão, ordenadas por data de compra.
func (r *DespesaCartaoRepository) ListarPorCartao(cartaoID string) ([]*domain.DespesaCartao, error) {
	var rows []despesaCartaoRow
	err := r.db.Select(&rows, `
		SELECT id, cartao_id, categoria_id, descricao, data_compra, valor_total, numero_parcelas, valor_parcela, fatura_mes, fatura_ano
		FROM despesas_cartao
		WHERE cartao_id = $1 AND deleted_at IS NULL
		ORDER BY data_compra DESC
	`, cartaoID)
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
func (r *DespesaCartaoRepository) ListarPorFatura(cartaoID string, mes, ano int) ([]*domain.DespesaCartao, error) {
	var rows []despesaCartaoRow
	err := r.db.Select(&rows, `
		SELECT id, cartao_id, categoria_id, descricao, data_compra, valor_total, numero_parcelas, valor_parcela, fatura_mes, fatura_ano
		FROM despesas_cartao
		WHERE cartao_id = $1 AND fatura_mes = $2 AND fatura_ano = $3 AND deleted_at IS NULL
		ORDER BY data_compra DESC
	`, cartaoID, mes, ano)
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
func (r *DespesaCartaoRepository) BuscarPorID(id string) (*domain.DespesaCartao, error) {
	var row despesaCartaoRow
	err := r.db.QueryRowx(`
		SELECT id, cartao_id, categoria_id, descricao, data_compra, valor_total, numero_parcelas, valor_parcela, fatura_mes, fatura_ano
		FROM despesas_cartao
		WHERE id = $1 AND deleted_at IS NULL
	`, id).StructScan(&row)
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
func (r *DespesaCartaoRepository) ListarPorFaturaGlobal(mes, ano int) ([]*domain.DespesaCartao, error) {
	var rows []despesaCartaoRow
	err := r.db.Select(&rows, `
		SELECT id, cartao_id, categoria_id, descricao, data_compra, valor_total, numero_parcelas, valor_parcela, fatura_mes, fatura_ano
		FROM despesas_cartao
		WHERE fatura_mes = $1 AND fatura_ano = $2 AND deleted_at IS NULL
		ORDER BY data_compra DESC
	`, mes, ano)
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
func (r *DespesaCartaoRepository) Excluir(id string) error {
	result, err := r.db.Exec(`
		UPDATE despesas_cartao
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`, id)
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
