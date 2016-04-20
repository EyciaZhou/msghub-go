package weibo

import (
	"encoding/json"
	"errors"
	"fmt"
	"git.eycia.me/eycia/msghub/generant"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
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

func ToJson(method string, url string, args url.Values, unmarshalTo interface{}) error {
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
	we := weiboError{}
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

type weiboUser struct {
	Name       string `json:"name"`
	Idstr      string `json:"idstr"`
	CoverImage string `json:"cover_image"`
}

type weiboPicUrl struct {
	ThumbnailPic string `json:"thumbnail_pic"`
}

type weiboPicUrls []weiboPicUrl

type weiboTweet struct {
	CreatedAt      string       `json:"created_at"`
	Id             int64        `json:"id"`
	Mid            string       `json:"mid"`
	Idstr          string       `json:"idstr"`
	Text           string       `json:"text"`
	PicUrls        weiboPicUrls `json:"pic_urls"`
	ThumbnailPic   string       `json:"thumbnail_pic"`
	RepostsCount   int64        `json:"reposts_count"`
	CommentsCount  int64        `json:"comments_count"`
	AttitudesCount int64        `json:"attitudes_count"`

	RetweetedStatus *weiboTweet `json:"retweeted_status"`

	User *weiboUser `json:"user"`
}

type weiboTimeline struct {
	Statuses []*weiboTweet `json:"statuses"`
}

type weiboError struct {
	Error     string `json:"error"`
	ErrorCode int    `json:"error_code"`
	Request   string `json:"request"`
}

var (
	weiboTimeLayout = time.RubyDate
)

const (
	BASE62_KEYBOARD = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func (p *weiboTweet) getMidBase62() string {
	mid := p.Id

	result := ""

	for mid > 0 {
		k := mid % 10000000
		mid /= 10000000

		for k > 0 {
			result = string(BASE62_KEYBOARD[k%62]) + result
			k /= 62
		}
	}

	return result
}

func (p *weiboTweet) GetSource() string {
	return "http://weibo.com/" + p.User.Idstr + "/" + p.getMidBase62()
}

func (p weiboPicUrls) Convert() []*generant.Image {
	if p == nil {
		return nil
	}

	result := make([]*generant.Image, len(p))
	for index, img_from := range p {
		result[index] = &generant.Image{
			URL: strings.Replace(img_from.ThumbnailPic, "thumbnail", "large", 1), //covert to big picture
		}
	}

	return result
}

func (p *weiboUser) Convert() (*generant.Author, error) {
	return &generant.Author{
		Name: p.Name,
		Uid: "weibo_" + p.Idstr,
		CovertSource: p.CoverImage,
	}
}

func (p *weiboTweet) Convert() (*generant.Message, error) {
	result := &generant.Message{}

	result.SnapTime = time.Now().Unix()

	pubtime, err := time.Parse(weiboTimeLayout, p.CreatedAt)
	if err != nil {
		return nil, err
	}
	result.PubTime = pubtime.Unix()
	result.Source = p.GetSource()
	result.Body = p.Text
	result.Title = ""                //no title, replace with author's name
	result.Subtitle = p.Text         //display on first screen
	result.CoverImg = p.ThumbnailPic //"" if not have
	result.Images = p.PicUrls.Convert()
	result.ReplyNumber = p.CommentsCount
	result.Replys = nil	//TODO
	result.ViewType = generant.VIEW_TYPE_PICTURES
	result.Topic = "weibo_" + p.User.Idstr
	result.Author = p.User.Convert()

	return result, nil
}

func Oauth2_GetTokenInfo() error {
	var Json map[string]interface{}

	err := ToJson("POST", URL_GET_TOKEN_INFO, url.Values{
		"access_token": {TOKEN},
	}, &Json)
	if err != nil {
		panic(err)
	}
	fmt.Println(Json)
	return nil
}

func FriendsTimelineFirstPage() ([]*weiboTweet, error) {
	args := url.Values{
		"access_token": {TOKEN},
	}

	wtl := weiboTimeline{}

	err := ToJson("GET", URL_FRIENDS_TIMELINE, args, &wtl)

	if err != nil {
		return nil, err
	}

	return wtl.Statuses, nil
}
