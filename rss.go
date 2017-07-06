package rockgo

import (
	"time"
	"encoding/xml"
)

type RssFeed struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	//Channel xml.Name `xml:"rss>channel"`

	Title       string `xml:"channel>title"`
	Link        string `xml:"channel>link"`
	Description string `xml:"channel>description"`

	Category            string `xml:"channel>category,omitempty"`
	Cloud               *Cloud `xml:"channel>cloud,omitempty"`
	Copyright           string `xml:"channel>copyright,omitempty"`
	Doc                 string `xml:"channel>doc,omitempty"`
	Generator           string `xml:"channel>generator,omitempty"`
	Image               *Image `xml:"channel>image,omitempty"`
	Language            string `xml:"channel>language,omitempty"`
	LastBuildDate       string `xml:"channel>lastBuildDate,omitempty"`
	LastBuildDateParsed *time.Time `xml:"-"`
	ManagingEditor      string `xml:"channel>managingEditor,omitempty"`
	PubDate             string `xml:"channel>pubDate,omitempty"`
	PubDateParsed       *time.Time `xml:"-"`
	Rating              string `xml:"channel>rating,omitempty"`
	SkipDays            []string `xml:"channel>skipDays,omitempty"`
	SkipHours           []string `xml:"channel>skipHours,omitempty"`
	TextInput           *TextInput `xml:"channel>textInput,omitempty"`
	Ttl                 int `xml:"channel>ttl,omitempty"`
	WebMaster           string `xml:"channel>webMaster,omitempty"`

	Items []*RssItem `xml:"channel>item"`
}

type RssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`

	Author   string `xml:"author,omitempty"`
	Category string `xml:"category,omitempty"`
	Comments string `xml:"comments,omitempty"`

	Enclosure     *Enclosure  `xml:"encloure,omitempty"`
	Guid          string `xml:"guid,omitempty"`
	PubDate       string `xml:"pubDate,omitempty"`
	PubDateParsed *time.Time `xml:"_"`
	Source        string `xml:"source,omitempty"`
}

type Cloud struct {
	Domain            string `xml:"domain,attr"`
	Port              string `xml:"port,attr"`
	Path              string `xml:"path,attr"`
	RegisterProcedure string `xml:"registerProcedure,attr"`
	Protocol          string `xml:"protocol,attr"`
}

type Image struct {
	Link  string `xml:"link"`
	Title string `xml:"title"`
	Url   string `xml:"url"`

	Description string `xml:"description,omitempty"`
	Width       int `xml:"width,omitempty"`
	Height      int `xml:"height,omitempty"`
}

type TextInput struct {
	Description string `xml:"description"`
	Name        string `xml:"name"`
	Link        string `xml:"link"`
	Title       string `xml:"title"`
}

type Enclosure struct {
	Length int `xml:"length,attr"`
	Type   string `xml:"type,attr"`
	Url    string `xml:"url,attr"`
}

type Source struct {
	Url  string `xml:"url,attr"`
	Name string `xml:",innerxml"`
}

func (rss *RssFeed) Parse(dataBytes []byte) error {
	err := xml.Unmarshal(dataBytes, &rss)
	return err
}

func (rss *RssFeed) ToRssXml() ([]byte, error) {
	xmlByte, err := xml.Marshal(rss)
	if err != nil {
		return nil, err
	}

	xmlByte = append([]byte(xml.Header), xmlByte...)
	return xmlByte, nil
}
