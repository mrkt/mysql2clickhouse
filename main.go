package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ChangeMysqlTableToClickHouse(tableCreate string, tableName string) string {
	var tables = tableCreate
	var rows = strings.Split(tables, "\n")
	var replaceTables = ""
	var i = 0
	for _, v := range rows {
		if strings.Index(v, "KEY") > -1 {
			continue
		}
		if strings.Index(v, ") ENGINE=InnoDB") > -1 {
			v = ") ENGINE = Memory"
		}
		var changeRow string
		changeRow = strings.ReplaceAll(v, " NOT NULL", "")
		changeRow = strings.ReplaceAll(changeRow, " NOT NULL", "")
		changeRow = strings.ReplaceAll(changeRow, "AUTO_INCREMENT", "")
		changeRow = strings.ReplaceAll(changeRow, "CHARACTER SET utf8mb4", "")
		changeRow = strings.ReplaceAll(changeRow, "CHARACTER SET utf8", "")
		changeRow = strings.ReplaceAll(changeRow, "ON UPDATE CURRENT_TIMESTAMP", "")
		changeRow = strings.ReplaceAll(changeRow, "CURRENT_TIMESTAMP", "")
		changeRow = strings.ReplaceAll(changeRow, "datetime DEFAULT NULL", " DateTime ")
		changeRow = strings.ReplaceAll(changeRow, " datetime ", " DateTime ")
		changeRow = strings.ReplaceAll(changeRow, " text ", " String ")
		changeRow = string(regexp.MustCompile("varchar\\(\\d+\\) DEFAULT NULL").ReplaceAll([]byte(changeRow), []byte("Nullable(String)")))
		changeRow = string(regexp.MustCompile("varchar\\(\\d+\\)").ReplaceAll([]byte(changeRow), []byte("String")))
		changeRow = string(regexp.MustCompile("DEFAULT \\d+").ReplaceAll([]byte(changeRow), []byte("")))

		changeColumns := strings.Split(changeRow, " ")
		if strings.Index(changeColumns[3], "int") > -1 || strings.Index(changeColumns[3], "bigint") > -1 {
			tmpCCNum := changeColumns[3]
			tmpCCNum = strings.ReplaceAll(tmpCCNum, "bigint", "")
			tmpCCNum = strings.ReplaceAll(tmpCCNum, "tinyint", "")
			tmpCCNum = strings.ReplaceAll(tmpCCNum, "int", "")
			tmpCCNum = strings.ReplaceAll(tmpCCNum, "(", "")
			tmpCCNum = strings.ReplaceAll(tmpCCNum, ")", "")
			length, _ := strconv.Atoi(tmpCCNum)
			var intType string
			if strings.Index(changeColumns[3], "bigint") > -1 {
				intType = "bigint"
			} else if strings.Index(changeColumns[3], "tinyint") > -1 {
				intType = "tinyint"
			} else {
				intType = "int"
			}
			if intType == "tinyint" {
				changeRow = strings.Replace(changeRow, intType+"("+tmpCCNum+")", "Int8", 1)
			} else {
				if length < 3 {
					changeRow = strings.Replace(changeRow, intType+"("+tmpCCNum+")", "Int8", 1)
				} else if length < 5 {
					changeRow = strings.Replace(changeRow, intType+"("+tmpCCNum+")", "Int16", 1)
				} else if length <= 9 {
					changeRow = strings.Replace(changeRow, intType+"("+tmpCCNum+")", "Int32", 1)
				} else {
					changeRow = strings.Replace(changeRow, intType+"("+tmpCCNum+")", "Int64", 1)
				}
			}
		}
		replaceTables += changeRow
		i++
	}
	if strings.Index(replaceTables, ",) ENGINE = Memory") > -1 {

		temp := replaceTables[0:strings.Index(replaceTables, ",) ENGINE = Memory")]
		replaceTables = temp + ") ENGINE = Memory "
	}

	replaceTables = strings.ReplaceAll(replaceTables, tableName, tableName+"_local")
	return replaceTables
}

func main() {
	//这里是用脚本跑出来后粘上去的。show create table XXX表名;
	//用mybatis连接数据库查询出来建表语句比较好
	createTable := "CREATE TABLE `test_table` (\n" +
		"  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键(自动递增)',\n" +
		"  `user_id` int(10) NOT NULL COMMENT '用户id',\n" +
		"  `status` tinyint(4) DEFAULT 0 COMMENT '用户名、账号',\n" +
		"  `user_name` varchar(255) DEFAULT NULL COMMENT '用户名、账号',\n" +
		"  `table_id` varchar(255) DEFAULT NULL COMMENT '当前的表名',\n" +
		"  `filed` varchar(500) DEFAULT NULL COMMENT '存放有权限的字段',\n" +
		"  `database_id` varchar(100) DEFAULT NULL COMMENT '主题数据库',\n" +
		"  `is_del` varchar(1) DEFAULT NULL COMMENT '删除状态  ',\n" +
		"  `created_date` datetime DEFAULT NULL COMMENT '创建时间',\n" +
		"  `created_by` varchar(255) DEFAULT NULL COMMENT '创建人id',\n" +
		"  `created_name` varchar(255) DEFAULT NULL COMMENT '创建人姓名',\n" +
		"  `updated_date` datetime DEFAULT NULL COMMENT '更新时间',\n" +
		"  `updated_by` varchar(255) DEFAULT NULL COMMENT '更新人id',\n" +
		"  `updated_name` varchar(255) DEFAULT NULL COMMENT '更新人name',\n" +
		"  `flag` varchar(1) DEFAULT NULL COMMENT '0:待审核  1:通过  2:未通过',\n" +
		"  `apply_reason` varchar(255) COMMENT '申请原因',\n" +
		"  `ch_name` varchar(50) DEFAULT NULL COMMENT '表中文名',\n" +
		"  `descri_table` varchar(50) DEFAULT NULL COMMENT '表描述',\n" +
		"  `owner` varchar(50) DEFAULT NULL COMMENT '表的所有者',\n" +
		"  `audit_advice` varchar(255) DEFAULT NULL COMMENT '审批意见',\n" +
		"  `db_comment` varchar(100) DEFAULT NULL COMMENT '业务名称（数据库描述）',\n" +
		"  `db_class` varchar(100) DEFAULT NULL COMMENT '业务分类',\n" +
		"  `table_record_id` varchar(50) DEFAULT NULL COMMENT '表的记录Id',\n" +
		"  PRIMARY KEY (`id`)\n" +
		") ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COMMENT='数据权限申请记录表'"
	res := ChangeMysqlTableToClickHouse(createTable, "test_table")
	fmt.Println(res)
}
