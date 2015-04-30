package main

import (
    "encoding/json"
    "fmt"
    "github.com/alxeg/flibooks/datastore"
    "github.com/alxeg/flibooks/inpx"
    "github.com/alxeg/flibooks/rest"
    flag "github.com/ogier/pflag"
    "log"
    "os"
    "path/filepath"
)

var (
    fileToParse  string
    dataDir      string
    searchTitle  string
    searchAuthor string
    listAuthor   uint
    limit        int
    getBook      uint
    save         bool
    listen       string
)

func init() {
    flag.StringVar(&fileToParse, "parse", "", "Parse inpx to the local database")
    flag.StringVar(&dataDir, "data-dir", "", "Folder to put database files")
    flag.StringVar(&searchTitle, "search-title", "", "Search books by their title")
    flag.StringVar(&searchAuthor, "search-author", "", "Search authors, or books by author if comes with search-title")
    flag.IntVar(&limit, "limit", 10, "Limit search results (-1 for no limit)")
    flag.UintVar(&listAuthor, "list-author", 0, "List all author's books by id")
    flag.UintVar(&getBook, "get-book", 0, "Get book by its id")
    flag.BoolVar(&save, "save", false, "Save book file to the disk")
    flag.StringVar(&listen, "listen", ":8000", "Set server listen address:port")

}

func printJson(object interface{}) {
    jsonBytes, err := json.MarshalIndent(object, "", "  ")
    if err == nil {
        fmt.Println(string(jsonBytes))
    } else {
        log.Fatalln("Invalid object")
    }

}

func main() {
    flag.Parse()

    if dataDir == "" {
        dataDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
    }

    store, err := datastore.NewDBStore(dataDir, fileToParse != "")
    if err != nil {
        log.Fatalln("Failed to open database")
    }

    if fileToParse != "" {
        log.Printf("Opening %s to parse data\n", fileToParse)
        inpx.ReadInpxFile(fileToParse, store)
    } else if searchTitle != "" {
        result, err := store.FindBooks(searchTitle, searchAuthor, limit)
        if err == nil && len(result) != 0 {
            printJson(result)
        } else {
            log.Println("Nothing found")
        }
    } else if searchAuthor != "" {
        result, err := store.FindAuthors(searchAuthor, limit)
        if err == nil && len(result) != 0 {
            printJson(result)
        } else {
            log.Println("Noone found")
        }
    } else if listAuthor > 0 {
        result, err := store.ListAuthorBooks(listAuthor)
        if err == nil && len(result) != 0 {
            printJson(result)
        } else {
            log.Println("Nothing found")
        }
    } else if getBook > 0 {
        result, err := store.GetBook(getBook)
        if err == nil {
            printJson(result)
            if save {
                err = inpx.UnzipBookFile(result, dataDir, true)
                if err != nil {
                    log.Fatalln("Failed to save file", err)
                }
            }
        } else {
            log.Println("Nothing found")
        }

    } else {
        fmt.Println("Additional parameters are:")
        flag.PrintDefaults()
        rest.NewRestService(listen, store).StartListen()
    }
}
