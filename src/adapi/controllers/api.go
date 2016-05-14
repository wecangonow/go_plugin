package controllers

import (
	"github.com/astaxie/beego"
	"strings"
	"lib"
	"encoding/json"
	"adapi/models"
)

type ApiController struct {
	beego.Controller
}

type AdRequest struct {
	Size   	 string
	Count  	 int
	Web    	 string
	Appid  	 string
	Uid    	 string
	Repeat 	 int
	Times  	 int
	JsType 	 int
	Callback string
}

func (u *ApiController) Get() {
	adReq    := AdRequest{}
	size     := u.GetString("size")
	count,_  := u.GetInt("count")
	web      := u.GetString("web")
	appid    := u.GetString("appid")
	uid      := u.GetString("uid")
	callback := u.GetString("callback")
	repeat,_ := u.GetInt("repeat")
	times,_  := u.GetInt("times")
	jsType,_ := u.GetInt("jsAdType")
	adReq.Appid    = appid
	adReq.Count    = count
	adReq.Repeat   = repeat
	adReq.Size     = size
	adReq.Uid      = uid
	adReq.Web      = web
	adReq.JsType   = jsType
	adReq.JsType   = jsType
	adReq.Times    = times
	adReq.Callback = callback

	remoteAddr := u.Ctx.Request.RemoteAddr
	remoteIp   := strings.Split(remoteAddr,":")[0]

    countryCode, _ := lib.IpToISOCode(remoteIp,"./static/GeoLite2-City.mmdb")

	if countryCode == "" {
		countryCode = "CN"
	}

	var response lib.Response
	data, code, msg := getAd(adReq, countryCode)

	response.Data = data
	response.Code = code
	response.Msg  = msg

	responseJson, err := json.Marshal(response)

	if err != nil {
		lib.ELogger.Error("Json marshal error:", err.Error())
	}


	lib.SetHeader(u.Ctx)  //设置http响应头
	if callback != "" {
		u.Ctx.WriteString(callback + "(" + string(responseJson) + ")")
	} else {
		u.Ctx.WriteString(string(responseJson))
	}

}

/**
 * 2.2
 * 广告调用逻辑
 *
 * @param adRequest  //前段请求get数据
 * @param countryCode //国家ISO码
 * @return data code msg
 */
func getAd(adRequest AdRequest, countryCode string) (data interface{}, code string , msg string) {
	code = "40004"
	msg  = "success"
	data = `{"data":"Hello world"}`
	plugin_info := lib.DataCache.Cache.Get("plugin-" + adRequest.Appid)
	if plugin_info == nil {
		plugin_info,_ = models.GetPPluginByAppId(adRequest.Appid)
		lib.DataCache.Cache.Put("plugin-" + adRequest.Appid, plugin_info, 5)
	}
	return plugin_info, code, msg
}



