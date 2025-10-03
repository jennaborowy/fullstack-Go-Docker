package models

import "time"

type Item struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Date      time.Time `json:"item_date" time_format:"2006-01-02"`
	Content   string    `json:"content"`
	ListID    int       `json:"list_id"`
	CreatedAt time.Time `json:"created_at" time_format:"2006-01-02"`
	UpdatedAt time.Time `json:"updated_at" time_format:"2006-01-02"`
}

// NewItem creates a new item
func NewItem(title string, date time.Time, content string, listID int) *Item {
	return &Item{
		Title:   title,
		Date:    date,
		Content: content,
		ListID:  listID,
	}
}
