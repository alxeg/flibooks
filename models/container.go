package models

type Container struct {
    ID       uint   `json:"-" gorm:"primary_key"`
    FileName string `json:"file_name" sql:"not null;unique_index"`
}
