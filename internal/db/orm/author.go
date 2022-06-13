package orm

type Author struct {
	ID   uint   `gorm:"primary_key"`
	Name string `json:"name" gorm:"not null;unique_index"`
}
