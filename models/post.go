package models

type Post struct {
	ID      int
	Title   string
	Author  *Author
	Content string
}
