package repository

import (
	"context"
	"smart-inventory/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryRepository interface {
	// read biasa (tanpa lock)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.InventoryItem, error)

	// 🔥 read + lock
	FindByIDForUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*domain.InventoryItem, error)

	// update dalam transaction
	UpdateTx(ctx context.Context, tx pgx.Tx, item *domain.InventoryItem) error
	FindAll(ctx context.Context, search string, limit, offset int) ([]domain.InventoryItem, error)
}

type inventoryRepository struct {
	db *pgxpool.Pool
}

func NewInventoryRepository(db *pgxpool.Pool) InventoryRepository {
	return &inventoryRepository{db}
}

func (r *inventoryRepository) FindAll(ctx context.Context, search string, limit, offset int) ([]domain.InventoryItem, error) {

	query := `
	SELECT id, name, sku, customer, physical_stock, reserved_stock
	FROM test_326.inventory_items
	WHERE ($1 = '' OR name ILIKE '%' || $1 || '%' 
		OR sku ILIKE '%' || $1 || '%' 
		OR customer ILIKE '%' || $1 || '%')
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, search, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.InventoryItem

	for rows.Next() {
		var item domain.InventoryItem

		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.SKU,
			&item.Customer,
			&item.PhysicalStock,
			&item.ReservedStock,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *inventoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.InventoryItem, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, name, sku, customer, physical_stock, reserved_stock
		FROM test_326.inventory_items
		WHERE id=$1`, id)

	var item domain.InventoryItem
	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.SKU,
		&item.Customer,
		&item.PhysicalStock,
		&item.ReservedStock,
	)

	if err != nil {
		return nil, err
	}

	return &item, nil
}
func (r *inventoryRepository) FindByIDForUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*domain.InventoryItem, error) {
	row := tx.QueryRow(ctx, `
		SELECT id, name, sku, customer, physical_stock, reserved_stock
		FROM test_326.inventory_items
		WHERE id=$1
		FOR UPDATE`, id)

	var item domain.InventoryItem
	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.SKU,
		&item.Customer,
		&item.PhysicalStock,
		&item.ReservedStock,
	)

	if err != nil {
		return nil, err
	}

	return &item, nil
}
func (r *inventoryRepository) UpdateTx(ctx context.Context, tx pgx.Tx, item *domain.InventoryItem) error {
	_, err := tx.Exec(ctx, `
		UPDATE test_326.inventory_items
		SET physical_stock=$1,
			reserved_stock=$2
		WHERE id=$3`,
		item.PhysicalStock,
		item.ReservedStock,
		item.ID,
	)

	return err
}
