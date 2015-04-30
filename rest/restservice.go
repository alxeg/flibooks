package rest

import (
    "github.com/alxeg/flibooks/datastore"
    "github.com/alxeg/flibooks/models"
    "github.com/emicklei/go-restful"
    "log"
    "net/http"
)

type RestService struct {
    listen    string
    dataStore datastore.DataStorer
    container *restful.Container
}

func (service RestService) registerBookResource(container *restful.Container) {
    ws := new(restful.WebService)
    ws.
        Path("/book").
        Consumes(restful.MIME_JSON).
        Produces(restful.MIME_JSON)

    ws.Route(ws.POST("/search").To(service.searchBooks).
        Doc("Search for the books").
        Operation("searchBooks").
        Returns(200, "OK", []models.Book{}))

    container.Add(ws)
}

func (service RestService) searchBooks(request *restful.Request, response *restful.Response) {
    search := models.Search{}
    request.ReadEntity(&search)

    result, err := service.dataStore.FindBooks(search.Title, search.Author, search.Limit)
    if err == nil && len(result) != 0 {
        response.WriteEntity(result)
    } else {
        response.AddHeader("Content-Type", "text/plain")
        response.WriteErrorString(http.StatusNotFound, "Nothing was found")
    }
}

func (service RestService) StartListen() {
    log.Println("Start listening on ", service.listen)
    server := &http.Server{Addr: service.listen, Handler: service.container}
    log.Fatal(server.ListenAndServe())
}

func NewRestService(listen string, dataStore datastore.DataStorer) *RestService {
    service := new(RestService)
    service.listen = listen
    service.dataStore = dataStore
    service.container = restful.NewContainer()
    service.container.Router(restful.CurlyRouter{})

    service.registerBookResource(service.container)

    return service
}
