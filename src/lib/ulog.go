package lib

import (
	"fmt"
	"github.com/astaxie/beego/logs"
)

var ELogger *logs.BeeLogger
var ALogger *logs.BeeLogger

func InitLog() {
	ELogger = logs.NewLogger(5)
	logConfig := fmt.Sprintf(`{"filename":"%s/error.log"}`, AppConfig.Logpath)
	ELogger.SetLogger("file", logConfig)
	ELogger.SetLevel(7)
	ELogger.EnableFuncCallDepth(true)

    initAlog()
}

func initAlog() {
	ALogger = logs.NewLogger(5)
	logConfig := fmt.Sprintf(`{"filename":"%s/access.log"}`, AppConfig.Logpath)
	ALogger.SetLogger("file", logConfig)
	ALogger.SetLevel(7)
	ALogger.EnableFuncCallDepth(true)

}

