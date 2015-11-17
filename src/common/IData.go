package common

type IData interface {
	Close()
	GetInfoCount() int
	GetData(count int)
}

type IMainData interface {
	IData
	GetStockCount() int
}
