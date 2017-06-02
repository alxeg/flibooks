package models

// Search data model
type Search struct {
	Title   string   `json:"title"`
	Author  string   `json:"author"`
	Series  string   `json:"series"`
	Limit   int      `json:"limit"`
	Deleted bool     `json:"deleted"`
	Langs   []string `json:"langs"`
}
