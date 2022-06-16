package api

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/jinzhu/copier"

	"github.com/alxeg/flibooks/internal/db"
	"github.com/alxeg/flibooks/internal/db/orm"
	"github.com/alxeg/flibooks/internal/services/convert"
	"github.com/alxeg/flibooks/pkg/inpx"
	"github.com/alxeg/flibooks/pkg/inpx/models"
	"github.com/alxeg/flibooks/pkg/utils"
)

var (
	allowedFormats = map[string]string{
		"epub": "application/epub+zip",
		"azw3": "application/vnd.amazon.ebook",
		"mobi": "application/vnd.amazon.mobi8-ebook",
	}
)

type RestService struct {
	listen    string
	dataDir   string
	dataStore db.DataStorer
	container *restful.Container
	converter convert.Converter
}

func (service RestService) registerBookResource(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/book").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/{bookId}").
		To(service.getBook).
		Doc("Get specific book info").
		Operation("getBook").
		Param(ws.PathParameter("bookId", "identifier of the book").DataType("int")).
		Returns(200, "OK", orm.Book{}))

	ws.Route(ws.GET("/langs").
		To(service.getLangs).
		Doc("Get all available books languages").
		Operation("getLangs").
		Returns(200, "OK", []string{"en"}))

	ws.Route(ws.GET("/{bookId}/download").
		To(service.downloadBook).
		Doc("Download book content").
		Operation("downloadBook").
		Param(ws.PathParameter("bookId", "identifier of the book").DataType("int")).
		Returns(200, "OK", orm.Book{}))

	ws.Route(ws.GET("/archive").
		To(service.downloadBooksArchive).
		Doc("Download books in single zip file").
		Operation("downloadBooksArchive").
		Returns(200, "OK", orm.Book{}))

	ws.Route(ws.POST("/search").
		To(service.searchBooks).
		Doc("Search for the books").
		Operation("searchBooks").
		Returns(200, "OK", []orm.Book{}))

	ws.Route(ws.POST("/series").
		To(service.searchSeries).
		Doc("Search for the books").
		Operation("searchBooks").
		Returns(200, "OK", []orm.Book{}))

	ws.Route(ws.GET("/lib/{libId}").
		To(service.getBooksByLibID).
		Doc("Get books by libId").
		Operation("getBooksByLibId").
		Param(ws.PathParameter("libId", "libId of the book").DataType("string")).
		Returns(200, "OK", []orm.Book{}))

	container.Add(ws)
}

func (service RestService) registerAuthorResource(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/author").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/{authorId}").
		To(service.getAuthor).
		Doc("Get author's info").
		Operation("getAuthor").
		Param(ws.PathParameter("authorId", "identifier of the author").DataType("int")).
		Returns(200, "OK", orm.Author{}))

	ws.Route(ws.GET("/{authorId}/books").
		To(service.listAuthorsBooks).
		Doc("Show author's books").
		Operation("listAuthorsBooks").
		Param(ws.PathParameter("authorId", "identifier of the author").DataType("int")).
		Returns(200, "OK", []orm.Book{}))

	ws.Route(ws.POST("/{authorId}/books").
		To(service.listAuthorsBooksPost).
		Doc("Show author's books").
		Operation("listAuthorsBooks").
		Param(ws.PathParameter("authorId", "identifier of the author").DataType("int")).
		Returns(200, "OK", []orm.Book{}))

	ws.Route(ws.POST("/search").
		To(service.searchAuthors).
		Doc("Search authors").
		Operation("searchAuthors").
		Returns(200, "OK", []orm.Author{}))

	container.Add(ws)
}

func (service RestService) getBook(request *restful.Request, response *restful.Response) {
	bookID, _ := strconv.ParseUint(request.PathParameter("bookId"), 0, 32)
	log.Println("Requesting book ", bookID)
	result, err := service.dataStore.GetBook(uint(bookID))
	if err == nil {
		response.WriteEntity(result)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Book wasn't found\n")
	}
}

