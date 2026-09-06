package rss

import (
	"encoding/xml"
	"slices"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`

	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`

	Items []Item `xml:"item,omitempty"`

	Language       string   `xml:"language,omitempty"`
	Copyright      string   `xml:"copyright,omitempty"`
	ManagingEditor string   `xml:"managingEditor,omitempty"`
	WebMaster      string   `xml:"webMaster,omitempty"`
	PubDate        string   `xml:"pubDate,omitempty"`
	LastBuildDate  string   `xml:"lastBuildDate,omitempty"`
	Categories     []string `xml:"category,omitempty"`
	Generator      string   `xml:"generator,omitempty"`
	Docs           string   `xml:"docs,omitempty"`
	Cloud          string   `xml:"cloud,omitempty"`
	Ttl            int      `xml:"ttl,omitempty"`
	Image          struct {
		Url   string `xml:"url,omitempty"`
		Title string `xml:"title,omitempty"`
		Link  string `xml:"link,omitempty"`

		Width  int `xml:"width,omitempty"`
		Height int `xml:"height,omitempty"`
	} `xml:"image,omitempty"`

	Rating    string `xml:"rating,omitempty"`
	TextInput struct {
		Title       string `xml:"title,omitempty"`
		Description string `xml:"description,omitempty"`
		Name        string `xml:"name,omitempty"`
		Link        string `xml:"link,omitempty"`
	} `xml:"textInput,omitempty"`
	Hour     []int  `xml:"skipHours>hour,omitempty"`
	SkipDays string `xml:"skipDays,omitempty"`
}

type Item struct {
	// Either title or description MUST be set
	Title       string `xml:"title,omitempty"`
	Description string `xml:"description,omitempty"`

	Author    string   `xml:"author,omitempty"`
	Category  []string `xml:"category,omitempty"`
	Comments  string   `xml:"comments,omitempty"`
	Enclosure string   `xml:"enclosure,omitempty"`
	Guid      string   `xml:"guid,omitempty"`
	Link      string   `xml:"link,omitempty"`
	PubDate   string   `xml:"pubDate,omitempty"`
	Source    string   `xml:"source,omitempty"`
}

func NewFeed() RSS {
	r := RSS{
		XMLName: xml.Name{Local: "rss"},
		Version: "2.0",
		Channel: Channel{
			Generator: "Locket v0.2",
		},
	}
	return r
}

func FeedFromXML(x []byte) (RSS, error) {
	var r RSS
	err := xml.Unmarshal(x, &r)
	return r, err
}

func (r RSS) GetXML() ([]byte, error) {
	return xml.MarshalIndent(r, "", "\t")
}

func (r *RSS) AddItem(raw []byte) error {
	var i Item
	err := xml.Unmarshal(raw, &i)
	if err == nil {
		r.Channel.Items = append(r.Channel.Items, i)
	}
	return err
}

func (r *RSS) DeleteItem(idx int) {
	if idx < 0 || idx >= len(r.Channel.Items) {
		return
	}
	r.Channel.Items = slices.Delete(r.Channel.Items, idx, idx)
}
