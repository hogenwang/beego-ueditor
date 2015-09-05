package routers

import (
	"github.com/astaxie/beego"
	"github.com/hogenwang/beego-ueditor/controllers"
)

func init() {
	//注册静态文件
	beego.SetStaticPath("/attach", "attach")
	beego.Router("/", &controllers.MainController{})
	//UEditor
	beego.Router("/ueditor/go/controller", &controllers.UEController{}, "*:UEditor")

}