func (service RestService) getBooksByLibID(request *restful.Request, response *restful.Response) {
	libID := request.PathParameter("libId")
	log.Println("Get books by libId ", libID)
	result, err := service.dataStore.FindBooksByLibID(libID)
	if err == nil && len(result) != 0 {
		response.WriteEntity(result)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Nothing was found\n")
	}
}

func (service RestService) downloadBook(request *restful.Request, response *restful.Response) {
	answerError := func(status int, reason string, params ...any) {
		response.AddHeader("Content-Type", "text/plain")
		errorStr := fmt.Sprint(reason+"\n", params)
		log.Println(errorStr)
		response.WriteErrorString(status, errorStr)
	}

	bookID, _ := strconv.ParseUint(request.PathParameter("bookId"), 0, 32)
	outFormat := request.QueryParameter("format")

	log.Println("Downloading book ", bookID)
	result, err := service.dataStore.GetBook(uint(bookID))
	if err != nil {
		answerError(http.StatusNotFound, "Book hasn't been found")
		return
	}
	outName := result.GetFullFilename()
	extBook := &models.Book{}
	copier.Copy(extBook, result)

	if _, ok := allowedFormats[outFormat]; !ok {
		response.AddHeader("Content-Type", "application/octet-stream")
		response.AddHeader("Content-disposition", "attachment; filename*=UTF-8''"+strings.Replace(url.QueryEscape(
			utils.ReplaceUnsupported(outName)), "+", "%20", -1))

		err := inpx.UnzipBookToWriter(service.dataDir, extBook, response)
		if err != nil {
			answerError(http.StatusNotFound, "Cannot retrieve the book file: %s", err.Error())
			return
		}
	} else {
		convName := strings.TrimSuffix(outName, filepath.Ext(outName)) + "." + outFormat
		tmpDir, _ := ioutil.TempDir("", "fliconvert")
		defer os.RemoveAll(tmpDir)

		srcPath := path.Join(tmpDir, "file.fb2")
		src, err := os.Create(srcPath)
		if err != nil {
			answerError(http.StatusNotFound, "Cannot open the src file for writing: %s", err.Error())
			return
		}
		err = func() error {
			defer src.Close()
			return inpx.UnzipBookToWriter(service.dataDir, extBook, src)
		}()
		if err != nil {
			answerError(http.StatusNotFound, "Cannot retrieve the book file: %s", err.Error())
			return
		}

		err = service.converter.Convert(srcPath, tmpDir, outFormat)
		if err != nil {
			answerError(http.StatusNotFound, "Cannot convert the book to %s format: %s", outFormat, err.Error())
			return
		}
		dstPath := path.Join(tmpDir, "file."+outFormat)
		fileBytes, err := ioutil.ReadFile(dstPath)
		if err != nil {
			answerError(http.StatusNotFound, "Cannot read the converted file: %s", outFormat, err.Error())
			return
		}

		response.AddHeader("Content-Type", allowedFormats[outFormat])
		response.AddHeader("Content-disposition", "attachment; filename*=UTF-8''"+strings.Replace(url.QueryEscape(
			utils.ReplaceUnsupported(convName)), "+", "%20", -1))
		outFileNfo, err := os.Stat(dstPath)
		if err == nil {
			response.AddHeader("Content-length", fmt.Sprintf("%d", outFileNfo.Size()))
		}

		response.Write(fileBytes)
	}
}

func (service RestService) downloadBooksArchive(request *restful.Request, response *restful.Response) {
	request.Request.ParseForm()
	ids := request.Request.Form["id"]
	if len(ids) > 0 {
		response.Header().Set("Content-Type", "application/zip")
		response.Header().Set("Content-disposition", "attachment; filename*=UTF-8''"+strings.Replace(url.QueryEscape(
			"flibooks-"+time.Now().Format("2006-01-02T15-04-05")+".zip"), "+", "%20", -1))
		zipWriter := zip.NewWriter(response)

		idsChan := make(chan string)
		done := make(chan bool)

		go func() {
			for {
				id, more := <-idsChan
				if more {
					bookID, _ := strconv.ParseUint(id, 0, 32)
					book, err := service.dataStore.GetBook(uint(bookID))
					if err == nil {
						zipHeader := &zip.FileHeader{Name: book.GetFullFilename(), Method: zip.Deflate, Flags: 0x800}
						entry, err := zipWriter.CreateHeader(zipHeader)
						// entry, err := zipWriter.Create(book.GetFullFilename())

						if err == nil {
							extBook := &models.Book{}
							copier.Copy(extBook, book)
							inpx.UnzipBookToWriter(service.dataDir, extBook, entry)
						} else {
							log.Println("Failed to compress ", book.GetFullFilename())
						}
					} else {
						log.Println("Failed to get book ", id)
					}
				} else {
					done <- true
					return
				}
			}
		}()
		for _, id := range ids {
			idsChan <- id
		}
		close(idsChan)
		<-done
		zipWriter.Close()
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "No parameters passed\n")
	}
}

