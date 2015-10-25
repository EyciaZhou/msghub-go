package netease_news

import (
	"encoding/json"
	"errors"
	"fmt"
	"git.eycia.me/eycia/msghub/generant"
	"git.eycia.me/eycia/msghub/netTools"
	"github.com/op/go-logging"
	"strings"
	"time"
)

var log = logging.MustGetLogger("netease_news")

type NeteaseNewsChannel struct {
}

/*
{
	"T1295501906343":[
		{
			"hasCover":false,
			"hasHead":1,
			"replyCount":4640,
			"hasImg":1,
			"digest":"男子要求4S店退熄火路虎未收回复，开3辆豪车堵店门口。",
			"hasIcon":false,
			"docid":"B29S9I3J00963VRO",
			"title":"买路虎5天熄火 男子开蝙蝠车堵4S店",
			"order":1,
			"priority":250,
			"lmodify":"2015-08-30 21:07:15",
			"boardid":"3g_bbs",
			"url_3w":"http://help.3g.163.com/15/0830/20/B29S9I3J00963VRO.html",
			"template":"manual",
			"votecount":3425,
			"alias":"",
			"cid":"",
			"url":"http://3g.163.com/ntes/15/0830/20/B29S9I3J00963VRO.html",
			"hasAD":1,
			"source":"钱江晚报",
			"subtitle":"",
			"imgsrc":"Uhttp://img2.cache.netease.com/3g/2015/8/30/20150830210154c4c41.jpg",
			"tname":"",
			"ename":"",
			"ptime":"2015-08-30 20:04:29"
		},
		{
			"url_3w":"http://help.3g.163.com/15/0830/18/B29LQEP400963VRO.html",
			"votecount":632,
			"replyCount":1112,
			"digest":"称王玉发当师政委时，谷俊山就是个不入流的营职干部。",
			"url":"http://3g.163.com/ntes/15/0830/18/B29LQEP400963VRO.html",
			"docid":"B29LQEP400963VRO",
			"title":"知情人否认王玉发涉谷俊山案",
			"source":"长安街知事",
			"priority":135,
			"lmodify":"2015-08-30 18:11:49",
			"imgsrc":"http://img4.cache.netease.com/3g/2015/8/30/201508301813396ec15.jpg",
			"subtitle":"",
			"boardid":"3g_bbs",
			"ptime":"2015-08-30 18:11:22"
			},
		}
	]

}
*/

func listURL(list string, page int) string {
	return fmt.Sprintf("http://c.3g.163.com/nc/article/list/%s/%d-%d.html", list, page*20, (page+1)*20)
}

func contentURL(id string) string {
	return fmt.Sprintf("http://c.m.163.com/nc/article/%s/full.html", id)
}

func ReplyURL(id string, boardId string) string {
	return fmt.Sprintf("http://comment.api.163.com/api/jsonp/post/list/hot/%s/%s/0/20/20/0/0", boardId, id)
}

func getNewsContent(id string) (*News, error) {
	url := contentURL(id)
	newsContentPain, err := netTools.Get(url)
	if err != nil {
		return nil, err
	}
	//log.Debug("NEWS_CONTENT_PAIN:[%v]", (string)(newsContentPain))
	var c map[string]*News

	err = json.Unmarshal(newsContentPain, &c)
	if err != nil {
		return nil, err
	}
	if _, can := c[id]; !can {
		return nil, errors.New("can't get news content because of  no news under id " + id)
	}
	//n := c[0][id]
	//log.Debug("%v", c[id])
	return c[id], nil
}

func getNewsReply(id string, boardId string) ([]Reply, error) {
	url := ReplyURL(id, boardId)
	newsReplyPain, err := netTools.Get(url)
	if err != nil {
		return nil, err
	}

	//log.Debug("URL:[%v]", url)
	//log.Debug("NEWS_REPLY_PAIN[%v]", (string)(newsReplyPain))
	var c reply_tmp
	err = json.Unmarshal(newsReplyPain, &c)
	if err != nil {
		return nil, err
	}
	return c.HotPosts, nil
}

