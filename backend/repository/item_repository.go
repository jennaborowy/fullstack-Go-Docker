// repository package provides data access logic
package repository

import (
	"database/sql"
	"fmt"
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
	rows, err := r.db.Query("SELECT id, title, date, content, list_id, created_at, updated_at FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ID, &item.Title, &item.Date, &item.Content, &item.ListID, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil

}

// GetByID retrieves a single item by its ID
func (r *ItemRepository) GetByID(id int) (*models.Item, error) {
	row := r.db.QueryRow("SELECT id, title, date, content, list_id, created_at, updated_at FROM items WHERE id = ?", id)

	var item models.Item
	err := row.Scan(&item.ID, &item.Title, &item.Date, &item.Content, &item.ListID, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("item not found")
		}
		return nil, err
	}

	return &item, nil
}

// GetItemsByListID gets an item by the list it is in
func (r *ItemRepository) GetByListID(listID int) (*models.Item, error) {
	row := r.db.QueryRow("SELECT id, title, date, content, list_id, created_at, updated_at FROM items WHERE list_id = ?", listID)
	var item models.Item
	err := row.Scan(&item.ID, &item.Title, &item.Date, &item.Content, &item.ListID, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("item not found")
		}
		return nil, err
	}

	return &item, nil
}

// DeleteItemByID deletes an item by ID
func (r *ItemRepository) DeleteItemByID(id int) error {
	res, err := r.db.Exec("DELETE FROM items WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no item found with id %d", id)
	}

	return nil
}

// CreateItem creates a new item with title, date, and content
func (r *ItemRepository) CreateItem(title string, date time.Time, content string, listID int) (int64, error) {
	res, err := r.db.Exec("INSERT INTO items (title, date, content, list_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)", title, date, content, listID, time.Now(), time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to create new item: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("could not obtain new id: %d", id)
	}

	return id, nil
}

// UpdateItem updates an item's title, date, and/or content
func (r *ItemRepository) UpdateItem(id int, title string, date time.Time, content string) error {
	res, err := r.db.Exec("UPDATE items Set title = ?, date = ?, content = ?, updated_at = ? WHERE id = ?", title, date, content, time.Now(), id)
	if err != nil {
		return fmt.Errorf("could not update item: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to check affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no item found with id %d", id)
	}

	return nil
}
