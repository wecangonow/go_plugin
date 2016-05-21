package controllers

import (
	"github.com/astaxie/beego"
	"strings"
	"lib"
	"encoding/json"
	"adapi/models"
	"github.com/astaxie/beego/orm"
	"strconv"
	"sort"
	"math/rand"
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
	jsType,_ := u.GetInt("jsAdType",0)
	adReq.Appid    = appid
	adReq.Count    = count
	adReq.Repeat   = repeat
	adReq.Size     = size
	adReq.Uid      = uid
	adReq.Web      = web
	adReq.JsType   = jsType
	adReq.Times    = times
	adReq.Callback = callback

	remoteAddr := u.Ctx.Request.RemoteAddr
	remoteIp   := strings.Split(remoteAddr,":")[0]

	countryCode, _ := lib.IpToISOCode(remoteIp,"./static/GeoLite2-City.mmdb")

	if countryCode == "" {
		countryCode = "CN"
	}

	data := getAd(adReq, countryCode)

	responseJson, err := json.Marshal(data)

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
func getAd(adRequest AdRequest, countryCode string) lib.Response {
	response := lib.Response{}
	plugin_info := lib.DataCache.Cache.Get("plugin-" + adRequest.Appid)
	if plugin_info == nil {
		plugin_info,_ = models.GetPPluginByAppId(adRequest.Appid)
		lib.DataCache.Cache.Put("plugin-" + adRequest.Appid, plugin_info, lib.AppConfig.Cachetime)
	}
	if adRequest.Appid == lib.AppConfig.Fabricateplugin {   //虚拟插件调用
		dau    := plugin_info.(*models.PPlugin).Dau
		isCall := plugin_info.(*models.PPlugin).IsCall
		if dau == "" || isCall == 1 {
			return getVirtualPluginAds(adRequest, countryCode)
		} else {
			response.Code = "40007"
			response.Msg  = "success"
			response.Data = ""
			return response
		}

	} else {  //实体插件
	}

	return response
}

//获取虚拟插件请求的广告
func getVirtualPluginAds(adRequest AdRequest, countryCode string) lib.Response {
	response := lib.Response{}
	web := adRequest.Web
	// 白名单验证--start
	whiteLists := lib.DataCache.Cache.Get("whiteLists")
	if whiteLists == nil {
		whiteLists, _ = models.GetMyAllPWhite()
		lib.DataCache.Cache.Put("whiteLists", whiteLists, lib.AppConfig.Cachetime)
	}
	for _, white := range whiteLists.(orm.ParamsList) {
		if white == web || strings.Contains(web, "." + white.(string)) {
			response.Data = ""
			response.Msg  = "success"
			response.Code = "40005"
			return response
		}
	}

	// 白名单验证--end

	// 网址管理验证--start
	websiteLists := lib.DataCache.Cache.Get("websiteLists")
	if websiteLists == nil {
		websiteLists, _ = models.GetMyAllPWebsite()
		lib.DataCache.Cache.Put("websiteLists", websiteLists, lib.AppConfig.Cachetime)
	}

	var webmanageList []*models.PWebsite
	for _, website := range websiteLists.(orm.ParamsList) {
		if website == web  {
			website, _ := models.GetPWebsiteByWeb(web)
			webmanageList = append(webmanageList, website)
		}

		if  strings.Contains(web, "." + website.(string)) {
			website, _ := models.GetPWebsiteByWeb(website.(string))
			webmanageList = append(webmanageList, website)

		}
	}
	class_slice := make([]int,0)
	if len(webmanageList) > 0 {
		for _, v := range webmanageList {
			class_slice = append(class_slice,v.ClassId)
		}

		return getAdData(adRequest, countryCode, class_slice)
	}
	// 黑名单逻辑--start

	blackLists := lib.DataCache.Cache.Get("blackLists")

	if blackLists == nil {
		blackLists, _ = models.GetMyAllPGrey()
		lib.DataCache.Cache.Put("blackLists", blackLists, lib.AppConfig.Cachetime)
	}
	for _, black := range blackLists.(orm.ParamsList) {
		if black == web || strings.Contains(web, "." + black.(string)) {
			return getAdData(adRequest, countryCode, class_slice)
		}
	}

	// 黑名单逻辑--end
	// 按比例调用逻辑--start
	myrand := lib.Random(1, 10)

	configPercent := lib.DataCache.Cache.Get("configPercent")

	if configPercent == nil {
		configPercent, _ = models.GetPConfigByName("grey_percentage")
		lib.DataCache.Cache.Put("configPercent", configPercent, lib.AppConfig.Cachetime)
	}

	//lib.ELogger.Info("myrand is %v", myrand)
	//lib.ELogger.Info("percent is %v", configPercent)

	config_str := configPercent.(*models.PConfig).ConfigValue
	config_int, _ := strconv.ParseInt(config_str, 10, 32)

	if myrand <= int(config_int) {
		return getAdData(adRequest, countryCode, class_slice)
	} else {
		response.Code = "40006"
		response.Msg  = "success"
		return response
	}

	// 按比例调用逻辑--end

	return response
}


func getAdData(adReuest AdRequest, countryCode string, class_slice []int) lib.Response {
	response := lib.Response{}
	 response.Data = deployLogic(adReuest, 2, countryCode, class_slice)
	//获取js广告
	if adReuest.Times == 1 {
		_ = getJsAdList(adReuest, countryCode, class_slice)
	}
	if adReuest.JsType == 2 || adReuest.JsType == 0 {
	}
//	if ($params['times'] == 1) {
//	//获取js广告数据
//	$JsAdList = self::JsAdList($params, $appconfig, $countryStr, $classArr);
//	if (!empty($JsAdList)) {
//	//js广告
//	$deployArr['jsAd']['url'] = '//'.Yii::$app->params['webUrl'] . Yii::$app->params['jsUrl'] . $JsAdList['id'] . '.js?time=' . $JsAdList['update_time'];
//	$deployArr['jsAd']['ad_id'] = $JsAdList['id'];
//	$deployArr['jsAd']['ad_name'] = $JsAdList['name'];
//	$deployArr['jsAd']['replace'] = $JsAdList['replace'];
//	$deployArr['jsAd']['ecpm'] = $JsAdList['ecpm'];
//	$deployArr['jsAd']['ad_type'] = 1;
//	$deployArr['jsAd']['AD_sort'] = $JsAdList['class'];
//	}
//	}
//	if($params['jsAdType']==2 || empty($params['jsAdType'])){
//	if (Yii::$app->params['fabricate_plugin'] === $params['appid']) {
//	//获取配置广告
//	if (!empty($JsAdList) && $JsAdList['replace'] == 2) {
//	//获取配置广告的广告
//	$deployAdList = self::DeployAdList($params, $appconfig, $countryStr, $classArr);
//	} elseif (empty($JsAdList)) {
//	//获取配置广告的广告
//	$deployAdList = self::DeployAdList($params, $appconfig, $countryStr, $classArr);
//	}
//	if (!empty($deployAdList)) {
//	$deployArr['deployAd'] = self::ConfigurationData($params, $deployAdList, $appconfig);
//	}
//	}
//	}
//	//添加前端传递来web地址在网址管理里面的分类值如果没有则为-1
//	$web = $params['web'];
//$web_info = LNbApi::checkWebManage($web);
//if(is_array($web_info)){
//if(isset($web_info['class_id'])){
//$deployArr['URL_type'] = $web_info['class_id'];
//} else {
//$deployArr['URL_type'] = -1;
//}
//} else {
//$deployArr['URL_type'] = -1;
//}
//
//if (empty($JsAdList) && empty($deployArr['deployAd'])) { //没有对应的配置广告
//return self::returnAction($callback, '', 'success', '40004');
//} else {
//return self::returnAction($callback, $deployArr, 'success', '20000');
//}

	return response
}


func getJsAdList(adReuest AdRequest, countryCode string, class_slice []int) string {
	deploy_logic := lib.AdConfig.DeployLogic

	for _, v := range deploy_logic {
		country := countryCode + "0"
		class   := ""
		where   := deployLogic(adReuest, int(v.(float64)), countryCode, class_slice)
		if int(v.(float64)) == 1 {
			if len(class_slice) > 0 {
				for _, class_id := range class_slice {
					class += strconv.Itoa(class_id) + "-"
				}
				class += "0"
			}

		}

        if int(v.(float64)) == 2 {
			if len(class_slice) > 0 {
				for _, class_id := range class_slice {
					class += strconv.Itoa(class_id) + "-"
				}
			}
		}

		key := "JsAd_" + strconv.Itoa(int(v.(float64))) + "_" + adReuest.Appid + "_" + country + "_" + class
		adList := lib.DataCache.Cache.Get(key)
		if adList == nil {
			adList, _ = models.GetAllPJsCodeByWhere(where)
			adList    = adList.([]models.PJsCode)
			lib.DataCache.Cache.Put(key, adList, lib.AppConfig.Cachetime)
		}
		for k, v := range adList.([]models.PJsCode) {
			webBlackStr := v.Webblack
			if webBlackStr != "" {
				webBlackSlice := strings.Split(webBlackStr,"\r\n")
				if len(webBlackSlice) > 0 {
					for _, web := range webBlackSlice {
						if adReuest.Web == web || strings.Contains(adReuest.Web, "." + web) {
							adList = append(adList.([]models.PJsCode)[:k], adList.([]models.PJsCode)[k+1:]...)
						}
					}
				}
			}
		}
		//lib.ELogger.Debug("adList before %v", adList)
		adList = randomByEcpm(adList.([]models.PJsCode), nil)
		//lib.ELogger.Debug("adList len %v", len(adList.([]models.PJsCode)))
		//lib.ELogger.Debug("adList after random %v", adList)
		adList = checkAdShowTimes(adReuest, adList.([]models.PJsCode), nil, 1)
		//lib.ELogger.Debug("key is %v", key)
		//lib.ELogger.Debug("web is %v", adReuest.Web)
		//lib.ELogger.Debug("class is %v", class)
		//lib.ELogger.Debug("country is %v", country)
		//lib.ELogger.Debug("where is %v", where)
	}
	return ""
}
func randomByEcpm(jsAdList []models.PJsCode, deployAdList []models.PAdList) interface{} {

	if jsAdList != nil {
		jsMap := map[float64][]models.PJsCode{}
		for _, v := range jsAdList {
			mapV, ok := jsMap[v.Ecpm]
			if ok {
				mapV = append(mapV, v)
				jsMap[v.Ecpm] = mapV
			} else {
				tmp := make([]models.PJsCode, 0)
				tmp = append(tmp, v)
				jsMap[v.Ecpm] = tmp
			}
		}

		keys := make([]float64, len(jsMap))
		i := 0
		for k, v := range jsMap {
			keys[i] = k
			i++
			for x := range v {
				j := rand.Intn(x + 1)
				v[x], v[j] = v[j], v[x]
			}
		}
		sort.Sort(sort.Reverse(sort.Float64Slice(keys)))

		ret := make([]models.PJsCode, 0)
		for _, v := range keys {
			for _, mapV := range jsMap[v] {
				ret = append(ret, mapV)
			}
		}

		return ret
	}

	if deployAdList != nil {
		adMap := map[float64][]models.PAdList{}
		for _, v := range deployAdList {
			mapV, ok := adMap[v.Ecpm]
			if ok {
				mapV = append(mapV, v)
				adMap[v.Ecpm] = mapV
			} else {
				tmp := make([]models.PAdList, 0)
				tmp = append(tmp, v)
				adMap[v.Ecpm] = tmp
			}
		}

		keys := make([]float64, len(adMap))
		i := 0
		for k, v := range adMap {
			keys[i] = k
			i++
			for x := range v {
				j := rand.Intn(x + 1)
				v[x], v[j] = v[j], v[x]
			}
		}
		sort.Sort(sort.Reverse(sort.Float64Slice(keys)))

		ret := make([]models.PAdList, 0)
		for _, v := range keys {
			for _, mapV := range adMap[v] {
				ret = append(ret, mapV)
			}
		}

		return ret

	}
	return ""
}
func checkAdShowTimes(adRequest AdRequest, jsAdList []models.PJsCode, deployAdList []models.PAdList, adType int) interface{} {
	//uuid := adRequest.Uid
	if adType == 2 {

	}


	return ""
}
//func oneAdCheck(jsAdList []models.PJsCode, uuid string, adType int) models.PJsCode {
//
//}

func getDeployAdList(adReuest AdRequest, countryCode string, class_slice []int) {

}
// 1 地区+类别+不限 2 地区+其他类别
func deployLogic(adRequest AdRequest, deployType int, countryCode string, class_slice []int) string {
	var where string
	switch deployType {
	case 1:
		where = " and (`country` like '%" + countryCode + "%' or `country`='0')"
		if len(class_slice) > 0 {
			where += " and (`class`='0'"
			for _, v := range class_slice {
				where += " or `class` like '%" + strconv.Itoa(v) + "%'"
			}
			where += ")"
			where += " and ("
			for k, v := range class_slice {
				if k == 0 {
					where += " `unclass` not like '%" + strconv.Itoa(v) + "%'"
				} else {
					where += " and `unclass` not like '%" + strconv.Itoa(v) + "%'"
				}
			}
			where += ")"
		} else {
			where += " and putin=2"
		}
		where += " and app_id=" + adRequest.Appid
	case 2:
		where = " and (`country` like '%" + countryCode + "%' or `country`='0')"
		if len(class_slice) > 0 {
			where += " and (`class`!='0'"
			for _, v := range class_slice {
				where += " or `class` like '%" + strconv.Itoa(v) + "%'"
			}
			where += ")"
			where += " and ("
			for k, v := range class_slice {
				if k == 0 {
					where += " `unclass` not like '%" + strconv.Itoa(v) + "%'"
				} else {
					where += " and `unclass` not like '%" + strconv.Itoa(v) + "%'"
				}
			}
			where += ")"
		} else {
			where += " and putin=2"
		}

		where += " and app_id=" + adRequest.Appid
	}

	return where
}
func getTruePluginAds() string {
	return "true addon"
}

