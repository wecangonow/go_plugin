package main

import (
	_ "adapi/docs"
	_ "adapi/routers"
	"lib"
	"github.com/astaxie/beego"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	lib.InitConfig()
	lib.InitLog()
	lib.InitOrm()
	go lib.Httpserver.StartHttp()
	beego.Run()
}
