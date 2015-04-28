package models

type Genre struct {
    ID        int    `json:"-"`
    GenreCode string `json:"genre_code" sql:"not null;unique_index"`
}
