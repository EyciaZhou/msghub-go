package generant_utils

import (
	"encoding/json"
	"errors"
	"github.com/EyciaZhou/msghub.go/generant/weibo/types"
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

func TryTimes(funct func() ([]byte, error), trytime int) (_bs []byte, err error) {
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

func handleResponse(resp *http.Response, err error) ([]byte, error) {
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

func GetToBytes(url string, args url.Values) (_bs []byte, err error) {
	return TryTimes(func() ([]byte, error) {
		if args != nil {
			url += "?" + args.Encode()
		}
		logrus.Debug("[GET]", url)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", _USER_AGENT)

		return handleResponse(client.Do(req))
	}, _TRY_TIME)
}

func PostToBytes(url string, args url.Values) (_bs []byte, err error) {
	return TryTimes(func() ([]byte, error) {
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

		return handleResponse(client.Do(req))
	}, _TRY_TIME)
}

func RequestToJson(method string, url string, args url.Values, unmarshalTo interface{}) error {
	var (
		bs  []byte
		err error
	)

	if method == "GET" {
		bs, err = GetToBytes(url, args)
	} else if method == "POST" {
		bs, err = PostToBytes(url, args)
	} else {
		return errors.New("unsupported method")
	}

	if err != nil {
		return err
	}

	//if weibo returns a error
	we := weibo_types.Error{}
	err1 := json.Unmarshal(bs, &we)
	if err1 != nil {
		return err
	}
	if we.Error != "" {
		return errors.New(we.Error)
	}

	err = json.Unmarshal(bs, unmarshalTo)
	if err != nil {
		return err
	}
	return nil
}
