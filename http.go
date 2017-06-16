package rockgo

import (
	"net/http"
	"io/ioutil"
	"os"
	"io"
	"net/url"
	"strings"
	"net"
	"time"

	"golang.org/x/net/proxy"
)

const (
	UserAgent_Android = "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.23 Mobile Safari/537.36"

	UserAgent_IOS = "Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1"

	UserAgent_Chrome_Web = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36"
)

func NewRockHttp() *RockHttp {
	rockhttp := &RockHttp{}

	return rockhttp
}

type RockHttp struct {
	http.Client
}

func (rockhttp *RockHttp) LoadResponse(request *http.Request) (*http.Response, error) {
	response, err := rockhttp.Do(request)

	//if response != nil {
	//	defer response.Body.Close()
	//}
	if err != nil {
		if response != nil {
			response.Body.Close()
		}
		return nil, err
	}

	return response, nil
}

func (rockhttp *RockHttp) DoRequest(method string, urlstr string, header *http.Header, body io.Reader) ([]byte, error, *http.Response) {
	request, err := http.NewRequest(method, urlstr, body)
	if err != nil {
		return nil, err, nil
	}
	return rockhttp.DoRequestBytes(request)
}

func (rockhttp *RockHttp) DoRequestBytes(request *http.Request) ([]byte, error, *http.Response) {
	response, err := rockhttp.Do(request)

	if response != nil {
		defer response.Body.Close()
	}

	if err != nil {
		return nil, err, response
	}

	resByte, err := ioutil.ReadAll(response.Body)
	return resByte, err, response
}

func (rockHttp *RockHttp) DoRequestFile(request *http.Request, filepath string) (string, error, *http.Response) {

	response, err := rockHttp.Do(request)

	if response != nil {
		defer response.Body.Close()
	}

	if err != nil {
		return "", err, response
	}

	outFile, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0644)
	defer outFile.Close()

	if err != nil {
		return "", err, response
	}

	_, err = io.Copy(outFile, response.Body)

	if err != nil {
		return "", err, response
	}
	return filepath, nil, response
}

func (rockhttp *RockHttp) DownloadFile(urlStr string, header *http.Header, filepath string) (string, error, *http.Response) {

	request, err := http.NewRequest("GET", urlStr, nil)

	if err != nil {
		return "", err, nil
	}
	if header != nil {
		request.Header = *header
	}
	return rockhttp.DoRequestFile(request, filepath)
}

func (rockhttp *RockHttp) PostForm(urlStr string, header *http.Header, data url.Values) ([]byte, error, *http.Response) {
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err, nil
	}
	if header != nil {
		req.Header = *header
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return rockhttp.DoRequestBytes(req)
}

func (rockhttp *RockHttp) PostData(urlStr string, header *http.Header, body io.Reader) ([]byte, error, *http.Response) {
	req, err := http.NewRequest("POST", urlStr, body)
	if err != nil {
		return nil, err, nil
	}

	if header != nil {
		req.Header = *header
	}
	return rockhttp.DoRequestBytes(req)
}

func (rockhttp *RockHttp) GetBytes(urlstr string, header *http.Header) ([]byte, error, *http.Response) {
	request, err := http.NewRequest("GET", urlstr, nil)

	if err != nil {
		return nil, err, nil
	}
	if header != nil {
		request.Header = *header
	}

	return rockhttp.DoRequestBytes(request)
}

//读取最终重定向 http地址
func (rockhttp *RockHttp) GetRediectUrl(urlstr string, header *http.Header) (string, error, *http.Response) {
	request, err := http.NewRequest("GET", urlstr, nil)

	if err != nil {
		return "", err, nil
	}
	if header != nil {
		request.Header = *header
	}

	response, err := rockhttp.Do(request)

	if response != nil {
		defer response.Body.Close()
	}

	if err != nil {
		return "", err, response
	}
	return response.Request.URL.String(), nil, response
}

func (rockhttp *RockHttp) SetProxy(urlStr string) error {
	proxyUrl, err := url.Parse(urlStr)
	if err != nil {
		//log.Println("SetProxy", err, proxyUrl)
		return err
	}
	transport := rockhttp.Transport.(*http.Transport)
	transport.Proxy = http.ProxyURL(proxyUrl)
	return nil
}

func (rockhttp *RockHttp) SetSocksProxy(urlStr string) error {
	proxyurl, err := url.Parse(urlStr)
	if err != nil {
		//log.Println("SetSocketProxy()", urlStr, err)
		return err
	}

	netDialer := &net.Dialer{
		Timeout:   60 * time.Second,
		KeepAlive: 60 * time.Second,
	}

	dialer, err := proxy.FromURL(proxyurl, netDialer)
	//dialer, err := 	proxy.SOCKS5("tcp", urlStr, nil, proxy.Direct)
	if err != nil {
		//log.logPrintln("can't connect to the proxy:", err)
		return err
	}
	transport := rockhttp.Transport.(*http.Transport)
	transport.Dial = dialer.Dial
	transport.DialContext = nil

	transport.TLSHandshakeTimeout = 60 * time.Second
	return nil
}
