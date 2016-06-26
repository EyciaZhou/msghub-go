package Interface

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/EyciaZhou/configparser"
	"github.com/EyciaZhou/picRouter/PicPipe"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

const (
	//VIEW_TYPE_NORMAL:
	//	some thing like a news, and some pictures inset the content
	VIEW_TYPE_NORMAL = 1

	//VIEW_TYPE_PICTURES:
	//	some thing like tweet,
	//	have series of picture. but short, can avoiding body content.
	VIEW_TYPE_PICTURES = 2
)

//Author:
//	Name: name to display
//	Uid: uuid of this author, recommend using PLUGIN-NAME{_TOPIC-NAME}_AUTHOR-NAME
//	AvatarUrl: url of Author's avatar
type Author struct {
	Name      string `json:"name"`
	Uid       string `json:"uid"`
	AvatarUrl string `json:"covert_source"` //can empty
}

//InsertIntoSQL:
//	Insert this author to sql database,
//	throw error if the author is null, or some error sql caused
func (t *Author) InsertIntoSQL() error {
	if t == nil {
		return errors.New("null author")
	}

	pid := sql.NullString{}

	if t.AvatarUrl != "" {
		_pid, err := insertImgUrlIntoQueue(t.AvatarUrl)
		if err != nil {
			return err
		}
		pid.String = _pid
		pid.Valid = true
	} else {
		pid.Valid = false
	}

	_, err := db.Exec(`
	INSERT INTO
				author (id, coverImg, name)
			VALUES
				(?,?,?)
			ON DUPLICATE KEY UPDATE
				coverImg = VALUES(coverImg),
				name = VALUES(name)
	`, t.Uid, pid, t.Name)
	return err
}

