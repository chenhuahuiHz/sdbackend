package dbproc

import (
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var LowDreamORM orm.Ormer
var sqlCache SqlCache

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
	sqlCache.init()
}

func SelectObjListByMainType(mainType int) (rows []ObjRow) {

	beego.Info("SelectObjListByMainType ...", mainType)

	stable := beego.AppConfig.String("dsdb::tbname")
	sql := ""
	if mainType > 0 {
		sql = fmt.Sprintf(`SELECT * FROM %s where main_type=%d`, stable, mainType)
	} else {
		sql = fmt.Sprintf(`SELECT * FROM %s`, stable) //for all
	}

	cache, exist := sqlCache.getObjListCache(sql)
	if exist {
		rows = cache
		beego.Info("hit cache for sql:", sql)
		return
	}

	if nil == LowDreamORM {
		beego.Error("SelectObjListByMainType failed: db not connected")
		return nil
	}
	num, err := LowDreamORM.Raw(sql).QueryRows(&rows)
	if err == nil {
		beego.Info(sql, "get item nums:", num)
		sqlCache.setObjListCache(sql, rows)
	}
	return rows
}

func StartBabyRecord(t int8) {
	beego.Info("StartBabyRecord ...", t)

	stable := beego.AppConfig.String("dsdb::tbbaby")
	sql := fmt.Sprintf(`insert into %s (type) values(%d)`, stable, t)

	if nil == LowDreamORM {
		beego.Error("StartBabyRecord failed: db not connected")
		return
	}
	_, err := LowDreamORM.Raw(sql).Exec()
	if err != nil {
		beego.Info("StartBabyRecord ... err:", err.Error())
	}
}

func StopBabyRecord(t int8) {
	beego.Info("StopBabyRecord ...", t)

	stable := beego.AppConfig.String("dsdb::tbbaby")
	// UPDATE t_baby_rcd SET stop_time=NOW(), cost_seconds=UNIX_TIMESTAMP(NOW())-UNIX_TIMESTAMP(start_time), state=1 WHERE state=0 and type=1 ORDER BY id desc LIMIT 1;
	sql := fmt.Sprintf(`UPDATE %s SET stop_time=NOW(), cost_seconds=UNIX_TIMESTAMP(NOW())-UNIX_TIMESTAMP(start_time), state=1 WHERE state=0 and type=%d ORDER BY id desc LIMIT 1;`, stable, t)

	if nil == LowDreamORM {
		beego.Error("StopBabyRecord failed: db not connected")
		return
	}
	_, err := LowDreamORM.Raw(sql).Exec()
	if err != nil {
		beego.Info("StopBabyRecord ... err:", err.Error())
	}
}

func SelectBabyRecordOfToday(t int8) (rows []BabyRow) {
	beego.Info("SelectBabyRecordOfToday ...", t)

	stable := beego.AppConfig.String("dsdb::tbbaby")

	// return last 15 items is better
	//sql := fmt.Sprintf(`SELECT * FROM %s where type=%d and start_time>='%s' ORDER BY id desc`, stable, t, time.Now().Format("2006-01-02")+" 00:00:00")
	sql := fmt.Sprintf(`SELECT * FROM %s where type=%d ORDER BY id desc limit 15`, stable, t)

	if nil == LowDreamORM {
		beego.Error("SelectBabyRecordOfToday failed: db not connected")
		return nil
	}
	num, err := LowDreamORM.Raw(sql).QueryRows(&rows)
	beego.Info(sql, "SelectBabyRecordOfToday get item nums:", num)
	if err != nil {
		beego.Info("SelectBabyRecordOfToday ... err:", err.Error())
	}
	return rows
}

func formatSeconds(value int) string {
	secondTime := value
	minuteTime := 0
	hourTime := 0

	if secondTime > 60 {
		minuteTime = secondTime / 60
		secondTime = secondTime % 60
		if minuteTime > 60 {
			hourTime = minuteTime / 60
			minuteTime = minuteTime % 60
		}
	}
	result := fmt.Sprintf(`%d秒`, secondTime)

	if minuteTime > 0 {
		result = fmt.Sprintf(`%d分%s`, minuteTime, result)
	}
	if hourTime > 0 {
		result = fmt.Sprintf(`%d小时%s`, hourTime, result)
	}
	return result
}

