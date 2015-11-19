/*
建立数据库的股票管理模块，是程序在一个新地方运行的基础
*/
package common

import "fmt"
import "sync"
import "time"
import "strings"
import _ "github.com/go-sql-driver/mysql"
import "database/sql"

const (
	USER       = "huangchen"
	PASSWORD   = "199212"
	MAIN_TABLE = "stockList"
)

type stockInfo struct {
	name    string
	url     string
	address int // 0 深圳 1 上海
}

type stockMgr struct {
	stockMapUrl string
	stockUrl    string
	updateMap   map[string]int
	mutex       sync.Mutex
}

var s_stockMgr *stockMgr

func getStockMgr() *stockMgr {
	init := func() {
		s_stockMgr = new(stockMgr)
		s_stockMgr.stockMapUrl = "http://quote.eastmoney.com/stocklist.html#sz" //http://quote.eastmoney.com/stocklist.html"
		s_stockMgr.stockUrl = "http://hq.sinajs.cn/list="
	}
	var once sync.Once
	once.Do(init)
	return s_stockMgr
}

func (mgr *stockMgr) open(dataBaseName string) (data *cData, err error) {
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
		cdata.stock_name = dataBaseName
		err = mgr.updateMain(cdata)
		if err != nil {
			fmt.Println(err)
			db.Close()
			return nil, err
		}
		err = startRealTimeUpdate()
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
	return cdata, err
}

func (mgr *stockMgr) InsertMainData(cdata *cData, codes, names []string, address []int) error {
	stmt, err := cdata.db.Prepare("INSERT INTO " + MAIN_TABLE + "(code, name, address) VALUES(?, ?, ?)")
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

func (mgr *stockMgr) stockLock(stock string, isWrite bool) {
label:
	mgr.mutex.Lock()
	_, ok := mgr.updateMap[stock]
	if ok {
		if mgr.updateMap[stock] == -1 {
			mgr.mutex.Unlock()
			//锁住则继续等待1秒钟
			time.Sleep(1e9)
			goto label
		}
		if isWrite {
			if mgr.updateMap[stock] != 0 {
				//锁住则继续等待1秒钟
				time.Sleep(1e9)
				goto label
			} else {
				mgr.updateMap[stock] = -1
			}
		} else {
			mgr.updateMap[stock]++
		}

	} else {
		if isWrite {
			mgr.updateMap[stock] = -1
		} else {
			mgr.updateMap[stock]++
		}
	}
	mgr.mutex.Unlock()
}

func (mgr *stockMgr) stockUnLock(stock string) {
	mgr.mutex.Lock()
	if mgr.updateMap[stock] == -1 {
		mgr.updateMap[stock] = 0
	} else {
		mgr.updateMap[stock]--
	}
	mgr.mutex.Unlock()
}

func (mgr *stockMgr) GetStockData(cdata *cData, infos []string, count int) (datas []map[string]string, err error) {
	mgr.stockLock(cdata.stock_name, true)
	defer mgr.stockUnLock(cdata.stock_name)
	isNeed, err := isNeedSupplement(cdata, 123)
	if err != nil {
		return nil, err
	}
	if isNeed {
		su := Newsupplement()
		err = su.start()
		if err != nil {
			return nil, err
		}
	}
	//获取数据
	return nil, nil
}

func (mgr *stockMgr) GetInfoCount(cdata *cData) int {
	mgr.stockLock(cdata.stock_name, false)
	defer mgr.stockUnLock(cdata.stock_name)
	rows, err := cdata.db.Query("select count(*) from " + cdata.stock_name)
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

func (mgr *stockMgr) GetStockCount(cdata *cData) int {
	mgr.stockLock(cdata.stock_name, false)
	defer mgr.stockUnLock(cdata.stock_name)

	rows, err := cdata.db.Query("select count(*) from " + MAIN_TABLE)
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

func (mgr *stockMgr) InsertData(cdata *cData, info map[string]interface{}) error {
	//时间，开盘、最高、收盘、最低,
	/*date := info[STOCKINFO_DATE].(int)
	if date == 0 {
		return Error("insert fmt error")
	}
	rows, err := data.db.Query("select * from "+data.stock_name+" where date = ?", date)
	if err != nil {
		return err
	}
	defer rows.Close()

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
