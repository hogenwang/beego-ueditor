package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "UEditor Demo"
	c.Data["Email"] = "hogenwang@vip.qq.com"
	c.TplNames = "index.tpl"
}
