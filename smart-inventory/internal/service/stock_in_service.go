package service

import (
	"context"
	"errors"
	"fmt"

	"smart-inventory/internal/domain"
	"smart-inventory/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StockInService struct {
	db      *pgxpool.Pool
	repo    repository.StockInRepository
	invRepo repository.InventoryRepository
}

func NewStockInService(db *pgxpool.Pool, r repository.StockInRepository, inv repository.InventoryRepository) *StockInService {
	return &StockInService{db, r, inv}
}
func (s *StockInService) Create(ctx context.Context, items []domain.StockInItem) (uuid.UUID, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	id := uuid.New()

	stock := &domain.StockIn{
		ID:     id,
		Status: "in_progress",
	}

	if err := s.repo.Create(ctx, tx, stock); err != nil {
		return uuid.Nil, err
	}

	for i := range items {
		items[i].StockInID = id
		items[i].ID = uuid.New()
	}

	if err := s.repo.AddItems(ctx, tx, items); err != nil {
		return uuid.Nil, err
	}

	if err := s.repo.AddLog(ctx, tx, id, "created", "Stock In created"); err != nil {
		return uuid.Nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
func (s *StockInService) UpdateStatus(ctx context.Context, id uuid.UUID, newStatus string) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 🔒 ambil items dalam tx
	items, err := s.repo.GetItems(ctx, tx, id)
	if err != nil {
		return err
	}

	// 🔥 VALIDASI STATUS (optional tapi disarankan)
	valid := map[string]bool{
		"created":     true,
		"in_progress": true,
		"done":        true,
		"canceled":    true,
	}
	if !valid[newStatus] {
		return errors.New("invalid status")
	}

	// 🔥 jika DONE → update inventory (SAFE)
	if newStatus == "done" {
		for _, item := range items {

			fmt.Println("PROCESS ITEM:", item.InventoryItemID)

			inv, err := s.invRepo.FindByIDForUpdate(ctx, tx, item.InventoryItemID)
			if err != nil {
				fmt.Println("ERROR FIND:", err)
				return err
			}

			fmt.Println("BEFORE:", inv.PhysicalStock)

			inv.PhysicalStock += item.Qty

			fmt.Println("AFTER:", inv.PhysicalStock)

			if err := s.invRepo.UpdateTx(ctx, tx, inv); err != nil {
				fmt.Println("ERROR UPDATE:", err)
				return err
			}
		}
		tx.Exec(ctx,
			`UPDATE stock_in SET status=$1, done_at=now() WHERE id=$2`,
			newStatus, id)
	}

	// ❌ prevent cancel jika sudah done (recommended)
	if newStatus == "canceled" {
		// bisa tambahkan check status sebelumnya di sini
	}

	if err := s.repo.UpdateStatus(ctx, tx, id, newStatus); err != nil {
		return err
	}

	if err := s.repo.AddLog(ctx, tx, id, newStatus, "status updated"); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
func (s *StockInService) GetAll(ctx context.Context) ([]domain.StockInWithItems, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	data, err := s.repo.GetAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	return data, tx.Commit(ctx)
}
