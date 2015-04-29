package inpx

import (
    "archive/zip"
    "github.com/alxeg/flibooks/models"
    "io"
    "log"
    "os"
    "path/filepath"
)

func UnzipBookFile(book *models.Book, targetFolder string, rename bool) (err error) {
    container := book.Container.FileName
    fileName := book.File + "." + book.Ext
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

            path := filepath.Join(targetFolder, outName)
            f, err := os.OpenFile(
                path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
            if err != nil {
                return err
            }
            defer f.Close()

            _, err = io.Copy(f, rc)
            break
        }
    }
    return err

}
