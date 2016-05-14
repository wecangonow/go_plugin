package routers

import (
	"adapi/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/api/call-ad",&controllers.ApiController{})
}
