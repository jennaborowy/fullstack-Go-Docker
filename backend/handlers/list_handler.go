// handlers package processes requests through the repositories
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jennaborowy/fullstack-Go-Docker/repository"
)

// ListHandler is used to process requests related to lists
type ListHandler struct {
	repo repository.ListRepositoryInterface
}

// NewListHandler creates and returns a new ListHandler
func NewListHandler(repo repository.ListRepositoryInterface) *ListHandler {
	return &ListHandler{repo: repo}
}

// GetLists gets every list without the individual items
func (h *ListHandler) GetLists(c *gin.Context) {
	lists, err := h.repo.GetAllLists()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lists)
}

// GetList returns a single list with items fetched
func (h *ListHandler) GetList(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid list ID"})
		return
	}

	list, err := h.repo.GetList(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

// CreateList creates a new list
func (h *ListHandler) CreateList(c *gin.Context) {
	var input struct {
		Title string `json:"title"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	list, err := h.repo.CreateList(input.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, list)
}

// UpdateListTitle updates the title of an existing list
func (h *ListHandler) UpdateListTitle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid list ID"})
		return
	}

	var input struct {
		Title string `json:"title"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	updatedList, err := h.repo.UpdateTitle(id, input.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedList)
}

// DeleteList deletes a list by ID
func (h *ListHandler) DeleteList(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid list ID"})
		return
	}

	if err := h.repo.DeleteList(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
