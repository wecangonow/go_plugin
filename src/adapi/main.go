package main

import (
	_ "adapi/docs"
	_ "adapi/routers"
	"lib"
	"github.com/astaxie/beego"
)

func main() {
	lib.InitConfig()
	lib.InitLog()
	lib.InitOrm()
	beego.Run()
}
