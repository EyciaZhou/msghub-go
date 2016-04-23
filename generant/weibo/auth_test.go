package weibo

import (
	"encoding/json"
	"testing"
)

var testTweet = `
{"created_at":"Mon Apr 18 12:08:18 +0800 2016","id":3965584557152638,"mid":"3965584557152638","idstr":"3965584557152638","text":"转发微博","source_allowclick":0,"source_type":1,"source":"<a href=\"http://app.weibo.com/t/feed/5yiHuw\" rel=\"nofollow\">iPhone 6 Plus</a>","favorited":false,"truncated":false,"in_reply_to_status_id":"","in_reply_to_user_id":"","in_reply_to_screen_name":"","pic_urls":[],"geo":null,"user":{"id":1401880315,"idstr":"1401880315","class":1,"screen_name":"左耳朵耗子","name":"左耳朵耗子","province":"11","city":"5","location":"北京 朝阳区","description":"芝兰生于深谷，不以无人而不芳；君子修道立德，不为困穷而改节。（ 酷壳：http://coolshell.cn/ ）","url":"http://coolshell.cn","profile_image_url":"http://tp4.sinaimg.cn/1401880315/50/40054262531/1","cover_image":"http://ww3.sinaimg.cn/crop.0.0.980.300/538efefbgw1ecas0gl5cwj20r808cgpb.jpg","cover_image_phone":"http://ww1.sinaimg.cn/crop.0.0.640.640.640/a1d3feabjw1ecat4uqw77j20hs0hsacp.jpg","profile_url":"haoel","domain":"haoel","weihao":"","gender":"m","followers_count":145397,"friends_count":1190,"pagefriends_count":1,"statuses_count":7916,"favourites_count":368,"created_at":"Wed Mar 24 14:57:34 +0800 2010","following":true,"allow_all_act_msg":true,"geo_enabled":true,"verified":true,"verified_type":0,"remark":"","ptype":8,"allow_all_comment":true,"avatar_large":"http://tp4.sinaimg.cn/1401880315/180/40054262531/1","avatar_hd":"http://tva3.sinaimg.cn/crop.27.27.337.337.1024/538efefbgw1eg77da7jggj20aw0aw743.jpg","verified_reason":"程序员，酷壳博主(CoolShell.cn) 微博签约自媒体","verified_trade":"1181","verified_reason_url":"","verified_source":"","verified_source_url":"","verified_state":0,"verified_level":2,"verified_reason_modified":"","verified_contact_name":"","verified_contact_email":"","verified_contact_mobile":"","follow_me":false,"online_status":0,"bi_followers_count":864,"lang":"zh-cn","star":0,"mbtype":0,"mbrank":0,"block_word":0,"block_app":0,"ability_tags":"内容资讯,电子商务,互联网技术,开发者","credit_score":80,"user_ability":24,"urank":31},"retweeted_status":{"created_at":"Mon Apr 18 08:42:16 +0800 2016","id":3965532702543548,"mid":"3965532702543548","idstr":"3965532702543548","text":"React Native 的官方最佳实践，花了两天实践翻译好了，应该是最好的 React Native 学习资料了：http://t.cn/RqowoJY  Facebook 2016 F8 App 的教程，从服务器端到 App，包括 Redux，Relay，GraphQL。教程涵盖：如何进行 app 技术选型，如何做跨平台设计，如何做 React Native 的测试。","textLength":264,"source_allowclick":0,"source_type":1,"source":"<a href=\"http://weibo.com/\" rel=\"nofollow\">微博 weibo.com</a>","favorited":false,"truncated":false,"in_reply_to_status_id":"","in_reply_to_user_id":"","in_reply_to_screen_name":"","pic_urls":[{"thumbnail_pic":"http://ww2.sinaimg.cn/thumbnail/599e230bjw1f30k2odlx6j20zj09paca.jpg"},{"thumbnail_pic":"http://ww2.sinaimg.cn/thumbnail/599e230bjw1f30k8v2zelj207e0b5dgl.jpg"}],"thumbnail_pic":"http://ww2.sinaimg.cn/thumbnail/599e230bjw1f30k2odlx6j20zj09paca.jpg","bmiddle_pic":"http://ww2.sinaimg.cn/bmiddle/599e230bjw1f30k2odlx6j20zj09paca.jpg","original_pic":"http://ww2.sinaimg.cn/large/599e230bjw1f30k2odlx6j20zj09paca.jpg","geo":null,"user":{"id":1503535883,"idstr":"1503535883","class":1,"screen_name":"廖祜秋liaohuqiu_秋百万","name":"廖祜秋liaohuqiu_秋百万","province":"400","city":"1","location":"海外 美国","description":"","url":"http://liaohuqiu.net","profile_image_url":"http://tp4.sinaimg.cn/1503535883/50/5734181567/1","profile_url":"liaohuqiu","domain":"liaohuqiu","weihao":"","gender":"m","followers_count":8512,"friends_count":719,"pagefriends_count":2,"statuses_count":1932,"favourites_count":81,"created_at":"Thu Jul 08 17:12:18 +0800 2010","following":false,"allow_all_act_msg":true,"geo_enabled":true,"verified":true,"verified_type":0,"remark":"","ptype":0,"allow_all_comment":true,"avatar_large":"http://tp4.sinaimg.cn/1503535883/180/5734181567/1","avatar_hd":"http://tva1.sinaimg.cn/crop.0.0.640.640.1024/599e230bjw8euxupc5w5mj20hs0hst9w.jpg","verified_reason":"  淘宝网  职员","verified_trade":"1182","verified_reason_url":"","verified_source":"","verified_source_url":"","verified_state":0,"verified_level":3,"verified_reason_modified":"","verified_contact_name":"","verified_contact_email":"","verified_contact_mobile":"","follow_me":false,"online_status":0,"bi_followers_count":474,"lang":"zh-cn","star":0,"mbtype":12,"mbrank":4,"block_word":1,"block_app":1,"credit_score":80,"user_ability":12,"urank":31},"reposts_count":341,"comments_count":69,"attitudes_count":67,"isLongText":false,"mlevel":0,"visible":{"type":0,"list_id":0},"biz_feature":0,"darwin_tags":[],"hot_weibo_tags":[],"text_tag_tips":[],"userType":0},"annotations":[{"mapi_request":true}],"reposts_count":5,"comments_count":0,"attitudes_count":3,"isLongText":false,"mlevel":0,"visible":{"type":0,"list_id":0},"biz_feature":0,"darwin_tags":[],"hot_weibo_tags":[],"text_tag_tips":[],"rid":"0_0_1_2789440108977006620","userType":0}`

func TestGetSource(t *testing.T) {
	var (
		tweet weiboTweet
	)
	err := json.Unmarshal(([]byte)(testTweet), &tweet)
	if err != nil {
		t.Error(err.Error())
	}

	if tweet.GetSource() != "http://weibo.com/1401880315/Drw9xu0J8" {
		t.Error(tweet.GetSource() + " not match http://weibo.com/1401880315/Drw9xu0J8")
	}
}
