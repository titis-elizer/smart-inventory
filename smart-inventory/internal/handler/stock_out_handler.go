package handler

import (
	"net/http"

	"smart-inventory/internal/domain"
	"smart-inventory/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StockOutHandler struct {
	service *service.StockOutService
}

func NewStockOutHandler(s *service.StockOutService) *StockOutHandler {
	return &StockOutHandler{s}
}

func (h *StockOutHandler) Create(c *gin.Context) {
	var body struct {
		Items []struct {
			InventoryItemID string `json:"inventory_item_id"`
			Qty             int    `json:"qty"`
		} `json:"items"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var items []domain.StockOutItem

	for _, i := range body.Items {
		id, err := uuid.Parse(i.InventoryItemID)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid item id"})
			return
		}

		items = append(items, domain.StockOutItem{
			InventoryItemID: id,
			Qty:             i.Qty,
		})
	}

	id, err := h.service.Create(c, items)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"id": id})
}
func (h *StockOutHandler) Allocate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Allocate(c, id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "allocated"})
}
func (h *StockOutHandler) InProgress(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.MarkInProgress(c, id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "in progress"})
}

func (h *StockOutHandler) Complete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Complete(c, id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "done"})
}
func (h *StockOutHandler) Cancel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Cancel(c, id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "canceled"})
}
func (h *StockOutHandler) GetAll(c *gin.Context) {
	data, err := h.service.GetAll(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, data)
}
