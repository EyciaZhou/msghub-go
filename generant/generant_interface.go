package generant

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("netease_news")

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
	ID          string   `json:"id"`
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
	From        string   `json:"from"`
	Type        string   `json:"type"`
}

func (m *Message) InsertIntoSQL(stmt *sql.Stmt) (sql.Result, error) {
	return stmt.Exec(m.SnapTime, m.PubTime, m.Source, m.Body, m.Title, m.Subtitle, m.CoverImg, m.ViewType, m.From, m.Type)
}

var (
	db *sql.DB
)

func GetStmtInsert() (*sql.Stmt, error) {
	return db.Prepare(
		`INSERT INTO
				msg (SnapTime, PubTime, SourceURL, Body, Title, SubTitle, CoverImg, ViewType, Frm, Typ)
			VALUES
				(?,?,?,?,?,?,?,?,?,?)
			ON DUPLICATE KEY UPDATE
				SnapTime = VALUES(SnapTime),
				PubTime = VALUES(PubTime),
				Body = VALUES(Body),
				Title = VALUES(Title),
				SubTitle = VALUES(SubTitle),
				CoverImg = VALUES(CoverImg),
				ViewType = VALUES(ViewType)`)
}

func CornCatch() {

}

func init() {
	logging.SetFormatter(logging.MustStringFormatter(
		"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
	))
	var err error
	db, err = sql.Open("mysql", "root:123456@tcp(q.dianm.in:3306)/msghub")
	if err != nil {
		log.Error(err.Error())
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	err = db.Ping()
	if err != nil {
		log.Error(err.Error())
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	log.Info("connected")
}
