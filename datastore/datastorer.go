package datastore

import (
    "github.com/alxeg/flibooks/models"
)

type DataStorer interface {
    PutBook(*models.Book) error
    UpdateBook(*models.Book) (*models.Book, error)
    FindBooks(title string, authors string, limit int) ([]models.Book, error)
    FindAuthors(author string, limit int) ([]models.Author, error)
    GetAuthor(authorId uint) (*models.Author, error)
    ListAuthorBooks(authorId uint, noDetails bool) ([]models.Book, error)
    GetBook(bookId uint) (*models.Book, error)
    IsContainerExist(fileName string) bool
    Close()
}
