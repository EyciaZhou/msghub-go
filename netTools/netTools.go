package netTools

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func GetWithoutAnyThing(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

var (
	client = &http.Client{
		Timeout: time.Duration(15 * time.Second),
	}
)

func GetByAndroid(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 4.3; Nexus 7 Build/JSS15Q) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2307.2 Safari/537.36`)

	var resp *http.Response

	for try := 0; try < 5; try++ {
		resp, err = client.Do(req)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func GetReturnsReader(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 4.3; Nexus 7 Build/JSS15Q) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2307.2 Safari/537.36`)

	var resp *http.Response

	for try := 0; try < 5; try++ {
		resp, err = client.Do(req)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

var (
	Get = GetByAndroid
)