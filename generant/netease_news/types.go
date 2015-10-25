package netease_news

type News struct {
	CoverURL string `json:"-"`
	URL      string `json:"-"`
	Priority int    `json:"-"`

	ID         string       `json:"docid"`
	ReplyCount int          `json:"replyCount"`
	Title      string       `json:"title"`
	SubTitle   string       `json:"digest"`
	BoardId    string       `json:"replyBoard"`
	PubTime    string       `json:"ptime"`
	SnapTime   string       `json:"-"`
	Body       string       `json:"body"`
	Images     []*NewsImage `json:"img"`

	Replys   []Reply `json:"-"`
	ViewType int     `json:"-"`
}

type PhotoSet struct {
	Priority   int    `json:"-"`
	ReplyCount int    `json:"-"`
	SnapTime   string `json:"-"`
	Body       string `json:"-"`

	ID       string `json:"postid"`
	CoverURL string `json:"cover"`
	URL      string `json:"url"`
	Title    string `json:"setname"`
	SubTitle string `json:"desc"`
	BoardId  string `json:"boardid"`
	PubTime  string `json:"createdate"`

	Images []*PhototSetImage `json:"photos"`

	Replys   []Reply `json:"-"`
	ViewType int     `json:"-"`
}

type PhototSetImage struct {
	Desc string `json:"note"`
	URL  string `json:"imgurl"`
}

type ReplyFloor struct {
	//Floor   int    `json:"-"`
	Time    string `json:"t"`
	Name    string `json:"f"`
	Content string `json:"b"`
	Digg    string `json:"v"`
}

type Reply map[string]*ReplyFloor

type reply_tmp struct {
	HotPosts []Reply `json:"hotPosts"`
}

type NewsImage struct {
	Ref   string `json:"ref"` //not set if not have
	Size  string `json:"pixel"`
	Title string `json:"alt"`
	URL   string `json:"src"`
}
