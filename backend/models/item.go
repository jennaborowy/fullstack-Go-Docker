package models

import "time"

type Item struct {
	ID      int
	Title   string
	Date    time.Time
	Content string
}
