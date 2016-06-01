package controllers

import (
	"adapi/models"
	"encoding/json"
	"errors"
	"lib"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type CallAdController struct {
	beego.Controller
}


type AdRequest struct {
	Size     string
	Count    int
	Web      string
	Appid    string
	Uid      string
	Repeat   string
	Times    int
	JsType   int
	Callback string
	Country  string
}

func (u *CallAdController) Get() {

	defer func() {
		u.Ctx.Request.Body.Close()
	}()

	size      := u.GetString("size","0_0")
	count, _  := u.GetInt("count")
	web       := u.GetString("web")
	appid     := u.GetString("appid")
	uid       := u.GetString("uid")
	callback  := u.GetString("plug_callback")
	repeat    := u.GetString("repeat", "")
	times, _  := u.GetInt("times")
	jsType, _ := u.GetInt("jsAdType", 0)
	//web = strings.Replace(web, "http://", "", -1)
	web = parseUrl(web)

	//remoteAddr := u.Ctx.Request.RemoteAddr
	//remoteIp   := strings.Split(remoteAddr, ":")[0]
	remoteIp := u.Ctx.Request.Header.Get("X-Forwarded-For")
	if remoteIp == "" {
		remoteIp = "124.205.66.66"
	}
	countryCode := ""
	if lib.DataCache.Cache.IsExist(remoteIp) {
		countryCode = lib.DataCache.Cache.Get(remoteIp).(string)
	} else {
		countryCode, _ = lib.IpToISOCode(remoteIp, "./static/GeoLite2-City.mmdb")
		lib.DataCache.Cache.Put(remoteIp, countryCode, lib.AppConfig.Cachetime)
	}

	adReq     := AdRequest{}
	adReq.Appid    = appid
	adReq.Count    = count
	adReq.Repeat   = repeat
	adReq.Size     = size
	adReq.Uid      = uid
	adReq.Web      = web
	adReq.JsType   = jsType
	adReq.Times    = times
	adReq.Callback = callback
	adReq.Country  = countryCode


	reqJson, _ := json.Marshal(adReq)
	lib.ALogger.Info("Callad: request form %v : %v:", remoteIp, string(reqJson))

	data := getAd(adReq, countryCode)

	responseJson, err := json.Marshal(data)

	if err != nil {
		lib.ELogger.Error("Json marshal error:", err.Error())
	}

	lib.SetHeader(u.Ctx) //设置http响应头
	lib.ALogger.Info("Callad: response to %v : %v ", remoteIp, string(responseJson))
	if callback != "" {
		u.Ctx.WriteString(callback + "(" + string(responseJson) + ")")
	} else {
		u.Ctx.WriteString(string(responseJson))
	}



}
func parseUrl(url_str string) string {

	host := strings.Replace(url_str, "http://", "", -1)

	return host
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
		plugin_info, _ = models.GetPPluginByAppId(adRequest.Appid)
		lib.DataCache.Cache.Put("plugin-"+adRequest.Appid, plugin_info, lib.AppConfig.Cachetime)
	}
	if adRequest.Appid == lib.AppConfig.Fabricateplugin { //虚拟插件调用
		dau := plugin_info.(*models.PPlugin).Dau
		isCall := plugin_info.(*models.PPlugin).IsCall
		if dau == "" || isCall == 1 {
			return getVirtualPluginAds(adRequest, countryCode)
		} else {
			response.Code = "40007"
			response.Msg = "success"
			response.Data = ""
			return response
		}

	} else { //实体插件
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
		//white = strings.Replace(white.(string), "http://", "", -1)
		white = parseUrl(white.(string))
		if white == web || strings.Contains(web, "."+white.(string)) {
			response.Data = ""
			response.Msg = "success"
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
		//website = strings.Replace(website.(string), "http://", "", -1)
		website = parseUrl(website.(string))
		if website == web {
			website, _ := models.GetPWebsiteByWeb(web)
			webmanageList = append(webmanageList, website)
		}

		if strings.Contains(web, "."+website.(string)) {
			website, _ := models.GetPWebsiteByWeb(website.(string))
			webmanageList = append(webmanageList, website)

		}
	}
	class_slice := make([]int, 0)
	if len(webmanageList) > 0 {
		for _, v := range webmanageList {
			class_slice = append(class_slice, v.ClassId)
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
		//black = strings.Replace(black.(string), "http://", "", -1)
		black = parseUrl(black.(string))
		if black == web || strings.Contains(web, "."+black.(string)) {
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

	config_str := configPercent.(*models.PConfig).ConfigValue
	config_int, _ := strconv.ParseInt(config_str, 10, 32)

	if myrand <= int(config_int) {
		return getAdData(adRequest, countryCode, class_slice)
	} else {
		response.Code = "40006"
		response.Msg = "success"
		return response
	}

	// 按比例调用逻辑--end

	return response
}

func getAdData(adReuest AdRequest, countryCode string, class_slice []int) lib.Response {
	response      := lib.Response{}
	jsAdOut       := models.PJsCode{}
	deployAdList  := []models.PAdList{}
	finalRes      := lib.FinalResponse{}
	finalRes.UrlType = -1

	//获取js广告
	if adReuest.Times == 1 {
		jsAd, err := getJsAdList(adReuest, countryCode, class_slice)
		jsAdOut = jsAd
		if err != nil {
			lib.ELogger.Info("Client get js ad error: %v", err.Error())
		} else {
			jsAdRes := lib.JsAdRes{}
			url := "https://" + lib.AppConfig.WebUrl + lib.AppConfig.JsUrl + strconv.Itoa(jsAd.Id) + ".js?time" + strconv.Itoa(jsAd.UpdateTime)
			jsAdRes.Url = url
			jsAdRes.Ad_type = 1
			jsAdRes.Ad_id = jsAd.Id
			jsAdRes.Ad_name = jsAd.Name
			jsAdRes.Ad_sort = jsAd.Class
			jsAdRes.Ecpm = jsAd.Ecpm
			jsAdRes.Replace = jsAd.Replace
			finalRes.JsAd = jsAdRes
		}

	}
	if adReuest.JsType == 2 || adReuest.JsType == 0 {
		if jsAdOut.Id == 0 {
			deployAdList = getDeployAdList(adReuest, countryCode, class_slice)
		} else {
			if jsAdOut.Replace == 2 {
				deployAdList = getDeployAdList(adReuest, countryCode, class_slice)
			}
		}
	}

	if jsAdOut.Id == 0 && len(deployAdList) == 0 {
		response.Code = "40004"
		response.Msg  = "success"
		return response
	}

	finalRes.DeployAd = configData(adReuest, deployAdList)
	response.Code = "20000"
	response.Msg  = "success"

	// 获取网址管理里面网站的分类--start
	// 没有则默认Url_type为-1
	websiteLists := lib.DataCache.Cache.Get("websiteLists")
	if websiteLists == nil {
		websiteLists, _ = models.GetMyAllPWebsite()
		lib.DataCache.Cache.Put("websiteLists", websiteLists, lib.AppConfig.Cachetime)
	}

	for _, website := range websiteLists.(orm.ParamsList) {
		//website = strings.Replace(website.(string), "http://", "", -1)
		website = parseUrl(website.(string))
		if website == adReuest.Web {
			website, _ := models.GetPWebsiteByWeb(adReuest.Web)
			//lib.ELogger.Debug("website is %v", website)
			if website.ClassId != 0 {
				finalRes.UrlType = website.ClassId
			}
		}

	}
	response.Data = finalRes

	return response
}

func getJsAdList(adReuest AdRequest, countryCode string, class_slice []int) (models.PJsCode, error) {
	deploy_logic := lib.AdConfig.DeployLogic
	ret := models.PJsCode{}
	for _, v := range deploy_logic {
		country := countryCode + "0"
		class := ""
		where := deployLogic(adReuest, int(v.(float64)), countryCode, class_slice)
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
			adList = adList.([]models.PJsCode)
			lib.DataCache.Cache.Put(key, adList, lib.AppConfig.Cachetime)
		}
		for k, v := range adList.([]models.PJsCode) {
			webBlackStr := v.Webblack
			if webBlackStr != "" {
				webBlackSlice := strings.Split(webBlackStr, "\r\n")
				if len(webBlackSlice) > 0 {
					for _, web := range webBlackSlice {
						//web = strings.Replace(web, "http://", "", -1)
						web = parseUrl(web)
						if adReuest.Web == web || strings.Contains(adReuest.Web, "."+web) {
							adList = append(adList.([]models.PJsCode)[:k], adList.([]models.PJsCode)[k+1:]...)
						}
					}
				}
			}
		}
		adList = randomByEcpm(adList.([]models.PJsCode), nil)
		adList = checkAdShowTimes(adReuest, adList.([]models.PJsCode), nil, 1)
		if adList.(models.PJsCode).AppId != "" {
			ret = adList.(models.PJsCode)
			return ret, nil
		}
	}
	err := errors.New("proper js ad not found")
	return ret, err
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
	uuid := adRequest.Uid
	if adType == 2 {
		if len(deployAdList) > 0 {
			for _, v := range deployAdList {
				countIndex := lib.AdCountIndex{}
				countIndex.Uuid = uuid
				countIndex.Ad_id = v.Id
				countIndex.Ad_type = adType
				ad_num := lib.GetAdCount(countIndex, "ad")
				user_ad_num := lib.GetAdCount(countIndex, "user")
				if ad_num < v.MaxShowNum {
					if user_ad_num < v.ShowNum {
						return v
					}
				}
			}
		}

	}

	ret := oneAdCheck(jsAdList, uuid, 1)

	return ret
}
func oneAdCheck(jsAdList []models.PJsCode, uuid string, adType int) models.PJsCode {
	if len(jsAdList) > 0 {
		for _, v := range jsAdList {
			countIndex := lib.AdCountIndex{}
			countIndex.Uuid = uuid
			countIndex.Ad_id = v.Id
			countIndex.Ad_type = adType
			ad_num := lib.GetAdCount(countIndex, "ad")
			user_ad_num := lib.GetAdCount(countIndex, "user")
			if ad_num < v.MaxShowNum {
				if user_ad_num < v.ShowNum {
					return v
				}
			}
		}

	}

	return models.PJsCode{}

}

func getDeployAdList(adReuest AdRequest, countryCode string, class_slice []int) []models.PAdList {
	ret := []models.PAdList{}
	deploy_logic := lib.AdConfig.DeployLogic
	sizeConfig   := lib.AdConfig.Size
	reqSize      := strings.Split(adReuest.Size, "_")
	repeat := make([]string,0)
	if adReuest.Repeat != "" {
		repeat = strings.Split(adReuest.Repeat, "_")
	}
	if len(reqSize) > 0 {
		for _, rv := range reqSize {
			for ck, cv := range sizeConfig {
				if rv == cv {
					for _, dlv := range deploy_logic {
						country := countryCode + "0"
						class := ""
						where := deployLogic(adReuest, int(dlv.(float64)), countryCode, class_slice)
						where += " and `size`=" + ck
						if len(repeat) > 0 {
							//lib.ELogger.Debug("repeat to string %v", strings.Join(repeat[:], ","))
							//lib.ELogger.Debug("repeat to string %v", repeat)
							//lib.ELogger.Debug("repeat len is %v ", len(repeat))


							inwhere := "(" + strings.Join(repeat[:], ",") + ")"
							where += " and `id` not in " + inwhere
						}
						if int(dlv.(float64)) == 1 {
							if len(class_slice) > 0 {
								for _, class_id := range class_slice {
									class += strconv.Itoa(class_id) + "-"
								}
								class += "0"
							}

						}

						if int(dlv.(float64)) == 2 {
							if len(class_slice) > 0 {
								for _, class_id := range class_slice {
									class += strconv.Itoa(class_id) + "-"
								}
							}
						}

						key := "DeployAd_" + strconv.Itoa(int(dlv.(float64))) + "_" + ck + "_" + adReuest.Appid + "_" + country + "_" + class + "_" + strings.Join(repeat[:], "-")

						adList := lib.DataCache.Cache.Get(key)
						if adList == nil {
							adList, _ = models.GetAllPAdListByWhere(where)
							adList = adList.([]models.PAdList)
							lib.DataCache.Cache.Put(key, adList, lib.AppConfig.Cachetime)
						}
						if len(adList.([]models.PAdList)) > 0 {
							for k, v := range adList.([]models.PAdList) {
								webBlackStr := v.Webblack
								if webBlackStr != "" {
									webBlackSlice := strings.Split(webBlackStr, "\r\n")
									if len(webBlackSlice) > 0 {
										for _, web := range webBlackSlice {
											//web = strings.Replace(web, "http://", "", -1)
											web = parseUrl(web)
											if adReuest.Web == web || strings.Contains(adReuest.Web, "."+web) {
												adList = append(adList.([]models.PAdList)[:k], adList.([]models.PAdList)[k+1:]...)
											}
										}
									}
								}
							}

							adList = randomByEcpm(nil, adList.([]models.PAdList))
							adList = checkAdShowTimes(adReuest, nil, adList.([]models.PAdList), 2)
							if adList.(models.PAdList).Id == 0 {
								continue
							} else {
								ret = append(ret, adList.(models.PAdList))
								break
							}
						}

					}
				}
			}

		}
	}

	return ret
}

func configData(adRequest AdRequest, adList []models.PAdList) []lib.DeployAdRes {
	ret := []lib.DeployAdRes{}

	if len(adList) > 0 {
		deployRes := lib.DeployAdRes{}
		for _, ad := range adList {
				size,_ := lib.AdConfig.Size[strconv.Itoa(ad.Size)]
				urlSlice := strings.Split(ad.StaticUrl, ".")
				htmlUrl := "//" + ad.StaticUrl + "/" + urlSlice[1] + "/" + strconv.Itoa(ad.Id) + ".html?time=" + strconv.Itoa(ad.UpdateTime) + "&host=" + adRequest.Web +
				"&ad_name=" + ad.Name + "&size=" + size.(string) + "&ad_id=" + strconv.Itoa(ad.Id) + "&ad_type=2#" + adRequest.Uid
				deployRes.Size = size.(string)
				deployRes.Ad_id = ad.Id
				deployRes.Ad_type = 2
				deployRes.Ecpm = ad.Ecpm
				deployRes.Ad_sort = ad.Class
				deployRes.Ad_name = ad.Name
				deployRes.Web = adRequest.Web
				deployRes.Url = htmlUrl
				ret = append(ret, deployRes)

		}
	}
	return ret
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
		//where += " and app_id=" + adRequest.Appid
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

		//where += " and app_id=" + adRequest.Appid
	}

	return where
}
func getTruePluginAds() string {
	return "true addon"
}
