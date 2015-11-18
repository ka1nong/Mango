package common

const (
	DATE  = "date"
	OPEN  = "open"
	HIGHT = "hight"
	CLOSE = "close"
	LOW   = "low"
)

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
	mgr := dataMgrInstance()
	data, err = mgr.open(dataName)
	return data, err
}

func GetIMainData() (mainData IMainData, err error) {
	data, err := GetIData(MAIN_TABLE)
	mainData = data.(IMainData)
	return mainData, err
}