//ReplyFloor:
//	not using now
type ReplyFloor struct {
	Floor   int    `json:"floor"`
	Time    int64  `json:"time"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Digg    int    `jsong:"digg"`
}

//Reply:
//	not using now
type Reply []*ReplyFloor

//Image:
//
type Image struct {
	Ref   string `json:"ref"` //not set if not have
	Desc  string `json:"desc"`
	Pixes string `json:"pixes"`
	URL   string `json:"url"`
}

func (img *Image) InsertIntoQueue() (string, error) {
	return insertImgUrlIntoQueue(img.URL)
}

func (img *Image) InsertIntoSQL(mid int64) (sql.Result, error) {
	pid, err := img.InsertIntoQueue()
	if err != nil {
		return nil, err
	}
	return StmtInsertRef.Exec(img.Ref, img.Desc, img.Pixes, pid, mid)
}

type Message struct {
	//ID          string   `json:"id"`
	SnapTime int64    `json:"snaptime"` //*   //lastmodify
	PubTime  int64    `json:"pubtime"`  //*
	Source   string   `json:"source"`   //*
	Body     string   `json:"body"`     //*
	Title    string   `json:"title"`    //*
	Subtitle string   `json:"subtitle"` //*
	CoverImg string   `json:"coverimg"` //if not have this field shoud be "" //*
	Images   []*Image `json:"images"`
	//ReplyNumber int64    `json:"replynumber"`
	//Replys      []Reply  `json:"replys"`
	ViewType int    `json:"viewtype"` //*
	Topic    string `json:"topic"`    //*
	//Version     string   `json:"version"`
	Tag    string  `json:"tag"`    //*
	Author *Author `json:"author"` //*
	//Priority    int      `json:"priority"`
}

func (m *Message) InsertIntoSQL() (sql.Result, error) {
	//insert cover img
	var coverImgId sql.NullString

	err := m.Author.InsertIntoSQL()
	if err != nil {
		return nil, err
	}

	if m.CoverImg != "" {
		var err error
		coverImgId.String, err = insertImgUrlIntoQueue(m.CoverImg)
		if err != nil {
			return nil, err
		}
		coverImgId.Valid = true
	} else {
		coverImgId.Valid = false
	}

	var TopicId sql.NullString
	if m.Topic != "" {
		TopicId.Valid, TopicId.String = true, m.Topic
	}

	res, err := StmtInsert.Exec(m.SnapTime, m.PubTime, m.Source, m.Body, m.Title, m.Subtitle, coverImgId, m.ViewType, m.Author.Uid, m.Tag, TopicId)

	if err != nil {
		log.Errorf("Error when insert Message[%v]\n error:[%s]", *m, err.Error())
		return nil, err
	}

	id, _ := res.LastInsertId()

	//if not modified, call SELECT to find the id
	if num, _ := res.RowsAffected(); num == 0 {
		row := StmtSelectMidFromURL.QueryRow(m.Source)
		err = row.Scan(&id)
		if err != nil {
			return nil, err
		}
	}

	for _, img := range m.Images {
		_, err = img.InsertIntoSQL(id)
		if err != nil {
			log.Warn(err.Error())
		}
	}

	return res, nil
}

func MsgsInsertIntoSQL(msgs []*Message) error {
	for _, msg := range msgs {
		_, err := msg.InsertIntoSQL()
		if err != nil {
			log.Warn(err.Error())
		}
	}
	return nil
}

var (
	picQueue pic.PicTaskPipe
)

func insertImgUrlIntoQueue(url string) (string, error) {
	task, err := picQueue.UpsertTask(url)
	if err != nil {
		return "", err
	}
	return task.Key, nil
}

var (
	db                   *sql.DB
	StmtInsert           *sql.Stmt
	StmtInsertRef        *sql.Stmt
	StmtSelectMidFromURL *sql.Stmt
	//	StmtTopicInsert      *sql.Stmt

	ErrorNotInvaildURL = errors.New("url is not invaild")
)

type config_t struct {
	PicRefTableName string `default:"picref"`
	MsgTableName    string `default:"msg"`
	TopicTableName  string `default:"topic"`

	DBAddress  string `default:"127.0.0.1"`
	DBPort     string `default:"3306"`
	DBName     string `default:"msghub"`
	DBUsername string `default:"root"`
	DBPassword string `default:"fmttm233"`
}

var (
	config config_t
)

func loadConfig() {
	configparser.AutoLoadConfig("interface", &config)
	configparser.ToJson(&config)
}

func Init() {
	var err error

	log.Info("Start Load Config")
	loadConfig()

	log.Info("Start Connect mysql")
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?collation=utf8mb4_general_ci", config.DBUsername, config.DBPassword, config.DBAddress, config.DBPort, config.DBName)
	db, err = sql.Open("mysql", url)
	if err != nil {
		log.Panic("Can't Connect DB REASON : " + err.Error())
		return
	}
	err = db.Ping()
	if err != nil {
		log.Panic("Can't Connect DB REASON : " + err.Error())
		return
	}
	log.Info("connected")

	picQueue = pic.NewMySQLPicPipeUseConnectedDB(db)
	if err != nil {
		log.Panic("Can't Connect DB REASON : " + err.Error())
		return
	}

	log.Info("Start prepare stmt")
	StmtInsert, err = db.Prepare(fmt.Sprintf(
		`INSERT INTO
				%s (SnapTime, PubTime, SourceURL, Body, Title, SubTitle, CoverImg, ViewType, AuthorId, Tag, Topic)
			VALUES
				(?,?,?,?,?,?,?,?,?,?,?)
			ON DUPLICATE KEY UPDATE
				SnapTime = VALUES(SnapTime),
				PubTime = VALUES(PubTime),
				Body = VALUES(Body),
				Title = VALUES(Title),
				SubTitle = VALUES(SubTitle),
				CoverImg = VALUES(CoverImg),
				ViewType = VALUES(ViewType),
				AuthorId = VALUES(AuthorId),
				Tag = VALUES(Tag),
				Topic = VALUES(Topic)`, config.MsgTableName))
	if err != nil {
		log.Panic(err.Error())
		return
	}

	/*StmtTopicInsert, err = db.Prepare(fmt.Sprintf(
	`INSERT INTO
			%s (id, Title, LastModify)
		VALUES
			(?,?,?)
		ON DUPLICATE KEY UPDATE
			Title = VALUES(Title),
			LastModify = VALUES(LastModify)`, config.TopicTableName))*/

	StmtInsertRef, err = db.Prepare(fmt.Sprintf(
		`INSERT INTO
				%s (Ref, Description, Pixes, pid, mid)
			VALUES
				(?,?,?,?,?)
			ON DUPLICATE KEY UPDATE
				Ref = VALUES(Ref),
				Pixes = VALUES(Pixes),
				Description = VALUES(Description)`,
		config.PicRefTableName))
	if err != nil {
		log.Panic(err.Error())
		return
	}

	StmtSelectMidFromURL, err = db.Prepare(fmt.Sprintf(`
	SELECT id FROM %s
		WHERE SourceURL=?;
	`, config.MsgTableName))
	if err != nil {
		log.Panic(err.Error())
		return
	}
}
