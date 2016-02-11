package generant

import (
	"database/sql"
	"errors"
	"fmt"
	"git.eycia.me/eycia/configparser"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zbindenren/logrus_mail"
	"os"
)

type Plugin func(args ...interface{}) (Generant, error)

type Generant interface {
	Catch()

	//Stop: maybe could continue this round catch, and insert into mysql
	Stop()

	//ForceStop: stop immediately, drop this round, only do some clean etc
	ForceStop()
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
	//ID          string   `json:"id"`
	SnapTime    int64    `json:"snaptime"` //*
	PubTime     int64    `json:"pubtime"`  //*
	Source      string   `json:"source"`   //*
	Body        string   `json:"body"`     //*
	Title       string   `json:"title"`    //*
	Subtitle    string   `json:"subtitle"` //*
	CoverImg    string   `json:"coverimg"` //if not have this field shoud be "" //*
	Images      []*Image `json:"images"`
	ReplyNumber int      `json:"replynumber"`
	Replys      []Reply  `json:"replys"`
	ViewType    int      `json:"viewtype"` //*
	Version     string   `json:"version"`
	From        string   `json:"from"` //*
	Channel     string   `json:"tag"`  //*
	Priority    int      `json:"priority"`
}

func (m *Message) InsertIntoSQL() (sql.Result, error) {
	//insert cover img
	var coverImgId interface{}

	if m.CoverImg != "" {
		var err error
		coverImgId, err = insertImgUrlIntoQueue(m.CoverImg)
		if err != nil {
			return nil, err
		}
	} else {
		coverImgId = nil
	}

	res, err := StmtInsert.Exec(m.SnapTime, m.PubTime, m.Source, m.Body, m.Title, m.Subtitle, coverImgId, m.ViewType, m.From, m.Channel)

	if err != nil {
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
		_, err := img.InsertIntoSQL(id)
		if err != nil {
			log.Errorf("Error raised when catch img of SOURCE[%s], URL is [%s] \n ERROR:[%s]", m.Source, img.URL, err.Error())
		}
	}

	return res, nil
}

func isUrl(ur string) bool {
	//TODO: judge the domain
	return true
}

func insertImgUrlIntoQueue(url string) (int64, error) {
	if isUrl(url) {
		res, err := StmtInsertImgToQueue.Exec(url, url)

		if err != nil {
			log.WithField("url", url).Error("Error when exec InsertImgToQueue's STMT, REASON : " + err.Error())
			return 0, err
		}

		rc, _ := res.RowsAffected()
		if rc != 1 {
			//duplicate
			var pid int64
			row := StmtSelectPidFromURL.QueryRow(url)
			err = row.Scan(&pid)
			if err != nil {
				return 0, err
			}
			return pid, nil
		}
		//inserted
		return res.LastInsertId()
	}

	return 0, ErrorNotInvaildURL
}

func (img *Image) InsertIntoQueue() (int64, error) {
	return insertImgUrlIntoQueue(img.URL)
}

func (img *Image) InsertIntoSQL(mid int64) (sql.Result, error) {
	pid, err := img.InsertIntoQueue()
	if err != nil {
		return nil, err
	}
	return StmtInsertRef.Exec(img.Ref, img.Desc, pid, mid)
}

var (
	db                   *sql.DB
	StmtInsert           *sql.Stmt
	StmtInsertRef        *sql.Stmt
	StmtInsertImgToQueue *sql.Stmt
	StmtSelectPidFromURL *sql.Stmt
	StmtSelectMidFromURL *sql.Stmt

	ErrorNotInvaildURL = errors.New("url is not invaild")
)

type Config struct {
	MailApplicationName string `default:"Generant_Interface"`
	MailSMTPAddress     string `default:"127.0.0.1"`
	MailSMTPPort        int    `default:"25"`
	MailFrom            string `default:"root@eycia.me"`
	MailTo              string `default:"zhou.eycia@gmail.com"`

	MailUsername string `default:"nomailusername"`
	MailPassword string `default:"nomailpassword"`

	QueueTableName  string `default:"pic_task_queue"`
	PicRefTableName string `default:"picref"`
	MsgTableName    string `default:"msg"`

	DBAddress  string `default:"127.0.0.1"`
	DBPort     string `default:"3306"`
	DBName     string `default:"msghub"`
	DBUsername string `default:"root"`
	DBPassword string `default:"fmttm233"`
}

