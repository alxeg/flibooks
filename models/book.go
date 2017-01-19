package models

import (
	"github.com/alxeg/flibooks/utils"
)

type Book struct {
	ID          uint      `gorm:"primary_key"`
	Container   Container `json:"container"`
	ContainerID uint      `json:"-"`
	Authors     []Author  `json:"authors" gorm:"many2many:book_authors;"`
	Genres      []Genre   `json:"genres"  gorm:"many2many:book_genres;"`
	Title       string    `json:"title" sql:"index"`
	Series      string    `json:"series,omitempty"`
	SerNo       string    `json:"ser_no,omitempty"`
	File        string    `json:"file"`
	FileSize    string    `json:"file_size"`
	LibId       string    `json:"lib_id" sql:"index"`
	Del         string    `json:"del"`
	Ext         string    `json:"ext"`
	Date        string    `json:"date"`
	Lang        string    `json:"lang"`
	Extra1      string    `json:"extra_1,omitempty"`
	Extra2      string    `json:"extra_2,omitempty"`
	Extra3      string    `json:"extra_3,omitempty"`
	Update      bool      `json:"-" sql:"-"`
}

func (book *Book) AfterFind() {
	book.Title = utils.UpperInitialAll(book.Title)
}
