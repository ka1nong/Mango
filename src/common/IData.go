package common

type IData interface {
	Close()
	GetInfoCount() int
	GetData(infos []string, count int) (datas []map[string]string, err error)
}

type IMainData interface {
	Close()
	GetStockCount() int
	GetRandomMainData() (stock map[string]string, err error) //伪随机
}

func GetIData(dataName string) (data IData, err error) {
	mgr := getStockMgr()
	cdata, err := mgr.open(dataName)
	return cdata.(IData), err
}

func GetIMainData() (mainData IMainData, err error) {
	cdata, err := GetIData(MAIN_TABLE)
	mainData = cdata.(IMainData)
	return mainData, err
}
