package handler

import (
	"smart-inventory/internal/domain"
	"smart-inventory/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StockInHandler struct {
	service *service.StockInService
}

func NewStockInHandler(s *service.StockInService) *StockInHandler {
	return &StockInHandler{s}
}
func (h *StockInHandler) Create(c *gin.Context) {
	var body struct {
		Items []struct {
			InventoryItemID string `json:"inventory_item_id"`
			Qty             int    `json:"qty"`
		} `json:"items"`
	}

	c.ShouldBindJSON(&body)

	var items []domain.StockInItem

	for _, i := range body.Items {
		id, _ := uuid.Parse(i.InventoryItemID)

		items = append(items, domain.StockInItem{
			InventoryItemID: id,
			Qty:             i.Qty,
		})
	}

	id, err := h.service.Create(c, items)

	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	c.JSON(200, gin.H{"id": id})
}
func (h *StockInHandler) UpdateStatus(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))

	var body struct {
		Status string
	}

	c.ShouldBindJSON(&body)

	err := h.service.UpdateStatus(c, id, body.Status)

	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	c.JSON(200, "updated")
}
func (h *StockInHandler) GetAll(c *gin.Context) {
	data, err := h.service.GetAll(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, data)
}
