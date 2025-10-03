package models

import "time"

type List struct {
	ID        int64
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Items     []Item
}

func NewList(title string, items []Item) *List {
	return &List{Title: title, Items: items}
}
