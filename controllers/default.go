package controllers

import (
	"github.com/astaxie/beego"
	"imdata/models/dbproc"
	//"encoding/json"
	//"strings"
	//"prototcp/typedefs"
)

type MainController struct {
	beego.Controller
}

//func (c *MainController) Get() {
//	c.TplName = "index.tpl"
//
//	beego.Info(c.Data)
//}

func (c *MainController) ObjList() {
	//c.TplName = "Test.tpl"
	beego.Info(c.Data)
	//sqlstr := c.GetString("sql")
	//date := strings.Replace(c.GetString("date"), "-", "", -1)
	//beego.Info(date, sqlstr)
	//
	//result := dbproc.SelectLotteryWithSql(sqlstr, date)
	//
	//data, err := json.Marshal(&result)
	//if err != nil {
	//	beego.Info(err)
	//}
	////beego.Info("data...", string(data))
	//c.Ctx.WriteString(string(data))
}
