package models

type Container struct {
    ID       int    `json:"-"`
    FileName string `json:"file_name" sql:"not null;unique_index"`
}
