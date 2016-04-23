package weibo_utils

import (
	"encoding/json"
	"errors"
	"github.com/EyciaZhou/msghub.go/generant/weibo/types"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GetToBytes(url string, args url.Values) (_bs []byte, err error) {
	url_full := url + "?" + args.Encode()
	logrus.Debug("[GET]", url_full)

	resp, err := http.Get(url_full)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}

func PostToBytes(url string, args url.Values) (_bs []byte, err error) {
	logrus.Debug("[GET]", url)
	resp, err := http.PostForm(url, args)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
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
