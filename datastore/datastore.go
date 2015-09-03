package datastore

import (
    "fmt"
    "github.com/alxeg/flibooks/models"
    "github.com/alxeg/flibooks/utils"
    _ "github.com/go-sql-driver/mysql"
    "github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"
    _ "os"
    "strings"
)

type dbStore struct {
    db gorm.DB
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

    result := []models.Book{}
    search := store.db.Select("distinct books.*").Table("books").
        Joins("left join book_authors on books.id=book_authors.book_id left join authors on authors.id=book_authors.author_id")
    for _, term := range utils.SplitBySeparators(strings.ToLower(title)) {
        search = search.Where("title LIKE ?", "%"+term+"%")
    }
    for _, term := range utils.SplitBySeparators(strings.ToLower(authors)) {
        search = search.Where("name LIKE ?", "%"+term+"%")
    }
    if !params.Deleted {
        search.Where("del = 0")
    }

    if len(params.Langs) > 0 {
        search.Where("lang in (" + strings.Join(params.Langs, ",") + ")")
    }

    if limit > 0 {
        search = search.Limit(limit)
    }
    search.Preload("Container").Order("title").Find(&result)

    result = store.fillBooksDetails(result, false)
    return result, nil
}

func (store *dbStore) FindBooksByLibId(libId string) ([]models.Book, error) {
    result := []models.Book{}
    store.db.Select("distinct books.*").Table("books").
        Where("lib_id = ?", libId).
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

func (store *dbStore) GetAuthor(authorId uint) (*models.Author, error) {
    result := new(models.Author)
    store.db.First(result, authorId)
    if result.ID > 0 {
        result.Name = utils.UpperInitialAll(result.Name)
        return result, nil
    } else {
        return nil, fmt.Errorf("No author found")
    }
}

func (store *dbStore) ListAuthorBooks(authorId uint, noDetails bool, params models.Search) ([]models.Book, error) {
    result := []models.Book{}
    search := store.db.Select("distinct books.*").Table("books").
        Joins("left join book_authors on books.id=book_authors.book_id left join authors on authors.id=book_authors.author_id")
    search.Where("authors.ID=?", authorId)
    if !params.Deleted {
        search.Where("del = 0")
    }
    if len(params.Langs) > 0 {
        search.Where("lang in (" + strings.Join(params.Langs, ",") + ")")
    }

    search.Preload("Container").Order("series, cast(ser_no as unsigned), title").Find(&result)
    if !noDetails {
        result = store.fillBooksDetails(result, false)
    }
    return result, nil
}

func (store *dbStore) GetBook(bookId uint) (*models.Book, error) {
    result := new(models.Book)
    store.db.Preload("Container").First(result, bookId)
    store.fillBookDetails(result, true)
    if result.ID > 0 {
        return result, nil
    } else {
        return nil, fmt.Errorf("No book found")
    }
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

func (store *dbStore) IsContainerExist(fileName string) bool {
    contObj := new(models.Container)
    store.db.Where("file_name = ?", fileName).First(&contObj)
    return contObj.ID > 0
}

func (store *dbStore) Close() {
}

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
