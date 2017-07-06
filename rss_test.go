package rockgo

import (
	"testing"
	"log"
)

func TestRssFeed(t *testing.T) {

	rockhttp := NewRockHttp()
	resBytes, err, _ := rockhttp.GetBytes("http://x.xiushuang.com/client/Sitemap/rss", nil)
	if err != nil {
		log.Fatalln(err)
	}
	rssfeed := &RssFeed{}
	err = rssfeed.Parse(resBytes)
	if err != nil {
		log.Fatalln(err)
	}
	//resjson, err := json.Marshal(rssfeed)

	xmlbyte, err := rssfeed.ToRssXml()
	if err != nil {
		log.Fatalln("转成xml失败", err)
	}
	log.Println(string(xmlbyte))
}
