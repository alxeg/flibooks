package models

import (
    "github.com/alxeg/flibooks/utils"
)

//AUTHOR;                    GENRE;           TITLE;
//Крыжановская,Вера,Ивановна:prose_history:Два сфинкса
//SERIES;SERNO;FILE;    SIZE;      LIBID;   DEL; EXT;  DATE;        LANG;N1 N2 N3
//       18309212416361830920fb22010-03-01ru

type Book struct {
    ID        int
    Container Container `json:"container"`
    Authors   []Author  `json:"authors" gorm:"many2many:book_authors;"`
    Genres    []Genre   `json:"genres"  gorm:"many2many:book_genres;"`
    Title     string    `json:"title"`
    Series    string    `json:"series"`
    SerNo     string    `json:"ser_no"`
    File      string    `json:"file"`
    FileSize  string    `json:"file_size"`
    LibId     string    `json:"lib_id"`
    Del       string    `json:"del"`
    Ext       string    `json:"ext"`
    Date      string    `json:"date"`
    Lang      string    `json:"lang"`
    Extra1    string    `json:"extra_1"`
    Extra2    string    `json:"extra_2"`
    Extra3    string    `json:"extra_3"`
}

func (book *Book) AfterFind() {
    book.Title = utils.UpperInitialAll(book.Title)
}
