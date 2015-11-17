/*
建立数据库的股票管理模块，是程序在一个新地方运行的基础
*/
package common

import "fmt"
import "sync"
import "time"

type Error string

func (e Error) Error() string {
	return string(e)
}

type stockInfo struct {
	name    string
	url     string
	address int // 0 深圳 1 上海
}

type stockMgr struct {
	stockMapUrl string
	stockUrl    string
	updateMap   map[string]bool
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

func (mgr *stockMgr) updateMain(cdata *cData) error {
	initMainStocks := func() {
		fmt.Println("update main")
		err := mgr.updateMainDatabase(cdata)
		if err != nil {
			fmt.Print("update main  error:")
			fmt.Println(err)
		}
	}
	var once sync.Once
	once.Do(initMainStocks)
	return nil
}

func (mgr *stockMgr) updateMainDatabase(cdata *cData) error {
	count := cdata.GetStockCount()
	fmt.Println("main stock count is:%d", count)
	if count == 0 {
		stocks, err := mgr.loadStockMap()
		if err != nil {
			return err
		}
		//建立本地库
		codes := make([]string, len(stocks))
		names := make([]string, len(stocks))
		address := make([]int, len(stocks))
		i := 0
		for key, info := range stocks {
			codes[i] = key
			names[i] = info.name
			address[i] = info.address
			if len(codes[i]) != 0 && len(info.name) != 0 {
				i++
			}
		}
		err = cdata.InsertMainData(codes[:i], names[:i], address[:i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (mgr *stockMgr) waitStockLock(stock string) {
label:
	mgr.mutex.Lock()
	value, ok := mgr.updateMap[stock]
	if ok && value {
		mgr.mutex.Unlock()
		//锁住则继续等待1秒钟
		time.Sleep(1e9)
		goto label
	}
	//没锁住则现在锁住
	mgr.updateMap[stock] = true
	mgr.mutex.Unlock()
}

func (mgr *stockMgr) stockUnLock(stock string) {
	mgr.mutex.Lock()
	mgr.updateMap[stock] = false
	mgr.mutex.Unlock()
}

func (mgr *stockMgr) GetStockData(cdata *cData, stock string, infos []string, count int) (datas []map[string]string, err error) {
	mgr.waitStockLock(stock)

	mgr.stockUnLock(stock)
	return nil, nil
}
