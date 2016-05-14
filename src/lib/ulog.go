package lib

import (
	"fmt"
	"github.com/astaxie/beego/logs"
)

var ELogger *logs.BeeLogger

func InitLog() { 
	ELogger = logs.NewLogger(5)
	logConfig := fmt.Sprintf(`{"filename":"%s/error.log"}`, AppConfig.Logpath)
	ELogger.SetLogger("file", logConfig)
	ELogger.SetLevel(7)
	ELogger.EnableFuncCallDepth(true)
}


