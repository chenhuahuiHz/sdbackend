package routers

import (
	"sdbackend/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/objlist", &controllers.MainController{}, "*:ObjList")
	beego.Router("/file/download", &controllers.MainController{}, "get:Download")
}
