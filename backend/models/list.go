package models

type List struct {
	ID    int64
	Title string
	Items []Item
}

func NewList(title string, items []Item) *List {
	return &List{Title: title, Items: items}
}
