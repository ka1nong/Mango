package common

import "download"
import "fmt"
import "io/ioutil"
import "strings"
import "mahonia"
import "sync"

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
		err = mgr.InsertMainData(cdata, codes[:i], names[:i], address[:i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (mgr *stockMgr) loadStockMap() (map[string]stockInfo, error) {
	dwMgr := download.Instance()
	localPath, err := dwMgr.Download(mgr.stockMapUrl)
	if err != nil {
		fmt.Println("load remote stock map error")
		return nil, err
	}
	remoteStockMap, err := mgr.getStockMapByParseFile(localPath)
	if err != nil {
		fmt.Println("parse remote stock map error")
		return nil, err
	}
	return remoteStockMap, nil
}

func (mgr *stockMgr) getStockMapByParseFile(localPath string) (stockMap map[string]stockInfo, err error) {
	buf, err := ioutil.ReadFile(localPath)
	if err != nil {
		return stockMap, Error("read remote stock map error")
	}

	enc := mahonia.NewEncoder("gbk")
	remoteStocksText := string(buf)
	//remoteStocksText = enc.ConvertString(remoteStocksText) //转换utf-8编码
	index := strings.Index(remoteStocksText, enc.ConvertString("股票代码查询一览表："))
	if index == -1 {
		return stockMap, Error("parse remote stock map error")
	}
	remoteStocksText = remoteStocksText[index:]
	index = strings.Index(remoteStocksText, string("<ul>"))
	if index == -1 {
		return stockMap, Error("parse remote stock map error")
	}
	shangHaiEndIndex := strings.Index(remoteStocksText, string("</ul>"))
	if shangHaiEndIndex == -1 || shangHaiEndIndex < index {
		return stockMap, Error("parse remote stock map error")
	}
	shangHaiStocksText := remoteStocksText[index:shangHaiEndIndex]

	remoteStocksText = remoteStocksText[shangHaiEndIndex+1:] //从加1开始查深圳股票，不然shenZhenEndIndex就是0
	index = strings.Index(remoteStocksText, string("<ul>"))
	if index == -1 {
		return stockMap, Error("parse remote stock map error")
	}
	shenZhenEndIndex := strings.Index(remoteStocksText, string("</ul>"))
	if shenZhenEndIndex == -1 || shenZhenEndIndex < index {
		return stockMap, Error("parse remote stock map error")
	}
	shenZhenStocksText := remoteStocksText[index:shenZhenEndIndex]

	shanghaiStocks, err := mgr.parseRemoteStockMapByUL(shangHaiStocksText)
	if err != nil {
		return stockMap, err
	}

	shenZhenStocks, err := mgr.parseRemoteStockMapByUL(shenZhenStocksText)
	if err != nil {
		return stockMap, err
	}

	stockMap = make(map[string]stockInfo)
	for key, info := range shanghaiStocks {
		info.address = 1
		stockMap[key] = info
	}
	//充分利用深圳股票的枚举是0
	for key, info := range shenZhenStocks {
		stockMap[key] = info
	}
	return stockMap, err
}

func (mgr *stockMgr) parseRemoteStockMapByUL(stocksText string) (stockMap map[string]stockInfo, err error) {

	stockMap = make(map[string]stockInfo)

	parseli := func(litext string) (err error) {
		index := strings.Index(litext, string("href="))
		if index == -1 {
			return Error("parse remote li error")
		}
		litext = litext[index+len(string("href=")):]
		index = strings.Index(litext, string(">"))
		if index == -1 {
			return Error("parse remote li error")
		}
		url := litext[1 : index-1] //去除两边的分号
		litext = litext[index+1:]
		index = strings.Index(litext, string("("))
		if index == -1 {
			return Error("parse remote li error")
		}
		name := litext[:index]
		litext = litext[index+1:]
		index = strings.Index(litext, string(")"))
		if index == -1 {
			return Error("parse remote li error")
		}
		stockid := litext[:index]
		if len(url) == 0 || len(name) == 0 || len(stockid) == 0 {
			return Error("parse remote li error")
		}

		stockMap[stockid] = stockInfo{name, url, 0}
		return err
	}

	for {
		liBeginIndex := strings.Index(stocksText, string("<li>"))
		liEndIndex := strings.Index(stocksText, string("</li>"))
		if liBeginIndex == -1 || liEndIndex == -1 || liEndIndex < liBeginIndex {
			break
		}
		text := stocksText[liBeginIndex:liEndIndex]
		err := parseli(text)
		if err != nil {
			return stockMap, err
		}
		stocksText = stocksText[liEndIndex+len(string("</li>")):]
	}
	return stockMap, err
}
