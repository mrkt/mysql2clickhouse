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
	//?????????????????????????????????????????????show create table XXX??????;
	//???mybatis????????????????????????????????????????????????
	createTable := "CREATE TABLE `test_table` (\n" +
		"  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '??????(????????????)',\n" +
		"  `user_id` int(10) NOT NULL COMMENT '??????id',\n" +
		"  `status` tinyint(4) DEFAULT 0 COMMENT '??????????????????',\n" +
		"  `user_name` varchar(255) DEFAULT NULL COMMENT '??????????????????',\n" +
		"  `table_id` varchar(255) DEFAULT NULL COMMENT '???????????????',\n" +
		"  `filed` varchar(500) DEFAULT NULL COMMENT '????????????????????????',\n" +
		"  `database_id` varchar(100) DEFAULT NULL COMMENT '???????????????',\n" +
		"  `is_del` varchar(1) DEFAULT NULL COMMENT '????????????  ',\n" +
		"  `created_date` datetime DEFAULT NULL COMMENT '????????????',\n" +
		"  `created_by` varchar(255) DEFAULT NULL COMMENT '?????????id',\n" +
		"  `created_name` varchar(255) DEFAULT NULL COMMENT '???????????????',\n" +
		"  `updated_date` datetime DEFAULT NULL COMMENT '????????????',\n" +
		"  `updated_by` varchar(255) DEFAULT NULL COMMENT '?????????id',\n" +
		"  `updated_name` varchar(255) DEFAULT NULL COMMENT '?????????name',\n" +
		"  `flag` varchar(1) DEFAULT NULL COMMENT '0:?????????  1:??????  2:?????????',\n" +
		"  `apply_reason` varchar(255) COMMENT '????????????',\n" +
		"  `ch_name` varchar(50) DEFAULT NULL COMMENT '????????????',\n" +
		"  `descri_table` varchar(50) DEFAULT NULL COMMENT '?????????',\n" +
		"  `owner` varchar(50) DEFAULT NULL COMMENT '???????????????',\n" +
		"  `audit_advice` varchar(255) DEFAULT NULL COMMENT '????????????',\n" +
		"  `db_comment` varchar(100) DEFAULT NULL COMMENT '?????????????????????????????????',\n" +
		"  `db_class` varchar(100) DEFAULT NULL COMMENT '????????????',\n" +
		"  `table_record_id` varchar(50) DEFAULT NULL COMMENT '????????????Id',\n" +
		"  PRIMARY KEY (`id`)\n" +
		") ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COMMENT='???????????????????????????'"
	res := ChangeMysqlTableToClickHouse(createTable, "test_table")
	fmt.Println(res)
}
