package models

type Container struct {
	ID       uint   `json:"-" gorm:"primary_key"`
	FileName string `json:"file_name" gorm:"not null;unique_index"`
}
