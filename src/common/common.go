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

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	USER       = "huangchen"
	PASSWORD   = "199212"
	MAIN_TABLE = "stockList"
)

type IData interface {
	Close()
	InsertData(info map[string]interface{}) error
	GetInfoCount() int
	GetData(count int)
}

type IMainData interface {
	IData
	GetStockCount() int
	InsertMainData(codes, names []string, address []int) error
}

type CData struct {
	db         *sql.DB
	stock_name string
}

//idata func-----------------------
func (data *CData) GetData(count int) {
}

func (data *CData) Close() {
	data.db.Close()
}

const (
	DATE  = "date"
	OPEN  = "open"
	HIGHT = "hight"
	CLOSE = "close"
	LOW   = "low"
)

func (data *CData) InsertData(info map[string]interface{}) error {
	//时间，开盘、最高、收盘、最低,
	date := info[DATE].(int)
	if date == 0 {
		return Error("insert fmt error")
	}
	rows, err := data.db.Query("select * from "+data.stock_name+" where date = ?", date)
	if err != nil {
		return err
	}
	defer rows.Close()
	/*
		tx, _ := data.db.Begin()
		//IF (SELECT * FROM ipstats WHERE date='192.168.0.1)' {
		//	    UPDATE ipstats SET clicks=clicks+1 WHERE ip='192.168.0.1';
		//	} else {
		//		 INSERT INTO ipstats (ip, clicks) VALUES ('192.168.0.1', 1);
		//}
		stmt, err := data.db.Prepare("INSERT INTO " + data.stock_name + "(date, open, hight, close, low) VALUES(?, ?, ?,?, ?)")
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}
		defer stmt.Close()
		//for i := 0; i < len(date); i++ {
		//	_, err = stmt.Exec(date[i], open[i], hight[i], close[i], low[i])
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}
		//}
		err = tx.Commit() //err := Tx.Rollback()
		return nil
	*/
	return nil
}

func (data *CData) GetInfoCount() int {
	rows, err := data.db.Query("select count(*) from " + data.stock_name)
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

func (data *CData) InsertMainData(codes, names []string, address []int) error {
	stmt, err := data.db.Prepare("INSERT INTO " + MAIN_TABLE + "(code, name, address) VALUES(?, ?, ?)")
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	for i := 0; i < len(codes); i++ {
		_, err = stmt.Exec(codes[i], names[i], address[i])
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
		stock_database := MAIN_TABLE + "(code VARCHAR(40)  PRIMARY KEY ,name VARCHAR(40) NOT NULL, address INTEGER)"
		_, err := db.Exec("create table if not exists " + stock_database)
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
	}

	cdata := new(CData)
	cdata.db = db
	cdata.stock_name = dataBaseName
	return IData(cdata), err
}
