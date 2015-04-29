package models

type Genre struct {
    ID        uint   `json:"-" gorm:"primary_key"`
    GenreCode string `json:"genre_code" sql:"not null;unique_index"`
}
