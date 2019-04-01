package rockgo

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"bytes"
	"context"
	"fmt"
	"golang.org/x/net/proxy"
	"mime/multipart"
	"net/textproto"
	"path/filepath"
)

const (
	UserAgent_Android = "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.23 Mobile Safari/537.36"

	UserAgent_IOS = "Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1"

	UserAgent_Chrome_Web = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36"
)

func NewRockHttp() *RockHttp {
	rockhttp := &RockHttp{}
	rockhttp.Transport = http.DefaultTransport
	//rockhttp.Timeout = 60 * time.Second

	return rockhttp
}

type RockHttp struct {
	http.Client
}

func (rockhttp *RockHttp) LoadResponseCtx(ctx context.Context, request *http.Request) (*http.Response, error) {
	return rockhttp.LoadResponse(request.WithContext(ctx))
}

func (rockhttp *RockHttp) LoadResponse(request *http.Request) (*http.Response, error) {

	response, err := rockhttp.Do(request)
	if err != nil {
		if response != nil {
			response.Body.Close()
		}
		return nil, err
	}
	return response, nil
}

func (rockhttp *RockHttp) DoRequestCtx(ctx context.Context, method string, urlstr string, header *http.Header, body io.Reader) ([]byte, error, *http.Response) {
	request, err := http.NewRequest(method, urlstr, body)
	if err != nil {
		return nil, err, nil
	}
	if header != nil {
		for key, _ := range *header {
			request.Header.Set(key, header.Get(key))
		}
		//request.Header = *header
	}
	return rockhttp.DoRequestBytes(request.WithContext(ctx))
}

func (rockhttp *RockHttp) DoRequest(method string, urlstr string, header *http.Header, body io.Reader) ([]byte, error, *http.Response) {

	return rockhttp.DoRequestCtx(context.TODO(), method, urlstr, header, body)
}

func (rockhttp *RockHttp) DoRequestBytesCtx(ctx context.Context, request *http.Request) ([]byte, error, *http.Response) {
	return rockhttp.DoRequestBytes(request.WithContext(ctx))
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

func (rockHttp *RockHttp) DoRequestFileCtx(ctx context.Context, request *http.Request, filepath string) (string, error, *http.Response) {
	return rockHttp.DoRequestFile(request.WithContext(ctx), filepath)
}

func (rockHttp *RockHttp) DoRequestFile(request *http.Request, outFilePath string) (string, error, *http.Response) {

	response, err := rockHttp.Do(request)

	if response != nil {
		defer response.Body.Close()
	}

	if err != nil {
		return "", err, response
	}

	outFilePathTemp := outFilePath + ".download"

	outFile, err := os.OpenFile(outFilePathTemp, os.O_CREATE|os.O_RDWR, 0644)
	//defer outFile.Close()

	if err != nil {
		outFile.Close()
		return "", err, response
	}

	_, err = io.Copy(outFile, response.Body)
	outFile.Close()
	if err != nil {
		return "", err, response
	}
	err = os.Rename(outFilePathTemp, outFilePath)
	if err != nil {
		return "", err, response
	}
	return outFilePath, nil, response
}

func (rockhttp *RockHttp) DownloadFileCtx(ctx context.Context, urlStr string, header *http.Header, filepath string) (string, error, *http.Response) {

	request, err := http.NewRequest("GET", urlStr, nil)

	if err != nil {
		return "", err, nil
	}
	if header != nil {
		request.Header = *header
	}
	return rockhttp.DoRequestFile(request.WithContext(ctx), filepath)
}

func (rockhttp *RockHttp) DownloadFile(urlStr string, header *http.Header, filepath string) (string, error, *http.Response) {
	return rockhttp.DownloadFileCtx(context.TODO(), urlStr, header, filepath)
}

func (rockhttp *RockHttp) PostFormCtx(ctx context.Context, urlStr string, header *http.Header, data url.Values) ([]byte, error, *http.Response) {
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err, nil
	}
	if header != nil {
		req.Header = *header
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	}
	return rockhttp.DoRequestBytes(req.WithContext(ctx))
}

func (rockhttp *RockHttp) PostForm(urlStr string, header *http.Header, data url.Values) ([]byte, error, *http.Response) {
	return rockhttp.PostFormCtx(context.TODO(), urlStr, header, data)
}

func (rockhttp *RockHttp) PostJson(ctx context.Context, urlStr string, header *http.Header, obj interface{}) ([]byte, error, *http.Response) {
	objJson, err := json.Marshal(obj)
	if err != nil {
		return nil, err, nil
	}
	req, err := http.NewRequest("POST", urlStr, bytes.NewReader(objJson))
	if err != nil {
		return nil, err, nil
	}
	if header != nil {
		req.Header = *header
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	}
	return rockhttp.DoRequestBytes(req.WithContext(ctx))

}
func (rockhttp *RockHttp) PostDataCtx(ctx context.Context, urlStr string, header *http.Header, body io.Reader) ([]byte, error, *http.Response) {
	req, err := http.NewRequest("POST", urlStr, body)
	if err != nil {
		return nil, err, nil
	}

	if header != nil {
		for key, _ := range *header {
			req.Header.Set(key, header.Get(key))
		}
	}
	return rockhttp.DoRequestBytes(req.WithContext(ctx))
}

