package generant

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/EyciaZhou/configparser"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zbindenren/logrus_mail"
	"io/ioutil"
	"path"
	"sync"
	"time"
)

type LoadConf func(raw []byte) (Generant, error)

type Generant interface {
	Catch()

	//Stop: maybe could continue this round catch, and insert into mysql
	Stop()

	//ForceStop: stop immediately, drop this round, only do some clean etc
	ForceStop()
}

var (
	pluginsMu    sync.Mutex
	regedPlugins = make(map[string]LoadConf)

	generants []Generant
)

func Register(name string, fLoadConf LoadConf) {
	pluginsMu.Lock()
	defer pluginsMu.Unlock()

	regedPlugins[name] = fLoadConf
}

const (
	VIEW_TYPE_NORMAL = 1
	VIEW_TYPE_PICTURES = 2
)

type Author struct {
	Name         string `json:"name"`
	Uid          string `json:"uid"`
	CoverSource string `json:"covert_source"`
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
	Ref   string `json:"ref"` //not set if not have
	Desc  string `json:"desc"`
	Pixes string `json:"pixes"`
	URL   string `json:"url"`
}

type Message struct {
	//ID          string   `json:"id"`
	SnapTime    int64    `json:"snaptime"` //*   //lastmodify
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
	Topic       string   `json:"topic"`
	Version     string   `json:"version"`
	Author      Author   `json:"author"` //*	//TODO to author
	Priority    int      `json:"priority"`
}

type Topic struct {
	Id         string     `json:id`
	Title      string     `json:"title"`
	Msgs       []*Message `json:"messages"`
	LastModify int64      `json:"lastmodify"`
}

func (t *Author) InsertIntoSQL() error {
	if t == nil {
		return errors.New("null author")
	}

	_, err := db.Exec(`
	INSERT INTO
				%s (id, coverImg, name)
			VALUES
				(?,?,?)
			ON DUPLICATE KEY UPDATE
				coverImg = VALUES(coverImg),
				name = VALUES(name)
	`, t.Uid, t.CoverSource, t.Name)

	if err != nil {
		return err
	}

	return nil
}

func (t *Topic) InsertIntoSQL() error {
	_, err := StmtTopicInsert.Exec(t.Id, t.Title, t.LastModify)
	if err != nil {
		return err
	}

	for _, msg := range t.Msgs {
		_, _ = msg.InsertIntoSQL()
	}

	return nil
}

