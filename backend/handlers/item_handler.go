// handlers package processes requests through the repositories
package handlers

// Note: should probably update to do more of the JSON binding

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jennaborowy/fullstack-Go-Docker/models"
	"github.com/jennaborowy/fullstack-Go-Docker/repository"
)

// ItemHandler is used to process requests related to items
type ItemHandler struct {
	repo repository.ItemRepositoryInterface
}

// NewItemHandler creates a new ItemHandler
func NewItemHandler(repo repository.ItemRepositoryInterface) *ItemHandler {
	return &ItemHandler{repo: repo}
}

// GetItems attempts to get all items
func (h *ItemHandler) GetItems(c *gin.Context) {
	items, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetItem attempts to get a single item by id
func (h *ItemHandler) GetItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	item, err := h.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeleteItem attempts to delete an item by id and returns no content
func (h *ItemHandler) DeleteItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.repo.DeleteItemByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// CreateItem attempts to create a new item and returns its ID
func (h *ItemHandler) CreateItem(c *gin.Context) {
	var input struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		ItemDate string `json:"item_date"`
		ListID   int    `json:"list_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Bind error: %v", err) // Add this
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Received item: %+v", input)

	itemDate, err := time.Parse("2006-01-02", input.ItemDate)
	if err != nil {
		log.Printf("Bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
		return
	}

	item, err := h.repo.CreateItem(input.Title, itemDate, input.Content, input.ListID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// UpdateItem attempts to edit an item's information and returns the updated item
func (h *ItemHandler) UpdateItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		ItemDate string `json:"item_date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	date, err := time.Parse("2006-01-02", req.ItemDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
		return
	}

	// fetch item to get listID
	existingItem, err := h.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.repo.UpdateItem(id, req.Title, date, req.Content)
	if err != nil {

		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedItem := models.NewItem(req.Title, date, req.Content, existingItem.ListID)

	c.JSON(http.StatusOK, updatedItem)

}
