package service

import (
	"context"
	"errors"
	"smart-inventory/internal/domain"
	"smart-inventory/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryService interface {
	AdjustStock(ctx context.Context, itemID uuid.UUID, qty int) error
	GetInventory(ctx context.Context, search string, page, limit int) ([]domain.InventoryItem, error)
}

type inventoryService struct {
	repo repository.InventoryRepository
	db   *pgxpool.Pool
}

func NewInventoryService(r repository.InventoryRepository, db *pgxpool.Pool) InventoryService {
	return &inventoryService{
		repo: r,
		db:   db,
	}
}
func (s *inventoryService) AdjustStock(ctx context.Context, itemID uuid.UUID, qty int) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 🔒 ambil + lock row
	item, err := s.repo.FindByIDForUpdate(ctx, tx, itemID)
	if err != nil {
		return err
	}

	// validasi
	if item.PhysicalStock+qty < 0 {
		return errors.New("stock cannot be negative")
	}

	// update
	item.PhysicalStock += qty

	if err := s.repo.UpdateTx(ctx, tx, item); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *inventoryService) GetInventory(ctx context.Context, search string, page, limit int) ([]domain.InventoryItem, error) {

	offset := (page - 1) * limit

	return s.repo.FindAll(ctx, search, limit, offset)
}
