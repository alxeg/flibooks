package models

import (
	"regexp"

	"github.com/alxeg/flibooks/pkg/utils"
)

var re = regexp.MustCompile(`^[\s\,]*(.*?)[\,\s]*$`)

type Book struct {
	Container   Container `json:"container"`
	ContainerID uint      `json:"-"`
	Authors     []Author  `json:"authors"`
	Genres      []Genre   `json:"genres"`
	Title       string    `json:"title"`
	Series      string    `json:"series,omitempty"`
	SerNo       string    `json:"ser_no,omitempty"`
	File        string    `json:"file"`
	FileSize    string    `json:"file_size"`
	LibId       string    `json:"lib_id"`
	Del         string    `json:"del"`
	Ext         string    `json:"ext"`
	Date        string    `json:"date"`
	Lang        string    `json:"lang"`
	Extra1      string    `json:"extra_1,omitempty"`
	Extra2      string    `json:"extra_2,omitempty"`
	Extra3      string    `json:"extra_3,omitempty"`
	Update      bool      `json:"-"`
}

func (book *Book) AfterFind() {
	book.Title = utils.UpperInitialAll(book.Title)
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
