package nenews_api

import (
//"encoding/json"
	"errors"
	"fmt"
	"github.com/EyciaZhou/msghub.go/Utiles"
	"github.com/EyciaZhou/msghub.go/generant/netease_news/types"
	"github.com/EyciaZhou/msghub.go/generant/utils"
	log "github.com/Sirupsen/logrus"
	"strings"
	"github.com/EyciaZhou/msghub.go/generant"
	"time"
)

func listURL(list string, page int) string {
	return fmt.Sprintf("http://c.3g.163.com/nc/article/list/%s/%d-%d.html", list, page*20, (page+1)*20)
}

func contentURL(id string) string {
	return fmt.Sprintf("http://c.m.163.com/nc/article/%s/full.html", id)
}

func replyURL(id string, boardId string) string {
	return fmt.Sprintf("http://comment.api.163.com/api/jsonp/post/list/hot/%s/%s/0/20/20/0/0", boardId, id)
}

func photoContentURL(splited1 string, splited2 string) string {
	return fmt.Sprintf("http://c.m.163.com/photo/api/set/%s/%s.json", splited1, splited2)
}

func specialURL(id string) string {
	return fmt.Sprintf("http://c.m.163.com/nc/special/%s.html", id)
}

/*
BaseController:
	basic api
 */
type BaseController struct {
}

/*
apiGetNewsReply:
	get news reply by (id, boardId)
 */
func (p BaseController) apiGetNewsReply(id string, boardId string) ([]nenews_types.Reply, error) {
	//pass reply
	return []nenews_types.Reply{}, nil

	var c nenews_types.Reply_tmp
	err := generant_utils.RequestToJson("GET", replyURL(id, boardId), nil, &c)
	if err != nil {
		return nil, err
	}

	return c.HotPosts, nil
}

/*
apiGetNormalNews:
	get full information about a normal news by its basic information
	basic information in key-value format
 */
func (p BaseController) apiGetNormalNews(item map[string]interface{}) (r *nenews_types.NormalNews, er error) {
	defer func() {
		if err := recover(); err != nil {
			r = nil
			er = Utiles.NewPanicError(err.(error))
		}
	}()

	//get id
	id := item["docid"].(string) //panic^2

	//get content
	var c map[string]*nenews_types.NormalNews
	err := generant_utils.RequestToJson("GET", contentURL(id), nil, &c)
	if err != nil {
		return nil, err
	}

	content := c[id]//panic

	//get comment
	reply, e := p.apiGetNewsReply(id, content.BoardId)
	if e != nil {
		return nil, errors.New(fmt.Sprintf("get news reply error, ERROR:[%v]", e.Error()))
	}
	content.Replys = reply

	can := true
	//get ext info
	if content.CoverURL, can = item["imgsrc"].(string); !can {
		content.CoverURL = ""
	}

	content.URL = item["url"].(string)          //panic
	content.SnapTime = item["lmodify"].(string) //panic
	return content, nil
}

/*
apiGetPhotosetNews:
	get full information about a photo set news by its basic information
	basic information in key-value format
 */
func (p BaseController) apiGetPhotosetNews(item map[string]interface{}) (r *nenews_types.PhotoSet, er error) {
	defer func() {
		if err := recover(); err != nil {
			r = nil
			er = Utiles.NewPanicError(err.(error))
		}
	}()

	//progress id
	photosetId := item["photosetID"].(string)

	pid_parted := strings.Split(photosetId, "|")
	if len(pid_parted) != 2 {
		return nil, errors.New(fmt.Sprintf("photosetID format unparseabe PHOTOTSETID:[%v]", photosetId))
	}

	//get info
	var content nenews_types.PhotoSet
	err := generant_utils.RequestToJson("GET", photoContentURL(pid_parted[0][4:], pid_parted[1]), nil, &content)
	if err != nil {
		return nil, err
	}

	//get comment
	reply, err := p.apiGetNewsReply(content.ID, content.BoardId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("get news reply error, ERROR:[%s]", err.Error()))
	}
	content.Replys = reply

	//add extra info
	content.Body = content.SubTitle
	content.ReplyCount = (int)(item["replyCount"].(float64))
	content.SnapTime = item["lmodify"].(string)

	return &content, nil
}

