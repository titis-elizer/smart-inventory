package domain

import (
	"time"

	"github.com/google/uuid"
)

type StockOut struct {
	ID     uuid.UUID
	Status string
}

type StockOutItem struct {
	ID              uuid.UUID
	StockOutID      uuid.UUID
	InventoryItemID uuid.UUID
	Qty             int
}
type StockOutWithItems struct {
	ID        uuid.UUID
	Status    string
	CreatedAt time.Time
	DoneAt    *time.Time
	Items     []StockOutItemWithProduct
}

type StockOutItemWithProduct struct {
	InventoryItemID uuid.UUID
	ProductName     string
	Qty             int
}
