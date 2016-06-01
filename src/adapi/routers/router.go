package routers

import (
	"adapi/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/api/call-ad",&controllers.CallAdController{})
	beego.Router("/api/report",&controllers.ReportController{})
}
