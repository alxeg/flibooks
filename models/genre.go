package models

type Genre struct {
    ID        int
    GenreCode string `json:"genre_code" sql:"not null;unique_index"`
}