func getNormalNews(item map[string]interface{}) (*generant.Message, error) {
	//log.Debug("ITEM%d: %v\n", i, item)
	//get id
	id, can := item["docid"].(string)
	if !can {
		return nil, errors.New(fmt.Sprintf("docid is not string or not set, DOCID:[%T : %v]", id, id))
	}

	//get content
	content, e := getNewsContent(id)
	if e != nil {
		return nil, errors.New(fmt.Sprintf("content can not get, ERROR:[%s]", e.Error()))
	}

	//log.Debug("NEWS_CONTENT[%v]", content)

	//get comment
	reply, e := getNewsReply(id, content.BoardId)
	if e != nil {
		return nil, errors.New(fmt.Sprintf("get news reply error, ERROR:[%v]", e.Error()))
	}
	content.Replys = reply

	//get ext info
	flag := true
	content.CoverURL, can = item["imgsrc"].(string)
	flag = flag && can
	content.URL, can = item["url"].(string)
	flag = flag && can
	content.SnapTime, can = item["lmodify"].(string)
	flag = flag && can
	pri, can := item["priority"].(float64)
	flag = flag && can
	content.Priority = int(pri)
	if !flag {
		return nil, errors.New(fmt.Sprintf("can't trans type int, IMGSRC:[%t], URL:[%t], PRIORITY:[%t]\n", item["imgsrc"], item["url"], item["priority"]))
	}
	content.ViewType = 1
	return content.ToMsg()
}

