package models

import "time"

type Item struct {
	ID        int
	Title     string
	Date      time.Time
	Content   string
	ListID    int
	CreatedAt time.Time
	UpdatedAt time.Time
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
