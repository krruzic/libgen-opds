package client

import (
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"reichard.io/libgen-opds/opds"
)

// Mime Type Mappings
var mimeMapping map[string]string = map[string]string{
	"epub": "application/epub+zip",
	"azw":  "application/vnd.amazon.mobi8-ebook",
	"mobi": "application/x-mobipocket-ebook",
	"pdf":  "application/pdf",
	"zip":  "application/zip",
	"txt":  "text/plain",
	"rtf":  "application/rtf",
	"htm":  "text/html",
	"html": "text/html",
	"doc":  "application/msword",
	"lit":  "application/x-ms-reader",
}

func GetPage(page string) io.ReadCloser {
	resp, _ := http.Get(page)
	return resp.Body
}

func StripText(text string) string {
	reMultiSpace := regexp.MustCompile(`\s+`)
	return reMultiSpace.ReplaceAllString(text, " ")
}

/* -------------------------------------------------------------------------- */
/* ---------------------------- Library Genesis ----------------------------- */
/* -------------------------------------------------------------------------- */

func ParseLibGenFiction(body io.ReadCloser) []opds.Entry {
	// Parse
	defer body.Close()
	doc, _ := goquery.NewDocumentFromReader(body)

	// Normalize Results
	var allEntries []opds.Entry
	doc.Find("table.catalog tbody > tr").Each(func(ix int, rawBook *goquery.Selection) {

		// Parse File Details
		fileItem := rawBook.Find("td:nth-child(5)")
		fileDesc := fileItem.Text()
		fileDescSplit := strings.Split(fileDesc, "/")
		fileType := strings.ToLower(strings.TrimSpace(fileDescSplit[0]))

		// Parse Upload Date
		uploadedRaw, _ := fileItem.Attr("title")
		uploadedDateRaw := strings.Split(uploadedRaw, "Uploaded at ")[1]
		uploadDate, _ := time.Parse("2006-01-02 15:04:05", uploadedDateRaw)

		// Parse MD5
		editHref, _ := rawBook.Find("td:nth-child(7) a").Attr("href")
		hrefArray := strings.Split(editHref, "/")
		id := hrefArray[len(hrefArray)-1]

		// Parse Other Details
		title := rawBook.Find("td:nth-child(3) p a").Text()
		author := rawBook.Find(".catalog_authors li a").Text()
		language := rawBook.Find("td:nth-child(4)").Text()
		series := rawBook.Find("td:nth-child(2)").Text()

		// Create Entry Item
		item := opds.Entry{
			Title:    StripText("[" + fileDesc + "] " + title),
			Language: StripText(language),
			Updated:  &uploadDate,
			Series: []opds.Serie{
				opds.Serie{
					Name: StripText(series),
				},
			},
			Author: []opds.Author{
				opds.Author{
					Name: StripText(author),
				},
			},
			Links: []opds.Link{
				opds.Link{
					Rel:      "http://opds-spec.org/acquisition",
					Href:     "./download?type=fiction&id=" + id,
					TypeLink: mimeMapping[fileType],
				},
			},
		}

		allEntries = append(allEntries, item)
	})

	// Return Results
	return allEntries
}

func ParseLibGenNonFiction(body io.ReadCloser) []opds.Entry {
	// Parse
	defer body.Close()
	doc, _ := goquery.NewDocumentFromReader(body)

	// Normalize Results
	var allEntries []opds.Entry
	doc.Find("table.c tbody > tr:nth-child(n + 2)").Each(func(ix int, rawBook *goquery.Selection) {

		// Parse Type & Size
		fileSize := strings.ToLower(strings.TrimSpace(rawBook.Find("td:nth-child(8)").Text()))
		fileType := strings.ToLower(strings.TrimSpace(rawBook.Find("td:nth-child(9)").Text()))
		fileDesc := fileType + " / " + fileSize

		// Parse MD5
		titleRaw := rawBook.Find("td:nth-child(3) [id]")
		editHref, _ := titleRaw.Attr("href")
		hrefArray := strings.Split(editHref, "?md5=")
		id := hrefArray[1]

		// Parse Other Details
		title := titleRaw.Text()
		author := rawBook.Find("td:nth-child(2)").Text()
		language := rawBook.Find("td:nth-child(7)").Text()
		series := rawBook.Find("td:nth-child(3) [href*='column=series']").Text()

		// Create Entry Item
		item := opds.Entry{
			Title:    StripText("[" + fileDesc + "] " + title),
			Language: StripText(language),
			Series: []opds.Serie{
				opds.Serie{
					Name: StripText(series),
				},
			},
			Author: []opds.Author{
				opds.Author{
					Name: StripText(author),
				},
			},
			Links: []opds.Link{
				opds.Link{
					Rel:      "http://opds-spec.org/acquisition",
					Href:     "./download?type=non-fiction&id=" + id,
					TypeLink: mimeMapping[fileType],
				},
			},
		}

		allEntries = append(allEntries, item)
	})

	// Return Results
	return allEntries
}

