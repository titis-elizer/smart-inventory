package repository

import (
	"context"
	"smart-inventory/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type StockOutRepository interface {
	Create(ctx context.Context, tx pgx.Tx, stock *domain.StockOut) error
	AddItems(ctx context.Context, tx pgx.Tx, items []domain.StockOutItem) error
	UpdateStatus(ctx context.Context, tx pgx.Tx, id uuid.UUID, status string) error
	GetItems(ctx context.Context, tx pgx.Tx, stockOutID uuid.UUID) ([]domain.StockOutItem, error)
	GetStatus(ctx context.Context, tx pgx.Tx, id uuid.UUID) (string, error)
	GetAll(ctx context.Context, tx pgx.Tx) ([]domain.StockOutWithItems, error)
}
type stockOutRepository struct{}

func NewStockOutRepository() StockOutRepository {
	return &stockOutRepository{}
}
func (r *stockOutRepository) Create(ctx context.Context, tx pgx.Tx, stock *domain.StockOut) error {
	_, err := tx.Exec(ctx,
		`INSERT INTO stock_out (id, status) VALUES ($1,$2)`,
		stock.ID, stock.Status)

	return err
}

func (r *stockOutRepository) AddItems(ctx context.Context, tx pgx.Tx, items []domain.StockOutItem) error {
	for _, item := range items {
		_, err := tx.Exec(ctx,
			`INSERT INTO stock_out_items (id, stock_out_id, inventory_item_id, qty)
			 VALUES ($1,$2,$3,$4)`,
			item.ID, item.StockOutID, item.InventoryItemID, item.Qty)

		if err != nil {
			return err
		}
	}
	return nil
}

func (r *stockOutRepository) GetItems(ctx context.Context, tx pgx.Tx, stockOutID uuid.UUID) ([]domain.StockOutItem, error) {
	rows, err := tx.Query(ctx,
		`SELECT id, inventory_item_id, qty
		 FROM stock_out_items WHERE stock_out_id=$1`,
		stockOutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.StockOutItem

	for rows.Next() {
		var item domain.StockOutItem
		item.StockOutID = stockOutID

		err := rows.Scan(&item.ID, &item.InventoryItemID, &item.Qty)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *stockOutRepository) UpdateStatus(ctx context.Context, tx pgx.Tx, id uuid.UUID, status string) error {
	_, err := tx.Exec(ctx,
		`UPDATE stock_out SET status=$1 WHERE id=$2`,
		status, id)

	return err
}

func (r *stockOutRepository) GetStatus(ctx context.Context, tx pgx.Tx, id uuid.UUID) (string, error) {
	row := tx.QueryRow(ctx,
		`SELECT status FROM stock_out WHERE id=$1`, id)

	var status string
	err := row.Scan(&status)

	return status, err
}
func (r *stockOutRepository) GetAll(ctx context.Context, tx pgx.Tx) ([]domain.StockOutWithItems, error) {
	rows, err := tx.Query(ctx, `
		SELECT 
			so.id,
			so.status,
			soi.inventory_item_id,
			soi.qty,
			ii.name,
			so.created_at,
			so.done_at
		FROM stock_out so
		JOIN stock_out_items soi ON soi.stock_out_id = so.id
		JOIN inventory_items ii ON ii.id = soi.inventory_item_id
		ORDER BY so.id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resultMap := map[uuid.UUID]*domain.StockOutWithItems{}

	for rows.Next() {
		var (
			id     uuid.UUID
			status string
			itemID uuid.UUID
			qty    int
			name   string
		)

		err := rows.Scan(&id, &status, &itemID, &qty, &name)
		if err != nil {
			return nil, err
		}

		if _, ok := resultMap[id]; !ok {
			resultMap[id] = &domain.StockOutWithItems{
				ID:     id,
				Status: status,
				Items:  []domain.StockOutItemWithProduct{},
			}
		}

		resultMap[id].Items = append(resultMap[id].Items,
			domain.StockOutItemWithProduct{
				InventoryItemID: itemID,
				ProductName:     name,
				Qty:             qty,
			})
	}

	var result []domain.StockOutWithItems
	for _, v := range resultMap {
		result = append(result, *v)
	}

	return result, nil
}
