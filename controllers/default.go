package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"sdbackend/models/dbproc"
	"os"
)

type MainController struct {
	beego.Controller
}


func (c *MainController) Get() {
	c.TplName = "index.html"
	beego.Info(c.Data)
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

func (c *MainController) Show() {
	filename := c.GetString("pic")
	beego.Info("show req for:", filename)

	ff, err := os.Open(filename)
	if err != nil {
		beego.Info(err)
		c.Ctx.WriteString(err.Error())
		return
	}

	defer ff.Close()
	buffer := make([]byte, 500000)
	n, _ := ff.Read(buffer)
	c.Ctx.Output.Body(buffer[:n])
}

func (c *MainController) Sd() {
	act := c.GetString("act")
	beego.Info("Sd req for:", act)

	// insert a record
	if act == "start" {
		dbproc.InsertBabyRecord(1, 0, "")
	} else if act == "stop" {
		dbproc.InsertBabyRecord(1, 1, "")
	}

	// select today results
	result := dbproc.SelectBabyRecordOfToday(1)
	data, err := json.Marshal(&result)
	if err != nil {
		beego.Info(err)
	}
	c.Ctx.WriteString(string(data))
}