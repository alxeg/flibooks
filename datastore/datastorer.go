package datastore

import (
    "github.com/alxeg/flibooks/models"
)

type DataStorer interface {
    PutBook(*models.Book) error
    FindBooksByTitle(title string, limit uint) ([]models.Book, error)
    FindBooksByAuthor(author string, limit uint) ([]models.Book, error)
    Close()
}
