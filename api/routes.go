package api

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/xml"
	"net/http"
	"net/url"
	"reichard.io/libgen-opds/client"
	"reichard.io/libgen-opds/opds"
	"time"
)

func (api *API) basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(api.Username))
			expectedPasswordHash := sha256.Sum256([]byte(api.Password))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func (api *API) RootHandler(w http.ResponseWriter, r *http.Request) {

	// Headers
	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	w.Header().Set("Content-Type", "application/xml")

	if r.Method == http.MethodHead {
		return
	}

	rootFeed := opds.Feed{
		Title:   "LibGen OPDS Bridge",
		Updated: time.Now().UTC(),
		Links: []opds.Link{
			opds.Link{
				Title:    "Search LibGen",
				Rel:      "search",
				TypeLink: "application/opensearchdescription+xml",
				Href:     "search.xml",
			},
		},
		Entries: []opds.Entry{
			opds.Entry{
				Title: "Goodreads - Most Read This Month",
				Content: &opds.Content{
					Content:     "Goodreads - Most Read This Month",
					ContentType: "text",
				},
				Links: []opds.Link{
					opds.Link{
						Href:     "./most-read?cadence=month",
						TypeLink: "application/atom+xml;type=feed;profile=opds-catalog",
					},
				},
			},
			opds.Entry{
				Title: "Goodreads - Most Read This Year",
				Content: &opds.Content{
					Content:     "Goodreads - Most Read This Year",
					ContentType: "text",
				},
				Links: []opds.Link{
					opds.Link{
						Href:     "./most-read?cadence=year",
						TypeLink: "application/atom+xml;type=feed;profile=opds-catalog",
					},
				},
			},
			// opds.Entry{
			// 	Title: "Check for Updates",
			// 	Content: &opds.Content{
			// 		Content:     "Check for Updates",
			// 		ContentType: "text",
			// 	},
			// 	Links: []opds.Link{
			// 		opds.Link{
			// 			Href:     "./update-check",
			// 			TypeLink: "application/atom+xml;type=feed;profile=opds-catalog",
			// 		},
			// 	},
			// },
		},
	}

	feedXML, _ := xml.Marshal(rootFeed)
	w.Write(feedXML)

}

func (api *API) SearchDescriptionHandler(w http.ResponseWriter, r *http.Request) {

	// Headers
	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	w.Header().Set("Content-Type", "application/xml")

	if r.Method == http.MethodHead {
		return
	}

	w.Write([]byte(`
		<OpenSearchDescription xmlns="http://a9.com/-/spec/opensearch/1.1/">
			<ShortName>Search LibGen</ShortName>
			<Description>Search LibGen</Description>
			<Url type="application/atom+xml;profile=opds-catalog;kind=acquisition" template="./search?query={searchTerms}"/>
		</OpenSearchDescription>`))

}

func (api *API) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	downloadType := r.URL.Query().Get("type")

	// Derive Info URL
	var infoURL string
	if downloadType == "fiction" {
		infoURL = "http://library.lol/fiction/" + id
	} else if downloadType == "non-fiction" {
		infoURL = "http://library.lol/main/" + id
	}

	// Parse & Derive Download URL
	body := client.GetPage(infoURL)
	downloadURL := client.ParseLibGenDownloadURL(body)

	// Redirect
	http.Redirect(w, r, downloadURL, 301)
}

func (api *API) GoodReadsMostReadHandler(w http.ResponseWriter, r *http.Request) {

	// Headers
	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	w.Header().Set("Content-Type", "application/xml")

	if r.Method == http.MethodHead {
		return
	}

	// Derive Duration
	var duration string
	if r.URL.Query().Get("cadence") == "month" {
		duration = "m"
	} else {
		duration = "y"
	}

	// Acquire & Parse Page Source
	body := client.GetPage("https://www.goodreads.com/book/most_read?category=all&country=US&duration=" + duration)
	allEntries := client.ParseGoodReads(body)

	// Build XML
	mostReadFeed := &opds.Feed{
		Title:   "GoodReads Most Read",
		Updated: time.Now().UTC(),
		Entries: allEntries,
	}
	feedXML, _ := xml.Marshal(mostReadFeed)

	// Serve
	w.Write(feedXML)

}

func (api *API) LibZMostPopularHandler(w http.ResponseWriter, r *http.Request) {
	// Headers
	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	w.Header().Set("Content-Type", "application/xml")

	if r.Method == http.MethodHead {
		return
	}

	// Acquire & Parse Page Source
	body := client.GetPage("https://usa1lib.org/popular.php")
	allEntries := client.ParseZLibPopular(body)

	// Build XML
	mostReadFeed := &opds.Feed{
		Title:   "ZLib Most Popular",
		Updated: time.Now().UTC(),
		Entries: allEntries,
	}
	feedXML, _ := xml.Marshal(mostReadFeed)

	// Serve
	w.Write(feedXML)
}

func (api *API) SearchHandler(w http.ResponseWriter, r *http.Request) {

	// Headers
	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	w.Header().Set("Content-Type", "application/xml")

	if r.Method == http.MethodHead {
		return
	}

	// Acquire Params
	query := r.URL.Query().Get("query")
	searchType := r.URL.Query().Get("type")

	var allEntries []opds.Entry
	feedTitle := "Search Results"

	if searchType == "fiction" {
		// Search Fiction
		url := "https://libgen.is/fiction/?q=" + url.QueryEscape(query) + "&language=English"
		body := client.GetPage(url)
		allEntries = client.ParseLibGenFiction(body)
	} else if searchType == "non-fiction" {
		// Search NonFiction
		url := "https://libgen.is/search.php?req=" + url.QueryEscape(query)
		body := client.GetPage(url)
		allEntries = client.ParseLibGenNonFiction(body)
	} else {
		// Offer Options
		feedTitle = "Select Search Type"
		allEntries = []opds.Entry{
			opds.Entry{
				Title: "Search Fiction",
				Content: &opds.Content{
					Content:     "Search Fiction",
					ContentType: "text",
				},
				Links: []opds.Link{
					opds.Link{
						Href:     "./search?type=fiction&query=" + url.QueryEscape(query),
						TypeLink: "application/atom+xml;type=feed;profile=opds-catalog",
					},
				},
			},
			opds.Entry{
				Title: "Search Non-Fiction",
				Content: &opds.Content{
					Content:     "Search Non-Fiction",
					ContentType: "text",
				},
				Links: []opds.Link{
					opds.Link{
						Href:     "./search?type=non-fiction&query=" + url.QueryEscape(query),
						TypeLink: "application/atom+xml;type=feed;profile=opds-catalog",
					},
				},
			},
		}
	}

	// Build XML
	searchFeed := &opds.Feed{
		Title:   feedTitle,
		Updated: time.Now().UTC(),
		Entries: allEntries,
	}
	feedXML, _ := xml.Marshal(searchFeed)

	// Serve
	w.Write(feedXML)

}
