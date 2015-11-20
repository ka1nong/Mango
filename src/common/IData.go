package common

type IData interface {
	Close()
	GetInfoCount() int
	GetData(infos []string, beingDate string, endData string, density int) (datas []map[string]string, err error)
}

type IMainData interface {
	Close()
	GetStockCount() int
	GetRandomMainData() (stock map[string]string, err error) //伪随机
	GetAllStockInfo() (stocks []map[string]string, err error)
}

type IStorage interface {
	Close()
	GetAllData(keyName string, valueName string) (key []interface{}, value []interface{}, err error)
	GetData(keyName string, valueName string) (key interface{}, value interface{}, err error)
	InsertData(keyName string, key interface{}, valueName string, value interface{}) error
}

func GetIData(dataName string) (data IData, err error) {
	mgr := getStockMgr()
	cdata, err := mgr.open(dataName)
	data = cdata
	return data, err
}

func GetIMainData() (mainData IMainData, err error) {
	cdata, err := GetIData(MAIN_TABLE)
	mainData = cdata.(IMainData)
	return mainData, err
}

func CreateStorage(name string) error {
	return nil
}
