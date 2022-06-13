package db

import (
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	models "github.com/alxeg/flibooks/internal/db/orm"
	"github.com/alxeg/flibooks/pkg/utils"
)

type dbStore struct {
	db *gorm.DB
}

func addParams(search *gorm.DB, params models.Search) *gorm.DB {
	if !params.Deleted {
		search = search.Where("del = 0")
	}

	if len(params.Langs) > 0 {
		for i := range params.Langs {
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

func (store *dbStore) FindBooks(params models.Search) ([]models.Book, error) {
	title := params.Title
	authors := params.Author
	limit := params.Limit

	// utils.PrintJson(params)

	result := []models.Book{}
	search := store.db.Select("books.*").Table("books").
		Joins("left join book_authors on books.id=book_authors.book_id left join authors on authors.id=book_authors.author_id")
	for _, term := range utils.SplitBySeparators(strings.ToLower(title)) {
		search = search.Where("title LIKE ?", "%"+term+"%")
	}
	for _, term := range utils.SplitBySeparators(strings.ToLower(authors)) {
		search = search.Where("name LIKE ?", "%"+term+"%")
	}

	search = addParams(search, params).Group("books.id")

	if limit > 0 {
		search = search.Limit(limit)
	}
	search.Preload("Container").Preload("Authors").Order("title").Find(&result)

	// result = store.fillBooksDetails(result, false)
	return result, nil
}

func (store *dbStore) FindBooksSeries(params models.Search) ([]models.Book, error) {
	title := params.Title
	series := params.Series
	limit := params.Limit

	result := []models.Book{}
	search := store.db.Select("books.*").Table("books").
		Joins("left join book_authors on books.id=book_authors.book_id left join authors on authors.id=book_authors.author_id")
	for _, term := range utils.SplitBySeparators(strings.ToLower(title)) {
		search = search.Where("title LIKE ?", "%"+term+"%")
	}
	for _, term := range utils.SplitBySeparators(strings.ToLower(series)) {
		search = search.Where("LOWER(series) LIKE ?", "%"+term+"%")
	}

	search = addParams(search, params).Group("books.id")

	if limit > 0 {
		search = search.Limit(limit)
	}
	search.Preload("Container").Order("series, cast(ser_no as unsigned), title").Find(&result)

	return result, nil
}

func (store *dbStore) FindBooksByLibID(libID string) ([]models.Book, error) {
	result := []models.Book{}
	store.db.Select("books.*").Table("books").
		Preload("Authors").
		Preload("Genres").
		Where("lib_id = ?", libID).
		Find(&result)

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
	return nil, errors.New("no author found")
}

func (store *dbStore) ListAuthorBooks(authorID uint, noDetails bool, params models.Search) ([]models.Book, error) {
	result := []models.Book{}
	search := store.db.Select("books.*").Table("books").
		Joins("left join book_authors on books.id=book_authors.book_id left join authors on authors.id=book_authors.author_id")
	search = search.Where("authors.ID=?", authorID)

	search = addParams(search, params).Group("books.id")

	search = search.Preload("Container")
	if !noDetails {
		search = search.Preload("Authors")
	}
	search.Order("series, cast(ser_no as unsigned), title").Find(&result)
	return result, nil
}

func (store *dbStore) GetBook(bookID uint) (*models.Book, error) {
	result := new(models.Book)
	store.db.Preload("Container").Preload("Authors").Preload("Genres").First(result, bookID)
	if result.ID > 0 {
		return result, nil
	}
	return nil, errors.New("no book found")
}

func (store *dbStore) IsBookExist(book *models.Book) (bool, error) {
	found := new(models.Book)
	res := store.db.Select("distinct books.*").Table("books").
		Joins("left join containers on containers.id = books.container_id").
		Where("lib_id = ? and file_name = ?", book.LibId, book.Container.FileName).
		First(found)

	return res.RowsAffected != 0, res.Error
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

// NewDBStore creates new instance of datastorer
func NewDBStore(dbType, connect, logLevel string) (DataStorer, error) {
	var dialector gorm.Dialector
	switch dbType {
	case "mysql":
		dialector = mysql.Open(connect)
	case "sqlite":
		dialector = sqlite.Open(connect)
	default:
		return nil, errors.New("unknown dbType")
	}

	logMap := map[string]logger.LogLevel{
		"Silent": logger.Silent,
		"Error":  logger.Error,
		"Warn":   logger.Warn,
		"Info":   logger.Info,
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second,      // Slow SQL threshold
				LogLevel:                  logMap[logLevel], // Log level
				IgnoreRecordNotFoundError: true,             // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,            // Disable color
			},
		),
	})

	if err == nil {
		db.DB()
		db.AutoMigrate(&models.Author{}, &models.Container{}, &models.Genre{}, &models.Book{})
		// db.LogMode(true)
	}
	result := new(dbStore)
	result.db = db

	return result, err
}
