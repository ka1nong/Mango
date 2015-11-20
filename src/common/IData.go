package common

import "strconv"

type IData interface {
	Close()
	GetInfoCount() int
	GetData(infos []string, beingDate int64, endData int64, density string) (datas []map[string]interface{}, err error)
}

type IMainData interface {
	Close()
	GetStockCount() int
	GetRandomMainData() (stock map[string]interface{}, err error) //伪随机
	GetAllStockInfo() (stocks []map[string]interface{}, err error)
}

type IStorage interface {
	Close()
	GetAllData(keyName string, valueName string) (key []interface{}, value []interface{}, err error)
	GetData(keyName string, valueName string) (key interface{}, value interface{}, err error)
	InsertData(keyName string, key interface{}, valueName string, value interface{}) error
}

func GetIData(code int) (data IData, err error) {
	mgr := getStockMgr()
	cdata, err := mgr.open(strconv.Itoa(code))
	data = cdata
	return data, err
}

func GetIMainData() (mainData IMainData, err error) {
	mgr := getStockMgr()
	cdata, err := mgr.open(MAIN_TABLE)
	if err != nil {
		return nil, err
	}
	mainData = cdata
	return mainData, err
}

func CreateStorage(name string) error {
	return nil
}