func (m *Message) InsertIntoSQL() (sql.Result, error) {
	//insert cover img
	var coverImgId sql.NullInt64

	if m.CoverImg != "" {
		var err error
		coverImgId.Int64, err = insertImgUrlIntoQueue(m.CoverImg)
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

	res, err := StmtInsert.Exec(m.SnapTime, m.PubTime, m.Source, m.Body, m.Title, m.Subtitle, coverImgId, m.ViewType, m.Author.Uid, TopicId)

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
		_, _ = img.InsertIntoSQL(id)
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
	return StmtInsertRef.Exec(img.Ref, img.Desc, img.Pixes, pid, mid)
}

var (
	db                   *sql.DB
	StmtInsert           *sql.Stmt
	StmtInsertRef        *sql.Stmt
	StmtInsertImgToQueue *sql.Stmt
	StmtSelectPidFromURL *sql.Stmt
	StmtSelectMidFromURL *sql.Stmt
	StmtTopicInsert      *sql.Stmt

	ErrorNotInvaildURL = errors.New("url is not invaild")
)

type Config struct {
	MailEnabled         bool   `default:"false"`
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
	TopicTableName  string `default:"topic"`

	DBAddress  string `default:"127.0.0.1"`
	DBPort     string `default:"3306"`
	DBName     string `default:"msghub"`
	DBUsername string `default:"root"`
	DBPassword string `default:"fmttm233"`

	ConfDir         string
	ConfFileNames   []string
	ConfPluginNames []string
}

var (
	config Config
)

func loadConfig() error {
	pluginsMu.Lock()
	defer pluginsMu.Unlock()

	var err error

	//load self first
	/*
		if len(os.Args) > 1 {
			log.Info("Loading config from ", os.Args[1])
			err = configparser.LoadConfPath(&config, os.Args[1])
		}

		if err != nil || len(os.Args) == 1 {
			log.Info("Loading config use default values")
			err = configparser.LoadConfDefault(&config)
		}
	*/

	configparser.AutoLoadConfig("generant", &config)
	//configparser.LoadConfFromFlag(&config)
	configparser.ToJson(&config)

	return err
}

func loadPluginConfig() error {
	//
	if generants != nil && len(generants) > 0 {
		return errors.New("can't load config twice")
	}

	if len(config.ConfFileNames) != len(config.ConfPluginNames) {
		return errors.New("the length of ConfFileNames and ConfPluginNames not same")
	}

	generants = make([]Generant, len(config.ConfPluginNames))

	config.ConfDir = path.Clean(config.ConfDir)

	log.Infof("%d plugin configs to load", len(config.ConfFileNames))

	for i, fn := range config.ConfFileNames {
		log.Infof("[%d/%d]...", i+1, len(config.ConfFileNames))
		bs, err := ioutil.ReadFile(config.ConfDir + "/" + fn)
		if err != nil {
			return err
		}

		if plugin, hv := regedPlugins[config.ConfPluginNames[i]]; !hv {
			return errors.New("Doesn't registed Plugin : " + config.ConfPluginNames[i])
		} else {
			generants[i], err = plugin(bs)
			if err != nil {
				return fmt.Errorf("Error when load %d th plugin : %s", i+1, err.Error())
			}
		}
	}

	return nil
}

func Init() {
	log.Info("Start Load Config")
	err := loadConfig()
	if err != nil {
		panic(err)
	}

	//process log's mail sending
	/*
		mailhook, err := logrus_mail.NewMailHook(config.MailApplicationName, config.MailSMTPAddress, config.MailSMTPPort, config.MailFrom, config.MailTo)
		if err == nil {
			log.AddHook(mailhook)
		} else {
			log.Error("Can't Hook mail, ERROR:", err.Error())
		}
	*/

	if config.MailEnabled {
		log.Info("Start Bind Mail Hook")
		mailhook_auth, err := logrus_mail.NewMailAuthHook(config.MailApplicationName, config.MailSMTPAddress, config.MailSMTPPort, config.MailFrom, config.MailTo,
			config.MailUsername, config.MailPassword)

		if err == nil {
			log.AddHook(mailhook_auth)
			log.Error("Don't Worry, just for send a email to test")
		} else {
			log.Panic("Can't Hook mail, ERROR:", err.Error())
		}
	}

	log.Info("Start Connect mysql")
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DBUsername, config.DBPassword, config.DBAddress, config.DBPort, config.DBName)
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

	log.Info("Start prepare stmt")
	StmtInsert, err = db.Prepare(fmt.Sprintf(
		`INSERT INTO
				%s (SnapTime, PubTime, SourceURL, Body, Title, SubTitle, CoverImg, ViewType, Frm, Topic)
			VALUES
				(?,?,?,?,?,?,?,?,?,?)
			ON DUPLICATE KEY UPDATE
				SnapTime = VALUES(SnapTime),
				PubTime = VALUES(PubTime),
				Body = VALUES(Body),
				Title = VALUES(Title),
				SubTitle = VALUES(SubTitle),
				CoverImg = VALUES(CoverImg),
				ViewType = VALUES(ViewType),
				Topic = VALUES(Topic)`, config.MsgTableName))
	if err != nil {
		log.Panic(err.Error())
		return
	}

	StmtTopicInsert, err = db.Prepare(fmt.Sprintf(
		`INSERT INTO
				%s (id, Title, LastModify)
			VALUES
				(?,?,?)
			ON DUPLICATE KEY UPDATE
				Title = VALUES(Title),
				LastModify = VALUES(LastModify)`, config.TopicTableName))

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
		log.Panic(err.Error())
		return
	}

	StmtSelectPidFromURL, err = db.Prepare(fmt.Sprintf(`
	SELECT id FROM %s
		WHERE url=?;
	`, config.QueueTableName))
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

	log.Info("Start load plugins's config")
	err = loadPluginConfig()
	if err != nil {
		log.Panic(err.Error())
	}

	log.Infof("Start fire plugins, %d plugins to fire", len(generants))
	for i, gen := range generants {
		log.Infof("[%d/%d]...", i+1, len(generants))
		go gen.Catch()
		log.Info("fired and start delay")
		time.Sleep(10 * time.Second)
	}

	log.Info("Init finished")
}
