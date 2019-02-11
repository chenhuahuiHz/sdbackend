/*
	定义所有查询的结构
*/
package dbproc

import (
	"strconv"

	"github.com/astaxie/beego"
)

func MyAtoi(s string) int {
	n, err := strconv.Atoi(s)
	if nil != err {
		n = 0
	}
	return n
}

const MAX_HIT_CNT = 1000

type SqlCache struct {
	objListCache map[string][]ObjRow
	hitCount     map[string]uint32
}

func (this *SqlCache) init() {
	this.objListCache = make(map[string][]ObjRow)
	this.hitCount = make(map[string]uint32)
}

func (this *SqlCache) setObjListCache(sql string, data []ObjRow) {
	if len(data) > 0 {
		this.objListCache[sql] = data
		this.hitCount[sql] = 0
	}
}

func (this *SqlCache) getObjListCache(sql string) (data []ObjRow, exist bool) {
	data, exist = this.objListCache[sql]
	if exist {
		this.hitCount[sql]++
		if this.hitCount[sql] > MAX_HIT_CNT {
			this.hitCount[sql] = 0
			delete(this.objListCache, sql)
			beego.Info("hit time expired, clear cache for sql:", sql)
		}
	}
	return
}

//`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
//`main_type` tinyint(4) unsigned DEFAULT '1' COMMENT '1-房子 2-车子 3-吃喝',
//`sub_type` tinyint(4) unsigned DEFAULT '1',
//`desc` varchar(128) DEFAULT '神秘物品',
//`price` int(11) DEFAULT '0',
//`tast` int(11) DEFAULT '0',
type ObjRow struct {
	Id       uint32
	MainType int8
	SubType  int8
	Desc     string
	Price    uint32
	Tast     uint32
}

// `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
// `type` tinyint(4) unsigned DEFAULT '1' COMMENT '1-吃奶 2-拉屎',
// `state` tinyint(4) unsigned DEFAULT '0' COMMENT '0-进行中 1-结束',
// `start_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
// `stop_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
// `cost_seconds` bigint(20) unsigned DEFAULT '0',
// `desc` varchar(128) DEFAULT '无',
type BabyRow struct {
	Id          uint32
	Type        int8
	State       int8
	StartTime   string
	StopTime    string
	CostSeconds uint32
	Desc        string
}

type StatistRow struct {
	Count uint32
	Sum   uint32
}
