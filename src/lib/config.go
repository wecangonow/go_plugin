package lib

import (
	"sync"

	"github.com/astaxie/beego/config"
)

var (
	once      sync.Once
	AppConfig Config
	AdConfig  config.ConfigContainer
)


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
}

func InitConfig() {
	once.Do(initAllConfig)
}


func initAllConfig() {
	cf, err := config.NewConfig("json","./conf/ad.conf")
	if err != nil {
		ELogger.Error("get json conf error:%s", err.Error())
	}

	AdConfig = cf


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
	AppConfig = Config{
		Logpath                   : logpath,
		Httpport                  : httpport,
		Appname                   : appname,
		Runmode                   : runmode,
		Autorender                : autorender,
		Copyrequestbody           : copyrequestbody,
		EnableDocs                : enabledocs,
		AccessControllAllowOrigin : access,
		DbConnectstr              :dbconn,
	}
}
