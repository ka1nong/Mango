/*
http://1.guotie.sinaapp.com/?p=424
http://www.cnblogs.com/top5/archive/2010/09/14/1825571.html
基础组件
*/
package common

import _ "github.com/go-sql-driver/mysql"
import "database/sql"
import "fmt"
import "strings"

const (
	USER       = "huangchen"
	PASSWORD   = "199212"
	MAIN_TABLE = "stockList"
)

type dataMgr struct {
}

var instance *dataMgr

func dataMgrInstance() *dataMgr {
	if instance == nil {
		instance = new(dataMgr)
	}
	return instance
}

func GetIData(dataName string) (data IData, err error) {
	mgr := dataMgrInstance()
	data, err = mgr.open(dataName)
	return data, err
}

func GetIMainData() (mainData IMainData, err error) {
	data, err := GetIData(MAIN_TABLE)
	mainData = data.(IMainData) //强转，可能会导致运行失败，如果类型不兼容
	return mainData, err
}

func (mgr *dataMgr) open(dataBaseName string) (data IData, err error) {
	db, err := sql.Open("mysql", USER+":"+PASSWORD+"@/stock_data")
	if err != nil {
		fmt.Println("database initialize error : ", err.Error())
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("database initialize error : ", err.Error())
		db.Close()
		return nil, err
	}

	cdata := new(cData)
	if strings.EqualFold(dataBaseName, MAIN_TABLE) {
		//主表如果不存在则创建
		stock_database := MAIN_TABLE + "(code VARCHAR(40)  PRIMARY KEY ,name VARCHAR(40) NOT NULL, address INTEGER)"
		_, err := db.Exec("create table if not exists " + stock_database)
		if err != nil {
			fmt.Println(err)
			db.Close()
			return nil, err
		}
		cdata.db = db
		stockMgr := getStockMgr()
		err = stockMgr.updateMain(cdata)
		if err != nil {
			fmt.Println(err)
			db.Close()
			return nil, err
		}
	} else {
		rows, err := db.Query("select code, name from " + MAIN_TABLE + "  where code = " + dataBaseName)
		if err != nil {
			fmt.Print(err)
			return nil, err
		}
		defer rows.Close()
		//http://table.finance.yahoo.com/table.csv?s=000001.sz
		//表名不能用纯数字，我加上stock
		//create table 时间，开盘、最高、收盘、最低,
		//社会需要的是熟工，而不是你的学习能力
		dataBaseName = dataBaseName + "stock"
		stock_database := dataBaseName + "(data  INTEGER  PRIMARY KEY, open VARCHAR(40),hight VARCHAR(40), low VARCHAR(40), close VARCHAR(40))ENGINE=INNODB"
		_, err = db.Exec("create table if not exists " + stock_database)
		if err != nil {
			fmt.Println(err)
			db.Close()
			return nil, err
		}
		cdata.db = db
		cdata.stock_name = dataBaseName
	}
	return IData(cdata), err
}
