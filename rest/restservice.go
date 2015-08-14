package rest

import (
    "github.com/alxeg/flibooks/datastore"
    "github.com/alxeg/flibooks/inpx"
    "github.com/alxeg/flibooks/models"
    "github.com/alxeg/flibooks/utils"
    "github.com/emicklei/go-restful"
    "log"
    "net/http"
    "net/url"
    "strconv"
    "strings"
)

type RestService struct {
    listen    string
    dataDir   string
    dataStore datastore.DataStorer
    container *restful.Container
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
        Returns(200, "OK", models.Book{}))

    ws.Route(ws.GET("/{bookId}/download").
        To(service.downloadBook).
        Doc("Download book content").
        Operation("downloadBook").
        Param(ws.PathParameter("bookId", "identifier of the book").DataType("int")).
        Returns(200, "OK", models.Book{}))

    ws.Route(ws.POST("/search").
        To(service.searchBooks).
        Doc("Search for the books").
        Operation("searchBooks").
        Returns(200, "OK", []models.Book{}))

    ws.Route(ws.GET("/lib/{libId}").
        To(service.getBooksByLibId).
        Doc("Get books by libId").
        Operation("getBooksByLibId").
        Param(ws.PathParameter("libId", "libId of the book").DataType("string")).
        Returns(200, "OK", []models.Book{}))

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
        Returns(200, "OK", models.Author{}))

    ws.Route(ws.GET("/{authorId}/books").
        To(service.listAuthorsBooks).
        Doc("Show author's books").
        Operation("listAuthorsBooks").
        Param(ws.PathParameter("authorId", "identifier of the author").DataType("int")).
        Returns(200, "OK", []models.Book{}))

    ws.Route(ws.POST("/search").
        To(service.searchAuthors).
        Doc("Search authors").
        Operation("searchAuthors").
        Returns(200, "OK", []models.Author{}))

    container.Add(ws)
}

func (service RestService) getBook(request *restful.Request, response *restful.Response) {
    bookId, _ := strconv.ParseUint(request.PathParameter("bookId"), 0, 32)
    log.Println("Requesting book ", bookId)
    result, err := service.dataStore.GetBook(uint(bookId))
    if err == nil {
        response.WriteEntity(result)
    } else {
        response.AddHeader("Content-Type", "text/plain")
        response.WriteErrorString(http.StatusNotFound, "Book wasn't found")
    }
}

func (service RestService) getBooksByLibId(request *restful.Request, response *restful.Response) {
    libId := request.PathParameter("libId")
    log.Println("Get books by libId ", libId)
    result, err := service.dataStore.FindBooksByLibId(libId)
    if err == nil && len(result) != 0 {
        response.WriteEntity(result)
    } else {
        response.AddHeader("Content-Type", "text/plain")
        response.WriteErrorString(http.StatusNotFound, "Nothing was found")
    }
}

func (service RestService) downloadBook(request *restful.Request, response *restful.Response) {
    bookId, _ := strconv.ParseUint(request.PathParameter("bookId"), 0, 32)
    log.Println("Downloading book ", bookId)
    result, err := service.dataStore.GetBook(uint(bookId))
    if err == nil {
        authors := ""
        for _, a := range result.Authors {
            authors = authors + a.Name
        }
        outName := authors + " - " + result.Title + "." + result.Ext

        response.AddHeader("Content-Type", "application/octet-stream")
        response.AddHeader("Content-disposition", "attachment; filename*=UTF-8''"+strings.Replace(url.QueryEscape(
            utils.ReplaceUnsupported(outName)), "+", "%20", -1))

        err := inpx.UnzipBookToWriter(service.dataDir, result, response)
        if err != nil {
            response.AddHeader("Content-Type", "text/plain")
            response.WriteErrorString(http.StatusNotFound, "Book wasn't found")
        }
    } else {
        response.AddHeader("Content-Type", "text/plain")
        response.WriteErrorString(http.StatusNotFound, "Book wasn't found")
    }
}

func (service RestService) searchBooks(request *restful.Request, response *restful.Response) {
    search := models.Search{}
    request.ReadEntity(&search)
    log.Println("Searching books ", search)

    result, err := service.dataStore.FindBooks(search.Title, search.Author, search.Limit)
    if err == nil && len(result) != 0 {
        response.WriteEntity(result)
    } else {
        response.AddHeader("Content-Type", "text/plain")
        response.WriteErrorString(http.StatusNotFound, "Nothing was found")
    }
}

func (service RestService) searchAuthors(request *restful.Request, response *restful.Response) {
    search := models.Search{}
    request.ReadEntity(&search)
    log.Println("Searching authors ", search)

    result, err := service.dataStore.FindAuthors(search.Author, search.Limit)
    if err == nil && len(result) != 0 {
        response.WriteEntity(result)
    } else {
        response.AddHeader("Content-Type", "text/plain")
        response.WriteErrorString(http.StatusNotFound, "Nothing was found")
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
        response.WriteErrorString(http.StatusNotFound, "No author was found")
    }
}

func (service RestService) listAuthorsBooks(request *restful.Request, response *restful.Response) {
    authorId, _ := strconv.ParseUint(request.PathParameter("authorId"), 0, 32)
    noDetails, _ := utils.ParseBool(request.QueryParameter("no-details"))

    log.Println("Requesting author's books ", authorId)

    result, err := service.dataStore.ListAuthorBooks(uint(authorId), noDetails)
    if err == nil {
        response.WriteEntity(result)
    } else {
        response.AddHeader("Content-Type", "text/plain")
        response.WriteErrorString(http.StatusNotFound, "No books was found")
    }
}

func (service RestService) StartListen() {
    log.Println("Start listening on ", service.listen)
    server := &http.Server{Addr: service.listen, Handler: service.container}
    log.Fatal(server.ListenAndServe())
}

func NewRestService(listen string, dataStore datastore.DataStorer, dataDir string) *RestService {
    service := new(RestService)
    service.listen = listen
    service.dataStore = dataStore
    service.dataDir = dataDir
    service.container = restful.NewContainer()
    service.container.Router(restful.CurlyRouter{})

    service.registerBookResource(service.container)
    service.registerAuthorResource(service.container)

    return service
}
