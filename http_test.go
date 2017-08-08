package rockgo

import (
	"testing"
	"log"
	"net/http"
	"crypto/tls"
)

func TestWeiPHP(t *testing.T) {
	rockhttp := NewRockHttp()
	rockhttp.SetProxy("http://127.0.0.1:8081")
	rockhttp.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	header := http.Header{}
	header.Set("user-agent", "golang_ua")
	header.Set("type", "test")
	_, err, response := rockhttp.GetBytes("https://news.uc.cn/", &header)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(response.Request.Header)

}
