package netease_news

import (
	//"encoding/json"
	"testing"
)

func TestGetNews(t *testing.T) {
	/*
		news, err := getNewsList(KEJI_ID, 0)
		if err != nil {
			t.Error(err.Error())
		}
		t.Logf("length:%v\n", len(news))
		for _, item := range news {
			msg := item
			pain, _ := json.MarshalIndent(item, "", "	")
			t.Log((string)(pain))

				t.Logf("viewType:%v\n", item.ViewType)
				t.Log("----------new item--------------")
				t.Logf("%v\n", *msg)
				t.Log("----------images--------------")
				for _, img := range msg.Images {
					t.Logf("	%v\n", *img)
				}
				t.Log("----------replys--------------")
				for _, reply := range msg.Replys {
					t.Logf("	%v\n", reply)
				}

		}
	*/
}

func TestStartCatch(t *testing.T) {
	StartCatch()
}
