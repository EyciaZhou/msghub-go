package netTools

import (
	"io/ioutil"
	"net/http"
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
	client = &http.Client{}
)

func GetByAndroid(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	// ...
	req.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 4.3; Nexus 7 Build/JSS15Q) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2307.2 Safari/537.36`)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

var (
	Get = GetByAndroid
)
