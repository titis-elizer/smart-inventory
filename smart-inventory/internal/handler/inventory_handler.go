package handler

import (
	"net/http"
	"strconv"

	"smart-inventory/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InventoryHandler struct {
	service service.InventoryService
}

func NewInventoryHandler(s service.InventoryService) *InventoryHandler {
	return &InventoryHandler{s}
}

func (h *InventoryHandler) AdjustStock(c *gin.Context) {
	var body struct {
		ItemID string `json:"item_id"`
		Qty    int    `json:"qty"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	id, _ := uuid.Parse(body.ItemID)

	err := h.service.AdjustStock(c, id, body.Qty)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, "success")
}
func (h *InventoryHandler) GetInventory(c *gin.Context) {

	search := c.Query("search")

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	data, err := h.service.GetInventory(c, search, page, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, data)
}
