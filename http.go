package rockgo

import (
	"net/http"
	"io/ioutil"
	"os"
	"io"
	"net/url"
	"strings"
)

type RockHttp struct {
	http.Client
}

func (rockhttp *RockHttp) LoadResponse(request *http.Request) (*http.Response, error) {
	response, err := rockhttp.Do(request)

	//if response != nil {
	//	defer response.Body.Close()
	//}
	if err != nil {
		if (response != nil) {
			response.Body.Close()
		}
		return nil, err
	}

	return response, nil
}

func (rockhttp *RockHttp) DoRequestBytes(request *http.Request) ([]byte, error) {
	response, err := rockhttp.Do(request)

	if response != nil {
		defer response.Body.Close()
	}

	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(response.Body)
}

func (rockHttp *RockHttp) DoRequestFile(request *http.Request, filepath string) (string, error) {

	response, err := rockHttp.Do(request)

	if response != nil {
		defer response.Body.Close()
	}

	if err != nil {
		return "", err
	}

	//if response.StatusCode < 200 || response.StatusCode >= 300 {
	//	return "", errors.New("StatusCode=" + response.Status + "_url=" + request.URL.String())
	//}

	outFile, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0644)
	defer outFile.Close()

	if err != nil {
		return "", err
	}

	_, err = io.Copy(outFile, response.Body)

	if err != nil {
		return "", err
	}
	return filepath, nil
}

func (rockhttp *RockHttp) DownloadFile(urlStr string, header *http.Header, filepath string) (string, error) {

	request, err := http.NewRequest("GET", urlStr, nil)

	if err != nil {
		return "", err
	}
	if header != nil {
		request.Header = *header
	}
	return rockhttp.DoRequestFile(request, filepath)
}

func (rockhttp *RockHttp) PostForm(urlStr string, header *http.Header, data url.Values) ([]byte, error) {
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	if header != nil {
		req.Header = *header
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return rockhttp.DoRequestBytes(req)
}

func (rockhttp *RockHttp) PostData(urlStr string, header *http.Header, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest("POST", urlStr, body)
	if err != nil {
		return nil, err
	}

	if header != nil {
		req.Header = *header
	}
	return rockhttp.DoRequestBytes(req)
}

//读取最终重定向 http地址
func (rockhttp *RockHttp) GetRediectUrl(urlstr string, header *http.Header) (string, error) {
	request, err := http.NewRequest("GET", urlstr, nil)

	if err != nil {
		return "", err
	}
	if header != nil {
		request.Header = *header
	}

	response, err := rockhttp.Do(request)

	if response != nil {
		defer response.Body.Close()
	}

	if err != nil {
		return "", err
	}
	return response.Request.URL.String(), nil
}
