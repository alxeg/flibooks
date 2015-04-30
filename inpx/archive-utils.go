package inpx

import (
    "archive/zip"
    "github.com/alxeg/flibooks/models"
    "io"
    "log"
    "os"
    "path/filepath"
)

func UnzipBookToWriter(book *models.Book, writer io.Writer) (err error) {
    container := book.Container.FileName
    fileName := book.File + "." + book.Ext

    r, err := zip.OpenReader(container)
    if err != nil {
        log.Fatalln("Failed to open container", container)
    }
    defer r.Close()
    for _, file := range r.File {
        if file.FileInfo().Name() == fileName {
            rc, err := file.Open()
            if err != nil {
                return err
            }
            defer rc.Close()

            _, err = io.Copy(writer, rc)

            break
        }
    }
    return err
}

func UnzipBookFile(book *models.Book, targetFolder string, rename bool) (err error) {
    var outName string
    if rename {
        authors := ""
        for _, a := range book.Authors {
            authors = authors + a.Name
        }
        outName = authors + " - " + book.Title + "." + book.Ext
    } else {
        outName = book.File + "." + book.Ext
    }

    path := filepath.Join(targetFolder, outName)
    f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        return err
    }
    defer f.Close()

    UnzipBookToWriter(book, f)

    return err

}
