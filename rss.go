package rockgo

import (
	"time"
	"encoding/xml"
	"sync"
	"github.com/jinzhu/copier"
)

var rssfeedPool *sync.Pool = &sync.Pool{New: func() interface{} {
	return &RssFeed{}
}}

func ObtainRssFeed() *RssFeed {
	rssfeed := rssfeedPool.Get().(*RssFeed)
	return rssfeed
}
func RecycleRssfeed(rssfeed *RssFeed) {
	rssfeed.Clear()
	rssfeedPool.Put(rssfeed)
}

type RssFeed struct {
	XMLName xml.Name `xml:"rss" json:"-"`
	Version string   `xml:"version,attr" json:"-"`
	//Channel xml.Name `xml:"rss>channel"`

	Title       string `xml:"channel>title" json:"title"`
	Link        string `xml:"channel>link" json:"link"`
	Description string `xml:"channel>description" json:"description"`

	Category            string `xml:"channel>category,omitempty" json:"category,omitempty"`
	Cloud               *Cloud `xml:"channel>cloud,omitempty" json:"cloud,omitempty"`
	Copyright           string `xml:"channel>copyright,omitempty" json:"copyright,omitempty"`
	Doc                 string `xml:"channel>doc,omitempty" json:"doc,omitempty"`
	Generator           string `xml:"channel>generator,omitempty" json:"generator,omitempty"`
	Image               *Image `xml:"channel>image,omitempty" json:"image,omitempty"`
	Language            string `xml:"channel>language,omitempty" json:"language,omitempty"`
	LastBuildDate       string `xml:"channel>lastBuildDate,omitempty" json:"lastBuildDate,omitempty"`
	LastBuildDateParsed *time.Time `xml:"-" json:"-"`
	ManagingEditor      string `xml:"channel>managingEditor,omitempty" json:"managingEditor,omitempty"`
	PubDate             string `xml:"channel>pubDate,omitempty" json:"pubDate,omitempty"`
	PubDateParsed       *time.Time `xml:"-" json:"-"`
	Rating              string `xml:"channel>rating,omitempty" json:"rating,omitempty"`
	SkipDays            []string `xml:"channel>skipDays,omitempty" json:"skipDays,omitempty"`
	SkipHours           []string `xml:"channel>skipHours,omitempty" json:"skipHours,omitempty"`
	TextInput           *TextInput `xml:"channel>textInput,omitempty" json:"textInput,omitempty"`
	Ttl                 int `xml:"channel>ttl,omitempty" json:"ttl,omitempty"`
	WebMaster           string `xml:"channel>webMaster,omitempty" json:"webMaster,omitempty"`

	Items []*RssItem `xml:"channel>item" json:"items"`
}

type RssItem struct {
	Title       string `xml:"title" json:"title"`
	Link        string `xml:"link" json:"link"`
	Description string `xml:"description" json:"description"`

	Author   string `xml:"author,omitempty" json:"author,omitempty"`
	Category string `xml:"category,omitempty" json:"category,omitempty"`
	Comments string `xml:"comments,omitempty" json:"comments,omitempty"`

	Enclosure     *Enclosure  `xml:"encloure,omitempty" json:"encloure,omitempty"`
	Guid          string `xml:"guid,omitempty" json:"guid,omitempty"`
	PubDate       string `xml:"pubDate,omitempty" json:"pubDate,omitempty"`
	PubDateParsed *time.Time `xml:"-" json:"-"`
	Source        string `xml:"source,omitempty" json:"source,omitempty"`
}

type Cloud struct {
	Domain            string `xml:"domain,attr" json:"domain"`
	Port              string `xml:"port,attr" json:"port"`
	Path              string `xml:"path,attr" json:"path"`
	RegisterProcedure string `xml:"registerProcedure,attr" json:"registerProcedure"`
	Protocol          string `xml:"protocol,attr" json:"protocol"`
}

type Image struct {
	Link  string `xml:"link" json:"link"`
	Title string `xml:"title" json:"title"`
	Url   string `xml:"url" json:"url"`

	Description string `xml:"description,omitempty" json:"description,omitempty"`
	Width       int `xml:"width,omitempty" json:"width,omitempty"`
	Height      int `xml:"height,omitempty" json:"height,omitempty"`
}

type TextInput struct {
	Description string `xml:"description" json:"description"`
	Name        string `xml:"name" json:"name"`
	Link        string `xml:"link" json:"link"`
	Title       string `xml:"title" json:"title"`
}

type Enclosure struct {
	Length int `xml:"length,attr" json:"length"`
	Type   string `xml:"type,attr" json:"type"`
	Url    string `xml:"url,attr" json:"url"`
}

type Source struct {
	Url  string `xml:"url,attr" json:"url"`
	Name string `xml:",innerxml" json:"name"`
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

func (rss *RssFeed) Clear() {
	rss.Cloud = nil
	rss.Image = nil
	rss.LastBuildDateParsed = nil
	rss.PubDateParsed = nil
	rss.TextInput = nil
	rss.Items = nil

	//rss.XMLName = nil
	rss.Version = ""
}

func (rss *RssFeed) Clone() *RssFeed {
	feed2 := RssFeed{}
	copier.Copy(&feed2, rss)
	return &feed2
}
