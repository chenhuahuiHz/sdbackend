package dbproc

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var LowDreamORM orm.Ormer

func InitSDSql() {
	user := beego.AppConfig.String("dsdb::user")
	passwd := beego.AppConfig.String("dsdb::passwd")
	host := beego.AppConfig.String("dsdb::urls")
	port, err := beego.AppConfig.Int("dsdb::port")
	dbname := beego.AppConfig.String("dsdb::dbname")
	if nil != err {
		port = 3306
	}
	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}

	beego.Info("init mysql ...", user, passwd, host, port, dbname)
	//orm.RegisterDriver("mysql", orm.DRMySQL)
	err = orm.RegisterDataBase("default", "mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, passwd, host, port, dbname))
	err = orm.RegisterDataBase("low_dream", "mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, passwd, host, port, dbname))
	if err != nil {
		beego.Error("init mysql db error.")
		return
	}
	
	beego.Info("init mysql ok")

	LowDreamORM = orm.NewOrm()
	LowDreamORM.Using("low_dream")
}

func SelectObjListByMainType(mainType int8) (rows []ObjRow){

	if nil == LowDreamORM {
		beego.Error("SelectObjListByMainType failed: db not connected")
		return nil
	}
	stable := beego.AppConfig.String("dsdb::tbname")

	num, err := LowDreamORM.Raw(fmt.Sprintf(`
		SELECT * FROM %s where main_type=%d
		`, stable, mainType)).QueryRows(&rows)

	if err == nil {
		fmt.Println("item nums: ", num)
	}

	return rows
}
