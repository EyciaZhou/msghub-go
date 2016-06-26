package netease_news

type NeteaseNewsChannel struct {
	Name string
	URL  string
	ID   string
}

var (
	//name from http://c.3g.163.com/nc/topicset/default.html
	channelsDefault map[string]*NeteaseNewsChannel = map[string]*NeteaseNewsChannel{
		"toutiao": {
			"头条",
			"http://c.3g.163.com/nc/article/headline/T1295501906343/0-20.html",
			"T1295501906343",
		},
		"keji": {
			"科技",
			"http://c.3g.163.com/nc/article/list/T1348649580692/0-20.html",
			"T1348649580692",
		},
		"shouji": {
			"手机",
			"http://c.3g.163.com/nc/article/list/T1348649654285/0-20.html",
			"T1348649654285",
		},
		"yule": {
			"娱乐",
			"http://c.3g.163.com/nc/article/list/T1348648517839/0-20.html",
			"T1348648517839",
		},
		"caijing": {
			"财经",
			"http://c.3g.163.com/nc/article/list/T1348648756099/0-20.html",
			"T1348648756099",
		},
		"youxi": {
			"游戏",
			"http://c.3g.163.com/nc/article/list/T1348648756099/0-20.html",
			"T1348648756099",
		},
		"lishi": {
			"历史",
			"http://c.3g.163.com/nc/article/list/T1368497029546/0-20.html",
			"T1368497029546",
		},
		"shehui": {
			"社会",
			"http://c.3g.163.com/nc/article/list/T1348648037603/0-20.html",
			"T1348648037603",
		},
		"junshi": {
			"军事",
			"http://c.3g.163.com/nc/article/list/T1348648141035/0-20.html",
			"T1348648141035",
		},
		"dianying": {
			"电影",
			"http://c.3g.163.com/nc/article/list/T1348648650048/0-20.html",
			"T1348648650048",
		},
		"dianshi": {
			"电视",
			"http://c.3g.163.com/nc/article/list/T1348648673314/0-20.html",
			"T1348648673314",
		},
		"tiyu": {
			"体育",
			"http://c.3g.163.com/nc/article/list/T1348649079062/0-20.html",
			"T1348649079062",
		},
		"nba": {
			"NBA",
			"http://c.3g.163.com/nc/article/list/T1348649145984/0-20.html",
			"T1348649145984",
		},
		"zuqiu": {
			"足球",
			"http://c.3g.163.com/nc/article/list/T1348649176279/0-20.html",
			"T1348649176279",
		},
		"lvyou": {
			"旅游",
			"http://c.3g.163.com/nc/article/list/T1348654204705/0-20.html",
			"T1348654204705",
		},
		"qiche": {
			"汽车",
			"http://c.3g.163.com/nc/article/list/T1348654060988/0-20.html",
			"T1348654060988",
		},
		"fangchan": {
			"房产",
			"http://c.3g.163.com/nc/article/list/T1348654085632/0-20.html",
			"T1348654085632",
		},
		"qingsongyike": {
			"轻松一刻",
			"http://c.3g.163.com/nc/article/list/T1350383429665/0-20.html",
			"T1350383429665",
		},
	}
)
