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
	//orm.RegisterDriver("mysql", orm.DRMySQL)
	err = orm.RegisterDataBase("default", "mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, passwd, host, port, dbname))
	err = orm.RegisterDataBase("low_dream", "mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, passwd, host, port, dbname))
	if err != nil {
		beego.Error("init mysql db error.")
		return
	}

	LowDreamORM = orm.NewOrm()
	LowDreamORM.Using("low_dream")
}

//func SelectAllLotteryForUser(userID int, date string) (rows []ObjRow){
//
//	if nil == LowDreamORM {
//		beego.Error("SelectAllLotteryForUser failed: db not connected")
//		return nil
//	}
//	stable := beego.AppConfig.String("dsdb::tbpool") + date
//
//	num, err := LowDreamORM.Raw(fmt.Sprintf(`
//		SELECT user_id,type,energy,energy,energy_give,
//		scene,gift_id,gift_count,gift_energy,gift_energygive,
//		gift_getter,time FROM %s where user_id=%d
//		`, stable, userID)).QueryRows(&rows)
//
//	if err == nil {
//		fmt.Println("item nums: ", num)
//	}
//
//	return rows
//}
