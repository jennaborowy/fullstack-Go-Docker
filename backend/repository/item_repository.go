// repository package provides data access logic
package repository

import (
	"database/sql"
	"time"

	"github.com/jennaborowy/fullstack-Go-Docker/models"
)

// ItemRepository handles CRUD operations for items
type ItemRepository struct {
	db *sql.DB
}

// NewItemRepository creates a new ItemRepository
func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

// GetAll retrieves all existing items from database
func (r *ItemRepository) GetAll() ([]models.Item, error) {
	rows, err := r.db.Query("SELECT id, title, date, content FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ID, &item.Title, &item.Date, &item.Content); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil

}

// DeleteItemByID deletes an item by ID
func (r *ItemRepository) DeleteItemByID(id int) {

}

// CreateItem creates a new item with title, date, and content
func (r *ItemRepository) CreateItem(title string, date time.Time, content string) {

}

// UpdateItem updates an item's title, date, and/or content
func (r *ItemRepository) UpdateItem(title string, date time.Time, content string) {

}
