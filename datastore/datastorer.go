package datastore

import (
    "github.com/alxeg/flibooks/models"
)

type DataStorer interface {
    PutBook(models.Book) error
    FindBooksByTitle(title string) ([]models.Book, error)
    FindBooksByAuthor(author string) ([]models.Book, error)
    Close()
}
