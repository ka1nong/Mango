/*
http://1.guotie.sinaapp.com/?p=424
基础组件
*/
package database

import _ "github.com/go-sql-driver/mysql"
import "database/sql"
import "fmt"

const (
	USER     = "huangchen"
	PASSWORD = "199212"
)

type IData interface {
	Close()
}
type IMainData interface {
	IData
}

type CData struct {
	db *sql.DB
}

func (data *CData) Close() {
	data.db.Close()
}

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

func (mgr *DataMgr) GetIMainData(dataName string) (mainData IMainData, err error) {
	data, err := mgr.GetIData(dataName)
	mainData = data
	return mainData, err
}

func (mgr *DataMgr) open(dataBaseName string) (data IData, err error) {
	db, err := sql.Open("mysql", USER+":"+PASSWORD+"@/dbname")
	if err != nil {
		fmt.Println("database initialize error : ", err.Error())
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("database initialize error : ", err.Error())
		return nil, err
	}

	cdata := new(CData)
	cdata.db = db
	return IData(cdata), err
}
