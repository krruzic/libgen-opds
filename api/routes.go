package api

import (
	"encoding/xml"
	"net/http"
	"net/url"
	"time"

	"reichard.io/libgen-opds/client"
	"reichard.io/libgen-opds/opds"
)

func (api *API) RootHandler(w http.ResponseWriter, r *http.Request) {

	rootFeed := opds.Feed{
		Title:   "LibGen OPDS Bridge",
		Updated: time.Now().UTC(),
		Links: []opds.Link{
			opds.Link{
				Rel:      "self",
				TypeLink: "application/atom+xml",
			},
			opds.Link{
				Title:    "Search LibGen",
				Rel:      "search",
				TypeLink: "application/opensearchdescription+xml",
				Href:     "search.xml",
			},
		},
		Entries: []opds.Entry{
			opds.Entry{
				Title: "Most Read This Month",
				Content: &opds.Content{
					Content:     "Most Read This Month",
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
				Title: "Most Read This Year",
				Content: &opds.Content{
					Content:     "Most Read This Year",
					ContentType: "text",
				},
				Links: []opds.Link{
					opds.Link{
						Href:     "./most-read?cadence=year",
						TypeLink: "application/atom+xml;type=feed;profile=opds-catalog",
					},
				},
			},
		},
	}

	feedXML, _ := xml.Marshal(rootFeed)
	w.Header().Set("Content-Type", "application/xml")
	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	w.Write(feedXML)

}

func (api *API) SearchDescriptionHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/xml")
	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	w.Write([]byte(`
		<OpenSearchDescription xmlns="http://a9.com/-/spec/opensearch/1.1/">
			<ShortName>Search LibGen</ShortName>
			<Description>Search LibGen</Description>
			<Url type="application/atom+xml;profile=opds-catalog;kind=acquisition" template="./search?query={searchTerms}"/>
		</OpenSearchDescription>`))

}

func (api *API) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	md5 := r.URL.Query().Get("md5")

	// Acquire & Parse Download URL
	body := client.GetPage("http://library.lol/fiction/" + md5)
	downloadURL := client.ParseDownloadURL(body)

	// Redirect
	http.Redirect(w, r, downloadURL, 301)
}

func (api *API) GoodReadsMostReadHandler(w http.ResponseWriter, r *http.Request) {

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
	w.Header().Set("Content-Type", "application/xml")
	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	w.Write(feedXML)

}

func (api *API) LibGenSearchHandler(w http.ResponseWriter, r *http.Request) {

	// Derive URL
	query := r.URL.Query().Get("query")
	mirror := "https://libgen.is"
	language := "English"
	format := "epub"
	libGenURL := mirror +
		"/fiction/?q=" +
		url.QueryEscape(query) +
		"&criteria=&" +
		"&language=" + language +
		"&format=" + format

	// Acquire & Parse Source
	body := client.GetPage(libGenURL)
	allEntries := client.ParseLibGen(body)

	// Build XML
	mostReadFeed := &opds.Feed{
		Title:   "Search Results",
		Updated: time.Now().UTC(),
		Entries: allEntries,
	}
	feedXML, _ := xml.Marshal(mostReadFeed)

	// Serve
	w.Header().Set("Content-Type", "application/xml")
	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	w.Write(feedXML)

}