func ParseLibGenDownloadURL(body io.ReadCloser) string {
	// Parse
	defer body.Close()
	doc, _ := goquery.NewDocumentFromReader(body)

	// Return Download URL
	// downloadURL, _ := doc.Find("#download [href*=cloudflare]").Attr("href")
	downloadURL, _ := doc.Find("#download h2 a").Attr("href")
	return downloadURL
}

/* -------------------------------------------------------------------------- */
/* ------------------------------- Z-Library  ------------------------------- */
/* -------------------------------------------------------------------------- */

func ParseZLibDownloadURL(body io.ReadCloser) string {
	// Parse
	defer body.Close()
	doc, _ := goquery.NewDocumentFromReader(body)

	// Return Download URL
	downloadURL, _ := doc.Find(".dlButton").Attr("href")
	return "https://usa1lib.org" + downloadURL
}

func ParseZLib(body io.ReadCloser) []opds.Entry {
	// Parse
	defer body.Close()
	doc, _ := goquery.NewDocumentFromReader(body)

	// Normalize Results
	var allEntries []opds.Entry

	doc.Find(".bookRow").Each(func(ix int, rawBook *goquery.Selection) {
		title := rawBook.Find("h3 a").Text()
		author := rawBook.Find(".authors a").Text()

		// Parse Download URL
		infoURL, _ := rawBook.Find("h3 a").Attr("href")
		id := strings.Replace(infoURL, "/book/", "", 1)

		// Parse Type & Size
		fileItemRaw := rawBook.Find(".property__file .property_value")
		fileItem := strings.Replace(fileItemRaw.Text(), ",", "/", 1)
		fileType := strings.ToLower(strings.TrimSpace(strings.Split(fileItem, "/")[0]))
		// fileSize = strings.Split(fileItem, ",")[1].Trim()

		// Create Entry Item
		item := opds.Entry{
			Title: StripText("[" + fileItem + "] " + title),
			Author: []opds.Author{
				opds.Author{
					Name: StripText(author),
				},
			},
			Links: []opds.Link{
				opds.Link{
					Rel:      "http://opds-spec.org/acquisition",
					Href:     "./download?type=zlib&id=" + id,
					TypeLink: mimeMapping[fileType],
				},
			},
		}

		allEntries = append(allEntries, item)
	})

	// Return Results
	return allEntries
}

// Doesnt Work... (Need to parse JSON in const data)
func ParseZLibPopular(body io.ReadCloser) []opds.Entry {
	// Parse
	defer body.Close()
	doc, _ := goquery.NewDocumentFromReader(body)

	// Normalize Results
	var allEntries []opds.Entry

	doc.Find(".brick a").Each(func(ix int, rawBook *goquery.Selection) {
		rawMeta := rawBook.Find(".popular")

		title, _ := rawMeta.Attr("data-title")
		author, _ := rawMeta.Attr("data-author")

		item := opds.Entry{
			Title: StripText(title),
			Author: []opds.Author{
				opds.Author{
					Name: StripText(author),
				},
			},
			Links: []opds.Link{
				opds.Link{
					Href:     "./search?query=" + url.QueryEscape(StripText(title)+" - "+StripText(author)),
					TypeLink: "application/atom+xml;type=feed;profile=opds-catalog",
				},
			},
		}

		allEntries = append(allEntries, item)
	})

	// Return Results
	return allEntries
}

/* -------------------------------------------------------------------------- */
/* ------------------------------- Good Reads ------------------------------- */
/* -------------------------------------------------------------------------- */

func ParseGoodReads(body io.ReadCloser) []opds.Entry {
	// Parse
	defer body.Close()
	doc, _ := goquery.NewDocumentFromReader(body)

	// Normalize Results
	var allEntries []opds.Entry

	doc.Find("[itemtype=\"http://schema.org/Book\"]").Each(func(ix int, rawBook *goquery.Selection) {
		title := rawBook.Find(".bookTitle span").Text()
		author := rawBook.Find(".authorName span").Text()

		item := opds.Entry{
			Title: StripText(title),
			Author: []opds.Author{
				opds.Author{
					Name: StripText(author),
				},
			},
			Links: []opds.Link{
				opds.Link{
					Href:     "./search?query=" + url.QueryEscape(StripText(title)+" - "+StripText(author)),
					TypeLink: "application/atom+xml;type=feed;profile=opds-catalog",
				},
			},
		}

		allEntries = append(allEntries, item)
	})

	// Return Results
	return allEntries
}
