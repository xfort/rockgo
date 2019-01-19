package rockgo

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"testing"
)

var rockhttpClient = NewRockHttp()

func TestWeiPHP(t *testing.T) {
	rockhttpClient.SetProxy("http://127.0.0.1:8081")
	rockhttpClient.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	header := http.Header{}
	header.Set("user-agent", "golang_ua")
	header.Set("type", "test")
	_, err, response := rockhttpClient.GetBytes("https://news.uc.cn/", &header)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(response.Request.Header)

}
func TestDouYin_Search(t *testing.T) {
	urlValue := url.Values{}
	urlValue.Set("keyword", "创意广告")
	urlValue.Set("offset", "0")
	urlValue.Set("count", "10")
	urlValue.Set("is_pull_refresh", "0")
	urlValue.Set("hot_search", "0")
	urlValue.Set("latitude", "0.0")
	urlValue.Set("longitude", "0.0")
	urlValue.Set("ts", "1547705335")
	urlValue.Set("js_sdk_version", "1.6.4")
	urlValue.Set("app_type", "normal")
	urlValue.Set("os_api", "23")
	urlValue.Set("device_type", "Redmi 4")
	urlValue.Set("device_platform", "android")
	urlValue.Set("ssmix", "a")
	urlValue.Set("iid", "56312689392")
	urlValue.Set("manifest_version_code", "390")
	urlValue.Set("dpi", "480")
	urlValue.Set("version_code", "390")
	urlValue.Set("app_name", "aweme")
	urlValue.Set("version_name", "3.9.0")
	urlValue.Set("openudid", "43b80838cd18aed6")
	urlValue.Set("device_id", "52909728013")
	urlValue.Set("resolution", "1080*1920")
	urlValue.Set("os_version", "6.0.1")
	urlValue.Set("language", "zh")
	urlValue.Set("device_brand", "Xiaomi")
	urlValue.Set("ac", "wifi")
	urlValue.Set("update_version_code", "3902")
	urlValue.Set("aid", "1128")
	urlValue.Set("channel", "aweGW")
	urlValue.Set("_rticket", "1547705336378")
	urlValue.Set("mcc_mnc", "46000")
	urlValue.Set("as", "a1e5314457df2c1b006622")
	urlValue.Set("cp", "1dfdc2527f0342bfe1_cMg")
	urlValue.Set("mas", "01603ed35eaf87727f871abe6d1aa993244c4c6c6c0c0c468cc64c")

}
