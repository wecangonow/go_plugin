package controllers

import (
	"lib"
	"github.com/astaxie/beego"
	"encoding/json"
	"strconv"
)

type ReportController struct {
	beego.Controller
}

type feedRes struct {
	Uuid            string  `json:"uuid"`
	Ad_id           string  `json:"ad_id"`
	Ad_type         string  `json:"ad_type"`
	Show_num        int     `json:"show_num"`
	Max_Show_num    int     `json:"max_show_num"`
}

type ret struct {
	AdShowTimes feedRes `json:"adShowTimes"`
}

func (u *ReportController) Get() {
	data := u.GetString("data")

	remoteIp := u.Ctx.Request.Header.Get("X-Forwarded-For")
	lib.ALogger.Info("Report: request form %v : %v:", remoteIp, string(data))

	var res map[string]interface{}
	err := json.Unmarshal([]byte(data), &res)


	if err != nil {
		lib.ELogger.Error("Json unmarshal error:", err)
	}
	event  := res["event"]
	uuid   := res["uuid"]
	params := res["params"].(map[string]interface{})
	ad_type, _:= strconv.Atoi(params["ad_type"].(string))
	ad_id,   _:= strconv.Atoi(params["ad_id"].(string))
	countIndex := lib.AdCountIndex{}
	countIndex.Uuid    = uuid.(string)
	countIndex.Ad_id   = ad_id
	countIndex.Ad_type = ad_type

	show_num     := 0
	max_show_num := 0
	if event == "get_ad" && ad_type == 1 {
		show_num, max_show_num = adFeedBack(countIndex)

	}
	if event == "download" {
		show_num, max_show_num = adFeedBack(countIndex)
	}

	feedRes := feedRes{}
	feedRes.Uuid    = uuid.(string)
	feedRes.Ad_id   = params["ad_id"].(string)
	feedRes.Ad_type = params["ad_type"].(string)
	feedRes.Show_num = show_num
	feedRes.Max_Show_num = max_show_num

	ret := ret{}
	ret.AdShowTimes = feedRes

	retJson, err := json.Marshal(ret)

	lib.SetHeader(u.Ctx) //设置http响应头
	lib.ALogger.Info("Report: response to %v : %v ", remoteIp, string(retJson))
	u.Ctx.WriteString(string(retJson))

}

func adFeedBack (countIndex lib.AdCountIndex) (show_num, max_show_num int){
	lib.IncrementUserAdCountByOne(countIndex)
	show_num     = lib.GetAdCount(countIndex, "user")
	max_show_num = lib.GetAdCount(countIndex, "ad")

	return show_num, max_show_num
}