func (rockhttp *RockHttp) PostData(urlStr string, header *http.Header, body io.Reader) ([]byte, error, *http.Response) {
	return rockhttp.PostDataCtx(context.TODO(), urlStr, header, body)
}

func (rockhttp *RockHttp) GetBytesCtx(ctx context.Context, urlstr string, header *http.Header) ([]byte, error, *http.Response) {
	request, err := http.NewRequest("GET", urlstr, nil)

	if err != nil {
		return nil, err, nil
	}
	if header != nil {
		for key, _ := range *header {
			request.Header.Set(key, header.Get(key))
		}

	}
	return rockhttp.DoRequestBytes(request.WithContext(ctx))
}

func (rockhttp *RockHttp) GetBytes(urlstr string, header *http.Header) ([]byte, error, *http.Response) {
	return rockhttp.GetBytesCtx(context.TODO(), urlstr, header)
}

func (rockhttp *RockHttp) GetRediectUrlCtx(ctx context.Context, urlstr string, header *http.Header) (string, error, *http.Response) {
	request, err := http.NewRequest("GET", urlstr, nil)

	if err != nil {
		return "", err, nil
	}
	if header != nil {
		request.Header = *header
	}

	response, err := rockhttp.Do(request.WithContext(ctx))

	if response != nil {
		defer response.Body.Close()
	}

	if err != nil {
		return "", err, response
	}
	return response.Request.URL.String(), nil, response
}

//读取最终重定向 http地址
func (rockhttp *RockHttp) GetRediectUrl(urlstr string, header *http.Header) (string, error, *http.Response) {
	return rockhttp.GetRediectUrlCtx(context.TODO(), urlstr, header)
}

func (rockhttp *RockHttp) SetProxy(urlStr string) error {
	proxyUrl, err := url.Parse(urlStr)
	if err != nil {
		//log.Println("SetProxy", err, proxyUrl)
		return err
	}

	transport := rockhttp.Transport.(*http.Transport)
	transport.Proxy = http.ProxyURL(proxyUrl)
	//transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
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

	transport.TLSHandshakeTimeout = 10 * time.Second
	return nil
}

func (rockhttp *RockHttp) DoUploadFile(urlstr string, filepathStr string, fileFieldname string, fileContentType string, header http.Header, formdata url.Values) ([]byte, error, *http.Response) {

	fileObj, err := os.Open(filepathStr)

	if err != nil {
		return nil, err, nil
	}
	defer fileObj.Close()

	fileContent, err := ioutil.ReadAll(fileObj)
	if err != nil {
		return nil, err, nil
	}

	bodyBuffer := &bytes.Buffer{}
	multiWriter := multipart.NewWriter(bodyBuffer)

	if formdata != nil && len(formdata) > 0 {
		for key, _ := range formdata {
			multiWriter.WriteField(key, formdata.Get(key))
		}
	}
	//writerObj, err := multiWriter.CreateFormFile(fileFieldname, filepath.Base(filepathStr))

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, quoteEscaper.Replace(fileFieldname), quoteEscaper.Replace(filepath.Base(filepathStr))))
	h.Set("Content-Type", fileContentType)

	writerObj, err := multiWriter.CreatePart(h)

	if err != nil {
		return nil, err, nil
	}

	_, err = writerObj.Write(fileContent)
	if err != nil {
		//log.Println("file _length", filelength, contentType)
		return nil, err, nil
	}
	contentType := multiWriter.FormDataContentType()
	multiWriter.Close()

	if header == nil {
		header = http.Header{}
	}
	header.Set("Content-Type", contentType)

	return rockhttp.DoRequestCtx(context.TODO(), "POST", urlstr, &header, bodyBuffer)
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

type RockHttpIn interface {
	DoRequest(ctx context.Context, url string, method string, header http.Header, body io.Reader) (error, *http.Response)

	DoRequestBytes(ctx context.Context, url string, method string, header http.Header, body io.Reader) ([]byte, error, *http.Response)

	DoDownloadFile(ctx context.Context, url string, method string, header http.Header, body io.Reader, targetfilepath string) (string, error, *http.Response)

	DoUploadFile(ctx context.Context, url string, method string, header http.Header, httpFile *http.File) (error, *http.Response)
	DoUploadFileBytes(ctx context.Context, url string, method string, header http.Header, httpFile *http.File) ([]byte, error, *http.Response)
}

type HttpFile struct {
	LocalFile       string
	FileFieldname   string
	FileName        string
	FileContentType string
}
