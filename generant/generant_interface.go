package generant

type Generant interface {
}

type ReplyFloor struct {
	Floor   int    `json:"floor"`
	Time    int64  `json:"time"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Digg    int    `jsong:"digg"`
}

type Reply []*ReplyFloor

type Image struct {
	Ref  string `json:"ref"` //not set if not have
	Desc string `json:"desc"`
	URL  string `json:"url"`
}

type Message struct {
	SnapTime    int64    `json:"snaptime"`
	PubTime     int64    `json:"pubtime"`
	Source      string   `json:"source"`
	Body        string   `json:"body"`
	Title       string   `json:"title"`
	Subtitle    string   `json:"subtitle"`
	CoverImg    string   `json:"coverimg"` //if not have this field shoud be ""
	Images      []*Image `json:"images"`
	ReplyNumber int      `json:"replynumber"`
	Replys      []Reply  `json:"replys"`
	ViewType    int      `json:"viewtype"`
	Version     string   `json:"version"`
}
