package models

type Search struct {
	Title   string   `json:"title"`
	Author  string   `json:"author"`
	Limit   int      `json:"limit"`
	Deleted bool     `json:"deleted"`
	Langs   []string `json:"langs"`
}
