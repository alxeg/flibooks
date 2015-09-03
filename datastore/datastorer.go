package datastore

import (
    "github.com/alxeg/flibooks/models"
)

type DataStorer interface {
    PutBook(*models.Book) error
    UpdateBook(*models.Book) (*models.Book, error)
    FindBooks(models.Search) ([]models.Book, error)
    FindBooksByLibId(libId string) ([]models.Book, error)
    FindAuthors(author string, limit int) ([]models.Author, error)
    GetAuthor(authorId uint) (*models.Author, error)
    ListAuthorBooks(authorId uint, noDetails bool, params models.Search) ([]models.Book, error)
    GetBook(bookId uint) (*models.Book, error)
    GetLangs() ([]string, error)
    IsContainerExist(fileName string) bool
    Close()
}