func getPhotosetNews(item map[string]interface{}) (*generant.Message, error) {
	/*
		{
			"docid":"9IG74V5H00963VRO_B2DM8O7QguohaoupdateDoc",
			"title":"十大最有趣的新发现的物种",
			"imgextra":[
				{"imgsrc":"http://img3.cache.netease.com/3g/2015/9/1/20150901073735e82ea.jpg"},
				{"imgsrc":"http://img5.cache.netease.com/3g/2015/9/1/20150901073737a8698.jpg"}
			],
			"replyCount":40,
			"skipID":"0AI20009|7470",
			"priority":75,
			"lmodify":"2015-09-01 07:36:08",
			"imgsrc":"http://img4.cache.netease.com/3g/2015/9/1/2015090107373304d02.jpg",
			"digest":"",
			"skipType":"photoset",
			"photosetID":"0AI20009|7470",
			"ptime":"2015-09-01 07:36:08"
		}
	*/
	//progress id
	photosetId, can := item["photosetID"].(string)

	if !can {
		return nil, errors.New(fmt.Sprintf("photosetId is not string or not set, photosetid:[%T : %v]", item["photosetID"], item["photosetID"]))
	}

	pid_parted := strings.Split(photosetId, "|")
	if len(pid_parted) != 2 {
		return nil, errors.New(fmt.Sprintf("photosetID format unparseabe PHOTOTSETID:[%v]", photosetId))
	}

	//log.Debug("%v", pid_parted)

	//get info
	url := fmt.Sprintf("http://c.m.163.com/photo/api/set/%s/%s.json", pid_parted[0][4:], pid_parted[1])
	contentPain, err := netTools.Get(url)

	//log.Debug(photosetId)
	//log.Debug("http://c.m.163.com/photo/api/set/%s/%s.json", pid_parted[0][4:], pid_parted[1])

	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't get photoset's info, URL:[%s], ERROR:[%s]", url, err.Error()))
	}

	var content PhotoSet
	err = json.Unmarshal(contentPain, &content)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("photoset's info unmarshal error: ERROR[%s]", err.Error()))
	}

	//log.Debug("%v", content)

	//get comment
	reply, err := getNewsReply(content.ID, content.BoardId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("get news reply error, ERROR:[%v]", err.Error()))
	}
	content.Replys = reply

	//add extra info
	content.Body = content.SubTitle

	pri, can := item["priority"].(float64)
	content.Priority = int(pri)
	if !can {
		return nil, errors.New(fmt.Sprintf("can't trans type int, PRIORITY:[%t]\n", item["priority"]))
	}

	rc, can := item["replyCount"].(float64)
	content.ReplyCount = (int)(rc)
	if !can {
		return nil, errors.New(fmt.Sprintf("can't trans type int,  REPLYCOUNT:[%t]\n", item["replyCount"]))
	}

	content.SnapTime, can = item["lmodify"].(string)

	content.ViewType = 2

	return content.ToMsg()
	/*
		http://c.m.163.com/photo/api/set/0009/7469.json
		{
			"postid":"PHOT079D00090AI2",
			"series":"",
			"clientadurl":"http://img4.cache.netease.com/photo/0096/2015-08-24/B1PSJPV46CB40096.png",
			"desc":"公用马桶圈每平方英寸的细菌数量超过1000个。然而，有些电子产品携带了没有“冲掉”的更多细菌。",
			"datatime":"2015-09-01 05:54:05",
			"createdate":"2015-08-31 22:25:30",
			"relatedids":[],
			"scover":"http://img3.cache.netease.com/photo/0009/2015-08-31/s_B2CMGB090AI20009.jpg",
			"autoid":"",
			"url":"http://tech.163.com/photoview/0AI20009/7469.html",
			"creator":"李德雄",
			"reporter":"",
			"photos":[
				{
					"timgurl":"http://img4.cache.netease.com/photo/0009/2015-08-31/t_B2CMGB090AI20009.jpg",
					"photohtml":"http://tech.163.com/photoview/0AI20009/7469.html#p=B2CMGB090AI20009",
					"newsurl":"#",
					"squareimgurl":"http://img3.cache.netease.com/photo/0009/2015-08-31/400x400_B2CMGB090AI20009.jpg",
					"cimgurl":"http://img4.cache.netease.com/photo/0009/2015-08-31/c_B2CMGB090AI20009.jpg",
					"imgtitle":"",
					"simgurl":"http://img3.cache.netease.com/photo/0009/2015-08-31/s_B2CMGB090AI20009.jpg",
					"note":"公用马桶圈每平方英寸的细菌数量超过1000个。然而，有些电子产品携带了没有“冲掉”的更多细菌。《福布斯》杂志网站刊文称，说起细菌，最恶心的莫过于卫生间里的细菌。公用马桶圈每平方英寸的细菌数量超过1000个。然而，有些电子产品携带了没有“冲掉”的更多细菌。做好恶心的准备！",
					"photoid":"B2CMGB090AI20009",
					"imgurl":"http://img4.cache.netease.com/photo/0009/2015-08-31/B2CMGB090AI20009.jpg"
				},
				{
					"timgurl":"http://img3.cache.netease.com/photo/0009/2015-08-31/t_B2CMGB4G0AI20009.jpg",
					"photohtml":"http://tech.163.com/photoview/0AI20009/7469.html#p=B2CMGB4G0AI20009",
					"newsurl":"#",
					"squareimgurl":"http://img4.cache.netease.com/photo/0009/2015-08-31/400x400_B2CMGB4G0AI20009.jpg",
					"cimgurl":"http://img3.cache.netease.com/photo/0009/2015-08-31/c_B2CMGB4G0AI20009.jpg",
					"imgtitle":"",
					"simgurl":"http://img3.cache.netease.com/photo/0009/2015-08-31/s_B2CMGB4G0AI20009.jpg",
					"note":"手机和大肠杆菌你的智能手机上有你的照片、音乐、联系人和大量的粪大肠杆菌。亚利桑那大学进行的研究发现，手机的带菌量是大多数马桶座的10倍。另一项研究发现，手机每平方英寸的细菌数量可能达到惊人的25107个。",
					"photoid":"B2CMGB4G0AI20009",
					"imgurl":"http://img3.cache.netease.com/photo/0009/2015-08-31/B2CMGB4G0AI20009.jpg"
				},
				{
					"timgurl":"http://img3.cache.netease.com/photo/0009/2015-08-31/t_B2CMGB7Q0AI20009.jpg","photohtml":"http://tech.163.com/photoview/0AI20009/7469.html#p=B2CMGB7Q0AI20009","newsurl":"#","squareimgurl":"http://img3.cache.netease.com/photo/0009/2015-08-31/400x400_B2CMGB7Q0AI20009.jpg","cimgurl":"http://img4.cache.netease.com/photo/0009/2015-08-31/c_B2CMGB7Q0AI20009.jpg","imgtitle":"","simgurl":"http://img4.cache.netease.com/photo/0009/2015-08-31/s_B2CMGB7Q0AI20009.jpg","note":"被污染的键盘键盘的带菌量是家用马桶座的3倍，公用马桶座的近3倍。其他研究发现，电脑键盘每平方英寸有3000个细菌，电脑鼠标有1600个。","photoid":"B2CMGB7Q0AI20009","imgurl":"http://img3.cache.netease.com/photo/0009/2015-08-31/B2CMGB7Q0AI20009.jpg"},
				{"
					timgurl":"http://img4.cache.netease.com/photo/0009/2015-08-31/t_B2CMGBBR0AI20009.jpg","photohtml":"http://tech.163.com/photoview/0AI20009/7469.html#p=B2CMGBBR0AI20009","newsurl":"#","squareimgurl":"http://img4.cache.netease.com/photo/0009/2015-08-31/400x400_B2CMGBBR0AI20009.jpg","cimgurl":"http://img4.cache.netease.com/photo/0009/2015-08-31/c_B2CMGBBR0AI20009.jpg","imgtitle":"","simgurl":"http://img3.cache.netease.com/photo/0009/2015-08-31/s_B2CMGBBR0AI20009.jpg","note":"别玩游戏了！带菌量是马桶座近5倍的游戏手柄，或许会遇到更可怕的东西（比如与萨菲罗斯对峙的那个Boss），然而，大肠杆菌是潜在的侵入者之一。","photoid":"B2CMGBBR0AI20009","imgurl":"http://img4.cache.netease.com/photo/0009/2015-08-31/B2CMGBBR0AI20009.jpg"},
				{
					"timgurl":"http://img4.cache.netease.com/photo/0009/2015-08-31/t_B2CMGBDU0AI20009.jpg","photohtml":"http://tech.163.com/photoview/0AI20009/7469.html#p=B2CMGBDU0AI20009","newsurl":"#","squareimgurl":"http://img4.cache.netease.com/photo/0009/2015-08-31/400x400_B2CMGBDU0AI20009.jpg","cimgurl":"http://img4.cache.netease.com/photo/0009/2015-08-31/c_B2CMGBDU0AI20009.jpg","imgtitle":"","simgurl":"http://img3.cache.netease.com/photo/0009/2015-08-31/s_B2CMGBDU0AI20009.jpg","note":"平板电脑、电子阅读器和大肠杆菌，天啊！在带菌量方面，把平板电脑和电子阅读器看成是功能更少、表面面积更大的智能手机。英国消费者杂志《Which？》对一部iPad进行测量，发现了600个金黄色葡萄球菌。","photoid":"B2CMGBDU0AI20009","imgurl":"http://img3.cache.netease.com/photo/0009/2015-08-31/B2CMGBDU0AI20009.jpg"},
				{
					"timgurl":"http://img4.cache.netease.com/photo/0009/2015-08-31/t_B2CMGBHP0AI20009.jpg","photohtml":"http://tech.163.com/photoview/0AI20009/7469.html#p=B2CMGBHP0AI20009","newsurl":"#","squareimgurl":"http://img4.cache.netease.com/photo/0009/2015-08-31/400x400_B2CMGBHP0AI20009.jpg","cimgurl":"http://img4.cache.netease.com/photo/0009/2015-08-31/c_B2CMGBHP0AI20009.jpg","imgtitle":"","simgurl":"http://img3.cache.netease.com/photo/0009/2015-08-31/s_B2CMGBHP0AI20009.jpg","note":"遥控器和附近的细菌电视遥控器可能比公用马桶座干净，但带菌量仍然略高于对家用马桶座的一些估算结果。一项研究发现，遥控器每平方英寸有70个细菌。","photoid":"B2CMGBHP0AI20009","imgurl":"http://img3.cache.netease.com/photo/0009/2015-08-31/B2CMGBHP0AI20009.jpg"},
				{
					"timgurl":"http://img3.cache.netease.com/photo/0009/2015-08-31/t_B2CMGBKJ0AI20009.jpg","photohtml":"http://tech.163.com/photoview/0AI20009/7469.html#p=B2CMGBKJ0AI20009","newsurl":"#","squareimgurl":"http://img3.cache.netease.com/photo/0009/2015-08-31/400x400_B2CMGBKJ0AI20009.jpg","cimgurl":"http://img3.cache.netease.com/photo/0009/2015-08-31/c_B2CMGBKJ0AI20009.jpg","imgtitle":"","simgurl":"http://img4.cache.netease.com/photo/0009/2015-08-31/s_B2CMGBKJ0AI20009.jpg","note":"避无可避？那么，如何避免细菌大量堆积？第一步：要记住，别在卫生间里使用手机（你可不想让手机上满是粪大肠杆菌）。第二步：使用设备前要洗手，或者至少使用免洗洗手液。第三步：使用适当的抹布擦拭，让你心爱的手持设备保持干净。如果你只做其中一步，请选择第二步。(参考消息网)","photoid":"B2CMGBKJ0AI20009","imgurl":"http://img3.cache.netease.com/photo/0009/2015-08-31/B2CMGBKJ0AI20009.jpg"}
			],
			"setname":"五种比马桶圈还脏的电子产品",
			"cover":"http://img4.cache.netease.com/photo/0009/2015-08-31/B2CMGB090AI20009.jpg",
			"commenturl":"http://comment.tech.163.com/photoview_bbs/PHOT079D00090AI2.html",
			"source":"",
			"settag":"马桶圈，电子产品，脏",
			"boardid":"photoview_bbs",
			"tcover":"http://img3.cache.netease.com/photo/0009/2015-08-31/t_B2CMGB090AI20009.jpg",
			"imgsum":"7"
		}
	*/
}

