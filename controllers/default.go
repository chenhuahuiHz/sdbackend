package controllers

import (
	"encoding/json"
	"os"
	"sdbackend/models/dbproc"

	"github.com/astaxie/beego"
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
	buffer := make([]byte, 1000000)
	n, _ := ff.Read(buffer)
	c.Ctx.Output.Body(buffer[:n])
}

func (c *MainController) Sd() {
	act := c.GetString("act")
	beego.Info("Sd req for:", act)

	// insert a record
	if act == "start" {
		// 先结束空闲记录
		dbproc.StopBabyRecord(0)
		// 再开始
		dbproc.StartBabyRecord(1)
	} else if act == "stop" {
		dbproc.StopBabyRecord(1)
		// 结束后插入一条空闲记录
		dbproc.StartBabyRecord(0)
	} else if act == "statist" {
		c.Ctx.WriteString(dbproc.StatistBabyRecord(1))
		return
	}

	// select today results
	babyRows := dbproc.SelectBabyRecordOfToday(1)
	data, err := json.Marshal(&babyRows)
	if err != nil {
		beego.Info(err)
	}
	c.Ctx.WriteString(string(data))
}

func (c *MainController) Qrchat() {
	txt := c.GetString("txt")
	beego.Info("Qrchat req for:", txt)
	dbproc.InsertQrchatRecord(txt)
	c.Ctx.WriteString("ok")
}
