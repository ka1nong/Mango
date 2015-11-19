package common

import _ "github.com/go-sql-driver/mysql"
import "database/sql"

//import "fmt"

type cData struct {
	db         *sql.DB
	stock_name string
}

//开始时间，结束时间。密度
func (data *cData) GetData(infos []string, count int) (datas []map[string]string, err error) {
	if count == 0 || len(infos) == 0 {
		return
	}
	return nil, Error("param error")
}

func (data *cData) Close() {
	data.db.Close()
}

func (data *cData) InsertData(info map[string]interface{}) error {
	mgr := getStockMgr()
	return mgr.InsertData(data, info)
}

func (data *cData) GetInfoCount() int {
	mgr := getStockMgr()
	return mgr.GetInfoCount(data)
}

func (data *cData) GetStockCount() int {
	mgr := getStockMgr()
	return mgr.GetStockCount(data)
}

func (data *cData) GetRandomMainData() (stock map[string]string, err error) {
	count := data.GetStockCount()
	if count != 0 {
	}
	return nil, err
}

//------------------------------------end---------------------------------------
