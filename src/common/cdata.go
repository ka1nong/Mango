package common

//提供数据解析和检查是否可用的服务
import _ "github.com/go-sql-driver/mysql"
import "database/sql"

//import "fmt"

type cData struct {
	db         *sql.DB
	stock_name string
}

//开始时间，结束时间。密度
func (data *cData) GetData(infos []string, beingDate string, endData string, density int) (datas []map[string]string, err error) {
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

//todo:这个函数消耗内存，要做缓存
func (data *cData) GetAllStockInfo() (stocks []map[string]string, err error) {
	return nil, err
}

//------------------------------------end---------------------------------------
