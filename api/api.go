package api

import (
	"log"
	"net/http"
	"time"
)

type API struct {
	Router *http.ServeMux
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	api.Router.ServeHTTP(w, r)
	log.Printf("- %s - %s %s - %v", r.RemoteAddr, r.Method, r.URL, time.Since(start))
}

func NewApi() *API {
	api := &API{
		Router: http.NewServeMux(),
	}

	api.Router.HandleFunc("/", api.RootHandler)
	api.Router.HandleFunc("/search.xml", api.SearchDescriptionHandler)
	api.Router.HandleFunc("/search", api.SearchHandler)
	api.Router.HandleFunc("/most-read", api.GoodReadsMostReadHandler)
	api.Router.HandleFunc("/most-popular", api.LibZMostPopularHandler)
	api.Router.HandleFunc("/download", api.DownloadHandler)

	return api
}
