package rockgo

import (
	"testing"
	"log"
	"encoding/json"
	"context"
	"time"
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
	resjson, err := json.Marshal(rssfeed)
	if err == nil {
		log.Println(string(resjson))
	}
	xmlbyte, err := rssfeed.ToRssXml()
	if err != nil {
		log.Fatalln("转成xml失败", err)
	}
	log.Println(string(xmlbyte))
	context.WithTimeout(context.Background(), 5*time.Second)
	ctx, cancelFunc := context.WithCancel()
	cancelFunc()
	select {
	case _, ok := <-ctx.Done():
		if !ok {

		}
	default:
	}
}
