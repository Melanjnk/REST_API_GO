package models

import "fmt"

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

func (b *Book) String() string {
	return fmt.Sprintf("ID: %d, Title: %s, Author: %s", b.ID, b.Title, b.Author)
}
