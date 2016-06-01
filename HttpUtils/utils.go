package HttpUtils

import (
	"encoding/json"
	"errors"
	"github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	client = &http.Client{
		Timeout: 60 * time.Second,
	}
	_USER_AGENT = `Mozilla/5.0 (Linux; Android 4.3; Nexus 7 Build/JSS15Q) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2307.2 Safari/537.36`
	_TRY_TIME   = 5
)

func tryTimes(funct func() ([]byte, error), trytime int) (_bs []byte, err error) {
	var (
		bs []byte
		e  error
	)

	for i := 1; i <= trytime; i++ {
		bs, e = funct()
		if err == nil {
			return bs, e
		}
	}
	return nil, e
}

func ReadResponse(resp *http.Response) ([]byte, error) {
	return readResponse(resp, nil)
}

func readResponse(resp *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status + " " + (string)(bs))
	}

	return bs, nil
}

func GETRequest(url string, args url.Values) (resp *http.Request, err error) {
	if args != nil {
		url += "?" + args.Encode()
	}
	logrus.Debug("[GET]", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", _USER_AGENT)

	return req, nil
}

func POSTRequest(url string, args url.Values) (resp *http.Request, err error) {
	logrus.Debug("[POST]", url)

	reader := (io.Reader)(nil)
	if args != nil {
		reader = strings.NewReader(args.Encode())
	}
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", _USER_AGENT)

	return req, nil
}

type ErrorChecker func(bs []byte) error

func JsonCheckError(method string, url string, args url.Values, checker ErrorChecker, unmarshalTo interface{}) error {
	var (
		request *http.Request
		err     error
	)

	if method == "GET" {
		request, err = GETRequest(url, args)
	} else if method == "POST" {
		request, err = POSTRequest(url, args)
	} else {
		return errors.New("unsupported method")
	}
	if err != nil {
		return err
	}

	bs, err := tryTimes(func() ([]byte, error) { return readResponse(client.Do(request)) }, _TRY_TIME)
	if err != nil {
		return err
	}

	if checker != nil {
		if err := checker(bs); err != nil {
			return err
		}
	}

	err = json.Unmarshal(bs, unmarshalTo)
	if err != nil {
		return err
	}
	return nil
}

func Json(method string, url string, args url.Values, unmarshalTo interface{}) error {
	return JsonCheckError(method, url, args, nil, unmarshalTo)
}
