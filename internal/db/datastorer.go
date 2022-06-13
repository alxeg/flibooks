package db

import (
	models "github.com/alxeg/flibooks/internal/db/orm"
)

// DataStorer interface for data layer
type DataStorer interface {
	PutBook(*models.Book) error
	IsBookExist(book *models.Book) (bool, error)
	FindBooks(models.Search) ([]models.Book, error)
	FindBooksSeries(models.Search) ([]models.Book, error)
	FindBooksByLibID(libID string) ([]models.Book, error)
	FindAuthors(author string, limit int) ([]models.Author, error)
	GetAuthor(authorID uint) (*models.Author, error)
	ListAuthorBooks(authorID uint, noDetails bool, params models.Search) ([]models.Book, error)
	GetBook(bookID uint) (*models.Book, error)
	GetLangs() ([]string, error)
	IsContainerExist(fileName string) bool
}