func getNewsList(listId string, page int) ([]*generant.Message, error) {
	newsListPain, err := netTools.Get(listURL(listId, page))
	if err != nil {
		return nil, err
	}
	var v map[string]([]map[string]interface{})
	err = json.Unmarshal(newsListPain, &v)
	if err != nil {
		return nil, err
	}

	//debug_ned, _ := json.MarshalIndent(v, "", "	")

	//log.Debug((string)(debug_ned))

	//log.Debug((string)(newsListPain))

	var result []*generant.Message

	baseInfoList := v[listId]

	var content *generant.Message

	for i, item := range baseInfoList {
		//judge skip type
		if typ, hv := item["skipType"]; hv {
			if skipTyp, can := typ.(string); !can {
				err = errors.New(fmt.Sprint("have skiptype but type is not string , SKIPID[%T : %v]", typ, typ))
			} else {
				switch skipTyp {
				case "photoset":
					content, err = getPhotosetNews(item)
				case "special":
				}
			}
		} else {
			content, err = getNormalNews(item)
		}
		if err != nil {
			log.Warning("at ITEM%d: %d", i, err.Error())
			continue
		}
		result = append(result, content)
		//log.Debug("%d\n", len(result))
	}
	return result, nil
}

/*
http://c.m.163.com/nc/article/B28K7NTL00031H2L/full.html


{
	"B28K7NTL00031H2L":{
		"body":"<!--IMG#0--><p>　　<strong>网易娱乐8月30日报道<\/strong> 据香港媒体报道，张柏芝昨晚获邀到上海担任《环球小姐中国区总决赛》评判，成首位华人女星出任此工作，据悉，她还会替冠军得主进行加冕。<\/p><p>　　首次担任评判的张柏芝说，做评审可能会比较严格，“美不一定要看外表，还要加入许多元素，例如性格、态度、内在美和智慧等。”而恰巧昨天也是张柏芝前夫谢霆锋的35岁生日，他在其facebook上贴出切蛋糕及疑为他亲手制作的多道美食的照片；相中见到有小孩入镜，亦见儿童用的矮桌椅，虽然未能见到样貌，但估计是他的一对宝贝儿子。另外，对于有报道指张柏芝阻霆锋迎娶女友王菲，张柏芝的经纪人在微博澄清，张柏芝衷心祝福他俩的恋情。<\/p>",
		"users":[],
		"replyCount":3546,
		"ydbaike":[],
		"link":[],
		"votes":[],
		"img":[
			{
				"ref":"<!--IMG#0-->",
				"pixel":"550*414",
				"alt":"谢霆锋切蛋糕",
				"src":"http://img1.cache.netease.com/ent/2015/8/30/20150830082244e02dd_550.jpg"
			}
		],
		"digest":"29日是张柏芝前夫谢霆锋的35岁生日，他在其facebook上贴出切蛋糕及疑为他亲手制作的多道美食的照片；相中见到有小孩入镜，被网友预测是他的一对宝贝儿子。",
		"docid":"B28K7NTL00031H2L",
		"title":"谢霆锋疑与两子庆35岁生日",
		"template":"normal1",
		"threadVote":0,
		"threadAgainst":0,
		"boboList":[],
		"replyBoard":"ent2_bbs",
		"source":"网易娱乐",
		"hasNext":false,
		"voicecomment":"off",
		"apps":[],
		"ptime":"2015-08-30 08:24:26"
	}
}
*/

