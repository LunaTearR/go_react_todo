package models

type Todo struct {
	ID    int    `json:"id" db:"id"`
	Title string `json:"title" db:"title"`
	Done  bool   `json:"done" db:"done"`
	Body  string `json:"body" db:"body"`
}