func (p BaseController) transArrayOfInterfaceToArrayOfMap(interfaces []interface{}) []map[string]interface{} {
	maps := make([]map[string]interface{}, len(interfaces))
	cnt := 0
	for i, _ := range interfaces {
		if mp, ok := interfaces[i].(map[string]interface{}); ok {
			maps[cnt] = mp
			cnt++
		}
	}
	return maps[:cnt]
}

/*
parseNewsList:
	parse list of news from map-interface format to
	struct format because of body info are not
	included in basic info, so there will a additional
	network request for each news.
 */
func (p BaseController) parseNewsList(baseInfoList []map[string]interface{}) ([]nenews_types.News, error) {
	var (
		err     error
		result  []nenews_types.News
		content nenews_types.News
	)

	for i, item := range baseInfoList {
		log.Infof("fetching [%d/%d]", i+1, len(baseInfoList))
		//judge skip type
		if typ, hv := item["skipType"]; hv {
			if skipTyp, can := typ.(string); !can {
				err = fmt.Errorf("have skiptype but type is not string , SKIPID[%T : %v]", typ, typ)
			} else {
				switch skipTyp {
				case "photoset":
					content, err = p.apiGetPhotosetNews(item)
				case "special":
					//TODO
					continue
				default:
					//TODO
					continue
				}
			}
		} else {
			content, err = p.apiGetNormalNews(item)
		}
		if err != nil {
			log.WithFields(log.Fields{
				"time":  "getNewsList",
				"error": err.Error(),
				"item":  item,
			}).Warn("error when change map to struct")
			continue
		}
		result = append(result, content)
	}
	log.Infof("[%d fetched/%d expect]", len(result), len(baseInfoList))
	return result, nil
}

type ChannController struct {
	basic BaseController

	channName string
	listId    string

	delayTime time.Duration
}

func NewChannController(channName string, listId string, delayTime time.Duration) *ChannController {
	return &ChannController{
		BaseController{},
		channName,
		listId,
		delayTime,
	}
}

func (p *ChannController) apiGetNewsChannel(page int) (r *nenews_types.Topic, er error) {
	defer func() {
		if err := recover(); err != nil {
			r = nil
			er = Utiles.NewPanicError(err.(error))
		}
	}()

	var v map[string]([]map[string]interface{})
	err := generant_utils.RequestToJson("GET", listURL(p.listId, page), nil, &v)
	if err != nil {
		return nil, err
	}

	newss, err := p.basic.parseNewsList(v[p.listId])	//panic
	if err != nil {
		return nil, err
	}

	return &nenews_types.Topic{
		Newss: newss,
		Id:    "nenews_list_" + p.listId,
		Title: p.channName,
	}, nil
}

func (p *ChannController) GetNew() (generant.CanConvertToTopic, error) {
	return p.apiGetNewsChannel(0)
}

func (p *ChannController) DelayBetweenCatchRound() time.Duration {
	return p.delayTime
}

type TopicController struct {
	basic BaseController

	specialId string

	delayTime time.Duration
}

func NewTopicController(specialId string, delayTime time.Duration) *TopicController {
	return &TopicController{
		BaseController{},
		specialId,

		delayTime,
	}
}

func (p *TopicController) apiGetSpecialList() (r *nenews_types.Topic, er error) {
	defer func() {
		if err := recover(); err != nil {
			r = nil
			er = Utiles.NewPanicError(err.(error))
		}
	}()

	var v map[string](map[string]interface{})
	err := generant_utils.RequestToJson("GET", specialURL(p.specialId), nil, &v)
	if err != nil {
		return nil, err
	}

	infos := v[p.specialId]                                                        //panic
	topics := p.basic.transArrayOfInterfaceToArrayOfMap(infos["topics"].([]interface{})) //panic^3

	newss := []nenews_types.News{}

	//a special including some topics
	for _, t := range topics {
		msgs, err := p.basic.parseNewsList(p.basic.transArrayOfInterfaceToArrayOfMap(t["docs"].([]interface{}))) //panic
		if err != nil {
			return nil, err
		}
		newss = append(newss, msgs...)
	}

	return &nenews_types.Topic{
		Id:    "nenews_specia_" + p.specialId,
		Title: infos["sname"].(string), //panic
		Newss: []nenews_types.News{},
	}, nil
}

func (p *TopicController) GetNew() (generant.CanConvertToTopic, error) {
	return p.apiGetSpecialList()
}

func (p *TopicController) DelayBetweenCatchRound() time.Duration {
	return p.delayTime
}