func catchOneTagOnce(tagid string, tag string, pagesOnce int) {
	stmt, err := generant.GetStmtInsert()
	if err != nil {
		log.Error("Error when get STMT, ERROR:[%s]", err.Error())
		return
	}
	defer stmt.Close()

	for i := 0; i < pagesOnce; i++ {
		news, err := getNewsList(tagid, i)
		if err != nil {
			log.Error("Error raisd when catch TAG[%s], PAGENO[%s], ERROR:[%s]", tag, i, err.Error())
			continue
		}
		log.Debug("%d %d", i, len(news))
		for _, item := range news {
			item.Type = tag
			if _, err := item.InsertIntoSQL(stmt); err != nil {
				log.Error("Error when insert into sql : ERROR:[%s]", err.Error())
			}
		}
	}
}

func catchOneTagDaemon(tagid string, tag string, pagesOnce int, delay time.Duration) {
	for {
		catchOneTagOnce(tagid, tag, pagesOnce)
		time.Sleep(delay)
	}
}

type NeteaseNewsCatchConfigure struct {
	Tag      string
	Tagid    string
	PageOnce int
	Delay    time.Duration
}

var (
	configures []NeteaseNewsCatchConfigure = []NeteaseNewsCatchConfigure{
		{KEJI_ID, "科技", 4, time.Minute * 10},
		{TOUTIAO_ID, "头条", 4, time.Minute * 10},
	}
)

func StartCatch() {
	for _, configure := range configures {
		go catchOneTagDaemon(configure.Tag, configure.Tagid, configure.PageOnce, configure.Delay)
		//time.Sleep(time.Minute)
		time.Sleep(time.Second * 10)
	}
}

func init() {
	logging.SetFormatter(logging.MustStringFormatter(
		"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
	))
}
