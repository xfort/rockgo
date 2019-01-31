package rockgo

import (
	"fmt"
	"testing"
)

var rockhttpClient = NewRockHttp()

func TestWeiPHP(t *testing.T) {
	//rockhttpClient.SetProxy("http://127.0.0.1:8081")
	//rockhttpClient.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	//
	//header := http.Header{}
	//header.Set("user-agent", "golang_ua")
	//header.Set("type", "test")
	//_, err, response := rockhttpClient.GetBytes("https://news.uc.cn/", &header)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//log.Println(response.Request.Header)

}
func TestRockHttp_DownloadFile(t *testing.T) {
	urlStr := "http://v3-dy-z.ixigua.com/607e13c6ea605e4e93a01ca1c181c7ee/5c52cfd6/video/m/2204ee0cdddebc842c2a1e776163d70078e11615994a00009f8f8febf987/?rc=M2pmeGVrNG46azMzM2kzM0ApQHRAbzc6NTk0OzgzNDs5OjQ0PDNAKXUpQGczdylAZmh1eXExZnNoaGRmOzRAZHNwNG5eZHI2Xy0tMC0vc3MtbyNvIzUyMDAtLS8tLTIwLjAtLi9pOmItbyM6YC1vI3BiZnJoXitqdDojLy5e%20https://www.iesdouyin.com/share/video/6652278615356476684/?region=CN&mid=6649252307286362888&u_code=i556h6ei&titleType=title"
	outPath, err, response := rockhttpClient.DownloadFile(urlStr, nil, "D:\\go\\code\\my\\src\\github.com\\xfort\\rockgo\\tmp\\test.mp4")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(outPath, response.Status)
}
