package api

import (
	"net/http"
)

type API struct {
	Router *http.ServeMux
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
