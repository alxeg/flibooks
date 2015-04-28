package main

import (
    "encoding/json"
    "fmt"
    "github.com/alxeg/flibooks/datastore"
    "github.com/alxeg/flibooks/inpx"
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
    limit        uint
)

func init() {
    flag.StringVar(&fileToParse, "parse", "", "Parse inpx to the local database")
    flag.StringVar(&dataDir, "data-dir", "", "Folder to put database files")
    flag.StringVar(&searchTitle, "search-title", "", "Search books by their title")
    flag.StringVar(&searchAuthor, "search-author", "", "Search books by author")
    flag.UintVar(&limit, "limit", 10, "Limit output results")
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
        result, err := store.FindBooksByTitle(searchTitle, limit)
        if err == nil && len(result) != 0 {
            jsonBytes, err := json.MarshalIndent(result, "", "  ")
            if err == nil {
                fmt.Println(string(jsonBytes))
            } else {
                log.Fatalln("Invalid object")
            }
        } else {
            log.Println("Nothing found")
        }
    } else {
        flag.PrintDefaults()
    }
}
