package repository

import (
	"context"
	"smart-inventory/internal/domain"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type StockInRepository interface {
	Create(ctx context.Context, tx pgx.Tx, stock *domain.StockIn) error
	AddItems(ctx context.Context, tx pgx.Tx, items []domain.StockInItem) error
	UpdateStatus(ctx context.Context, tx pgx.Tx, id uuid.UUID, status string) error
	GetItems(ctx context.Context, tx pgx.Tx, stockInID uuid.UUID) ([]domain.StockInItem, error)
	AddLog(ctx context.Context, tx pgx.Tx, stockInID uuid.UUID, status string, note string) error
	GetAll(ctx context.Context, tx pgx.Tx) ([]domain.StockInWithItems, error)
}
type stockInRepository struct{}

func NewStockInRepository() StockInRepository {
	return &stockInRepository{}
}

func (r *stockInRepository) Create(ctx context.Context, tx pgx.Tx, stock *domain.StockIn) error {
	_, err := tx.Exec(ctx,
		`INSERT INTO stock_in (id, status) VALUES ($1,$2)`,
		stock.ID, stock.Status)
	return err
}

func (r *stockInRepository) AddItems(ctx context.Context, tx pgx.Tx, items []domain.StockInItem) error {
	for _, item := range items {
		_, err := tx.Exec(ctx,
			`INSERT INTO stock_in_items (id, stock_in_id, inventory_item_id, qty)
			 VALUES ($1,$2,$3,$4)`,
			item.ID, item.StockInID, item.InventoryItemID, item.Qty)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *stockInRepository) UpdateStatus(ctx context.Context, tx pgx.Tx, id uuid.UUID, status string) error {
	_, err := tx.Exec(ctx,
		`UPDATE stock_in SET status=$1 WHERE id=$2`,
		status, id)
	return err
}

func (r *stockInRepository) GetItems(ctx context.Context, tx pgx.Tx, stockInID uuid.UUID) ([]domain.StockInItem, error) {
	rows, err := tx.Query(ctx,
		`SELECT id, inventory_item_id, qty FROM stock_in_items WHERE stock_in_id=$1`,
		stockInID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.StockInItem
	for rows.Next() {
		var item domain.StockInItem
		item.StockInID = stockInID
		rows.Scan(&item.ID, &item.InventoryItemID, &item.Qty)
		items = append(items, item)
	}
	return items, nil
}

func (r *stockInRepository) AddLog(ctx context.Context, tx pgx.Tx, stockInID uuid.UUID, status string, note string) error {
	_, err := tx.Exec(ctx,
		`INSERT INTO stock_in_logs (id, stock_in_id, status, note)
		 VALUES (gen_random_uuid(), $1, $2, $3)`,
		stockInID, status, note)
	return err
}
func (r *stockInRepository) GetAll(ctx context.Context, tx pgx.Tx) ([]domain.StockInWithItems, error) {
	rows, err := tx.Query(ctx, `
	SELECT 
		si.id,
		si.status,
		si.created_at,
		si.done_at,
		sii.inventory_item_id,
		sii.qty,
		ii.name
	FROM stock_in si
	JOIN stock_in_items sii ON sii.stock_in_id = si.id
	JOIN inventory_items ii ON ii.id = sii.inventory_item_id
	ORDER BY si.id DESC
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resultMap := map[uuid.UUID]*domain.StockInWithItems{}

	for rows.Next() {
		var (
			id        uuid.UUID
			status    string
			createdAt time.Time
			doneAt    *time.Time
			itemID    uuid.UUID
			qty       int
			name      string
		)

		err := rows.Scan(
			&id,
			&status,
			&createdAt,
			&doneAt,
			&itemID,
			&qty,
			&name,
		)
		if err != nil {
			return nil, err
		}

		if _, ok := resultMap[id]; !ok {
			resultMap[id] = &domain.StockInWithItems{
				ID:        id,
				Status:    status,
				CreatedAt: createdAt,
				DoneAt:    doneAt,
				Items:     []domain.StockInItemWithProduct{},
			}
		}

		resultMap[id].Items = append(resultMap[id].Items,
			domain.StockInItemWithProduct{
				InventoryItemID: itemID,
				ProductName:     name,
				Qty:             qty,
			})
	}

	var result []domain.StockInWithItems
	for _, v := range resultMap {
		result = append(result, *v)
	}

	return result, nil
}
