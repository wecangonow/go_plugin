package lib

import (
	"sync"

	"github.com/astaxie/beego/config"
)

var (
	once      sync.Once
	AppConfig Config
	AdConfig  AdConf
)

type AdConf struct{

	DeployLogic []interface{}
	Size        map[string]interface{}
}

type Config struct {
	Httpport                  int
	Appname                   string
	Runmode                   string
	Autorender                bool
	Copyrequestbody           bool
	EnableDocs                bool
	Logpath                   string
	AccessControllAllowOrigin string
	DbConnectstr              string
	Cachetime                 int64
	CountCachetime            int64
	Fabricateplugin           string
	WebUrl                    string
	JsUrl                     string
}

func InitConfig() {
	once.Do(initAllConfig)
}


func initAllConfig() {
	cf, err := config.NewConfig("json","./conf/ad.conf")
	if err != nil {
		ELogger.Error("get json conf error:%s", err.Error())
	}

	var ss []interface{}
	var size map[string]interface{}

	if sp, err := cf.DIY("deploy_logic"); err != nil {
		panic(err)
	} else if m, ok := sp.([]interface{}); ok {
		ss = m
	}
	if sp2, err2 := cf.DIY("size"); err2 != nil {
		ELogger.Error("get size conf err is :", err2)
	} else {
		size = sp2.(map[string]interface{})
	}

	AdConfig.DeployLogic = ss
	AdConfig.Size = size


	cf, err = config.NewConfig("ini", "./conf/app.conf")
	if err != nil {
		ELogger.Error("get ini conf error:%s", err.Error())
	}

	httpport        := cf.DefaultInt("httpport",0)
	appname         := cf.DefaultString("appname", "")
	runmode         := cf.DefaultString("runmode", "")
	autorender      := cf.DefaultBool("autorender", false)
	copyrequestbody := cf.DefaultBool("copyrequestbody", false)
	enabledocs      := cf.DefaultBool("enabledocs", false)
	logpath         := cf.DefaultString("logpath", "")
	access          := cf.DefaultString("access_control_allow_origin", "")
	dbconn          := cf.DefaultString("dbconnect", "")
	cache_time      := cf.DefaultInt64("cache_time",0)
	countcache_time := cf.DefaultInt64("count_cache_time",0)
	plugin_appid    := cf.DefaultString("fabricate_plugin","")
	weburl          := cf.DefaultString("weburl","")
	jsurl           := cf.DefaultString("jsurl","")
	AppConfig = Config{
		Logpath                   : logpath,
		Httpport                  : httpport,
		Appname                   : appname,
		Runmode                   : runmode,
		Autorender                : autorender,
		Copyrequestbody           : copyrequestbody,
		EnableDocs                : enabledocs,
		AccessControllAllowOrigin : access,
		DbConnectstr              : dbconn,
		Cachetime                 : cache_time,
		CountCachetime            : countcache_time,
		Fabricateplugin           : plugin_appid,
		WebUrl                    : weburl,
		JsUrl                     : jsurl,
	}
}
