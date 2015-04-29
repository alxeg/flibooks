package datastore

import (
    "github.com/alxeg/flibooks/models"
)

type DataStorer interface {
    PutBook(*models.Book) error
    FindBooks(title string, authors string, limit uint) ([]models.Book, error)
    FindAuthors(author string, limit uint) ([]models.Author, error)
    Close()
}