var (
	config Config
)

func loadConfig() {
	var err error

	//load
	if len(os.Args) > 1 {
		log.Info("Loading config from ", os.Args[1])
		err = configparser.LoadConfPath(&config, os.Args[1])
	}

	if err != nil || len(os.Args) == 1 {
		log.Info("Loading config use default values")
		err = configparser.LoadConfDefault(&config)
	}

	if err != nil {
		panic(err)
	}
}

func init() {
	loadConfig()

	//process log's mail sending
	/*
		mailhook, err := logrus_mail.NewMailHook(config.MailApplicationName, config.MailSMTPAddress, config.MailSMTPPort, config.MailFrom, config.MailTo)
		if err == nil {
			log.AddHook(mailhook)
		} else {
			log.Error("Can't Hook mail, ERROR:", err.Error())
		}
	*/
	mailhook_auth, err := logrus_mail.NewMailAuthHook(config.MailApplicationName, config.MailSMTPAddress, config.MailSMTPPort, config.MailFrom, config.MailTo,
		config.MailUsername, config.MailPassword)

	if err == nil {
		log.AddHook(mailhook_auth)
		log.Error("Don't Worry, just for send a email to test")
	} else {
		log.Error("Can't Hook mail, ERROR:", err.Error())
	}

	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DBUsername, config.DBPassword, config.DBAddress, config.DBPort, config.DBName)
	db, err = sql.Open("mysql", url)
	if err != nil {
		log.Error("Can't Connect DB REASON : " + err.Error())
		return
	}
	err = db.Ping()
	if err != nil {
		log.Error("Can't Connect DB REASON : " + err.Error())
		return
	}
	log.Info("connected")

	StmtInsert, err = db.Prepare(fmt.Sprintf(
		`INSERT INTO
				%s (SnapTime, PubTime, SourceURL, Body, Title, SubTitle, CoverImg, ViewType, Frm, Tag)
			VALUES
				(?,?,?,?,?,?,?,?,?,?)
			ON DUPLICATE KEY UPDATE
				SnapTime = VALUES(SnapTime),
				PubTime = VALUES(PubTime),
				Body = VALUES(Body),
				Title = VALUES(Title),
				SubTitle = VALUES(SubTitle),
				CoverImg = VALUES(CoverImg),
				ViewType = VALUES(ViewType)`, config.MsgTableName))
	if err != nil {
		log.Error("Error when get STMT, ERROR:[%s]", err.Error())
		return
	}

	StmtInsertRef, err = db.Prepare(fmt.Sprintf(
		`INSERT INTO
				%s (Ref, Description, pid, mid)
			VALUES
				(?,?,?,?)
			ON DUPLICATE KEY UPDATE
				Ref = VALUES(Ref),
				Description = VALUES(Description)`,
		config.PicRefTableName))
	if err != nil {
		panic(err)
		return
	}

	StmtInsertImgToQueue, err = db.Prepare(fmt.Sprintf(`
	INSERT INTO
			%s (url, status, owner)
		SELECT
				?, 0, 0
			FROM DUAL
			WHERE NOT EXISTS (SELECT 1 FROM %s WHERE url=?);
	`, config.QueueTableName, config.QueueTableName))
	/*
		StmtInsertImgToQueue, err = db.Prepare(fmt.Sprintf(`
		INSERT INTO
				%s (url, status)
			VALUES
				(?,?);`, config.QueueTableName))*/
	if err != nil {
		panic(err)
		return
	}

	StmtSelectPidFromURL, err = db.Prepare(fmt.Sprintf(`
	SELECT id FROM %s
		WHERE url=?;
	`, config.QueueTableName))
	if err != nil {
		panic(err)
		return
	}

	StmtSelectMidFromURL, err = db.Prepare(fmt.Sprintf(`
	SELECT id FROM %s
		WHERE SourceURL=?;
	`, config.MsgTableName))
	if err != nil {
		panic(err)
		return
	}
}
