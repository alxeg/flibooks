package orm

import (
	"regexp"

	"gorm.io/gorm"

	"github.com/alxeg/flibooks/pkg/utils"
)

var re = regexp.MustCompile(`^[\s\,]*(.*?)[\,\s]*$`)

type Book struct {
	ID          uint      `gorm:"primary_key"`
	Container   Container `json:"container" gorm:"foreignKey:ContainerID"`
	ContainerID uint      `json:"-"`
	Authors     []Author  `json:"authors" gorm:"many2many:book_authors;"`
	Genres      []Genre   `json:"genres"  gorm:"many2many:book_genres;"`
	Title       string    `json:"title" gorm:"index"`
	Series      string    `json:"series,omitempty"`
	SerNo       string    `json:"ser_no,omitempty"`
	File        string    `json:"file"`
	FileSize    string    `json:"file_size"`
	LibId       string    `json:"lib_id" gorm:"index"`
	Del         string    `json:"del"`
	Ext         string    `json:"ext"`
	Date        string    `json:"date"`
	Lang        string    `json:"lang"`
	Extra1      string    `json:"extra_1,omitempty"`
	Extra2      string    `json:"extra_2,omitempty"`
	Extra3      string    `json:"extra_3,omitempty"`
	Update      bool      `json:"-" gorm:"-"`
}

func (book *Book) AfterFind(db *gorm.DB) error {
	book.Title = utils.UpperInitialAll(book.Title)
	return nil
}

func (book *Book) GetFullFilename() string {
	authors := ""
	for _, a := range book.Authors {
		authors = authors + a.Name
	}
	authors = re.ReplaceAllString(authors, "$1")
	authRunes := []rune(authors)
	if len(authRunes) > 100 {
		authors = string(authRunes[0:100]) + `…`
	}
	outName := authors + " - "
	if book.SerNo != "" && book.SerNo != "0" {
		if len(book.SerNo) == 1 {
			book.SerNo = "0" + book.SerNo
		}
		outName = outName + "[" + book.SerNo + "] "
	}

	return outName + book.Title + "." + book.Ext
}