func (service RestService) searchBooks(request *restful.Request, response *restful.Response) {
	search := orm.Search{}
	request.ReadEntity(&search)
	log.Println("Searching books ", search)

	result, err := service.dataStore.FindBooks(search)
	if err == nil && len(result) != 0 {
		response.WriteEntity(result)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Nothing was found\n")
	}
}

func (service RestService) searchSeries(request *restful.Request, response *restful.Response) {
	search := orm.Search{}
	request.ReadEntity(&search)
	log.Println("Searching book series ", search)

	result, err := service.dataStore.FindBooksSeries(search)
	if err == nil && len(result) != 0 {
		response.WriteEntity(result)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Nothing was found\n")
	}
}

func (service RestService) getLangs(request *restful.Request, response *restful.Response) {
	log.Println("Getting languages")

	result, err := service.dataStore.GetLangs()
	if err == nil && len(result) != 0 {
		response.WriteEntity(result)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Nothing was found\n")
	}
}

func (service RestService) searchAuthors(request *restful.Request, response *restful.Response) {
	search := orm.Search{}
	request.ReadEntity(&search)
	log.Println("Searching authors ", search)

	result, err := service.dataStore.FindAuthors(search.Author, search.Limit)
	if err == nil && len(result) != 0 {
		response.WriteEntity(result)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Nothing was found\n")
	}
}

func (service RestService) getAuthor(request *restful.Request, response *restful.Response) {
	authorId, _ := strconv.ParseUint(request.PathParameter("authorId"), 0, 32)
	log.Println("Requesting author ", authorId)

	result, err := service.dataStore.GetAuthor(uint(authorId))
	if err == nil {
		response.WriteEntity(result)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "No author was found\n")
	}
}

func (service RestService) listAuthorsBooks(request *restful.Request, response *restful.Response) {
	authorId, _ := strconv.ParseUint(request.PathParameter("authorId"), 0, 32)
	noDetails, _ := utils.ParseBool(request.QueryParameter("no-details"))

	log.Println("Requesting author's books ", authorId)

	result, err := service.dataStore.ListAuthorBooks(uint(authorId), noDetails, orm.Search{})
	if err == nil {
		response.WriteEntity(result)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "No books was found\n")
	}
}

func (service RestService) listAuthorsBooksPost(request *restful.Request, response *restful.Response) {
	authorId, _ := strconv.ParseUint(request.PathParameter("authorId"), 0, 32)
	noDetails, _ := utils.ParseBool(request.QueryParameter("no-details"))
	search := orm.Search{}
	request.ReadEntity(&search)

	log.Println("Requesting author's books ", authorId)

	result, err := service.dataStore.ListAuthorBooks(uint(authorId), noDetails, search)
	if err == nil {
		response.WriteEntity(result)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "No books was found\n")
	}
}

func (service RestService) StartListen() {
	log.Println("Start listening on ", service.listen)
	server := &http.Server{Addr: service.listen, Handler: service.container}
	log.Fatal(server.ListenAndServe())
}

func NewRestService(listen string, dataStore db.DataStorer, dataDir string, converter convert.Converter) RestServer {
	service := new(RestService)
	service.listen = listen
	service.dataStore = dataStore
	service.dataDir = dataDir
	service.container = restful.NewContainer()
	service.container.Router(restful.CurlyRouter{})
	service.converter = converter

	service.registerBookResource(service.container)
	service.registerAuthorResource(service.container)

	return service
}
