package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"sdbackend/models/dbproc"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) ObjList() {
	result := dbproc.SelectObjListByMainType(dbproc.MyAtoi(c.GetString("maintype")))
	data, err := json.Marshal(&result)
	if err != nil {
		beego.Info(err)
	}
	c.Ctx.WriteString(string(data))
}

func (c *MainController) Download() {
	filename := c.GetString("filename")
	beego.Info("download req for:", filename)
	c.Ctx.Output.Download(filename)
}