package stockmanger

import "download"
import "fmt"

type StockMgr struct {
	stockUrl string
}

func NewStockMgr() *StockMgr {
	mgr := new(StockMgr)
	mgr.stockUrl = "http://quote.eastmoney.com/stocklist.html"
	return mgr
}

func (mgr *StockMgr) Start() error {

	fmt.Println("stock manger start run")

	dwMgr := download.Instance()
	dwMgr.Download(mgr.stockUrl)

	fmt.Println("stock manger end run")
	return nil
}
