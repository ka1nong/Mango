/*
建立数据库的股票管理模块，是程序在一个新地方运行的基础
*/
package common

import "download"
import "fmt"
import "io/ioutil"
import "strings"
import "mahonia"
import "os"
import "strconv"

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
	stocks      map[string]stockInfo
}

/*var instance *stockMgr

func Instance() *stockMgr {
	if instance == nil {
		instance = new(stockMgr)
		instance.stockMapUrl = "http://quote.eastmoney.com/stocklist.html#sz" //http://quote.eastmoney.com/stocklist.html"
		instance.stockUrl = "http://hq.sinajs.cn/list="
	}
	return instance
}
*/
func newStockMgr() *stockMgr {
	mgr := &stockMgr{}
	mgr.stockMapUrl = "http://quote.eastmoney.com/stocklist.html#sz" //http://quote.eastmoney.com/stocklist.html"
	mgr.stockUrl = "http://hq.sinajs.cn/list="
	return mgr
}

func (mgr *stockMgr) updateMain() error {
	return nil
}

func (mgr *stockMgr) Start() error {
	fmt.Println("stock manger start run")
	err := mgr.updateMainDatabase()
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = mgr.updateStockSpecificDatabase()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("stock manger end run")
	return nil
}

func (mgr *stockMgr) updateMainDatabase() error {
	//	dataMgr := dataMgrInstance()
	iMainData, err := GetIMainData()
	if err != nil {
		fmt.Println("local haven't main table")
		return err
	}
	defer iMainData.Close()
	count := iMainData.GetStockCount()
	fmt.Println(count)
	err = mgr.loadStockMap() //要把数据库中的数据导到stocks中
	if err != nil {
		return err
	}

	if count == 0 {
		err = mgr.loadStockMap()
		if err != nil {
			return err
		}
		//建立本地库
		codes := make([]string, len(mgr.stocks))
		names := make([]string, len(mgr.stocks))
		address := make([]int, len(mgr.stocks))
		i := 0
		for key, info := range mgr.stocks {
			codes[i] = key
			names[i] = info.name
			address[i] = info.address
			if len(codes[i]) != 0 && len(info.name) != 0 {
				i++
			}
		}
		//err = iMainData.InsertMainData(codes[:i], names[:i], address[:i])
		if err != nil {
			return err
		}

	} else {
		err = mgr.loadStockMap()
		if err != nil {
			//远程下载失败，但是本地存在，则不认为是错误
			return nil
		}

		//todo:比较一下是否要用远程库
	}
	return nil
}

func (mgr *stockMgr) updateStockSpecificDatabase() error {
	//	dataMgr := dataMgrInstance()
	iMainData, err := GetIMainData()
	if err != nil {
		fmt.Println("local haven't main table")
		return err
	}
	defer iMainData.Close()

	count := 0
	for key, _ := range mgr.stocks {
		idata, err := GetIData(key)
		if err != nil {
			count++
			continue
		}
		defer idata.Close()

		count := idata.GetInfoCount()
		//小于一个月则需要去下载
		if count < 30 {
			var address string
			if mgr.stocks[key].address == 0 {
				address = "sz"
			} else {
				address = "ss"
			}
			url := "http://table.finance.yahoo.com/table.csv?s=" + key + "." + address
			dwMgr := download.Instance()
			localPath, err := dwMgr.Download(url)
			if err != nil {
				return err
			}
			err = mgr.getStockInfoByParseFile(localPath)
			//idata.InsertData()
		}
	}

	return nil
}

func (mgr *stockMgr) loadSpecificStockInfo(stockCode string) error {
	dwMgr := download.Instance()
	localPath, err := dwMgr.Download(mgr.stockUrl + string("sh") + stockCode)
	buf, err := ioutil.ReadFile(localPath)
	if err != nil {
		fmt.Println("download specific stock info error .the code is %s", stockCode)
	}
	stockText := string(buf)
	index := strings.Index(stockText, string("=\""))
	if index == -1 || index > len(stockText)-1 {
		fmt.Println("download specific stock info error .the code is %s", stockCode)
	}
	stockText = stockText[index+len(string("=\"")) : len(stockText)-3]
	stockInfos := strings.Split(stockText, string(","))
	for _, value := range stockInfos {
		fmt.Println(value)
	}

	//http://blog.sciencenet.cn/home.php?mod=space&uid=461456&do=blog&id=455211
	return nil
}

func (mgr *stockMgr) loadStockMap() error {
	dwMgr := download.Instance()
	localPath, err := dwMgr.Download(mgr.stockMapUrl)
	bHasRemoteStockMap := true
	if err != nil {
		fmt.Println("load remote stock map error")
		bHasRemoteStockMap = false
	}
	remoteStockMap, err := mgr.getStockMapByParseFile(localPath)
	if err != nil {
		fmt.Println("parse remote stock map error")
		bHasRemoteStockMap = false
	}

	localStockMap, err := mgr.loadLocalStockMap()
	if err != nil {
		//不存在本地缓存或本地缓存加载失败
		if bHasRemoteStockMap {
			//使用远程库
			mgr.stocks = remoteStockMap
		} else {
			return Error("load Stock Map error")
		}
	} else {
		if bHasRemoteStockMap {
			//都存在则进行比较
			localStocksCount := len(localStockMap)
			remoteStoksCount := len(remoteStockMap)
			//尼玛，没找到abs函数,
			if localStocksCount > remoteStoksCount {
				if localStocksCount-remoteStoksCount < 50 {
					mgr.stocks = localStockMap
				} else {
					return Error("local stock map count is very big")
				}
			} else {
				if remoteStoksCount-localStocksCount < 50 {
					mgr.stocks = remoteStockMap
				} else {
					return Error("remote stock map count is very big")
				}
			}

		} else {
			//使用本地缓存
			mgr.stocks = localStockMap
		}
	}

	err = mgr.saveLocalStockMap()
	if err != nil {
		return err
	}
	return nil
}

func (mgr *stockMgr) getStockInfoByParseFile(localPath string) error {
	buf, err := ioutil.ReadFile(localPath)
	if err != nil {
		return Error("read stock info error")
	}
	stockInfoText := string(buf)
	strInfos := strings.Split(stockInfoText, "\n")
	if len(strInfos) < 2 {
		return Error("parse stock info error")
	}

	for i := 1; i < len(strInfos); i++ {
		_ = strings.Split(strInfos[i], ",")

	}
	return nil
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

func (mgr *stockMgr) loadLocalStockMap() (stockMap map[string]stockInfo, err error) {
	//暂时还不实现本地化存储

	return stockMap, Error("unknown eror")
}

func (mgr *stockMgr) saveLocalStockMap() error {
	stocksPath := "../stocksPath.txt"
	if _, err := os.Stat(stocksPath); err == nil {
		os.Remove(stocksPath)
	}
	fout, err := os.Create(stocksPath)
	if err != nil {
		return Error("saveLoaclStockMap error")
	}
	defer fout.Close()
	for key, info := range mgr.stocks {
		text := strconv.Itoa(len(key)) + string("|") + key + strconv.Itoa(len(info.name)) + string("|") + info.name + strconv.Itoa(len(info.url)) + string("|") + info.url
		_, err := fout.WriteString(text)
		if err != nil {
			return Error("saveLoaclStockMap error")
		}
	}
	return nil
}
