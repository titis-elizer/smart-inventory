package domain

import (
	"time"

	"github.com/google/uuid"
)

type StockIn struct {
	ID     uuid.UUID
	Status string
}

type StockInItem struct {
	ID              uuid.UUID
	StockInID       uuid.UUID
	InventoryItemID uuid.UUID
	Qty             int
}
type StockInWithItems struct {
	ID        uuid.UUID
	Status    string
	CreatedAt time.Time
	DoneAt    *time.Time
	Items     []StockInItemWithProduct
}

type StockInItemWithProduct struct {
	InventoryItemID uuid.UUID
	ProductName     string
	Qty             int
}
type StockInCreateItem struct {
	InventoryItemID string
	ProductName     string
	Qty             int
}
