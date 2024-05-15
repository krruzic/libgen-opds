package api

import (
	"log"
	"net/http"
	"os"
	"time"
)

type API struct {
	Router   *http.ServeMux
	Username string
	Password string
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	api.Router.ServeHTTP(w, r)
	log.Printf("- %s - %s %s - %v", r.RemoteAddr, r.Method, r.URL, time.Since(start))
}

func NewApi() *API {
	username := os.Getenv("API_USERNAME")
	if username == "" {
		log.Fatal("environment variable API_USERNAME must be set!")
	}
	password := os.Getenv("API_PASSWORD")
	if password == "" {
		log.Fatal("environment variable API_PASSWORD must be set!")
	}

	api := &API{
		Router:   http.NewServeMux(),
		Username: username,
		Password: password,
	}

	api.Router.HandleFunc("/", api.basicAuth(api.RootHandler))
	api.Router.HandleFunc("/search.xml", api.basicAuth(api.SearchDescriptionHandler))
	api.Router.HandleFunc("/search", api.basicAuth(api.SearchHandler))
	api.Router.HandleFunc("/most-read", api.basicAuth(api.GoodReadsMostReadHandler))
	api.Router.HandleFunc("/most-popular", api.basicAuth(api.LibZMostPopularHandler))
	api.Router.HandleFunc("/download", api.basicAuth(api.DownloadHandler))

	return api
}
