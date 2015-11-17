package common

/*
import "download"
import "fmt"
import "io/ioutil"
import "strings"

//import "mahonia"

func (mgr *stockMgr) updateStockSpecificDatabase() error {
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

*/
