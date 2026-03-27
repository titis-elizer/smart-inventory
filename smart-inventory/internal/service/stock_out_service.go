package service

import (
	"context"
	"errors"

	"smart-inventory/internal/domain"
	"smart-inventory/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StockOutService struct {
	db      *pgxpool.Pool
	repo    repository.StockOutRepository
	invRepo repository.InventoryRepository
}

func NewStockOutService(db *pgxpool.Pool, r repository.StockOutRepository, inv repository.InventoryRepository) *StockOutService {
	return &StockOutService{db, r, inv}
}
func (s *StockOutService) Create(ctx context.Context, items []domain.StockOutItem) (uuid.UUID, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	id := uuid.New()

	// 🔥 langsung allocated
	stock := &domain.StockOut{
		ID:     id,
		Status: "allocated",
	}

	if err := s.repo.Create(ctx, tx, stock); err != nil {
		return uuid.Nil, err
	}

	for i := range items {
		items[i].StockOutID = id
		items[i].ID = uuid.New()
	}

	// 🔒 VALIDASI + RESERVE SEKALIGUS
	for _, item := range items {

		inv, err := s.invRepo.FindByIDForUpdate(ctx, tx, item.InventoryItemID)
		if err != nil {
			return uuid.Nil, err
		}

		available := inv.PhysicalStock - inv.ReservedStock

		if item.Qty > available {
			return uuid.Nil, errors.New("insufficient stock")
		}

		inv.ReservedStock += item.Qty

		if err := s.invRepo.UpdateTx(ctx, tx, inv); err != nil {
			return uuid.Nil, err
		}
	}

	if err := s.repo.AddItems(ctx, tx, items); err != nil {
		return uuid.Nil, err
	}

	return id, tx.Commit(ctx)
}
func (s *StockOutService) Allocate(ctx context.Context, id uuid.UUID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	items, err := s.repo.GetItems(ctx, tx, id)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return errors.New("no items found")
	}

	for _, item := range items {

		// 🔒 lock inventory row
		inv, err := s.invRepo.FindByIDForUpdate(ctx, tx, item.InventoryItemID)
		if err != nil {
			return err
		}

		available := inv.PhysicalStock - inv.ReservedStock

		if available < item.Qty {
			return errors.New("insufficient stock")
		}

		inv.ReservedStock += item.Qty

		if err := s.invRepo.UpdateTx(ctx, tx, inv); err != nil {
			return err
		}
	}

	if err := s.repo.UpdateStatus(ctx, tx, id, "allocated"); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
func (s *StockOutService) MarkInProgress(ctx context.Context, id uuid.UUID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.repo.UpdateStatus(ctx, tx, id, "in_progress"); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
func (s *StockOutService) Complete(ctx context.Context, id uuid.UUID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	items, err := s.repo.GetItems(ctx, tx, id)
	if err != nil {
		return err
	}

	for _, item := range items {

		inv, err := s.invRepo.FindByIDForUpdate(ctx, tx, item.InventoryItemID)
		if err != nil {
			return err
		}

		// 🔥 FINAL MOVEMENT
		inv.PhysicalStock -= item.Qty
		inv.ReservedStock -= item.Qty

		if inv.PhysicalStock < 0 || inv.ReservedStock < 0 {
			return errors.New("invalid stock state")
		}

		if err := s.invRepo.UpdateTx(ctx, tx, inv); err != nil {
			return err
		}
	}

	if err := s.repo.UpdateStatus(ctx, tx, id, "done"); err != nil {
		return err
	}
	tx.Exec(ctx,
		`UPDATE stock_out SET status=$1, done_at=now() WHERE id=$2`,
		"done", id)

	return tx.Commit(ctx)
}
func (s *StockOutService) Cancel(ctx context.Context, id uuid.UUID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	items, err := s.repo.GetItems(ctx, tx, id)
	if err != nil {
		return err
	}

	for _, item := range items {

		inv, err := s.invRepo.FindByIDForUpdate(ctx, tx, item.InventoryItemID)
		if err != nil {
			return err
		}

		// 🔥 rollback reservation
		inv.ReservedStock -= item.Qty

		if inv.ReservedStock < 0 {
			inv.ReservedStock = 0
		}

		if err := s.invRepo.UpdateTx(ctx, tx, inv); err != nil {
			return err
		}
	}

	if err := s.repo.UpdateStatus(ctx, tx, id, "canceled"); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
func isValidTransition(current, next string) bool {
	valid := map[string][]string{
		"draft":       {"allocated", "canceled"},
		"allocated":   {"in_progress", "canceled"},
		"in_progress": {"done", "canceled"},
	}

	for _, v := range valid[current] {
		if v == next {
			return true
		}
	}
	return false
}
func (s *StockOutService) GetAll(ctx context.Context) ([]domain.StockOutWithItems, error) {
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
