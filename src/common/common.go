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

type IData interface {
	Close()
}
type IMainData interface {
	IData
	GetStockCount() int
	InsertData(codes, names []string) error
}

type CData struct {
	db *sql.DB
}

//idata func
func (data *CData) Close() {
	data.db.Close()
}

//-----------------------end----------------------------------
//main func
func (data *CData) GetStockCount() int {
	rows, err := data.db.Query("select count(*) from " + MAIN_TABLE)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	var id int
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			fmt.Println(err)
			return 0
		}
	}
	return id
}

func (data *CData) InsertData(codes, names []string) error {
	stmt, err := data.db.Prepare("INSERT INTO " + MAIN_TABLE + "(code, name) VALUES(?, ?)")
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	for i := 0; i < len(codes); i++ {
		_, err = stmt.Exec(codes[i], names[i])
		if err != nil {
			continue
		}
	}
	return nil
}

//------------------------------------end---------------------------------------
type DataMgr struct {
}

var instance *DataMgr

func Instance() *DataMgr {
	if instance == nil {
		instance = new(DataMgr)
	}
	return instance
}

func (mgr *DataMgr) GetIData(dataName string) (data IData, err error) {
	data, err = mgr.open(dataName)
	return data, err
}

func (mgr *DataMgr) GetIMainData() (mainData IMainData, err error) {
	data, err := mgr.GetIData(MAIN_TABLE)
	mainData = data.(IMainData) //强转，可能会导致运行失败，如果类型不兼容
	return mainData, err
}

func (mgr *DataMgr) open(dataBaseName string) (data IData, err error) {
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

	if strings.EqualFold(dataBaseName, MAIN_TABLE) {
		//主表如果不存在则创建
		stock_database := MAIN_TABLE + "(code VARCHAR(40)  PRIMARY KEY ,name VARCHAR(40) NOT NULL)"
		_, err := db.Exec("create table if not exists " + stock_database)
		if err != nil {
			fmt.Println(err)
			db.Close()
			return nil, err
		}
	} else {
		//"select code, name from" + MAIN_TABLE + "  where code = " + dataBaseName
		rows, err := db.Query("select code, name from stockList  where code = " + dataBaseName)
		if err != nil {
			fmt.Print(err)
			return nil, err
		}
		defer rows.Close()
		//create table
		stock_database := dataBaseName + "(time VARCHAR(40)  PRIMARY KEY ,kaipan VARCHAR(40) NOT NULL)"
		_, err = db.Exec("create table if not exists " + stock_database)
		if err != nil {
			fmt.Println(err)
			db.Close()
			return nil, err
		}
	}

	cdata := new(CData)
	cdata.db = db
	return IData(cdata), err
}
