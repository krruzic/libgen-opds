// https://github.com/opds-community/libopds2-go/blob/master/opds1/opds1.go
package opds

import (
	"encoding/xml"
	"time"
)

// Feed root element for acquisition or navigation feed
type Feed struct {
	XMLName      xml.Name  `xml:"feed"`
	ID           string    `xml:"id,omitempty",`
	Title        string    `xml:"title,omitempty"`
	Updated      time.Time `xml:"updated,omitempty"`
	Entries      []Entry   `xml:"entry,omitempty"`
	Links        []Link    `xml:"link,omitempty"`
	TotalResults int       `xml:"totalResults,omitempty"`
	ItemsPerPage int       `xml:"itemsPerPage,omitempty"`
}

// Link link to different resources
type Link struct {
	Rel                 string                `xml:"rel,attr"`
	Href                string                `xml:"href,attr"`
	TypeLink            string                `xml:"type,attr"`
	Title               string                `xml:"title,attr,omitempty"`
	FacetGroup          string                `xml:"facetGroup,attr,omitempty"`
	Count               int                   `xml:"count,attr,omitempty"`
	Price               *Price                `xml:"price,omitempty"`
	IndirectAcquisition []IndirectAcquisition `xml:"indirectAcquisition"`
}

// Author represent the feed author or the entry author
type Author struct {
	Name string `xml:"name"`
	URI  string `xml:"uri,omitempty"`
}

// Entry an atom entry in the feed
type Entry struct {
	Title      string     `xml:"title,omitempty"`
	ID         string     `xml:"id,omitempty"`
	Identifier string     `xml:"identifier,omitempty"`
	Updated    *time.Time `xml:"updated,omitempty"`
	Rights     string     `xml:"rights,omitempty"`
	Publisher  string     `xml:"publisher,omitempty"`
	Author     []Author   `xml:"author,omitempty"`
	Language   string     `xml:"language,omitempty"`
	Issued     string     `xml:"issued,omitempty"`
	Published  *time.Time `xml:"published,omitempty"`
	Category   []Category `xml:"category,omitempty"`
	Links      []Link     `xml:"link,omitempty"`
	Summary    *Content   `xml:"summary,omitempty"`
	Content    *Content   `xml:"content,omitempty"`
	Series     []Serie    `xml:"series,omitempty"`
}

// Content content tag in an entry, the type will be html or text
type Content struct {
	Content     string `xml:",cdata"`
	ContentType string `xml:"type,attr"`
}

// Category represent the book category with scheme and term to machine
// handling
type Category struct {
	Scheme string `xml:"scheme,attr"`
	Term   string `xml:"term,attr"`
	Label  string `xml:"label,attr"`
}

// Price represent the book price
type Price struct {
	CurrencyCode string  `xml:"currencycode,attr,omitempty"`
	Value        float64 `xml:",cdata"`
}

// IndirectAcquisition represent the link mostly for buying or borrowing
// a book
type IndirectAcquisition struct {
	TypeAcquisition     string                `xml:"type,attr"`
	IndirectAcquisition []IndirectAcquisition `xml:"indirectAcquisition"`
}

// Serie store serie information from schema.org
type Serie struct {
	Name     string  `xml:"name,attr,omitempty"`
	URL      string  `xml:"url,attr,omitempty"`
	Position float32 `xml:"position,attr,omitempty"`
}
