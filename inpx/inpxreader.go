package inpx

import (
    "archive/zip"
    "bufio"
    "fmt"
    "github.com/alxeg/flibooks/datastore"
    "github.com/alxeg/flibooks/models"
    "io"
    "log"
    "runtime"
    "strings"
)

var (
    versionFile  = "version.info"
    numProcesses = runtime.NumCPU() * 2
)

func ReadInpxFile(dataFile string, store datastore.DataStorer) (err error) {
    r, err := zip.OpenReader(dataFile)
    if err != nil {
        log.Fatal(err)
        return err
    }
    defer r.Close()

    files := make(chan *zip.File, 10)
    done := make(chan struct{})
    results := make(chan *models.Book, 10)

    go func() {
        for _, file := range r.File {
            files <- file
        }
        close(files)
    }()

    log.Println("Paralleling in ", numProcesses)
    for i := 0; i < numProcesses; i++ {
        go processInp(files, results, done)
    }

    waitAndProcessResults(done, results, store)

    store.Close()

    return nil
}

func processInp(files <-chan *zip.File, results chan<- *models.Book, done chan<- struct{}) {

    defer func() {
        done <- struct{}{}
    }()

    for file := range files {

        if !strings.HasSuffix(file.FileInfo().Name(), ".inp") {
            continue
        }

        log.Println("Processing file ", file.FileInfo().Name())
        rc, err := file.Open()
        if err == nil {
            reader := bufio.NewReader(rc)
            for {
                line, err := reader.ReadString('\n')
                if line != "" {
                    book, bookErr := processBook(line)
                    if bookErr == nil {
                        book.Container = models.Container{FileName: strings.Replace(file.FileInfo().Name(), ".inp", ".zip", -1)}
                        results <- book
                    } else {
                        log.Println(bookErr)
                    }
                }
                if err != nil {
                    if err != io.EOF {
                        log.Println("failed to finish reading the file:", err)
                    }
                    break
                }
            }

            rc.Close()
        }

        if err != nil {
            log.Println("Error occured while reading file", err)
        }
        log.Println("Done with ", file.FileInfo().Name())
    }

}

func trimSlice(in []string) []string {
    for len(in) > 0 && in[len(in)-1] == "" {
        in = in[:len(in)-1]
    }
    return in
}

func processBook(line string) (book *models.Book, err error) {
    elements := strings.Split(line, string([]byte{0x04}))
    if len(elements) < 12 {
        return book, fmt.Errorf("Illegal number of elements")
    }

    book = new(models.Book)
    for _, author := range trimSlice(strings.Split(elements[0], ":")) {
        book.Authors = append(book.Authors, models.Author{Name: strings.ToLower(author)})
    }
    for _, genre := range trimSlice(strings.Split(elements[1], ":")) {
        book.Genres = append(book.Genres, models.Genre{GenreCode: genre})
    }
    book.Title = strings.ToLower(elements[2])
    book.Series = elements[3]
    book.SerNo = elements[4]
    book.File = elements[5]
    book.FileSize = elements[6]
    book.LibId = elements[7]
    book.Del = elements[8]
    book.Ext = elements[9]
    book.Date = elements[10]
    book.Lang = elements[11]
    return book, nil
}

func waitAndProcessResults(done <-chan struct{}, results <-chan *models.Book, store datastore.DataStorer) {
    for working := numProcesses; working > 0; {
        select {
        case book := <-results:
            store.PutBook(book)
        case <-done:
            log.Println("Gorutine finished work")
            working--
        }
    }

    for {
        select { // Nonblocking
        case book := <-results:
            store.PutBook(book)
        default:
            return
        }
    }
}
