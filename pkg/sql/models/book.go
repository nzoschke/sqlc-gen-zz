package models

type Book struct {
	Genres    []string `json:"genres"`
	Title     string   `json:"title"`
	Published bool     `json:"published"`
}
