// repository package provides data access logic
package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jennaborowy/fullstack-Go-Docker/models"
)

type ListRepositoryInterface interface {
	CreateList(title string) (*models.List, error)
	GetList(id int) (*models.List, error)
	GetAllLists() ([]models.List, error)
	UpdateTitle(id int, title string) (*models.List, error)
	DeleteList(id int) error
}

// ListRepository handles CRUD operations for lists of items
type ListRepository struct {
	db *sql.DB
}

// NewListRepository creates a new ListRepository
func NewListRepository(db *sql.DB) *ListRepository {
	return &ListRepository{db: db}
}

// CreateList creates a new list and returns its ID
func (r *ListRepository) CreateList(title string) (*models.List, error) {
	list := &models.List{}
	err := r.db.QueryRow(
		"INSERT INTO lists (title) VALUES ($1) RETURNING id, title, created_at, updated_at",
		title,
	).Scan(&list.ID, &list.Title, &list.CreatedAt, &list.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("could not obtain new id: %w", err)
	}
	return list, nil
}

// GetList retrieves a list by ID with its items
func (r *ListRepository) GetList(id int) (*models.List, error) {
	list := &models.List{}
	// Get the list info
	row := r.db.QueryRow("SELECT id, title, created_at, updated_at FROM lists WHERE id = $1", id)
	if err := row.Scan(&list.ID, &list.Title, &list.CreatedAt, &list.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("list not found")
		}
		return nil, fmt.Errorf("failed to scan list: %w", err)
	}

	// Get items for this list
	itemsRows, err := r.db.Query("SELECT id, title, item_date, content, list_id, created_at, updated_at FROM items WHERE list_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}
	defer itemsRows.Close()

	for itemsRows.Next() {
		var item models.Item
		if err := itemsRows.Scan(&item.ID, &item.Title, &item.Date, &item.Content, &item.ListID, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		list.Items = append(list.Items, item)
	}

	// could instead change this to do a left join for a single query, rather than two

	return list, nil
}

// GetAllLists retrieves all lists without their items
func (r *ListRepository) GetAllLists() ([]models.List, error) {
	rows, err := r.db.Query("SELECT id, title, created_at, updated_at FROM lists")
	if err != nil {
		return nil, fmt.Errorf("failed to query lists: %w", err)
	}
	defer rows.Close()

	lists := []models.List{}
	for rows.Next() {
		var l models.List
		if err := rows.Scan(&l.ID, &l.Title, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan list: %w", err)
		}
		lists = append(lists, l)
	}
	return lists, nil
}

// UpdateList updates the title of a list
func (r *ListRepository) UpdateTitle(id int, title string) (*models.List, error) {
	res, err := r.db.Exec("UPDATE lists SET title = $1, updated_at = $2 WHERE id = $3", title, time.Now(), id)
	if err != nil {
		return nil, fmt.Errorf("failed to update list: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("list not found")
	}

	// Query the updated list
	list := &models.List{}
	row := r.db.QueryRow("SELECT id, title, created_at, updated_at FROM lists WHERE id = $1", id)
	if err := row.Scan(&list.ID, &list.Title, &list.CreatedAt, &list.UpdatedAt); err != nil {
		return nil, fmt.Errorf("failed to fetch updated list: %w", err)
	}

	return list, nil
}

// DeleteList deletes a list and optionally its items
func (r *ListRepository) DeleteList(id int) error {
	// delete items first
	_, _ = r.db.Exec("DELETE FROM items WHERE list_id = $1", id)

	res, err := r.db.Exec("DELETE FROM lists WHERE id = $1", id)
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
