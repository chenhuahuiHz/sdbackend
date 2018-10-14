/*
	定义所有查询的结构
*/
package dbproc

import (
	"strconv"
)

func MyAtoi(s string) int {
	n, err := strconv.Atoi(s)
	if nil != err {
		n = 0
	}
	return n
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
