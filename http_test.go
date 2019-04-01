package rockgo

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
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
	//urlStr := "http://v3-dy-z.ixigua.com/607e13c6ea605e4e93a01ca1c181c7ee/5c52cfd6/video/m/2204ee0cdddebc842c2a1e776163d70078e11615994a00009f8f8febf987/?rc=M2pmeGVrNG46azMzM2kzM0ApQHRAbzc6NTk0OzgzNDs5OjQ0PDNAKXUpQGczdylAZmh1eXExZnNoaGRmOzRAZHNwNG5eZHI2Xy0tMC0vc3MtbyNvIzUyMDAtLS8tLTIwLjAtLi9pOmItbyM6YC1vI3BiZnJoXitqdDojLy5e%20https://www.iesdouyin.com/share/video/6652278615356476684/?region=CN&mid=6649252307286362888&u_code=i556h6ei&titleType=title"
	//outPath, err, response := rockhttpClient.DownloadFile(urlStr, nil, "D:\\go\\code\\my\\src\\github.com\\xfort\\rockgo\\tmp\\test.mp4")
	//if err != nil {
	//	t.Error(err)
	//}
	//fmt.Println(outPath, response.Status)
}

func TestProxy(t *testing.T) {
	fileData, err := ioutil.ReadFile("D:\\go\\code\\my\\src\\github.com\\xfort\\rockgo\\tmp\\all.txt")
	if err != nil {
		t.Fatal(err)
	}
	rockhttpClient.Timeout = 5 * time.Second
	fileScanner := bufio.NewScanner(bytes.NewReader(fileData))
	index := 0
	httpHeader := http.Header{}
	httpHeader.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
	for fileScanner.Scan() {
		ipaddr := fileScanner.Text()
		if ipaddr == "" {
			continue
		}
		proxyUrl := `socks5://`+ipaddr
		index++
		log.Println(index, proxyUrl)
		rockhttpClient.SetProxy(proxyUrl)
		targetUrl := "https://www.youtube.com/"
		startUTC := time.Now().Unix()

		_, err, resObj := rockhttpClient.GetBytes(targetUrl, &httpHeader)
		if resObj != nil && resObj.StatusCode == http.StatusOK {
			log.Println(proxyUrl, "__", time.Now().Unix()-startUTC)
//log.Println(string(data))
		}

		if err != nil {
			//log.Println(err)
		}
	}

}
