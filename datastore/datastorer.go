package datastore

import (
	"github.com/alxeg/flibooks/models"
)

// DataStorer interface for data layer
type DataStorer interface {
	PutBook(*models.Book) error
	UpdateBook(*models.Book) (*models.Book, error)
	FindBooks(models.Search) ([]models.Book, error)
	FindBooksSeries(models.Search) ([]models.Book, error)
	FindBooksByLibID(libID string) ([]models.Book, error)
	FindAuthors(author string, limit int) ([]models.Author, error)
	GetAuthor(authorID uint) (*models.Author, error)
	ListAuthorBooks(authorID uint, noDetails bool, params models.Search) ([]models.Book, error)
	GetBook(bookID uint) (*models.Book, error)
	GetLangs() ([]string, error)
	IsContainerExist(fileName string) bool
	Close()
}
