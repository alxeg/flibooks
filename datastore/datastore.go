package datastore

import (
    "github.com/alxeg/flibooks/models"
    "github.com/alxeg/flibooks/utils"
    "github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"
    "os"
    "strings"
)

type dbStore struct {
    db    gorm.DB
    reset bool
}

func (store dbStore) PutBook(book models.Book) (err error) {
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

func (store dbStore) FindBooksByTitle(title string) ([]models.Book, error) {
    result := []models.Book{}
    search := store.db.Order("title")
    for _, term := range utils.SplitBySeparators(title) {
        search = search.Where("title LIKE ?", "%"+strings.ToLower(term)+"%")
    }
    search.Find(&result)
    for i, book := range result {
        store.db.Model(&book).Related(&book.Authors, "Authors")
        result[i].Authors = book.Authors
        store.db.Model(&book).Related(&book.Genres, "Genres")
        result[i].Genres = book.Genres
    }
    return result, nil
}

func (store dbStore) FindBooksByAuthor(author string) ([]models.Book, error) {
    result := []models.Book{}

    return result, nil
}

func (store dbStore) Close() {
    if store.reset {
    }
}

func NewDBStore(dbPath string, reset bool) (DataStorer, error) {
    dataPath := dbPath + "/fli-data.db"
    if reset {
        os.Remove(dataPath)
    }
    db, err := gorm.Open("sqlite3", dataPath)
    if err == nil {
        db.DB()
        db.AutoMigrate(&models.Author{}, &models.Container{}, &models.Genre{}, &models.Book{})
        // db.LogMode(true)
    }

    return dbStore{db, reset}, err
}
