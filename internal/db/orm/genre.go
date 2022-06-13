package orm

type Genre struct {
	ID        uint   `json:"-" gorm:"primary_key"`
	GenreCode string `json:"genre_code" gorm:"not null;unique_index"`
}
