package models

type Author struct {
    ID   int
    Name string `json:"name" sql:"not null;unique_index"`
}
