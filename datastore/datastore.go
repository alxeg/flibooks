package datastore

import (
	"fmt"
	_ "os"
	"strings"

	"github.com/alxeg/flibooks/models"
	"github.com/alxeg/flibooks/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type dbStore struct {
	db *gorm.DB
}

func addParams(search *gorm.DB, params models.Search) *gorm.DB {
	if !params.Deleted {
		search = search.Where("del = 0")
	}

	if len(params.Langs) > 0 {
		for i, _ := range params.Langs {
			params.Langs[i] = "'" + params.Langs[i] + "'"
		}
		search = search.Where("lang in (" + strings.Join(params.Langs, ",") + ")")
	}
	return search
}

func (store *dbStore) PutBook(book *models.Book) (err error) {
	tx := store.db.Begin()

	store.db.FirstOrCreate(&book.Container, book.Container)
	authors := []models.Author{}
	for _, author := range book.Authors {
		filledAuthor := models.Author{}
		store.db.FirstOrCreate(&filledAuthor, author)
		authors = append(authors, filledAuthor)
	}
	book.Authors = authors

	genres := []models.Genre{}
	for _, genre := range book.Genres {
		filledGenre := models.Genre{}
		store.db.FirstOrCreate(&filledGenre, genre)
		genres = append(genres, filledGenre)
	}
	book.Genres = genres

	store.db.Create(&book)

	tx.Commit()

	return err
}

func (store *dbStore) fillBookDetails(book *models.Book, fillGenres bool) {
	store.db.Select("authors.*").Model(book).Related(&book.Authors, "Authors")
	for j, a := range book.Authors {
		book.Authors[j].Name = utils.UpperInitialAll(a.Name)
	}
	if fillGenres {
		store.db.Select("genres.*").Model(book).Related(&book.Genres, "Genres")
	}
}

func (store *dbStore) fillBooksDetails(books []models.Book, fillGenres bool) []models.Book {
	for i, _ := range books {
		store.fillBookDetails(&books[i], fillGenres)
	}

	return books
}

func (store *dbStore) FindBooks(params models.Search) ([]models.Book, error) {
	title := params.Title
	authors := params.Author
	limit := params.Limit

	// utils.PrintJson(params)

	result := []models.Book{}
	search := store.db.Select("distinct books.*").Table("books").
		Joins("left join book_authors on books.id=book_authors.book_id left join authors on authors.id=book_authors.author_id")
	for _, term := range utils.SplitBySeparators(strings.ToLower(title)) {
		search = search.Where("title LIKE ?", "%"+term+"%")
	}
	for _, term := range utils.SplitBySeparators(strings.ToLower(authors)) {
		search = search.Where("name LIKE ?", "%"+term+"%")
	}

	search = addParams(search, params)

	if limit > 0 {
		search = search.Limit(limit)
	}
	search.Preload("Container").Order("title").Find(&result)

	result = store.fillBooksDetails(result, false)
	return result, nil
}

func (store *dbStore) FindBooksSeries(params models.Search) ([]models.Book, error) {
	title := params.Title
	series := params.Series
	limit := params.Limit

	result := []models.Book{}
	search := store.db.Select("distinct books.*").Table("books").
		Joins("left join book_authors on books.id=book_authors.book_id left join authors on authors.id=book_authors.author_id")
	for _, term := range utils.SplitBySeparators(strings.ToLower(title)) {
		search = search.Where("title LIKE ?", "%"+term+"%")
	}
	for _, term := range utils.SplitBySeparators(strings.ToLower(series)) {
		search = search.Where("series LIKE ?", "%"+term+"%")
	}

	search = addParams(search, params)

	if limit > 0 {
		search = search.Limit(limit)
	}
	search.Preload("Container").Order("title").Find(&result)

	result = store.fillBooksDetails(result, false)
	return result, nil
}

func (store *dbStore) FindBooksByLibID(libID string) ([]models.Book, error) {
	result := []models.Book{}
	store.db.Select("distinct books.*").Table("books").
		Where("lib_id = ?", libID).
		Find(&result)
	result = store.fillBooksDetails(result, true)
	return result, nil
}

func (store *dbStore) FindAuthors(author string, limit int) ([]models.Author, error) {
	result := []models.Author{}
	search := store.db.Order("name")
	for _, term := range utils.SplitBySeparators(strings.ToLower(author)) {
		search = search.Where("name LIKE ?", "%"+term+"%")
	}
	if limit > 0 {
		search = search.Limit(limit)
	}
	search.Find(&result)
	for i, a := range result {
		result[i].Name = utils.UpperInitialAll(a.Name)
	}
	return result, nil
}

func (store *dbStore) GetAuthor(authorID uint) (*models.Author, error) {
	result := new(models.Author)
	store.db.First(result, authorID)
	if result.ID > 0 {
		result.Name = utils.UpperInitialAll(result.Name)
		return result, nil
	}
	return nil, fmt.Errorf("No author found")
}

func (store *dbStore) ListAuthorBooks(authorID uint, noDetails bool, params models.Search) ([]models.Book, error) {
	result := []models.Book{}
	search := store.db.Select("distinct books.*").Table("books").
		Joins("left join book_authors on books.id=book_authors.book_id left join authors on authors.id=book_authors.author_id")
	search = search.Where("authors.ID=?", authorID)

	search = addParams(search, params)

	search.Preload("Container").Order("series, cast(ser_no as unsigned), title").Find(&result)
	if !noDetails {
		result = store.fillBooksDetails(result, false)
	}
	return result, nil
}

func (store *dbStore) GetBook(bookID uint) (*models.Book, error) {
	result := new(models.Book)
	store.db.Preload("Container").First(result, bookID)
	store.fillBookDetails(result, true)
	if result.ID > 0 {
		return result, nil
	}
	return nil, fmt.Errorf("No book found")
}

func (store *dbStore) UpdateBook(book *models.Book) (*models.Book, error) {
	found := new(models.Book)
	store.db.Select("distinct books.*").Table("books").
		Joins("left join containers on containers.id = books.container_id").
		Where("lib_id = ? and file_name = ?", book.LibId, book.Container.FileName).
		First(found)
	book.ID = found.ID
	book.ContainerID = found.ContainerID
	book.Container = models.Container{}
	if found != book {
		store.db.Save(book)
	}
	return book, nil
}

func (store *dbStore) GetLangs() ([]string, error) {
	var result []string
	found := []models.Book{}
	store.db.Select("distinct books.lang").
		Table("books").Where("lang <> ''").
		Order("lang").
		Find(&found)

	for _, book := range found {
		result = append(result, book.Lang)
	}
	return result, nil
}

func (store *dbStore) IsContainerExist(fileName string) bool {
	contObj := new(models.Container)
	store.db.Where("file_name = ?", fileName).First(&contObj)
	return contObj.ID > 0
}

func (store *dbStore) Close() {
}

// NewDBStore creates new instance of datastorer
func NewDBStore(config *models.DBConfig) (DataStorer, error) {
	db, err := gorm.Open(config.DBType, config.DBParams)
	if err == nil {
		db.DB()
		db.AutoMigrate(&models.Author{}, &models.Container{}, &models.Genre{}, &models.Book{})
		// db.LogMode(true)
	}
	result := new(dbStore)
	result.db = db

	return result, err
}