func StatistBabyRecord(t int8) (txt string) {

	if nil == LowDreamORM {
		beego.Error("SelectBabyRecordOfToday failed: db not connected")
		return "Error： db not connected"
	}

	// "昨天共饲养%d次，总时长%s，平均每次%s\n"
	// "今天已饲养%d次，总时长%s，平均每次%s"
	today := time.Now().Format("2006-01-02") + " 00:00:00"
	yestorday := time.Now().AddDate(0, 0, -1).Format("2006-01-02") + " 00:00:00"

	yestordaySql := fmt.Sprintf(`SELECT count(0) as count, sum(cost_seconds) as sum FROM %s WHERE type = %d and state=1 and start_time >= '%s' and start_time < '%s'`,
		beego.AppConfig.String("dsdb::tbbaby"), t, yestorday, today)
	todaySql := fmt.Sprintf(`SELECT count(0) as count, sum(cost_seconds) as sum FROM %s WHERE type = %d and state=1 and start_time >= '%s'`,
		beego.AppConfig.String("dsdb::tbbaby"), t, today)

	var yestordayResult StatistRow
	err := LowDreamORM.Raw(yestordaySql).QueryRow(&yestordayResult)
	if err != nil {
		beego.Info("StatistBabyRecord ... err:", err.Error())
		return err.Error()
	}

	var todayResult StatistRow
	err = LowDreamORM.Raw(todaySql).QueryRow(&todayResult)
	if err != nil {
		beego.Info("StatistBabyRecord ... err:", err.Error())
		return err.Error()
	}

	beego.Info("StatistBabyRecord ... ", yestordayResult, todayResult)

	avgYestorday, avgToday := 0, 0
	if yestordayResult.Count > 0 {
		avgYestorday = int(yestordayResult.Sum / yestordayResult.Count)
	}
	if todayResult.Count > 0 {
		avgToday = int(todayResult.Sum / todayResult.Count)
	}

	var ysetordayInter StatistRow
	err = LowDreamORM.Raw(fmt.Sprintf(`SELECT count(0) as count, sum(cost_seconds) as sum FROM %s WHERE type = 0 and state=1 and start_time >= '%s' and start_time < '%s'`,
		beego.AppConfig.String("dsdb::tbbaby"), yestorday, today)).QueryRow(&ysetordayInter)
	if err != nil {
		beego.Info("StatistBabyRecord ... err:", err.Error())
		return err.Error()
	}

	var todayInter StatistRow
	err = LowDreamORM.Raw(fmt.Sprintf(`SELECT count(0) as count, sum(cost_seconds) as sum FROM %s WHERE type = 0 and state=1 and start_time >= '%s'`,
		beego.AppConfig.String("dsdb::tbbaby"), today)).QueryRow(&todayInter)
	if err != nil {
		beego.Info("StatistBabyRecord ... err:", err.Error())
		return err.Error()
	}
	avgInterYestorday, avgInterToday := 0, 0
	if ysetordayInter.Count > 0 {
		avgInterYestorday = int(ysetordayInter.Sum / ysetordayInter.Count)
	}
	if todayInter.Count > 0 {
		avgInterToday = int(todayInter.Sum / todayInter.Count)
	}

	// txt = fmt.Sprintf(`<font color='red'>昨天</font>共饲养%d次，总时长%s，平均每次%s。<br><font color='green'>今天</font>已饲养%d次，总时长%s，平均每次%s。`,
	// 	yestordayResult.Count, formatSeconds(int(yestordayResult.Sum)), formatSeconds(avgYestorday),
	// 	todayResult.Count, formatSeconds(int(todayResult.Sum)), formatSeconds(avgToday))
	txt = fmt.Sprintf("昨天共饲养%d次，总时长%s，平均每次%s，平均间隔%s。\r\n今天已饲养%d次，总时长%s，平均每次%s，平均间隔%s。",
		yestordayResult.Count, formatSeconds(int(yestordayResult.Sum)), formatSeconds(avgYestorday), formatSeconds(avgInterYestorday),
		todayResult.Count, formatSeconds(int(todayResult.Sum)), formatSeconds(avgToday), formatSeconds(avgInterToday))

	return txt
}

func InsertQrchatRecord(txt string) {
	if nil == LowDreamORM {
		beego.Error("InsertQrchatRecord failed: db not connected")
		return
	}
	beego.Info("InsertQrchatRecord ...", txt)
	stable := beego.AppConfig.String("dsdb::tbqrchat")
	sql := fmt.Sprintf(`insert into %s (txt) values('%s')`, stable, txt)
	_, err := LowDreamORM.Raw(sql).Exec()
	if err != nil {
		beego.Info("InsertQrchatRecord ... err:", err.Error())
	}
}
