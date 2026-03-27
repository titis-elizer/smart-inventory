package domain

import "github.com/google/uuid"

type InventoryItem struct {
	ID            uuid.UUID
	Name          string
	SKU           string
	Customer      *string
	PhysicalStock int
	ReservedStock int
}
