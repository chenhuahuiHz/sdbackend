package routers

import (
	"sdbackend/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/objlist", &controllers.MainController{}, "*:ObjList")
	beego.Router("/file/download", &controllers.MainController{}, "get:Download")
	beego.Router("/show", &controllers.MainController{}, "get:Show")
	beego.Router("/sd", &controllers.MainController{}, "*:Sd")
	beego.Router("/qrchat", &controllers.MainController{}, "*:Qrchat")
}
