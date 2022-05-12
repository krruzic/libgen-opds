package client

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"reichard.io/libgen-opds/opds"
)

func GetPage(page string) io.ReadCloser {
	resp, _ := http.Get(page)
	return resp.Body
}

func ParseDownloadURL(body io.ReadCloser) string {
	// Parse
	defer body.Close()
	doc, _ := goquery.NewDocumentFromReader(body)

	// Return Download URL
	downloadURL, _ := doc.Find("#download [href*=cloud]").Attr("href")
	return downloadURL
}

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
			Title: title,
			Author: []opds.Author{
				opds.Author{
					Name: author,
				},
			},
			Links: []opds.Link{
				opds.Link{
					Href:     "./search?query=" + url.QueryEscape(title+" - "+author),
					TypeLink: "application/atom+xml;type=feed;profile=opds-catalog",
				},
			},
		}

		allEntries = append(allEntries, item)
	})

	// Return Results
	return allEntries
}

func ParseLibGen(body io.ReadCloser) []opds.Entry {
	// Parse
	defer body.Close()
	doc, _ := goquery.NewDocumentFromReader(body)

	// Normalize Results
	var allEntries []opds.Entry
	doc.Find("table.catalog tbody > tr").Each(func(ix int, rawBook *goquery.Selection) {

		fileItem := rawBook.Find("td:nth-child(5)")
		// fileType = fileItem.Text().Split("/")[0].Trim().Lower()
		// fileSize = fileItem.Text().Split("/")[1].Trim()

		// Parse Upload Date
		uploadedRaw, _ := fileItem.Attr("title")
		uploadedDateRaw := strings.Split(uploadedRaw, "Uploaded at ")[1]
		uploadDate, _ := time.Parse("2006-01-02 15:04:05", uploadedDateRaw)

		// Parse MD5
		editHref, _ := rawBook.Find("td:nth-child(7) a").Attr("href")
		hrefArray := strings.Split(editHref, "/")
		md5 := hrefArray[len(hrefArray)-1]

		// Create Entry Item
		item := opds.Entry{
			Title:    rawBook.Find("td:nth-child(3)").Text(),
			Language: rawBook.Find("td:nth-child(4)").Text(),
			Updated:  &uploadDate,
			Series: []opds.Serie{
				opds.Serie{
					Name: rawBook.Find("td:nth-child(2)").Text(),
				},
			},
			Author: []opds.Author{
				opds.Author{
					Name: rawBook.Find(".catalog_authors li a").Text(),
				},
			},
			Links: []opds.Link{
				opds.Link{
					Rel:      "acquisition",
					Href:     "./download?md5=" + md5,
					TypeLink: "application/epub+zip",
				},
			},
		}

		allEntries = append(allEntries, item)
	})

	// Return Results
	return allEntries
}
