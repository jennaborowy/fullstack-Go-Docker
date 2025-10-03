// repository package provides data access logic
package repository

import (
	"database/sql"
	"fmt"

	"github.com/jennaborowy/fullstack-Go-Docker/models"
)

// ListRepository handles CRUD operations for lists of items
type ListRepository struct {
	db *sql.DB
}

// NewListRepository creates a new ListRepository
func NewListRepository(db *sql.DB) *ListRepository {
	return &ListRepository{db: db}
}

// CreateList creates a new list and returns its ID
func (r *ListRepository) CreateList(title string) (int64, error) {
	res, err := r.db.Exec(
		"INSERT INTO lists (title) VALUES (?)",
		title,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create list: %w", err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	return id, nil
}

// GetList retrieves a list by ID with its items
func (r *ListRepository) GetList(id int) (*models.List, error) {
	list := &models.List{}
	// Get the list info
	row := r.db.QueryRow("SELECT id, title FROM lists WHERE id = ?", id)
	if err := row.Scan(&list.ID, &list.Title); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("list not found")
		}
		return nil, fmt.Errorf("failed to scan list: %w", err)
	}

	// Get items for this list
	itemsRows, err := r.db.Query("SELECT id, title, date, content, list_id FROM items WHERE list_id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}
	defer itemsRows.Close()

	for itemsRows.Next() {
		var item models.Item
		if err := itemsRows.Scan(&item.ID, &item.Title, &item.Date, &item.Content, &item.ListID); err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		list.Items = append(list.Items, item)
	}

	return list, nil
}

// GetAllLists retrieves all lists without their items
func (r *ListRepository) GetAllLists() ([]models.List, error) {
	rows, err := r.db.Query("SELECT id, title FROM lists")
	if err != nil {
		return nil, fmt.Errorf("failed to query lists: %w", err)
	}
	defer rows.Close()

	lists := []models.List{}
	for rows.Next() {
		var l models.List
		if err := rows.Scan(&l.ID, &l.Title); err != nil {
			return nil, fmt.Errorf("failed to scan list: %w", err)
		}
		lists = append(lists, l)
	}
	return lists, nil
}

// UpdateList updates the title of a list
func (r *ListRepository) UpdateTitle(id int, title string) error {
	res, err := r.db.Exec("UPDATE lists SET title = ? WHERE id = ?", title, id)
	if err != nil {
		return fmt.Errorf("failed to update list: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("list not found")
	}

	return nil
}

// DeleteList deletes a list and optionally its items
func (r *ListRepository) DeleteList(id int) error {
	// delete items first
	_, _ = r.db.Exec("DELETE FROM items WHERE list_id = ?", id)

	res, err := r.db.Exec("DELETE FROM lists WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete list: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("list not found")
	}

	return nil
}
