package main

import (
	"time"
)

// Playground struct
type Playground struct {
	ID         string    `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
	Articles   []Article `json:"articles"`
}

// Article struct
type Article struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Comments []Comment `json:"comments"`
}

// Comment struct
type Comment struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}
